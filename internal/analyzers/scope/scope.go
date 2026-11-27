// Package scope reports the breadth of a change: files by class, lines,
// symbols by status, public interface changes, name-based affected callers,
// and whether changed units form connected groups in the dependency graph.
// Breadth is reported, not judged; only unrelated formatting/rename churn is
// scored.
package scope

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

// Analyzer is the change-scope analyzer.
type Analyzer struct{}

// New returns the analyzer.
func New() *Analyzer { return &Analyzer{} }

func (Analyzer) Descriptor() analyzer.Descriptor {
	return analyzer.Descriptor{
		ID: "scope", Dimension: analyzer.DimMaintainability, Title: "Change scope and locality", Version: "1",
		Requires:    []lang.Capability{lang.CapSymbols},
		Description: "Files, lines and symbols changed by class; exported interface changes; name-based affected callers; locality of changed units in the dependency graph; formatting-only churn.",
		Classes:     []model.FileClass{model.ClassSource},
	}
}

func (an Analyzer) Analyze(ctx *analyzer.Context) (*analyzer.Result, error) {
	res := &analyzer.Result{}
	cfg := ctx.Config.Scope

	// Files and lines by class.
	byClass := map[model.FileClass][3]int{} // files, added, deleted
	for _, fc := range ctx.Change.Files {
		v := byClass[fc.Class]
		v[0]++
		v[1] += fc.Added
		v[2] += fc.Deleted
		byClass[fc.Class] = v
	}
	var classes []string
	for c := range byClass {
		classes = append(classes, string(c))
	}
	sort.Strings(classes)
	for _, c := range classes {
		v := byClass[model.FileClass(c)]
		res.Measurements = append(res.Measurements,
			analyzer.Count("files_changed", c, v[0], "files"),
			analyzer.Count("lines_added", c, v[1], "lines"),
			analyzer.Count("lines_deleted", c, v[2], "lines"),
		)
	}

	// Symbols by status (production source only).
	d := an.Descriptor()
	status := map[string]int{}
	var formattingUnits = map[string]int{}
	var behaviouralUnits = map[string]bool{}
	var apiAdded, apiRemoved, apiChanged int
	for _, sc := range ctx.Change.Symbols {
		p := sc.Primary()
		side := "target"
		if sc.Target == nil {
			side = "base"
		}
		if !ctx.SymbolAnalyzable(d, side, p) {
			continue
		}
		status[sc.Status]++
		switch sc.Status {
		case change.SymFormatting, change.SymMoved, change.SymRenamed:
			formattingUnits[p.Package]++
		default:
			behaviouralUnits[p.Package] = true
		}
		switch {
		case sc.Status == change.SymAdded && p.Exported:
			apiAdded++
			res.Findings = append(res.Findings, apiFinding(ctx, sc, "scope.api_added", "exported symbol added: "+p.DisplayName()))
		case sc.Status == change.SymDeleted && p.Exported:
			apiRemoved++
			f := apiFinding(ctx, sc, "scope.api_removed", "exported symbol removed: "+p.DisplayName())
			f.Measurements = []analyzer.Measurement{analyzer.Count("referencing_files_in_base", "name-based", callers(ctx, "base", p), "files")}
			res.Findings = append(res.Findings, f)
		case sc.Status == change.SymRenamed && (sc.Base.Exported || sc.Target.Exported):
			apiChanged++
			res.Findings = append(res.Findings, apiFinding(ctx, sc, "scope.api_renamed", fmt.Sprintf("exported symbol renamed/moved: %s → %s", sc.Base.DisplayName(), sc.Target.DisplayName())))
		case sc.Status == change.SymModified && sc.SignatureChanged && sc.Base.Exported:
			apiChanged++
			f := apiFinding(ctx, sc, "scope.api_signature_changed", fmt.Sprintf("exported signature changed: %s %s → %s", p.DisplayName(), sc.Base.Signature, sc.Target.Signature))
			f.Measurements = []analyzer.Measurement{analyzer.Count("referencing_files", "name-based", callers(ctx, "target", p), "files")}
			res.Findings = append(res.Findings, f)
		}
	}
	for _, st := range []string{change.SymAdded, change.SymDeleted, change.SymModified, change.SymRenamed, change.SymMoved, change.SymFormatting} {
		res.Measurements = append(res.Measurements, analyzer.Count("symbols_"+st, "source", status[st], "symbols"))
	}
	res.Measurements = append(res.Measurements,
		analyzer.Count("exported_symbols_added", "source", apiAdded, "symbols"),
		analyzer.Count("exported_symbols_removed", "source", apiRemoved, "symbols"),
		analyzer.Count("exported_symbols_changed", "source", apiChanged, "symbols"),
	)

	// Locality: group changed units (packages) connected within 2 hops in
	// the union of base and target graphs.
	groups, units := locality(ctx, d)
	// The main change is the group with the most behavioural units.
	behavioural := func(g []string) int {
		n := 0
		for _, u := range g {
			if behaviouralUnits[u] {
				n++
			}
		}
		return n
	}
	sort.SliceStable(groups, func(i, j int) bool { return behavioural(groups[i]) > behavioural(groups[j]) })
	res.Measurements = append(res.Measurements,
		analyzer.Count("units_changed", "source", len(units), "units"),
		analyzer.Count("disconnected_unit_groups", "source", len(groups), "groups"),
	)
	churn := 0
	if len(groups) > 1 {
		for _, g := range groups[1:] { // groups[0] is the main change
			onlyFormatting := true
			for _, u := range g {
				if behaviouralUnits[u] {
					onlyFormatting = false
				}
			}
			kind, dir, msg := "scope.disconnected_group", analyzer.DirInfo, "changed units not within two dependency hops of the main change: "+strings.Join(g, ", ")
			conf := "low"
			if onlyFormatting {
				for _, u := range g {
					churn += formattingUnits[u]
				}
				kind, dir = "scope.unrelated_churn", analyzer.DirDegradation
				msg = "formatting/rename-only changes in units unrelated to the main change: " + strings.Join(g, ", ")
				conf = "medium"
			}
			res.Findings = append(res.Findings, analyzer.Finding{
				Kind: kind, Basis: analyzer.BasisHeuristic, Direction: dir, Message: msg, Confidence: conf,
				Caveat:   "locality uses the adapter's dependency graph only; runtime wiring, config and tests are invisible to it",
				Evidence: unitEvidence(ctx, g),
			})
		}
	}
	res.Measurements = append(res.Measurements, analyzer.Count("unrelated_churn_symbols", "source", churn, "symbols"))

	score := 0
	switch {
	case churn >= 2*cfg.FormattingOnlyChurnAt:
		score = -2
		res.ScoreRule = fmt.Sprintf("unrelated formatting/rename-only symbols ≥ %d → -2", 2*cfg.FormattingOnlyChurnAt)
	case churn >= cfg.FormattingOnlyChurnAt || (churn > 0 && len(groups) > 1):
		score = -1
		res.ScoreRule = "formatting/rename-only churn in units unrelated to the main change → -1"
	default:
		res.ScoreRule = "no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)"
	}
	res.Score = &score
	res.Rationale = append(res.Rationale,
		"affected callers are counted by identifier name (qualified by package name across packages), without type resolution",
		"a focused change scores 0: scope has no improvement direction",
	)
	res.Summary = fmt.Sprintf("%d files changed (%d source); symbols +%d/-%d/~%d (renamed %d, formatting-only %d); exported API +%d/-%d/~%d; %d unit group(s)",
		len(ctx.Change.Files), byClass[model.ClassSource][0], status[change.SymAdded], status[change.SymDeleted], status[change.SymModified],
		status[change.SymRenamed], status[change.SymFormatting], apiAdded, apiRemoved, apiChanged, len(groups))
	return res, nil
}

