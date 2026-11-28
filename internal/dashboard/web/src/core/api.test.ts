import { normalizeGraph, type GraphData } from "./api";
import { compare } from "../lib/changes";

describe("normalizeGraph", () => {
  it("turns null lists from a clean commit into arrays that the comparison can read", () => {
    const clean = { repository: "r", modules: null, callables: null, clones: null, imports: null } as unknown as GraphData;
    const g = normalizeGraph(clean);
    expect(g.callables).toEqual([]);
    expect(() => compare(g, normalizeGraph(clean))).not.toThrow();
  });
});
