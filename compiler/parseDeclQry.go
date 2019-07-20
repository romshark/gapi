package compiler

func (c *Compiler) parseDeclQry(node *node32) error {
	// Read name
	nd := skipUntil(node.up, ruleWrd)
	newQueryEndpointName := c.getSrc(nd)

	newQueryEndpoint := &QueryEndpoint{
		Src:  src(node),
		Name: newQueryEndpointName,
	}

	// Parse properties
	nd = skipUntil(nd, rulePrms)
	if err := c.parseParams(
		newQueryEndpoint,
		nd,
	); err != nil {
		return err
	}

	// Try to define the endpoint
	c.defineGraphNode(newQueryEndpoint)

	return nil
}
