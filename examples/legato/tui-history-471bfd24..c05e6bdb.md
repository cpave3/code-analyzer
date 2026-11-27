# CodeEval history: legato-tui

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) 2026-11-27T15:46:34+11:00
- **Target:** `c05e6bdb0e338410fcc12a5b07748e2d05382a7c` (c05e6bdb0e338410fcc12a5b07748e2d05382a7c) 2026-11-27T20:01:42+11:00

- each step compares a first-parent commit with its first parent; for merge commits this is the whole merged branch, counted once
- step results are observed contributions to metrics, not proof that a workflow or person caused a quality change
- step reports do not sum to the net report: repository-wide measurements (duplication, graph) interact across steps

## First-parent steps

| # | Commit | Subject (untrusted) | Files | +/- | Go files | duplication | complexity | dependencies | obsolete | scope | correctness |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | `27265986fd97` | test(tui): pin how plan verdicts reach the conductor and the bus (i95qqzwb) | 1 | +82/-0 | 1 | 0 | 0 | 0 | 0 | 0 | n/a |
| 2 | `c05e6bdb0e33` | refactor(tui): engine behaviour reaches the TUI through service (ADR 3, i95qqzwb) | 7 | +200/-18 | 4 | 0 | 0 | +1 | 0 | 0 | n/a |
| net | `c05e6bdb0e33` | base → target (net) | 8 | +282/-18 | 5 | 0 | 0 | +1 | 0 | 0 | n/a |

## Key measurements per step

| # | Commit | dup Δ (struct redundancy) | modified cc Δ | added callables (cc) | deleted callables (cc) | hotspots before→after | edges +/- | cycles Δ | obsolete kept | API +/-/~ |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | `27265986fd97` | 0 | 0 | 0 (0) | 0 (0) | 0→0 | +0/-0 | 0 | 0 | +0/-0/~0 |
| 2 | `c05e6bdb0e33` | 0 | 0 | 1 (1) | 0 (0) | 0→0 | +0/-1 | 0 | 0 | +3/-0/~1 |
| net | `c05e6bdb0e33` | 0 | 0 | 1 (1) | 0 (0) | 0→0 | +0/-1 | 0 | 0 | +3/-0/~1 |


---

# CodeEval report: legato-tui

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `c05e6bdb0e338410fcc12a5b07748e2d05382a7c` (c05e6bdb0e338410fcc12a5b07748e2d05382a7c) tree `1c6eab010c50` 2026-11-27T20:01:42+11:00
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `3cef63633359769f` · tool 0.1.0-poc
- **Change:** 8 files, +282/-18 lines, 14 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | ok | 0 | structural redundancy 260→260 (Δ+0), exact 10→10; 0 near-structural candidate(s) |
| complexity | complexity | ok | 0 | modified callables net cyclomatic Δ+0; 1 added (cc 1), 0 deleted (cc 0); hotspots 0→0 |
| architecture_coupling | dependencies | ok | +1 | internal edges 200→199 (+0/-1); cycles 0→0; rule violations 15→14 |
| maintainability_scope | obsolete | ok | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | ok | 0 | 8 files changed (2 source); symbols +3/-0/~4 (renamed 0, formatting-only 0); exported API +3/-0/~1; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 3 test files changed, +193/-1 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 2210 | 5 | 258 |  |
| javascript | no | 33 | 0 | 0 | unsupported: no adapter |
| kotlin | no | 57 | 0 | 0 | unsupported: no adapter |
| python | no | 10 | 0 | 0 | unsupported: no adapter |
| rust | no | 34 | 0 | 0 | unsupported: no adapter |
| shell | no | 15 | 0 | 0 | unsupported: no adapter |
| sql | no | 160 | 0 | 0 | unsupported: no adapter |
| swift | no | 2 | 0 | 0 | unsupported: no adapter |
| typescript | no | 999 | 0 | 0 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
docs 3 (+39/-3); source 2 (+50/-14); test 3 (+193/-1); 

