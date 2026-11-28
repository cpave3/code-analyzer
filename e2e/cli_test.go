// Package e2e exercises the real codeeval binary as a subprocess against
// fixture Git repositories built from testdata/fixtures.
package e2e

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/fixturerepo"
	"github.com/cameron/code-evaluator/internal/report"
)

var (
	binary    string
	fixtures  = filepath.Join("..", "testdata", "fixtures")
	repos     = map[string]string{} // scenario -> repo path
	refs      = map[string]map[string]string{}
	scenarios = map[string]fixturerepo.Scenario{}
	workDir   string
)

func TestMain(m *testing.M) {
	code, err := setup(m)
	if err != nil {
		fmt.Fprintln(os.Stderr, "setup:", err)
		code = 1
	}
	if workDir != "" {
		os.RemoveAll(workDir)
	}
	os.Exit(code)
}

func setup(m *testing.M) (int, error) {
	var err error
	workDir, err = os.MkdirTemp("", "codeeval-e2e-")
	if err != nil {
		return 1, err
	}
	binary = filepath.Join(workDir, "codeeval")
	build := exec.Command("go", "build", "-o", binary, "../cmd/codeeval")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return 1, fmt.Errorf("build: %w", err)
	}
	list, err := fixturerepo.Load(fixtures)
	if err != nil {
		return 1, err
	}
	for _, s := range list {
		dst := filepath.Join(workDir, "repos", s.Name)
		r, err := s.BuildRepo(dst)
		if err != nil {
			return 1, fmt.Errorf("%s: %w", s.Name, err)
		}
		repos[s.Name], refs[s.Name], scenarios[s.Name] = dst, r, s
	}
	return m.Run(), nil
}

type run struct {
	code   int
	stdout string
	stderr string
}

func cli(t *testing.T, args ...string) run {
	t.Helper()
	cmd := exec.Command(binary, args...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return run{code, out.String(), errb.String()}
}

// evaluate runs `codeeval evaluate` for a scenario variant and decodes JSON.
func evaluate(t *testing.T, scenario, base, target string, extra ...string) *report.Report {
	t.Helper()
	args := []string{"evaluate", "--repo", repos[scenario], "--base", base, "--target", target, "--format", "json", "--no-cache"}
	if c := scenarios[scenario].Config; c != "" {
		args = append(args, "--config", c)
	}
	args = append(args, extra...)
	r := cli(t, args...)
	if r.code != 0 {
		t.Fatalf("exit %d: %s", r.code, r.stderr)
	}
	var rep report.Report
	if err := json.Unmarshal([]byte(r.stdout), &rep); err != nil {
		t.Fatalf("decode: %v\n%s", err, r.stdout)
	}
	return &rep
}

func result(t *testing.T, r *report.Report, id string) *analyzer.Result {
	t.Helper()
	for _, res := range r.Results {
		if res.Analyzer == id {
			return res
		}
	}
	t.Fatalf("no result for %s", id)
	return nil
}

func score(t *testing.T, r *report.Report, id string) int {
	t.Helper()
	res := result(t, r, id)
	if res.Score == nil {
		t.Fatalf("%s: score is null (status %s: %s)", id, res.Status, res.Summary)
	}
	return *res.Score
}

func findings(res *analyzer.Result, kind string) []analyzer.Finding {
	var out []analyzer.Finding
	for _, f := range res.Findings {
		if f.Kind == kind {
			out = append(out, f)
		}
	}
	return out
}

func measurement(t *testing.T, res *analyzer.Result, metric, scope string) analyzer.Measurement {
	t.Helper()
	for _, m := range res.Measurements {
		if m.Metric == metric && (scope == "" || m.Scope == scope) {
			return m
		}
	}
	t.Fatalf("%s: no measurement %s/%s", res.Analyzer, metric, scope)
	return analyzer.Measurement{}
}

// lineOf returns the 1-based line of the first occurrence of needle in a
// fixture file, so evidence assertions track the fixture source.
func lineOf(t *testing.T, path, needle string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for i, l := range strings.Split(string(data), "\n") {
		if strings.Contains(l, needle) {
			return i + 1
		}
	}
	t.Fatalf("%q not in %s", needle, path)
	return 0
}

func hasEvidence(f analyzer.Finding, side, file, symbol string, line int) bool {
	for _, e := range f.Evidence {
		if e.Side == side && e.File == file && (symbol == "" || e.Symbol == symbol) && (line == 0 || e.StartLine == line) {
			return true
		}
	}
	return false
}

func TestFixtureVariantsPreserveBehaviour(t *testing.T) {
	if testing.Short() {
		t.Skip("runs go test in every fixture variant")
	}
	for name, s := range scenarios {
		for _, v := range append([]string{"base"}, s.Variants...) {
			if reason := s.NoCompile(v); reason != "" {
				t.Logf("%s/%s: not built by design: %s", name, v, reason)
				continue
			}
			t.Run(name+"/"+v, func(t *testing.T) {
				dir := t.TempDir()
				if err := s.Materialize(v, dir); err != nil {
					t.Fatal(err)
				}
				cmd := exec.Command("go", "test", "./...")
				cmd.Dir = dir
				cmd.Env = append(os.Environ(), "GOWORK=off", "GOFLAGS=-mod=mod", "GOTOOLCHAIN=local")
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("go test failed: %v\n%s", err, out)
				}
			})
		}
	}
}

