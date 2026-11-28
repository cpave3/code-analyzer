package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"

	"github.com/cameron/code-evaluator/internal/combined"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/all"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// analyzeCmd runs v2 debt and defect rules together.
func analyzeCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var c common
	c.register(fs)
	rev := fs.String("rev", "", "revision to analyse (snapshot mode)")
	base := fs.String("base", "", "base revision (change mode)")
	target := fs.String("target", "", "target revision (change mode)")
	worktree := fs.Bool("working-tree", false, "change mode: tracked working-tree changes against base")
	debtFail := fs.String("debt-fail-on", "none", "change mode: exit 3 when the v2 check fails: new-debt, net-debt, gates, none")
	defFail := fs.String("defects-fail-on", "", "exit 3 on defects: gate, any, none (default: gate in change mode, none in snapshot mode)")
	files := fs.Int("files", 25, "rows in the per-file table (0 = all)")
	top := fs.Int("top", 50, "debt factors and defect findings listed per section (0 = all)")
	detailed := fs.Int("detailed", 5, "debt factors with full remediation text (snapshot mode)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	change := *base != ""
	if fs.NArg() > 0 || (*rev != "") == change || (change && (*target == "") == !*worktree) || (!change && (*target != "" || *worktree)) {
		fmt.Fprintln(stderr, "analyze needs --rev R, or --base B with exactly one of --target T or --working-tree")
		return 2
	}
	switch *debtFail {
	case "new-debt", "net-debt", "gates", "none":
	default:
		fmt.Fprintln(stderr, "--debt-fail-on must be new-debt, net-debt, gates or none")
		return 2
	}
	if !change && *debtFail != "none" {
		fmt.Fprintln(stderr, "--debt-fail-on applies to change mode only")
		return 2
	}
	if *defFail == "" {
		*defFail = rules.FailNone
		if change {
			*defFail = rules.FailGate
		}
	}
	if *defFail != rules.FailGate && *defFail != rules.FailAny && *defFail != rules.FailNone {
		fmt.Fprintln(stderr, "--defects-fail-on must be gate, any or none")
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	langs := languages.Default()
	reg := all.Default(langs)
	for _, id := range cfg.Defects.Disable {
		if !knownRule(reg, id) {
			fmt.Fprintf(stderr, "unknown rule or category %q in config defects.disable\n", id)
			return 2
		}
	}
	scope := scorev2.Scope{Paths: c.paths}
	opt := rules.Options{InScope: scope.Match}
	w := scorev2.DefaultWeights()
	var out *combined.Report
	if !change {
		sha, snap, err := ev.ResolveSnapshot(*rev)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		rv := report.Revision{Ref: *rev, Commit: sha, Tree: snap.Tree, CommitTime: ev.CommitTime(sha)}
		ex := scorev2.NewExplainReport(ev.RepoName(), cfg.Hash(), rv, snap, scorev2.Analyze(snap, cfg, w), scope, 0)
		d := rules.NewReport(reg, nil, rules.Run(reg, snap, cfg, langs, opt), *defFail)
		d.Revision, d.Repository, d.ConfigHash, d.Scope = &rv, ev.RepoName(), cfg.Hash(), scopeList(c.paths)
		out = combined.Snapshot(ex, d)
		if *top > 0 && len(ex.Factors) > *top {
			ex.Omitted += len(ex.Factors) - *top
			ex.Factors = ex.Factors[:*top]
		}
		d.Limit(*top)
	} else {
		bsha, bsnap, err := ev.ResolveSnapshot(*base)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		var tsnap *index.Snapshot
		trev := report.Revision{Ref: *target}
		if *worktree {
			tsnap, err = ev.WorktreeSnapshot()
			trev = report.Revision{Ref: "WORKTREE", Commit: "WORKTREE"}
		} else {
			var tsha string
			tsha, tsnap, err = ev.ResolveSnapshot(*target)
			trev.Commit, trev.CommitTime = tsha, ev.CommitTime(tsha)
		}
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		trev.Tree = tsnap.Tree
		brev := report.Revision{Ref: *base, Commit: bsha, Tree: bsnap.Tree, CommitTime: ev.CommitTime(bsha)}
		cs, err := ev.ChangeSet(bsnap, tsnap)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		ck := scorev2.NewCheckReport(ev.RepoName(), cfg.Hash(), brev, trev, bsnap, tsnap, scorev2.Analyze(bsnap, cfg, w), scorev2.Analyze(tsnap, cfg, w), cs, scope, *debtFail)
		d := rules.NewReport(reg, rules.Run(reg, bsnap, cfg, langs, opt), rules.Run(reg, tsnap, cfg, langs, opt), *defFail)
		d.Base, d.Target, d.Repository, d.ConfigHash, d.Scope = &brev, &trev, ev.RepoName(), cfg.Hash(), scopeList(c.paths)
		out = combined.Change(ck, d)
		ck.Limit(*top)
		d.Limit(*top)
	}
	var md bytes.Buffer
	combined.WriteMarkdown(&md, out, *files, *detailed)
	if code := c.emit(out, md.Bytes(), stdout, stderr); code != 0 {
		return code
	}
	if out.Failed {
		return 3
	}
	return 0
}

func scopeList(p pathList) []string {
	if p == nil {
		return []string{}
	}
	return p
}
