package scorev2

import (
	"math"
	"sort"

	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/model"
	"github.com/cameron/code-evaluator/internal/report"
)

// Header is common report metadata.
type Header struct {
	Schema     string      `json:"schema"`
	Scoring    string      `json:"scoringVersion"`
	Tool       report.Tool `json:"tool"`
	ConfigHash string      `json:"configHash"`
	Repository string      `json:"repository"`
	Weights    Weights     `json:"weights"`
	Thresholds Thresholds  `json:"thresholds"`
	Notes      []string    `json:"notes"`
}

// SnapshotReport is the absolute measurement of one revision.
type SnapshotReport struct {
	Header
	Revision   report.Revision `json:"revision"`
	Coverage   Coverage        `json:"coverage"`
	Codebase   Unit            `json:"codebase"`
	Languages  []Unit          `json:"languages"`
	Modules    []Unit          `json:"modules"`
	Callables  []CallableDebt  `json:"topCallables"`
	Clones     []CloneGroup    `json:"topCloneGroups"`
	Violations []Violation     `json:"violations"`
	Cycles     [][]string      `json:"cycles"`
}

// ChangedCallables decomposes complexity in the callables a change touched.
type ChangedCallables struct {
	Modified           int    `json:"modified"`
	ModifiedDecisions  [2]int `json:"modifiedDecisionsBeforeAfter"`
	ModifiedCyclomatic [2]int `json:"modifiedCyclomaticBeforeAfter"`
	Added              int    `json:"added"`
	AddedDecisions     int    `json:"addedDecisions"`
	Deleted            int    `json:"deleted"`
	DeletedDecisions   int    `json:"deletedDecisions"`
	NetDecisionPoints  int    `json:"netDecisionPoints"` // changed scope, after − before
	NetBasePoints      int    `json:"netFunctionBasePoints"`
}

// ModuleDelta is one module's change.
type ModuleDelta struct {
	ID          string   `json:"id"`
	Language    string   `json:"language"`
	DebtBefore  float64  `json:"debtBefore"`
	DebtAfter   float64  `json:"debtAfter"`
	DebtDelta   Debt     `json:"debtDelta"`
	ScoreBefore *float64 `json:"scoreBefore"`
	ScoreAfter  *float64 `json:"scoreAfter"`
}

// LanguageLedger is one language's before/after totals.
type LanguageLedger struct {
	Language string `json:"language"`
	Before   Unit   `json:"before"`
	After    Unit   `json:"after"`
	Ledger   Ledger `json:"ledger"`
}

func languageLedgers(before, after *Analysis) []LanguageLedger {
	bm, am := map[string]Unit{}, map[string]Unit{}
	for _, u := range before.Languages {
		bm[u.ID] = u
	}
	for _, u := range after.Languages {
		am[u.ID] = u
	}
	keys := map[string]bool{}
	for k := range bm {
		keys[k] = true
	}
	for k := range am {
		keys[k] = true
	}
	out := []LanguageLedger{}
	for _, k := range sortedKeys(keys) {
		b, a := bm[k], am[k]
		out = append(out, LanguageLedger{Language: k, Before: roundUnit(b), After: roundUnit(a), Ledger: ledger(b, a)})
	}
	return out
}

// Gate is an advisory pass/fail check.
type Gate struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
	Detail string `json:"detail"`
}

// Ledger is after − before for the codebase.
type Ledger struct {
	Components  Components `json:"components"`
	Debt        Debt       `json:"debt"`
	DebtPerKLOC *float64   `json:"debtPerKLOC"`
	Score       *float64   `json:"score"`
}

// ChangeSize is the size of the change in measured production source.
type ChangeSize struct {
	ChangedSourceFiles      int            `json:"changedSourceFiles"`
	ChangedSupportedFiles   int            `json:"changedSupportedSourceFiles"`
	ChangedSupportedLines   int            `json:"changedSupportedSourceLines"` // added + deleted
	UnsupportedChangedFiles map[string]int `json:"unsupportedChangedSourceFilesByLanguage"`
}

