// Changes: the history and the measured window (drag across the strip, or
// pick the last N commits or a date range), stepping through commits, what
// each step changed, metrics over the window and step changes by
// provenance group.
import { useCallback, useMemo, useState } from "react";
import type { Commit, Meta, Metrics } from "../core/api";
import { EChart, type EChartsOption } from "../components/EChart";
import { Next, Pause, Play, Prev } from "../components/icons";
import { usePlayback } from "../components/Playback";
import { ErrorBanner, Skeleton } from "../components/states";
import { useDefects, useMeta, useMetrics, useSelection, useSetWindow } from "../core/data";
import { fieldValue, keyColors, stepKey } from "../lib/provenance";
import { fmt, short, signed, stamp } from "../lib/format";
import { change, COMPONENTS, dateWindow, debtChange, groupRows, lastN, recentre, usable } from "../lib/history";
import { useTheme, type Tokens } from "../core/theme";

export default function Changes() {
  const meta = useMeta().data!;
  const metrics = useMetrics();
  const { index, select } = useSelection();
  const [fields, setFields] = useState<string[]>(["model", "harness"].filter((f) => meta.groupFields.includes(f)));
  const keys = useMemo(() => meta.commits.map((c) => stepKey(c, fields)), [meta, fields]);
  const colours = useMemo(() => keyColors(keys), [keys]);
  const win = metrics.data?.window;
  if (!win)
    return (
      <main className="page">
        {metrics.isError ? <ErrorBanner title="Can't reach codeeval serve" error={metrics.error} onRetry={() => metrics.refetch()} /> : <Skeleton height={500} />}
      </main>
    );
  return (
    <main className="page">
      <History meta={meta} keys={keys} colours={colours} from={win.from} to={win.to} index={index} select={select} measured={metrics.data!.measuredOverall} />
      <div className="row">
        <div className="main-col">
          <CommitPanel meta={meta} index={index} byIndex={metrics.byIndex} select={select} from={win.from} to={win.to} fields={fields} />
          <StepBars from={win.from} to={win.to} index={index} byIndex={metrics.byIndex} commits={meta.commits} select={select} />
          <MetricCharts meta={meta} from={win.from} to={win.to} index={index} byIndex={metrics.byIndex} keys={keys} colours={colours} select={select} />
        </div>
        <div className="side-col">
          <GroupPicker meta={meta} fields={fields} setFields={setFields} keys={keys} colours={colours} />
        </div>
      </div>
      <GroupTable meta={meta} from={win.from} to={win.to} byIndex={metrics.byIndex} fields={fields} colours={colours} />
    </main>
  );
}

