package compiler

import (
	"log"

	"github.com/pkg/errors"
)

// AST represents the abstract GAPI syntax tree
type AST struct {
	SchemaName string
}

// Compiler represents a GAPI compiler
type Compiler struct{}

// NewCompiler creates a new compiler instance
func NewCompiler() (*Compiler, error) {
	c := &Compiler{}
	return c, nil
}

// Compile compiles
func (c *Compiler) Compile(source string) (*AST, error) {
	parser := GAPIParser{
		Buffer: source,
	}

	// Initialize parser
	if err := parser.Init(); err != nil {
		return nil, errors.Wrap(err, "parser init")
	}

	// Parse source
	if err := parser.Parse(); err != nil {
		log.Fatalf("parser: %s", err)
	}

	ast := &AST{}
	root := parser.AST()

	// Extract schema name
	schemaNameTkn := root.up.up.next.next.token32
	ast.SchemaName = source[schemaNameTkn.begin:schemaNameTkn.end]

	//TODO: Perform semantic analysis

	return ast, nil
}

// Analize analizes the parser output
func (c *Compiler) Analize() error {
	return nil
}
