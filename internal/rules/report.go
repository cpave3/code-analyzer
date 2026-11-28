package rules

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/model"
	"github.com/cameron/code-evaluator/internal/report"
)

// Schema of the defects report.
const Schema = "codeeval.defects/v1"

// Fail-on policies.
const (
	FailGate = "gate" // gate findings (in a change: introduced ones)
	FailAny  = "any"  // any finding (in a change: introduced ones)
	FailNone = "none"
)

// RuleSummary describes one rule in a report.
type RuleSummary struct {
	Rule
	Implementations map[string]string `json:"implementations"` // language -> note
	Findings        int               `json:"findings"`
	Introduced      int               `json:"introduced,omitempty"`
	Resolved        int               `json:"resolved,omitempty"`
}

// Report is the output of `codeeval defects`.
type Report struct {
	Schema       string           `json:"schema"`
	RulesVersion string           `json:"rulesVersion"`
	Tool         report.Tool      `json:"tool"`
	ConfigHash   string           `json:"configHash"`
	Repository   string           `json:"repository"`
	Mode         string           `json:"mode"` // snapshot | change
	Revision     *report.Revision `json:"revision,omitempty"`
	Base         *report.Revision `json:"base,omitempty"`
	Target       *report.Revision `json:"target,omitempty"`
	Scope        []string         `json:"scope"`
	Selected     []string         `json:"selectedRules,omitempty"`
	Rules        []RuleSummary    `json:"rules"`
	Files        []LanguageFiles  `json:"files"`
	Coverage     []Coverage       `json:"coverage"`
	// Findings: all findings in snapshot mode; introduced ones in change mode.
	Findings   []Finding `json:"findings"`
	Resolved   []Finding `json:"resolved,omitempty"`
	Persisting int       `json:"persisting,omitempty"`
	Totals     struct {
		Gate     int `json:"gate"`
		Advisory int `json:"advisory"`
	} `json:"totals"`
	FailOn      string             `json:"failOn"`
	Failed      bool               `json:"failed"`
	Reason      string             `json:"reason"`
	Omitted     int                `json:"findingsOmitted,omitempty"`
	Notes       []string           `json:"notes"`
	Diagnostics []model.Diagnostic `json:"diagnostics,omitempty"`
}

// Notes common to every report.
func notes() []string {
	return []string{
		"findings are concrete, syntax-level defects; they are not a quality score and are not part of scoring/v1 or scoring/v2 debt",
		"no type checking, building or test execution: each rule states what its implementation for a language covers and misses",
		"gate findings are proven from syntax and repository files; advisory findings need a human decision",
		"in change mode findings are matched by rule, file and a line-independent key; a finding in a moved file counts as resolved and introduced",
	}
}

// NewReport assembles a report. In change mode pass before (base result)
// and after; in snapshot mode before is nil.
func NewReport(reg *Registry, before, after *Result, failOn string) *Report {
	r := &Report{Schema: Schema, RulesVersion: Version, Tool: report.Tool{Name: "codeeval", Version: report.ToolVersion}, FailOn: failOn, Notes: notes()}
	r.Files, r.Coverage, r.Diagnostics = after.Files, after.Coverage, after.Diagnostics
	counts := map[string]*RuleSummary{}
	checked := map[string]bool{}
	for _, c := range after.Coverage {
		if c.Status == "checked" {
			checked[c.Rule] = true
		}
	}
	for _, rule := range reg.Rules() {
		if !checked[rule.ID] {
			continue
		}
		s := &RuleSummary{Rule: rule, Implementations: map[string]string{}}
		for _, l := range reg.Languages(rule.ID) {
			impl, _ := reg.Impl(rule.ID, l)
			s.Implementations[l] = impl.Note
		}
		counts[rule.ID] = s
	}
	for _, f := range after.Findings {
		if s := counts[f.Rule]; s != nil {
			s.Findings++
		}
	}
	if before == nil {
		r.Mode = "snapshot"
		r.Findings = after.Findings
	} else {
		r.Mode = "change"
		r.Findings, r.Resolved, r.Persisting = Compare(before.Findings, after.Findings)
		for _, f := range r.Findings {
			counts[f.Rule].Introduced++
		}
		for _, f := range r.Resolved {
			if s := counts[f.Rule]; s != nil {
				s.Resolved++
			}
		}
	}
	for _, rule := range reg.Rules() {
		if s := counts[rule.ID]; s != nil {
			r.Rules = append(r.Rules, *s)
		}
	}
	if r.Findings == nil {
		r.Findings = []Finding{}
	}
	for _, f := range r.Findings {
		if f.Severity == Gate {
			r.Totals.Gate++
		} else {
			r.Totals.Advisory++
		}
	}
	what := "findings"
	if r.Mode == "change" {
		what = "introduced findings"
	}
	switch failOn {
	case FailGate:
		r.Failed = r.Totals.Gate > 0
		r.Reason = fmt.Sprintf("%d gate %s", r.Totals.Gate, what)
	case FailAny:
		r.Failed = len(r.Findings) > 0
		r.Reason = fmt.Sprintf("%d %s", len(r.Findings), what)
	default:
		r.Reason = "fail-on none"
	}
	return r
}

// Limit keeps at most n findings (0 = all), gate findings first.
func (r *Report) Limit(n int) {
	if n > 0 && len(r.Findings) > n {
		r.Omitted = len(r.Findings) - n
		r.Findings = r.Findings[:n]
	}
}

