package compiler

import (
	"fmt"
)

// parseStrFlds parses the fields of a struct declaration
// returns true if the fields are valid, otherwise returns false
func (c *Compiler) parseStrFlds(
	structType *TypeStruct,
	first *node32,
) (valid bool, err error) {
	valid = true
	nd := first
	for nd != nil {
		var newField *StructField
		var ndFldTp *node32
		ndField := nd

		// Read field name
		nd = ndField.up
		fieldName := c.getSrc(nd)

		// Verify field identifier
		if err := verifyStructFieldIdent(fieldName); err != nil {
			c.err(cErr{
				ErrStructFieldIllegalIdent,
				fmt.Sprintf(
					"invalid struct field identifier at %d:%d: %s",
					nd.begin,
					nd.end,
					err,
				),
			})
			valid = false
			goto NEXT
		}

		// Check for redeclared fields
		if field := structType.FieldByName(fieldName); field != nil {
			c.err(cErr{
				ErrStructFieldRedecl,
				fmt.Sprintf(
					"Redeclaration of struct field %s at %d:%d "+
						"(previously declared at %d:%d)",
					fieldName,
					nd.begin,
					nd.end,
					field.Begin,
					field.End,
				),
			})
			valid = false
			goto NEXT
		}

		// Add field
		newField = &StructField{
			Src: Src{
				Begin: nd.begin,
				End:   nd.end,
			},
			Struct: structType,
			Name:   fieldName,
			Type:   nil, // Deferred
		}
		newField.GraphID = c.defineGraphNode(newField)
		structType.Fields = append(structType.Fields, newField)

		nd = skipUntil(nd.next, ruleTp)
		ndFldTp = nd

		// Parse the field type in deferred mode
		c.deferJob(func() error {
			fieldType, err := c.parseType(ndFldTp)
			if err != nil {
				c.err(err)
			}

			// Set the field type
			newField.Type = fieldType

			return nil
		})

	NEXT:
		nd = skipUntil(ndField.next, ruleStFld)
	}
	return
}
