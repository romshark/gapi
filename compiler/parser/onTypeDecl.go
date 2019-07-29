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

	// Issue a new type ID
	pr.lastIssuedTypeID += TypeID(1)
	newID := pr.lastIssuedTypeID

	// Register the newly defined type
	pr.mod.Types = append(pr.mod.Types, newType)
	pr.typeByID[newID] = newType
	pr.typeByName[name] = newType

	// Set ID and define in the schema model
	switch t := newType.(type) {
	case *TypeAlias:
		t.terminalType.ID = newID
		pr.mod.AliasTypes = append(pr.mod.AliasTypes, newType)
	case *TypeEnum:
		t.terminalType.ID = newID
		pr.mod.EnumTypes = append(pr.mod.EnumTypes, newType)
	case *TypeUnion:
		t.terminalType.ID = newID
		pr.mod.UnionTypes = append(pr.mod.UnionTypes, newType)
	case *TypeStruct:
		t.terminalType.ID = newID
		pr.mod.StructTypes = append(pr.mod.StructTypes, newType)
	case *TypeResolver:
		t.terminalType.ID = newID
		pr.mod.ResolverTypes = append(pr.mod.ResolverTypes, newType)
	}
}
