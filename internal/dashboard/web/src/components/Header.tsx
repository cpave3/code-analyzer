import { useEffect, useMemo, useRef, useState } from "react";
import { NavLink, useLocation } from "react-router";
import { day, short } from "../lib/format";
import { recentre } from "../lib/history";
import { progress, progressText, useMeta, useMetrics, useSelection, useSetWindow } from "../data";
import { useTheme } from "../theme";
import { Chevron, Info, Logo, Moon, Sun } from "./icons";

const VIEWS = [
  { to: "/", label: "Overview", end: true },
  { to: "/modules", label: "Modules", end: false },
  { to: "/changes", label: "Changes", end: false },
  { to: "/graph", label: "Graph", end: false },
];

export function Header({ onAbout }: { onAbout: () => void }) {
  const { search } = useLocation();
  const keep = new URLSearchParams(search).get("commit");
  const qs = keep ? "?commit=" + keep : "";
  const metrics = useMetrics();
  const { dark, toggle } = useTheme();
  const p = progress(metrics.data);
  return (
    <header className="topbar">
      <div className="brand">
        <Logo />
        CodeEval
      </div>
      <nav aria-label="Views" className="tabs">
        {VIEWS.map((v) => (
          <NavLink key={v.to} to={v.to + qs} end={v.end} className={({ isActive }) => (isActive ? "active" : "")}>
            {v.label}
          </NavLink>
        ))}
      </nav>
      <div className="spacer" />
      {metrics.data && p < 1 && (
        <div className="progress" role="status">
          <div className="bar">
            <div style={{ width: `${100 * p}%` }} />
          </div>
          <span>{progressText(metrics.data)}</span>
        </div>
      )}
      <CommitPicker />
      <button className="btn icon" aria-label={dark ? "Use light theme" : "Use dark theme"} title={dark ? "Light theme" : "Dark theme"} onClick={toggle}>
        {dark ? <Sun /> : <Moon />}
      </button>
      <button className="btn icon" aria-label="About these numbers" title="About these numbers" onClick={onAbout}>
        <Info />
      </button>
    </header>
  );
}

function CommitPicker() {
  const meta = useMeta().data;
  const win = useMetrics().data?.window;
  const { index, commit, select } = useSelection();
  const setWindow = useSetWindow();
  const [open, setOpen] = useState(false);
  const [q, setQ] = useState("");
  const box = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!open) return;
    const close = (e: MouseEvent) => {
      if (!box.current?.contains(e.target as Node)) setOpen(false);
    };
    const esc = (e: KeyboardEvent) => e.key === "Escape" && setOpen(false);
    addEventListener("mousedown", close);
    addEventListener("keydown", esc);
    return () => {
      removeEventListener("mousedown", close);
      removeEventListener("keydown", esc);
    };
  }, [open]);
  const list = useMemo(() => {
    if (!meta) return [];
    const needle = q.trim().toLowerCase();
    const out = [];
    for (let i = meta.commits.length - 1; i >= 0 && out.length < 300; i--) {
      const c = meta.commits[i];
      if (!needle || c.sha.startsWith(needle) || c.subject.toLowerCase().includes(needle) || c.author.toLowerCase().includes(needle)) out.push(c);
    }
    return out;
  }, [meta, q]);
  if (!meta || !commit) return null;
  const pick = (i: number) => {
    setOpen(false);
    if (win && (i < win.from || i > win.to)) {
      const [from, to] = recentre(i, win.from, win.to, meta.commits.length);
      setWindow.mutate({ from, to });
      select(i, true);
    } else select(i);
  };
  return (
    <div ref={box} style={{ position: "relative" }}>
      <button className="btn" aria-haspopup="dialog" aria-expanded={open} onClick={() => setOpen((o) => !o)} title="Choose the commit every view shows">
        <span className="dim" style={{ fontSize: 13 }}>
          Commit
        </span>
        <span className="mono">{short(commit.sha)}</span>
        <span className="dim" style={{ fontSize: 13 }}>
          {index + 1}/{meta.commits.length} · {day(commit.time)}
        </span>
        <Chevron />
      </button>
      {open && (
        <div className="popover" role="dialog" aria-label="Choose a commit">
          <input type="search" autoFocus placeholder="SHA, subject or author" aria-label="Find a commit" value={q} onChange={(e) => setQ(e.target.value)} />
          <div className="commit-list">
            {list.map((c) => {
              const inWin = win && c.index >= win.from && c.index <= win.to;
              return (
                <button key={c.sha} className={c.index === index ? "on" : ""} onClick={() => pick(c.index)}>
                  <span className="mono">{short(c.sha)}</span>
                  <span className="subj">{c.subject}</span>
                  <span className="dim" style={{ fontSize: 12, whiteSpace: "nowrap" }}>
                    {inWin ? day(c.time) : "not measured"}
                  </span>
                </button>
              );
            })}
            {!list.length && <p className="dim">No commit matches.</p>}
          </div>
        </div>
      )}
    </div>
  );
}
