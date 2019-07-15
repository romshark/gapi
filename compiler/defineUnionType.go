package compiler

import (
	"github.com/pkg/errors"
)

func (c *Compiler) defineUnionType(
	source string,
	ast *AST,
	node *node32,
) error {
	current := node.up.next.next
	newUnionTypeName := getSrc(source, current)

	newType := &TypeUnion{
		typeBaseInfo: typeBaseInfo{
			src:  src(node),
			name: newUnionTypeName,
		},
		Types: make(map[string]Type),
	}

	// Parse types
	current = current.next.next.up.next.next
	for {
		referencedTypeName := source[current.begin:current.end]

		// Check for duplicate values
		if _, isDefined := newType.Types[referencedTypeName]; isDefined {
			return errors.Errorf(
				"Multiple references to the same type (%s) "+
					"in union type %s at %d:%d ",
				referencedTypeName,
				newUnionTypeName,
				node.begin,
				node.end,
			)
		}

		// Mark type for deferred checking
		newType.Types[referencedTypeName] = nil

		next := current.next.next
		if next == nil || next.pegRule == ruleBLKE {
			break
		}
		current = next
	}

	if len(newType.Types) < 2 {
		return errors.Errorf(
			"Union %s requires at least two types at %d:%d",
			newUnionTypeName,
			node.begin,
			node.end,
		)
	}

	c.deferJob(func() error {
		// Ensure all referenced types are defined
		for name := range newType.Types {
			reg := ast.typeByName(name)
			if reg == nil {
				return errors.Errorf(
					"Undefined type %s referenced "+
						"in union type %s at %d:%d ",
					name,
					newUnionTypeName,
					node.begin,
					node.end,
				)
			}
			newType.Types[name] = reg
		}

		return nil
	})

	// Try to define the type
	return ast.defineType(newType)
}
