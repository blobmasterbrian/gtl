// Queries: every method here leaves the map untouched.

package hashmap

import (
	"maps"

	"github.com/blobmasterbrian/gtl/algorithm"
	"github.com/blobmasterbrian/gtl/container/slice"
)

// Get returns the value stored for k and whether k is present.
//
// Time O(1) expected.
func (h *HashMap[K, V]) Get(k K) (V, bool) {
	v, ok := h.m[k]
	return v, ok
}

// Contains reports whether k is a key of the map.
//
// Time O(1) expected.
func (h *HashMap[K, V]) Contains(k K) bool {
	_, ok := h.m[k]
	return ok
}

// Len returns the number of entries in the map.
//
// Time O(1).
func (h *HashMap[K, V]) Len() int {
	return len(h.m)
}

// All returns the map's pairs as a lazy pipeline, in no particular order.
//
// Time O(1) to create and O(n) to walk, space O(1).
func (h *HashMap[K, V]) All() algorithm.Stream2[K, V] {
	return algorithm.From2(maps.All(h.m))
}

// Keys returns the map's keys as a lazy pipeline, in no particular order.
//
// Time O(1) to create and O(n) to walk, space O(1).
func (h *HashMap[K, V]) Keys() algorithm.Stream[K] {
	return algorithm.From(maps.Keys(h.m))
}

// Values returns the map's values as a lazy pipeline, in no particular order.
//
// Time O(1) to create and O(n) to walk, space O(1).
func (h *HashMap[K, V]) Values() algorithm.Stream[V] {
	return algorithm.From(maps.Values(h.m))
}

// Keep returns a new map of the entries for which pred returns true.
//
// Time O(n) expected, space O(k).
func (h *HashMap[K, V]) Keep(pred func(K, V) bool) *HashMap[K, V] {
	return Collect(h.All().Keep(pred))
}

// Reject returns a new map of the entries for which pred returns false.
//
// Time O(n) expected, space O(k).
func (h *HashMap[K, V]) Reject(pred func(K, V) bool) *HashMap[K, V] {
	return Collect(h.All().Reject(pred))
}

// Map returns a new map of f applied to each entry. If f gives two entries
// the same key, one of them is kept; which one is unspecified.
//
// Time O(n) expected, space O(n), sized up front.
func (h *HashMap[K, V]) Map[K2 comparable, V2 any](f func(K, V) (K2, V2)) *HashMap[K2, V2] {
	out := &HashMap[K2, V2]{m: make(map[K2]V2, len(h.m))}
	for k, v := range h.m {
		k2, v2 := f(k, v)
		out.m[k2] = v2
	}
	return out
}

// MapKeys returns a new map with f applied to each key and the values
// unchanged. If f gives two keys the same result, one of their entries is
// kept; which one is unspecified.
//
// Time O(n) expected, space O(n), sized up front.
func (h *HashMap[K, V]) MapKeys[K2 comparable](f func(K) K2) *HashMap[K2, V] {
	out := &HashMap[K2, V]{m: make(map[K2]V, len(h.m))}
	for k, v := range h.m {
		out.m[f(k)] = v
	}
	return out
}

// MapValues returns a new map with the same keys and f applied to each value.
//
// Time O(n) expected, space O(n), sized up front.
func (h *HashMap[K, V]) MapValues[U any](f func(V) U) *HashMap[K, U] {
	out := &HashMap[K, U]{m: make(map[K]U, len(h.m))}
	for k, v := range h.m {
		out.m[k] = f(v)
	}
	return out
}

// Collapse returns f applied to each entry, collapsing the entries to a slice
// in no particular order.
//
// Time O(n), space O(n) in a single allocation.
func (h *HashMap[K, V]) Collapse[U any](f func(K, V) U) slice.Slice[U] {
	out := make(slice.Slice[U], 0, len(h.m))
	for k, v := range h.m {
		out = append(out, f(k, v))
	}
	return out
}
