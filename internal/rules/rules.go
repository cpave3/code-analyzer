// Package rules is the defect-rule framework: concrete, provable defects
// (dead code, lost errors, dependency integrity, architecture policy,
// test-oracle defects), reported as findings rather than folded into a score.
//
// A Rule is defined once, language-neutrally (ID, category, severity,
// summary, limits). Each language contributes its own Impl of a rule: a Check
// function that inspects that language's files through a Context. Rules and
// implementations are registered separately, so a new rule can start with one
// language and gain others later, and a new language can implement any subset
// of existing rules. The CLI lists the rule × language matrix (`codeeval
// rules`), and every finding names the rule and language implementation that
// produced it.
//
// To add a rule: Define it in internal/rules/catalog, then add an Impl in each
// language package that can support it (internal/rules/golang,
// internal/rules/typescript, or internal/rules/generic for checks that only
// need the dependency graph). Register nothing else: internal/rules/all wires
// the packages together.
package rules

import (
	"fmt"
	"sort"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/memo"
	"github.com/cameron/code-evaluator/internal/model"
)

// Version identifies the rule set's semantics. Bump it whenever a rule's
// detection changes so results stay comparable.
const Version = "rules/v1"

// Severity says how a finding should be used.
type Severity string

const (
	// Gate findings are proven from syntax and repository files; a change
	// that introduces one fails `codeeval defects` by default.
	Gate Severity = "gate"
	// Advisory findings need a human decision (ambiguous fallbacks, unused
	// public API, dependencies used outside source).
	Advisory Severity = "advisory"
)

// Categories.
const (
	CatDeadCode     = "dead-code"
	CatLostErrors   = "lost-errors"
	CatDependencies = "dependencies"
	CatArchitecture = "architecture"
	CatTestOracle   = "test-oracle"
)

// Rule is the language-neutral definition of a defect rule.
type Rule struct {
	ID       string   `json:"id"` // "<category>/<name>"
	Category string   `json:"category"`
	Severity Severity `json:"severity"` // default; an Impl may lower individual findings to advisory
	Summary  string   `json:"summary"`
	Limits   []string `json:"limits,omitempty"`
}

// Impl is one language's implementation of a rule.
type Impl struct {
	Rule     string
	Language string
	// Note states what this implementation covers and misses, in words that
	// go into reports.
	Note string
	// FileLocal marks a check whose findings for a file depend only on that
	// file's content. Such checks run file by file and their findings are
	// cached by content, so unchanged files are not re-checked.
	FileLocal bool
	Check     func(*Context)
}

// Registry holds rule definitions and their language implementations.
type Registry struct {
	rules []Rule
	byID  map[string]int
	impls map[string]map[string]Impl // rule -> language -> impl
	prep  map[string][]func(*File)   // language -> per-file fact builders
}

// Prepare registers a per-file step for a language that warms the cached
// facts its cross-file rules read (identifier counts, declarations,
// imports ...). Run calls it once per file, in the same pass as the
// file-local checks, so a file is parsed at most once per run.
func (r *Registry) Prepare(language string, fn func(*File)) {
	if r.prep == nil {
		r.prep = map[string][]func(*File){}
	}
	r.prep[language] = append(r.prep[language], fn)
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{byID: map[string]int{}, impls: map[string]map[string]Impl{}}
}

// Define adds a rule definition.
func (r *Registry) Define(rule Rule) error {
	if rule.ID == "" || rule.Category == "" || rule.Summary == "" {
		return fmt.Errorf("rule %q: id, category and summary are required", rule.ID)
	}
	if !strings.HasPrefix(rule.ID, rule.Category+"/") {
		return fmt.Errorf("rule %q: id must start with its category %q", rule.ID, rule.Category)
	}
	if rule.Severity != Gate && rule.Severity != Advisory {
		return fmt.Errorf("rule %q: severity must be gate or advisory", rule.ID)
	}
	if _, ok := r.byID[rule.ID]; ok {
		return fmt.Errorf("rule %q already defined", rule.ID)
	}
	r.byID[rule.ID] = len(r.rules)
	r.rules = append(r.rules, rule)
	return nil
}

// Implement adds a language implementation of a defined rule.
func (r *Registry) Implement(impl Impl) error {
	if _, ok := r.byID[impl.Rule]; !ok {
		return fmt.Errorf("implementation for undefined rule %q", impl.Rule)
	}
	if impl.Language == "" || impl.Check == nil {
		return fmt.Errorf("rule %q: implementation needs a language and a check", impl.Rule)
	}
	m := r.impls[impl.Rule]
	if m == nil {
		m = map[string]Impl{}
		r.impls[impl.Rule] = m
	}
	if _, ok := m[impl.Language]; ok {
		return fmt.Errorf("rule %q already implemented for %s", impl.Rule, impl.Language)
	}
	m[impl.Language] = impl
	return nil
}

