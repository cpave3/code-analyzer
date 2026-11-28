// Modules: the folder structure as a treemap you drill into (the previous
// page's Structure view) beside that folder's report (explain --path, plus
// on-demand defect checks).
import { useVirtualizer } from "@tanstack/react-virtual";
import { useMemo, useRef, useState } from "react";
import { Link, useLocation, useNavigate, useParams } from "react-router";
import type { ExplainReport, Factor, GraphData, Metrics } from "../api";
import { EChart, type EChartsOption } from "../components/EChart";
import { EmptyState, ErrorBanner, Skeleton } from "../components/states";
import { useDefects, useExplain, useGraph, useMeta, useMetrics, useSelection } from "../data";
import { compare, moduleKey, type Comparison } from "../lib/changes";
import { fmt, pct, short } from "../lib/format";
import { COMPONENTS, usable } from "../lib/history";
import { BANDS, band, byFile, bySubfolder, componentOf, crumbs, KIND_NAMES } from "../lib/report";
import { buildTree, DENSITY_STEPS, densityBand, findNode, parentPath, sizeOf, type SizeBy, type TreeNode } from "../lib/tree";
import { useTheme, type Tokens } from "../theme";
import { CommitTitle } from "./Overview";

type ColourBy = "density" | "component";

const NO_FACTORS: Factor[] = [];

export default function Modules() {
  const path = (useParams()["*"] ?? "").replace(/\/+$/, "");
  const meta = useMeta().data!;
  const { byIndex } = useMetrics();
  const { index, commit } = useSelection();
  const m = byIndex[index];
  const graph = useGraph(commit?.sha);
  const prevSha = index > 0 && usable(byIndex[index - 1]) ? meta.commits[index - 1].sha : undefined;
  const prev = useGraph(prevSha);
  const [sizeBy, setSizeBy] = useState<SizeBy>("lines");
  const [colourBy, setColourBy] = useState<ColourBy>("density");
  const [lang, setLang] = useState("");
  if (!commit) return null;
  const langs = [...new Set((graph.data?.modules ?? []).map((x) => x.lang))].sort();
  return (
    <main className="page">
      <CommitTitle commit={commit} scoring={meta.scoring} />
      <div style={{ display: "flex", flexWrap: "wrap", alignItems: "center", gap: "10px 16px" }}>
        <nav aria-label="Folder" className="mono" style={{ display: "flex", gap: 6, alignItems: "center", flexWrap: "wrap" }}>
          {crumbs(path).map((c, i, all) => (
            <span key={c.path} style={{ display: "inline-flex", gap: 6 }}>
              {i > 0 && <span className="dim">/</span>}
              {i === all.length - 1 ? <span style={{ fontWeight: 600 }}>{c.name}</span> : <Link to={{ pathname: "/modules/" + c.path, search: location.search }}>{c.name}</Link>}
            </span>
          ))}
        </nav>
        <div className="spacer" />
        <label className="inline">
          Size by
          <select value={sizeBy} onChange={(e) => setSizeBy(e.target.value as SizeBy)}>
            <option value="lines">Source lines</option>
            <option value="debt">Debt</option>
          </select>
        </label>
        <label className="inline">
          Colour by
          <select value={colourBy} onChange={(e) => setColourBy(e.target.value as ColourBy)}>
            <option value="density">Debt per kLOC</option>
            <option value="component">Dominant component</option>
          </select>
        </label>
        {langs.length > 1 && (
          <label className="inline">
            Language
            <select value={lang} onChange={(e) => setLang(e.target.value)}>
              <option value="">All</option>
              {langs.map((l) => (
                <option key={l} value={l}>
                  {l}
                </option>
              ))}
            </select>
          </label>
        )}
      </div>
      {graph.isError ? (
        <ErrorBanner title="The module map could not be built" error={graph.error} onRetry={() => graph.refetch()} />
      ) : (
        <div className="row">
          <div className="main-col">
            <ModuleMap graph={graph.data} prev={prevSha ? (prev.data ?? null) : null} path={path} sizeBy={sizeBy} colourBy={colourBy} lang={lang} />
            <ScopeTables sha={commit.sha} path={path} />
          </div>
          <div className="side-col" style={{ flexBasis: 420 }}>
            <FolderReport sha={commit.sha} path={path} repo={m} />
          </div>
        </div>
      )}
    </main>
  );
}

