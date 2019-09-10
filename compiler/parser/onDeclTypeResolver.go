package parser

import parser "github.com/romshark/llparser"

// onDeclTypeResolver is executed when a resolver type declaration is matched
func (pr *Parser) onDeclTypeResolver(frag parser.Fragment) error {
	/* elems := frag.Elements()

	// Instantiate type
	newType := &TypeResolver{
		terminalType: terminalType{
			Src:  frag,
			Name: string(elems[0].Src()),
		},
	}

	offset := uint(0)
	var el parser.Fragment
	for {
		// Traverse all fields
		el, offset = findElement(frag.Elements(), FragStrField, offset)
		if el == nil {
			break
		}
		offset++

		fieldItems := el.Elements()
		field := &StructField{
			Src:    el,
			Struct: newType,
			Name:   string(fieldItems[0].Src()),
		}
		newType.Fields = append(newType.Fields, field)
		if !pr.onGraphNode(field) {
			continue
		}

		// Defer parsing and setting the type of the field
		pr.deferJob(func() {
			pr.parseType(
				fieldItems[2],
				func(t Type) { field.Type = t },
			)
		})
	}

	// Define the type
	pr.onTypeDecl(newType) */
	return nil
}
