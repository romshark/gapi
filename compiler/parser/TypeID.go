package parser

// TypeID represents a unique type identifier
type TypeID int

const (
	// TypeIDPrimitiveNone represents the unique identifier of the primitive
	// none type
	TypeIDPrimitiveNone TypeID = 1

	// TypeIDPrimitiveBool represents the unique identifier of the primitive
	// boolean type
	TypeIDPrimitiveBool TypeID = 2

	// TypeIDPrimitiveByte represents the unique identifier of the primitive
	// byte type
	TypeIDPrimitiveByte TypeID = 3

	// TypeIDPrimitiveInt32 represents the unique identifier of the primitive
	// 32-bit signed integer type
	TypeIDPrimitiveInt32 TypeID = 4

	// TypeIDPrimitiveUint32 represents the unique identifier of the primitive
	// 32-bit unsigned integer type
	TypeIDPrimitiveUint32 TypeID = 5

	// TypeIDPrimitiveInt64 represents the unique identifier of the primitive
	// 64-bit signed integer type
	TypeIDPrimitiveInt64 TypeID = 6

	// TypeIDPrimitiveUint64 represents the unique identifier of the primitive
	// 64-bit unsigned integer type
	TypeIDPrimitiveUint64 TypeID = 7

	// TypeIDPrimitiveFloat64 represents the unique identifier of the primitive
	// 64-bit floating point number type
	TypeIDPrimitiveFloat64 TypeID = 8

	// TypeIDPrimitiveString represents the unique identifier of the primitive
	// UTF8 encoded string type
	TypeIDPrimitiveString TypeID = 9

	// TypeIDPrimitiveTime represents the unique identifier of the primitive
	// RFC3339 encoded time type
	TypeIDPrimitiveTime TypeID = 10

	// TypeIDUserTypeOffset represents the ID offset for user-defined types
	// (the ID of the first user-defined type starts with 100)
	TypeIDUserTypeOffset TypeID = 99
)
