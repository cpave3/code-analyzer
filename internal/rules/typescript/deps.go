package typescript

import (
	"encoding/json"
	"path"
	"sort"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"

	tslang "github.com/cameron/code-evaluator/internal/lang/typescript"
	"github.com/cameron/code-evaluator/internal/model"
	"github.com/cameron/code-evaluator/internal/rules"
)

func resolver(c *rules.Context) *tslang.Resolver {
	return c.Memo("ts:resolver", func() any { return tslang.NewResolver(c.AllPaths, c.ReadFile) }).(*tslang.Resolver)
}

// ---- unused exports ----

type importIndex struct {
	imported map[string]bool            // file is imported by some other file
	wildcard map[string]bool            // namespace / export * / dynamic: any export may be used
	names    map[string]map[string]bool // file -> imported export names
}

// importFact is one import, re-export or dynamic import of a file.
type importFact struct {
	spec     string
	names    []string // imported or re-exported names ("default" for a default import)
	wildcard bool     // namespace import, export *, or dynamic import(): any export may be used
}

// fileImports lists a file's import facts, cached by content.
func fileImports(f *rules.File) []importFact {
	return f.Memo("ts:imports", func() any {
		var out []importFact
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			switch n.Kind() {
			case "import_statement", "export_statement":
				src := n.ChildByFieldName("source")
				if src == nil {
					return n.Kind() == "export_statement" // export declarations may contain dynamic imports
				}
				fact := importFact{spec: strings.Trim(f.Text(src), "\"'`")}
				if n.Kind() == "export_statement" && !rules.Contains(n, "export_clause") {
					fact.wildcard = true // export * from
				}
				rules.Walk(n, func(x *sitter.Node) bool {
					switch x.Kind() {
					case "namespace_import", "namespace_export":
						fact.wildcard = true
					case "import_specifier", "export_specifier":
						fact.names = append(fact.names, f.Text(x.ChildByFieldName("name")))
						return false
					case "import_clause":
						for _, ch := range rules.Named(x) {
							if ch.Kind() == "identifier" {
								fact.names = append(fact.names, "default")
							}
						}
					}
					return true
				})
				out = append(out, fact)
				return false
			case "call_expression":
				fn, args := n.ChildByFieldName("function"), rules.Named(n.ChildByFieldName("arguments"))
				if fn != nil && len(args) > 0 && args[0].Kind() == "string" && (fn.Kind() == "import" || f.Text(fn) == "require") {
					out = append(out, importFact{spec: strings.Trim(f.Text(args[0]), "\"'`"), wildcard: true})
				}
			}
			return true
		})
		return out
	}).([]importFact)
}

func buildImportIndex(c *rules.Context) *importIndex {
	return c.Memo("ts:importIndex", func() any {
		r := resolver(c)
		ix := &importIndex{imported: map[string]bool{}, wildcard: map[string]bool{}, names: map[string]map[string]bool{}}
		for _, f := range c.Files {
			for _, fact := range fileImports(f) {
				t := r.Resolve(f.Path, fact.spec).Target
				if t == "" || t == f.Path {
					continue
				}
				ix.imported[t] = true
				if fact.wildcard {
					ix.wildcard[t] = true
				}
				if len(fact.names) > 0 && ix.names[t] == nil {
					ix.names[t] = map[string]bool{}
				}
				for _, n := range fact.names {
					ix.names[t][n] = true
				}
			}
		}
		return ix
	}).(*importIndex)
}

type exportFact struct {
	name string
	line int
	uses int // identifier occurrences of the name in the module (the declaration is one)
}

