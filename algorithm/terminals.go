// Terminals: every method here consumes the stream and returns a value.

package algorithm

import "slices"

// Each calls f for each element.
//
// Time O(n), space O(1).
func (s Stream[T]) Each(f func(T)) {
	for v := range s {
		f(v)
	}
}

// Reduce folds the stream into a single value, starting from init.
//
// Time O(n), space O(1) beyond the accumulator.
func (s Stream[T]) Reduce[A any](init A, f func(A, T) A) A {
	acc := init
	for v := range s {
		acc = f(acc, v)
	}
	return acc
}

// Find returns the first element for which pred returns true.
//
// Time O(n) in the worst case, stopping at the first match and stopping the
// source with it. Space O(1).
func (s Stream[T]) Find(pred func(T) bool) (T, bool) {
	for v := range s {
		if pred(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Exists reports whether pred returns true for any element.
//
// Time O(n) in the worst case, stopping at the first match and stopping the
// source with it. Space O(1).
func (s Stream[T]) Exists(pred func(T) bool) bool {
	_, ok := s.Find(pred)
	return ok
}

// Satisfies reports whether pred returns true for every element.
//
// Time O(n) in the worst case, stopping at the first failure and stopping the
// source with it. Space O(1).
func (s Stream[T]) Satisfies(pred func(T) bool) bool {
	for v := range s {
		if !pred(v) {
			return false
		}
	}
	return true
}

// Count returns the number of elements.
//
// Time O(n), space O(1). There is no shortcut: a sequence carries no length.
func (s Stream[T]) Count() int {
	n := 0
	for range s {
		n++
	}
	return n
}

// Collect runs the pipeline and returns its elements as a slice.
//
// Time O(n), space O(n). The slice is grown by append, so its capacity can
// exceed n.
func (s Stream[T]) Collect() []T {
	return slices.Collect(s.Seq())
}

// Each calls f for each pair.
//
// Time O(n), space O(1).
func (s Stream2[K, V]) Each(f func(K, V)) {
	for k, v := range s {
		f(k, v)
	}
}

// Count returns the number of pairs.
//
// Time O(n), space O(1). There is no shortcut: a sequence carries no length.
func (s Stream2[K, V]) Count() int {
	n := 0
	for range s {
		n++
	}
	return n
}