function History({ meta, keys, colours, from, to, index, select, measured }: { meta: Meta; keys: string[]; colours: Map<string, string>; from: number; to: number; index: number; select: (i: number, force?: boolean) => void; measured: number }) {
  const { tokens } = useTheme();
  const setWindow = useSetWindow();
  const n = meta.commits.length;
  const [last, setLast] = useState(Math.min(50, n));
  const [a, setA] = useState(meta.commits[0].time.slice(0, 10));
  const [b, setB] = useState(meta.commits[n - 1].time.slice(0, 10));
  const [msg, setMsg] = useState("");
  const move = useCallback(
    (lo: number, hi: number) => {
      setMsg("");
      setWindow.mutate({ from: lo, to: hi });
      if (index < lo || index > hi) select(hi, true);
    },
    [setWindow, index, select],
  );
  const option = useMemo<EChartsOption>(
    () => ({
      animation: false,
      grid: { left: 6, right: 6, top: 6, bottom: 22 },
      xAxis: {
        type: "category",
        data: meta.commits.map((c) => c.time.slice(0, 10)),
        axisTick: { show: false },
        axisLine: { show: false },
        axisLabel: { color: tokens.dim, fontSize: 10, interval: Math.max(0, Math.floor(n / 6) - 1), hideOverlap: true },
      },
      yAxis: { type: "value", show: false, max: 1 },
      tooltip: {
        trigger: "axis",
        confine: true,
        formatter: (p: { dataIndex: number }[]) => {
          const c = meta.commits[p[0].dataIndex];
          return `<b>${esc(c.subject)}</b><br><span style="font-family:monospace">${short(c.sha)}</span> · ${c.time.slice(0, 10)}<br>${esc(keys[c.index])}`;
        },
      },
      brush: { toolbox: [], xAxisIndex: 0, brushType: "lineX", brushMode: "single", brushStyle: { color: tokens.accent + "2e", borderColor: tokens.accent, borderWidth: 1.5 }, transformable: false, throttleType: "debounce", throttleDelay: 200, outOfBrush: { colorAlpha: 0.35 } },
      series: [
        {
          type: "bar",
          barCategoryGap: n > 200 ? "0%" : "12%",
          data: meta.commits.map((c, i) => ({ value: 1, itemStyle: { color: i === index ? tokens.fg : colours.get(keys[c.index]) } })),
          large: n > 2000,
        },
      ],
    }),
    [meta, keys, colours, tokens, index, n],
  );
  const after = useCallback(
    (c: { dispatchAction: (a: never) => void }) => {
      c.dispatchAction({ type: "takeGlobalCursor", key: "brush", brushOption: { brushType: "lineX", brushMode: "single" } } as never);
      c.dispatchAction({ type: "brush", areas: [{ brushType: "lineX", xAxisIndex: 0, coordRange: [from, to] }] } as never);
    },
    [from, to],
  );
  const onEvents = useMemo(
    () => ({
      brushEnd: (p: unknown) => {
        const area = (p as { areas?: { coordRange?: number[] }[] }).areas?.[0];
        const r = area?.coordRange;
        if (!r) return;
        const lo = Math.round(Math.min(r[0], r[1])),
          hi = Math.round(Math.max(r[0], r[1]));
        if (lo === from && hi === to) return;
        if (lo === hi) {
          if (lo >= from && lo <= to) select(lo);
          else {
            const [x, y] = recentre(lo, from, to, n);
            setWindow.mutate({ from: x, to: y });
            select(lo, true);
          }
          return;
        }
        move(lo, hi);
      },
    }),
    [from, to, n, move, select, setWindow],
  );
  return (
    <section className="card" aria-labelledby="hist">
      <div className="head">
        <h2 id="hist">History</h2>
        <span className="note">
          {fmt(n)} first-parent commits · measuring {fmt(from + 1)}–{fmt(to + 1)} ({meta.commits[from].time.slice(0, 10)} → {meta.commits[to].time.slice(0, 10)}) · drag across the strip to choose
        </span>
      </div>
      <div className="body" style={{ display: "flex", flexDirection: "column", gap: 12 }}>
        <EChart option={option} after={after} onEvents={onEvents} height={70} label="Every commit in the history, coloured by provenance; the highlighted span is the measured window" />
        <form
          style={{ display: "flex", flexWrap: "wrap", gap: "10px 18px", alignItems: "center", fontSize: 13 }}
          onSubmit={(e) => {
            e.preventDefault();
          }}
        >
          <span style={{ display: "inline-flex", gap: 6, alignItems: "center" }}>
            <label htmlFor="lastN">Last</label>
            <input id="lastN" type="number" min={1} max={n} value={last} onChange={(e) => setLast(+e.target.value)} style={{ width: 80 }} />
            commits
            <button type="button" className="btn small" onClick={() => move(...lastN(last || 1, n))}>
              Measure
            </button>
          </span>
          <span style={{ display: "inline-flex", gap: 6, alignItems: "center", flexWrap: "wrap" }}>
            <label htmlFor="fromDate">or from</label>
            <input id="fromDate" type="date" value={a} onChange={(e) => setA(e.target.value)} />
            <label htmlFor="toDate">to</label>
            <input id="toDate" type="date" value={b} onChange={(e) => setB(e.target.value)} />
            <button
              type="button"
              className="btn small"
              onClick={() => {
                const w = dateWindow(meta.commits, a, b);
                if (w) move(w[0], w[1]);
                else setMsg("No commits in that date range.");
              }}
            >
              Measure
            </button>
          </span>
          <span className="dim" role="status">
            {msg || (setWindow.isError ? "The window could not be changed: " + String(setWindow.error) : `${fmt(measured)} commits measured so far`)}
          </span>
        </form>
      </div>
    </section>
  );
}

