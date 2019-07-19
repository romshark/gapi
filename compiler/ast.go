package compiler

// GraphNodeID represents a unique graph node identifier
type GraphNodeID int

// ParamID represents a unique parameter identifier
type ParamID int

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
	ResolverTypes  []Type
	QueryEndpoints []QueryEndpoint
	Mutations      []Mutation
	GraphNodes     []GraphNode
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

	resolverTypes := make([]Type, len(ast.ResolverTypes))
	copy(resolverTypes, ast.ResolverTypes)

	queryEndpoints := make([]QueryEndpoint, len(ast.QueryEndpoints))
	copy(queryEndpoints, ast.QueryEndpoints)

	mutations := make([]Mutation, len(ast.Mutations))
	copy(mutations, ast.Mutations)

	graphNodes := make([]GraphNode, len(ast.GraphNodes))
	copy(graphNodes, ast.GraphNodes)

	return &AST{
		SchemaName:     ast.SchemaName,
		Types:          types,
		AliasTypes:     aliasTypes,
		EnumTypes:      enumTypes,
		UnionTypes:     unionTypes,
		StructTypes:    structTypes,
		ResolverTypes:  resolverTypes,
		QueryEndpoints: queryEndpoints,
		Mutations:      mutations,
		GraphNodes:     graphNodes,
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

// FindTypeByID returns a type given its unique identifier or nil
// if no type is identified by the given identifier
func (ast *AST) FindTypeByID(id TypeID) Type {
	for _, tp := range ast.Types {
		if tp.TypeID() == id {
			return tp
		}
	}
	return nil
}

// FindParameterByID returns a parameter given its unique identifier or nil
// if no parameter is identified by the given identifier
func (ast *AST) FindParameterByID(id ParamID) *Parameter {
	for _, rsv := range ast.ResolverTypes {
		for _, prop := range rsv.(*TypeResolver).Properties {
			for _, param := range prop.Parameters {
				if param.ID == id {
					return param
				}
			}
		}
	}
	//TODO: search in queries
	//TODO: search in mutations
	//TODO: search in subscriptions
	return nil
}