// Step is one first-parent step's ledger.
type Step struct {
	Index     int              `json:"index"`
	Parent    string           `json:"parent"`
	Commit    string           `json:"commit"`
	Subject   string           `json:"subject"`
	Size      ChangeSize       `json:"size"`
	Ledger    Ledger           `json:"ledger"`
	Languages []LanguageLedger `json:"languages"`
	Rate      *Debt            `json:"debtPerChangedKLOC"`
	Failed    []string         `json:"gatesFailed"`
}

// StepCheck verifies telescoping.
type StepCheck struct {
	ChainStartsAtBase bool    `json:"chainStartsAtBase"`
	SumOfStepDebt     float64 `json:"sumOfStepDebtDelta"`
	NetDebt           float64 `json:"netDebtDelta"`
	Telescopes        bool    `json:"telescopes"`
}

// DeltaReport compares two revisions.
type DeltaReport struct {
	Header
	Base          report.Revision  `json:"base"`
	Target        report.Revision  `json:"target"`
	CoverageAfter Coverage         `json:"coverageAfter"`
	Size          ChangeSize       `json:"size"`
	Before        Unit             `json:"before"`
	After         Unit             `json:"after"`
	Ledger        Ledger           `json:"ledger"`
	Languages     []LanguageLedger `json:"languages"`
	Rate          *Debt            `json:"debtPerChangedKLOC"`
	Changed       ChangedCallables `json:"changedCallables"`
	Gates         []Gate           `json:"gates"`
	Modules       []ModuleDelta    `json:"modules"`
	Steps         []Step           `json:"steps,omitempty"`
	StepCheck     *StepCheck       `json:"stepCheck,omitempty"`
}

func notes() []string {
	return []string{
		"experimental scoring/v2: separate from and not comparable to scoring/v1 −2..+2 scores; v1 reports are unchanged",
		"debt points = complexity (Σ max(0, cyclomatic − H + 1): debt exactly when a callable is a hotspot) + nesting (Σ max(0, depth − N + 1) × weight) + duplication (structural-clone redundant tokens ÷ tokens-per-point) + architecture (configured-rule-violating edges and units in cycles × weights)",
		"debt is additive across modules and telescopes across commits (Δ = after − before of one snapshot function); telescoping is mathematical consistency under a fixed model and config, not empirical calibration",
		"debt is an uncalibrated, size-sensitive structural proxy: deleting features or code, or reshaping syntax (e.g. splitting functions below the thresholds), can lower it without improving engineering; do not optimise it alone — keep the component and raw evidence, behaviour/test, coverage and scope checks",
		"score = 100 × h ÷ (h + debt per 1,000 supported source lines), h = halfScoreDensity: a size-normalised view for comparing modules or codebases; it can be diluted by adding simple code and is never a target on its own",
		"weights are uncalibrated defaults; the score is a structural index, not a measure of correctness, design quality or behaviour",
		"only languages with adapters are measured; see coverage — unsupported source is not scored and not neutral",
		"architecture debt reflects only the configured import-level rules; rules the configuration cannot express are not measured",
		"complexity growth below the hotspot thresholds carries no debt; it is reported as decision points and checked by the existing-callable branch gate",
		"gates are advisory checks, not merge decisions",
	}
}

func header(schema, repo, cfgHash string, a *Analysis) Header {
	return Header{Schema: schema, Scoring: Version, Tool: report.Tool{Name: "codeeval", Version: report.ToolVersion},
		ConfigHash: cfgHash, Repository: repo, Weights: a.Weights, Thresholds: a.Thresholds, Notes: notes()}
}

func roundUnit(u Unit) Unit {
	u.Debt = u.Debt.rounded()
	u.Components.RedundantTokens = r2(u.Components.RedundantTokens)
	return u
}

// NewSnapshotReport builds the absolute report. top limits the callable and
// clone lists (0 = all).
func NewSnapshotReport(repo, cfgHash string, rev report.Revision, a *Analysis, top int) *SnapshotReport {
	r := &SnapshotReport{Header: header(SnapshotSchema, repo, cfgHash, a), Revision: rev, Coverage: a.Coverage,
		Codebase: roundUnit(a.Codebase), Violations: a.Violations, Cycles: a.Cycles}
	if r.Violations == nil {
		r.Violations = []Violation{}
	}
	if r.Cycles == nil {
		r.Cycles = [][]string{}
	}
	for _, m := range a.Modules {
		r.Modules = append(r.Modules, roundUnit(m))
	}
	for _, l := range a.Languages {
		r.Languages = append(r.Languages, roundUnit(l))
	}
	r.Callables = limit(a.Callables, top)
	r.Clones = limit(a.Clones, top)
	return r
}

