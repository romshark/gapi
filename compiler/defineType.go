package compiler

import "fmt"

// defineType returns an error if the type name is already reserved
func (c *Compiler) defineType(newType Type) {
	src := newType.Source()
	name := newType.Name()

	// Verify the type identifier
	if err := verifyCapitalizedCamelCase(name); err != nil {
		c.err(cErr{
			ErrTypeIllegalIdent,
			fmt.Sprintf("illegal type identifier %d:%d: %s",
				src.Begin,
				src.End,
				err,
			),
		})
		return
	}

	// Check for collisions with reserved primitive types
	if stdTypeByName(name) != nil {
		c.err(cErr{
			ErrTypeRedecl,
			fmt.Sprintf(
				"Redeclaration of type %s at %d:%d (reserved primitive type)",
				name,
				src.Begin,
				src.End,
			),
		})
		return
	}

	// Check for collisions with other user-defined types
	if reservedBy := c.ast.FindTypeByName("", name); reservedBy != nil {
		reservedBySrcNode := reservedBy.Source()
		c.err(cErr{
			ErrTypeRedecl,
			fmt.Sprintf("Redeclaration of type %s at %d:%d "+
				"(previous declaration: %d:%d (%s))",
				name,
				src.Begin,
				src.End,
				reservedBySrcNode.Begin,
				reservedBySrcNode.End,
				reservedBy.Category(),
			),
		})
		return
	}

	// Increment last issued type ID
	c.lastIssuedTypeID += TypeID(1)
	newID := c.lastIssuedTypeID

	// Define a new type
	c.ast.Types = append(c.ast.Types, newType)
	c.typeByID[newID] = newType
	c.typeByName[name] = newType

	// Set ID and define in sub-category
	switch t := newType.(type) {
	case *TypeAlias:
		t.id = newID
		c.ast.AliasTypes = append(c.ast.AliasTypes, newType)
	case *TypeEnum:
		t.id = newID
		c.ast.EnumTypes = append(c.ast.EnumTypes, newType)
	case *TypeUnion:
		t.id = newID
		c.ast.UnionTypes = append(c.ast.UnionTypes, newType)
	case *TypeStruct:
		t.id = newID
		c.ast.StructTypes = append(c.ast.StructTypes, newType)
	case *TypeResolver:
		t.id = newID
		c.ast.ResolverTypes = append(c.ast.ResolverTypes, newType)
	}
}
