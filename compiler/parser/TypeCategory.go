package parser

// TypeCategory represents the category of a type
type TypeCategory string

const (
	// TypeCategoryAnonymous represents an anonymous type
	TypeCategoryAnonymous TypeCategory = "anonymous"

	// TypeCategoryPrimitive represents a standard scalar type
	TypeCategoryPrimitive TypeCategory = "primitive"

	// TypeCategoryUserDefined represents all user-defined types
	TypeCategoryUserDefined TypeCategory = "user-defined"

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
