package parser

func (pr *Parser) findTypeByDesignation(designation []rune) Type {
	str := string(designation)

	// Search in all categories including primitives
	if t := stdTypeByName(str); t != nil {
		return t
	}

	// Search in user-defined types
	return pr.typeByName[str]
}
