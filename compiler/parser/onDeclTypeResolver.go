package parser

import (
	"fmt"

	parser "github.com/romshark/llparser"
)

// onDeclTypeResolver is executed when a resolver type declaration is matched
func (pr *Parser) onDeclTypeResolver(frag parser.Fragment) error {
	elems := frag.Elements()
	byName := map[string]parser.Fragment{}

	// Instantiate type
	newType := &TypeResolver{
		terminalType: terminalType{
			Src:  frag,
			Name: string(elems[0].Src()),
		},
	}

	offset := uint(0)
	var propEl parser.Fragment
	for {
		// Traverse all properties
		propEl, offset = findElement(frag.Elements(), FragRsvProp, offset)
		if propEl == nil {
			break
		}
		offset++

		propItems := propEl.Elements()
		prop := &ResolverProperty{
			Src:      propEl,
			Resolver: newType,
			Name:     string(propItems[0].Src()),
		}
		newType.Properties = append(newType.Properties, prop)
		if !pr.onGraphNode(prop) {
			continue
		}

		// Check for redeclarations
		if defined, isDefined := byName[prop.Name]; isDefined {
			pr.err(&pErr{
				at:   propEl.Begin(),
				code: ErrResolverPropRedecl,
				message: fmt.Sprintf(
					"Redeclaration of resolver property %s "+
						"(previously declared at %s)",
					prop.Name,
					defined.Begin(),
				),
			})
			return nil
		}
		byName[prop.Name] = propEl

		// Evaluate parameters if any
		params, _ := findElement(propEl.Elements(), FragParams, 0)
		if params != nil {
			pr.parseParams(params, prop)
		}

		// Defer parsing and setting the type of the prop
		pr.deferJob(func() {
			pr.parseType(
				propItems[len(propItems)-1],
				func(t Type) { prop.Type = t },
			)
		})
	}

	// Define the type
	pr.onTypeDecl(newType)
	return nil
}
