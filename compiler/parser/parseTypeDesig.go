package parser

import "fmt"

// parseTypeDesig parses a type reference fragment.
// returns an incomplete type that's completed
func (pr *Parser) parseTypeDesig(
	lex *Lexer,
	onTypeResolved func(Type),
) Fragment {
	var tp Type
	var previousTp Type
	frags := []Fragment{}

	appendTp := func(t Type) {
		if tp == nil {
			tp = t
		} else {
			switch v := previousTp.(type) {
			case *TypeOptional:
				v.StoreType = t
			case *TypeList:
				v.StoreType = t
			}
		}
		previousTp = t
	}

	// Parse chain
SCAN_LOOP:
	for {
		tk, err := lex.NextSkip(Skip{FragTkSpace})
		if pr.err(err) {
			return nil
		}
		if tk == nil {
			// Unexpected EOF
			pr.err(&pErr{
				at:      lex.Cursor(),
				code:    ErrSyntax,
				message: "unexpected end of file",
			})
			return nil
		}

		switch tk.id {
		case FragTkSymList:
			// List container type
			appendTp(&TypeList{})
			frags = append(frags, tk)

		case FragTkSymOpt:
			// Optional container type
			// Ensure the previous type in the chain was not also an optional
			if previousTp != nil {
				if _, tailIsOpt := previousTp.(*TypeOptional); tailIsOpt {
					// Illegal optionals chain detected
					// (Optional type of optional types)
					pr.err(&pErr{
						at:   tk.begin,
						code: ErrTypeOptChain,
						message: fmt.Sprintf(
							"illegal chain of optionals " +
								"(optional type of optional types)",
						),
					})
					return nil
				}
			}
			appendTp(&TypeOptional{})
			frags = append(frags, tk)

		case FragTkLatinAlphanum:
			// Terminal type reached

			// Make sure the type identifier is legal
			if err := capitalizedCamelCase(tk.src); err != nil {
				pr.err(&pErr{
					at:   tk.begin,
					code: ErrSyntax,
					message: fmt.Sprintf(
						"illegal type identifier: %s",
						err,
					),
				})
				return nil
			}

			tk.id = FragTkIdnType
			frags = append(frags, tk)

			pr.deferJob(func() {
				// Make sure the terminal type is defined
				terminalType := pr.findTypeByDesignation(tk.src)
				if terminalType == nil {
					pr.err(&pErr{
						at:   tk.begin,
						code: ErrTypeUndef,
						message: fmt.Sprintf(
							"terminal type %s is undefined",
							tk.src,
						),
					})
					return
				}
				appendTp(terminalType)

				// Reference the terminal type in the type chain
				for t := tp; t != nil; {
					if v, isOpt := t.(*TypeOptional); isOpt {
						v.Terminal = previousTp
						t = v.StoreType
						continue
					}
					if v, isList := t.(*TypeList); isList {
						v.Terminal = previousTp
						t = v.StoreType
						continue
					}
					break
				}

				_, terminalIsNone := previousTp.(TypeStdNone)
				if _, isNone := tp.(TypeStdNone); !isNone && terminalIsNone {
					pr.err(&pErr{
						at:      frags[0].Begin(),
						code:    ErrSyntax,
						message: "illegal None-type",
					})
				}

				if tp.TerminalType() != nil {
					tp = pr.onAnonymousType(tp)
				}
				onTypeResolved(tp)

			})
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
	}

	return NewConstruct(lex, FragType, frags...)
}