## duplication — Callable duplication delta (`duplication`)

**Status:** ok · **Score:** 0 · **Rule:** net structural redundancy Δ+0: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 260→260 (Δ+0), exact 10→10; 0 near-structural candidate(s)

Coverage: 2 changed code file(s) in scope, 2 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 3

- callables with ≥40 tokens in production source are compared repository-wide (base 6762, target 6763)
- introduced 0 redundant instance(s), removed 0; 0 pre-existing group(s) touched but unchanged in size

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 10 | 10 | 0 | – |
| structural_redundant_callables | repository | 260 | 260 | 0 | – |
| callables_considered | repository | 6762 | 6763 | +1 | – |
| redundant_instances_introduced | change | – | – | – | 0 |
| redundant_instances_removed | change | – | – | – | 0 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

## complexity — Callable complexity deltas (`complexity`)

**Status:** ok · **Score:** 0 · **Rule:** no hotspot change and \|modified net cyclomatic\| below threshold → 0

modified callables net cyclomatic Δ+0; 1 added (cc 1), 0 deleted (cc 0); hotspots 0→0

Coverage: 2 changed code file(s) in scope, 2 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 3

- modified callables: 0 increased, 0 decreased; net cyclomatic Δ+0
- added callables: 1 (cyclomatic 1); deleted callables: 0 (cyclomatic 0)
- changed-scope cyclomatic 12→13 (Δ+1), hotspots 0→0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 12 | 12 | 0 | – |
| decision_points_sum | modified_callables | 9 | 9 | 0 | – |
| lines_sum | modified_callables | 105 | 103 | -2 | – |
| max_nesting | modified_callables | 1 | 1 | 0 | – |
| cyclomatic_sum | changed_scope | 12 | 13 | +1 | – |
| max_cyclomatic | changed_scope | 4 | 4 | 0 | – |
| hotspot_callables | changed_scope | 0 | 0 | 0 | – |
| callables | changed_scope | 3 | 4 | +1 | – |
| callables_modified |  | – | – | – | 3 |
| callables_added |  | – | – | – | 1 |
| callables_deleted |  | – | – | – | 0 |
| cyclomatic_sum_added | added_callables | – | – | – | 1 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 0 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 0 |
| hotspots_added | added_callables | – | – | – | 0 |

Findings (2):

- **complexity.modified** [neutral; measured] App.handlePlanApprove: cyclomatic 4→4, nesting 1→1, lines 26→25
  - base `471bfd24a696` internal/tui/app.go:2474-2499 `App.handlePlanApprove` (before)
  - target `c05e6bdb0e33` internal/tui/app.go:2475-2499 `App.handlePlanApprove` (after)
- **complexity.modified** [neutral; measured] App.handlePlanReject: cyclomatic 4→4, nesting 1→1, lines 26→25
  - base `471bfd24a696` internal/tui/app.go:2501-2526 `App.handlePlanReject` (before)
  - target `c05e6bdb0e33` internal/tui/app.go:2501-2525 `App.handlePlanReject` (after)

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** ok · **Score:** +1 · **Rule:** rule violation resolved, or edges removed with none added → +1

internal edges 200→199 (+0/-1); cycles 0→0; rule violations 15→14

Coverage: 2 changed code file(s) in scope, 2 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 3

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| internal_units | repository | 86 | 86 | 0 | – |
| internal_edges | repository | 200 | 199 | -1 | – |
| external_dependency_refs | repository | 1028 | 1028 | 0 | – |
| cycles | repository | 0 | 0 | 0 | – |
| rule_violations | repository | 15 | 14 | -1 | – |
| edges_added | change | – | – | – | 0 |
| edges_removed | change | – | – | – | 1 |

Findings (17):

- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/tui → internal/engine/ipc
  - base `471bfd24a696` internal/tui/app.go:15-15 (removed import)
