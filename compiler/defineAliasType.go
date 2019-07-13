package compiler

import "github.com/pkg/errors"

func (c *Compiler) defineAliasType(
	source string,
	ast *AST,
	node *node32,
) error {
	typeName := getSrc(source, node.up.next.next)
	aliasedTypeName := getSrc(source, node.up.next.next.next.next.next.next)

	// Try to define the type
	if err := ast.defineType(&TypeAlias{
		typeBaseInfo: typeBaseInfo{
			src:  src(node),
			name: typeName,
		},
	}); err != nil {
		return err
	}

	c.deferJob(func() error {
		// Ensure the aliased type exists after all types have been defined
		if !ast.isTypeNameDefined(aliasedTypeName) {
			return errors.Errorf(
				"undefined type %s aliased by %s at %d:%d",
				aliasedTypeName,
				typeName,
				node.begin,
				node.end,
			)
		}
		return nil
	})

	return nil
}
