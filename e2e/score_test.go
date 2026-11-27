package e2e

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/scorev2"
)

func scoreArgs(scenario string, args ...string) []string {
	out := append([]string{"score", "--repo", repos[scenario], "--format", "json", "--no-cache"}, args...)
	if c := scenarios[scenario].Config; c != "" {
		out = append(out, "--config", c)
	}
	return out
}

func snapshotV2(t *testing.T, scenario, rev string) *scorev2.SnapshotReport {
	t.Helper()
	r := cli(t, scoreArgs(scenario, "--rev", rev, "--top", "0")...)
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	var s scorev2.SnapshotReport
	if err := json.Unmarshal([]byte(r.stdout), &s); err != nil {
		t.Fatalf("decode: %v\n%s", err, r.stdout)
	}
	return &s
}

func deltaV2(t *testing.T, repo, cfg, base, target string, extra ...string) *scorev2.DeltaReport {
	t.Helper()
	args := []string{"score", "--repo", repo, "--base", base, "--target", target, "--format", "json", "--no-cache"}
	if cfg != "" {
		args = append(args, "--config", cfg)
	}
	r := cli(t, append(args, extra...)...)
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	var d scorev2.DeltaReport
	if err := json.Unmarshal([]byte(r.stdout), &d); err != nil {
		t.Fatalf("decode: %v\n%s", err, r.stdout)
	}
	return &d
}

func near(a, b float64) bool { return math.Abs(a-b) < 0.011 }

// Module debt is additive: modules sum to the codebase.
func TestV2ModuleDebtSumsToCodebase(t *testing.T) {
	for _, sc := range []string{"reuse", "consolidate", "complexity", "layering"} {
		for _, v := range append([]string{"base"}, scenarios[sc].Variants...) {
			s := snapshotV2(t, sc, v)
			if s.Schema != scorev2.SnapshotSchema || s.Scoring != scorev2.Version || s.Revision.Commit != refs[sc][v] {
				t.Fatalf("%s/%s: header %+v %+v", sc, v, s.Header.Schema, s.Revision)
			}
			sum := 0.0
			for _, m := range s.Modules {
				sum += m.Debt.Total
			}
			if !near(sum, s.Codebase.Debt.Total) {
				t.Errorf("%s/%s: modules sum %.2f, codebase %.2f", sc, v, sum, s.Codebase.Debt.Total)
			}
			if s.Codebase.Score == nil || *s.Codebase.Score <= 0 || *s.Codebase.Score > 100 {
				t.Errorf("%s/%s: score %v out of range", sc, v, s.Codebase.Score)
			}
		}
	}
}

// The absolute measurement separates the paired fixture variants in the
// expected direction on the component each scenario exercises.
func TestV2SnapshotSeparatesFixturePairs(t *testing.T) {
	cases := []struct {
		scenario, component string
		get                 func(*scorev2.SnapshotReport) float64
	}{
		{"reuse", "duplication", func(s *scorev2.SnapshotReport) float64 { return s.Codebase.Debt.Duplication }},
		{"consolidate", "duplication", func(s *scorev2.SnapshotReport) float64 { return s.Codebase.Debt.Duplication }},
		{"complexity", "complexity+nesting", func(s *scorev2.SnapshotReport) float64 { return s.Codebase.Debt.Complexity + s.Codebase.Debt.Nesting }},
		{"layering", "architecture", func(s *scorev2.SnapshotReport) float64 { return s.Codebase.Debt.Architecture }},
	}
	for _, c := range cases {
		good, bad := c.get(snapshotV2(t, c.scenario, "good")), c.get(snapshotV2(t, c.scenario, "bad"))
		if !(bad > good) {
			t.Errorf("%s: %s debt bad %.2f should exceed good %.2f", c.scenario, c.component, bad, good)
		}
	}
	l := snapshotV2(t, "layering", "bad")
	if len(l.Violations) != 1 || l.Violations[0].From != "internal/store" || l.Violations[0].To != "internal/cli" {
		t.Errorf("layering bad should list the store → cli violation: %+v", l.Violations)
	}
}