func limit[T any](s []T, n int) []T {
	if s == nil {
		return []T{}
	}
	if n > 0 && len(s) > n {
		return s[:n]
	}
	return s
}

func ledger(before, after Unit) Ledger {
	l := Ledger{Components: after.Components.Sub(before.Components), Debt: after.Debt.Sub(before.Debt).rounded()}
	l.Components.RedundantTokens = r2(l.Components.RedundantTokens)
	if before.DebtPerKLOC != nil && after.DebtPerKLOC != nil {
		l.DebtPerKLOC = fp(r2(*after.DebtPerKLOC - *before.DebtPerKLOC))
	}
	if before.Score != nil && after.Score != nil {
		l.Score = fp(r2(*after.Score - *before.Score))
	}
	return l
}

func size(cs *change.Set) ChangeSize {
	s := ChangeSize{UnsupportedChangedFiles: map[string]int{}}
	for _, f := range cs.Files {
		if f.Class != model.ClassSource || f.Language == "" {
			continue
		}
		s.ChangedSourceFiles++
		if f.Supported {
			s.ChangedSupportedFiles++
			s.ChangedSupportedLines += f.Added + f.Deleted
		} else {
			s.UnsupportedChangedFiles[f.Language]++
		}
	}
	return s
}

func rate(d Debt, s ChangeSize) *Debt {
	if s.ChangedSupportedLines == 0 {
		return nil
	}
	r := d.Scale(1000 / float64(s.ChangedSupportedLines)).rounded()
	return &r
}

func isSource(snap *index.Snapshot, s *model.Symbol) bool {
	f := snap.Files[s.File]
	return s.Callable() && f != nil && f.Class == model.ClassSource && f.Supported()
}

func changed(base, target *index.Snapshot, cs *change.Set) ChangedCallables {
	var c ChangedCallables
	before, after, nb, na := 0, 0, 0, 0
	for _, sc := range cs.Symbols {
		b := sc.Base != nil && isSource(base, sc.Base)
		t := sc.Target != nil && isSource(target, sc.Target)
		switch {
		case b && t:
			if sc.Status == change.SymFormatting {
				continue
			}
			c.Modified++
			c.ModifiedDecisions[0] += sc.Base.Metrics.Decisions
			c.ModifiedDecisions[1] += sc.Target.Metrics.Decisions
			c.ModifiedCyclomatic[0] += sc.Base.Metrics.Cyclomatic
			c.ModifiedCyclomatic[1] += sc.Target.Metrics.Cyclomatic
			before += sc.Base.Metrics.Decisions
			after += sc.Target.Metrics.Decisions
			nb++
			na++
		case t:
			c.Added++
			c.AddedDecisions += sc.Target.Metrics.Decisions
			after += sc.Target.Metrics.Decisions
			na++
		case b:
			c.Deleted++
			c.DeletedDecisions += sc.Base.Metrics.Decisions
			before += sc.Base.Metrics.Decisions
			nb++
		}
	}
	c.NetDecisionPoints = after - before
	c.NetBasePoints = na - nb
	return c
}

func gates(l Ledger, c ChangedCallables) []Gate {
	g := func(name string, ok bool, detail string) Gate { return Gate{name, ok, detail} }
	return []Gate{
		g("no_new_exact_duplicates", l.Components.ExactRedundant <= 0, itoa(l.Components.ExactRedundant)+" exact redundant callables"),
		g("no_new_structural_duplicates", l.Components.StructRedundant <= 0, itoa(l.Components.StructRedundant)+" structural redundant callables"),
		g("no_new_rule_violations", l.Components.RuleViolations <= 0, itoa(l.Components.RuleViolations)+" rule-violating edges"),
		g("no_new_cycles", l.Components.CycleUnits <= 0, itoa(l.Components.CycleUnits)+" units in cycles"),
		g("no_branch_growth_in_existing_callables", c.ModifiedDecisions[1] <= c.ModifiedDecisions[0],
			itoa(c.ModifiedDecisions[1]-c.ModifiedDecisions[0])+" decision points in modified pre-existing callables"),
		g("no_hotspot_excess_growth", l.Components.CyclomaticExcess+l.Components.NestingExcess <= 0,
			itoa(l.Components.CyclomaticExcess)+" cyclomatic excess, "+itoa(l.Components.NestingExcess)+" nesting excess"),
	}
}

