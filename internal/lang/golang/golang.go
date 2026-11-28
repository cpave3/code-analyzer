// Package golang is the Go language adapter, built on Tree-sitter.
//
// It extracts functions, methods and type declarations; per-callable
// cyclomatic complexity, decision count, maximum nesting, length and statement
// count; exact/structural fingerprints; name-based identifier references; and
// a package import graph resolved against go.mod module paths.
//
// It deliberately does not type-check, resolve interface dispatch, honour
// build tags, or build a call graph.
package golang

import (
	"bytes"
	"path"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tsgo "github.com/tree-sitter/tree-sitter-go/bindings/go"

	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/lang/fingerprint"
	"github.com/cameron/code-evaluator/internal/lang/sitterutil"
	"github.com/cameron/code-evaluator/internal/model"
)

// Adapter implements lang.Adapter for Go.
type Adapter struct {
	language *sitter.Language
}

// New returns the Go adapter.
func New() *Adapter {
	return &Adapter{language: sitter.NewLanguage(tsgo.Language())}
}

func (a *Adapter) ID() string           { return "go" }
func (a *Adapter) Version() string      { return "go-adapter/3+tree-sitter-go/0.25.0" }
func (a *Adapter) Extensions() []string { return []string{".go"} }

func (a *Adapter) Capabilities() map[lang.Capability]string {
	return map[lang.Capability]string{
		lang.CapSymbols:      "functions, methods and type declarations (Tree-sitter)",
		lang.CapComplexity:   "cyclomatic = 1 + if/for/case/comm-case/&&/||; nesting counts if/for/switch/select/func-literal, else-if not nested",
		lang.CapFingerprints: "exact token text and structural kinds (identifiers/literals abstracted), comments ignored",
		lang.CapDependencies: "package import graph resolved via nearest go.mod; test files excluded",
		lang.CapReferences:   "name-based identifier references; qualifier mapped to import by alias or last path element; no type resolution",
	}
}

func (a *Adapter) Limits() []string {
	return []string{
		"no type checking: references are matched by name, not resolved",
		"no call graph or interface dispatch resolution",
		"build tags and GOOS/GOARCH file variants are all indexed together",
		"import cycles are compile errors in Go; detected cycles indicate non-building code",
	}
}

func (a *Adapter) IsTestPath(p string) bool { return strings.HasSuffix(p, "_test.go") }

func (a *Adapter) PackageOf(p string) string { return path.Dir(p) }

var generatedRe = regexp.MustCompile(`(?m)^// Code generated .* DO NOT EDIT\.$`)

// Parse parses Go source.
func (a *Adapter) Parse(src []byte) (*model.ParsedFile, error) {
	tree, err := sitterutil.Parse(a.language, src)
	if err != nil {
		return nil, err
	}
	defer tree.Close()
	root := tree.RootNode()

	pf := &model.ParsedFile{Refs: map[string]int{}, Lines: bytes.Count(src, []byte{'\n'}) + 1}
	// Generated header must appear before the package clause.
	if i := bytes.Index(src, []byte("\npackage ")); i >= 0 || bytes.HasPrefix(src, []byte("package ")) {
		head := src
		if i >= 0 {
			head = src[:i+1]
		} else {
			head = nil
		}
		pf.Generated = generatedRe.Match(head)
	}

	p := &fileParser{src: src, pf: pf, errs: sitterutil.Errors{Language: "go"}}
	if root.HasError() {
		p.errs.Collect(root)
		pf.ParseErrors = p.errs.Diags
	}
	p.walkTop(root)
	sort.SliceStable(pf.Symbols, func(i, j int) bool { return pf.Symbols[i].StartLine < pf.Symbols[j].StartLine })
	return pf, nil
}

type fileParser struct {
	src  []byte
	pf   *model.ParsedFile
	errs sitterutil.Errors
}

var line, endLine = sitterutil.Line, sitterutil.EndLine

func (p *fileParser) text(n *sitter.Node) string { return n.Utf8Text(p.src) }

func (p *fileParser) inError(start, end int) bool { return p.errs.In(start, end) }

