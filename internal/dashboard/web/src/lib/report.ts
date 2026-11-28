// Derivations for a folder report (an explain report for one path).
import type { DefectFinding, Factor } from "../core/api";

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
  hotspots: number; // distinct callables with complexity or nesting debt
  worst: string;
}

// byFile rolls the factors up per file (cycle and rule factors land on the
// file that holds the import).
export function byFile(factors: Factor[]): FileRow[] {
  const m = new Map<string, FileRow>();
  const hot = new Set<string>(); // a callable over both thresholds has two factors
  for (const f of factors) {
    const file = f.location.file || f.module;
    const r = m.get(file) ?? m.set(file, { file, debt: 0, items: 0, hotspots: 0, worst: "" }).get(file)!;
    r.debt += f.debt;
    r.items++;
    const callable = `${file}:${f.location.startLine}:${f.location.symbol ?? ""}`;
    if ((f.kind === "complexity" || f.kind === "nesting") && !hot.has(callable)) {
      hot.add(callable);
      r.hotspots++;
    }
    r.worst ||= f.location.symbol || f.module;
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

// agentPrompt is a factor written as a task for a coding agent.
export function agentPrompt(f: Factor): string {
  const rem = f.remediation;
  const loc = f.location;
  const lines = [
    `Reduce the ${KIND_NAMES[f.kind] ?? f.kind} debt of ${loc.symbol || f.module} at ${loc.file}:${loc.startLine}.`,
    f.explanation,
    rem.strategy && `Fix: ${rem.strategy}. ${rem.target}`,
    ...rem.steps.map((s) => "- " + s),
    rem.acceptance.length ? "Done when: " + rem.acceptance.join("; ") : "",
    rem.caution ? "Caution: " + rem.caution : "",
  ];
  return lines.filter(Boolean).join("\n");
}

// groupFindings groups defect findings by severity and rule, gate rules
// first, then the most frequent.
export function groupFindings(findings: DefectFinding[]): DefectFinding[][] {
  const g = new Map<string, DefectFinding[]>();
  for (const f of findings) {
    const k = f.severity + "|" + f.rule;
    g.set(k, [...(g.get(k) ?? []), f]);
  }
  const rank = (xs: DefectFinding[]) => (xs[0].severity === "gate" ? 0 : 1);
  return [...g.values()].sort((a, b) => rank(a) - rank(b) || b.length - a.length);
}

export type SortKey = "debt" | "cc" | "nest" | "lines" | "name";

const METRIC: Record<Exclude<SortKey, "debt" | "name">, string> = { cc: "cyclomatic", nest: "nesting", lines: "lines" };

// sortedItems filters factors by symbol, file or kind and sorts them, largest
// first (by name: alphabetically).
export function sortedItems(factors: Factor[], query: string, sort: SortKey): Factor[] {
  const q = query.trim().toLowerCase();
  const hit = (f: Factor) => (f.location.symbol ?? "").toLowerCase().includes(q) || f.location.file.toLowerCase().includes(q) || f.kind.includes(q);
  const xs = q ? factors.filter(hit) : [...factors];
  if (sort === "name") return xs.sort((a, b) => (a.location.symbol ?? "").localeCompare(b.location.symbol ?? ""));
  const key = (f: Factor) => (sort === "debt" ? f.debt : (f.metrics?.[METRIC[sort]] ?? -1));
  return xs.sort((a, b) => key(b) - key(a));
}
