package parser

import "fmt"

// onGraphNode returns true if the registration of a new graph node
// was successful, otherwise returns false
func (pr *Parser) onGraphNode(newNode GraphNode) bool {
	// Prepare
	var errCodeRedecl ErrCode
	targetType := "<unknown>"

	switch newNode.(type) {
	case *StructField:
		errCodeRedecl = ErrStructFieldRedecl
		targetType = "struct"
	case *ResolverProperty:
		errCodeRedecl = ErrResolverPropRedecl
		targetType = "resolver property"
	case *Query:
		errCodeRedecl = ErrGraphRootNodeRedecl
		targetType = "graph root node"
	case *Mutation:
		errCodeRedecl = ErrGraphRootNodeRedecl
		targetType = "graph root node"
	}

	// Check for redeclaration
	newNodeName := newNode.GraphNodeName()
	if defined, isDef := pr.graphNodeByName[newNodeName]; isDef {
		definedSrc := defined.Source()
		pr.err(&pErr{
			at:   newNode.Source().Begin(),
			code: errCodeRedecl,
			message: fmt.Sprintf(
				"Redeclaration of %s %s (previously declared at %s)",
				targetType,
				newNodeName,
				definedSrc.Begin(),
			),
		})
		return false
	}

	// Assign unique identifier and register node
	pr.lastIssuedGraphID += GraphNodeID(1)
	newID := pr.lastIssuedGraphID
	switch newNode := newNode.(type) {
	case *StructField:
		newNode.GraphID = newID
	case *ResolverProperty:
		newNode.GraphID = newID
	case *Query:
		newNode.GraphID = newID
		pr.mod.QueryEndpoints = append(pr.mod.QueryEndpoints, newNode)
	case *Mutation:
		newNode.GraphID = newID
		pr.mod.Mutations = append(pr.mod.Mutations, newNode)
	}

	pr.mod.GraphNodes = append(pr.mod.GraphNodes, newNode)
	pr.graphNodeByID[newNode.GraphNodeID()] = newNode
	pr.graphNodeByName[newNodeName] = newNode

	return true
}
