// Package typescript is the TypeScript language adapter, built on Tree-sitter
// (tree-sitter-typescript, both the TypeScript and TSX grammars).
//
// It extracts function declarations, classes and their methods, interfaces,
// type aliases, enums, and top-level variables initialised with an arrow
// function or function expression (React components and hooks, often wrapped
// in memo/forwardRef); per-callable cyclomatic complexity, decisions,
// nesting, length and statements; exact/structural fingerprints; name-based
// references; and a directory-level import graph resolved through relative
// specifiers, tsconfig "paths" aliases and workspace package names.
//
// It does not type-check, resolve re-export chains beyond the first file,
// follow node_modules, or model runtime/dynamic module loading.
package typescript

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tsts "github.com/tree-sitter/tree-sitter-typescript/bindings/go"

	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/lang/fingerprint"
	"github.com/cameron/code-evaluator/internal/model"
)

// Adapter implements lang.Adapter for TypeScript.
type Adapter struct {
	ts, tsx *sitter.Language
}

// New returns the TypeScript adapter.
func New() *Adapter {
	return &Adapter{ts: sitter.NewLanguage(tsts.LanguageTypescript()), tsx: sitter.NewLanguage(tsts.LanguageTSX())}
}

func (a *Adapter) ID() string           { return "typescript" }
func (a *Adapter) Version() string      { return "ts-adapter/1+tree-sitter-typescript/0.23.2" }
func (a *Adapter) Extensions() []string { return []string{".ts", ".tsx", ".mts", ".cts"} }

func (a *Adapter) Capabilities() map[lang.Capability]string {
	return map[lang.Capability]string{
		lang.CapSymbols:      "functions, classes, methods, interfaces, type aliases, enums, and top-level const/let/var bound to arrow functions or function expressions (incl. memo/forwardRef wrappers)",
		lang.CapComplexity:   "cyclomatic = 1 + if/for/for-in-of/while/do/case/catch/ternary/&&/||/??; nesting counts if/loops/switch/try-catch and nested function literals (inner callbacks count toward the enclosing callable); else-if not nested",
		lang.CapFingerprints: "exact token text and structural kinds (identifiers/literals/JSX text abstracted), comments ignored",
		lang.CapDependencies: "directory-level graph from import/export-from/import()/require resolved via relative paths, tsconfig paths and workspace package names; type-only imports included; test files excluded",
		lang.CapReferences:   "name-based identifier references; member access recorded as object.property; no type resolution",
	}
}

func (a *Adapter) Limits() []string {
	return []string{
		"no type checking: references are matched by name, not resolved",
		"inner callbacks and nested function literals are measured as part of their enclosing callable (as Go func literals are)",
		"graph units are directories; directory cycles are legal in TypeScript and indicate coupling, not build errors",
		"node_modules, package exports maps and re-export chains are not followed",
		"each file is parsed with the TSX grammar, falling back to the TypeScript grammar when that yields fewer syntax errors",
	}
}

var testPathRe = regexp.MustCompile(`(^|/)(__tests__|__mocks__|tests?|e2e)/|\.(test|spec|stories)\.[cm]?tsx?$`)

func (a *Adapter) IsTestPath(p string) bool { return testPathRe.MatchString(p) }

func (a *Adapter) PackageOf(p string) string { return path.Dir(p) }

var generatedRe = regexp.MustCompile(`(?i)(@generated|code generated .* do not edit|auto-generated)`)

// Parse parses TypeScript or TSX source.
func (a *Adapter) Parse(src []byte) (*model.ParsedFile, error) {
	pf, errs, err := a.parseWith(a.tsx, src)
	if err != nil {
		return nil, err
	}
	if errs > 0 {
		if pf2, errs2, err2 := a.parseWith(a.ts, src); err2 == nil && errs2 < errs {
			pf = pf2
		}
	}
	return pf, nil
}

func (a *Adapter) parseWith(language *sitter.Language, src []byte) (*model.ParsedFile, int, error) {
	parser := sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(language); err != nil {
		return nil, 0, err
	}
	tree := parser.Parse(src, nil)
	if tree == nil {
		return nil, 0, fmt.Errorf("tree-sitter returned no tree")
	}
	defer tree.Close()
	root := tree.RootNode()
	pf := &model.ParsedFile{Refs: map[string]int{}, Lines: bytes.Count(src, []byte{'\n'}) + 1}
	head := src
	if len(head) > 600 {
		head = head[:600]
	}
	pf.Generated = generatedRe.Match(head)
	p := &fileParser{src: src, pf: pf}
	if root.HasError() {
		p.collectErrors(root)
	}
	p.walkTop(root)
	p.refs(root)
	sort.SliceStable(pf.Symbols, func(i, j int) bool { return pf.Symbols[i].StartLine < pf.Symbols[j].StartLine })
	return pf, len(pf.ParseErrors), nil
}