// Rules returns the definitions sorted by ID.
func (r *Registry) Rules() []Rule {
	out := append([]Rule{}, r.rules...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Rule returns a definition by ID.
func (r *Registry) Rule(id string) (Rule, bool) {
	i, ok := r.byID[id]
	if !ok {
		return Rule{}, false
	}
	return r.rules[i], true
}

// Impl returns a rule's implementation for a language.
func (r *Registry) Impl(rule, language string) (Impl, bool) {
	impl, ok := r.impls[rule][language]
	return impl, ok
}

// Languages returns the languages implementing a rule, sorted.
func (r *Registry) Languages(rule string) []string {
	var out []string
	for l := range r.impls[rule] {
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}

// Finding is one defect.
type Finding struct {
	Rule     string   `json:"rule"`
	Category string   `json:"category"`
	Severity Severity `json:"severity"`
	Language string   `json:"language"`
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Symbol   string   `json:"symbol,omitempty"` // enclosing declaration
	Message  string   `json:"message"`
	// Key identifies the finding within its file without line numbers, so a
	// finding that only moved is matched across revisions.
	Key string `json:"key"`
	// Change is set when comparing revisions: introduced | resolved.
	Change string `json:"change,omitempty"`
}

// File is a source or test file given to a rule. Its content and syntax
// tree are loaded on first use and released after the file's checks, so a
// run holds about one tree at a time.
type File struct {
	Path     string
	Blob     string // content id; keys cached per-file results ("" = no caching)
	Class    model.FileClass
	Parsed   *model.ParsedFile
	Language string

	src    []byte
	tree   *sitter.Tree
	loaded bool
	load   func() ([]byte, *sitter.Tree, error)
	cache  *memo.Cache
	local  map[string]any
}

// NewFile returns a file with content and tree already loaded (tests and
// tools that check files outside a snapshot).
func NewFile(path string, class model.FileClass, parsed *model.ParsedFile, src []byte, tree *sitter.Tree) *File {
	return &File{Path: path, Class: class, Parsed: parsed, src: src, tree: tree, loaded: true}
}

func (f *File) ensure() {
	if f.loaded {
		return
	}
	f.loaded = true
	if f.load != nil {
		f.src, f.tree, _ = f.load()
	}
}

// Source returns the file's content.
func (f *File) Source() []byte { f.ensure(); return f.src }

// Root returns the root syntax node (nil when the file could not be parsed).
func (f *File) Root() *sitter.Node {
	f.ensure()
	if f.tree == nil {
		return nil
	}
	return f.tree.RootNode()
}

// release frees the syntax tree; the next use reloads it.
func (f *File) release() {
	if f.tree != nil {
		f.tree.Close()
	}
	f.src, f.tree, f.loaded = nil, nil, f.load == nil
}

// Memo returns build()'s value for this file's content, computed once and
// shared across runs through the cache (consecutive commits share most
// files). key must name everything the value depends on besides the
// content (the language is added).
func (f *File) Memo(key string, build func() any) any {
	if f.Blob != "" && f.cache != nil {
		return f.cache.Get(f.Language+"\x00"+f.Blob+"\x00"+key, build)
	}
	if v, ok := f.local[key]; ok {
		return v
	}
	v := build()
	if f.local == nil {
		f.local = map[string]any{}
	}
	f.local[key] = v
	return v
}

// Test reports whether the file is a test file.
func (f *File) Test() bool { return f.Class == model.ClassTest }

// Text returns a node's source text.
func (f *File) Text(n *sitter.Node) string {
	if n == nil {
		return ""
	}
	return n.Utf8Text(f.Source())
}

// Context is what a rule implementation sees: the files of one language in a
// snapshot (source, test and generated files without syntax errors; findings
// in generated files are dropped; vendored and excluded files are absent), every path in the snapshot, a reader
// for any file (manifests, configs), the language's dependency graph and the
// effective configuration.
type Context struct {
	Rule     Rule
	Language string
	Files    []*File
	AllPaths []string
	ReadFile func(string) ([]byte, error)
	Graph    *lang.Graph // nil when the adapter builds none
	Config   *config.Config

	findings []Finding
	diags    []model.Diagnostic
	cache    map[string]any
}

// Report records a finding. Rule, category, language and (when empty)
// severity are filled in from the context.
func (c *Context) Report(f Finding) {
	f.Rule, f.Category, f.Language = c.Rule.ID, c.Rule.Category, c.Language
	if f.Severity == "" {
		f.Severity = c.Rule.Severity
	}
	if f.Key == "" {
		f.Key = f.Message
	}
	c.findings = append(c.findings, f)
}

// Diagnose records a limitation met while checking (unreadable manifest ...).
func (c *Context) Diagnose(d model.Diagnostic) {
	d.Analyzer, d.Language = "rules:"+c.Rule.ID, c.Language
	c.diags = append(c.diags, d)
}

// Memo returns a value shared by all rules of one language in one run, built
// on first use (a resolver, parsed manifests ...).
func (c *Context) Memo(key string, build func() any) any {
	if v, ok := c.cache[key]; ok {
		return v
	}
	v := build()
	c.cache[key] = v
	return v
}

// NewContext builds a context outside a snapshot run, for unit tests and
// tools that check individual files. read may be nil.
func NewContext(rule Rule, language string, files []*File, allPaths []string, read func(string) ([]byte, error), cfg *config.Config) *Context {
	if read == nil {
		read = func(p string) ([]byte, error) { return nil, fmt.Errorf("%s: no reader", p) }
	}
	return &Context{Rule: rule, Language: language, Files: files, AllPaths: allPaths, ReadFile: read, Config: cfg, cache: map[string]any{}}
}

// Findings returns what the context's checks reported.
func (c *Context) Findings() []Finding { return c.findings }

// At reports a finding at node n of file f inside symbol, keyed
// "symbol|key". sev "" uses the rule's severity.
func (c *Context) At(f *File, n *sitter.Node, symbol, key, message string, sev Severity) {
	c.Report(Finding{File: f.Path, Line: Line(n), Symbol: symbol, Key: symbol + "|" + key, Message: message, Severity: sev})
}
