# CodeEval report: legato-b3

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `b9696a9c71129006afb881645501d29520dbbcfd` (b9696a9c71129006afb881645501d29520dbbcfd) tree `29c8697082bd` 2026-11-27T19:49:14+11:00
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `3cef63633359769f` · tool 0.1.0-poc
- **Change:** 7 files, +208/-54 lines, 11 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | ok | 0 | structural redundancy 260→260 (Δ+0), exact 10→10; 0 near-structural candidate(s) |
| complexity | complexity | ok | +1 | modified callables net cyclomatic Δ-16; 4 added (cc 23), 0 deleted (cc 0); hotspots 3→2 |
| architecture_coupling | dependencies | ok | 0 | internal edges 200→200 (+0/-0); cycles 0→0; rule violations 15→15 |
| maintainability_scope | obsolete | ok | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | ok | 0 | 7 files changed (4 source); symbols +4/-0/~4 (renamed 0, formatting-only 0); exported API +0/-0/~0; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 3 test files changed, +116/-0 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 2208 | 7 | 262 |  |
| javascript | no | 33 | 0 | 0 | unsupported: no adapter |
| kotlin | no | 57 | 0 | 0 | unsupported: no adapter |
| python | no | 10 | 0 | 0 | unsupported: no adapter |
| rust | no | 34 | 0 | 0 | unsupported: no adapter |
| shell | no | 15 | 0 | 0 | unsupported: no adapter |
| sql | no | 160 | 0 | 0 | unsupported: no adapter |
| swift | no | 2 | 0 | 0 | unsupported: no adapter |
| typescript | no | 999 | 0 | 0 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
source 4 (+92/-54); test 3 (+116/-0); 

## duplication — Callable duplication delta (`duplication`)

**Status:** ok · **Score:** 0 · **Rule:** net structural redundancy Δ+0: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 260→260 (Δ+0), exact 10→10; 0 near-structural candidate(s)

Coverage: 4 changed code file(s) in scope, 4 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 0

- callables with ≥40 tokens in production source are compared repository-wide (base 6762, target 6766)
- introduced 0 redundant instance(s), removed 0; 0 pre-existing group(s) touched but unchanged in size

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 10 | 10 | 0 | – |
| structural_redundant_callables | repository | 260 | 260 | 0 | – |
| callables_considered | repository | 6762 | 6766 | +4 | – |
| redundant_instances_introduced | change | – | – | – | 0 |
| redundant_instances_removed | change | – | – | – | 0 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

## complexity — Callable complexity deltas (`complexity`)

**Status:** ok · **Score:** +1 · **Rule:** pre-existing hotspotΔ&lt;0 or modified net cyclomatic ≤ -3 → +1

modified callables net cyclomatic Δ-16; 4 added (cc 23), 0 deleted (cc 0); hotspots 3→2

Coverage: 4 changed code file(s) in scope, 4 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 0

- modified callables: 0 increased, 4 decreased; net cyclomatic Δ-16
- added callables: 4 (cyclomatic 23); deleted callables: 0 (cyclomatic 0)
- changed-scope cyclomatic 64→71 (Δ+7), hotspots 3→2
- complexity removed from modified callables (Δ-16) is matched by 4 added callables (cyclomatic 23): this is extraction/redistribution, not net removal

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 64 | 48 | -16 | – |
| decision_points_sum | modified_callables | 60 | 44 | -16 | – |
| lines_sum | modified_callables | 206 | 160 | -46 | – |
| max_nesting | modified_callables | 4 | 4 | 0 | – |
| cyclomatic_sum | changed_scope | 64 | 71 | +7 | – |
| max_cyclomatic | changed_scope | 21 | 18 | -3 | – |
| hotspot_callables | changed_scope | 3 | 2 | -1 | – |
| callables | changed_scope | 4 | 8 | +4 | – |
| callables_modified |  | – | – | – | 4 |
| callables_added |  | – | – | – | 4 |
| callables_deleted |  | – | – | – | 0 |
| cyclomatic_sum_added | added_callables | – | – | – | 23 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 0 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | -1 |
| hotspots_added | added_callables | – | – | – | 0 |

Findings (4):

- **complexity.modified** [improvement; measured] upgrader.mapText: cyclomatic 13→8, nesting 4→3, lines 31→21
  - base `471bfd24a696` internal/engine/workflow/upgrade.go:528-558 `upgrader.mapText` (before)
  - target `b9696a9c7112` internal/engine/workflow/upgrade.go:528-548 `upgrader.mapText` (after)
- **complexity.modified** [improvement; measured] agentService.historyConversation: cyclomatic 9→6, nesting 2→1, lines 24→14
  - base `471bfd24a696` internal/service/agent_resume.go:190-213 `agentService.historyConversation` (before)
  - target `b9696a9c7112` internal/service/agent_resume.go:190-203 `agentService.historyConversation` (after)
- **complexity.modified** [improvement; measured] SessionHistoryService.Recover: cyclomatic 21→16, nesting 4→4, lines 75→59
  - base `471bfd24a696` internal/service/session_recovery.go:28-102 `SessionHistoryService.Recover` (before)
  - target `b9696a9c7112` internal/service/session_recovery.go:28-86 `SessionHistoryService.Recover` (after)
