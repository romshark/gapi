package compiler

import (
	"fmt"
)

// parseType tries to parse a type node.
// Should only be used in deferred functons after all types have been defined
// because it checks the type for its definition and validity
func (c *Compiler) parseType(node *node32) (Type, Error) {
	var tp Type
	var tailTp Type
	setStoredType := func(c, t Type) {
		switch v := c.(type) {
		case *TypeOptional:
			v.StoreType = t
		case *TypeList:
			v.StoreType = t
		}
	}
	appendTp := func(t Type) {
		if tp == nil {
			tp = t
			tailTp = t
		} else {
			setStoredType(tailTp, t)
			tailTp = t
		}
	}

LOOP:
	for current := node.up; current != nil; current = current.next {
		switch current.pegRule {
		case ruleONU:
			// Optional container
			// Ensure the previous type in the chain was not also an optional
			if tailTp != nil {
				if _, tailIsOpt := tailTp.(*TypeOptional); tailIsOpt {
					// Illegal optionals chain detected
					// (Optional type of optional types)
					return nil, cErr{
						ErrTypeOptChain,
						fmt.Sprintf(
							"illegal chain of optionals "+
								"(optional type of optional types) at %d:%d",
							current.begin,
							current.end,
						),
					}
				}
			}
			appendTp(&TypeOptional{})

		case ruleOLI:
			// List container
			appendTp(&TypeList{})

		case ruleWrd:
			// Terminal type
			terminalTypeName := c.getSrc(current)
			terminalType := c.ast.FindTypeByName("", terminalTypeName)
			if terminalType == nil {
				return nil, cErr{
					ErrTypeUndef,
					fmt.Sprintf(
						"terminal type %s is undefined "+
							"in type declaration at %d:%d",
						terminalTypeName,
						node.begin,
						node.end,
					),
				}
			}
			appendTp(terminalType)
			break LOOP
		}
	}

	// Reference the terminal type in type chains
	for t := tp; t != nil; {
		if v, isOpt := t.(*TypeOptional); isOpt {
			v.Terminal = tailTp
			t = v.StoreType
			continue
		}
		if v, isList := t.(*TypeList); isList {
			v.Terminal = tailTp
			t = v.StoreType
			continue
		}
		break
	}

	return tp, nil
}
