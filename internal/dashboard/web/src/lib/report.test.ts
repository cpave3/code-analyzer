import type { Factor } from "../api";
import { band, byFile, bySubfolder, crumbs } from "./report";

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
