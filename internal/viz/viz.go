// Package viz renders a scoring/v2 snapshot as a self-contained, interactive
// HTML graph (POC): modules, debt-carrying callables, clone groups and import
// edges, in a force-directed "knowledge graph" view. The page has no external
// dependencies, so it works offline.
package viz

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

//go:embed graph.html
var page string

// Schema identifies the graph data format.
const Schema = "codeeval.graph/v2"

// Module is a module node.
type Module struct {
	ID          string       `json:"id"`
	Language    string       `json:"lang"`
	Group       string       `json:"group"` // top-level folder
	Lines       int          `json:"lines"`
	Callables   int          `json:"callables"`
	Hotspots    int          `json:"hotspots"`
	Debt        scorev2.Debt `json:"debt"`
	DebtPerKLOC float64      `json:"density"`
	Score       *float64     `json:"score"`
}

// Callable is a callable node (it carries complexity debt or is a clone member).
type Callable struct {
	ID         string  `json:"id"`
	Module     string  `json:"module"` // module node id
	Name       string  `json:"name"`
	File       string  `json:"file"`
	Line       int     `json:"line"`
	Cyclomatic int     `json:"cc"`
	Nesting    int     `json:"nest"`
	Lines      int     `json:"lines"`
	Debt       float64 `json:"debt"` // complexity + nesting debt
	CloneDebt  float64 `json:"cloneDebt"`
}

// Clone is a structural clone group.
type Clone struct {
	Exact   bool     `json:"exact"`
	Tokens  float64  `json:"tokens"`
	Debt    float64  `json:"debt"`
	Members []string `json:"members"` // callable ids
}

// Import is a module-to-module edge.
type Import struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Violation string `json:"violation,omitempty"` // rule name
	InCycle   bool   `json:"cycle,omitempty"`
}

// Data is the full graph.
type Data struct {
	Schema     string           `json:"schema"`
	Scoring    string           `json:"scoring"`
	Tool       string           `json:"tool"`
	Repository string           `json:"repository"`
	Revision   report.Revision  `json:"revision"`
	Weights    scorev2.Weights  `json:"weights"`
	Coverage   scorev2.Coverage `json:"coverage"`
	Codebase   scorev2.Unit     `json:"codebase"`
	Languages  []scorev2.Unit   `json:"languages"`
	Modules    []Module         `json:"modules"`
	Callables  []Callable       `json:"callables"`
	Clones     []Clone          `json:"clones"`
	Imports    []Import         `json:"imports"`
	Cycles     [][]string       `json:"cycles"`
	Notes      []string         `json:"notes"`
}

func modKey(lang, id string) string { return lang + ":" + id }

func group(id string) string {
	if i := strings.Index(id, "/"); i > 0 {
		return id[:i]
	}
	return id
}

func r2(v float64) float64 { return float64(int64(v*100+0.5)) / 100 }

