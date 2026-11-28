// The graphology graphs the Graph view draws with Sigma, and how each node
// and edge is styled for the current selection, search and options.
import Graph from "graphology";
import forceAtlas2 from "graphology-layout-forceatlas2";
import type { GraphData, GraphModule } from "../core/api";
import type { Tokens } from "../core/theme";
import { functionKey, moduleKey, type Comparison } from "./changes";
import type { DepGraph } from "./deps";
import { moduleLabel } from "./format";
import { densityBand } from "./tree";

export type SizeBy = "lines" | "debt";

export interface NetworkOpts {
  sizeBy: SizeBy;
  functions: boolean;
  clones: boolean;
  imports: boolean;
  minDebt: number;
}

export interface NodeAttrs {
  x: number;
  y: number;
  size: number;
  color: string;
  label: string;
  ring?: string;
  bad?: boolean;
  kind: "group" | "module" | "function";
}

export interface EdgeAttrs {
  size: number;
  color: string;
  bad?: boolean;
  arrow?: boolean;
}

export type SGraph = Graph<NodeAttrs, EdgeAttrs>;

// depsGraph draws the layered group graph (positions come from buildDeps).
export function depsGraph(deps: DepGraph, sizeBy: SizeBy, t: Tokens): SGraph {
  const g: SGraph = new Graph({ type: "directed", multi: false });
  const value = (n: DepGraph["groups"][number]) => (sizeBy === "debt" ? n.debt : n.lines);
  const max = Math.max(1, ...deps.groups.map(value));
  for (const n of deps.groups) {
    const label = (n.expandable ? "▸ " : "") + n.label;
    g.addNode(n.key, { x: n.x, y: -n.y, size: 5 + 13 * Math.sqrt(value(n) / max), color: t.ramp[densityBand(n.density)], label, bad: n.inCycle, kind: "group" });
  }
  for (const e of deps.edges) {
    const bad = e.cycle || e.rule;
    g.mergeEdge(e.from, e.to, { size: 0.8 + Math.log2(1 + e.count), color: bad ? t.bad : t.dim, bad, arrow: true });
  }
  return g;
}

// ringColour is the outline for an item's change since the previous commit.
export function ringColour(cmp: Comparison | null, key: string, t: Tokens): string | undefined {
  const st = cmp?.states.get(key)?.state;
  return { up: t.bad, down: t.good, new: t.accent, same: undefined }[st ?? "same"];
}

// seeded returns a deterministic pseudo-random generator, so a layout is
// stable until the viewer asks for a new one.
export function seeded(seed: number): () => number {
  let s = 1 + seed * 7919;
  return () => ((s = (s * 16807) % 2147483647) - 1) / 2147483646;
}

// networkGraph is the force layout: modules with debt (or ≥ 150 lines) and
// the functions with the most debt, pulled together by imports and clone
// links.
export function networkGraph(data: GraphData, o: NetworkOpts, t: Tokens, cmp: Comparison | null, seed: number): SGraph {
  const g: SGraph = new Graph({ type: "directed", multi: false });
  const rnd = seeded(seed);
  addModules(g, data, o.sizeBy, t, cmp, rnd);
  if (o.functions) addFunctions(g, data, o.minDebt, t, cmp, rnd);
  if (o.functions && o.clones) addClones(g, data, t);
  if (o.imports) addImports(g, data, t);
  if (g.order > 1) forceAtlas2.assign(g, { iterations: 220, settings: { ...forceAtlas2.inferSettings(g), barnesHutOptimize: g.order > 400, gravity: 1, scalingRatio: 12, adjustSizes: true } });
  return g;
}

const shown = (m: GraphModule) => m.debt.total > 0 || m.lines >= 150;

function addModules(g: SGraph, data: GraphData, sizeBy: SizeBy, t: Tokens, cmp: Comparison | null, rnd: () => number) {
  const value = (m: GraphModule) => (sizeBy === "debt" ? m.debt.total : m.lines);
  const max = Math.max(1, ...data.modules.map(value));
  const bad = new Set<string>();
  for (const e of data.imports.filter((x) => x.violation || x.cycle)) bad.add(e.from).add(e.to);
  for (const m of data.modules.filter(shown)) {
    const color = t.ramp[densityBand(m.density)];
    g.addNode(m.id, { x: rnd() * 100, y: rnd() * 100, size: 3 + 11 * Math.sqrt(value(m) / max), color, label: moduleLabel(m.id), ring: ringColour(cmp, moduleKey(m.id), t), bad: bad.has(m.id), kind: "module" });
  }
}

