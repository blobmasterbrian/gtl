package slice

import (
	"slices"
	"strconv"
	"testing"
)

func even(n int) bool { return n%2 == 0 }

func TestNewCopies(t *testing.T) {
	src := []int{1, 2, 3}
	s := New(src...)
	s[0] = 99
	if want := []int{1, 2, 3}; !slices.Equal(src, want) {
		t.Errorf("New aliased its input: %v", src)
	}
	if want := (Slice[int]{99, 2, 3}); !slices.Equal(s, want) {
		t.Errorf("got %v, want %v", s, want)
	}
}

func TestCollectAll(t *testing.T) {
	s := Collect(slices.Values([]int{1, 2, 3})) // from an iter.Seq
	if want := (Slice[int]{1, 2, 3}); !slices.Equal(s, want) {
		t.Errorf("Collect = %v, want %v", s, want)
	}
	if got := s.All().Collect(); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("All yielded %v", got)
	}
	if got := Collect(s.All()); !slices.Equal(got, s) { // from a Stream
		t.Errorf("Collect(All) = %v, want %v", got, s)
	}
}

func TestWrap(t *testing.T) {
	raw := []int{1, 2, 3, 4, 5, 6}

	if got, want := Wrap(&raw).Keep(even), (Slice[int]{2, 4, 6}); !slices.Equal(got, want) {
		t.Errorf("Keep through Wrap = %v, want %v", got, want)
	}
	if len(raw) != 6 {
		t.Errorf("a query through Wrap changed raw: %v", raw)
	}

	Wrap(&raw).Retain(even)
	if want := []int{2, 4, 6}; !slices.Equal(raw, want) {
		t.Errorf("Retain through Wrap left raw as %v, want %v", raw, want)
	}
}

func TestQueries(t *testing.T) {
	in := Slice[int]{1, 2, 3, 4, 5}

	if got, want := in.Keep(even), (Slice[int]{2, 4}); !slices.Equal(got, want) {
		t.Errorf("Keep = %v, want %v", got, want)
	}
	if got, want := in.Reject(even), (Slice[int]{1, 3, 5}); !slices.Equal(got, want) {
		t.Errorf("Reject = %v, want %v", got, want)
	}
	if got, want := in.Map(strconv.Itoa), (Slice[string]{"1", "2", "3", "4", "5"}); !slices.Equal(got, want) {
		t.Errorf("Map = %v, want %v", got, want)
	}
	if got := in.Reduce(0, func(acc, n int) int { return acc + n }); got != 15 {
		t.Errorf("Reduce = %d, want 15", got)
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

	var seen []int
	in.Each(func(n int) { seen = append(seen, n) })
	if !slices.Equal(seen, in) {
		t.Errorf("Each saw %v, want %v", seen, in)
	}
	if want := (Slice[int]{1, 2, 3, 4, 5}); !slices.Equal(in, want) {
		t.Errorf("a query modified its receiver: %v", in)
	}
}

func TestChain(t *testing.T) {
	got := Slice[int]{1, 2, 3, 4, 5}.
		Reject(even).
		Map(func(n int) string { return strconv.Itoa(n * 10) })
	if want := (Slice[string]{"10", "30", "50"}); !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	lazy := Slice[int]{1, 2, 3, 4, 5}.All().Keep(even).Collect()
	if want := []int{2, 4}; !slices.Equal(lazy, want) {
		t.Errorf("lazy chain = %v, want %v", lazy, want)
	}
}

func TestMutations(t *testing.T) {
	tests := []struct {
		name string
		in   Slice[int]
		fn   func(*Slice[int], func(int) bool)
		want Slice[int]
	}{
		{"Retain", Slice[int]{1, 2, 3, 4, 5, 6}, (*Slice[int]).Retain, Slice[int]{2, 4, 6}},
		{"Retain none", Slice[int]{1, 3, 5}, (*Slice[int]).Retain, Slice[int]{}},
		{"Retain all", Slice[int]{2, 4}, (*Slice[int]).Retain, Slice[int]{2, 4}},
		{"Retain empty", Slice[int]{}, (*Slice[int]).Retain, Slice[int]{}},
		{"Delete", Slice[int]{1, 2, 3, 4, 5, 6}, (*Slice[int]).Delete, Slice[int]{1, 3, 5}},
		{"Delete none", Slice[int]{1, 3, 5}, (*Slice[int]).Delete, Slice[int]{1, 3, 5}},
		{"Delete all", Slice[int]{2, 4}, (*Slice[int]).Delete, Slice[int]{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.in
			before := cap(s)
			tt.fn(&s, even)

			if !slices.Equal(s, tt.want) {
				t.Errorf("got %v, want %v", s, tt.want)
			}
			if cap(s) != before {
				t.Errorf("cap changed from %d to %d: not in place", before, cap(s))
			}
			for i, v := range s[len(s):cap(s)] {
				if v != 0 {
					t.Errorf("tail[%d] = %d, want cleared", i, v)
				}
			}
		})
	}

	s := Slice[int]{1, 2, 3}
	s.Apply(func(n int) int { return n * 10 })
	if want := (Slice[int]{10, 20, 30}); !slices.Equal(s, want) {
		t.Errorf("Apply left %v, want %v", s, want)
	}
}
