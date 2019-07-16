package compiler

// TypeCategory represents the category of a type
type TypeCategory string

const (
	// TypeCategoryPrimitive represents a standard scalar type
	TypeCategoryPrimitive TypeCategory = "primitive"

	// TypeCategoryUserDefined represents all user-defined types
	TypeCategoryUserDefined TypeCategory = "user-defined"

	// TypeCategoryAlias represents an alias type
	TypeCategoryAlias TypeCategory = "alias"

	// TypeCategoryEnum represents an enum type
	TypeCategoryEnum TypeCategory = "enum"

	// TypeCategoryStruct represents a struct type
	TypeCategoryStruct TypeCategory = "struct"

	// TypeCategoryUnion represents a union type
	TypeCategoryUnion TypeCategory = "union"

	// TypeCategoryTrait represents a trait type
	TypeCategoryTrait TypeCategory = "trait"

	// TypeCategoryResolver represents a resolver type
	TypeCategoryResolver TypeCategory = "resolver"
)
