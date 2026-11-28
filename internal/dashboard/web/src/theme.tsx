// Light and dark themes: the system setting by default, or a choice the
// viewer made (kept in this browser). Charts read the CSS tokens.
import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from "react";

export type Mode = "system" | "light" | "dark";

interface ThemeValue {
  mode: Mode;
  dark: boolean;
  toggle: () => void;
  tokens: Tokens;
}

export interface Tokens {
  bg: string;
  panel: string;
  fg: string;
  dim: string;
  line: string;
  track: string;
  accent: string;
  good: string;
  bad: string;
  components: { complexity: string; nesting: string; duplication: string; architecture: string };
  ramp: string[];
}

const ThemeContext = createContext<ThemeValue | null>(null);

function readMode(): Mode {
  try {
    const v = localStorage.getItem("codeeval.theme");
    return v === "light" || v === "dark" ? v : "system";
  } catch {
    return "system";
  }
}

function systemDark(): boolean {
  return typeof matchMedia === "function" && matchMedia("(prefers-color-scheme: dark)").matches;
}

function readTokens(): Tokens {
  const s = getComputedStyle(document.documentElement);
  const v = (n: string) => s.getPropertyValue(n).trim();
  return {
    bg: v("--bg"),
    panel: v("--panel"),
    fg: v("--fg"),
    dim: v("--dim"),
    line: v("--line"),
    track: v("--track"),
    accent: v("--accent"),
    good: v("--good"),
    bad: v("--bad"),
    components: { complexity: v("--c1"), nesting: v("--c2"), duplication: v("--c3"), architecture: v("--c4") },
    ramp: [v("--d0"), v("--d1"), v("--d2"), v("--d3"), v("--d4")],
  };
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setMode] = useState<Mode>(readMode);
  const [sys, setSys] = useState(systemDark);
  useEffect(() => {
    if (typeof matchMedia !== "function") return;
    const mq = matchMedia("(prefers-color-scheme: dark)");
    const on = () => setSys(mq.matches);
    mq.addEventListener("change", on);
    return () => mq.removeEventListener("change", on);
  }, []);
  const dark = mode === "dark" || (mode === "system" && sys);
  const [tokens, setTokens] = useState<Tokens>(readTokens);
  useEffect(() => {
    const el = document.documentElement;
    if (mode === "system") el.removeAttribute("data-theme");
    else el.setAttribute("data-theme", mode);
    setTokens(readTokens());
  }, [mode, sys]);
  const value = useMemo<ThemeValue>(
    () => ({
      mode,
      dark,
      tokens,
      toggle: () => {
        const next: Mode = dark ? "light" : "dark";
        setMode(next);
        try {
          localStorage.setItem("codeeval.theme", next);
        } catch {
          /* private window: the choice lasts for this page only */
        }
      },
    }),
    [mode, dark, tokens],
  );
  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme(): ThemeValue {
  const v = useContext(ThemeContext);
  if (!v) throw new Error("useTheme outside ThemeProvider");
  return v;
}