func TestReuseVersusReimplementation(t *testing.T) {
	good := evaluate(t, "reuse", "base", "good")
	bad := evaluate(t, "reuse", "base", "bad")
	if s := score(t, good, "duplication"); s != 0 {
		t.Errorf("good duplication score %d, want 0", s)
	}
	if s := score(t, bad, "duplication"); s != -1 {
		t.Errorf("bad duplication score %d, want -1", s)
	}
	fs := findings(result(t, bad, "duplication"), "duplication.introduced_structural")
	if len(fs) != 1 {
		t.Fatalf("want 1 introduced_structural finding, got %+v", result(t, bad, "duplication").Findings)
	}
	f := fs[0]
	slugLine := lineOf(t, filepath.Join(fixtures, "reuse/base/text/slug.go"), "func Slugify")
	tagLine := lineOf(t, filepath.Join(fixtures, "reuse/bad/tags/tags.go"), "func slugTag")
	if !hasEvidence(f, "target", "text/slug.go", "Slugify", slugLine) || !hasEvidence(f, "target", "tags/tags.go", "slugTag", tagLine) {
		t.Errorf("evidence should point at Slugify:%d and slugTag:%d: %+v", slugLine, tagLine, f.Evidence)
	}
	if f.Caveat == "" || !strings.Contains(f.Message, "reuse candidate") {
		t.Errorf("structural finding must carry caveat and reuse-candidate wording: %+v", f)
	}
	for _, e := range f.Evidence {
		if e.Commit != refs["reuse"]["bad"] {
			t.Errorf("evidence commit %s, want pinned target %s", e.Commit, refs["reuse"]["bad"])
		}
	}
	if bad.Base.Commit != refs["reuse"]["base"] || bad.Target.Commit != refs["reuse"]["bad"] || bad.Base.Tree == "" {
		t.Errorf("revisions not pinned: %+v %+v", bad.Base, bad.Target)
	}
}

func TestConsolidationVersusObsoleteCopies(t *testing.T) {
	good := evaluate(t, "consolidate", "base", "good")
	bad := evaluate(t, "consolidate", "base", "bad")
	if s := score(t, good, "duplication"); s != 1 {
		t.Errorf("good duplication %d, want +1", s)
	}
	if s := score(t, bad, "duplication"); s != -1 {
		t.Errorf("bad duplication %d, want -1", s)
	}
	removed := findings(result(t, good, "duplication"), "duplication.removed")
	if len(removed) != 1 {
		t.Fatalf("want one removed finding: %+v", result(t, good, "duplication").Findings)
	}
	fmLine := lineOf(t, filepath.Join(fixtures, "consolidate/base/orders/orders.go"), "func formatMoney")
	if !hasEvidence(removed[0], "base", "orders/orders.go", "formatMoney", fmLine) || removed[0].Evidence[0].Commit != refs["consolidate"]["base"] {
		t.Errorf("removed code must carry base evidence: %+v", removed[0].Evidence)
	}
	if s := score(t, good, "obsolete"); s != 0 {
		t.Errorf("good obsolete %d, want 0", s)
	}
	if s := score(t, bad, "obsolete"); s != -1 {
		t.Errorf("bad obsolete %d, want -1", s)
	}
	orphans := findings(result(t, bad, "obsolete"), "obsolete.left_unreferenced")
	names := map[string]bool{}
	for _, f := range orphans {
		names[f.Evidence[1].Symbol] = true
		if f.Basis != analyzer.BasisHeuristic {
			t.Errorf("obsolete finding must be heuristic: %+v", f)
		}
	}
	if !names["formatMoney"] || !names["formatAmount"] || len(orphans) != 2 {
		t.Errorf("want formatMoney and formatAmount orphaned, got %+v", orphans)
	}
	// Consolidation into a renamed function is a rename, not a deletion + addition.
	var renamed int
	for _, s := range good.Change.Symbols {
		if s.Status == "renamed" {
			renamed++
		}
	}
	if renamed != 1 {
		t.Errorf("want one renamed symbol (formatMoney → Format), got %+v", good.Change.Symbols)
	}
}

