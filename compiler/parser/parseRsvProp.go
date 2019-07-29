package parser

// parseRsvProp parses a resolver property
func (pr *Parser) parseRsvProp(
	lex *Lexer,
	resolver *TypeResolver,
) *ResolverProperty {
	// Read property name
	fName, err := readWord(
		lex,
		"property identifier",
		FragTkIdnProp,
		lowerCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	newProp := &ResolverProperty{
		Resolver: resolver,
		Name:     fName.src,
	}

	// Parse parameters
	fParams, params, parsed := pr.parseOptParams(lex, newProp)
	if !parsed {
		return nil
	}
	newProp.Parameters = params

	// Read type and set it when it's determined
	fType := pr.parseTypeDesig(lex, func(t Type) { newProp.Type = t })
	if fType == nil {
		return nil
	}

	if fParams != nil {
		newProp.Src = NewConstruct(lex, FragRsvProp,
			fName,
			fParams,
			fType,
		)
	} else {
		newProp.Src = NewConstruct(lex, FragRsvProp,
			fName,
			fType,
		)
	}

	// Define the graph node
	pr.defineGraphNode(newProp)

	return newProp
}
