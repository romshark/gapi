package parser

// Type represents an abstract type implementation
type Type interface {
	Source() Fragment

	// String returns the designation of the type
	String() string

	// TerminalType returns the terminal type or nil if the current
	// type is already the terminal type
	TerminalType() Type

	// TypeID returns the type's unique identifier
	TypeID() TypeID

	// IsPure returns true if the type represents a data-only (pure) type
	IsPure() bool
}

type terminalType struct {
	Src  Fragment
	Name string
	ID   TypeID
}

func (i terminalType) Source() Fragment   { return i.Src }
func (i terminalType) String() string     { return i.Name }
func (i terminalType) TerminalType() Type { return nil }
func (i terminalType) TypeID() TypeID     { return i.ID }

/****************************************************************
	Alias
****************************************************************/

// TypeAlias represents an alias type implementation
type TypeAlias struct {
	terminalType
	AliasedType Type
}

// IsPure returns true if the aliased type is pure
func (t *TypeAlias) IsPure() bool { return t.AliasedType.IsPure() }

/****************************************************************
	Union
****************************************************************/

// TypeUnion represents an alias type implementation
type TypeUnion struct {
	terminalType
	Types []Type
}

// IsPure returns true if all option types are pure
func (t *TypeUnion) IsPure() bool {
	for _, optionType := range t.Types {
		if !optionType.IsPure() {
			return false
		}
	}
	return true
}

/****************************************************************
	Enumeration
****************************************************************/

// EnumValue represents an enumeration value
type EnumValue struct {
	Src  Fragment
	Name string
	Enum *TypeEnum
}

// TypeEnum represents a standard scalar type implementation
type TypeEnum struct {
	terminalType
	Values []*EnumValue
}

// IsPure always returns true for enumeration types
func (t *TypeEnum) IsPure() bool { return true }

/****************************************************************
	Optional
****************************************************************/

// TypeOptional represents an optional type implementation
type TypeOptional struct {
	Src       Fragment
	ID        TypeID
	StoreType Type
	Terminal  Type
}

// Source implements the Type interface
func (t *TypeOptional) Source() Fragment { return t.Src }

// Name implements the Type interface
func (t *TypeOptional) Name() string { return "?" + t.StoreType.String() }

// String implements the Type interface
func (t *TypeOptional) String() string { return stringifyType(t) }

// TerminalType implements the Type interface
func (t *TypeOptional) TerminalType() Type { return t.Terminal }

// TypeID returns the type's unique identifier
func (t *TypeOptional) TypeID() TypeID { return t.ID }

// IsPure returns true if the terminal type is pure
func (t *TypeOptional) IsPure() bool { return t.Terminal.IsPure() }

/****************************************************************
	List
****************************************************************/

// TypeList represents a list type implementation
type TypeList struct {
	Src       Fragment
	ID        TypeID
	StoreType Type
	Terminal  Type
}

// Source implements the Type interface
func (t *TypeList) Source() Fragment { return t.Src }

// Name implements the Type interface
func (t *TypeList) Name() string { return "[]" + t.StoreType.String() }

// String implements the Type interface
func (t *TypeList) String() string { return stringifyType(t) }

// TerminalType implements the Type interface
func (t *TypeList) TerminalType() Type { return t.Terminal }

// TypeID returns the type's unique identifier
func (t *TypeList) TypeID() TypeID { return t.ID }

// IsPure returns true if the terminal type is pure
func (t *TypeList) IsPure() bool { return t.Terminal.IsPure() }

/****************************************************************
	Standard Bool
****************************************************************/

// TypeStdNone represents a standard scalar type implementation
type TypeStdNone struct{}

// Source implements the Type interface
func (t TypeStdNone) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdNone) Name() string { return "None" }

// Category implements the Type interface
func (t TypeStdNone) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdNone) String() string { return "None" }

// TerminalType implements the Type interface
func (t TypeStdNone) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdNone) TypeID() TypeID { return TypeIDPrimitiveNone }

// IsPure always returns false for None primitives
func (t TypeStdNone) IsPure() bool { return false }

/****************************************************************
	Standard Bool
****************************************************************/

// TypeStdBool represents a standard scalar type implementation
type TypeStdBool struct{}

// Source implements the Type interface
func (t TypeStdBool) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdBool) Name() string { return "Bool" }

// Category implements the Type interface
func (t TypeStdBool) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdBool) String() string { return "Bool" }

// TerminalType implements the Type interface
func (t TypeStdBool) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdBool) TypeID() TypeID { return TypeIDPrimitiveBool }

// IsPure always returns true for Bool primitives
func (t TypeStdBool) IsPure() bool { return true }

/****************************************************************
	Standard Byte
****************************************************************/

// TypeStdByte represents a standard scalar type implementation
type TypeStdByte struct{}

// Source implements the Type interface
func (t TypeStdByte) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdByte) Name() string { return "Byte" }

// Category implements the Type interface
func (t TypeStdByte) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdByte) String() string { return "Byte" }

// TerminalType implements the Type interface
func (t TypeStdByte) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdByte) TypeID() TypeID { return TypeIDPrimitiveByte }

// IsPure always returns true for Byte primitives
func (t TypeStdByte) IsPure() bool { return true }

/****************************************************************
	Standard Int32
****************************************************************/

