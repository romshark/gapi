package parser

import parser "github.com/romshark/llparser"

// GraphNodeID represents a unique graph node identifier
type GraphNodeID int

// ParamID represents a unique parameter identifier
type ParamID int

// GraphNode represents a graph node
type GraphNode interface {
	Source() parser.Fragment
	GraphNodeID() GraphNodeID
	Parent() Type
	NodeName() string
	GraphNodeName() string
}

// SchemaModel represents a schema model
type SchemaModel struct {
	SchemaName     string
	Types          []Type
	AliasTypes     []Type
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
func (mod *SchemaModel) Clone() *SchemaModel {
	if mod == nil {
		return nil
	}

	types := make([]Type, len(mod.Types))
	copy(types, mod.Types)

	aliasTypes := make([]Type, len(mod.AliasTypes))
	copy(aliasTypes, mod.AliasTypes)

	enumTypes := make([]Type, len(mod.EnumTypes))
	copy(enumTypes, mod.EnumTypes)

	unionTypes := make([]Type, len(mod.UnionTypes))
	copy(unionTypes, mod.UnionTypes)

	structTypes := make([]Type, len(mod.StructTypes))
	copy(structTypes, mod.StructTypes)

	resolverTypes := make([]Type, len(mod.ResolverTypes))
	copy(resolverTypes, mod.ResolverTypes)

	anonymousTypes := make([]Type, len(mod.AnonymousTypes))
	copy(anonymousTypes, mod.AnonymousTypes)

	queryEndpoints := make([]*Query, len(mod.QueryEndpoints))
	copy(queryEndpoints, mod.QueryEndpoints)

	mutations := make([]*Mutation, len(mod.Mutations))
	copy(mutations, mod.Mutations)

	graphNodes := make([]GraphNode, len(mod.GraphNodes))
	copy(graphNodes, mod.GraphNodes)

	return &SchemaModel{
		SchemaName:     mod.SchemaName,
		Types:          types,
		AliasTypes:     aliasTypes,
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
func (mod *SchemaModel) FindTypeByDesignation(name string) Type {
	if t := stdTypeByName(name); t != nil {
		return t
	}
	for _, tp := range mod.Types {
		if tp.String() == name {
			return tp
		}
	}
	return nil
}

// FindGraphNodeByID returns a graph node given its unique identifier
func (mod *SchemaModel) FindGraphNodeByID(id GraphNodeID) GraphNode {
	for _, nd := range mod.GraphNodes {
		if nd.GraphNodeID() == id {
			return nd
		}
	}
	return nil
}

// FindTypeByID returns a type given its unique identifier or nil
// if no type is identified by the given identifier
func (mod *SchemaModel) FindTypeByID(id TypeID) Type {
	for _, tp := range mod.Types {
		if tp.TypeID() == id {
			return tp
		}
	}
	return nil
}

// FindParameterByID returns a parameter given its unique identifier or nil
// if no parameter is identified by the given identifier
func (mod *SchemaModel) FindParameterByID(id ParamID) *Parameter {
	for _, rsv := range mod.ResolverTypes {
		for _, prop := range rsv.(*TypeResolver).Properties {
			for _, param := range prop.Parameters {
				if param.ID == id {
					return param
				}
			}
		}
	}
	for _, qry := range mod.QueryEndpoints {
		for _, param := range qry.Parameters {
			if param.ID == id {
				return param
			}
		}
	}
	for _, mut := range mod.Mutations {
		for _, param := range mut.Parameters {
			if param.ID == id {
				return param
			}
		}
	}
	//TODO: search in subscriptions
	return nil
}
