package parser

import "fmt"

// parseParams parses a parameter
func (pr *Parser) parseParam(lex *Lexer, target GraphNode) *Parameter {
	// Read name
	fName, err := readWord(
		lex,
		"parameter identifier",
		FragTkIdnParam,
		lowerCamelCase,
	)
	if pr.err(err) {
		return nil
	}

	// Create a new parameter instance
	newParam := &Parameter{
		Target: target,
		Name:   fName.src,
	}

	// Read type and set it when it's determined
	fType := pr.parseTypeDesig(lex, func(t Type) {
		// Make sure the type of the parameter is pure
		if !t.IsPure() {
			pr.err(&pErr{
				at:      fName.begin,
				code:    ErrParamImpure,
				message: fmt.Sprintf("parameter of impure type %s", t),
			})
		}

		newParam.Type = t
	})
	if fType == nil {
		return nil
	}

	newParam.Src = NewConstruct(lex, FragParam,
		fName,
		fType,
	)

	// Define the graph node
	pr.onParameter(newParam)

	return newParam
}