// Functions are coloured by their own debt: under 5, 5–15, 15 or more.
export function functionColour(debt: number, t: Tokens): string {
  return debt < 5 ? t.ramp[1] : debt < 15 ? t.ramp[3] : t.ramp[4];
}

function addFunctions(g: SGraph, data: GraphData, minDebt: number, t: Tokens, cmp: Comparison | null, rnd: () => number) {
  for (const c of data.callables) {
    const d = c.debt + c.cloneDebt;
    if (d < minDebt || !g.hasNode(c.module) || g.hasNode(c.id)) continue;
    const p = g.getNodeAttributes(c.module);
    g.addNode(c.id, { x: p.x + rnd() * 6 - 3, y: p.y + rnd() * 6 - 3, size: 1.5 + Math.sqrt(d) * 0.4, color: functionColour(d, t), label: c.name, ring: ringColour(cmp, functionKey(c.file, c.name), t), kind: "function" });
    g.mergeEdge(c.id, c.module, { size: 0.5, color: t.line });
  }
}

function addClones(g: SGraph, data: GraphData, t: Tokens) {
  for (const cg of data.clones) {
    const [first, ...rest] = cg.members;
    for (const b of rest.filter((m) => g.hasNode(first) && g.hasNode(m) && !g.hasEdge(first, m))) g.addEdge(first, b, { size: 1.2, color: t.components.duplication });
  }
}

function addImports(g: SGraph, data: GraphData, t: Tokens) {
  for (const e of data.imports) {
    if (!g.hasNode(e.from) || !g.hasNode(e.to) || g.hasEdge(e.from, e.to)) continue;
    const bad = !!(e.violation || e.cycle);
    g.addEdge(e.from, e.to, { size: bad ? 2 : 0.8, color: bad ? t.bad : t.dim, bad, arrow: true });
  }
}

export interface StyleState {
  focus: string | null; // the picked node, else the hovered one
  picked: string | null;
  marks: boolean;
  labels: boolean;
  find: string;
}

// nodeStyle restyles a node for the selection, search and options: items
// outside the focus's neighbourhood or the search fade out.
export function nodeStyle(g: SGraph, key: string, attrs: NodeAttrs, s: StyleState, t: Tokens): Record<string, unknown> {
  const out: Record<string, unknown> = { labelColor: t.fg, backdropColor: t.panel };
  out.color = s.marks && attrs.bad ? t.bad : attrs.ring === t.bad ? mix(attrs.color, t.bad) : attrs.color;
  if (!s.labels && key !== s.focus) out.label = "";
  const needle = s.find.trim().toLowerCase();
  const missed = needle !== "" && !attrs.label.toLowerCase().includes(needle);
  const unrelated = !!s.focus && key !== s.focus && !g.areNeighbors(key, s.focus);
  if (missed || unrelated) Object.assign(out, { color: t.track, label: "" });
  if (key === s.picked) Object.assign(out, { highlighted: true, forceLabel: true });
  return out;
}

// edgeStyle draws arrowheads, hides edges away from the focus and colours
// the focus's imports (accent) and importers (duplication colour).
export function edgeStyle(g: SGraph, key: string, attrs: EdgeAttrs, s: StyleState, t: Tokens): Record<string, unknown> {
  const out: Record<string, unknown> = { head: attrs.arrow ? "arrow" : "none" };
  if (!s.marks && attrs.bad) out.color = t.dim;
  if (!s.focus) return out;
  const [a, b] = g.extremities(key);
  if (a !== s.focus && b !== s.focus) out.visibility = "hidden";
  else if (!attrs.bad) out.color = a === s.focus ? t.accent : t.components.duplication;
  return out;
}

// mix blends two hex colours half and half.
export function mix(a: string, b: string): string {
  const hex = /^#[0-9a-f]{6}$/i;
  if (!hex.test(a) || !hex.test(b)) return a;
  const p = (h: string) => [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16));
  const [x, y] = [p(a), p(b)];
  return "#" + x.map((v, i) => Math.round((v + y[i]) / 2).toString(16).padStart(2, "0")).join("");
}
