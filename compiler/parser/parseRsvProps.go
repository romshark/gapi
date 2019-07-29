package parser

import "fmt"

// parseRsvProps parses the properties block of a resolver declaration
func (pr *Parser) parseRsvProps(
	lex *Lexer,
	resolver *TypeResolver,
) (Fragment, []*ResolverProperty) {
	// Read '{'
	fBlockBegin, err := readToken(
		lex,
		FragTkBlk,
		"resolver properties block opening '{'",
	)
	if pr.err(err) {
		return nil, nil
	}

	frags := []Fragment{fBlockBegin}
	byName := map[string]*Token{}
	props := []*ResolverProperty{}

	// Parse properties
SCAN_LOOP:
	for {
		peeker := lex.New()
		// Peek for 1 token to find out whether
		// the block ended or a new property began
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
			// A property
			newProp := pr.parseRsvProp(lex, resolver)
			if newProp == nil {
				return nil, nil
			}
			frags = append(frags, newProp.Src)
			props = append(props, newProp)
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

		propName := tk.src

		// Check for redeclarations
		if defined, isDefined := byName[propName]; isDefined {
			pr.err(&pErr{
				at:   tk.begin,
				code: ErrResolverPropRedecl,
				message: fmt.Sprintf(
					"Redeclaration of resolver property %s "+
						"(previously declared at %s)",
					propName,
					defined.begin,
				),
			})
			return nil, nil
		}

		byName[propName] = tk
	}

	// Make sure there's at least 1 property
	if len(props) < 1 {
		pr.err(&pErr{
			at:   fBlockBegin.begin,
			code: ErrResolverNoProps,
			message: fmt.Sprintf(
				"resolver %s is missing properties",
				resolver.Name,
			),
		})
		return nil, nil
	}

	return NewConstruct(lex, FragRsvProps, frags...), props
}
