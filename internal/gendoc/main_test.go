package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadmesAreCurrent(t *testing.T) {
	root, err := moduleRoot()
	if err != nil {
		t.Fatal(err)
	}
	want, err := render(root)
	if err != nil {
		t.Fatal(err)
	}
	for readme, content := range want {
		got, err := os.ReadFile(readme)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != content {
			rel, _ := filepath.Rel(root, readme)
			t.Errorf("%s is stale: run go generate ./...", rel)
		}
	}
}

func TestComplexity(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		want string
	}{
		{"none", "F does a thing.\n", ""},
		{"one sentence", "F does a thing.\n\nTime O(n), space O(1).\n", "Time O(n), space O(1)."},
		{"qualifier dropped", "F does a thing.\n\nTime O(n), space O(k). The result is grown by append.\n", "Time O(n), space O(k)."},
		{"space sentence kept", "F does a thing.\n\nTime O(n) in the worst case, stopping early. Space O(1).\n", "Time O(n) in the worst case, stopping early. Space O(1)."},
		{"wrapped lines", "F does a thing.\n\nTime O(n),\nspace O(1).\n", "Time O(n), space O(1)."},
		{"later paragraph", "F does a thing.\n\nMore detail.\n\nTime O(1).\n", "Time O(1)."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := complexity(tt.doc); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// loadModule writes each source as the only file of a package under a
// temporary module named m, then loads and verifies them. The map key is the
// package's directory relative to the module root.
func loadModule(t *testing.T, sources map[string]string) ([]*pkg, error) {
	t.Helper()
	root := t.TempDir()
	var pkgs []*pkg
	for _, rel := range sortedKeys(sources) {
		dir := filepath.Join(root, rel)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "x.go"), []byte(sources[rel]), 0o644); err != nil {
			t.Fatal(err)
		}
		p, err := loadPackage(root, dir)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, p)
	}
	return pkgs, verifyWrappers(pkgs, "m")
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

func TestMissingNoteFails(t *testing.T) {
	_, err := loadModule(t, map[string]string{"x": `package x

// F does a thing.
func F() {}
`})
	if err == nil || !strings.Contains(err.Error(), "F has no complexity note") {
		t.Fatalf("got %v, want a missing-note error for F", err)
	}
}

func TestMissingMethodNoteFails(t *testing.T) {
	_, err := loadModule(t, map[string]string{"x": `package x

// T is a type.
type T struct{}

// M does a thing.
func (T) M() {}
`})
	if err == nil || !strings.Contains(err.Error(), "T.M has no complexity note") {
		t.Fatalf("got %v, want a missing-note error for T.M", err)
	}
}

func TestWrapperMismatchFails(t *testing.T) {
	_, err := loadModule(t, map[string]string{"x": `package x

// F does a thing.
//
// Time O(n), space O(1).
func F() {}

// T is a type.
type T struct{}

// F does a thing.
//
// Time O(1), space O(1).
func (T) F() { F() }
`})
	if err == nil || !strings.Contains(err.Error(), "T.F states") {
		t.Fatalf("got %v, want a mismatch error for T.F", err)
	}
}

func TestWrapperMatchPasses(t *testing.T) {
	pkgs, err := loadModule(t, map[string]string{"x": `package x

// F does a thing.
//
// Time O(n), space O(1). With a qualifier.
func F() {}

// T is a type.
type T struct{}

// F does a thing.
//
// Time O(n), space O(1).
func (T) F() { F() }
`})
	if err != nil {
		t.Fatal(err)
	}
	if got := pkgs[0].sections[1].rows[0].complexity; got != "Time O(n), space O(1)." {
		t.Errorf("T.F complexity = %q", got)
	}
}

func TestCrossPackageWrapperMismatchFails(t *testing.T) {
	_, err := loadModule(t, map[string]string{
		"a": `package a

// F does a thing.
//
// Time O(n), space O(1).
func F() {}
`,
		"b": `package b

import "m/a"

// T is a type.
type T struct{}

// F does a thing.
//
// Time O(1), space O(1).
func (T) F() { a.F() }
`,
	})
	if err == nil || !strings.Contains(err.Error(), "wraps a.F") {
		t.Fatalf("got %v, want a mismatch error naming a.F", err)
	}
}

func TestNonDelegatingCallIsNotAWrapper(t *testing.T) {
	_, err := loadModule(t, map[string]string{"x": `package x

// F does a thing.
//
// Time O(n), space O(1).
func F() int { return 0 }

// G does more.
//
// Time O(n), space O(n).
func G(int) int { return 0 }

// T is a type.
type T struct{}

// F does a thing and more.
//
// Time O(n), space O(n).
func (T) F() int { return G(F()) }
`})
	if err != nil {
		t.Fatalf("got %v, want no error: the outer call is G, not F", err)
	}
}
