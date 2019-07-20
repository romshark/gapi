package compiler

import "fmt"

// defineGraphNode returns a new unique graph node identifier
func (c *Compiler) defineGraphNode(newNode GraphNode) {
	// Prepare
	src := newNode.Source()
	var errCodeIdent ErrCode
	var errCodeRedecl ErrCode
	nodeType := "<unknown>"
	targetType := "<unknown>"

	switch newNode.(type) {
	case *StructField:
		errCodeIdent = ErrStructFieldIllegalIdent
		errCodeRedecl = ErrStructFieldRedecl
		nodeType = "struct field"
		targetType = "struct"
	case *ResolverProperty:
		errCodeIdent = ErrResolverPropIllegalIdent
		errCodeRedecl = ErrResolverPropRedecl
		nodeType = "resolver property"
		targetType = "resolver property"
	case *QueryEndpoint:
		errCodeIdent = ErrQryEndpointIllegalIdent
		errCodeRedecl = ErrGraphRootNodeRedecl
		nodeType = "query endpoint"
		targetType = "graph root node"
	case *Mutation:
		errCodeIdent = ErrMutEndpointIllegalIdent
		errCodeRedecl = ErrGraphRootNodeRedecl
		nodeType = "mutation endpoint"
		targetType = "graph root node"
	}

	// Verify node identifier
	if err := verifyLowerCamelCase(newNode.NodeName()); err != nil {
		c.err(cErr{
			errCodeIdent,
			fmt.Sprintf("illegal %s identifier %d:%d: %s",
				nodeType,
				src.Begin,
				src.End,
				err,
			),
		})
		return
	}

	// Check for redeclaration
	newNodeName := newNode.GraphNodeName()
	if defined, isDef := c.graphNodeByName[newNodeName]; isDef {
		definedSrc := defined.Source()
		c.err(cErr{
			errCodeRedecl,
			fmt.Sprintf(
				"Redeclaration of %s %s at %d:%d "+
					"(previously declared at %d:%d)",
				targetType,
				newNodeName,
				src.Begin,
				src.End,
				definedSrc.Begin,
				definedSrc.End,
			),
		})
		return
	}

	// Assign unique identifier
	c.lastIssuedGraphID += GraphNodeID(1)
	newID := c.lastIssuedGraphID
	switch newNode := newNode.(type) {
	case *StructField:
		newNode.GraphID = newID
	case *QueryEndpoint:
		newNode.GraphID = newID
	case *Mutation:
		newNode.GraphID = newID
	case *ResolverProperty:
		newNode.GraphID = newID
	}

	c.ast.GraphNodes = append(c.ast.GraphNodes, newNode)
	c.graphNodeByID[newNode.GraphNodeID()] = newNode
}
