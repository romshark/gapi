package compiler

// defineGraphNode returns a new unique graph node identifier
func (c *Compiler) defineGraphNode(newNode GraphNode) GraphNodeID {
	// Update the last issued ID
	c.lastIssuedGraphID += GraphNodeID(1)
	newID := c.lastIssuedGraphID

	c.ast.GraphNodes = append(c.ast.GraphNodes, newNode)

	return newID
}