// Build assembles graph data from a snapshot and its v2 analysis.
func Build(repo string, rev report.Revision, snap *index.Snapshot, cfg *config.Config, a *scorev2.Analysis) *Data {
	d := &Data{Schema: Schema, Scoring: scorev2.Version, Tool: report.ToolVersion, Repository: repo, Revision: rev,
		Weights: a.Weights, Coverage: a.Coverage, Codebase: a.Codebase, Languages: a.Languages, Cycles: a.Cycles,
		Notes: []string{
			"POC view of scoring/v2-experimental debt: node size and colour are uncalibrated structural proxies, not quality verdicts",
			"callables shown carry complexity/nesting debt (hotspots) or belong to a structural clone group; other callables are omitted",
			"TypeScript complexity includes inline callbacks and JSX conditionals; cross-language sizes are not directly comparable",
			"import edges are module (directory/package) level; architecture rules are only those configured",
		}}
	if d.Cycles == nil {
		d.Cycles = [][]string{}
	}
	langOf := map[string]string{}
	for _, m := range a.Modules {
		k := modKey(m.Language, m.ID)
		langOf[m.ID] = m.Language
		mm := Module{ID: k, Language: m.Language, Group: group(m.ID), Lines: m.Components.SourceLines, Callables: m.Components.Callables,
			Hotspots: m.Components.Hotspots, Debt: m.Debt, Score: m.Score}
		mm.Debt = scorev2.Debt{Complexity: r2(m.Debt.Complexity), Nesting: r2(m.Debt.Nesting), Duplication: r2(m.Debt.Duplication), Architecture: r2(m.Debt.Architecture), Total: r2(m.Debt.Total)}
		if m.DebtPerKLOC != nil {
			mm.DebtPerKLOC = *m.DebtPerKLOC
		}
		d.Modules = append(d.Modules, mm)
	}
	fileLang := func(file string) string {
		if f := snap.Files[file]; f != nil {
			return f.Language
		}
		return ""
	}
	calls := map[string]*Callable{}
	cid := func(file string, line int, name string) string { return fmt.Sprintf("%s:%d:%s", file, line, name) }
	for _, c := range a.Callables {
		id := cid(c.File, c.StartLine, c.Name)
		calls[id] = &Callable{ID: id, Module: modKey(fileLang(c.File), c.Module), Name: c.Name, File: c.File, Line: c.StartLine,
			Cyclomatic: c.Cyclomatic, Nesting: c.Nesting, Lines: c.Lines, Debt: c.Debt}
	}
	// metrics for clone members that are not hotspots
	metrics := map[string]*Callable{}
	for _, s := range snap.Symbols {
		if s.Callable() {
			metrics[cid(s.File, s.StartLine, s.DisplayName())] = &Callable{Cyclomatic: s.Metrics.Cyclomatic, Nesting: s.Metrics.MaxNesting, Lines: s.Metrics.Lines}
		}
	}
	for _, g := range a.Clones {
		cl := Clone{Exact: g.Exact, Tokens: g.RedundantTokens, Debt: g.Debt}
		n := float64(len(g.Members))
		for _, m := range g.Members {
			id := cid(m.File, m.StartLine, m.Name)
			c := calls[id]
			if c == nil {
				c = &Callable{ID: id, Module: modKey(fileLang(m.File), m.Module), Name: m.Name, File: m.File, Line: m.StartLine}
				if mt := metrics[id]; mt != nil {
					c.Cyclomatic, c.Nesting, c.Lines = mt.Cyclomatic, mt.Nesting, mt.Lines
				}
				calls[id] = c
			}
			c.CloneDebt = r2(c.CloneDebt + float64(m.Tokens)*(n-1)/n/a.Weights.TokensPerPoint)
			cl.Members = append(cl.Members, id)
		}
		d.Clones = append(d.Clones, cl)
	}
	for _, c := range calls {
		d.Callables = append(d.Callables, *c)
	}
	sort.Slice(d.Callables, func(i, j int) bool { return d.Callables[i].ID < d.Callables[j].ID })

	inCycle := map[string]bool{}
	for _, c := range a.Cycles {
		for _, u := range c {
			inCycle[u] = true
		}
	}
	viol := map[[2]string]string{}
	for _, v := range a.Violations {
		viol[[2]string{modKey(v.Language, v.From), modKey(v.Language, v.To)}] = v.Rule
	}
	var gids []string
	for id := range snap.Graphs {
		gids = append(gids, id)
	}
	sort.Strings(gids)
	for _, gid := range gids {
		g := snap.Graphs[gid]
		var froms []string
		for f := range g.Edges {
			froms = append(froms, f)
		}
		sort.Strings(froms)
		for _, f := range froms {
			var tos []string
			for t := range g.Edges[f] {
				tos = append(tos, t)
			}
			sort.Strings(tos)
			for _, t := range tos {
				fk, tk := modKey(gid, f), modKey(gid, t)
				d.Imports = append(d.Imports, Import{From: fk, To: tk, Violation: viol[[2]string{fk, tk}], InCycle: inCycle[f] && inCycle[t] && sameCycle(a.Cycles, f, t)})
			}
		}
	}
	if d.Imports == nil {
		d.Imports = []Import{}
	}
	if d.Clones == nil {
		d.Clones = []Clone{}
	}
	return d
}

func sameCycle(cycles [][]string, a, b string) bool {
	for _, c := range cycles {
		ha, hb := false, false
		for _, u := range c {
			ha = ha || u == a
			hb = hb || u == b
		}
		if ha && hb {
			return true
		}
	}
	return false
}

// Trim keeps the n callables with the largest total debt (complexity, nesting
// and clone share) and drops clone groups left with fewer than two members.
// Modules and import edges are kept. n <= 0 keeps everything.
func (d *Data) Trim(n int) {
	if n <= 0 || len(d.Callables) <= n {
		return
	}
	cs := append([]Callable{}, d.Callables...)
	sort.SliceStable(cs, func(i, j int) bool { return cs[i].Debt+cs[i].CloneDebt > cs[j].Debt+cs[j].CloneDebt })
	keep := map[string]bool{}
	for _, c := range cs[:n] {
		keep[c.ID] = true
	}
	var out []Callable
	for _, c := range d.Callables {
		if keep[c.ID] {
			out = append(out, c)
		}
	}
	d.Callables = out
	var clones []Clone
	for _, g := range d.Clones {
		var ms []string
		for _, m := range g.Members {
			if keep[m] {
				ms = append(ms, m)
			}
		}
		if len(ms) >= 2 {
			g.Members = ms
			clones = append(clones, g)
		}
	}
	if clones == nil {
		clones = []Clone{}
	}
	d.Clones = clones
	d.Notes = append(d.Notes, fmt.Sprintf("trimmed view: the %d callables with the largest debt are shown", n))
}

// WriteHTML writes the self-contained page.
func WriteHTML(w io.Writer, d *Data) error {
	js, err := json.Marshal(d)
	if err != nil {
		return err
	}
	// json.Marshal escapes <, > and & so the data cannot close the script.
	out := strings.Replace(page, "[[.Title]]", html.EscapeString(fmt.Sprintf("Debt graph: %s @ %.12s", d.Repository, d.Revision.Commit)), 1)
	out = strings.Replace(out, "[[.Data]]", string(js), 1)
	_, err = io.WriteString(w, out)
	return err
}