func revLine(rv *report.Revision) string {
	c := rv.Commit
	if len(c) > 12 {
		c = c[:12]
	}
	return fmt.Sprintf("`%s` (%s)", rv.Ref, c)
}

func esc(s string) string { return strings.ReplaceAll(strings.ReplaceAll(s, "|", "\\|"), "\n", " ") }

// WriteMarkdown renders the report.
func WriteMarkdown(w io.Writer, r *Report) {
	fmt.Fprintf(w, "# Defects — %s\n\n", r.Repository)
	if r.Mode == "change" {
		fmt.Fprintf(w, "Change %s → %s · %s · config %s\n\n", revLine(r.Base), revLine(r.Target), r.RulesVersion, r.ConfigHash)
	} else {
		fmt.Fprintf(w, "Revision %s · %s · config %s\n\n", revLine(r.Revision), r.RulesVersion, r.ConfigHash)
	}
	if len(r.Scope) > 0 {
		fmt.Fprintf(w, "Scope: %s\n\n", strings.Join(r.Scope, ", "))
	}
	status := "PASS"
	if r.Failed {
		status = "FAIL"
	}
	fmt.Fprintf(w, "**%s** (fail-on %s): %s.", status, r.FailOn, r.Reason)
	if r.Mode == "change" {
		fmt.Fprintf(w, " Introduced: %d gate, %d advisory · resolved: %d · unchanged: %d.", r.Totals.Gate, r.Totals.Advisory, len(r.Resolved), r.Persisting)
	} else {
		fmt.Fprintf(w, " %d gate, %d advisory.", r.Totals.Gate, r.Totals.Advisory)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| Rule | Severity | Languages | Findings |")
	fmt.Fprintln(w, "|---|---|---|---|")
	for _, s := range r.Rules {
		langs := make([]string, 0, len(s.Implementations))
		for l := range s.Implementations {
			langs = append(langs, l)
		}
		sort.Strings(langs)
		n := fmt.Sprint(s.Findings)
		if r.Mode == "change" {
			n = fmt.Sprintf("%d (+%d / −%d)", s.Findings, s.Introduced, s.Resolved)
		}
		fmt.Fprintf(w, "| `%s` | %s | %s | %s |\n", s.ID, s.Severity, strings.Join(langs, ", "), n)
	}
	fmt.Fprintln(w)
	title := "Findings"
	if r.Mode == "change" {
		title = "Introduced findings"
	}
	writeFindings(w, title, r.Findings)
	if r.Omitted > 0 {
		fmt.Fprintf(w, "%d more findings omitted (--top).\n\n", r.Omitted)
	}
	if r.Mode == "change" && len(r.Resolved) > 0 {
		writeFindings(w, "Resolved findings", r.Resolved)
	}
	fmt.Fprintln(w, "## Coverage")
	fmt.Fprintln(w)
	for _, lf := range r.Files {
		fmt.Fprintf(w, "- %s: %d files checked", lf.Language, lf.Checked)
		if lf.SkippedSyntaxErrors > 0 {
			fmt.Fprintf(w, ", %d skipped for syntax errors", lf.SkippedSyntaxErrors)
		}
		fmt.Fprintln(w)
	}
	for _, c := range r.Coverage {
		if c.Status != "checked" {
			fmt.Fprintf(w, "- `%s` for %s: %s\n", c.Rule, c.Language, strings.ReplaceAll(c.Status, "_", " "))
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Limits")
	fmt.Fprintln(w)
	for _, n := range r.Notes {
		fmt.Fprintf(w, "- %s\n", n)
	}
	fmt.Fprintln(w, "- per-rule coverage notes are in the JSON report (`rules[].implementations`) and `codeeval rules`")
	if len(r.Diagnostics) > 0 {
		fmt.Fprintf(w, "- %d diagnostics in the JSON report\n", len(r.Diagnostics))
	}
}

func writeFindings(w io.Writer, title string, fs []Finding) {
	fmt.Fprintf(w, "## %s\n\n", title)
	if len(fs) == 0 {
		fmt.Fprintln(w, "None.")
		fmt.Fprintln(w)
		return
	}
	fmt.Fprintln(w, "| Severity | Location | Rule | Finding |")
	fmt.Fprintln(w, "|---|---|---|---|")
	for _, f := range fs {
		loc := fmt.Sprintf("%s:%d", f.File, f.Line)
		if f.Symbol != "" {
			loc += " (" + f.Symbol + ")"
		}
		fmt.Fprintf(w, "| %s | %s | `%s` | %s |\n", f.Severity, esc(loc), f.Rule, esc(f.Message))
	}
	fmt.Fprintln(w)
}

// WriteCatalog lists rules and their language implementations.
func WriteCatalog(w io.Writer, reg *Registry, languages []string) {
	fmt.Fprintf(w, "Defect rules (%s). Languages: %s\n\n", Version, strings.Join(languages, ", "))
	for _, rule := range reg.Rules() {
		fmt.Fprintf(w, "%s  [%s]\n    %s\n", rule.ID, rule.Severity, rule.Summary)
		for _, l := range languages {
			if impl, ok := reg.Impl(rule.ID, l); ok {
				fmt.Fprintf(w, "    + %-11s %s\n", l, impl.Note)
			} else {
				fmt.Fprintf(w, "    - %-11s not implemented\n", l)
			}
		}
		for _, lim := range rule.Limits {
			fmt.Fprintf(w, "    ! %s\n", lim)
		}
	}
}