func apiFinding(ctx *analyzer.Context, sc change.SymbolChange, kind, msg string) analyzer.Finding {
	var ev []analyzer.Evidence
	if sc.Base != nil {
		ev = append(ev, ctx.Ev("base", sc.Base, "before"))
	}
	if sc.Target != nil {
		ev = append(ev, ctx.Ev("target", sc.Target, "after"))
	}
	return analyzer.Finding{Kind: kind, Basis: analyzer.BasisMeasured, Direction: analyzer.DirInfo, Message: msg, Evidence: ev}
}

// callers counts files (other than the declaring file) referencing a symbol
// by name: unqualified within its package, "pkg.Name" elsewhere.
func callers(ctx *analyzer.Context, side string, s *model.Symbol) int {
	snap := ctx.Target
	if side == "base" {
		snap = ctx.Base
	}
	declFile := snap.Files[s.File]
	if declFile == nil || declFile.Parsed == nil {
		return 0
	}
	pkgName := declFile.Parsed.PackageName
	n := 0
	for _, p := range snap.Paths {
		f := snap.Files[p]
		if f.Parsed == nil || p == s.File || f.Language != declFile.Language {
			continue
		}
		same := f.Adapter.PackageOf(p) == s.Package
		if (same && f.Parsed.Refs[s.Name] > 0) || (!same && f.Parsed.Refs[pkgName+"."+s.Name] > 0) {
			n++
		}
	}
	return n
}

