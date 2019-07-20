package compiler

import (
	"fmt"
	"strings"
)

// ErrCode represents a compiler error code
type ErrCode int

const (
	_ ErrCode = iota

	// ErrSyntax represents a syntax error
	ErrSyntax

	// ErrSchemaIllegalIdent indicates an illegal schema identifier
	ErrSchemaIllegalIdent

	// ErrQryEndpointIllegalIdent indicates an illegal query endpoint identifier
	ErrQryEndpointIllegalIdent

	// ErrMutEndpointIllegalIdent indicates an illegal mutation endpoint
	// identifier
	ErrMutEndpointIllegalIdent

	// ErrGraphRootNodeRedecl indicates a redeclared graph root node
	ErrGraphRootNodeRedecl

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

	// ErrResolverPropIllegalIdent indicates an illegal resolver property
	// identifier
	ErrResolverPropIllegalIdent

	// ErrResolverPropRedecl indicates a redeclared resolver property
	ErrResolverPropRedecl

	// ErrResolverNoProps indicates an empty resolver type missing properties
	ErrResolverNoProps

	// ErrParamIllegalIdent indicates an illegal parameter identifier
	ErrParamIllegalIdent

	// ErrParamImpure indicates a parameter of a non-data (impure) type
	ErrParamImpure

	// ErrResolverPropParamRedecl indicates a redeclared resolver property
	// parameter
	ErrResolverPropParamRedecl
)

// String stringifies the error code
func (c ErrCode) String() string {
	switch c {
	case ErrSyntax:
		return "Syntax"
	case ErrSchemaIllegalIdent:
		return "SchemaIllegalIdent"
	case ErrQryEndpointIllegalIdent:
		return "QryEndpointIllegalIdent"
	case ErrMutEndpointIllegalIdent:
		return "MutEndpointIllegalIdent"
	case ErrGraphRootNodeRedecl:
		return "GraphRootNodeRedecl"
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
	case ErrResolverPropIllegalIdent:
		return "ResolverPropIllegalIdent"
	case ErrResolverPropRedecl:
		return "ResolverPropRedecl"
	case ErrResolverNoProps:
		return "ResolverNoProps"
	case ErrParamIllegalIdent:
		return "ParamIllegalIdent"
	case ErrParamImpure:
		return "ParamImpure"
	case ErrResolverPropParamRedecl:
		return "ResolverPropParamRedecl"
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

// CompilationErr represents a compilation error
type CompilationErr struct {
	Errors []Error
}

func (e CompilationErr) Error() string {
	s := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		s[i] = fmt.Sprintf("%d: %s(%s)", i+1, err.Code(), err.Message())
	}
	return fmt.Sprintf(
		"%d compilation errors: [%s]",
		len(s),
		strings.Join(s, "; "),
	)
}
