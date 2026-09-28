// Command gendoc writes the function tables in every README.md from the
// package sources.
//
// Each package README and the top-level README hold a region between
// "<!-- gendoc:begin -->" and "<!-- gendoc:end -->" markers. The region lists
// the package's exported functions, grouped under the title of the header
// comment of the file that declares them, followed by each exported type with
// its constructors and methods. Everything outside the region is left as is.
//
// A row's description is the first sentence of the doc comment. Its
// complexity is taken from the doc comment's paragraph that begins with
// "Time": the leading sentences that begin with "Time" or "Space", leaving any
// qualifiers that follow to go doc. Every exported function and method must
// have such a paragraph. A method whose whole body is a call to a function of
// the same name, in its own package or another package of this module, must
// state the same complexity as that function. Generation fails otherwise.
package main

//go:generate go run .

import (
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

const (
	beginMarker = "<!-- gendoc:begin -->"
	endMarker   = "<!-- gendoc:end -->"
)

func main() {
	root, err := moduleRoot()
	if err != nil {
		fatal(err)
	}
	files, err := render(root)
	if err != nil {
		fatal(err)
	}
	for p, content := range files {
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			fatal(err)
		}
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "gendoc:", err)
	os.Exit(1)
}

// moduleRoot walks up from the working directory to the directory that holds
// go.mod.
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.mod not found")
		}
		dir = parent
	}
}

// modulePath returns the module path declared in root's go.mod.
func modulePath(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(rest), nil
		}
	}
	return "", errors.New("go.mod: no module directive")
}

type row struct {
	name       string
	desc       string
	complexity string
}

type section struct {
	title string
	desc  string // shown under the title when not empty
	what  string // column heading for the names
	rows  []row
}

// wrapper is a method whose whole body is a call to a function of the same
// name.
type wrapper struct {
	pos        token.Position
	name       string // Type.Method
	complexity string
	importPath string // package of the wrapped function; "" for the method's own package
	fn         string
}

type pkg struct {
	name         string
	rel          string // directory relative to the module root, slash-separated
	synopsis     string
	sections     []section
	complexities map[string]string // package-level function name to its complexity
	wrappers     []wrapper
}

// render returns the full content of every README that holds a gendoc region,
// keyed by path, with the region regenerated.
func render(root string) (map[string]string, error) {
	pkgs, err := loadPackages(root)
	if err != nil {
		return nil, err
	}
	mod, err := modulePath(root)
	if err != nil {
		return nil, err
	}
	if err := verifyWrappers(pkgs, mod); err != nil {
		return nil, err
	}

	out := make(map[string]string)
	var top strings.Builder
	for _, p := range pkgs {
		fmt.Fprintf(&top, "### %s\n\n%s ([README](%s/README.md))\n\n", p.rel, p.synopsis, p.rel)
		top.WriteString(p.markdown(4))

		readme := filepath.Join(root, p.rel, "README.md")
		content, err := replaceRegion(readme, p.markdown(2))
		if err != nil {
			return nil, err
		}
		out[readme] = content
	}

	readme := filepath.Join(root, "README.md")
	content, err := replaceRegion(readme, top.String())
	if err != nil {
		return nil, err
	}
	out[readme] = content
	return out, nil
}

// loadPackages returns every library package under root, in path order.
// Commands and anything under an internal directory are left out.
func loadPackages(root string) ([]*pkg, error) {
	var dirs []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if name := d.Name(); p != root && (strings.HasPrefix(name, ".") || name == "internal" || name == "testdata") {
			return filepath.SkipDir
		}
		dirs = append(dirs, p)
		return nil
	})
	if err != nil {
		return nil, err
	}
	slices.Sort(dirs)

	var pkgs []*pkg
	for _, dir := range dirs {
		p, err := loadPackage(root, dir)
		if err != nil {
			return nil, err
		}
		if p != nil {
			pkgs = append(pkgs, p)
		}
	}
	return pkgs, nil
}