func itoa(v int) string {
	s := ""
	if v > 0 {
		s = "+"
	}
	return s + fmtInt(v)
}

func fmtInt(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var b []byte
	for v > 0 {
		b = append([]byte{byte('0' + v%10)}, b...)
		v /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

func moduleDeltas(before, after *Analysis) []ModuleDelta {
	type key struct{ lang, id string }
	bm, am := map[key]Unit{}, map[key]Unit{}
	for _, u := range before.Modules {
		bm[key{u.Language, u.ID}] = u
	}
	for _, u := range after.Modules {
		am[key{u.Language, u.ID}] = u
	}
	seen := map[key]bool{}
	var out []ModuleDelta
	add := func(k key) {
		if seen[k] {
			return
		}
		seen[k] = true
		b, a := bm[k], am[k]
		d := a.Debt.Sub(b.Debt).rounded()
		if d.Total == 0 && d.Complexity == 0 && d.Duplication == 0 && d.Architecture == 0 && d.Nesting == 0 {
			return
		}
		out = append(out, ModuleDelta{ID: k.id, Language: k.lang, DebtBefore: r2(b.Debt.Total), DebtAfter: r2(a.Debt.Total), DebtDelta: d, ScoreBefore: b.Score, ScoreAfter: a.Score})
	}
	for k := range bm {
		add(k)
	}
	for k := range am {
		add(k)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].DebtDelta.Total != out[j].DebtDelta.Total {
			return out[i].DebtDelta.Total > out[j].DebtDelta.Total
		}
		if out[i].Language != out[j].Language {
			return out[i].Language < out[j].Language
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// NewDeltaReport builds the comparison report.
func NewDeltaReport(repo, cfgHash string, baseRev, targetRev report.Revision, base, target *index.Snapshot, before, after *Analysis, cs *change.Set) *DeltaReport {
	r := &DeltaReport{Header: header(DeltaSchema, repo, cfgHash, after), Base: baseRev, Target: targetRev, CoverageAfter: after.Coverage,
		Size: size(cs), Before: roundUnit(before.Codebase), After: roundUnit(after.Codebase)}
	r.Ledger = ledger(before.Codebase, after.Codebase)
	r.Languages = languageLedgers(before, after)
	r.Rate = rate(after.Codebase.Debt.Sub(before.Codebase.Debt), r.Size)
	r.Changed = changed(base, target, cs)
	r.Gates = gates(r.Ledger, r.Changed)
	r.Modules = moduleDeltas(before, after)
	if r.Modules == nil {
		r.Modules = []ModuleDelta{}
	}
	return r
}

// NewStep builds one step's ledger.
func NewStep(i int, parent, commit, subject string, base, target *index.Snapshot, before, after *Analysis, cs *change.Set) Step {
	s := Step{Index: i, Parent: parent, Commit: commit, Subject: subject, Size: size(cs), Ledger: ledger(before.Codebase, after.Codebase), Languages: languageLedgers(before, after)}
	s.Rate = rate(after.Codebase.Debt.Sub(before.Codebase.Debt), s.Size)
	s.Failed = []string{}
	for _, g := range gates(s.Ledger, changed(base, target, cs)) {
		if !g.Passed {
			s.Failed = append(s.Failed, g.Name)
		}
	}
	return s
}

// CheckSteps verifies that step debt deltas sum to the net delta.
func CheckSteps(r *DeltaReport, startsAtBase bool) {
	sum := 0.0
	for _, s := range r.Steps {
		sum += s.Ledger.Debt.Total
	}
	r.StepCheck = &StepCheck{ChainStartsAtBase: startsAtBase, SumOfStepDebt: r2(sum), NetDebt: r.Ledger.Debt.Total,
		Telescopes: math.Abs(sum-r.Ledger.Debt.Total) <= 0.005*float64(len(r.Steps)+1)+1e-9}
}
