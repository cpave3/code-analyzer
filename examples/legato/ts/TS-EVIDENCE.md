# Legato with TypeScript support (codeeval 0.2.0)

This is the first measurement of Legato's TypeScript, using the new
TypeScript adapter. It is a **new measurement with a new tool version**
(0.2.0-poc). The Go-only v1 reports made with the frozen `5e0c679` build
(0.1.0-poc), and the v2 Go-only reports, are unchanged and remain the
provenance of record for those runs. Numbers here that include TypeScript
are **not comparable** to the Go-only numbers.

## Setup

| | |
|---|---|
| Tool | codeeval 0.2.0-poc from this repo's working tree. Adapters: `go-adapter/3+tree-sitter-go/0.25.0` and `ts-adapter/1+tree-sitter-typescript/0.23.2` |
| Scoring | `scoring/v1` for `evaluate`/`history`; `scoring/v2-experimental` for `score` |
| Config | `legato-architecture.json` (configHash `3cef63633359769f`). Its rules cover only Go paths; **no TypeScript architecture rules are configured** |
| Revisions | `9a61e9f8` → `471bfd24` (the original two-day window) and `471bfd24` → `bb13a424` (the cleanup stack, which changed no TypeScript) |
| Objects | The existing private `--no-checkout` clone `.tmp/legato-stack`. No network; Legato's own scripts, `tsc` and `node_modules` were never run. The shared Legato checkout was untouched (HEAD `471bfd24`, status fingerprint `3ab8895f206d466e` before and after) |
| Run | 2026-11-27T21:14+11:00. v1 net 0.95 s, v1 history (23 merges) 10 s, v2 snapshot about 0.7 s, v2 window with steps 9 s |
| Dependency added | `github.com/tree-sitter/tree-sitter-typescript` v0.23.2, fetched once through the Go module proxy with checksum verification (only this tool's `go.mod`) |

Reports in this folder: `v1-net-9a61e9f8..471bfd24.{json,md}`,
`v1-history-9a61e9f8..471bfd24.{json,md}`, `v2-snapshot-471bfd24.*`,
`v2-snapshot-bb13a424.*`, `v2-delta-window-9a61e9f8..471bfd24.*` (with
per-merge steps and per-language ledgers) and
`v2-delta-stack-471bfd24..bb13a424.*`.

## Coverage

| | Go only (0.1.0) | Go and TypeScript (0.2.0) |
|---|---:|---:|
| Production source files measured at `471bfd24` | 988 of 1,846 | **1,553 of 1,846** |
| Share of production source bytes | 61.8% | **90.0%** |
| Changed source files measured in the window | 201 of 296 | **282 of 293**: 11 still unsupported (SQL 7, JS 2, Kotlin 2). The total drops from 296 to 293 because the TS adapter classifies 3 more paths as tests (`tests/`, `e2e/`, `__mocks__`, `.stories.`) |
| Merges with no measurable code | 6 (#403, #407, #420, #423, #428, #431) | **0** |

Still not measured: SQL 160 files, Kotlin 54, Rust 34, JavaScript 23,
shell 15, Python 5, Swift 2. TypeScript **test** files (434) are classified
as tests, so they appear in coverage but not in duplication, complexity or
debt. TypeScript at `471bfd24`: 565 source files, 78,035 lines and 2,670
callables.

**Parse quality.** 43 TypeScript files at `471bfd24` have syntax errors under
the grammar. 39 are test files and **4 are source files**:
`web/src/components/AgentStudioPanel.tsx`, `web/src/lib/board-types.ts`,
`web/src/pages/PlanPage.tsx` and `web/src/pages/ReviewTourPage.tsx`. These
come from gaps in tree-sitter-typescript 0.23.2, not from Legato's code:
- A semicolon-free interface member whose name starts with `in` (for
  example `in_review: number`) is misread as the `in` operator.
- `typeof import("…")` inside type arguments is not parsed.

Symbols in error regions are excluded from duplication. The files are
reported as `parse_error` diagnostics. **Import resolution:** 0 unresolved
internal imports. An alias into `node_modules` (`desktop`'s `react` →
`node_modules/@types/react`) is treated as external.

## v1 window `9a61e9f8 → 471bfd24` with TypeScript

| Analyzer | Go only (frozen 0.1.0) | Go and TypeScript (0.2.0) | TypeScript's share (difference) |
|---|---|---|---|
| duplication | −2: structural 243→260, exact 8→10 | −2: structural 277→295, exact 17→20 | structural 34→35 (+1), exact 9→10 (+1) |
| complexity | −2: modified net +194; 569 added; 84 new hotspots | −2: modified net **+659**; 740 added (cc 3,834); 110 new hotspots | modified net **+465** (70 modified TS callables: 60 up, 1 down); 171 added (cc 858); 26 new hotspots |
| dependencies | 0: edges 193→200, cycles 0→0 | 0: edges 262→279, cycles 3→3 | TS edges 69→79 (+11/−1); 3 **pre-existing** TS directory cycles; no TS rules |
| obsolete / scope | 0 / 0 | 0 / 0 | |

**Merges that were unscored before are now scored** (v1 history, TypeScript
only): #403 complexity −2; #407 0; #420 complexity −1; #423 duplication −1
and complexity −2; #428 complexity −1; #431 complexity −1, obsolete +1.

**Largest TypeScript complexity growth in existing callables:**

| Callable | cc | Nesting | Lines | Location (target) |
|---|---|---|---|---|
| `startGlance` | 204 → **320** (+116) | 4 | 666 → 970 | `packages/even-hub/src/app.ts:538` |
| `Outline` | 172 → 204 (+32) | 6 | 383 → 454 | `web/src/components/workflow-editor/WorkflowOutline.tsx:194` |
| `WorkflowEditPage` | 167 → 193 (+26) | 4 | 364 → 430 | `web/src/pages/WorkflowEditPage.tsx:166` |
| `Graph` | 20 → 44 (+24) | 2 → 3 | 86 → 131 | `web/src/components/workflow-graph/WorkflowGraph.tsx:392` |
| `RunsTab` | 9 → 32 (+23) | 1 → 2 | 14 → 65 | `web/src/pages/WorkflowsPage.tsx:158` |
| `runNodeExtras` | 14 → 32 (+18) | 3 → 4 | 21 → 39 | `web/src/lib/workflowRunGraph.ts:75` |
| `SpawnAgentDialog` | 236 → 250 (+14) | 5 | 941 → 947 | `web/src/components/SpawnAgentDialog.tsx:104` |

Largest reduction: `WorkflowsTab` 25 → 14 (`web/src/pages/WorkflowsPage.tsx`).
**New TypeScript hotspots (26)**, the largest being `WorkflowCanvas` cc 38
(nesting 4), `PinControls` 33, `WorkflowList` 28 and `WorkflowTagsDialog` 26.
**New exact duplicate:** `recommended` in
`packages/even-hub/src/questions.ts:113` is identical to the pre-existing
`web/src/lib/questionDrafts.ts:28`, a cross-package copy similar to the Go
`retryAfter` case.

## v2 absolute snapshot `471bfd24` (TypeScript unchanged at `bb13a424`)

| Language | Score | Debt | complexity | nesting | duplication | architecture | Lines | Debt/kLOC | Hotspots |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Go | 33.93 | 11,611.30 | 10,027 | 618 | 816.30 | 150 | 238,533 | 48.68 | 1,288 |
| TypeScript | 19.21 | 8,203.45 | 7,599 | 242 | 102.45 | 260 | 78,035 | 105.13 | 488 |
| Codebase | 28.54 | 19,814.75 | | | | | 316,568 | | |

The TypeScript debt density is about twice Go's, but **the two are not
directly comparable**. The TS adapter measures inline callbacks and JSX
conditionals (`cond && <X/>`, ternaries) as part of the enclosing React
component, so large components carry very high cyclomatic counts. The
largest TypeScript callable debts are all components or app entry points:
`NativeAgentPanel` cc 331 (1,184 lines), `startGlance` 320, `BoardPage` 297,
`SpawnAgentDialog` 250, `AgentsPage` 247, `Outline` 204 and
`WorkflowEditPage` 193.

Largest TypeScript modules by debt: `web/src/components` (2,178, 19.9k
lines), `web/src/pages` (1,595 at **204 debt/kLOC**, 7.8k lines),
`web/src/lib` (1,218), `web/src/components/workflow-editor` (504) and
`web/src/components/board` (387).

**TypeScript architecture debt (260)** is 13 directories inside 3
directory-level import cycles. There are no configured rule violations,
because no TS rules exist. The cycles are:
`browser-extension/src ↔ browser-extension/src/scope-source`;
`desktop/src ↔ desktop/src/components`; and a 9-directory cycle across
`web/src/components`, `components/agent-graph`, `components/needsYou`,
`components/settings`, `components/ui`, `components/workflow-graph`,
`web/src/hooks`, `web/src/lib` and `web/src/voice`. Directory cycles are
legal in TypeScript. They indicate coupling, not build errors, and whether
they matter is a design judgement.

**Largest TypeScript clone groups:** a local `ok(res)` helper exactly
duplicated in 6 files under `web/src/lib` (`jobClones`, `mcpConnections`,
`repositories`, `repositoryGuidance`, `resources`, `skills`; 325 redundant
tokens), and structurally identical `fetch*`/`list*` wrappers across
`configSettings`, `reeve`, `devices` and `workers`.

## v2 window delta `9a61e9f8 → 471bfd24` by language

| Language | Δ debt | complexity | nesting | duplication | architecture | Δ lines | Δ decision points | Δ hotspots |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| Go | +828.00 | +752 | +20 | +56.00 | 0 | +14,451 | +2,586 | +91 |
| TypeScript | **+599.55** | +587 | +10 | +2.55 | 0 | +4,046 | +1,128 | +29 |
| Total | **+1,427.55** | +1,339 | +30 | +58.55 | 0 | | | |

The 23 per-merge steps sum exactly to the net (+1,427.55). That is
consistency under this fixed model and config, not calibration. Largest
TypeScript debt increases by merge:

| Merge | TS Δ debt |
|---|---:|
| #423 | +142.55 |
| #418 | +78.00 |
| #408 | +77.00 |
| #403 | +70.00 |
| #395 | +63.00 |
| #407 | +55.00 |
| #430 | +41.00 |
| #425 | +40.00 |

The cleanup stack (`471bfd24 → bb13a424`) changed no TypeScript, so its
TypeScript delta is 0 and its Go delta is unchanged at −224.10.

## What this adds, and its limits

**Adds:**
- Legato's web, desktop, browser-extension and Even Hub code is now in the
  picture. Every merge in the window is measurable.
- TypeScript complexity growth in the window (+465 modified cc) is more than
  double Go's (+194), concentrated in `startGlance` and the workflow editor
  and graph components.
- One more cross-package exact copy (`recommended`) and a six-way `ok()`
  helper copy surfaced.
- Three pre-existing directory cycles are named, including a 9-directory
  cycle in `web/src`.

**Limits:**
- **Cross-language numbers are not comparable.** TypeScript complexity
  includes inline callbacks and JSX conditionals.
- **4 source files have grammar-gap parse errors.**
- No TypeScript architecture rules are configured, so architecture debt is
  cycles only.
- **The graph is directory-level:** it follows neither node_modules,
  package `exports` maps nor re-export chains.
- No type resolution; references are matched by name.
- **Weights and thresholds are uncalibrated for TypeScript.** The Go
  thresholds (cc 10, nesting 4) were reused.
- Neither the v2 score nor debt should be optimised on its own.
- No correctness claim.
