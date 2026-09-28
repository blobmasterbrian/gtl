// Mutations: every method here modifies the map in place.

package hashmap

import "github.com/blobmasterbrian/gtl/algorithm"

// Set stores v under k, replacing any value already stored there.
//
// Time O(1) expected, amortized over growth.
func (h *HashMap[K, V]) Set(k K, v V) {
	if h.m == nil {
		h.m = make(map[K]V)
	}
	h.m[k] = v
}

// Remove removes the entry for k. Removing an absent key does nothing.
//
// Time O(1) expected.
func (h *HashMap[K, V]) Remove(k K) {
	delete(h.m, k)
}

// Insert stores every pair of src, which may be an iter.Seq2 or an
// algorithm.Stream2, replacing any values already stored under their keys.
//
// Time O(n) expected in the pairs of src, space O(n).
func (h *HashMap[K, V]) Insert[S algorithm.Seq2[K, V]](src S) {
	for k, v := range src {
		h.Set(k, v)
	}
}

// Retain keeps the entries for which pred returns true and deletes the rest.
//
// Time O(n) expected, space O(1).
func (h *HashMap[K, V]) Retain(pred func(K, V) bool) {
	for k, v := range h.m {
		if !pred(k, v) {
			delete(h.m, k)
		}
	}
}

// Delete deletes the entries for which pred returns true.
//
// Time O(n) expected, space O(1).
func (h *HashMap[K, V]) Delete(pred func(K, V) bool) {
	for k, v := range h.m {
		if pred(k, v) {
			delete(h.m, k)
		}
	}
}
