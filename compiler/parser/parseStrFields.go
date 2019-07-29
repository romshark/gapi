package parser

import "fmt"

// parseStrFields parses the fields block of a struct declaration
func (pr *Parser) parseStrFields(
	lex *Lexer,
	structType *TypeStruct,
) (Fragment, []*StructField) {
	// Read '{'
	fBlockBegin, err := readToken(
		lex,
		FragTkBlk,
		"struct fields block opening '{'",
	)
	if pr.err(err) {
		return nil, nil
	}

	frags := []Fragment{fBlockBegin}
	byName := map[string]*Token{}
	fields := []*StructField{}

	// Parse fields
SCAN_LOOP:
	for {
		peeker := lex.New()
		// Peek for 1 token to find out whether
		// the block ended or a new field began
		tk, err := peeker.NextSkip(Skip{FragTkSpace})
		if pr.err(err) {
			return nil, nil
		}
		if tk == nil {
			// Unexpected EOF
			pr.err(&pErr{
				at:      peeker.Cursor(),
				code:    ErrSyntax,
				message: "unexpected end of file",
			})
			return nil, nil
		}

		switch tk.id {
		case FragTkLatinAlphanum:
			// A field
			newField := pr.parseStrField(lex, structType)
			if newField == nil {
				return nil, nil
			}
			frags = append(frags, newField.Src)
			fields = append(fields, newField)
		case FragTkBlkEnd:
			// End of the block
			frags = append(frags, tk)
			_, _ = lex.NextSkip(Skip{FragTkSpace})
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

		fieldName := tk.src

		// Check for redeclarations
		if defined, isDefined := byName[fieldName]; isDefined {
			pr.err(&pErr{
				at:   tk.begin,
				code: ErrStructFieldRedecl,
				message: fmt.Sprintf(
					"Redeclaration of struct field %s "+
						"(previously declared at %s)",
					fieldName,
					defined.begin,
				),
			})
			return nil, nil
		}

		byName[fieldName] = tk
	}

	// Make sure there's at least 1 field
	if len(fields) < 1 {
		pr.err(&pErr{
			at:   fBlockBegin.begin,
			code: ErrStructNoFields,
			message: fmt.Sprintf(
				"struct %s is missing fields",
				structType.Name,
			),
		})
		return nil, nil
	}

	return NewConstruct(lex, FragStrFields, frags...), fields
}
