package compiler

import (
	"fmt"
)

func (c *Compiler) parseDeclMut(node *node32) error {
	// Read name
	nd := skipUntil(node.up, ruleWrd)
	newMutationName := c.getSrc(nd)

	newMutation := &Mutation{
		Src:  src(node),
		Name: newMutationName,
	}

	// Parse properties
	if err := c.parseParams(
		newMutation,
		skipUntil(nd, rulePrms),
	); err != nil {
		return err
	}

	nd = skipUntil(nd.next, ruleTp)
	resultTypeName := c.getSrc(nd)

	c.deferJob(func() error {
		// Ensure the type of the mutation exists
		resultType := c.findTypeByName(resultTypeName)
		if resultType == nil {
			c.err(cErr{ErrTypeUndef, fmt.Sprintf(
				"undefined type %s referenced by mutation %s at %d:%d",
				resultTypeName,
				newMutationName,
				node.begin,
				node.end,
			)})
			return nil
		}

		// Reference the aliased type
		newMutation.Type = resultType

		return nil
	})

	// Try to define the endpoint
	c.defineGraphNode(newMutation)

	return nil
}
