package compiler

func (c *Compiler) findTypeByName(name string) Type {
	// Search in all categories including primitives
	if t := stdTypeByName(name); t != nil {
		return t
	}
	return c.typeByName[name]
}
