package parser

import (
	"fmt"
)

// parseParams parses a block-list of parameters
func (pr *Parser) parseParams(
	lex *Lexer,
	target GraphNode,
) (Fragment, []*Parameter) {
	// Read block opening '('
	fBlockOpening, err := readToken(
		lex,
		FragTkPar,
		"parameter block opening '('",
	)
	if pr.err(err) {
		return nil, nil
	}

	frags := []Fragment{fBlockOpening}
	byName := map[string]*Token{}
	params := []*Parameter{}

	// Parse parameters
SCAN_LOOP:
	for {
		peeker := lex.New()
		// Peek for 1 token to find out whether
		// the block ended or a new parameter began
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
			// A parameter
			newParam := pr.parseParam(lex, target)
			if newParam == nil {
				return nil, nil
			}
			frags = append(frags, newParam.Src)
			params = append(params, newParam)
		case FragTkParEnd:
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

		paramName := tk.src

		// Check for redeclarations
		if defined, isDefined := byName[paramName]; isDefined {
			pr.err(&pErr{
				at:   tk.begin,
				code: ErrParamRedecl,
				message: fmt.Sprintf(
					"Redeclaration of parameter %s "+
						"(previously declared at %s)",
					paramName,
					defined.begin,
				),
			})
			return nil, nil
		}

		byName[paramName] = tk

		// Skip separator if any
		sepTk, err := lex.New().NextSkip(Skip{FragTkSpace})
		if pr.err(err) {
			return nil, nil
		}
		if sepTk.id == FragTkSymSep {
			separator, err := lex.NextSkip(Skip{FragTkSpace})
			if pr.err(err) {
				return nil, nil
			}
			frags = append(frags, separator)
		}
	}

	return NewConstruct(lex, FragParams, frags...), params
}
