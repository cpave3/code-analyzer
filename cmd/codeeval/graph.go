package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cameron/code-evaluator/internal/analyzers"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/scorev2"
	"github.com/cameron/code-evaluator/internal/viz"
)

// graphCmd writes the POC interactive debt graph for one revision.
func graphCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("graph", flag.ContinueOnError)
	fs.SetOutput(stderr)
	repo := fs.String("repo", ".", "repository path")
	rev := fs.String("rev", "", "revision to visualise")
	cfgPath := fs.String("config", "", "config file")
	output := fs.String("output", "", "HTML output path")
	jsonOut := fs.String("json", "", "optional graph data JSON path")
	cacheDir := fs.String("cache-dir", "", "index cache directory")
	noCache := fs.Bool("no-cache", false, "in-memory index only")
	maxCallables := fs.Int("max-callables", 0, "keep only the N callables with the largest debt (0 = all)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *rev == "" || *output == "" || fs.NArg() > 0 {
		fmt.Fprintln(stderr, "graph needs --rev and --output")
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
	ev, err := engine.New(engine.Options{RepoPath: *repo, Config: cfg, CacheDir: dir, Langs: languages.Default(), Analyzers: analyzers.Default()})
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	defer ev.Close()
	for _, p := range []string{*output, *jsonOut} {
		if p != "" && inside(ev.RepoPath(), p) {
			fmt.Fprintf(stderr, "error: refusing to write %s inside the evaluated repository\n", p)
			return 1
		}
	}
	sha, snap, err := ev.ResolveSnapshot(*rev)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	a := scorev2.Analyze(snap, cfg, scorev2.DefaultWeights())
	d := viz.Build(ev.RepoName(), report.Revision{Ref: *rev, Commit: sha, Tree: snap.Tree, CommitTime: ev.CommitTime(sha)}, snap, cfg, a)
	d.Trim(*maxCallables)
	var page bytes.Buffer
	if err := viz.WriteHTML(&page, d); err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	if err := writeFile(*output, page.Bytes()); err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	if *jsonOut != "" {
		js, _ := json.MarshalIndent(d, "", "  ")
		if err := writeFile(*jsonOut, append(js, '\n')); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}
	fmt.Fprintf(stdout, "wrote %s (%d modules, %d callables, %d clone groups, %d import edges)\n", *output, len(d.Modules), len(d.Callables), len(d.Clones), len(d.Imports))
	return 0
}
