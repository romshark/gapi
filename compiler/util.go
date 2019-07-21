package compiler

import (
	"sort"
)

func sortTypesByName(types []Type) {
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name() < types[j].Name()
	})
}

func sortQueryEndpointsByName(endpoints []*QueryEndpoint) {
	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].GraphNodeName() < endpoints[j].GraphNodeName()
	})
}

func sortMutationsByName(mutations []*Mutation) {
	sort.Slice(mutations, func(i, j int) bool {
		return mutations[i].GraphNodeName() < mutations[j].GraphNodeName()
	})
}

func stringifyType(t Type) (name string) {
	if t == nil {
		return
	}
	for {
		if v, isOptional := t.(*TypeOptional); isOptional {
			name += "?"
			t = v.StoreType
			continue
		}
		if v, isList := t.(*TypeList); isList {
			name += "[]"
			t = v.StoreType
			continue
		}
		name += t.String()
		break
	}
	return
}

// skipUntil skips all tokens beginning at the start node until the expected
// rule is reached and returns the matching node
func skipUntil(start *node32, expectedRule pegRule) *node32 {
	for current := start; current != nil; {
		if current.pegRule == expectedRule {
			return current
		}
		current = current.next
	}
	return nil
}