func TestExtractionVersusNesting(t *testing.T) {
	good := evaluate(t, "complexity", "base", "good")
	bad := evaluate(t, "complexity", "base", "bad")
	g, b := result(t, good, "complexity"), result(t, bad, "complexity")
	if s := score(t, good, "complexity"); s != 1 {
		t.Errorf("good complexity %d, want +1: %s", s, g.ScoreRule)
	}
	if s := score(t, bad, "complexity"); s != -1 {
		t.Errorf("bad complexity %d, want -1: %s", s, b.ScoreRule)
	}
	mb, mg := measurement(t, b, "cyclomatic_sum", "modified_callables"), measurement(t, g, "cyclomatic_sum", "modified_callables")
	if *mb.Before != 4 || *mb.After != 10 || *mg.Before != 4 || *mg.After != 1 {
		t.Errorf("Discount cyclomatic: bad %d→%d (want 4→10), good %d→%d (want 4→1)", *mb.Before, *mb.After, *mg.Before, *mg.After)
	}
	if v := *measurement(t, g, "callables_added", "").Value; v != 3 {
		t.Errorf("good should add 3 helpers, got %d", v)
	}
	// Extraction must not be reported as net removal.
	total := measurement(t, g, "cyclomatic_sum", "changed_scope")
	if *total.Delta <= 0 {
		t.Errorf("changed-scope total should increase for extraction, got %+v", total)
	}
	if !strings.Contains(strings.Join(g.Rationale, " "), "extraction/redistribution") {
		t.Errorf("rationale should explain extraction: %v", g.Rationale)
	}
	if h := measurement(t, b, "hotspot_callables", "changed_scope"); *h.Before != 0 || *h.After != 1 {
		t.Errorf("bad should create one hotspot: %+v", h)
	}
	discLine := lineOf(t, filepath.Join(fixtures, "complexity/bad/pricing/pricing.go"), "func Discount")
	mods := findings(b, "complexity.modified")
	if len(mods) != 1 || !hasEvidence(mods[0], "target", "pricing/pricing.go", "Discount", discLine) || !hasEvidence(mods[0], "base", "pricing/pricing.go", "Discount", 0) {
		t.Errorf("modified finding should cite Discount before/after: %+v", mods)
	}
}

func TestArchitectureRuleViolation(t *testing.T) {
	good := evaluate(t, "layering", "base", "good")
	bad := evaluate(t, "layering", "base", "bad")
	if s := score(t, good, "dependencies"); s != 0 {
		t.Errorf("good architecture %d, want 0", s)
	}
	if s := score(t, bad, "dependencies"); s != -1 {
		t.Errorf("bad architecture %d, want -1", s)
	}
	v := findings(result(t, bad, "dependencies"), "dependencies.rule_violation_introduced")
	importLine := lineOf(t, filepath.Join(fixtures, "layering/bad/internal/store/dump.go"), `"example.com/layering/internal/cli"`)
	if len(v) != 1 || !hasEvidence(v[0], "target", "internal/store/dump.go", "", importLine) {
		t.Errorf("violation should cite the import line %d: %+v", importLine, v)
	}
	// Without the config, the same patch is not judged on rules.
	unconfigured := cli(t, "evaluate", "--repo", repos["layering"], "--base", "base", "--target", "bad", "--format", "json", "--no-cache")
	var r report.Report
	if err := json.Unmarshal([]byte(unconfigured.stdout), &r); err != nil {
		t.Fatal(err)
	}
	if s := score(t, &r, "dependencies"); s != 0 {
		t.Errorf("unconfigured rules should not produce violations, got %d", s)
	}
	if !strings.Contains(strings.Join(result(t, &r, "dependencies").Rationale, " "), "no architecture rules configured") {
		t.Error("should say rules are not configured")
	}
}

