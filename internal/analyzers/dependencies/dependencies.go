// Package dependencies compares the internal dependency graphs that language
// adapters build for base and target: edges added/removed, fan-in/fan-out of
// touched units, cycles introduced/removed, and configured architecture rule
// violations.
package dependencies

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

// Analyzer is the dependency/architecture analyzer.
type Analyzer struct{}

// New returns the analyzer.
func New() *Analyzer { return &Analyzer{} }

func (Analyzer) Descriptor() analyzer.Descriptor {
	return analyzer.Descriptor{
		ID: "dependencies", Dimension: analyzer.DimArchitecture, Title: "Dependency graph delta", Version: "1",
		Requires:    []lang.Capability{lang.CapDependencies},
		Description: "Internal dependency edges, fan-in/fan-out, strongly connected cycles and configured forbidden-dependency rules, before and after.",
		Classes:     []model.FileClass{model.ClassSource, model.ClassGenerated},
	}
}

type edge struct{ from, to string }

func edges(g *lang.Graph) map[edge][]model.Location {
	out := map[edge][]model.Location{}
	if g == nil {
		return out
	}
	for f, m := range g.Edges {
		for t, locs := range m {
			out[edge{f, t}] = locs
		}
	}
	return out
}

func sortedEdges(m map[edge]bool) []edge {
	var out []edge
	for e := range m {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].from != out[j].from {
			return out[i].from < out[j].from
		}
		return out[i].to < out[j].to
	})
	return out
}

