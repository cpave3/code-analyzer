package main

import (
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

// Shared command plumbing: repository/config/cache flags, report output
// flags, error reporting and revision resolution.

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

// conn holds the flags every command that opens a repository has.
type conn struct {
	repo, cfgPath, cacheDir string
	noCache, timing         bool
	untracked               bool
}

// worktreeFlags registers --working-tree and --untracked.
func (k *conn) worktreeFlags(fs *flag.FlagSet, worktree *bool, usage string) {
	fs.BoolVar(worktree, "working-tree", false, usage)
	fs.BoolVar(&k.untracked, "untracked", true, "with --working-tree: also read untracked files that are not ignored (false = tracked files only)")
}

func (k *conn) register(fs *flag.FlagSet) {
	fs.StringVar(&k.repo, "repo", ".", "repository path")
	fs.StringVar(&k.cfgPath, "config", "", "config file")
	fs.StringVar(&k.cacheDir, "cache-dir", "", "index cache directory")
	fs.BoolVar(&k.noCache, "no-cache", false, "in-memory index only")
}

// open loads the config and opens the evaluator, refusing to write any of
// outputs inside the evaluated repository.
func (k *conn) open(stderr io.Writer, outputs ...string) (*engine.Evaluator, *config.Config, int) {
	cfg, err := config.Load(k.cfgPath)
	if err != nil {
		return nil, nil, fail(stderr, err)
	}
	dir := k.cacheDir
	if k.noCache {
		dir = ""
	} else if dir == "" {
		if ucd, err := os.UserCacheDir(); err == nil {
			dir = filepath.Join(ucd, "codeeval")
		}
	}
	ev, err := engine.New(engine.Options{RepoPath: k.repo, Config: cfg, CacheDir: dir, Langs: languages.Default(), Analyzers: analyzers.Default(), Timing: k.timing, Untracked: k.untracked})
	if err != nil {
		return nil, nil, fail(stderr, err)
	}
	for _, p := range outputs {
		if p != "" && inside(ev.RepoPath(), p) {
			ev.Close()
			fmt.Fprintf(stderr, "error: refusing to write %s inside the evaluated repository\n", p)
			return nil, nil, 1
		}
	}
	return ev, cfg, 0
}

// out holds the report output flags.
type out struct{ output, summary, format string }

func (o *out) register(fs *flag.FlagSet) {
	fs.StringVar(&o.output, "output", "", "JSON report path")
	fs.StringVar(&o.summary, "summary", "", "Markdown report path")
	fs.StringVar(&o.format, "format", "md", "stdout format: md, json, none")
}

func (o *out) valid(stderr io.Writer) bool {
	if o.format != "md" && o.format != "json" && o.format != "none" {
		fmt.Fprintln(stderr, "--format must be md, json or none")
		return false
	}
	return true
}

// emit writes the JSON and Markdown reports to the requested places.
func (o *out) emit(payload any, md []byte, stdout, stderr io.Writer) int {
	js, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fail(stderr, err)
	}
	js = append(js, '\n')
	for _, w := range []struct {
		path string
		data []byte
	}{{o.output, js}, {o.summary, md}} {
		if w.path != "" {
			if err := writeFile(w.path, w.data); err != nil {
				return fail(stderr, err)
			}
		}
	}
	switch o.format {
	case "md":
		stdout.Write(md)
	case "json":
		stdout.Write(js)
	}
	return 0
}

// common is conn + out + a --path scope, for the v2 and rules commands.
type common struct {
	conn
	out
	paths pathList
}

func (c *common) register(fs *flag.FlagSet) {
	c.conn.register(fs)
	c.out.register(fs)
	fs.Var(&c.paths, "path", "scope: a path prefix (module, application folder) or glob; repeatable or comma-separated")
}

func (c *common) open(stderr io.Writer) (*engine.Evaluator, *config.Config, int) {
	if !c.valid(stderr) {
		return nil, nil, 2
	}
	return c.conn.open(stderr, c.output, c.summary)
}

func (c *common) scopeList() []string {
	if c.paths == nil {
		return []string{}
	}
	return c.paths
}

// fail reports an error and returns exit code 1.
func fail(stderr io.Writer, err error) int {
	fmt.Fprintln(stderr, "error:", err)
	return 1
}

// side is a resolved revision with its snapshot.
type side struct {
	rev  report.Revision
	snap *index.Snapshot
}

func snapshotOf(ev *engine.Evaluator, ref string) (side, error) {
	sha, snap, err := ev.ResolveSnapshot(ref)
	if err != nil {
		return side{}, err
	}
	return side{rev: report.Revision{Ref: ref, Commit: sha, Tree: snap.Tree, CommitTime: ev.CommitTime(sha)}, snap: snap}, nil
}

// targetOf resolves a target revision, or the working tree.
func targetOf(ev *engine.Evaluator, ref string, worktree bool) (side, error) {
	if !worktree {
		return snapshotOf(ev, ref)
	}
	snap, err := ev.WorktreeSnapshot()
	if err != nil {
		return side{}, err
	}
	return side{rev: report.Revision{Ref: "WORKTREE", Commit: "WORKTREE", Tree: snap.Tree}, snap: snap}, nil
}

// changePair resolves base and target (or the working tree).
func changePair(ev *engine.Evaluator, base, target string, worktree bool) (side, side, error) {
	b, err := snapshotOf(ev, base)
	if err != nil {
		return side{}, side{}, err
	}
	t, err := targetOf(ev, target, worktree)
	return b, t, err
}

// modeArgs validates "--rev R, or --base B with exactly one of --target T
// or --working-tree" and reports whether change mode was chosen.
func modeArgs(fs *flag.FlagSet, rev, base, target string, worktree bool) (change, ok bool) {
	change = base != ""
	ok = fs.NArg() == 0 && (rev != "") != change && !(change && (target == "") == !worktree) && !(!change && (target != "" || worktree))
	return change, ok
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

// checkReport builds the v2 check report for a resolved change.
func checkReport(ev *engine.Evaluator, cfg *config.Config, b, t side, scope scorev2.Scope, failOn string) (*scorev2.CheckReport, error) {
	w := scorev2.DefaultWeights()
	cs, err := ev.ChangeSet(b.snap, t.snap)
	if err != nil {
		return nil, err
	}
	return scorev2.NewCheckReport(ev.RepoName(), cfg.Hash(), b.rev, t.rev, b.snap, t.snap, scorev2.Analyze(b.snap, cfg, w), scorev2.Analyze(t.snap, cfg, w), cs, scope, failOn), nil
}
