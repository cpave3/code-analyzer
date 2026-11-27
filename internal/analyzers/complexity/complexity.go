// Package complexity measures per-callable cyclomatic complexity, decisions,
// nesting and length before and after the change, and separates complexity
// moved into new helpers from complexity actually removed.
package complexity

import (
	"fmt"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

// Analyzer is the complexity analyzer.
type Analyzer struct{}

// New returns the analyzer.
func New() *Analyzer { return &Analyzer{} }

func (Analyzer) Descriptor() analyzer.Descriptor {
	return analyzer.Descriptor{
		ID: "complexity", Dimension: analyzer.DimComplexity, Title: "Callable complexity deltas", Version: "1",
		Requires:    []lang.Capability{lang.CapSymbols, lang.CapComplexity},
		Description: "Cyclomatic complexity, decisions, max nesting and length for modified, added and deleted callables in production source.",
		Classes:     []model.FileClass{model.ClassSource},
	}
}

type agg struct{ cc, dec, lines, n, maxCC, maxNest, hot int }

func (a *agg) add(m *model.Metrics, hot bool) {
	a.cc += m.Cyclomatic
	a.dec += m.Decisions
	a.lines += m.Lines
	a.n++
	if m.Cyclomatic > a.maxCC {
		a.maxCC = m.Cyclomatic
	}
	if m.MaxNesting > a.maxNest {
		a.maxNest = m.MaxNesting
	}
	if hot {
		a.hot++
	}
}

func (an Analyzer) Analyze(ctx *analyzer.Context) (*analyzer.Result, error) {
	d := an.Descriptor()
	cfg := ctx.Config.Complexity
	hot := func(m *model.Metrics) bool {
		return m.Cyclomatic >= cfg.HotspotCyclomatic || m.MaxNesting >= cfg.HotspotNesting
	}
	res := &analyzer.Result{}
	var modBefore, modAfter, added, deleted agg
	var scopeBefore, scopeAfter agg // changed scope: modified+deleted (before) vs modified+added (after)
	var errSyms int
	increased, decreased := 0, 0
	for _, sc := range ctx.Change.Symbols {
		p := sc.Primary()
		if !p.Callable() {
			continue
		}
		if (sc.Target != nil && sc.Target.InErrorRegion) || (sc.Base != nil && sc.Base.InErrorRegion) {
			errSyms++
		}
		switch sc.Status {
		case change.SymModified:
			if !ctx.SymbolAnalyzable(d, "base", sc.Base) || !ctx.SymbolAnalyzable(d, "target", sc.Target) {
				continue
			}
			b, t := sc.Base.Metrics, sc.Target.Metrics
			modBefore.add(b, hot(b))
			modAfter.add(t, hot(t))
			scopeBefore.add(b, hot(b))
			scopeAfter.add(t, hot(t))
			if *b == *t {
				continue
			}
			dir := analyzer.DirNeutral
			switch {
			case t.Cyclomatic > b.Cyclomatic || (t.Cyclomatic == b.Cyclomatic && t.MaxNesting > b.MaxNesting):
				dir = analyzer.DirDegradation
				increased++
			case t.Cyclomatic < b.Cyclomatic || (t.Cyclomatic == b.Cyclomatic && t.MaxNesting < b.MaxNesting):
				dir = analyzer.DirImprovement
				decreased++
			}
			res.Findings = append(res.Findings, analyzer.Finding{
				Kind: "complexity.modified", Basis: analyzer.BasisMeasured, Direction: dir,
				Message:      fmt.Sprintf("%s: cyclomatic %d→%d, nesting %d→%d, lines %d→%d", sc.Target.DisplayName(), b.Cyclomatic, t.Cyclomatic, b.MaxNesting, t.MaxNesting, b.Lines, t.Lines),
				Measurements: metricRows(b, t),
				Evidence:     []analyzer.Evidence{ctx.Ev("base", sc.Base, "before"), ctx.Ev("target", sc.Target, "after")},
			})
		case change.SymAdded:
			if !ctx.SymbolAnalyzable(d, "target", sc.Target) {
				continue
			}
			t := sc.Target.Metrics
			added.add(t, hot(t))
			scopeAfter.add(t, hot(t))
			if hot(t) {
				res.Findings = append(res.Findings, analyzer.Finding{
					Kind: "complexity.added_hotspot", Basis: analyzer.BasisMeasured, Direction: analyzer.DirDegradation,
					Message:  fmt.Sprintf("new callable %s has cyclomatic %d, nesting %d (hotspot thresholds %d/%d)", sc.Target.DisplayName(), t.Cyclomatic, t.MaxNesting, cfg.HotspotCyclomatic, cfg.HotspotNesting),
					Evidence: []analyzer.Evidence{ctx.Ev("target", sc.Target, "added")},
				})
			}
		case change.SymDeleted:
			if !ctx.SymbolAnalyzable(d, "base", sc.Base) {
				continue
			}
			b := sc.Base.Metrics
			deleted.add(b, hot(b))
			scopeBefore.add(b, hot(b))
			if hot(b) {
				res.Findings = append(res.Findings, analyzer.Finding{
					Kind: "complexity.deleted_hotspot", Basis: analyzer.BasisMeasured, Direction: analyzer.DirImprovement,
					Message:  fmt.Sprintf("deleted callable %s had cyclomatic %d, nesting %d", sc.Base.DisplayName(), b.Cyclomatic, b.MaxNesting),
					Evidence: []analyzer.Evidence{ctx.Ev("base", sc.Base, "deleted")},
				})
			}
		}
	}
	m := func(metric, scope string, b, a int) {
		res.Measurements = append(res.Measurements, analyzer.Measure(metric, scope, b, a, ""))
	}
	m("cyclomatic_sum", "modified_callables", modBefore.cc, modAfter.cc)
	m("decision_points_sum", "modified_callables", modBefore.dec, modAfter.dec)
	m("lines_sum", "modified_callables", modBefore.lines, modAfter.lines)
	m("max_nesting", "modified_callables", modBefore.maxNest, modAfter.maxNest)
	m("cyclomatic_sum", "changed_scope", scopeBefore.cc, scopeAfter.cc)
	m("max_cyclomatic", "changed_scope", scopeBefore.maxCC, scopeAfter.maxCC)
	m("hotspot_callables", "changed_scope", scopeBefore.hot, scopeAfter.hot)
	m("callables", "changed_scope", scopeBefore.n, scopeAfter.n)
	res.Measurements = append(res.Measurements,
		analyzer.Count("callables_modified", "", modAfter.n, "count"),
		analyzer.Count("callables_added", "", added.n, "count"),
		analyzer.Count("callables_deleted", "", deleted.n, "count"),
		analyzer.Count("cyclomatic_sum_added", "added_callables", added.cc, ""),
		analyzer.Count("cyclomatic_sum_deleted", "deleted_callables", deleted.cc, ""),
	)

	modNet := modAfter.cc - modBefore.cc
	totalNet := scopeAfter.cc - scopeBefore.cc
	res.Rationale = append(res.Rationale,
		fmt.Sprintf("modified callables: %d increased, %d decreased; net cyclomatic Δ%+d", increased, decreased, modNet),
		fmt.Sprintf("added callables: %d (cyclomatic %d); deleted callables: %d (cyclomatic %d)", added.n, added.cc, deleted.n, deleted.cc),
		fmt.Sprintf("changed-scope cyclomatic %d→%d (Δ%+d), hotspots %d→%d", scopeBefore.cc, scopeAfter.cc, totalNet, scopeBefore.hot, scopeAfter.hot),
	)
	if modNet < 0 && added.n > 0 && added.cc >= -modNet {
		res.Rationale = append(res.Rationale, fmt.Sprintf("complexity removed from modified callables (Δ%d) is matched by %d added callables (cyclomatic %d): this is extraction/redistribution, not net removal", modNet, added.n, added.cc))
	}
	if errSyms > 0 {
		res.Diagnostics = append(res.Diagnostics, model.Diagnostic{Severity: model.SevWarning, Code: "metrics_in_error_region", Message: fmt.Sprintf("%d changed callables overlap syntax errors; their metrics are unreliable", errSyms)})
	}

	// Hotspot transitions in pre-existing code (modified callables becoming or
	// ceasing to be hotspots, deleted hotspots) weigh more than hotspots in
	// newly added callables, which are capped at -1: new functionality is
	// reported, not treated as significant degradation by itself.
	existingHot := modAfter.hot - modBefore.hot - deleted.hot
	degradeExisting := existingHot > 0 || modNet >= cfg.DegradeAt
	degrade := degradeExisting || added.hot > 0
	improve := existingHot < 0 || modNet <= -cfg.DegradeAt
	score := 0
	switch {
	case degrade && improve:
		res.ScoreRule = "mixed: degradation and improvement signals both present → 0"
	case degrade:
		score = -1
		res.ScoreRule = fmt.Sprintf("pre-existing hotspotΔ>0, modified net cyclomatic ≥ %d, or new hotspot callables → -1", cfg.DegradeAt)
		if existingHot >= 2 || modNet >= cfg.SignificantDegradeAt {
			score = -2
			res.ScoreRule = fmt.Sprintf("pre-existing hotspotΔ ≥ 2 or modified net cyclomatic ≥ %d → -2", cfg.SignificantDegradeAt)
		}
	case improve:
		score = 1
		res.ScoreRule = fmt.Sprintf("pre-existing hotspotΔ<0 or modified net cyclomatic ≤ -%d → +1", cfg.DegradeAt)
		if (existingHot <= -2 || modNet <= -cfg.SignificantDegradeAt) && totalNet < 0 {
			score = 2
			res.ScoreRule = fmt.Sprintf("(pre-existing hotspotΔ ≤ -2 or modified net ≤ -%d) and changed-scope total decreased → +2", cfg.SignificantDegradeAt)
		}
	default:
		res.ScoreRule = "no hotspot change and |modified net cyclomatic| below threshold → 0"
	}
	res.Measurements = append(res.Measurements,
		analyzer.Count("hotspot_delta_preexisting", "modified_and_deleted", existingHot, "callables"),
		analyzer.Count("hotspots_added", "added_callables", added.hot, "callables"),
	)
	res.Score = &score
	res.Summary = fmt.Sprintf("modified callables net cyclomatic Δ%+d; %d added (cc %d), %d deleted (cc %d); hotspots %d→%d",
		modNet, added.n, added.cc, deleted.n, deleted.cc, scopeBefore.hot, scopeAfter.hot)
	return res, nil
}

func metricRows(b, t *model.Metrics) []analyzer.Measurement {
	return []analyzer.Measurement{
		analyzer.Measure("cyclomatic", "", b.Cyclomatic, t.Cyclomatic, ""),
		analyzer.Measure("decision_points", "", b.Decisions, t.Decisions, ""),
		analyzer.Measure("max_nesting", "", b.MaxNesting, t.MaxNesting, ""),
		analyzer.Measure("lines", "", b.Lines, t.Lines, "lines"),
		analyzer.Measure("statements", "", b.Statements, t.Statements, ""),
	}
}
