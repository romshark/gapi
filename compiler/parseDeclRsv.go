package compiler

import (
	"fmt"
)

func (c *Compiler) parseDeclRsv(node *node32) error {
	// Read name
	nd := skipUntil(node.up, ruleWrd)
	newResolverTypeName := c.getSrc(nd)

	if err := verifyTypeName(newResolverTypeName); err != nil {
		c.err(cErr{
			ErrTypeIllegalIdent,
			fmt.Sprintf("illegal resolver type identifier %d:%d: %s",
				nd.begin,
				nd.end,
				err,
			),
		})
		return nil
	}

	newResolver := &TypeResolver{
		terminalType: terminalType{
			src:  src(node),
			name: newResolverTypeName,
		},
		Properties: make([]*ResolverProperty, 0),
	}

	// Parse properties
	nd = skipUntil(nd, ruleBlkRv)
	nd = skipUntil(nd.up, ruleRvPrp)
	checkProps, err := c.parseRsvProps(newResolver, nd)
	if err != nil {
		return err
	}

	if checkProps && len(newResolver.Properties) < 1 {
		c.err(cErr{
			ErrResolverNoProps,
			fmt.Sprintf(
				"resolver %s is missing properties at %d:%d",
				newResolverTypeName,
				node.begin,
				node.end,
			),
		})
	}

	// Try to define the type
	typeID, typeDefErr := c.defineType(newResolver)
	if typeDefErr != nil {
		c.err(typeDefErr)
	}
	newResolver.id = typeID

	return nil
}
