package compiler

// Type represents an abstract type implementation
type Type interface {
	Name() string
	Src() Src
	Category() TypeCategory
}

type typeBaseInfo struct {
	src  Src
	name string
}

func (i typeBaseInfo) Src() Src {
	return i.src
}

func (i typeBaseInfo) Name() string {
	return i.name
}

/****************************************************************
	Alias
****************************************************************/

// TypeAlias represents an alias type implementation
type TypeAlias struct {
	typeBaseInfo
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
	typeBaseInfo
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
	typeBaseInfo
	Values map[string]EnumValue
}

// Category implements the Type interface
func (t *TypeEnum) Category() TypeCategory {
	return TypeCategoryEnum
}

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

/****************************************************************
	Struct
****************************************************************/

// StructField represents a struct field
type StructField struct {
	Src
	Name string
	Type Type
}

// TypeStruct represents a standard scalar type implementation
type TypeStruct struct {
	typeBaseInfo
	Fields []StructField
}

// Category implements the Type interface
func (t *TypeStruct) Category() TypeCategory {
	return TypeCategoryStruct
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
	typeBaseInfo
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
	typeBaseInfo
	Pure       bool
	Properties []ResolverProperty
}

// Category implements the Type interface
func (t *TypeTrait) Category() TypeCategory {
	return TypeCategoryTrait
}
