package compiler

import (
	"fmt"
)

// parseUniTps parses the types of a union type declaration
// returns true if the types are valid, otherwise returns false
func (c *Compiler) parseUniTps(
	unionType *TypeUnion,
	first *node32,
) (valid bool, err error) {
	valid = true
	nd := first
	for nd != nil {
		referencedTypeName := c.getSrc(nd)

		if err := verifyTypeName(referencedTypeName); err != nil {
			c.err(cErr{
				ErrTypeIllegalIdent,
				fmt.Sprintf(
					"invalid union option-type identifier at %d:%d: %s",
					nd.begin,
					nd.end,
					err,
				),
			})
			valid = false
			goto NEXT
		}

		// Check for duplicate values
		if _, isDefined := unionType.Types[referencedTypeName]; isDefined {
			c.err(cErr{
				ErrUnionRedund,
				fmt.Sprintf(
					"multiple references to the same type (%s) "+
						"in union type %s at %d:%d ",
					referencedTypeName,
					unionType.name,
					nd.begin,
					nd.end,
				),
			})
			valid = false
			goto NEXT
		}

		// Mark type for deferred checking
		unionType.Types[referencedTypeName] = nil

	NEXT:
		nd = skipUntil(nd.next, ruleWrd)
	}
	return
}