func locality(ctx *analyzer.Context, d analyzer.Descriptor) ([][]string, []string) {
	unitSet := map[string]string{} // unit -> language
	for _, fc := range ctx.Change.Files {
		if fc.Class != model.ClassSource {
			continue
		}
		ad := ctx.Langs.Adapter(fc.Language)
		if ad == nil {
			continue
		}
		unitSet[ad.PackageOf(fc.Path)] = fc.Language
	}
	var units []string
	for u := range unitSet {
		units = append(units, u)
	}
	sort.Strings(units)
	adj := map[string]map[string]bool{}
	link := func(a, b string) {
		if adj[a] == nil {
			adj[a] = map[string]bool{}
		}
		adj[a][b] = true
	}
	for _, snap := range []map[string]*lang.Graph{ctx.Base.Graphs, ctx.Target.Graphs} {
		for _, g := range snap {
			for f, m := range g.Edges {
				for t := range m {
					link(f, t)
					link(t, f)
				}
			}
		}
	}
	near := func(a, b string) bool {
		if adj[a][b] {
			return true
		}
		for m := range adj[a] {
			if adj[m][b] {
				return true
			}
		}
		return false
	}
	parent := map[string]string{}
	var find func(string) string
	find = func(x string) string {
		if parent[x] == "" || parent[x] == x {
			return x
		}
		parent[x] = find(parent[x])
		return parent[x]
	}
	for i, a := range units {
		for _, b := range units[i+1:] {
			if unitSet[a] == unitSet[b] && near(a, b) {
				ra, rb := find(a), find(b)
				if ra != rb {
					if ra < rb {
						parent[rb] = ra
					} else {
						parent[ra] = rb
					}
				}
			}
		}
	}
	byRoot := map[string][]string{}
	for _, u := range units {
		r := find(u)
		byRoot[r] = append(byRoot[r], u)
	}
	var groups [][]string
	for _, g := range byRoot {
		groups = append(groups, g)
	}
	sort.Slice(groups, func(i, j int) bool {
		if len(groups[i]) != len(groups[j]) {
			return len(groups[i]) > len(groups[j])
		}
		return groups[i][0] < groups[j][0]
	})
	return groups, units
}

func unitEvidence(ctx *analyzer.Context, units []string) []analyzer.Evidence {
	in := map[string]bool{}
	for _, u := range units {
		in[u] = true
	}
	var ev []analyzer.Evidence
	for _, sc := range ctx.Change.Symbols {
		p := sc.Primary()
		if !in[p.Package] {
			continue
		}
		side := "target"
		if sc.Target == nil {
			side = "base"
		}
		ev = append(ev, ctx.Ev(side, p, sc.Status))
		if len(ev) >= 10 {
			break
		}
	}
	if len(ev) == 0 {
		for _, fc := range ctx.Change.Files {
			if ad := ctx.Langs.Adapter(fc.Language); ad != nil && in[ad.PackageOf(fc.Path)] {
				ev = append(ev, ctx.EvLine("target", fc.Path, 1, "changed file"))
			}
		}
	}
	return ev
}
