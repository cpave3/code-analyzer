// Overview: what to look at first for the selected commit.
import { useMemo } from "react";
import { Link, useNavigate } from "react-router";
import type { Commit, ExplainReport, GraphData, Metrics } from "../api";
import { EChart, type EChartsOption } from "../components/EChart";
import { EmptyState, ErrorBanner, Skeleton } from "../components/states";
import { useExplain, useGraph, useMeta, useMetrics, useSelection } from "../data";
import { compare } from "../lib/changes";
import { day, fmt, moduleLabel, pct, short, signed, stamp } from "../lib/format";
import { COMPONENTS, usable } from "../lib/history";
import { componentOf } from "../lib/report";
import { useTheme, type Tokens } from "../theme";

export function Overview() {
  const meta = useMeta().data!;
  const { byIndex, data: mr } = useMetrics();
  const { index, commit } = useSelection();
  const m = byIndex[index];
  const explain = useExplain(commit?.sha, "");
  const graph = useGraph(usable(m) ? commit?.sha : undefined);
  const prevSha = index > 0 && usable(byIndex[index - 1]) ? meta.commits[index - 1].sha : undefined;
  const prevGraph = useGraph(prevSha);
  if (!commit) return null;
  return (
    <main className="page">
      <CommitTitle commit={commit} scoring={meta.scoring} />
      {m?.error ? (
        <ErrorBanner title={`Commit ${short(commit.sha)} could not be measured`} error={new Error(m.error)} />
      ) : !usable(m) ? (
        <MeasuringCommit />
      ) : m.lines === 0 ? (
        <NothingMeasurable commit={commit} explain={explain.data} />
      ) : (
        <>
          <Totals m={m} prev={byIndex[index - 1]} explain={explain.data} />
          <div className="row">
            <div className="main-col">
              <FixFirst explain={explain.data} error={explain.error} loading={explain.isPending} retry={() => explain.refetch()} />
              <TopModules graph={graph.data} loading={graph.isPending} />
            </div>
            <div className="side-col">
              <ByComponent m={m} />
              {mr && <Trend from={mr.window.from} to={mr.window.to} index={index} byIndex={byIndex} commits={meta.commits} />}
              {graph.data && <Changed now={graph.data} prev={prevSha ? (prevGraph.data ?? null) : null} hasPrev={!!prevSha} />}
              {explain.data && <NotMeasured explain={explain.data} />}
            </div>
          </div>
        </>
      )}
    </main>
  );
}

