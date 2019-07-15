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

	if err := verifyTypeName(newUnionTypeName); err != nil {
		return errors.Errorf(
			"invalid union type identifier %d:%d: %s",
			current.begin,
			current.end,
			err,
		)
	}

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

		if err := verifyTypeName(referencedTypeName); err != nil {
			return errors.Errorf(
				"invalid union option-type identifier at %d:%d: %s",
				current.begin,
				current.end,
				err,
			)
		}

		// Check for duplicate values
		if _, isDefined := newType.Types[referencedTypeName]; isDefined {
			return errors.Errorf(
				"multiple references to the same type (%s) "+
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
			"union %s requires at least two types at %d:%d",
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
					"undefined type %s referenced "+
						"in union type %s at %d:%d ",
					name,
					newUnionTypeName,
					node.begin,
					node.end,
				)
			}
			if name == newUnionTypeName {
				return errors.Errorf(
					"union type %s references itself at %d:%d",
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
