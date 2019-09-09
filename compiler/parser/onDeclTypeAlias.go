package parser

import (
	parser "github.com/romshark/llparser"
)

// onDeclTypeAlias is executed when an alias type declaration is matched
func (pr *Parser) onDeclTypeAlias(frag parser.Fragment) error {
	// Instantiate type
	newType := &TypeAlias{
		terminalType: terminalType{
			Src:  frag,
			Name: string(frag.Elements()[0].Src()),
		},
	}

	if !pr.parseType(
		findElement(frag, FragType, 2),
		func(t Type) { newType.AliasedType = t },
	) {
		return nil
	}

	// Define the type
	pr.onTypeDecl(newType)
	return nil
}
