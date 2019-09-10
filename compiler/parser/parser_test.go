package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/romshark/gapi/compiler/parser"
	"github.com/romshark/gapi/internal/intset"
	"github.com/stretchr/testify/require"
)

type SchemaModel = *parser.SchemaModel
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
	modInspector func(SchemaModel),
) {
	// Initialize parser
	pr, err := parser.NewParser()
	require.NoError(t, err)
	require.NotNil(t, pr)

	// Compile
	mainFrag, err := pr.Parse(src(source))
	require.NoError(t, err)
	require.NotNil(t, mainFrag)
	mod := pr.SchemaModel()
	require.NotNil(t, mod)

	verifyModel(t, mod)

	// Inspect SchemaModel
	modInspector(mod)
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
			mainFrag, err := pr.Parse(src(errCase.Src))
			require.Error(t, err)
			require.Nil(t, mainFrag)
			require.IsType(t, parser.ParseErr{}, err)
			actualErrs := err.(parser.ParseErr).Errors
			require.True(t, len(actualErrs) > 0)
			require.Nil(t, pr.SchemaModel())

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

func verifyModel(t *testing.T, mod SchemaModel) {
	typeIDs := intset.NewIntSet()
	graphNodeIDs := intset.NewIntSet()
	paramIDs := intset.NewIntSet()

	// Ensure type ID uniqueness
	for _, tp := range mod.Types {
		id := tp.TypeID()
		require.NotEqual(t, parser.TypeIDUserTypeOffset, id)
		require.False(t, typeIDs.Has(int(id)))
		typeIDs.Insert(int(id))

		// Ensure correct type ID mapping
		require.Equal(t, tp, mod.FindTypeByID(id))
	}

	// Ensure graph node ID uniqueness
	for _, str := range mod.StructTypes {
		for _, fld := range str.(*parser.TypeStruct).Fields {
			intID := int(fld.GraphNodeID())
			require.False(t, graphNodeIDs.Has(intID))
			graphNodeIDs.Insert(intID)
		}
	}
	for _, rsv := range mod.ResolverTypes {
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
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalName2": ErrCase{
			Src: `schema illegal_Name
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalName3": ErrCase{
			Src: `schema IllegalName
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
	})
}

// TestDeclTypeErrs tests type declaration errors (generic errors)
func TestDeclTypeErrs(t *testing.T) {
	testCases := map[string]ErrCase{
		"IllegalName": ErrCase{
			Src: `schema test
			_illegalName = enum { e }
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalName2": ErrCase{
			Src: `schema test
			illegal_Name = enum { e }
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalName3": ErrCase{
			Src: `schema test
			Illegal_Name = enum { e }
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"RedeclUserType": ErrCase{
			Src: `schema test
			X = enum { a b }
			X = String
			q = query Bool`,
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
				%s = enum { e }
				q = query Bool`,
				primTypeName,
			),
			Errs: []ErrCode{parser.ErrTypeRedecl},
		}
	}

	testErrs(t, testCases)
}

// TestModAliases tests alias type declaration in SchemaModel
func TestModAliases(t *testing.T) {
	src := `schema test
	
	A1 = String
	A2 = Uint32
	A3 = A1
	q = query Bool
	`

	test(t, src, func(mod SchemaModel) {
		require.Len(t, mod.QueryEndpoints, 1)
		require.Len(t, mod.Mutations, 0)
		require.Len(t, mod.Types, 3)
	})
}

