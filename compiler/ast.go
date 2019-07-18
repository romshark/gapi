package compiler

// GraphNodeID represents a unique graph node identifier
type GraphNodeID uint64

// GraphNode represents a graph node
type GraphNode interface {
	GraphNodeID() GraphNodeID
	Parent() Type
	GraphNodeName() string
}

// AST represents the abstract GAPI syntax tree
type AST struct {
	SchemaName     string
	Types          []Type
	AliasTypes     []Type
	EnumTypes      []Type
	UnionTypes     []Type
	StructTypes    []Type
	QueryEndpoints []QueryEndpoint
	Mutations      []Mutation
	GraphNodes     []GraphNode
	TypeByID       map[TypeID]Type
}

// Clone returns a copy of the abstract syntax tree
func (ast *AST) Clone() *AST {
	if ast == nil {
		return nil
	}

	types := make([]Type, len(ast.Types))
	copy(types, ast.Types)

	aliasTypes := make([]Type, len(ast.AliasTypes))
	copy(aliasTypes, ast.AliasTypes)

	enumTypes := make([]Type, len(ast.EnumTypes))
	copy(enumTypes, ast.EnumTypes)

	unionTypes := make([]Type, len(ast.UnionTypes))
	copy(unionTypes, ast.UnionTypes)

	structTypes := make([]Type, len(ast.StructTypes))
	copy(structTypes, ast.StructTypes)

	queryEndpoints := make([]QueryEndpoint, len(ast.QueryEndpoints))
	copy(queryEndpoints, ast.QueryEndpoints)

	mutations := make([]Mutation, len(ast.Mutations))
	copy(mutations, ast.Mutations)

	graphNodes := make([]GraphNode, len(ast.GraphNodes))
	copy(graphNodes, ast.GraphNodes)

	typeByID := make(map[TypeID]Type, len(ast.TypeByID))
	for k, v := range ast.TypeByID {
		typeByID[k] = v
	}

	return &AST{
		SchemaName:     ast.SchemaName,
		Types:          types,
		AliasTypes:     aliasTypes,
		EnumTypes:      enumTypes,
		UnionTypes:     unionTypes,
		StructTypes:    structTypes,
		QueryEndpoints: queryEndpoints,
		Mutations:      mutations,
		GraphNodes:     graphNodes,
		TypeByID:       typeByID,
	}
}

// FindTypeByName returns a type given its category and name
func (ast *AST) FindTypeByName(category TypeCategory, name string) Type {
	findUserDefined := func() Type {
		for _, tp := range ast.Types {
			if tp.Name() == name {
				return tp
			}
		}
		return nil
	}

	switch category {
	case TypeCategoryPrimitive:
		// Search in primitives only
		return stdTypeByName(name)
	case TypeCategoryUserDefined:
		// Search in all categories including primitives
		return findUserDefined()
	case "":
		// Search in all user-defined types
		if t := stdTypeByName(name); t != nil {
			return t
		}
		return findUserDefined()
	default:
		// Search in specific category
		for _, tp := range ast.Types {
			if tp.Category() == category && tp.Name() == name {
				return tp
			}
		}
	}
	return nil
}

// FindGraphNodeByID returns a graph node given its unique identifier
func (ast *AST) FindGraphNodeByID(id GraphNodeID) GraphNode {
	for _, nd := range ast.GraphNodes {
		if nd.GraphNodeID() == id {
			return nd
		}
	}
	return nil
}
