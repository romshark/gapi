package parser

import (
	"encoding/json"
)

// JSONSchemaModel represents the schema JSON model
type JSONSchemaModel struct {
	SchemaName     string                   `json:"schema-name"`
	AliasTypes     []JSONModelAliasType     `json:"alias-types"`
	EnumTypes      []JSONModelEnumType      `json:"enum-types"`
	UnionTypes     []JSONModelUnionType     `json:"union-types"`
	StructTypes    []JSONModelStructType    `json:"struct-types"`
	ResolverTypes  []JSONModelResolverType  `json:"resolver-types"`
	AnonymousTypes []JSONModelAnonymousType `json:"anonymous-types"`
	QueryEndpoints []JSONModelQueryEndpoint `json:"query-endpoints"`
	Mutations      []JSONModelMutation      `json:"mutations"`
}

// JSONModelAliasType represents the JSON model of an alias type
type JSONModelAliasType struct {
	Name          string `json:"name"`
	ID            int    `json:"id"`
	AliasedTypeID int    `json:"aliased-type-id"`
}

// JSONModelEnumType represents the JSON model of an enum type
type JSONModelEnumType struct {
	Name   string   `json:"name"`
	ID     int      `json:"id"`
	Values []string `json:"values"`
}

// JSONModelUnionType represents the JSON model of a union type
type JSONModelUnionType struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	OptionTypes []int  `json:"option-types"`
}

// JSONModelStructField represents the JSON model of a struct field
type JSONModelStructField struct {
	Name        string `json:"name"`
	Type        int    `json:"type"`
	GraphNodeID int    `json:"graph-node-id"`
}

// JSONModelStructType represents the JSON model of a struct type
type JSONModelStructType struct {
	Name   string                 `json:"name"`
	ID     int                    `json:"id"`
	Fields []JSONModelStructField `json:"fields"`
}

// JSONModelParameter represents the JSON model of a parameter
type JSONModelParameter struct {
	Name         string `json:"name"`
	Type         int    `json:"type"`
	GraphParamID int    `json:"graph-param-id"`
}

// JSONModelResolverProperty represents the JSON model of a resolver property
type JSONModelResolverProperty struct {
	Name        string               `json:"name"`
	Type        int                  `json:"type"`
	GraphNodeID int                  `json:"graph-node-id"`
	Parameters  []JSONModelParameter `json:"parameters"`
}

// JSONModelResolverType represents the JSON model of a resolver type
type JSONModelResolverType struct {
	Name       string                      `json:"name"`
	ID         int                         `json:"id"`
	Properties []JSONModelResolverProperty `json:"properties"`
}

// JSONModelAnonymousType represents the JSON model of an anonymous type
type JSONModelAnonymousType struct {
	Designation string `json:"designation"`
	ID          int    `json:"id"`
}

// JSONModelQueryEndpoint represents the JSON model of a query endpoint
type JSONModelQueryEndpoint struct {
	Name        string               `json:"name"`
	Type        int                  `json:"type"`
	GraphNodeID int                  `json:"graph-node-id"`
	Parameters  []JSONModelParameter `json:"parameters"`
}

// JSONModelMutation represents the JSON model of a mutation
type JSONModelMutation struct {
	Name        string               `json:"name"`
	Type        int                  `json:"type"`
	GraphNodeID int                  `json:"graph-node-id"`
	Parameters  []JSONModelParameter `json:"parameters"`
}

