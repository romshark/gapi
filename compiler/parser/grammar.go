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

	symParOpen := termEx{
		Kind:        FragTkSymParOpen,
		Expectation: []rune("("),
	}

	symParClose := termEx{
		Kind:        FragTkSymParClose,
		Expectation: []rune(")"),
	}

	symSep := termEx{
		Kind:        FragTkSymSep,
		Expectation: []rune(","),
	}

	symOpt := termEx{
		Kind:        FragTkSymOpt,
		Expectation: []rune("?"),
	}

	symList := termEx{
		Kind:        FragTkSymList,
		Expectation: []rune("[]"),
	}

	keywordSchema := termEx{
		Kind:        FragTkKwdScm,
		Expectation: []rune("schema"),
	}

	keywordStruct := termEx{
		Kind:        FragTkKwdStr,
		Expectation: []rune("struct"),
	}

	keywordEnum := termEx{
		Kind:        FragTkKwdEnm,
		Expectation: []rune("enum"),
	}

	/* keywordResolver := termEx{
		Kind:        FragTkKwdRsv,
		Expectation: []rune("resolver"),
	} */

	keywordQuery := termEx{
		Kind:        FragTkKwdQry,
		Expectation: []rune("query"),
	}

	optSpace := opt{Pattern: term(misc.FrSpace)}

	// File header
	ruleFileHeader := &parser.Rule{
		Kind:        FragDeclSchema,
		Designation: "file header",
		Pattern: seq{
			keywordSchema,
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
			seq{
				keywordEnum,
				optSpace,
				symBlkOpen,
				onePlus{Pattern: seq{
					optSpace,
					&parser.Rule{
						Kind:        FragTkIdnEnumVal,
						Designation: "enum value",
						Pattern: checked{
							Designation: "enum value",
							Fn:          lowerCamelCase,
						},
					},
				}},
				optSpace,
				symBlkClose,
			},
		},
		Action: pr.onDeclTypeEnum,
	}

	// Struct type declaration
	ruleDeclTypeStruct := &parser.Rule{
		Kind:        FragDeclStr,
		Designation: "struct type declaration",
		Pattern: seq{
			checked{
				Designation: "struct type name",
				Fn:          capitalizedCamelCase,
			},
			optSpace,
			symEq,
			optSpace,
			seq{
				keywordStruct,
				optSpace,
				symBlkOpen,
				onePlus{Pattern: seq{
					optSpace,
					&parser.Rule{
						Kind:        FragStrField,
						Designation: "struct field",
						Pattern: seq{
							checked{
								Designation: "struct field name",
								Fn:          lowerCamelCase,
							},
							term(misc.FrSpace),
							ruleTypeDesig,
						},
					},
				}},
				optSpace,
				symBlkClose,
			},
		},
		Action: pr.onDeclTypeStruct,
	}

	ruleParameter := &parser.Rule{
		Kind:        FragParam,
		Designation: "parameter",
		Pattern: seq{
			// Parameter body
			checked{
				Designation: "parameter name",
				Fn:          lowerCamelCase,
			},
			term(misc.FrSpace),
			ruleTypeDesig,
		},
	}

	// Parameter list
	ruleParameterList := &parser.Rule{
		Kind:        FragParams,
		Designation: "parameters list",
	}
	ruleParameterList.Pattern = either{
		seq{
			ruleParameter,
			optSpace,
			symSep,
			optSpace,
			ruleParameterList,
		},
		seq{
			ruleParameter,
			opt{Pattern: symSep},
		},
	}

	// Parameters
	ruleParameters := &parser.Rule{
		Kind:        FragParams,
		Designation: "parameters",
		Pattern: seq{
			symParOpen,
			optSpace,
			zeroPlus{Pattern: ruleParameterList},
			optSpace,
			symParClose,
		},
	}

	// Resolver type declaration
	ruleQueryDecl := &parser.Rule{
		Kind:        FragDeclQry,
		Designation: "query endpoint declaration",
		Pattern: seq{
			checked{
				Designation: "query endpoint name",
				Fn:          lowerCamelCase,
			},
			optSpace,
			symEq,
			optSpace,
			keywordQuery,
			optSpace,
			opt{Pattern: ruleParameters},
			optSpace,
			ruleTypeDesig,
		},
		Action: pr.onDeclQuery,
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
					ruleDeclTypeStruct,
					ruleQueryDecl,
				},
			}},
			optSpace,
		},
	}
}
