package parser

import (
	"fmt"
)

// parseEnmVals parses the values of an enumeration type declaration
func (pr *Parser) parseEnmVals(
	lex *Lexer,
	enum *TypeEnum,
) (Fragment, []*EnumValue) {
	// Read '{'
	fBlockBegin, err := readToken(lex, FragTkBlk, "enum value block '{'")
	if pr.err(err) {
		return nil, nil
	}

	frags := []Fragment{fBlockBegin}
	byName := map[string]*Token{}
	values := []*EnumValue{}

	// Parse values
SCAN_LOOP:
	for {
		tk, err := lex.NextSkip(Skip{FragTkSpace})
		if pr.err(err) {
			return nil, nil
		}
		if tk == nil {
			// Unexpected EOF
			pr.err(&pErr{
				at:      lex.Cursor(),
				code:    ErrSyntax,
				message: "unexpected end of file",
			})
			return nil, nil
		}

		frags = append(frags, tk)

		switch tk.id {
		case FragTkLatinAlphanum:
			// An enum value
			tk.id = FragTkEnmVal

			// Make sure the enum value identifier is legal
			if err := lowerCamelCase(tk.src); err != nil {
				pr.err(&pErr{
					at:   tk.begin,
					code: ErrSyntax,
					message: fmt.Sprintf(
						"illegal enum value identifier: %s",
						err,
					),
				})
				return nil, nil
			}
		case FragTkBlkEnd:
			// End of the block
			break SCAN_LOOP
		default:
			// Unexpected token
			pr.err(&pErr{
				at:      tk.begin,
				code:    ErrSyntax,
				message: fmt.Sprintf("unexpected token '%s'", tk.src),
			})
			return nil, nil
		}

		value := tk.src

		// Check for duplicate values
		if defined, isDefined := byName[value]; isDefined {
			pr.err(&pErr{
				at:   tk.begin,
				code: ErrEnumValRedecl,
				message: fmt.Sprintf(
					"Redeclaration of enum value %s "+
						"(previously declared at %s)",
					value,
					defined.begin,
				),
			})
			continue
		}

		byName[value] = tk

		// Add enum value
		values = append(values, &EnumValue{
			Src:  tk,
			Name: value,
			Enum: enum,
		})
	}

	// Make sure there's at least 1 value
	if len(values) < 1 {
		pr.err(&pErr{
			at:   fBlockBegin.begin,
			code: ErrEnumNoVal,
			message: fmt.Sprintf(
				"enum %s is missing values",
				enum.terminalType.TypeName,
			),
		})
		return nil, nil
	}

	return NewConstruct(lex, FragEnmVals, frags...), values
}
