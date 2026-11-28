// Package attribution groups per-step quality deltas by commit provenance
// (agent, harness, workflow, ...). It is observational: steps differ in
// task, size and difficulty, so differences between groups are not effects
// of the agent or workflow.
package attribution

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/provenance"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// Schema of the attribution report.
const Schema = "codeeval.attribution/v1"

// Defects summarises defect findings a step introduced and resolved.
type Defects struct {
	IntroducedGate     int `json:"introducedGate"`
	IntroducedAdvisory int `json:"introducedAdvisory"`
	Resolved           int `json:"resolved"`
}

func (d *Defects) add(o *Defects) {
	d.IntroducedGate += o.IntroducedGate
	d.IntroducedAdvisory += o.IntroducedAdvisory
	d.Resolved += o.Resolved
}

// Step is one first-parent step with its provenance.
type Step struct {
	Index     int                 `json:"index"`
	Parent    string              `json:"parent"`
	Commit    string              `json:"commit"`
	Subject   string              `json:"subject"`
	Group     string              `json:"group"`
	Commits   []provenance.Record `json:"commits"` // commits the step brings in (a merge: the merged branch too)
	Size      scorev2.ChangeSize  `json:"size"`
	DebtDelta scorev2.Debt        `json:"debtDelta"`
	Decisions int                 `json:"decisionPointsDelta"`
	Defects   *Defects            `json:"defects,omitempty"`
}

// Group aggregates the steps sharing a provenance key.
type Group struct {
	Key          string       `json:"key"`
	Steps        int          `json:"steps"`
	Commits      int          `json:"commits"`
	ChangedLines int          `json:"changedSupportedSourceLines"`
	DebtDelta    scorev2.Debt `json:"debtDelta"`
	DebtPerKLOC  *float64     `json:"debtDeltaPerChangedKLOC"`
	Decisions    int          `json:"decisionPointsDelta"`
	Defects      *Defects     `json:"defects,omitempty"`
	GatePerKLOC  *float64     `json:"introducedGatePerChangedKLOC,omitempty"`
	StepsDebtUp  int          `json:"stepsWithDebtIncrease"`
	StepsDebtDn  int          `json:"stepsWithDebtDecrease"`
}

// Report is the output of `codeeval attribution`.
type Report struct {
	Schema       string          `json:"schema"`
	Tool         report.Tool     `json:"tool"`
	Scoring      string          `json:"scoringVersion"`
	RulesVersion string          `json:"rulesVersion,omitempty"`
	ConfigHash   string          `json:"configHash"`
	Repository   string          `json:"repository"`
	Base         report.Revision `json:"base"`
	Target       report.Revision `json:"target"`
	GroupBy      []string        `json:"groupBy"`
	Inferred     bool            `json:"coAuthorInference"`
	Sources      map[string]int  `json:"commitsBySource"`
	Groups       []Group         `json:"groups"`
	Steps        []Step          `json:"steps"`
	NetDebt      float64         `json:"netDebtDeltaOfSteps"`
	Notes        []string        `json:"notes"`
}

func r2(v float64) float64 { return math.Round(v*100) / 100 }

func addDebt(a, b scorev2.Debt) scorev2.Debt {
	return scorev2.Debt{Complexity: a.Complexity + b.Complexity, Nesting: a.Nesting + b.Nesting, Duplication: a.Duplication + b.Duplication,
		Architecture: a.Architecture + b.Architecture, Total: a.Total + b.Total}
}

// GroupKey is the key of a step: its commits' shared key, or Mixed. Merge
// commits without provenance are bookkeeping and do not take part, unless
// the step has nothing else.
func GroupKey(commits []provenance.Record, groupBy []string) string {
	key, n := "", 0
	for _, c := range commits {
		if c.Merge && c.Source == provenance.SourceNone {
			continue
		}
		k := c.Key(groupBy)
		if n == 0 {
			key = k
		} else if k != key {
			return provenance.Mixed
		}
		n++
	}
	if n == 0 && len(commits) > 0 {
		return commits[0].Key(groupBy)
	}
	return key
}

