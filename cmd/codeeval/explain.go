package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"

	"github.com/cameron/code-evaluator/internal/scorev2"
)

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
	s, err := snapshotOf(ev, *rev)
	if err != nil {
		return fail(stderr, err)
	}
	a := scorev2.Analyze(s.snap, cfg, scorev2.DefaultWeights())
	r := scorev2.NewExplainReport(ev.RepoName(), cfg.Hash(), s.rev, s.snap, a, scorev2.Scope{Paths: c.paths}, *top)
	if r.Totals.SourceFiles == 0 && len(c.paths) > 0 {
		fmt.Fprintf(stderr, "warning: no measured source files match %v\n", []string(c.paths))
	}
	var md bytes.Buffer
	scorev2.WriteExplainMarkdown(&md, r, *detailed)
	return c.emit(r, md.Bytes(), stdout, stderr)
}

var checkFailOn = map[string]bool{"new-debt": true, "net-debt": true, "gates": true, "none": true}

// checkCmd reports debt factors a change introduced or worsened (PR check).
func checkCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var c common
	c.register(fs)
	base := fs.String("base", "", "base revision (e.g. the merge base with main)")
	target := fs.String("target", "", "target revision")
	var worktree = new(bool)
	c.worktreeFlags(fs, worktree, "check working-tree changes against base")
	failOn := fs.String("fail-on", "new-debt", "exit 3 when: new-debt (any factor introduced or worsened in scope), net-debt (debt in scope rises), gates (any advisory gate fails), none")
	top := fs.Int("top", 0, "factors listed per section (0 = all)")
	budget := fs.Float64("debt-budget", 0, debtBudgetUsage)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *budget < 0 {
		fmt.Fprintln(stderr, "--debt-budget must not be negative")
		return 2
	}
	if *base == "" || (*target == "") == !*worktree || fs.NArg() > 0 {
		fmt.Fprintln(stderr, "check needs --base and exactly one of --target or --working-tree")
		return 2
	}
	if !checkFailOn[*failOn] {
		fmt.Fprintln(stderr, "--fail-on must be new-debt, net-debt, gates or none")
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	b, t, err := changePair(ev, *base, *target, *worktree)
	if err != nil {
		return fail(stderr, err)
	}
	r, err := checkReport(ev, cfg, b, t, scorev2.Scope{Paths: c.paths}, *failOn, *budget)
	if err != nil {
		return fail(stderr, err)
	}
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
