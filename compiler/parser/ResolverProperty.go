package parser

// ResolverProperty represents a resolver property
type ResolverProperty struct {
	Src        Fragment
	Resolver   *TypeResolver
	Name       string
	GraphID    GraphNodeID
	Type       Type
	Parameters []*Parameter
}

// Source returns the source location of the declaration
func (pr *ResolverProperty) Source() Fragment { return pr.Src }

// GraphNodeID returns the unique graph node identifier of the resolver prop
func (pr *ResolverProperty) GraphNodeID() GraphNodeID { return pr.GraphID }

// Parent returns the parent resolver type of the resolver prop
func (pr *ResolverProperty) Parent() Type { return pr.Resolver }

// NodeName returns the property name
func (pr *ResolverProperty) NodeName() string { return pr.Name }

// GraphNodeName returns the graph node name
func (pr *ResolverProperty) GraphNodeName() string {
	return pr.Resolver.TypeName + "." + pr.Name
}
