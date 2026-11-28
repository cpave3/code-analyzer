// The ECharts treemap behind the Modules page: which folder is shown for a
// path, and each tile's size, colour, outline and label.
import type { EChartsCoreOption } from "echarts/core";
import type { GraphData } from "../core/api";
import type { Tokens } from "../core/theme";
import { moduleKey, type ChangeState, type Comparison } from "./changes";
import { fmt } from "./format";
import { COMPONENTS } from "./history";
import { densityBand, findNode, parentPath, sizeOf, type SizeBy, type TreeNode } from "./tree";

export type ColourBy = "density" | "component";

export interface TileOpts {
  sizeBy: SizeBy;
  colourBy: ColourBy;
  lang: string; // "" = every language
  focus: string; // the path being viewed; its tile is outlined
  marked: Set<string>; // module ids in an import cycle or breaking a rule
  cmp: Comparison | null;
}

// shownFolder is the folder the map shows for path: the folder itself, the
// folder holding a module tile, or the nearest existing ancestor.
export function shownFolder(tree: TreeNode, path: string): TreeNode {
  const n = findNode(tree, path);
  if (n && !n.module) return n;
  let p = parentPath(path);
  while (p && !findNode(tree, p)) p = parentPath(p);
  return findNode(tree, p) ?? tree;
}

// markedModules lists modules at either end of a cycle or rule-breaking import.
export function markedModules(g: GraphData | undefined): Set<string> {
  const out = new Set<string>();
  for (const e of (g?.imports ?? []).filter((x) => x.cycle || x.violation)) out.add(e.from).add(e.to);
  return out;
}

export function tileColour(n: TreeNode, colourBy: ColourBy, t: Tokens): string {
  if (colourBy !== "component" || !n.module) return t.ramp[densityBand(n.density)];
  const [k, v] = COMPONENTS.map((c) => [c, n.module!.debt[c]] as const).sort((a, b) => b[1] - a[1])[0];
  return v > 0 ? t.components[k] : t.track;
}

const RING: Record<ChangeState, keyof Tokens | null> = { up: "bad", down: "good", new: "accent", same: null };

export function tileRing(n: TreeNode, o: TileOpts, t: Tokens): string | undefined {
  if (!n.module) return undefined;
  if (n.path === o.focus) return t.fg;
  const key = RING[o.cmp?.states.get(moduleKey(n.module.id))?.state ?? "same"];
  return key ? (t[key] as string) : undefined;
}

// inkFor picks dark or light text for a tile colour by its luminance.
export function inkFor(hex: string): string {
  const m = /^#([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})$/i.exec(hex);
  if (!m) return "#111";
  const lin = (x: string) => {
    const c = parseInt(x, 16) / 255;
    return c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4;
  };
  return 0.2126 * lin(m[1]) + 0.7152 * lin(m[2]) + 0.0722 * lin(m[3]) > 0.3 ? "#111" : "#fff";
}

export interface Tile {
  name: string;
  path: string;
  bad: boolean;
  value: number;
  debt: number;
  lines: number;
  density: number;
  state?: ChangeState;
  children?: Tile[];
  itemStyle: Record<string, unknown>;
  label?: { color: string };
}

function moduleTile(c: TreeNode, o: TileOpts, t: Tokens): Tile {
  const color = tileColour(c, o.colourBy, t);
  const ring = tileRing(c, o, t);
  const bad = o.marked.has(c.module!.id);
  return {
    ...facts(c, o),
    name: bad ? "↻ " + c.name : c.name,
    bad,
    state: o.cmp?.states.get(moduleKey(c.module!.id))?.state,
    itemStyle: { color, borderColor: ring ?? t.panel, borderWidth: ring ? 3 : 1 },
    label: { color: inkFor(color) },
  };
}

function facts(c: TreeNode, o: TileOpts) {
  return { name: c.name, path: c.path, bad: false, value: sizeOf(c, o.sizeBy), debt: c.debt, lines: c.lines, density: c.density };
}

// tiles turns a folder's children into treemap tiles; folders left empty by
// the language filter are dropped.
export function tiles(n: TreeNode, o: TileOpts, t: Tokens): Tile[] {
  const out: Tile[] = [];
  for (const c of n.children) {
    if (c.module) {
      if (!o.lang || c.module.lang === o.lang) out.push(moduleTile(c, o, t));
      continue;
    }
    const children = tiles(c, o, t);
    if (children.length) out.push({ ...facts(c, o), children, itemStyle: { borderColor: t.line, color: t.bg } });
  }
  return out;
}

const esc = (s: string) => s.replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" })[c]!);

const STATE_TEXT: Record<ChangeState, string> = { up: " · more debt than the previous commit", down: " · less debt than the previous commit", new: " · new in this commit", same: "" };

export function tileTooltip(d: Tile): string {
  const bad = d.bad ? " · in an import cycle or breaking a rule" : "";
  return `<b>${esc(d.path || "(repository)")}</b><br>${fmt(d.lines)} lines · debt ${fmt(d.debt, 1)} · ${fmt(d.density, 1)} per kLOC${STATE_TEXT[d.state ?? "same"]}${bad}`;
}

export function treemapOption(folder: TreeNode, o: TileOpts, t: Tokens): EChartsCoreOption {
  return {
    animation: false,
    tooltip: { confine: true, formatter: (p: { data: Tile }) => tileTooltip(p.data) },
    series: [
      {
        type: "treemap",
        data: tiles(folder, o, t),
        roam: false,
        nodeClick: false,
        breadcrumb: { show: false },
        width: "100%",
        height: "100%",
        top: 0,
        left: 0,
        squareRatio: 1.2,
        label: { show: true, formatter: "{b}", color: "#111", fontSize: 11, overflow: "truncate" },
        upperLabel: { show: true, height: 18, color: t.fg, fontSize: 11, fontWeight: 600 },
        itemStyle: { gapWidth: 2 },
        levels: [{ itemStyle: { borderWidth: 0, gapWidth: 3 } }, { itemStyle: { borderWidth: 1, gapWidth: 2 } }, { itemStyle: { gapWidth: 1 } }],
      },
    ],
  };
}
