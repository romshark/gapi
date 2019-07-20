package compiler

import (
	"fmt"
)

func (c *Compiler) parseDeclAls(node *node32) error {
	// Read name
	nd := skipUntil(node.up, ruleWrd)
	newAliasTypeName := c.getSrc(nd)

	if err := verifyTypeName(newAliasTypeName); err != nil {
		c.err(cErr{
			ErrTypeIllegalIdent,
			fmt.Sprintf("illegal type identifier at %d:%d: %s",
				nd.begin,
				nd.end,
				err,
			),
		})
		return nil
	}

	nd = skipUntil(nd.next, ruleWrd)
	aliasedTypeName := c.getSrc(nd)

	if err := verifyTypeName(aliasedTypeName); err != nil {
		c.err(cErr{
			ErrTypeIllegalIdent,
			fmt.Sprintf("illegal type identifier at %d:%d: %s",
				nd.begin,
				nd.end,
				err,
			),
		})
		return nil
	}

	newType := &TypeAlias{
		terminalType: terminalType{
			src:  src(node),
			name: newAliasTypeName,
		},
	}

	// Try to define the type
	typeID, err := c.defineType(newType)
	if err != nil {
		c.err(err)
		return nil
	}
	newType.id = typeID

	c.deferJob(func() error {
		// Ensure the aliased type exists after all types have been defined
		aliasedType := c.ast.FindTypeByName("", aliasedTypeName)
		if aliasedType == nil {
			c.err(cErr{ErrTypeUndef, fmt.Sprintf(
				"undefined type %s aliased by %s at %d:%d",
				aliasedTypeName,
				newAliasTypeName,
				node.begin,
				node.end,
			)})
			return nil
		}

		// Reference the aliased type
		newType.AliasedType = aliasedType

		return nil
	})

	return nil
}
