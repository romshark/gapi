package compiler

import (
	"log"

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

// QueryEndpoint represents a query endpoint
type QueryEndpoint struct {
	Src
	Name string
	Vars []Variable
	Type Type
}

// Mutation represents a mutation endpoint
type Mutation struct {
	Src
	Name string
	Vars []Variable
	Type Type
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

func getSrc(source string, token *node32) string {
	return source[token.begin:token.end]
}

// Compiler represents a GAPI compiler
type Compiler struct {
	parser       GAPIParser
	errors       []Error
	deferredJobs []job
	ast          *AST
}

// NewCompiler creates a new compiler instance
func NewCompiler(source string) (*Compiler, error) {
	c := &Compiler{
		parser: GAPIParser{
			Buffer: source,
		},
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
	c.errors = append(c.errors, err)
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
		Types: make(map[string]Type),
	}

	// Extract schema name
	current := root.up
	if current.pegRule == ruleSpOpt {
		// Ignore space before schema declaration
		current = current.next
	}
	c.ast.SchemaName = getSrc(c.parser.Buffer, current.up.next.next)

	if err := verifySchemaName(c.ast.SchemaName); err != nil {
		c.err(cErr{ErrSchemaIllegalIdent, err.Error()})
	}

	// Read all declarations
	var handler func(node *node32) error
	for current = root.up; current != nil; current = current.next {
		switch current.pegRule {
		case ruleDclAl:
			// Alias type declaration
			handler = c.defineAliasType
		case ruleDclEn:
			// Enum type declaration
			handler = c.defineEnumType
		case ruleDclRv:
			// Resolver type declaration
			log.Print("Resolver type declaration")
		case ruleDclSt:
			// Struct type declaration
			log.Print("Struct type declaration")
		case ruleDclTr:
			// Trait type declaration
			log.Print("Trait type declaration")
		case ruleDclUn:
			// Union type declaration
			handler = c.defineUnionType
		case ruleDclQr:
			// Query declaration
			log.Print("Query declaration")
		case ruleDclMt:
			// Mutation declaration
			log.Print("Mutation declaration")
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

	if len(c.errors) > 0 {
		return errors.Errorf("%d compiler errors", len(c.errors))
	}

	return nil
}
