// Command codeeval evaluates the engineering-quality delta of a Git change in
// the context of its repository. It reads the repository strictly read-only.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cameron/code-evaluator/internal/analyzers"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/report"
)

const usage = `codeeval — repository-aware code change evaluator (POC)

Usage:
  codeeval evaluate --repo PATH --base REV (--target REV | --working-tree) [flags]
  codeeval history  --repo PATH --base REV --target REV [flags]
  codeeval score    --repo PATH --rev REV [flags]                 (experimental scoring/v2: absolute snapshot)
  codeeval score    --repo PATH --base REV --target REV [--steps] [flags]   (scoring/v2 delta)
  codeeval explain  --repo PATH --rev REV [--path P ...] [--top N]                     (scoring/v2 debt factors for a scope)
  codeeval check    --repo PATH --base REV (--target REV | --working-tree) [--path P ...] [--fail-on new-debt|net-debt|gates|none] [--top N]
                                                                                    (PR check: factors introduced/worsened; exit 3 on failure)
  codeeval graph    --repo PATH --rev REV --output FILE.html [--json FILE] [--max-callables N] [--config FILE]   (POC debt graph view)
  codeeval defects  --repo PATH (--rev REV | --base REV (--target REV | --working-tree)) [--rule ID ...] [--path P ...] [--fail-on gate|any|none] [--top N]   (defect rules, rules/v1)
  codeeval analyze  --repo PATH (--rev REV | --base REV (--target REV | --working-tree)) [--path P ...] [--debt-fail-on ...] [--defects-fail-on ...]   (v2 debt + defects in one report)
  codeeval rules                                                          (defect rules × languages)
  codeeval attribution --repo PATH --base REV --target REV [--group-by agent,harness,workflow] [--defects] [--infer-co-authors=false]   (step deltas by commit provenance)
  codeeval languages
  codeeval analyzers
  codeeval version

Flags for evaluate/history:
  --config FILE      evaluation config (JSON, configVersion 1); defaults otherwise
  --output FILE      write the JSON report
  --summary FILE     write the Markdown summary
  --format FMT       stdout format: md (default), json, none
  --cache-dir DIR    SQLite index location (default: user cache dir/codeeval)
  --no-cache         use an in-memory index only
  --timing           include timing metadata (makes output non-deterministic)

Additional flags for score:
  --modules N        modules shown in the Markdown summary (default 25; JSON has all)
  --top N            callables and clone groups listed (default 25; 0 = all)
  --steps            with --base/--target: add a ledger for each first-parent step

The evaluated repository is never modified: no checkout, reset, stash, fetch,
hooks or external diff drivers; outputs and cache may not be placed inside it.
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "evaluate":
		return evaluate(args[1:], stdout, stderr, false)
	case "history":
		return evaluate(args[1:], stdout, stderr, true)
	case "score":
		return scoreCmd(args[1:], stdout, stderr)
	case "explain":
		return explainCmd(args[1:], stdout, stderr)
	case "check":
		return checkCmd(args[1:], stdout, stderr)
	case "graph":
		return graphCmd(args[1:], stdout, stderr)
	case "defects":
		return defectsCmd(args[1:], stdout, stderr)
	case "rules":
		return listRules(stdout)
	case "analyze":
		return analyzeCmd(args[1:], stdout, stderr)
	case "attribution":
		return attributionCmd(args[1:], stdout, stderr)
	case "languages":
		return listLanguages(stdout)
	case "analyzers":
		return listAnalyzers(stdout)
	case "version":
		fmt.Fprintf(stdout, "codeeval %s (%s, report %s)\n", report.ToolVersion, config.ScoringVersion, report.Schema)
		return 0
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return 0
	}
	fmt.Fprintf(stderr, "unknown command %q\n\n%s", args[0], usage)
	return 2
}

func evaluate(args []string, stdout, stderr io.Writer, history bool) int {
	name := "evaluate"
	if history {
		name = "history"
	}
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	repo := fs.String("repo", ".", "repository path")
	base := fs.String("base", "", "base revision")
	target := fs.String("target", "", "target revision")
	worktree := fs.Bool("working-tree", false, "evaluate tracked working-tree changes against base")
	cfgPath := fs.String("config", "", "config file")
	output := fs.String("output", "", "JSON report path")
	summary := fs.String("summary", "", "Markdown summary path")
	format := fs.String("format", "md", "stdout format: md, json, none")
	cacheDir := fs.String("cache-dir", "", "index cache directory")
	noCache := fs.Bool("no-cache", false, "in-memory index only")
	timing := fs.Bool("timing", false, "include timing metadata")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "unexpected arguments: %v\n", fs.Args())
		return 2
	}
	if *base == "" {
		fmt.Fprintln(stderr, "--base is required")
		return 2
	}
	if history && (*target == "" || *worktree) {
		fmt.Fprintln(stderr, "history needs --target and does not support --working-tree")
		return 2
	}
	if !history && (*target == "") == !*worktree {
		fmt.Fprintln(stderr, "exactly one of --target or --working-tree is required")
		return 2
	}
	if *format != "md" && *format != "json" && *format != "none" {
		fmt.Fprintln(stderr, "--format must be md, json or none")
		return 2
	}
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	dir := *cacheDir
	if *noCache {
		dir = ""
	} else if dir == "" {
		if ucd, err := os.UserCacheDir(); err == nil {
			dir = filepath.Join(ucd, "codeeval")
		}
	}
	ev, err := engine.New(engine.Options{RepoPath: *repo, Config: cfg, CacheDir: dir, Langs: languages.Default(), Analyzers: analyzers.Default(), Timing: *timing})
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	defer ev.Close()
	for _, p := range []string{*output, *summary} {
		if p != "" && inside(ev.RepoPath(), p) {
			fmt.Fprintf(stderr, "error: refusing to write %s inside the evaluated repository\n", p)
			return 1
		}
	}

	var payload any
	var md bytes.Buffer
	if history {
		h, err := ev.History(*base, *target, func(s string) { fmt.Fprintln(stderr, "codeeval:", s) })
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		payload = h
		report.WriteHistoryMarkdown(&md, h)
		md.WriteString("\n---\n\n")
		report.WriteMarkdown(&md, h.Net)
	} else {
		r, err := ev.Evaluate(*base, *target, *worktree)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		payload = r
		report.WriteMarkdown(&md, r)
	}
	js, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	js = append(js, '\n')
	if *output != "" {
		if err := writeFile(*output, js); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}
	if *summary != "" {
		if err := writeFile(*summary, md.Bytes()); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}
	switch *format {
	case "md":
		stdout.Write(md.Bytes())
	case "json":
		stdout.Write(js)
	}
	return 0
}

func inside(repo, p string) bool {
	abs, err := filepath.Abs(p)
	if err != nil {
		return true
	}
	if resolved, err := filepath.EvalSymlinks(filepath.Dir(abs)); err == nil {
		abs = filepath.Join(resolved, filepath.Base(abs))
	}
	rel, err := filepath.Rel(repo, abs)
	return err == nil && !strings.HasPrefix(rel, "..")
}

func writeFile(p string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func listLanguages(w io.Writer) int {
	reg := languages.Default()
	fmt.Fprintln(w, "Supported languages (adapters):")
	for _, a := range reg.Adapters() {
		fmt.Fprintf(w, "  %s  %s  extensions %v\n", a.ID(), a.Version(), a.Extensions())
		caps := a.Capabilities()
		for _, c := range lang.AllCapabilities {
			if note, ok := caps[c]; ok {
				fmt.Fprintf(w, "    + %-13s %s\n", c, note)
			} else {
				fmt.Fprintf(w, "    - %-13s not supported\n", c)
			}
		}
		for _, l := range a.Limits() {
			fmt.Fprintf(w, "    ! %s\n", l)
		}
	}
	fmt.Fprintln(w, "Recognised but unsupported (counted as coverage gaps, never scored):")
	for _, d := range reg.DeclaredLanguages() {
		fmt.Fprintf(w, "  %s  extensions %v  %s\n", d.ID, d.Extensions, d.Note)
	}
	return 0
}

func listAnalyzers(w io.Writer) int {
	for _, a := range analyzers.Default().All() {
		d := a.Descriptor()
		fmt.Fprintf(w, "%-12s dimension=%s version=%s requires=%v classes=%v\n    %s\n", d.ID, d.Dimension, d.Version, d.Requires, d.Classes, d.Description)
	}
	return 0
}