// TestDeclAliasTypeErrs tests alias type declaration errors
func TestDeclAliasTypeErrs(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"IllegalTypeName": ErrCase{
			Src: `schema test
			illegalName = String
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			Illegal_Name = String
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			_IllegalName = String
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalAliasedTypeName": ErrCase{
			Src: `schema test
			A = illegalName
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalAliasedTypeName2": ErrCase{
			Src: `schema test
			A = Illegal_Name
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalAliasedTypeName3": ErrCase{
			Src: `schema test
			A = _IllegalName
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"UndefinedAliasedType": ErrCase{
			Src: `schema test
			A = Undefined
			q = query Bool`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"DirectAliasCycle": ErrCase{
			Src: `schema test
			A = A
			q = query Bool`,
			Errs: []ErrCode{parser.ErrAliasRecurs},
		},
		"IndirectAliasCycle1": ErrCase{
			Src: `schema test
			A = B
			B = A
			q = query Bool`,
			Errs: []ErrCode{parser.ErrAliasRecurs},
		},
		"IndirectAliasCycle2": ErrCase{
			Src: `schema test
			G = H
			H = String
			F = C
			A = B
			B = C
			C = D
			D = A
			q = query Bool`,
			Errs: []ErrCode{parser.ErrAliasRecurs},
		},
		"MultipleIndirectAliasesCycles": ErrCase{
			Src: `schema test
			A = A
			B = C
			C = D
			D = B
			H = K
			K = I
			I = K
			q = query Bool`,
			Errs: []ErrCode{
				parser.ErrAliasRecurs,
				parser.ErrAliasRecurs,
				parser.ErrAliasRecurs,
			},
		},
	})
}

