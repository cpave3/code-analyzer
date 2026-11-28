// Command codeeval evaluates the engineering-quality delta of a Git change in
// the context of its repository. It reads the repository strictly read-only.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cameron/code-evaluator/internal/analyzers"
	"github.com/cameron/code-evaluator/internal/config"
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
  codeeval check    --repo PATH --base REV (--target REV | --working-tree) [--path P ...] [--fail-on new-debt|net-debt|gates|none] [--debt-budget N] [--top N]
                                                                                    (PR check: factors introduced/worsened; exit 3 on failure)
  codeeval graph    --repo PATH --rev REV --output FILE.html [--json FILE] [--max-callables N] [--config FILE]   (POC debt graph view)
  codeeval defects  --repo PATH (--rev REV | --base REV (--target REV | --working-tree)) [--rule ID ...] [--path P ...] [--fail-on gate|any|none] [--top N]   (defect rules, rules/v1)
  codeeval analyze  --repo PATH (--rev REV | --base REV (--target REV | --working-tree)) [--path P ...] [--debt-fail-on ...] [--debt-budget N] [--defects-fail-on ...]   (v2 debt + defects in one report)
  codeeval serve    [PATH] [--last N | --since DATE --until DATE] [--base REV] [--target REV] [--max N] [--addr 127.0.0.1:8765]   (web dashboard of the checkout's history: quality over time, graph per commit, provenance)
  codeeval hook     install|uninstall --tool claude-code [--bin PATH] [--config FILE] [--dry-run]   (add the analyze hooks to the project in the current directory)
  codeeval rules                                                          (defect rules × languages)
  codeeval attribution --repo PATH --base REV --target REV [--group-by model,harness,workflow] [--defects] [--infer-co-authors=false]   (step deltas by commit provenance)
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
	case "hook":
		return hookCmd(args[1:], stdout, stderr)
	case "serve":
		return serveCmd(args[1:], stdout, stderr)
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
	var k conn
	var o out
	k.register(fs)
	o.register(fs)
	base := fs.String("base", "", "base revision")
	target := fs.String("target", "", "target revision")
	var worktree = new(bool)
	k.worktreeFlags(fs, worktree, "evaluate working-tree changes against base")
	fs.BoolVar(&k.timing, "timing", false, "include timing metadata")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if msg := evaluateArgs(fs, *base, *target, *worktree, history); msg != "" {
		fmt.Fprintln(stderr, msg)
		return 2
	}
	if !o.valid(stderr) {
		return 2
	}
	ev, _, code := k.open(stderr, o.output, o.summary)
	if ev == nil {
		return code
	}
	defer ev.Close()
	var payload any
	var md bytes.Buffer
	if history {
		h, err := ev.History(*base, *target, func(s string) { fmt.Fprintln(stderr, "codeeval:", s) })
		if err != nil {
			return fail(stderr, err)
		}
		payload = h
		report.WriteHistoryMarkdown(&md, h)
		md.WriteString("\n---\n\n")
		report.WriteMarkdown(&md, h.Net)
	} else {
		r, err := ev.Evaluate(*base, *target, *worktree)
		if err != nil {
			return fail(stderr, err)
		}
		payload = r
		report.WriteMarkdown(&md, r)
	}
	return o.emit(payload, md.Bytes(), stdout, stderr)
}

// evaluateArgs returns a usage error for evaluate/history, or "".
func evaluateArgs(fs *flag.FlagSet, base, target string, worktree, history bool) string {
	switch {
	case fs.NArg() > 0:
		return fmt.Sprintf("unexpected arguments: %v", fs.Args())
	case base == "":
		return "--base is required"
	case history && (target == "" || worktree):
		return "history needs --target and does not support --working-tree"
	case !history && (target == "") == !worktree:
		return "exactly one of --target or --working-tree is required"
	}
	return ""
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
