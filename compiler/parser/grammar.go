package parser

import (
	parser "github.com/romshark/llparser"
	"github.com/romshark/llparser/misc"
)

type seq = parser.Sequence
type opt = parser.Optional
type onePlus = parser.OneOrMore
type zeroPlus = parser.ZeroOrMore
type termEx = parser.TermExact
type term = parser.Term
type checked = parser.Checked
type either = parser.Either

// Grammar returns the language grammar
func (pr *Parser) Grammar() *parser.Rule {
	symEq := termEx{
		Kind:        FragTkSymEq,
		Expectation: []rune("="),
	}

	symBlkOpen := termEx{
		Kind:        FragTkSymBlockOpen,
		Expectation: []rune("{"),
	}

	symBlkClose := termEx{
		Kind:        FragTkSymBlockClose,
		Expectation: []rune("}"),
	}

	symOpt := termEx{
		Kind:        FragTkSymOpt,
		Expectation: []rune("?"),
	}

	symList := termEx{
		Kind:        FragTkSymList,
		Expectation: []rune("[]"),
	}

	optSpace := opt{Pattern: term(misc.FrSpace)}

	// File header
	ruleFileHeader := &parser.Rule{
		Kind:        FragDeclSchema,
		Designation: "file header",
		Pattern: seq{
			termEx{
				Kind:        FragTkKwdScm,
				Expectation: []rune("schema"),
			},
			term(misc.FrSpace),
			checked{
				Designation: "schema name",
				Fn:          lowerCamelCase,
			},
		},
		Action: pr.onFileHeader,
	}

	// Type designation
	ruleTypeDesig := &parser.Rule{
		Kind:        FragType,
		Designation: "type designation",
		Pattern: seq{
			zeroPlus{Pattern: either{
				seq{
					symOpt,
					symList,
				},
				symOpt,
				symList,
			}},
			checked{
				Designation: "alias type name",
				Fn:          capitalizedCamelCase,
			},
		},
	}

	// Alias type declaration
	ruleDeclTypeAlias := &parser.Rule{
		Kind:        FragDeclAls,
		Designation: "alias type declaration",
		Pattern: seq{
			checked{
				Designation: "alias type name",
				Fn:          capitalizedCamelCase,
			},
			optSpace,
			symEq,
			optSpace,
			ruleTypeDesig,
		},
		Action: pr.onDeclTypeAlias,
	}

	// Enum values block
	ruleEnumValBlock := &parser.Rule{
		Kind:        100,
		Designation: "enum values block",
		Pattern: seq{
			symBlkOpen,
			seq{
				optSpace,
				checked{
					Designation: "enum value",
					Fn:          lowerCamelCase,
				},
			},
			optSpace,
			symBlkClose,
		},
	}

	// Enum type declaration
	ruleDeclTypeEnum := &parser.Rule{
		Kind:        FragDeclEnm,
		Designation: "enum type declaration",
		Pattern: seq{
			checked{
				Designation: "enum type name",
				Fn:          capitalizedCamelCase,
			},
			optSpace,
			symEq,
			optSpace,
			ruleEnumValBlock,
		},
		Action: pr.onDeclTypeEnum,
	}

	// File rule
	return &parser.Rule{
		Designation: "schema file",
		Pattern: seq{
			ruleFileHeader,
			zeroPlus{Pattern: seq{
				term(misc.FrSpace),
				either{
					ruleDeclTypeAlias,
					ruleDeclTypeEnum,
				},
			}},
			optSpace,
		},
	}
}