const esc = (s: string) => s.replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" })[c]!);

function CommitPanel({ meta, index, byIndex, select, from, to, fields }: { meta: Meta; index: number; byIndex: (Metrics | undefined)[]; select: (i: number) => void; from: number; to: number; fields: string[] }) {
  const { playing, toggle } = usePlayback();
  const c = meta.commits[index];
  const shown = new Set(["model", "harness", ...meta.groupFields]);
  const recs = c.provenance.filter((r) => !(r.merge && r.source === "none"));
  const chips: [string, string][] = [];
  for (const f of shown) {
    const vs = [...new Set(recs.map((r) => fieldValue(r, f)).filter(Boolean))];
    if (vs.length) chips.push([f, vs.join(" | ")]);
  }
  const sources = [...new Set(c.provenance.map((r) => r.source))].join(", ");
  return (
    <section className="card" aria-labelledby="commit">
      <div className="head" style={{ alignItems: "center" }}>
        <h2 id="commit">Commit</h2>
        <div style={{ display: "flex", gap: 6, alignItems: "center", flex: 1, minWidth: 260 }}>
          <button className="btn icon" aria-label="Previous commit (←)" title="Previous commit (←)" onClick={() => select(index - 1)} disabled={index <= from}>
            <Prev />
          </button>
          <button className="btn" onClick={toggle} title="Play or pause (space)">
            {playing ? <Pause /> : <Play />} {playing ? "Pause" : "Play"}
          </button>
          <button className="btn icon" aria-label="Next commit (→)" title="Next commit (→)" onClick={() => select(index + 1)} disabled={index >= to}>
            <Next />
          </button>
          <input type="range" aria-label="Commit in the window" min={from} max={to} value={index} onChange={(e) => select(+e.target.value)} style={{ flex: 1, accentColor: "var(--accent)" }} />
          <span className="mono dim num">
            {index + 1}/{meta.commits.length}
          </span>
        </div>
      </div>
      <div className="body" style={{ display: "flex", flexDirection: "column", gap: 10 }}>
        <div style={{ fontWeight: 600 }}>{c.subject}</div>
        <div className="dim" style={{ fontSize: 13 }}>
          <span className="mono">{c.sha.slice(0, 12)}</span> · {c.author} · {stamp(c.time)}
          {c.merge ? " · merge" : ""} · {c.provenance.length} commit{c.provenance.length === 1 ? "" : "s"} in this step
        </div>
        <div className="chips">
          {chips.length ? (
            chips.map(([f, v]) => (
              <span className="chip" key={f}>
                {f}: <b>{v}</b>
              </span>
            ))
          ) : (
            <span className="chip dim">no provenance</span>
          )}
          <span className="chip dim">source: {sources}</span>
          <span className="chip dim">group: {stepKey(c, fields)}</span>
        </div>
        <Deltas meta={meta} index={index} byIndex={byIndex} />
      </div>
    </section>
  );
}

interface DeltaSpec {
  label: string;
  get: (x: Metrics) => number | null | undefined;
  d: number;
  goodDown: boolean | null; // null: neither direction is better
}

const DELTAS: DeltaSpec[] = [
  { label: "Debt", get: (x) => x.debt.total, d: 1, goodDown: true },
  { label: "Score", get: (x) => x.score, d: 1, goodDown: false },
  { label: "Hotspots", get: (x) => x.hotspots, d: 0, goodDown: true },
  { label: "Decision points", get: (x) => x.decisionPoints, d: 0, goodDown: true },
  { label: "Source lines", get: (x) => x.lines, d: 0, goodDown: null },
];

const nearPairs = (x: Metrics) => (x.detail ? (x.nearMisses ?? 0) + (x.fragments ?? 0) : null);

// tone is the class for a change: bad when it moves the wrong way.
function tone(v: number | null, goodDown: boolean | null): string {
  if (v == null || Math.abs(v) < 1e-9 || goodDown == null) return "";
  return v > 0 === goodDown ? "bad" : "good";
}