func TestCycleIntroducedAndRemoved(t *testing.T) {
	in := evaluate(t, "cycle", "base", "introduce")
	dec := evaluate(t, "cycle", "base", "decouple")
	rm := evaluate(t, "cycle-removal", "base", "remove")
	if s := score(t, in, "dependencies"); s != -2 {
		t.Errorf("introduce %d, want -2", s)
	}
	if len(findings(result(t, in, "dependencies"), "dependencies.cycle_introduced")) != 1 {
		t.Errorf("want cycle_introduced finding: %+v", result(t, in, "dependencies").Findings)
	}
	if s := score(t, dec, "dependencies"); s != 1 {
		t.Errorf("decouple %d, want +1", s)
	}
	if s := score(t, rm, "dependencies"); s != 2 {
		t.Errorf("remove %d, want +2", s)
	}
	cyc := measurement(t, result(t, rm, "dependencies"), "cycles", "repository")
	if *cyc.Before != 1 || *cyc.After != 0 {
		t.Errorf("cycles %+v", cyc)
	}
}

func TestFocusedVersusUnrelatedChurn(t *testing.T) {
	good := evaluate(t, "scope", "base", "good")
	bad := evaluate(t, "scope", "base", "bad")
	if s := score(t, good, "scope"); s != 0 {
		t.Errorf("good scope %d, want 0", s)
	}
	if s := score(t, bad, "scope"); s != -1 {
		t.Errorf("bad scope %d, want -1: %s", s, result(t, bad, "scope").ScoreRule)
	}
	if v := *measurement(t, result(t, bad, "scope"), "symbols_formatting_only", "").Value; v != 5 {
		t.Errorf("want 5 formatting-only symbols (Header, Footer, pageOf, Subject, Body), got %d", v)
	}
	if len(findings(result(t, bad, "scope"), "scope.unrelated_churn")) != 2 {
		t.Errorf("want churn findings for report and notify: %+v", result(t, bad, "scope").Findings)
	}
	for _, id := range []string{"complexity", "duplication"} {
		if score(t, good, id) != 0 || score(t, bad, id) != 0 {
			t.Errorf("%s should be neutral for both scope variants", id)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("unchanged", func(t *testing.T) {
		r := evaluate(t, "edge", "base", "base")
		if r.Change.FilesChanged != 0 || len(r.Change.Symbols) != 0 {
			t.Fatalf("unchanged baseline should have no changes: %+v", r.Change)
		}
		for _, res := range r.Results {
			if res.Analyzer == "correctness" {
				if res.Score != nil || res.Status != "not_assessed" {
					t.Errorf("correctness must be not_assessed: %+v", res)
				}
				continue
			}
			if res.Score == nil || *res.Score != 0 || res.Status != "ok" {
				t.Errorf("%s: want ok/0, got %s/%v", res.Analyzer, res.Status, res.Score)
			}
		}
	})
	t.Run("rename-only", func(t *testing.T) {
		r := evaluate(t, "edge", "base", "rename-only")
		if len(r.Change.Files) != 1 || r.Change.Files[0].Status != "renamed" || r.Change.Files[0].OldPath != "util/util.go" {
			t.Fatalf("want one renamed file: %+v", r.Change.Files)
		}
		for _, s := range r.Change.Symbols {
			if s.Status != "moved" {
				t.Errorf("symbol %s: %s, want moved", s.Name, s.Status)
			}
		}
		for _, id := range []string{"duplication", "complexity", "dependencies", "scope", "obsolete"} {
			if s := score(t, r, id); s != 0 {
				t.Errorf("%s: %d, want 0 for rename-only", id, s)
			}
		}
	})
	t.Run("rename-symbol", func(t *testing.T) {
		r := evaluate(t, "edge", "base", "rename-symbol")
		st := map[string]string{}
		for _, s := range r.Change.Symbols {
			st[s.Name] = s.Status
			if s.Status == "renamed" && s.BaseName != "helper" {
				t.Errorf("renamed from %q, want helper", s.BaseName)
			}
		}
		if st["normalize"] != "renamed" || st["Greet"] != "modified" || len(st) != 2 {
			t.Errorf("statuses %v", st)
		}
		if score(t, r, "duplication") != 0 || score(t, r, "complexity") != 0 || score(t, r, "obsolete") != 0 {
			t.Error("renaming must not be penalised")
		}
	})
	t.Run("added-deleted", func(t *testing.T) {
		r := evaluate(t, "edge", "base", "added-deleted")
		c := result(t, r, "complexity")
		if *measurement(t, c, "callables_added", "").Value != 1 || *measurement(t, c, "callables_deleted", "").Value != 1 {
			t.Errorf("want 1 added/1 deleted: %+v", c.Measurements)
		}
		var del *report.SymbolRow
		for i, s := range r.Change.Symbols {
			if s.Status == "deleted" {
				del = &r.Change.Symbols[i]
			}
		}
		legacy := lineOf(t, filepath.Join(fixtures, "edge/base/util/util.go"), "func Legacy")
		if del == nil || del.Name != "Legacy" || del.BaseFile != "util/util.go" || del.BaseStartLine != legacy {
			t.Errorf("deleted symbol needs base location line %d: %+v", legacy, del)
		}
		api := findings(result(t, r, "scope"), "scope.api_removed")
		if len(api) != 1 {
			t.Errorf("exported Legacy removal should be an API finding: %+v", result(t, r, "scope").Findings)
		}
	})
	t.Run("parse-error", func(t *testing.T) {
		r := evaluate(t, "edge", "base", "parse-error")
		found := false
		for _, d := range r.Diagnostics.Change {
			if d.Code == "parse_error" && d.File == "util/util.go" && d.Side == "target" {
				found = true
			}
		}
		if !found {
			t.Errorf("want parse_error diagnostic: %+v", r.Diagnostics)
		}
		c := result(t, r, "complexity")
		if c.Status != "partial" || c.Coverage.WithParseErrors != 1 {
			t.Errorf("complexity should be partial with a parse error: %s %+v", c.Status, c.Coverage)
		}
	})
	t.Run("unsupported-only", func(t *testing.T) {
		r := evaluate(t, "edge", "base", "unsupported-only")
		for _, id := range []string{"duplication", "complexity", "dependencies", "scope", "obsolete"} {
			res := result(t, r, id)
			if res.Status != "unsupported" || res.Score != nil {
				t.Errorf("%s: want unsupported/null, got %s/%v", id, res.Status, res.Score)
			}
		}
		for _, l := range r.Languages {
			if l.ID == "kotlin" && (l.Supported || l.FilesChanged != 1) {
				t.Errorf("kotlin coverage wrong: %+v", l)
			}
		}
	})
	t.Run("odd-filename", func(t *testing.T) {
		r := evaluate(t, "edge", "base", "odd-filename")
		if len(r.Change.Files) != 1 || r.Change.Files[0].Path != "util/odd name ü.go" || r.Change.Files[0].Status != "added" {
			t.Fatalf("odd filename not preserved: %+v", r.Change.Files)
		}
		if len(r.Change.Symbols) != 1 || r.Change.Symbols[0].File != "util/odd name ü.go" || r.Change.Symbols[0].Name != "Shout" {
			t.Errorf("symbol in odd filename: %+v", r.Change.Symbols)
		}
	})
	t.Run("mixed", func(t *testing.T) {
		r := evaluate(t, "edge", "base", "mixed")
		c := result(t, r, "complexity")
		if c.Status != "partial" || c.Score == nil || c.Coverage.UnsupportedByLanguage["kotlin"] != 1 {
			t.Errorf("mixed change: want partial with a score and Kotlin counted: %s %v %+v", c.Status, c.Score, c.Coverage)
		}
	})
}

func TestDeterministicOutput(t *testing.T) {
	cache := t.TempDir()
	args := []string{"evaluate", "--repo", repos["consolidate"], "--base", "base", "--target", "bad", "--format", "json", "--cache-dir", cache}
	cold := cli(t, args...)
	warm := cli(t, args...)
	nocache := cli(t, "evaluate", "--repo", repos["consolidate"], "--base", "base", "--target", "bad", "--format", "json", "--no-cache")
	if cold.code != 0 || cold.stdout != warm.stdout || cold.stdout != nocache.stdout {
		t.Fatalf("output differs between runs (cold/warm/no-cache)")
	}
	if _, err := os.Stat(filepath.Join(cache, "index.sqlite")); err != nil {
		t.Errorf("SQLite index not created: %v", err)
	}
}

// repoState fingerprints everything a mutation could touch.
func repoState(t *testing.T, dir string) string {
	t.Helper()
	h := sha256.New()
	for _, args := range [][]string{{"rev-parse", "HEAD"}, {"for-each-ref"}, {"status", "--porcelain=v1", "--untracked-files=all"}, {"stash", "list"}, {"worktree", "list", "--porcelain"}} {
		cmd := exec.Command("git", append([]string{"--no-optional-locks"}, args...)...)
		cmd.Dir = dir
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		h.Write(out)
	}
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, _ := d.Info()
		data, _ := os.ReadFile(p)
		fmt.Fprintf(h, "%s %d %s\n", p, info.ModTime().UnixNano(), sha256.Sum256(data))
		return nil
	})
	return hex.EncodeToString(h.Sum(nil))
}

