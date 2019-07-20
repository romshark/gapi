package compiler

import "fmt"

// parseParams parses a list of parameters if there are any
func (c *Compiler) parseParams(
	target GraphNode,
	node *node32,
) error {
	nodeCur := node
	if nodeCur == nil {
		// No parameters
		return nil
	}
	if nodeCur.pegRule != rulePrms {
		return fmt.Errorf(
			"unexpected rule: %s (expected %s)",
			rul3s[nodeCur.pegRule],
			rul3s[rulePrms],
		)
	}
	nodeCur = skipUntil(nodeCur.up, rulePrmsBd)
	nodeCur = skipUntil(nodeCur.up, rulePrm)
	for nodeCur != nil {
		nodeParam := nodeCur

		// Read parameter name
		nodeCur = skipUntil(nodeCur.up, ruleWrd)
		paramName := c.getSrc(nodeCur)

		// Read parameter type
		nodeCur = skipUntil(nodeCur.next, ruleTp)
		nodeParamType := nodeCur

		newParam := &Parameter{
			Src:    src(nodeParam),
			Target: target,
			Name:   paramName,
			Type:   nil, // deferred
		}

		// Parse the parameter type in deferred mode
		c.deferJob(func() error {
			paramType, err := c.parseType(nodeParamType)
			if err != nil {
				c.err(err)
			}

			// Set the parameter type
			newParam.Type = paramType

			return nil
		})

		// Reference the new parameter in its target graph node
		switch target := target.(type) {
		case *QueryEndpoint:
			target.Parameters = append(target.Parameters, newParam)
		case *ResolverProperty:
			target.Parameters = append(target.Parameters, newParam)
		case *Mutation:
			target.Parameters = append(target.Parameters, newParam)
		}

		// Define the new parameter
		c.defineParameter(newParam)

		// Pick next parameter if any
		nodeCur = skipUntil(nodeParam.next, rulePrm)
		if nodeCur == nil {
			nodeCur = skipUntil(nodeParam.next, rulePrmsBd)
			if nodeCur != nil {
				nodeCur = nodeCur.up
			}
		}
	}
	return nil
}
