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
	Types []Type
}

// Category implements the Type interface
func (t *TypeUnion) Category() TypeCategory {
	return TypeCategoryUnion
}

/****************************************************************
	Enumeration
****************************************************************/

// TypeEnum represents a standard scalar type implementation
type TypeEnum struct {
	typeBaseInfo
	Values []string
}

// Category implements the Type interface
func (t *TypeEnum) Category() TypeCategory {
	return TypeCategoryEnum
}

/****************************************************************
	Standard Scalar
****************************************************************/

// TypeScalarStd represents a standard scalar type implementation
type TypeScalarStd struct {
	typeBaseInfo
}

// Category implements the Type interface
func (t *TypeScalarStd) Category() TypeCategory {
	return TypeCategoryScalarStd
}

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
