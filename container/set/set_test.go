package set

import (
	"slices"
	"strconv"
	"testing"
)

func even(n int) bool { return n%2 == 0 }

func sorted(s *Set[int]) []int { return slices.Sorted(s.All().Seq()) }

func TestBasics(t *testing.T) {
	s := New(1, 2, 2, 3)
	if s.Len() != 3 {
		t.Errorf("Len = %d, want 3", s.Len())
	}
	if !s.Contains(2) || s.Contains(4) {
		t.Error("Contains gave the wrong answer")
	}

	s.Add(4)
	s.Add(4)
	s.Remove(1)
	s.Remove(99)
	if got, want := sorted(s), []int{2, 3, 4}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestZeroValue(t *testing.T) {
	var s Set[string]
	if s.Len() != 0 || s.Contains("a") {
		t.Error("zero value is not empty")
	}
	if got := s.All().Collect(); len(got) != 0 {
		t.Errorf("zero value yielded %v", got)
	}
	s.Add("a")
	if !s.Contains("a") {
		t.Error("Add on zero value did not take")
	}
}

func TestCollectInsert(t *testing.T) {
	s := Collect(slices.Values([]int{3, 1, 3, 2})) // from an iter.Seq
	if got, want := sorted(s), []int{1, 2, 3}; !slices.Equal(got, want) {
		t.Errorf("Collect = %v, want %v", got, want)
	}

	s.Insert(slices.Values([]int{2, 5}))
	if got, want := sorted(s), []int{1, 2, 3, 5}; !slices.Equal(got, want) {
		t.Errorf("after Insert = %v, want %v", got, want)
	}

	s.Insert(New(7, 8).All()) // from a Stream
	if got, want := sorted(s), []int{1, 2, 3, 5, 7, 8}; !slices.Equal(got, want) {
		t.Errorf("after Insert from a Stream = %v, want %v", got, want)
	}
}

func TestAlgorithms(t *testing.T) {
	s := New(1, 2, 3, 4, 5)

	if got, want := sorted(s.Keep(even)), []int{2, 4}; !slices.Equal(got, want) {
		t.Errorf("Keep = %v, want %v", got, want)
	}
	if got, want := sorted(s.Reject(even)), []int{1, 3, 5}; !slices.Equal(got, want) {
		t.Errorf("Reject = %v, want %v", got, want)
	}
	if s.Len() != 5 {
		t.Error("Keep or Reject modified the receiver")
	}

	labels := s.Map(func(n int) string { return strconv.Itoa(n % 2) })
	if got, want := slices.Sorted(labels.All().Seq()), []string{"0", "1"}; !slices.Equal(got, want) {
		t.Errorf("Map = %v, want %v", got, want)
	}

	sum := s.All().Keep(even).Reduce(0, func(acc, n int) int { return acc + n })
	if sum != 6 {
		t.Errorf("lazy chain = %d, want 6", sum)
	}
}

func TestRetainDelete(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*Set[int], func(int) bool)
		want []int
	}{
		{"Retain", (*Set[int]).Retain, []int{2, 4}},
		{"Delete", (*Set[int]).Delete, []int{1, 3, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(1, 2, 3, 4, 5)
			tt.fn(s, even)
			if got := sorted(s); !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
