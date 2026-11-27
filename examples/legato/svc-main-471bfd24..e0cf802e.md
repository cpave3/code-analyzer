# CodeEval report: legato-svc

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `e0cf802e15cf3c2007285cbdc43e759073b1ebbb` (e0cf802e15cf3c2007285cbdc43e759073b1ebbb) tree `fad541164613` 2026-11-27T19:53:11+11:00
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `3cef63633359769f` · tool 0.1.0-poc
- **Change:** 11 files, +615/-121 lines, 19 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | ok | +2 | structural redundancy 260→255 (Δ-5), exact 10→10; 0 near-structural candidate(s) |
| complexity | complexity | ok | +2 | modified callables net cyclomatic Δ-29; 4 added (cc 14), 0 deleted (cc 0); hotspots 0→0 |
| architecture_coupling | dependencies | ok | 0 | internal edges 200→200 (+0/-0); cycles 0→0; rule violations 15→15 |
| maintainability_scope | obsolete | ok | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | ok | 0 | 11 files changed (7 source); symbols +5/-0/~11 (renamed 0, formatting-only 0); exported API +0/-0/~0; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 3 test files changed, +156/-0 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 2210 | 10 | 389 |  |
| javascript | no | 33 | 0 | 0 | unsupported: no adapter |
| kotlin | no | 57 | 0 | 0 | unsupported: no adapter |
| python | no | 10 | 0 | 0 | unsupported: no adapter |
| rust | no | 34 | 0 | 0 | unsupported: no adapter |
| shell | no | 15 | 0 | 0 | unsupported: no adapter |
| sql | no | 160 | 0 | 0 | unsupported: no adapter |
| swift | no | 2 | 0 | 0 | unsupported: no adapter |
| typescript | no | 999 | 0 | 0 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
docs 1 (+347/-0); source 7 (+112/-121); test 3 (+156/-0); 

## duplication — Callable duplication delta (`duplication`)

**Status:** ok · **Score:** +2 · **Rule:** net structural redundancy Δ-5: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 260→255 (Δ-5), exact 10→10; 0 near-structural candidate(s)

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 1

- callables with ≥40 tokens in production source are compared repository-wide (base 6762, target 6757)
- introduced 1 redundant instance(s), removed 6; 0 pre-existing group(s) touched but unchanged in size

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 10 | 10 | 0 | – |
| structural_redundant_callables | repository | 260 | 255 | -5 | – |
| callables_considered | repository | 6762 | 6757 | -5 | – |
| redundant_instances_introduced | change | – | – | – | 1 |
| redundant_instances_removed | change | – | – | – | 6 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

Findings (6):

- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [claudeUsageError opencodeUsageError]
  - target `e0cf802e15cf` internal/service/subscription_usage_readers.go:64-70 `claudeUsageError` (new or changed (modified))
  - target `e0cf802e15cf` internal/service/subscription_usage_readers.go:207-213 `opencodeUsageError` (new or changed (modified))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.removed** [improvement; measured] clone group shrank from 3 to 0 (removed/changed: [claimChatBots claimMQTTPublisher claimUsagePoller])
  - base `471bfd24a696` cmd/legato/chatbots_runtime.go:68-82 `claimChatBots` (removed or changed (modified))
  - base `471bfd24a696` cmd/legato/mqtt_runtime.go:27-41 `claimMQTTPublisher` (removed or changed (modified))
  - base `471bfd24a696` cmd/legato/usage_runtime.go:32-46 `claimUsagePoller` (removed or changed (modified))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 0 (removed/changed: [FakeReviewGitHub.Draft FakeReviewGitHub.Push])
  - base `471bfd24a696` internal/handofftest/reviews.go:141-150 `FakeReviewGitHub.Draft` (removed or changed (modified))
  - base `471bfd24a696` internal/handofftest/reviews.go:830-839 `FakeReviewGitHub.Push` (removed or changed (modified))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 0 (removed/changed: [FakeReviewGitHub.SetAuthor FakeReviewGitHub.Describe])
  - base `471bfd24a696` internal/handofftest/reviews.go:360-369 `FakeReviewGitHub.SetAuthor` (removed or changed (modified))
  - base `471bfd24a696` internal/handofftest/reviews.go:373-382 `FakeReviewGitHub.Describe` (removed or changed (modified))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 0 (removed/changed: [PullRequestModes.PullRequestMode PullRequestModes.ReadyAfterReview])
  - base `471bfd24a696` internal/service/pull_request_modes.go:23-33 `PullRequestModes.PullRequestMode` (removed or changed (modified))
  - base `471bfd24a696` internal/service/pull_request_modes.go:45-55 `PullRequestModes.ReadyAfterReview` (removed or changed (modified))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 0 (removed/changed: [claudeUsageError opencodeUsageError])
  - base `471bfd24a696` internal/service/subscription_usage_readers.go:52-73 `claudeUsageError` (removed or changed (modified))
  - base `471bfd24a696` internal/service/subscription_usage_readers.go:176-197 `opencodeUsageError` (removed or changed (modified))

