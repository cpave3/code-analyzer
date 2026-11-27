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
)

// scoreCmd runs the experimental scoring/v2 model: an absolute snapshot
// (--rev) or a delta (--base/--target, optionally per first-parent step).
func scoreCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("score", flag.ContinueOnError)
	fs.SetOutput(stderr)
	repo := fs.String("repo", ".", "repository path")
	rev := fs.String("rev", "", "revision to score (absolute snapshot)")
	base := fs.String("base", "", "base revision (delta)")
	target := fs.String("target", "", "target revision (delta)")
	steps := fs.Bool("steps", false, "add per first-parent step ledgers to a delta")
	cfgPath := fs.String("config", "", "config file")
	output := fs.String("output", "", "JSON report path")
	summary := fs.String("summary", "", "Markdown summary path")
	format := fs.String("format", "md", "stdout format: md, json, none")
	cacheDir := fs.String("cache-dir", "", "index cache directory")
	noCache := fs.Bool("no-cache", false, "in-memory index only")
	modules := fs.Int("modules", 25, "modules shown in Markdown")
	top := fs.Int("top", 25, "callables and clone groups listed (0 = all)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "unexpected arguments: %v\n", fs.Args())
		return 2
	}
	snapshotMode := *rev != ""
	if snapshotMode == (*base != "" || *target != "") || (!snapshotMode && (*base == "" || *target == "")) {
		fmt.Fprintln(stderr, "score needs either --rev, or both --base and --target")
		return 2
	}
	if snapshotMode && *steps {
		fmt.Fprintln(stderr, "--steps needs --base and --target")
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
	ev, err := engine.New(engine.Options{RepoPath: *repo, Config: cfg, CacheDir: dir, Langs: languages.Default(), Analyzers: analyzers.Default()})
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
	w := scorev2.DefaultWeights()
	revOf := func(ref, sha, tree string) report.Revision {
		return report.Revision{Ref: ref, Commit: sha, Tree: tree, CommitTime: ev.CommitTime(sha)}
	}

	var payload any
	var md bytes.Buffer
	if snapshotMode {
		sha, snap, err := ev.ResolveSnapshot(*rev)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		a := scorev2.Analyze(snap, cfg, w)
		r := scorev2.NewSnapshotReport(ev.RepoName(), cfg.Hash(), revOf(*rev, sha, snap.Tree), a, *top)
		payload = r
		scorev2.WriteSnapshotMarkdown(&md, r, *modules)
	} else {
		bsha, bsnap, err := ev.ResolveSnapshot(*base)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		before := scorev2.Analyze(bsnap, cfg, w)
		tsha, tsnap, err := ev.ResolveSnapshot(*target)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		after := scorev2.Analyze(tsnap, cfg, w)
		cs, err := ev.ChangeSet(bsnap, tsnap)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		r := scorev2.NewDeltaReport(ev.RepoName(), cfg.Hash(), revOf(*base, bsha, bsnap.Tree), revOf(*target, tsha, tsnap.Tree), bsnap, tsnap, before, after, cs)
		if *steps {
			pairs, err := ev.FirstParentSteps(bsha, tsha)
			if err != nil {
				fmt.Fprintln(stderr, "error:", err)
				return 1
			}
			for i, p := range pairs {
				fmt.Fprintf(stderr, "codeeval: v2 step %d/%d %.12s\n", i+1, len(pairs), p[1])
				_, ps, err := ev.ResolveSnapshot(p[0])
				if err != nil {
					fmt.Fprintln(stderr, "error:", err)
					return 1
				}
				pa := scorev2.Analyze(ps, cfg, w)
				_, cs2, err := ev.ResolveSnapshot(p[1])
				if err != nil {
					fmt.Fprintln(stderr, "error:", err)
					return 1
				}
				ca := scorev2.Analyze(cs2, cfg, w)
				set, err := ev.ChangeSet(ps, cs2)
				if err != nil {
					fmt.Fprintln(stderr, "error:", err)
					return 1
				}
				r.Steps = append(r.Steps, scorev2.NewStep(i+1, p[0], p[1], ev.Subject(p[1]), ps, cs2, pa, ca, set))
			}
			scorev2.CheckSteps(r, len(pairs) > 0 && pairs[0][0] == bsha)
		}
		payload = r
		scorev2.WriteDeltaMarkdown(&md, r, *modules)
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