func TestEvaluatedRepositoryIsNotModified(t *testing.T) {
	dir := repos["edge"]
	before := repoState(t, dir)
	evaluate(t, "edge", "base", "mixed")
	cli(t, "history", "--repo", dir, "--base", "base", "--target", "mixed", "--format", "none", "--no-cache")
	if after := repoState(t, dir); after != before {
		t.Fatal("evaluated repository changed")
	}
	r := cli(t, "evaluate", "--repo", dir, "--base", "base", "--target", "mixed", "--output", filepath.Join(dir, "out.json"), "--no-cache")
	if r.code == 0 || !strings.Contains(r.stderr, "inside the evaluated repository") {
		t.Errorf("writing output inside the repo must be refused: %d %s", r.code, r.stderr)
	}
	r = cli(t, "evaluate", "--repo", dir, "--base", "base", "--target", "mixed", "--cache-dir", filepath.Join(dir, ".cache"))
	if r.code == 0 {
		t.Error("cache inside the repo must be refused")
	}
	if after := repoState(t, dir); after != before {
		t.Fatal("evaluated repository changed by refused runs")
	}
}

func gitT(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=T", "-c", "user.email=t@example.invalid", "-c", "core.hooksPath=/dev/null", "-c", "commit.gpgsign=false"}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2026-01-02T00:00:00Z", "GIT_COMMITTER_DATE=2026-01-02T00:00:00Z", "GIT_CONFIG_GLOBAL=/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestWorkingTreeMode(t *testing.T) {
	dir := t.TempDir()
	if _, err := scenarios["edge"].BuildRepo(dir); err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(filepath.Join(fixtures, "edge/added-deleted/util/util.go"))
	if err := os.WriteFile(filepath.Join(dir, "util/util.go"), content, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "private-notes.txt"), []byte("do-not-read-me"), 0o600); err != nil {
		t.Fatal(err)
	}
	before := repoState(t, dir)
	r := cli(t, "evaluate", "--repo", dir, "--base", "base", "--working-tree", "--untracked=false", "--format", "json", "--no-cache")
	if r.code != 0 {
		t.Fatal(r.stderr)
	}
	if strings.Contains(r.stdout, "private-notes") || strings.Contains(r.stdout, "do-not-read-me") {
		t.Fatal("untracked file leaked into a --untracked=false report")
	}
	var rep report.Report
	if err := json.Unmarshal([]byte(r.stdout), &rep); err != nil {
		t.Fatal(err)
	}
	if rep.Mode != "worktree" || rep.Target.Commit != "WORKTREE" || !strings.HasPrefix(rep.Target.Tree, "worktree:") {
		t.Errorf("worktree identity: %+v", rep.Target)
	}
	if len(rep.Change.Files) != 1 || rep.Change.Files[0].Path != "util/util.go" {
		t.Errorf("want only util/util.go changed: %+v", rep.Change.Files)
	}
	if *measurement(t, result(t, &rep, "complexity"), "callables_added", "").Value != 1 {
		t.Error("worktree change should add Classify")
	}
	if after := repoState(t, dir); after != before {
		t.Fatal("worktree evaluation modified the repository")
	}
}

