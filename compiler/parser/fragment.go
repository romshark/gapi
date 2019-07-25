package parser

import "fmt"

// Cursor represents a source-code location
type Cursor struct {
	Index  uint32
	Line   uint32
	Column uint32
	File   *File
}

// String returns the stringified cursor
func (crs Cursor) String() string {
	if crs.File == nil {
		return "<unknown>"
	}
	return fmt.Sprintf(
		"%s:%d:%d",
		crs.File.Name,
		crs.Line,
		crs.Column,
	)
}

// Fragment represents a typed source code fragment
type Fragment interface {
	Begin() Cursor
	End() Cursor
	Src() string
	Elements() []Fragment
}

// FragID represents a GAPI source code fragment identifier
type FragID int32

const (
	_ FragID = iota

	/* TERMINALS (Primitive) */

	// FragTkSpace represents a space sequence token
	// (any combinations of spaces, tabs & line-breaks)
	FragTkSpace

	// FragTkBlk represents a '{' token
	FragTkBlk

	// FragTkBlkEnd represents a '}' token
	FragTkBlkEnd

	// FragTkPar represents a '(' token
	FragTkPar

	// FragTkParEnd represents a ')' token
	FragTkParEnd

	// FragTkSeq represents a sequence of non-space characters
	FragTkSeq

	// FragTkLatinAlphanum represents an alpha-numeric word token ([a-zA-Z])
	FragTkLatinAlphanum

	// FragTkMemAcc represents a member-accessor '.' token
	FragTkMemAcc

	// FragTkDocLineInit represents a documentation line initiator '#'
	FragTkDocLineInit

	// FragTkSymSep represents a separator ',' token
	FragTkSymSep

	// FragTkSymEq represents an equals symbol '='
	FragTkSymEq

	// FragTkSymOpt represents an optionality symbol '?'
	FragTkSymOpt

	// FragTkSymList represents a list symbol '[]'
	FragTkSymList

	/* TERMINALS (Keywords) */

	// FragTkKwdScm represents a schema declaration keyword fragment
	FragTkKwdScm

	// FragTkKwdEnm represents an enum type declaration keyword fragment
	FragTkKwdEnm

	// FragTkKwdAls represents a alias type declaration keyword fragment
	FragTkKwdAls

	// FragTkKwdUnn represents a union type declaration keyword fragment
	FragTkKwdUnn

	// FragTkKwdStr represents a struct type declaration keyword fragment
	FragTkKwdStr

	// FragTkKwdRsv represents a resolver type declaration keyword fragment
	FragTkKwdRsv

	// FragTkKwdTrt represents a trait type declaration keyword fragment
	FragTkKwdTrt

	// FragTkKwdQry represents a query endpoint declaration keyword fragment
	FragTkKwdQry

	// FragTkKwdMut represents a mutation endpoint declaration keyword fragment
	FragTkKwdMut

	// FragTkKwdSub represents a subscription endpoint declaration keyword
	// fragment
	FragTkKwdSub

	// FragTkIdnScm represents a schema identifier fragment
	FragTkIdnScm

	// FragTkIdnType represents a type identifier fragment
	FragTkIdnType

	// FragTkIdnProp represents a property identifier fragment
	FragTkIdnProp

	// FragTkIdnFld represents a struct field identifier fragment
	FragTkIdnFld

	// FragTkIdnParam represents a parameter identifier fragment
	FragTkIdnParam

	// FragTkIdnEnumVal represents an enum value identifier fragment
	FragTkIdnEnumVal

	// FragTkEnmVal represents an enum value fragment
	FragTkEnmVal

	/* NON-TERMINALS (declarations) */

	// FragScmFile represents a schema file fragment
	FragScmFile

	// FragDeclSchema represents a schema declaration fragment
	FragDeclSchema

	// FragDeclAls represents an alias type declaration fragment
	FragDeclAls

	// FragDeclEnm represents an enum type declaration fragment
	FragDeclEnm

	// FragDeclUnn represents a union type declaration fragment
	FragDeclUnn

	// FragDeclStr represents a struct type declaration fragment
	FragDeclStr

	// FragDeclRsv represents a resolver type declaration fragment
	FragDeclRsv

	// FragDeclTrt represents a trait type declaration fragment
	FragDeclTrt

	// FragDeclQry represents a query endpoint declaration fragment
	FragDeclQry

	// FragDeclMut represents a mutation endpoint type declaration fragment
	FragDeclMut

	// FragDeclSub represents a subscription endpoint type declaration fragment
	FragDeclSub

	// FragEnmVals represents an enum values block fragment
	FragEnmVals

	// FragRsvProps represents a resolver properties block fragment
	FragRsvProps

	// FragRsvProp represents a resolver property fragment
	FragRsvProp

	// FragParams represents a parameter block-list fragment
	FragParams

	// FragParam represents a parameter fragment
	FragParam

	// FragStrFields represents a struct fields block fragment
	FragStrFields

	// FragStrField represents a struct field fragment
	FragStrField

	// FragUnnOpts represents the option-types block of a union type
	// declaration fragment
	FragUnnOpts

	// FragType represents a type definition
	FragType
)