// fileExports lists a module's named function exports, cached by content.
func fileExports(f *rules.File) []exportFact {
	return f.Memo("ts:exports", func() any {
		var out []exportFact
		root := f.Root()
		if root == nil {
			return out
		}
		uses := countIdents(f, root, useKinds)
		for _, st := range rules.Named(root) {
			if st.Kind() != "export_statement" {
				continue
			}
			d := st.ChildByFieldName("declaration")
			if d == nil || isDefaultExport(st) {
				continue // default exports follow framework conventions; not judged
			}
			var names []*sitter.Node
			switch d.Kind() {
			case "function_declaration", "generator_function_declaration":
				names = append(names, d.ChildByFieldName("name"))
			case "lexical_declaration":
				for _, v := range rules.Named(d) {
					if v.Kind() == "variable_declarator" && functionValue(v.ChildByFieldName("value")) {
						names = append(names, v.ChildByFieldName("name"))
					}
				}
			}
			for _, nm := range names {
				if nm != nil && nm.Kind() == "identifier" {
					out = append(out, exportFact{f.Text(nm), rules.Line(nm), uses[f.Text(nm)]})
				}
			}
		}
		return out
	}).([]exportFact)
}

// fileAmbient lists `declare module "x"` names, cached by content.
func fileAmbient(f *rules.File) []string {
	return f.Memo("ts:ambient", func() any {
		var out []string
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() == "module" {
				if nm := n.ChildByFieldName("name"); nm != nil && nm.Kind() == "string" {
					out = append(out, strings.Trim(f.Text(nm), "\"'"))
				}
			}
			return n.Kind() == "program" || n.Kind() == "ambient_declaration" || n.Kind() == "export_statement" || n.Kind() == "statement_block"
		})
		return out
	}).([]string)
}

// prepare warms the per-file facts the cross-file rules read.
func prepare(f *rules.File) { fileImports(f); fileExports(f); fileAmbient(f) }

func unusedExport(c *rules.Context) {
	ix := buildImportIndex(c)
	for _, f := range c.Files {
		if f.Test() || isDeclarationFile(f.Path) || !ix.imported[f.Path] || ix.wildcard[f.Path] {
			continue
		}
		for _, e := range fileExports(f) {
			if ix.names[f.Path][e.name] || e.uses > 1 {
				continue
			}
			c.Report(rules.Finding{File: f.Path, Line: e.line, Symbol: e.name, Key: "export:" + e.name,
				Message: "exported function " + e.name + " is imported by no other file and unused in its module (other files import this module by name)"})
		}
	}
}

// ---- unresolved imports ----

var jsLike = []string{"", ".js", ".jsx", ".mjs", ".cjs", ".json", ".d.ts", "/index.js", "/index.jsx", "/index.mjs"}

func anyFile(r *tslang.Resolver, cand string) bool {
	cand = path.Clean(cand)
	for _, e := range jsLike {
		if r.Exists(cand + e) {
			return true
		}
	}
	return false
}

func unresolvedImport(c *rules.Context) {
	r := resolver(c)
	for _, f := range c.Files {
		if f.Parsed == nil {
			continue
		}
		for _, imp := range f.Parsed.Imports {
			spec := imp.Path
			clean := spec
			if i := strings.IndexAny(clean, "?#"); i >= 0 {
				clean = clean[:i]
			}
			res := r.Resolve(f.Path, spec)
			switch {
			case tslang.IsRelative(spec):
				if res.Target != "" || anyFile(r, path.Join(path.Dir(f.Path), clean)) {
					continue
				}
				c.Report(rules.Finding{File: f.Path, Line: imp.Line, Key: "import:" + spec,
					Message: "relative import \"" + spec + "\" resolves to no file in the snapshot"})
			case res.Kind == tslang.KindAlias || res.Kind == tslang.KindWorkspace:
				if res.Target != "" {
					continue
				}
				c.Report(rules.Finding{File: f.Path, Line: imp.Line, Key: "import:" + spec, Severity: rules.Advisory,
					Message: res.Kind + " import \"" + spec + "\" resolves to no TypeScript file in the snapshot (it may point at build output or JavaScript)"})
			}
		}
	}
}

// ---- manifests ----

type manifest struct {
	file     string
	dir      string
	name     string
	declared map[string]bool
	deps     []string // "dependencies" only
	scripts  string
	raw      string
}

type manifests struct {
	byDir   map[string]*manifest
	list    []*manifest
	ambient map[string]bool // declare module "x"
}

