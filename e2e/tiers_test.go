package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

type tierReport struct {
	Failed  bool   `json:"failed"`
	Checks  string `json:"checks"`
	Clones  any    `json:"clones"`
	Reasons []string
	Defects struct {
		Findings []struct{ Rule, File, Key, Severity string } `json:"findings"`
		Coverage []struct{ Rule, Language, Status string }    `json:"coverage"`
	} `json:"defects"`
}

func tier(t *testing.T, checks string, wantCode int) (tierReport, []string) {
	t.Helper()
	r := cli(t, "analyze", "--repo", repos["defects"], "--config", scenarios["defects"].Config, "--no-cache",
		"--base", "base", "--target", "bad", "--checks", checks, "--top", "0", "--format", "json")
	if r.code != wantCode {
		t.Fatalf("--checks %s: exit %d (want %d): %s", checks, r.code, wantCode, r.stderr)
	}
	var rep tierReport
	if err := json.Unmarshal([]byte(r.stdout), &rep); err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, f := range rep.Defects.Findings {
		ids = append(ids, f.Rule+"|"+f.File+"|"+f.Key)
	}
	sort.Strings(ids)
	return rep, ids
}

// The quick and broad tiers split the defect rules without losing or
// duplicating a finding: together they report exactly what --checks all does.
func TestAnalyzeTiersSplitTheRules(t *testing.T) {
	_, all := tier(t, "all", 3)
	quick, q := tier(t, "quick", 3)
	broad, b := tier(t, "broad", 3)
	got := append(append([]string{}, q...), b...)
	sort.Strings(got)
	if strings.Join(got, "\n") != strings.Join(all, "\n") {
		t.Errorf("quick + broad != all\nquick %v\nbroad %v\nall   %v", q, b, all)
	}
	if len(q) == 0 || len(b) == 0 {
		t.Errorf("fixture should have findings in both tiers: quick %d broad %d", len(q), len(b))
	}
	if quick.Checks != "quick" || quick.Clones != nil {
		t.Errorf("quick: checks %q, clones %v (want none)", quick.Checks, quick.Clones)
	}
	if broad.Clones == nil {
		t.Error("broad should include near-duplication")
	}
	deferred := func(r tierReport) int {
		n := 0
		for _, c := range r.Defects.Coverage {
			if c.Status == "deferred" {
				n++
			}
		}
		return n
	}
	if deferred(quick) == 0 || deferred(broad) == 0 {
		t.Errorf("each tier must list the other's rules as deferred: quick %d broad %d", deferred(quick), deferred(broad))
	}
	for _, reason := range broad.Reasons {
		if strings.HasPrefix(reason, "debt") {
			t.Errorf("broad must not fail on debt: %v", broad.Reasons)
		}
	}
	if r := cli(t, "analyze", "--repo", repos["defects"], "--rev", "bad", "--checks", "quick"); r.code != 2 {
		t.Errorf("--checks quick in snapshot mode: exit %d", r.code)
	}
}

// The Stop hook runs the quick tier in the foreground and the broad tier in
// the background; the broad result is reported at the next stop, once.
func TestStopHookTiers(t *testing.T) {
	dir := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		c := exec.Command("git", append([]string{"-c", "core.hooksPath=/dev/null"}, args...)...)
		c.Dir = dir
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git("clone", "-q", repos["defects"], ".")
	git("checkout", "-q", "base")
	git("checkout", "-q", "bad", "--", ".") // the agent's uncommitted work
	hooks, _ := filepath.Abs(filepath.Join("..", "integrations", "claude-code", "hooks"))
	cfg, _ := filepath.Abs(scenarios["defects"].Config)
	hook := func(name, input string) run {
		t.Helper()
		c := exec.Command("bash", filepath.Join(hooks, name))
		c.Dir = dir
		c.Stdin = strings.NewReader(input)
		c.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+dir, "CODEEVAL_BIN="+binary, "CODEEVAL_CONFIG="+cfg, "CODEEVAL_DEBT_FAIL_ON=none")
		var r run
		out, err := c.Output()
		r.stdout = string(out)
		if ee, ok := err.(*exec.ExitError); ok {
			r.code, r.stderr = ee.ExitCode(), string(ee.Stderr)
		}
		return r
	}
	hook("codeeval-session-base.sh", `{"source":"startup"}`)
	broad := filepath.Join(dir, ".git", "codeeval", "default", "broad")

	first := hook("codeeval-stop-analyze.sh", `{}`)
	if first.code != 2 || !strings.Contains(first.stderr, "Quick check") || strings.Contains(first.stderr, "Background check") {
		t.Fatalf("first stop should block on the quick tier only: %d\n%s", first.code, first.stderr)
	}
	deadline := time.Now().Add(60 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(broad, "exit")); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("background check did not finish")
		}
		time.Sleep(100 * time.Millisecond)
	}
	second := hook("codeeval-stop-analyze.sh", `{}`)
	if second.code != 2 || !strings.Contains(second.stderr, "Background check") || !strings.Contains(second.stderr, "Broad check") || strings.Contains(second.stderr, "Quick check") {
		t.Fatalf("second stop should report the background result only (the quick one was shown): %d\n%s", second.code, second.stderr)
	}
	if r := hook("codeeval-stop-analyze.sh", `{"stop_hook_active":true}`); r.code != 0 {
		t.Errorf("stop_hook_active must pass: %d", r.code)
	}
	// The second stop started another background run; a new session stops
	// that codeeval process and clears the state.
	time.Sleep(50 * time.Millisecond) // let the wrapper record the pid
	pidb, _ := os.ReadFile(filepath.Join(broad, "pid"))
	hook("codeeval-session-base.sh", `{"source":"startup"}`)
	if pid := strings.TrimSpace(string(pidb)); pid != "" {
		time.Sleep(200 * time.Millisecond)
		if exec.Command("kill", "-0", pid).Run() == nil {
			t.Errorf("background codeeval %s still running after a new session", pid)
		}
	}
	if _, err := os.Stat(broad); err == nil {
		t.Error("a new session should clear the background state")
	}
}