// loadPackage reads the non-test Go files in dir. It returns nil when dir
// holds no library package.
func loadPackage(root, dir string) (*pkg, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	var files []*ast.File
	headers := make(map[string]string)            // source path to its header title
	imports := make(map[string]map[string]string) // source path to import name to import path
	var headerOrder []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		src := filepath.Join(dir, name)
		f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
		if title, ok := header(f); ok {
			headers[src] = title
			headerOrder = append(headerOrder, title)
		}
		imports[src] = make(map[string]string)
		for _, imp := range f.Imports {
			ipath := strings.Trim(imp.Path.Value, `"`)
			iname := path.Base(ipath)
			if imp.Name != nil {
				iname = imp.Name.Name
			}
			imports[src][iname] = ipath
		}
	}
	if len(files) == 0 || files[0].Name.Name == "main" {
		return nil, nil
	}

	rel, err := filepath.Rel(root, dir)
	if err != nil {
		return nil, err
	}
	rel = filepath.ToSlash(rel)
	dp, err := doc.NewFromFiles(fset, files, rel, doc.PreserveAST) // keep function bodies for wrapped
	if err != nil {
		return nil, err
	}

	p := &pkg{name: dp.Name, rel: rel, synopsis: dp.Synopsis(dp.Doc), complexities: make(map[string]string)}

	note := func(f *doc.Func, name string) (string, error) {
		c := complexity(f.Doc)
		if c == "" {
			return "", fmt.Errorf("%s: %s has no complexity note", fset.Position(f.Decl.Pos()), name)
		}
		return c, nil
	}

	for _, f := range dp.Funcs {
		c, err := note(f, f.Name)
		if err != nil {
			return nil, err
		}
		p.complexities[f.Name] = c
	}
	for _, t := range dp.Types {
		for _, f := range t.Funcs {
			c, err := note(f, f.Name)
			if err != nil {
				return nil, err
			}
			p.complexities[f.Name] = c
		}
	}

	groups := make(map[string][]row)
	for _, f := range dp.Funcs {
		src := fset.Position(f.Decl.Pos()).Filename
		title, ok := headers[src]
		if !ok {
			title = strings.TrimSuffix(filepath.Base(src), ".go")
			if !slices.Contains(headerOrder, title) {
				headerOrder = append(headerOrder, title)
			}
		}
		groups[title] = append(groups[title], row{f.Name, dp.Synopsis(f.Doc), p.complexities[f.Name]})
	}
	for _, title := range headerOrder {
		if rows := groups[title]; len(rows) > 0 {
			p.sections = append(p.sections, section{title: title, what: "Function", rows: rows})
		}
	}

	for _, t := range dp.Types {
		s := section{title: typeName(t), desc: dp.Synopsis(t.Doc), what: "Method"}
		for _, f := range t.Funcs {
			s.rows = append(s.rows, row{f.Name, dp.Synopsis(f.Doc), p.complexities[f.Name]})
		}
		for _, m := range t.Methods {
			name := t.Name + "." + m.Name
			c, err := note(m, name)
			if err != nil {
				return nil, err
			}
			pos := fset.Position(m.Decl.Pos())
			if qual, fn, ok := wrapped(m.Decl); ok && fn == m.Name {
				w := wrapper{pos: pos, name: name, complexity: c, fn: fn}
				if qual != "" {
					w.importPath, ok = imports[pos.Filename][qual]
				}
				if ok {
					p.wrappers = append(p.wrappers, w)
				}
			}
			s.rows = append(s.rows, row{m.Name, dp.Synopsis(m.Doc), c})
		}
		p.sections = append(p.sections, s)
	}
	return p, nil
}

// verifyWrappers returns an error for the first wrapper whose complexity
// differs from the function it wraps. Functions outside the module are not
// checked.
func verifyWrappers(pkgs []*pkg, mod string) error {
	byPath := make(map[string]*pkg)
	for _, p := range pkgs {
		byPath[path.Join(mod, p.rel)] = p
	}
	for _, p := range pkgs {
		for _, w := range p.wrappers {
			target, fn := p, w.fn
			if w.importPath != "" {
				if target = byPath[w.importPath]; target == nil {
					continue
				}
				fn = target.name + "." + w.fn
			}
			if c, ok := target.complexities[w.fn]; ok && c != w.complexity {
				return fmt.Errorf("%s: %s states %q but wraps %s, which states %q", w.pos, w.name, w.complexity, fn, c)
			}
		}
	}
	return nil
}

