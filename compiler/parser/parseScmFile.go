package parser

import "fmt"

// parseScmFile parses a schema file
func (pr *Parser) parseScmFile(lex *Lexer) Fragment {
	fDeclScm := pr.parseDeclScm(lex)
	if fDeclScm == nil {
		return nil
	}

	frags := []Fragment{fDeclScm}

	// Read declarations by peeking for 1 token
	for {
		tk, err := lex.New().NextSkip(Skip{FragTkSpace})
		if pr.err(err) {
			return nil
		}
		if tk == nil {
			break
		}

		frags = append(frags, tk)

		var frag Fragment
		switch tk.id {
		case FragTkLatinAlphanum:
			// A keyword?
			switch tk.src {
			case KeywordAlias:
				// Alias type declaration
				if f := pr.parseDeclAls(lex); f != nil {
					frag = f.Src
				} else {
					return nil
				}
			case KeywordEnum:
				// Enum type declaration
				if f := pr.parseDeclEnm(lex); f != nil {
					frag = f.Src
				} else {
					return nil
				}
			case KeywordUnion:
				// Union type declaration
				if f := pr.parseDeclUnn(lex); f != nil {
					frag = f.Src
				} else {
					return nil
				}
			case KeywordStruct:
				// Struct type declaration
				if f := pr.parseDeclStr(lex); f != nil {
					frag = f.Src
				} else {
					return nil
				}
			case KeywordResolver:
				// Resolver type declaration
				if f := pr.parseDeclRsv(lex); f != nil {
					frag = f.Src
				} else {
					return nil
				}
			case KeywordTrait:
				// Trait type declaration
				panic("trait types are not yet implemented")
			case KeywordQuery:
				// Query endpoint declaration
				if f := pr.parseDeclQry(lex); f != nil {
					frag = f.Src
				} else {
					return nil
				}
			case KeywordMutation:
				// Mutation endpoint declaration
				if f := pr.parseDeclMut(lex); f != nil {
					frag = f.Src
				} else {
					return nil
				}
			case KeywordSubscription:
				// Subscription endpoint declaration
				panic("subscriptions are not yet implemented")
			default:
				pr.err(&pErr{
					at:   tk.begin,
					code: ErrSyntax,
					message: fmt.Sprintf(
						"unexpected token '%s', expected a declaration",
						tk.src,
					),
				})
				return nil
			}
			frags = append(frags, frag)
		}
	}

	return NewConstruct(lex, FragScmFile, frags...)
}
