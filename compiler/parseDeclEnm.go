package compiler

import (
	"fmt"
)

func (c *Compiler) parseDeclEnm(node *node32) error {
	// Read name
	nd := skipUntil(node.up, ruleWrd)
	newEnumTypeName := c.getSrc(nd)

	newType := &TypeEnum{
		terminalType: terminalType{
			src:  src(node),
			name: newEnumTypeName,
		},
		Values: make(map[string]EnumValue),
	}

	// Parse values
	nd = skipUntil(nd.next, ruleBlkEn)
	nd = skipUntil(nd.up, ruleWrd)
	checkVals, err := c.parseEnmVals(newType, nd)
	if err != nil {
		return err
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
	c.defineType(newType)

	return nil
}