// Finish computes groups and totals from steps.
func (r *Report) Finish() {
	byKey := map[string]*Group{}
	r.Sources = map[string]int{}
	seen := map[string]bool{}
	for _, s := range r.Steps {
		g := byKey[s.Group]
		if g == nil {
			g = &Group{Key: s.Group}
			byKey[s.Group] = g
		}
		g.Steps++
		g.Commits += len(s.Commits)
		g.ChangedLines += s.Size.ChangedSupportedLines
		g.DebtDelta = addDebt(g.DebtDelta, s.DebtDelta)
		g.Decisions += s.Decisions
		switch {
		case s.DebtDelta.Total > 0.005:
			g.StepsDebtUp++
		case s.DebtDelta.Total < -0.005:
			g.StepsDebtDn++
		}
		if s.Defects != nil {
			if g.Defects == nil {
				g.Defects = &Defects{}
			}
			g.Defects.add(s.Defects)
		}
		r.NetDebt += s.DebtDelta.Total
		for _, c := range s.Commits {
			if !seen[c.Commit] {
				seen[c.Commit] = true
				r.Sources[c.Source]++
			}
		}
	}
	r.NetDebt = r2(r.NetDebt)
	for _, g := range byKey {
		g.DebtDelta = scorev2.Debt{Complexity: r2(g.DebtDelta.Complexity), Nesting: r2(g.DebtDelta.Nesting), Duplication: r2(g.DebtDelta.Duplication),
			Architecture: r2(g.DebtDelta.Architecture), Total: r2(g.DebtDelta.Total)}
		if g.ChangedLines > 0 {
			v := r2(g.DebtDelta.Total * 1000 / float64(g.ChangedLines))
			g.DebtPerKLOC = &v
			if g.Defects != nil {
				d := r2(float64(g.Defects.IntroducedGate) * 1000 / float64(g.ChangedLines))
				g.GatePerKLOC = &d
			}
		}
		r.Groups = append(r.Groups, *g)
	}
	sort.Slice(r.Groups, func(i, j int) bool {
		if r.Groups[i].Steps != r.Groups[j].Steps {
			return r.Groups[i].Steps > r.Groups[j].Steps
		}
		return r.Groups[i].Key < r.Groups[j].Key
	})
	r.Notes = []string{
		"observational: groups differ in tasks, size, difficulty and timing; a difference between groups is not an effect of the agent, harness or workflow",
		"provenance is self-reported text in commit messages and notes (or inferred from Co-Authored-By lines when enabled); it is not verified",
		"steps are first-parent commits; a merge step includes every commit it brings in (a merge commit without provenance does not count), and a step whose commits disagree is grouped as (mixed)",
		"debt is scoring/v2 (uncalibrated structural proxy); per-step deltas sum exactly to the net change of the steps",
		"small groups are anecdotes: compare groups only with enough steps and similar kinds of work",
	}
}

// WriteMarkdown renders the report.
func WriteMarkdown(w io.Writer, r *Report) {
	fmt.Fprintf(w, "# Attribution — %s\n\n", r.Repository)
	fmt.Fprintf(w, "`%s` (%.12s) → `%s` (%.12s) · grouped by %s · %s", r.Base.Ref, r.Base.Commit, r.Target.Ref, r.Target.Commit, strings.Join(r.GroupBy, ", "), r.Scoring)
	if r.RulesVersion != "" {
		fmt.Fprintf(w, " + %s", r.RulesVersion)
	}
	fmt.Fprintf(w, " · config %s\n\n", r.ConfigHash)
	var src []string
	for k, v := range r.Sources {
		src = append(src, fmt.Sprintf("%s %d", k, v))
	}
	sort.Strings(src)
	fmt.Fprintf(w, "%d steps, net debt Δ %+.2f. Commit provenance sources: %s.\n\n", len(r.Steps), r.NetDebt, strings.Join(src, ", "))
	defects := len(r.Groups) > 0 && r.Groups[0].Defects != nil
	head := "| Group | Steps | Commits | Changed lines | Debt Δ | Debt Δ / changed kLOC | Steps debt ↑/↓ | Decision points Δ |"
	sep := "|---|---|---|---|---|---|---|---|"
	if defects {
		head += " Gate defects introduced | Advisory introduced | Resolved | Gate / changed kLOC |"
		sep += "---|---|---|---|"
	}
	fmt.Fprintln(w, head)
	fmt.Fprintln(w, sep)
	for _, g := range r.Groups {
		per := "–"
		if g.DebtPerKLOC != nil {
			per = fmt.Sprintf("%+.2f", *g.DebtPerKLOC)
		}
		fmt.Fprintf(w, "| %s | %d | %d | %d | %+.2f | %s | %d/%d | %+d |", g.Key, g.Steps, g.Commits, g.ChangedLines, g.DebtDelta.Total, per, g.StepsDebtUp, g.StepsDebtDn, g.Decisions)
		if defects && g.Defects != nil {
			gp := "–"
			if g.GatePerKLOC != nil {
				gp = fmt.Sprintf("%.2f", *g.GatePerKLOC)
			}
			fmt.Fprintf(w, " %d | %d | %d | %s |", g.Defects.IntroducedGate, g.Defects.IntroducedAdvisory, g.Defects.Resolved, gp)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Steps")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| # | Commit | Group | Debt Δ | Changed lines | Subject |")
	fmt.Fprintln(w, "|---|---|---|---|---|---|")
	for _, s := range r.Steps {
		subj := strings.ReplaceAll(strings.ReplaceAll(s.Subject, "|", "\\|"), "`", "'")
		if len(subj) > 70 {
			subj = subj[:67] + "..."
		}
		fmt.Fprintf(w, "| %d | %.10s | %s | %+.2f | %d | %s |\n", s.Index, s.Commit, s.Group, s.DebtDelta.Total, s.Size.ChangedSupportedLines, subj)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "## Limits")
	fmt.Fprintln(w)
	for _, n := range r.Notes {
		fmt.Fprintf(w, "- %s\n", n)
	}
}
