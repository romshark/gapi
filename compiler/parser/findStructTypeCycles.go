package parser

import "github.com/romshark/gapi/internal/intset"

type structTypeCycle struct {
	fields []*StructField
}

// String stringifies the cycle
func (c structTypeCycle) String() string {
	first := c.fields[0]
	if len(c.fields) < 2 {
		return first.GraphNodeName() + " -> " + first.Parent().String()
	}

	s := first.GraphNodeName()
	for _, n := range c.fields[1:] {
		s += " -> " + n.GraphNodeName()
	}
	return s + " -> " + first.Parent().String()
}

func typeInChain(t *TypeStruct, typeChain []*TypeStruct) bool {
	for _, i := range typeChain {
		if i == t {
			return true
		}
	}
	return false
}

func findStructTypeCycle(
	cycleReg intset.IntSet,
	fieldChain []*StructField,
	typeChain []*TypeStruct,
	field *StructField,
) (cycle *structTypeCycle) {
	tp, ok := field.Type.(*TypeStruct)
	if !ok || cycleReg.Has(int(field.GraphID)) {
		return
	}

	fieldChain = append(fieldChain, field)

	// Check for indirect recursion
	if typeInChain(tp, typeChain) {
		if len(fieldChain) > 1 {
			// Backtrack the cycle start field
			for rev := len(fieldChain) - 1; rev >= 0; rev-- {
				if fieldChain[rev].Parent() == tp {
					fieldChain = fieldChain[rev:]
					break
				}
			}
		}

		cycle = &structTypeCycle{fields: fieldChain}

		// Mark all fields of the cycle as cyclic
		// so they don't have to be checked later
		for _, n := range cycle.fields {
			cycleReg.Insert(int(n.GraphID))
		}
		return
	}

	// Extend chain
	typeChain = append(typeChain, tp)

	// Traverse sub-fields (if any)
	if _, ok := field.Type.(*TypeStruct); ok {
		for _, subField := range tp.Fields {
			if _, isStruct := subField.Type.(*TypeStruct); isStruct {
				if cycle = findStructTypeCycle(
					cycleReg,
					fieldChain,
					typeChain,
					subField,
				); cycle != nil {
					return
				}
			}
		}
	}
	return
}

// findStructTypeCycles returns all recursive struct type cycles
// or nil if there are none
func (pr *Parser) findStructTypeCycles() (cycles []structTypeCycle) {
	cycleReg := intset.NewIntSet()

	// Collect all fields of struct type
	for _, tp := range pr.ast.StructTypes {
		for _, fd := range tp.(*TypeStruct).Fields {
			typeChainReg := intset.NewIntSet()
			typeChainReg.Insert(int(fd.Parent().TypeID()))
			if cy := findStructTypeCycle(
				cycleReg,
				nil,
				[]*TypeStruct{tp.(*TypeStruct)},
				fd,
			); cy != nil {
				cycles = append(cycles, *cy)
			}
		}
	}
	return
}
