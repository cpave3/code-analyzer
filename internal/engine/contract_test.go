package engine_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/analyzers"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/model"
)

// toyAdapter is a minimal second language: each "def NAME" line is a
// function. It declares only the symbols capability, so analyzers that need
// complexity, fingerprints or dependencies must report it as unsupported.
type toyAdapter struct{}

func (toyAdapter) ID() string           { return "toy" }
func (toyAdapter) Version() string      { return "toy/1" }
func (toyAdapter) Extensions() []string { return []string{".toy"} }
func (toyAdapter) Capabilities() map[lang.Capability]string {
	return map[lang.Capability]string{lang.CapSymbols: "def lines"}
}
func (toyAdapter) Limits() []string          { return []string{"toy"} }
func (toyAdapter) IsTestPath(string) bool    { return false }
func (toyAdapter) PackageOf(p string) string { return p }
func (toyAdapter) Parse(src []byte) (*model.ParsedFile, error) {
	pf := &model.ParsedFile{}
	for i, l := range strings.Split(string(src), "\n") {
		if name, ok := strings.CutPrefix(l, "def "); ok {
			pf.Symbols = append(pf.Symbols, model.Symbol{Name: name, Kind: model.KindFunction, StartLine: i + 1, EndLine: i + 1, ExactFP: name})
		}
	}
	return pf, nil
}
func (toyAdapter) BuildGraph(lang.GraphInput) (*lang.Graph, []model.Diagnostic) { return nil, nil }

// symbolCount is a new analyzer added without touching existing ones.
type symbolCount struct{}

func (symbolCount) Descriptor() analyzer.Descriptor {
	return analyzer.Descriptor{ID: "symbol-count", Dimension: analyzer.DimMaintainability, Title: "Symbols added", Version: "1",
		Requires: []lang.Capability{lang.CapSymbols}, Classes: []model.FileClass{model.ClassSource}}
}

func (symbolCount) Analyze(ctx *analyzer.Context) (*analyzer.Result, error) {
	n := 0
	for _, s := range ctx.Change.Symbols {
		if s.Status == "added" {
			n++
		}
	}
	zero := 0
	return &analyzer.Result{Score: &zero, ScoreRule: "informational", Summary: fmt.Sprintf("%d added", n),
		Measurements: []analyzer.Measurement{analyzer.Count("added", "", n, "symbols")}}, nil
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=T", "-c", "user.email=t@example.invalid", "-c", "commit.gpgsign=false", "-c", "core.hooksPath=/dev/null"}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
}

func TestLanguageAndAnalyzerRegistrationContract(t *testing.T) {
	dir := t.TempDir()
	write := func(p, s string) {
		os.MkdirAll(filepath.Dir(filepath.Join(dir, p)), 0o755)
		if err := os.WriteFile(filepath.Join(dir, p), []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, dir, "init", "-q")
	write("go.mod", "module example.com/toy\n\ngo 1.22\n")
	write("a/a.go", "package a\n\nfunc A() int { return 1 }\n")
	write("lib/x.toy", "def alpha\n")
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-qm", "base")
	git(t, dir, "tag", "base")
	write("lib/x.toy", "def alpha\ndef beta\n")
	write("a/a.go", "package a\n\nfunc A() int { return 1 }\n\nfunc B() int { return 2 }\n")
	git(t, dir, "commit", "-qam", "target")

	// The CLI's registries, extended through the same public API.
	langs := languages.Default()
	if err := langs.Register(toyAdapter{}); err != nil {
		t.Fatal(err)
	}
	if err := langs.Register(toyAdapter{}); err == nil {
		t.Fatal("duplicate language registration must fail")
	}
	reg := analyzers.Default()
	if err := reg.Register(symbolCount{}); err != nil {
		t.Fatal(err)
	}
	if err := reg.Register(symbolCount{}); err == nil {
		t.Fatal("duplicate analyzer registration must fail")
	}

	ev, err := engine.New(engine.Options{RepoPath: dir, Config: config.Default(), Langs: langs, Analyzers: reg})
	if err != nil {
		t.Fatal(err)
	}
	defer ev.Close()
	rep, err := ev.Evaluate("base", "HEAD", false)
	if err != nil {
		t.Fatal(err)
	}
	byID := map[string]*analyzer.Result{}
	for _, r := range rep.Results {
		byID[r.Analyzer] = r
	}
	sc := byID["symbol-count"]
	if sc == nil || *sc.Measurements[0].Value != 2 || sc.Status != analyzer.StatusOK {
		t.Fatalf("new analyzer should see Go and toy symbols: %+v", sc)
	}
	// Existing analyzers see the toy file as unsupported for capabilities
	// it lacks, without any change to their code.
	cx := byID["complexity"]
	if cx.Status != analyzer.StatusPartial || cx.Coverage.UnsupportedByLanguage["toy"] != 1 || cx.Coverage.Analyzed != 1 {
		t.Fatalf("complexity coverage: %s %+v", cx.Status, cx.Coverage)
	}
	if sc := byID["scope"]; sc.Coverage.Analyzed != 2 {
		t.Fatalf("scope (symbols only) should analyse both languages: %+v", sc.Coverage)
	}
	var toy bool
	for _, l := range rep.Languages {
		if l.ID == "toy" && l.Supported && l.FilesChanged == 1 {
			toy = true
		}
	}
	if !toy {
		t.Fatalf("toy language coverage missing: %+v", rep.Languages)
	}
}
