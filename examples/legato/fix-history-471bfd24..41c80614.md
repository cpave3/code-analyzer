# CodeEval history: legato-fix

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) 2026-11-27T15:46:34+11:00
- **Target:** `41c8061476b6592b5ea4faa192d68a2a803b6c59` (41c8061476b6592b5ea4faa192d68a2a803b6c59) 2026-11-27T17:33:03+11:00
- **Run:** started 2026-11-27T17:35:18+11:00, 1.4s total, 0 files parsed, cache 8834 hits / 0 misses

- each step compares a first-parent commit with its first parent; for merge commits this is the whole merged branch, counted once
- step results are observed contributions to metrics, not proof that a workflow or person caused a quality change
- step reports do not sum to the net report: repository-wide measurements (duplication, graph) interact across steps

## First-parent steps

| # | Commit | Subject (untrusted) | Files | +/- | Go files | duplication | complexity | dependencies | obsolete | scope | correctness |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | `6f4e46168009` | refactor(server): share the invalidation broadcast loop (ah7qmrsa) | 4 | +91/-34 | 4 | +1 | +1 | 0 | 0 | 0 | n/a |
| 2 | `b8ebb812d635` | refactor(jira): one doJSON for POST and PUT writes (ah7qmrsa) | 5 | +63/-19 | 5 | +1 | 0 | 0 | 0 | 0 | n/a |
| 3 | `41c8061476b6` | refactor(service): plan images use pathWithin (ah7qmrsa) | 2 | +16/-6 | 2 | 0 | 0 | 0 | 0 | 0 | n/a |
| net | `41c8061476b6` | base → target (net) | 11 | +170/-59 | 11 | +1 | +1 | 0 | 0 | 0 | n/a |

## Key measurements per step

| # | Commit | dup Δ (struct redundancy) | modified cc Δ | added callables (cc) | deleted callables (cc) | hotspots before→after | edges +/- | cycles Δ | obsolete kept | API +/-/~ |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 1 | `6f4e46168009` | -1 | -6 | 1 (4) | 0 (0) | 0→0 | +0/-0 | 0 | 0 | +0/-0/~0 |
| 2 | `b8ebb812d635` | -1 | 0 | 1 (2) | 2 (4) | 1→1 | +0/-0 | 0 | 0 | +0/-0/~0 |
| 3 | `41c8061476b6` | 0 | -2 | 0 (0) | 0 (0) | 0→0 | +0/-0 | 0 | 0 | +0/-0/~0 |
| net | `41c8061476b6` | -2 | -8 | 2 (6) | 2 (4) | 1→1 | +0/-0 | 0 | 0 | +0/-0/~0 |


---

# CodeEval report: legato-fix

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `41c8061476b6592b5ea4faa192d68a2a803b6c59` (41c8061476b6592b5ea4faa192d68a2a803b6c59) tree `555c4b3e8792` 2026-11-27T17:33:03+11:00
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `3cef63633359769f` · tool 0.1.0-poc
- **Change:** 11 files, +170/-59 lines, 18 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | ok | +1 | structural redundancy 260→258 (Δ-2), exact 10→10; 0 near-structural candidate(s) |
| complexity | complexity | ok | +1 | modified callables net cyclomatic Δ-8; 2 added (cc 6), 2 deleted (cc 4); hotspots 1→1 |
| architecture_coupling | dependencies | ok | 0 | internal edges 200→200 (+0/-0); cycles 0→0; rule violations 15→15 |
| maintainability_scope | obsolete | ok | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | ok | 0 | 11 files changed (7 source); symbols +2/-2/~7 (renamed 0, formatting-only 0); exported API +0/-0/~0; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 4 test files changed, +129/-3 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 2209 | 11 | 229 |  |
| javascript | no | 33 | 0 | 0 | unsupported: no adapter |
| kotlin | no | 57 | 0 | 0 | unsupported: no adapter |
| python | no | 10 | 0 | 0 | unsupported: no adapter |
| rust | no | 34 | 0 | 0 | unsupported: no adapter |
| shell | no | 15 | 0 | 0 | unsupported: no adapter |
| sql | no | 160 | 0 | 0 | unsupported: no adapter |
| swift | no | 2 | 0 | 0 | unsupported: no adapter |
| typescript | no | 999 | 0 | 0 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
source 7 (+41/-56); test 4 (+129/-3); 

## duplication — Callable duplication delta (`duplication`)

**Status:** ok · **Score:** +1 · **Rule:** net structural redundancy Δ-2: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 260→258 (Δ-2), exact 10→10; 0 near-structural candidate(s)

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:4}, non-code 0

