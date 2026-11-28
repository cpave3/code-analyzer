// Types and fetchers for the codeeval serve JSON API (internal/dashboard/api.go).

export interface Debt {
  complexity: number;
  nesting: number;
  duplication: number;
  architecture: number;
  total: number;
}

export interface ProvenanceRecord {
  commit: string;
  model?: string[];
  harness?: string[];
  meta?: Record<string, string>;
  source: string;
  merge?: boolean;
}

export interface Commit {
  index: number;
  sha: string;
  parent: string;
  subject: string; // untrusted text
  author: string;
  time: string;
  merge: boolean;
  provenance: ProvenanceRecord[];
}

export interface Meta {
  repository: string;
  tool: string;
  scoring: string;
  rules: string;
  clones: string;
  defects: boolean;
  nearDuplication: boolean;
  groupFields: string[];
  notes: string[];
  thresholds: { cyclomatic: number; nesting: number; minCloneTokens: number };
  weights: Record<string, number>;
  commits: Commit[];
}

export interface Metrics {
  done: boolean;
  error?: string;
  debt: Debt;
  debtPerKLOC: number | null;
  score: number | null;
  lines: number;
  files: number;
  callables: number;
  hotspots: number;
  decisionPoints: number;
  clonedCallables: number;
  ruleViolations: number;
  cycleUnits: number;
  changedLines: number;
  stepDone: boolean;
  detail: boolean;
  detailStep: boolean;
  gate?: number;
  advisory?: number;
  introducedGate?: number;
  introducedAdvisory?: number;
  resolved?: number;
  nearMisses?: number;
  fragments?: number;
}

export interface MetricsResponse {
  window: { from: number; to: number };
  done: number;
  stepsDone: number;
  detailDone: number;
  detailStepsDone: number;
  total: number;
  measuredOverall: number;
  cache: { entries: number; hits: number; misses: number };
  metrics: Record<string, Metrics>;
}

export interface GraphModule {
  id: string;
  lang: string;
  group: string;
  lines: number;
  callables: number;
  hotspots: number;
  debt: Debt;
  density: number;
  score: number | null;
}

export interface GraphCallable {
  id: string;
  module: string;
  name: string;
  file: string;
  line: number;
  cc: number;
  nest: number;
  lines: number;
  debt: number;
  cloneDebt: number;
}

export interface GraphImport {
  from: string;
  to: string;
  count?: number;
  sites?: string[];
  violation?: string;
  cycle?: boolean;
}

export interface GraphData {
  repository: string;
  modules: GraphModule[];
  callables: GraphCallable[];
  clones: { exact: boolean; tokens: number; debt: number; members: string[] }[];
  imports: GraphImport[];
  roots?: { dir: string; kind: string }[];
  cycles?: string[][];
  notes?: string[];
}

export interface Location {
  file: string;
  startLine: number;
  endLine?: number;
  symbol?: string;
}

export interface Factor {
  id: string;
  kind: string;
  module: string;
  language: string;
  location: Location;
  metrics?: Record<string, number>;
  debt: number;
  shareOfScopeDebt: number;
  cumulativeShare: number;
  explanation: string;
  related?: Location[];
  remediation: { strategy: string; target: string; steps: string[]; acceptance: string[]; caution?: string };
}

export interface LanguageCoverage {
  language: string;
  supported: boolean;
  files: number;
  bytes: number;
}

export interface ExplainReport {
  scope: { paths: string[] };
  coverage: {
    languages: LanguageCoverage[];
    supportedSourceFiles: number;
    sourceFiles: number;
    supportedByteShare: number;
    parseErrorFiles: number;
  };
  totals: {
    sourceFiles: number;
    sourceLines: number;
    callables: number;
    hotspots: number;
    debt: Debt;
    debtPerKLOC: number | null;
    score: number | null;
    modules: string[];
  };
  factors: Factor[] | null;
  factorsOmitted: number;
  concentration: { factors: number; factorsFor50Percent: number; factorsFor80Percent: number };
  guidance: string[];
  notes: string[];
}

export interface DefectFinding {
  rule: string;
  severity: string; // gate | advisory
  file: string;
  line: number;
  symbol?: string;
  message: string;
}

export interface DefectReport {
  findings: DefectFinding[] | null;
  suppressed?: DefectFinding[] | null;
  totals: { gate: number; advisory: number };
}

// normalizeGraph turns the lists Go encodes as null when empty (a commit
// with no debt-carrying functions, no clones or no imports) into arrays.
export function normalizeGraph(d: GraphData): GraphData {
  return { ...d, modules: d.modules ?? [], callables: d.callables ?? [], clones: d.clones ?? [], imports: d.imports ?? [] };
}

export class HttpError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
  }
}

async function getJSON<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) throw new HttpError(res.status, (await res.text()).trim() || res.statusText);
  return (await res.json()) as T;
}

export const api = {
  meta: () => getJSON<Meta>("/api/meta"),
  metrics: () => getJSON<MetricsResponse>("/api/metrics"),
  graph: async (sha: string) => normalizeGraph(await getJSON<GraphData>(`/api/graph?commit=${sha}`)),
  explain: (sha: string, path: string) => getJSON<ExplainReport>(`/api/explain?commit=${sha}&path=${encodeURIComponent(path)}`),
  // cached=true answers null (204) when the commit has not been checked yet,
  // so the page can offer the check instead of starting it.
  defects: async (sha: string, path: string, cached: boolean): Promise<DefectReport | null> => {
    const res = await fetch(`/api/defects?commit=${sha}&path=${encodeURIComponent(path)}${cached ? "&cached=1" : ""}`);
    if (res.status === 204) return null;
    if (!res.ok) throw new HttpError(res.status, (await res.text()).trim() || res.statusText);
    return (await res.json()) as DefectReport;
  },
  setWindow: async (from: number, to: number) => {
    const res = await fetch("/api/window", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ from, to }) });
    if (!res.ok) throw new HttpError(res.status, (await res.text()).trim() || res.statusText);
    return (await res.json()) as { from: number; to: number };
  },
};