func loadManifests(c *rules.Context) *manifests {
	return c.Memo("ts:manifests", func() any {
		ms := &manifests{byDir: map[string]*manifest{}, ambient: map[string]bool{}}
		for _, p := range c.AllPaths {
			if path.Base(p) != "package.json" || strings.Contains(p, "node_modules/") {
				continue
			}
			data, err := c.ReadFile(p)
			if err != nil {
				continue
			}
			var pj struct {
				Name     string            `json:"name"`
				Deps     map[string]string `json:"dependencies"`
				Dev      map[string]string `json:"devDependencies"`
				Peer     map[string]string `json:"peerDependencies"`
				Optional map[string]string `json:"optionalDependencies"`
				Scripts  map[string]string `json:"scripts"`
			}
			if err := json.Unmarshal(data, &pj); err != nil {
				c.Diagnose(model.Diagnostic{Severity: model.SevWarning, Code: "package_json_invalid", Message: err.Error(), File: p})
				continue
			}
			m := &manifest{file: p, dir: path.Dir(p), name: pj.Name, declared: map[string]bool{}, raw: string(data)}
			for _, set := range []map[string]string{pj.Deps, pj.Dev, pj.Peer, pj.Optional} {
				for k := range set {
					m.declared[k] = true
				}
			}
			for k := range pj.Deps {
				m.deps = append(m.deps, k)
			}
			sort.Strings(m.deps)
			for _, s := range pj.Scripts {
				m.scripts += s + "\n"
			}
			ms.byDir[m.dir] = m
			ms.list = append(ms.list, m)
		}
		for _, f := range c.Files {
			for _, name := range fileAmbient(f) {
				ms.ambient[name] = true
			}
		}
		return ms
	}).(*manifests)
}

// chain returns the manifests from a file's directory up to the root.
func (ms *manifests) chain(file string) []*manifest {
	var out []*manifest
	for dir := path.Dir(file); ; dir = path.Dir(dir) {
		if m := ms.byDir[dir]; m != nil {
			out = append(out, m)
		}
		if dir == "." || dir == "/" {
			return out
		}
	}
}

var nodeBuiltins = map[string]bool{}

func init() {
	for _, b := range strings.Fields(`assert async_hooks buffer child_process cluster console constants crypto dgram diagnostics_channel dns domain events fs
		http http2 https inspector module net os path perf_hooks process punycode querystring readline repl stream string_decoder sys timers tls
		trace_events tty url util v8 vm wasi worker_threads zlib test sqlite`) {
		nodeBuiltins[b] = true
	}
}

func typesName(pkg string) string {
	if strings.HasPrefix(pkg, "@") {
		return "@types/" + strings.Replace(strings.TrimPrefix(pkg, "@"), "/", "__", 1)
	}
	return "@types/" + pkg
}

func exempt(spec, name string) bool {
	if strings.Contains(spec, ":") || strings.HasPrefix(spec, "~") || strings.HasPrefix(spec, "$") || strings.HasPrefix(spec, "#") || strings.HasPrefix(spec, "/") {
		return true // node:, bun:, virtual:, URLs, subpath imports, absolute paths
	}
	if name == "@" || strings.HasPrefix(name, "@/") || nodeBuiltins[strings.Split(name, "/")[0]] {
		return true
	}
	return false
}

func undeclaredDependency(c *rules.Context) {
	r := resolver(c)
	ms := loadManifests(c)
	for _, f := range c.Files {
		if f.Parsed == nil {
			continue
		}
		chain := ms.chain(f.Path)
		if len(chain) == 0 {
			continue
		}
		for _, imp := range f.Parsed.Imports {
			res := r.Resolve(f.Path, imp.Path)
			if res.Kind != tslang.KindExternal || exempt(imp.Path, res.Package) || ms.ambient[imp.Path] || ms.ambient[res.Package] {
				continue
			}
			declared := false
			for _, m := range chain {
				if m.declared[res.Package] || m.declared[typesName(res.Package)] || m.name == res.Package {
					declared = true
					break
				}
			}
			if declared {
				continue
			}
			c.Report(rules.Finding{File: f.Path, Line: imp.Line, Key: "import:" + res.Package,
				Message: "package \"" + res.Package + "\" is imported but not declared in " + chain[0].file + " or any package.json above it"})
		}
	}
}

