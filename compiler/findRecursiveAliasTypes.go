package compiler

import "sort"

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
func (ast *AST) findAliasTypeCycles() (cycles []aliasTypeCycle) {
	// Keeps track of all nodes being part of a cycle
	// so they don't need to be checked repeatedly
	cycleReg := make(map[Type]struct{})

	// Remember the nodes to be checked
	toBeChecked := make(map[Type]struct{}, len(ast.AliasTypes))
	for _, n := range ast.AliasTypes {
		toBeChecked[n] = struct{}{}
	}

	// Check all nodes to be checked
	for len(toBeChecked) > 0 {
		// Get any node that's still to be checked
		var node Type
		for node = range toBeChecked {
			break
		}
		delete(toBeChecked, node)

		// chain keeps track of the order of nodes in the current chain
		chain := []Type{node}

		// chainReg keeps track of all nodes in the current chain
		chainReg := map[Type]struct{}{node: struct{}{}}

		// Traverse the path until there's no more aliased type
		next := node.(*TypeAlias).AliasedType
		for next != nil && next.Category() == TypeCategoryAlias {
			delete(toBeChecked, next)
			if _, inCycleReg := cycleReg[next]; inCycleReg {
				break
			}
			if _, recursive := chainReg[next]; recursive {
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
					cycleReg[n] = struct{}{}
				}
				cycles = append(cycles, cycle)
				break
			}

			// Continue to traverse the path and update the chain
			chain = append(chain, next)
			chainReg[next] = struct{}{}

			next = next.(*TypeAlias).AliasedType
		}
	}

	return
}
