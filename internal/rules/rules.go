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
	Note  string
	Check func(*Context)
}

// Registry holds rule definitions and their language implementations.
type Registry struct {
	rules []Rule
	byID  map[string]int
	impls map[string]map[string]Impl // rule -> language -> impl
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

// File is a source or test file given to a rule, with its syntax tree.
type File struct {
	Path   string
	Class  model.FileClass
	Parsed *model.ParsedFile
	Src    []byte
	Tree   *sitter.Tree
}

// Root returns the root syntax node.
func (f *File) Root() *sitter.Node { return f.Tree.RootNode() }

// Test reports whether the file is a test file.
func (f *File) Test() bool { return f.Class == model.ClassTest }

// Text returns a node's source text.
func (f *File) Text(n *sitter.Node) string {
	if n == nil {
		return ""
	}
	return n.Utf8Text(f.Src)
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
