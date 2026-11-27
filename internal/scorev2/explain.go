package scorev2

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/model"
	"github.com/cameron/code-evaluator/internal/report"
)

// Schemas for factor reports.
const (
	ExplainSchema = "codeeval.explain/v2"
	CheckSchema   = "codeeval.check/v2"
)

// Factor kinds.
const (
	KindComplexity = "complexity" // cyclomatic at/above the hotspot threshold
	KindNesting    = "nesting"    // nesting at/above the nesting threshold
	KindDuplicate  = "duplication"
	KindViolation  = "rule_violation"
	KindCycle      = "cycle"
)

// Location points at source.
type Location struct {
	File      string `json:"file"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
}

// Remediation is deterministic, rule-based guidance for fixing a factor.
type Remediation struct {
	Strategy   string   `json:"strategy"`
	Target     string   `json:"target"`
	Steps      []string `json:"steps"`
	Acceptance []string `json:"acceptance"`
	Caution    string   `json:"caution,omitempty"`
}

// Factor is one item contributing debt to a scope.
type Factor struct {
	ID              string         `json:"id"`
	Kind            string         `json:"kind"`
	Module          string         `json:"module"`
	Language        string         `json:"language"`
	Location        Location       `json:"location"`
	Metrics         map[string]int `json:"metrics,omitempty"`
	Debt            float64        `json:"debt"`
	Share           float64        `json:"shareOfScopeDebt"` // 0–1
	CumulativeShare float64        `json:"cumulativeShare"`
	Explanation     string         `json:"explanation"`
	Related         []Location     `json:"related,omitempty"`
	Remediation     Remediation    `json:"remediation"`
	// check mode only
	Change    string  `json:"change,omitempty"`    // introduced | worsened
	DebtDelta float64 `json:"debtDelta,omitempty"` // after − before for this item
}

// Scope selects files by path prefix or glob.
type Scope struct {
	Paths []string `json:"paths"` // empty = whole codebase
}

// Match reports whether a repository path is in scope. A pattern without
// glob characters matches the path itself or anything below it.
func (s Scope) Match(p string) bool {
	if len(s.Paths) == 0 {
		return true
	}
	for _, pat := range s.Paths {
		pat = strings.TrimSuffix(strings.TrimPrefix(pat, "./"), "/")
		if strings.ContainsAny(pat, "*?") {
			if config.Match(pat, p) {
				return true
			}
			continue
		}
		if pat == "" || pat == "." || p == pat || strings.HasPrefix(p, pat+"/") {
			return true
		}
	}
	return false
}

// ScopeTotals summarises the scope.
type ScopeTotals struct {
	SourceFiles int      `json:"sourceFiles"`
	SourceLines int      `json:"sourceLines"`
	Callables   int      `json:"callables"`
	Hotspots    int      `json:"hotspots"`
	Debt        Debt     `json:"debt"`
	DebtPerKLOC *float64 `json:"debtPerKLOC"`
	Score       *float64 `json:"score"`
	Modules     []string `json:"modules"`
}

// ExplainReport lists the factors behind a scope's debt.
type ExplainReport struct {
	Header
	Revision report.Revision `json:"revision"`
	Scope    Scope           `json:"scope"`
	Coverage Coverage        `json:"coverage"`
	Totals   ScopeTotals     `json:"totals"`
	Factors  []Factor        `json:"factors"`
	Omitted  int             `json:"factorsOmitted"`
	// Concentration: how many of the largest factors make up 50% / 80% of the scope's debt.
	Concentration struct {
		Factors  int `json:"factors"`
		For50Pct int `json:"factorsFor50Percent"`
		For80Pct int `json:"factorsFor80Percent"`
	} `json:"concentration"`
	Guidance []string `json:"guidance"`
}

func guidance() []string {
	return []string{
		"preserve behaviour: add or confirm characterisation tests for a factor's code before changing it; the evaluator does not run code",
		"reduce decision points, do not only move them: a split that leaves the same branching in helpers removes hotspot debt but not logic (check the decision-point ledger with `codeeval check`)",
		"keep helpers cohesive and named for one purpose; avoid creating new hotspot helpers (cc ≥ threshold) while removing a root's debt",
		"do not delete features, inline configuration into data tables, or reshape syntax solely to lower the number; debt is an uncalibrated structural proxy",
		"verify each fix by re-running `codeeval explain` on the same scope (the factor should disappear or shrink) and `codeeval check` on the change (gates pass, no new factors)",
	}
}

func remediation(f *Factor, t Thresholds) Remediation {
	switch f.Kind {
	case KindComplexity:
		cc := f.Metrics["cyclomatic"]
		return Remediation{
			Strategy: "reduce branching in this callable",
			Target:   fmt.Sprintf("cyclomatic ≤ %d (now %d) removes this factor; each point below %d removes 1 debt point", t.HotspotCyclomatic-1, cc, cc),
			Steps: []string{
				"identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function",
				"replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data",
				"remove duplicated or dead branches; merge conditions that lead to the same outcome",
				"keep each extracted helper below the hotspot threshold",
			},
			Acceptance: []string{
				fmt.Sprintf("`%s` cyclomatic ≤ %d, or a justified won't-fix", f.Location.Symbol, t.HotspotCyclomatic-1),
				"no new hotspot callables introduced in the same change",
				"decision points across the changed callables do not increase",
				"existing and characterisation tests pass",
			},
			Caution: "large roots usually need several focused changes; split work per responsibility and verify between steps",
		}
	case KindNesting:
		return Remediation{
			Strategy: "flatten nesting",
			Target:   fmt.Sprintf("max nesting ≤ %d (now %d)", t.HotspotNesting-1, f.Metrics["nesting"]),
			Steps: []string{
				"invert conditions into guard clauses that return or continue early",
				"extract the innermost loop or branch body into a named function",
				"replace nested callbacks with named functions where the language allows",
			},
			Acceptance: []string{fmt.Sprintf("`%s` max nesting ≤ %d", f.Location.Symbol, t.HotspotNesting-1), "tests pass; behaviour unchanged"},
		}
	case KindDuplicate:
		exact := f.Metrics["exact"] == 1
		r := Remediation{
			Strategy: "consolidate duplicated callables",
			Target:   fmt.Sprintf("one implementation instead of %d copies", f.Metrics["members"]),
			Steps: []string{
				"confirm the copies have the same semantics (structural clones may differ in literals or names that matter)",
				"move one implementation to a module every caller may import under the configured architecture rules",
				"replace the other copies with calls; delete the copies",
			},
			Acceptance: []string{"the clone group shrinks to one member (or is removed)", "no new architecture rule violations", "tests for every caller pass"},
		}
		if !exact {
			r.Caution = "structural similarity is a candidate signal, not proof of equivalence; if semantics differ or layering forbids a shared home, record a won't-fix with the reason"
		}
		return r
	case KindViolation:
		return Remediation{
			Strategy:   "remove the forbidden dependency",
			Target:     "no import from " + f.Metrics2("from") + " to " + f.Metrics2("to"),
			Steps:      []string{"find what the importing unit uses from the target", "route it through a permitted layer (an interface owned by the importer, or a service it may depend on)", "remove the import"},
			Acceptance: []string{"the rule violation disappears", "no new cycle or violation elsewhere"},
			Caution:    "confirm the rule reflects intent before changing code; documented rules can drift from accepted practice",
		}
	case KindCycle:
		return Remediation{
			Strategy:   "break the dependency cycle",
			Target:     "no strongly connected component containing this unit",
			Steps:      []string{"pick the cycle edge with the fewest import sites (listed as related locations)", "invert or move the shared type/function so that edge disappears"},
			Acceptance: []string{"the unit is no longer part of a cycle", "no new cycles introduced"},
			Caution:    "in TypeScript directory cycles are legal; treat as coupling to review, not a defect",
		}
	}
	return Remediation{}
}

// Metrics2 holds string metadata for violation factors.
func (f *Factor) Metrics2(k string) string {
	for _, r := range f.Related {
		if r.Symbol == k {
			return r.File
		}
	}
	return ""
}

// factors builds every debt factor of a snapshot within scope.
func factors(snap *index.Snapshot, a *Analysis, scope Scope) []Factor {
	var out []Factor
	w := a.Weights
	for _, c := range a.Callables {
		if !scope.Match(c.File) {
			continue
		}
		m := map[string]int{"cyclomatic": c.Cyclomatic, "decisions": c.Decisions, "nesting": c.Nesting, "lines": c.Lines}
		loc := Location{File: c.File, StartLine: c.StartLine, EndLine: c.EndLine, Symbol: c.Name}
		if c.CCDebt > 0 {
			out = append(out, Factor{ID: "cc:" + c.File + ":" + c.Name, Kind: KindComplexity, Module: c.Module, Language: c.Language, Location: loc, Metrics: m, Debt: c.CCDebt,
				Explanation: fmt.Sprintf("cyclomatic %d is %d point(s) at or above the hotspot threshold %d (%d decision points, %d lines)", c.Cyclomatic, c.Cyclomatic-a.Thresholds.HotspotCyclomatic+1, a.Thresholds.HotspotCyclomatic, c.Decisions, c.Lines)})
		}
		if c.NestDebt > 0 {
			out = append(out, Factor{ID: "nest:" + c.File + ":" + c.Name, Kind: KindNesting, Module: c.Module, Language: c.Language, Location: loc, Metrics: m, Debt: c.NestDebt,
				Explanation: fmt.Sprintf("nesting depth %d is %d level(s) at or above the threshold %d (×%g)", c.Nesting, c.Nesting-a.Thresholds.HotspotNesting+1, a.Thresholds.HotspotNesting, w.NestingExcess)})
		}
	}
	for _, g := range a.Clones {
		n := float64(len(g.Members))
		var in []CloneMember
		var rel []Location
		for _, m := range g.Members {
			rel = append(rel, Location{File: m.File, StartLine: m.StartLine, EndLine: m.EndLine, Symbol: m.Name})
			if scope.Match(m.File) {
				in = append(in, m)
			}
		}
		if len(in) == 0 {
			continue
		}
		debt := 0.0
		for _, m := range in {
			debt += float64(m.Tokens) * (n - 1) / n / w.TokensPerPoint
		}
		ex := 0
		if g.Exact {
			ex = 1
		}
		first := in[0]
		kind := "structural"
		if g.Exact {
			kind = "exact"
		}
		lang := ""
		if f := snap.Files[first.File]; f != nil {
			lang = f.Language
		}
		out = append(out, Factor{ID: "dup:" + first.File + ":" + first.Name, Kind: KindDuplicate, Module: first.Module, Language: lang,
			Location: Location{File: first.File, StartLine: first.StartLine, EndLine: first.EndLine, Symbol: first.Name},
			Metrics:  map[string]int{"members": len(g.Members), "membersInScope": len(in), "exact": ex, "tokens": first.Tokens},
			Debt:     debt, Related: rel,
			Explanation: fmt.Sprintf("%s clone group of %d callables (%d in scope); %.0f redundant tokens attributed to this scope", kind, len(g.Members), len(in), debt*w.TokensPerPoint)})
	}
	for _, v := range a.Violations {
		if !scope.Match(v.File) {
			continue
		}
		out = append(out, Factor{ID: "rule:" + v.From + "->" + v.To, Kind: KindViolation, Module: v.From, Language: v.Language,
			Location: Location{File: v.File, StartLine: v.Line}, Debt: w.RuleViolation,
			Related:     []Location{{Symbol: "from", File: v.From}, {Symbol: "to", File: v.To}},
			Explanation: fmt.Sprintf("import %s → %s violates configured rule %q", v.From, v.To, v.Rule)})
	}
	for _, comp := range a.Cycles {
		for _, u := range comp {
			if !scope.Match(u) && !(len(scope.Paths) > 0 && anyPrefix(scope, u)) {
				continue
			}
			var rel []Location
			lang := ""
			for gid, g := range snap.Graphs {
				if !g.Nodes[u] {
					continue
				}
				lang = gid
				for to, locs := range g.Edges[u] {
					if contains(comp, to) {
						rel = append(rel, Location{File: locs[0].File, StartLine: locs[0].Line, Symbol: fmt.Sprintf("%s → %s (%d import site(s))", u, to, len(locs))})
					}
				}
			}
			sort.Slice(rel, func(i, j int) bool { return rel[i].Symbol < rel[j].Symbol })
			out = append(out, Factor{ID: "cycle:" + u, Kind: KindCycle, Module: u, Language: lang, Location: Location{File: u}, Debt: w.CycleUnit, Related: rel,
				Metrics:     map[string]int{"cycleSize": len(comp)},
				Explanation: fmt.Sprintf("unit is in a %d-unit dependency cycle: %s", len(comp), strings.Join(comp, " ↔ "))})
		}
	}
	for i := range out {
		out[i].Debt = r2(out[i].Debt)
		out[i].Remediation = remediation(&out[i], a.Thresholds)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Debt != out[j].Debt {
			return out[i].Debt > out[j].Debt
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func anyPrefix(s Scope, unit string) bool {
	for _, p := range s.Paths {
		p = strings.TrimSuffix(p, "/")
		if strings.HasPrefix(unit, p) {
			return true
		}
	}
	return false
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

func shares(fs []Factor, total float64) {
	cum := 0.0
	for i := range fs {
		if total > 0 {
			fs[i].Share = r4(fs[i].Debt / total)
			cum += fs[i].Debt / total
			fs[i].CumulativeShare = r4(cum)
		}
	}
}

// NewExplainReport builds the factor report for a scope. top limits the
// factors listed (0 = all); totals always cover every factor.
func NewExplainReport(repo, cfgHash string, rev report.Revision, snap *index.Snapshot, a *Analysis, scope Scope, top int) *ExplainReport {
	r := &ExplainReport{Header: header(ExplainSchema, repo, cfgHash, a), Revision: rev, Scope: scope, Coverage: a.Coverage, Guidance: guidance()}
	if r.Scope.Paths == nil {
		r.Scope.Paths = []string{}
	}
	fs := factors(snap, a, scope)
	var t Components
	mods := map[string]bool{}
	for _, p := range snap.Paths {
		f := snap.Files[p]
		if f.Class != model.ClassSource || !f.Supported() || !scope.Match(p) {
			continue
		}
		t.SourceFiles++
		t.SourceLines += f.Parsed.Lines
		mods[f.Adapter.PackageOf(p)] = true
	}
	for _, s := range snap.Symbols {
		f := snap.Files[s.File]
		if s.Callable() && f != nil && f.Class == model.ClassSource && f.Supported() && scope.Match(s.File) {
			t.Callables++
			if s.Metrics.Cyclomatic >= a.Thresholds.HotspotCyclomatic || s.Metrics.MaxNesting >= a.Thresholds.HotspotNesting {
				t.Hotspots++
			}
		}
	}
	var d Debt
	for _, f := range fs {
		switch f.Kind {
		case KindComplexity:
			d.Complexity += f.Debt
		case KindNesting:
			d.Nesting += f.Debt
		case KindDuplicate:
			d.Duplication += f.Debt
		default:
			d.Architecture += f.Debt
		}
	}
	d.Total = d.Complexity + d.Nesting + d.Duplication + d.Architecture
	r.Totals = ScopeTotals{SourceFiles: t.SourceFiles, SourceLines: t.SourceLines, Callables: t.Callables, Hotspots: t.Hotspots, Debt: d.rounded(), Modules: sortedKeys(mods)}
	if t.SourceLines > 0 {
		dens := d.Total / (float64(t.SourceLines) / 1000)
		r.Totals.DebtPerKLOC = fp(r2(dens))
		r.Totals.Score = fp(r2(100 * a.Weights.HalfScoreDensity / (a.Weights.HalfScoreDensity + dens)))
	}
	shares(fs, d.Total)
	r.Concentration.Factors = len(fs)
	for i, f := range fs {
		if r.Concentration.For50Pct == 0 && f.CumulativeShare >= .5 {
			r.Concentration.For50Pct = i + 1
		}
		if r.Concentration.For80Pct == 0 && f.CumulativeShare >= .8 {
			r.Concentration.For80Pct = i + 1
		}
	}
	r.Factors = limit(fs, top)
	r.Omitted = len(fs) - len(r.Factors)
	return r
}

// ---- PR check ----

// CheckReport lists factors a change introduced or worsened.
type CheckReport struct {
	Header
	Base            report.Revision  `json:"base"`
	Target          report.Revision  `json:"target"`
	Scope           Scope            `json:"scope"`
	Size            ChangeSize       `json:"size"`
	Ledger          Ledger           `json:"ledger"`
	Changed         ChangedCallables `json:"changedCallables"`
	Gates           []Gate           `json:"gates"`
	Factors         []Factor         `json:"factors"`  // introduced or worsened in this change
	Improved        []Factor         `json:"improved"` // factors this change reduced or removed (base-side)
	ScopeDebtBefore float64          `json:"scopeDebtBefore"`
	ScopeDebtAfter  float64          `json:"scopeDebtAfter"`
	Failed          bool             `json:"failed"`
	FailOn          string           `json:"failOn"`
	Guidance        []string         `json:"guidance"`
}

// NewCheckReport compares two snapshots and lists factor changes within scope.
func NewCheckReport(repo, cfgHash string, baseRev, targetRev report.Revision, base, target *index.Snapshot, before, after *Analysis, cs *change.Set, scope Scope, failOn string) *CheckReport {
	r := &CheckReport{Header: header(CheckSchema, repo, cfgHash, after), Base: baseRev, Target: targetRev, Scope: scope, FailOn: failOn, Guidance: guidance()}
	if r.Scope.Paths == nil {
		r.Scope.Paths = []string{}
	}
	r.Size = size(cs)
	r.Ledger = ledger(before.Codebase, after.Codebase)
	r.Changed = changed(base, target, cs)
	r.Gates = gates(r.Ledger, r.Changed)

	// Map target symbol keys to base keys for callables the change touched.
	baseKeyOf := map[string]string{}
	changedTarget := map[string]bool{}
	for _, sc := range cs.Symbols {
		if sc.Target != nil {
			changedTarget[sc.Target.Key] = true
			if sc.Base != nil {
				baseKeyOf[sc.Target.Key] = sc.Base.Key
			}
		}
	}
	bf := map[string]CallableDebt{}
	for _, c := range before.Callables {
		bf[c.Key] = c
	}
	tf := factors(target, after, scope)
	bfs := factors(base, before, scope)
	for _, f := range tf {
		r.ScopeDebtAfter += f.Debt
	}
	for _, f := range bfs {
		r.ScopeDebtBefore += f.Debt
	}
	r.ScopeDebtBefore, r.ScopeDebtAfter = r2(r.ScopeDebtBefore), r2(r.ScopeDebtAfter)
	baseByID := map[string]Factor{}
	for _, f := range bfs {
		baseByID[f.ID] = f
	}
	targetByID := map[string]Factor{}
	for _, f := range tf {
		targetByID[f.ID] = f
	}
	afterCall := map[string]CallableDebt{}
	for _, c := range after.Callables {
		afterCall[c.File+":"+c.Name] = c
	}
	baseClone := map[string]int{}
	for _, g := range before.Clones {
		baseClone[g.Key] = len(g.Members)
	}
	cloneKeyOf := map[string]string{}
	for _, g := range after.Clones {
		for _, m := range g.Members {
			cloneKeyOf["dup:"+m.File+":"+m.Name] = g.Key
		}
	}
	for _, f := range tf {
		switch f.Kind {
		case KindComplexity, KindNesting:
			c := afterCall[f.Location.File+":"+f.Location.Symbol]
			if !changedTarget[c.Key] {
				continue
			}
			prev := 0.0
			var bc *CallableDebt
			if bk, ok := baseKeyOf[c.Key]; ok {
				if b, ok := bf[bk]; ok {
					bc = &b
					if f.Kind == KindComplexity {
						prev = b.CCDebt
					} else {
						prev = b.NestDebt
					}
				}
			}
			if f.Debt > prev {
				f.DebtDelta = r2(f.Debt - prev)
				f.Change = "worsened"
				if prev == 0 {
					f.Change = "introduced"
				}
				if bc != nil && prev > 0 {
					m := map[string]int{}
					for k, v := range f.Metrics {
						m[k] = v
					}
					m["baseCyclomatic"], m["baseNesting"] = bc.Cyclomatic, bc.Nesting
					f.Metrics = m
					acc := append([]string{}, f.Remediation.Acceptance...)
					if f.Kind == KindComplexity {
						f.Remediation.Target = fmt.Sprintf("at minimum, back to the pre-change cyclomatic %d (now %d); ideally ≤ %d, which removes the factor", bc.Cyclomatic, c.Cyclomatic, after.Thresholds.HotspotCyclomatic-1)
						acc[0] = fmt.Sprintf("`%s` cyclomatic ≤ %d (its pre-change value), or a justified won't-fix", c.Name, bc.Cyclomatic)
					} else {
						f.Remediation.Target = fmt.Sprintf("at minimum, back to the pre-change nesting %d (now %d); ideally ≤ %d", bc.Nesting, c.Nesting, after.Thresholds.HotspotNesting-1)
						acc[0] = fmt.Sprintf("`%s` max nesting ≤ %d (its pre-change value)", c.Name, bc.Nesting)
					}
					f.Remediation.Acceptance = acc
				}
				r.Factors = append(r.Factors, f)
			}
		case KindDuplicate:
			k := cloneKeyOf[f.ID]
			if n := f.Metrics["members"]; n > baseClone[k] {
				f.Change, f.DebtDelta = "introduced", f.Debt
				if baseClone[k] > 1 {
					f.Change = "worsened"
				}
				r.Factors = append(r.Factors, f)
			}
		default:
			if _, ok := baseByID[f.ID]; !ok {
				f.Change, f.DebtDelta = "introduced", f.Debt
				r.Factors = append(r.Factors, f)
			}
		}
	}
	// Improvements: base factors that disappeared or shrank (by id).
	for _, f := range bfs {
		if t, ok := targetByID[f.ID]; !ok || t.Debt < f.Debt {
			after := 0.0
			if ok {
				after = t.Debt
			}
			f.DebtDelta = r2(after - f.Debt)
			f.Change = "reduced"
			if !ok {
				f.Change = "removed"
			}
			if f.Kind == KindComplexity || f.Kind == KindNesting || f.Kind == KindDuplicate || !ok {
				r.Improved = append(r.Improved, f)
			}
		}
	}
	sort.SliceStable(r.Factors, func(i, j int) bool { return r.Factors[i].DebtDelta > r.Factors[j].DebtDelta })
	sort.SliceStable(r.Improved, func(i, j int) bool { return r.Improved[i].DebtDelta < r.Improved[j].DebtDelta })
	if r.Factors == nil {
		r.Factors = []Factor{}
	}
	if r.Improved == nil {
		r.Improved = []Factor{}
	}
	switch failOn {
	case "new-debt":
		r.Failed = len(r.Factors) > 0
	case "net-debt":
		r.Failed = r.ScopeDebtAfter > r.ScopeDebtBefore+0.001
	case "gates":
		for _, g := range r.Gates {
			if !g.Passed {
				r.Failed = true
			}
		}
	}
	return r
}

// Limit keeps the first n introduced/worsened and improved factors (0 = all).
func (r *CheckReport) Limit(n int) {
	r.Factors, r.Improved = limit(r.Factors, n), limit(r.Improved, n)
}

// ShortPath trims a path for display.
func ShortPath(p string) string {
	if len(p) <= 60 {
		return p
	}
	return "…/" + path.Base(path.Dir(p)) + "/" + path.Base(p)
}
