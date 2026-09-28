// Queries: every method here leaves the slice untouched.

package slice

import (
	"slices"

	"github.com/blobmasterbrian/gtl/algorithm"
)

// All returns the elements as a lazy pipeline, in order.
//
// Time O(1) to create and O(n) to walk, space O(1).
func (s Slice[T]) All() algorithm.Stream[T] {
	return algorithm.From(slices.Values(s))
}

// Keep returns the elements for which pred returns true.
//
// Time O(n), space O(k). The result is grown by append, so its capacity can
// exceed k.
func (s Slice[T]) Keep(pred func(T) bool) Slice[T] {
	var out Slice[T]
	for _, v := range s {
		if pred(v) {
			out = append(out, v)
		}
	}
	return out
}

// Reject returns the elements for which pred returns false.
//
// Time O(n), space O(k). The result is grown by append, so its capacity can
// exceed k.
func (s Slice[T]) Reject(pred func(T) bool) Slice[T] {
	var out Slice[T]
	for _, v := range s {
		if !pred(v) {
			out = append(out, v)
		}
	}
	return out
}

// Map returns f applied to each element.
//
// Time O(n), space O(n) in a single allocation.
func (s Slice[T]) Map[U any](f func(T) U) Slice[U] {
	out := make(Slice[U], len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

// Each calls f for each element.
//
// Time O(n), space O(1).
func (s Slice[T]) Each(f func(T)) {
	for _, v := range s {
		f(v)
	}
}

// Reduce folds the slice into a single value, starting from init.
//
// Time O(n), space O(1) beyond the accumulator.
func (s Slice[T]) Reduce[A any](init A, f func(A, T) A) A {
	acc := init
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

// Find returns the first element for which pred returns true.
//
// Time O(n) in the worst case, stopping at the first match. Space O(1).
func (s Slice[T]) Find(pred func(T) bool) (T, bool) {
	for _, v := range s {
		if pred(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Exists reports whether pred returns true for any element.
//
// Time O(n) in the worst case, stopping at the first match. Space O(1).
func (s Slice[T]) Exists(pred func(T) bool) bool {
	_, ok := s.Find(pred)
	return ok
}

// Satisfies reports whether pred returns true for every element.
//
// Time O(n) in the worst case, stopping at the first failure. Space O(1).
func (s Slice[T]) Satisfies(pred func(T) bool) bool {
	for _, v := range s {
		if !pred(v) {
			return false
		}
	}
	return true
}
