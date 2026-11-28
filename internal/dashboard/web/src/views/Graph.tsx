// Graph: the previous page's Dependencies (who uses what, layered) and
// Network (force layout) views, drawn with Sigma.js (WebGL) on graphology.
// The Structure (folders) view is the Modules page's treemap.
import { useEffect, useMemo, useRef, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router";
import Sigma from "sigma";
import { extremityArrow, pathLine } from "sigma/rendering";
import type { GraphCallable, GraphData, GraphModule } from "../core/api";
import { ErrorBanner, Skeleton } from "../components/states";
import { useGraph, useMeta, useMetrics, useSelection } from "../core/data";
import { compare, functionKey, moduleKey, type Comparison } from "../lib/changes";
import { depsGraph, edgeStyle, networkGraph, nodeStyle, type EdgeAttrs, type NodeAttrs, type SizeBy, type StyleState } from "../lib/sigmaGraph";
import { buildDeps, type DepGraph } from "../lib/deps";
import { fmt, moduleLabel, short, signed } from "../lib/format";
import { COMPONENTS, usable } from "../lib/history";
import { useTheme } from "../core/theme";

type Mode = "deps" | "network";

interface Opts {
  mode: Mode;
  sizeBy: SizeBy;
  marks: boolean;
  functions: boolean;
  clones: boolean;
  imports: boolean;
  labels: boolean;
  minDebt: number;
  find: string;
}

export default function GraphView() {
  const meta = useMeta().data!;
  const { byIndex } = useMetrics();
  const { index, commit } = useSelection();
  const graph = useGraph(commit?.sha);
  const prevSha = index > 0 && usable(byIndex[index - 1]) ? meta.commits[index - 1].sha : undefined;
  const prev = useGraph(prevSha);
  const [o, setO] = useState<Opts>({ mode: "deps", sizeBy: "lines", marks: true, functions: true, clones: true, imports: true, labels: true, minDebt: 5, find: "" });
  const [expanded, setExpanded] = useState<Set<string>>(new Set());
  const [picked, setPicked] = useState<string | null>(null);
  const [layoutSeq, setLayoutSeq] = useState(0);
  const set = (p: Partial<Opts>) => setO((x) => ({ ...x, ...p }));
  const deps = useMemo(() => (graph.data ? buildDeps(graph.data, expanded) : null), [graph.data, expanded]);
  const cmp = useMemo(() => (graph.data && prev.data ? compare(graph.data, prev.data) : null), [graph.data, prev.data]);
  const switchMode = (mode: Mode) => {
    setPicked(null);
    set({ mode });
  };
  const expand = (key: string) => {
    setExpanded((s) => new Set([...s, key]));
    setPicked(null);
  };
  const collapse = (key: string) => {
    let k = key.replace(/\/\(own files\)$/, "");
    while (k && !expanded.has(k)) k = k.includes("/") ? k.slice(0, k.lastIndexOf("/")) : "";
    if (!k) return;
    setExpanded((s) => new Set([...s].filter((e) => e !== k && !e.startsWith(k + "/"))));
    setPicked(null);
  };
  if (!commit) return null;
  return (
    <div style={{ display: "flex", flexWrap: "wrap", minHeight: "calc(100vh - 61px)" }}>
      <Rail o={o} set={set} switchMode={switchMode} graph={graph.data} deps={deps} onCollapseAll={() => setExpanded(new Set())} onRelayout={() => setLayoutSeq((n) => n + 1)} />
      <section aria-label="Graph" style={{ flex: "999 1 560px", minWidth: 0, position: "relative", minHeight: 640, background: "var(--bg)" }}>
        {graph.isError ? (
          <div style={{ padding: 24 }}>
            <ErrorBanner title="The graph for this commit could not be built" error={graph.error} onRetry={() => graph.refetch()} />
          </div>
        ) : !graph.data || !deps ? (
          <div style={{ padding: 24 }} role="status">
            <p className="dim">Building the module graph for {short(commit.sha)}…</p>
            <Skeleton height={520} />
          </div>
        ) : (
          <SigmaCanvas data={graph.data} deps={deps} o={o} cmp={cmp} picked={picked} setPicked={setPicked} onExpand={expand} layoutSeq={layoutSeq} />
        )}
      </section>
      <aside aria-label="Details" style={{ flex: "0 1 360px", minWidth: 280, borderLeft: "1px solid var(--line)", background: "var(--panel)", padding: 18, display: "flex", flexDirection: "column", gap: 14, fontSize: 13 }}>
        {graph.data && deps && (o.mode === "deps" ? <GroupDetail deps={deps} picked={picked} setPicked={setPicked} expanded={expanded} onExpand={expand} onCollapse={collapse} /> : <NodeDetail data={graph.data} cmp={cmp} picked={picked} setPicked={setPicked} />)}
      </aside>
    </div>
  );
}

function Rail({ o, set, switchMode, graph, deps, onCollapseAll, onRelayout }: { o: Opts; set: (p: Partial<Opts>) => void; switchMode: (m: Mode) => void; graph?: GraphData; deps: DepGraph | null; onCollapseAll: () => void; onRelayout: () => void }) {
  const { tokens } = useTheme();
  const { search } = useLocation();
  return (
    <aside aria-label="Graph controls" style={{ flex: "0 1 260px", minWidth: 220, borderRight: "1px solid var(--line)", background: "var(--panel)", padding: 18, display: "flex", flexDirection: "column", gap: 16, fontSize: 13 }}>
      <fieldset style={{ border: 0, padding: 0, margin: 0, display: "flex", flexDirection: "column", gap: 6 }}>
        <legend className="dim" style={{ fontSize: 12, marginBottom: 6 }}>
          View
        </legend>
        <label style={{ display: "flex", gap: 8 }}>
          <input type="radio" name="mode" checked={o.mode === "deps"} onChange={() => switchMode("deps")} /> Dependencies (who uses what)
        </label>
        <label style={{ display: "flex", gap: 8 }}>
          <input type="radio" name="mode" checked={o.mode === "network"} onChange={() => switchMode("network")} /> Network (force layout)
        </label>
        <Link to={{ pathname: "/modules", search }} style={{ fontSize: 12.5 }}>
          Structure (folders) is on Modules →
        </Link>
      </fieldset>
      <label style={{ display: "flex", flexDirection: "column", gap: 6 }}>
        <span className="dim" style={{ fontSize: 12 }}>
          Find
        </span>
        <input type="search" placeholder={o.mode === "deps" ? "Group" : "Module or function"} value={o.find} onChange={(e) => set({ find: e.target.value })} />
      </label>
      <label style={{ display: "flex", flexDirection: "column", gap: 6 }}>
        <span className="dim" style={{ fontSize: 12 }}>
          Size by
        </span>
        <select value={o.sizeBy} onChange={(e) => set({ sizeBy: e.target.value as SizeBy })}>
          <option value="lines">Source lines</option>
          <option value="debt">Debt</option>
        </select>
      </label>
      <fieldset style={{ border: 0, padding: 0, margin: 0, display: "flex", flexDirection: "column", gap: 6 }}>
        <legend className="dim" style={{ fontSize: 12, marginBottom: 6 }}>
          Show
        </legend>
        <label style={{ display: "flex", gap: 8 }}>
          <input type="checkbox" checked={o.marks} onChange={(e) => set({ marks: e.target.checked })} /> Mark import cycles and rule breaks
        </label>
        <label style={{ display: "flex", gap: 8 }}>
          <input type="checkbox" checked={o.labels} onChange={(e) => set({ labels: e.target.checked })} /> Labels
        </label>
        {o.mode === "network" && (
          <>
            <label style={{ display: "flex", gap: 8 }}>
              <input type="checkbox" checked={o.functions} onChange={(e) => set({ functions: e.target.checked })} /> Functions with the most debt
            </label>
            <label style={{ display: "flex", gap: 8 }}>
              <input type="checkbox" checked={o.clones} onChange={(e) => set({ clones: e.target.checked })} /> Clone links
            </label>
            <label style={{ display: "flex", gap: 8 }}>
              <input type="checkbox" checked={o.imports} onChange={(e) => set({ imports: e.target.checked })} /> Import edges
            </label>
          </>
        )}
      </fieldset>
      {o.mode === "network" && (
        <label style={{ display: "flex", flexDirection: "column", gap: 6 }}>
          <span className="dim" style={{ fontSize: 12 }}>
            Min function debt · {o.minDebt}
          </span>
          <input type="range" min={0} max={60} value={o.minDebt} onChange={(e) => set({ minDebt: +e.target.value })} style={{ accentColor: "var(--accent)" }} />
        </label>
      )}
      <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
        {o.mode === "deps" ? (
          <button className="btn small" onClick={onCollapseAll}>
            Collapse all
          </button>
        ) : (
          <button className="btn small" onClick={onRelayout}>
            Re-run layout
          </button>
        )}
      </div>
      <div className="dim" style={{ marginTop: "auto", fontSize: 12, display: "flex", flexDirection: "column", gap: 8 }}>
        {graph && deps && <div>{o.mode === "deps" ? `${fmt(deps.groups.length)} groups · ${fmt(deps.edges.length)} links between them` : `${fmt(graph.modules.length)} modules · ${fmt(graph.callables.length)} functions with the most debt · ${fmt(graph.imports.length)} imports`}</div>}
        <div style={{ display: "flex", alignItems: "center", gap: 6, flexWrap: "wrap" }}>
          <span>Debt per kLOC</span>
          {tokens.ramp.map((c) => (
            <span key={c} className="sw" style={{ background: c, marginRight: 0 }} />
          ))}
        </div>
        <div>
          <span style={{ color: tokens.bad }}>■</span> in an import cycle or breaking a rule · ring: <span style={{ color: tokens.bad }}>more debt</span>, <span style={{ color: tokens.good }}>less</span>, <span style={{ color: tokens.accent }}>new</span> since the previous commit
        </div>
        <div>Scroll to zoom, drag to pan; click for details{o.mode === "deps" ? ", double-click a group to expand it" : ""}.</div>
      </div>
    </aside>
  );
}

function SigmaCanvas({ data, deps, o, cmp, picked, setPicked, onExpand, layoutSeq }: { data: GraphData; deps: DepGraph; o: Opts; cmp: Comparison | null; picked: string | null; setPicked: (k: string | null) => void; onExpand: (k: string) => void; layoutSeq: number }) {
  const { tokens } = useTheme();
  const el = useRef<HTMLDivElement>(null);
  const sigma = useRef<Sigma<NodeAttrs, EdgeAttrs> | null>(null);
  const hover = useRef<string | null>(null);
  const style = useRef<() => StyleState>(() => ({ focus: null, picked: null, marks: true, labels: true, find: "" }));
  style.current = () => ({ focus: picked ?? hover.current, picked, marks: o.marks, labels: o.labels, find: o.find });
  const theme = useRef(tokens);
  theme.current = tokens;
  const handlers = useRef({ setPicked, onExpand });
  handlers.current = { setPicked, onExpand };
  // Rebuild the graph when the data, mode or what is shown changes; other
  // options only restyle it through the reducers.
  const g = useMemo(
    () => (o.mode === "deps" ? depsGraph(deps, o.sizeBy, tokens) : networkGraph(data, o, tokens, cmp, layoutSeq)),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [data, deps, o.mode, o.sizeBy, o.functions, o.clones, o.imports, o.minDebt, tokens, cmp, layoutSeq],
  );
  useEffect(() => {
    const s = new Sigma<NodeAttrs, EdgeAttrs>(g, el.current!, {
      primitives: { edges: { paths: [pathLine()], extremities: [extremityArrow({ lengthRatio: 3, widthRatio: 2 })] } },
      settings: { renderEdgeLabels: false, labelRenderedSizeThreshold: 5, stagePadding: 40, itemSizesReference: "screen" } as never,
      nodeReducer: (key, d, attrs) => ({ ...d, ...nodeStyle(g, key, attrs, style.current(), theme.current) }),
      edgeReducer: (key, d, attrs) => ({ ...d, ...edgeStyle(g, key, attrs, style.current(), theme.current) }),
    });
    sigma.current = s;
    const refresh = () => s.refresh({ skipIndexation: true } as never);
    s.on("clickNode", ({ node }: { node: string }) => handlers.current.setPicked(node));
    s.on("doubleClickNode", (e: { node: string; preventSigmaDefault?: () => void }) => {
      e.preventSigmaDefault?.();
      if (g.getNodeAttribute(e.node, "kind") === "group") handlers.current.onExpand(e.node);
    });
    s.on("clickStage", () => handlers.current.setPicked(null));
    s.on("enterNode", ({ node }: { node: string }) => {
      hover.current = node;
      refresh();
    });
    s.on("leaveNode", () => {
      hover.current = null;
      refresh();
    });
    return () => {
      s.kill();
      sigma.current = null;
    };
  }, [g]);
  useEffect(() => {
    sigma.current?.refresh({ skipIndexation: true } as never);
  }, [picked, o.marks, o.labels, o.find, tokens]);
  const camera = () => sigma.current?.getCamera();
  return (
    <>
      <div ref={el} style={{ position: "absolute", inset: 0 }} role="img" aria-label={o.mode === "deps" ? "Layered dependency graph of module groups" : "Force layout of modules and the functions with the most debt"} />
      <div style={{ position: "absolute", right: 14, bottom: 14, display: "flex", flexDirection: "column", gap: 4 }}>
        <button className="btn icon" aria-label="Zoom in" onClick={() => camera()?.zoomIn({ duration: 200 })}>
          +
        </button>
        <button className="btn icon" aria-label="Zoom out" onClick={() => camera()?.zoomOut({ duration: 200 })}>
          −
        </button>
        <button className="btn icon" aria-label="Fit" title="Fit" onClick={() => camera()?.reset({ duration: 200 })}>
          ⤢
        </button>
      </div>
    </>
  );
}

function GroupDetail({ deps, picked, setPicked, expanded, onExpand, onCollapse }: { deps: DepGraph; picked: string | null; setPicked: (k: string) => void; expanded: Set<string>; onExpand: (k: string) => void; onCollapse: (k: string) => void }) {
  const { search } = useLocation();
  const g = deps.groups.find((x) => x.key === picked);
  if (!g) return <p className="dim">Click a group for what it uses and what uses it. Double-click (or Expand) to open it into subfolders.</p>;
  const uses = deps.edges.filter((e) => e.from === g.key).sort((a, b) => b.count - a.count);
  const usedBy = deps.edges.filter((e) => e.to === g.key).sort((a, b) => b.count - a.count);
  const parentExpanded = [...expanded].some((k) => g.key.startsWith(k + "/"));
  const list = (es: typeof uses, dir: "out" | "in") =>
    es.length ? (
      <ul style={{ margin: 0, paddingLeft: 16, display: "flex", flexDirection: "column", gap: 6 }}>
        {es.slice(0, 12).map((e) => {
          const other = dir === "out" ? e.to : e.from;
          return (
            <li key={other}>
              <button className="mono" style={{ background: "none", border: 0, padding: 0, color: "var(--accent)", cursor: "pointer", textAlign: "left" }} onClick={() => setPicked(other)}>
                {other}
              </button>{" "}
              · <b>{e.count}</b> import{e.count === 1 ? "" : "s"}
              {e.cycle && <span className="bad"> · in a cycle</span>}
              {e.rule && <span className="bad"> · breaks a rule</span>}
              {e.sites.length > 0 && (
                <div className="mono dim" style={{ fontSize: 11.5 }}>
                  {e.sites.join(" · ")}
                </div>
              )}
            </li>
          );
        })}
      </ul>
    ) : (
      <p className="dim" style={{ margin: 0 }}>
        none
      </p>
    );
  return (
    <>
      <h2 className="mono" style={{ margin: 0, fontSize: 15, fontWeight: 600, wordBreak: "break-all" }}>
        {g.key}
      </h2>
      <div>
        {fmt(g.lines)} lines · debt {fmt(g.debt, 1)} ({fmt(g.density, 1)} per kLOC) · {g.modules.length} module{g.modules.length === 1 ? "" : "s"}
        {g.inCycle && <span className="bad"> · in an import cycle at this level</span>}
      </div>
      <div style={{ display: "flex", gap: 8, flexWrap: "wrap" }}>
        {g.expandable && (
          <button className="btn small" onClick={() => onExpand(g.key)}>
            Expand into subfolders
          </button>
        )}
        {parentExpanded && (
          <button className="btn small" onClick={() => onCollapse(g.key)}>
            Collapse to parent
          </button>
        )}
        <Link className="btn small" style={{ textDecoration: "none" }} to={{ pathname: "/modules/" + g.key.replace(/\/\(own files\)$/, ""), search }}>
          Open folder report
        </Link>
      </div>
      <div>
        <h3>Uses ({uses.length})</h3>
        {list(uses, "out")}
      </div>
      <div>
        <h3>Used by ({usedBy.length})</h3>
        {list(usedBy, "in")}
      </div>
    </>
  );
}

function NodeDetail({ data, cmp, picked, setPicked }: { data: GraphData; cmp: Comparison | null; picked: string | null; setPicked: (k: string) => void }) {
  const mod = data.modules.find((m) => m.id === picked);
  if (mod) return <ModuleDetail data={data} cmp={cmp} mod={mod} setPicked={setPicked} />;
  const fn = data.callables.find((c) => c.id === picked);
  if (fn) return <FunctionDetail data={data} cmp={cmp} fn={fn} setPicked={setPicked} />;
  return <p className="dim">Click a circle for its details. Modules are sized by source lines or debt and coloured by debt per kLOC; small circles are the functions with the most debt.</p>;
}

// ChangeText says how an item's debt changed since the previous commit.
function ChangeText({ cmp, k }: { cmp: Comparison | null; k: string }) {
  const { tokens } = useTheme();
  const c = cmp?.states.get(k);
  if (!c) return null;
  if (c.state === "new") return <span style={{ color: tokens.accent }}>new in this commit</span>;
  if (c.state === "same") return <span className="dim">unchanged since the previous commit</span>;
  return <span className={c.delta > 0 ? "bad" : "good"}>{signed(c.delta, 1)} debt since the previous commit</span>;
}

function ModuleDetail({ data, cmp, mod, setPicked }: { data: GraphData; cmp: Comparison | null; mod: GraphModule; setPicked: (k: string) => void }) {
  const t = useMeta().data!.thresholds;
  const { tokens } = useTheme();
  const nav = useNavigate();
  const { search } = useLocation();
    const label = moduleLabel(mod.id);
    const fns = data.callables.filter((c) => c.module === mod.id).sort((a, b) => b.debt + b.cloneDebt - (a.debt + a.cloneDebt)).slice(0, 8);
    const out = data.imports.filter((e) => e.from === mod.id).length;
    const inn = data.imports.filter((e) => e.to === mod.id).length;
    return (
      <>
        <h2 className="mono" style={{ margin: 0, fontSize: 15, fontWeight: 600, wordBreak: "break-all" }}>
          {label} <span className="dim" style={{ fontSize: 12 }}>{mod.lang} module</span>
        </h2>
        <div><ChangeText cmp={cmp} k={moduleKey(mod.id)} /></div>
        <div className="kpis compact" style={{ gridTemplateColumns: "repeat(2, minmax(0, 1fr))" }}>
          <div className="kpi">
            <div className="label">Debt</div>
            <div className="value">{fmt(mod.debt.total, 1)}</div>
          </div>
          <div className="kpi">
            <div className="label">Per kLOC</div>
            <div className="value">{fmt(mod.density, 1)}</div>
          </div>
        </div>
        <table className="data" style={{ fontSize: 12.5 }}>
          <tbody>
            {COMPONENTS.map((k) => (
              <tr key={k}>
                <td style={{ paddingLeft: 0 }}>
                  <span className="sw" style={{ background: tokens.components[k] }} />
                  {k}
                </td>
                <td className="r" style={{ paddingRight: 0 }}>
                  {fmt(mod.debt[k], 1)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        <div>
          {fmt(mod.lines)} source lines · {mod.callables} functions, {mod.hotspots} hotspots (cc ≥ {t.cyclomatic} or depth ≥ {t.nesting}) · score {fmt(mod.score, 1)} · imports {out} module{out === 1 ? "" : "s"}, imported by {inn}
        </div>
        {fns.length > 0 && (
          <div>
            <h3>Functions with the most debt</h3>
            <div className="chips" style={{ marginTop: 6 }}>
              {fns.map((f) => (
                <button key={f.id} className="chip" onClick={() => setPicked(f.id)}>
                  {f.name} <b>{fmt(f.debt + f.cloneDebt, 1)}</b>
                </button>
              ))}
            </div>
          </div>
        )}
        <div>
          <button className="btn small" onClick={() => nav({ pathname: "/modules/" + label, search })}>
            Open folder report
          </button>
        </div>
      </>
    );
  }

function FunctionDetail({ data, cmp, fn, setPicked }: { data: GraphData; cmp: Comparison | null; fn: GraphCallable; setPicked: (k: string) => void }) {
  const t = useMeta().data!.thresholds;
    const cx = Math.max(0, fn.cc - t.cyclomatic + 1);
    const nx = fn.nest >= t.nesting ? 2 * (fn.nest - t.nesting + 1) : 0;
    const partners = data.clones.filter((g) => g.members.includes(fn.id)).flatMap((g) => g.members.filter((m) => m !== fn.id));
    return (
      <>
        <h2 className="mono" style={{ margin: 0, fontSize: 15, fontWeight: 600 }}>
          {fn.name}
        </h2>
        <div className="mono dim" style={{ fontSize: 12 }}>
          {fn.file}:{fn.line}
        </div>
        <div>
          in module{" "}
          <button className="mono" style={{ background: "none", border: 0, padding: 0, color: "var(--accent)", cursor: "pointer" }} onClick={() => setPicked(fn.module)}>
            {moduleLabel(fn.module)}
          </button>{" "}
          · <ChangeText cmp={cmp} k={functionKey(fn.file, fn.name)} />
        </div>
        <div>
          Debt <b>{fmt(fn.debt + fn.cloneDebt, 1)}</b> · cc {fn.cc} · nesting {fn.nest} · {fmt(fn.lines)} lines
        </div>
        <ul style={{ margin: 0, paddingLeft: 18 }}>
          {cx > 0 && (
            <li>
              cyclomatic complexity {fn.cc} ≥ {t.cyclomatic}: <b>{cx}</b> complexity point{cx === 1 ? "" : "s"}
            </li>
          )}
          {nx > 0 && (
            <li>
              nesting depth {fn.nest} ≥ {t.nesting}: <b>{nx}</b> nesting point{nx === 1 ? "" : "s"}
            </li>
          )}
          {fn.cloneDebt > 0 && (
            <li>
              member of a structural clone group: <b>{fmt(fn.cloneDebt, 1)}</b> duplication points
            </li>
          )}
          {!cx && !nx && !fn.cloneDebt && <li>below the hotspot thresholds</li>}
        </ul>
        {partners.length > 0 && (
          <div>
            <h3>Structural clone of</h3>
            <div className="chips" style={{ marginTop: 6 }}>
              {partners.map((p) => (
                <button key={p} className="chip mono" onClick={() => setPicked(p)}>
                  {p.split(":").pop()}
                </button>
              ))}
            </div>
          </div>
        )}
        <div className="dim">Reduce branching or nesting, or reuse one copy, to remove these points; do not reshape code only to lower the number.</div>
      </>
    );
  }
