package compiler

// StructField represents a struct field
type StructField struct {
	Src
	Struct  *TypeStruct
	GraphID GraphNodeID
	Name    string
	Type    Type
}

// Source returns the source location of the declaration
func (sf *StructField) Source() Src { return sf.Src }

// GraphNodeID returns the unique graph node identifier of the struct field
func (sf *StructField) GraphNodeID() GraphNodeID { return sf.GraphID }

// Parent returns the parent struct type of the struct field
func (sf *StructField) Parent() Type { return sf.Struct }

// NodeName returns the property name
func (sf *StructField) NodeName() string { return sf.Name }

// GraphNodeName returns the graph node name
func (sf *StructField) GraphNodeName() string {
	return sf.Struct.name + "." + sf.Name
}
