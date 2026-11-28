package e2e

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/combined"
)

func analyze(t *testing.T, wantCode int, args ...string) (map[string]any, string) {
	t.Helper()
	base := []string{"analyze", "--repo", repos["defects"], "--config", scenarios["defects"].Config, "--no-cache"}
	r := cli(t, append(append(base, args...), "--format", "json")...)
	if r.code != wantCode {
		t.Fatalf("exit %d (want %d): %s", r.code, wantCode, r.stderr)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(r.stdout), &m); err != nil {
		t.Fatal(err)
	}
	if m["schema"] != combined.Schema {
		t.Fatalf("schema %v", m["schema"])
	}
	md := cli(t, append(base, args...)...)
	return m, md.stdout
}

func TestAnalyzeCombinesDebtAndDefects(t *testing.T) {
	m, md := analyze(t, 0, "--rev", "bad")
	defects := m["defects"].(map[string]any)
	if g := defects["totals"].(map[string]any)["gate"].(float64); g == 0 {
		t.Error("no gate defects in bad")
	}
	if m["debt"].(map[string]any)["schema"] != "codeeval.explain/v2" {
		t.Errorf("debt section %v", m["debt"].(map[string]any)["schema"])
	}
	files := m["files"].([]any)
	if len(files) < 2 {
		t.Fatalf("files %v", files)
	}
	for i := 1; i < len(files); i++ {
		if files[i-1].(map[string]any)["gateDefects"].(float64) < files[i].(map[string]any)["gateDefects"].(float64) {
			t.Errorf("ranking not by gate defects at %d: %v", i, files)
		}
	}
	if cl, ok := m["clones"].(map[string]any); !ok || cl["version"] != "clones/v1" || cl["mode"] != "all" {
		t.Errorf("clones section: %v", m["clones"])
	}
	for _, want := range []string{"# Duplication beyond structural clones", "## Summary", "## Files to look at first", "# Debt (scoring/v2)", "# Defects (rules/v1)", "# Limits of the combined view"} {
		if !strings.Contains(md, want) {
			t.Errorf("markdown lacks %q", want)
		}
	}

	ch, cmd := analyze(t, 3, "--base", "base", "--target", "bad", "--top", "1")
	if cl := ch["clones"].(map[string]any); cl["mode"] != "introduced" {
		t.Errorf("change clones mode %v", cl["mode"])
	}
	if ch["mode"] != "change" || ch["failed"] != true || ch["factorsIntroducedOrWorsened"] == nil {
		t.Errorf("change: mode %v failed %v worse %v", ch["mode"], ch["failed"], ch["factorsIntroducedOrWorsened"])
	}
	if fp, _ := ch["failureFingerprint"].(string); len(fp) != 12 || !strings.Contains(cmd, "Failure fingerprint: `"+fp+"`") {
		t.Errorf("failure fingerprint %q missing from JSON or Markdown", fp)
	}
	again, _ := analyze(t, 3, "--base", "base", "--target", "bad", "--top", "5")
	if again["failureFingerprint"] != ch["failureFingerprint"] {
		t.Error("fingerprint depends on --top or is not deterministic")
	}
	if !strings.Contains(cmd, "FAIL: defects:") {
		t.Errorf("markdown status:\n%s", cmd[:200])
	}
	analyze(t, 0, "--base", "base", "--target", "clean")
	if r := cli(t, "analyze", "--repo", repos["defects"], "--rev", "bad", "--debt-fail-on", "new-debt"); r.code != 2 {
		t.Errorf("--debt-fail-on in snapshot mode: exit %d", r.code)
	}
}
