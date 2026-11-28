// Derivations over the measured window: per-commit step changes and the
// "step changes by provenance group" table.
import type { Commit, Debt, Metrics } from "../api";
import { stepKey } from "./provenance";

export type MetricsByIndex = (Metrics | undefined)[];

export const COMPONENTS = ["complexity", "nesting", "duplication", "architecture"] as const;
export type Component = (typeof COMPONENTS)[number];

export function usable(m: Metrics | undefined): m is Metrics {
  return !!m && m.done && !m.error;
}

// change returns get(after) − get(before) for commit i against its parent
// in the history, or null when either side is not measured yet.
export function change(ms: MetricsByIndex, i: number, get: (m: Metrics) => number | null | undefined): number | null {
  const a = ms[i - 1],
    b = ms[i];
  if (!usable(a) || !usable(b)) return null;
  const x = get(a),
    y = get(b);
  return x == null || y == null ? null : y - x;
}

export function debtChange(ms: MetricsByIndex, i: number): Debt | null {
  const a = ms[i - 1],
    b = ms[i];
  if (!usable(a) || !usable(b)) return null;
  return {
    complexity: b.debt.complexity - a.debt.complexity,
    nesting: b.debt.nesting - a.debt.nesting,
    duplication: b.debt.duplication - a.debt.duplication,
    architecture: b.debt.architecture - a.debt.architecture,
    total: b.debt.total - a.debt.total,
  };
}

export interface GroupRow {
  key: string;
  steps: number;
  lines: number;
  debt: number;
  up: number;
  down: number;
  hotspots: number;
  gate: number;
  advisory: number;
  resolved: number;
  near: number;
}

// groupRows sums the window's steps per provenance key, most steps first.
export function groupRows(commits: Commit[], ms: MetricsByIndex, from: number, to: number, fields: string[]): GroupRow[] {
  const g = new Map<string, GroupRow>();
  for (let i = Math.max(1, from); i <= to; i++) {
    const m = ms[i],
      p = ms[i - 1];
    if (!usable(m) || !usable(p) || !m.stepDone) continue;
    const key = stepKey(commits[i], fields);
    let a = g.get(key);
    if (!a) g.set(key, (a = { key, steps: 0, lines: 0, debt: 0, up: 0, down: 0, hotspots: 0, gate: 0, advisory: 0, resolved: 0, near: 0 }));
    const dd = m.debt.total - p.debt.total;
    a.steps++;
    a.lines += m.changedLines;
    a.debt += dd;
    a.hotspots += m.hotspots - p.hotspots;
    if (dd > 0.005) a.up++;
    else if (dd < -0.005) a.down++;
    a.gate += m.introducedGate ?? 0;
    a.advisory += m.introducedAdvisory ?? 0;
    a.resolved += m.resolved ?? 0;
    if (m.detail && p.detail) a.near += (m.nearMisses ?? 0) + (m.fragments ?? 0) - ((p.nearMisses ?? 0) + (p.fragments ?? 0));
  }
  return [...g.values()].sort((x, y) => y.steps - x.steps || x.key.localeCompare(y.key));
}

// clampWindow orders and clamps a requested window to the history.
export function clampWindow(from: number, to: number, n: number): [number, number] {
  if (from > to) [from, to] = [to, from];
  return [Math.max(0, Math.min(n - 1, from)), Math.max(0, Math.min(n - 1, to))];
}

// lastN is the window holding the newest k commits.
export function lastN(k: number, n: number): [number, number] {
  return clampWindow(n - Math.max(1, k), n - 1, n);
}

// dateWindow is the window of commits whose day falls in [a, b] (YYYY-MM-DD),
// or null when none does.
export function dateWindow(commits: Commit[], a: string, b: string): [number, number] | null {
  const lo = commits.findIndex((c) => c.time.slice(0, 10) >= a);
  let hi = -1;
  for (let i = commits.length - 1; i >= 0; i--)
    if (commits[i].time.slice(0, 10) <= b) {
      hi = i;
      break;
    }
  return lo < 0 || hi < 0 || lo > hi ? null : [lo, hi];
}

// recentre moves a window of the same size so that it is centred on i.
export function recentre(i: number, from: number, to: number, n: number): [number, number] {
  const size = to - from,
    half = Math.floor(size / 2);
  let lo = i - half;
  lo = Math.max(0, Math.min(n - 1 - size, lo));
  return [lo, lo + size];
}
