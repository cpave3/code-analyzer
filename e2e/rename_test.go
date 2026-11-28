package e2e

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/scorev2"
)

// A callable renamed, or whose body moved to a new name, keeps its debt: the
// check reports only what the change added to it (the Legato PR #458 case).
func TestCheckFollowsRenamedCallables(t *testing.T) {
	check := func(target, failOn string) (scorev2.CheckReport, int) {
		t.Helper()
		r := cli(t, "check", "--repo", repos["rename"], "--base", "base", "--target", target, "--fail-on", failOn, "--format", "json", "--no-cache")
		var c scorev2.CheckReport
		if err := json.Unmarshal([]byte(r.stdout), &c); err != nil {
			t.Fatalf("%s: %v\n%s", target, err, r.stderr)
		}
		return c, r.code
	}
	ids := func(fs []scorev2.Factor) string {
		var out []string
		for _, f := range fs {
			out = append(out, f.Change+" "+f.ID)
		}
		return strings.Join(out, ", ")
	}
	for _, tc := range []struct{ target, name string }{{"wrapper", "Service.verifyPane"}, {"renamed", "Service.VerifyPane"}} {
		c, code := check(tc.target, "new-debt")
		var cc *scorev2.Factor
		for i, f := range c.Factors {
			if f.Change != "worsened" {
				t.Errorf("%s: %s should be worsened, not %s", tc.target, f.ID, f.Change)
			}
			if f.Kind == scorev2.KindComplexity {
				cc = &c.Factors[i]
			}
		}
		if cc == nil || cc.Location.Symbol != tc.name || cc.DebtDelta != 1 || cc.Metrics["baseCyclomatic"] != 15 {
			t.Errorf("%s: want %s worsened by 1 point from cyclomatic 15, got %s", tc.target, tc.name, ids(c.Factors))
		}
		if len(c.Improved) != 0 {
			t.Errorf("%s: the old callable's debt moved, it was not removed: %s", tc.target, ids(c.Improved))
		}
		if code != 3 {
			t.Errorf("%s: one added point is still new debt: exit %d", tc.target, code)
		}
	}
	// After the split the moved body is below the threshold: its debt is
	// reported as reduced under the new name, nothing as introduced.
	c, code := check("split", "new-debt")
	if len(c.Factors) != 0 || code != 0 {
		t.Errorf("split: nothing introduced or worsened, got %s (exit %d)", ids(c.Factors), code)
	}
	if got := ids(c.Improved); !strings.Contains(got, "removed cc:delivery/delivery.go:Service.VerifyPaneCommand") {
		t.Errorf("split: want the old hotspot reported as removed, got %s", got)
	}
}
