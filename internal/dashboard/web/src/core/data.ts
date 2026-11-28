// Data hooks over the serve API, and the selected commit (kept in the URL
// as ?commit=SHA so any view can be linked).
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useCallback, useMemo } from "react";
import { useSearchParams } from "react-router";
import { api, type Metrics, type MetricsResponse } from "./api";

export function useMeta() {
  return useQuery({ queryKey: ["meta"], queryFn: api.meta, staleTime: Infinity, retry: 1 });
}

export function busy(r: MetricsResponse | undefined): boolean {
  if (!r) return true;
  return r.done < r.total || r.stepsDone < r.total || r.detailDone < r.total || r.detailStepsDone < r.total;
}

// progress is the share of the window's measuring work that is done (0–1).
export function progress(r: MetricsResponse | undefined): number {
  if (!r || !r.total) return 0;
  return (r.done + r.stepsDone + r.detailDone + r.detailStepsDone) / (4 * r.total);
}

export function progressText(r: MetricsResponse | undefined): string {
  if (!r) return "Loading…";
  if (r.done < r.total) return `Measuring debt: ${r.done} of ${r.total} commits, newest first`;
  if (r.stepsDone < r.total) return `Comparing steps: ${r.stepsDone} of ${r.total}`;
  if (r.detailDone < r.total) return `Checking defects and near-duplicates: ${r.detailDone} of ${r.total}`;
  if (r.detailStepsDone < r.total) return `Comparing defects: ${r.detailStepsDone} of ${r.total}`;
  return `${r.total} commits measured`;
}

export function useMetrics() {
  const q = useQuery({
    queryKey: ["metrics"],
    queryFn: api.metrics,
    refetchInterval: (query) => (busy(query.state.data) ? 1500 : 5000),
    retry: 0,
  });
  const n = useMeta().data?.commits.length ?? 0;
  const byIndex = useMemo(() => {
    const out: (Metrics | undefined)[] = new Array(n);
    for (const [k, m] of Object.entries(q.data?.metrics ?? {})) out[+k] = m;
    return out;
  }, [q.data, n]);
  return { ...q, byIndex };
}

export function useSetWindow() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ from, to }: { from: number; to: number }) => api.setWindow(from, to),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["metrics"] }),
  });
}

// useSelection returns the selected commit index (inside the measured
// window; the newest measured commit by default) and a setter.
export function useSelection() {
  const [params, setParams] = useSearchParams();
  const meta = useMeta().data;
  const win = useMetrics().data?.window;
  const sha = params.get("commit");
  const index = useMemo(() => {
    if (!meta) return -1;
    const found = sha ? meta.commits.findIndex((c) => c.sha === sha || c.sha.startsWith(sha)) : -1;
    if (found >= 0) return found;
    return win ? win.to : meta.commits.length - 1;
  }, [meta, sha, win]);
  const select = useCallback(
    (i: number, force = false) => {
      if (!meta) return;
      const lo = force ? 0 : (win?.from ?? 0),
        hi = force ? meta.commits.length - 1 : (win?.to ?? meta.commits.length - 1);
      const j = Math.max(lo, Math.min(hi, i));
      setParams(
        (p) => {
          const n = new URLSearchParams(p);
          n.set("commit", meta.commits[j].sha.slice(0, 12));
          return n;
        },
        { replace: true },
      );
    },
    [meta, win, setParams],
  );
  return { index, commit: meta && index >= 0 ? meta.commits[index] : undefined, select };
}

export function useGraph(sha: string | undefined) {
  return useQuery({ queryKey: ["graph", sha], queryFn: () => api.graph(sha!), enabled: !!sha, staleTime: Infinity, retry: 1 });
}

export function useExplain(sha: string | undefined, path: string) {
  return useQuery({ queryKey: ["explain", sha, path], queryFn: () => api.explain(sha!, path), enabled: !!sha, staleTime: Infinity, retry: 1 });
}

// useDefects reads a commit's defect findings for a path without running
// the rules; run() checks the whole commit once, after which every path of
// that commit is instant.
export function useDefects(sha: string | undefined, path: string) {
  const qc = useQueryClient();
  const q = useQuery({ queryKey: ["defects", sha, path], queryFn: () => api.defects(sha!, path, true), enabled: !!sha, staleTime: Infinity, retry: 0 });
  const run = useMutation({
    mutationFn: () => api.defects(sha!, path, false),
    onSuccess: (r) => {
      qc.setQueryData(["defects", sha, path], r);
      qc.invalidateQueries({ queryKey: ["defects", sha] });
    },
  });
  return { ...q, run };
}
