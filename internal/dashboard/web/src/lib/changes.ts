// Commit-to-commit comparison of two module graphs: which modules and
// functions gained or lost debt, which are new and which are gone.
import type { GraphData } from "../core/api";
import { moduleLabel } from "./format";

export type ChangeState = "new" | "up" | "down" | "same";

export interface Changed {
  id: string; // "m:<module id>" or "c:<file>#<name>"
  kind: "module" | "function";
  label: string;
  debt: number;
  before: number | null;
  delta: number;
  state: ChangeState;
}

// Functions are identified by file and name, not line, so one that only
// moved keeps its identity between commits.
const fnId = (file: string, name: string) => "c:" + file + "#" + name;

function items(d: GraphData): Map<string, { kind: "module" | "function"; label: string; debt: number }> {
  const out = new Map<string, { kind: "module" | "function"; label: string; debt: number }>();
  for (const m of d.modules) out.set("m:" + m.id, { kind: "module", label: moduleLabel(m.id), debt: m.debt.total });
  for (const c of d.callables) out.set(fnId(c.file, c.name), { kind: "function", label: c.name, debt: c.debt + c.cloneDebt });
  return out;
}

export interface Comparison {
  states: Map<string, Changed>;
  biggest: Changed[]; // largest |delta| first
  gone: { id: string; kind: "module" | "function"; label: string; debt: number }[];
  counts: Record<ChangeState, number>;
}

export function compare(now: GraphData, prev: GraphData | null): Comparison {
  const a = items(now);
  const b = prev ? items(prev) : null;
  const states = new Map<string, Changed>();
  const counts: Record<ChangeState, number> = { new: 0, up: 0, down: 0, same: 0 };
  for (const [id, x] of a) {
    const before = b ? (b.get(id)?.debt ?? null) : x.debt;
    const state: ChangeState = before == null ? "new" : x.debt > before + 0.01 ? "up" : x.debt < before - 0.01 ? "down" : "same";
    counts[state]++;
    states.set(id, { id, ...x, before, delta: before == null ? x.debt : x.debt - before, state });
  }
  const biggest = [...states.values()]
    .filter((c) => c.state !== "same" && Math.abs(c.delta) >= 0.01)
    .sort((p, q) => Math.abs(q.delta) - Math.abs(p.delta))
    .slice(0, 8);
  const gone = b
    ? [...b.entries()]
        .filter(([id, x]) => !a.has(id) && x.debt > 0)
        .map(([id, x]) => ({ id, ...x }))
        .sort((p, q) => q.debt - p.debt)
        .slice(0, 4)
    : [];
  return { states, biggest, gone, counts };
}

export const moduleKey = (moduleId: string) => "m:" + moduleId;
export const functionKey = fnId;
