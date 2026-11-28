package scorev2

import (
	"fmt"
	"io"
	"strings"
)

func locStr(l Location) string {
	if l.StartLine == 0 {
		return l.File
	}
	if l.EndLine > l.StartLine {
		return fmt.Sprintf("%s:%d-%d", l.File, l.StartLine, l.EndLine)
	}
	return fmt.Sprintf("%s:%d", l.File, l.StartLine)
}

func writeFactor(w io.Writer, i int, f Factor, detail bool) {
	title := f.Location.Symbol
	if title == "" {
		title = f.Module
	}
	change := ""
	if f.Change != "" {
		change = fmt.Sprintf(" · **%s** (Δ %+.2f)", f.Change, f.DebtDelta)
	}
	share := ""
	if f.Share > 0 {
		share = fmt.Sprintf(" · %.1f%% of scope (cumulative %.1f%%)", 100*f.Share, 100*f.CumulativeShare)
	}
	fmt.Fprintf(w, "### %d. %s `%s` — %.2f debt%s%s\n\n", i, f.Kind, esc(title), f.Debt, share, change)
	fmt.Fprintf(w, "- **Where:** `%s` (module `%s`, %s)\n", locStr(f.Location), f.Module, f.Language)
	fmt.Fprintf(w, "- **Why:** %s\n", f.Explanation)
	if f.Kind == KindViolation {
		fmt.Fprintf(w, "- **Edge:** `%s` → `%s`\n", f.Metrics2("from"), f.Metrics2("to"))
	} else if len(f.Related) > 0 {
		var rs []string
		for _, r := range f.Related {
			if f.Kind == KindCycle {
				rs = append(rs, fmt.Sprintf("%s at `%s`", r.Symbol, locStr(r)))
			} else {
				rs = append(rs, fmt.Sprintf("`%s` `%s`", r.Symbol, locStr(r)))
			}
		}
		fmt.Fprintf(w, "- **Related:** %s\n", strings.Join(rs, "; "))
	}
	if !detail {
		fmt.Fprintln(w)
		return
	}
	r := f.Remediation
	fmt.Fprintf(w, "- **Fix:** %s. Target: %s.\n", r.Strategy, r.Target)
	for _, s := range r.Steps {
		fmt.Fprintf(w, "  - %s\n", s)
	}
	fmt.Fprintf(w, "- **Done when:** %s\n", strings.Join(r.Acceptance, "; "))
	if r.Caution != "" {
		fmt.Fprintf(w, "- **Caution:** %s\n", r.Caution)
	}
	fmt.Fprintln(w)
}

// WriteExplainMarkdown renders the factor report. detailed limits how many
// factors carry full remediation text (the rest are listed compactly).
func WriteExplainMarkdown(w io.Writer, r *ExplainReport, detailed int) {
	scope := "whole codebase"
	if len(r.Scope.Paths) > 0 {
		scope = strings.Join(r.Scope.Paths, ", ")
	}
	fmt.Fprintf(w, "# Debt factors: %s @ `%s` — %s\n\n", r.Repository, short(r.Revision.Commit), scope)
	writeHeader(w, r.Header)
	writeCoverage(w, r.Coverage)
	t := r.Totals
	fmt.Fprintf(w, "**Scope:** %d source files, %d lines, %d callables (%d hotspots) in %d module(s).\n\n", t.SourceFiles, t.SourceLines, t.Callables, t.Hotspots, len(t.Modules))
	fmt.Fprint(w, "| Debt | complexity | nesting | duplication | architecture | Debt/kLOC | Score |\n|---:|---:|---:|---:|---:|---:|---:|\n")
	fmt.Fprintf(w, "| **%s** | %s | %s | %s | %s | %s | %s |\n\n", f2(t.Debt.Total), f2(t.Debt.Complexity), f2(t.Debt.Nesting), f2(t.Debt.Duplication), f2(t.Debt.Architecture), opt(t.DebtPerKLOC), opt(t.Score))
	n := len(r.Factors)
	if c := r.Concentration; c.Factors > 0 {
		fmt.Fprintf(w, "**Concentration:** %d factor(s) in scope; the largest %d carry 50%% of the debt and the largest %d carry 80%%.\n\n", c.Factors, c.For50Pct, c.For80Pct)
	}
	fmt.Fprintf(w, "## Factors (%d listed", n)
	if r.Omitted > 0 {
		fmt.Fprintf(w, ", %d smaller omitted", r.Omitted)
	}
	fmt.Fprint(w, "; largest first)\n\n")
	for i, f := range r.Factors {
		writeFactor(w, i+1, f, detailed <= 0 || i < detailed)
	}
	fmt.Fprint(w, "## Guidance for fixing\n\n")
	for _, g := range r.Guidance {
		fmt.Fprintf(w, "- %s\n", g)
	}
	fmt.Fprintln(w)
	writeNotes(w, r.Header)
}

