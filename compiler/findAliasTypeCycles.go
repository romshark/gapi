package compiler

import (
	"sort"

	"github.com/romshark/gapi/internal/intset"
)

type aliasTypeCycle struct {
	nodes []Type
}

// newAliasTypeCycle creates a new alphabetically sorted cycle
func newAliasTypeCycle(nodes []Type) aliasTypeCycle {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name() < nodes[j].Name()
	})
	return aliasTypeCycle{nodes: nodes}
}

// String stringifies the cycle
func (c aliasTypeCycle) String() string {
	first := c.nodes[0]
	if len(c.nodes) < 2 {
		return first.Name() + " -> " + first.Name()
	}

	s := first.Name()
	for _, n := range c.nodes[1:] {
		s += " -> " + n.Name()
	}
	return s + " -> " + first.Name()
}

// findAliasTypeCycles returns all recursive alias type cycles
// or nil if there are none
func (c *Compiler) findAliasTypeCycles() (cycles []aliasTypeCycle) {
	// Keeps track of all nodes being part of a cycle
	// so they don't need to be checked repeatedly
	cycleReg := intset.NewIntSet()

	// Remember the nodes to be checked
	toBeChecked := intset.NewIntSet()
	for _, n := range c.ast.AliasTypes {
		toBeChecked.Insert(int(n.TypeID()))
	}

	// Check all nodes to be checked
	for {
		// Get any node that's still to be checked
		node := c.typeByID[TypeID(toBeChecked.Take())]
		if node == nil {
			break
		}

		// chain keeps track of the order of nodes in the current chain
		chain := []Type{node}

		// chainReg keeps track of all nodes in the current chain
		chainReg := intset.NewIntSet()
		chainReg.Insert(int(node.TypeID()))

		// Traverse the path until there's no more aliased type
		next := node.(*TypeAlias).AliasedType
		for next != nil && next.Category() == TypeCategoryAlias {
			toBeChecked.Remove(int(next.TypeID()))
			if cycleReg.Has(int(next.TypeID())) {
				break
			}
			if chainReg.Has(int(next.TypeID())) {
				// Cycle detected, backtrack to the cycle start node
				for rev := len(chain) - 1; rev >= 0; rev-- {
					if chain[rev] == next {
						chain = chain[rev:]
						break
					}
				}

				cycle := newAliasTypeCycle(chain)
				// Mark all nodes of the cycle as cyclic
				// so they don't have to be checked later
				for _, n := range cycle.nodes {
					cycleReg.Insert(int(n.TypeID()))
				}
				cycles = append(cycles, cycle)
				break
			}

			// Continue to traverse the path and update the chain
			chain = append(chain, next)
			chainReg.Insert(int(next.TypeID()))

			next = next.(*TypeAlias).AliasedType
		}
	}

	return
}
