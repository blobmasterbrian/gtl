// Mutations: every method here modifies the slice in place. A method that
// changes the length has a pointer receiver and updates the caller's slice
// header itself.

package slice

// Apply replaces each element with f applied to it.
//
// Time O(n), space O(1).
func (s Slice[T]) Apply(f func(T) T) {
	for i := range s {
		s[i] = f(s[i])
	}
}

// Retain keeps the elements for which pred returns true and shortens the
// slice to fit. The vacated tail of the backing array is cleared.
//
// Time O(n), space O(1).
func (s *Slice[T]) Retain(pred func(T) bool) {
	n := 0
	for _, v := range *s {
		if pred(v) {
			(*s)[n] = v
			n++
		}
	}
	clear((*s)[n:])
	*s = (*s)[:n]
}

// Delete removes the elements for which pred returns true and shortens the
// slice to fit. The vacated tail of the backing array is cleared.
//
// Time O(n), space O(1).
func (s *Slice[T]) Delete(pred func(T) bool) {
	n := 0
	for _, v := range *s {
		if !pred(v) {
			(*s)[n] = v
			n++
		}
	}
	clear((*s)[n:])
	*s = (*s)[:n]
}