function Kpi({ label, value, sub, cls = "" }: { label: string; value: string; sub: string; cls?: string }) {
  return (
    <div className="kpi">
      <div className="label">{label}</div>
      <div className={"value " + cls}>{value}</div>
      <div className="sub">{sub}</div>
    </div>
  );
}

function DeltaKpi({ spec, byIndex, index }: { spec: DeltaSpec; byIndex: (Metrics | undefined)[]; index: number }) {
  const v = change(byIndex, index, spec.get);
  const now = spec.get(byIndex[index]!);
  return <Kpi label={spec.label} value={signed(v, spec.d)} sub={now != null ? "now " + fmt(now, spec.d) : ""} cls={tone(v, spec.goodDown)} />;
}

function Deltas({ meta, index, byIndex }: { meta: Meta; index: number; byIndex: (Metrics | undefined)[] }) {
  const m = byIndex[index];
  if (m?.error) return <ErrorBanner title="This commit could not be measured" error={new Error(m.error)} />;
  if (!usable(m))
    return (
      <p className="dim" role="status">
        Measuring this commit…
      </p>
    );
  const changed = index === 0 ? "–" : m.stepDone ? fmt(m.changedLines) : "…";
  return (
    <div className="kpis compact" style={{ gridTemplateColumns: "repeat(auto-fit, minmax(140px, 1fr))" }}>
      {DELTAS.map((spec) => (
        <DeltaKpi key={spec.label} spec={spec} byIndex={byIndex} index={index} />
      ))}
      <Kpi label="Changed lines" value={changed} sub="measured source" />
      {meta.defects ? <DefectsIntroduced m={m} /> : <DefectCheck sha={meta.commits[index].sha} />}
      {meta.nearDuplication && (m.detail ? <DeltaKpi spec={{ label: "Near-dup pairs", get: nearPairs, d: 0, goodDown: true }} byIndex={byIndex} index={index} /> : <Kpi label="Near-dup pairs" value="…" sub="checking…" />)}
    </div>
  );
}

// DefectsIntroduced shows what the server's per-commit defect check found
// (with --defects).
function DefectsIntroduced({ m }: { m: Metrics }) {
  if (m.introducedGate == null) return <Kpi label="Defects introduced" value={m.detailStep ? "–" : "…"} sub={m.detailStep ? "" : "checking…"} />;
  const sub = `${m.introducedAdvisory ?? 0} advisory · ${m.resolved ?? 0} resolved · now ${m.gate ?? 0} gate`;
  return <Kpi label="Defects introduced" value={`${m.introducedGate} gate`} sub={sub} cls={m.introducedGate > 0 ? "bad" : ""} />;
}

// DefectCheck offers the whole-commit defect check when the server does not
// check every commit; it shows the totals once checked.
function DefectCheck({ sha }: { sha: string }) {
  const q = useDefects(sha, "");
  return (
    <div className="kpi">
      <div className="label">Defect findings</div>
      {q.data ? (
        <>
          <div className={"value " + (q.data.totals.gate ? "bad" : "")}>{q.data.totals.gate} gate</div>
          <div className="sub">{q.data.totals.advisory} advisory · whole commit</div>
        </>
      ) : q.run.isPending ? (
        <div className="sub" role="status">
          Checking…
        </div>
      ) : q.run.isError ? (
        <div className="sub bad">Check failed</div>
      ) : (
        <div>
          <button className="btn small" onClick={() => q.run.mutate()} disabled={q.isPending}>
            Check this commit
          </button>
        </div>
      )}
    </div>
  );
}