// TypeStdInt32 represents a standard scalar type implementation
type TypeStdInt32 struct{}

// Source implements the Type interface
func (t TypeStdInt32) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdInt32) Name() string { return "Int32" }

// Category implements the Type interface
func (t TypeStdInt32) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdInt32) String() string { return "Int32" }

// TerminalType implements the Type interface
func (t TypeStdInt32) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdInt32) TypeID() TypeID { return TypeIDPrimitiveInt32 }

// IsPure always returns true for Int32 primitives
func (t TypeStdInt32) IsPure() bool { return true }

/****************************************************************
	Standard Uint32
****************************************************************/

// TypeStdUint32 represents a standard scalar type implementation
type TypeStdUint32 struct{}

// Source implements the Type interface
func (t TypeStdUint32) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdUint32) Name() string { return "Uint32" }

// Category implements the Type interface
func (t TypeStdUint32) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdUint32) String() string { return "Uint32" }

// TerminalType implements the Type interface
func (t TypeStdUint32) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdUint32) TypeID() TypeID { return TypeIDPrimitiveUint32 }

// IsPure always returns true for Uint32 primitives
func (t TypeStdUint32) IsPure() bool { return true }

/****************************************************************
	Standard Int64
****************************************************************/

// TypeStdInt64 represents a standard scalar type implementation
type TypeStdInt64 struct{}

// Source implements the Type interface
func (t TypeStdInt64) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdInt64) Name() string { return "Int64" }

// Category implements the Type interface
func (t TypeStdInt64) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdInt64) String() string { return "Int64" }

// TerminalType implements the Type interface
func (t TypeStdInt64) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdInt64) TypeID() TypeID { return TypeIDPrimitiveInt64 }

// IsPure always returns true for Int64 primitives
func (t TypeStdInt64) IsPure() bool { return true }

/****************************************************************
	Standard Uint64
****************************************************************/

// TypeStdUint64 represents a standard scalar type implementation
type TypeStdUint64 struct{}

// Source implements the Type interface
func (t TypeStdUint64) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdUint64) Name() string { return "Uint64" }

// Category implements the Type interface
func (t TypeStdUint64) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdUint64) String() string { return "Uint64" }

// TerminalType implements the Type interface
func (t TypeStdUint64) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdUint64) TypeID() TypeID { return TypeIDPrimitiveUint64 }

// IsPure always returns true for Uint64 primitives
func (t TypeStdUint64) IsPure() bool { return true }

/****************************************************************
	Standard Float64
****************************************************************/

// TypeStdFloat64 represents a standard scalar type implementation
type TypeStdFloat64 struct{}

// Source implements the Type interface
func (t TypeStdFloat64) Source() Fragment { return nil }

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

// TypeID returns the type's unique identifier
func (t TypeStdFloat64) TypeID() TypeID { return TypeIDPrimitiveFloat64 }

// IsPure always returns true for Float64 primitives
func (t TypeStdFloat64) IsPure() bool { return true }

/****************************************************************
	Standard String
****************************************************************/

// TypeStdString represents a standard scalar type implementation
type TypeStdString struct{}

// Source implements the Type interface
func (t TypeStdString) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdString) Name() string { return "String" }

// Category implements the Type interface
func (t TypeStdString) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdString) String() string { return "String" }

// TerminalType implements the Type interface
func (t TypeStdString) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdString) TypeID() TypeID { return TypeIDPrimitiveString }

// IsPure always returns true for String primitives
func (t TypeStdString) IsPure() bool { return true }

/****************************************************************
	Standard Time
****************************************************************/

// TypeStdTime represents a standard scalar type implementation
type TypeStdTime struct{}

// Source implements the Type interface
func (t TypeStdTime) Source() Fragment { return nil }

// Name implements the Type interface
func (t TypeStdTime) Name() string { return "Time" }

// Category implements the Type interface
func (t TypeStdTime) Category() TypeCategory { return TypeCategoryPrimitive }

// String implements the Type interface
func (t TypeStdTime) String() string { return "Time" }

// TerminalType implements the Type interface
func (t TypeStdTime) TerminalType() Type { return nil }

// TypeID returns the type's unique identifier
func (t TypeStdTime) TypeID() TypeID { return TypeIDPrimitiveTime }

// IsPure always returns true for Time primitives
func (t TypeStdTime) IsPure() bool { return true }

/****************************************************************
	Struct
****************************************************************/

// TypeStruct represents a standard scalar type implementation
type TypeStruct struct {
	terminalType
	Fields []*StructField
}

// Category implements the Type interface
func (t *TypeStruct) Category() TypeCategory { return TypeCategoryStruct }

// IsPure always returns true for struct types
func (t *TypeStruct) IsPure() bool { return true }

/****************************************************************
	Resolver
****************************************************************/

// Parameter represents either a resolver property-, a query-, a mutation-
// or a subscription parameter
type Parameter struct {
	Src    Fragment
	Target GraphNode
	Name   string
	ID     ParamID
	Type   Type
}

// TypeResolver represents a resolver type
type TypeResolver struct {
	terminalType
	Properties []*ResolverProperty
}

// Category implements the Type interface
func (t *TypeResolver) Category() TypeCategory { return TypeCategoryResolver }

// IsPure always returns false for resolver types
func (t *TypeResolver) IsPure() bool { return false }

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
