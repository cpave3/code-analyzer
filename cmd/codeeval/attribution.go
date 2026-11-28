package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/cameron/code-evaluator/internal/attribution"
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
	groupBy := fs.String("group-by", "agent,harness", "provenance fields to group by: agent, harness, or any Agent-<Key> metadata key (e.g. workflow)")
	withDefects := fs.Bool("defects", false, "also count defect findings each step introduced and resolved (runs rules on every step)")
	infer := fs.Bool("infer-co-authors", true, "infer the agent from Co-Authored-By lines naming a known coding agent when a commit has no provenance")
	notesRef := fs.String("notes-ref", provenance.NotesRef, "notes ref to read provenance notes from")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *base == "" || *target == "" || fs.NArg() > 0 {
		fmt.Fprintln(stderr, "attribution needs --base and --target")
		return 2
	}
	var group []string
	for _, g := range strings.Split(*groupBy, ",") {
		if g = strings.ToLower(strings.TrimSpace(g)); g != "" {
			group = append(group, g)
		}
	}
	if len(c.paths) > 0 {
		fmt.Fprintln(stderr, "attribution measures whole-codebase step deltas; --path is not supported")
		return 2
	}
	if len(group) == 0 {
		fmt.Fprintln(stderr, "--group-by needs at least one field")
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	bsha, bsnap, err := ev.ResolveSnapshot(*base)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	tsha, tsnap, err := ev.ResolveSnapshot(*target)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	pairs, err := ev.FirstParentSteps(bsha, tsha)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	msgs, err := ev.CommitsBetween(bsha, tsha, *notesRef)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	records := map[string]provenance.Record{}
	for _, m := range msgs {
		rec := provenance.Parse(m.SHA, m.Message, m.Note, *infer)
		rec.Merge = len(m.Parents) > 1
		records[m.SHA] = rec
	}
	r := &attribution.Report{Schema: attribution.Schema, Tool: report.Tool{Name: "codeeval", Version: report.ToolVersion}, Scoring: scorev2.Version,
		ConfigHash: cfg.Hash(), Repository: ev.RepoName(), GroupBy: group, Inferred: *infer,
		Base:   report.Revision{Ref: *base, Commit: bsha, Tree: bsnap.Tree, CommitTime: ev.CommitTime(bsha)},
		Target: report.Revision{Ref: *target, Commit: tsha, Tree: tsnap.Tree, CommitTime: ev.CommitTime(tsha)}}
	langs := languages.Default()
	reg := all.Default(langs)
	if *withDefects {
		r.RulesVersion = rules.Version
	}
	w := scorev2.DefaultWeights()
	var prevRules *rules.Result
	prevSHA := ""
	for i, p := range pairs {
		fmt.Fprintf(stderr, "codeeval: attribution step %d/%d %.12s\n", i+1, len(pairs), p[1])
		_, ps, err := ev.ResolveSnapshot(p[0])
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		_, cs, err := ev.ResolveSnapshot(p[1])
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		set, err := ev.ChangeSet(ps, cs)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		pa, ca := scorev2.Analyze(ps, cfg, w), scorev2.Analyze(cs, cfg, w)
		st := scorev2.NewStep(i+1, p[0], p[1], ev.Subject(p[1]), ps, cs, pa, ca, set)
		shas, err := ev.RevList(p[0], p[1])
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		step := attribution.Step{Index: i + 1, Parent: p[0], Commit: p[1], Subject: st.Subject, Size: st.Size, DebtDelta: st.Ledger.Debt,
			Decisions: st.Ledger.Components.DecisionPoints}
		for _, s := range shas {
			rec, ok := records[s]
			if !ok {
				rec = provenance.Record{Commit: s, Source: provenance.SourceNone}
			}
			step.Commits = append(step.Commits, rec)
		}
		step.Group = attribution.GroupKey(step.Commits, group)
		if *withDefects {
			opt := rules.Options{}
			before := prevRules
			if before == nil || prevSHA != p[0] {
				before = rules.Run(reg, ps, cfg, langs, opt)
			}
			after := rules.Run(reg, cs, cfg, langs, opt)
			in, out, _ := rules.Compare(before.Findings, after.Findings)
			d := &attribution.Defects{Resolved: len(out)}
			for _, f := range in {
				if f.Severity == rules.Gate {
					d.IntroducedGate++
				} else {
					d.IntroducedAdvisory++
				}
			}
			step.Defects = d
			prevRules, prevSHA = after, p[1]
		}
		r.Steps = append(r.Steps, step)
	}
	r.Finish()
	var md bytes.Buffer
	attribution.WriteMarkdown(&md, r)
	return c.emit(r, md.Bytes(), stdout, stderr)
}