- **dependencies.fan_change** [info; measured] [go] internal/tui fan-out 19→18, fan-in 1→1
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "config is standalone, no internal imports (CLAUDE.md prose)" already violated before the change: config → internal/engine/macros (not scored)
  - target `c05e6bdb0e33` config/config.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/events (not scored)
  - target `c05e6bdb0e33` internal/tui/app.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/macros (not scored)
  - target `c05e6bdb0e33` internal/tui/app.go:15-15 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/store (not scored)
  - target `c05e6bdb0e33` internal/tui/app.go:16-16 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/workflow (not scored)
  - target `c05e6bdb0e33` internal/tui/app.go:17-17 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/board → internal/engine/store (not scored)
  - target `c05e6bdb0e33` internal/tui/board/model.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/github (not scored)
  - target `c05e6bdb0e33` internal/tui/overlay/link_pr.go:10-10 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/macros (not scored)
  - target `c05e6bdb0e33` internal/tui/overlay/macro_picker.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/workflow (not scored)
  - target `c05e6bdb0e33` internal/tui/overlay/workflow_trigger.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/store (not scored)
  - target `c05e6bdb0e33` internal/tui/overlay/workspace.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/plan → internal/engine/store (not scored)
  - target `c05e6bdb0e33` internal/tui/plan/model.go:13-13 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/report → internal/engine/analytics (not scored)
  - target `c05e6bdb0e33` internal/tui/report/model.go:12-12 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/git (not scored)
  - target `c05e6bdb0e33` internal/tui/review/messages.go:4-4 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/store (not scored)
  - target `c05e6bdb0e33` internal/tui/review/model.go:17-17 (import)
- **dependencies.rule_violation_resolved** [improvement; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" no longer violated: internal/tui → internal/engine/ipc
  - base `471bfd24a696` internal/tui/app.go:15-15 (removed import)

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** ok · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 2 changed code file(s) in scope, 2 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 3

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** ok · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

8 files changed (2 source); symbols +3/-0/~4 (renamed 0, formatting-only 0); exported API +3/-0/~1; 1 unit group(s)

Coverage: 2 changed code file(s) in scope, 2 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 3

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | docs | – | – | – | 3 |
| lines_added | docs | – | – | – | 39 |
| lines_deleted | docs | – | – | – | 3 |
| files_changed | source | – | – | – | 2 |
| lines_added | source | – | – | – | 50 |
| lines_deleted | source | – | – | – | 14 |
| files_changed | test | – | – | – | 3 |
| lines_added | test | – | – | – | 193 |
| lines_deleted | test | – | – | – | 1 |
| symbols_added | source | – | – | – | 3 |
| symbols_deleted | source | – | – | – | 0 |
| symbols_modified | source | – | – | – | 4 |
| symbols_renamed | source | – | – | – | 0 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 0 |
| exported_symbols_added | source | – | – | – | 3 |
| exported_symbols_removed | source | – | – | – | 0 |
| exported_symbols_changed | source | – | – | – | 1 |
| units_changed | source | – | – | – | 2 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

Findings (4):

- **scope.api_added** [info; measured] exported symbol added: EventStream
  - target `c05e6bdb0e33` internal/service/event_stream.go:11-14 `EventStream` (after)
- **scope.api_added** [info; measured] exported symbol added: PlanVerdict
  - target `c05e6bdb0e33` internal/service/event_stream.go:20-25 `PlanVerdict` (after)
- **scope.api_added** [info; measured] exported symbol added: SendPlanVerdict
  - target `c05e6bdb0e33` internal/service/event_stream.go:29-37 `SendPlanVerdict` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: NewApp (svc service.BoardService, syncSvc service.SyncService, agentSvc service.AgentService, prSvc service.PRTrackingService, reportSvc service.ReportService, icons theme.Icons, bus *events.Bus, editor string, workspaces []service.Workspace, tmux service.TmuxManager, workDir string, swarmSvc service.SwarmService, macrosList []macros.Macro, reviewServices ...review.Service) App → (svc service.BoardService, syncSvc service.SyncService, agentSvc service.AgentService, prSvc service.PRTrackingService, reportSvc service.ReportService, icons theme.Icons, bus service.EventStream, editor string, workspaces []service.Workspace, tmux service.TmuxManager, workDir string, swarmSvc service.SwarmService, macrosList []macros.Macro, reviewServices ...review.Service) App
  - base `471bfd24a696` internal/tui/app.go:392-444 `NewApp` (before)
  - target `c05e6bdb0e33` internal/tui/app.go:393-445 `NewApp` (after)

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 3 test files changed, +193/-1 lines; test churn is not evidence of correctness)

