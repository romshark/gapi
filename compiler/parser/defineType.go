package parser

import "fmt"

// defineType returns an error if the type name is already reserved
func (pr *Parser) defineType(newType Type) {
	src := newType.Source()
	name := newType.Name()

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
				"(previous declaration: %s (%s))",
				name,
				reservedBySrcNode.Begin(),
				reservedBy.Category(),
			),
		})
		return
	}

	// Increment last issued type ID
	pr.lastIssuedTypeID += TypeID(1)
	newID := pr.lastIssuedTypeID

	// Define a new type
	pr.ast.Types = append(pr.ast.Types, newType)
	pr.typeByID[newID] = newType
	pr.typeByName[name] = newType

	// Set ID and define in sub-category
	switch t := newType.(type) {
	case *TypeAlias:
		t.terminalType.ID = newID
		pr.ast.AliasTypes = append(pr.ast.AliasTypes, newType)
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
