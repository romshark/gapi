package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/romshark/gapi/compiler/parser"
	"github.com/romshark/gapi/internal/intset"
	"github.com/stretchr/testify/require"
)

type AST = *parser.AST
type ErrCode = parser.ErrCode

func src(src string) parser.SourceFile {
	return parser.SourceFile{
		File: parser.File{
			Name: "test.schema",
			Path: "/tests/",
		},
		Src: src,
	}
}

func test(
	t *testing.T,
	source string,
	astInspector func(AST),
) {
	// Initialize parser
	pr, err := parser.NewParser()
	require.NoError(t, err)
	require.NotNil(t, pr)

	// Compile
	require.NoError(t, pr.Parse(src(source)))
	ast := pr.AST()
	require.NotNil(t, ast)

	verifyAST(t, ast)

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
			// Initialize parser
			pr, err := parser.NewParser()
			require.NoError(t, err)
			require.NotNil(t, pr)

			// Parse
			require.Error(t, pr.Parse(src(errCase.Src)))
			actualErrs := pr.Errors()
			require.True(t, len(actualErrs) > 0)
			require.Nil(t, pr.AST())

			type Err struct {
				Code ErrCode
				Name string
			}

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
				actualErrMsgs[i] = err.Error()
			}
			errMsgs := "Actual errors:\n" + strings.Join(actualErrMsgs, ";\n")

			require.Equal(t, expectedCodes, actualCodes, errMsgs)
		})
	}
}

func verifyAST(t *testing.T, ast AST) {
	typeIDs := intset.NewIntSet()
	graphNodeIDs := intset.NewIntSet()
	paramIDs := intset.NewIntSet()

	// Ensure type ID uniqueness
	for _, tp := range ast.Types {
		id := tp.TypeID()
		require.NotEqual(t, parser.TypeIDUserTypeOffset, id)
		require.False(t, typeIDs.Has(int(id)))
		typeIDs.Insert(int(id))

		// Ensure correct type ID mapping
		require.Equal(t, tp, ast.FindTypeByID(id))
	}

	// Ensure graph node ID uniqueness
	for _, str := range ast.StructTypes {
		for _, fld := range str.(*parser.TypeStruct).Fields {
			intID := int(fld.GraphNodeID())
			require.False(t, graphNodeIDs.Has(intID))
			graphNodeIDs.Insert(intID)
		}
	}
	for _, rsv := range ast.ResolverTypes {
		for _, prop := range rsv.(*parser.TypeResolver).Properties {
			intID := int(prop.GraphNodeID())
			require.False(t, graphNodeIDs.Has(intID))
			graphNodeIDs.Insert(intID)

			// Ensure parameter ID uniqueness
			for _, param := range prop.Parameters {
				intID := int(param.ID)
				require.False(t, paramIDs.Has(intID))
				paramIDs.Insert(intID)
			}
		}
	}
}

// TestDeclSchemaErrs tests schema declaration errors
func TestDeclSchemaErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalName": ErrCase{
			Src: `schema _illegalName
			enum E { e }`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalName2": ErrCase{
			Src: `schema illegal_Name
			enum E { e }`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalName3": ErrCase{
			Src: `schema IllegalName
			enum E { e }`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
	})
}

// TestDeclTypeErrs tests type declaration errors (generic errors)
func TestDeclTypeErrs(t *testing.T) {
	testCases := map[string]ErrCase{
		"IllegalName": ErrCase{
			Src: `schema test
			enum _illegalName { e }`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalName2": ErrCase{
			Src: `schema test
			enum illegal_Name { e }`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalName3": ErrCase{
			Src: `schema test
			enum Illegal_Name { e }`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"RedeclUserType": ErrCase{
			Src: `schema test
			enum X { a b }
			alias X = String`,
			Errs: []ErrCode{parser.ErrTypeRedecl},
		},
	}

	// Test primitive type redeclaration
	primitiveTypeNames := []string{
		"None",
		"Bool",
		"Byte",
		"Int32",
		"Uint32",
		"Int64",
		"Uint64",
		"Float64",
		"String",
		"Time",
	}
	for _, primTypeName := range primitiveTypeNames {
		testCases[fmt.Sprintf("RedeclPrimitive(%s)", primTypeName)] = ErrCase{
			Src: fmt.Sprintf(
				`schema tst
				enum %s { e }`,
				primTypeName,
			),
			Errs: []ErrCode{parser.ErrTypeRedecl},
		}
	}

	testErrs(t, testCases)
}

// TestASTAliases tests alias type declaration in AST
func TestASTAliases(t *testing.T) {
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
			Type        parser.Type
			AliasedType parser.Type
		}
		expected := []Expectation{
			Expectation{"A1", a1, parser.TypeStdString{}},
			Expectation{"A2", a2, parser.TypeStdUint32{}},
			Expectation{"A3", a3, a1},
		}

		for _, expec := range expected {
			require.Equal(t, expec.Name, expec.Type.Name())
			require.Equal(t, parser.TypeCategoryAlias, a1.Category())
			require.IsType(t, &parser.TypeAlias{}, a1)
			require.Equal(
				t,
				expec.AliasedType,
				expec.Type.(*parser.TypeAlias).AliasedType,
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
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			alias Illegal_Name = String`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			alias _IllegalName = String`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalAliasedTypeName": ErrCase{
			Src: `schema test
			alias A = illegalName`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalAliasedTypeName2": ErrCase{
			Src: `schema test
			alias A = Illegal_Name`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalAliasedTypeName3": ErrCase{
			Src: `schema test
			alias A = _IllegalName`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"UndefinedAliasedType": ErrCase{
			Src: `schema test
			alias A = Undefined`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"DirectAliasCycle": ErrCase{
			Src: `schema test
			alias A = A`,
			Errs: []ErrCode{parser.ErrAliasRecurs},
		},
		"IndirectAliasCycle1": ErrCase{
			Src: `schema test
			alias A = B
			alias B = A`,
			Errs: []ErrCode{parser.ErrAliasRecurs},
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
			Errs: []ErrCode{parser.ErrAliasRecurs},
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
				parser.ErrAliasRecurs,
				parser.ErrAliasRecurs,
				parser.ErrAliasRecurs,
			},
		},
	})
}