var textExts = map[string]bool{".js": true, ".jsx": true, ".mjs": true, ".cjs": true, ".vue": true, ".svelte": true, ".astro": true,
	".css": true, ".scss": true, ".sass": true, ".less": true, ".html": true, ".json": true, ".yaml": true, ".yml": true, ".toml": true}

func unusedDependency(c *rules.Context) {
	ms := loadManifests(c)
	imports := importedPackages(c)
	for _, m := range ms.list {
		if len(m.deps) == 0 || !m.governsAny(c.Files) {
			continue // no checked TypeScript (e.g. fixtures under testdata/): nothing to compare against
		}
		pending := m.unimported(imports)
		pending = m.unmentioned(c, pending)
		for _, d := range pending {
			c.Report(rules.Finding{File: m.file, Line: m.depLine(d), Key: "dependency:" + d,
				Message: "dependency \"" + d + "\" is declared in " + m.file + " but nothing under " + m.dir + " imports or names it"})
		}
	}
}

// importedPackages maps each file to the package names it imports.
func importedPackages(c *rules.Context) map[string][]string {
	out := map[string][]string{}
	for _, f := range c.Files {
		if f.Parsed == nil {
			continue
		}
		for _, imp := range f.Parsed.Imports {
			name, _ := tslang.PackageName(imp.Path)
			out[f.Path] = append(out[f.Path], name)
		}
	}
	return out
}

func (m *manifest) under(p string) bool { return m.dir == "." || strings.HasPrefix(p, m.dir+"/") }

func (m *manifest) governsAny(files []*rules.File) bool {
	for _, f := range files {
		if m.under(f.Path) {
			return true
		}
	}
	return false
}

// unimported lists runtime dependencies no governed TypeScript file imports
// (@types packages and names used in package scripts are skipped).
func (m *manifest) unimported(imports map[string][]string) []string {
	used := map[string]bool{}
	for p, names := range imports {
		if m.under(p) {
			for _, n := range names {
				used[n] = true
			}
		}
	}
	var out []string
	for _, d := range m.deps {
		if !used[d] && !strings.HasPrefix(d, "@types/") && !strings.Contains(m.scripts, d) {
			out = append(out, d)
		}
	}
	return out
}

// unmentioned keeps the dependencies no other text file under the package
// names in quotes (JS, CSS, HTML, configs).
func (m *manifest) unmentioned(c *rules.Context, pending []string) []string {
	for _, p := range c.AllPaths {
		if len(pending) == 0 {
			break
		}
		if !m.under(p) || p == m.file || strings.Contains(p, "node_modules/") || !textExts[path.Ext(p)] || path.Base(p) == "package-lock.json" {
			continue
		}
		data, err := c.ReadFile(p)
		if err != nil || len(data) > 512<<10 {
			continue
		}
		pending = withoutMentioned(string(data), pending)
	}
	return pending
}

func withoutMentioned(text string, deps []string) []string {
	kept := deps[:0]
	for _, d := range deps {
		if !(strings.Contains(text, "\""+d) || strings.Contains(text, "'"+d) || strings.Contains(text, "`"+d)) {
			kept = append(kept, d)
		}
	}
	return kept
}

// depLine is the line of a dependency in the "dependencies" block.
func (m *manifest) depLine(d string) int {
	if i := strings.Index(m.raw, "\"dependencies\""); i >= 0 {
		if j := strings.Index(m.raw[i:], "\""+d+"\""); j >= 0 {
			return strings.Count(m.raw[:i+j], "\n") + 1
		}
	}
	return 1
}

func isDefaultExport(st *sitter.Node) bool {
	for i := uint(0); i < st.ChildCount(); i++ {
		if st.Child(i).Kind() == "default" {
			return true
		}
	}
	return false
}