## complexity — Callable complexity deltas (`complexity`)

**Status:** ok · **Score:** +2 · **Rule:** (pre-existing hotspotΔ ≤ -2 or modified net ≤ -10) and changed-scope total decreased → +2

modified callables net cyclomatic Δ-29; 4 added (cc 14), 0 deleted (cc 0); hotspots 0→0

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 1

- modified callables: 0 increased, 11 decreased; net cyclomatic Δ-29
- added callables: 4 (cyclomatic 14); deleted callables: 0 (cyclomatic 0)
- changed-scope cyclomatic 46→31 (Δ-15), hotspots 0→0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 46 | 17 | -29 | – |
| decision_points_sum | modified_callables | 35 | 6 | -29 | – |
| lines_sum | modified_callables | 151 | 46 | -105 | – |
| max_nesting | modified_callables | 2 | 1 | -1 | – |
| cyclomatic_sum | changed_scope | 46 | 31 | -15 | – |
| max_cyclomatic | changed_scope | 8 | 4 | -4 | – |
| hotspot_callables | changed_scope | 0 | 0 | 0 | – |
| callables | changed_scope | 11 | 15 | +4 | – |
| callables_modified |  | – | – | – | 11 |
| callables_added |  | – | – | – | 4 |
| callables_deleted |  | – | – | – | 0 |
| cyclomatic_sum_added | added_callables | – | – | – | 14 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 0 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 0 |
| hotspots_added | added_callables | – | – | – | 0 |

Findings (11):

- **complexity.modified** [improvement; measured] claimChatBots: cyclomatic 4→1, nesting 2→0, lines 15→3
  - base `471bfd24a696` cmd/legato/chatbots_runtime.go:68-82 `claimChatBots` (before)
  - target `e0cf802e15cf` cmd/legato/chatbots_runtime.go:65-67 `claimChatBots` (after)
- **complexity.modified** [improvement; measured] claimMQTTPublisher: cyclomatic 4→1, nesting 2→0, lines 15→3
  - base `471bfd24a696` cmd/legato/mqtt_runtime.go:27-41 `claimMQTTPublisher` (before)
  - target `e0cf802e15cf` cmd/legato/mqtt_runtime.go:25-27 `claimMQTTPublisher` (after)
- **complexity.modified** [improvement; measured] claimUsagePoller: cyclomatic 4→1, nesting 2→0, lines 15→3
  - base `471bfd24a696` cmd/legato/usage_runtime.go:32-46 `claimUsagePoller` (before)
  - target `e0cf802e15cf` cmd/legato/usage_runtime.go:30-32 `claimUsagePoller` (after)
- **complexity.modified** [improvement; measured] FakeReviewGitHub.Draft: cyclomatic 3→1, nesting 2→1, lines 10→3
  - base `471bfd24a696` internal/handofftest/reviews.go:141-150 `FakeReviewGitHub.Draft` (before)
  - target `e0cf802e15cf` internal/handofftest/reviews.go:153-155 `FakeReviewGitHub.Draft` (after)
