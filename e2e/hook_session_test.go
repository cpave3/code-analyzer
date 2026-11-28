package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// hookRepo is a clone of a fixture repository with the Claude Code hook
// templates run against it the way Claude Code runs them.
type hookRepo struct {
	t     *testing.T
	dir   string
	env   []string
	hooks string
}

func newHookRepo(t *testing.T, scenario, rev string, env ...string) *hookRepo {
	t.Helper()
	h := &hookRepo{t: t, dir: t.TempDir()}
	h.hooks, _ = filepath.Abs(filepath.Join("..", "integrations", "claude-code", "hooks"))
	h.git("clone", "-q", repos[scenario], ".")
	h.git("checkout", "-q", rev)
	h.env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+h.dir, "CODEEVAL_BIN="+binary, "CODEEVAL_STOP_CHECKS=quick", "CODEEVAL_DEFECTS_FAIL_ON=none")
	if cfg := scenarios[scenario].Config; cfg != "" {
		abs, _ := filepath.Abs(cfg)
		h.env = append(h.env, "CODEEVAL_CONFIG="+abs)
	}
	for _, e := range env {
		if e != "" {
			h.env = append(h.env, e)
		}
	}
	return h
}

func (h *hookRepo) git(args ...string) string {
	h.t.Helper()
	c := exec.Command("git", append([]string{"-c", "core.hooksPath=/dev/null", "-c", "user.name=t", "-c", "user.email=t@example.com"}, args...)...)
	c.Dir = h.dir
	out, err := c.CombinedOutput()
	if err != nil {
		h.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func (h *hookRepo) run(name, input string, env ...string) run {
	h.t.Helper()
	c := exec.Command("bash", filepath.Join(h.hooks, name))
	c.Dir = h.dir
	c.Stdin = strings.NewReader(input)
	c.Env = append(append([]string{}, h.env...), env...)
	var r run
	out, err := c.Output()
	r.stdout = string(out)
	if ee, ok := err.(*exec.ExitError); ok {
		r.code, r.stderr = ee.ExitCode(), string(ee.Stderr)
	} else if err != nil {
		h.t.Fatal(err)
	}
	return r
}

func (h *hookRepo) start(session, source string) run {
	return h.run("codeeval-session-base.sh", `{"session_id":"`+session+`","source":"`+source+`"}`)
}

func (h *hookRepo) stop(session string, env ...string) run {
	return h.run("codeeval-stop-analyze.sh", `{"session_id":"`+session+`","stop_hook_active":false}`, env...)
}

func (h *hookRepo) state(session, file string) string {
	b, _ := os.ReadFile(filepath.Join(h.dir, ".git", "codeeval", session, file))
	return strings.TrimSpace(string(b))
}

// The agent's change: the rename fixture's wrapper variant, which adds one
// debt point to a renamed callable.
func (h *hookRepo) addOnePoint(commit bool) {
	h.git("checkout", "-q", "wrapper", "--", ".")
	if commit {
		h.git("commit", "-q", "-am", "wrapper")
	}
}

// CE-3: two sessions in one checkout keep their own base, block budget and
// shown findings.
func TestStopHookKeepsEachSessionsState(t *testing.T) {
	h := newHookRepo(t, "rename", "base", "CODEEVAL_DEBT_BUDGET=0")
	x := h.git("rev-parse", "HEAD")
	h.start("A", "startup")
	h.addOnePoint(true)
	if r := h.stop("A"); r.code != 2 {
		t.Fatalf("A's first stop should block on the added point: %d\n%s", r.code, r.stderr)
	}
	h.start("B", "startup")
	h.start("A", "resume")
	if got := h.state("A", "session-base"); got != x {
		t.Errorf("A's base after B started and A resumed = %s, want %s", got, x)
	}
	if r := h.stop("A"); r.code != 0 {
		t.Errorf("B's start must not clear A's shown findings: A blocked again\n%s", r.stderr)
	}
	if got := h.state("B", "session-base"); got != h.git("rev-parse", "HEAD") {
		t.Errorf("B's base = %q, want HEAD", got)
	}
}

// A stop that passes moves the session's base to HEAD, so a later turn is
// not blocked on committed work an earlier turn's check already passed.
func TestStopHookAdvancesTheBaseOnAPass(t *testing.T) {
	h := newHookRepo(t, "rename", "base")
	h.start("A", "startup")
	h.addOnePoint(true)
	if r := h.stop("A", "CODEEVAL_DEBT_BUDGET=5"); r.code != 0 {
		t.Fatalf("+1 within a budget of 5 should pass: %d\n%s", r.code, r.stderr)
	}
	if got, head := h.state("A", "checked-base"), h.git("rev-parse", "HEAD"); got != head {
		t.Fatalf("checked-base = %q, want HEAD %s", got, head)
	}
	// A later turn, with no budget: the committed point was already checked.
	if r := h.stop("A", "CODEEVAL_DEBT_BUDGET=0"); r.code != 0 {
		t.Errorf("a later turn must not be blocked on work an earlier pass covered: %d\n%s", r.code, r.stderr)
	}
	// A new session start forgets it.
	h.start("A", "startup")
	if got := h.state("A", "checked-base"); got != "" {
		t.Errorf("a startup should clear checked-base, got %s", got)
	}
}

// A stop that blocks keeps the base, so the next stop still sees the change.
func TestStopHookKeepsTheBaseOnABlock(t *testing.T) {
	h := newHookRepo(t, "rename", "base", "CODEEVAL_DEBT_BUDGET=0")
	h.start("A", "startup")
	h.addOnePoint(false)
	if r := h.stop("A"); r.code != 2 {
		t.Fatalf("should block: %d", r.code)
	}
	if got := h.state("A", "checked-base"); got != "" {
		t.Errorf("a blocked stop must not move the base, got %s", got)
	}
}

// The hook's default budget lets a change add a little debt.
func TestStopHookDefaultBudget(t *testing.T) {
	h := newHookRepo(t, "rename", "base")
	h.start("A", "startup")
	h.addOnePoint(false)
	if r := h.stop("A"); r.code != 0 {
		t.Errorf("+1 is within the default budget: %d\n%s", r.code, r.stderr)
	}
}

// Once the branch's pull request is approved, a finding is reported to the
// user instead of sending the agent back to work.
func TestStopHookIsAdvisoryOnAnApprovedPR(t *testing.T) {
	fake := t.TempDir()
	gh := "#!/bin/sh\n[ \"$1 $2\" = \"pr view\" ] && echo \"$FAKE_DECISION\"\n"
	if err := os.WriteFile(filepath.Join(fake, "gh"), []byte(gh), 0o755); err != nil {
		t.Fatal(err)
	}
	ghEnv := "CODEEVAL_GH=" + filepath.Join(fake, "gh")
	for _, tc := range []struct {
		decision, mode string
		want           int
	}{
		{"APPROVED", "", 0},
		{"REVIEW_REQUIRED", "", 2},
		{"APPROVED", "CODEEVAL_HOOK_WHEN_APPROVED=block", 2},
	} {
		h := newHookRepo(t, "rename", "base", "CODEEVAL_DEBT_BUDGET=0", ghEnv, "FAKE_DECISION="+tc.decision, tc.mode)
		h.start("A", "startup")
		h.addOnePoint(false)
		r := h.stop("A")
		if r.code != tc.want {
			t.Errorf("%s %s: exit %d, want %d\n%s", tc.decision, tc.mode, r.code, tc.want, r.stderr)
			continue
		}
		if tc.want == 0 && (!strings.Contains(r.stdout, `"systemMessage"`) || !strings.Contains(r.stdout, "approved")) {
			t.Errorf("an advisory stop should tell the user: %q", r.stdout)
		}
		if tc.want == 0 && h.state("A", "hook-state.blocks") != "" {
			t.Error("an advisory stop is not a block")
		}
	}
}

// SessionStart tells the agent how to run the same check before it pushes.
func TestSessionStartGivesThePrePushCheck(t *testing.T) {
	h := newHookRepo(t, "rename", "base")
	r := h.start("A", "startup")
	for _, want := range []string{binary + "' analyze", "--working-tree", "--checks quick", "--debt-budget 10", "git merge-base HEAD"} {
		if !strings.Contains(r.stdout, want) {
			t.Errorf("SessionStart output lacks %q:\n%s", want, r.stdout)
		}
	}
}
