package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/cameron/code-evaluator/internal/attribution"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/provenance"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/all"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// attributionCmd groups first-parent step deltas by commit provenance.
func attributionCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("attribution", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var c common
	c.register(fs)
	base := fs.String("base", "", "base revision (exclusive)")
	target := fs.String("target", "", "target revision")
	groupBy := fs.String("group-by", "model,harness", "provenance fields to group by: model, harness, or any Agent-<Key> metadata key (e.g. workflow)")
	withDefects := fs.Bool("defects", false, "also count defect findings each step introduced and resolved (runs rules on every step)")
	infer := fs.Bool("infer-co-authors", true, "infer the agent from Co-Authored-By lines naming a known coding agent when a commit has no provenance")
	notesRef := fs.String("notes-ref", provenance.NotesRef, "notes ref to read provenance notes from")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	group := splitFields(*groupBy)
	if msg := attributionArgs(fs, *base, *target, c.paths, group); msg != "" {
		fmt.Fprintln(stderr, msg)
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	r, err := attribute(ev, cfg, *base, *target, *notesRef, group, *infer, *withDefects, stderr)
	if err != nil {
		return fail(stderr, err)
	}
	var md bytes.Buffer
	attribution.WriteMarkdown(&md, r)
	return c.emit(r, md.Bytes(), stdout, stderr)
}

func splitFields(s string) []string {
	var out []string
	for _, g := range strings.Split(s, ",") {
		if g = strings.ToLower(strings.TrimSpace(g)); g != "" {
			out = append(out, g)
		}
	}
	return out
}

func attributionArgs(fs *flag.FlagSet, base, target string, paths pathList, group []string) string {
	switch {
	case base == "" || target == "" || fs.NArg() > 0:
		return "attribution needs --base and --target"
	case len(paths) > 0:
		return "attribution measures whole-codebase step deltas; --path is not supported"
	case len(group) == 0:
		return "--group-by needs at least one field"
	}
	return ""
}

// attributer measures steps and keeps the previous step's rule results so
// each snapshot's rules run once.
type attributer struct {
	ev      *engine.Evaluator
	cfg     *config.Config
	langs   *lang.Registry
	reg     *rules.Registry
	records map[string]provenance.Record
	group   []string
	defects bool
	prev    *rules.Result
	prevSHA string
}

func attribute(ev *engine.Evaluator, cfg *config.Config, base, target, notesRef string, group []string, infer, defects bool, stderr io.Writer) (*attribution.Report, error) {
	b, err := snapshotOf(ev, base)
	if err != nil {
		return nil, err
	}
	t, err := snapshotOf(ev, target)
	if err != nil {
		return nil, err
	}
	pairs, err := ev.FirstParentSteps(b.rev.Commit, t.rev.Commit)
	if err != nil {
		return nil, err
	}
	msgs, err := ev.CommitsBetween(b.rev.Commit, t.rev.Commit, notesRef)
	if err != nil {
		return nil, err
	}
	at := &attributer{ev: ev, cfg: cfg, langs: languages.Default(), records: map[string]provenance.Record{}, group: group, defects: defects}
	at.reg = all.Default(at.langs)
	for _, m := range msgs {
		rec := provenance.Parse(m.SHA, m.Message, m.Note, infer)
		rec.Merge = len(m.Parents) > 1
		at.records[m.SHA] = rec
	}
	r := &attribution.Report{Schema: attribution.Schema, Tool: report.Tool{Name: "codeeval", Version: report.ToolVersion}, Scoring: scorev2.Version,
		ConfigHash: cfg.Hash(), Repository: ev.RepoName(), GroupBy: group, Inferred: infer, Base: b.rev, Target: t.rev}
	if defects {
		r.RulesVersion = rules.Version
	}
	for i, p := range pairs {
		fmt.Fprintf(stderr, "codeeval: attribution step %d/%d %.12s\n", i+1, len(pairs), p[1])
		step, err := at.step(i+1, p)
		if err != nil {
			return nil, err
		}
		r.Steps = append(r.Steps, step)
	}
	r.Finish()
	return r, nil
}

func (at *attributer) step(index int, p [2]string) (attribution.Step, error) {
	st, err := v2Step(at.ev, at.cfg, index, p)
	if err != nil {
		return attribution.Step{}, err
	}
	shas, err := at.ev.RevList(p[0], p[1])
	if err != nil {
		return attribution.Step{}, err
	}
	step := attribution.Step{Index: index, Parent: p[0], Commit: p[1], Subject: st.Subject, Size: st.Size, DebtDelta: st.Ledger.Debt,
		Decisions: st.Ledger.Components.DecisionPoints}
	for _, s := range shas {
		rec, ok := at.records[s]
		if !ok {
			rec = provenance.Record{Commit: s, Source: provenance.SourceNone}
		}
		step.Commits = append(step.Commits, rec)
	}
	step.Group = attribution.GroupKey(step.Commits, at.group)
	if at.defects {
		d, err := at.defectDelta(p)
		if err != nil {
			return attribution.Step{}, err
		}
		step.Defects = d
	}
	return step, nil
}

func (at *attributer) defectDelta(p [2]string) (*attribution.Defects, error) {
	before := at.prev
	if before == nil || at.prevSHA != p[0] {
		_, ps, err := at.ev.ResolveSnapshot(p[0])
		if err != nil {
			return nil, err
		}
		before = rules.Run(at.reg, ps, at.cfg, at.langs, rules.Options{})
	}
	_, cs, err := at.ev.ResolveSnapshot(p[1])
	if err != nil {
		return nil, err
	}
	after := rules.Run(at.reg, cs, at.cfg, at.langs, rules.Options{})
	in, out, _ := rules.Compare(before.Findings, after.Findings)
	d := &attribution.Defects{Resolved: len(out)}
	for _, f := range in {
		if f.Severity == rules.Gate {
			d.IntroducedGate++
		} else {
			d.IntroducedAdvisory++
		}
	}
	at.prev, at.prevSHA = after, p[1]
	return d, nil
}