Coverage: 3 changed code file(s) in scope, 3 analysed, 0 with parse errors, unsupported {}, excluded by class {source:2}, non-code 3

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 3 |
| test_lines_added | context_only | – | – | – | 193 |
| test_lines_deleted | context_only | – | – | – | 1 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| added | type | EventStream | internal/service/event_stream.go:11-14 |  | source |
| added | type | PlanVerdict | internal/service/event_stream.go:20-25 |  | source |
| added | function | SendPlanVerdict | internal/service/event_stream.go:29-37 |  | source |
| modified | type | App | internal/tui/app.go:167-270 | internal/tui/app.go:168-271 | source |
| modified | function | NewApp | internal/tui/app.go:393-445 | internal/tui/app.go:392-444 | source |
| modified | method | App.handlePlanApprove | internal/tui/app.go:2475-2499 | internal/tui/app.go:2474-2499 | source |
| modified | method | App.handlePlanReject | internal/tui/app.go:2501-2525 | internal/tui/app.go:2501-2526 | source |
| modified | function | newAgentUIApp | internal/tui/app_agent_ui_test.go:58-64 | internal/tui/app_agent_ui_test.go:58-64 | test |
| added | function | verdictListener | internal/tui/app_plan_verdict_test.go:16-31 |  | test |
| added | function | receive | internal/tui/app_plan_verdict_test.go:33-43 |  | test |
| added | function | TestPlanVerdictsReachTheConductorAndTheBus | internal/tui/app_plan_verdict_test.go:47-82 |  | test |
| added | function | set | internal/tui/architecture_test.go:38-44 |  | test |
| added | function | eventsValue | internal/tui/architecture_test.go:46-49 |  | test |
| added | function | TestTUIUsesOnlyEngineValueTypesAndPureHelpers | internal/tui/architecture_test.go:51-110 |  | test |

## Diagnostics

No diagnostics for changed files.

## Limitations

- POC: measurements of structure, not proof of design quality; lower metrics do not by themselves prove better architecture
- scores are per analyzer and directional; there is deliberately no weighted overall score
- correctness confidence is not assessed; the evaluated repository's code is never executed
- duplication, complexity and obsolete-code analyzers measure production source only; tests, generated, vendor and docs are counted in scope and coverage
- structural/near-structural similarity marks candidates, not semantic equivalence or proven missed reuse
- go: no type checking: references are matched by name, not resolved
- go: no call graph or interface dispatch resolution
- go: build tags and GOOS/GOARCH file variants are all indexed together
- go: import cycles are compile errors in Go; detected cycles indicate non-building code
- csharp: no adapter — files are counted as unsupported coverage, never scored
- java: no adapter — files are counted as unsupported coverage, never scored
- javascript: no adapter — files are counted as unsupported coverage, never scored
- kotlin: no adapter — files are counted as unsupported coverage, never scored
- python: no adapter — files are counted as unsupported coverage, never scored
- rust: no adapter — files are counted as unsupported coverage, never scored
- shell: no adapter — files are counted as unsupported coverage, never scored
- sql: no adapter — files are counted as unsupported coverage, never scored
- swift: no adapter — files are counted as unsupported coverage, never scored
- typescript: no adapter — files are counted as unsupported coverage, never scored