// MarshalJSON marshal the schema model into its JSON representation
func (mod *SchemaModel) MarshalJSON() ([]byte, error) {
	copyParams := func(ps []*Parameter) []JSONModelParameter {
		v := make([]JSONModelParameter, len(ps))
		for i, p := range ps {
			v[i] = JSONModelParameter{
				Name:         p.Name,
				Type:         int(p.Type.TypeID()),
				GraphParamID: int(p.ID),
			}
		}
		return v
	}

	model := &JSONSchemaModel{
		SchemaName:     mod.SchemaName,
		AliasTypes:     make([]JSONModelAliasType, len(mod.AliasTypes)),
		EnumTypes:      make([]JSONModelEnumType, len(mod.EnumTypes)),
		UnionTypes:     make([]JSONModelUnionType, len(mod.UnionTypes)),
		StructTypes:    make([]JSONModelStructType, len(mod.StructTypes)),
		ResolverTypes:  make([]JSONModelResolverType, len(mod.ResolverTypes)),
		AnonymousTypes: make([]JSONModelAnonymousType, len(mod.AnonymousTypes)),
		QueryEndpoints: make([]JSONModelQueryEndpoint, len(mod.QueryEndpoints)),
		Mutations:      make([]JSONModelMutation, len(mod.Mutations)),
	}

	// Alias types
	for i, t := range mod.AliasTypes {
		v := t.(*TypeAlias)

		model.AliasTypes[i] = JSONModelAliasType{
			Name:          v.Name,
			ID:            int(v.ID),
			AliasedTypeID: int(v.AliasedType.TypeID()),
		}
	}

	// Enum types
	for i, t := range mod.EnumTypes {
		v := t.(*TypeEnum)

		// Values
		vals := make([]string, len(v.Values))
		for i, val := range v.Values {
			vals[i] = val.Name
		}

		model.EnumTypes[i] = JSONModelEnumType{
			Name:   v.Name,
			ID:     int(v.ID),
			Values: vals,
		}
	}

	// Union types
	for i, t := range mod.UnionTypes {
		v := t.(*TypeUnion)

		// Option types
		opts := make([]int, len(v.Types))
		for i, opt := range v.Types {
			opts[i] = int(opt.TypeID())
		}

		model.UnionTypes[i] = JSONModelUnionType{
			Name:        v.Name,
			ID:          int(v.ID),
			OptionTypes: opts,
		}
	}

	// Struct types
	for i, t := range mod.StructTypes {
		v := t.(*TypeStruct)

		// Fields
		fields := make([]JSONModelStructField, len(v.Fields))
		for i, fld := range v.Fields {
			fields[i] = JSONModelStructField{
				Name:        fld.Name,
				Type:        int(fld.Type.TypeID()),
				GraphNodeID: int(fld.GraphID),
			}
		}

		model.StructTypes[i] = JSONModelStructType{
			Name:   v.Name,
			ID:     int(v.ID),
			Fields: fields,
		}
	}

	// Resolver types
	for i, t := range mod.ResolverTypes {
		v := t.(*TypeResolver)

		// Properties
		props := make([]JSONModelResolverProperty, len(v.Properties))
		for i, fld := range v.Properties {
			props[i] = JSONModelResolverProperty{
				Name:        fld.Name,
				Type:        int(fld.Type.TypeID()),
				GraphNodeID: int(fld.GraphID),
				Parameters:  copyParams(fld.Parameters),
			}
		}

		model.ResolverTypes[i] = JSONModelResolverType{
			Name:       v.Name,
			ID:         int(v.ID),
			Properties: props,
		}
	}

	// Anonymous types
	for i, t := range mod.AnonymousTypes {
		model.AnonymousTypes[i] = JSONModelAnonymousType{
			Designation: t.String(),
			ID:          int(t.TypeID()),
		}
	}

	// Query endpoints
	for i, q := range mod.QueryEndpoints {
		model.QueryEndpoints[i] = JSONModelQueryEndpoint{
			Name:        q.Name,
			GraphNodeID: int(q.GraphID),
			Parameters:  copyParams(q.Parameters),
			Type:        int(q.Type.TypeID()),
		}
	}

	// Mutations
	for i, m := range mod.Mutations {
		model.Mutations[i] = JSONModelMutation{
			Name:        m.Name,
			GraphNodeID: int(m.GraphID),
			Parameters:  copyParams(m.Parameters),
			Type:        int(m.Type.TypeID()),
		}
	}

	return json.Marshal(model)
}