type fileParser struct {
	src    []byte
	pf     *model.ParsedFile
	errors []int
}

func line(n *sitter.Node) int    { return int(n.StartPosition().Row) + 1 }
func endLine(n *sitter.Node) int { return int(n.EndPosition().Row) + 1 }

func (p *fileParser) text(n *sitter.Node) string { return n.Utf8Text(p.src) }

func (p *fileParser) collectErrors(n *sitter.Node) {
	if len(p.pf.ParseErrors) >= 20 {
		return
	}
	if n.IsError() || n.IsMissing() {
		msg := "syntax error"
		if n.IsMissing() {
			msg = "missing " + n.Kind()
		}
		p.pf.ParseErrors = append(p.pf.ParseErrors, model.Diagnostic{Severity: model.SevError, Code: "parse_error", Message: msg, Line: line(n), Language: "typescript"})
		p.errors = append(p.errors, line(n))
		return
	}
	if !n.HasError() {
		return
	}
	for i := uint(0); i < n.ChildCount(); i++ {
		p.collectErrors(n.Child(i))
	}
}

func (p *fileParser) inError(start, end int) bool {
	for _, l := range p.errors {
		if l >= start && l <= end {
			return true
		}
	}
	return false
}

func unquote(s string) string {
	if len(s) >= 2 && (s[0] == '"' || s[0] == '\'' || s[0] == '`') {
		return s[1 : len(s)-1]
	}
	return s
}

func (p *fileParser) walkTop(root *sitter.Node) {
	for i := uint(0); i < root.ChildCount(); i++ {
		p.statement(root.Child(i), false)
	}
	p.dynamicImports(root)
}

// statement handles one top-level statement; exported is set inside export.
func (p *fileParser) statement(n *sitter.Node, exported bool) {
	switch n.Kind() {
	case "import_statement":
		if s := n.ChildByFieldName("source"); s != nil {
			p.pf.Imports = append(p.pf.Imports, model.Import{Path: unquote(p.text(s)), Line: line(n)})
		}
	case "export_statement":
		if s := n.ChildByFieldName("source"); s != nil {
			p.pf.Imports = append(p.pf.Imports, model.Import{Path: unquote(p.text(s)), Line: line(n)})
		}
		if d := n.ChildByFieldName("declaration"); d != nil {
			p.statement(d, true)
		}
		if v := n.ChildByFieldName("value"); v != nil {
			// export default <expression>
			if fn := functionOf(v); fn != nil {
				p.callable(fn, n, "default", "", true)
			}
		}
	case "function_declaration", "generator_function_declaration":
		name := "default"
		if nn := n.ChildByFieldName("name"); nn != nil {
			name = p.text(nn)
		}
		p.callable(n, n, name, "", exported)
	case "class_declaration", "abstract_class_declaration":
		p.class(n, exported)
	case "interface_declaration", "type_alias_declaration", "enum_declaration":
		p.typeDecl(n, exported)
	case "lexical_declaration", "variable_declaration":
		for i := uint(0); i < n.ChildCount(); i++ {
			d := n.Child(i)
			if d.Kind() != "variable_declarator" {
				continue
			}
			nameNode, value := d.ChildByFieldName("name"), d.ChildByFieldName("value")
			if nameNode == nil || value == nil || nameNode.Kind() != "identifier" {
				continue
			}
			if fn := functionOf(value); fn != nil {
				p.callable(fn, n, p.text(nameNode), "", exported)
			}
		}
	case "ambient_declaration":
		// declare ... : declarations only, no bodies.
	}
}