// Untracked files are read by default in worktree mode (not ignored ones);
// --untracked=false restores tracked-only.
func TestWorkingTreeUntracked(t *testing.T) {
	dir := t.TempDir()
	if _, err := scenarios["edge"].BuildRepo(dir); err != nil {
		t.Fatal(err)
	}
	extra := "package util\n\nfunc UntrackedHelper(s string) string {\n\tif s == \"\" {\n\t\treturn \"x\"\n\t}\n\treturn s\n}\n"
	os.WriteFile(filepath.Join(dir, "util/extra.go"), []byte(extra), 0o644)
	os.WriteFile(filepath.Join(dir, "util/ignored.go"), []byte("package util\n\nfunc IgnoredSecretHelper() {}\n"), 0o644)
	os.WriteFile(filepath.Join(dir, ".git/info/exclude"), []byte("util/ignored.go\n"), 0o644)
	os.WriteFile(filepath.Join(dir, "private-notes.txt"), []byte("do-not-read-me"), 0o600)
	before := repoState(t, dir)
	run := func(extraArgs ...string) string {
		args := append([]string{"evaluate", "--repo", dir, "--base", "base", "--working-tree", "--format", "json", "--no-cache"}, extraArgs...)
		r := cli(t, args...)
		if r.code != 0 {
			t.Fatal(r.stderr)
		}
		return r.stdout
	}
	def := run()
	if !strings.Contains(def, "util/extra.go") || !strings.Contains(def, "untracked files that are not ignored") {
		t.Error("untracked source file not analysed by default")
	}
	if strings.Contains(def, "IgnoredSecretHelper") || strings.Contains(def, "util/ignored.go") || strings.Contains(def, "do-not-read-me") {
		t.Error("ignored file or private content leaked")
	}
	tracked := run("--untracked=false")
	if strings.Contains(tracked, "util/extra.go") || !strings.Contains(tracked, "tracked files only") {
		t.Error("--untracked=false still read untracked files")
	}
	d := cli(t, "defects", "--repo", dir, "--base", "base", "--working-tree", "--format", "json", "--no-cache")
	if !strings.Contains(d.stdout, "util/extra.go") && d.code > 3 {
		t.Errorf("defects worktree: %d %s", d.code, d.stderr)
	}
	if after := repoState(t, dir); after != before {
		t.Fatal("worktree evaluation modified the repository")
	}
}

