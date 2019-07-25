package parser

import (
	"github.com/romshark/gapi/internal/intset"
)

type structTypeCycle struct {
	fields []GraphNode
}

// newStructTypeCycle creates a new cycle
func newStructTypeCycle(fields []GraphNode) structTypeCycle {
	return structTypeCycle{fields: fields}
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

// findStructTypeCycles returns all recursive struct type cycles
// or nil if there are none
func (pr *Parser) findStructTypeCycles() (cycles []structTypeCycle) {
	// Keeps track of all fields being part of a cycle
	// so they don't need to be checked repeatedly
	cycleReg := intset.NewIntSet()

	// Remember the fields to be checked
	toBeChecked := intset.NewIntSet()
	for _, n := range pr.ast.StructTypes {
		for _, f := range n.(*TypeStruct).Fields {
			if f.Type != nil && f.Type.Category() == TypeCategoryStruct {
				toBeChecked.Insert(int(f.GraphNodeID()))
			}
		}
	}

	// Check all fields to be checked
	for {
		// Get any field that's still to be checked
		intID := toBeChecked.Take()
		if intID < 1 {
			break
		}
		id := GraphNodeID(intID)
		rootField := pr.graphNodeByID[id].(*StructField)

		// Don't traverse fields that are already part of some cycles
		if cycleReg.Has(int(rootField.GraphNodeID())) {
			continue
		}

		// chain keeps track of the order of fields in the current chain
		chain := []GraphNode{rootField}

		// chainReg keeps track of all fields in the current chain
		chainReg := intset.NewIntSet()
		chainReg.Insert(int(rootField.Struct.ID))

		// Traverse all fields until there's no more fields left
		currentType := rootField.Type.(*TypeStruct)
		fieldIndex := 0
	FIELD_TRAVERSAL:
		for {
			// Select the next field of the current struct type
			if len(currentType.Fields) < fieldIndex+1 {
				break FIELD_TRAVERSAL
			}
			currentField := currentType.Fields[fieldIndex]

			// Don't traverse fields that are already part of some cycles
			if cycleReg.Has(int(currentField.GraphNodeID())) {
				break FIELD_TRAVERSAL
			}

			// Don't traverse this field any more
			toBeChecked.Remove(int(currentField.GraphNodeID()))

			// Ignore fields of non-struct type
			if currentField.Type.Category() != TypeCategoryStruct {
				break FIELD_TRAVERSAL
			}

			// Add this field to the chain
			chain = append(chain, currentField)

			fieldType := currentField.Type.(*TypeStruct)

			// Check if this field's type is already in the type chain registry
			if chainReg.Has(int(fieldType.ID)) {
				// Cycle detected, backtrack to the cycle start field
				for rev := len(chain) - 1; rev >= 0; rev-- {
					if chain[rev].Parent() == rootField.Parent() {
						chain = chain[rev:]
						break
					}
				}

				cycle := newStructTypeCycle(chain)
				// Mark all fields of the cycle as cyclic
				// so they don't have to be checked later
				for _, n := range cycle.fields {
					cycleReg.Insert(int(n.GraphNodeID()))
				}
				cycles = append(cycles, cycle)
				break FIELD_TRAVERSAL
			}

			// Add this field's type to the type chain registry
			chainReg.Insert(int(fieldType.ID))

			// Follow the path up the struct
			currentType = fieldType
			fieldIndex = 0
		}
	}

	return
}
