// Package scorev2 is the experimental second scoring model. It sits beside
// the frozen scoring/v1 analyzers and does not change them.
//
// v2 measures an absolute snapshot of a codebase (and of each module in it)
// as debt points from four transparent components, then derives deltas as
// after − before. Because every component is a function of a single
// snapshot, deltas telescope: the per-step deltas along a first-parent chain
// sum exactly to the net delta. Debt is additive across modules: module debts
// sum to the codebase debt.
//
//   - complexity debt: cyclomatic points at or above the hotspot threshold, Σ max(0, cc − H + 1)
//   - nesting debt:    nesting levels at or above the nesting threshold, Σ max(0, depth − N + 1)
//
// so a callable carries complexity or nesting debt exactly when it is a
// hotspot (cc ≥ H or depth ≥ N), and growth below the thresholds carries none
// (it is reported separately as decision points).
//   - duplication debt: redundant tokens in structural clone groups / tokens-per-point
//   - architecture debt: configured rule-violating edges and units in cycles
//
// The 0–100 score is a size-normalised view of total debt for comparing
// modules or codebases of different sizes; debt itself is the quantity to
// optimise within one codebase, because adding code cannot dilute it.
// Weights are uncalibrated defaults and are recorded in every report.
package scorev2

import (
	"math"
	"sort"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

// Versions and schemas.
const (
	Version        = "scoring/v2-experimental"
	SnapshotSchema = "codeeval.snapshot/v2"
	DeltaSchema    = "codeeval.delta/v2"
)

// Weights convert raw excess into debt points. They are uncalibrated defaults.
type Weights struct {
	CyclomaticExcess float64 `json:"cyclomaticExcessPoint"`   // debt per cyclomatic point at or above the hotspot threshold
	NestingExcess    float64 `json:"nestingExcessLevel"`      // debt per nesting level at or above the nesting threshold
	TokensPerPoint   float64 `json:"redundantTokensPerPoint"` // redundant clone tokens per debt point
	RuleViolation    float64 `json:"ruleViolationEdge"`       // debt per configured-rule-violating edge
	CycleUnit        float64 `json:"cycleUnit"`               // debt per unit inside a dependency cycle
	HalfScoreDensity float64 `json:"halfScoreDensity"`        // debt per 1,000 source lines at which the score is 50
}

// DefaultWeights returns the v2 defaults.
func DefaultWeights() Weights {
	return Weights{CyclomaticExcess: 1, NestingExcess: 2, TokensPerPoint: 20, RuleViolation: 10, CycleUnit: 20, HalfScoreDensity: 25}
}

// Thresholds copied from the configuration.
type Thresholds struct {
	HotspotCyclomatic int `json:"hotspotCyclomatic"`
	HotspotNesting    int `json:"hotspotNesting"`
	MinTokens         int `json:"minCloneTokens"`
}

// Components are raw, integer-or-token measurements of a unit.
type Components struct {
	SourceFiles      int     `json:"sourceFiles"`
	SourceLines      int     `json:"sourceLines"`
	Callables        int     `json:"callables"`
	CyclomaticTotal  int     `json:"cyclomaticTotal"`
	DecisionPoints   int     `json:"decisionPoints"`
	Hotspots         int     `json:"hotspots"`
	CyclomaticExcess int     `json:"cyclomaticExcess"`
	NestingExcess    int     `json:"nestingExcess"`
	ClonedCallables  int     `json:"clonedCallables"`               // members of structural clone groups (size ≥ 2)
	ExactRedundant   int     `json:"exactRedundant,omitempty"`      // Σ(group size − 1), exact groups (codebase only)
	StructRedundant  int     `json:"structuralRedundant,omitempty"` // Σ(group size − 1), structural groups (codebase only)
	RedundantTokens  float64 `json:"redundantTokens"`
	RuleViolations   int     `json:"ruleViolations"`
	CycleUnits       int     `json:"cycleUnits"`
	InternalEdges    int     `json:"internalEdges"` // codebase: all edges; module: fan-out
}

// Sub returns c − o.
func (c Components) Sub(o Components) Components {
	return Components{
		SourceFiles: c.SourceFiles - o.SourceFiles, SourceLines: c.SourceLines - o.SourceLines,
		Callables: c.Callables - o.Callables, CyclomaticTotal: c.CyclomaticTotal - o.CyclomaticTotal,
		DecisionPoints: c.DecisionPoints - o.DecisionPoints, Hotspots: c.Hotspots - o.Hotspots,
		CyclomaticExcess: c.CyclomaticExcess - o.CyclomaticExcess, NestingExcess: c.NestingExcess - o.NestingExcess,
		ClonedCallables: c.ClonedCallables - o.ClonedCallables, ExactRedundant: c.ExactRedundant - o.ExactRedundant,
		StructRedundant: c.StructRedundant - o.StructRedundant, RedundantTokens: c.RedundantTokens - o.RedundantTokens,
		RuleViolations: c.RuleViolations - o.RuleViolations, CycleUnits: c.CycleUnits - o.CycleUnits,
		InternalEdges: c.InternalEdges - o.InternalEdges,
	}
}

func (c *Components) add(o Components) {
	c.SourceFiles += o.SourceFiles
	c.SourceLines += o.SourceLines
	c.Callables += o.Callables
	c.CyclomaticTotal += o.CyclomaticTotal
	c.DecisionPoints += o.DecisionPoints
	c.Hotspots += o.Hotspots
	c.CyclomaticExcess += o.CyclomaticExcess
	c.NestingExcess += o.NestingExcess
	c.ClonedCallables += o.ClonedCallables
	c.RedundantTokens += o.RedundantTokens
	c.RuleViolations += o.RuleViolations
	c.CycleUnits += o.CycleUnits
}

// Debt is debt points by component.
type Debt struct {
	Complexity   float64 `json:"complexity"`
	Nesting      float64 `json:"nesting"`
	Duplication  float64 `json:"duplication"`
	Architecture float64 `json:"architecture"`
	Total        float64 `json:"total"`
}

// Sub returns d − o.
func (d Debt) Sub(o Debt) Debt {
	return Debt{d.Complexity - o.Complexity, d.Nesting - o.Nesting, d.Duplication - o.Duplication, d.Architecture - o.Architecture, d.Total - o.Total}
}

// Scale returns d × f.
func (d Debt) Scale(f float64) Debt {
	return Debt{d.Complexity * f, d.Nesting * f, d.Duplication * f, d.Architecture * f, d.Total * f}
}

func (d Debt) rounded() Debt {
	return Debt{r2(d.Complexity), r2(d.Nesting), r2(d.Duplication), r2(d.Architecture), r2(d.Total)}
}

// DebtOf converts components into debt.
func DebtOf(c Components, w Weights) Debt {
	d := Debt{
		Complexity:   float64(c.CyclomaticExcess) * w.CyclomaticExcess,
		Nesting:      float64(c.NestingExcess) * w.NestingExcess,
		Duplication:  c.RedundantTokens / w.TokensPerPoint,
		Architecture: float64(c.RuleViolations)*w.RuleViolation + float64(c.CycleUnits)*w.CycleUnit,
	}
	d.Total = d.Complexity + d.Nesting + d.Duplication + d.Architecture
	return d
}

// Unit is a scored module or the whole codebase.
type Unit struct {
	ID          string     `json:"id"`
	Language    string     `json:"language,omitempty"`
	Components  Components `json:"components"`
	Debt        Debt       `json:"debt"`
	DebtPerKLOC *float64   `json:"debtPerKLOC"`
	Score       *float64   `json:"score"` // 0–100, null when there is no supported source
}

func finishUnit(u *Unit, w Weights) {
	u.Debt = DebtOf(u.Components, w)
	if u.Components.SourceLines > 0 {
		dens := u.Debt.Total / (float64(u.Components.SourceLines) / 1000)
		s := 100 * w.HalfScoreDensity / (w.HalfScoreDensity + dens)
		u.DebtPerKLOC, u.Score = fp(r2(dens)), fp(r2(s))
	}
}

// CallableDebt is one callable's complexity and nesting debt.
type CallableDebt struct {
	Module     string  `json:"module"`
	File       string  `json:"file"`
	Name       string  `json:"name"`
	StartLine  int     `json:"startLine"`
	EndLine    int     `json:"endLine"`
	Cyclomatic int     `json:"cyclomatic"`
	Nesting    int     `json:"nesting"`
	Lines      int     `json:"lines"`
	Debt       float64 `json:"debt"`
}

// CloneMember is one member of a clone group.
type CloneMember struct {
	Module    string `json:"module"`
	File      string `json:"file"`
	Name      string `json:"name"`
	StartLine int    `json:"startLine"`
	Tokens    int    `json:"tokens"`
}

// CloneGroup is a structural clone group with its debt.
type CloneGroup struct {
	Exact           bool          `json:"exact"`
	RedundantTokens float64       `json:"redundantTokens"`
	Debt            float64       `json:"debt"`
	Members         []CloneMember `json:"members"`
}

// Violation is a rule-violating edge.
type Violation struct {
	Language string `json:"language"`
	Rule     string `json:"rule"`
	From     string `json:"from"`
	To       string `json:"to"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// LanguageCoverage counts production source per language.
type LanguageCoverage struct {
	Language  string `json:"language"`
	Supported bool   `json:"supported"`
	Files     int    `json:"files"`
	Bytes     int64  `json:"bytes"`
}

// Coverage says how much production source v2 could measure.
type Coverage struct {
	Languages            []LanguageCoverage `json:"languages"`
	SupportedSourceFiles int                `json:"supportedSourceFiles"`
	SourceFiles          int                `json:"sourceFiles"`
	SupportedByteShare   float64            `json:"supportedByteShare"` // 0–1
	ParseErrorFiles      int                `json:"parseErrorFiles"`
}

// Analysis is the full v2 measurement of one snapshot.
type Analysis struct {
	Weights    Weights
	Thresholds Thresholds
	Coverage   Coverage
	Codebase   Unit
	Modules    []Unit
	Callables  []CallableDebt // all callables with debt > 0, worst first
	Clones     []CloneGroup   // all structural groups, worst first
	Violations []Violation
	Cycles     [][]string
}

// Analyze measures a snapshot.
func Analyze(snap *index.Snapshot, cfg *config.Config, w Weights) *Analysis {
	a := &Analysis{Weights: w, Thresholds: Thresholds{cfg.Complexity.HotspotCyclomatic, cfg.Complexity.HotspotNesting, cfg.Duplication.MinTokens}}
	mods := map[string]*Unit{}
	mod := func(f *index.File) *Unit {
		id := f.Adapter.PackageOf(f.Path)
		key := f.Language + "\x00" + id
		u := mods[key]
		if u == nil {
			u = &Unit{ID: id, Language: f.Language}
			mods[key] = u
		}
		return u
	}

	// Coverage and size.
	byLang := map[string]*LanguageCoverage{}
	var supBytes, allBytes int64
	for _, p := range snap.Paths {
		f := snap.Files[p]
		if f.Class != model.ClassSource || f.Language == "" {
			continue
		}
		lc := byLang[f.Language]
		if lc == nil {
			lc = &LanguageCoverage{Language: f.Language}
			byLang[f.Language] = lc
		}
		lc.Files++
		lc.Bytes += f.Size
		a.Coverage.SourceFiles++
		allBytes += f.Size
		if !f.Supported() {
			continue
		}
		lc.Supported = true
		a.Coverage.SupportedSourceFiles++
		supBytes += f.Size
		if len(f.Parsed.ParseErrors) > 0 {
			a.Coverage.ParseErrorFiles++
		}
		u := mod(f)
		u.Components.SourceFiles++
		u.Components.SourceLines += f.Parsed.Lines
	}
	for _, lc := range byLang {
		a.Coverage.Languages = append(a.Coverage.Languages, *lc)
	}
	sort.Slice(a.Coverage.Languages, func(i, j int) bool { return a.Coverage.Languages[i].Language < a.Coverage.Languages[j].Language })
	if allBytes > 0 {
		a.Coverage.SupportedByteShare = r4(float64(supBytes) / float64(allBytes))
	}

	// Complexity and nesting.
	H, N := cfg.Complexity.HotspotCyclomatic, cfg.Complexity.HotspotNesting
	for _, s := range snap.Symbols {
		f := snap.Files[s.File]
		if !s.Callable() || f == nil || f.Class != model.ClassSource || !f.Supported() || !lang.Has(f.Adapter, lang.CapComplexity) {
			continue
		}
		m := s.Metrics
		u := mod(f)
		u.Components.Callables++
		u.Components.CyclomaticTotal += m.Cyclomatic
		u.Components.DecisionPoints += m.Decisions
		if m.Cyclomatic >= H || m.MaxNesting >= N {
			u.Components.Hotspots++
		}
		ce, ne := max0(m.Cyclomatic-H+1), max0(m.MaxNesting-N+1)
		u.Components.CyclomaticExcess += ce
		u.Components.NestingExcess += ne
		if d := float64(ce)*w.CyclomaticExcess + float64(ne)*w.NestingExcess; d > 0 {
			a.Callables = append(a.Callables, CallableDebt{Module: u.ID, File: s.File, Name: s.DisplayName(), StartLine: s.StartLine, EndLine: s.EndLine,
				Cyclomatic: m.Cyclomatic, Nesting: m.MaxNesting, Lines: m.Lines, Debt: d})
		}
	}
	sort.SliceStable(a.Callables, func(i, j int) bool {
		if a.Callables[i].Debt != a.Callables[j].Debt {
			return a.Callables[i].Debt > a.Callables[j].Debt
		}
		if a.Callables[i].File != a.Callables[j].File {
			return a.Callables[i].File < a.Callables[j].File
		}
		return a.Callables[i].StartLine < a.Callables[j].StartLine
	})

	// Duplication: the same population as the v1 duplication analyzer.
	var pop []*model.Symbol
	for _, s := range snap.Symbols {
		f := snap.Files[s.File]
		if !s.Callable() || s.Tokens < cfg.Duplication.MinTokens || s.InErrorRegion || f == nil || !f.Supported() {
			continue
		}
		if !(f.Class == model.ClassSource || (cfg.Duplication.IncludeTests && f.Class == model.ClassTest)) || !lang.Has(f.Adapter, lang.CapFingerprints) {
			continue
		}
		pop = append(pop, s)
	}
	exact, structural := map[string][]*model.Symbol{}, map[string][]*model.Symbol{}
	for _, s := range pop {
		if s.ExactFP != "" {
			exact[s.Kind+":"+s.ExactFP] = append(exact[s.Kind+":"+s.ExactFP], s)
		}
		if s.StructFP != "" {
			structural[s.Kind+":"+s.StructFP] = append(structural[s.Kind+":"+s.StructFP], s)
		}
	}
	for _, g := range exact {
		if len(g) > 1 {
			a.Codebase.Components.ExactRedundant += len(g) - 1
		}
	}
	keys := make([]string, 0, len(structural))
	for k, g := range structural {
		if len(g) > 1 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		g := structural[k]
		sort.Slice(g, func(i, j int) bool { return g[i].Key < g[j].Key })
		n := float64(len(g))
		a.Codebase.Components.StructRedundant += len(g) - 1
		cg := CloneGroup{Exact: true}
		for _, s := range g {
			f := snap.Files[s.File]
			u := mod(f)
			share := float64(s.Tokens) * (n - 1) / n
			u.Components.ClonedCallables++
			u.Components.RedundantTokens += share
			cg.RedundantTokens += share
			cg.Members = append(cg.Members, CloneMember{Module: u.ID, File: s.File, Name: s.DisplayName(), StartLine: s.StartLine, Tokens: s.Tokens})
			if s.ExactFP != g[0].ExactFP {
				cg.Exact = false
			}
		}
		cg.Debt = r2(cg.RedundantTokens / w.TokensPerPoint)
		cg.RedundantTokens = r2(cg.RedundantTokens)
		a.Clones = append(a.Clones, cg)
	}
	sort.SliceStable(a.Clones, func(i, j int) bool { return a.Clones[i].RedundantTokens > a.Clones[j].RedundantTokens })

	// Architecture: rule-violating edges (attributed to the importing unit)
	// and units inside cycles.
	gids := make([]string, 0, len(snap.Graphs))
	for id := range snap.Graphs {
		gids = append(gids, id)
	}
	sort.Strings(gids)
	for _, gid := range gids {
		g := snap.Graphs[gid]
		unitOf := func(id string) *Unit {
			key := gid + "\x00" + id
			u := mods[key]
			if u == nil {
				u = &Unit{ID: id, Language: gid}
				mods[key] = u
			}
			return u
		}
		var froms []string
		for from := range g.Edges {
			froms = append(froms, from)
		}
		sort.Strings(froms)
		for _, from := range froms {
			var tos []string
			for to := range g.Edges[from] {
				tos = append(tos, to)
			}
			sort.Strings(tos)
			unitOf(from).Components.InternalEdges += len(tos)
			a.Codebase.Components.InternalEdges += len(tos)
			for _, to := range tos {
				for _, r := range cfg.Architecture.Rules {
					if config.Match(r.From, from) && config.Match(r.To, to) {
						unitOf(from).Components.RuleViolations++
						loc := g.Edges[from][to][0]
						a.Violations = append(a.Violations, Violation{Language: gid, Rule: r.Name, From: from, To: to, File: loc.File, Line: loc.Line})
					}
				}
			}
		}
		for _, comp := range sccs(g) {
			a.Cycles = append(a.Cycles, comp)
			for _, id := range comp {
				unitOf(id).Components.CycleUnits++
			}
		}
	}

	// Modules and totals (codebase-only fields were set above).
	mkeys := make([]string, 0, len(mods))
	for k := range mods {
		mkeys = append(mkeys, k)
	}
	sort.Strings(mkeys)
	a.Codebase.ID = "codebase"
	for _, k := range mkeys {
		u := mods[k]
		finishUnit(u, w)
		a.Codebase.Components.add(u.Components)
		a.Modules = append(a.Modules, *u)
	}
	finishUnit(&a.Codebase, w)
	sort.SliceStable(a.Modules, func(i, j int) bool {
		if a.Modules[i].Debt.Total != a.Modules[j].Debt.Total {
			return a.Modules[i].Debt.Total > a.Modules[j].Debt.Total
		}
		return a.Modules[i].ID < a.Modules[j].ID
	})
	return a
}

func sccs(g *lang.Graph) [][]string {
	var nodes []string
	for n := range g.Nodes {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)
	adj := map[string][]string{}
	for f, m := range g.Edges {
		for t := range m {
			adj[f] = append(adj[f], t)
		}
		sort.Strings(adj[f])
	}
	idx, low, on := map[string]int{}, map[string]int{}, map[string]bool{}
	var stack []string
	var out [][]string
	i := 0
	var strong func(v string)
	strong = func(v string) {
		idx[v], low[v] = i, i
		i++
		stack = append(stack, v)
		on[v] = true
		for _, w := range adj[v] {
			if _, ok := idx[w]; !ok {
				strong(w)
				low[v] = min(low[v], low[w])
			} else if on[w] {
				low[v] = min(low[v], idx[w])
			}
		}
		if low[v] == idx[v] {
			var comp []string
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				on[w] = false
				comp = append(comp, w)
				if w == v {
					break
				}
			}
			if len(comp) > 1 {
				sort.Strings(comp)
				out = append(out, comp)
			}
		}
	}
	for _, n := range nodes {
		if _, ok := idx[n]; !ok {
			strong(n)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i][0] < out[j][0] })
	return out
}

func max0(v int) int {
	if v < 0 {
		return 0
	}
	return v
}

func r2(v float64) float64  { return math.Round(v*100) / 100 }
func r4(v float64) float64  { return math.Round(v*10000) / 10000 }
func fp(v float64) *float64 { return &v }
