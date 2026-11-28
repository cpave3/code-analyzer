// Package combined joins the scoring/v2 debt view (complexity, nesting,
// duplication, architecture) and the rules/v1 defect findings into one
// report, with a per-file view that ranks files by both. It adds no score
// of its own: debt and defects measure different things and are not summed.
package combined

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/clones"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// Schema of the combined report.
const Schema = "codeeval.analysis/v1"

// FileRow joins one file's debt factors and defect findings.
type FileRow struct {
	File     string   `json:"file"`
	Debt     float64  `json:"debt"`                 // snapshot: factor debt; change: debt delta of introduced/worsened factors
	Factors  int      `json:"factors"`              // debt factors located in the file
	Gate     int      `json:"gateDefects"`          // snapshot: all; change: introduced
	Advisory int      `json:"advisoryDefects"`      // snapshot: all; change: introduced
	Clones   int      `json:"nearOrFragmentClones"` // near-miss and fragment pairs touching the file (change: introduced)
	Kinds    []string `json:"kinds"`                // factor kinds and defect categories present
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
	// Clones: near-miss and fragment duplication (clones/v1, unscored). In
	// change mode the pairs are only those the change introduced.
	Clones *CloneSection `json:"clones"`
	// FactorsWorse / FactorsBetter count debt factors before any --top limit
	// (change mode).
	FactorsWorse  int `json:"factorsIntroducedOrWorsened,omitempty"`
	FactorsBetter int `json:"factorsReducedOrRemoved,omitempty"`
	// Checks is what ran (change mode): all, quick (v2 debt and single-file
	// defect rules on the changed files) or broad (cross-file defect rules
	// and near-duplication; debt shown, never failing).
	Checks  string   `json:"checks,omitempty"`
	Failed  bool     `json:"failed"`
	Reasons []string `json:"failReasons"`
	// FailureFingerprint identifies what failed (introduced gate defects and
	// introduced/worsened debt factors), so a hook can avoid blocking twice
	// on the same items. Empty when the analysis passes.
	FailureFingerprint string           `json:"failureFingerprint,omitempty"`
	Notes              []string         `json:"notes"`
	Revision           *report.Revision `json:"revision,omitempty"`
	Base               *report.Revision `json:"base,omitempty"`
	Target             *report.Revision `json:"target,omitempty"`

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
	sortRows(out)
	return out
}

func sortRows(out []FileRow) {
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
		if a.Clones != b.Clones {
			return a.Clones > b.Clones
		}
		return a.File < b.File
	})
}

func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

// CloneSection holds near-miss and fragment clones.
type CloneSection struct {
	*clones.Result
	Mode string `json:"mode"` // all | introduced
}

// AddClones attaches clone pairs (all, or introduced in change mode) and
// counts them per file.
func (r *Report) AddClones(res *clones.Result) {
	mode := "all"
	if r.Mode == "change" {
		mode = "introduced"
	}
	r.Clones = &CloneSection{Result: res, Mode: mode}
	by := map[string]int{}
	for _, ps := range [][]clones.Pair{res.NearMisses, res.Fragments} {
		for _, p := range ps {
			by[p.A.File]++
			if p.B.File != p.A.File {
				by[p.B.File]++
			}
		}
	}
	seen := map[string]bool{}
	for i := range r.Files {
		r.Files[i].Clones = by[r.Files[i].File]
		seen[r.Files[i].File] = true
	}
	for f, n := range by {
		if !seen[f] {
			r.Files = append(r.Files, FileRow{File: f, Clones: n, Kinds: []string{"near-duplication"}})
		}
	}
	for i := range r.Files {
		if r.Files[i].Clones > 0 {
			addKind(&r.Files[i], "near-duplication")
			sort.Strings(r.Files[i].Kinds)
		}
	}
	sortRows(r.Files)
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
	r.fingerprint(nil, false)
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
	r.fingerprint(ck.Factors, ck.Failed)
	return r
}

