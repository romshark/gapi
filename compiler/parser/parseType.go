package parser

import (
	"fmt"

	parser "github.com/romshark/llparser"
	"github.com/romshark/llparser/misc"
)

func (pr *Parser) parseType(
	frag parser.Fragment,
	onDetermined func(Type),
) bool {
	var tp Type
	var previousTp Type

	appendTp := func(t Type) {
		if tp == nil {
			tp = t
		} else {
			switch v := previousTp.(type) {
			case *TypeOptional:
				v.StoreType = t
			case *TypeList:
				v.StoreType = t
			}
		}
		previousTp = t
	}

	for _, elem := range frag.Elements() {
		switch elem.Kind() {
		case FragTkSymList:
			// List container type
			appendTp(&TypeList{})

		case FragTkSymOpt:
			// Optional container type
			// Ensure the previous type in the chain was not also an optional
			if previousTp != nil {
				if _, tailIsOpt := previousTp.(*TypeOptional); tailIsOpt {
					// Illegal optionals chain detected
					// (Optional type of optional types)
					pr.err(&pErr{
						at:   elem.Begin(),
						code: ErrTypeOptChain,
						message: fmt.Sprintf(
							"illegal chain of optionals " +
								"(optional type of optional types)",
						),
					})
					return false
				}
			}
			appendTp(&TypeOptional{})

		case misc.FrWord:
			// Terminal type reached
			pr.deferJob(func() {
				// Make sure the terminal type is defined
				terminalType := pr.findTypeByDesignation(elem.Src())
				if terminalType == nil {
					pr.err(&pErr{
						at:   elem.Begin(),
						code: ErrTypeUndef,
						message: fmt.Sprintf(
							"terminal type %s is undefined",
							string(elem.Src()),
						),
					})
					return
				}
				appendTp(terminalType)

				// Reference the terminal type in the type chain
				for t := tp; t != nil; {
					if v, isOpt := t.(*TypeOptional); isOpt {
						v.Terminal = previousTp
						t = v.StoreType
						continue
					}
					if v, isList := t.(*TypeList); isList {
						v.Terminal = previousTp
						t = v.StoreType
						continue
					}
					break
				}

				_, terminalIsNone := previousTp.(TypeStdNone)
				if _, isNone := tp.(TypeStdNone); !isNone && terminalIsNone {
					pr.err(&pErr{
						at:      frag.Begin(),
						code:    ErrSyntax,
						message: "illegal None-type",
					})
					return
				}

				if tp.TerminalType() != nil {
					tp = pr.onAnonymousType(tp)
				}
				onDetermined(tp)

			})
		default:
			panic(fmt.Errorf("unexpected token in type designation: %s", elem))
		}
	}

	return true
}
