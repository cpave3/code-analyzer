package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/languages"
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
	var worktree = new(bool)
	c.worktreeFlags(fs, worktree, "change mode: check working-tree changes against base")
	var only pathList
	fs.Var(&only, "rule", "rule ID or category to run (repeatable, comma-separated); default all")
	failOn := fs.String("fail-on", "", "exit 3 when: gate (gate findings; in change mode introduced ones), any, none. Default: gate in change mode, none in snapshot mode")
	top := fs.Int("top", 0, "findings listed (0 = all)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	change, ok := modeArgs(fs, *rev, *base, *target, *worktree)
	if !ok {
		fmt.Fprintln(stderr, "defects needs --rev R, or --base B with exactly one of --target T or --working-tree")
		return 2
	}
	if !defectsFailOn(failOn, change, "--fail-on", stderr) {
		return 2
	}
	ev, cfg, code := c.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	langs := languages.Default()
	reg := all.Default(langs)
	if !rulesKnown(reg, cfg, only, stderr) {
		return 2
	}
	opt := rules.Options{Only: only, InScope: scorev2.Scope{Paths: c.paths}.Match}
	var r *rules.Report
	if !change {
		s, err := snapshotOf(ev, *rev)
		if err != nil {
			return fail(stderr, err)
		}
		r = rules.NewReport(reg, nil, rules.Run(reg, s.snap, cfg, langs, opt), *failOn)
		r.Revision = &s.rev
	} else {
		b, t, err := changePair(ev, *base, *target, *worktree)
		if err != nil {
			return fail(stderr, err)
		}
		r = rules.NewReport(reg, rules.Run(reg, b.snap, cfg, langs, opt), rules.Run(reg, t.snap, cfg, langs, opt), *failOn)
		r.Base, r.Target = &b.rev, &t.rev
	}
	r.Repository, r.ConfigHash, r.Selected, r.Scope = ev.RepoName(), cfg.Hash(), only, c.scopeList()
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

// defectsFailOn applies the mode default and validates a defects fail-on flag.
func defectsFailOn(v *string, change bool, flagName string, stderr io.Writer) bool {
	if *v == "" {
		*v = rules.FailNone
		if change {
			*v = rules.FailGate
		}
	}
	if *v != rules.FailGate && *v != rules.FailAny && *v != rules.FailNone {
		fmt.Fprintln(stderr, flagName+" must be gate, any or none")
		return false
	}
	return true
}

// rulesKnown checks --rule selections and config defects.disable entries.
func rulesKnown(reg *rules.Registry, cfg *config.Config, only []string, stderr io.Writer) bool {
	for _, id := range append(append([]string{}, only...), cfg.Defects.Disable...) {
		if !knownRule(reg, id) {
			fmt.Fprintf(stderr, "unknown rule or category %q (see `codeeval rules`)\n", id)
			return false
		}
	}
	return true
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
