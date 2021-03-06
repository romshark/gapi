package parser

import (
	"github.com/romshark/gapi/internal/intset"
)

type aliasTypeCycle struct {
	nodes []*TypeAlias
}

// String stringifies the cycle
func (c aliasTypeCycle) String() string {
	first := c.nodes[0]
	if len(c.nodes) < 2 {
		return first.Name + " -> " + first.Name
	}

	s := first.Name
	for _, n := range c.nodes[1:] {
		s += " -> " + n.Name
	}
	return s + " -> " + first.Name
}

// findAliasTypeCycles returns all recursive alias type cycles
// or nil if there are none
func (pr *Parser) findAliasTypeCycles() (cycles []aliasTypeCycle) {
	// Keeps track of all nodes being part of a cycle
	// so they don't need to be checked repeatedly
	cycleReg := intset.NewIntSet()

	// Remember the nodes to be checked
	toBeChecked := intset.NewIntSet()
	for _, n := range pr.mod.AliasTypes {
		toBeChecked.Insert(int(n.TypeID()))
	}

	// Check all nodes to be checked
	for {
		// Get any node that's still to be checked
		nodei := pr.typeByID[TypeID(toBeChecked.Take())]
		if nodei == nil {
			break
		}
		node := nodei.(*TypeAlias)

		// chain keeps track of the order of nodes in the current chain
		chain := []*TypeAlias{node}

		// chainReg keeps track of all nodes in the current chain
		chainReg := intset.NewIntSet()
		chainReg.Insert(int(node.TypeID()))

		// Traverse the path until there's no more aliased type
		next := node.AliasedType
		for {
			_, isAlias := next.(*TypeAlias)
			if next == nil || !isAlias {
				break
			}

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

				cycle := aliasTypeCycle{nodes: chain}
				// Mark all nodes of the cycle as cyclic
				// so they don't have to be checked later
				for _, n := range cycle.nodes {
					cycleReg.Insert(int(n.TypeID()))
				}
				cycles = append(cycles, cycle)
				break
			}

			// Continue to traverse the path and update the chain
			chain = append(chain, next.(*TypeAlias))
			chainReg.Insert(int(next.TypeID()))

			next = next.(*TypeAlias).AliasedType
		}
	}

	return
}
