package parser

import "fmt"

// onTypeDecl is executed when a type declaration was parsed.
// It check the new type for name collisions and redeclarations
// and registers it in the parser's context if necessary
func (pr *Parser) onTypeDecl(newType Type) {
	src := newType.Source()
	name := newType.String()

	// Check for collisions with reserved primitive types
	if stdTypeByName(name) != nil {
		pr.err(&pErr{
			at:   src.Begin(),
			code: ErrTypeRedecl,
			message: fmt.Sprintf(
				"Redeclaration of type %s (reserved primitive type)",
				name,
			),
		})
		return
	}

	// Check for collisions with other user-defined types
	if reservedBy, isReserved := pr.typeByName[name]; isReserved {
		reservedBySrcNode := reservedBy.Source()
		pr.err(&pErr{
			at:   src.Begin(),
			code: ErrTypeRedecl,
			message: fmt.Sprintf("Redeclaration of type %s "+
				"(previous declaration at %s)",
				name,
				reservedBySrcNode.Begin(),
			),
		})
		return
	}

	// Don't issue type IDs to alias types
	if alias, isAlias := newType.(*TypeAlias); isAlias {
		// Issue a new alias type ID
		pr.lastIssuedAliasTypeID += TypeID(1)
		newID := pr.lastIssuedAliasTypeID
		alias.ID = newID

		pr.aliasByName[name] = alias
		pr.aliasByID[newID] = alias
		return
	}

	// Issue a new type ID
	pr.lastIssuedTypeID += TypeID(1)
	newID := pr.lastIssuedTypeID

	// Register the newly defined type
	pr.ast.Types = append(pr.ast.Types, newType)
	pr.typeByID[newID] = newType
	pr.typeByName[name] = newType

	// Set ID and define in the AST
	switch t := newType.(type) {
	case *TypeEnum:
		t.terminalType.ID = newID
		pr.ast.EnumTypes = append(pr.ast.EnumTypes, newType)
	case *TypeUnion:
		t.terminalType.ID = newID
		pr.ast.UnionTypes = append(pr.ast.UnionTypes, newType)
	case *TypeStruct:
		t.terminalType.ID = newID
		pr.ast.StructTypes = append(pr.ast.StructTypes, newType)
	case *TypeResolver:
		t.terminalType.ID = newID
		pr.ast.ResolverTypes = append(pr.ast.ResolverTypes, newType)
	}
}
