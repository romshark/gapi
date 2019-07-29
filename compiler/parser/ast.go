package parser

// GraphNodeID represents a unique graph node identifier
type GraphNodeID int

// ParamID represents a unique parameter identifier
type ParamID int

// GraphNode represents a graph node
type GraphNode interface {
	Source() Fragment
	GraphNodeID() GraphNodeID
	Parent() Type
	NodeName() string
	GraphNodeName() string
}

// AST represents the abstract GAPI syntax tree
type AST struct {
	SchemaName     string
	Types          []Type
	EnumTypes      []Type
	UnionTypes     []Type
	StructTypes    []Type
	ResolverTypes  []Type
	AnonymousTypes []Type
	QueryEndpoints []*Query
	Mutations      []*Mutation
	GraphNodes     []GraphNode
}

// Clone returns a copy of the abstract syntax tree
func (ast *AST) Clone() *AST {
	if ast == nil {
		return nil
	}

	types := make([]Type, len(ast.Types))
	copy(types, ast.Types)

	enumTypes := make([]Type, len(ast.EnumTypes))
	copy(enumTypes, ast.EnumTypes)

	unionTypes := make([]Type, len(ast.UnionTypes))
	copy(unionTypes, ast.UnionTypes)

	structTypes := make([]Type, len(ast.StructTypes))
	copy(structTypes, ast.StructTypes)

	resolverTypes := make([]Type, len(ast.ResolverTypes))
	copy(resolverTypes, ast.ResolverTypes)

	anonymousTypes := make([]Type, len(ast.AnonymousTypes))
	copy(anonymousTypes, ast.AnonymousTypes)

	queryEndpoints := make([]*Query, len(ast.QueryEndpoints))
	copy(queryEndpoints, ast.QueryEndpoints)

	mutations := make([]*Mutation, len(ast.Mutations))
	copy(mutations, ast.Mutations)

	graphNodes := make([]GraphNode, len(ast.GraphNodes))
	copy(graphNodes, ast.GraphNodes)

	return &AST{
		SchemaName:     ast.SchemaName,
		Types:          types,
		EnumTypes:      enumTypes,
		UnionTypes:     unionTypes,
		StructTypes:    structTypes,
		ResolverTypes:  resolverTypes,
		AnonymousTypes: anonymousTypes,
		QueryEndpoints: queryEndpoints,
		Mutations:      mutations,
		GraphNodes:     graphNodes,
	}
}

// FindTypeByDesignation returns a type given its designation
func (ast *AST) FindTypeByDesignation(name string) Type {
	if t := stdTypeByName(name); t != nil {
		return t
	}
	for _, tp := range ast.Types {
		if tp.String() == name {
			return tp
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
	for _, qry := range ast.QueryEndpoints {
		for _, param := range qry.Parameters {
			if param.ID == id {
				return param
			}
		}
	}
	for _, mut := range ast.Mutations {
		for _, param := range mut.Parameters {
			if param.ID == id {
				return param
			}
		}
	}
	//TODO: search in subscriptions
	return nil
}
