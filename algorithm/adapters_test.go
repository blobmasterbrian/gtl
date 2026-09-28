package algorithm

import (
	"maps"
	"slices"
	"strconv"
	"testing"
)

func TestAdapters(t *testing.T) {
	in := From(slices.Values([]int{1, 2, 3, 4, 5}))
	tests := []struct {
		name string
		got  Stream[int]
		want []int
	}{
		{"Keep", in.Keep(even), []int{2, 4}},
		{"Reject", in.Reject(even), []int{1, 3, 5}},
		{"Take", in.Take(2), []int{1, 2}},
		{"Take past end", in.Take(10), []int{1, 2, 3, 4, 5}},
		{"Take zero", in.Take(0), nil},
		{"Drop", in.Drop(3), []int{4, 5}},
		{"Drop past end", in.Drop(10), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.got.Collect(); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	got := From(slices.Values([]int{1, 2, 3})).Map(strconv.Itoa).Collect()
	if want := []string{"1", "2", "3"}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAdaptersAreLazy(t *testing.T) {
	src, yielded := counting(100)
	calls := 0
	pipeline := From(src).Keep(even).Map(func(n int) int { calls++; return n * n })

	if *yielded != 0 || calls != 0 {
		t.Fatalf("building the pipeline did work: yielded %d, calls %d", *yielded, calls)
	}

	got := pipeline.Take(2).Collect()
	if want := []int{4, 16}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if *yielded != 4 {
		t.Errorf("source yielded %d elements, want 4", *yielded)
	}
	if calls != 2 {
		t.Errorf("f called %d times, want 2", calls)
	}
}

func TestBreakStopsSource(t *testing.T) {
	src, yielded := counting(100)
	for n := range From(src).Keep(even) {
		if n == 6 {
			break
		}
	}
	if *yielded != 6 {
		t.Errorf("source yielded %d elements, want 6", *yielded)
	}
}

func TestStream2Adapters(t *testing.T) {
	in := map[string]int{"a": 1, "b": 2, "c": 3}
	odd := func(_ string, v int) bool { return v%2 == 1 }

	got := maps.Collect(From2(maps.All(in)).Keep(odd).Seq())
	if want := map[string]int{"a": 1, "c": 3}; !maps.Equal(got, want) {
		t.Errorf("Keep = %v, want %v", got, want)
	}

	got = maps.Collect(From2(maps.All(in)).Reject(odd).Seq())
	if want := map[string]int{"b": 2}; !maps.Equal(got, want) {
		t.Errorf("Reject = %v, want %v", got, want)
	}

	swapped := maps.Collect(From2(maps.All(in)).Map(func(k string, v int) (int, string) { return v, k }).Seq())
	if want := map[int]string{1: "a", 2: "b", 3: "c"}; !maps.Equal(swapped, want) {
		t.Errorf("Map = %v, want %v", swapped, want)
	}

	labels := From2(maps.All(in)).Collapse(func(k string, v int) string { return k + strconv.Itoa(v) }).Collect()
	slices.Sort(labels)
	if want := []string{"a1", "b2", "c3"}; !slices.Equal(labels, want) {
		t.Errorf("Collapse = %v, want %v", labels, want)
	}
}