function StepBars({ from, to, index, byIndex, commits, select }: { from: number; to: number; index: number; byIndex: (Metrics | undefined)[]; commits: Commit[]; select: (i: number) => void }) {
  const { tokens } = useTheme();
  const { option, net } = useMemo(() => {
    const idx: number[] = [];
    for (let i = Math.max(1, from); i <= to; i++) idx.push(i);
    const ch = idx.map((i) => debtChange(byIndex, i));
    const net = ch.reduce((a, x) => a + (x?.total ?? 0), 0);
    const option: EChartsOption = {
      animation: false,
      grid: { left: 8, right: 8, top: 10, bottom: 24, containLabel: true },
      tooltip: {
        trigger: "axis",
        confine: true,
        axisPointer: { type: "shadow" },
        formatter: (p: { dataIndex: number }[]) => {
          const i = idx[p[0].dataIndex];
          const x = ch[p[0].dataIndex];
          const c = commits[i];
          return `<b>${esc(c.subject)}</b><br><span style="font-family:monospace">${short(c.sha)}</span><br>` + (x ? COMPONENTS.map((k) => `${k}: ${signed(x[k], 2)}`).join("<br>") + `<br><b>total ${signed(x.total, 2)}</b>` : "not measured yet");
        },
      },
      xAxis: { type: "category", data: idx.map((i) => short(commits[i].sha)), axisLabel: { show: false }, axisTick: { show: false }, axisLine: { lineStyle: { color: tokens.line } } },
      yAxis: { type: "value", splitLine: { lineStyle: { color: tokens.line } }, axisLabel: { color: tokens.dim, fontSize: 10 } },
      series: COMPONENTS.map((k) => ({
        name: k,
        type: "bar",
        stack: "d",
        data: ch.map((x, j) => ({ value: x ? Math.round(x[k] * 100) / 100 : 0, itemStyle: { color: tokens.components[k], opacity: idx[j] === index ? 1 : 0.85, borderColor: idx[j] === index ? tokens.fg : undefined, borderWidth: idx[j] === index ? 1 : 0 } })),
        barMaxWidth: 24,
      })),
    };
    return { option, net };
  }, [from, to, byIndex, commits, tokens, index]);
  const onEvents = useMemo(() => ({ click: (p: unknown) => select(Math.max(1, from) + (p as { dataIndex: number }).dataIndex) }), [from, select]);
  return (
    <section className="card" aria-labelledby="steps">
      <div className="head">
        <h2 id="steps">Debt change per commit</h2>
        <span className="note">Steps sum to the net: {signed(net, 2)}</span>
        <div className="spacer" />
        <Legend tokens={tokens} />
      </div>
      <div className="body">
        <EChart option={option} onEvents={onEvents} height={220} label={`Debt change per commit across the window; net ${signed(net, 2)}`} />
      </div>
    </section>
  );
}

function Legend({ tokens }: { tokens: Tokens }) {
  return (
    <span className="legend">
      {COMPONENTS.map((k) => (
        <span key={k}>
          <span className="sw" style={{ background: tokens.components[k] }} />
          {k}
        </span>
      ))}
    </span>
  );
}

interface ChartSpec {
  title: string;
  stack?: boolean;
  needs?: "defects" | "nearDuplication";
  d?: number;
  series: { name: string; color: (t: Tokens) => string; get: (m: Metrics) => number | null | undefined }[];
}

const CHARTS: ChartSpec[] = [
  { title: "Debt by component", stack: true, d: 0, series: COMPONENTS.map((k) => ({ name: k, color: (t: Tokens) => t.components[k], get: (m: Metrics) => m.debt[k] })) },
  { title: "Score (0–100, comparison only)", d: 1, series: [{ name: "score", color: (t) => t.good, get: (m) => m.score }] },
  { title: "Debt per kLOC", d: 1, series: [{ name: "debt/kLOC", color: (t) => t.components.nesting, get: (m) => m.debtPerKLOC }] },
  {
    title: "Hotspots and source size",
    series: [
      { name: "hotspots", color: (t) => t.components.architecture, get: (m) => m.hotspots },
      { name: "kLOC × 10", color: (t) => t.accent, get: (m) => m.lines / 100 },
    ],
  },
  {
    title: "Defect findings (rules/v1)",
    needs: "defects",
    series: [
      { name: "gate", color: (t) => t.bad, get: (m) => m.gate },
      { name: "advisory", color: (t) => t.components.nesting, get: (m) => m.advisory },
    ],
  },
  {
    title: "Near duplication (clones/v1, unscored)",
    needs: "nearDuplication",
    series: [
      { name: "near-miss pairs", color: (t) => t.components.duplication, get: (m) => m.nearMisses },
      { name: "fragment pairs", color: (t) => t.accent, get: (m) => m.fragments },
    ],
  },
];

