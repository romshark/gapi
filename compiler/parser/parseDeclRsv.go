package parser

func (pr *Parser) parseDeclRsv(lex *Lexer) *TypeResolver {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordResolver,
		FragTkKwdRsv,
		"keyword",
	)
	if pr.err(err) {
		return nil
	}

	// Read type ID
	fType, err := readWord(
		lex,
		"resolver type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Create a new resolver type instance
	newResolver := &TypeResolver{
		terminalType: terminalType{
			TypeName: fType.src,
		},
	}

	// Parse properties
	fProps, props := pr.parseRsvProps(lex, newResolver)
	if fProps == nil {
		return nil
	}
	newResolver.Properties = props

	newResolver.Src = NewConstruct(lex, FragDeclRsv,
		fDeclKeyword,
		fType,
		fProps,
	)

	// Define the type
	pr.defineType(newResolver)

	return newResolver
}
