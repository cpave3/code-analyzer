// Package lang defines the language adapter contract and its registry.
//
// A language adapter turns source files into the shared model (symbols,
// ranges, normalized fingerprints, imports, name-based references) and builds
// the dependency graph it can honestly resolve. Every adapter declares its
// capabilities; analyzers that need a capability an adapter lacks report
// partial or unsupported coverage instead of a neutral score.
//
// To add a language (for example TypeScript or C#), implement Adapter in a new
// package and register it in internal/languages. No analyzer changes are
// needed.
package lang

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/model"
)

// Capability is something an adapter can extract.
type Capability string

const (
	CapSymbols      Capability = "symbols"      // declarations with ranges
	CapComplexity   Capability = "complexity"   // per-callable metrics
	CapFingerprints Capability = "fingerprints" // exact/structural fingerprints + minhash
	CapDependencies Capability = "dependencies" // module/package dependency graph
	CapReferences   Capability = "references"   // name-based identifier references
)

// AllCapabilities in stable order.
var AllCapabilities = []Capability{CapSymbols, CapComplexity, CapFingerprints, CapDependencies, CapReferences}

// Adapter is the language plug-in contract.
type Adapter interface {
	// ID is a short stable identifier, e.g. "go".
	ID() string
	// Version changes whenever parse output changes; it keys the parse cache.
	Version() string
	// Extensions handled by this adapter, including the dot (".go").
	Extensions() []string
	// Capabilities maps each supported capability to a fidelity note that is
	// copied into reports (e.g. "name-based, unresolved").
	Capabilities() map[Capability]string
	// Limits lists what the adapter deliberately does not resolve.
	Limits() []string
	// IsTestPath classifies test files by path.
	IsTestPath(p string) bool
	// PackageOf returns the grouping unit used for symbol identity and the
	// dependency graph (a directory for Go, a file/module for TypeScript).
	PackageOf(p string) string
	// Parse extracts a path-independent result from file content. A syntax
	// error is reported in ParsedFile.ParseErrors, not as an error return;
	// an error return means the adapter could not run at all.
	Parse(src []byte) (*model.ParsedFile, error)
	// BuildGraph builds the dependency graph over the snapshot's files of this
	// language.
	BuildGraph(in GraphInput) (*Graph, []model.Diagnostic)
}

// GraphFile is a parsed, classified file given to BuildGraph.
type GraphFile struct {
	Path   string
	Class  model.FileClass
	Parsed *model.ParsedFile
}

// GraphInput gives an adapter what it needs to build a graph without access
// to the repository: its files, all paths in the snapshot, and a reader for
// build files (go.mod, tsconfig.json, *.csproj ...).
type GraphInput struct {
	Files    []GraphFile
	AllPaths []string
	ReadFile func(p string) ([]byte, error)
}

// Graph is a directed dependency graph between internal units.
type Graph struct {
	Nodes    map[string]bool
	Edges    map[string]map[string][]model.Location // from -> to -> import sites
	External map[string]map[string]bool             // from -> external dependency
	Note     string                                 // fidelity statement
}

// NewGraph returns an empty graph.
func NewGraph() *Graph {
	return &Graph{Nodes: map[string]bool{}, Edges: map[string]map[string][]model.Location{}, External: map[string]map[string]bool{}}
}

// AddEdge records from -> to with an evidence location.
func (g *Graph) AddEdge(from, to string, loc model.Location) {
	if from == to {
		return
	}
	g.Nodes[from], g.Nodes[to] = true, true
	m := g.Edges[from]
	if m == nil {
		m = map[string][]model.Location{}
		g.Edges[from] = m
	}
	m[to] = append(m[to], loc)
}

// Declared is a language the registry recognises but has no adapter for.
// Files in it are counted and reported as unsupported, never scored.
type Declared struct {
	ID         string
	Extensions []string
	Note       string
}

// Registry holds adapters and declared-but-unsupported languages.
type Registry struct {
	adapters []Adapter
	declared []Declared
	byExt    map[string]Adapter
	declExt  map[string]string
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{byExt: map[string]Adapter{}, declExt: map[string]string{}}
}

// Register adds an adapter. Duplicate IDs or extensions are rejected.
func (r *Registry) Register(a Adapter) error {
	for _, e := range r.adapters {
		if e.ID() == a.ID() {
			return fmt.Errorf("language %q already registered", a.ID())
		}
	}
	for _, ext := range a.Extensions() {
		ext = strings.ToLower(ext)
		if other, ok := r.byExt[ext]; ok {
			return fmt.Errorf("extension %s already handled by %s", ext, other.ID())
		}
	}
	for _, ext := range a.Extensions() {
		ext = strings.ToLower(ext)
		r.byExt[ext] = a
		delete(r.declExt, ext) // a real adapter supersedes a declaration
	}
	r.adapters = append(r.adapters, a)
	return nil
}

// Declare records a known language without an adapter.
func (r *Registry) Declare(d Declared) {
	r.declared = append(r.declared, d)
	for _, ext := range d.Extensions {
		ext = strings.ToLower(ext)
		if _, ok := r.byExt[ext]; !ok {
			r.declExt[ext] = d.ID
		}
	}
}

// Adapters returns registered adapters sorted by ID.
func (r *Registry) Adapters() []Adapter {
	out := append([]Adapter{}, r.adapters...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// DeclaredLanguages returns declared unsupported languages sorted by ID,
// excluding ones that later gained an adapter.
func (r *Registry) DeclaredLanguages() []Declared {
	var out []Declared
	for _, d := range r.declared {
		if r.Adapter(d.ID) == nil {
			out = append(out, d)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Adapter returns the adapter with the given ID, or nil.
func (r *Registry) Adapter(id string) Adapter {
	for _, a := range r.adapters {
		if a.ID() == id {
			return a
		}
	}
	return nil
}

// Detect returns the language ID for a path and its adapter (nil when the
// language is declared but unsupported, or unknown with id "").
func (r *Registry) Detect(p string) (string, Adapter) {
	ext := strings.ToLower(path.Ext(p))
	if a, ok := r.byExt[ext]; ok {
		return a.ID(), a
	}
	if id, ok := r.declExt[ext]; ok {
		return id, nil
	}
	return "", nil
}

// Has reports whether an adapter supports a capability.
func Has(a Adapter, c Capability) bool {
	if a == nil {
		return false
	}
	_, ok := a.Capabilities()[c]
	return ok
}
