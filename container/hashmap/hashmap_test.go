package hashmap

import (
	"maps"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func oddValue(_ string, v int) bool { return v%2 == 1 }

func entries(h *HashMap[string, int]) map[string]int { return maps.Collect(h.All().Seq()) }

func TestNew(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}
	h := New(src, map[string]int{"b": 20, "c": 3})

	if want := (map[string]int{"a": 1, "b": 20, "c": 3}); !maps.Equal(entries(h), want) {
		t.Errorf("New = %v, want %v", entries(h), want)
	}
	h.Set("a", 99)
	if src["a"] != 1 {
		t.Errorf("New aliased its input: %v", src)
	}
	if empty := New[string, int](); empty.Len() != 0 {
		t.Errorf("New() has %d entries", empty.Len())
	}
}

func TestBasics(t *testing.T) {
	h := New(map[string]int{"a": 1, "b": 2})

	if v, ok := h.Get("a"); !ok || v != 1 {
		t.Errorf("Get(a) = %d, %v, want 1, true", v, ok)
	}
	if v, ok := h.Get("z"); ok || v != 0 {
		t.Errorf("Get(z) = %d, %v, want 0, false", v, ok)
	}
	if !h.Contains("b") || h.Contains("z") || h.Len() != 2 {
		t.Error("Contains or Len gave the wrong answer")
	}

	h.Set("c", 3)
	h.Set("a", 10)
	h.Remove("b")
	h.Remove("z")
	if want := (map[string]int{"a": 10, "c": 3}); !maps.Equal(entries(h), want) {
		t.Errorf("got %v, want %v", entries(h), want)
	}
}

func TestZeroValue(t *testing.T) {
	var h HashMap[string, int]
	if h.Len() != 0 || h.Contains("a") {
		t.Error("zero value is not empty")
	}
	if got := h.All().Count(); got != 0 {
		t.Errorf("zero value yielded %d pairs", got)
	}
	h.Set("a", 1)
	if !h.Contains("a") {
		t.Error("Set on zero value did not take")
	}
}

func TestWrapShares(t *testing.T) {
	raw := map[string]int{"a": 1, "b": 2}
	h := Wrap(raw)

	h.Set("c", 3)
	if raw["c"] != 3 {
		t.Errorf("Set through Wrap not visible in raw: %v", raw)
	}
	h.Retain(oddValue)
	if want := (map[string]int{"a": 1, "c": 3}); !maps.Equal(raw, want) {
		t.Errorf("Retain through Wrap left raw as %v, want %v", raw, want)
	}
}

func TestCollectInsert(t *testing.T) {
	h := Collect(maps.All(map[string]int{"a": 1, "b": 2})) // from an iter.Seq2
	if want := (map[string]int{"a": 1, "b": 2}); !maps.Equal(entries(h), want) {
		t.Errorf("Collect = %v, want %v", entries(h), want)
	}

	h.Insert(maps.All(map[string]int{"b": 20, "c": 3}))
	if want := (map[string]int{"a": 1, "b": 20, "c": 3}); !maps.Equal(entries(h), want) {
		t.Errorf("after Insert = %v, want %v", entries(h), want)
	}

	h.Insert(New(map[string]int{"d": 4}).All()) // from a Stream2
	if !h.Contains("d") {
		t.Error("Insert from a Stream2 did not take")
	}
}

func TestKeysValues(t *testing.T) {
	h := New(map[string]int{"a": 1, "b": 2, "c": 3})

	if got, want := slices.Sorted(h.Keys().Seq()), []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Errorf("Keys = %v, want %v", got, want)
	}
	if got, want := slices.Sorted(h.Values().Seq()), []int{1, 2, 3}; !slices.Equal(got, want) {
		t.Errorf("Values = %v, want %v", got, want)
	}
}

func TestQueries(t *testing.T) {
	h := New(map[string]int{"a": 1, "b": 2, "c": 3})

	if got, want := entries(h.Keep(oddValue)), (map[string]int{"a": 1, "c": 3}); !maps.Equal(got, want) {
		t.Errorf("Keep = %v, want %v", got, want)
	}
	if got, want := entries(h.Reject(oddValue)), (map[string]int{"b": 2}); !maps.Equal(got, want) {
		t.Errorf("Reject = %v, want %v", got, want)
	}
	if h.Len() != 3 {
		t.Error("Keep or Reject modified the receiver")
	}

	swapped := h.Map(func(k string, v int) (int, string) { return v, k })
	if got, want := maps.Collect(swapped.All().Seq()), (map[int]string{1: "a", 2: "b", 3: "c"}); !maps.Equal(got, want) {
		t.Errorf("Map = %v, want %v", got, want)
	}

	upper := h.MapKeys(strings.ToUpper)
	if want := (map[string]int{"A": 1, "B": 2, "C": 3}); !maps.Equal(entries(upper), want) {
		t.Errorf("MapKeys = %v, want %v", entries(upper), want)
	}

	labels := h.Collapse(func(k string, v int) string { return k + strconv.Itoa(v) })
	slices.Sort(labels)
	if want := []string{"a1", "b2", "c3"}; !slices.Equal(labels, want) {
		t.Errorf("Collapse = %v, want %v", labels, want)
	}

	doubled := h.MapValues(func(v int) int { return v * 2 })
	if want := (map[string]int{"a": 2, "b": 4, "c": 6}); !maps.Equal(entries(doubled), want) {
		t.Errorf("MapValues = %v, want %v", entries(doubled), want)
	}
	if got, ok := h.MapValues(strconv.Itoa).Get("a"); !ok || got != "1" {
		t.Errorf("MapValues to a new type: got %q, %v", got, ok)
	}
}

func TestRetainDelete(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*HashMap[string, int], func(string, int) bool)
		want map[string]int
	}{
		{"Retain", (*HashMap[string, int]).Retain, map[string]int{"a": 1, "c": 3}},
		{"Delete", (*HashMap[string, int]).Delete, map[string]int{"b": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(map[string]int{"a": 1, "b": 2, "c": 3})
			tt.fn(h, oddValue)
			if !maps.Equal(entries(h), tt.want) {
				t.Errorf("got %v, want %v", entries(h), tt.want)
			}
		})
	}
}
