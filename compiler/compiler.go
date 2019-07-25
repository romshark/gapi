package compiler

import "github.com/romshark/gapi/compiler/parser"

// Compile compiles the source file returning an abstract syntax tree
func Compile(source parser.SourceFile) (*parser.AST, error) {
	parser, err := parser.NewParser()
	if err != nil {
		return nil, err
	}
	if err := parser.Parse(source); err != nil {
		return nil, err
	}
	return parser.AST(), nil
}
