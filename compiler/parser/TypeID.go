package parser

// TypeID represents a unique type identifier
type TypeID int

const (
	// TypeIDOptional represents the unique identifier of the generic optional
	// type
	TypeIDOptional TypeID = 1

	// TypeIDList represents the unique identifier of the generic list type
	TypeIDList TypeID = 2

	// TypeIDPrimitiveNone represents the unique identifier of the primitive
	// none type
	TypeIDPrimitiveNone TypeID = 3

	// TypeIDPrimitiveBool represents the unique identifier of the primitive
	// boolean type
	TypeIDPrimitiveBool TypeID = 4

	// TypeIDPrimitiveByte represents the unique identifier of the primitive
	// byte type
	TypeIDPrimitiveByte TypeID = 5

	// TypeIDPrimitiveInt32 represents the unique identifier of the primitive
	// 32-bit signed integer type
	TypeIDPrimitiveInt32 TypeID = 6

	// TypeIDPrimitiveUint32 represents the unique identifier of the primitive
	// 32-bit unsigned integer type
	TypeIDPrimitiveUint32 TypeID = 7

	// TypeIDPrimitiveInt64 represents the unique identifier of the primitive
	// 64-bit signed integer type
	TypeIDPrimitiveInt64 TypeID = 8

	// TypeIDPrimitiveUint64 represents the unique identifier of the primitive
	// 64-bit unsigned integer type
	TypeIDPrimitiveUint64 TypeID = 9

	// TypeIDPrimitiveFloat64 represents the unique identifier of the primitive
	// 64-bit floating point number type
	TypeIDPrimitiveFloat64 TypeID = 10

	// TypeIDPrimitiveString represents the unique identifier of the primitive
	// UTF8 encoded string type
	TypeIDPrimitiveString TypeID = 11

	// TypeIDPrimitiveTime represents the unique identifier of the primitive
	// RFC3339 encoded time type
	TypeIDPrimitiveTime TypeID = 12

	// TypeIDUserTypeOffset represents the ID offset for user-defined types
	// (the ID of the first user-defined type starts with 100)
	TypeIDUserTypeOffset TypeID = 99
)