function treemapData(n: TreeNode, sizeBy: SizeBy, colourBy: ColourBy, t: Tokens, cmp: Comparison | null, focus: string, lang: string, marked: Set<string>): unknown[] {
  const kids = n.children.filter((c) => !lang || !c.module || c.module.lang === lang);
  return kids
    .map((c) => {
      const children = c.module ? [] : treemapData(c, sizeBy, colourBy, t, cmp, focus, lang, marked);
      if (!c.module && !children.length) return null;
      let color = t.ramp[densityBand(c.density)];
      if (colourBy === "component" && c.module) {
        const top = COMPONENTS.map((k) => [k, c.module!.debt[k]] as const).sort((a, b) => b[1] - a[1])[0];
        color = top[1] > 0 ? t.components[top[0]] : t.track;
      }
      const st = c.module && cmp ? cmp.states.get(moduleKey(c.module.id))?.state : undefined;
      const ring = c.path === focus && c.module ? t.fg : st === "up" ? t.bad : st === "down" ? t.good : st === "new" ? t.accent : undefined;
      const bad = !!c.module && marked.has(c.module.id);
      return {
        name: bad ? "↻ " + c.name : c.name,
        path: c.path,
        bad,
        value: sizeOf(c, sizeBy),
        debt: c.debt,
        lines: c.lines,
        density: c.density,
        state: st,
        children: children.length ? children : undefined,
        itemStyle: c.module ? { color, borderColor: ring ?? t.panel, borderWidth: ring ? 3 : 1 } : { borderColor: t.line, color: t.bg },
        label: c.module ? { color: inkFor(color) } : undefined,
      };
    })
    .filter(Boolean);
}

