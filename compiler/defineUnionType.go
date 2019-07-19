package compiler

import "fmt"

func (c *Compiler) defineUnionType(node *node32) error {
	current := node.up.next.next
	newUnionTypeName := c.getSrc(current)

	if err := verifyTypeName(newUnionTypeName); err != nil {
		c.err(cErr{
			ErrTypeIllegalIdent,
			fmt.Sprintf(
				"invalid union type identifier %d:%d: %s",
				current.begin,
				current.end,
				err,
			),
		})
		return nil
	}

	newType := &TypeUnion{
		terminalType: terminalType{
			src:  src(node),
			name: newUnionTypeName,
		},
		Types: make(map[string]Type),
	}

	// Parse types
	current = current.next.next.up.next.next
	checkOpts := true
	for {
		referencedTypeName := c.parser.Buffer[current.begin:current.end]

		if err := verifyTypeName(referencedTypeName); err != nil {
			c.err(cErr{
				ErrTypeIllegalIdent,
				fmt.Sprintf(
					"invalid union option-type identifier at %d:%d: %s",
					current.begin,
					current.end,
					err,
				),
			})
			checkOpts = false
			goto NEXT
		}

		// Check for duplicate values
		if _, isDefined := newType.Types[referencedTypeName]; isDefined {
			c.err(cErr{
				ErrUnionRedund,
				fmt.Sprintf(
					"multiple references to the same type (%s) "+
						"in union type %s at %d:%d ",
					referencedTypeName,
					newUnionTypeName,
					node.begin,
					node.end,
				),
			})
			checkOpts = false
			goto NEXT
		}

		// Mark type for deferred checking
		newType.Types[referencedTypeName] = nil

	NEXT:
		next := current.next.next
		if next == nil || next.pegRule == ruleBLKE {
			break
		}
		current = next
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
			reg := c.ast.FindTypeByName("", name)
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
	typeID, err := c.defineType(newType)
	if err != nil {
		c.err(err)
	}
	newType.id = typeID

	return nil
}