export function CommitTitle({ commit, scoring }: { commit: Commit; scoring: string }) {
  return (
    <div className="titlebar">
      <h1 title={commit.subject} style={{ maxWidth: "100%", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
        {commit.subject}
      </h1>
      <span className="mono dim">
        {short(commit.sha)} · {commit.author} · {stamp(commit.time)}
        {commit.merge ? " · merge" : ""}
      </span>
      <span className="pill">{scoring}</span>
    </div>
  );
}

function MeasuringCommit() {
  return (
    <>
      <div className="kpis">
        {[0, 1, 2, 3, 4].map((i) => (
          <Skeleton key={i} height={88} />
        ))}
      </div>
      <p className="dim" role="status">
        Measuring this commit; the numbers appear here as soon as it is done.
      </p>
      <Skeleton height={360} />
    </>
  );
}

function NothingMeasurable({ commit, explain }: { commit: Commit; explain?: ExplainReport }) {
  const langs = explain?.coverage.languages.filter((l) => !l.supported) ?? [];
  return (
    <EmptyState title={`No Go, TypeScript or C# source at ${short(commit.sha)}`}>
      {langs.length > 0 && (
        <>
          <div>CodeEval found {fmt(langs.reduce((a, l) => a + l.files, 0))} source files it can't measure:</div>
          <div className="chips">
            {langs.map((l) => (
              <span className="chip" key={l.language}>
                {l.language} {fmt(l.files)}
              </span>
            ))}
          </div>
        </>
      )}
      <div>Choose another commit from the Commit menu, or open Changes to measure a different part of the history.</div>
    </EmptyState>
  );
}

function Totals({ m, prev, explain }: { m: Metrics; prev?: Metrics; explain?: ExplainReport }) {
  const d = usable(prev) ? m.debt.total - prev.debt.total : null;
  const cov = explain?.coverage;
  const items = [
    { label: "Debt", value: fmt(m.debt.total, 2), sub: d == null ? "first measured commit" : `${signed(d, 2)} since the previous commit` },
    { label: "Debt per kLOC", value: fmt(m.debtPerKLOC, 2), sub: `${fmt(m.lines)} source lines` },
    { label: "Hotspots", value: fmt(m.hotspots), sub: `of ${fmt(m.callables)} callables` },
    { label: "Score (experimental)", value: fmt(m.score, 2), sub: "0–100 · for comparison, not a target" },
    { label: "Measured", value: cov ? pct(cov.supportedByteShare) : "…", sub: cov ? `of source bytes · ${fmt(cov.supportedSourceFiles)} of ${fmt(cov.sourceFiles)} files` : "reading coverage" },
  ];
  return (
    <section aria-label="Totals" className="kpis">
      {items.map((k) => (
        <div className="kpi" key={k.label}>
          <div className="label">{k.label}</div>
          <div className="value">{k.value}</div>
          <div className="sub">{k.sub}</div>
        </div>
      ))}
    </section>
  );
}

function FixFirst({ explain, error, loading, retry }: { explain?: ExplainReport; error: unknown; loading: boolean; retry: () => void }) {
  const { tokens } = useTheme();
  const top = (explain?.factors ?? []).slice(0, 8);
  const max = top[0]?.debt ?? 1;
  const conc = explain?.concentration;
  return (
    <section className="card" aria-labelledby="fixfirst">
      <div className="head">
        <h2 id="fixfirst">Fix first</h2>
        <span className="note">{conc && conc.factors ? `${fmt(conc.factorsFor50Percent)} of ${fmt(conc.factors)} debt items carry half the debt` : "Largest debt items"}</span>
        <div className="spacer" />
        <Link to={{ pathname: "/modules", search: location.search }} style={{ fontSize: 13 }}>
          All items →
        </Link>
      </div>
      {loading ? (
        <div className="body">
          <Skeleton height={300} />
        </div>
      ) : error ? (
        <div className="body">
          <ErrorBanner title="The debt items could not be loaded" error={error} onRetry={retry} />
        </div>
      ) : !top.length ? (
        <div className="body dim">No debt items at this commit.</div>
      ) : (
        <div className="scroll-x">
          <table className="data">
            <thead>
              <tr>
                <th scope="col">Item</th>
                <th scope="col">Module</th>
                <th scope="col" className="r">
                  cc
                </th>
                <th scope="col" className="r">
                  nest
                </th>
                <th scope="col" style={{ width: 210 }}>
                  Debt
                </th>
              </tr>
            </thead>
            <tbody>
              {top.map((f) => (
                <tr key={f.id}>
                  <td>
                    <div className="mono" style={{ fontWeight: 500 }}>
                      {f.location.symbol || f.module}
                    </div>
                    <div className="mono dim" style={{ fontSize: 11.5 }}>
                      {f.location.file}
                      {f.location.startLine ? ":" + f.location.startLine : ""}
                    </div>
                  </td>
                  <td className="mono dim" style={{ fontSize: 12 }}>
                    <Link to={{ pathname: "/modules/" + f.module, search: location.search }}>{f.module}</Link>
                  </td>
                  <td className="r">{f.metrics?.cyclomatic ?? "–"}</td>
                  <td className="r">{f.metrics?.nesting ?? "–"}</td>
                  <td>
                    <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
                      <div className="meter">
                        <div style={{ width: `${(100 * f.debt) / max}%`, background: tokens.components[componentOf(f.kind)] }} />
                      </div>
                      <span className="num" style={{ width: 56, textAlign: "right", fontWeight: 600 }}>
                        {fmt(f.debt, 2)}
                      </span>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

function ByComponent({ m }: { m: Metrics }) {
  const { tokens } = useTheme();
  const total = m.debt.total || 1;
  const raw: Record<string, string> = {
    complexity: `${fmt(m.hotspots)} hotspots`,
    nesting: "deep nesting",
    duplication: `${fmt(m.clonedCallables)} cloned callables`,
    architecture: `${fmt(m.ruleViolations)} rule edges, ${fmt(m.cycleUnits)} in cycles`,
  };
  const rows = COMPONENTS.map((k) => ({ k, v: m.debt[k] })).sort((a, b) => b.v - a.v);
  return (
    <section className="card pad" aria-labelledby="bycomp">
      <h2 id="bycomp">Debt by component</h2>
      <div className="stack" aria-hidden="true">
        {rows
          .filter((r) => r.v > 0)
          .map((r) => (
            <div key={r.k} style={{ width: `${(100 * r.v) / total}%`, background: tokens.components[r.k] }} />
          ))}
      </div>
      <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
        {rows.map((r) => (
          <div key={r.k} style={{ display: "flex", alignItems: "baseline", gap: 10, fontSize: 13 }}>
            <span className="sw" style={{ background: tokens.components[r.k], alignSelf: "center" }} />
            <span style={{ flex: 1, textTransform: "capitalize" }}>
              {r.k} <span className="dim" style={{ textTransform: "none" }}>· {raw[r.k]}</span>
            </span>
            <span className="num" style={{ fontWeight: 600 }}>
              {fmt(r.v, 2)}
            </span>
            <span className="num dim" style={{ width: 48, textAlign: "right" }}>
              {pct(r.v / total)}
            </span>
          </div>
        ))}
      </div>
    </section>
  );
}

function trendOption(t: Tokens, pts: [number, number][], index: number): EChartsOption {
  return {
    animation: false,
    grid: { left: 4, right: 8, top: 10, bottom: 4, containLabel: true },
    xAxis: { type: "value", min: pts[0]?.[0], max: pts[pts.length - 1]?.[0], show: false },
    yAxis: { type: "value", scale: true, splitLine: { lineStyle: { color: t.line } }, axisLabel: { color: t.dim, fontSize: 10 } },
    tooltip: { trigger: "axis", valueFormatter: (v: unknown) => fmt(v as number, 2) },
    series: [
      {
        type: "line",
        data: pts,
        showSymbol: false,
        lineStyle: { color: t.accent, width: 2 },
        itemStyle: { color: t.accent },
        markLine: { symbol: "none", silent: true, label: { show: false }, lineStyle: { color: t.dim, type: "dashed" }, data: [{ xAxis: index }] },
      },
    ],
  };
}

function Trend({ from, to, index, byIndex, commits }: { from: number; to: number; index: number; byIndex: (Metrics | undefined)[]; commits: Commit[] }) {
  const { tokens } = useTheme();
  const pts = useMemo(() => {
    const out: [number, number][] = [];
    for (let i = from; i <= to; i++) {
      const m = byIndex[i];
      if (usable(m)) out.push([i, m.debt.total]);
    }
    return out;
  }, [from, to, byIndex]);
  const option = useMemo(() => trendOption(tokens, pts, index), [tokens, pts, index]);
  const first = pts[0],
    last = pts[pts.length - 1];
  const net = first && last ? last[1] - first[1] : null;
  return (
    <section className="card pad" aria-labelledby="trend">
      <div style={{ display: "flex", alignItems: "baseline", gap: 8 }}>
        <h2 id="trend" style={{ margin: 0, fontSize: 15, fontWeight: 650 }}>
          Debt over {to - from + 1} commits
        </h2>
        <div className="spacer" />
        {net != null && <span className={"num " + (net > 0 ? "bad" : net < 0 ? "good" : "")} style={{ fontSize: 13, fontWeight: 600 }}>{signed(net, 2)}</span>}
      </div>
      {pts.length > 1 ? <EChart option={option} height={130} label={`Debt across the measured window, ${first ? fmt(first[1], 0) : ""} to ${last ? fmt(last[1], 0) : ""}`} /> : <Skeleton height={130} />}
      <div className="num dim" style={{ display: "flex", justifyContent: "space-between", fontSize: 12 }}>
        <span>
          {short(commits[from].sha)} · {day(commits[from].time)}
        </span>
        <Link to={{ pathname: "/changes", search: location.search }}>Open Changes →</Link>
      </div>
    </section>
  );
}

function Changed({ now, prev, hasPrev }: { now: GraphData; prev: GraphData | null; hasPrev: boolean }) {
  const c = useMemo(() => (prev ? compare(now, prev) : null), [now, prev]);
  return (
    <section className="card pad" aria-labelledby="changed">
      <h2 id="changed">Changed in this commit</h2>
      {!hasPrev ? (
        <div className="dim" style={{ fontSize: 13 }}>
          The previous commit is outside the measured window.
        </div>
      ) : !c ? (
        <Skeleton height={60} />
      ) : !c.biggest.length && !c.gone.length ? (
        <div className="dim" style={{ fontSize: 13 }}>
          No module or function gained or lost debt.
        </div>
      ) : (
        <div style={{ display: "flex", flexDirection: "column", gap: 6, fontSize: 13 }}>
          {c.biggest.map((x) => (
            <div key={x.id} style={{ display: "flex", gap: 8, alignItems: "baseline" }}>
              <span className="dim" style={{ width: 54, fontSize: 12 }}>
                {x.kind === "module" ? "module" : "function"}
              </span>
              <span className="mono" style={{ flex: 1, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                {x.label}
              </span>
              <span className={"num " + (x.delta > 0 ? "bad" : "good")} style={{ fontWeight: 600 }}>
                {x.state === "new" ? "new " + fmt(x.delta, 1) : signed(x.delta, 1)}
              </span>
            </div>
          ))}
          {c.gone.map((x) => (
            <div key={x.id} className="dim" style={{ display: "flex", gap: 8 }}>
              <span style={{ width: 54, fontSize: 12 }}>{x.kind === "module" ? "module" : "function"}</span>
              <span className="mono" style={{ flex: 1 }}>
                {x.label}
              </span>
              <span>gone (was {fmt(x.debt, 1)})</span>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}

function NotMeasured({ explain }: { explain: ExplainReport }) {
  const cov = explain.coverage;
  const langs = cov.languages.filter((l) => !l.supported).sort((a, b) => b.files - a.files);
  if (!langs.length && !cov.parseErrorFiles) return null;
  return (
    <section className="card pad" aria-labelledby="cov">
      <h2 id="cov">Not measured</h2>
      <div className="dim" style={{ fontSize: 13 }}>
        {fmt(cov.sourceFiles - cov.supportedSourceFiles)} of {fmt(cov.sourceFiles)} source files
        {cov.parseErrorFiles ? ` · ${fmt(cov.parseErrorFiles)} with parse errors` : ""}
      </div>
      <div className="chips">
        {langs.map((l) => (
          <span className="chip" key={l.language}>
            {l.language} {fmt(l.files)}
          </span>
        ))}
      </div>
    </section>
  );
}

function TopModules({ graph, loading }: { graph?: GraphData; loading: boolean }) {
  const { tokens } = useTheme();
  const nav = useNavigate();
  const mods = useMemo(() => [...(graph?.modules ?? [])].sort((a, b) => b.debt.total - a.debt.total).slice(0, 10), [graph]);
  const max = mods[0]?.debt.total || 1;
  return (
    <section className="card" aria-labelledby="mods">
      <div className="head">
        <h2 id="mods">Modules with the most debt</h2>
        <span className="note">{graph ? `${Math.min(10, graph.modules.length)} of ${graph.modules.length} · bar split by component` : ""}</span>
        <div className="spacer" />
        <Link to={{ pathname: "/modules", search: location.search }} style={{ fontSize: 13 }}>
          Open the module map →
        </Link>
      </div>
      {loading ? (
        <div className="body">
          <Skeleton height={260} />
        </div>
      ) : (
        <div className="scroll-x">
          <table className="data">
            <thead>
              <tr>
                <th scope="col">Module</th>
                <th scope="col" style={{ width: "45%" }}>
                  <span className="sr-only">Debt by component</span>
                </th>
                <th scope="col" className="r">
                  Debt
                </th>
                <th scope="col" className="r">
                  per kLOC
                </th>
              </tr>
            </thead>
            <tbody>
              {mods.map((m) => {
                const path = moduleLabel(m.id);
                return (
                  <tr key={m.id} className="click" onClick={() => nav({ pathname: "/modules/" + path, search: location.search })}>
                    <td className="mono" style={{ fontSize: 12.5 }}>
                      <Link to={{ pathname: "/modules/" + path, search: location.search }} onClick={(e) => e.stopPropagation()} style={{ color: "inherit" }}>
                        {path}
                      </Link>
                    </td>
                    <td>
                      <div style={{ display: "flex", height: 12, width: `${(100 * m.debt.total) / max}%`, gap: 1 }}>
                        {COMPONENTS.map((k) => (m.debt[k] > 0 ? <div key={k} title={`${k} ${fmt(m.debt[k], 1)}`} style={{ flex: m.debt[k], background: tokens.components[k] }} /> : null))}
                      </div>
                    </td>
                    <td className="r" style={{ fontWeight: 600 }}>
                      {fmt(m.debt.total, 2)}
                    </td>
                    <td className="r" style={{ color: m.density >= 105 ? tokens.bad : tokens.dim }}>
                      {fmt(m.density, 1)}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

