package scorev2

import (
	"fmt"
	"io"
	"strings"
)

func f2(v float64) string { return fmt.Sprintf("%.2f", v) }

func sf(v float64) string { return fmt.Sprintf("%+.2f", v) }

func si(v int) string { return fmt.Sprintf("%+d", v) }

func opt(p *float64) string {
	if p == nil {
		return "n/a"
	}
	return f2(*p)
}

func optSigned(p *float64) string {
	if p == nil {
		return "n/a"
	}
	return sf(*p)
}

func short(s string) string {
	if len(s) > 12 {
		return s[:12]
	}
	return s
}

func esc(s string) string { return strings.ReplaceAll(s, "|", `\|`) }

func writeHeader(w io.Writer, h Header) {
	fmt.Fprintf(w, "- **Scoring:** `%s` (experimental) · config `%s` · tool %s\n", h.Scoring, h.ConfigHash, h.Tool.Version)
	fmt.Fprintf(w, "- **Weights (uncalibrated):** cyclomatic excess ×%g, nesting excess ×%g, %g redundant tokens = 1 point, rule-violating edge ×%g, unit in cycle ×%g; score half-point at %g debt/kLOC\n",
		h.Weights.CyclomaticExcess, h.Weights.NestingExcess, h.Weights.TokensPerPoint, h.Weights.RuleViolation, h.Weights.CycleUnit, h.Weights.HalfScoreDensity)
	fmt.Fprintf(w, "- **Thresholds:** hotspot cyclomatic %d, nesting %d, clone min tokens %d\n\n", h.Thresholds.HotspotCyclomatic, h.Thresholds.HotspotNesting, h.Thresholds.MinTokens)
}

func writeCoverage(w io.Writer, c Coverage) {
	fmt.Fprintf(w, "**Coverage:** %d of %d production source files measured (%.1f%% of source bytes)", c.SupportedSourceFiles, c.SourceFiles, 100*c.SupportedByteShare)
	var un []string
	for _, l := range c.Languages {
		if !l.Supported {
			un = append(un, fmt.Sprintf("%s %d", l.Language, l.Files))
		}
	}
	if len(un) > 0 {
		fmt.Fprintf(w, "; not measured: %s", strings.Join(un, ", "))
	}
	if c.ParseErrorFiles > 0 {
		fmt.Fprintf(w, "; %d file(s) with parse errors", c.ParseErrorFiles)
	}
	fmt.Fprint(w, "\n\n")
}

func writeNotes(w io.Writer, h Header) {
	fmt.Fprint(w, "## Notes\n\n")
	for _, n := range h.Notes {
		fmt.Fprintf(w, "- %s\n", n)
	}
}

