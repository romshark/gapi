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
		require.IsType(t, &compiler.TypeAlias{}, t1)
		require.Equal(
			t,
			compiler.TypeStdString{},
			t1.(*compiler.TypeAlias).AliasedType,
		)

		require.Equal(t, "T2", t2.Name())
		require.Equal(t, compiler.TypeCategoryAlias, t2.Category())
		require.IsType(t, &compiler.TypeAlias{}, t2)
		require.Equal(
			t,
			compiler.TypeStdUint32{},
			t2.(*compiler.TypeAlias).AliasedType,
		)

		require.Equal(t, "T3", t3.Name())
		require.Equal(t, compiler.TypeCategoryAlias, t3.Category())
		require.IsType(t, &compiler.TypeAlias{}, t3)
		require.Equal(t, t1, t3.(*compiler.TypeAlias).AliasedType)
	})
}

// TestDeclEnumTypes tests enum type declaration
func TestDeclEnumTypes(t *testing.T) {
	src := `schema test
	
	enum E1 {
		oneVal
	}
	enum E2 {
		foo
		bar
	}
	enum E3 {
		foo1
		bar2
		baz3
	}
	`

	test(t, src, func(ast AST) {
		require.Len(t, ast.QueryEndpoints, 0)
		require.Len(t, ast.Mutations, 0)

		expected := map[string][]string{
			"E1": []string{"oneVal"},
			"E2": []string{"foo", "bar"},
			"E3": []string{"foo1", "bar2", "baz3"},
		}

		require.Len(t, ast.Types, len(expected))
		for name, vals := range expected {
			require.Contains(t, ast.Types, name)
			tp := ast.Types[name]
			require.Equal(t, name, tp.Name())
			require.Equal(t, compiler.TypeCategoryEnum, tp.Category())
			require.IsType(t, &compiler.TypeEnum{}, tp)
			tpe := tp.(*compiler.TypeEnum)
			for _, val := range vals {
				require.Contains(t, tpe.Values, val)
			}
		}
	})
}

// TestDeclUnionTypes tests union type declaration
func TestDeclUnionTypes(t *testing.T) {
	src := `schema test
	
	union U1 {
		String
		Uint32
	}
	union U2 {
		Uint32
		Float64
		String
	}
	union U3 {
		String
		Float64
		Int32
		Int64
	}
	`

	test(t, src, func(ast AST) {
		require.Len(t, ast.QueryEndpoints, 0)
		require.Len(t, ast.Mutations, 0)

		expected := map[string][]compiler.Type{
			"U1": []compiler.Type{
				compiler.TypeStdString{},
				compiler.TypeStdUint32{},
			},
			"U2": []compiler.Type{
				compiler.TypeStdUint32{},
				compiler.TypeStdFloat64{},
				compiler.TypeStdString{},
			},
			"U3": []compiler.Type{
				compiler.TypeStdString{},
				compiler.TypeStdFloat64{},
				compiler.TypeStdInt32{},
				compiler.TypeStdInt64{},
			},
		}

		require.Len(t, ast.Types, len(expected))
		for name, expectedReferencedTypes := range expected {
			require.Contains(t, ast.Types, name)
			tp := ast.Types[name]
			require.Equal(t, name, tp.Name())
			require.Equal(t, compiler.TypeCategoryUnion, tp.Category())
			require.IsType(t, &compiler.TypeUnion{}, tp)
			tpe := tp.(*compiler.TypeUnion)
			for _, referencedType := range expectedReferencedTypes {
				require.Contains(t, tpe.Types, referencedType.Name())
			}
		}
	})
}

func TestDeclUnionTypesErr(t *testing.T) {
	t.Run("OneTypeUnion", func(t *testing.T) {
		testErr(t, `schema test
		union U1 {
			String
		}
		`)
	})

	t.Run("MultiReferencedType", func(t *testing.T) {
		testErr(t, `schema test
		union U1 {
			String
			String
		}
		`)
	})

	t.Run("UndefinedType", func(t *testing.T) {
		testErr(t, `schema test
		union U1 {
			Undefined
		}
		`)
	})
}
