// The About drawer: versions, weights, thresholds and the limits of the
// numbers (the previous page's notes and "How to read" text).
import { useEffect } from "react";
import { useMeta, useMetrics } from "../core/data";
import { fmt } from "../lib/format";
import { Close } from "./icons";

export function About({ onClose }: { onClose: () => void }) {
  const meta = useMeta().data;
  const m = useMetrics().data;
  useEffect(() => {
    const esc = (e: KeyboardEvent) => e.key === "Escape" && onClose();
    addEventListener("keydown", esc);
    return () => removeEventListener("keydown", esc);
  }, [onClose]);
  if (!meta) return null;
  const t = meta.thresholds;
  const w = meta.weights;
  return (
    <>
      <div className="drawer-backdrop" onClick={onClose} />
      <aside className="drawer" aria-label="About these numbers">
        <div style={{ display: "flex", alignItems: "center" }}>
          <h2 style={{ flex: 1 }}>About these numbers</h2>
          <button className="btn icon" aria-label="Close" onClick={onClose}>
            <Close />
          </button>
        </div>
        <div className="chips">
          <span className="chip">tool {meta.tool}</span>
          <span className="chip">{meta.scoring}</span>
          <span className="chip">{meta.rules}</span>
          <span className="chip">{meta.clones}</span>
          <span className="chip">defects {meta.defects ? "measured for every commit" : "checked on demand"}</span>
          <span className="chip">near-duplication {meta.nearDuplication ? "on" : "off"}</span>
        </div>
        <section>
          <h3>Debt</h3>
          <ul>
            <li>
              A function with cyclomatic complexity ≥ {t.cyclomatic} adds cc − {t.cyclomatic - 1} points × {w.cyclomaticExcessPoint}.
            </li>
            <li>
              Nesting depth ≥ {t.nesting} adds {w.nestingExcessLevel} × (depth − {t.nesting - 1}).
            </li>
            <li>
              Copies in a structural clone group (≥ {t.minCloneTokens} tokens) add 1 point per {w.redundantTokensPerPoint} redundant tokens.
            </li>
            <li>
              An import that breaks a configured rule adds {w.ruleViolationEdge}; a module in an import cycle adds {w.cycleUnit}.
            </li>
            <li>The score is 100 × h ÷ (h + debt per 1,000 lines), h = {w.halfScoreDensity}. Use it to compare, not as a target.</li>
          </ul>
        </section>
        <section>
          <h3>Limits</h3>
          <ul>
            {meta.notes.map((n) => (
              <li key={n}>{n}</li>
            ))}
            <li>Severity bands (critical ≥ 25, high ≥ 10, medium ≥ 4 debt points per item) are a reading aid, not calibrated.</li>
          </ul>
        </section>
        {m && (
          <section>
            <h3>This session</h3>
            <ul>
              <li>
                {fmt(m.measuredOverall)} commits measured so far; per-file results cached: {fmt(m.cache.entries)} ({fmt(m.cache.hits)} reused).
              </li>
            </ul>
          </section>
        )}
      </aside>
    </>
  );
}