// WriteSnapshotMarkdown renders the absolute report. modules limits the
// module table (0 = all).
func WriteSnapshotMarkdown(w io.Writer, r *SnapshotReport, modules int) {
	fmt.Fprintf(w, "# CodeEval v2 snapshot: %s @ `%s`\n\n", r.Repository, short(r.Revision.Commit))
	fmt.Fprintf(w, "- **Revision:** `%s` tree `%s` %s\n", r.Revision.Commit, short(r.Revision.Tree), r.Revision.CommitTime)
	writeHeader(w, r.Header)
	writeCoverage(w, r.Coverage)
	c, d := r.Codebase.Components, r.Codebase.Debt
	fmt.Fprint(w, "## Codebase\n\n")
	fmt.Fprintf(w, "| Score (0–100) | Debt | Debt / kLOC | Source lines | Callables |\n|---:|---:|---:|---:|---:|\n| **%s** | **%s** | %s | %d | %d |\n\n",
		opt(r.Codebase.Score), f2(d.Total), opt(r.Codebase.DebtPerKLOC), c.SourceLines, c.Callables)
	fmt.Fprint(w, "| Component | Raw | Debt |\n|---|---|---:|\n")
	fmt.Fprintf(w, "| complexity | cyclomatic excess %d (hotspots %d; decision points %d, cyclomatic total %d) | %s |\n", c.CyclomaticExcess, c.Hotspots, c.DecisionPoints, c.CyclomaticTotal, f2(d.Complexity))
	fmt.Fprintf(w, "| nesting | nesting excess %d | %s |\n", c.NestingExcess, f2(d.Nesting))
	fmt.Fprintf(w, "| duplication | redundant tokens %s (structural redundant %d, exact %d, cloned callables %d) | %s |\n", f2(c.RedundantTokens), c.StructRedundant, c.ExactRedundant, c.ClonedCallables, f2(d.Duplication))
	fmt.Fprintf(w, "| architecture | rule-violating edges %d, units in cycles %d (internal edges %d) | %s |\n\n", c.RuleViolations, c.CycleUnits, c.InternalEdges, f2(d.Architecture))

	if len(r.Languages) > 1 {
		fmt.Fprint(w, "## By language\n\n| Language | Score | Debt | Debt/kLOC | Lines | Callables | Hotspots | Complexity | Nesting | Duplication | Architecture |\n|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|\n")
		for _, l := range r.Languages {
			fmt.Fprintf(w, "| %s | %s | %s | %s | %d | %d | %d | %s | %s | %s | %s |\n", l.ID, opt(l.Score), f2(l.Debt.Total), opt(l.DebtPerKLOC), l.Components.SourceLines, l.Components.Callables, l.Components.Hotspots,
				f2(l.Debt.Complexity), f2(l.Debt.Nesting), f2(l.Debt.Duplication), f2(l.Debt.Architecture))
		}
		fmt.Fprint(w, "\nComplexity is measured per language adapter; cyclomatic counts are not directly comparable across languages (see adapter capability notes).\n\n")
	}
	n := len(r.Modules)
	if modules > 0 && n > modules {
		n = modules
	}
	fmt.Fprintf(w, "## Modules (top %d of %d by debt)\n\n", n, len(r.Modules))
	fmt.Fprint(w, "| Module | Score | Debt | Debt/kLOC | Lines | Complexity | Nesting | Duplication | Architecture |\n|---|---:|---:|---:|---:|---:|---:|---:|---:|\n")
	for _, m := range r.Modules[:n] {
		fmt.Fprintf(w, "| `%s` | %s | %s | %s | %d | %s | %s | %s | %s |\n", esc(m.ID), opt(m.Score), f2(m.Debt.Total), opt(m.DebtPerKLOC), m.Components.SourceLines,
			f2(m.Debt.Complexity), f2(m.Debt.Nesting), f2(m.Debt.Duplication), f2(m.Debt.Architecture))
	}
	fmt.Fprint(w, "\n## Largest callable debts\n\n| Callable | Location | cc | nesting | lines | debt |\n|---|---|---:|---:|---:|---:|\n")
	for _, c := range r.Callables {
		fmt.Fprintf(w, "| `%s` | %s:%d-%d | %d | %d | %d | %s |\n", esc(c.Name), esc(c.File), c.StartLine, c.EndLine, c.Cyclomatic, c.Nesting, c.Lines, f2(c.Debt))
	}
	fmt.Fprint(w, "\n## Largest clone groups\n\n| Kind | Redundant tokens | Debt | Members |\n|---|---:|---:|---|\n")
	for _, g := range r.Clones {
		kind := "structural"
		if g.Exact {
			kind = "exact"
		}
		var ms []string
		for _, m := range g.Members {
			ms = append(ms, fmt.Sprintf("`%s` %s:%d", esc(m.Name), esc(m.File), m.StartLine))
		}
		fmt.Fprintf(w, "| %s | %s | %s | %s |\n", kind, f2(g.RedundantTokens), f2(g.Debt), strings.Join(ms, "; "))
	}
	if len(r.Violations) > 0 {
		fmt.Fprint(w, "\n## Rule-violating edges\n\n| Rule | Edge | Import |\n|---|---|---|\n")
		for _, v := range r.Violations {
			fmt.Fprintf(w, "| %s | `%s` → `%s` | %s:%d |\n", esc(v.Rule), v.From, v.To, esc(v.File), v.Line)
		}
	}
	if len(r.Cycles) > 0 {
		fmt.Fprint(w, "\n## Cycles\n\n")
		for _, c := range r.Cycles {
			fmt.Fprintf(w, "- %s\n", strings.Join(c, " ↔ "))
		}
	}
	fmt.Fprintln(w)
	writeNotes(w, r.Header)
}

