package compiler

import "fmt"

// defineType returns an error if the type name is already reserved
func (c *Compiler) defineType(newType Type) (TypeID, Error) {
	// Check for collisions with reserved primitive types
	srcNode := newType.Src()
	name := newType.Name()
	if stdTypeByName(name) != nil {
		return 0, cErr{
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
	if reservedBy := c.ast.FindTypeByName("", name); reservedBy != nil {
		reservedBySrcNode := reservedBy.Src()
		return 0, cErr{
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

	// Increment last issued type ID
	c.lastIssuedTypeID += TypeID(1)

	// Define a new type
	c.ast.Types = append(c.ast.Types, newType)
	c.typeByID[c.lastIssuedTypeID] = newType

	// Define in sub-category
	switch newType.Category() {
	case TypeCategoryAlias:
		c.ast.AliasTypes = append(c.ast.AliasTypes, newType)
	case TypeCategoryEnum:
		c.ast.EnumTypes = append(c.ast.EnumTypes, newType)
	case TypeCategoryUnion:
		c.ast.UnionTypes = append(c.ast.UnionTypes, newType)
	case TypeCategoryStruct:
		c.ast.StructTypes = append(c.ast.StructTypes, newType)
	}

	return c.lastIssuedTypeID, nil
}