// String stringifies the fragment identifier
func (id FragID) String() string {
	switch id {
	case FragTkSpace:
		return "TkSpace"
	case FragTkBlk:
		return "TkBlk"
	case FragTkBlkEnd:
		return "TkBlkEnd"
	case FragTkPar:
		return "TkPar"
	case FragTkParEnd:
		return "TkParEnd"
	case FragTkSeq:
		return "TkSeq"
	case FragTkLatinAlphanum:
		return "TkLatinAlphanum"
	case FragTkMemAcc:
		return "TkMemAcc"
	case FragTkDocLineInit:
		return "TkDocLineInit"
	case FragTkSymSep:
		return "TkSymSep"
	case FragTkSymEq:
		return "TkSymEq"
	case FragTkSymOpt:
		return "TkSymOpt"
	case FragTkSymList:
		return "TkSymList"
	case FragTkKwdScm:
		return "TkKwdScm"
	case FragTkKwdEnm:
		return "TkKwdEnm"
	case FragTkKwdAls:
		return "TkKwdAls"
	case FragTkKwdUnn:
		return "TkKwdUnn"
	case FragTkKwdStr:
		return "TkKwdStr"
	case FragTkKwdRsv:
		return "TkKwdRsv"
	case FragTkKwdTrt:
		return "TkKwdTrt"
	case FragTkKwdQry:
		return "TkKwdQry"
	case FragTkKwdMut:
		return "TkKwdMut"
	case FragTkKwdSub:
		return "TkKwdSub"
	case FragTkIdnScm:
		return "TkIdnScm"
	case FragTkIdnType:
		return "TkIdnType"
	case FragTkIdnProp:
		return "TkIdnProp"
	case FragTkIdnFld:
		return "TkIdnFld"
	case FragTkIdnParam:
		return "TkIdnParam"
	case FragTkIdnEnumVal:
		return "TkIdnEnumVal"
	case FragTkEnmVal:
		return "TkEnmVal"
	case FragScmFile:
		return "ScmFile"
	case FragDeclSchema:
		return "DeclSchema"
	case FragDeclAls:
		return "DeclAls"
	case FragDeclEnm:
		return "DeclEnm"
	case FragDeclUnn:
		return "DeclUnn"
	case FragDeclStr:
		return "DeclStr"
	case FragDeclRsv:
		return "DeclRsv"
	case FragDeclTrt:
		return "DeclTrt"
	case FragDeclQry:
		return "DeclQry"
	case FragDeclMut:
		return "DeclMut"
	case FragDeclSub:
		return "DeclSub"
	case FragEnmVals:
		return "EnmVals"
	case FragRsvProps:
		return "RsvProps"
	case FragRsvProp:
		return "RsvProp"
	case FragParams:
		return "Params"
	case FragParam:
		return "Param"
	case FragStrFields:
		return "StrFields"
	case FragStrField:
		return "StrField"
	case FragUnnOpts:
		return "UnnOpts"
	case FragType:
		return "Type"
	}
	return ""
}
