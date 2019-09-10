package parser

import (
	parser "github.com/romshark/llparser"
	"github.com/romshark/llparser/misc"
)

const (
	_ parser.FragmentKind = misc.FrSign + iota

	/* TERMINALS (Primitive) */

	// FragTkBlk represents a '{' token
	FragTkBlk

	// FragTkBlkEnd represents a '}' token
	FragTkBlkEnd

	// FragTkSeq represents a sequence of non-space characters
	FragTkSeq

	// FragTkMemAcc represents a member-accessor '.' token
	FragTkMemAcc

	// FragTkDocLineInit represents a documentation line initiator '#'
	FragTkDocLineInit

	// FragTkSym represents any special character token
	FragTkSym

	// FragTkSymSep represents a separator ',' token
	FragTkSymSep

	// FragTkSymEq represents an equals symbol '='
	FragTkSymEq

	// FragTkSymBlockOpen represents a block opening symbol '{'
	FragTkSymBlockOpen

	// FragTkSymBlockClose represents a block closing symbol '}'
	FragTkSymBlockClose

	// FragTkSymParOpen represents a '(' token
	FragTkSymParOpen

	// FragTkSymParClose represents a ')' token
	FragTkSymParClose

	// FragTkSymOpt represents an optionality symbol '?'
	FragTkSymOpt

	// FragTkSymList represents a list symbol '[]'
	FragTkSymList

	/* TERMINALS (Keywords) */

	// FragTkKwdScm represents a schema declaration keyword fragment
	FragTkKwdScm

	// FragTkKwdEnm represents an enum type declaration keyword fragment
	FragTkKwdEnm

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