function MetricCharts(props: { meta: Meta; from: number; to: number; index: number; byIndex: (Metrics | undefined)[]; keys: string[]; colours: Map<string, string>; select: (i: number) => void }) {
  const charts = CHARTS.filter((c) => !c.needs || props.meta[c.needs]);
  return (
    <section className="card" aria-labelledby="mbc">
      <div className="head">
        <h2 id="mbc">Metrics by commit</h2>
        <span className="note">Points are coloured by provenance group; click one to select that commit</span>
      </div>
      <div className="body" style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(320px, 1fr))", gap: 18 }}>
        {charts.map((c) => (
          <MetricChart key={c.title} spec={c} {...props} />
        ))}
      </div>
    </section>
  );
}

function MetricChart({ spec, from, to, index, byIndex, keys, colours, select, meta }: { spec: ChartSpec; meta: Meta; from: number; to: number; index: number; byIndex: (Metrics | undefined)[]; keys: string[]; colours: Map<string, string>; select: (i: number) => void }) {
  const { tokens } = useTheme();
  const option = useMemo<EChartsOption>(() => {
    const idx: number[] = [];
    for (let i = from; i <= to; i++) if (usable(byIndex[i])) idx.push(i);
    const lineSeries = spec.series.map((s, k) => ({
      name: s.name,
      type: "line",
      stack: spec.stack ? "s" : undefined,
      areaStyle: spec.stack ? { opacity: 0.55 } : undefined,
      data: idx.map((i) => {
        const v = s.get(byIndex[i]!);
        const marker = (spec.stack ? k === spec.series.length - 1 : k === 0) && v != null;
        return {
          value: [i, v == null ? null : Math.round(v * 100) / 100],
          symbolSize: marker ? (i === index ? 9 : 5) : 0,
          itemStyle: marker ? { color: colours.get(keys[i]), borderColor: i === index ? tokens.fg : "transparent", borderWidth: 2 } : undefined,
        };
      }),
      symbol: "circle",
      showSymbol: true,
      connectNulls: true,
      lineStyle: { color: s.color(tokens), width: spec.stack ? 0 : 1.6 },
      itemStyle: { color: s.color(tokens) },
      color: s.color(tokens),
    }));
    return {
      animation: false,
      grid: { left: 4, right: 8, top: 8, bottom: 4, containLabel: true },
      tooltip: {
        trigger: "axis",
        confine: true,
        formatter: (p: { value: [number, number | null]; seriesName: string }[]) => {
          const i = p[0].value[0];
          const c = meta.commits[i];
          return `<b>${esc(c.subject)}</b><br><span style="font-family:monospace">${short(c.sha)}</span> · ${esc(keys[i])}<br>` + p.map((x) => `${x.seriesName}: <b>${fmt(x.value[1], spec.d ?? 0)}</b>`).join("<br>");
        },
      },
      xAxis: { type: "value", min: from, max: to, show: false },
      yAxis: { type: "value", scale: !spec.stack, splitNumber: 3, splitLine: { lineStyle: { color: tokens.line } }, axisLabel: { color: tokens.dim, fontSize: 10 } },
      series: [...lineSeries, { type: "line", data: [], markLine: { symbol: "none", silent: true, label: { show: false }, lineStyle: { color: tokens.dim, type: "dashed" }, data: [{ xAxis: index }] } }],
    };
  }, [spec, from, to, index, byIndex, keys, colours, tokens, meta]);
  const onEvents = useMemo(() => ({ click: (p: unknown) => select((p as { value: [number] }).value[0]) }), [select]);
  return (
    <div style={{ minWidth: 0 }}>
      <div style={{ display: "flex", gap: 8, alignItems: "baseline", flexWrap: "wrap", fontSize: 13 }}>
        <span style={{ fontWeight: 600 }}>{spec.title}</span>
        <span className="spacer" />
        <span className="legend">
          {spec.series.map((s) => (
            <span key={s.name}>
              <span className="sw" style={{ background: s.color(tokens) }} />
              {s.name}
            </span>
          ))}
        </span>
      </div>
      <EChart option={option} onEvents={onEvents} height={150} label={spec.title + " across the measured window"} />
    </div>
  );
}