- **complexity.modified** [improvement; measured] TicketCreateExecutor.Start: cyclomatic 21→18, nesting 3→3, lines 76→66
  - base `471bfd24a696` internal/service/workflow_ticket.go:171-246 `TicketCreateExecutor.Start` (before)
  - target `b9696a9c7112` internal/service/workflow_ticket.go:171-236 `TicketCreateExecutor.Start` (after)

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** ok · **Score:** 0 · **Rule:** no cycles or rule violations changed → 0 (new edges alone are reported, not judged: features legitimately add dependencies)

internal edges 200→200 (+0/-0); cycles 0→0; rule violations 15→15

Coverage: 4 changed code file(s) in scope, 4 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 0

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| internal_units | repository | 86 | 86 | 0 | – |
| internal_edges | repository | 200 | 200 | 0 | – |
| external_dependency_refs | repository | 1028 | 1028 | 0 | – |
| cycles | repository | 0 | 0 | 0 | – |
| rule_violations | repository | 15 | 15 | 0 | – |
| edges_added | change | – | – | – | 0 |
| edges_removed | change | – | – | – | 0 |

Findings (15):

- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "config is standalone, no internal imports (CLAUDE.md prose)" already violated before the change: config → internal/engine/macros (not scored)
  - target `b9696a9c7112` config/config.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/events (not scored)
  - target `b9696a9c7112` internal/tui/app.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/ipc (not scored)
  - target `b9696a9c7112` internal/tui/app.go:15-15 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/macros (not scored)
  - target `b9696a9c7112` internal/tui/app.go:16-16 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/store (not scored)
  - target `b9696a9c7112` internal/tui/app.go:17-17 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/workflow (not scored)
  - target `b9696a9c7112` internal/tui/app.go:18-18 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/board → internal/engine/store (not scored)
  - target `b9696a9c7112` internal/tui/board/model.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/github (not scored)
  - target `b9696a9c7112` internal/tui/overlay/link_pr.go:10-10 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/macros (not scored)
  - target `b9696a9c7112` internal/tui/overlay/macro_picker.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/workflow (not scored)
  - target `b9696a9c7112` internal/tui/overlay/workflow_trigger.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/store (not scored)
  - target `b9696a9c7112` internal/tui/overlay/workspace.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/plan → internal/engine/store (not scored)
  - target `b9696a9c7112` internal/tui/plan/model.go:13-13 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/report → internal/engine/analytics (not scored)
  - target `b9696a9c7112` internal/tui/report/model.go:12-12 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/git (not scored)
  - target `b9696a9c7112` internal/tui/review/messages.go:4-4 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/store (not scored)
  - target `b9696a9c7112` internal/tui/review/model.go:17-17 (import)

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** ok · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 4 changed code file(s) in scope, 4 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** ok · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

7 files changed (4 source); symbols +4/-0/~4 (renamed 0, formatting-only 0); exported API +0/-0/~0; 1 unit group(s)

Coverage: 4 changed code file(s) in scope, 4 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 0

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | source | – | – | – | 4 |
| lines_added | source | – | – | – | 92 |
| lines_deleted | source | – | – | – | 54 |
| files_changed | test | – | – | – | 3 |
| lines_added | test | – | – | – | 116 |
| lines_deleted | test | – | – | – | 0 |
| symbols_added | source | – | – | – | 4 |
| symbols_deleted | source | – | – | – | 0 |
| symbols_modified | source | – | – | – | 4 |
| symbols_renamed | source | – | – | – | 0 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 0 |
| exported_symbols_added | source | – | – | – | 0 |
| exported_symbols_removed | source | – | – | – | 0 |
| exported_symbols_changed | source | – | – | – | 0 |
| units_changed | source | – | – | – | 2 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 3 test files changed, +116/-0 lines; test churn is not evidence of correctness)

Coverage: 3 changed code file(s) in scope, 3 analysed, 0 with parse errors, unsupported {}, excluded by class {source:4}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 3 |
| test_lines_added | context_only | – | – | – | 116 |
| test_lines_deleted | context_only | – | – | – | 0 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| modified | method | upgrader.mapText | internal/engine/workflow/upgrade.go:528-548 | internal/engine/workflow/upgrade.go:528-558 | source |
| added | function | wholeReferences | internal/engine/workflow/upgrade.go:553-569 |  | source |
| added | function | TestAStepRewriteMapsOnlyWholeReferences | internal/engine/workflow/upgrade_test.go:239-279 |  | test |
| added | function | TestHistoryConversationTakesTheNewestRunsAccount | internal/service/agent_history_account_test.go:14-53 |  | test |
| modified | method | agentService.historyConversation | internal/service/agent_resume.go:190-203 | internal/service/agent_resume.go:190-213 | source |
| added | method | agentService.latestRunAccount | internal/service/agent_resume.go:208-220 |  | source |
| modified | method | SessionHistoryService.Recover | internal/service/session_recovery.go:28-86 | internal/service/session_recovery.go:28-102 | source |
| added | method | SessionHistoryService.storageDirsFor | internal/service/session_recovery.go:91-110 |  | source |
| modified | method | TicketCreateExecutor.Start | internal/service/workflow_ticket.go:171-236 | internal/service/workflow_ticket.go:171-246 | source |
| added | method | TicketCreateExecutor.resolveCreateAssignee | internal/service/workflow_ticket.go:242-258 |  | source |
| added | function | TestTicketCreateWithAnAssigneeAWriterCannotSetCreatesNothing | internal/service/workflow_ticket_test.go:886-900 |  | test |

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
