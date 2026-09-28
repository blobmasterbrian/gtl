// Mutations: every method here modifies the set in place.

package set

import "github.com/blobmasterbrian/gtl/algorithm"

// Add adds v to the set.
//
// Time O(1) expected, amortized over growth.
func (s *Set[T]) Add(v T) {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	s.m[v] = struct{}{}
}

// Remove removes v from the set. Removing an absent value does nothing.
//
// Time O(1) expected.
func (s *Set[T]) Remove(v T) {
	delete(s.m, v)
}

// Insert adds every element of src, which may be an iter.Seq or an
// algorithm.Stream, to the set.
//
// Time O(n) expected in the elements of src, space O(n).
func (s *Set[T]) Insert[S algorithm.Seq[T]](src S) {
	for v := range src {
		s.Add(v)
	}
}

// Retain keeps the values for which pred returns true and removes the rest.
//
// Time O(n) expected, space O(1).
func (s *Set[T]) Retain(pred func(T) bool) {
	for v := range s.m {
		if !pred(v) {
			delete(s.m, v)
		}
	}
}

// Delete removes the values for which pred returns true.
//
// Time O(n) expected, space O(1).
func (s *Set[T]) Delete(pred func(T) bool) {
	for v := range s.m {
		if pred(v) {
			delete(s.m, v)
		}
	}
}
