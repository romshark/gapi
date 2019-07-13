package compiler

import "github.com/pkg/errors"

func (c *Compiler) defineAliasType(
	source string,
	ast *AST,
	node *node32,
) error {
	typeName := getSrc(source, node.up.next.next)
	aliasedTypeName := getSrc(source, node.up.next.next.next.next.next.next)

	newType := &TypeAlias{
		typeBaseInfo: typeBaseInfo{
			src:  src(node),
			name: typeName,
		},
	}

	// Try to define the type
	if err := ast.defineType(newType); err != nil {
		return err
	}

	c.deferJob(func() error {
		// Ensure the aliased type exists after all types have been defined
		aliasedType := ast.typeByName(aliasedTypeName)
		if aliasedType == nil {
			return errors.Errorf(
				"undefined type %s aliased by %s at %d:%d",
				aliasedTypeName,
				typeName,
				node.begin,
				node.end,
			)
		}

		// Reference the aliased type
		newType.AliasedType = aliasedType

		return nil
	})

	return nil
}
