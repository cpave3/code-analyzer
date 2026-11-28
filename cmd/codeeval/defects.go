package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/all"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// defectsCmd runs the defect rules on a revision or a change.
func defectsCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("defects", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var c common
	c.register(fs)
	rev := fs.String("rev", "", "revision to check (snapshot mode)")
	base := fs.String("base", "", "base revision (change mode)")
	target := fs.String("target", "", "target revision (change mode)")
	worktree := fs.Bool("working-tree", false, "change mode: check tracked working-tree changes against base")
	var only pathList
	fs.Var(&only, "rule", "rule ID or category to run (repeatable, comma-separated); default all")
	failOn := fs.String("fail-on", "", "exit 3 when: gate (gate findings; in change mode introduced ones), any, none. Default: gate in change mode, none in snapshot mode")
	top := fs.Int("top", 0, "findings listed (0 = all)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	change := *base != ""
	if fs.NArg() > 0 || (*rev != "") == change || (change && (*target == "") == !*worktree) || (!change && (*target != "" || *worktree)) {
		fmt.Fprintln(stderr, "defects needs --rev R, or --base B with exactly one of --target T or --working-tree")
		return 2
	}
	if *failOn == "" {
		*failOn = rules.FailNone
		if change {
			*failOn = rules.FailGate
		}
	}
	if *failOn != rules.FailGate && *failOn != rules.FailAny && *failOn != rules.FailNone {
		fmt.Fprintln(stderr, "--fail-on must be gate, any or none")
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	langs := languages.Default()
	reg := all.Default(langs)
	for _, id := range append(append([]string{}, only...), cfg.Defects.Disable...) {
		if !knownRule(reg, id) {
			fmt.Fprintf(stderr, "unknown rule or category %q (see `codeeval rules`)\n", id)
			return 2
		}
	}
	scope := scorev2.Scope{Paths: c.paths}
	opt := rules.Options{Only: only, InScope: scope.Match}

	var r *rules.Report
	if !change {
		sha, snap, err := ev.ResolveSnapshot(*rev)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		res := rules.Run(reg, snap, cfg, langs, opt)
		r = rules.NewReport(reg, nil, res, *failOn)
		r.Revision = &report.Revision{Ref: *rev, Commit: sha, Tree: snap.Tree, CommitTime: ev.CommitTime(sha)}
	} else {
		bsha, bsnap, err := ev.ResolveSnapshot(*base)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		var tsnap *index.Snapshot
		trev := &report.Revision{Ref: *target}
		if *worktree {
			tsnap, err = ev.WorktreeSnapshot()
			trev = &report.Revision{Ref: "WORKTREE", Commit: "WORKTREE"}
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
		before := rules.Run(reg, bsnap, cfg, langs, opt)
		after := rules.Run(reg, tsnap, cfg, langs, opt)
		r = rules.NewReport(reg, before, after, *failOn)
		r.Base = &report.Revision{Ref: *base, Commit: bsha, Tree: bsnap.Tree, CommitTime: ev.CommitTime(bsha)}
		r.Target = trev
	}
	r.Repository, r.ConfigHash, r.Selected = ev.RepoName(), cfg.Hash(), only
	r.Scope = c.paths
	if r.Scope == nil {
		r.Scope = []string{}
	}
	r.Limit(*top)
	var md bytes.Buffer
	rules.WriteMarkdown(&md, r)
	if code := c.emit(r, md.Bytes(), stdout, stderr); code != 0 {
		return code
	}
	if r.Failed {
		return 3
	}
	return 0
}

func knownRule(reg *rules.Registry, id string) bool {
	for _, r := range reg.Rules() {
		if r.ID == id || strings.HasPrefix(r.ID, id+"/") {
			return true
		}
	}
	return false
}

func listRules(w io.Writer) int {
	langs := languages.Default()
	var ids []string
	for _, a := range langs.Adapters() {
		ids = append(ids, a.ID())
	}
	rules.WriteCatalog(w, all.Default(langs), ids)
	return 0
}
