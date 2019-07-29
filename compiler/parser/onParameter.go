package parser

// onParameter defines a new parameter assigning it a unique identifier
func (pr *Parser) onParameter(newParam *Parameter) {
	// Register a new parameter
	pr.lastIssuedParamID += ParamID(1)
	newParam.ID = pr.lastIssuedParamID

	pr.paramByID[newParam.ID] = newParam
}