func (p *fileParser) walkTop(root *sitter.Node) {
	for i := uint(0); i < root.ChildCount(); i++ {
		n := root.Child(i)
		switch n.Kind() {
		case "package_clause":
			for j := uint(0); j < n.ChildCount(); j++ {
				if c := n.Child(j); c.Kind() == "package_identifier" {
					p.pf.PackageName = p.text(c)
				}
			}
		case "import_declaration":
			p.imports(n)
		case "function_declaration":
			p.function(n, "")
		case "method_declaration":
			p.function(n, p.receiverType(n.ChildByFieldName("receiver")))
		case "type_declaration":
			p.typeDecl(n)
		}
		p.refs(n)
	}
}

func (p *fileParser) imports(n *sitter.Node) {
	var visit func(*sitter.Node)
	visit = func(n *sitter.Node) {
		if n.Kind() == "import_spec" {
			imp := model.Import{Line: line(n)}
			if pathNode := n.ChildByFieldName("path"); pathNode != nil {
				imp.Path = strings.Trim(p.text(pathNode), "\"`")
			}
			if name := n.ChildByFieldName("name"); name != nil {
				imp.Name = p.text(name)
			}
			if imp.Path != "" {
				p.pf.Imports = append(p.pf.Imports, imp)
			}
			return
		}
		for i := uint(0); i < n.ChildCount(); i++ {
			visit(n.Child(i))
		}
	}
	visit(n)
}

func (p *fileParser) receiverType(recv *sitter.Node) string {
	if recv == nil {
		return ""
	}
	// Find the first type_identifier under the receiver parameter list.
	var found string
	var visit func(*sitter.Node)
	visit = func(n *sitter.Node) {
		if found != "" {
			return
		}
		if n.Kind() == "type_identifier" {
			found = p.text(n)
			return
		}
		for i := uint(0); i < n.ChildCount(); i++ {
			visit(n.Child(i))
		}
	}
	visit(recv)
	return found
}

func exported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

var spaceRe = regexp.MustCompile(`\s+`)

func (p *fileParser) normalize(n *sitter.Node) string {
	if n == nil {
		return ""
	}
	return spaceRe.ReplaceAllString(p.text(n), " ")
}

func (p *fileParser) function(n *sitter.Node, container string) {
	nameNode := n.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	name := p.text(nameNode)
	sym := model.Symbol{
		Name:      name,
		Container: container,
		Kind:      model.KindFunction,
		StartLine: line(n),
		EndLine:   endLine(n),
	}
	if container != "" {
		sym.Kind = model.KindMethod
		sym.Exported = exported(name) && exported(container)
	} else {
		sym.Exported = exported(name)
	}
	params, result := n.ChildByFieldName("parameters"), n.ChildByFieldName("result")
	sig := p.normalize(n.ChildByFieldName("type_parameters")) + p.normalize(params)
	if result != nil {
		sig += " " + p.normalize(result)
	}
	sym.Signature = sig

	var fb fingerprint.Builder
	for _, part := range []*sitter.Node{n.ChildByFieldName("type_parameters"), params, result, n.ChildByFieldName("body")} {
		if part != nil {
			p.tokens(part, &fb)
		}
	}
	sym.ExactFP, sym.StructFP, sym.Tokens, sym.MinHash = fb.Exact(), fb.Struct(), fb.Tokens(), fb.MinHash()

	m := &model.Metrics{Lines: sym.EndLine - sym.StartLine + 1}
	if body := n.ChildByFieldName("body"); body != nil {
		p.measure(body, 0, m)
	}
	m.Cyclomatic = 1 + m.Decisions
	sym.Metrics = m
	sym.InErrorRegion = p.inError(sym.StartLine, sym.EndLine)
	p.pf.Symbols = append(p.pf.Symbols, sym)
}

