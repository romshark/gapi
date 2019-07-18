package compiler

// ErrCode represents a compiler error code
type ErrCode int

const (
	_ ErrCode = iota

	// ErrSyntax represents a syntax error
	ErrSyntax

	// ErrSchemaIllegalIdent indicates an illegal schema identifier
	ErrSchemaIllegalIdent

	// ErrTypeUndef indicates an undefined referenced type
	ErrTypeUndef

	// ErrTypeRedecl indicates a redeclared type
	ErrTypeRedecl

	// ErrTypeIllegalIdent indicates an illegal type identifier
	ErrTypeIllegalIdent

	// ErrTypeOptChain indicates an illegal chain of optionals
	ErrTypeOptChain

	// ErrAliasRecurs indicates a recursive alias type
	ErrAliasRecurs

	// ErrEnumNoVal indicates an empty enumeration missing values
	ErrEnumNoVal

	// ErrEnumValIllegalIdent indicates an illegal enum value identifier
	ErrEnumValIllegalIdent

	// ErrEnumValRedecl indicates a redeclared enum value
	ErrEnumValRedecl

	// ErrUnionRedund indicates redundant union option types
	ErrUnionRedund

	// ErrUnionSelfref indicates a self-referencing union type
	ErrUnionSelfref

	// ErrUnionMissingOpts indicates a union type with too little option types
	ErrUnionMissingOpts

	// ErrUnionIncludesNone indicates a union type including the None type
	ErrUnionIncludesNone

	// ErrStructFieldIllegalIdent indicates an illegal struct field identifier
	ErrStructFieldIllegalIdent

	// ErrStructFieldRedecl indicates a redeclared struct field
	ErrStructFieldRedecl

	// ErrStructNoFields indicates an empty struct type missing fields
	ErrStructNoFields

	// ErrStructRecurs indicates a recursive struct type
	ErrStructRecurs
)

// String stringifies the error code
func (c ErrCode) String() string {
	switch c {
	case ErrSyntax:
		return "Syntax"
	case ErrSchemaIllegalIdent:
		return "SchemaIllegalIdent"
	case ErrTypeUndef:
		return "TypeUndef"
	case ErrTypeRedecl:
		return "TypeRedecl"
	case ErrTypeIllegalIdent:
		return "TypeIllegalIdent"
	case ErrTypeOptChain:
		return "TypeOptChain"
	case ErrAliasRecurs:
		return "AliasRecurs"
	case ErrEnumNoVal:
		return "EnumNoVal"
	case ErrEnumValIllegalIdent:
		return "EnumValIllegalIdent"
	case ErrEnumValRedecl:
		return "EnumValRedecl"
	case ErrUnionRedund:
		return "UnionRedund"
	case ErrUnionSelfref:
		return "UnionSelfref"
	case ErrUnionMissingOpts:
		return "UnionMissingOpts"
	case ErrUnionIncludesNone:
		return "UnionIncludesNone"
	case ErrStructFieldIllegalIdent:
		return "StructFieldIllegalIdent"
	case ErrStructFieldRedecl:
		return "StructFieldRedecl"
	case ErrStructNoFields:
		return "StructNoFields"
	case ErrStructRecurs:
		return "StructRecurs"
	}
	return ""
}

// Error represents a compiler error
type Error interface {
	Code() ErrCode
	Message() string
	Error() string
}

type cErr struct {
	code ErrCode
	msg  string
}

func (e cErr) Error() string   { return e.code.String() + " " + e.msg }
func (e cErr) Code() ErrCode   { return e.code }
func (e cErr) Message() string { return e.msg }