// header returns the title of a file's header comment: the text before the
// first colon of a comment that sits above the package clause without being
// the package doc.
func header(f *ast.File) (string, bool) {
	if f.Doc != nil || len(f.Comments) == 0 {
		return "", false
	}
	c := f.Comments[0]
	if c.End() > f.Package {
		return "", false
	}
	title, _, _ := strings.Cut(c.Text(), ":")
	return strings.TrimSpace(title), true
}

// typeName returns the type's name with its type parameters, as in Set[T].
func typeName(t *doc.Type) string {
	spec := t.Decl.Specs[0].(*ast.TypeSpec)
	if spec.TypeParams == nil {
		return t.Name
	}
	var names []string
	for _, field := range spec.TypeParams.List {
		for _, n := range field.Names {
			names = append(names, n.Name)
		}
	}
	return t.Name + "[" + strings.Join(names, ", ") + "]"
}

// complexity returns the leading sentences of the doc's complexity paragraph
// that begin with Time or Space, or "" when the doc has no such paragraph.
func complexity(text string) string {
	for _, para := range strings.Split(text, "\n\n") {
		para = strings.Join(strings.Fields(para), " ")
		if !strings.HasPrefix(para, "Time ") {
			continue
		}
		var kept []string
		for para != "" {
			sentence, rest, _ := strings.Cut(para, ". ")
			if !strings.HasSuffix(sentence, ".") {
				sentence += "."
			}
			if !strings.HasPrefix(sentence, "Time ") && !strings.HasPrefix(sentence, "Space ") {
				break
			}
			kept = append(kept, sentence)
			para = rest
		}
		return strings.Join(kept, " ")
	}
	return ""
}

// wrapped returns the function a method delegates to when its whole body is
// "return F(...)" or "F(...)": the package qualifier, empty for an unqualified
// call, and the function name.
func wrapped(decl *ast.FuncDecl) (qual, name string, ok bool) {
	if decl.Body == nil || len(decl.Body.List) != 1 {
		return "", "", false
	}
	var expr ast.Expr
	switch st := decl.Body.List[0].(type) {
	case *ast.ReturnStmt:
		if len(st.Results) != 1 {
			return "", "", false
		}
		expr = st.Results[0]
	case *ast.ExprStmt:
		expr = st.X
	default:
		return "", "", false
	}
	call, isCall := expr.(*ast.CallExpr)
	if !isCall {
		return "", "", false
	}
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		return "", fn.Name, true
	case *ast.SelectorExpr:
		if x, isIdent := fn.X.(*ast.Ident); isIdent {
			return x.Name, fn.Sel.Name, true
		}
	}
	return "", "", false
}

// markdown renders the package's sections with headings at the given level.
func (p *pkg) markdown(level int) string {
	var b strings.Builder
	h := strings.Repeat("#", level)
	for _, s := range p.sections {
		fmt.Fprintf(&b, "%s %s\n\n", h, s.title)
		if s.desc != "" {
			fmt.Fprintf(&b, "%s\n\n", s.desc)
		}
		if len(s.rows) == 0 { // a type with no functions or methods, such as a constraint
			continue
		}
		fmt.Fprintf(&b, "| %s | Description | Complexity |\n|---|---|---|\n", s.what)
		for _, r := range s.rows {
			fmt.Fprintf(&b, "| `%s` | %s | %s |\n", r.name, r.desc, r.complexity)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// replaceRegion returns the file's content with the text between the gendoc
// markers replaced.
func replaceRegion(readme, content string) (string, error) {
	data, err := os.ReadFile(readme)
	if err != nil {
		return "", err
	}
	text := string(data)
	start := strings.Index(text, beginMarker)
	stop := strings.Index(text, endMarker)
	if start < 0 || stop < 0 || stop < start {
		return "", fmt.Errorf("%s: gendoc markers missing or out of order", readme)
	}
	start += len(beginMarker)
	return text[:start] + "\n\n" + content + text[stop:], nil
}