func (p *fileParser) typeDecl(n *sitter.Node) {
	for i := uint(0); i < n.ChildCount(); i++ {
		spec := n.Child(i)
		if spec.Kind() != "type_spec" && spec.Kind() != "type_alias" {
			continue
		}
		nameNode := spec.ChildByFieldName("name")
		if nameNode == nil {
			continue
		}
		name := p.text(nameNode)
		var fb fingerprint.Builder
		if t := spec.ChildByFieldName("type"); t != nil {
			p.tokens(t, &fb)
		}
		sym := model.Symbol{
			Name: name, Kind: model.KindType, Exported: exported(name),
			StartLine: line(spec), EndLine: endLine(spec),
			ExactFP: fb.Exact(), StructFP: fb.Struct(), Tokens: fb.Tokens(),
		}
		if t := spec.ChildByFieldName("type"); t != nil {
			sym.Signature = t.Kind()
		}
		sym.InErrorRegion = p.inError(sym.StartLine, sym.EndLine)
		p.pf.Symbols = append(p.pf.Symbols, sym)
	}
}

var literalKinds = map[string]bool{
	"interpreted_string_literal": true, "raw_string_literal": true, "int_literal": true,
	"float_literal": true, "imaginary_literal": true, "rune_literal": true,
}

var identKinds = map[string]bool{
	"identifier": true, "field_identifier": true, "type_identifier": true,
	"package_identifier": true, "label_name": true,
}

// tokens feeds leaf tokens into the fingerprint builder.
func (p *fileParser) tokens(n *sitter.Node, fb *fingerprint.Builder) {
	kind := n.Kind()
	if kind == "comment" {
		return
	}
	if literalKinds[kind] {
		fb.Add(p.text(n), "LIT")
		return
	}
	if n.ChildCount() == 0 {
		switch {
		case identKinds[kind]:
			fb.Add(p.text(n), "ID")
		case kind == "true" || kind == "false" || kind == "nil" || kind == "iota":
			fb.Add(kind, "LIT")
		default:
			fb.Add(p.text(n), kind)
		}
		return
	}
	for i := uint(0); i < n.ChildCount(); i++ {
		p.tokens(n.Child(i), fb)
	}
}

var nestingKinds = map[string]bool{
	"if_statement": true, "for_statement": true, "expression_switch_statement": true,
	"type_switch_statement": true, "select_statement": true, "func_literal": true,
}

var measureSpec = &sitterutil.MeasureSpec{
	Decision: func(n *sitter.Node) bool {
		switch n.Kind() {
		case "if_statement", "for_statement", "expression_case", "type_case", "communication_case":
			return true
		}
		return sitterutil.LogicalOp(n, "&&", "||")
	},
	Statement: func(kind string) bool {
		return (strings.HasSuffix(kind, "_statement") && kind != "block") || kind == "short_var_declaration" || kind == "var_declaration" || kind == "const_declaration"
	},
	Nests: func(n *sitter.Node) bool {
		// "else if" continues the chain rather than nesting deeper.
		return nestingKinds[n.Kind()] && !(n.Kind() == "if_statement" && n.Parent() != nil && n.Parent().Kind() == "if_statement")
	},
}

// measure walks a function body accumulating decisions, nesting and statements.
func (p *fileParser) measure(n *sitter.Node, depth int, m *model.Metrics) {
	sitterutil.Measure(n, depth, m, measureSpec)
}

// refs records identifier uses, skipping declaration names.
func (p *fileParser) refs(n *sitter.Node) {
	var visit func(n *sitter.Node, parent *sitter.Node)
	visit = func(n, parent *sitter.Node) {
		switch n.Kind() {
		case "comment", "interpreted_string_literal", "raw_string_literal":
			return
		case "selector_expression":
			operand, field := n.ChildByFieldName("operand"), n.ChildByFieldName("field")
			if operand != nil && field != nil && operand.Kind() == "identifier" {
				p.pf.Refs[p.text(operand)+"."+p.text(field)]++
			}
			if field != nil {
				p.pf.Refs[p.text(field)]++
			}
			if operand != nil {
				visit(operand, n)
			}
			return
		case "qualified_type":
			pkg, name := n.ChildByFieldName("package"), n.ChildByFieldName("name")
			if pkg != nil && name != nil {
				p.pf.Refs[p.text(pkg)+"."+p.text(name)]++
			}
			return
		case "identifier", "type_identifier", "field_identifier":
			if parent != nil && isDeclName(parent, n) {
				return
			}
			p.pf.Refs[p.text(n)]++
			return
		}
		for i := uint(0); i < n.ChildCount(); i++ {
			visit(n.Child(i), n)
		}
	}
	visit(n, nil)
}