function ModuleMap({ graph, prev, path, sizeBy, colourBy, lang }: { graph?: GraphData; prev: GraphData | null; path: string; sizeBy: SizeBy; colourBy: ColourBy; lang: string }) {
  const { tokens } = useTheme();
  const nav = useNavigate();
  const { search } = useLocation();
  const tree = useMemo(() => (graph ? buildTree(graph.modules) : null), [graph]);
  const cmp = useMemo(() => (graph && prev ? compare(graph, prev) : null), [graph, prev]);
  const marked = useMemo(() => {
    const out = new Set<string>();
    for (const e of graph?.imports ?? []) if (e.cycle || e.violation) out.add(e.from).add(e.to);
    return out;
  }, [graph]);
  const shown = useMemo(() => {
    if (!tree) return null;
    const n = findNode(tree, path);
    if (n && !n.module) return n;
    if (n?.module) return findNode(tree, parentPath(path)) ?? tree;
    let p = parentPath(path);
    while (p && !findNode(tree, p)) p = parentPath(p);
    return findNode(tree, p) ?? tree;
  }, [tree, path]);
  const option = useMemo<EChartsOption | null>(() => {
    if (!shown) return null;
    return {
      animation: false,
      tooltip: {
        confine: true,
        formatter: (p: { data: { path: string; debt: number; lines: number; density: number; state?: string; bad?: boolean } }) => {
          const d = p.data;
          const bad = d.bad ? " · in an import cycle or breaking a rule" : "";
          const st = d.state === "up" ? " · more debt than the previous commit" : d.state === "down" ? " · less debt than the previous commit" : d.state === "new" ? " · new in this commit" : "";
          return `<b>${esc(d.path || "(repository)")}</b><br>${fmt(d.lines)} lines · debt ${fmt(d.debt, 1)} · ${fmt(d.density, 1)} per kLOC${st}${bad}`;
        },
      },
      series: [
        {
          type: "treemap",
          data: treemapData(shown, sizeBy, colourBy, tokens, cmp, path, lang, marked),
          roam: false,
          nodeClick: false,
          breadcrumb: { show: false },
          width: "100%",
          height: "100%",
          top: 0,
          left: 0,
          squareRatio: 1.2,
          label: { show: true, formatter: "{b}", color: "#111", fontSize: 11, overflow: "truncate" },
          upperLabel: { show: true, height: 18, color: tokens.fg, fontSize: 11, fontWeight: 600 },
          itemStyle: { gapWidth: 2 },
          levels: [{ itemStyle: { borderWidth: 0, gapWidth: 3 } }, { itemStyle: { borderWidth: 1, gapWidth: 2 } }, { itemStyle: { gapWidth: 1 } }],
        },
      ],
    };
  }, [shown, sizeBy, colourBy, tokens, cmp, path, lang, marked]);
  const onClick = (p: unknown) => {
    const d = (p as { data?: { path?: string } }).data;
    if (d?.path != null) nav({ pathname: "/modules/" + d.path, search });
  };
  return (
    <section className="card" aria-labelledby="map">
      <div className="head">
        <h2 id="map">Module map</h2>
        <span className="note">{tree ? `${fmt(graph!.modules.length)} modules · click a box to open it` : ""}</span>
      </div>
      <div className="body" style={{ display: "flex", flexDirection: "column", gap: 10 }}>
        {option ? <EChart option={option} height={520} onEvents={{ click: onClick }} label="Treemap of folders and modules, sized by source lines or debt and coloured by debt" /> : <Skeleton height={520} />}
        <div className="legend" style={{ alignItems: "center" }}>
          {colourBy === "density" ? (
            <>
              <span>Debt per kLOC</span>
              {tokens.ramp.map((c, i) => (
                <span key={c}>
                  <span className="sw" style={{ background: c }} />
                  {i === 0 ? `< ${DENSITY_STEPS[0]}` : i === tokens.ramp.length - 1 ? `≥ ${DENSITY_STEPS[i - 1]}` : `${DENSITY_STEPS[i - 1]}–${DENSITY_STEPS[i]}`}
                </span>
              ))}
            </>
          ) : (
            COMPONENTS.map((k) => (
              <span key={k}>
                <span className="sw" style={{ background: tokens.components[k] }} />
                {k}
              </span>
            ))
          )}
          {marked.size > 0 && <span>↻ in an import cycle or breaking a rule (see Graph)</span>}
          {cmp && (
            <span>
              Outline vs previous commit: <span style={{ color: tokens.bad }}>{cmp.counts.up} more debt</span> · <span style={{ color: tokens.good }}>{cmp.counts.down} less</span> · <span style={{ color: tokens.accent }}>{cmp.counts.new} new</span>
            </span>
          )}
        </div>
      </div>
    </section>
  );
}

// inkFor picks dark or light text for a tile colour by its luminance.
function inkFor(hex: string): string {
  const m = /^#([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})$/i.exec(hex);
  if (!m) return "#111";
  const [r, g, b] = [m[1], m[2], m[3]].map((x) => parseInt(x, 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4));
  return 0.2126 * r + 0.7152 * g + 0.0722 * b > 0.3 ? "#111" : "#fff";
}

const esc = (s: string) => s.replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" })[c]!);

