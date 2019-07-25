package parser

// Mutation represents a mutation endpoint
type Mutation struct {
	Src        Fragment
	Name       string
	GraphID    GraphNodeID
	Parameters []*Parameter
	Type       Type
}

// Source returns the source location of the declaration
func (mt *Mutation) Source() Fragment { return mt.Src }

// GraphNodeID returns the mutation endpoint's unique graph node identifier
func (mt *Mutation) GraphNodeID() GraphNodeID { return mt.GraphID }

// NodeName returns the property name
func (mt *Mutation) NodeName() string { return mt.Name }

// GraphNodeName returns the mutation endpoint's graph node name
func (mt *Mutation) GraphNodeName() string { return mt.Name }

// Parent returns nil indicating root
func (mt *Mutation) Parent() Type { return nil }
