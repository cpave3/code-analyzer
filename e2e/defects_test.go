package e2e

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
)

func defects(t *testing.T, wantCode int, args ...string) *rules.Report {
	t.Helper()
	full := append([]string{"defects", "--repo", repos["defects"], "--config", scenarios["defects"].Config, "--format", "json", "--no-cache"}, args...)
	r := cli(t, full...)
	if r.code != wantCode {
		t.Fatalf("exit %d (want %d): %s\n%s", r.code, wantCode, r.stderr, r.stdout)
	}
	var rep rules.Report
	if err := json.Unmarshal([]byte(r.stdout), &rep); err != nil {
		t.Fatalf("decode: %v\n%s", err, r.stdout)
	}
	if rep.Schema != rules.Schema || rep.RulesVersion != rules.Version {
		t.Fatalf("schema %s %s", rep.Schema, rep.RulesVersion)
	}
	return &rep
}

func TestDefectsBaseIsClean(t *testing.T) {
	rep := defects(t, 0, "--rev", "base")
	if len(rep.Findings) != 0 {
		t.Fatalf("clean base has findings: %+v", rep.Findings)
	}
	for _, lf := range rep.Files {
		if lf.Checked == 0 {
			t.Errorf("%s: no files checked", lf.Language)
		}
	}
}

// Every rule × language pair must detect its defect in the bad variant.
func TestDefectsEveryRuleAndLanguageDetects(t *testing.T) {
	rep := defects(t, 0, "--rev", "bad")
	seen := map[string]bool{}
	for _, f := range rep.Findings {
		seen[f.Rule+"|"+f.Language] = true
	}
	for _, rule := range catalog.Rules {
		for _, l := range []string{"go", "typescript"} {
			if !seen[rule.ID+"|"+l] {
				t.Errorf("%s found nothing in %s", rule.ID, l)
			}
		}
	}
	for _, f := range rep.Findings {
		if f.Rule == catalog.UnusedExport || f.Rule == catalog.UnusedDependency || f.Rule == catalog.ErrorToSuccess {
			if f.Severity != rules.Advisory {
				t.Errorf("%s must be advisory: %+v", f.Rule, f)
			}
		}
	}
}

func TestDefectsChangeGate(t *testing.T) {
	bad := defects(t, 3, "--base", "base", "--target", "bad")
	if !bad.Failed || bad.Mode != "change" || bad.Totals.Gate == 0 || bad.Persisting != 0 {
		t.Fatalf("bad change: failed=%v mode=%s totals=%+v persisting=%d", bad.Failed, bad.Mode, bad.Totals, bad.Persisting)
	}
	for _, f := range bad.Findings {
		if f.Change != "introduced" {
			t.Errorf("finding without change marker: %+v", f)
		}
	}
	clean := defects(t, 0, "--base", "base", "--target", "clean")
	if clean.Failed || len(clean.Findings) != 0 {
		t.Fatalf("clean change introduced findings: %+v", clean.Findings)
	}
	// Reverting bad resolves everything and introduces nothing.
	back := defects(t, 0, "--base", "bad", "--target", "base")
	if len(back.Findings) != 0 || len(back.Resolved) != len(defects(t, 0, "--rev", "bad").Findings) {
		t.Fatalf("revert: introduced %d resolved %d", len(back.Findings), len(back.Resolved))
	}
}

func TestDefectsSelectionScopeAndFailOn(t *testing.T) {
	rep := defects(t, 0, "--rev", "bad", "--rule", "test-oracle", "--path", "web")
	if len(rep.Findings) != 2 {
		t.Fatalf("want the 2 TypeScript test-oracle findings, got %+v", rep.Findings)
	}
	for _, f := range rep.Findings {
		if !strings.HasPrefix(f.Rule, "test-oracle/") || !strings.HasPrefix(f.File, "web/") {
			t.Errorf("out of selection: %+v", f)
		}
	}
	defects(t, 3, "--rev", "bad", "--fail-on", "gate")
	defects(t, 0, "--base", "base", "--target", "bad", "--fail-on", "none")
	for _, args := range [][]string{
		{"--rev", "bad", "--base", "base"},
		{"--base", "base"},
		{"--rev", "bad", "--rule", "no-such-rule"},
		{"--rev", "bad", "--fail-on", "sometimes"},
	} {
		full := append([]string{"defects", "--repo", repos["defects"], "--no-cache"}, args...)
		if r := cli(t, full...); r.code != 2 {
			t.Errorf("%v: exit %d, want 2", args, r.code)
		}
	}
}

func TestDefectsMarkdownAndRulesCommand(t *testing.T) {
	r := cli(t, "defects", "--repo", repos["defects"], "--config", scenarios["defects"].Config, "--base", "base", "--target", "bad", "--no-cache")
	if r.code != 3 || !strings.Contains(r.stdout, "**FAIL** (fail-on gate)") || !strings.Contains(r.stdout, "## Introduced findings") {
		t.Fatalf("markdown:\n%s", r.stdout)
	}
	list := cli(t, "rules")
	if list.code != 0 || strings.Contains(list.stdout, "not implemented") {
		t.Fatalf("rules:\n%s", list.stdout)
	}
	for _, rule := range catalog.Rules {
		if !strings.Contains(list.stdout, rule.ID+"  ["+string(rule.Severity)+"]") {
			t.Errorf("rules output lacks %s", rule.ID)
		}
	}
}
