package parser

import (
	parser "github.com/romshark/llparser"
)

// onDeclTypeEnum is executed when an enum type declaration is matched
func (pr *Parser) onDeclTypeEnum(frag parser.Fragment) error {
	// Instantiate type
	newType := &TypeEnum{
		terminalType: terminalType{
			Src:  frag,
			Name: string(frag.Elements()[0].Src()),
		},
	}

	// Define the type
	pr.onTypeDecl(newType)
	return nil
}
