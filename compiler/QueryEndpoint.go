package compiler

// QueryEndpoint represents a query endpoint
type QueryEndpoint struct {
	Src
	Name       string
	GraphID    GraphNodeID
	Parameters []*Parameter
	Type       Type
}

// Source returns the source location of the declaration
func (qe *QueryEndpoint) Source() Src { return qe.Src }

// GraphNodeID returns the query endpoint's unique graph node identifier
func (qe *QueryEndpoint) GraphNodeID() GraphNodeID { return qe.GraphID }

// NodeName returns the query endpoint's name
func (qe *QueryEndpoint) NodeName() string { return qe.Name }

// GraphNodeName returns the query endpoint's graph node name
func (qe *QueryEndpoint) GraphNodeName() string { return qe.Name }

// Parent returns nil indicating root
func (qe *QueryEndpoint) Parent() Type { return nil }
