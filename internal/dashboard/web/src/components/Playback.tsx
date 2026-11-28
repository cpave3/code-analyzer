// Stepping through the window: ← and → move one commit, space plays.
import { createContext, useCallback, useContext, useEffect, useRef, useState, type ReactNode } from "react";
import { useMetrics, useSelection } from "../data";

interface PlaybackValue {
  playing: boolean;
  toggle: () => void;
}
const Ctx = createContext<PlaybackValue>({ playing: false, toggle: () => {} });

export function PlaybackProvider({ children }: { children: ReactNode }) {
  const { index, select } = useSelection();
  const win = useMetrics().data?.window;
  const [playing, setPlaying] = useState(false);
  const at = useRef(index);
  at.current = index;
  const toggle = useCallback(() => {
    setPlaying((p) => {
      if (!p && win && at.current >= win.to) select(win.from);
      return !p;
    });
  }, [win, select]);
  useEffect(() => {
    if (!playing || !win) return;
    const t = setInterval(() => {
      if (at.current >= win.to) setPlaying(false);
      else select(at.current + 1);
    }, 1300);
    return () => clearInterval(t);
  }, [playing, win, select]);
  useEffect(() => {
    const on = (e: KeyboardEvent) => {
      const t = e.target as HTMLElement;
      if (t.closest("input:not([type=range]), select, textarea, [contenteditable]") || e.metaKey || e.ctrlKey || e.altKey) return;
      if (e.key === "ArrowLeft") select(at.current - 1);
      else if (e.key === "ArrowRight") select(at.current + 1);
      else if (e.key === " " && !t.closest("button, a")) {
        e.preventDefault();
        toggle();
      }
    };
    addEventListener("keydown", on);
    return () => removeEventListener("keydown", on);
  }, [select, toggle]);
  return <Ctx.Provider value={{ playing, toggle }}>{children}</Ctx.Provider>;
}

export const usePlayback = () => useContext(Ctx);
