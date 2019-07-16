package compiler

import (
	"fmt"
)

// AST represents the abstract GAPI syntax tree
type AST struct {
	Types          []Type
	AliasTypes     []Type
	EnumTypes      []Type
	UnionTypes     []Type
	QueryEndpoints []QueryEndpoint
	Mutations      []Mutation
	SchemaName     string
}

// Clone returns a copy of the abstract syntax tree
func (ast *AST) Clone() *AST {
	if ast == nil {
		return nil
	}

	types := make([]Type, len(ast.Types))
	copy(types, ast.Types)

	aliasTypes := make([]Type, len(ast.AliasTypes))
	copy(aliasTypes, ast.AliasTypes)

	enumTypes := make([]Type, len(ast.EnumTypes))
	copy(enumTypes, ast.EnumTypes)

	unionTypes := make([]Type, len(ast.UnionTypes))
	copy(unionTypes, ast.UnionTypes)

	queryEndpoints := make([]QueryEndpoint, len(ast.QueryEndpoints))
	copy(queryEndpoints, ast.QueryEndpoints)

	mutations := make([]Mutation, len(ast.Mutations))
	copy(mutations, ast.Mutations)

	return &AST{
		Types:          types,
		AliasTypes:     aliasTypes,
		EnumTypes:      enumTypes,
		UnionTypes:     unionTypes,
		QueryEndpoints: queryEndpoints,
		Mutations:      mutations,
		SchemaName:     ast.SchemaName,
	}
}

// FindTypeByName returns a type given its category and name
func (ast *AST) FindTypeByName(category TypeCategory, name string) Type {
	findUserDefined := func() Type {
		for _, tp := range ast.Types {
			if tp.Name() == name {
				return tp
			}
		}
		return nil
	}

	switch category {
	case TypeCategoryPrimitive:
		// Search in primitives only
		return stdTypeByName(name)
	case TypeCategoryUserDefined:
		// Search in all categories including primitives
		return findUserDefined()
	case "":
		// Search in all user-defined types
		if t := stdTypeByName(name); t != nil {
			return t
		}
		return findUserDefined()
	default:
		// Search in specific category
		for _, tp := range ast.Types {
			if tp.Category() == category && tp.Name() == name {
				return tp
			}
		}
	}
	return nil
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
	if reservedBy := ast.FindTypeByName("", name); reservedBy != nil {
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
	ast.Types = append(ast.Types, newType)

	// Define in sub-category
	switch newType.Category() {
	case TypeCategoryAlias:
		ast.AliasTypes = append(ast.AliasTypes, newType)
	case TypeCategoryEnum:
		ast.EnumTypes = append(ast.EnumTypes, newType)
	case TypeCategoryUnion:
		ast.UnionTypes = append(ast.UnionTypes, newType)
	}

	return nil
}
