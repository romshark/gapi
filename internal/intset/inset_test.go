package intset_test

import (
	"testing"

	"github.com/romshark/gapi/internal/intset"
	"github.com/stretchr/testify/require"
)

// TestInsert tests Insert, Has and Len
func TestInsert(t *testing.T) {
	set := intset.NewIntSet()
	sz := 10

	for i := 1; i <= sz; i++ {
		require.False(t, set.Has(i))
		require.True(t, set.Insert(i))
		require.True(t, set.Has(i))
	}
	require.Equal(t, sz, set.Len())
}

// TestRemove tests Remove, Has and Len
func TestRemove(t *testing.T) {
	set := intset.NewIntSet()
	sz := 10

	// Insert
	for i := 1; i <= sz; i++ {
		require.True(t, set.Insert(i))
	}
	require.Equal(t, sz, set.Len())

	// Remove
	for i := sz; i >= 1; i-- {
		require.True(t, set.Has(i))
		require.True(t, set.Remove(i))
		require.False(t, set.Has(i))
	}
	require.Equal(t, 0, set.Len())
}

// TestTake tests Take, Has and Len
func TestTake(t *testing.T) {
	set := intset.NewIntSet()
	sz := 10

	// Insert
	for i := 1; i <= sz; i++ {
		set.Insert(i)
	}
	require.Equal(t, sz, set.Len())

	// Take all
	for i := 1; i <= sz; i++ {
		require.True(t, set.Has(i))
		took := set.Take()
		require.Equal(t, i, took)
		require.False(t, set.Has(i))
		require.Equal(t, sz-i, set.Len())
	}
	require.Equal(t, 0, set.Len())
}
