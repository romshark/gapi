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

func isStandardTypeName(name string) bool {
	switch name {
	case "None":
	case "Bool":
	case "Byte":
	case "Int32":
	case "Uint32":
	case "Int64":
	case "Uint64":
	case "Float64":
	case "String":
	case "Time":
	default:
		return false
	}
	return true
}

type job func() error

func getSrc(source string, token *node32) string {
	return source[token.begin:token.end]
}

// Compiler represents a GAPI compiler
type Compiler struct {
	deferredJobs []job
}

// NewCompiler creates a new compiler instance
func NewCompiler() (*Compiler, error) {
	c := &Compiler{}
	return c, nil
}

func (c *Compiler) deferJob(job job) {
	c.deferredJobs = append(c.deferredJobs, job)
}

// Compile compiles
func (c *Compiler) Compile(src string) (*AST, error) {
	parser := GAPIParser{
		Buffer: src,
	}

	// Initialize parser
	if err := parser.Init(); err != nil {
		return nil, errors.Wrap(err, "parser init")
	}

	// Parse source
	if err := parser.Parse(); err != nil {
		log.Fatalf("parser: %s", err)
	}

	// Get parse-tree
	root := parser.AST()
	ast := &AST{
		Types: make(map[string]Type),
	}

	// Extract schema name
	ast.SchemaName = getSrc(src, root.up.up.next.next)

	// Read all declarations
	current := root.up
	var handler func(src string, ast *AST, node *node32) error
	for ; current != nil; current = current.next {
		switch current.pegRule {
		case ruleDclAl:
			// Alias type declaration
			handler = c.defineAliasType
		case ruleDclEn:
			// Enum type declaration
			log.Print("Enum type declaration")
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
			log.Print("Union type declaration")
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

		if err := handler(src, ast, current); err != nil {
			return nil, err
		}
	}

	// Executed all deferred jobs
	for _, job := range c.deferredJobs {
		if err := job(); err != nil {
			return nil, err
		}
	}

	//TODO: Perform semantic analysis

	return ast, nil
}