func TestHistoryAttributesFirstParentSteps(t *testing.T) {
	dir := t.TempDir()
	r, err := scenarios["edge"].BuildRepo(dir)
	if err != nil {
		t.Fatal(err)
	}
	gitT(t, dir, "checkout", "-q", "-B", "main", "base")
	gitT(t, dir, "checkout", "-q", "-b", "feature")
	content, _ := os.ReadFile(filepath.Join(fixtures, "edge/added-deleted/util/util.go"))
	os.WriteFile(filepath.Join(dir, "util/util.go"), content, 0o644)
	gitT(t, dir, "commit", "-qam", "feature: classify")
	os.WriteFile(filepath.Join(dir, "web/app.ts"), []byte("export const x = 1;\n"), 0o644)
	gitT(t, dir, "commit", "-qam", "feature: web")
	gitT(t, dir, "checkout", "-q", "main")
	gitT(t, dir, "merge", "-q", "--no-ff", "-m", "Merge feature", "feature")
	// Rename helper on top of the merged feature (Classify stays).
	cur, _ := os.ReadFile(filepath.Join(dir, "util/util.go"))
	os.WriteFile(filepath.Join(dir, "util/util.go"), []byte(strings.ReplaceAll(string(cur), "helper", "normalize")), 0o644)
	gitT(t, dir, "commit", "-qam", "rename helper")
	head := gitT(t, dir, "rev-parse", "HEAD")

	out := cli(t, "history", "--repo", dir, "--base", r["base"], "--target", head, "--format", "json", "--no-cache")
	if out.code != 0 {
		t.Fatal(out.stderr)
	}
	var h report.History
	if err := json.Unmarshal([]byte(out.stdout), &h); err != nil {
		t.Fatal(err)
	}
	if len(h.Steps) != 2 {
		t.Fatalf("want 2 first-parent steps (merge, rename), got %d", len(h.Steps))
	}
	if h.Steps[0].Parent != r["base"] || h.Steps[0].Subject != "Merge feature" || h.Steps[1].Subject != "rename helper" {
		t.Errorf("steps: %+v / %+v", h.Steps[0].Subject, h.Steps[1].Subject)
	}
	if h.Net == nil || h.Net.Base.Commit != r["base"] || h.Net.Target.Commit != head {
		t.Fatal("net report missing or unpinned")
	}
	// The merge step carries the whole feature branch (Go + TS), once.
	if h.Steps[0].Report.Change.FilesChanged != 2 || h.Steps[1].Report.Change.FilesChanged != 1 {
		t.Errorf("files per step: %d, %d", h.Steps[0].Report.Change.FilesChanged, h.Steps[1].Report.Change.FilesChanged)
	}
}

