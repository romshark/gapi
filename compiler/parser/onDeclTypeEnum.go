package parser

import (
	"fmt"

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

	values := map[string]parser.Fragment{}
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

		// Check for redefinitions
		if defined, isDefined := values[value.Name]; isDefined {
			pr.err(&pErr{
				at:   valueEl.Begin(),
				code: ErrEnumValRedecl,
				message: fmt.Sprintf(
					"Redeclaration of enum value %s "+
						"(previously declared at %s)",
					value.Name,
					defined.Begin(),
				),
			})
			return nil
		}
		values[value.Name] = valueEl

		newType.Values = append(newType.Values, value)
	}

	// Check for values
	if len(newType.Values) < 1 {
		pr.err(&pErr{
			at:   frag.Begin(),
			code: ErrEnumNoVal,
			message: fmt.Sprintf(
				"Enum %s is missing values",
				newType.Name,
			),
		})
	}

	// Define the type
	pr.onTypeDecl(newType)
	return nil
}
