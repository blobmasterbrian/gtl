// Package hashmap provides an unordered collection of key/value pairs.
package hashmap

import (
	"maps"

	"github.com/blobmasterbrian/gtl/algorithm"
)

// HashMap is an unordered collection of key/value pairs with unique keys. The
// zero value is an empty map ready to use.
//
// A HashMap is backed by a Go map, so single-key operations are expected
// O(1), amortized over growth. Complexity notes use n for the number of
// entries in the input and k for the number in the result.
type HashMap[K comparable, V any] struct {
	m map[K]V
}

// New returns a map containing the entries of every map in sources, a later
// map overriding an earlier one where keys repeat. The entries are copied.
//
// Time O(n), space O(n), sized up front.
func New[K comparable, V any](sources ...map[K]V) *HashMap[K, V] {
	size := 0
	for _, src := range sources {
		size += len(src)
	}
	h := &HashMap[K, V]{m: make(map[K]V, size)}
	for _, src := range sources {
		maps.Copy(h.m, src)
	}
	return h
}

// Collect returns a map containing the pairs of src, which may be an
// iter.Seq2 or an algorithm.Stream2. A later pair overrides an earlier one
// with the same key.
//
// Time O(n) expected, space O(n). A sequence carries no length, so the map
// grows as pairs arrive.
func Collect[S algorithm.Seq2[K, V], K comparable, V any](src S) *HashMap[K, V] {
	h := New[K, V]()
	h.Insert(src)
	return h
}

// Wrap gives a map that is not a HashMap the container's methods without
// copying it. The two share their entries: a change through either is
// visible through the other.
//
// Time O(1).
func Wrap[K comparable, V any](m map[K]V) *HashMap[K, V] {
	return &HashMap[K, V]{m: m}
}
