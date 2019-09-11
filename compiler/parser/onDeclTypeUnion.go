package parser

import (
	"fmt"

	parser "github.com/romshark/llparser"
)

// onDeclTypeUnion is executed when a union type declaration is matched
func (pr *Parser) onDeclTypeUnion(frag parser.Fragment) error {
	// Instantiate type
	newType := &TypeUnion{
		terminalType: terminalType{
			Src:  frag,
			Name: string(frag.Elements()[0].Src()),
		},
	}

	types := map[string]parser.Fragment{}
	offset := uint(0)
	var typeEl parser.Fragment
	for {
		// Traverse all types
		typeEl, offset = findElement(
			frag.Elements(),
			FragType,
			offset,
		)
		if typeEl == nil {
			break
		}
		offset++

		typeDesignation := string(typeEl.Src())

		// Check for redefinitions
		if defined, isDefined := types[typeDesignation]; isDefined {
			pr.err(&pErr{
				at:   typeEl.Begin(),
				code: ErrEnumValRedecl,
				message: fmt.Sprintf(
					"Redundant union option type %s "+
						"(previously declared at %s)",
					typeDesignation,
					defined.Begin(),
				),
			})
			return nil
		}
		types[typeDesignation] = typeEl

		// Ensure the union doesn't reference itself as an option
		if typeDesignation == newType.Name {
			pr.err(&pErr{
				at:   frag.Begin(),
				code: ErrUnionRecurs,
				message: fmt.Sprintf(
					"Union type %s references itself "+
						"as one of its options",
					newType,
				),
			})
			return nil
		}

		typeElement := typeEl
		pr.deferJob(func() {
			pr.parseType(typeElement, func(tp Type) {
				// Ensure the union type doesn't include None as an option
				if _, isNone := tp.(TypeStdNone); isNone {
					pr.err(&pErr{
						at:   typeElement.Begin(),
						code: ErrUnionIncludesNone,
						message: fmt.Sprintf(
							"Union type %s includes the None primitive",
							newType,
						),
					})
					return
				}

				newType.Types = append(newType.Types, tp)
			})
		})
	}

	// Make sure there's at least 2 options
	if len(types) < 2 {
		pr.err(&pErr{
			at:   frag.Begin(),
			code: ErrUnionMissingOpts,
			message: fmt.Sprintf(
				"Union %s is missing type options (%d out of at least 2)",
				newType.Name,
				len(types),
			),
		})
		return nil
	}

	// Define the type
	pr.onTypeDecl(newType)
	return nil
}