// WriteCheckMarkdown renders the PR check, written to be read by people and agents.
func WriteCheckMarkdown(w io.Writer, r *CheckReport) {
	status := "PASS"
	if r.Failed {
		status = "FAIL"
	} else if r.FailOn == "none" {
		status = "ADVISORY"
	}
	scope := "whole repository"
	if len(r.Scope.Paths) > 0 {
		scope = strings.Join(r.Scope.Paths, ", ")
	}
	fmt.Fprintf(w, "# codeeval check: %s (`%s` → `%s`, %s)\n\n", status, short(r.Base.Commit), short(r.Target.Commit), scope)
	budget := ""
	if r.DebtBudget > 0 && (r.FailOn == "net-debt" || r.FailOn == "new-debt") {
		budget = fmt.Sprintf(" (budget %g point(s): the change may add up to that much debt)", r.DebtBudget)
	}
	fmt.Fprintf(w, "- **Fail on:** `%s`%s · scoring `%s` (experimental, uncalibrated)\n", r.FailOn, budget, r.Scoring)
	fmt.Fprintf(w, "- **Change:** %d measured source file(s), %d changed lines", r.Size.ChangedSupportedFiles, r.Size.ChangedSupportedLines)
	if len(r.Size.UnsupportedChangedFiles) > 0 {
		var un []string
		for l, n := range r.Size.UnsupportedChangedFiles {
			un = append(un, fmt.Sprintf("%s %d", l, n))
		}
		fmt.Fprintf(w, "; not measured: %s", strings.Join(sortStrings(un), ", "))
	}
	fmt.Fprintln(w)
	l, c := r.Ledger, r.Changed
	fmt.Fprintf(w, "- **Debt in scope:** %.2f → %.2f (%+.2f)\n", r.ScopeDebtBefore, r.ScopeDebtAfter, r.ScopeDebtAfter-r.ScopeDebtBefore)
	fmt.Fprintf(w, "- **Repository debt:** %+.2f (complexity %+.2f, nesting %+.2f, duplication %+.2f, architecture %+.2f)\n", l.Debt.Total, l.Debt.Complexity, l.Debt.Nesting, l.Debt.Duplication, l.Debt.Architecture)
	fmt.Fprintf(w, "- **Changed callables:** decision points %+d, function base points %+d; modified pre-existing %d (decisions %d → %d), added %d\n\n",
		c.NetDecisionPoints, c.NetBasePoints, c.Modified, c.ModifiedDecisions[0], c.ModifiedDecisions[1], c.Added)
	fmt.Fprint(w, "| Gate | Result |\n|---|---|\n")
	for _, g := range r.Gates {
		res := "pass"
		if !g.Passed {
			res = "**fail** — " + g.Detail
		}
		fmt.Fprintf(w, "| %s | %s |\n", g.Name, res)
	}
	fmt.Fprintln(w)
	if len(r.Factors) == 0 {
		fmt.Fprint(w, "## No debt factors introduced or worsened in scope\n\n")
	} else {
		fmt.Fprintf(w, "## Fix these (%d introduced or worsened; largest first)\n\n", len(r.Factors))
		for i, f := range r.Factors {
			writeFactor(w, i+1, f, true)
		}
	}
	if len(r.Improved) > 0 {
		fmt.Fprintf(w, "## Improved by this change (%d)\n\n", len(r.Improved))
		for _, f := range r.Improved {
			fmt.Fprintf(w, "- %s `%s` at `%s`: %s (Δ %+.2f)\n", f.Kind, esc(f.Location.Symbol+f.Metrics2("")), locStr(f.Location), f.Change, f.DebtDelta)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprint(w, "## How to fix\n\n")
	for _, g := range r.Guidance {
		fmt.Fprintf(w, "- %s\n", g)
	}
	fmt.Fprint(w, "\nThis check measures structure only. It does not run code or establish correctness.\n")
}
