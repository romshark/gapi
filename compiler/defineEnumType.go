package compiler

import (
	"log"

	"github.com/pkg/errors"
)

func (c *Compiler) defineEnumType(
	source string,
	ast *AST,
	node *node32,
) error {
	current := node.up.next.next
	typeName := getSrc(source, node.up.next.next)

	newType := &TypeEnum{
		typeBaseInfo: typeBaseInfo{
			src:  src(node),
			name: typeName,
		},
		Values: make(map[string]EnumValue),
	}

	// Parse values
	current = current.next.next.up.next.next
	for {
		valueName := source[current.begin:current.end]

		// Check for duplicate values
		if defined, isDefined := newType.Values[valueName]; isDefined {
			return errors.Errorf(
				"Redeclaration of enum value %s at %d:%d "+
					"(previously declared at %d:%d)",
				valueName,
				current.begin,
				current.end,
				defined.Begin,
				defined.End,
			)
		}

		// Add enum value
		newType.Values[valueName] = EnumValue{
			Src: Src{
				Begin: current.begin,
				End:   current.end,
			},
			Name: valueName,
		}

		next := current.next.next
		if next == nil || next.pegRule == ruleBLKE {
			break
		}
		current = next
	}
	log.Print(newType.Values)

	// Try to define the type
	return ast.defineType(newType)
}
