package algorithm

import (
	"maps"
	"slices"
	"testing"
)

func TestTerminals(t *testing.T) {
	in := From(slices.Values([]int{1, 2, 3, 4}))

	if got := in.Reduce(0, func(acc, n int) int { return acc + n }); got != 10 {
		t.Errorf("Reduce = %d, want 10", got)
	}
	if got, ok := in.Find(even); !ok || got != 2 {
		t.Errorf("Find = %d, %v, want 2, true", got, ok)
	}
	if got, ok := in.Find(func(n int) bool { return n > 10 }); ok || got != 0 {
		t.Errorf("Find = %d, %v, want 0, false", got, ok)
	}
	if !in.Exists(even) || in.Exists(func(n int) bool { return n > 10 }) {
		t.Error("Exists gave the wrong answer")
	}
	if in.Satisfies(even) || !in.Satisfies(func(n int) bool { return n > 0 }) {
		t.Error("Satisfies gave the wrong answer")
	}
	if got := in.Count(); got != 4 {
		t.Errorf("Count = %d, want 4", got)
	}

	var seen []int
	in.Each(func(n int) { seen = append(seen, n) })
	if want := []int{1, 2, 3, 4}; !slices.Equal(seen, want) {
		t.Errorf("Each saw %v, want %v", seen, want)
	}
}

func TestFindStopsSource(t *testing.T) {
	src, yielded := counting(100)
	From(src).Find(func(n int) bool { return n == 3 })
	if *yielded != 3 {
		t.Errorf("source yielded %d elements, want 3", *yielded)
	}
}

func TestStream2Terminals(t *testing.T) {
	in := From2(maps.All(map[string]int{"a": 1, "b": 2}))

	if got := in.Count(); got != 2 {
		t.Errorf("Count = %d, want 2", got)
	}
	sum := 0
	in.Each(func(_ string, v int) { sum += v })
	if sum != 3 {
		t.Errorf("Each summed to %d, want 3", sum)
	}
}