// TestModEnums tests enum type declaration in SchemaModel
func TestModEnums(t *testing.T) {
	src := `schema test
	
	E1 = enum {
		oneVal
	}
	E2 = enum {
		foo
		bar
	}
	E3 = enum {
		foo1
		bar2
		baz3
	}
	q = query Bool`

	test(t, src, func(mod SchemaModel) {
		require.Len(t, mod.QueryEndpoints, 1)
		require.Len(t, mod.Mutations, 0)
		require.Len(t, mod.Types, 3)

		e1 := mod.EnumTypes[0]
		e2 := mod.EnumTypes[1]
		e3 := mod.EnumTypes[2]

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
			require.Equal(t, expec.Name, expec.Type.String())
			require.NotNil(t, mod.FindTypeByDesignation(expec.Name))
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
			illegalName = enum {
				foo
				bar
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			_IllegalName = enum {
				foo
				bar
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			Illegal_Name = enum {
				foo
				bar
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"RedundantValue": ErrCase{
			Src: `schema test
			E = enum {
				foo
				foo
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrEnumValRedecl},
		},
		"NoValues": ErrCase{
			Src: `schema test
			E = enum {}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrEnumNoVal},
		},
		"IllegalValueIdentifier": ErrCase{
			Src: `schema test
			E = enum {
				_foo
				_bar
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalValueIdentifier2": ErrCase{
			Src: `schema test
			E = enum {
				1foo
				2bar
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalValueIdentifier3": ErrCase{
			Src: `schema test
			E = enum {
				fo_o
				ba_r
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
	})
}

// TestModUnions tests union type declarations in SchemaModel
func TestModUnions(t *testing.T) {
	src := `schema test
	
	U1 = union {
		String
		Uint32
	}
	U2 = union {
		Uint32
		Float64
		String
	}
	U3 = union {
		String
		Float64
		Int32
		Int64
	}
	query q Bool`

	test(t, src, func(mod SchemaModel) {
		require.Len(t, mod.QueryEndpoints, 1)
		require.Len(t, mod.Mutations, 0)
		require.Len(t, mod.Types, 3)

		u1 := mod.UnionTypes[0]
		u2 := mod.UnionTypes[1]
		u3 := mod.UnionTypes[2]

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
			require.Equal(t, expec.Name, expec.Type.String())
			require.IsType(t, &parser.TypeUnion{}, expec.Type)
			tpe := expec.Type.(*parser.TypeUnion)

			containsType := func(
				expected string,
				types []parser.Type,
			) bool {
				for _, tp := range types {
					if tp.String() == expected {
						return true
					}
				}
				return false
			}

			for i, referencedType := range expec.Types {
				require.True(t, containsType(
					referencedType.String(),
					tpe.Types,
				))
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
			illegalName = union {
				String
				Int32
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			_IllegalName = union {
				String
				Int32
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			Illegal_Name = union {
				String
				Int32
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"OneTypeUnion": ErrCase{
			Src: `schema test
			U = union {
				String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrUnionMissingOpts},
		},
		"RedundantOptionType": ErrCase{
			Src: `schema test
			U = union {
				String
				String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrUnionRedund},
		},
		"UndefinedType": ErrCase{
			Src: `schema test
			U = union {
				String
				Undefined
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"SelfReference": ErrCase{
			Src: `schema test
			U = union {
				Int32
				U
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrUnionRecurs},
		},
		"NonTypeElements": ErrCase{
			Src: `schema test
			U = union {
				foo
				bar
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"NonTypeElements2": ErrCase{
			Src: `schema test
			U = union {
				_foo
				_bar
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IncludesNone": ErrCase{
			Src: `schema test
			U = union {
				Int32
				None
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrUnionIncludesNone},
		},
	})
}

// TestModStructs tests struct type declarations in SchemaModel
func TestModStructs(t *testing.T) {
	src := `schema test
	S1 = struct {
		x String
	}
	S2 = struct {
		x []S2
		y ?S2
	}
	S3 = struct {
		optional ?String
		list []Float64
		matrix [][]Int64
		matrix3D [][][]Int64
		optionalList ?[]Int32
		listOfOptionals []?Int32
		optionalListOfOptionals ?[]?Int32
		optionalListOfOptionalListsOfOptionals ?[]?[]?String
	}
	q = query Bool`

	test(t, src, func(mod SchemaModel) {
		require.Len(t, mod.QueryEndpoints, 1)
		require.Len(t, mod.Mutations, 0)
		require.Len(t, mod.Types, 3+10)

		s1 := mod.StructTypes[0]
		s2 := mod.StructTypes[1]
		s3 := mod.StructTypes[2]

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
			require.Equal(t, expec.Name, expec.Type.String())
			require.IsType(t, &parser.TypeStruct{}, expec.Type)
			structType := expec.Type.(*parser.TypeStruct)

			// Make sure fields match the expectations
			for i, field := range expec.Fields {
				actualField := structType.Fields[i]
				require.Equal(t, field.Name, actualField.Name)
				require.Equal(t, field.Type.String(), actualField.Type.String())
				require.Equal(t, field.GraphID, actualField.GraphID)
				require.Equal(t, structType, expec.Type)

				// Make sure graph node IDs are unique
				require.NotContains(t, graphNodes, actualField.GraphID)
				graphNodes[actualField.GraphID] = actualField
			}
		}

		// Make sure the graph nodes are registered correctly
		require.Len(t, mod.GraphNodes, len(graphNodes)+1)
		for id, field := range graphNodes {
			node := mod.FindGraphNodeByID(id)
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
			illegalName = struct {
				foo String
				bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			_IllegalName = struct {
				foo String
				bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			Illegal_Name = struct {
				foo String
				bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"NoFields": ErrCase{
			Src: `schema test
			S = struct {}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructNoFields},
		},
		"RedundantField": ErrCase{
			Src: `schema test
			S = struct {
				foo String
				foo String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldRedecl},
		},
		"IllegalFieldIdentifier": ErrCase{
			Src: `schema test
			S = struct {
				_foo String
				_bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalFieldIdentifier2": ErrCase{
			Src: `schema test
			S = struct {
				1foo String
				2bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalFieldIdentifier3": ErrCase{
			Src: `schema test
			S = struct {
				fo_o String
				ba_r String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"RecursDirect": ErrCase{
			Src: `schema test
			S = struct {
				s S
			}
			q = query Bool`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // S.s -> S
			},
		},
		"RecursIndirect": ErrCase{
			Src: `schema test
			X = struct {
				s S
			}
			S = struct {
				x X
			}
			q = query Bool`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // X.s -> S.x -> X
			},
		},
		"RecursIndirect2": ErrCase{
			Src: `schema test
			Y = struct {
				s S
			}
			X = struct {
				y Y
			}
			S = struct {
				x X
			}
			q = query Bool`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // Y.s -> S.x -> X.y -> Y
			},
		},
		"RecursIndirect3": ErrCase{
			Src: `schema test
			Y = struct {
				s S
				z S
			}
			X = struct {
				y Y
				s S
			}
			S = struct {
				x X
				y Y
			}
			q = query Bool`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // S.x -> X.y -> Y.s -> S
				parser.ErrStructRecurs, // S.y -> Y.z -> S
			},
		},
		"RecursMultiple": ErrCase{
			Src: `schema test
			A = struct {
				a A
			}
			B = struct {
				b B
			}
			X = struct {
				y Y
			}
			Y = struct {
				x X
			}
			q = query Bool`,
			Errs: []ErrCode{
				parser.ErrStructRecurs, // A.a -> A
				parser.ErrStructRecurs, // B.b -> B
				parser.ErrStructRecurs, // X.y -> Y.x -> X
			},
		},
	})
}

// TestModResolvers tests resolver type declarations in SchemaModel
func TestModResolvers(t *testing.T) {
	src := `schema test
	R1 = resolver {
		x String
	}
	R2 = resolver {
		r R1
		x R2
		y []R2
		z ?R3
	}
	R3 = resolver {
		optional ?String
		list []Float64
		matrix [][]Int64
		matrix3D [][][]Int64
		optionalList ?[]Int32
		listOfOptionals []?Int32
		optionalListOfOptionals ?[]?Int32
		optionalListOfOptionalListsOfOptionals ?[]?[]?String
	}
	R4 = resolver {
		x(x Int32) Int32
		y(x Int32, y ?String, z ?[]Bool) String
	}
	q = query Bool`

	test(t, src, func(mod SchemaModel) {
		require.Len(t, mod.QueryEndpoints, 1)
		require.Len(t, mod.Mutations, 0)
		require.Len(t, mod.Types, 4+11)

		r1 := mod.ResolverTypes[0]
		r2 := mod.ResolverTypes[1]
		r3 := mod.ResolverTypes[2]
		r4 := mod.ResolverTypes[3]

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
			require.Equal(t, expec.Name, expec.Type.String())
			require.IsType(t, &parser.TypeResolver{}, expec.Type)
			resolverType := expec.Type.(*parser.TypeResolver)

			// Make sure properties match expectations
			require.Len(t, resolverType.Properties, len(expec.Props))
			for i, prop := range expec.Props {
				actualProp := resolverType.Properties[i]
				require.Equal(t, prop.Name, actualProp.Name)
				require.Equal(
					t,
					prop.Type.String(),
					actualProp.Type.String(),
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
						param.Type.String(),
						actualParam.Type.String(),
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
		require.Len(t, mod.GraphNodes, len(graphNodes)+1)
		for id, prop := range graphNodes {

			node := mod.FindGraphNodeByID(id)
			require.NotNil(t, node, "graph node (%d) not found in SchemaModel", id)
			require.Equal(t, id, node.GraphNodeID())
			require.Equal(t, prop.Resolver, node.Parent())
		}

		// Make sure parameters are registered correctly
		for id, p := range parameters {
			param := mod.FindParameterByID(id)
			require.NotNil(t, param, "parameter (%d) not found in SchemaModel", id)
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
			illegalName = resolver {
				foo String
				bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			_IllegalName = resolver {
				foo String
				bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			Illegal_Name = resolver {
				foo String
				bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"NoProps": ErrCase{
			Src: `schema test
			S = resolver {}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrResolverNoProps},
		},
		"RedundantProp": ErrCase{
			Src: `schema test
			S = resolver {
				foo String
				foo String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrResolverPropRedecl},
		},
		"IllegalPropIdentifier": ErrCase{
			Src: `schema test
			S = resolver {
				_foo String
				_bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropIdentifier2": ErrCase{
			Src: `schema test
			S = resolver {
				1foo String
				2bar String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropIdentifier3": ErrCase{
			Src: `schema test
			S = resolver {
				fo_o String
				ba_r String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropParamIdentifier": ErrCase{
			Src: `schema test
			S = resolver {
				foo(_foo String) String
				bar(_bar String) String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropParamIdentifier2": ErrCase{
			Src: `schema test
			S = resolver {
				foo(1foo String) String
				bar(2bar String) String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalPropParamIdentifier3": ErrCase{
			Src: `schema test
			S = resolver {
				foo(fo_o String) String
				bar(ba_r String) String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"RedundantPropParam": ErrCase{
			Src: `schema test
			S = resolver {
				foo(foo String, foo Int32) String
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrParamRedecl},
		},
	})
}

// TestModQueries tests query declarations in SchemaModel
func TestModQueries(t *testing.T) {
	src := `schema test
	Foo = struct {
		foo String
	}
	Bar = resolver {
		bar String
	}
	foo = query Foo
	bar = query Bar
	str = query String
	foo2 = query(foo Foo) Foo
	bar2 = query(bar Int32, baz Float64) Bar
	baz = query(
		first Int32,
		second Bool,
		third Uint64,
	) String
	`

	test(t, src, func(mod SchemaModel) {
		require.Len(t, mod.Types, 2)
		require.Len(t, mod.Mutations, 0)
		require.Len(t, mod.QueryEndpoints, 6)

		require.Len(t, mod.StructTypes, 1)
		tFoo := mod.StructTypes[0]

		require.Len(t, mod.ResolverTypes, 1)
		tBar := mod.ResolverTypes[0]

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
		require.Len(t, mod.QueryEndpoints, len(expected))
		for i1, expec := range expected {
			require.IsType(t, parser.Query{}, expec)

			actual := mod.QueryEndpoints[i1]
			require.Equal(t, expec.Name, actual.Name)
			require.Equal(t, expec.GraphID, actual.GraphID)
			require.Equal(t, expec.Type, actual.Type)

			// Make sure the graph nodes are registered correctly
			foundNode := mod.FindGraphNodeByID(expec.GraphID)
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
				regParam := mod.FindParameterByID(param.ID)
				require.Equal(t, actualParam, regParam)
			}
		}
	})
}

// TestModMutations tests mutation declarations in SchemaModel
func TestModMutations(t *testing.T) {
	src := `schema test
	Foo = struct {
		foo String
	}
	Bar = resolver {
		bar String
	}
	foo = mutation Foo
	bar = mutation Bar
	str = mutation String
	foo2 = mutation(foo Foo) Foo
	bar2 = mutation(bar Int32, baz Float64) Bar
	baz = mutation(
		first Int32,
		second Bool,
		third Uint64,
	) String
	`

	test(t, src, func(mod SchemaModel) {
		require.Len(t, mod.Types, 2)
		require.Len(t, mod.Mutations, 6)
		require.Len(t, mod.QueryEndpoints, 0)

		require.Len(t, mod.StructTypes, 1)
		tFoo := mod.StructTypes[0]

		require.Len(t, mod.ResolverTypes, 1)
		tBar := mod.ResolverTypes[0]

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
		require.Len(t, mod.Mutations, len(expected))
		for i1, expec := range expected {
			require.IsType(t, parser.Mutation{}, expec)

			actual := mod.Mutations[i1]
			require.Equal(t, expec.Name, actual.Name)
			require.Equal(t, expec.GraphID, actual.GraphID)
			require.Equal(t, expec.Type, actual.Type)

			// Make sure the graph nodes are registered correctly
			foundNode := mod.FindGraphNodeByID(expec.GraphID)
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
				regParam := mod.FindParameterByID(param.ID)
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
			IllegalName = query String
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName2": ErrCase{
			Src: `schema test
			_illegalName = query String
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
		"IllegalTypeName3": ErrCase{
			Src: `schema test
			illegal_Name = query String
			q = query Bool`,
			Errs: []ErrCode{parser.ErrSyntax},
		},
	})
}

// TestParamImpureType tests specifying parameters of non-pure types
func TestParamImpureType(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"Resolver": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			q = query(param R) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureUnion": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			U = union {
				R
				Int32
			}
			q = query(param U) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureAliasToResolver": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			A = R
			q = query(param A) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureUnionOfImpureAlias": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			A = R
			U = union {
				Int32
				A
			}
			q = query(param U) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureOptional": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			q = query(param ?R) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureList": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			q = query(param []R) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"ImpureOptionalList": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			q = query(param ?[]R) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
		"None": ErrCase{
			Src: `schema test
			q = query(param None) String`,
			Errs: []ErrCode{parser.ErrParamImpure},
		},
	})
}

// TestStructImpureFieldType tests specifying struct fields of non-pure types
func TestStructImpureFieldType(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"Resolver": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			S = struct {
				f R
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureUnion": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			U = union {
				R
				Int32
			}
			S = struct {
				f U
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureAliasToResolver": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			A = R
			S = struct {
				f A
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureUnionOfImpureAlias": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			A = R
			U = union {
				Int32
				A
			}
			S = struct {
				f U
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureOptional": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			S = struct {
				f ?R
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureList": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			S = struct {
				f []R
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"ImpureOptionalList": ErrCase{
			Src: `schema test
			R = resolver {
				x Int32
			}
			S = struct {
				f ?[]R
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
		"None": ErrCase{
			Src: `schema test
			S = struct {
				f None
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrStructFieldImpure},
		},
	})
}

// TestTypeErr tests property type errors
func TestTypeErr(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"StructFieldOptionalChain": ErrCase{
			Src: `schema test
			S = struct {
				optChain ??T
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrTypeOptChain},
		},
		"StructFieldOptionalChain2": ErrCase{
			Src: `schema test
			S = struct {
				optChain []?[]??T
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrTypeOptChain},
		},
		"ResolverPropOptionalChain": ErrCase{
			Src: `schema test
			R = resolver {
				optChain ??T
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrTypeOptChain},
		},
		"ResolverPropOptionalChain2": ErrCase{
			Src: `schema test
			R = resolver {
				optChain []?[]??T
			}
			q = query Bool`,
			Errs: []ErrCode{parser.ErrTypeOptChain},
		},
		"QueryUndefinedType": ErrCase{
			Src: `schema test
			q = query Undefined`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"QueryParamUndefinedType": ErrCase{
			Src: `schema test
			q = query(x Undefined) String`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"MutationUndefinedType": ErrCase{
			Src: `schema test
			m = mutation Undefined`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
		"MutationParamUndefinedType": ErrCase{
			Src: `schema test
			m = mutation(x Undefined) String`,
			Errs: []ErrCode{parser.ErrTypeUndef},
		},
	})
}

// TestNoEndpoints expects the parser to fail in case
// no query and mutation endpoints are declared
func TestNoEndpoints(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"StructFieldOptionalChain": ErrCase{
			Src: `schema test
			Foo = struct {
				foo String
			}
			Bar = resolver {
				bar String
			}`,
			Errs: []ErrCode{parser.ErrNoEndpoints},
		},
	})
}

// TestIllegalNoneTypes expects the parser to fail in case
// of illegal None-types
func TestIllegalNoneTypes(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"NoneQueryEndpoints": ErrCase{
			Src: `schema test
			q1 = query ?None
			q2 = query []None
			q3 = query []?None
			q4 = query ?[]?None`,
			Errs: []ErrCode{
				parser.ErrSyntax,
				parser.ErrSyntax,
				parser.ErrSyntax,
				parser.ErrSyntax,
			},
		},
	})
}

// TestGraphRootNodeRedecl expects the parser to fail in case
// of a name conflict between multiple graph root nodes
func TestGraphRootNodeRedecl(t *testing.T) {
	testErrs(t, map[string]ErrCase{
		"Queries": ErrCase{
			Src: `schema test
			q = query String
			q = query Int32`,
			Errs: []ErrCode{parser.ErrGraphRootNodeRedecl},
		},
		"Mutations": ErrCase{
			Src: `schema test
			m = mutation String
			m = mutation Int32`,
			Errs: []ErrCode{parser.ErrGraphRootNodeRedecl},
		},
		"QueryMutation": ErrCase{
			Src: `schema test
			q = query String
			q = mutation Int32`,
			Errs: []ErrCode{parser.ErrGraphRootNodeRedecl},
		},
	})
}
