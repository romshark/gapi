package compiler

// Type represents an abstract type implementation
type Type interface {
	Name() string
	Src() Src
	Category() TypeCategory
	String() string

	// TerminalType returns the terminal type or nil if the current
	// type is already the terminal type
	TerminalType() Type
}

type terminalType struct {
	src  Src
	name string
}

func (i terminalType) Src() Src           { return i.src }
func (i terminalType) Name() string       { return i.name }
func (i terminalType) String() string     { return i.name }
func (i terminalType) TerminalType() Type { return nil }

/****************************************************************
	Alias
****************************************************************/

// TypeAlias represents an alias type implementation
type TypeAlias struct {
	terminalType
	AliasedType Type
}

// Category implements the Type interface
func (t *TypeAlias) Category() TypeCategory {
	return TypeCategoryAlias
}

/****************************************************************
	Union
****************************************************************/

// TypeUnion represents an alias type implementation
type TypeUnion struct {
	terminalType
	Types map[string]Type
}

// Category implements the Type interface
func (t *TypeUnion) Category() TypeCategory {
	return TypeCategoryUnion
}

/****************************************************************
	Enumeration
****************************************************************/

// EnumValue represents an enumeration value
type EnumValue struct {
	Src
	Name string
}

// TypeEnum represents a standard scalar type implementation
type TypeEnum struct {
	terminalType
	Values map[string]EnumValue
}

// Category implements the Type interface
func (t *TypeEnum) Category() TypeCategory {
	return TypeCategoryEnum
}

/****************************************************************
	Optional
****************************************************************/

// TypeOptional represents an optional type implementation
type TypeOptional struct {
	StoreType Type
	Terminal  Type
}

// Src implements the Type interface
func (t *TypeOptional) Src() Src { return Src{} }

// Name implements the Type interface
func (t *TypeOptional) Name() string { return "?" + t.StoreType.Name() }

// Category implements the Type interface
func (t *TypeOptional) Category() TypeCategory { return TypeCategoryOptional }

// String implements the Type interface
func (t *TypeOptional) String() string { return stringifyType(t) }

// TerminalType implements the Type interface
func (t *TypeOptional) TerminalType() Type { return t.Terminal }

/****************************************************************
	List
****************************************************************/

// TypeList represents a list type implementation
type TypeList struct {
	StoreType Type
	Terminal  Type
}

// Src implements the Type interface
func (t *TypeList) Src() Src { return Src{} }

// Name implements the Type interface
func (t *TypeList) Name() string { return "[]" + t.StoreType.Name() }

// Category implements the Type interface
func (t *TypeList) Category() TypeCategory { return TypeCategoryList }

// String implements the Type interface
func (t *TypeList) String() string { return stringifyType(t) }

// TerminalType implements the Type interface
func (t *TypeList) TerminalType() Type { return t.Terminal }

/****************************************************************
	Standard Bool
****************************************************************/

// TypeStdNone represents a standard scalar type implementation
type TypeStdNone struct{}

// Src implements the Type interface
func (t TypeStdNone) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdNone) Name() string { return "None" }

// Category implements the Type interface
func (t TypeStdNone) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdNone) String() string { return "None" }

// TerminalType implements the Type interface
func (t TypeStdNone) TerminalType() Type { return nil }

/****************************************************************
	Standard Bool
****************************************************************/

// TypeStdBool represents a standard scalar type implementation
type TypeStdBool struct{}

// Src implements the Type interface
func (t TypeStdBool) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdBool) Name() string { return "Bool" }

// Category implements the Type interface
func (t TypeStdBool) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdBool) String() string { return "Bool" }

// TerminalType implements the Type interface
func (t TypeStdBool) TerminalType() Type { return nil }

/****************************************************************
	Standard Byte
****************************************************************/

// TypeStdByte represents a standard scalar type implementation
type TypeStdByte struct{}

// Src implements the Type interface
func (t TypeStdByte) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdByte) Name() string { return "Byte" }

// Category implements the Type interface
func (t TypeStdByte) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdByte) String() string { return "Byte" }

// TerminalType implements the Type interface
func (t TypeStdByte) TerminalType() Type { return nil }

/****************************************************************
	Standard Int32
****************************************************************/

// TypeStdInt32 represents a standard scalar type implementation
type TypeStdInt32 struct{}

// Src implements the Type interface
func (t TypeStdInt32) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdInt32) Name() string { return "Int32" }

// Category implements the Type interface
func (t TypeStdInt32) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdInt32) String() string { return "Int32" }

// TerminalType implements the Type interface
func (t TypeStdInt32) TerminalType() Type { return nil }

/****************************************************************
	Standard Uint32
****************************************************************/