func isDeclName(parent, n *sitter.Node) bool {
	switch parent.Kind() {
	case "function_declaration", "method_declaration", "type_spec", "type_alias":
		if name := parent.ChildByFieldName("name"); name != nil && name.StartByte() == n.StartByte() {
			return true
		}
	}
	return false
}

// BuildGraph builds the package import graph. Nodes are repository-relative
// package directories.
func (a *Adapter) BuildGraph(in lang.GraphInput) (*lang.Graph, []model.Diagnostic) {
	g := lang.NewGraph()
	g.Note = "Go package imports (non-test files), internal packages resolved via nearest go.mod"
	var diags []model.Diagnostic

	// Discover modules.
	type mod struct{ dir, path string }
	var mods []mod
	for _, p := range in.AllPaths {
		if path.Base(p) != "go.mod" || strings.Contains(p, "vendor/") || strings.Contains(p, "testdata/") {
			continue
		}
		data, err := in.ReadFile(p)
		if err != nil {
			diags = append(diags, model.Diagnostic{Severity: model.SevWarning, Code: "gomod_unreadable", Message: err.Error(), File: p, Language: "go"})
			continue
		}
		mp := modulePath(data)
		if mp == "" {
			diags = append(diags, model.Diagnostic{Severity: model.SevWarning, Code: "gomod_no_module", Message: "go.mod has no module directive", File: p, Language: "go"})
			continue
		}
		mods = append(mods, mod{dir: path.Dir(p), path: mp})
	}
	sort.Slice(mods, func(i, j int) bool { return len(mods[i].dir) > len(mods[j].dir) })
	if len(mods) == 0 {
		diags = append(diags, model.Diagnostic{Severity: model.SevWarning, Code: "gomod_missing", Message: "no go.mod found; internal imports cannot be resolved, only directory nodes are recorded", Language: "go"})
	}
	moduleOf := func(dir string) (mod, bool) {
		for _, m := range mods {
			if m.dir == "." || dir == m.dir || strings.HasPrefix(dir, m.dir+"/") {
				return m, true
			}
		}
		return mod{}, false
	}
	importToDir := map[string]string{}
	inGraph := func(f lang.GraphFile) bool {
		return f.Class == model.ClassSource || f.Class == model.ClassGenerated
	}
	for _, f := range in.Files {
		if !inGraph(f) {
			continue
		}
		dir := path.Dir(f.Path)
		g.Nodes[dir] = true
		if m, ok := moduleOf(dir); ok {
			ip := m.path
			if dir != m.dir {
				rel := strings.TrimPrefix(dir, m.dir+"/")
				if m.dir == "." {
					rel = dir
				}
				ip = m.path + "/" + rel
			}
			importToDir[ip] = dir
		}
	}
	for _, f := range in.Files {
		if !inGraph(f) || f.Parsed == nil {
			continue
		}
		from := path.Dir(f.Path)
		for _, imp := range f.Parsed.Imports {
			if to, ok := importToDir[imp.Path]; ok {
				g.AddEdge(from, to, model.Location{File: f.Path, Line: imp.Line})
				continue
			}
			internal := false
			for _, m := range mods {
				if imp.Path == m.path || strings.HasPrefix(imp.Path, m.path+"/") {
					internal = true
				}
			}
			if internal {
				diags = append(diags, model.Diagnostic{Severity: model.SevInfo, Code: "import_unresolved", Message: "internal import has no indexed package: " + imp.Path, File: f.Path, Line: imp.Line, Language: "go"})
				continue
			}
			ext := g.External[from]
			if ext == nil {
				ext = map[string]bool{}
				g.External[from] = ext
			}
			ext[imp.Path] = true
		}
	}
	return g, diags
}

var moduleRe = regexp.MustCompile(`(?m)^module\s+("?)([^\s"]+)("?)\s*$`)

func modulePath(gomod []byte) string {
	m := moduleRe.FindSubmatch(gomod)
	if m == nil {
		return ""
	}
	return string(m[2])
}

// ParseSyntax returns the Tree-sitter syntax tree (lang.SyntaxParser).
func (a *Adapter) ParseSyntax(src []byte) (*sitter.Tree, error) {
	return sitterutil.Parse(a.language, src)
}
