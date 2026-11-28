import type { GraphData, GraphModule } from "../core/api";
import type { Tokens } from "../core/theme";
import { compare } from "./changes";
import { inkFor, markedModules, shownFolder, tileTooltip, tiles } from "./treemap";
import { buildTree } from "./tree";

const t = {
  bg: "#000000", panel: "#111111", fg: "#ffffff", dim: "#777777", line: "#333333", track: "#222222", accent: "#0000ff", good: "#00ff00", bad: "#ff0000",
  components: { complexity: "#000001", nesting: "#000002", duplication: "#000003", architecture: "#000004" },
  ramp: ["#f0f0f0", "#c0c0c0", "#808080", "#404040", "#101010"],
} as Tokens;
const mod = (id: string, lang: string, debt: number, nesting = 0): GraphModule => ({ id, lang, group: "", lines: 1000, callables: 1, hotspots: 0, debt: { complexity: debt, nesting, duplication: 0, architecture: 0, total: debt + nesting }, density: debt + nesting, score: null });
const g: GraphData = { repository: "r", modules: [mod("go:internal/a", "go", 10), mod("typescript:internal/b", "typescript", 0, 200), mod("go:cmd/x", "go", 0)], callables: [], clones: [], imports: [{ from: "go:internal/a", to: "go:cmd/x", cycle: true }] };
const tree = buildTree(g.modules);
const base = { sizeBy: "lines" as const, colourBy: "density" as const, lang: "", focus: "", marked: markedModules(g), cmp: null };

describe("treemap tiles", () => {
  it("shows a module's folder, or the nearest existing ancestor", () => {
    expect(shownFolder(tree, "internal/a").path).toBe("internal");
    expect(shownFolder(tree, "internal/nope/deeper").path).toBe("internal");
    expect(shownFolder(tree, "internal").path).toBe("internal");
  });
  it("filters by language and drops folders left empty", () => {
    const names = tiles(tree, { ...base, lang: "typescript" }, t).map((x) => x.name);
    expect(names).toEqual(["internal"]);
    expect(tiles(tree, { ...base, lang: "typescript" }, t)[0].children!.map((x) => x.name)).toEqual(["b"]);
  });
  it("marks cycle modules, outlines the focus and changes, and colours by component", () => {
    const prev = { ...g, modules: [mod("go:internal/a", "go", 5), mod("typescript:internal/b", "typescript", 0, 200), mod("go:cmd/x", "go", 0)] };
    const inner = tiles(tree, { ...base, cmp: compare(g, prev), focus: "internal/b", colourBy: "component" }, t).find((x) => x.name === "internal")!.children!;
    const a = inner.find((x) => x.path === "internal/a")!;
    const b = inner.find((x) => x.path === "internal/b")!;
    expect(a.name).toBe("↻ a");
    expect(a.itemStyle.borderColor).toBe(t.bad);
    expect(b.itemStyle.borderColor).toBe(t.fg);
    expect(b.itemStyle.color).toBe(t.components.nesting);
    expect(tileTooltip(a)).toContain("more debt than the previous commit");
  });
  it("picks readable label ink", () => {
    expect(inkFor("#f0f0f0")).toBe("#111");
    expect(inkFor("#101010")).toBe("#fff");
  });
});
