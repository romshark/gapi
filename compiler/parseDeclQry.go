package compiler

import (
	"fmt"
)

func (c *Compiler) parseDeclQry(node *node32) error {
	// Read name
	nd := skipUntil(node.up, ruleWrd)
	newQueryEndpointName := c.getSrc(nd)

	newQueryEndpoint := &QueryEndpoint{
		Src:  src(node),
		Name: newQueryEndpointName,
	}

	// Parse properties
	if err := c.parseParams(
		newQueryEndpoint,
		skipUntil(nd, rulePrms),
	); err != nil {
		return err
	}

	nd = skipUntil(nd.next, ruleTp)
	resultTypeName := c.getSrc(nd)

	c.deferJob(func() error {
		// Ensure the type of the query endpoint exists
		resultType := c.findTypeByName(resultTypeName)
		if resultType == nil {
			c.err(cErr{ErrTypeUndef, fmt.Sprintf(
				"undefined type %s referenced by query endpoint %s at %d:%d",
				resultTypeName,
				newQueryEndpointName,
				node.begin,
				node.end,
			)})
			return nil
		}

		// Reference the aliased type
		newQueryEndpoint.Type = resultType

		return nil
	})

	// Try to define the endpoint
	c.defineGraphNode(newQueryEndpoint)

	return nil
}