// functionOf returns the function node for an arrow function, a function
// expression, or a call whose first argument is one (memo(), forwardRef()).
func functionOf(v *sitter.Node) *sitter.Node {
	for v != nil && (v.Kind() == "parenthesized_expression" || v.Kind() == "as_expression" || v.Kind() == "satisfies_expression") {
		v = v.NamedChild(0)
	}
	if v == nil {
		return nil
	}
	switch v.Kind() {
	case "arrow_function", "function_expression", "function", "generator_function":
		return v
	case "call_expression":
		args := v.ChildByFieldName("arguments")
		if args == nil {
			return nil
		}
		for i := uint(0); i < args.NamedChildCount(); i++ {
			if fn := functionOf(args.NamedChild(i)); fn != nil {
				return fn
			}
			break
		}
	}
	return nil
}

var spaceRe = regexp.MustCompile(`\s+`)

func (p *fileParser) normalize(n *sitter.Node) string {
	if n == nil {
		return ""
	}
	return spaceRe.ReplaceAllString(p.text(n), " ")
}

// callable records a function-like symbol. fn has the parameters and body;
// span is the node whose lines the symbol covers (the whole declaration).
func (p *fileParser) callable(fn, span *sitter.Node, name, container string, exported bool) {
	sym := model.Symbol{Name: name, Container: container, Kind: model.KindFunction, Exported: exported, StartLine: line(span), EndLine: endLine(span)}
	if container != "" {
		sym.Kind = model.KindMethod
	}
	params := fn.ChildByFieldName("parameters")
	if params == nil {
		params = fn.ChildByFieldName("parameter")
	}
	tp, rt, body := fn.ChildByFieldName("type_parameters"), fn.ChildByFieldName("return_type"), fn.ChildByFieldName("body")
	sym.Signature = p.normalize(tp) + p.normalize(params)
	if rt != nil {
		sym.Signature += p.normalize(rt)
	}
	var fb fingerprint.Builder
	for _, part := range []*sitter.Node{tp, params, rt, body} {
		if part != nil {
			p.tokens(part, &fb)
		}
	}
	sym.ExactFP, sym.StructFP, sym.Tokens, sym.MinHash = fb.Exact(), fb.Struct(), fb.Tokens(), fb.MinHash()
	m := &model.Metrics{Lines: sym.EndLine - sym.StartLine + 1}
	if body != nil {
		p.measure(body, 0, m)
	}
	m.Cyclomatic = 1 + m.Decisions
	sym.Metrics = m
	sym.InErrorRegion = p.inError(sym.StartLine, sym.EndLine)
	p.pf.Symbols = append(p.pf.Symbols, sym)
}

func (p *fileParser) class(n *sitter.Node, exported bool) {
	nameNode := n.ChildByFieldName("name")
	name := "default"
	if nameNode != nil {
		name = p.text(nameNode)
	}
	var fb fingerprint.Builder
	if h := n.ChildByFieldName("body"); h != nil {
		// Class type fingerprint over member names/shapes only.
		for i := uint(0); i < h.NamedChildCount(); i++ {
			if nn := h.NamedChild(i).ChildByFieldName("name"); nn != nil {
				fb.Add(p.text(nn), "ID")
			}
		}
	}
	p.pf.Symbols = append(p.pf.Symbols, model.Symbol{Name: name, Kind: model.KindType, Exported: exported, StartLine: line(n), EndLine: endLine(n),
		ExactFP: fb.Exact(), StructFP: fb.Struct(), Tokens: fb.Tokens(), Signature: "class", InErrorRegion: p.inError(line(n), endLine(n))})
	body := n.ChildByFieldName("body")
	if body == nil {
		return
	}
	for i := uint(0); i < body.NamedChildCount(); i++ {
		m := body.NamedChild(i)
		nn := m.ChildByFieldName("name")
		if nn == nil {
			continue
		}
		mname := p.text(nn)
		private := strings.HasPrefix(mname, "#")
		for j := uint(0); j < m.ChildCount(); j++ {
			if c := m.Child(j); c.Kind() == "accessibility_modifier" && (p.text(c) == "private" || p.text(c) == "protected") {
				private = true
			}
		}
		switch m.Kind() {
		case "method_definition":
			p.callable(m, m, mname, name, exported && !private)
		case "public_field_definition":
			if fn := functionOf(m.ChildByFieldName("value")); fn != nil {
				p.callable(fn, m, mname, name, exported && !private)
			}
		}
	}
}

