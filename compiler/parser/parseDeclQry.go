package parser

func (pr *Parser) parseDeclQry(lex *Lexer) *Query {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordQuery,
		FragTkKwdQry,
		"keyword",
	)
	if pr.err(err) {
		return nil
	}

	// Read endpoint name
	fName, err := readWord(
		lex,
		"endpoint name",
		FragTkIdnProp,
		lowerCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Create a new query endpoint instance
	newQuery := &Query{
		Name: fName.src,
	}

	// Parse parameters
	fParams, params, parsed := pr.parseOptParams(lex, newQuery)
	if !parsed {
		return nil
	}
	newQuery.Parameters = params

	// Read type ID
	fType := pr.parseTypeDesig(lex, func(t Type) {
		if _, isNone := t.(TypeStdNone); isNone {
			pr.err(&pErr{
				at:      fDeclKeyword.begin,
				code:    ErrSyntax,
				message: "Query endpoint resolves to None",
			})
		}
		newQuery.Type = t
	})
	if fType == nil {
		return nil
	}

	newQuery.Src = NewConstruct(lex, FragDeclQry,
		fDeclKeyword,
		fName,
		fParams,
		fType,
	)

	// Define the endpoint
	if !pr.onGraphNode(newQuery) {
		return nil
	}

	return newQuery
}