function FolderReport({ sha, path, repo }: { sha: string; path: string; repo?: Metrics }) {
  const { tokens } = useTheme();
  const r = useExplain(sha, path);
  if (r.isPending)
    return (
      <section className="card pad" aria-busy="true">
        <h2 className="mono">{path || "(whole repository)"}</h2>
        <p className="dim" role="status">
          Analysing this folder…
        </p>
        <Skeleton height={320} />
      </section>
    );
  if (r.isError)
    return (
      <section className="card pad">
        <ErrorBanner title="The folder report could not be built" error={r.error} onRetry={() => r.refetch()} />
      </section>
    );
  const e = r.data;
  const t = e.totals;
  if (!t.sourceFiles) return <EmptyState title={`No measured source under ${path || "the repository"}`}>Pick a coloured box in the map, or go up a level.</EmptyState>;
  const factors = e.factors ?? [];
  const d = t.debt;
  const repoDensity = repo?.debtPerKLOC;
  return (
    <section className="card" aria-labelledby="scope">
      <div className="body" style={{ display: "flex", flexDirection: "column", gap: 12, borderBottom: `1px solid ${tokens.line}` }}>
        <div style={{ display: "flex", alignItems: "baseline", gap: 8, flexWrap: "wrap" }}>
          <h2 id="scope" className="mono" style={{ margin: 0, fontSize: 15, fontWeight: 600 }}>
            {path || "(whole repository)"}
          </h2>
          <span className="dim" style={{ fontSize: 12 }}>
            {fmt(t.sourceFiles)} files · {fmt(t.sourceLines)} lines · {fmt(t.modules.length)} module{t.modules.length === 1 ? "" : "s"}
          </span>
        </div>
        <div className="kpis compact" style={{ gridTemplateColumns: "repeat(3, minmax(0, 1fr))" }}>
          <div className="kpi">
            <div className="label">Debt</div>
            <div className="value">{fmt(d.total, 2)}</div>
          </div>
          <div className="kpi">
            <div className="label">Per kLOC</div>
            <div className="value">{fmt(t.debtPerKLOC, 1)}</div>
            {repoDensity != null && t.debtPerKLOC != null && path && (
              <div className="sub" style={{ color: t.debtPerKLOC > repoDensity ? tokens.bad : tokens.good }}>
                {t.debtPerKLOC > repoDensity ? `${fmt(t.debtPerKLOC / repoDensity, 1)}× the repository` : `below the repository's ${fmt(repoDensity, 1)}`}
              </div>
            )}
          </div>
          <div className="kpi">
            <div className="label">Hotspots</div>
            <div className="value">{fmt(t.hotspots)}</div>
            <div className="sub">of {fmt(t.callables)} callables</div>
          </div>
        </div>
        {d.total > 0 && (
          <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
            <div className="stack" aria-hidden="true">
              {COMPONENTS.filter((k) => d[k] > 0).map((k) => (
                <div key={k} style={{ width: `${(100 * d[k]) / d.total}%`, background: tokens.components[k] }} />
              ))}
            </div>
            <div className="legend">
              {COMPONENTS.filter((k) => d[k] > 0).map((k) => (
                <span key={k}>
                  <span className="sw" style={{ background: tokens.components[k] }} />
                  {k} {fmt(d[k], 1)} ({pct(d[k] / d.total, 0)})
                </span>
              ))}
            </div>
          </div>
        )}
        {factors.length > 1 && <Concentration e={e} />}
        <Subfolders factors={factors} path={path} total={t.sourceLines} />
      </div>
      <Items factors={factors} />
      <Defects sha={sha} path={path} />
      <div className="body" style={{ display: "flex", flexDirection: "column", gap: 4 }}>
        <details className="item">
          <summary>How to fix well</summary>
          <div className="body">
            <ul>
              {e.guidance.map((g) => (
                <li key={g}>{g}</li>
              ))}
            </ul>
          </div>
        </details>
      </div>
    </section>
  );
}

