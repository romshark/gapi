package parser

import "fmt"

// parseUnnOpts parses the type options block of a union type declaration
func (pr *Parser) parseUnnOpts(
	lex *Lexer,
	unionType *TypeUnion,
	onTypesResolved func(t []Type),
) Fragment {
	// Read '{'
	fBlockBegin, err := readToken(
		lex,
		FragTkBlk,
		"union option-types block '{'",
	)
	if pr.err(err) {
		return nil
	}

	frags := []Fragment{fBlockBegin}
	byDescription := map[string]*Token{}
	typeOptions := []Type{}
	typeOptionsNum := 0

	// Parse type options
SCAN_LOOP:
	for {
		peeker := lex.New()
		// Peek for 1 token to find out whether
		// the block ended or a new type began
		tk, err := peeker.NextSkip(Skip{FragTkSpace})
		if pr.err(err) {
			return nil
		}
		if tk == nil {
			// Unexpected EOF
			pr.err(&pErr{
				at:      peeker.Cursor(),
				code:    ErrSyntax,
				message: "unexpected end of file",
			})
			return nil
		}

		frags = append(frags, tk)

		switch tk.id {
		case FragTkLatinAlphanum:
			// A type (terminal type)
		case FragTkSymList:
			// A type (list of...)
		case FragTkSymOpt:
			// A type (optional ...)
		case FragTkBlkEnd:
			// End of the block
			_, _ = lex.NextSkip(Skip{FragTkSpace})
			break SCAN_LOOP
		default:
			// Unexpected token
			pr.err(&pErr{
				at:      tk.begin,
				code:    ErrSyntax,
				message: fmt.Sprintf("unexpected token '%s'", tk.src),
			})
			return nil
		}

		// Parse the type and remember it
		typeOptionsNum++
		fOption := pr.parseTypeDesig(lex, func(t Type) {
			typeOptions = append(typeOptions, t)
			typeOptionsNum--
			// Execute onTypesResolved if/when all types are resolved
			if typeOptionsNum < 1 {
				terminalTypeName := t.String()

				// Ensure the union doesn't reference itself as an option
				if terminalTypeName == unionType.Name {
					pr.err(&pErr{
						at:   tk.begin,
						code: ErrUnionRecurs,
						message: fmt.Sprintf(
							"union type %s references itself "+
								"as one of its options",
							unionType,
						),
					})
					return
				}

				// Ensure the union type doesn't include None as an option
				if _, isNone := t.(TypeStdNone); isNone {
					pr.err(&pErr{
						at:   tk.begin,
						code: ErrUnionIncludesNone,
						message: fmt.Sprintf(
							"union type %s includes the None primitive",
							unionType,
						),
					})
					return
				}

				onTypesResolved(typeOptions)
			}
		})
		if fOption == nil {
			return nil
		}

		// Check for duplicate options
		optTypeDesc := fOption.Src()
		if _, isDefined := byDescription[optTypeDesc]; isDefined {
			pr.err(&pErr{
				at:   fOption.Begin(),
				code: ErrUnionRedund,
				message: fmt.Sprintf(
					"redelaration of option-type %s in union type %s",
					optTypeDesc,
					unionType,
				),
			})
			return nil
		}
		byDescription[optTypeDesc] = tk
	}

	// Make sure there's at least 2 options
	if len(byDescription) < 2 {
		pr.err(&pErr{
			at:   fBlockBegin.begin,
			code: ErrUnionMissingOpts,
			message: fmt.Sprintf(
				"union %s is missing options",
				unionType,
			),
		})
		return nil
	}

	return NewConstruct(lex, FragUnnOpts, frags...)
}
