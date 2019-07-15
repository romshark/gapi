package compiler

import (
	"fmt"
)

// AST represents the abstract GAPI syntax tree
type AST struct {
	Types          map[string]Type
	QueryEndpoints map[string]QueryEndpoint
	Mutations      map[string]Mutation
	SchemaName     string
}

// Clone returns a copy of the abstract syntax tree
func (ast *AST) Clone() *AST {
	if ast == nil {
		return nil
	}

	types := make(map[string]Type, len(ast.Types))
	for k, v := range ast.Types {
		types[k] = v
	}

	queryEndpoints := make(map[string]QueryEndpoint, len(ast.QueryEndpoints))
	for k, v := range ast.QueryEndpoints {
		queryEndpoints[k] = v
	}

	mutations := make(map[string]Mutation, len(ast.Mutations))
	for k, v := range ast.Mutations {
		mutations[k] = v
	}

	return &AST{
		Types:          types,
		QueryEndpoints: queryEndpoints,
		Mutations:      mutations,
		SchemaName:     ast.SchemaName,
	}
}

func (ast *AST) typeByName(name string) Type {
	if t := stdTypeByName(name); t != nil {
		return t
	}
	return ast.Types[name]
}

// defineType returns an error if the type name is already reserved
func (ast *AST) defineType(newType Type) Error {
	// Check for collisions with reserved primitive types
	srcNode := newType.Src()
	name := newType.Name()
	if stdTypeByName(name) != nil {
		return cErr{
			ErrTypeRedecl,
			fmt.Sprintf(
				"Redeclaration of type %s at %d:%d (reserved primitive type)",
				name,
				srcNode.Begin,
				srcNode.End,
			),
		}
	}

	// Check for collisions with other user-defined types
	if reservedBy, reserved := ast.Types[name]; reserved {
		reservedBySrcNode := reservedBy.Src()
		return cErr{
			ErrTypeRedecl,
			fmt.Sprintf("Redeclaration of type %s at %d:%d "+
				"(previous declaration: %d:%d (%s))",
				name,
				srcNode.Begin,
				srcNode.End,
				reservedBySrcNode.Begin,
				reservedBySrcNode.End,
				reservedBy.Category(),
			),
		}
	}

	// Define
	ast.Types[name] = newType

	return nil
}
