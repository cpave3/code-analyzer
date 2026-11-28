// The Dependencies view: modules grouped by project or package root (go.mod,
// package.json, .csproj; else top folder), aggregated import edges between
// groups, and a layered layout where a group sits above everything it uses.
// Groups in an import cycle share a layer; their edges are cycle edges.
// Ported from the previous dashboard's buildDeps/layoutDeps.
import type { GraphData } from "../api";
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
  const groups = new Map<string, DepGroup>();
  for (const m of d.modules) {
    const label = labelOf.get(m.id)!;
    const k = groupKey(label, roots, expanded);
    let g = groups.get(k);
    if (!g) {
      const parts = k.split("/");
      groups.set(k, (g = { key: k, label: parts.length > 2 ? "…/" + parts.slice(-2).join("/") : k, lines: 0, debt: 0, density: 0, modules: [], expandable: false, inCycle: false, layer: 0, x: 0, y: 0 }));
    }
    g.lines += m.lines;
    g.debt += m.debt.total;
    g.modules.push(label);
  }
  for (const g of groups.values()) {
    g.density = g.lines ? (g.debt * 1000) / g.lines : 0;
    g.expandable = !g.key.endsWith(OWN) && g.modules.some((m) => m !== g.key);
  }
  const edges = new Map<string, DepEdge>();
  for (const e of d.imports) {
    const a = labelOf.get(e.from),
      b = labelOf.get(e.to);
    if (a == null || b == null) continue;
    const ga = groupKey(a, roots, expanded),
      gb = groupKey(b, roots, expanded);
    if (ga === gb) continue;
    const k = ga + "\u0000" + gb;
    let x = edges.get(k);
    if (!x) edges.set(k, (x = { from: ga, to: gb, count: 0, sites: [], rule: false, cycle: false }));
    x.count += e.count ?? 1;
    if (e.violation) x.rule = true;
    for (const s of e.sites ?? []) if (x.sites.length < 5) x.sites.push(s);
  }
  const out = { groups: [...groups.values()], edges: [...edges.values()] };
  layout(out);
  return out;
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
      if (index[w] === undefined) {
        strong(w);
        low[v] = Math.min(low[v], low[w]);
      } else if (on[w]) low[v] = Math.min(low[v], index[w]);
    }
    if (low[v] === index[v]) {
      let w: number;
      do {
        w = stack.pop()!;
        on[w] = false;
        comp[w] = v;
      } while (w !== v);
    }
  };
  ns.forEach((_, v) => {
    if (index[v] === undefined) strong(v);
  });
  const size = new Map<number, number>();
  ns.forEach((_, v) => size.set(comp[v], (size.get(comp[v]) ?? 0) + 1));
  for (const e of g.edges) e.cycle = comp[idx.get(e.from)!] === comp[idx.get(e.to)!];
  ns.forEach((n, v) => (n.inCycle = (size.get(comp[v]) ?? 0) > 1));
  const succ = new Map<number, Set<number>>(),
    indeg = new Map<number, number>();
  for (const e of g.edges) {
    if (e.cycle) continue;
    const a = comp[idx.get(e.from)!],
      b = comp[idx.get(e.to)!];
    if (!succ.has(a)) succ.set(a, new Set());
    succ.get(a)!.add(b);
  }
  for (const s of succ.values()) for (const b of s) indeg.set(b, (indeg.get(b) ?? 0) + 1);
  const layerOf = new Map<number, number>();
  const queue = [...new Set(comp)].filter((c) => !indeg.get(c));
  for (const c of queue) layerOf.set(c, 0);
  while (queue.length) {
    const c = queue.shift()!;
    for (const d of succ.get(c) ?? []) {
      layerOf.set(d, Math.max(layerOf.get(d) ?? 0, (layerOf.get(c) ?? 0) + 1));
      indeg.set(d, indeg.get(d)! - 1);
      if (indeg.get(d) === 0) queue.push(d);
    }
  }
  const layers: number[][] = [];
  ns.forEach((n, v) => {
    n.layer = layerOf.get(comp[v]) ?? 0;
    (layers[n.layer] ??= []).push(v);
  });
  const order: number[] = [];
  for (const L of layers) if (L) L.sort((a, b) => comp[a] - comp[b] || ns[a].key.localeCompare(ns[b].key)).forEach((v, i) => (order[v] = i));
  for (let pass = 0; pass < 6; pass++) {
    const down = pass % 2 === 0;
    const seq = down ? layers.slice(1) : layers.slice(0, -1).reverse();
    for (const L of seq) {
      if (!L) continue;
      const bary = new Map<number, number>();
      for (const v of L) {
        const near = (down ? inn[v] : out[v]).filter((w) => ns[w].layer === ns[v].layer + (down ? -1 : 1));
        bary.set(v, near.length ? near.reduce((a, w) => a + order[w], 0) / near.length : order[v]);
      }
      L.sort((a, b) => bary.get(a)! - bary.get(b)! || comp[a] - comp[b]).forEach((v, i) => (order[v] = i));
    }
  }
  const perRow = 12,
    gapX = 190,
    rowH = 90;
  let y = 0;
  for (const L of layers) {
    if (!L) continue;
    for (let s = 0; s < L.length; s += perRow) {
      const row = L.slice(s, s + perRow);
      row.forEach((v, i) => {
        ns[v].x = (i - (row.length - 1) / 2) * gapX;
        ns[v].y = y;
      });
      y += rowH;
    }
    y += 30;
  }
}
