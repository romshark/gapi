package compiler_test

import (
	"strings"
	"testing"

	"github.com/romshark/gapi/compiler"
	"github.com/romshark/gapi/internal/intset"
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
	c, err := compiler.NewCompiler(source)
	require.NoError(t, err)

	// Compile
	require.NoError(t, c.Compile())
	ast := c.AST()
	require.NotNil(t, ast)

	// Ensure type ID uniqueness
	typeIDs := intset.NewIntSet()
	for _, tp := range ast.Types {
		id := tp.TypeID()
		require.NotEqual(t, compiler.TypeIDUserTypeOffset, id)
		require.False(t, typeIDs.Has(int(id)))
		typeIDs.Insert(int(id))
	}

	// Inspect AST
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

			actualErrMsgs := make([]string, len(actualErrs))
			for i, err := range actualErrs {
				actualErrMsgs[i] = err.Code().String() + ": " + err.Message()
			}
			errMsgs := "Actual errors:\n" + strings.Join(actualErrMsgs, ";\n")

			require.Equal(t, expectedCodes, actualCodes, errMsgs)
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
	
	alias A1 = String
	alias A2 = Uint32
	alias A3 = A1
	`

	test(t, src, func(ast AST) {
		require.Len(t, ast.QueryEndpoints, 0)
		require.Len(t, ast.Mutations, 0)

		require.Len(t, ast.Types, 3)
		a1 := ast.Types[0]
		a2 := ast.Types[1]
		a3 := ast.Types[2]

		type Expectation struct {
			Name        string
			Type        compiler.Type
			AliasedType compiler.Type
		}
		expected := []Expectation{
			Expectation{"A1", a1, compiler.TypeStdString{}},
			Expectation{"A2", a2, compiler.TypeStdUint32{}},
			Expectation{"A3", a3, a1},
		}

		for _, expec := range expected {
			require.Equal(t, expec.Name, expec.Type.Name())
			require.Equal(t, compiler.TypeCategoryAlias, a1.Category())
			require.IsType(t, &compiler.TypeAlias{}, a1)
			require.Equal(
				t,
				expec.AliasedType,
				expec.Type.(*compiler.TypeAlias).AliasedType,
			)
		}
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
		"DirectAliasCycle": ErrCase{
			Src: `schema test
			alias A = A`,
			Errs: []ErrCode{compiler.ErrAliasRecurs},
		},
		"IndirectAliasCycle1": ErrCase{
			Src: `schema test
			alias A = B
			alias B = A`,
			Errs: []ErrCode{compiler.ErrAliasRecurs},
		},
		"IndirectAliasCycle2": ErrCase{
			Src: `schema test
			alias G = H
			alias H = String
			alias F = C
			alias A = B
			alias B = C
			alias C = D
			alias D = A`,
			Errs: []ErrCode{compiler.ErrAliasRecurs},
		},
		"MultipleIndirectAliasesCycles": ErrCase{
			Src: `schema test
			alias A = A
			alias B = C
			alias C = D
			alias D = B
			alias H = K
			alias K = I
			alias I = K`,
			Errs: []ErrCode{
				compiler.ErrAliasRecurs,
				compiler.ErrAliasRecurs,
				compiler.ErrAliasRecurs,
			},
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
		require.Len(t, ast.Types, 3)

		e1 := ast.EnumTypes[0]
		e2 := ast.EnumTypes[1]
		e3 := ast.EnumTypes[2]

		type Expectation struct {
			Name   string
			Type   compiler.Type
			Values []string
		}
		expected := []Expectation{
			Expectation{"E1", e1, []string{"oneVal"}},
			Expectation{"E2", e2, []string{"foo", "bar"}},
			Expectation{"E3", e3, []string{"foo1", "bar2", "baz3"}},
		}

		for _, expec := range expected {
			require.Equal(t, expec.Name, expec.Type.Name())
			require.NotNil(t, ast.FindTypeByName("", expec.Name))
			require.Equal(t, compiler.TypeCategoryEnum, expec.Type.Category())
			require.IsType(t, &compiler.TypeEnum{}, expec.Type)
			tpe := expec.Type.(*compiler.TypeEnum)
			for _, val := range expec.Values {
				require.Contains(t, tpe.Values, val)
				require.Equal(t, tpe.Values[val].Name, val)
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
			Errs: []ErrCode{compiler.ErrEnumNoVal},
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
		require.Len(t, ast.Types, 3)

		u1 := ast.UnionTypes[0]
		u2 := ast.UnionTypes[1]
		u3 := ast.UnionTypes[2]

		type Expectation struct {
			Name  string
			Type  compiler.Type
			Types []compiler.Type
		}
		expected := []Expectation{
			Expectation{"U1", u1, []compiler.Type{
				compiler.TypeStdString{},
				compiler.TypeStdUint32{},
			}},
			Expectation{"U2", u2, []compiler.Type{
				compiler.TypeStdUint32{},
				compiler.TypeStdFloat64{},
				compiler.TypeStdString{},
			}},
			Expectation{"U3", u3, []compiler.Type{
				compiler.TypeStdString{},
				compiler.TypeStdFloat64{},
				compiler.TypeStdInt32{},
				compiler.TypeStdInt64{},
			}},
		}
		for _, expec := range expected {
			require.Equal(t, expec.Name, expec.Type.Name())
			require.Equal(t, compiler.TypeCategoryUnion, expec.Type.Category())
			require.IsType(t, &compiler.TypeUnion{}, expec.Type)
			tpe := expec.Type.(*compiler.TypeUnion)
			for _, referencedType := range expec.Types {
				require.Contains(t, tpe.Types, referencedType.Name())
				require.Equal(
					t,
					referencedType,
					tpe.Types[referencedType.Name()],
				)
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
		"IncludesNone": ErrCase{
			Src: `schema test
			union U {
				Int32
				None
			}`,
			Errs: []ErrCode{compiler.ErrUnionIncludesNone},
		},
	})
}

// TestDeclStructTypes tests struct type declaration
func TestDeclStructTypes(t *testing.T) {
	src := `schema test
	struct S1 {
		x String
	}
	struct S2 {
		x Uint32
		y Float64
	}
	struct S3 {
		optional ?String
		list []Float64
		matrix [][]Int64
		matrix3D [][][]Int64
		optionalList ?[]Int32
		listOfOptionals []?Int32
		optionalListOfOptionals ?[]?Int32
		optionalListOfOptionalListsOfOptionals ?[]?[]?String
	}
	`

	test(t, src, func(ast AST) {
		require.Len(t, ast.QueryEndpoints, 0)
		require.Len(t, ast.Mutations, 0)
		require.Len(t, ast.Types, 3)

		s1 := ast.StructTypes[0]
		s2 := ast.StructTypes[1]
		s3 := ast.StructTypes[2]

		type Expectation struct {
			Name   string
			Type   compiler.Type
			Fields []compiler.StructField
		}
		expected := []Expectation{
			Expectation{"S1", s1, []compiler.StructField{
				compiler.StructField{
					Name:    "x",
					GraphID: compiler.GraphNodeID(1),
					Type:    compiler.TypeStdString{},
				},
			}},
			Expectation{"S2", s2, []compiler.StructField{
				compiler.StructField{
					Name:    "x",
					GraphID: compiler.GraphNodeID(2),
					Type:    compiler.TypeStdUint32{},
				},
				compiler.StructField{
					Name:    "y",
					GraphID: compiler.GraphNodeID(3),
					Type:    compiler.TypeStdFloat64{},
				},
			}},
			Expectation{"S3", s3, []compiler.StructField{
				compiler.StructField{
					Name:    "optional",
					GraphID: compiler.GraphNodeID(4),
					Type: &compiler.TypeOptional{
						Terminal:  compiler.TypeStdString{},
						StoreType: compiler.TypeStdString{},
					},
				},
				compiler.StructField{
					Name:    "list",
					GraphID: compiler.GraphNodeID(5),
					Type: &compiler.TypeList{
						Terminal:  compiler.TypeStdFloat64{},
						StoreType: compiler.TypeStdFloat64{},
					},
				},
				compiler.StructField{
					Name:    "matrix",
					GraphID: compiler.GraphNodeID(6),
					Type: &compiler.TypeList{
						Terminal: compiler.TypeStdInt64{},
						StoreType: &compiler.TypeList{
							Terminal:  compiler.TypeStdInt64{},
							StoreType: compiler.TypeStdInt64{},
						},
					},
				},
				compiler.StructField{
					Name:    "matrix3D",
					GraphID: compiler.GraphNodeID(7),
					Type: &compiler.TypeList{
						Terminal: compiler.TypeStdInt64{},
						StoreType: &compiler.TypeList{
							Terminal: compiler.TypeStdInt64{},
							StoreType: &compiler.TypeList{
								Terminal:  compiler.TypeStdInt64{},
								StoreType: compiler.TypeStdInt64{},
							},
						},
					},
				},
				compiler.StructField{
					Name:    "optionalList",
					GraphID: compiler.GraphNodeID(8),
					Type: &compiler.TypeOptional{
						Terminal: compiler.TypeStdInt32{},
						StoreType: &compiler.TypeList{
							Terminal:  compiler.TypeStdInt32{},
							StoreType: compiler.TypeStdInt32{},
						},
					},
				},
				compiler.StructField{
					Name:    "listOfOptionals",
					GraphID: compiler.GraphNodeID(9),
					Type: &compiler.TypeList{
						Terminal: compiler.TypeStdInt32{},
						StoreType: &compiler.TypeOptional{
							Terminal:  compiler.TypeStdInt32{},
							StoreType: compiler.TypeStdInt32{},
						},
					},
				},
				compiler.StructField{
					Name:    "optionalListOfOptionals",
					GraphID: compiler.GraphNodeID(10),
					Type: &compiler.TypeOptional{
						Terminal: compiler.TypeStdInt32{},
						StoreType: &compiler.TypeList{
							Terminal: compiler.TypeStdInt32{},
							StoreType: &compiler.TypeOptional{
								Terminal:  compiler.TypeStdInt32{},
								StoreType: compiler.TypeStdInt32{},
							},
						},
					},
				},
				compiler.StructField{
					Name:    "optionalListOfOptionalListsOfOptionals",
					GraphID: compiler.GraphNodeID(11),
					Type: &compiler.TypeOptional{
						Terminal: compiler.TypeStdString{},
						StoreType: &compiler.TypeList{
							Terminal: compiler.TypeStdString{},
							StoreType: &compiler.TypeOptional{
								Terminal: compiler.TypeStdString{},
								StoreType: &compiler.TypeList{
									Terminal: compiler.TypeStdString{},
									StoreType: &compiler.TypeOptional{
										Terminal:  compiler.TypeStdString{},
										StoreType: compiler.TypeStdString{},
									},
								},
							},
						},
					},
				},
			}},
		}
		graphNodes := make(map[compiler.GraphNodeID]*compiler.StructField)
		for _, expec := range expected {
			require.Equal(t, expec.Name, expec.Type.Name())
			require.Equal(t, compiler.TypeCategoryStruct, expec.Type.Category())
			require.IsType(t, &compiler.TypeStruct{}, expec.Type)
			structType := expec.Type.(*compiler.TypeStruct)

			// Make sure fields match the expectations
			for i, field := range expec.Fields {
				actualField := structType.Fields[i]
				require.Equal(t, field.Name, actualField.Name)
				require.Equal(
					t,
					field.Type,
					actualField.Type,
					"unexpected type %s for field %s of struct type %s "+
						"(expected: %s)",
					field.Type.String(),
					actualField.Name,
					expec.Name,
					field.Type.String(),
				)
				require.Equal(t, field.GraphID, actualField.GraphID)
				require.Equal(t, structType, expec.Type)

				// Make sure graph node IDs are unique
				require.NotContains(t, graphNodes, actualField.GraphID)
				graphNodes[actualField.GraphID] = actualField
			}
		}

		// Make sure the graph nodes are registered correctly
		require.Len(t, ast.GraphNodes, len(graphNodes))
		for id, field := range graphNodes {
			node := ast.FindGraphNodeByID(id)
			require.NotNil(t, node, "graph node (%d) not found", id)
			require.Equal(t, id, node.GraphNodeID())
			require.Equal(t, field.Struct, node.Parent())
		}
	})
}

// TestDeclStructTypeErrs tests struct type declaration errors
func TestDeclStructTypeErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalTypeName": ErrCase{
			Src: `schema test
			struct illegalName {
				foo String
				bar String
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			struct _IllegalName {
				foo String
				bar String
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			struct Illegal_Name {
				foo String
				bar String
			}`,
			Errs: []ErrCode{compiler.ErrTypeIllegalIdent},
		},
		"NoFields": ErrCase{
			Src: `schema test
			struct S {}`,
			Errs: []ErrCode{compiler.ErrStructNoFields},
		},
		"RedundantField": ErrCase{
			Src: `schema test
			struct S {
				foo String
				foo String
			}`,
			Errs: []ErrCode{compiler.ErrStructFieldRedecl},
		},
		"IllegalFieldIdentifier": ErrCase{
			Src: `schema test
			struct S {
				_foo String
				_bar String
			}`,
			Errs: []ErrCode{
				compiler.ErrStructFieldIllegalIdent,
				compiler.ErrStructFieldIllegalIdent,
			},
		},
		"IllegalFieldIdentifier2": ErrCase{
			Src: `schema test
			struct S {
				1foo String
				2bar String
			}`,
			Errs: []ErrCode{
				compiler.ErrStructFieldIllegalIdent,
				compiler.ErrStructFieldIllegalIdent,
			},
		},
		"IllegalFieldIdentifier3": ErrCase{
			Src: `schema test
			struct S {
				fo_o String
				ba_r String
			}`,
			Errs: []ErrCode{
				compiler.ErrStructFieldIllegalIdent,
				compiler.ErrStructFieldIllegalIdent,
			},
		},
	})
}
