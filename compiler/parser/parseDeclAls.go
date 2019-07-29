package parser

func (pr *Parser) parseDeclAls(lex *Lexer) *TypeAlias {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordAlias,
		FragTkKwdAls,
		"keyword",
	)
	if pr.err(err) {
		return nil
	}

	// Read type ID
	fTypeID, err := readWord(
		lex,
		"alias type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Read '='
	fSymEq, err := readToken(lex, FragTkSymEq, "equals sign")
	if pr.err(err) {
		return nil
	}

	// Instantiate type
	newType := &TypeAlias{
		terminalType: terminalType{
			Name: fTypeID.src,
		},
	}

	// Read aliased type
	fType := pr.parseTypeDesig(lex, func(t Type) { newType.AliasedType = t })
	if fType == nil {
		return nil
	}

	newType.Src = NewConstruct(lex, FragDeclAls,
		fDeclKeyword,
		fTypeID,
		fSymEq,
		fType,
	)

	// Define the type
	pr.onTypeDecl(newType)

	return newType
}