// fingerprint hashes the identities of the failing items.
func (r *Report) fingerprint(factors []scorev2.Factor, debtFailed bool) {
	if !r.Failed {
		return
	}
	var ids []string
	if debtFailed {
		for _, f := range factors {
			ids = append(ids, "factor:"+f.ID)
		}
	}
	if r.Defects.Failed {
		for _, f := range r.Defects.Findings {
			if f.Severity == rules.Gate || r.Defects.FailOn == rules.FailAny {
				ids = append(ids, "defect:"+f.Rule+"|"+f.File+"|"+f.Key)
			}
		}
	}
	sort.Strings(ids)
	h := sha256.Sum256([]byte(strings.Join(ids, "\n")))
	r.FailureFingerprint = hex.EncodeToString(h[:6])
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
func WriteMarkdown(w io.Writer, r *Report, files, detailed, clonesTop int) {
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
	fmt.Fprintf(w, "**%s**\n\n", status)
	switch r.Checks {
	case "quick":
		fmt.Fprintln(w, "Quick check: v2 debt, and single-file defect rules on the changed files. Cross-file defect rules and near-duplication did not run (`--checks broad`).")
		fmt.Fprintln(w)
	case "broad":
		fmt.Fprintln(w, "Broad check: cross-file defect rules and near-duplication. Debt is shown but does not fail here (`--checks quick` covers it).")
		fmt.Fprintln(w)
	}
	if r.FailureFingerprint != "" {
		fmt.Fprintf(w, "Failure fingerprint: `%s`\n\n", r.FailureFingerprint)
	}
	fmt.Fprintf(w, "## Summary\n\n| Measure | Value |\n|---|---|\n")
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
		if r.Clones != nil {
			fmt.Fprintf(w, "| Near duplication (unscored) | %d near-miss pairs, %d fragment pairs |\n", len(r.Clones.NearMisses), len(r.Clones.Fragments))
		}
	}
	if ck := r.check; ck != nil {
		l := ck.Ledger.Debt
		fmt.Fprintf(w, "| Change size | %d changed supported source lines in %d files |\n", ck.Size.ChangedSupportedLines, ck.Size.ChangedSupportedFiles)
		fmt.Fprintf(w, "| Debt Δ (codebase) | %+.2f: complexity %+.2f, nesting %+.2f, duplication %+.2f, architecture %+.2f |\n", l.Total, l.Complexity, l.Nesting, l.Duplication, l.Architecture)
		fmt.Fprintf(w, "| Debt in scope | %.2f → %.2f |\n", ck.ScopeDebtBefore, ck.ScopeDebtAfter)
		fmt.Fprintf(w, "| Debt factors | %d introduced or worsened, %d reduced or removed |\n", r.FactorsWorse, r.FactorsBetter)
		fmt.Fprintf(w, "| Decision points (changed scope) | %+d |\n", ck.Changed.NetDecisionPoints)
		fmt.Fprintf(w, "| Defects | introduced %d gate, %d advisory · resolved %d · unchanged %d |\n", d.Totals.Gate, d.Totals.Advisory, len(d.Resolved), d.Persisting)
		if r.Clones != nil {
			fmt.Fprintf(w, "| Near duplication introduced (unscored) | %d near-miss pairs, %d fragment pairs |\n", len(r.Clones.NearMisses), len(r.Clones.Fragments))
		}
	}
	fmt.Fprintln(w)
	title, debtCol := "## Files to look at first", "Debt"
	if r.Mode == "change" {
		title, debtCol = "## Files this change made worse", "Debt Δ (introduced/worsened)"
	}
	fmt.Fprintf(w, "%s\n\nRanked by gate defects, then debt, then advisory defects, then near-duplication pairs.\n\n", title)
	if len(r.Files) == 0 {
		fmt.Fprintln(w, "None.")
	} else {
		fmt.Fprintf(w, "| File | %s | Debt factors | Max cc | Gate defects | Advisory | Near-dup pairs | Kinds |\n|---|---|---|---|---|---|---|---|\n", debtCol)
		for i, f := range r.Files {
			if files > 0 && i == files {
				fmt.Fprintf(w, "\n%d more files in the JSON report.\n", len(r.Files)-files)
				break
			}
			cc := "–"
			if f.MaxCyclo > 0 {
				cc = fmt.Sprint(f.MaxCyclo)
			}
			fmt.Fprintf(w, "| `%s` | %.2f | %d | %s | %d | %d | %d | %s |\n", f.File, f.Debt, f.Factors, cc, f.Gate, f.Advisory, f.Clones, strings.Join(f.Kinds, ", "))
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
	if r.Clones != nil {
		writeClones(w, r.Clones, clonesTop)
	}
	fmt.Fprintln(w, "# Limits of the combined view")
	fmt.Fprintln(w)
	for _, n := range r.Notes {
		fmt.Fprintf(w, "- %s\n", n)
	}
}

func writeClones(w io.Writer, c *CloneSection, top int) {
	what := "all"
	if c.Mode == "introduced" {
		what = "introduced by the change"
	}
	fmt.Fprintf(w, "# Duplication beyond structural clones (%s, unscored, %s)\n\n", c.Version, what)
	var counts []string
	for l, n := range c.Callables {
		counts = append(counts, fmt.Sprintf("%d %s", n, l))
	}
	sort.Strings(counts)
	fmt.Fprintf(w, "Compared %s callables (≥ %d tokens). Near miss: similarity ≥ %.2f on abstracted tokens, structure not identical. Fragment: ≥ %d identical abstracted tokens shared by different code.\n\n",
		strings.Join(counts, ", "), c.Options.MinTokens, c.Options.Similarity, c.Options.FragmentTokens)
	for _, sec := range []struct {
		title string
		ps    []clones.Pair
	}{{"Near-miss clones", c.NearMisses}, {"Fragment clones", c.Fragments}} {
		fmt.Fprintf(w, "## %s (%d)\n\n", sec.title, len(sec.ps))
		if len(sec.ps) == 0 {
			fmt.Fprintln(w, "None.")
			fmt.Fprintln(w)
			continue
		}
		fmt.Fprintln(w, "| A | B | Shared tokens | Similarity |")
		fmt.Fprintln(w, "|---|---|---|---|")
		for i, p := range sec.ps {
			if top > 0 && i == top {
				fmt.Fprintf(w, "\n%d more in the JSON report.\n", len(sec.ps)-top)
				break
			}
			sim := "–"
			if p.Similarity > 0 {
				sim = fmt.Sprintf("%.2f", p.Similarity)
			}
			fmt.Fprintf(w, "| `%s` %s:%d-%d | `%s` %s:%d-%d | %d | %s |\n", p.A.Name, p.A.File, p.A.StartLine, p.A.EndLine, p.B.Name, p.B.File, p.B.StartLine, p.B.EndLine, p.Shared, sim)
		}
		fmt.Fprintln(w)
	}
	for _, n := range c.Notes {
		fmt.Fprintf(w, "- %s\n", n)
	}
	fmt.Fprintln(w)
}
