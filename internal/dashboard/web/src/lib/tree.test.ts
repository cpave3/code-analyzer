import type { GraphModule } from "../api";
import { buildTree, densityBand, findNode, parentPath } from "./tree";

const mod = (id: string, lines: number, debt: number): GraphModule => ({
  id,
  lang: "go",
  group: id.split(":")[1].split("/")[0],
  lines,
  callables: 1,
  hotspots: 0,
  debt: { complexity: debt, nesting: 0, duplication: 0, architecture: 0, total: debt },
  density: (debt * 1000) / lines,
  score: null,
});

describe("buildTree", () => {
  const root = buildTree([mod("go:internal/lang", 100, 10), mod("go:internal/lang/golang", 300, 30), mod("go:cmd/codeeval", 200, 0)]);

  it("sums lines and debt up the folders", () => {
    expect(root.lines).toBe(600);
    expect(root.debt).toBe(40);
  });

  it("merges single-child folder chains and keeps a folder's own module as a tile", () => {
    const names = root.children.map((c) => c.name).sort();
    expect(names).toEqual(["cmd/codeeval", "internal/lang"]);
    const lang = root.children.find((c) => c.name === "internal/lang")!;
    expect(lang.module).toBeUndefined();
    expect(lang.children.map((c) => c.name).sort()).toEqual(["(own files)", "golang"]);
    const cmd = root.children.find((c) => c.name === "cmd/codeeval")!;
    expect(cmd.module?.id).toBe("go:cmd/codeeval");
  });

  it("finds folders and module tiles by path", () => {
    expect(findNode(root, "internal/lang")?.children.length).toBe(2);
    expect(findNode(root, "internal/lang/golang")?.module?.id).toBe("go:internal/lang/golang");
    expect(findNode(root, "nope")).toBeNull();
    expect(findNode(root, "")).toBe(root);
  });
});

describe("helpers", () => {
  it("parentPath", () => {
    expect(parentPath("a/b/c")).toBe("a/b");
    expect(parentPath("a")).toBe("");
  });
  it("densityBand", () => {
    expect(densityBand(10)).toBe(0);
    expect(densityBand(70)).toBe(2);
    expect(densityBand(400)).toBe(4);
  });
});
