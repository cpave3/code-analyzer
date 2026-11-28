// Package combined joins the scoring/v2 debt view (complexity, nesting,
// duplication, architecture) and the rules/v1 defect findings into one
// report, with a per-file view that ranks files by both. It adds no score
// of its own: debt and defects measure different things and are not summed.
package combined

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// Schema of the combined report.
const Schema = "codeeval.analysis/v1"

// FileRow joins one file's debt factors and defect findings.
type FileRow struct {
	File     string   `json:"file"`
	Debt     float64  `json:"debt"`            // snapshot: factor debt; change: debt delta of introduced/worsened factors
	Factors  int      `json:"factors"`         // debt factors located in the file
	Gate     int      `json:"gateDefects"`     // snapshot: all; change: introduced
	Advisory int      `json:"advisoryDefects"` // snapshot: all; change: introduced
	Kinds    []string `json:"kinds"`           // factor kinds and defect categories present
	MaxCyclo int      `json:"maxCyclomatic,omitempty"`
}

// Report is the output of `codeeval analyze`.
type Report struct {
	Schema     string        `json:"schema"`
	Tool       report.Tool   `json:"tool"`
	Repository string        `json:"repository"`
	Mode       string        `json:"mode"` // snapshot | change
	Scope      []string      `json:"scope"`
	Files      []FileRow     `json:"files"` // ranked: gate defects, then debt, then advisory
	Debt       any           `json:"debt"`  // *scorev2.ExplainReport or *scorev2.CheckReport
	Defects    *rules.Report `json:"defects"`
	// FactorsWorse / FactorsBetter count debt factors before any --top limit
	// (change mode).
	FactorsWorse  int              `json:"factorsIntroducedOrWorsened,omitempty"`
	FactorsBetter int              `json:"factorsReducedOrRemoved,omitempty"`
	Failed        bool             `json:"failed"`
	Reasons       []string         `json:"failReasons"`
	Notes         []string         `json:"notes"`
	Revision      *report.Revision `json:"revision,omitempty"`
	Base          *report.Revision `json:"base,omitempty"`
	Target        *report.Revision `json:"target,omitempty"`

	explain *scorev2.ExplainReport
	check   *scorev2.CheckReport
}

func notes() []string {
	return []string{
		"two separate measurements, not one score: v2 debt is an uncalibrated structural proxy (where complexity, nesting, duplication and rule violations sit); defects are concrete syntax-level findings",
		"the per-file table ranks by gate defects, then debt, then advisory defects; it says where to look, not how good a file is",
		"neither measures correctness, behaviour or test quality beyond the listed test-oracle rules",
	}
}

func addKind(row *FileRow, k string) {
	for _, x := range row.Kinds {
		if x == k {
			return
		}
	}
	row.Kinds = append(row.Kinds, k)
}

func rows(factors []scorev2.Factor, findings []rules.Finding, change bool) []FileRow {
	by := map[string]*FileRow{}
	get := func(f string) *FileRow {
		if by[f] == nil {
			by[f] = &FileRow{File: f, Kinds: []string{}}
		}
		return by[f]
	}
	for _, f := range factors {
		if f.Location.File == "" {
			continue
		}
		row := get(f.Location.File)
		if change {
			row.Debt += f.DebtDelta
		} else {
			row.Debt += f.Debt
		}
		row.Factors++
		addKind(row, f.Kind)
		if c := f.Metrics["cyclomatic"]; c > row.MaxCyclo {
			row.MaxCyclo = c
		}
	}
	for _, f := range findings {
		row := get(f.File)
		if f.Severity == rules.Gate {
			row.Gate++
		} else {
			row.Advisory++
		}
		addKind(row, f.Category)
	}
	out := make([]FileRow, 0, len(by))
	for _, r := range by {
		r.Debt = float64(int(r.Debt*100+0.5*sign(r.Debt))) / 100
		sort.Strings(r.Kinds)
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Gate != b.Gate {
			return a.Gate > b.Gate
		}
		if a.Debt != b.Debt {
			return a.Debt > b.Debt
		}
		if a.Advisory != b.Advisory {
			return a.Advisory > b.Advisory
		}
		return a.File < b.File
	})
	return out
}