// Deltas telescope: per-step ledgers along the first-parent chain sum to the
// net ledger, and the delta equals after − before of the two snapshots.
func TestV2DeltaTelescopesAcrossSteps(t *testing.T) {
	dir := t.TempDir()
	s := scenarios["complexity"]
	r, err := s.BuildRepo(dir)
	if err != nil {
		t.Fatal(err)
	}
	gitT(t, dir, "checkout", "-q", "-B", "main", "base")
	for _, v := range []string{"bad", "good", "bad"} {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if e.Name() != ".git" {
				os.RemoveAll(filepath.Join(dir, e.Name()))
			}
		}
		if err := s.Materialize(v, dir); err != nil {
			t.Fatal(err)
		}
		gitT(t, dir, "add", "-A")
		gitT(t, dir, "commit", "-qm", "apply "+v)
	}
	head := gitT(t, dir, "rev-parse", "HEAD")
	d := deltaV2(t, dir, "", r["base"], head, "--steps")
	if d.Schema != scorev2.DeltaSchema || len(d.Steps) != 3 || d.StepCheck == nil {
		t.Fatalf("want 3 steps and a step check: %+v", d.StepCheck)
	}
	if !d.StepCheck.Telescopes || !d.StepCheck.ChainStartsAtBase {
		t.Errorf("steps must telescope: %+v", d.StepCheck)
	}
	sum := 0.0
	for _, st := range d.Steps {
		sum += st.Ledger.Debt.Total
	}
	if !near(sum, d.Ledger.Debt.Total) {
		t.Errorf("Σ steps %.2f, net %.2f", sum, d.Ledger.Debt.Total)
	}
	if !near(d.After.Debt.Total-d.Before.Debt.Total, d.Ledger.Debt.Total) {
		t.Errorf("ledger is not after − before")
	}
	// bad → good → bad: the middle step improves, the others degrade.
	if !(d.Steps[0].Ledger.Debt.Total > 0 && d.Steps[1].Ledger.Debt.Total < 0 && d.Steps[2].Ledger.Debt.Total > 0) {
		t.Errorf("step directions: %.2f %.2f %.2f", d.Steps[0].Ledger.Debt.Total, d.Steps[1].Ledger.Debt.Total, d.Steps[2].Ledger.Debt.Total)
	}
	if !near(d.Steps[0].Ledger.Debt.Total, -d.Steps[1].Ledger.Debt.Total) {
		t.Errorf("applying and reverting the same change should cancel: %.2f vs %.2f", d.Steps[0].Ledger.Debt.Total, d.Steps[1].Ledger.Debt.Total)
	}
}

// The decision-point versus base-point split and the advisory gates.
func TestV2DeltaDecomposesChangedCallables(t *testing.T) {
	good := deltaV2(t, repos["complexity"], "", "base", "good")
	bad := deltaV2(t, repos["complexity"], "", "base", "bad")
	if good.Changed.NetBasePoints <= 0 {
		t.Errorf("extraction adds functions, so base points should rise: %+v", good.Changed)
	}
	if bad.Changed.ModifiedDecisions[1] <= bad.Changed.ModifiedDecisions[0] {
		t.Errorf("nesting variant should add decision points to the existing callable: %+v", bad.Changed)
	}
	gate := func(d *scorev2.DeltaReport, name string) bool {
		for _, g := range d.Gates {
			if g.Name == name {
				return g.Passed
			}
		}
		t.Fatalf("gate %s missing", name)
		return false
	}
	if gate(bad, "no_branch_growth_in_existing_callables") {
		t.Error("bad variant should fail the existing-callable branch gate")
	}
	dup := deltaV2(t, repos["reuse"], "", "base", "bad")
	if gate(dup, "no_new_structural_duplicates") || dup.Ledger.Components.StructRedundant != 1 {
		t.Errorf("reimplementation should add one structural duplicate: %+v", dup.Ledger.Components)
	}
	lay := deltaV2(t, repos["layering"], scenarios["layering"].Config, "base", "bad")
	if gate(lay, "no_new_rule_violations") || lay.Ledger.Debt.Architecture <= 0 {
		t.Errorf("layering violation should raise architecture debt: %+v", lay.Ledger.Debt)
	}
}

func TestV2DeterministicAndReadOnly(t *testing.T) {
	before := repoState(t, repos["consolidate"])
	a := cli(t, scoreArgs("consolidate", "--base", "base", "--target", "bad", "--steps")...)
	b := cli(t, scoreArgs("consolidate", "--base", "base", "--target", "bad", "--steps")...)
	if a.code != 0 || a.stdout != b.stdout {
		t.Fatalf("score output must be deterministic (exit %d)", a.code)
	}
	if after := repoState(t, repos["consolidate"]); after != before {
		t.Error("score must not modify the evaluated repository")
	}
	md := cli(t, "score", "--repo", repos["consolidate"], "--rev", "bad", "--no-cache").stdout
	for _, want := range []string{"scoring/v2-experimental", "## Codebase", "## Modules", "uncalibrated", "not comparable to scoring/v1"} {
		if !strings.Contains(md, want) {
			t.Errorf("snapshot markdown missing %q", want)
		}
	}
	if r := cli(t, "score", "--repo", repos["edge"], "--rev", "base", "--base", "base", "--no-cache"); r.code != 2 {
		t.Errorf("--rev with --base must be rejected, got %d", r.code)
	}
	if r := cli(t, "score", "--repo", repos["edge"], "--rev", "base", "--steps", "--no-cache"); r.code != 2 {
		t.Errorf("--steps without a range must be rejected, got %d", r.code)
	}
}
