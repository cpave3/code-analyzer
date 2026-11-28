// The Dependencies view: modules grouped by project or package root (go.mod,
// package.json, .csproj; else top folder), aggregated import edges between
// groups, and a layered layout where a group sits above everything it uses.
// Groups in an import cycle share a layer; their edges are cycle edges.
// Ported from the previous dashboard's buildDeps/layoutDeps.
import type { GraphData } from "../core/api";
import { moduleLabel } from "./format";

export interface DepGroup {
  key: string;
  label: string;
  lines: number;
  debt: number;
  density: number;
  modules: string[]; // module labels
  expandable: boolean;
  inCycle: boolean;
  layer: number;
  x: number;
  y: number;
}

export interface DepEdge {
  from: string;
  to: string;
  count: number;
  sites: string[];
  rule: boolean;
  cycle: boolean;
}

export interface DepGraph {
  groups: DepGroup[];
  edges: DepEdge[];
}

const OWN = "/(own files)";

export function groupKey(label: string, roots: string[], expanded: Set<string>): string {
  let key = label.split("/")[0];
  for (const r of roots)
    if (label === r || label.startsWith(r + "/")) {
      key = r;
      break;
    }
  while (expanded.has(key)) {
    if (label === key) return key + OWN;
    key = key + "/" + label.slice(key.length + 1).split("/")[0];
  }
  return key;
}

export function buildDeps(d: GraphData, expanded: Set<string>): DepGraph {
  const roots = (d.roots ?? [])
    .map((r) => r.dir)
    .filter((x) => x !== ".")
    .sort((a, b) => b.length - a.length);
  const labelOf = new Map(d.modules.map((m) => [m.id, moduleLabel(m.id)]));
  const keyOf = (label: string) => groupKey(label, roots, expanded);
  const out = { groups: groupModules(d, labelOf, keyOf), edges: groupEdges(d, labelOf, keyOf) };
  layout(out);
  return out;
}

function newGroup(key: string): DepGroup {
  const parts = key.split("/");
  const label = parts.length > 2 ? "…/" + parts.slice(-2).join("/") : key;
  return { key, label, lines: 0, debt: 0, density: 0, modules: [], expandable: false, inCycle: false, layer: 0, x: 0, y: 0 };
}

function groupModules(d: GraphData, labelOf: Map<string, string>, keyOf: (label: string) => string): DepGroup[] {
  const groups = new Map<string, DepGroup>();
  for (const m of d.modules) {
    const label = labelOf.get(m.id)!;
    const k = keyOf(label);
    const g = groups.get(k) ?? groups.set(k, newGroup(k)).get(k)!;
    g.lines += m.lines;
    g.debt += m.debt.total;
    g.modules.push(label);
  }
  for (const g of groups.values()) {
    g.density = g.lines ? (g.debt * 1000) / g.lines : 0;
    g.expandable = !g.key.endsWith(OWN) && g.modules.some((m) => m !== g.key);
  }
  return [...groups.values()];
}

function groupEdges(d: GraphData, labelOf: Map<string, string>, keyOf: (label: string) => string): DepEdge[] {
  const edges = new Map<string, DepEdge>();
  for (const e of d.imports) {
    const a = labelOf.get(e.from),
      b = labelOf.get(e.to);
    if (a == null || b == null) continue;
    const ga = keyOf(a),
      gb = keyOf(b);
    if (ga === gb) continue;
    const k = ga + "\u0000" + gb;
    const x = edges.get(k) ?? edges.set(k, { from: ga, to: gb, count: 0, sites: [], rule: false, cycle: false }).get(k)!;
    x.count += e.count ?? 1;
    x.rule ||= !!e.violation;
    x.sites.push(...(e.sites ?? []).slice(0, 5 - x.sites.length));
  }
  return [...edges.values()];
}

// layout: Tarjan SCCs share a layer; the condensation is layered by longest
// path (users above what they use); barycentre sweeps reduce crossings;
// long layers wrap onto extra rows.
function layout(g: DepGraph) {
  const ns = g.groups,
    idx = new Map(ns.map((n, i) => [n.key, i]));
  const out: number[][] = ns.map(() => []);
  const inn: number[][] = ns.map(() => []);
  for (const e of g.edges) {
    out[idx.get(e.from)!].push(idx.get(e.to)!);
    inn[idx.get(e.to)!].push(idx.get(e.from)!);
  }
  const comp = components(out);
  markCycles(g, idx, comp);
  const layers = assignLayers(g, idx, comp);
  const order = orderLayers(layers, comp, ns, inn, out);
  place(layers, ns, order);
}

