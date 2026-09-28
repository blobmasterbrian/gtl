// Package slice provides a slice with the library's algorithms as methods.
package slice

import (
	"slices"

	"github.com/blobmasterbrian/gtl/algorithm"
)

// Slice is a []T with the library's eager algorithms as methods, so a
// pipeline can be written as a chain. All switches a chain to lazy
// evaluation.
//
// Methods with a pointer receiver modify the slice in place; the rest return a
// new value. Complexity notes use n for the number of elements in the input
// and k for the number in the result.
type Slice[T any] []T

// New returns a slice containing items. The items are copied, so the result
// never aliases a slice expanded into the call.
//
// Time O(n), space O(n), sized up front.
func New[T any](items ...T) Slice[T] {
	return slices.Clone(items)
}

// Collect returns a slice containing the elements of src, which may be an
// iter.Seq or an algorithm.Stream.
//
// Time O(n), space O(n). The slice is grown by append, so its capacity can
// exceed n.
func Collect[S algorithm.Seq[T], T any](src S) Slice[T] {
	return algorithm.From(src).Collect()
}

// Wrap gives a slice that is not a Slice the container's methods without
// copying it. It takes the slice's address so that methods which change the
// length update the caller's slice.
//
// Time O(1), no allocation.
func Wrap[T any](s *[]T) *Slice[T] {
	return (*Slice[T])(s)
}