// TypeStdUint32 represents a standard scalar type implementation
type TypeStdUint32 struct{}

// Src implements the Type interface
func (t TypeStdUint32) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdUint32) Name() string { return "Uint32" }

// Category implements the Type interface
func (t TypeStdUint32) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdUint32) String() string { return "Uint32" }

// TerminalType implements the Type interface
func (t TypeStdUint32) TerminalType() Type { return nil }

/****************************************************************
	Standard Int64
****************************************************************/

// TypeStdInt64 represents a standard scalar type implementation
type TypeStdInt64 struct{}

// Src implements the Type interface
func (t TypeStdInt64) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdInt64) Name() string { return "Int64" }

// Category implements the Type interface
func (t TypeStdInt64) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdInt64) String() string { return "Int64" }

// TerminalType implements the Type interface
func (t TypeStdInt64) TerminalType() Type { return nil }

/****************************************************************
	Standard Uint64
****************************************************************/

// TypeStdUint64 represents a standard scalar type implementation
type TypeStdUint64 struct{}

// Src implements the Type interface
func (t TypeStdUint64) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdUint64) Name() string { return "Uint64" }

// Category implements the Type interface
func (t TypeStdUint64) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdUint64) String() string { return "Uint64" }

// TerminalType implements the Type interface
func (t TypeStdUint64) TerminalType() Type { return nil }

/****************************************************************
	Standard Float64
****************************************************************/

// TypeStdFloat64 represents a standard scalar type implementation
type TypeStdFloat64 struct{}

// Src implements the Type interface
func (t TypeStdFloat64) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdFloat64) Name() string { return "Float64" }

// Category implements the Type interface
func (t TypeStdFloat64) Category() TypeCategory {
	return TypeCategoryPrimitive
}

// String implements the Type interface
func (t TypeStdFloat64) String() string { return "Float64" }

// TerminalType implements the Type interface
func (t TypeStdFloat64) TerminalType() Type { return nil }

/****************************************************************
	Standard String
****************************************************************/

// TypeStdString represents a standard scalar type implementation
type TypeStdString struct{}

// Src implements the Type interface
func (t TypeStdString) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdString) Name() string { return "String" }

// Category implements the Type interface
func (t TypeStdString) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdString) String() string { return "String" }

// TerminalType implements the Type interface
func (t TypeStdString) TerminalType() Type { return nil }

/****************************************************************
	Standard Time
****************************************************************/

// TypeStdTime represents a standard scalar type implementation
type TypeStdTime struct{}

// Src implements the Type interface
func (t TypeStdTime) Src() Src { return Src{} }

// Name implements the Type interface
func (t TypeStdTime) Name() string { return "Time" }

// Category implements the Type interface
func (t TypeStdTime) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdTime) String() string { return "Time" }

// TerminalType implements the Type interface
func (t TypeStdTime) TerminalType() Type { return nil }

/****************************************************************
	Struct
****************************************************************/

// StructField represents a struct field
type StructField struct {
	Src
	Struct  *TypeStruct
	GraphID GraphNodeID
	Name    string
	Type    Type
}

// GraphNodeID returns the unique graph node identifier of the struct field
func (sf *StructField) GraphNodeID() GraphNodeID {
	return sf.GraphID
}

// Parent returns the parent struct type of the struct field
func (sf *StructField) Parent() Type {
	return sf.Struct
}

// GraphNodeName returns the graph node name
func (sf *StructField) GraphNodeName() string {
	return sf.Struct.name + "." + sf.Name
}

// TypeStruct represents a standard scalar type implementation
type TypeStruct struct {
	terminalType
	Fields []*StructField
}

// Category implements the Type interface
func (t *TypeStruct) Category() TypeCategory {
	return TypeCategoryStruct
}

// FieldByName returns a field given its name
func (t *TypeStruct) FieldByName(name string) *StructField {
	for _, field := range t.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

/****************************************************************
	Resolver
****************************************************************/

// Variable represents a variable
type Variable struct {
	Src
	Name string
	Type Type
}

// ResolverProperty represents a resolver property
type ResolverProperty struct {
	Name      string
	Type      Type
	Variables []Variable
}

// TypeResolver represents a standard scalar type implementation
type TypeResolver struct {
	terminalType
	Properties []ResolverProperty
}

// Category implements the Type interface
func (t *TypeResolver) Category() TypeCategory {
	return TypeCategoryResolver
}

/****************************************************************
	Trait
****************************************************************/

// TypeTrait represents a standard scalar type implementation
type TypeTrait struct {
	terminalType
	Pure       bool
	Properties []ResolverProperty
}

// Category implements the Type interface
func (t *TypeTrait) Category() TypeCategory {
	return TypeCategoryTrait
}