func (p *fileParser) typeDecl(n *sitter.Node, exported bool) {
	nameNode := n.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	var fb fingerprint.Builder
	for _, f := range []string{"body", "value"} {
		if b := n.ChildByFieldName(f); b != nil {
			p.tokens(b, &fb)
		}
	}
	p.pf.Symbols = append(p.pf.Symbols, model.Symbol{Name: p.text(nameNode), Kind: model.KindType, Exported: exported, StartLine: line(n), EndLine: endLine(n),
		ExactFP: fb.Exact(), StructFP: fb.Struct(), Tokens: fb.Tokens(), Signature: strings.TrimSuffix(n.Kind(), "_declaration"), InErrorRegion: p.inError(line(n), endLine(n))})
}

// dynamicImports records import("x") and require("x") with literal specifiers.
func (p *fileParser) dynamicImports(n *sitter.Node) {
	if n.Kind() == "call_expression" {
		fn, args := n.ChildByFieldName("function"), n.ChildByFieldName("arguments")
		if fn != nil && args != nil && (fn.Kind() == "import" || (fn.Kind() == "identifier" && p.text(fn) == "require")) && args.NamedChildCount() > 0 {
			if a := args.NamedChild(0); a.Kind() == "string" {
				p.pf.Imports = append(p.pf.Imports, model.Import{Path: unquote(p.text(a)), Line: line(n)})
			}
		}
	}
	for i := uint(0); i < n.ChildCount(); i++ {
		p.dynamicImports(n.Child(i))
	}
}

var literalKinds = map[string]bool{"string": true, "template_string": true, "number": true, "regex": true, "jsx_text": true}

var identKinds = map[string]bool{
	"identifier": true, "property_identifier": true, "type_identifier": true, "shorthand_property_identifier": true,
	"shorthand_property_identifier_pattern": true, "private_property_identifier": true, "statement_identifier": true,
}

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
		case kind == "true" || kind == "false" || kind == "null" || kind == "undefined":
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
	"if_statement": true, "for_statement": true, "for_in_statement": true, "while_statement": true, "do_statement": true,
	"switch_statement": true, "catch_clause": true, "arrow_function": true, "function_expression": true, "function": true, "generator_function": true,
}

func (p *fileParser) measure(n *sitter.Node, depth int, m *model.Metrics) {
	kind := n.Kind()
	switch kind {
	case "if_statement", "for_statement", "for_in_statement", "while_statement", "do_statement", "switch_case", "catch_clause", "ternary_expression":
		m.Decisions++
	case "binary_expression":
		if op := n.ChildByFieldName("operator"); op != nil {
			if o := op.Kind(); o == "&&" || o == "||" || o == "??" {
				m.Decisions++
			}
		}
	}
	if strings.HasSuffix(kind, "_statement") && kind != "statement_block" && kind != "expression_statement" {
		m.Statements++
	} else if kind == "expression_statement" || kind == "lexical_declaration" || kind == "variable_declaration" {
		m.Statements++
	}
	childDepth := depth
	if nestingKinds[kind] {
		isElseIf := kind == "if_statement" && n.Parent() != nil && n.Parent().Kind() == "else_clause"
		if !isElseIf {
			childDepth = depth + 1
			if childDepth > m.MaxNesting {
				m.MaxNesting = childDepth
			}
		}
	}
	for i := uint(0); i < n.ChildCount(); i++ {
		p.measure(n.Child(i), childDepth, m)
	}
}

