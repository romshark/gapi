package parser

import (
	"fmt"
	"strings"

	parser "github.com/romshark/llparser"
)

// ErrCode represents a compiler error code
type ErrCode int

const (
	_ ErrCode = iota

	// ErrSyntax represents a syntax error
	ErrSyntax

	// ErrGraphRootNodeRedecl indicates a redeclared graph root node
	ErrGraphRootNodeRedecl

	// ErrTypeUndef indicates an undefined referenced type
	ErrTypeUndef

	// ErrTypeRedecl indicates a redeclared type
	ErrTypeRedecl

	// ErrTypeOptChain indicates an illegal chain of optionals
	ErrTypeOptChain

	// ErrAliasRecurs indicates a recursive alias type
	ErrAliasRecurs

	// ErrEnumNoVal indicates an empty enumeration missing values
	ErrEnumNoVal

	// ErrEnumValRedecl indicates a redeclared enum value
	ErrEnumValRedecl

	// ErrUnionRedund indicates redundant union option types
	ErrUnionRedund

	// ErrUnionRecurs indicates a self-referencing union type
	ErrUnionRecurs

	// ErrUnionMissingOpts indicates a union type with too little option types
	ErrUnionMissingOpts

	// ErrUnionIncludesNone indicates a union type including the None type
	ErrUnionIncludesNone

	// ErrStructFieldRedecl indicates a redeclared struct field
	ErrStructFieldRedecl

	// ErrStructFieldImpure indicates a struct field of impure (non-data) type
	ErrStructFieldImpure

	// ErrStructNoFields indicates an empty struct type missing fields
	ErrStructNoFields

	// ErrStructRecurs indicates a recursive struct type
	ErrStructRecurs

	// ErrResolverPropRedecl indicates a redeclared resolver property
	ErrResolverPropRedecl

	// ErrResolverNoProps indicates an empty resolver type missing properties
	ErrResolverNoProps

	// ErrParamImpure indicates a parameter of a non-data (impure) type
	ErrParamImpure

	// ErrParamRedecl indicates a redeclared resolver property parameter
	ErrParamRedecl

	// ErrNoEndpoints indicates the absence of any API endpoints
	ErrNoEndpoints
)

// String stringifies the error code
func (c ErrCode) String() string {
	switch c {
	case ErrSyntax:
		return "Syntax"
	case ErrGraphRootNodeRedecl:
		return "GraphRootNodeRedecl"
	case ErrTypeUndef:
		return "TypeUndef"
	case ErrTypeRedecl:
		return "TypeRedecl"
	case ErrTypeOptChain:
		return "TypeOptChain"
	case ErrAliasRecurs:
		return "AliasRecurs"
	case ErrEnumNoVal:
		return "EnumNoVal"
	case ErrEnumValRedecl:
		return "EnumValRedecl"
	case ErrUnionRedund:
		return "UnionRedund"
	case ErrUnionRecurs:
		return "UnionSelfref"
	case ErrUnionMissingOpts:
		return "UnionMissingOpts"
	case ErrUnionIncludesNone:
		return "UnionIncludesNone"
	case ErrStructFieldRedecl:
		return "StructFieldRedecl"
	case ErrStructFieldImpure:
		return "StructFieldImpure"
	case ErrStructNoFields:
		return "StructNoFields"
	case ErrStructRecurs:
		return "StructRecurs"
	case ErrResolverPropRedecl:
		return "ResolverPropRedecl"
	case ErrResolverNoProps:
		return "ResolverNoProps"
	case ErrParamImpure:
		return "ParamImpure"
	case ErrParamRedecl:
		return "ParamRedecl"
	case ErrNoEndpoints:
		return "NoEndpoints"
	}
	return ""
}

// Error represents a generic parser error
type Error interface {
	error

	// Code returns the error code
	Code() ErrCode

	// Message returns the error message without the location appended
	Message() string

	// At returns the error position in the source code
	At() parser.Cursor
}

// pErr represents a syntax error
type pErr struct {
	code    ErrCode
	message string
	at      parser.Cursor
}

func (err *pErr) Error() string {
	return fmt.Sprintf(
		"%s: %s at %s",
		err.code,
		err.message,
		err.at.String(),
	)
}

// Code returns the error code
func (err *pErr) Code() ErrCode { return err.code }

// Message returns the error message without the location appended
func (err *pErr) Message() string { return err.message }

// At returns the error position in the source code
func (err *pErr) At() parser.Cursor { return err.at }

// ParseErr represents a parsing error
type ParseErr struct {
	Errors []Error
}

func (e ParseErr) Error() string {
	s := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		s[i] = fmt.Sprintf("%d: %s", i+1, err.Error())
	}
	return fmt.Sprintf(
		"%d compilation errors: [%s]",
		len(s),
		strings.Join(s, "; "),
	)
}
