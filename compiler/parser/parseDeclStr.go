package parser

func (pr *Parser) parseDeclStr(lex *Lexer) *TypeStruct {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordStruct,
		FragTkKwdStr,
		"keyword",
	)
	if pr.err(err) {
		return nil
	}

	// Read type ID
	fType, err := readWord(
		lex,
		"struct type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Create a new struct type instance
	newStruct := &TypeStruct{
		terminalType: terminalType{
			TypeName: fType.src,
		},
	}

	// Parse fields
	fFields, fields := pr.parseStrFields(lex, newStruct)
	if fFields == nil {
		return nil
	}
	newStruct.Fields = fields

	newStruct.Src = NewConstruct(lex, FragDeclStr,
		fDeclKeyword,
		fType,
		fFields,
	)

	// Define the type
	pr.defineType(newStruct)

	return newStruct
}
