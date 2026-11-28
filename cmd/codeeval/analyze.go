package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"

	"github.com/cameron/code-evaluator/internal/clones"
	"github.com/cameron/code-evaluator/internal/combined"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/all"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// analyzeCmd runs v2 debt, defect rules and near-duplication together.
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
	change, ok := modeArgs(fs, *rev, *base, *target, *worktree)
	if !ok {
		fmt.Fprintln(stderr, "analyze needs --rev R, or --base B with exactly one of --target T or --working-tree")
		return 2
	}
	if msg := analyzeDebtArgs(*debtFail, change); msg != "" {
		fmt.Fprintln(stderr, msg)
		return 2
	}
	if !defectsFailOn(defFail, change, "--defects-fail-on", stderr) {
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	a := analysis{ev: ev, cfg: cfg, langs: languages.Default(), scope: scorev2.Scope{Paths: c.paths}, defFail: *defFail, top: *top}
	a.reg = all.Default(a.langs)
	if !rulesKnown(a.reg, cfg, nil, stderr) {
		return 2
	}
	var out *combined.Report
	var err error
	if change {
		out, err = a.change(*base, *target, *worktree, *debtFail)
	} else {
		out, err = a.snapshot(*rev)
	}
	if err != nil {
		return fail(stderr, err)
	}
	out.Defects.Repository, out.Defects.ConfigHash, out.Defects.Scope = ev.RepoName(), cfg.Hash(), c.scopeList()
	var md bytes.Buffer
	combined.WriteMarkdown(&md, out, *files, *detailed, *top)
	if code := c.emit(out, md.Bytes(), stdout, stderr); code != 0 {
		return code
	}
	if out.Failed {
		return 3
	}
	return 0
}

func analyzeDebtArgs(debtFail string, change bool) string {
	switch {
	case !checkFailOn[debtFail]:
		return "--debt-fail-on must be new-debt, net-debt, gates or none"
	case !change && debtFail != "none":
		return "--debt-fail-on applies to change mode only"
	}
	return ""
}

type analysis struct {
	ev      *engine.Evaluator
	cfg     *config.Config
	langs   *lang.Registry
	reg     *rules.Registry
	scope   scorev2.Scope
	defFail string
	top     int
}

func (a analysis) opt() rules.Options { return rules.Options{InScope: a.scope.Match} }

func (a analysis) snapshot(rev string) (*combined.Report, error) {
	s, err := snapshotOf(a.ev, rev)
	if err != nil {
		return nil, err
	}
	ex := scorev2.NewExplainReport(a.ev.RepoName(), a.cfg.Hash(), s.rev, s.snap, scorev2.Analyze(s.snap, a.cfg, scorev2.DefaultWeights()), a.scope, 0)
	d := rules.NewReport(a.reg, nil, rules.Run(a.reg, s.snap, a.cfg, a.langs, a.opt()), a.defFail)
	d.Revision = &s.rev
	out := combined.Snapshot(ex, d)
	out.AddClones(clones.Analyze(s.snap, a.cfg, a.langs, a.scope.Match))
	if a.top > 0 && len(ex.Factors) > a.top {
		ex.Omitted += len(ex.Factors) - a.top
		ex.Factors = ex.Factors[:a.top]
	}
	d.Limit(a.top)
	return out, nil
}

func (a analysis) change(base, target string, worktree bool, debtFail string) (*combined.Report, error) {
	b, t, err := changePair(a.ev, base, target, worktree)
	if err != nil {
		return nil, err
	}
	ck, err := checkReport(a.ev, a.cfg, b, t, a.scope, debtFail)
	if err != nil {
		return nil, err
	}
	d := rules.NewReport(a.reg, rules.Run(a.reg, b.snap, a.cfg, a.langs, a.opt()), rules.Run(a.reg, t.snap, a.cfg, a.langs, a.opt()), a.defFail)
	d.Base, d.Target = &b.rev, &t.rev
	out := combined.Change(ck, d)
	bc, tc := clones.Analyze(b.snap, a.cfg, a.langs, a.scope.Match), clones.Analyze(t.snap, a.cfg, a.langs, a.scope.Match)
	tc.NearMisses, tc.Fragments = clones.Introduced(bc.NearMisses, tc.NearMisses), clones.Introduced(bc.Fragments, tc.Fragments)
	out.AddClones(tc)
	ck.Limit(a.top)
	d.Limit(a.top)
	return out, nil
}
