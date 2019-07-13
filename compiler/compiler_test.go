package compiler_test

import (
	"testing"

	"github.com/romshark/gapi/compiler"
	"github.com/stretchr/testify/require"
)

type AST *compiler.AST

func test(
	t *testing.T,
	source string,
	astInspector func(AST),
) {
	// Initialize compiler
	compiler, err := compiler.NewCompiler()
	require.NoError(t, err)

	// Compile
	ast, err := compiler.Compile(source)
	require.NoError(t, err)
	require.NotNil(t, ast)

	astInspector(ast)
}

func testErr(t *testing.T, source string) {
	// Initialize compiler
	compiler, err := compiler.NewCompiler()
	require.NoError(t, err)

	// Compile
	ast, err := compiler.Compile(source)
	require.Error(t, err)
	require.Nil(t, ast)
}

// TestDeclAliasTypes tests alias type declaration
func TestDeclAliasTypes(t *testing.T) {
	src := `schema test
	
	alias T1 = String
	alias T2 = Uint32
	alias T3 = T1
	`

	test(t, src, func(ast AST) {
		require.Len(t, ast.Types, 3)
		require.Len(t, ast.QueryEndpoints, 0)
		require.Len(t, ast.Mutations, 0)
		require.Contains(t, ast.Types, "T1")
		require.Contains(t, ast.Types, "T2")
		require.Contains(t, ast.Types, "T3")

		t1 := ast.Types["T1"]
		t2 := ast.Types["T2"]
		t3 := ast.Types["T3"]

		require.Equal(t, "T1", t1.Name())
		require.Equal(t, compiler.TypeCategoryAlias, t1.Category())

		require.Equal(t, "T2", t2.Name())
		require.Equal(t, compiler.TypeCategoryAlias, t2.Category())

		require.Equal(t, "T3", t3.Name())
		require.Equal(t, compiler.TypeCategoryAlias, t3.Category())
	})
}
