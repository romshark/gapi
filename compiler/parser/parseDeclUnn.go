package parser

func (pr *Parser) parseDeclUnn(lex *Lexer) *TypeUnion {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordUnion,
		FragTkKwdUnn,
		"keyword",
	)
	if pr.err(err) {
		return nil
	}

	// Read type ID
	fType, err := readWord(
		lex,
		"union type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Create a new resolver type instance
	newUnion := &TypeUnion{
		terminalType: terminalType{
			Name: fType.src,
		},
	}

	// Parse option types and set them when they're resolved
	fOpts := pr.parseUnnOpts(lex, newUnion, func(ts []Type) {
		newUnion.Types = ts
	})
	if fOpts == nil {
		return nil
	}

	newUnion.Src = NewConstruct(lex, FragDeclUnn,
		fDeclKeyword,
		fType,
		fOpts,
	)

	// Define the type
	pr.onTypeDecl(newUnion)

	return newUnion
}
