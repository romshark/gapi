package compiler

import "fmt"

// parseRsvPropParams parses the parameters of a resolver property
// if there are any
func (c *Compiler) parseRsvPropParams(
	prop *ResolverProperty,
	first *node32,
) error {
	nodeCur := first
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
			Target: prop,
			Name:   paramName,
			Type:   nil, // deferred
		}

		// Verify parameter identifier
		if err := verifyParameterIdent(paramName); err != nil {
			c.err(cErr{
				ErrParamIllegalIdent,
				fmt.Sprintf(
					"invalid resolver property parameter identifier "+
						"at %d:%d: %s",
					nodeCur.begin,
					nodeCur.end,
					err,
				),
			})
			goto NEXT_PARAM
		}

		// Check for redeclared parameters
		if param := prop.ParamByName(paramName); param != nil {
			c.err(cErr{
				ErrResolverPropParamRedecl,
				fmt.Sprintf(
					"Redeclaration of resolver property (%s) parameter %s "+
						"at %d:%d (previously declared at %d:%d)",
					prop.Name,
					paramName,
					nodeCur.begin,
					nodeCur.end,
					param.Begin,
					param.End,
				),
			})
			goto NEXT_PARAM
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

		// Define the parameter and reference it in the resolver property
		newParam.ID = c.defineParameter(newParam)
		prop.Parameters = append(prop.Parameters, newParam)

	NEXT_PARAM:
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
