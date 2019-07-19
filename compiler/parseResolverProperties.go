package compiler

import "fmt"

// parseResolverProperties parses the properties of a resolver declaration
// reports returns true if the properties are valid, otherwise returns false
func (c *Compiler) parseResolverProperties(
	resolver *TypeResolver,
	node *node32,
) (valid bool, err error) {
	valid = true
	for node != nil {
		nodeProp := node
		node = skipUntil(node.up, ruleWrd)
		propName := c.getSrc(node)

		var newProp *ResolverProperty

		// Verify property identifier
		if err := verifyResolverPropIdent(propName); err != nil {
			c.err(cErr{
				ErrResolverPropIllegalIdent,
				fmt.Sprintf(
					"invalid resolver property identifier at %d:%d: %s",
					node.begin,
					node.end,
					err,
				),
			})
			valid = false
			goto NEXT_PROP
		}

		// Check for redeclared properties
		if prop := resolver.PropertyByName(propName); prop != nil {
			c.err(cErr{
				ErrResolverPropRedecl,
				fmt.Sprintf(
					"Redeclaration of resolver property %s at %d:%d "+
						"(previously declared at %d:%d)",
					propName,
					node.begin,
					node.end,
					prop.Begin,
					prop.End,
				),
			})
			valid = false
			goto NEXT_PROP
		}

		// Add property
		newProp = &ResolverProperty{
			Src: Src{
				Begin: node.begin,
				End:   node.end,
			},
			Resolver: resolver,
			Name:     propName,
			Type:     nil, // Deferred
		}
		newProp.GraphID = c.defineGraphNode(newProp)
		resolver.Properties = append(resolver.Properties, newProp)

		// Parse property parameters
		if err = c.parseResolverPropertyParameters(
			newProp,
			skipUntil(node.next, rulePrms),
		); err != nil {
			return
		}

		// Parse the property type in deferred mode
		c.deferJob(func() error {
			// Read property type
			nodePropType := skipUntil(nodeProp.up, ruleTp)
			propType, err := c.parseType(nodePropType)
			if err != nil {
				c.err(err)
			}

			// Set the property type
			newProp.Type = propType

			return nil
		})

	NEXT_PROP:
		node = skipUntil(nodeProp.next, ruleRvPrp)
	}
	return
}
