// Number and text formatting shared by every view.

const nf = (d: number) => new Intl.NumberFormat("en-AU", { minimumFractionDigits: d, maximumFractionDigits: d });

export function fmt(v: number | null | undefined, d = 0): string {
  if (v == null || Number.isNaN(v)) return "–";
  return nf(d).format(v);
}

// signed formats a change with an explicit sign and a true minus.
export function signed(v: number | null | undefined, d = 0): string {
  if (v == null || Number.isNaN(v)) return "–";
  if (Math.abs(v) < 0.5 * 10 ** -d) return fmt(0, d);
  return (v > 0 ? "+" : "−") + fmt(Math.abs(v), d);
}

export const short = (sha: string) => sha.slice(0, 8);

export function day(iso: string): string {
  const t = new Date(iso);
  return Number.isNaN(t.getTime()) ? iso.slice(0, 10) : t.toLocaleDateString("en-AU", { day: "numeric", month: "short", year: "numeric" });
}

export function stamp(iso: string): string {
  return iso.replace("T", " ").slice(0, 16);
}

// moduleLabel drops the language prefix of a graph module id ("go:internal/x").
export function moduleLabel(id: string): string {
  const i = id.indexOf(":");
  return i >= 0 ? id.slice(i + 1) : id;
}

export const pct = (v: number, d = 1) => fmt(100 * v, d) + "%";
