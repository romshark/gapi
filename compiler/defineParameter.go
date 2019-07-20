package compiler

import (
	"fmt"
)

// defineParameter returns validates the name, checks for redeclaration
// and defines a new parameter assigning it a unique identifier
func (c *Compiler) defineParameter(newParam *Parameter) {
	targetParams := c.paramsByName[newParam.Target]

	// Verify parameter identifier
	if err := verifyLowerCamelCase(newParam.Name); err != nil {
		c.err(cErr{
			ErrParamIllegalIdent,
			fmt.Sprintf(
				"invalid parameter identifier at %d:%d: %s",
				newParam.Begin,
				newParam.End,
				err,
			),
		})
		return
	}

	// Check for redeclared parameters
	if targetParams != nil {
		if defined, isDefined := targetParams[newParam.Name]; isDefined {
			c.err(cErr{
				ErrResolverPropParamRedecl,
				fmt.Sprintf(
					"Redeclaration parameter %s "+
						"at %d:%d (previously declared at %d:%d)",
					newParam.Name,
					newParam.Begin,
					newParam.End,
					defined.Begin,
					defined.End,
				),
			})
			return
		}
	}

	// Register a new parameter
	c.lastIssuedParamID += ParamID(1)
	newParam.ID = c.lastIssuedParamID

	if targetParams == nil {
		c.paramsByName[newParam.Target] = map[string]*Parameter{
			newParam.Name: newParam,
		}
	} else {
		c.paramsByName[newParam.Target][newParam.Name] = newParam
	}

	c.paramByID[newParam.ID] = newParam
}
