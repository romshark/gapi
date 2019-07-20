package compiler

import (
	"fmt"
)

func (c *Compiler) parseDeclEnm(node *node32) error {
	nd := skipUntil(node.up, ruleWrd)
	newEnumTypeName := c.getSrc(nd)

	if err := verifyTypeName(newEnumTypeName); err != nil {
		c.err(cErr{
			ErrTypeIllegalIdent,
			fmt.Sprintf("illegal enum type identifier %d:%d: %s",
				nd.begin,
				nd.end,
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
	typeID, typeDefErr := c.defineType(newType)
	if typeDefErr != nil {
		c.err(typeDefErr)
	}
	newType.id = typeID

	return nil
}