- **complexity.modified** [improvement; measured] FakeReviewGitHub.SetAuthor: cyclomatic 3→1, nesting 2→1, lines 10→3
  - base `471bfd24a696` internal/handofftest/reviews.go:360-369 `FakeReviewGitHub.SetAuthor` (before)
  - target `e0cf802e15cf` internal/handofftest/reviews.go:365-367 `FakeReviewGitHub.SetAuthor` (after)
- **complexity.modified** [improvement; measured] FakeReviewGitHub.Describe: cyclomatic 3→1, nesting 2→1, lines 10→3
  - base `471bfd24a696` internal/handofftest/reviews.go:373-382 `FakeReviewGitHub.Describe` (before)
  - target `e0cf802e15cf` internal/handofftest/reviews.go:371-373 `FakeReviewGitHub.Describe` (after)
- **complexity.modified** [improvement; measured] FakeReviewGitHub.Push: cyclomatic 3→1, nesting 2→1, lines 10→3
  - base `471bfd24a696` internal/handofftest/reviews.go:830-839 `FakeReviewGitHub.Push` (before)
  - target `e0cf802e15cf` internal/handofftest/reviews.go:821-823 `FakeReviewGitHub.Push` (after)
- **complexity.modified** [improvement; measured] PullRequestModes.PullRequestMode: cyclomatic 3→2, nesting 1→1, lines 11→7
  - base `471bfd24a696` internal/service/pull_request_modes.go:23-33 `PullRequestModes.PullRequestMode` (before)
  - target `e0cf802e15cf` internal/service/pull_request_modes.go:23-29 `PullRequestModes.PullRequestMode` (after)
- **complexity.modified** [improvement; measured] PullRequestModes.ReadyAfterReview: cyclomatic 3→2, nesting 1→0, lines 11→4
  - base `471bfd24a696` internal/service/pull_request_modes.go:45-55 `PullRequestModes.ReadyAfterReview` (before)
  - target `e0cf802e15cf` internal/service/pull_request_modes.go:56-59 `PullRequestModes.ReadyAfterReview` (after)
- **complexity.modified** [improvement; measured] claudeUsageError: cyclomatic 8→3, nesting 1→1, lines 22→7
  - base `471bfd24a696` internal/service/subscription_usage_readers.go:52-73 `claudeUsageError` (before)
  - target `e0cf802e15cf` internal/service/subscription_usage_readers.go:64-70 `claudeUsageError` (after)
- **complexity.modified** [improvement; measured] opencodeUsageError: cyclomatic 8→3, nesting 1→1, lines 22→7
  - base `471bfd24a696` internal/service/subscription_usage_readers.go:176-197 `opencodeUsageError` (before)
  - target `e0cf802e15cf` internal/service/subscription_usage_readers.go:207-213 `opencodeUsageError` (after)

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** ok · **Score:** 0 · **Rule:** no cycles or rule violations changed → 0 (new edges alone are reported, not judged: features legitimately add dependencies)

internal edges 200→200 (+0/-0); cycles 0→0; rule violations 15→15

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 1

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
  - target `e0cf802e15cf` config/config.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/events (not scored)
  - target `e0cf802e15cf` internal/tui/app.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/ipc (not scored)
  - target `e0cf802e15cf` internal/tui/app.go:15-15 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/macros (not scored)
  - target `e0cf802e15cf` internal/tui/app.go:16-16 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/store (not scored)
  - target `e0cf802e15cf` internal/tui/app.go:17-17 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/workflow (not scored)
  - target `e0cf802e15cf` internal/tui/app.go:18-18 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/board → internal/engine/store (not scored)
  - target `e0cf802e15cf` internal/tui/board/model.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/github (not scored)
  - target `e0cf802e15cf` internal/tui/overlay/link_pr.go:10-10 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/macros (not scored)
  - target `e0cf802e15cf` internal/tui/overlay/macro_picker.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/workflow (not scored)
  - target `e0cf802e15cf` internal/tui/overlay/workflow_trigger.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/store (not scored)
  - target `e0cf802e15cf` internal/tui/overlay/workspace.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/plan → internal/engine/store (not scored)
  - target `e0cf802e15cf` internal/tui/plan/model.go:13-13 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/report → internal/engine/analytics (not scored)
  - target `e0cf802e15cf` internal/tui/report/model.go:12-12 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/git (not scored)
  - target `e0cf802e15cf` internal/tui/review/messages.go:4-4 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/store (not scored)
  - target `e0cf802e15cf` internal/tui/review/model.go:17-17 (import)

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** ok · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 1

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** ok · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

