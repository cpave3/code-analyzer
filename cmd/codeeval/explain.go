package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cameron/code-evaluator/internal/analyzers"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// pathList is a repeatable, comma-separated --path flag.
type pathList []string

func (p *pathList) String() string { return strings.Join(*p, ",") }
func (p *pathList) Set(v string) error {
	for _, s := range strings.Split(v, ",") {
		if s = strings.TrimSpace(s); s != "" {
			*p = append(*p, s)
		}
	}
	return nil
}

type common struct {
	repo, cfgPath, output, summary, format, cacheDir string
	noCache                                          bool
	paths                                            pathList
}

func (c *common) register(fs *flag.FlagSet) {
	fs.StringVar(&c.repo, "repo", ".", "repository path")
	fs.StringVar(&c.cfgPath, "config", "", "config file")
	fs.StringVar(&c.output, "output", "", "JSON report path")
	fs.StringVar(&c.summary, "summary", "", "Markdown report path")
	fs.StringVar(&c.format, "format", "md", "stdout format: md, json, none")
	fs.StringVar(&c.cacheDir, "cache-dir", "", "index cache directory")
	fs.BoolVar(&c.noCache, "no-cache", false, "in-memory index only")
	fs.Var(&c.paths, "path", "scope: a path prefix (module, application folder) or glob; repeatable or comma-separated")
}

func (c *common) open(stderr io.Writer) (*engine.Evaluator, *config.Config, int) {
	if c.format != "md" && c.format != "json" && c.format != "none" {
		fmt.Fprintln(stderr, "--format must be md, json or none")
		return nil, nil, 2
	}
	cfg, err := config.Load(c.cfgPath)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return nil, nil, 1
	}
	dir := c.cacheDir
	if c.noCache {
		dir = ""
	} else if dir == "" {
		if ucd, err := os.UserCacheDir(); err == nil {
			dir = filepath.Join(ucd, "codeeval")
		}
	}
	ev, err := engine.New(engine.Options{RepoPath: c.repo, Config: cfg, CacheDir: dir, Langs: languages.Default(), Analyzers: analyzers.Default()})
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return nil, nil, 1
	}
	for _, p := range []string{c.output, c.summary} {
		if p != "" && inside(ev.RepoPath(), p) {
			ev.Close()
			fmt.Fprintf(stderr, "error: refusing to write %s inside the evaluated repository\n", p)
			return nil, nil, 1
		}
	}
	return ev, cfg, 0
}

func (c *common) emit(payload any, md []byte, stdout, stderr io.Writer) int {
	js, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	js = append(js, '\n')
	if c.output != "" {
		if err := writeFile(c.output, js); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}
	if c.summary != "" {
		if err := writeFile(c.summary, md); err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
	}
	switch c.format {
	case "md":
		stdout.Write(md)
	case "json":
		stdout.Write(js)
	}
	return 0
}

// explainCmd lists the debt factors of a module, application, path or glob.
func explainCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("explain", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var c common
	c.register(fs)
	rev := fs.String("rev", "", "revision to explain")
	top := fs.Int("top", 50, "factors listed (0 = all)")
	detailed := fs.Int("detailed", 15, "factors with full remediation text in Markdown (0 = all)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *rev == "" || fs.NArg() > 0 {
		fmt.Fprintln(stderr, "explain needs --rev [--path P ...]")
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	sha, snap, err := ev.ResolveSnapshot(*rev)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	a := scorev2.Analyze(snap, cfg, scorev2.DefaultWeights())
	r := scorev2.NewExplainReport(ev.RepoName(), cfg.Hash(), report.Revision{Ref: *rev, Commit: sha, Tree: snap.Tree, CommitTime: ev.CommitTime(sha)}, snap, a, scorev2.Scope{Paths: c.paths}, *top)
	if r.Totals.SourceFiles == 0 && len(c.paths) > 0 {
		fmt.Fprintf(stderr, "warning: no measured source files match %v\n", []string(c.paths))
	}
	var md bytes.Buffer
	scorev2.WriteExplainMarkdown(&md, r, *detailed)
	return c.emit(r, md.Bytes(), stdout, stderr)
}

// checkCmd reports debt factors a change introduced or worsened (PR check).
func checkCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var c common
	c.register(fs)
	base := fs.String("base", "", "base revision (e.g. the merge base with main)")
	target := fs.String("target", "", "target revision")
	worktree := fs.Bool("working-tree", false, "check tracked working-tree changes against base")
	failOn := fs.String("fail-on", "new-debt", "exit 3 when: new-debt (any factor introduced or worsened in scope), net-debt (debt in scope rises), gates (any advisory gate fails), none")
	top := fs.Int("top", 0, "factors listed per section (0 = all)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *base == "" || (*target == "") == !*worktree || fs.NArg() > 0 {
		fmt.Fprintln(stderr, "check needs --base and exactly one of --target or --working-tree")
		return 2
	}
	if *failOn != "new-debt" && *failOn != "net-debt" && *failOn != "gates" && *failOn != "none" {
		fmt.Fprintln(stderr, "--fail-on must be new-debt, net-debt, gates or none")
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	w := scorev2.DefaultWeights()
	bsha, bsnap, err := ev.ResolveSnapshot(*base)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	before := scorev2.Analyze(bsnap, cfg, w)
	var tsnap *index.Snapshot
	targetRev := report.Revision{Ref: *target}
	if *worktree {
		tsnap, err = ev.WorktreeSnapshot()
		targetRev = report.Revision{Ref: "WORKTREE", Commit: "WORKTREE"}
	} else {
		var tsha string
		tsha, tsnap, err = ev.ResolveSnapshot(*target)
		targetRev.Commit, targetRev.CommitTime = tsha, ev.CommitTime(tsha)
	}
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	targetRev.Tree = tsnap.Tree
	after := scorev2.Analyze(tsnap, cfg, w)
	cs, err := ev.ChangeSet(bsnap, tsnap)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	r := scorev2.NewCheckReport(ev.RepoName(), cfg.Hash(), report.Revision{Ref: *base, Commit: bsha, Tree: bsnap.Tree, CommitTime: ev.CommitTime(bsha)}, targetRev, bsnap, tsnap, before, after, cs, scorev2.Scope{Paths: c.paths}, *failOn)
	r.Limit(*top)
	var md bytes.Buffer
	scorev2.WriteCheckMarkdown(&md, r)
	if code := c.emit(r, md.Bytes(), stdout, stderr); code != 0 {
		return code
	}
	if r.Failed {
		return 3
	}
	return 0
}
