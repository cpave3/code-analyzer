package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// scoreCmd runs the experimental scoring/v2 model: an absolute snapshot
// (--rev) or a delta (--base/--target, optionally per first-parent step).
func scoreCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("score", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var k conn
	var o out
	k.register(fs)
	o.register(fs)
	rev := fs.String("rev", "", "revision to score (absolute snapshot)")
	base := fs.String("base", "", "base revision (delta)")
	target := fs.String("target", "", "target revision (delta)")
	steps := fs.Bool("steps", false, "add per first-parent step ledgers to a delta")
	modules := fs.Int("modules", 25, "modules shown in Markdown")
	top := fs.Int("top", 25, "callables and clone groups listed (0 = all)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if msg := scoreArgs(fs, *rev, *base, *target, *steps); msg != "" {
		fmt.Fprintln(stderr, msg)
		return 2
	}
	if !o.valid(stderr) {
		return 2
	}
	ev, cfg, code := k.open(stderr, o.output, o.summary)
	if ev == nil {
		return code
	}
	defer ev.Close()
	var md bytes.Buffer
	if *rev != "" {
		s, err := snapshotOf(ev, *rev)
		if err != nil {
			return fail(stderr, err)
		}
		r := scorev2.NewSnapshotReport(ev.RepoName(), cfg.Hash(), s.rev, scorev2.Analyze(s.snap, cfg, scorev2.DefaultWeights()), *top)
		scorev2.WriteSnapshotMarkdown(&md, r, *modules)
		return o.emit(r, md.Bytes(), stdout, stderr)
	}
	r, err := scoreDelta(ev, cfg, *base, *target, *steps, stderr)
	if err != nil {
		return fail(stderr, err)
	}
	scorev2.WriteDeltaMarkdown(&md, r, *modules)
	return o.emit(r, md.Bytes(), stdout, stderr)
}

func scoreArgs(fs *flag.FlagSet, rev, base, target string, steps bool) string {
	switch {
	case fs.NArg() > 0:
		return fmt.Sprintf("unexpected arguments: %v", fs.Args())
	case (rev != "") == (base != "" || target != "") || (rev == "" && (base == "" || target == "")):
		return "score needs either --rev, or both --base and --target"
	case rev != "" && steps:
		return "--steps needs --base and --target"
	}
	return ""
}

// scoreDelta builds the v2 delta report, with per first-parent step
// ledgers when steps is set.
func scoreDelta(ev *engine.Evaluator, cfg *config.Config, base, target string, steps bool, stderr io.Writer) (*scorev2.DeltaReport, error) {
	w := scorev2.DefaultWeights()
	b, t, err := changePair(ev, base, target, false)
	if err != nil {
		return nil, err
	}
	cs, err := ev.ChangeSet(b.snap, t.snap)
	if err != nil {
		return nil, err
	}
	r := scorev2.NewDeltaReport(ev.RepoName(), cfg.Hash(), b.rev, t.rev, b.snap, t.snap, scorev2.Analyze(b.snap, cfg, w), scorev2.Analyze(t.snap, cfg, w), cs)
	if !steps {
		return r, nil
	}
	pairs, err := ev.FirstParentSteps(b.rev.Commit, t.rev.Commit)
	if err != nil {
		return nil, err
	}
	for i, p := range pairs {
		fmt.Fprintf(stderr, "codeeval: v2 step %d/%d %.12s\n", i+1, len(pairs), p[1])
		st, err := v2Step(ev, cfg, i+1, p)
		if err != nil {
			return nil, err
		}
		r.Steps = append(r.Steps, st)
	}
	scorev2.CheckSteps(r, len(pairs) > 0 && pairs[0][0] == b.rev.Commit)
	return r, nil
}

// v2Step measures one first-parent step (parent, commit).
func v2Step(ev *engine.Evaluator, cfg *config.Config, index int, p [2]string) (scorev2.Step, error) {
	w := scorev2.DefaultWeights()
	_, ps, err := ev.ResolveSnapshot(p[0])
	if err != nil {
		return scorev2.Step{}, err
	}
	_, cs, err := ev.ResolveSnapshot(p[1])
	if err != nil {
		return scorev2.Step{}, err
	}
	set, err := ev.ChangeSet(ps, cs)
	if err != nil {
		return scorev2.Step{}, err
	}
	return scorev2.NewStep(index, p[0], p[1], ev.Subject(p[1]), ps, cs, scorev2.Analyze(ps, cfg, w), scorev2.Analyze(cs, cfg, w), set), nil
}
