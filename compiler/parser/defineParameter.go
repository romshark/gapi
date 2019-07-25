package parser

// defineParameter returns validates the name, checks for redeclaration
// and defines a new parameter assigning it a unique identifier
func (pr *Parser) defineParameter(newParam *Parameter) {
	// Register a new parameter
	pr.lastIssuedParamID += ParamID(1)
	newParam.ID = pr.lastIssuedParamID

	pr.paramByID[newParam.ID] = newParam
}
