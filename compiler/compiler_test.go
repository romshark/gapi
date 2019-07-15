package compiler_test

import (
	"testing"

	"github.com/romshark/gapi/compiler"
	"github.com/stretchr/testify/require"
)

type AST = *compiler.AST
type ErrCode = compiler.ErrCode

func test(
	t *testing.T,
	source string,
	astInspector func(AST),
) {
	// Initialize compiler
	compiler, err := compiler.NewCompiler(source)
	require.NoError(t, err)

	// Compile
	require.NoError(t, compiler.Compile())
	ast := compiler.AST()
	require.NotNil(t, ast)
	astInspector(ast)
}

type ErrCase struct {
	// Src defines the source-code
	Src string

	// Errs defines all expected compiler error codes
	Errs []ErrCode
}

func testErrs(t *testing.T, cases map[string]ErrCase) {
	for tst, errCase := range cases {
		if len(errCase.Errs) < 1 {
			panic("missing expected errors in error test case")
		}
		t.Run(tst, func(t *testing.T) {
			// Initialize compiler
			compiler, err := compiler.NewCompiler(errCase.Src)
			require.NoError(t, err)

			// Compile
			require.Error(t, compiler.Compile())
			require.Nil(t, compiler.AST())

			type Err struct {
				Code ErrCode
				Name string
			}

			actualErrs := compiler.Errors()
			actualCodes := make([]Err, len(actualErrs))
			for i, actulErr := range actualErrs {
				c := actulErr.Code()
				actualCodes[i] = Err{
					Code: c,
					Name: c.String(),
				}
			}

			expectedCodes := make([]Err, len(errCase.Errs))
			for i, expectedCode := range errCase.Errs {
				expectedCodes[i] = Err{
					Code: expectedCode,
					Name: expectedCode.String(),
				}
			}

			require.Equal(t, expectedCodes, actualCodes)
		})
	}
}

// TestDeclSchemaErrs tests schema declaration errors
func TestDeclSchemaErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalName": ErrCase{
			Src: `schema _illegalName
			enum E { e }`,
			Errs: []ErrCode{compiler.ErrSchemaIllegalIdent},
		},
		"IllegalName2": ErrCase{
			Src: `schema illegal_Name
			enum E { e }`,
			Errs: []ErrCode{compiler.ErrSchemaIllegalIdent},
		},
		"IllegalName3": ErrCase{
			Src: `schema IllegalName
			enum E { e }`,
			Errs: []ErrCode{compiler.ErrSchemaIllegalIdent},
		},
	})
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

// TestDeclAliasTypeErrs tests alias type declaration errors
func TestDeclAliasTypeErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalTypeName": ErrCase{
			Src: `schema test
			alias illegalName = String`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			alias Illegal_Name = String`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			alias _IllegalName = String`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalAliasedTypeName": ErrCase{
			Src: `schema test
			alias A = illegalName`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalAliasedTypeName2": ErrCase{
			Src: `schema test
			alias A = Illegal_Name`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalAliasedTypeName3": ErrCase{
			Src: `schema test
			alias A = _IllegalName`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"UndefinedAliasedType": ErrCase{
			Src: `schema test
			alias A = Undefined`,
			Errs: []ErrCode{compiler.ErrTypeUndef},
		},
		"SelfReference": ErrCase{
			Src: `schema test
			alias A = A`,
			Errs: []ErrCode{compiler.ErrAliasRecur},
		},
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

// TestDeclEnumTypeErrs tests enum type declaration errors
func TestDeclEnumTypeErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalTypeName": ErrCase{
			Src: `schema test
			enum illegalName {
				foo
				bar
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			enum _IllegalName {
				foo
				bar
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			enum Illegal_Name {
				foo
				bar
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"RedundantValue": ErrCase{
			Src: `schema test
			enum E {
				foo
				foo
			}`,
			Errs: []ErrCode{compiler.ErrEnumValRedecl},
		},
		"NoValues": ErrCase{
			Src: `schema test
			enum E {}`,
			Errs: []ErrCode{compiler.ErrSyntax},
		},
		"IllegalValueIdentifier": ErrCase{
			Src: `schema test
			enum E {
				_foo
				_bar
			}`,
			Errs: []ErrCode{
				compiler.ErrEnumValIllegalIdent,
				compiler.ErrEnumValIllegalIdent,
			},
		},
		"IllegalValueIdentifier2": ErrCase{
			Src: `schema test
			enum E {
				1foo
				2bar
			}`,
			Errs: []ErrCode{
				compiler.ErrEnumValIllegalIdent,
				compiler.ErrEnumValIllegalIdent,
			},
		},
		"IllegalValueIdentifier3": ErrCase{
			Src: `schema test
			enum E {
				fo_o
				ba_r
			}`,
			Errs: []ErrCode{
				compiler.ErrEnumValIllegalIdent,
				compiler.ErrEnumValIllegalIdent,
			},
		},
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

// TestDeclUnionTypeErrs tests union type declaration errors
func TestDeclUnionTypeErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalTypeName": ErrCase{
			Src: `schema test
			union illegalName {
				String
				Int32
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			union _IllegalName {
				String
				Int32
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			union Illegal_Name {
				String
				Int32
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"OneTypeUnion": ErrCase{
			Src: `schema test
			union U {
				String
			}`,
			Errs: []ErrCode{compiler.ErrUnionMissingOpts},
		},
		"RedundantOptionType": ErrCase{
			Src: `schema test
			union U {
				String
				String
			}`,
			Errs: []ErrCode{compiler.ErrUnionRedund},
		},
		"UndefinedType": ErrCase{
			Src: `schema test
			union U {
				String
				Undefined
			}`,
			Errs: []ErrCode{compiler.ErrTypeUndef},
		},
		"SelfReference": ErrCase{
			Src: `schema test
			union U {
				Int32
				U
			}`,
			Errs: []ErrCode{compiler.ErrUnionSelfref},
		},
		"NonTypeElements": ErrCase{
			Src: `schema test
			union U {
				foo
				bar
			}`,
			Errs: []ErrCode{
				compiler.ErrTypeIllegalIdent,
				compiler.ErrTypeIllegalIdent,
			},
		},
		"NonTypeElements2": ErrCase{
			Src: `schema test
			union U {
				_foo
				_bar
			}`,
			Errs: []ErrCode{
				compiler.ErrTypeIllegalIdent,
				compiler.ErrTypeIllegalIdent,
			},
		},
	})
}
