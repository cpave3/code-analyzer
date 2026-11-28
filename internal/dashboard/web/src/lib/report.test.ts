import type { Factor } from "../core/api";
import { agentPrompt, band, byFile, bySubfolder, crumbs, groupFindings, sortedItems } from "./report";

const f = (file: string, debt: number, kind = "complexity", symbol = "fn"): Factor =>
  ({ id: file + symbol, kind, module: "m", language: "go", location: { file, startLine: 1, symbol }, debt, shareOfScopeDebt: 0, cumulativeShare: 0, explanation: "", remediation: { strategy: "", target: "", steps: [], acceptance: [] } }) as Factor;

describe("folder report", () => {
  it("bands debt", () => {
    expect([band(30), band(12), band(4), band(1)]).toEqual(["critical", "high", "medium", "low"]);
  });
  it("sums debt by subfolder below the path", () => {
    const rows = bySubfolder([f("internal/a/x.go", 5), f("internal/a/y.go", 2), f("internal/b/z.go", 9), f("internal/top.go", 1)], "internal");
    expect(rows).toEqual([
      { name: "b/", path: "internal/b", debt: 9 },
      { name: "a/", path: "internal/a", debt: 7 },
      { name: "(files directly here)", path: null, debt: 1 },
    ]);
  });
  it("rolls factors up per file, worst file first", () => {
    const rows = byFile([f("a.go", 2, "complexity", "small"), f("b.go", 9, "duplication"), f("a.go", 8, "nesting", "big")]);
    expect(rows[0]).toMatchObject({ file: "a.go", debt: 10, items: 2, hotspots: 2 });
    expect(rows[1]).toMatchObject({ file: "b.go", hotspots: 0 });
  });
  it("builds breadcrumbs", () => {
    expect(crumbs("a/b").map((c) => c.path)).toEqual(["", "a", "a/b"]);
  });
});

describe("items and findings", () => {
  it("filters and sorts items", () => {
    const a = f("x/a.go", 5, "complexity", "alpha"),
      b = f("x/b.go", 9, "nesting", "beta");
    a.metrics = { cyclomatic: 30 };
    b.metrics = { cyclomatic: 12 };
    expect(sortedItems([a, b], "", "debt").map((x) => x.debt)).toEqual([9, 5]);
    expect(sortedItems([a, b], "", "cc").map((x) => x.debt)).toEqual([5, 9]);
    expect(sortedItems([a, b], "BETA", "debt")).toEqual([b]);
    expect(sortedItems([b, a], "", "name")[0]).toBe(a);
  });
  it("groups findings with gate rules first", () => {
    const d = (rule: string, severity: string) => ({ rule, severity, file: "a.go", line: 1, message: "" });
    const g = groupFindings([d("r1", "advisory"), d("r1", "advisory"), d("r2", "gate")]);
    expect(g.map((x) => [x[0].rule, x.length])).toEqual([["r2", 1], ["r1", 2]]);
  });
  it("writes a factor as an agent prompt", () => {
    const x = f("a.go", 12, "complexity", "Run");
    x.remediation = { strategy: "reduce branching", target: "cc ≤ 9", steps: ["extract"], acceptance: ["cc ≤ 9"] };
    expect(agentPrompt(x)).toContain("Reduce the complexity debt of Run at a.go:1.");
    expect(agentPrompt(x)).toContain("- extract");
  });
});
