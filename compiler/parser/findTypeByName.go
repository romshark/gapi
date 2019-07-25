package parser

func (pr *Parser) findTypeByName(name string) Type {
	// Search in all categories including primitives
	if t := stdTypeByName(name); t != nil {
		return t
	}
	return pr.typeByName[name]
}
