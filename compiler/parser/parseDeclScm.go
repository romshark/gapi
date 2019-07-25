package parser

func (pr *Parser) parseDeclScm(lex *Lexer) (Fragment, string) {
	// Read schema declaration keyword
	fDeclScmKeyword, err := readWordExact(
		lex,
		KeywordSchema,
		FragTkKwdScm,
		"schema declaration keyword",
	)
	if pr.err(err) {
		return nil, ""
	}

	// Read schema name
	fName, err := readWord(
		lex,
		"schema identifier",
		FragTkIdnScm,
		lowerCamelCase,
	)
	if pr.err(err) {
		return nil, ""
	}

	return NewConstruct(lex, FragDeclSchema,
		fDeclScmKeyword,
		fName,
	), fName.src
}
