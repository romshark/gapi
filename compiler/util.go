package compiler

import (
	"sort"
)

func sortTypesByName(types []Type) {
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name() < types[j].Name()
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
