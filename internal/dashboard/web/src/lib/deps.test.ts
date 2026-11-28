import type { GraphData, GraphModule } from "../api";
import { buildDeps, groupKey } from "./deps";

const mod = (path: string, lines = 100): GraphModule => ({ id: "go:" + path, lang: "go", group: "", lines, callables: 1, hotspots: 0, debt: { complexity: 0, nesting: 0, duplication: 0, architecture: 0, total: 0 }, density: 0, score: null });

const data: GraphData = {
  repository: "r",
  modules: [mod("cmd/app"), mod("internal/a"), mod("internal/b"), mod("web/src/x"), mod("web/src/y")],
  callables: [],
  clones: [],
  imports: [
    { from: "go:cmd/app", to: "go:internal/a", count: 2, sites: ["cmd/app/main.go:3"] },
    { from: "go:internal/a", to: "go:internal/b", count: 1 },
    { from: "go:web/src/x", to: "go:web/src/y", count: 1 },
    { from: "go:web/src/y", to: "go:web/src/x", count: 1 },
  ],
  roots: [{ dir: ".", kind: "go.mod" }, { dir: "web", kind: "package.json" }],
};

describe("groupKey", () => {
  it("groups by project root, else top folder, and follows expanded groups", () => {
    expect(groupKey("web/src/x", ["web"], new Set())).toBe("web");
    expect(groupKey("internal/a", ["web"], new Set())).toBe("internal");
    expect(groupKey("internal/a", ["web"], new Set(["internal"]))).toBe("internal/a");
    expect(groupKey("internal", ["web"], new Set(["internal"]))).toBe("internal/(own files)");
  });
});

describe("buildDeps", () => {
  it("aggregates imports between groups and layers users above what they use", () => {
    const g = buildDeps(data, new Set());
    const key = (k: string) => g.groups.find((x) => x.key === k)!;
    expect(g.edges.find((e) => e.from === "cmd" && e.to === "internal")).toMatchObject({ count: 2, cycle: false });
    expect(key("cmd").layer).toBeLessThan(key("internal").layer);
    expect(key("cmd").y).toBeLessThan(key("internal").y);
  });

  it("marks groups in an import cycle once expanded to that level", () => {
    const g = buildDeps(data, new Set(["web", "web/src"]));
    expect(g.groups.filter((x) => x.inCycle).map((x) => x.key).sort()).toEqual(["web/src/x", "web/src/y"]);
    expect(g.edges.filter((e) => e.cycle)).toHaveLength(2);
  });
});