// TestASTEnums tests enum type declaration in AST
func TestASTEnums(t *testing.T) {
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
			Type   parser.Type
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
			require.Equal(t, parser.TypeCategoryEnum, expec.Type.Category())
			require.IsType(t, &parser.TypeEnum{}, expec.Type)
			tpe := expec.Type.(*parser.TypeEnum)

			containsVal := func(
				expected string,
				vals []*parser.EnumValue,
			) bool {
				for _, val := range vals {
					if val.Name == expected {
						return true
					}
				}
				return false
			}

			for i, val := range expec.Values {
				require.True(t, containsVal(val, tpe.Values))
				require.Equal(t, tpe.Values[i].Name, val)
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
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			enum _IllegalName {
				foo
				bar
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			enum Illegal_Name {
				foo
				bar
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"RedundantValue": ErrCase{
			Src: `schema test
			enum E {
				foo
				foo
			}`,
			Errs: []ErrCode{parser.ErrEnumValRedecl},
		},
		"NoValues": ErrCase{
			Src: `schema test
			enum E {}`,
			Errs: []ErrCode{parser.ErrEnumNoVal},
		},
		"IllegalValueIdentifier": ErrCase{
			Src: `schema test
			enum E {
				_foo
				_bar
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalValueIdentifier2": ErrCase{
			Src: `schema test
			enum E {
				1foo
				2bar
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalValueIdentifier3": ErrCase{
			Src: `schema test
			enum E {
				fo_o
				ba_r
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
	})
}

// TestASTUnions tests union type declarations in AST
func TestASTUnions(t *testing.T) {
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
			Type  parser.Type
			Types []parser.Type
		}
		expected := []Expectation{
			Expectation{"U1", u1, []parser.Type{
				parser.TypeStdString{},
				parser.TypeStdUint32{},
			}},
			Expectation{"U2", u2, []parser.Type{
				parser.TypeStdUint32{},
				parser.TypeStdFloat64{},
				parser.TypeStdString{},
			}},
			Expectation{"U3", u3, []parser.Type{
				parser.TypeStdString{},
				parser.TypeStdFloat64{},
				parser.TypeStdInt32{},
				parser.TypeStdInt64{},
			}},
		}
		for _, expec := range expected {
			require.Equal(t, expec.Name, expec.Type.Name())
			require.Equal(t, parser.TypeCategoryUnion, expec.Type.Category())
			require.IsType(t, &parser.TypeUnion{}, expec.Type)
			tpe := expec.Type.(*parser.TypeUnion)

			containsType := func(
				expected string,
				types []parser.Type,
			) bool {
				for _, tp := range types {
					if tp.Name() == expected {
						return true
					}
				}
				return false
			}

			for i, referencedType := range expec.Types {
				require.True(t, containsType(referencedType.Name(), tpe.Types))
				require.Equal(t, referencedType, tpe.Types[i])
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
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			union _IllegalName {
				String
				Int32
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			union Illegal_Name {
				String
				Int32
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"OneTypeUnion": ErrCase{
			Src: `schema test
			union U {
				String
			}`,
			Errs: []ErrCode{parser.ErrUnionMissingOpts},
		},
		"RedundantOptionType": ErrCase{
			Src: `schema test
			union U {
				String
				String
			}`,
			Errs: []ErrCode{parser.ErrUnionRedund},
		},
		"UndefinedType": ErrCase{
			Src: `schema test
			union U {
				String
				Undefined
			}`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"SelfReference": ErrCase{
			Src: `schema test
			union U {
				Int32
				U
			}`,
			Errs: []ErrCode{parser.ErrUnionRecurs},
		},
		"NonTypeElements": ErrCase{
			Src: `schema test
			union U {
				foo
				bar
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"NonTypeElements2": ErrCase{
			Src: `schema test
			union U {
				_foo
				_bar
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IncludesNone": ErrCase{
			Src: `schema test
			union U {
				Int32
				None
			}`,
			Errs: []ErrCode{parser.ErrUnionIncludesNone},
		},
	})
}

// TestASTStructs tests struct type declarations in AST
func TestASTStructs(t *testing.T) {
	src := `schema test
	struct S1 {
		x String
	}
	struct S2 {
		x []S2
		y ?S2
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
			Type   parser.Type
			Fields []parser.StructField
		}
		expected := []Expectation{
			Expectation{"S1", s1, []parser.StructField{
				parser.StructField{
					Name:    "x",
					GraphID: 1,
					Type:    parser.TypeStdString{},
				},
			}},
			Expectation{"S2", s2, []parser.StructField{
				parser.StructField{
					Name:    "x",
					GraphID: 2,
					Type: &parser.TypeList{
						Terminal:  s2,
						StoreType: s2,
					},
				},
				parser.StructField{
					Name:    "y",
					GraphID: 3,
					Type: &parser.TypeOptional{
						Terminal:  s2,
						StoreType: s2,
					},
				},
			}},
			Expectation{"S3", s3, []parser.StructField{
				parser.StructField{
					Name:    "optional",
					GraphID: 4,
					Type: &parser.TypeOptional{
						Terminal:  parser.TypeStdString{},
						StoreType: parser.TypeStdString{},
					},
				},
				parser.StructField{
					Name:    "list",
					GraphID: 5,
					Type: &parser.TypeList{
						Terminal:  parser.TypeStdFloat64{},
						StoreType: parser.TypeStdFloat64{},
					},
				},
				parser.StructField{
					Name:    "matrix",
					GraphID: 6,
					Type: &parser.TypeList{
						Terminal: parser.TypeStdInt64{},
						StoreType: &parser.TypeList{
							Terminal:  parser.TypeStdInt64{},
							StoreType: parser.TypeStdInt64{},
						},
					},
				},
				parser.StructField{
					Name:    "matrix3D",
					GraphID: 7,
					Type: &parser.TypeList{
						Terminal: parser.TypeStdInt64{},
						StoreType: &parser.TypeList{
							Terminal: parser.TypeStdInt64{},
							StoreType: &parser.TypeList{
								Terminal:  parser.TypeStdInt64{},
								StoreType: parser.TypeStdInt64{},
							},
						},
					},
				},
				parser.StructField{
					Name:    "optionalList",
					GraphID: 8,
					Type: &parser.TypeOptional{
						Terminal: parser.TypeStdInt32{},
						StoreType: &parser.TypeList{
							Terminal:  parser.TypeStdInt32{},
							StoreType: parser.TypeStdInt32{},
						},
					},
				},
				parser.StructField{
					Name:    "listOfOptionals",
					GraphID: 9,
					Type: &parser.TypeList{
						Terminal: parser.TypeStdInt32{},
						StoreType: &parser.TypeOptional{
							Terminal:  parser.TypeStdInt32{},
							StoreType: parser.TypeStdInt32{},
						},
					},
				},
				parser.StructField{
					Name:    "optionalListOfOptionals",
					GraphID: 10,
					Type: &parser.TypeOptional{
						Terminal: parser.TypeStdInt32{},
						StoreType: &parser.TypeList{
							Terminal: parser.TypeStdInt32{},
							StoreType: &parser.TypeOptional{
								Terminal:  parser.TypeStdInt32{},
								StoreType: parser.TypeStdInt32{},
							},
						},
					},
				},
				parser.StructField{
					Name:    "optionalListOfOptionalListsOfOptionals",
					GraphID: 11,
					Type: &parser.TypeOptional{
						Terminal: parser.TypeStdString{},
						StoreType: &parser.TypeList{
							Terminal: parser.TypeStdString{},
							StoreType: &parser.TypeOptional{
								Terminal: parser.TypeStdString{},
								StoreType: &parser.TypeList{
									Terminal: parser.TypeStdString{},
									StoreType: &parser.TypeOptional{
										Terminal:  parser.TypeStdString{},
										StoreType: parser.TypeStdString{},
									},
								},
							},
						},
					},
				},
			}},
		}
		graphNodes := make(map[parser.GraphNodeID]*parser.StructField)
		for _, expec := range expected {
			require.Equal(t, expec.Name, expec.Type.Name())
			require.Equal(t, parser.TypeCategoryStruct, expec.Type.Category())
			require.IsType(t, &parser.TypeStruct{}, expec.Type)
			structType := expec.Type.(*parser.TypeStruct)

			// Make sure fields match the expectations
			for i, field := range expec.Fields {
				actualField := structType.Fields[i]
				require.Equal(t, field.Name, actualField.Name)
				require.Equal(t, field.Type.Name(), actualField.Type.Name())
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
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			struct _IllegalName {
				foo String
				bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			struct Illegal_Name {
				foo String
				bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"NoFields": ErrCase{
			Src: `schema test
			struct S {}`,
			Errs: []ErrCode{parser.ErrStructNoFields},
		},
		"RedundantField": ErrCase{
			Src: `schema test
			struct S {
				foo String
				foo String
			}`,
			Errs: []ErrCode{parser.ErrStructFieldRedecl},
		},
		"IllegalFieldIdentifier": ErrCase{
			Src: `schema test
			struct S {
				_foo String
				_bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalFieldIdentifier2": ErrCase{
			Src: `schema test
			struct S {
				1foo String
				2bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalFieldIdentifier3": ErrCase{
			Src: `schema test
			struct S {
				fo_o String
				ba_r String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"RecursDirect": ErrCase{
			Src: `schema test
			struct S {
				s S
			}`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // S.s -> S
			},
		},
		"RecursIndirect": ErrCase{
			Src: `schema test
			struct X {
				s S
			}
			struct S {
				x X
			}`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // X.s -> S.x -> X
			},
		},
		"RecursIndirect2": ErrCase{
			Src: `schema test
			struct Y {
				s S
			}
			struct X {
				y Y
			}
			struct S {
				x X
			}`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // Y.s -> S.x -> X.y -> Y
			},
		},
		"RecursIndirect3": ErrCase{
			Src: `schema test
			struct Y {
				s S
				z S
			}
			struct X {
				y Y
				s S
			}
			struct S {
				x X
				y Y
			}`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // Y.s -> S.x -> X.y -> Y
			},
		},
		"RecursMultiple": ErrCase{
			Src: `schema test
			struct A {
				a A
			}
			struct B {
				b B
			}
			struct X {
				y Y
			}
			struct Y {
				x X
			}`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // A.a -> A
				parser.ErrStructRecurs, // B.b -> B
				parser.ErrStructRecurs, // X.y -> Y.x -> X
			},
		},
	})
}

// TestASTResolvers tests resolver type declarations in AST
func TestASTResolvers(t *testing.T) {
	src := `schema test
	resolver R1 {
		x String
	}
	resolver R2 {
		r R1
		x R2
		y []R2
		z ?R3
	}
	resolver R3 {
		optional ?String
		list []Float64
		matrix [][]Int64
		matrix3D [][][]Int64
		optionalList ?[]Int32
		listOfOptionals []?Int32
		optionalListOfOptionals ?[]?Int32
		optionalListOfOptionalListsOfOptionals ?[]?[]?String
	}
	resolver R4 {
		x(x Int32) Int32
		y(x Int32, y ?String, z ?[]Bool) String
	}
	`

	test(t, src, func(ast AST) {
		require.Len(t, ast.QueryEndpoints, 0)
		require.Len(t, ast.Mutations, 0)
		require.Len(t, ast.Types, 4)

		r1 := ast.ResolverTypes[0]
		r2 := ast.ResolverTypes[1]
		r3 := ast.ResolverTypes[2]
		r4 := ast.ResolverTypes[3]

		type Expectation struct {
			Name  string
			Type  parser.Type
			Props []parser.ResolverProperty
		}
		expected := []Expectation{
			Expectation{"R1", r1, []parser.ResolverProperty{
				parser.ResolverProperty{
					Name:    "x",
					GraphID: 1,
					Type:    parser.TypeStdString{},
				},
			}},
			Expectation{"R2", r2, []parser.ResolverProperty{
				parser.ResolverProperty{
					Name:    "r",
					GraphID: 2,
					Type:    r1,
				},
				parser.ResolverProperty{
					Name:    "x",
					GraphID: 3,
					Type:    r2,
				},
				parser.ResolverProperty{
					Name:    "y",
					GraphID: 4,
					Type: &parser.TypeList{
						Terminal:  r2,
						StoreType: r2,
					},
				},
				parser.ResolverProperty{
					Name:    "z",
					GraphID: 5,
					Type: &parser.TypeOptional{
						Terminal:  r3,
						StoreType: r3,
					},
				},
			}},
			Expectation{"R3", r3, []parser.ResolverProperty{
				parser.ResolverProperty{
					Name:    "optional",
					GraphID: 6,
					Type: &parser.TypeOptional{
						Terminal:  parser.TypeStdString{},
						StoreType: parser.TypeStdString{},
					},
				},
				parser.ResolverProperty{
					Name:    "list",
					GraphID: 7,
					Type: &parser.TypeList{
						Terminal:  parser.TypeStdFloat64{},
						StoreType: parser.TypeStdFloat64{},
					},
				},
				parser.ResolverProperty{
					Name:    "matrix",
					GraphID: 8,
					Type: &parser.TypeList{
						Terminal: parser.TypeStdInt64{},
						StoreType: &parser.TypeList{
							Terminal:  parser.TypeStdInt64{},
							StoreType: parser.TypeStdInt64{},
						},
					},
				},
				parser.ResolverProperty{
					Name:    "matrix3D",
					GraphID: 9,
					Type: &parser.TypeList{
						Terminal: parser.TypeStdInt64{},
						StoreType: &parser.TypeList{
							Terminal: parser.TypeStdInt64{},
							StoreType: &parser.TypeList{
								Terminal:  parser.TypeStdInt64{},
								StoreType: parser.TypeStdInt64{},
							},
						},
					},
				},
				parser.ResolverProperty{
					Name:    "optionalList",
					GraphID: 10,
					Type: &parser.TypeOptional{
						Terminal: parser.TypeStdInt32{},
						StoreType: &parser.TypeList{
							Terminal:  parser.TypeStdInt32{},
							StoreType: parser.TypeStdInt32{},
						},
					},
				},
				parser.ResolverProperty{
					Name:    "listOfOptionals",
					GraphID: 11,
					Type: &parser.TypeList{
						Terminal: parser.TypeStdInt32{},
						StoreType: &parser.TypeOptional{
							Terminal:  parser.TypeStdInt32{},
							StoreType: parser.TypeStdInt32{},
						},
					},
				},
				parser.ResolverProperty{
					Name:    "optionalListOfOptionals",
					GraphID: 12,
					Type: &parser.TypeOptional{
						Terminal: parser.TypeStdInt32{},
						StoreType: &parser.TypeList{
							Terminal: parser.TypeStdInt32{},
							StoreType: &parser.TypeOptional{
								Terminal:  parser.TypeStdInt32{},
								StoreType: parser.TypeStdInt32{},
							},
						},
					},
				},
				parser.ResolverProperty{
					Name:    "optionalListOfOptionalListsOfOptionals",
					GraphID: 13,
					Type: &parser.TypeOptional{
						Terminal: parser.TypeStdString{},
						StoreType: &parser.TypeList{
							Terminal: parser.TypeStdString{},
							StoreType: &parser.TypeOptional{
								Terminal: parser.TypeStdString{},
								StoreType: &parser.TypeList{
									Terminal: parser.TypeStdString{},
									StoreType: &parser.TypeOptional{
										Terminal:  parser.TypeStdString{},
										StoreType: parser.TypeStdString{},
									},
								},
							},
						},
					},
				},
			}},
			Expectation{"R4", r4, []parser.ResolverProperty{
				parser.ResolverProperty{
					Name:    "x",
					GraphID: 14,
					Type:    parser.TypeStdInt32{},
					Parameters: []*parser.Parameter{
						&parser.Parameter{
							Name: "x",
							ID:   1,
							Type: parser.TypeStdInt32{},
						},
					},
				},
				parser.ResolverProperty{
					Name:    "y",
					GraphID: 15,
					Type:    parser.TypeStdString{},
					Parameters: []*parser.Parameter{
						&parser.Parameter{
							Name: "x",
							ID:   2,
							Type: parser.TypeStdInt32{},
						},
						&parser.Parameter{
							Name: "y",
							ID:   3,
							Type: &parser.TypeOptional{
								Terminal:  parser.TypeStdString{},
								StoreType: parser.TypeStdString{},
							},
						},
						&parser.Parameter{
							Name: "z",
							ID:   4,
							Type: &parser.TypeOptional{
								Terminal: parser.TypeStdBool{},
								StoreType: &parser.TypeList{
									Terminal:  parser.TypeStdBool{},
									StoreType: parser.TypeStdBool{},
								},
							},
						},
					},
				},
			}},
		}
		graphNodes := make(map[parser.GraphNodeID]*parser.ResolverProperty)
		parameters := make(map[parser.ParamID]*parser.Parameter)
		for _, expec := range expected {
			require.Equal(t, expec.Name, expec.Type.Name())
			require.Equal(
				t,
				parser.TypeCategoryResolver,
				expec.Type.Category(),
			)
			require.IsType(t, &parser.TypeResolver{}, expec.Type)
			resolverType := expec.Type.(*parser.TypeResolver)

			// Make sure properties match expectations
			require.Len(t, resolverType.Properties, len(expec.Props))
			for i, prop := range expec.Props {
				actualProp := resolverType.Properties[i]
				require.Equal(t, prop.Name, actualProp.Name)
				require.Equal(
					t,
					prop.Type,
					actualProp.Type,
					"unexpected type %s for property %s of resolver type %s "+
						"(expected: %s)",
					actualProp.Type,
					actualProp.Name,
					expec.Name,
					prop.Type,
				)
				require.Equal(
					t,
					prop.GraphID,
					actualProp.GraphID,
					"unexpected graph ID %d for property %s "+
						"of resolver type %s (expected: %d)",
					actualProp.GraphID,
					actualProp.Name,
					expec.Name,
					prop.GraphID,
				)
				require.Equal(t, resolverType, expec.Type)

				// Make sure graph node IDs are unique
				require.NotContains(t, graphNodes, actualProp.GraphID)
				graphNodes[actualProp.GraphID] = actualProp

				// Make sure property parameters match expectations
				require.Len(t, actualProp.Parameters, len(prop.Parameters))
				for j, param := range prop.Parameters {
					actualParam := actualProp.Parameters[j]
					require.Equal(t, param.Name, actualParam.Name)
					require.Equal(t, param.ID, actualParam.ID)
					require.Equal(
						t,
						param.Type,
						actualParam.Type,
						"unexpected type %s for parameter %s "+
							"of property %s of resolver type %s",
						actualParam.Type,
						param.Name,
						actualProp.Name,
						expec.Name,
					)
					require.IsType(
						t,
						&parser.ResolverProperty{},
						actualParam.Target,
					)
					require.Equal(
						t,
						actualProp,
						actualParam.Target.(*parser.ResolverProperty),
					)
					parameters[actualParam.ID] = actualParam
				}
			}
		}

		// Make sure the graph nodes are registered correctly
		require.Len(t, ast.GraphNodes, len(graphNodes))
		for id, prop := range graphNodes {

			node := ast.FindGraphNodeByID(id)
			require.NotNil(t, node, "graph node (%d) not found in AST", id)
			require.Equal(t, id, node.GraphNodeID())
			require.Equal(t, prop.Resolver, node.Parent())
		}

		// Make sure parameters are registered correctly
		for id, p := range parameters {
			param := ast.FindParameterByID(id)
			require.NotNil(t, param, "parameter (%d) not found in AST", id)
			require.Equal(t, id, param.ID)
			require.IsType(t, &parser.ResolverProperty{}, param.Target)
			require.Equal(t, p.Target, param.Target)
		}
	})
}

// TestDeclResolverTypeErrs tests resolver type declaration errors
func TestDeclResolverTypeErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalTypeName": ErrCase{
			Src: `schema test
			resolver illegalName {
				foo String
				bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			resolver _IllegalName {
				foo String
				bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			resolver Illegal_Name {
				foo String
				bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"NoProps": ErrCase{
			Src: `schema test
			resolver S {}`,
			Errs: []ErrCode{parser.ErrResolverNoProps},
		},
		"RedundantProp": ErrCase{
			Src: `schema test
			resolver S {
				foo String
				foo String
			}`,
			Errs: []ErrCode{parser.ErrResolverPropRedecl},
		},
		"IllegalPropIdentifier": ErrCase{
			Src: `schema test
			resolver S {
				_foo String
				_bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropIdentifier2": ErrCase{
			Src: `schema test
			resolver S {
				1foo String
				2bar String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropIdentifier3": ErrCase{
			Src: `schema test
			resolver S {
				fo_o String
				ba_r String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropParamIdentifier": ErrCase{
			Src: `schema test
			resolver S {
				foo(_foo String) String
				bar(_bar String) String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropParamIdentifier2": ErrCase{
			Src: `schema test
			resolver S {
				foo(1foo String) String
				bar(2bar String) String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropParamIdentifier3": ErrCase{
			Src: `schema test
			resolver S {
				foo(fo_o String) String
				bar(ba_r String) String
			}`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"RedundantPropParam": ErrCase{
			Src: `schema test
			resolver S {
				foo(foo String, foo Int32) String
			}`,
			Errs: []ErrCode{parser.ErrParamRedecl},
		},
	})
}

// TestASTQueries tests query declarations in AST
func TestASTQueries(t *testing.T) {
	src := `schema test
	struct Foo {
		foo String
	}
	resolver Bar {
		bar String
	}
	query foo Foo
	query bar Bar
	query str String
	query foo2(foo Foo) Foo
	query bar2(bar Int32, baz Float64) Bar
	query baz(
		first Int32,
		second Bool,
		third Uint64,
	) String
	`

	test(t, src, func(ast AST) {
		require.Len(t, ast.Types, 2)
		require.Len(t, ast.Mutations, 0)
		require.Len(t, ast.QueryEndpoints, 6)

		require.Len(t, ast.StructTypes, 1)
		tFoo := ast.StructTypes[0]

		require.Len(t, ast.ResolverTypes, 1)
		tBar := ast.ResolverTypes[0]

		expected := []parser.Query{
			parser.Query{
				GraphID: 4,
				Name:    "bar",
				Type:    tBar,
			},
			parser.Query{
				GraphID: 7,
				Name:    "bar2",
				Type:    tBar,
				Parameters: []*parser.Parameter{
					&parser.Parameter{
						ID:   2,
						Name: "bar",
						Type: parser.TypeStdInt32{},
					},
					&parser.Parameter{
						ID:   3,
						Name: "baz",
						Type: parser.TypeStdFloat64{},
					},
				},
			},
			parser.Query{
				GraphID: 8,
				Name:    "baz",
				Type:    parser.TypeStdString{},
				Parameters: []*parser.Parameter{
					&parser.Parameter{
						ID:   4,
						Name: "first",
						Type: parser.TypeStdInt32{},
					},
					&parser.Parameter{
						ID:   5,
						Name: "second",
						Type: parser.TypeStdBool{},
					},
					&parser.Parameter{
						ID:   6,
						Name: "third",
						Type: parser.TypeStdUint64{},
					},
				},
			},
			parser.Query{
				GraphID: 3,
				Name:    "foo",
				Type:    tFoo,
			},
			parser.Query{
				GraphID: 6,
				Name:    "foo2",
				Type:    tFoo,
				Parameters: []*parser.Parameter{
					&parser.Parameter{
						ID:   1,
						Name: "foo",
						Type: tFoo,
					},
				},
			},
			parser.Query{
				GraphID: 5,
				Name:    "str",
				Type:    parser.TypeStdString{},
			},
		}
		require.Len(t, ast.QueryEndpoints, len(expected))
		for i1, expec := range expected {
			require.IsType(t, parser.Query{}, expec)

			actual := ast.QueryEndpoints[i1]
			require.Equal(t, expec.Name, actual.Name)
			require.Equal(t, expec.GraphID, actual.GraphID)
			require.Equal(t, expec.Type, actual.Type)

			// Make sure the graph nodes are registered correctly
			foundNode := ast.FindGraphNodeByID(expec.GraphID)
			require.Equal(t, actual, foundNode)

			// Make sure parameters match expectations
			require.Len(t, actual.Parameters, len(expec.Parameters))
			for i2, param := range expec.Parameters {
				actualParam := actual.Parameters[i2]
				require.Equal(t, param.Name, actualParam.Name)
				require.Equal(t, param.ID, actualParam.ID)
				require.Equal(t, param.Type, actualParam.Type)
				require.IsType(
					t,
					&parser.Query{},
					actualParam.Target,
				)
				require.Equal(
					t,
					actual,
					actualParam.Target.(*parser.Query),
				)

				// Make sure parameters are registered correctly
				regParam := ast.FindParameterByID(param.ID)
				require.Equal(t, actualParam, regParam)
			}
		}
	})
}

// TestASTMutations tests mutation declarations in AST
func TestASTMutations(t *testing.T) {
	src := `schema test
	struct Foo {
		foo String
	}
	resolver Bar {
		bar String
	}
	mutation foo Foo
	mutation bar Bar
	mutation str String
	mutation foo2(foo Foo) Foo
	mutation bar2(bar Int32, baz Float64) Bar
	mutation baz(
		first Int32,
		second Bool,
		third Uint64,
	) String
	`

	test(t, src, func(ast AST) {
		require.Len(t, ast.Types, 2)
		require.Len(t, ast.Mutations, 6)
		require.Len(t, ast.QueryEndpoints, 0)

		require.Len(t, ast.StructTypes, 1)
		tFoo := ast.StructTypes[0]

		require.Len(t, ast.ResolverTypes, 1)
		tBar := ast.ResolverTypes[0]

		expected := []parser.Mutation{
			parser.Mutation{
				GraphID: 4,
				Name:    "bar",
				Type:    tBar,
			},
			parser.Mutation{
				GraphID: 7,
				Name:    "bar2",
				Type:    tBar,
				Parameters: []*parser.Parameter{
					&parser.Parameter{
						ID:   2,
						Name: "bar",
						Type: parser.TypeStdInt32{},
					},
					&parser.Parameter{
						ID:   3,
						Name: "baz",
						Type: parser.TypeStdFloat64{},
					},
				},
			},
			parser.Mutation{
				GraphID: 8,
				Name:    "baz",
				Type:    parser.TypeStdString{},
				Parameters: []*parser.Parameter{
					&parser.Parameter{
						ID:   4,
						Name: "first",
						Type: parser.TypeStdInt32{},
					},
					&parser.Parameter{
						ID:   5,
						Name: "second",
						Type: parser.TypeStdBool{},
					},
					&parser.Parameter{
						ID:   6,
						Name: "third",
						Type: parser.TypeStdUint64{},
					},
				},
			},
			parser.Mutation{
				GraphID: 3,
				Name:    "foo",
				Type:    tFoo,
			},
			parser.Mutation{
				GraphID: 6,
				Name:    "foo2",
				Type:    tFoo,
				Parameters: []*parser.Parameter{
					&parser.Parameter{
						ID:   1,
						Name: "foo",
						Type: tFoo,
					},
				},
			},
			parser.Mutation{
				GraphID: 5,
				Name:    "str",
				Type:    parser.TypeStdString{},
			},
		}
		require.Len(t, ast.Mutations, len(expected))
		for i1, expec := range expected {
			require.IsType(t, parser.Mutation{}, expec)

			actual := ast.Mutations[i1]
			require.Equal(t, expec.Name, actual.Name)
			require.Equal(t, expec.GraphID, actual.GraphID)
			require.Equal(t, expec.Type, actual.Type)

			// Make sure the graph nodes are registered correctly
			foundNode := ast.FindGraphNodeByID(expec.GraphID)
			require.Equal(t, actual, foundNode)

			// Make sure parameters match expectations
			require.Len(t, actual.Parameters, len(expec.Parameters))
			for i2, param := range expec.Parameters {
				actualParam := actual.Parameters[i2]
				require.Equal(t, param.Name, actualParam.Name)
				require.Equal(t, param.ID, actualParam.ID)
				require.Equal(t, param.Type, actualParam.Type)
				require.IsType(
					t,
					&parser.Mutation{},
					actualParam.Target,
				)
				require.Equal(
					t,
					actual,
					actualParam.Target.(*parser.Mutation),
				)

				// Make sure parameters are registered correctly
				regParam := ast.FindParameterByID(param.ID)
				require.Equal(t, actualParam, regParam)
			}
		}
	})
}

// TestDeclQueryErrs tests query endpoint declaration errors
func TestDeclQueryErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalTypeName": ErrCase{
			Src: `schema test
			query IllegalName String`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			query _illegalName String`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			query illegal_Name String`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
	})
}

// TestParamImpureType tests specifying parameters of non-pure types
func TestParamImpureType(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"Resolver": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			query q(param R) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureUnion": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			union U {
				R
				Int32
			}
			query q(param U) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureAliasToResolver": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			alias A = R
			query q(param A) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureUnionOfImpureAlias": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			alias A = R
			union U {
				Int32
				A
			}
			query q(param U) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureOptional": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			query q(param ?R) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureList": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			query q(param []R) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureOptionalList": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			query q(param ?[]R) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"None": ErrCase{
			Src: `schema test
			query q(param None) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
	})
}

// TestStructImpureFieldType tests specifying struct fields of non-pure types
func TestStructImpureFieldType(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"Resolver": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			struct S {
				f R
			}`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureUnion": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			union U {
				R
				Int32
			}
			struct S {
				f U
			}`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureAliasToResolver": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			alias A = R
			struct S {
				f A
			}`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureUnionOfImpureAlias": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			alias A = R
			union U {
				Int32
				A
			}
			struct S {
				f U
			}`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureOptional": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			struct S {
				f ?R
			}`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureList": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			struct S {
				f []R
			}`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureOptionalList": ErrCase{
			Src: `schema test
			resolver R {
				x Int32
			}
			struct S {
				f ?[]R
			}`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"None": ErrCase{
			Src: `schema test
			struct S {
				f None
			}`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
	})
}

// TestTypeErr tests property type errors
func TestTypeErr(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"StructFieldOptionalChain": ErrCase{
			Src: `schema test
			struct S {
				optChain ??T
			}`,
			Errs: []ErrCode{parser.ErrTypeOptChain},
		},
		"StructFieldOptionalChain2": ErrCase{
			Src: `schema test
			struct S {
				optChain []?[]??T
			}`,
			Errs: []ErrCode{parser.ErrTypeOptChain},
		},
		"ResolverPropOptionalChain": ErrCase{
			Src: `schema test
			resolver R {
				optChain ??T
			}`,
			Errs: []ErrCode{parser.ErrTypeOptChain},
		},
		"ResolverPropOptionalChain2": ErrCase{
			Src: `schema test
			resolver R {
				optChain []?[]??T
			}`,
			Errs: []ErrCode{parser.ErrTypeOptChain},
		},
		"QueryUndefinedType": ErrCase{
			Src: `schema test
			query q Undefined`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"QueryParamUndefinedType": ErrCase{
			Src: `schema test
			query q(x Undefined) String`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"MutationUndefinedType": ErrCase{
			Src: `schema test
			mutation m Undefined`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"MutationParamUndefinedType": ErrCase{
			Src: `schema test
			mutation m(x Undefined) String`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
	})
}
