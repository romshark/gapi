package parser

import parser "github.com/romshark/llparser"

// Query represents a query endpoint
type Query struct {
	Src        parser.Fragment
	Name       string
	GraphID    GraphNodeID
	Parameters []*Parameter
	Type       Type
}

// Source returns the source location of the declaration
func (qe *Query) Source() parser.Fragment { return qe.Src }

// GraphNodeID returns the query endpoint's unique graph node identifier
func (qe *Query) GraphNodeID() GraphNodeID { return qe.GraphID }

// NodeName returns the query endpoint's name
func (qe *Query) NodeName() string { return qe.Name }

// GraphNodeName returns the query endpoint's graph node name
func (qe *Query) GraphNodeName() string { return qe.Name }

// Parent returns nil indicating root
func (qe *Query) Parent() Type { return nil }
