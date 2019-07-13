package compiler

import "github.com/pkg/errors"

// AST represents the abstract GAPI syntax tree
type AST struct {
	Types          map[string]Type
	QueryEndpoints map[string]QueryEndpoint
	Mutations      map[string]Mutation
	SchemaName     string
}

func (ast *AST) isTypeNameDefined(name string) bool {
	if isStandardTypeName(name) {
		return true
	}
	_, defined := ast.Types[name]
	return defined
}

// defineType returns an error if the type name is already reserved
func (ast *AST) defineType(newType Type) error {
	// Check for collisions with reserved standard types
	srcNode := newType.Src()
	name := newType.Name()
	if isStandardTypeName(name) {
		return errors.Errorf(
			"Redeclaration of type %s at %d:%d (reserved standard type)",
			name,
			srcNode.Begin,
			srcNode.End,
		)
	}

	// Check for collisions with other user-defined types
	if reservedBy, reserved := ast.Types[name]; reserved {
		reservedBySrcNode := reservedBy.Src()
		return errors.Errorf(
			"Redeclaration of type %s at %d:%d "+
				"(previous declaration: %d:%d (%s))",
			name,
			srcNode.Begin,
			srcNode.End,
			reservedBySrcNode.Begin,
			reservedBySrcNode.End,
			reservedBy.Category(),
		)
	}

	// Define
	ast.Types[name] = newType

	return nil
}