function Concentration({ e }: { e: ExplainReport }) {
  const { tokens } = useTheme();
  const fs = e.factors ?? NO_FACTORS;
  const option = useMemo<EChartsOption>(() => {
    const pts = fs.map((f, i) => [i + 1, Math.round(1000 * f.cumulativeShare) / 10]);
    const c = e.concentration;
    return {
      animation: false,
      grid: { left: 4, right: 10, top: 8, bottom: 4, containLabel: true },
      xAxis: { type: "value", min: 1, max: fs.length, axisLabel: { color: tokens.dim, fontSize: 10 }, splitLine: { show: false } },
      yAxis: { type: "value", min: 0, max: 100, interval: 50, axisLabel: { color: tokens.dim, fontSize: 10, formatter: "{value}%" }, splitLine: { lineStyle: { color: tokens.line, type: "dashed" } } },
      tooltip: { trigger: "axis", formatter: (p: { value: number[] }[]) => `largest ${p[0].value[0]} items: ${p[0].value[1]}% of the debt` },
      series: [
        {
          type: "line",
          data: pts,
          showSymbol: false,
          lineStyle: { color: tokens.accent, width: 2 },
          markPoint: {
            symbol: "circle",
            symbolSize: 7,
            itemStyle: { color: tokens.accent },
            label: { position: "right", color: tokens.dim, fontSize: 11, formatter: (p: { value: string }) => p.value },
            data: [
              { coord: [c.factorsFor50Percent, 50], value: `${c.factorsFor50Percent} → 50%` },
              { coord: [c.factorsFor80Percent, 80], value: `${c.factorsFor80Percent} → 80%` },
            ],
          },
        },
      ],
    };
  }, [fs, e.concentration, tokens]);
  return (
    <div>
      <div className="dim" style={{ fontSize: 12, marginBottom: 2 }}>
        Share of the debt carried by the largest items ({fmt(e.concentration.factors)} in all)
      </div>
      <EChart option={option} height={110} label={`${e.concentration.factorsFor50Percent} of ${e.concentration.factors} items carry half the debt`} />
    </div>
  );
}

function Subfolders({ factors, path, total }: { factors: Factor[]; path: string; total: number }) {
  const { tokens } = useTheme();
  const subs = useMemo(() => bySubfolder(factors, path), [factors, path]);
  if (subs.length < 2) return null;
  const max = subs[0].debt || 1;
  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
      <div className="dim" style={{ fontSize: 12 }}>
        Where in this folder
      </div>
      {subs.slice(0, 10).map((s) => (
        <div key={s.name} style={{ display: "grid", gridTemplateColumns: "minmax(0, 150px) 1fr 60px", gap: 8, alignItems: "center", fontSize: 12.5 }}>
          <span className="mono" style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }} title={s.name}>
            {s.path ? <Link to={{ pathname: "/modules/" + s.path, search: location.search }}>{s.name}</Link> : s.name}
          </span>
          <div className="meter">
            <div style={{ width: `${(100 * s.debt) / max}%`, background: tokens.ramp[densityBand((s.debt * 1000) / Math.max(1, total))] }} />
          </div>
          <span className="num" style={{ textAlign: "right" }}>
            {fmt(s.debt, 1)}
          </span>
        </div>
      ))}
    </div>
  );
}

function Items({ factors }: { factors: Factor[] }) {
  const { tokens } = useTheme();
  const groups = useMemo(() => {
    const g = new Map(BANDS.map((b) => [b, [] as Factor[]]));
    for (const f of factors) g.get(band(f.debt))!.push(f);
    return g;
  }, [factors]);
  if (!factors.length) return <div className="body dim">No debt items here.</div>;
  return (
    <div className="body" style={{ display: "flex", flexDirection: "column", gap: 8, borderBottom: `1px solid ${tokens.line}` }}>
      <h3>Debt items by severity</h3>
      {BANDS.filter((b) => groups.get(b)!.length).map((b) => {
        const xs = groups.get(b)!;
        return (
          <details key={b} className="item" open={b === "critical" || b === "high"}>
            <summary>
              <span className={"sev " + b}>{b}</span>
              {fmt(xs.length)} item{xs.length === 1 ? "" : "s"} · {fmt(xs.reduce((a, f) => a + f.debt, 0), 1)} debt
            </summary>
            <div style={{ paddingLeft: 12 }}>
              {xs.slice(0, 40).map((f) => (
                <FactorItem key={f.id} f={f} />
              ))}
              {xs.length > 40 && <p className="dim">{fmt(xs.length - 40)} more in the table under the map.</p>}
            </div>
          </details>
        );
      })}
    </div>
  );
}

