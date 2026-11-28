import { useCallback, useEffect, useRef, useState, type RefObject } from "react";
import type { Commit } from "../core/api";
import { NavLink, useLocation } from "react-router";
import { day, short } from "../lib/format";
import { findCommits, recentre } from "../lib/history";
import { progress, progressText, useMeta, useMetrics, useSelection, useSetWindow } from "../core/data";
import { useTheme } from "../core/theme";
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

// useDismiss closes a popover on an outside click or Escape.
function useDismiss(open: boolean, close: () => void, box: RefObject<HTMLElement | null>) {
  useEffect(() => {
    if (!open) return;
    const outside = (e: MouseEvent) => !box.current?.contains(e.target as Node) && close();
    const esc = (e: KeyboardEvent) => e.key === "Escape" && close();
    addEventListener("mousedown", outside);
    addEventListener("keydown", esc);
    return () => {
      removeEventListener("mousedown", outside);
      removeEventListener("keydown", esc);
    };
  }, [open, close, box]);
}

// usePick selects a commit; one outside the measured window first moves a
// same-size window there.
function usePick() {
  const meta = useMeta().data;
  const win = useMetrics().data?.window;
  const { select } = useSelection();
  const setWindow = useSetWindow();
  return (i: number) => {
    const inside = !win || (i >= win.from && i <= win.to);
    if (inside || !meta) return select(i);
    const [from, to] = recentre(i, win.from, win.to, meta.commits.length);
    setWindow.mutate({ from, to });
    select(i, true);
  };
}

function CommitPicker() {
  const meta = useMeta().data;
  const { index, commit } = useSelection();
  const pick = usePick();
  const [open, setOpen] = useState(false);
  const [q, setQ] = useState("");
  const box = useRef<HTMLDivElement>(null);
  const close = useCallback(() => setOpen(false), []);
  useDismiss(open, close, box);
  if (!meta || !commit) return null;
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
          <CommitList
            commits={findCommits(meta.commits, q, 300)}
            index={index}
            onPick={(i) => {
              close();
              pick(i);
            }}
          />
        </div>
      )}
    </div>
  );
}

function CommitList({ commits, index, onPick }: { commits: Commit[]; index: number; onPick: (i: number) => void }) {
  const win = useMetrics().data?.window;
  const measured = (c: Commit) => !win || (c.index >= win.from && c.index <= win.to);
  return (
    <div className="commit-list">
      {commits.map((c) => (
        <button key={c.sha} className={c.index === index ? "on" : ""} onClick={() => onPick(c.index)}>
          <span className="mono">{short(c.sha)}</span>
          <span className="subj">{c.subject}</span>
          <span className="dim" style={{ fontSize: 12, whiteSpace: "nowrap" }}>
            {measured(c) ? day(c.time) : "not measured"}
          </span>
        </button>
      ))}
      {!commits.length && <p className="dim">No commit matches.</p>}
    </div>
  );
}
