package compiler

import (
	"github.com/pkg/errors"
)

func (c *Compiler) defineEnumType(
	source string,
	ast *AST,
	node *node32,
) error {
	current := node.up.next.next
	newEnumTypeName := getSrc(source, node.up.next.next)

	if err := verifyTypeName(newEnumTypeName); err != nil {
		return errors.Errorf(
			"invalid enum type identifier %d:%d: %s",
			current.begin,
			current.end,
			err,
		)
	}

	newType := &TypeEnum{
		typeBaseInfo: typeBaseInfo{
			src:  src(node),
			name: newEnumTypeName,
		},
		Values: make(map[string]EnumValue),
	}

	// Parse values
	current = current.next.next.up.next.next
	for {
		valueName := source[current.begin:current.end]

		if err := verifyEnumValue(valueName); err != nil {
			return errors.Errorf(
				"invalid enum value identifier at %d:%d: %s",
				current.begin,
				current.end,
				err,
			)
		}

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

	// Try to define the type
	return ast.defineType(newType)
}