func TestCLIContract(t *testing.T) {
	if r := cli(t, "evaluate", "--repo", repos["edge"]); r.code != 2 {
		t.Errorf("missing --base should exit 2, got %d", r.code)
	}
	if r := cli(t, "evaluate", "--repo", repos["edge"], "--base", "base", "--target", "nope-not-a-ref", "--no-cache"); r.code != 1 || !strings.Contains(r.stderr, "cannot resolve") {
		t.Errorf("bad ref: %d %s", r.code, r.stderr)
	}
	if r := cli(t, "evaluate", "--repo", repos["edge"], "--base", "--upload-pack=x", "--target", "base", "--no-cache"); r.code == 0 {
		t.Error("option-like revision must be rejected")
	}
	langs := cli(t, "languages").stdout
	for _, want := range []string{"go ", "typescript ", "ts-adapter", "csharp", "dependencies", "never scored"} {
		if !strings.Contains(langs, want) {
			t.Errorf("languages output missing %q", want)
		}
	}
	an := cli(t, "analyzers").stdout
	for _, id := range []string{"duplication", "complexity", "dependencies", "scope", "obsolete", "correctness"} {
		if !strings.Contains(an, id) {
			t.Errorf("analyzers missing %s", id)
		}
	}
	md := cli(t, "evaluate", "--repo", repos["reuse"], "--base", "base", "--target", "bad", "--no-cache").stdout
	for _, want := range []string{"no overall score", "duplication.introduced_structural", "tags/tags.go", "## Limitations"} {
		if !strings.Contains(md, want) {
			t.Errorf("markdown summary missing %q", want)
		}
	}
}

func TestTypeScriptScenarios(t *testing.T) {
	good := evaluate(t, "ts", "base", "reuse-good")
	bad := evaluate(t, "ts", "base", "reuse-bad")
	if score(t, good, "duplication") != 0 || score(t, bad, "duplication") != -1 {
		t.Errorf("TS reuse pair: good %d bad %d, want 0 / -1", score(t, good, "duplication"), score(t, bad, "duplication"))
	}
	fs := findings(result(t, bad, "duplication"), "duplication.introduced_structural")
	if len(fs) != 1 || !hasEvidence(fs[0], "target", "src/text/slug.ts", "slugify", 1) || !hasEvidence(fs[0], "target", "src/tags/tags.ts", "slugTag", 1) {
		t.Errorf("TS clone evidence: %+v", fs)
	}
	for _, l := range bad.Languages {
		if l.ID == "typescript" && (!l.Supported || l.FilesChanged != 1) {
			t.Errorf("typescript should be supported: %+v", l)
		}
	}
	nest := evaluate(t, "ts", "base", "nest-bad")
	extract := evaluate(t, "ts", "base", "extract-good")
	if score(t, nest, "complexity") >= 0 || score(t, extract, "complexity") < 0 {
		t.Errorf("TS complexity pair: nest %d extract %d", score(t, nest, "complexity"), score(t, extract, "complexity"))
	}
	cyc := evaluate(t, "ts", "base", "cycle")
	if score(t, cyc, "dependencies") != -2 || len(findings(result(t, cyc, "dependencies"), "dependencies.cycle_introduced")) != 1 {
		t.Errorf("TS directory cycle should score -2: %+v", result(t, cyc, "dependencies"))
	}
	// v2 snapshot sees the TypeScript clone too.
	g, b := snapshotV2(t, "ts", "reuse-good"), snapshotV2(t, "ts", "reuse-bad")
	if !(b.Codebase.Debt.Duplication > g.Codebase.Debt.Duplication) {
		t.Errorf("v2 TS duplication debt: bad %.2f good %.2f", b.Codebase.Debt.Duplication, g.Codebase.Debt.Duplication)
	}
}
