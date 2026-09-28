package algorithm

import "iter"

func even(n int) bool { return n%2 == 0 }

// counting returns a sequence of 1..n that records how many elements it
// yielded, so tests can check that a pipeline stopped when it should have.
func counting(n int) (iter.Seq[int], *int) {
	yielded := new(int)
	return func(yield func(int) bool) {
		for i := 1; i <= n; i++ {
			*yielded++
			if !yield(i) {
				return
			}
		}
	}, yielded
}
