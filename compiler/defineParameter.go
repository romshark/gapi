package compiler

// defineParameter returns a new unique parameter identifier
func (c *Compiler) defineParameter(newParam *Parameter) ParamID {
	// Update the last issued ID
	c.lastIssuedParamID += ParamID(1)
	newID := c.lastIssuedParamID

	c.paramByID[newID] = newParam

	return newID
}
