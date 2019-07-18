package intset

import "golang.org/x/tools/container/intsets"

// IntSet implements a sparse integer set
type IntSet struct {
	set *intsets.Sparse
}

// NewIntSet creates a new sparse integer set instance
func NewIntSet() IntSet { return IntSet{&intsets.Sparse{}} }

// Insert inserts v into the set and reports whether the set grew
func (t IntSet) Insert(v int) bool { return t.set.Insert(v) }

// Remove removes v from the set and reports whether the set shrank
func (t IntSet) Remove(v int) bool { return t.set.Remove(v) }

// Has returns true if v is in the set, otherwise returns false
func (t IntSet) Has(v int) bool { return t.set.Has(v) }

// Take takes the smallest integer from the set removing it
func (t IntSet) Take() int {
	var min int
	if t.set.TakeMin(&min) {
		return min
	}
	return -1
}
