// Queries: every method here leaves the set untouched.

package set

import (
	"maps"

	"github.com/blobmasterbrian/gtl/algorithm"
)

// Contains reports whether v is in the set.
//
// Time O(1) expected.
func (s *Set[T]) Contains(v T) bool {
	_, ok := s.m[v]
	return ok
}

// Len returns the number of values in the set.
//
// Time O(1).
func (s *Set[T]) Len() int {
	return len(s.m)
}

// All returns the set's values as a lazy pipeline, in no particular order.
//
// Time O(1) to create and O(n) to walk, space O(1).
func (s *Set[T]) All() algorithm.Stream[T] {
	return algorithm.From(maps.Keys(s.m))
}

// Keep returns a new set of the values for which pred returns true.
//
// Time O(n) expected, space O(k).
func (s *Set[T]) Keep(pred func(T) bool) *Set[T] {
	return Collect(s.All().Keep(pred))
}

// Reject returns a new set of the values for which pred returns false.
//
// Time O(n) expected, space O(k).
func (s *Set[T]) Reject(pred func(T) bool) *Set[T] {
	return Collect(s.All().Reject(pred))
}

// Map returns a new set of f applied to each value.
//
// Time O(n) expected, space O(k), with k at most n after deduplication.
func (s *Set[T]) Map[U comparable](f func(T) U) *Set[U] {
	return Collect(s.All().Map(f))
}
