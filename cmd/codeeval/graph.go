package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/cameron/code-evaluator/internal/scorev2"
	"github.com/cameron/code-evaluator/internal/viz"
)

// graphCmd writes the POC interactive debt graph for one revision.
func graphCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("graph", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var k conn
	k.register(fs)
	rev := fs.String("rev", "", "revision to visualise")
	output := fs.String("output", "", "HTML output path")
	jsonOut := fs.String("json", "", "optional graph data JSON path")
	maxCallables := fs.Int("max-callables", 0, "keep only the N callables with the largest debt (0 = all)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *rev == "" || *output == "" || fs.NArg() > 0 {
		fmt.Fprintln(stderr, "graph needs --rev and --output")
		return 2
	}
	ev, cfg, code := k.open(stderr, *output, *jsonOut)
	if ev == nil {
		return code
	}
	defer ev.Close()
	s, err := snapshotOf(ev, *rev)
	if err != nil {
		return fail(stderr, err)
	}
	d := viz.Build(ev.RepoName(), s.rev, s.snap, cfg, scorev2.Analyze(s.snap, cfg, scorev2.DefaultWeights()))
	d.Trim(*maxCallables)
	var page bytes.Buffer
	if err := viz.WriteHTML(&page, d); err != nil {
		return fail(stderr, err)
	}
	if err := writeFile(*output, page.Bytes()); err != nil {
		return fail(stderr, err)
	}
	if *jsonOut != "" {
		js, _ := json.MarshalIndent(d, "", "  ")
		if err := writeFile(*jsonOut, append(js, '\n')); err != nil {
			return fail(stderr, err)
		}
	}
	fmt.Fprintf(stdout, "wrote %s (%d modules, %d callables, %d clone groups, %d import edges)\n", *output, len(d.Modules), len(d.Callables), len(d.Clones), len(d.Imports))
	return 0
}
