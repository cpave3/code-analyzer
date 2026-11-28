import type { Commit, Metrics } from "../api";
import { change, clampWindow, dateWindow, debtChange, groupRows, lastN, recentre } from "./history";

const debt = (t: number) => ({ complexity: t, nesting: 0, duplication: 0, architecture: 0, total: t });
const m = (t: number, extra: Partial<Metrics> = {}): Metrics =>
  ({ done: true, stepDone: true, detail: false, detailStep: false, debt: debt(t), debtPerKLOC: 1, score: 50, lines: 1000, files: 1, callables: 1, hotspots: 1, decisionPoints: 1, clonedCallables: 0, ruleViolations: 0, cycleUnits: 0, changedLines: 10, ...extra }) as Metrics;
const commit = (i: number, model: string, time = "2026-11-01T10:00:00Z"): Commit => ({ index: i, sha: "s" + i, parent: "", subject: "", author: "", time, merge: false, provenance: [{ commit: "s" + i, model: [model], source: "trailer" }] });

describe("step changes", () => {
  it("compares a commit with its parent once both are measured", () => {
    const ms = [m(10), m(7), undefined];
    expect(change(ms, 1, (x) => x.debt.total)).toBe(-3);
    expect(change(ms, 2, (x) => x.debt.total)).toBeNull();
    expect(debtChange(ms, 1)?.complexity).toBe(-3);
  });
  it("skips commits whose measurement failed", () => {
    expect(change([m(1), m(2, { error: "parse" })], 1, (x) => x.debt.total)).toBeNull();
  });
});

describe("groupRows", () => {
  it("sums steps per provenance group inside the window", () => {
    const cs = [commit(0, "a"), commit(1, "a"), commit(2, "b"), commit(3, "a")];
    const ms = [m(10), m(12, { changedLines: 5, introducedGate: 1 }), m(11), m(11)];
    const rows = groupRows(cs, ms, 1, 3, ["model"]);
    expect(rows.map((r) => r.key)).toEqual(["model=a", "model=b"]);
    expect(rows[0]).toMatchObject({ steps: 2, lines: 15, debt: 2, up: 1, down: 0, gate: 1 });
    expect(rows[1]).toMatchObject({ steps: 1, debt: -1, down: 1 });
  });
});

describe("windows", () => {
  it("clamps and orders", () => expect(clampWindow(9, -3, 5)).toEqual([0, 4]));
  it("takes the newest k", () => expect(lastN(3, 10)).toEqual([7, 9]));
  it("finds commits by day", () => {
    const cs = [commit(0, "a", "2026-11-01T00:00:00Z"), commit(1, "a", "2026-11-03T00:00:00Z"), commit(2, "a", "2026-11-05T00:00:00Z")];
    expect(dateWindow(cs, "2026-11-02", "2026-11-04")).toEqual([1, 1]);
    expect(dateWindow(cs, "2026-12-01", "2026-12-02")).toBeNull();
  });
  it("recentres a same-size window, inside the history", () => {
    expect(recentre(2, 10, 19, 100)).toEqual([0, 9]);
    expect(recentre(50, 10, 19, 100)).toEqual([46, 55]);
    expect(recentre(99, 10, 19, 100)).toEqual([90, 99]);
  });
});
