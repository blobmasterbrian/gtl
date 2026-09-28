// Package set provides an unordered collection of unique values.
package set

import "github.com/blobmasterbrian/gtl/algorithm"

// Set is an unordered collection of unique values. The zero value is an empty
// set ready to use.
//
// A Set is backed by a Go map, so single-value operations are expected O(1),
// amortized over growth. Complexity notes use n for the number of values in
// the input and k for the number in the result.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New returns a set containing items.
//
// Time O(n), space O(n), sized up front.
func New[T comparable](items ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(items))}
	for _, v := range items {
		s.m[v] = struct{}{}
	}
	return s
}

// Collect returns a set containing the elements of src, which may be an
// iter.Seq or an algorithm.Stream.
//
// Time O(n) expected, space O(n). A sequence carries no length, so the set
// grows as elements arrive.
func Collect[S algorithm.Seq[T], T comparable](src S) *Set[T] {
	s := New[T]()
	s.Insert(src)
	return s
}
