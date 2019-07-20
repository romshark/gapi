package compiler

import (
	"fmt"
)

func (c *Compiler) parseDeclStr(node *node32) error {
	nd := skipUntil(node.up, ruleWrd)
	newStructTypeName := c.getSrc(nd)

	newType := &TypeStruct{
		terminalType: terminalType{
			src:  src(node),
			name: newStructTypeName,
		},
		Fields: make([]*StructField, 0),
	}

	// Parse fields
	nd = skipUntil(nd.next, ruleBlkSt)
	nd = skipUntil(nd.up, ruleStFld)
	checkFields, err := c.parseStrFlds(newType, nd)
	if err != nil {
		return err
	}

	if checkFields && len(newType.Fields) < 1 {
		c.err(cErr{
			ErrStructNoFields,
			fmt.Sprintf(
				"struct %s is missing fields at %d:%d",
				newStructTypeName,
				node.begin,
				node.end,
			),
		})
	}

	// Try to define the type
	c.defineType(newType)

	return nil
}
