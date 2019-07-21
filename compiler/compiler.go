package compiler

import (
	"fmt"
	"log"
	"sync"

	"github.com/pkg/errors"
)

// Src represents the source node meta information
type Src struct {
	Begin uint32
	End   uint32
}

func src(node *node32) Src {
	return Src{
		Begin: node.begin,
		End:   node.end,
	}
}

func stdTypeByName(name string) Type {
	switch name {
	case "None":
		return TypeStdNone{}
	case "Bool":
		return TypeStdBool{}
	case "Byte":
		return TypeStdByte{}
	case "Int32":
		return TypeStdInt32{}
	case "Uint32":
		return TypeStdUint32{}
	case "Int64":
		return TypeStdInt64{}
	case "Uint64":
		return TypeStdUint64{}
	case "Float64":
		return TypeStdFloat64{}
	case "String":
		return TypeStdString{}
	case "Time":
		return TypeStdTime{}
	default:
		return nil
	}
}

type job func() error

func (c *Compiler) getSrc(token *node32) string {
	return c.parser.Buffer[token.begin:token.end]
}

// Compiler represents a GAPI compiler
type Compiler struct {
	parser            GAPIParser
	errors            []Error
	errorsLock        *sync.Mutex
	deferredJobs      []job
	ast               *AST
	lastIssuedGraphID GraphNodeID
	lastIssuedTypeID  TypeID
	lastIssuedParamID ParamID
	paramsByName      map[GraphNode]map[string]*Parameter
	graphNodeByName   map[string]GraphNode
	typeByName        map[string]Type
	typeByID          map[TypeID]Type
	graphNodeByID     map[GraphNodeID]GraphNode
	paramByID         map[ParamID]*Parameter
}

// NewCompiler creates a new compiler instance
func NewCompiler(source string) (*Compiler, error) {
	c := &Compiler{
		errorsLock: &sync.Mutex{},
		parser: GAPIParser{
			Buffer: source,
		},
		lastIssuedTypeID: TypeIDUserTypeOffset,
		paramsByName:     make(map[GraphNode]map[string]*Parameter),
		graphNodeByName:  make(map[string]GraphNode),
		typeByName:       make(map[string]Type),
		typeByID:         make(map[TypeID]Type),
		graphNodeByID:    make(map[GraphNodeID]GraphNode),
		paramByID:        make(map[ParamID]*Parameter),
	}

	// Initialize parser
	if err := c.parser.Init(); err != nil {
		return nil, errors.Wrap(err, "parser init")
	}

	return c, nil
}

func (c *Compiler) deferJob(job job) {
	c.deferredJobs = append(c.deferredJobs, job)
}

func (c *Compiler) err(err Error) {
	if err == nil {
		panic("nil error")
	}
	if err.Code() == 0 {
		panic("invalid error code (0)")
	}
	c.errorsLock.Lock()
	c.errors = append(c.errors, err)
	c.errorsLock.Unlock()
}

// resetState resets the compiler
func (c *Compiler) resetState() {
	c.errors = nil
	c.deferredJobs = nil
	c.ast = nil
}

// Errors returns a copy of the list of all compiler errors
func (c *Compiler) Errors() []Error {
	errs := make([]Error, len(c.errors))
	copy(errs, c.errors)
	return errs
}

// AST returns a copy of the abstract syntax tree
func (c *Compiler) AST() *AST {
	if len(c.errors) > 0 {
		return nil
	}
	return c.ast.Clone()
}

// Compile compiles
func (c *Compiler) Compile() error {
	c.resetState()

	// Parse source
	if err := c.parser.Parse(); err != nil {
		// Tokenization errors are fatal because we're missing the parse tree
		c.err(cErr{ErrSyntax, err.Error()})
		return errors.Errorf("tokenization error: %s", err)
	}

	// Get parse-tree
	root := c.parser.AST()
	c.ast = &AST{
		Types:          make([]Type, 0),
		AliasTypes:     make([]Type, 0),
		EnumTypes:      make([]Type, 0),
		UnionTypes:     make([]Type, 0),
		QueryEndpoints: make([]*QueryEndpoint, 0),
		Mutations:      make([]*Mutation, 0),
	}

	// Extract schema name
	current := root.up
	if current.pegRule == ruleSpOpt {
		// Ignore space before schema declaration
		current = current.next
	}
	c.ast.SchemaName = c.getSrc(current.up.next.next)

	if err := verifyLowerCamelCase(c.ast.SchemaName); err != nil {
		c.err(cErr{
			ErrSchemaIllegalIdent,
			fmt.Sprintf(
				"invalid schema identifier at %d:%d: %s",
				current.begin,
				current.end,
				err,
			),
		})
	}

	// Read all declarations
	var handler func(node *node32) error
	for current = root.up; current != nil; current = current.next {
		switch current.pegRule {
		case ruleDclAl:
			// Alias type declaration
			handler = c.parseDeclAls
		case ruleDclEn:
			// Enum type declaration
			handler = c.parseDeclEnm
		case ruleDclRv:
			// Resolver type declaration
			handler = c.parseDeclRsv
		case ruleDclSt:
			// Struct type declaration
			handler = c.parseDeclStr
		case ruleDclTr:
			// Trait type declaration
			log.Print("Trait type declaration")
		case ruleDclUn:
			// Union type declaration
			handler = c.parseDeclUni
		case ruleDclQr:
			// Query declaration
			handler = c.parseDeclQry
		case ruleDclMt:
			// Mutation declaration
			handler = c.parseDeclMut
		case ruleDclSb:
			// Subscription declaration
			log.Print("Subscription declaration")
		default:
			continue
		}

		if fatalErr := handler(current); fatalErr != nil {
			return errors.Wrap(fatalErr, "fatal compiler error")
		}
	}

	// Executed all deferred jobs
	for _, job := range c.deferredJobs {
		if fatalErr := job(); fatalErr != nil {
			return errors.Wrap(fatalErr, "fatal compiler error")
		}
	}

	wg := &sync.WaitGroup{}

	// Sort everything by name (ascending)
	wg.Add(8)
	go func() { sortTypesByName(c.ast.Types); wg.Done() }()
	go func() { sortTypesByName(c.ast.AliasTypes); wg.Done() }()
	go func() { sortTypesByName(c.ast.EnumTypes); wg.Done() }()
	go func() { sortTypesByName(c.ast.UnionTypes); wg.Done() }()
	go func() { sortTypesByName(c.ast.StructTypes); wg.Done() }()
	go func() { sortTypesByName(c.ast.ResolverTypes); wg.Done() }()
	go func() { sortQueryEndpointsByName(c.ast.QueryEndpoints); wg.Done() }()
	go func() { sortMutationsByName(c.ast.Mutations); wg.Done() }()
	//TODO: sort trait types
	wg.Wait()

	wg.Add(2)
	go func() {
		// Find all recursive alias type cycles
		defer wg.Done()
		cycles := c.findAliasTypeCycles()
		for _, cycle := range cycles {
			c.err(cErr{
				ErrAliasRecurs,
				fmt.Sprintf("Recursive alias type cycle: %s", cycle.String()),
			})
		}
	}()
	go func() {
		// Find all recursive struct type cycles
		defer wg.Done()
		cycles := c.findStructTypeCycles()
		for _, cycle := range cycles {
			c.err(cErr{
				ErrStructRecurs,
				fmt.Sprintf("Recursive struct type cycle: %s", cycle.String()),
			})
		}
	}()
	wg.Wait()

	if len(c.errors) > 0 {
		return CompilationErr{c.Errors()}
	}

	return nil
}
