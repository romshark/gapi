package parser

func (pr *Parser) findTypeByDesignation(designation string) Type {
	// Search in all categories including primitives
	if t := stdTypeByName(designation); t != nil {
		return t
	}

	// Search in alias types
	if alias, isAlias := pr.aliasByName[designation]; isAlias {
		return alias
	}

	// Search in user-defined types
	return pr.typeByName[designation]
}
