package algorithm

import "iter"

// Seq is satisfied by every push-iterator function type over single values,
// including iter.Seq and Stream, so a function taking an S accepts either.
type Seq[T any] interface{ ~func(func(T) bool) }

// Seq2 is the Seq counterpart for pairs, satisfied by iter.Seq2 and Stream2.
type Seq2[K, V any] interface{ ~func(func(K, V) bool) }

// Stream is a sequence with this package's pipeline operations as methods.
// It has the shape of an iter.Seq and can be ranged over directly.
type Stream[T any] iter.Seq[T]

// Stream2 is the Stream counterpart for pairs, such as a map's entries or a
// slice's index/value pairs.
type Stream2[K, V any] iter.Seq2[K, V]

// From wraps a sequence as a Stream.
//
// Time O(1), no allocation.
func From[S Seq[T], T any](s S) Stream[T] {
	return Stream[T](s)
}

// From2 wraps a sequence of pairs as a Stream2.
//
// Time O(1), no allocation.
func From2[S Seq2[K, V], K, V any](s S) Stream2[K, V] {
	return Stream2[K, V](s)
}

// Seq returns the stream as a plain iter.Seq, for functions that take one.
//
// Time O(1), no allocation.
func (s Stream[T]) Seq() iter.Seq[T] {
	return iter.Seq[T](s)
}

// Seq returns the stream as a plain iter.Seq2, for functions that take one.
//
// Time O(1), no allocation.
func (s Stream2[K, V]) Seq() iter.Seq2[K, V] {
	return iter.Seq2[K, V](s)
}
