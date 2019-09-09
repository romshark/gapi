package parser

import parser "github.com/romshark/llparser"

// onDeclTypeAlias is executed when an alias type declaration is matched
func (pr *Parser) onDeclTypeAlias(frag parser.Fragment) error {
	elems := frag.Elements()

	// Instantiate type
	newType := &TypeAlias{
		terminalType: terminalType{
			Src:  frag,
			Name: string(elems[0].Src()),
		},
	}

	typeDesigFrag, _ := findElement(elems, FragType, 2)

	if !pr.parseType(
		typeDesigFrag,
		func(t Type) { newType.AliasedType = t },
	) {
		return nil
	}

	// Define the type
	pr.onTypeDecl(newType)
	return nil
}
