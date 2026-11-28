// Derivations for a folder report (an explain report for one path).
import type { Factor } from "../api";

export type Band = "critical" | "high" | "medium" | "low";
export const BANDS: Band[] = ["critical", "high", "medium", "low"];

// Severity bands group v2 debt points per item; a reading aid, not calibrated.
export function band(debt: number): Band {
  return debt >= 25 ? "critical" : debt >= 10 ? "high" : debt >= 4 ? "medium" : "low";
}

export const KIND_NAMES: Record<string, string> = {
  complexity: "complexity",
  nesting: "nesting",
  duplication: "duplication",
  rule_violation: "rule violation",
  cycle: "import cycle",
};

// componentOf maps a factor kind to the debt component it adds to.
export function componentOf(kind: string): "complexity" | "nesting" | "duplication" | "architecture" {
  if (kind === "complexity" || kind === "nesting" || kind === "duplication") return kind;
  return "architecture";
}

// bySubfolder sums factor debt by the next path segment below path.
export function bySubfolder(factors: Factor[], path: string): { name: string; path: string | null; debt: number }[] {
  const sub = new Map<string, number>();
  for (const f of factors) {
    const file = f.location.file || f.module;
    const rest = path && file.startsWith(path + "/") ? file.slice(path.length + 1) : path ? file : file;
    const seg = rest.includes("/") ? rest.split("/")[0] + "/" : "(files directly here)";
    sub.set(seg, (sub.get(seg) ?? 0) + f.debt);
  }
  return [...sub.entries()]
    .sort((a, b) => b[1] - a[1])
    .map(([name, debt]) => ({ name, debt, path: name.endsWith("/") ? (path ? path + "/" : "") + name.slice(0, -1) : null }));
}

export interface FileRow {
  file: string;
  debt: number;
  items: number;
  hotspots: number; // complexity or nesting items
  worst: string;
}

// byFile rolls the factors up per file (cycle and rule factors land on the
// file that holds the import).
export function byFile(factors: Factor[]): FileRow[] {
  const m = new Map<string, FileRow>();
  for (const f of factors) {
    const file = f.location.file || f.module;
    let r = m.get(file);
    if (!r) m.set(file, (r = { file, debt: 0, items: 0, hotspots: 0, worst: "" }));
    r.debt += f.debt;
    r.items++;
    if (f.kind === "complexity" || f.kind === "nesting") r.hotspots++;
    if (!r.worst) r.worst = f.location.symbol || f.module;
  }
  return [...m.values()].sort((a, b) => b.debt - a.debt || a.file.localeCompare(b.file));
}

export function crumbs(path: string): { name: string; path: string }[] {
  const out = [{ name: "repository", path: "" }];
  let acc = "";
  for (const p of path ? path.split("/") : []) {
    acc = acc ? acc + "/" + p : p;
    out.push({ name: p, path: acc });
  }
  return out;
}
