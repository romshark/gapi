package parser

import parser "github.com/romshark/llparser"

// onDeclTypeEnum is executed when an enum type declaration is matched
func (pr *Parser) onDeclTypeEnum(frag parser.Fragment) error {
	// Instantiate type
	newType := &TypeEnum{
		terminalType: terminalType{
			Src:  frag,
			Name: string(frag.Elements()[0].Src()),
		},
	}

	offset := uint(0)
	var valueEl parser.Fragment
	for {
		// Traverse all values
		valueEl, offset = findElement(
			frag.Elements(),
			FragTkIdnEnumVal,
			offset,
		)
		if valueEl == nil {
			break
		}
		offset++

		fieldItems := valueEl.Elements()
		value := &EnumValue{
			Src:  valueEl,
			Enum: newType,
			Name: string(fieldItems[0].Src()),
		}
		newType.Values = append(newType.Values, value)
	}

	// Define the type
	pr.onTypeDecl(newType)
	return nil
}