// WriteDeltaMarkdown renders the comparison report.
func WriteDeltaMarkdown(w io.Writer, r *DeltaReport, modules int) {
	fmt.Fprintf(w, "# CodeEval v2 delta: %s `%s` → `%s`\n\n", r.Repository, short(r.Base.Commit), short(r.Target.Commit))
	fmt.Fprintf(w, "- **Base:** `%s` tree `%s` %s\n- **Target:** `%s` tree `%s` %s\n", r.Base.Commit, short(r.Base.Tree), r.Base.CommitTime, r.Target.Commit, short(r.Target.Tree), r.Target.CommitTime)
	writeHeader(w, r.Header)
	writeCoverage(w, r.CoverageAfter)
	fmt.Fprintf(w, "**Change size:** %d production source file(s), %d measured (%d changed lines)", r.Size.ChangedSourceFiles, r.Size.ChangedSupportedFiles, r.Size.ChangedSupportedLines)
	var un []string
	for l, n := range r.Size.UnsupportedChangedFiles {
		un = append(un, fmt.Sprintf("%s %d", l, n))
	}
	if len(un) > 0 {
		fmt.Fprintf(w, "; not measured: %s", strings.Join(sortStrings(un), ", "))
	}
	fmt.Fprint(w, "\n\n")

	b, a, l := r.Before, r.After, r.Ledger
	fmt.Fprint(w, "## Ledger (after − before; telescopes across commits)\n\n")
	fmt.Fprint(w, "| | Before | After | Δ | Δ per changed kLOC |\n|---|---:|---:|---:|---:|\n")
	for _, row := range debtRows {
		rate := "n/a"
		if r.Rate != nil {
			rate = sf(row.get(*r.Rate))
		}
		delta := sf(row.get(l.Debt))
		if row.bold {
			delta = "**" + delta + "**"
		}
		fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n", row.label, f2(row.get(b.Debt)), f2(row.get(a.Debt)), delta, rate)
	}
	fmt.Fprintf(w, "| score (0–100) | %s | %s | %s | |\n", opt(b.Score), opt(a.Score), optSigned(l.Score))
	fmt.Fprintf(w, "| debt / kLOC | %s | %s | %s | |\n\n", opt(b.DebtPerKLOC), opt(a.DebtPerKLOC), optSigned(l.DebtPerKLOC))

	if len(r.Languages) > 1 {
		fmt.Fprint(w, "| Language | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |\n|---|---|---:|---:|---:|---:|---:|---|\n")
		for _, x := range r.Languages {
			fmt.Fprintf(w, "| %s | %s → %s | %s | %s | %s | %s | %s | %s → %s |\n", x.Language, f2(x.Before.Debt.Total), f2(x.After.Debt.Total), sf(x.Ledger.Debt.Total),
				sf(x.Ledger.Debt.Complexity), sf(x.Ledger.Debt.Nesting), sf(x.Ledger.Debt.Duplication), sf(x.Ledger.Debt.Architecture), opt(x.Before.Score), opt(x.After.Score))
		}
		fmt.Fprintln(w)
	}
	bc, ac, lc := b.Components, a.Components, l.Components
	fmt.Fprint(w, "| Raw component | Before | After | Δ |\n|---|---:|---:|---:|\n")
	rows := []struct {
		n    string
		x, y int
	}{
		{"source lines", bc.SourceLines, ac.SourceLines}, {"callables", bc.Callables, ac.Callables},
		{"decision points", bc.DecisionPoints, ac.DecisionPoints}, {"cyclomatic total", bc.CyclomaticTotal, ac.CyclomaticTotal},
		{"hotspots", bc.Hotspots, ac.Hotspots}, {"cyclomatic excess", bc.CyclomaticExcess, ac.CyclomaticExcess},
		{"nesting excess", bc.NestingExcess, ac.NestingExcess}, {"exact redundant callables", bc.ExactRedundant, ac.ExactRedundant},
		{"structural redundant callables", bc.StructRedundant, ac.StructRedundant}, {"rule-violating edges", bc.RuleViolations, ac.RuleViolations},
		{"units in cycles", bc.CycleUnits, ac.CycleUnits}, {"internal edges", bc.InternalEdges, ac.InternalEdges},
	}
	for _, row := range rows {
		fmt.Fprintf(w, "| %s | %d | %d | %s |\n", row.n, row.x, row.y, si(row.y-row.x))
	}
	fmt.Fprintf(w, "| redundant tokens | %s | %s | %s |\n\n", f2(bc.RedundantTokens), f2(ac.RedundantTokens), sf(lc.RedundantTokens))

	c := r.Changed
	fmt.Fprint(w, "## Changed callables: decision points versus function base points\n\n")
	fmt.Fprintf(w, "- modified pre-existing: %d, decision points %d → %d (%s), cyclomatic %d → %d\n", c.Modified, c.ModifiedDecisions[0], c.ModifiedDecisions[1], si(c.ModifiedDecisions[1]-c.ModifiedDecisions[0]), c.ModifiedCyclomatic[0], c.ModifiedCyclomatic[1])
	fmt.Fprintf(w, "- added: %d (decision points %d); deleted: %d (decision points %d)\n", c.Added, c.AddedDecisions, c.Deleted, c.DeletedDecisions)
	fmt.Fprintf(w, "- **changed scope: decision points %s, function base points %s** (cyclomatic = decisions + 1 per function)\n\n", si(c.NetDecisionPoints), si(c.NetBasePoints))

	fmt.Fprint(w, "## Advisory gates\n\n| Gate | Result | Detail |\n|---|---|---|\n")
	for _, g := range r.Gates {
		res := "pass"
		if !g.Passed {
			res = "**fail**"
		}
		fmt.Fprintf(w, "| %s | %s | %s |\n", g.Name, res, g.Detail)
	}
	n := len(r.Modules)
	if modules > 0 && n > modules {
		n = modules
	}
	label := "worst first"
	if modules > 0 && len(r.Modules) > modules {
		label = fmt.Sprintf("the %d worst and %d best", modules-modules/2, modules/2)
	}
	fmt.Fprintf(w, "\n## Module changes (%d changed; %s)\n\n", len(r.Modules), label)
	fmt.Fprint(w, "| Module | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |\n|---|---|---:|---:|---:|---:|---:|---|\n")
	show := r.Modules[:n]
	if modules > 0 && len(r.Modules) > modules {
		// worst and best halves
		h := modules / 2
		show = append(append([]ModuleDelta{}, r.Modules[:modules-h]...), r.Modules[len(r.Modules)-h:]...)
	}
	for _, m := range show {
		fmt.Fprintf(w, "| `%s` | %s → %s | %s | %s | %s | %s | %s | %s → %s |\n", esc(m.ID), f2(m.DebtBefore), f2(m.DebtAfter), sf(m.DebtDelta.Total),
			sf(m.DebtDelta.Complexity), sf(m.DebtDelta.Nesting), sf(m.DebtDelta.Duplication), sf(m.DebtDelta.Architecture), opt(m.ScoreBefore), opt(m.ScoreAfter))
	}
	if len(r.Steps) > 0 {
		fmt.Fprint(w, "\n## First-parent steps\n\n| # | Commit | Subject (untrusted) | Changed lines | Δ debt | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Δ score | Gates failed |\n|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|\n")
		for _, s := range r.Steps {
			fmt.Fprintf(w, "| %d | `%s` | %s | %d | %s | %s | %s | %s | %s | %s | %s |\n", s.Index, short(s.Commit), esc(trunc(s.Subject, 70)), s.Size.ChangedSupportedLines,
				sf(s.Ledger.Debt.Total), sf(s.Ledger.Debt.Complexity), sf(s.Ledger.Debt.Nesting), sf(s.Ledger.Debt.Duplication), sf(s.Ledger.Debt.Architecture), optSigned(s.Ledger.Score), strings.Join(s.Failed, ", "))
		}
		if sc := r.StepCheck; sc != nil {
			fmt.Fprintf(w, "\n**Telescoping check:** Σ step Δdebt %s vs net %s → %v (chain starts at base: %v)\n", sf(sc.SumOfStepDebt), sf(sc.NetDebt), sc.Telescopes, sc.ChainStartsAtBase)
		}
	}
	fmt.Fprintln(w)
	writeNotes(w, r.Header)
}

func trunc(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	return string([]rune(s)[:n-1]) + "…"
}

func sortStrings(s []string) []string {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
	return s
}

// debtRows are the ledger rows of the delta report, in order.
var debtRows = []struct {
	label string
	bold  bool
	get   func(Debt) float64
}{
	{"**Total debt**", true, func(d Debt) float64 { return d.Total }},
	{"complexity debt", false, func(d Debt) float64 { return d.Complexity }},
	{"nesting debt", false, func(d Debt) float64 { return d.Nesting }},
	{"duplication debt", false, func(d Debt) float64 { return d.Duplication }},
	{"architecture debt", false, func(d Debt) float64 { return d.Architecture }},
}
