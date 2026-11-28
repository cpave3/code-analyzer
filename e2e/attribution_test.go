package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/attribution"
)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=Fixture", "-c", "user.email=fixture@example.com", "-c", "commit.gpgsign=false", "-c", "core.hooksPath=/dev/null"}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2026-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2026-01-01T00:00:00Z")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

const simple = "package p\n\nfunc Parse(s string) int { return len(s) }\n"

const complexParse = `package p

func Parse(s string) int {
	n := 0
	for _, r := range s {
		if r == 'a' { n++ } else if r == 'b' { n += 2 } else if r == 'c' { n += 3 } else if r == 'd' { n += 4 }
		if r == 'e' || r == 'f' || r == 'g' || r == 'h' || r == 'i' { n-- }
	}
	return n
}
`

func commit(t *testing.T, dir, file, content, msg string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, file), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "-m", msg)
	return git(t, dir, "rev-parse", "HEAD")
}

func TestAttributionByProvenance(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q", "-b", "main")
	commit(t, dir, "go.mod", "module example.com/p\n\ngo 1.22\n", "init")
	base := commit(t, dir, "p.go", simple, "base")
	commit(t, dir, "p.go", complexParse, "Add parser\n\nAgent-Model: model-a\nAgent-Harness: harness-x\nAgent-Workflow: v1\n")
	commit(t, dir, "p.go", simple, "Simplify parser\n\nAgent-Model: model-a\nAgent-Harness: harness-x\nAgent-Workflow: v2\n")
	git(t, dir, "checkout", "-q", "-b", "feature")
	commit(t, dir, "q.go", "package p\n\nfunc Q() int { return 1 }\n", "Feature part 1\n\nAgent-Model: model-b\nAgent-Workflow: v2\n")
	noted := commit(t, dir, "r.go", "package p\n\nfunc R() int { return 2 }\n", "Feature part 2")
	// a bare "Agent:" (pre-0.7 format) is still read as the model
	git(t, dir, "notes", "--ref=agent-provenance", "add", "-m", "Agent: model-b\nAgent-Workflow: v2", noted)
	git(t, dir, "checkout", "-q", "main")
	git(t, dir, "merge", "-q", "--no-ff", "-m", "Merge feature", "feature")
	commit(t, dir, "p.go", simple+"\n// typo fixed\n", "Fix typo\n\nCo-authored-by: Jane Doe <jane@example.com>\n")

	r := cli(t, "attribution", "--repo", dir, "--base", base, "--target", "main", "--group-by", "model,workflow", "--defects", "--format", "json", "--no-cache")
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	var rep attribution.Report
	if err := json.Unmarshal([]byte(r.stdout), &rep); err != nil {
		t.Fatalf("decode: %v\n%s", err, r.stdout)
	}
	if rep.Schema != attribution.Schema || len(rep.Steps) != 4 {
		t.Fatalf("schema %s steps %d", rep.Schema, len(rep.Steps))
	}
	groups := map[string]attribution.Group{}
	for _, g := range rep.Groups {
		groups[g.Key] = g
	}
	v1, v2 := groups["model=model-a, workflow=v1"], groups["model=model-a, workflow=v2"]
	if v1.Steps != 1 || v1.DebtDelta.Total <= 0 || v2.Steps != 1 || v2.DebtDelta.Total >= 0 || v1.DebtDelta.Total != -v2.DebtDelta.Total {
		t.Errorf("workflow groups: v1 %+v v2 %+v", v1, v2)
	}
	if b := groups["model=model-b, workflow=v2"]; b.Steps != 1 || b.Commits != 3 {
		t.Errorf("merge step (merge commit + message trailer + note) not attributed to model-b: %+v\n%+v", b, rep.Groups)
	}
	if h := groups["model=(none), workflow=(none)"]; h.Steps != 1 {
		t.Errorf("human commit group: %+v", h)
	}
	if rep.Sources["message"] != 3 || rep.Sources["note"] != 1 || rep.Sources["none"] != 2 {
		t.Errorf("sources %+v", rep.Sources)
	}
	sum := 0.0
	for _, s := range rep.Steps {
		sum += s.DebtDelta.Total
		if s.Defects == nil {
			t.Errorf("step %d without defects", s.Index)
		}
	}
	if sum != rep.NetDebt || rep.NetDebt != 0 {
		t.Errorf("step debt sum %v net %v, want 0 (complexity added then removed)", sum, rep.NetDebt)
	}
	if bad := cli(t, "attribution", "--repo", dir, "--base", base, "--target", "main", "--path", "p.go"); bad.code != 2 {
		t.Errorf("--path accepted: exit %d", bad.code)
	}
}
