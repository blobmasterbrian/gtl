package algorithm

import (
	"slices"
	"strconv"
	"testing"
)

func TestChain(t *testing.T) {
	src, yielded := counting(100)

	got := From(src).
		Keep(even).
		Map(func(n int) string { return strconv.Itoa(n * 10) }).
		Take(3).
		Collect()
	if want := []string{"20", "40", "60"}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if *yielded != 6 {
		t.Errorf("source yielded %d elements, want 6", *yielded)
	}
}

func TestFromAcceptsEitherShape(t *testing.T) {
	xs := []int{1, 2, 3}
	plain := From(slices.Values(xs)) // an iter.Seq
	again := From(plain)             // a Stream
	if got := again.Collect(); !slices.Equal(got, xs) {
		t.Errorf("From(Stream) = %v, want %v", got, xs)
	}
}

func TestSeqCrossesToStandardLibrary(t *testing.T) {
	s := From(slices.Values([]int{3, 1, 2}))
	if got := slices.Sorted(s.Seq()); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("slices.Sorted(s.Seq()) = %v", got)
	}
}

func TestStreamIsRangeable(t *testing.T) {
	n := 0
	for range From(slices.Values([]int{1, 2, 3})).Drop(1) {
		n++
	}
	if n != 2 {
		t.Errorf("ranged over %d elements, want 2", n)
	}
}
