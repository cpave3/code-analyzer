package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/analyzer"
)

// MaxFindingsPerAnalyzer bounds findings shown in Markdown (JSON has all).
const MaxFindingsPerAnalyzer = 25

func fmtScore(s *int) string {
	if s == nil {
		return "n/a"
	}
	if *s > 0 {
		return fmt.Sprintf("+%d", *s)
	}
	return fmt.Sprintf("%d", *s)
}

func fmtPtr(p *int) string {
	if p == nil {
		return "–"
	}
	return fmt.Sprint(*p)
}

func fmtDelta(p *int) string {
	if p == nil {
		return "–"
	}
	if *p > 0 {
		return fmt.Sprintf("+%d", *p)
	}
	return fmt.Sprint(*p)
}

// esc neutralises Markdown table/HTML characters in untrusted text.
func esc(s string) string {
	r := strings.NewReplacer("|", "\\|", "<", "&lt;", ">", "&gt;", "\n", " ", "`", "'")
	return r.Replace(s)
}

func short(s string) string {
	if len(s) > 12 {
		return s[:12]
	}
	return s
}

// WriteMarkdown renders the human-readable summary.
func WriteMarkdown(w io.Writer, r *Report) {
	fmt.Fprintf(w, "# CodeEval report: %s\n\n", esc(r.Repository))
	fmt.Fprintf(w, "- **Base:** `%s` (%s) tree `%s` %s\n", r.Base.Commit, esc(r.Base.Ref), short(r.Base.Tree), r.Base.CommitTime)
	fmt.Fprintf(w, "- **Target:** `%s` (%s) tree `%s` %s\n", r.Target.Commit, esc(r.Target.Ref), short(r.Target.Tree), r.Target.CommitTime)
	fmt.Fprintf(w, "- **Mode:** %s · schema `%s` · %s · config v%d `%s` · tool %s\n", r.Mode, r.Schema, r.ScoringVersion, r.ConfigVersion, r.ConfigHash, r.Tool.Version)
	if r.Timing != nil {
		fmt.Fprintf(w, "- **Run:** started %s, %.1fs, %d files parsed, cache %d hits / %d misses\n", r.Timing.StartedAt, float64(r.Timing.ElapsedMillis)/1000, r.Timing.FilesParsed, r.Timing.CacheHits, r.Timing.CacheMisses)
	}
	fmt.Fprintf(w, "- **Change:** %d files, +%d/-%d lines, %d changed symbols\n\n", r.Change.FilesChanged, r.Change.LinesAdded, r.Change.LinesDeleted, len(r.Change.Symbols))
	fmt.Fprintf(w, "Scale: %s\n\n", r.ScoreScale)

	fmt.Fprintln(w, "## Scores by dimension (independent; no overall score)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| Dimension | Analyzer | Status | Score | Summary |")
	fmt.Fprintln(w, "|---|---|---|---:|---|")
	for _, d := range r.Dimensions {
		if len(d.Analyzers) == 0 {
			fmt.Fprintf(w, "| %s | – | no analyzer | n/a | |\n", d.Dimension)
		}
		for _, a := range d.Analyzers {
			fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n", d.Dimension, a.Analyzer, a.Status, fmtScore(a.Score), esc(a.Summary))
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "## Language coverage")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| Language | Supported | Files in target | Changed files | Changed lines | Note |")
	fmt.Fprintln(w, "|---|---|---:|---:|---:|---|")
	for _, l := range r.Languages {
		sup := "no"
		if l.Supported {
			sup = "yes"
		}
		fmt.Fprintf(w, "| %s | %s | %d | %d | %d | %s |\n", l.ID, sup, l.FilesTotal, l.FilesChanged, l.LinesChanged, esc(l.Note))
	}
	fmt.Fprintln(w)

	classes := map[string][3]int{}
	for _, f := range r.Change.Files {
		k := string(f.Class)
		v := classes[k]
		v[0]++
		v[1] += f.Added
		v[2] += f.Deleted
		classes[k] = v
	}
	var keys []string
	for k := range classes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Fprintln(w, "Changed files by class: ")
	for _, k := range keys {
		v := classes[k]
		fmt.Fprintf(w, "%s %d (+%d/-%d); ", k, v[0], v[1], v[2])
	}
	fmt.Fprint(w, "\n\n")

	for _, res := range r.Results {
		writeResult(w, res)
	}

	if len(r.Change.Symbols) > 0 {
		fmt.Fprintln(w, "## Changed symbols")
		fmt.Fprintln(w)
		fmt.Fprintln(w, "| Status | Kind | Symbol | Location | Base location | Class |")
		fmt.Fprintln(w, "|---|---|---|---|---|---|")
		limit := 200
		for i, s := range r.Change.Symbols {
			if i >= limit {
				fmt.Fprintf(w, "\n… %d more in JSON\n", len(r.Change.Symbols)-limit)
				break
			}
			loc, bloc := "", ""
			if s.File != "" {
				loc = fmt.Sprintf("%s:%d-%d", s.File, s.StartLine, s.EndLine)
			}
			if s.BaseFile != "" {
				bloc = fmt.Sprintf("%s:%d-%d", s.BaseFile, s.BaseStartLine, s.BaseEndLine)
			}
			name := s.Name
			if s.BaseName != "" {
				name = s.BaseName + " → " + s.Name
			}
			fmt.Fprintf(w, "| %s | %s | %s | %s | %s | %s |\n", s.Status, s.Kind, esc(name), esc(loc), esc(bloc), s.Class)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "## Diagnostics")
	fmt.Fprintln(w)
	if len(r.Diagnostics.Change) == 0 {
		fmt.Fprintln(w, "No diagnostics for changed files.")
	}
	for i, d := range r.Diagnostics.Change {
		if i >= 50 {
			fmt.Fprintf(w, "- … %d more in JSON\n", len(r.Diagnostics.Change)-50)
			break
		}
		loc := d.File
		if d.Line > 0 {
			loc = fmt.Sprintf("%s:%d", d.File, d.Line)
		}
		fmt.Fprintf(w, "- **%s** `%s` %s %s %s\n", d.Severity, d.Code, d.Side, esc(loc), esc(d.Message))
	}
	if len(r.Diagnostics.SnapshotCountByCode) > 0 {
		var ks []string
		for k := range r.Diagnostics.SnapshotCountByCode {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		fmt.Fprint(w, "\nSnapshot-wide diagnostic counts: ")
		for _, k := range ks {
			fmt.Fprintf(w, "`%s`=%d ", k, r.Diagnostics.SnapshotCountByCode[k])
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Limitations")
	fmt.Fprintln(w)
	for _, l := range r.Limitations {
		fmt.Fprintf(w, "- %s\n", esc(l))
	}
}

func writeResult(w io.Writer, res *analyzer.Result) {
	fmt.Fprintf(w, "## %s — %s (`%s`)\n\n", res.Dimension, res.Title, res.Analyzer)
	fmt.Fprintf(w, "**Status:** %s · **Score:** %s · **Rule:** %s\n\n", res.Status, fmtScore(res.Score), esc(res.ScoreRule))
	fmt.Fprintf(w, "%s\n\n", esc(res.Summary))
	c := res.Coverage
	fmt.Fprintf(w, "Coverage: %d changed code file(s) in scope, %d analysed, %d with parse errors, unsupported %v, excluded by class %v, non-code %d\n\n",
		c.ChangedCodeFiles, c.Analyzed, c.WithParseErrors, mapStr(c.UnsupportedByLanguage), mapStr(c.ExcludedByClass), c.NonCodeFiles)
	for _, r := range res.Rationale {
		fmt.Fprintf(w, "- %s\n", esc(r))
	}
	if len(res.Rationale) > 0 {
		fmt.Fprintln(w)
	}
	if len(res.Measurements) > 0 {
		fmt.Fprintln(w, "| Metric | Scope | Before | After | Δ | Value |")
		fmt.Fprintln(w, "|---|---|---:|---:|---:|---:|")
		for _, m := range res.Measurements {
			fmt.Fprintf(w, "| %s | %s | %s | %s | %s | %s |\n", m.Metric, esc(m.Scope), fmtPtr(m.Before), fmtPtr(m.After), fmtDelta(m.Delta), fmtPtr(m.Value))
		}
		fmt.Fprintln(w)
	}
	if len(res.Findings) > 0 {
		fmt.Fprintf(w, "Findings (%d):\n\n", len(res.Findings))
		for i, f := range res.Findings {
			if i >= MaxFindingsPerAnalyzer {
				fmt.Fprintf(w, "- … %d more in JSON\n", len(res.Findings)-MaxFindingsPerAnalyzer)
				break
			}
			tag := f.Basis
			if f.Confidence != "" {
				tag += ", confidence " + f.Confidence
			}
			fmt.Fprintf(w, "- **%s** [%s; %s] %s\n", f.Kind, f.Direction, tag, esc(f.Message))
			for _, e := range f.Evidence {
				sym := ""
				if e.Symbol != "" {
					sym = " `" + esc(e.Symbol) + "`"
				}
				fmt.Fprintf(w, "  - %s `%s` %s:%d-%d%s (%s)\n", e.Side, short(e.Commit), esc(e.File), e.StartLine, e.EndLine, sym, esc(e.Role))
			}
			if f.Caveat != "" {
				fmt.Fprintf(w, "  - _caveat: %s_\n", esc(f.Caveat))
			}
		}
		fmt.Fprintln(w)
	}
	for _, d := range res.Diagnostics {
		fmt.Fprintf(w, "- diagnostic **%s** `%s`: %s\n", d.Severity, d.Code, esc(d.Message))
	}
	if len(res.Diagnostics) > 0 {
		fmt.Fprintln(w)
	}
}

func mapStr(m map[string]int) string {
	if len(m) == 0 {
		return "{}"
	}
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	var parts []string
	for _, k := range ks {
		parts = append(parts, fmt.Sprintf("%s:%d", k, m[k]))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// WriteHistoryMarkdown renders a first-parent series.
func WriteHistoryMarkdown(w io.Writer, h *History) {
	fmt.Fprintf(w, "# CodeEval history: %s\n\n", esc(h.Net.Repository))
	fmt.Fprintf(w, "- **Base:** `%s` (%s) %s\n- **Target:** `%s` (%s) %s\n", h.Base.Commit, esc(h.Base.Ref), h.Base.CommitTime, h.Target.Commit, esc(h.Target.Ref), h.Target.CommitTime)
	if h.Timing != nil {
		fmt.Fprintf(w, "- **Run:** started %s, %.1fs total, %d files parsed, cache %d hits / %d misses\n", h.Timing.StartedAt, float64(h.Timing.ElapsedMillis)/1000, h.Timing.FilesParsed, h.Timing.CacheHits, h.Timing.CacheMisses)
	}
	fmt.Fprintln(w)
	for _, n := range h.Notes {
		fmt.Fprintf(w, "- %s\n", esc(n))
	}
	fmt.Fprintln(w)
	var ids []string
	for _, r := range h.Net.Results {
		ids = append(ids, r.Analyzer)
	}
	fmt.Fprintln(w, "## First-parent steps")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "| # | Commit | Subject (untrusted) | Files | +/- | Go files | %s |\n", strings.Join(ids, " | "))
	fmt.Fprintf(w, "|---|---|---|---:|---:|---:|%s\n", strings.Repeat("---:|", len(ids)))
	row := func(label, commit, subject string, r *Report) {
		goFiles := 0
		for _, l := range r.Languages {
			if l.ID == "go" {
				goFiles = l.FilesChanged
			}
		}
		var cells []string
		for _, res := range r.Results {
			cells = append(cells, fmtScore(res.Score))
		}
		fmt.Fprintf(w, "| %s | `%s` | %s | %d | +%d/-%d | %d | %s |\n", label, short(commit), esc(subject), r.Change.FilesChanged, r.Change.LinesAdded, r.Change.LinesDeleted, goFiles, strings.Join(cells, " | "))
	}
	for _, s := range h.Steps {
		row(fmt.Sprint(s.Index), s.Commit, s.Subject, s.Report)
	}
	row("net", h.Net.Target.Commit, "base → target (net)", h.Net)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Key measurements per step")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| # | Commit | dup Δ (struct redundancy) | modified cc Δ | added callables (cc) | deleted callables (cc) | hotspots before→after | edges +/- | cycles Δ | obsolete kept | API +/-/~ |")
	fmt.Fprintln(w, "|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|")
	key := func(label, commit string, r *Report) {
		fmt.Fprintf(w, "| %s | `%s` | %s | %s | %s (%s) | %s (%s) | %s→%s | +%s/-%s | %s | %s | +%s/-%s/~%s |\n", label, short(commit),
			delta(r, "duplication", "structural_redundant_callables", ""),
			delta(r, "complexity", "cyclomatic_sum", "modified_callables"),
			val(r, "complexity", "callables_added"), val(r, "complexity", "cyclomatic_sum_added"),
			val(r, "complexity", "callables_deleted"), val(r, "complexity", "cyclomatic_sum_deleted"),
			before(r, "complexity", "hotspot_callables"), after(r, "complexity", "hotspot_callables"),
			val(r, "dependencies", "edges_added"), val(r, "dependencies", "edges_removed"),
			delta(r, "dependencies", "cycles", ""),
			val(r, "obsolete", "callables_left_unreferenced"),
			val(r, "scope", "exported_symbols_added"), val(r, "scope", "exported_symbols_removed"), val(r, "scope", "exported_symbols_changed"))
	}
	for _, s := range h.Steps {
		key(fmt.Sprint(s.Index), s.Commit, s.Report)
	}
	key("net", h.Net.Target.Commit, h.Net)
	fmt.Fprintln(w)
}

func find(r *Report, analyzerID, metric, scope string) *analyzer.Measurement {
	for _, res := range r.Results {
		if res.Analyzer != analyzerID {
			continue
		}
		for i, m := range res.Measurements {
			if m.Metric == metric && (scope == "" || m.Scope == scope) {
				return &res.Measurements[i]
			}
		}
	}
	return nil
}

func delta(r *Report, a, m, s string) string {
	if x := find(r, a, m, s); x != nil {
		return fmtDelta(x.Delta)
	}
	return "–"
}
func val(r *Report, a, m string) string {
	if x := find(r, a, m, ""); x != nil {
		return fmtPtr(x.Value)
	}
	return "–"
}
func before(r *Report, a, m string) string {
	if x := find(r, a, m, ""); x != nil {
		return fmtPtr(x.Before)
	}
	return "–"
}
func after(r *Report, a, m string) string {
	if x := find(r, a, m, ""); x != nil {
		return fmtPtr(x.After)
	}
	return "–"
}