func (an Analyzer) Analyze(ctx *analyzer.Context) (*analyzer.Result, error) {
	res := &analyzer.Result{}
	langs := map[string]bool{}
	for id := range ctx.Base.Graphs {
		langs[id] = true
	}
	for id := range ctx.Target.Graphs {
		langs[id] = true
	}
	ids := make([]string, 0, len(langs))
	for id := range langs {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	// Units touched by the change (graph node ids are adapter-defined).
	touched := map[string]map[string]bool{}
	for _, fc := range ctx.Change.Files {
		ad := ctx.Langs.Adapter(fc.Language)
		if ad == nil || !lang.Has(ad, lang.CapDependencies) {
			continue
		}
		if touched[fc.Language] == nil {
			touched[fc.Language] = map[string]bool{}
		}
		touched[fc.Language][ad.PackageOf(fc.Path)] = true
		if fc.OldPath != "" {
			touched[fc.Language][ad.PackageOf(fc.OldPath)] = true
		}
	}

	var newEdges, goneEdges, cyclesIn, cyclesOut, newViol, goneViol int
	var totB, totA, nodesB, nodesA, extB, extA, violB, violA, cycB, cycA int
	for _, id := range ids {
		bg, tg := ctx.Base.Graphs[id], ctx.Target.Graphs[id]
		be, te := edges(bg), edges(tg)
		totB += len(be)
		totA += len(te)
		if bg != nil {
			nodesB += len(bg.Nodes)
			extB += externalCount(bg)
		}
		if tg != nil {
			nodesA += len(tg.Nodes)
			extA += externalCount(tg)
			res.Rationale = append(res.Rationale, id+": "+tg.Note)
		}
		added, gone := map[edge]bool{}, map[edge]bool{}
		for e := range te {
			if _, ok := be[e]; !ok {
				added[e] = true
			}
		}
		for e := range be {
			if _, ok := te[e]; !ok {
				gone[e] = true
			}
		}
		newEdges += len(added)
		goneEdges += len(gone)
		for _, e := range sortedEdges(added) {
			loc := te[e][0]
			res.Findings = append(res.Findings, analyzer.Finding{
				Kind: "dependencies.edge_added", Basis: analyzer.BasisMeasured, Direction: analyzer.DirInfo,
				Message:  fmt.Sprintf("[%s] new dependency %s → %s", id, e.from, e.to),
				Evidence: []analyzer.Evidence{ctx.EvLine("target", loc.File, loc.Line, "import")},
			})
		}
		for _, e := range sortedEdges(gone) {
			loc := be[e][0]
			res.Findings = append(res.Findings, analyzer.Finding{
				Kind: "dependencies.edge_removed", Basis: analyzer.BasisMeasured, Direction: analyzer.DirInfo,
				Message:  fmt.Sprintf("[%s] dependency removed %s → %s", id, e.from, e.to),
				Evidence: []analyzer.Evidence{ctx.EvLine("base", loc.File, loc.Line, "removed import")},
			})
		}

		// Fan-in/fan-out of touched units.
		var units []string
		for u := range touched[id] {
			units = append(units, u)
		}
		sort.Strings(units)
		for _, u := range units {
			fob, foa := fanOut(be, u), fanOut(te, u)
			fib, fia := fanIn(be, u), fanIn(te, u)
			if fob != foa || fib != fia {
				res.Findings = append(res.Findings, analyzer.Finding{
					Kind: "dependencies.fan_change", Basis: analyzer.BasisMeasured, Direction: analyzer.DirInfo,
					Message: fmt.Sprintf("[%s] %s fan-out %d→%d, fan-in %d→%d", id, u, fob, foa, fib, fia),
					Measurements: []analyzer.Measurement{
						analyzer.Measure("fan_out", u, fob, foa, "units"),
						analyzer.Measure("fan_in", u, fib, fia, "units"),
					},
					Evidence: []analyzer.Evidence{},
				})
			}
		}

		// Cycles.
		bc, tc := cycles(be), cycles(te)
		cycB += len(bc)
		cycA += len(tc)
		for k, members := range tc {
			if _, ok := bc[k]; !ok {
				cyclesIn++
				res.Findings = append(res.Findings, analyzer.Finding{
					Kind: "dependencies.cycle_introduced", Basis: analyzer.BasisMeasured, Direction: analyzer.DirDegradation,
					Message:  fmt.Sprintf("[%s] dependency cycle introduced among %s", id, strings.Join(members, ", ")),
					Evidence: cycleEvidence(ctx, "target", te, members),
				})
			}
		}
		for k, members := range bc {
			if _, ok := tc[k]; !ok {
				cyclesOut++
				res.Findings = append(res.Findings, analyzer.Finding{
					Kind: "dependencies.cycle_removed", Basis: analyzer.BasisMeasured, Direction: analyzer.DirImprovement,
					Message:  fmt.Sprintf("[%s] dependency cycle removed among %s", id, strings.Join(members, ", ")),
					Evidence: cycleEvidence(ctx, "base", be, members),
				})
			}
		}

		// Architecture rules.
		vb, va := violations(ctx.Config, be), violations(ctx.Config, te)
		violB += len(vb)
		violA += len(va)
		for _, v := range sortedViolations(va) {
			if _, ok := vb[v]; ok {
				loc := te[v.e][0]
				res.Findings = append(res.Findings, analyzer.Finding{
					Kind: "dependencies.rule_violation_preexisting", Basis: analyzer.BasisMeasured, Direction: analyzer.DirInfo,
					Message:  fmt.Sprintf("[%s] rule %q already violated before the change: %s → %s (not scored)", id, v.rule, v.e.from, v.e.to),
					Evidence: []analyzer.Evidence{ctx.EvLine("target", loc.File, loc.Line, "import")},
				})
				continue
			}
			newViol++
			loc := te[v.e][0]
			res.Findings = append(res.Findings, analyzer.Finding{
				Kind: "dependencies.rule_violation_introduced", Basis: analyzer.BasisMeasured, Direction: analyzer.DirDegradation,
				Message:  fmt.Sprintf("[%s] rule %q violated: %s → %s", id, v.rule, v.e.from, v.e.to),
				Evidence: []analyzer.Evidence{ctx.EvLine("target", loc.File, loc.Line, "import")},
			})
		}
		for _, v := range sortedViolations(vb) {
			if _, ok := va[v]; ok {
				continue
			}
			goneViol++
			loc := be[v.e][0]
			res.Findings = append(res.Findings, analyzer.Finding{
				Kind: "dependencies.rule_violation_resolved", Basis: analyzer.BasisMeasured, Direction: analyzer.DirImprovement,
				Message:  fmt.Sprintf("[%s] rule %q no longer violated: %s → %s", id, v.rule, v.e.from, v.e.to),
				Evidence: []analyzer.Evidence{ctx.EvLine("base", loc.File, loc.Line, "removed import")},
			})
		}
	}
	res.Measurements = append(res.Measurements,
		analyzer.Measure("internal_units", "repository", nodesB, nodesA, "units"),
		analyzer.Measure("internal_edges", "repository", totB, totA, "edges"),
		analyzer.Measure("external_dependency_refs", "repository", extB, extA, "unit→external pairs"),
		analyzer.Measure("cycles", "repository", cycB, cycA, "strongly connected components"),
		analyzer.Measure("rule_violations", "repository", violB, violA, "edges"),
		analyzer.Count("edges_added", "change", newEdges, "edges"),
		analyzer.Count("edges_removed", "change", goneEdges, "edges"),
	)
	if len(ctx.Config.Architecture.Rules) == 0 {
		res.Rationale = append(res.Rationale, "no architecture rules configured: boundary violations not assessed")
	}
	score := 0
	degrade := cyclesIn > 0 || newViol > 0
	improve := cyclesOut > 0 || goneViol > 0 || (goneEdges > 0 && newEdges == 0)
	switch {
	case degrade && improve:
		res.ScoreRule = "mixed: degradation and improvement signals both present → 0"
	case cyclesIn > 0 || newViol >= 2:
		score = -2
		res.ScoreRule = "cycle introduced or ≥2 new rule violations → -2"
	case newViol == 1:
		score = -1
		res.ScoreRule = "one new rule violation → -1"
	case cyclesOut > 0:
		score = 2
		res.ScoreRule = "cycle removed, none introduced → +2"
	case improve:
		score = 1
		res.ScoreRule = "rule violation resolved, or edges removed with none added → +1"
	default:
		res.ScoreRule = "no cycles or rule violations changed → 0 (new edges alone are reported, not judged: features legitimately add dependencies)"
	}
	res.Score = &score
	res.Summary = fmt.Sprintf("internal edges %d→%d (+%d/-%d); cycles %d→%d; rule violations %d→%d", totB, totA, newEdges, goneEdges, cycB, cycA, violB, violA)
	return res, nil
}

func externalCount(g *lang.Graph) int {
	n := 0
	for _, m := range g.External {
		n += len(m)
	}
	return n
}

func fanOut(es map[edge][]model.Location, u string) int {
	n := 0
	for e := range es {
		if e.from == u {
			n++
		}
	}
	return n
}

func fanIn(es map[edge][]model.Location, u string) int {
	n := 0
	for e := range es {
		if e.to == u {
			n++
		}
	}
	return n
}

// cycles returns strongly connected components of size > 1 keyed by their
// sorted member list.
func cycles(es map[edge][]model.Location) map[string][]string {
	adj := map[string][]string{}
	nodes := map[string]bool{}
	for e := range es {
		adj[e.from] = append(adj[e.from], e.to)
		nodes[e.from], nodes[e.to] = true, true
	}
	var order []string
	for n := range nodes {
		order = append(order, n)
	}
	sort.Strings(order)
	for _, n := range order {
		sort.Strings(adj[n])
	}
	idx, low := map[string]int{}, map[string]int{}
	on := map[string]bool{}
	var stack []string
	i := 0
	out := map[string][]string{}
	var strong func(v string)
	strong = func(v string) {
		idx[v], low[v] = i, i
		i++
		stack = append(stack, v)
		on[v] = true
		for _, w := range adj[v] {
			if _, ok := idx[w]; !ok {
				strong(w)
				if low[w] < low[v] {
					low[v] = low[w]
				}
			} else if on[w] && idx[w] < low[v] {
				low[v] = idx[w]
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
				out[strings.Join(comp, "\x00")] = comp
			}
		}
	}
	for _, n := range order {
		if _, ok := idx[n]; !ok {
			strong(n)
		}
	}
	return out
}

func cycleEvidence(ctx *analyzer.Context, side string, es map[edge][]model.Location, members []string) []analyzer.Evidence {
	in := map[string]bool{}
	for _, m := range members {
		in[m] = true
	}
	var ev []analyzer.Evidence
	m := map[edge]bool{}
	for e := range es {
		if in[e.from] && in[e.to] {
			m[e] = true
		}
	}
	for _, e := range sortedEdges(m) {
		loc := es[e][0]
		ev = append(ev, ctx.EvLine(side, loc.File, loc.Line, "cycle edge "+e.from+" → "+e.to))
	}
	return ev
}

type violation struct {
	rule string
	e    edge
}

func violations(c *config.Config, es map[edge][]model.Location) map[violation]bool {
	out := map[violation]bool{}
	for _, r := range c.Architecture.Rules {
		for e := range es {
			if config.Match(r.From, e.from) && config.Match(r.To, e.to) {
				out[violation{r.Name, e}] = true
			}
		}
	}
	return out
}

func sortedViolations(m map[violation]bool) []violation {
	var out []violation
	for v := range m {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].rule != out[j].rule {
			return out[i].rule < out[j].rule
		}
		if out[i].e.from != out[j].e.from {
			return out[i].e.from < out[j].e.from
		}
		return out[i].e.to < out[j].e.to
	})
	return out
}