func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

// Snapshot combines an explain report (all factors, top limit 0) with a
// snapshot defects report.
func Snapshot(ex *scorev2.ExplainReport, d *rules.Report) *Report {
	r := &Report{Schema: Schema, Tool: report.Tool{Name: "codeeval", Version: report.ToolVersion}, Repository: ex.Repository, Mode: "snapshot",
		Scope: ex.Scope.Paths, Debt: ex, Defects: d, Notes: notes(), explain: ex, Revision: d.Revision, Reasons: []string{}}
	r.Files = rows(ex.Factors, d.Findings, false)
	if d.Failed {
		r.Failed = true
		r.Reasons = append(r.Reasons, "defects: "+d.Reason)
	}
	return r
}

// Change combines a check report with a change-mode defects report.
func Change(ck *scorev2.CheckReport, d *rules.Report) *Report {
	r := &Report{Schema: Schema, Tool: report.Tool{Name: "codeeval", Version: report.ToolVersion}, Repository: ck.Repository, Mode: "change",
		Scope: ck.Scope.Paths, Debt: ck, Defects: d, Notes: notes(), check: ck, Base: d.Base, Target: d.Target, Reasons: []string{}}
	r.Files = rows(ck.Factors, d.Findings, true)
	r.FactorsWorse, r.FactorsBetter = len(ck.Factors), len(ck.Improved)
	if ck.Failed {
		r.Failed = true
		r.Reasons = append(r.Reasons, "debt: check failed (fail-on "+ck.FailOn+")")
	}
	if d.Failed {
		r.Failed = true
		r.Reasons = append(r.Reasons, "defects: "+d.Reason)
	}
	return r
}

