package parser

func (pr *Parser) parseDeclMut(lex *Lexer) *Mutation {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordMutation,
		FragTkKwdMut,
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
	newMutation := &Mutation{
		Name: fName.src,
	}

	// Parse parameters
	fParams, params, parsed := pr.parseOptParams(lex, newMutation)
	if !parsed {
		return nil
	}
	newMutation.Parameters = params

	// Read type ID
	fType := pr.parseTypeDesig(lex, func(t Type) { newMutation.Type = t })
	if fType == nil {
		return nil
	}

	newMutation.Src = NewConstruct(lex, FragDeclMut,
		fDeclKeyword,
		fName,
		fParams,
		fType,
	)

	// Define the endpoint
	if !pr.onGraphNode(newMutation) {
		return nil
	}

	return newMutation
}