func (p *fileParser) refs(n *sitter.Node) {
	var visit func(n, parent *sitter.Node)
	visit = func(n, parent *sitter.Node) {
		switch n.Kind() {
		case "comment", "string", "template_string", "regex", "import_statement":
			return
		case "member_expression":
			obj, prop := n.ChildByFieldName("object"), n.ChildByFieldName("property")
			if obj != nil && prop != nil && obj.Kind() == "identifier" {
				p.pf.Refs[p.text(obj)+"."+p.text(prop)]++
			}
			if prop != nil {
				p.pf.Refs[p.text(prop)]++
			}
			if obj != nil {
				visit(obj, n)
			}
			return
		case "identifier", "type_identifier", "shorthand_property_identifier":
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
	case "function_declaration", "generator_function_declaration", "class_declaration", "abstract_class_declaration",
		"interface_declaration", "type_alias_declaration", "enum_declaration", "variable_declarator", "method_definition":
		if name := parent.ChildByFieldName("name"); name != nil && name.StartByte() == n.StartByte() {
			return true
		}
	}
	return false
}

// ---- dependency graph ----

var tsExts = []string{".ts", ".tsx", ".mts", ".cts", ".d.ts"}

type alias struct {
	prefix, suffix string // pattern split at '*'
	wild           bool
	targets        []string // repository-relative, with '*' placeholder
}

type tsconfig struct {
	dir     string
	aliases []alias
}

var lineComment = regexp.MustCompile(`(?m)^\s*//.*$|([^:"'])//[^"\n]*$`)
var blockComment = regexp.MustCompile(`(?s)/\*.*?\*/`)
var trailingComma = regexp.MustCompile(`,(\s*[}\]])`)

// jsonc decodes JSON with comments and trailing commas (tsconfig style).
func jsonc(data []byte, v any) error {
	s := blockComment.ReplaceAll(data, nil)
	s = lineComment.ReplaceAll(s, []byte("$1"))
	s = trailingComma.ReplaceAll(s, []byte("$1"))
	return json.Unmarshal(s, v)
}

type rawTSConfig struct {
	Extends         any `json:"extends"`
	CompilerOptions struct {
		BaseURL string              `json:"baseUrl"`
		Paths   map[string][]string `json:"paths"`
	} `json:"compilerOptions"`
}

func loadTSConfig(p string, read func(string) ([]byte, error), depth int) (baseDir string, paths map[string][]string) {
	data, err := read(p)
	if err != nil || depth > 5 {
		return "", nil
	}
	var rc rawTSConfig
	if jsonc(data, &rc) != nil {
		return "", nil
	}
	dir := path.Dir(p)
	if ext, ok := rc.Extends.(string); ok && (strings.HasPrefix(ext, "./") || strings.HasPrefix(ext, "../")) {
		ep := path.Clean(path.Join(dir, ext))
		if !strings.HasSuffix(ep, ".json") {
			ep += ".json"
		}
		baseDir, paths = loadTSConfig(ep, read, depth+1)
	}
	if rc.CompilerOptions.Paths != nil {
		base := dir
		if rc.CompilerOptions.BaseURL != "" {
			base = path.Clean(path.Join(dir, rc.CompilerOptions.BaseURL))
		}
		baseDir, paths = base, rc.CompilerOptions.Paths
	} else if rc.CompilerOptions.BaseURL != "" && baseDir == "" {
		baseDir = path.Clean(path.Join(dir, rc.CompilerOptions.BaseURL))
	}
	return baseDir, paths
}

// BuildGraph builds the directory-level import graph.
func (a *Adapter) BuildGraph(in lang.GraphInput) (*lang.Graph, []model.Diagnostic) {
	g := lang.NewGraph()
	g.Note = "TypeScript directory imports (non-test files): relative specifiers, tsconfig paths and workspace packages; node_modules are external"
	var diags []model.Diagnostic
	exists := map[string]bool{}
	for _, p := range in.AllPaths {
		exists[p] = true
	}

	// tsconfig aliases, applied to files under the config's directory.
	var configs []tsconfig
	// workspace packages: package.json "name" → directory.
	pkgs := map[string]string{}
	for _, p := range in.AllPaths {
		if strings.Contains(p, "node_modules/") {
			continue
		}
		base := path.Base(p)
		if strings.HasPrefix(base, "tsconfig") && strings.HasSuffix(base, ".json") {
			baseDir, paths := loadTSConfig(p, in.ReadFile, 0)
			if len(paths) == 0 {
				continue
			}
			c := tsconfig{dir: path.Dir(p)}
			keys := make([]string, 0, len(paths))
			for k := range paths {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				al := alias{}
				if i := strings.Index(k, "*"); i >= 0 {
					al.prefix, al.suffix, al.wild = k[:i], k[i+1:], true
				} else {
					al.prefix = k
				}
				for _, t := range paths[k] {
					al.targets = append(al.targets, path.Clean(path.Join(baseDir, t)))
				}
				c.aliases = append(c.aliases, al)
			}
			configs = append(configs, c)
		}
		if base == "package.json" {
			data, err := in.ReadFile(p)
			if err != nil {
				continue
			}
			var pj struct {
				Name string `json:"name"`
			}
			if json.Unmarshal(data, &pj) == nil && pj.Name != "" {
				pkgs[pj.Name] = path.Dir(p)
			}
		}
	}
	sort.Slice(configs, func(i, j int) bool { return len(configs[i].dir) > len(configs[j].dir) })

	resolveFile := func(cand string) string {
		cand = path.Clean(cand)
		if exists[cand] && isTS(cand) {
			return cand
		}
		stem := cand
		for _, js := range []string{".js", ".jsx", ".mjs", ".cjs"} {
			if strings.HasSuffix(cand, js) {
				stem = strings.TrimSuffix(cand, js)
			}
		}
		for _, e := range tsExts {
			if exists[stem+e] {
				return stem + e
			}
		}
		for _, e := range tsExts {
			if exists[cand+"/index"+e] {
				return cand + "/index" + e
			}
		}
		return ""
	}
	isAsset := func(spec string) bool {
		ext := path.Ext(spec)
		switch ext {
		case ".css", ".scss", ".sass", ".less", ".json", ".svg", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".woff", ".woff2", ".ttf", ".md", ".html", ".wasm", ".txt":
			return true
		}
		return strings.Contains(spec, "?")
	}
	inGraph := func(f lang.GraphFile) bool { return f.Class == model.ClassSource || f.Class == model.ClassGenerated }

	for _, f := range in.Files {
		if !inGraph(f) {
			continue
		}
		g.Nodes[path.Dir(f.Path)] = true
	}
	for _, f := range in.Files {
		if !inGraph(f) || f.Parsed == nil {
			continue
		}
		from := path.Dir(f.Path)
		for _, imp := range f.Parsed.Imports {
			spec := imp.Path
			target := ""
			switch {
			case strings.HasPrefix(spec, "./") || strings.HasPrefix(spec, "../") || spec == "." || spec == "..":
				if isAsset(spec) {
					continue
				}
				target = resolveFile(path.Join(from, spec))
				if target == "" {
					diags = append(diags, model.Diagnostic{Severity: model.SevInfo, Code: "import_unresolved", Message: "relative import has no indexed TypeScript file: " + spec, File: f.Path, Line: imp.Line, Language: "typescript"})
					continue
				}
			default:
				matched := false
				for _, c := range configs {
					if c.dir != "." && !strings.HasPrefix(f.Path, c.dir+"/") {
						continue
					}
					for _, al := range c.aliases {
						var star string
						if al.wild {
							if !strings.HasPrefix(spec, al.prefix) || !strings.HasSuffix(spec, al.suffix) || len(spec) < len(al.prefix)+len(al.suffix) {
								continue
							}
							star = spec[len(al.prefix) : len(spec)-len(al.suffix)]
						} else if spec != al.prefix {
							continue
						}
						matched = true
						external := true
						for _, t := range al.targets {
							if !strings.Contains(t, "node_modules/") {
								external = false
							}
						}
						if external {
							// An alias into node_modules (e.g. react → node_modules/@types/react)
							// names an external package, not internal code.
							matched = false
							break
						}
						for _, t := range al.targets {
							if target = resolveFile(strings.Replace(t, "*", star, 1)); target != "" {
								break
							}
						}
						break
					}
					if matched {
						break
					}
				}
				if !matched {
					name, sub := packageName(spec)
					if dir, ok := pkgs[name]; ok {
						matched = true
						for _, cand := range []string{path.Join(dir, sub), path.Join(dir, "src", sub)} {
							if target = resolveFile(cand); target != "" {
								break
							}
						}
					}
				}
				if !matched {
					if isAsset(spec) {
						continue
					}
					ext := g.External[from]
					if ext == nil {
						ext = map[string]bool{}
						g.External[from] = ext
					}
					name, _ := packageName(spec)
					ext[name] = true
					continue
				}
				if target == "" {
					diags = append(diags, model.Diagnostic{Severity: model.SevInfo, Code: "import_unresolved", Message: "aliased or workspace import has no indexed TypeScript file: " + spec, File: f.Path, Line: imp.Line, Language: "typescript"})
					continue
				}
			}
			g.AddEdge(from, path.Dir(target), model.Location{File: f.Path, Line: imp.Line})
		}
	}
	return g, diags
}

func isTS(p string) bool {
	for _, e := range tsExts {
		if strings.HasSuffix(p, e) {
			return true
		}
	}
	return false
}

// packageName splits a bare specifier into package name and subpath.
func packageName(spec string) (string, string) {
	parts := strings.Split(spec, "/")
	if strings.HasPrefix(spec, "@") && len(parts) >= 2 {
		return parts[0] + "/" + parts[1], strings.Join(parts[2:], "/")
	}
	return parts[0], strings.Join(parts[1:], "/")
}