- callables with ≥40 tokens in production source are compared repository-wide (base 6762, target 6760)
- introduced 0 redundant instance(s), removed 2; 0 pre-existing group(s) touched but unchanged in size

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 10 | 10 | 0 | – |
| structural_redundant_callables | repository | 260 | 258 | -2 | – |
| callables_considered | repository | 6762 | 6760 | -2 | – |
| redundant_instances_introduced | change | – | – | – | 0 |
| redundant_instances_removed | change | – | – | – | 2 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

Findings (2):

- **duplication.removed** [improvement; measured] clone group shrank from 2 to 0 (removed/changed: [Client.doPost Client.doPut])
  - base `471bfd24a696` internal/engine/jira/client.go:157-164 `Client.doPost` (removed or changed (deleted))
  - base `471bfd24a696` internal/engine/jira/users.go:57-64 `Client.doPut` (removed or changed (deleted))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 0 (removed/changed: [Server.StartSubscriptionUsageEvents Server.StartServerThemeEvents])
  - base `471bfd24a696` internal/server/subscription_usage.go:70-87 `Server.StartSubscriptionUsageEvents` (removed or changed (modified))
  - base `471bfd24a696` internal/server/theme_preference.go:23-40 `Server.StartServerThemeEvents` (removed or changed (modified))

## complexity — Callable complexity deltas (`complexity`)

**Status:** ok · **Score:** +1 · **Rule:** pre-existing hotspotΔ&lt;0 or modified net cyclomatic ≤ -3 → +1

modified callables net cyclomatic Δ-8; 2 added (cc 6), 2 deleted (cc 4); hotspots 1→1

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:4}, non-code 0

- modified callables: 0 increased, 3 decreased; net cyclomatic Δ-8
- added callables: 2 (cyclomatic 6); deleted callables: 2 (cyclomatic 4)
- changed-scope cyclomatic 41→35 (Δ-6), hotspots 1→1

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 37 | 29 | -8 | – |
| decision_points_sum | modified_callables | 30 | 22 | -8 | – |
| lines_sum | modified_callables | 138 | 104 | -34 | – |
| max_nesting | modified_callables | 3 | 2 | -1 | – |
| cyclomatic_sum | changed_scope | 41 | 35 | -6 | – |
| max_cyclomatic | changed_scope | 10 | 10 | 0 | – |
| hotspot_callables | changed_scope | 1 | 1 | 0 | – |
| callables | changed_scope | 9 | 9 | 0 | – |
| callables_modified |  | – | – | – | 7 |
| callables_added |  | – | – | – | 2 |
| callables_deleted |  | – | – | – | 2 |
| cyclomatic_sum_added | added_callables | – | – | – | 6 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 4 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 0 |
| hotspots_added | added_callables | – | – | – | 0 |

Findings (3):

- **complexity.modified** [improvement; measured] Server.StartSubscriptionUsageEvents: cyclomatic 4→1, nesting 3→0, lines 18→3
  - base `471bfd24a696` internal/server/subscription_usage.go:70-87 `Server.StartSubscriptionUsageEvents` (before)
  - target `41c8061476b6` internal/server/subscription_usage.go:69-71 `Server.StartSubscriptionUsageEvents` (after)
- **complexity.modified** [improvement; measured] Server.StartServerThemeEvents: cyclomatic 4→1, nesting 3→0, lines 18→3
  - base `471bfd24a696` internal/server/theme_preference.go:23-40 `Server.StartServerThemeEvents` (before)
  - target `41c8061476b6` internal/server/theme_preference.go:22-24 `Server.StartServerThemeEvents` (after)
- **complexity.modified** [improvement; measured] planBundleImage: cyclomatic 8→6, nesting 2→2, lines 26→22
  - base `471bfd24a696` internal/service/plan.go:72-97 `planBundleImage` (before)
  - target `41c8061476b6` internal/service/plan.go:72-93 `planBundleImage` (after)

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** ok · **Score:** 0 · **Rule:** no cycles or rule violations changed → 0 (new edges alone are reported, not judged: features legitimately add dependencies)

