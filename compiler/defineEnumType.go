package compiler

import (
	"fmt"
)

func (c *Compiler) defineEnumType(node *node32) error {
	current := node.up.next.next
	newEnumTypeName := getSrc(c.parser.Buffer, node.up.next.next)

	if err := verifyTypeName(newEnumTypeName); err != nil {
		c.err(cErr{
			ErrTypeIllegalIdent,
			fmt.Sprintf("illegal enum type identifier %d:%d: %s",
				current.begin,
				current.end,
				err,
			),
		})
		return nil
	}

	newType := &TypeEnum{
		terminalType: terminalType{
			src:  src(node),
			name: newEnumTypeName,
		},
		Values: make(map[string]EnumValue),
	}

	checkVals := true

	// Parse values
	current = current.next.next.up.next.next
	for current != nil {
		valueName := c.parser.Buffer[current.begin:current.end]

		if err := verifyEnumValue(valueName); err != nil {
			c.err(cErr{
				ErrEnumValIllegalIdent,
				fmt.Sprintf(
					"invalid enum value identifier at %d:%d: %s",
					current.begin,
					current.end,
					err,
				),
			})
			checkVals = false
			goto NEXT
		}

		// Check for duplicate values
		if defined, isDefined := newType.Values[valueName]; isDefined {
			c.err(cErr{
				ErrEnumValRedecl,
				fmt.Sprintf(
					"Redeclaration of enum value %s at %d:%d "+
						"(previously declared at %d:%d)",
					valueName,
					current.begin,
					current.end,
					defined.Begin,
					defined.End,
				),
			})
			checkVals = false
			goto NEXT
		}

		// Add enum value
		newType.Values[valueName] = EnumValue{
			Src: Src{
				Begin: current.begin,
				End:   current.end,
			},
			Name: valueName,
		}

	NEXT:
		next := current.next.next
		if next == nil || next.pegRule == ruleBLKE {
			break
		}
		current = next
	}

	if checkVals && len(newType.Values) < 1 {
		c.err(cErr{
			ErrEnumNoVal,
			fmt.Sprintf(
				"enum %s is missing values at %d:%d",
				newEnumTypeName,
				node.begin,
				node.end,
			),
		})
		return nil
	}

	// Try to define the type
	typeID, err := c.defineType(newType)
	if err != nil {
		c.err(err)
	}
	newType.id = typeID

	return nil
}
