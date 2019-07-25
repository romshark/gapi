package parser

// parseOptParams parses an optional block-list of parameters if there are any
func (pr *Parser) parseOptParams(
	lex *Lexer,
	target GraphNode,
) (Fragment, []*Parameter, bool) {
	// Peek for 1 token to find out whether there is a parameter block
	next, err := lex.New().Next()
	if pr.err(err) {
		return nil, nil, false
	}
	if next == nil || next.id != FragTkPar {
		// Not an opening parenthesis, no parameters
		return nil, nil, true
	}

	frag, params := pr.parseParams(lex, target)
	if frag == nil {
		return nil, nil, false
	}
	return frag, params, true
}