function FactorItem({ f }: { f: Factor }) {
  const { tokens } = useTheme();
  const rem = f.remediation;
  const loc = f.location;
  const [copied, setCopied] = useState(false);
  const prompt = () => {
    const text = [
      `Reduce the ${KIND_NAMES[f.kind] ?? f.kind} debt of ${loc.symbol || f.module} at ${loc.file}:${loc.startLine}.`,
      f.explanation,
      rem.strategy && `Fix: ${rem.strategy}. ${rem.target}`,
      ...rem.steps.map((s) => "- " + s),
      rem.acceptance.length ? "Done when: " + rem.acceptance.join("; ") : "",
      rem.caution ? "Caution: " + rem.caution : "",
    ]
      .filter(Boolean)
      .join("\n");
    navigator.clipboard?.writeText(text).then(
      () => setCopied(true),
      () => setCopied(false),
    );
  };
  return (
    <details className="item">
      <summary>
        <span className="tag" style={{ background: tokens.components[componentOf(f.kind)] }}>
          {KIND_NAMES[f.kind] ?? f.kind}
        </span>
        <b className="mono" style={{ fontWeight: 600, fontSize: 12.5 }}>
          {loc.symbol || f.module}
        </b>
        <span className="spacer" />
        <span className="num" style={{ fontWeight: 600 }}>
          {fmt(f.debt, 1)}
        </span>
        <span className="num dim" style={{ fontSize: 12 }}>
          {pct(f.shareOfScopeDebt)}
        </span>
      </summary>
      <div className="body">
        <div className="mono dim" style={{ fontSize: 12 }}>
          {loc.file}
          {loc.startLine ? `:${loc.startLine}${loc.endLine ? "–" + loc.endLine : ""}` : ""}
        </div>
        <div>
          <b>Why</b> · {f.explanation}
        </div>
        {rem.strategy && (
          <div>
            <b>Fix</b> · {rem.strategy}
            {rem.target ? <span className="dim"> — {rem.target}</span> : null}
          </div>
        )}
        {rem.steps.length > 0 && (
          <ul>
            {rem.steps.map((s) => (
              <li key={s}>{s}</li>
            ))}
          </ul>
        )}
        {rem.acceptance.length > 0 && (
          <div>
            <b>Done when</b> · {rem.acceptance.join("; ")}
          </div>
        )}
        {rem.caution && <div className="dim">Caution: {rem.caution}</div>}
        {f.related && f.related.length > 0 && (
          <div className="mono dim" style={{ fontSize: 12 }}>
            Related: {f.related.map((x) => `${x.file}:${x.startLine}${x.symbol ? " " + x.symbol : ""}`).join(" · ")}
          </div>
        )}
        <div>
          <button className="btn small" onClick={prompt}>
            {copied ? "Copied" : "Copy as agent prompt"}
          </button>
        </div>
      </div>
    </details>
  );
}

function Defects({ sha, path }: { sha: string; path: string }) {
  const q = useDefects(sha, path);
  const r = q.data;
  const groups = useMemo(() => {
    const g = new Map<string, NonNullable<typeof r>["findings"]>();
    for (const f of r?.findings ?? []) {
      const k = f.severity + "|" + f.rule;
      g.set(k, [...(g.get(k) ?? []), f]);
    }
    return [...g.values()].sort((a, b) => (a![0].severity === "gate" ? 0 : 1) - (b![0].severity === "gate" ? 0 : 1) || b!.length - a!.length);
  }, [r]);
  return (
    <div className="body" style={{ display: "flex", flexDirection: "column", gap: 8 }}>
      <h3>Defect findings</h3>
      {q.isPending || q.run.isPending ? (
        <p className="dim" role="status">
          {q.run.isPending ? "Checking this commit…" : "Looking for an earlier check…"}
        </p>
      ) : q.isError || q.run.isError ? (
        <ErrorBanner title="The defect check failed" error={q.error ?? q.run.error} onRetry={() => q.run.mutate()} />
      ) : r === null ? (
        <>
          <p className="dim" style={{ margin: 0, fontSize: 13 }}>
            Not checked for {short(sha)} yet. The rules run once for the whole commit (about 30 s on a large repository the first time); every folder then reuses the result.
          </p>
          <div>
            <button className="btn" onClick={() => q.run.mutate()}>
              Check defects for this commit
            </button>
          </div>
        </>
      ) : r ? (
        <>
          <p style={{ margin: 0, fontSize: 13 }}>
            <b className={r.totals.gate ? "bad" : ""}>{fmt(r.totals.gate)} gate</b> · {fmt(r.totals.advisory)} advisory
            {r.suppressed?.length ? ` · ${fmt(r.suppressed.length)} suppressed` : ""} <span className="dim">(gate = provable from syntax; advisory = needs a human decision)</span>
          </p>
          {!groups.length && <p className="dim">No findings in this folder.</p>}
          {groups.map((xs) => (
            <details key={xs![0].severity + xs![0].rule} className="item">
              <summary>
                <span className={"sev " + xs![0].severity}>{xs![0].severity}</span>
                <b>{xs![0].rule}</b>
                <span className="dim">· {xs!.length}</span>
              </summary>
              <div className="body">
                <ul>
                  {xs!.slice(0, 40).map((f, i) => (
                    <li key={i}>
                      <span className="mono">
                        {f.file}:{f.line}
                      </span>
                      {f.symbol ? ` (${f.symbol})` : ""} — {f.message}
                    </li>
                  ))}
                </ul>
                {xs!.length > 40 && <p className="dim">{fmt(xs!.length - 40)} more</p>}
              </div>
            </details>
          ))}
        </>
      ) : null}
    </div>
  );
}

