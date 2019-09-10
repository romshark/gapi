package parser

import parser "github.com/romshark/llparser"

// onDeclMutation is executed when a mutation endpoint declaration is matched
func (pr *Parser) onDeclMutation(frag parser.Fragment) error {
	elems := frag.Elements()

	newEndpoint := &Mutation{
		Src:  frag,
		Name: string(elems[0].Src()),
	}

	paramsElem, _ := findElement(elems, FragParams, 2)
	if paramsElem != nil {
		pr.parseParams(paramsElem, newEndpoint)
	}

	if !pr.parseType(
		elems[len(elems)-1],
		func(t Type) { newEndpoint.Type = t },
	) {
		return nil
	}

	// Define the endpoint
	pr.onGraphNode(newEndpoint)
	return nil
}
