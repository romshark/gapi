package compiler

import (
	"fmt"
)

func (c *Compiler) defineStructType(node *node32) error {
	current := node.up.next.next
	newStructTypeName := getSrc(c.parser.Buffer, node.up.next.next)

	if err := verifyTypeName(newStructTypeName); err != nil {
		c.err(cErr{
			ErrTypeIllegalIdent,
			fmt.Sprintf("illegal struct type identifier %d:%d: %s",
				current.begin,
				current.end,
				err,
			),
		})
		return nil
	}

	newType := &TypeStruct{
		terminalType: terminalType{
			src:  src(node),
			name: newStructTypeName,
		},
		Fields: make([]*StructField, 0),
	}

	// Parse fields
	current = current.next.next.up.next.next
	for {
		field := current
		fieldNameNode := field.up
		fieldTypeNode := fieldNameNode.next.next
		fieldName := c.parser.Buffer[fieldNameNode.begin:fieldNameNode.end]

		var newField *StructField

		// Verify field identifier
		if err := verifyStructFieldIdent(fieldName); err != nil {
			c.err(cErr{
				ErrStructFieldIllegalIdent,
				fmt.Sprintf(
					"invalid struct field identifier at %d:%d: %s",
					current.begin,
					current.end,
					err,
				),
			})
			goto NEXT
		}

		// Check for redeclared fields
		if field := newType.FieldByName(fieldName); field != nil {
			c.err(cErr{
				ErrStructFieldRedecl,
				fmt.Sprintf(
					"Redeclaration of struct field %s at %d:%d "+
						"(previously declared at %d:%d)",
					fieldName,
					current.begin,
					current.end,
					field.Begin,
					field.End,
				),
			})
			goto NEXT
		}

		// Add field
		newField = &StructField{
			Src: Src{
				Begin: current.begin,
				End:   current.end,
			},
			Name: fieldName,
			Type: nil, // Deferred
		}
		newType.Fields = append(newType.Fields, newField)

		// Parse the field type in deferred mode
		c.deferJob(func() error {
			fieldType, err := c.parseType(fieldTypeNode)
			if err != nil {
				c.err(err)
			}

			// Set the field type
			newField.Type = fieldType

			return nil
		})

	NEXT:
		next := current.next.next
		if next == nil || next.pegRule == ruleBLKE {
			break
		}
		current = next
	}

	// Try to define the type
	if err := c.ast.defineType(newType); err != nil {
		c.err(err)
	}

	return nil
}