type SortKey = "debt" | "cc" | "nest" | "lines" | "name";

function ScopeTables({ sha, path }: { sha: string; path: string }) {
  const r = useExplain(sha, path);
  const [tab, setTab] = useState<"items" | "files">("items");
  const [q, setQ] = useState("");
  const [sort, setSort] = useState<SortKey>("debt");
  const factors = r.data?.factors ?? NO_FACTORS;
  const files = useMemo(() => byFile(factors), [factors]);
  const items = useMemo(() => {
    const needle = q.trim().toLowerCase();
    const xs = needle ? factors.filter((f) => (f.location.symbol ?? "").toLowerCase().includes(needle) || f.location.file.toLowerCase().includes(needle) || f.kind.includes(needle)) : factors;
    const key = (f: Factor) => (sort === "debt" ? f.debt : sort === "cc" ? (f.metrics?.cyclomatic ?? -1) : sort === "nest" ? (f.metrics?.nesting ?? -1) : sort === "lines" ? (f.metrics?.lines ?? -1) : 0);
    return sort === "name" ? [...xs].sort((a, b) => (a.location.symbol ?? "").localeCompare(b.location.symbol ?? "")) : [...xs].sort((a, b) => key(b) - key(a));
  }, [factors, q, sort]);
  const shownFiles = useMemo(() => {
    const needle = q.trim().toLowerCase();
    return needle ? files.filter((f) => f.file.toLowerCase().includes(needle)) : files;
  }, [files, q]);
  if (!r.data) return null;
  return (
    <section className="card" aria-label="Items and files in this folder">
      <div className="head" style={{ alignItems: "center" }}>
        <div role="tablist" style={{ display: "flex", gap: 4 }}>
          <button role="tab" aria-selected={tab === "items"} className={"btn small" + (tab === "items" ? " primary" : "")} onClick={() => setTab("items")}>
            Items ({fmt(factors.length)})
          </button>
          <button role="tab" aria-selected={tab === "files"} className={"btn small" + (tab === "files" ? " primary" : "")} onClick={() => setTab("files")}>
            Files ({fmt(files.length)})
          </button>
        </div>
        <div className="spacer" />
        <input type="search" aria-label="Filter" placeholder="Filter by name, file or kind" value={q} onChange={(e) => setQ(e.target.value)} style={{ width: 240 }} />
        <span className="dim num" style={{ fontSize: 12 }}>
          {fmt(tab === "items" ? items.length : shownFiles.length)} shown
        </span>
      </div>
      {tab === "items" ? <ItemTable items={items} sort={sort} setSort={setSort} /> : <FileTable files={shownFiles} />}
    </section>
  );
}

const ROW = 34;

