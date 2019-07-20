package compiler

import "fmt"

// parseEnmVals parses the values of an enumeration type declaration
// returns true if the values are valid, otherwise returns false
func (c *Compiler) parseEnmVals(
	enumType *TypeEnum,
	first *node32,
) (valid bool, err error) {
	valid = true
	nd := first

	// Parse values
	for nd != nil {
		valueName := c.getSrc(nd)

		if err := verifyLowerCamelCase(valueName); err != nil {
			c.err(cErr{
				ErrEnumValIllegalIdent,
				fmt.Sprintf(
					"invalid enum value identifier at %d:%d: %s",
					nd.begin,
					nd.end,
					err,
				),
			})
			valid = false
			goto NEXT
		}

		// Check for duplicate values
		if defined, isDefined := enumType.Values[valueName]; isDefined {
			c.err(cErr{
				ErrEnumValRedecl,
				fmt.Sprintf(
					"Redeclaration of enum value %s at %d:%d "+
						"(previously declared at %d:%d)",
					valueName,
					nd.begin,
					nd.end,
					defined.Begin,
					defined.End,
				),
			})
			valid = false
			goto NEXT
		}

		// Add enum value
		enumType.Values[valueName] = EnumValue{
			Src: Src{
				Begin: nd.begin,
				End:   nd.end,
			},
			Name: valueName,
		}

	NEXT:
		nd = skipUntil(nd.next, ruleWrd)
	}
	return
}
