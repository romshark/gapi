package parser

func (pr *Parser) parseDeclEnm(lex *Lexer) *TypeEnum {
	// Read keyword
	fDeclKeyword, err := readWordExact(
		lex,
		KeywordEnum,
		FragTkKwdEnm,
		"keyword",
	)
	if pr.err(err) {
		return nil
	}

	// Read type identifier
	fTypeID, err := readWord(
		lex,
		"enum type identifier",
		FragTkIdnType,
		capitalizedCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Instantiate type
	newType := &TypeEnum{
		terminalType: terminalType{
			TypeName: fTypeID.src,
		},
	}

	// Parse enum values
	fVals, values := pr.parseEnmVals(lex, newType)
	if fVals == nil {
		return nil
	}
	newType.Values = values

	newType.Src = NewConstruct(lex, FragDeclEnm,
		fDeclKeyword,
		fTypeID,
		fVals,
	)

	// Define the type
	pr.defineType(newType)

	return newType
}
