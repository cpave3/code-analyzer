package csharp

import (
	"testing"

	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

const src = `using System;
using Shop.Domain.Models;
using static System.Math;
using Alias = Shop.Other.Thing;

namespace Shop.Api;

public sealed class Foo<T> : Base where T : class
{
    private readonly int _x;
    public int Auto { get; set; }
    public int Calc => _x > 0 ? 1 : 2;
    public Foo(int x) { _x = x; }
    public int Run(int a, string b)
    {
        if (a > 0 && b != null) { return 1; } else if (a < 0 || b == "") { return 2; }
        for (var i = 0; i < a; i++) { foreach (var c in b) { while (true) { break; } } }
        switch (a) { case 1: return 3; default: break; }
        try { Do(); } catch (Exception) { }
        var z = a switch { 1 => "one", _ => "other" };
        Func<int, int> f = q => q + 1;
        return b?.Length ?? 0;
    }
    class Inner { void N() { } }
}
internal record Point(int X, int Y);
public interface IFoo { int Run(); }
`

func parse(t *testing.T) *model.ParsedFile {
	t.Helper()
	pf, err := New().Parse([]byte(src))
	if err != nil || len(pf.ParseErrors) > 0 {
		t.Fatalf("parse: %v %+v", err, pf.ParseErrors)
	}
	return pf
}

func TestSymbolsAndMetrics(t *testing.T) {
	pf := parse(t)
	syms := map[string]model.Symbol{}
	for _, s := range pf.Symbols {
		syms[s.DisplayName()] = s
	}
	for _, want := range []string{"Foo", "Foo.Calc", "Foo.Foo", "Foo.Run", "Foo.Inner", "Foo.Inner.N", "Point", "IFoo"} {
		if _, ok := syms[want]; !ok {
			t.Errorf("missing symbol %s; have %v", want, keys(syms))
		}
	}
	if _, ok := syms["Foo.Auto"]; ok {
		t.Error("auto-property measured as a callable")
	}
	run := syms["Foo.Run"]
	// decisions: if, &&, else-if, ||, for, foreach, while, case 1 section, catch, switch arm "1" (not _), ??, ternary none -> 11
	if run.Metrics == nil || run.Metrics.Decisions != 11 || run.Metrics.Cyclomatic != 12 {
		t.Errorf("Run metrics %+v, want 11 decisions", run.Metrics)
	}
	if run.Metrics.MaxNesting != 3 {
		t.Errorf("Run nesting %d, want 3 (for > foreach > while)", run.Metrics.MaxNesting)
	}
	if !syms["Foo.Run"].Exported || syms["Point"].Exported || syms["Foo.Inner.N"].Exported {
		t.Errorf("exported flags: Run %v Point %v N %v", syms["Foo.Run"].Exported, syms["Point"].Exported, syms["Foo.Inner.N"].Exported)
	}
	if c := syms["Foo.Calc"]; c.Metrics == nil || c.Metrics.Decisions != 1 {
		t.Errorf("expression-bodied property: %+v", c.Metrics)
	}
}

func keys(m map[string]model.Symbol) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestImportsAndNamespace(t *testing.T) {
	pf := parse(t)
	if pf.PackageName != "Shop.Api" {
		t.Errorf("namespace %q", pf.PackageName)
	}
	got := map[string]string{}
	for _, i := range pf.Imports {
		got[i.Path] = i.Name
	}
	want := map[string]string{"System": "", "Shop.Domain.Models": "", "System.Math": "static", "Shop.Other.Thing": "Alias"}
	for k, v := range want {
		if n, ok := got[k]; !ok || n != v {
			t.Errorf("import %s: got %q present=%v, want %q (all %v)", k, n, ok, v, got)
		}
	}
}

func TestTestPathsAndGraph(t *testing.T) {
	a := New()
	for p, want := range map[string]bool{
		"svc/Shop.Api.Tests/FooTests.cs": true, "tests/Foo.cs": true, "svc/Shop.Api/FooTests.cs": true,
		"svc/Shop.Api/Foo.cs": false, "src/Contest/Foo.cs": false,
	} {
		if a.IsTestPath(p) != want {
			t.Errorf("IsTestPath(%s) = %v", p, !want)
		}
	}
	file := func(path, ns string, imports ...string) lang.GraphFile {
		pf := &model.ParsedFile{PackageName: ns}
		for _, i := range imports {
			pf.Imports = append(pf.Imports, model.Import{Path: i, Line: 1})
		}
		return lang.GraphFile{Path: path, Class: model.ClassSource, Parsed: pf}
	}
	g, _ := a.BuildGraph(lang.GraphInput{Files: []lang.GraphFile{
		file("api/A.cs", "Shop.Api", "Shop.Domain", "System.Text"),
		file("domain/D.cs", "Shop.Domain"),
		file("domain/models/M.cs", "Shop.Domain.Models", "Shop.Api"),
	}})
	if len(g.Edges["api"]["domain"]) != 1 || len(g.Edges["domain/models"]["api"]) != 1 {
		t.Errorf("edges %+v", g.Edges)
	}
	if !g.External["api"]["System.Text"] {
		t.Errorf("external %+v", g.External)
	}
}