function GroupPicker({ meta, fields, setFields, keys, colours }: { meta: Meta; fields: string[]; setFields: (f: string[]) => void; keys: string[]; colours: Map<string, string> }) {
  const shown = [...new Set(keys)];
  return (
    <section className="card pad" aria-labelledby="gp">
      <h2 id="gp">Colour by provenance</h2>
      <div className="chips">
        {meta.groupFields.map((f) => (
          <label key={f} className="chip" style={{ cursor: "pointer" }}>
            <input
              type="checkbox"
              checked={fields.includes(f)}
              onChange={(e) => {
                const next = e.target.checked ? [...fields, f] : fields.filter((x) => x !== f);
                setFields(next.length ? next : ["model"]);
              }}
            />
            {f}
          </label>
        ))}
      </div>
      <div className="chips">
        {shown.map((k) => (
          <span key={k} className="chip">
            <span className="sw" style={{ background: colours.get(k) }} />
            {k}
          </span>
        ))}
      </div>
    </section>
  );
}

function GroupTable({ meta, from, to, byIndex, fields, colours }: { meta: Meta; from: number; to: number; byIndex: (Metrics | undefined)[]; fields: string[]; colours: Map<string, string> }) {
  const rows = useMemo(() => groupRows(meta.commits, byIndex, from, to, fields), [meta, byIndex, from, to, fields]);
  return (
    <section className="card" aria-labelledby="gt">
      <div className="head">
        <h2 id="gt">Step changes by provenance group</h2>
      </div>
      {!rows.length ? (
        <div className="body dim">Waiting for at least two measured commits.</div>
      ) : (
        <div className="scroll-x">
          <table className="data">
            <thead>
              <tr>
                <th scope="col">Group</th>
                <th scope="col" className="r">
                  Steps
                </th>
                <th scope="col" className="r">
                  Changed lines
                </th>
                <th scope="col" className="r">
                  Debt Δ
                </th>
                <th scope="col" className="r">
                  Δ / changed kLOC
                </th>
                <th scope="col" className="r">
                  ↑ / ↓
                </th>
                <th scope="col" className="r">
                  Hotspots Δ
                </th>
                {meta.defects && (
                  <>
                    <th scope="col" className="r">
                      Gate introduced
                    </th>
                    <th scope="col" className="r">
                      Advisory
                    </th>
                    <th scope="col" className="r">
                      Resolved
                    </th>
                  </>
                )}
                {meta.nearDuplication && (
                  <th scope="col" className="r">
                    Near-dup Δ
                  </th>
                )}
              </tr>
            </thead>
            <tbody>
              {rows.map((a) => (
                <tr key={a.key}>
                  <td>
                    <span className="sw" style={{ background: colours.get(a.key) }} />
                    {a.key}
                  </td>
                  <td className="r">{a.steps}</td>
                  <td className="r">{fmt(a.lines)}</td>
                  <td className={"r " + (a.debt > 0 ? "bad" : a.debt < 0 ? "good" : "")}>{signed(a.debt, 1)}</td>
                  <td className="r">{a.lines ? signed((a.debt * 1000) / a.lines, 1) : "–"}</td>
                  <td className="r">
                    {a.up} / {a.down}
                  </td>
                  <td className="r">{signed(a.hotspots)}</td>
                  {meta.defects && (
                    <>
                      <td className="r">{a.gate}</td>
                      <td className="r">{a.advisory}</td>
                      <td className="r">{a.resolved}</td>
                    </>
                  )}
                  {meta.nearDuplication && <td className="r">{signed(a.near)}</td>}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      <div className="body dim" style={{ fontSize: 12 }}>
        Provenance is self-reported (or inferred from Co-Authored-By lines); differences between groups are observational, not causal.
      </div>
    </section>
  );
}
