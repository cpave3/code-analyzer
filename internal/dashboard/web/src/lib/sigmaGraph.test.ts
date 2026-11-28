import type { GraphData } from "../core/api";
import type { Tokens } from "../core/theme";
import { compare } from "./changes";
import { buildDeps } from "./deps";
import { depsGraph, edgeStyle, functionColour, mix, networkGraph, nodeStyle, ringColour } from "./sigmaGraph";

const t: Tokens = {
  bg: "#000000", panel: "#111111", fg: "#ffffff", dim: "#777777", line: "#333333", track: "#222222", accent: "#0000ff", good: "#00ff00", bad: "#ff0000",
  components: { complexity: "#000001", nesting: "#000002", duplication: "#000003", architecture: "#000004" },
  ramp: ["#100000", "#200000", "#300000", "#400000", "#500000"],
};
const mod = (id: string, lines: number, debt: number) => ({ id, lang: "go", group: "", lines, callables: 1, hotspots: 0, debt: { complexity: debt, nesting: 0, duplication: 0, architecture: 0, total: debt }, density: (debt * 1000) / lines, score: null });
const data: GraphData = {
  repository: "r",
  modules: [mod("go:a", 400, 40), mod("go:b", 300, 0), mod("go:tiny", 10, 0)],
  callables: [
    { id: "a.go:1:big", module: "go:a", name: "big", file: "a.go", line: 1, cc: 30, nest: 3, lines: 80, debt: 21, cloneDebt: 0 },
    { id: "a.go:9:small", module: "go:a", name: "small", file: "a.go", line: 9, cc: 11, nest: 2, lines: 20, debt: 2, cloneDebt: 0 },
  ],
  clones: [],
  imports: [{ from: "go:a", to: "go:b", count: 1 }, { from: "go:b", to: "go:a", count: 1, cycle: true }],
};
const opts = { sizeBy: "lines" as const, functions: true, clones: true, imports: true, minDebt: 5 };

describe("networkGraph", () => {
  it("keeps modules with debt or ≥ 150 lines and functions above the minimum debt", () => {
    const g = networkGraph(data, opts, t, null, 0);
    expect(g.hasNode("go:a") && g.hasNode("go:b")).toBe(true);
    expect(g.hasNode("go:tiny")).toBe(false);
    expect(g.hasNode("a.go:1:big")).toBe(true);
    expect(g.hasNode("a.go:9:small")).toBe(false);
    expect(g.getNodeAttribute("go:a", "bad")).toBe(true);
    expect(g.getEdgeAttribute("go:b", "go:a", "bad")).toBe(true);
  });
  it("is deterministic for a seed", () => {
    const a = networkGraph(data, opts, t, null, 3),
      b = networkGraph(data, opts, t, null, 3);
    expect(a.getNodeAttribute("go:a", "x")).toBe(b.getNodeAttribute("go:a", "x"));
  });
  it("outlines changes since the previous commit", () => {
    const prev = { ...data, modules: [mod("go:a", 400, 30), mod("go:b", 300, 0)] };
    const g = networkGraph(data, opts, t, compare(data, prev), 0);
    expect(g.getNodeAttribute("go:a", "ring")).toBe(t.bad);
    expect(ringColour(null, "m:go:a", t)).toBeUndefined();
  });
});

describe("styles", () => {
  const g = depsGraph(buildDeps(data, new Set()), "lines", t);
  const base = { focus: null, picked: null, marks: true, labels: true, find: "" };
  it("colours groups in a cycle red when marks are on", () => {
    const k = g.nodes()[0];
    expect(nodeStyle(g, k, g.getNodeAttributes(k), base, t).color).toBe(g.getNodeAttribute(k, "bad") ? t.bad : g.getNodeAttribute(k, "color"));
  });
  it("fades nodes the search misses and hides edges away from the focus", () => {
    const n = networkGraph(data, { ...opts, functions: false }, t, null, 0);
    expect(nodeStyle(n, "go:b", n.getNodeAttributes("go:b"), { ...base, find: "zzz" }, t)).toMatchObject({ color: t.track, label: "" });
    const e = n.edge("go:a", "go:b")!;
    expect(edgeStyle(n, e, n.getEdgeAttributes(e), { ...base, focus: "go:a" }, t)).toMatchObject({ head: "arrow", color: t.accent });
  });
  it("helpers", () => {
    expect(mix("#000000", "#ffffff")).toBe("#808080");
    expect(functionColour(20, t)).toBe(t.ramp[4]);
  });
});
