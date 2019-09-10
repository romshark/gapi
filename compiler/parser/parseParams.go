package parser

import (
	"fmt"

	parser "github.com/romshark/llparser"
)

// parseParams parses a block-list of parameters
func (pr *Parser) parseParams(
	frag parser.Fragment,
	target GraphNode,
) {
	byName := map[string]parser.Fragment{}
	elems := frag.Elements()
	paramsListFrag, _ := findElement(elems, FragParams, 0)

	parseParam := func(frag parser.Fragment) {
		param := &Parameter{
			Src:    frag,
			Name:   string(frag.Elements()[0].Src()),
			Target: target,
		}

		// Check for redeclarations
		if defined, isDefined := byName[param.Name]; isDefined {
			pr.err(&pErr{
				at:   frag.Begin(),
				code: ErrParamRedecl,
				message: fmt.Sprintf(
					"Redeclaration of parameter %s "+
						"(previously declared at %s)",
					param.Name,
					defined.Begin(),
				),
			})
			return
		}

		pr.deferJob(func() {
			pr.parseType(
				frag.Elements()[2],
				func(t Type) { param.Type = t },
			)
		})

		byName[param.Name] = frag
		pr.onParameter(param)

		switch v := target.(type) {
		case *Query:
			v.Parameters = append(v.Parameters, param)
		case *Mutation:
			v.Parameters = append(v.Parameters, param)
		case *ResolverProperty:
			v.Parameters = append(v.Parameters, param)
		}
	}

	var parseGroup func(parser.Fragment)
	parseGroup = func(frag parser.Fragment) {
		for _, el := range frag.Elements() {
			switch el.Kind() {
			case FragParam:
				parseParam(el)
			case FragParams:
				parseGroup(el)
			}
		}
	}

	parseGroup(paramsListFrag)
}