function ItemTable({ items, sort, setSort }: { items: Factor[]; sort: SortKey; setSort: (s: SortKey) => void }) {
  const box = useRef<HTMLDivElement>(null);
  const v = useVirtualizer({ count: items.length, getScrollElement: () => box.current, estimateSize: () => ROW, overscan: 12 });
  const cols = "minmax(0, 2.2fr) 110px 56px 56px 64px 76px";
  const H = ({ k, label, right }: { k: SortKey; label: string; right?: boolean }) => (
    <button className="dim" style={{ background: "none", border: 0, padding: 0, cursor: "pointer", textAlign: right ? "right" : "left", fontWeight: sort === k ? 650 : 500 }} aria-sort={sort === k ? "descending" : "none"} onClick={() => setSort(k)}>
      {label}
      {sort === k ? " ▾" : ""}
    </button>
  );
  return (
    <div role="table" aria-rowcount={items.length} style={{ fontSize: 13 }}>
      <div role="row" style={{ display: "grid", gridTemplateColumns: cols, gap: 8, padding: "8px 18px", borderBottom: "1px solid var(--line)" }}>
        <H k="name" label="Item" />
        <span className="dim">Kind</span>
        <H k="cc" label="cc" right />
        <H k="nest" label="nest" right />
        <H k="lines" label="lines" right />
        <H k="debt" label="debt" right />
      </div>
      <div ref={box} style={{ height: Math.min(440, Math.max(ROW, items.length * ROW)), overflow: "auto" }}>
        <div style={{ height: v.getTotalSize(), position: "relative" }}>
          {v.getVirtualItems().map((r) => {
            const f = items[r.index];
            return (
              <div key={f.id} role="row" style={{ position: "absolute", top: 0, left: 0, right: 0, transform: `translateY(${r.start}px)`, height: ROW, display: "grid", gridTemplateColumns: cols, gap: 8, alignItems: "center", padding: "0 18px", borderBottom: "1px solid var(--line)" }}>
                <span style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }} title={`${f.location.file}:${f.location.startLine}`}>
                  <span className="mono">{f.location.symbol || f.module}</span>{" "}
                  <span className="mono dim" style={{ fontSize: 11.5 }}>
                    {f.location.file.split("/").pop()}:{f.location.startLine}
                  </span>
                </span>
                <span className="dim">{KIND_NAMES[f.kind] ?? f.kind}</span>
                <span className="num" style={{ textAlign: "right" }}>
                  {f.metrics?.cyclomatic ?? "–"}
                </span>
                <span className="num" style={{ textAlign: "right" }}>
                  {f.metrics?.nesting ?? "–"}
                </span>
                <span className="num" style={{ textAlign: "right" }}>
                  {f.metrics?.lines ?? "–"}
                </span>
                <span className="num" style={{ textAlign: "right", fontWeight: 600 }}>
                  {fmt(f.debt, 2)}
                </span>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

function FileTable({ files }: { files: ReturnType<typeof byFile> }) {
  return (
    <div className="scroll-x" style={{ maxHeight: 440, overflowY: "auto" }}>
      <table className="data">
        <thead>
          <tr>
            <th scope="col">File</th>
            <th scope="col">Largest item</th>
            <th scope="col" className="r">
              Items
            </th>
            <th scope="col" className="r">
              Hotspots
            </th>
            <th scope="col" className="r">
              Debt ▾
            </th>
          </tr>
        </thead>
        <tbody>
          {files.slice(0, 500).map((f) => (
            <tr key={f.file}>
              <td className="mono" style={{ fontSize: 12.5 }}>
                {f.file}
              </td>
              <td className="mono dim" style={{ fontSize: 12 }}>
                {f.worst}
              </td>
              <td className="r">{fmt(f.items)}</td>
              <td className="r">{fmt(f.hotspots)}</td>
              <td className="r" style={{ fontWeight: 600 }}>
                {fmt(f.debt, 2)}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {files.length > 500 && <p className="dim" style={{ padding: "0 18px" }}>{fmt(files.length - 500)} more files; filter to narrow the list.</p>}
    </div>
  );
}