// components labels each node with its strongly connected component (Tarjan).
export function components(out: number[][]): number[] {
  let counter = 0;
  const index: number[] = [],
    low: number[] = [],
    on: boolean[] = [],
    stack: number[] = [],
    comp: number[] = [];
  const strong = (v: number) => {
    index[v] = low[v] = counter++;
    stack.push(v);
    on[v] = true;
    for (const w of out[v]) {
      if (index[w] === undefined) strong(w);
      if (on[w]) low[v] = Math.min(low[v], index[w] === undefined ? low[w] : Math.min(low[w], index[w]));
    }
    if (low[v] !== index[v]) return;
    let w: number;
    do {
      w = stack.pop()!;
      on[w] = false;
      comp[w] = v;
    } while (w !== v);
  };
  out.forEach((_, v) => index[v] === undefined && strong(v));
  return comp;
}

function markCycles(g: DepGraph, idx: Map<string, number>, comp: number[]) {
  const size = new Map<number, number>();
  for (const c of comp) size.set(c, (size.get(c) ?? 0) + 1);
  for (const e of g.edges) e.cycle = comp[idx.get(e.from)!] === comp[idx.get(e.to)!];
  g.groups.forEach((n, v) => (n.inCycle = size.get(comp[v])! > 1));
}

// assignLayers layers the condensation by longest path from its sources.
function assignLayers(g: DepGraph, idx: Map<string, number>, comp: number[]): number[][] {
  const succ = new Map<number, Set<number>>(),
    indeg = new Map<number, number>();
  for (const e of g.edges.filter((x) => !x.cycle)) {
    const a = comp[idx.get(e.from)!],
      b = comp[idx.get(e.to)!];
    const s = succ.get(a) ?? succ.set(a, new Set()).get(a)!;
    if (!s.has(b)) indeg.set(b, (indeg.get(b) ?? 0) + 1);
    s.add(b);
  }
  const layerOf = new Map<number, number>();
  const queue = [...new Set(comp)].filter((c) => !indeg.get(c));
  queue.forEach((c) => layerOf.set(c, 0));
  while (queue.length) {
    const c = queue.shift()!;
    for (const d of succ.get(c) ?? []) {
      layerOf.set(d, Math.max(layerOf.get(d) ?? 0, layerOf.get(c)! + 1));
      indeg.set(d, indeg.get(d)! - 1);
      if (indeg.get(d) === 0) queue.push(d);
    }
  }
  const layers: number[][] = [];
  g.groups.forEach((n, v) => {
    n.layer = layerOf.get(comp[v]) ?? 0;
    (layers[n.layer] ??= []).push(v);
  });
  return layers;
}

// orderLayers orders each layer by component, then by barycentre sweeps
// against the neighbouring layers to reduce crossings.
function orderLayers(layers: number[][], comp: number[], ns: DepGroup[], inn: number[][], out: number[][]): number[] {
  const order: number[] = [];
  const renumber = (L: number[]) => L.forEach((v, i) => (order[v] = i));
  for (const L of layers.filter(Boolean)) renumber(L.sort((a, b) => comp[a] - comp[b] || ns[a].key.localeCompare(ns[b].key)));
  for (let pass = 0; pass < 6; pass++) {
    const down = pass % 2 === 0;
    const seq = (down ? layers.slice(1) : layers.slice(0, -1).reverse()).filter(Boolean);
    for (const L of seq) {
      const bary = new Map(L.map((v) => [v, barycentre(v, down ? inn[v] : out[v], ns, order, down ? -1 : 1)]));
      renumber(L.sort((a, b) => bary.get(a)! - bary.get(b)! || comp[a] - comp[b]));
    }
  }
  return order;
}

function barycentre(v: number, neighbours: number[], ns: DepGroup[], order: number[], step: number): number {
  const near = neighbours.filter((w) => ns[w].layer === ns[v].layer + step);
  return near.length ? near.reduce((a, w) => a + order[w], 0) / near.length : order[v];
}

// place puts each layer on rows of at most perRow groups, centred.
function place(layers: number[][], ns: DepGroup[], order: number[]) {
  const perRow = 12,
    gapX = 190,
    rowH = 90;
  let y = 0;
  for (const L of layers.filter(Boolean)) {
    const sorted = [...L].sort((a, b) => order[a] - order[b]);
    for (let s = 0; s < sorted.length; s += perRow) {
      const row = sorted.slice(s, s + perRow);
      row.forEach((v, i) => Object.assign(ns[v], { x: (i - (row.length - 1) / 2) * gapX, y }));
      y += rowH;
    }
    y += 30;
  }
}