internal edges 200→200 (+0/-0); cycles 0→0; rule violations 15→15

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:4}, non-code 0

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
  - target `41c8061476b6` config/config.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/events (not scored)
  - target `41c8061476b6` internal/tui/app.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/ipc (not scored)
  - target `41c8061476b6` internal/tui/app.go:15-15 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/macros (not scored)
  - target `41c8061476b6` internal/tui/app.go:16-16 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/store (not scored)
  - target `41c8061476b6` internal/tui/app.go:17-17 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/workflow (not scored)
  - target `41c8061476b6` internal/tui/app.go:18-18 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/board → internal/engine/store (not scored)
  - target `41c8061476b6` internal/tui/board/model.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/github (not scored)
  - target `41c8061476b6` internal/tui/overlay/link_pr.go:10-10 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/macros (not scored)
  - target `41c8061476b6` internal/tui/overlay/macro_picker.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/workflow (not scored)
  - target `41c8061476b6` internal/tui/overlay/workflow_trigger.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/store (not scored)
  - target `41c8061476b6` internal/tui/overlay/workspace.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/plan → internal/engine/store (not scored)
  - target `41c8061476b6` internal/tui/plan/model.go:13-13 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/report → internal/engine/analytics (not scored)
  - target `41c8061476b6` internal/tui/report/model.go:12-12 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/git (not scored)
  - target `41c8061476b6` internal/tui/review/messages.go:4-4 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/store (not scored)
  - target `41c8061476b6` internal/tui/review/model.go:17-17 (import)

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** ok · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:4}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** ok · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

11 files changed (7 source); symbols +2/-2/~7 (renamed 0, formatting-only 0); exported API +0/-0/~0; 1 unit group(s)

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:4}, non-code 0

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | source | – | – | – | 7 |
| lines_added | source | – | – | – | 41 |
| lines_deleted | source | – | – | – | 56 |
| files_changed | test | – | – | – | 4 |
| lines_added | test | – | – | – | 129 |
| lines_deleted | test | – | – | – | 3 |
| symbols_added | source | – | – | – | 2 |
| symbols_deleted | source | – | – | – | 2 |
| symbols_modified | source | – | – | – | 7 |
| symbols_renamed | source | – | – | – | 0 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 0 |
| exported_symbols_added | source | – | – | – | 0 |
| exported_symbols_removed | source | – | – | – | 0 |
| exported_symbols_changed | source | – | – | – | 0 |
| units_changed | source | – | – | – | 3 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 4 test files changed, +129/-3 lines; test churn is not evidence of correctness)

Coverage: 4 changed code file(s) in scope, 4 analysed, 0 with parse errors, unsupported {}, excluded by class {source:7}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 4 |
| test_lines_added | context_only | – | – | – | 129 |
| test_lines_deleted | context_only | – | – | – | 3 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| modified | method | Client.DoTransition | internal/engine/jira/client.go:140-147 | internal/engine/jira/client.go:140-147 | source |
| deleted | method | Client.doPost → Client.doPost |  | internal/engine/jira/client.go:157-164 | source |
| added | method | Client.doJSON | internal/engine/jira/client.go:158-165 |  | source |
| modified | method | Client.Assign | internal/engine/jira/users.go:41-54 | internal/engine/jira/users.go:42-55 | source |
| deleted | method | Client.doPut → Client.doPut |  | internal/engine/jira/users.go:57-64 | source |
| formatting_only | function | TestCreateIssueAssignsAsItCreates | internal/engine/jira/users_test.go:94-103 | internal/engine/jira/users_test.go:91-100 | test |
| added | function | TestWritesSendJSONWithBasicAuth | internal/engine/jira/users_test.go:106-127 |  | test |
| added | function | TestAssignRateLimitedIsRetriedWithItsBody | internal/engine/jira/users_test.go:130-150 |  | test |
| modified | method | Client.CreateIssue | internal/engine/jira/write.go:93-129 | internal/engine/jira/write.go:92-128 | source |
| modified | method | Client.AddComment | internal/engine/jira/write.go:132-148 | internal/engine/jira/write.go:131-147 | source |
| modified | type | recordedRequest | internal/engine/jira/write_test.go:15-19 | internal/engine/jira/write_test.go:15-18 | test |
| modified | function | jiraServer | internal/engine/jira/write_test.go:21-36 | internal/engine/jira/write_test.go:20-34 | test |
| added | method | Server.broadcastInvalidations | internal/server/invalidation_events.go:12-29 |  | source |
| added | function | TestInvalidationBroadcastersStopCleanly | internal/server/invalidation_events_test.go:12-60 |  | test |
| modified | method | Server.StartSubscriptionUsageEvents | internal/server/subscription_usage.go:69-71 | internal/server/subscription_usage.go:70-87 | source |
| modified | method | Server.StartServerThemeEvents | internal/server/theme_preference.go:22-24 | internal/server/theme_preference.go:23-40 | source |
| modified | function | planBundleImage | internal/service/plan.go:72-93 | internal/service/plan.go:72-97 | source |
| added | function | TestPlanImageInASiblingFolderSharingTheBundlesPrefixIsOutside | internal/service/plan_images_test.go:418-428 |  | test |

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
