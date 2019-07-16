package compiler

import "sort"

func sortTypesByName(types []Type) {
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name() < types[j].Name()
	})
}
