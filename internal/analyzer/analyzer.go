// Package analyzer defines the analyzer contract, the shared result schema and
// the analyzer registry.
//
// An analyzer receives the indexed base and target snapshots and the change
// context, and returns one Result for one dimension. Analyzers are independent
// of each other, of the workflow that produced the patch, and of specific
// languages: they consume the shared model and declare which adapter
// capabilities they need. The engine computes coverage from those
// declarations, so an analyzer never silently scores code it cannot see.
package analyzer

import (
	"fmt"
	"sort"

	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

// Dimensions, in report order.
const (
	DimDuplication     = "duplication"
	DimComplexity      = "complexity"
	DimArchitecture    = "architecture_coupling"
	DimMaintainability = "maintainability_scope"
	DimCorrectness     = "correctness_confidence"
)

// DimensionOrder is the stable report order.
var DimensionOrder = []string{DimDuplication, DimComplexity, DimArchitecture, DimMaintainability, DimCorrectness}

// Statuses.
const (
	StatusOK          = "ok"
	StatusPartial     = "partial"     // some changed code could not be analysed
	StatusUnsupported = "unsupported" // no changed code could be analysed
	StatusError       = "error"       // the analyzer failed
	StatusNotAssessed = "not_assessed"
)

// Bases for findings.
const (
	BasisMeasured  = "measured"  // a deterministic measurement
	BasisHeuristic = "heuristic" // an interpretation that may be wrong
)

// Directions.
const (
	DirDegradation = "degradation"
	DirImprovement = "improvement"
	DirNeutral     = "neutral"
	DirInfo        = "info"
)

// Descriptor describes an analyzer.
type Descriptor struct {
	ID          string            `json:"id"`
	Dimension   string            `json:"dimension"`
	Title       string            `json:"title"`
	Version     string            `json:"version"`
	Requires    []lang.Capability `json:"requires"`
	Description string            `json:"description"`
	// Classes of changed files the analyzer measures; other classes are
	// reported as excluded in coverage.
	Classes []model.FileClass `json:"classes"`
}

// Analyzer is the analyzer contract.
type Analyzer interface {
	Descriptor() Descriptor
	Analyze(ctx *Context) (*Result, error)
}

// Context is the input to an analyzer.
type Context struct {
	Base, Target *index.Snapshot
	Change       *change.Set
	Config       *config.Config
	Langs        *lang.Registry
	Coverage     Coverage // computed by the engine for this analyzer
}

// Evidence points at source.
type Evidence struct {
	Side      string `json:"side"` // base | target
	Commit    string `json:"commit"`
	File      string `json:"file"`
	Symbol    string `json:"symbol,omitempty"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Role      string `json:"role,omitempty"`
}

// Measurement is a raw before/after value.
type Measurement struct {
	Metric string `json:"metric"`
	Scope  string `json:"scope,omitempty"`
	Before *int   `json:"before"`
	After  *int   `json:"after"`
	Delta  *int   `json:"delta"`
	Value  *int   `json:"value,omitempty"` // single-valued counts (before/after/delta are null)
	Unit   string `json:"unit,omitempty"`
	Note   string `json:"note,omitempty"`
}

// Finding is an individual observation with evidence.
type Finding struct {
	Kind         string        `json:"kind"`
	Basis        string        `json:"basis"`
	Direction    string        `json:"direction"`
	Message      string        `json:"message"`
	Similarity   *float64      `json:"similarity,omitempty"`
	Confidence   string        `json:"confidence,omitempty"` // high|medium|low, heuristics only
	Caveat       string        `json:"caveat,omitempty"`
	Measurements []Measurement `json:"measurements,omitempty"`
	Evidence     []Evidence    `json:"evidence"`
}

// Coverage says how much of the change an analyzer could see.
type Coverage struct {
	ChangedCodeFiles      int            `json:"changedCodeFiles"`
	Analyzed              int            `json:"analyzed"`
	UnsupportedByLanguage map[string]int `json:"unsupportedByLanguage,omitempty"`
	WithParseErrors       int            `json:"withParseErrors"`
	ExcludedByClass       map[string]int `json:"excludedByClass,omitempty"`
	NonCodeFiles          int            `json:"nonCodeFiles"`
}

// Unsupported totals unsupported changed code files.
func (c Coverage) Unsupported() int {
	n := 0
	for _, v := range c.UnsupportedByLanguage {
		n += v
	}
	return n
}

// Result is the common analyzer output schema.
type Result struct {
	Analyzer     string             `json:"analyzer"`
	Dimension    string             `json:"dimension"`
	Title        string             `json:"title"`
	Version      string             `json:"version"`
	Status       string             `json:"status"`
	Score        *int               `json:"score"` // -2..+2; null when not scored
	ScoreRule    string             `json:"scoreRule"`
	Summary      string             `json:"summary"`
	Rationale    []string           `json:"rationale,omitempty"`
	Measurements []Measurement      `json:"measurements"`
	Findings     []Finding          `json:"findings"`
	Coverage     Coverage           `json:"coverage"`
	Diagnostics  []model.Diagnostic `json:"diagnostics,omitempty"`
}

// IntP returns a pointer to v.
func IntP(v int) *int { return &v }

// Measure builds a before/after measurement.
func Measure(metric, scope string, before, after int, unit string) Measurement {
	return Measurement{Metric: metric, Scope: scope, Before: IntP(before), After: IntP(after), Delta: IntP(after - before), Unit: unit}
}

// Count builds a single-valued measurement of the change itself.
func Count(metric, scope string, value int, unit string) Measurement {
	return Measurement{Metric: metric, Scope: scope, Value: IntP(value), Unit: unit}
}

// Ev builds evidence for a symbol on a side.
func (c *Context) Ev(side string, s *model.Symbol, role string) Evidence {
	commit := c.Target.Commit
	if side == "base" {
		commit = c.Base.Commit
	}
	return Evidence{Side: side, Commit: commit, File: s.File, Symbol: s.DisplayName(), StartLine: s.StartLine, EndLine: s.EndLine, Role: role}
}

// EvLine builds evidence for a line.
func (c *Context) EvLine(side, file string, line int, role string) Evidence {
	commit := c.Target.Commit
	if side == "base" {
		commit = c.Base.Commit
	}
	return Evidence{Side: side, Commit: commit, File: file, StartLine: line, EndLine: line, Role: role}
}

// InClasses reports whether a class is in the analyzer's measured classes.
func (d Descriptor) InClasses(c model.FileClass) bool {
	for _, x := range d.Classes {
		if x == c {
			return true
		}
	}
	return false
}

// Registry holds analyzers.
type Registry struct {
	list []Analyzer
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry { return &Registry{} }

// Register adds an analyzer; IDs must be unique and dimensions known.
func (r *Registry) Register(a Analyzer) error {
	d := a.Descriptor()
	if d.ID == "" {
		return fmt.Errorf("analyzer has empty id")
	}
	known := false
	for _, dim := range DimensionOrder {
		if dim == d.Dimension {
			known = true
		}
	}
	if !known {
		return fmt.Errorf("analyzer %s: unknown dimension %q", d.ID, d.Dimension)
	}
	for _, e := range r.list {
		if e.Descriptor().ID == d.ID {
			return fmt.Errorf("analyzer %q already registered", d.ID)
		}
	}
	r.list = append(r.list, a)
	return nil
}

// All returns analyzers in dimension order, then by ID.
func (r *Registry) All() []Analyzer {
	out := append([]Analyzer{}, r.list...)
	rank := map[string]int{}
	for i, d := range DimensionOrder {
		rank[d] = i
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i].Descriptor(), out[j].Descriptor()
		if rank[a.Dimension] != rank[b.Dimension] {
			return rank[a.Dimension] < rank[b.Dimension]
		}
		return a.ID < b.ID
	})
	return out
}

// SortFindings orders findings deterministically.
func SortFindings(fs []Finding) {
	key := func(f Finding) (string, string, int, string) {
		if len(f.Evidence) == 0 {
			return f.Kind, "", 0, f.Message
		}
		return f.Kind, f.Evidence[0].File, f.Evidence[0].StartLine, f.Message
	}
	sort.SliceStable(fs, func(i, j int) bool {
		a1, a2, a3, a4 := key(fs[i])
		b1, b2, b3, b4 := key(fs[j])
		if a1 != b1 {
			return a1 < b1
		}
		if a2 != b2 {
			return a2 < b2
		}
		if a3 != b3 {
			return a3 < b3
		}
		return a4 < b4
	})
}

// Analyzable reports whether a snapshot file is measured by an analyzer: its
// class is in scope, it parsed, and its adapter has every required capability.
func (c *Context) Analyzable(d Descriptor, f *index.File) bool {
	if f == nil || !f.Supported() || !d.InClasses(f.Class) {
		return false
	}
	for _, cap := range d.Requires {
		if !lang.Has(f.Adapter, cap) {
			return false
		}
	}
	return true
}

// SymbolAnalyzable reports whether a symbol's file on the given side is
// analyzable.
func (c *Context) SymbolAnalyzable(d Descriptor, side string, s *model.Symbol) bool {
	snap := c.Target
	if side == "base" {
		snap = c.Base
	}
	return s != nil && c.Analyzable(d, snap.Files[s.File])
}

// ScoreFromThresholds maps a signed quantity where positive is worse onto the
// -2..+2 scale using symmetric thresholds.
func ScoreFromThresholds(v, at, significantAt int) int {
	switch {
	case v >= significantAt:
		return -2
	case v >= at:
		return -1
	case v <= -significantAt:
		return 2
	case v <= -at:
		return 1
	}
	return 0
}