11 files changed (7 source); symbols +5/-0/~11 (renamed 0, formatting-only 0); exported API +0/-0/~0; 1 unit group(s)

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:3}, non-code 1

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | docs | – | – | – | 1 |
| lines_added | docs | – | – | – | 347 |
| lines_deleted | docs | – | – | – | 0 |
| files_changed | source | – | – | – | 7 |
| lines_added | source | – | – | – | 112 |
| lines_deleted | source | – | – | – | 121 |
| files_changed | test | – | – | – | 3 |
| lines_added | test | – | – | – | 156 |
| lines_deleted | test | – | – | – | 0 |
| symbols_added | source | – | – | – | 5 |
| symbols_deleted | source | – | – | – | 0 |
| symbols_modified | source | – | – | – | 11 |
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

not assessed (context: 3 test files changed, +156/-0 lines; test churn is not evidence of correctness)

Coverage: 3 changed code file(s) in scope, 3 analysed, 0 with parse errors, unsupported {}, excluded by class {source:7}, non-code 1

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 3 |
| test_lines_added | context_only | – | – | – | 156 |
| test_lines_deleted | context_only | – | – | – | 0 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| modified | function | claimChatBots | cmd/legato/chatbots_runtime.go:65-67 | cmd/legato/chatbots_runtime.go:68-82 | source |
| modified | function | claimMQTTPublisher | cmd/legato/mqtt_runtime.go:25-27 | cmd/legato/mqtt_runtime.go:27-41 | source |
| added | function | claimDataDirLock | cmd/legato/runtime_lock.go:16-30 |  | source |
| added | function | TestRuntimeRolesLockTheirOwnFiles | cmd/legato/runtime_lock_test.go:13-45 |  | test |
| modified | function | claimUsagePoller | cmd/legato/usage_runtime.go:30-32 | cmd/legato/usage_runtime.go:32-46 | source |
| added | method | FakeReviewGitHub.editPR | internal/handofftest/reviews.go:141-150 |  | source |
| modified | method | FakeReviewGitHub.Draft | internal/handofftest/reviews.go:153-155 | internal/handofftest/reviews.go:141-150 | source |
| modified | method | FakeReviewGitHub.SetAuthor | internal/handofftest/reviews.go:365-367 | internal/handofftest/reviews.go:360-369 | source |
| modified | method | FakeReviewGitHub.Describe | internal/handofftest/reviews.go:371-373 | internal/handofftest/reviews.go:373-382 | source |
| modified | method | FakeReviewGitHub.Push | internal/handofftest/reviews.go:821-823 | internal/handofftest/reviews.go:830-839 | source |
| modified | method | PullRequestModes.PullRequestMode | internal/service/pull_request_modes.go:23-29 | internal/service/pull_request_modes.go:23-33 | source |
| added | method | PullRequestModes.rules | internal/service/pull_request_modes.go:34-44 |  | source |
| modified | method | PullRequestModes.ReadyAfterReview | internal/service/pull_request_modes.go:56-59 | internal/service/pull_request_modes.go:45-55 | source |
| added | function | TestPullRequestModesAnswerFromTheRulesOrNothing | internal/service/pull_request_modes_test.go:105-140 |  | test |
| added | function | TestProviderUsageErrorsGiveReasons | internal/service/subscription_usage_errors_test.go:17-72 |  | test |
| modified | function | claudeUsageError | internal/service/subscription_usage_readers.go:64-70 | internal/service/subscription_usage_readers.go:52-73 | source |
| added | type | usageErrorReason | internal/service/subscription_usage_readers.go:74-77 |  | source |
| added | function | usageError | internal/service/subscription_usage_readers.go:82-92 |  | source |
| modified | function | opencodeUsageError | internal/service/subscription_usage_readers.go:207-213 | internal/service/subscription_usage_readers.go:176-197 | source |

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
