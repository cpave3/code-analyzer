// Provenance grouping (mirrors internal/attribution.GroupKey and the
// previous dashboard): each step is keyed by the chosen fields of the
// commits it brings in.
import type { Commit, ProvenanceRecord } from "../core/api";

export function fieldValue(r: ProvenanceRecord, f: string): string {
  if (f === "model" || f === "agent") return (r.model ?? []).join("+");
  if (f === "harness") return (r.harness ?? []).join("+");
  return r.meta?.[f] ?? "";
}

export function recordKey(r: ProvenanceRecord, fields: string[]): string {
  return fields.map((f) => f + "=" + (fieldValue(r, f) || "(none)")).join(", ");
}

export function stepKey(c: Commit, fields: string[]): string {
  const rs = (c.provenance ?? []).filter((r) => !(r.merge && r.source === "none"));
  if (!rs.length) return c.provenance?.length ? recordKey(c.provenance[0], fields) : "(no commits)";
  const k = recordKey(rs[0], fields);
  return rs.every((r) => recordKey(r, fields) === k) ? k : "(mixed)";
}

// Categorical palette, checked for contrast on both themes; greys are
// reserved for mixed and unattributed steps.
const PALETTE = ["#3b82f6", "#e0730b", "#0f9d8a", "#c0368a", "#7c5cd6", "#8a6a12", "#2a8fc9", "#c2410c"];
const MIXED = "#8b8f98";
const NONE = "#a7abb3";

export function isNoneKey(k: string): boolean {
  return /^[^=]*=\(none\)(, [^=]*=\(none\))*$/.test(k) || k === "(no commits)";
}

// keyColors assigns each key a colour in first-seen order (stable for a
// given history and field choice).
export function keyColors(keys: string[]): Map<string, string> {
  const out = new Map<string, string>();
  let n = 0;
  for (const k of keys) {
    if (out.has(k)) continue;
    out.set(k, k.includes("(mixed)") ? MIXED : isNoneKey(k) ? NONE : PALETTE[n++ % PALETTE.length]);
  }
  return out;
}
