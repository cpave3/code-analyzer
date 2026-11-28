package e2e

import (
	"strings"
	"testing"
)

// --debt-budget lets a change add up to that many debt points before
// net-debt or new-debt fails; the rename fixture's wrapper adds exactly 1.
func TestDebtBudget(t *testing.T) {
	for _, tc := range []struct {
		cmd, failOn, budget string
		want                int
	}{
		{"check", "net-debt", "", 3},
		{"check", "net-debt", "0.5", 3},
		{"check", "net-debt", "1", 0},
		{"check", "new-debt", "1", 0},
		{"check", "new-debt", "0.5", 3},
		{"check", "gates", "100", 3}, // a budget never excuses a failed gate
		{"analyze", "net-debt", "1", 0},
		{"analyze", "net-debt", "0", 3},
	} {
		args := []string{tc.cmd, "--repo", repos["rename"], "--base", "base", "--target", "wrapper", "--no-cache", "--format", "md"}
		if tc.cmd == "check" {
			args = append(args, "--fail-on", tc.failOn)
		} else {
			args = append(args, "--debt-fail-on", tc.failOn, "--checks", "quick")
		}
		if tc.budget != "" {
			args = append(args, "--debt-budget", tc.budget)
		}
		r := cli(t, args...)
		if r.code != tc.want {
			t.Errorf("%s --fail-on %s --debt-budget %q: exit %d, want %d\n%s%s", tc.cmd, tc.failOn, tc.budget, r.code, tc.want, r.stdout, r.stderr)
		}
		if tc.cmd == "check" && tc.failOn != "gates" && tc.budget != "" && !strings.Contains(r.stdout, "budget "+tc.budget+" point") {
			t.Errorf("%s: the report should state the budget:\n%s", tc.budget, r.stdout)
		}
	}
	if r := cli(t, "check", "--repo", repos["rename"], "--base", "base", "--target", "wrapper", "--debt-budget", "-1"); r.code != 2 {
		t.Errorf("a negative budget should be refused: exit %d", r.code)
	}
}
