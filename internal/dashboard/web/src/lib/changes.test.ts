import type { GraphData } from "../api";
import { compare } from "./changes";

const g = (mods: [string, number][], fns: [string, string, number][]): GraphData => ({
  repository: "r",
  modules: mods.map(([id, d]) => ({ id, lang: "go", group: "", lines: 100, callables: 1, hotspots: 0, debt: { complexity: d, nesting: 0, duplication: 0, architecture: 0, total: d }, density: d * 10, score: null })),
  callables: fns.map(([file, name, d], i) => ({ id: file + ":" + i + ":" + name, module: "", name, file, line: i * 10, cc: 12, nest: 2, lines: 30, debt: d, cloneDebt: 0 })),
  clones: [],
  imports: [],
});

describe("compare", () => {
  it("classifies modules and functions against the previous commit", () => {
    const prev = g([["go:a", 10], ["go:b", 5], ["go:gone", 7]], [["a.go", "f", 4]]);
    const now = g([["go:a", 12], ["go:b", 5], ["go:c", 3]], [["a.go", "f", 1]]);
    const c = compare(now, prev);
    expect(c.states.get("m:go:a")).toMatchObject({ state: "up", delta: 2 });
    expect(c.states.get("m:go:b")?.state).toBe("same");
    expect(c.states.get("m:go:c")).toMatchObject({ state: "new", delta: 3 });
    expect(c.states.get("c:a.go#f")).toMatchObject({ state: "down", delta: -3 });
    expect(c.biggest.map((x) => x.id)).toEqual(["m:go:c", "c:a.go#f", "m:go:a"]);
    expect(c.gone.map((x) => x.label)).toEqual(["gone"]);
    expect(c.counts).toEqual({ new: 1, up: 1, down: 1, same: 1 });
  });

  it("treats a function that only moved as the same function", () => {
    const prev = g([], [["a.go", "f", 4]]);
    const now = g([], [["a.go", "f", 4]]);
    now.callables[0].line = 99;
    expect(compare(now, prev).states.get("c:a.go#f")?.state).toBe("same");
  });

  it("reports nothing as changed without a previous commit", () => {
    const c = compare(g([["go:a", 1]], []), null);
    expect(c.biggest).toEqual([]);
    expect(c.counts.same).toBe(1);
  });
});
