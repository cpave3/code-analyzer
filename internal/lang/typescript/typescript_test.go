package typescript

import (
	"testing"

	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

const sample = `import { useState } from "react";
import type { Item } from "./types";
export { helper } from "../util/helper";

// A comment that mentions if (x) && y.
export function classify(n: number): string {
  if (n < 0) {
    return "neg";
  } else if (n === 0) {
    return "zero";
  }
  for (const x of [1, 2]) {
    while (x > n && n !== 3) { n++; }
  }
  return n > 10 ? "big" : n ?? "small";
}

export const Panel = ({ items }: { items: Item[] }) => {
  const [open, setOpen] = useState(false);
  const onClick = () => { if (open) { setOpen(false); } };
  return <div onClick={onClick}>{open && items.map((i) => <span>{i.name}</span>)}</div>;
};

const Memoised = memo(function Inner(props: Props) {
  switch (props.kind) {
    case "a": return 1;
    case "b": return 2;
    default: return 0;
  }
});

class Store {
  private cache = new Map<string, number>();
  get(key: string): number | undefined {
    try { return this.cache.get(key); } catch (e) { return undefined; }
  }
  private reset = () => { this.cache.clear(); };
}

export interface Item { name: string }
export type Kind = "a" | "b";
enum Mode { On, Off }
const lazy = () => import("./lazy");
`

func find(pf *model.ParsedFile, name string) *model.Symbol {
	for i := range pf.Symbols {
		if pf.Symbols[i].DisplayName() == name {
			return &pf.Symbols[i]
		}
	}
	return nil
}

func TestParseSymbolsMetricsImports(t *testing.T) {
	a := New()
	pf, err := a.Parse([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	if len(pf.ParseErrors) != 0 {
		t.Fatalf("unexpected parse errors: %+v", pf.ParseErrors)
	}
	want := map[string]struct {
		kind     string
		exported bool
		cc, nest int
	}{
		"classify":    {model.KindFunction, true, 8, 2},  // if, else-if, for, while, &&, ternary, ?? = 7 decisions; for>while nesting 2
		"Panel":       {model.KindFunction, true, 3, 2},  // inner if + && (callbacks count toward Panel)
		"Memoised":    {model.KindFunction, false, 3, 1}, // two cases inside the memo() function
		"Store.get":   {model.KindMethod, false, 2, 1},   // catch
		"Store.reset": {model.KindMethod, false, 1, 0},
		"lazy":        {model.KindFunction, false, 1, 0},
	}
	for name, w := range want {
		s := find(pf, name)
		if s == nil {
			t.Errorf("missing symbol %s; have %v", name, names(pf))
			continue
		}
		if s.Kind != w.kind || s.Exported != w.exported || s.Metrics == nil || s.Metrics.Cyclomatic != w.cc || s.Metrics.MaxNesting != w.nest {
			t.Errorf("%s: kind %s exported %v metrics %+v; want %+v", name, s.Kind, s.Exported, s.Metrics, w)
		}
		if s.ExactFP == "" || s.StructFP == "" {
			t.Errorf("%s: missing fingerprints", name)
		}
	}
	for _, typ := range []string{"Store", "Item", "Kind", "Mode"} {
		if s := find(pf, typ); s == nil || s.Kind != model.KindType {
			t.Errorf("type %s missing or wrong kind: %+v", typ, s)
		}
	}
	var specs []string
	for _, imp := range pf.Imports {
		specs = append(specs, imp.Path)
	}
	for _, w := range []string{"react", "./types", "../util/helper", "./lazy"} {
		if !contains(specs, w) {
			t.Errorf("import %s missing: %v", w, specs)
		}
	}
	if pf.Refs["useState"] == 0 || pf.Refs["this.cache"] != 0 && pf.Refs["cache"] == 0 {
		t.Errorf("refs missing: %v", pf.Refs)
	}
}

func TestStructuralFingerprintIgnoresNames(t *testing.T) {
	a := New()
	x, _ := a.Parse([]byte("function f(a: number) { if (a > 1) { return a * 2; } return 0; }"))
	y, _ := a.Parse([]byte("function g(b: number) { if (b > 7) { return b * 3; } return 0; }"))
	if x.Symbols[0].StructFP != y.Symbols[0].StructFP || x.Symbols[0].ExactFP == y.Symbols[0].ExactFP {
		t.Error("renamed copy should share the structural but not the exact fingerprint")
	}
}

func TestAngleBracketAssertionFallsBackToTS(t *testing.T) {
	pf, _ := New().Parse([]byte("function f(x: unknown) { const n = <number>x; return n + 1; }"))
	if len(pf.ParseErrors) != 0 || find(pf, "f") == nil {
		t.Errorf("TS-only syntax should parse via the TypeScript grammar: %+v", pf.ParseErrors)
	}
}

func TestBuildGraphResolvesRelativeAliasAndWorkspace(t *testing.T) {
	a := New()
	files := map[string]string{
		"web/src/app/App.tsx":       `import { x } from "../util"; import { v } from "@voice/rec"; import "./App.css"; import { h } from "@acme/hub/lib/h"; import React from "react";`,
		"web/src/util/index.ts":     `export const x = 1;`,
		"web/src/voice/rec.ts":      `export const v = 1;`,
		"packages/hub/src/lib/h.ts": `export const h = 1;`,
		"web/src/app/App.test.tsx":  `import { x } from "../util";`,
	}
	read := map[string]string{
		"web/tsconfig.json": `{ // comment
  "compilerOptions": { "baseUrl": ".", "paths": { "@voice/*": ["src/voice/*"], "react": ["./node_modules/@types/react"], }, },
}`,
		"packages/hub/package.json": `{"name": "@acme/hub"}`,
	}
	var all []string
	var gf []lang.GraphFile
	for p, src := range files {
		all = append(all, p)
		pf, _ := a.Parse([]byte(src))
		class := model.ClassSource
		if a.IsTestPath(p) {
			class = model.ClassTest
		}
		gf = append(gf, lang.GraphFile{Path: p, Class: class, Parsed: pf})
	}
	for p := range read {
		all = append(all, p)
	}
	g, diags := a.BuildGraph(lang.GraphInput{Files: gf, AllPaths: all, ReadFile: func(p string) ([]byte, error) { return []byte(read[p]), nil }})
	for _, to := range []string{"web/src/util", "web/src/voice", "packages/hub/src/lib"} {
		if _, ok := g.Edges["web/src/app"][to]; !ok {
			t.Errorf("missing edge web/src/app → %s: %v", to, g.Edges["web/src/app"])
		}
	}
	if len(diags) != 0 {
		t.Errorf("an alias into node_modules is external, not unresolved: %+v", diags)
	}
	if len(g.Edges["web/src/app"]) != 3 || !g.External["web/src/app"]["react"] {
		t.Errorf("edges %v external %v", g.Edges["web/src/app"], g.External["web/src/app"])
	}
}

func names(pf *model.ParsedFile) []string {
	var out []string
	for _, s := range pf.Symbols {
		out = append(out, s.DisplayName())
	}
	return out
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
