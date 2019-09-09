package compiler

import "github.com/romshark/gapi/compiler/parser"

// Compile compiles the source file returning an abstract syntax tree
func Compile(source parser.SourceFile) (*parser.SchemaModel, error) {
	parser, err := parser.NewParser()
	if err != nil {
		return nil, err
	}
	if _, err = parser.Parse(source); err != nil {
		return nil, err
	}
	return parser.SchemaModel(), nil
}
