package compiler

import "fmt"

func (c *Compiler) parseDeclUni(node *node32) error {
	// Read name
	nd := skipUntil(node.up, ruleWrd)
	newUnionTypeName := c.getSrc(nd)

	newType := &TypeUnion{
		terminalType: terminalType{
			src:  src(node),
			name: newUnionTypeName,
		},
		Types: make(map[string]Type),
	}

	// Parse types
	nd = skipUntil(nd.next, ruleBlkUn)
	nd = skipUntil(nd.up, ruleWrd)
	checkOpts, err := c.parseUniTps(newType, nd)
	if err != nil {
		return err
	}

	if checkOpts && len(newType.Types) < 2 {
		c.err(cErr{
			ErrUnionMissingOpts,
			fmt.Sprintf(
				"union %s requires at least two types at %d:%d",
				newUnionTypeName,
				node.begin,
				node.end,
			),
		})
		return nil
	}

	c.deferJob(func() error {
		// Ensure all referenced types are defined and legal
		for name := range newType.Types {
			reg := c.findTypeByName(name)
			if reg == nil {
				c.err(cErr{
					ErrTypeUndef,
					fmt.Sprintf(
						"undefined type %s referenced "+
							"in union type %s at %d:%d ",
						name,
						newUnionTypeName,
						node.begin,
						node.end,
					),
				})
				continue
			}
			if name == newUnionTypeName {
				c.err(cErr{
					ErrUnionSelfref,
					fmt.Sprintf(
						"union type %s references itself at %d:%d",
						newUnionTypeName,
						node.begin,
						node.end,
					),
				})
				continue
			}
			if _, isNone := reg.(TypeStdNone); isNone {
				c.err(cErr{
					ErrUnionIncludesNone,
					fmt.Sprintf(
						"union type %s includes the None primitive at %d:%d",
						newUnionTypeName,
						node.begin,
						node.end,
					),
				})
				continue
			}
			newType.Types[name] = reg
		}

		return nil
	})

	// Try to define the type
	c.defineType(newType)

	return nil
}