// demote shifts Markdown headings one level down and drops the first H1.
func demote(md []byte) string {
	var b strings.Builder
	first := true
	for _, line := range strings.Split(string(md), "\n") {
		if strings.HasPrefix(line, "# ") && first {
			first = false
			continue
		}
		if strings.HasPrefix(line, "#") {
			line = "#" + line
		}
		b.WriteString(line + "\n")
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func short(rv *report.Revision) string {
	if rv == nil {
		return ""
	}
	c := rv.Commit
	if len(c) > 12 {
		c = c[:12]
	}
	return fmt.Sprintf("`%s` (%s)", rv.Ref, c)
}

// WriteMarkdown renders the combined report. files limits the per-file
// table, detailed the debt factors with full remediation text.
func WriteMarkdown(w io.Writer, r *Report, files, detailed int) {
	fmt.Fprintf(w, "# Analysis — %s\n\n", r.Repository)
	d := r.Defects
	if r.Mode == "change" {
		fmt.Fprintf(w, "Change %s → %s · %s + %s · config %s\n\n", short(r.Base), short(r.Target), scorev2.Version, d.RulesVersion, d.ConfigHash)
	} else {
		fmt.Fprintf(w, "Revision %s · %s + %s · config %s\n\n", short(r.Revision), scorev2.Version, d.RulesVersion, d.ConfigHash)
	}
	if len(r.Scope) > 0 {
		fmt.Fprintf(w, "Scope: %s\n\n", strings.Join(r.Scope, ", "))
	}
	status := "PASS"
	if r.Failed {
		status = "FAIL: " + strings.Join(r.Reasons, "; ")
	}
	fmt.Fprintf(w, "**%s**\n\n## Summary\n\n| Measure | Value |\n|---|---|\n", status)
	if ex := r.explain; ex != nil {
		t := ex.Totals
		fmt.Fprintf(w, "| Size | %d source files, %d lines, %d callables |\n", t.SourceFiles, t.SourceLines, t.Callables)
		fmt.Fprintf(w, "| Debt (v2) | %.2f total: complexity %.2f, nesting %.2f, duplication %.2f, architecture %.2f |\n", t.Debt.Total, t.Debt.Complexity, t.Debt.Nesting, t.Debt.Duplication, t.Debt.Architecture)
		if t.DebtPerKLOC != nil && t.Score != nil {
			fmt.Fprintf(w, "| Debt per kLOC / score | %.2f / %.1f (comparison only) |\n", *t.DebtPerKLOC, *t.Score)
		}
		fmt.Fprintf(w, "| Hotspots | %d callables at or above the thresholds |\n", t.Hotspots)
		fmt.Fprintf(w, "| Concentration | %d factors; %d make up 50%% of debt, %d make up 80%% |\n", ex.Concentration.Factors, ex.Concentration.For50Pct, ex.Concentration.For80Pct)
		fmt.Fprintf(w, "| Defects | %d gate, %d advisory |\n", d.Totals.Gate, d.Totals.Advisory)
	}
	if ck := r.check; ck != nil {
		l := ck.Ledger.Debt
		fmt.Fprintf(w, "| Change size | %d changed supported source lines in %d files |\n", ck.Size.ChangedSupportedLines, ck.Size.ChangedSupportedFiles)
		fmt.Fprintf(w, "| Debt Δ (codebase) | %+.2f: complexity %+.2f, nesting %+.2f, duplication %+.2f, architecture %+.2f |\n", l.Total, l.Complexity, l.Nesting, l.Duplication, l.Architecture)
		fmt.Fprintf(w, "| Debt in scope | %.2f → %.2f |\n", ck.ScopeDebtBefore, ck.ScopeDebtAfter)
		fmt.Fprintf(w, "| Debt factors | %d introduced or worsened, %d reduced or removed |\n", r.FactorsWorse, r.FactorsBetter)
		fmt.Fprintf(w, "| Decision points (changed scope) | %+d |\n", ck.Changed.NetDecisionPoints)
		fmt.Fprintf(w, "| Defects | introduced %d gate, %d advisory · resolved %d · unchanged %d |\n", d.Totals.Gate, d.Totals.Advisory, len(d.Resolved), d.Persisting)
	}
	fmt.Fprintln(w)
	title, debtCol := "## Files to look at first", "Debt"
	if r.Mode == "change" {
		title, debtCol = "## Files this change made worse", "Debt Δ (introduced/worsened)"
	}
	fmt.Fprintf(w, "%s\n\nRanked by gate defects, then debt, then advisory defects.\n\n", title)
	if len(r.Files) == 0 {
		fmt.Fprintln(w, "None.")
	} else {
		fmt.Fprintf(w, "| File | %s | Debt factors | Max cc | Gate defects | Advisory | Kinds |\n|---|---|---|---|---|---|---|\n", debtCol)
		for i, f := range r.Files {
			if files > 0 && i == files {
				fmt.Fprintf(w, "\n%d more files in the JSON report.\n", len(r.Files)-files)
				break
			}
			cc := "–"
			if f.MaxCyclo > 0 {
				cc = fmt.Sprint(f.MaxCyclo)
			}
			fmt.Fprintf(w, "| `%s` | %.2f | %d | %s | %d | %d | %s |\n", f.File, f.Debt, f.Factors, cc, f.Gate, f.Advisory, strings.Join(f.Kinds, ", "))
		}
	}
	fmt.Fprintln(w)
	var sec bytes.Buffer
	fmt.Fprintln(w, "# Debt (scoring/v2)")
	if r.explain != nil {
		scorev2.WriteExplainMarkdown(&sec, r.explain, detailed)
	} else {
		scorev2.WriteCheckMarkdown(&sec, r.check)
	}
	fmt.Fprintln(w, demote(sec.Bytes()))
	sec.Reset()
	fmt.Fprintln(w, "# Defects (rules/v1)")
	rules.WriteMarkdown(&sec, d)
	fmt.Fprintln(w, demote(sec.Bytes()))
	fmt.Fprintln(w, "# Limits of the combined view")
	fmt.Fprintln(w)
	for _, n := range r.Notes {
		fmt.Fprintf(w, "- %s\n", n)
	}
}
