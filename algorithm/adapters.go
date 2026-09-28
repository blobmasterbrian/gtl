// Adapters: every method here returns a new stream and does no work until
// that stream is consumed.

package algorithm

// Keep returns the elements for which pred returns true.
//
// Time O(n), space O(1).
func (s Stream[T]) Keep(pred func(T) bool) Stream[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if pred(v) && !yield(v) {
				return
			}
		}
	}
}

// Reject returns the elements for which pred returns false.
//
// Time O(n), space O(1).
func (s Stream[T]) Reject(pred func(T) bool) Stream[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !pred(v) && !yield(v) {
				return
			}
		}
	}
}

// Map returns f applied to each element.
//
// Time O(n), space O(1).
func (s Stream[T]) Map[U any](f func(T) U) Stream[U] {
	return func(yield func(U) bool) {
		for v := range s {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Take returns the first count elements, or all of them if there are fewer.
// The source is not advanced past the last element taken.
//
// Time O(min(n, count)), space O(1).
func (s Stream[T]) Take(count int) Stream[T] {
	return func(yield func(T) bool) {
		if count <= 0 {
			return
		}
		i := 0
		for v := range s {
			if !yield(v) {
				return
			}
			if i++; i == count {
				return
			}
		}
	}
}

// Drop returns the stream without its first count elements.
//
// Time O(n), space O(1). The dropped elements are still walked: a sequence
// has no random access.
func (s Stream[T]) Drop(count int) Stream[T] {
	return func(yield func(T) bool) {
		i := 0
		for v := range s {
			if i < count {
				i++
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Keep returns the pairs for which pred returns true.
//
// Time O(n), space O(1).
func (s Stream2[K, V]) Keep(pred func(K, V) bool) Stream2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if pred(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Reject returns the pairs for which pred returns false.
//
// Time O(n), space O(1).
func (s Stream2[K, V]) Reject(pred func(K, V) bool) Stream2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range s {
			if !pred(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Map returns f applied to each pair.
//
// Time O(n), space O(1).
func (s Stream2[K, V]) Map[K2, V2 any](f func(K, V) (K2, V2)) Stream2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}

// Collapse returns f applied to each pair, collapsing the pairs to single
// values.
//
// Time O(n), space O(1).
func (s Stream2[K, V]) Collapse[U any](f func(K, V) U) Stream[U] {
	return func(yield func(U) bool) {
		for k, v := range s {
			if !yield(f(k, v)) {
				return
			}
		}
	}
}
