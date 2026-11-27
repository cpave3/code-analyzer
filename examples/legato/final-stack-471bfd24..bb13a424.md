# CodeEval report: legato-stack

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `bb13a424dc73f3eeb98f8bce43db932a9d19f7e6` (bb13a424dc73f3eeb98f8bce43db932a9d19f7e6) tree `25b725950d6a` 2026-11-27T20:29:04+11:00
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `3cef63633359769f` · tool 0.1.0-poc
- **Change:** 83 files, +4497/-774 lines, 260 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | ok | +2 | structural redundancy 260→253 (Δ-7), exact 10→8; 0 near-structural candidate(s) |
| complexity | complexity | ok | 0 | modified callables net cyclomatic Δ-211; 43 added (cc 245), 5 deleted (cc 20); hotspots 29→26 |
| architecture_coupling | dependencies | ok | +1 | internal edges 200→199 (+5/-6); cycles 0→0; rule violations 15→11 |
| maintainability_scope | obsolete | ok | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | ok | 0 | 83 files changed (37 source); symbols +50/-6/~63 (renamed 7, formatting-only 3); exported API +6/-3/~12; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 38 test files changed, +2988/-79 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 2226 | 75 | 4765 |  |
| javascript | no | 33 | 0 | 0 | unsupported: no adapter |
| kotlin | no | 57 | 0 | 0 | unsupported: no adapter |
| python | no | 10 | 0 | 0 | unsupported: no adapter |
| rust | no | 34 | 0 | 0 | unsupported: no adapter |
| shell | no | 15 | 0 | 0 | unsupported: no adapter |
| sql | no | 160 | 0 | 0 | unsupported: no adapter |
| swift | no | 2 | 0 | 0 | unsupported: no adapter |
| typescript | no | 999 | 0 | 0 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
docs 8 (+495/-11); source 37 (+1014/-684); test 38 (+2988/-79); 

## duplication — Callable duplication delta (`duplication`)

**Status:** ok · **Score:** +2 · **Rule:** net structural redundancy Δ-7: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 260→253 (Δ-7), exact 10→8; 0 near-structural candidate(s)

Coverage: 37 changed code file(s) in scope, 37 analysed, 0 with parse errors, unsupported {}, excluded by class {test:38}, non-code 8

- callables with ≥40 tokens in production source are compared repository-wide (base 6762, target 6787)
- introduced 1 redundant instance(s), removed 8; 0 pre-existing group(s) touched but unchanged in size

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 10 | 8 | -2 | – |
| structural_redundant_callables | repository | 260 | 253 | -7 | – |
| callables_considered | repository | 6762 | 6787 | +25 | – |
| redundant_instances_introduced | change | – | – | – | 1 |
| redundant_instances_removed | change | – | – | – | 8 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

Findings (8):

- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [claudeUsageError opencodeUsageError]
  - target `bb13a424dc73` internal/service/subscription_usage_readers.go:64-70 `claudeUsageError` (new or changed (modified))
  - target `bb13a424dc73` internal/service/subscription_usage_readers.go:207-213 `opencodeUsageError` (new or changed (modified))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.removed** [improvement; measured] clone group shrank from 3 to 0 (removed/changed: [claimChatBots claimMQTTPublisher claimUsagePoller])
  - base `471bfd24a696` cmd/legato/chatbots_runtime.go:68-82 `claimChatBots` (removed or changed (modified))
  - base `471bfd24a696` cmd/legato/mqtt_runtime.go:27-41 `claimMQTTPublisher` (removed or changed (modified))
  - base `471bfd24a696` cmd/legato/usage_runtime.go:32-46 `claimUsagePoller` (removed or changed (modified))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 1 (removed/changed: [retryAfter])
  - base `471bfd24a696` internal/engine/claudeusage/claudeusage.go:197-209 `retryAfter` (remaining copy)
  - base `471bfd24a696` internal/engine/opencodeusage/opencodeusage.go:185-197 `retryAfter` (removed or changed (deleted))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 1 (removed/changed: [openRegular openRegular])
  - base `471bfd24a696` internal/engine/claudeusage/open.go:17-28 `openRegular` (removed or changed (deleted))
  - base `471bfd24a696` internal/engine/codexusage/open.go:17-28 `openRegular` (removed or changed (deleted))
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

**Status:** ok · **Score:** 0 · **Rule:** mixed: degradation and improvement signals both present → 0

modified callables net cyclomatic Δ-211; 43 added (cc 245), 5 deleted (cc 20); hotspots 29→26

Coverage: 37 changed code file(s) in scope, 37 analysed, 0 with parse errors, unsupported {}, excluded by class {test:38}, non-code 8

- modified callables: 0 increased, 33 decreased; net cyclomatic Δ-211
- added callables: 43 (cyclomatic 245); deleted callables: 5 (cyclomatic 20)
- changed-scope cyclomatic 816→830 (Δ+14), hotspots 29→26
- complexity removed from modified callables (Δ-211) is matched by 43 added callables (cyclomatic 245): this is extraction/redistribution, not net removal

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 796 | 585 | -211 | – |
| decision_points_sum | modified_callables | 741 | 530 | -211 | – |
| lines_sum | modified_callables | 2443 | 1920 | -523 | – |
| max_nesting | modified_callables | 4 | 4 | 0 | – |
| cyclomatic_sum | changed_scope | 816 | 830 | +14 | – |
| max_cyclomatic | changed_scope | 107 | 96 | -11 | – |
| hotspot_callables | changed_scope | 29 | 26 | -3 | – |
| callables | changed_scope | 60 | 98 | +38 | – |
| callables_modified |  | – | – | – | 55 |
| callables_added |  | – | – | – | 43 |
| callables_deleted |  | – | – | – | 5 |
| cyclomatic_sum_added | added_callables | – | – | – | 245 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 20 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | -8 |
| hotspots_added | added_callables | – | – | – | 5 |

Findings (40):

- **complexity.added_hotspot** [degradation; measured] new callable agentService.checkNamedAccountIsolation has cyclomatic 22, nesting 2 (hotspot thresholds 10/4)
  - target `bb13a424dc73` internal/service/agent_accounts.go:156-215 `agentService.checkNamedAccountIsolation` (added)
- **complexity.added_hotspot** [degradation; measured] new callable agentService.resumeConversation has cyclomatic 12, nesting 2 (hotspot thresholds 10/4)
  - target `bb13a424dc73` internal/service/agent_resume.go:180-206 `agentService.resumeConversation` (added)
- **complexity.added_hotspot** [degradation; measured] new callable withClaudeAccounts has cyclomatic 10, nesting 2 (hotspot thresholds 10/4)
  - target `bb13a424dc73` internal/service/agent_sandbox.go:443-470 `withClaudeAccounts` (added)
- **complexity.added_hotspot** [degradation; measured] new callable reviewRoundDone has cyclomatic 13, nesting 2 (hotspot thresholds 10/4)
  - target `bb13a424dc73` internal/service/review_pass_run.go:1233-1249 `reviewRoundDone` (added)
- **complexity.added_hotspot** [degradation; measured] new callable ReviewPassService.readWatchRound has cyclomatic 15, nesting 2 (hotspot thresholds 10/4)
  - target `bb13a424dc73` internal/service/review_watch.go:392-421 `ReviewPassService.readWatchRound` (added)
- **complexity.modified** [improvement; measured] claimChatBots: cyclomatic 4→1, nesting 2→0, lines 15→3
  - base `471bfd24a696` cmd/legato/chatbots_runtime.go:68-82 `claimChatBots` (before)
  - target `bb13a424dc73` cmd/legato/chatbots_runtime.go:65-67 `claimChatBots` (after)
- **complexity.modified** [improvement; measured] claimMQTTPublisher: cyclomatic 4→1, nesting 2→0, lines 15→3
  - base `471bfd24a696` cmd/legato/mqtt_runtime.go:27-41 `claimMQTTPublisher` (before)
  - target `bb13a424dc73` cmd/legato/mqtt_runtime.go:25-27 `claimMQTTPublisher` (after)
- **complexity.modified** [improvement; measured] claimUsagePoller: cyclomatic 4→1, nesting 2→0, lines 15→3
  - base `471bfd24a696` cmd/legato/usage_runtime.go:32-46 `claimUsagePoller` (before)
  - target `bb13a424dc73` cmd/legato/usage_runtime.go:30-32 `claimUsagePoller` (after)
- **complexity.modified** [improvement; measured] checkReviewLoopStep: cyclomatic 29→22, nesting 4→3, lines 61→47
  - base `471bfd24a696` internal/engine/workflow/review_loop.go:149-209 `checkReviewLoopStep` (before)
  - target `bb13a424dc73` internal/engine/workflow/review_loop.go:149-195 `checkReviewLoopStep` (after)
- **complexity.modified** [improvement; measured] upgrader.mapText: cyclomatic 13→8, nesting 4→3, lines 31→21
  - base `471bfd24a696` internal/engine/workflow/upgrade.go:528-558 `upgrader.mapText` (before)
  - target `bb13a424dc73` internal/engine/workflow/upgrade.go:528-548 `upgrader.mapText` (after)
- **complexity.modified** [improvement; measured] FakeReviewGitHub.Draft: cyclomatic 3→1, nesting 2→1, lines 10→3
  - base `471bfd24a696` internal/handofftest/reviews.go:141-150 `FakeReviewGitHub.Draft` (before)
  - target `bb13a424dc73` internal/handofftest/reviews.go:153-155 `FakeReviewGitHub.Draft` (after)
- **complexity.modified** [improvement; measured] FakeReviewGitHub.SetAuthor: cyclomatic 3→1, nesting 2→1, lines 10→3
  - base `471bfd24a696` internal/handofftest/reviews.go:360-369 `FakeReviewGitHub.SetAuthor` (before)
  - target `bb13a424dc73` internal/handofftest/reviews.go:365-367 `FakeReviewGitHub.SetAuthor` (after)
- **complexity.modified** [improvement; measured] FakeReviewGitHub.Describe: cyclomatic 3→1, nesting 2→1, lines 10→3
  - base `471bfd24a696` internal/handofftest/reviews.go:373-382 `FakeReviewGitHub.Describe` (before)
  - target `bb13a424dc73` internal/handofftest/reviews.go:371-373 `FakeReviewGitHub.Describe` (after)
- **complexity.modified** [improvement; measured] FakeReviewGitHub.Push: cyclomatic 3→1, nesting 2→1, lines 10→3
  - base `471bfd24a696` internal/handofftest/reviews.go:830-839 `FakeReviewGitHub.Push` (before)
  - target `bb13a424dc73` internal/handofftest/reviews.go:821-823 `FakeReviewGitHub.Push` (after)
- **complexity.modified** [improvement; measured] agentService.resolveClaudeAccount: cyclomatic 55→28, nesting 3→3, lines 131→68
  - base `471bfd24a696` internal/service/agent_accounts.go:67-197 `agentService.resolveClaudeAccount` (before)
  - target `bb13a424dc73` internal/service/agent_accounts.go:67-134 `agentService.resolveClaudeAccount` (after)
- **complexity.modified** [improvement; measured] agentService.prepareTaskResume: cyclomatic 48→28, nesting 3→2, lines 127→81
  - base `471bfd24a696` internal/service/agent_resume.go:60-186 `agentService.prepareTaskResume` (before)
  - target `bb13a424dc73` internal/service/agent_resume.go:60-140 `agentService.prepareTaskResume` (after)
- **complexity.modified** [improvement; measured] agentService.historyConversation: cyclomatic 9→6, nesting 2→1, lines 24→14
  - base `471bfd24a696` internal/service/agent_resume.go:190-213 `agentService.historyConversation` (before)
  - target `bb13a424dc73` internal/service/agent_resume.go:233-246 `agentService.historyConversation` (after)
- **complexity.modified** [improvement; measured] agentService.sandboxPrefix: cyclomatic 67→58, nesting 3→3, lines 216→192
  - base `471bfd24a696` internal/service/agent_sandbox.go:48-263 `agentService.sandboxPrefix` (before)
  - target `bb13a424dc73` internal/service/agent_sandbox.go:49-240 `agentService.sandboxPrefix` (after)
- **complexity.modified** [improvement; measured] PullRequestModes.PullRequestMode: cyclomatic 3→2, nesting 1→1, lines 11→7
  - base `471bfd24a696` internal/service/pull_request_modes.go:23-33 `PullRequestModes.PullRequestMode` (before)
  - target `bb13a424dc73` internal/service/pull_request_modes.go:23-29 `PullRequestModes.PullRequestMode` (after)
- **complexity.modified** [improvement; measured] PullRequestModes.ReadyAfterReview: cyclomatic 3→2, nesting 1→0, lines 11→4
  - base `471bfd24a696` internal/service/pull_request_modes.go:45-55 `PullRequestModes.ReadyAfterReview` (before)
  - target `bb13a424dc73` internal/service/pull_request_modes.go:56-59 `PullRequestModes.ReadyAfterReview` (after)
- **complexity.modified** [improvement; measured] OpenReviewLoopClones: cyclomatic 22→20, nesting 3→3, lines 72→71
  - base `471bfd24a696` internal/service/review_loop_open.go:20-91 `OpenReviewLoopClones` (before)
  - target `bb13a424dc73` internal/service/review_loop_open.go:20-90 `OpenReviewLoopClones` (after)
- **complexity.modified** [improvement; measured] ReviewRunCaller.openSteps: cyclomatic 10→6, nesting 3→2, lines 29→25
  - base `471bfd24a696` internal/service/review_loop_steps.go:502-530 `ReviewRunCaller.openSteps` (before)
  - target `bb13a424dc73` internal/service/review_loop_steps.go:502-526 `ReviewRunCaller.openSteps` (after)
- **complexity.modified** [improvement; measured] ReviewRunCaller.SettleSteps: cyclomatic 12→8, nesting 3→2, lines 35→31
  - base `471bfd24a696` internal/service/review_loop_steps.go:622-656 `ReviewRunCaller.SettleSteps` (before)
  - target `bb13a424dc73` internal/service/review_loop_steps.go:635-665 `ReviewRunCaller.SettleSteps` (after)
- **complexity.modified** [improvement; measured] ReviewRunCaller.closeOutputs: cyclomatic 20→15, nesting 3→3, lines 58→44
  - base `471bfd24a696` internal/service/review_loop_steps.go:707-764 `ReviewRunCaller.closeOutputs` (before)
  - target `bb13a424dc73` internal/service/review_loop_steps.go:716-759 `ReviewRunCaller.closeOutputs` (after)
- **complexity.modified** [improvement; measured] ReviewRunCaller.settleObservers: cyclomatic 11→7, nesting 3→2, lines 30→26
  - base `471bfd24a696` internal/service/review_loop_steps.go:931-960 `ReviewRunCaller.settleObservers` (before)
  - target `bb13a424dc73` internal/service/review_loop_steps.go:941-966 `ReviewRunCaller.settleObservers` (after)
- … 15 more in JSON

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** ok · **Score:** +1 · **Rule:** rule violation resolved, or edges removed with none added → +1

internal edges 200→199 (+5/-6); cycles 0→0; rule violations 15→11

Coverage: 37 changed code file(s) in scope, 37 analysed, 0 with parse errors, unsupported {}, excluded by class {test:38}, non-code 8

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| internal_units | repository | 86 | 87 | +1 | – |
| internal_edges | repository | 200 | 199 | -1 | – |
| external_dependency_refs | repository | 1028 | 1030 | +2 | – |
| cycles | repository | 0 | 0 | 0 | – |
| rule_violations | repository | 15 | 11 | -4 | – |
| edges_added | change | – | – | – | 5 |
| edges_removed | change | – | – | – | 6 |

Findings (36):

- **dependencies.edge_added** [info; measured] [go] new dependency internal/engine/claudeusage → internal/leaf/retryafter
  - target `bb13a424dc73` internal/engine/claudeusage/claudeusage.go:25-25 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency internal/engine/claudeusage → internal/leaf/regularfile
  - target `bb13a424dc73` internal/engine/claudeusage/credentials.go:11-11 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency internal/engine/codexusage → internal/leaf/regularfile
  - target `bb13a424dc73` internal/engine/codexusage/logs.go:13-13 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency internal/engine/opencodeusage → internal/leaf/retryafter
  - target `bb13a424dc73` internal/engine/opencodeusage/opencodeusage.go:27-27 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency internal/tui/overlay → config
  - target `bb13a424dc73` internal/tui/overlay/macro_picker.go:5-5 (import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed config → internal/engine/macros
  - base `471bfd24a696` config/config.go:14-14 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/server → internal/engine/macros
  - base `471bfd24a696` internal/server/macros.go:12-12 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/service → internal/engine/macros
  - base `471bfd24a696` internal/service/macros.go:12-12 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/tui → internal/engine/ipc
  - base `471bfd24a696` internal/tui/app.go:15-15 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/tui → internal/engine/macros
  - base `471bfd24a696` internal/tui/app.go:16-16 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/tui/overlay → internal/engine/macros
  - base `471bfd24a696` internal/tui/overlay/macro_picker.go:9-9 (removed import)
- **dependencies.fan_change** [info; measured] [go] config fan-out 1→0, fan-in 6→7
- **dependencies.fan_change** [info; measured] [go] internal/engine/claudeusage fan-out 0→2, fan-in 2→2
- **dependencies.fan_change** [info; measured] [go] internal/engine/codexusage fan-out 0→1, fan-in 2→2
- **dependencies.fan_change** [info; measured] [go] internal/engine/macros fan-out 0→0, fan-in 5→0
- **dependencies.fan_change** [info; measured] [go] internal/engine/opencodeusage fan-out 0→1, fan-in 2→2
- **dependencies.fan_change** [info; measured] [go] internal/leaf/regularfile fan-out 0→0, fan-in 0→2
- **dependencies.fan_change** [info; measured] [go] internal/leaf/retryafter fan-out 0→0, fan-in 0→2
- **dependencies.fan_change** [info; measured] [go] internal/server fan-out 18→17, fan-in 2→2
- **dependencies.fan_change** [info; measured] [go] internal/service fan-out 48→47, fan-in 15→15
- **dependencies.fan_change** [info; measured] [go] internal/tui fan-out 19→17, fan-in 1→1
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/events (not scored)
  - target `bb13a424dc73` internal/tui/app.go:15-15 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/store (not scored)
  - target `bb13a424dc73` internal/tui/app.go:16-16 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/workflow (not scored)
  - target `bb13a424dc73` internal/tui/app.go:17-17 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/board → internal/engine/store (not scored)
  - target `bb13a424dc73` internal/tui/board/model.go:9-9 (import)
- … 11 more in JSON

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** ok · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 37 changed code file(s) in scope, 37 analysed, 0 with parse errors, unsupported {}, excluded by class {test:38}, non-code 8

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** ok · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

83 files changed (37 source); symbols +50/-6/~63 (renamed 7, formatting-only 3); exported API +6/-3/~12; 1 unit group(s)

Coverage: 37 changed code file(s) in scope, 37 analysed, 0 with parse errors, unsupported {}, excluded by class {test:38}, non-code 8

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | docs | – | – | – | 8 |
| lines_added | docs | – | – | – | 495 |
| lines_deleted | docs | – | – | – | 11 |
| files_changed | source | – | – | – | 37 |
| lines_added | source | – | – | – | 1014 |
| lines_deleted | source | – | – | – | 684 |
| files_changed | test | – | – | – | 38 |
| lines_added | test | – | – | – | 2988 |
| lines_deleted | test | – | – | – | 79 |
| symbols_added | source | – | – | – | 50 |
| symbols_deleted | source | – | – | – | 6 |
| symbols_modified | source | – | – | – | 63 |
| symbols_renamed | source | – | – | – | 7 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 3 |
| exported_symbols_added | source | – | – | – | 6 |
| exported_symbols_removed | source | – | – | – | 3 |
| exported_symbols_changed | source | – | – | – | 12 |
| units_changed | source | – | – | – | 13 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

Findings (21):

- **scope.api_added** [info; measured] exported symbol added: Macro.Validate
  - target `bb13a424dc73` config/macros.go:92-99 `Macro.Validate` (after)
- **scope.api_added** [info; measured] exported symbol added: SplitFavoriteMacros
  - target `bb13a424dc73` config/macros.go:114-126 `SplitFavoriteMacros` (after)
- **scope.api_added** [info; measured] exported symbol added: Open
  - target `bb13a424dc73` internal/leaf/regularfile/regularfile.go:25-36 `Open` (after)
- **scope.api_added** [info; measured] exported symbol added: EventStream
  - target `bb13a424dc73` internal/service/event_stream.go:11-14 `EventStream` (after)
- **scope.api_added** [info; measured] exported symbol added: PlanVerdict
  - target `bb13a424dc73` internal/service/event_stream.go:20-25 `PlanVerdict` (after)
- **scope.api_added** [info; measured] exported symbol added: SendPlanVerdict
  - target `bb13a424dc73` internal/service/event_stream.go:29-37 `SendPlanVerdict` (after)
- **scope.api_removed** [info; measured] exported symbol removed: Macro.Validate
  - base `471bfd24a696` internal/engine/macros/macros.go:91-98 `Macro.Validate` (before)
- **scope.api_removed** [info; measured] exported symbol removed: ListResult
  - base `471bfd24a696` internal/engine/macros/macros.go:101-103 `ListResult` (before)
- **scope.api_removed** [info; measured] exported symbol removed: SplitFavorites
  - base `471bfd24a696` internal/engine/macros/macros.go:118-130 `SplitFavorites` (before)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: retryAfter → Parse
  - base `471bfd24a696` internal/engine/claudeusage/claudeusage.go:197-209 `retryAfter` (before)
  - target `bb13a424dc73` internal/leaf/retryafter/retryafter.go:19-31 `Parse` (after)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: Macro → Macro
  - base `471bfd24a696` internal/engine/macros/macros.go:21-37 `Macro` (before)
  - target `bb13a424dc73` config/macros.go:22-38 `Macro` (after)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: Macro.UnmarshalYAML → Macro.UnmarshalYAML
  - base `471bfd24a696` internal/engine/macros/macros.go:63-74 `Macro.UnmarshalYAML` (before)
  - target `bb13a424dc73` config/macros.go:64-75 `Macro.UnmarshalYAML` (after)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: Macro.UnmarshalJSON → Macro.UnmarshalJSON
  - base `471bfd24a696` internal/engine/macros/macros.go:77-88 `Macro.UnmarshalJSON` (before)
  - target `bb13a424dc73` config/macros.go:78-89 `Macro.UnmarshalJSON` (after)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: AnyFavorite → AnyFavoriteMacro
  - base `471bfd24a696` internal/engine/macros/macros.go:106-113 `AnyFavorite` (before)
  - target `bb13a424dc73` config/macros.go:102-109 `AnyFavoriteMacro` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: Server.SetMacros (m []macros.Macro) → (m []config.Macro)
  - base `471bfd24a696` internal/server/macros.go:24-26 `Server.SetMacros` (before)
  - target `bb13a424dc73` internal/server/macros.go:23-25 `Server.SetMacros` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: BuiltinMacros () []macros.Macro → () []config.Macro
  - base `471bfd24a696` internal/service/macros.go:71-75 `BuiltinMacros` (before)
  - target `bb13a424dc73` internal/service/macros.go:70-74 `BuiltinMacros` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: MacroService.Macros (ctx context.Context) ([]macros.Macro, error) → (ctx context.Context) ([]config.Macro, error)
  - base `471bfd24a696` internal/service/macros.go:218-228 `MacroService.Macros` (before)
  - target `bb13a424dc73` internal/service/macros.go:217-227 `MacroService.Macros` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: MacroService.Save (ctx context.Context, name string, m macros.Macro, by string) error → (ctx context.Context, name string, m config.Macro, by string) error
  - base `471bfd24a696` internal/service/macros.go:248-268 `MacroService.Save` (before)
  - target `bb13a424dc73` internal/service/macros.go:247-267 `MacroService.Save` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: NewApp (svc service.BoardService, syncSvc service.SyncService, agentSvc service.AgentService, prSvc service.PRTrackingService, reportSvc service.ReportService, icons theme.Icons, bus *events.Bus, editor string, workspaces []service.Workspace, tmux service.TmuxManager, workDir string, swarmSvc service.SwarmService, macrosList []macros.Macro, reviewServices ...review.Service) App → (svc service.BoardService, syncSvc service.SyncService, agentSvc service.AgentService, prSvc service.PRTrackingService, reportSvc service.ReportService, icons theme.Icons, bus service.EventStream, editor string, workspaces []service.Workspace, tmux service.TmuxManager, workDir string, swarmSvc service.SwarmService, macrosList []config.Macro, reviewServices ...review.Service) App
  - base `471bfd24a696` internal/tui/app.go:392-444 `NewApp` (before)
  - target `bb13a424dc73` internal/tui/app.go:393-445 `NewApp` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: App.SetMacroSource (source func(context.Context) ([]macros.Macro, error)) → (source func(context.Context) ([]config.Macro, error))
  - base `471bfd24a696` internal/tui/app.go:2881-2883 `App.SetMacroSource` (before)
  - target `bb13a424dc73` internal/tui/app.go:2880-2882 `App.SetMacroSource` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: NewMacroPicker (macrosList []macros.Macro) MacroPickerOverlay → (macrosList []config.Macro) MacroPickerOverlay
  - base `471bfd24a696` internal/tui/overlay/macro_picker.go:33-36 `NewMacroPicker` (before)
  - target `bb13a424dc73` internal/tui/overlay/macro_picker.go:33-36 `NewMacroPicker` (after)

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 38 test files changed, +2988/-79 lines; test churn is not evidence of correctness)

Coverage: 38 changed code file(s) in scope, 38 analysed, 0 with parse errors, unsupported {}, excluded by class {source:37}, non-code 8

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 38 |
| test_lines_added | context_only | – | – | – | 2988 |
| test_lines_deleted | context_only | – | – | – | 79 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| modified | function | claimChatBots | cmd/legato/chatbots_runtime.go:65-67 | cmd/legato/chatbots_runtime.go:68-82 | source |
| modified | function | claimMQTTPublisher | cmd/legato/mqtt_runtime.go:25-27 | cmd/legato/mqtt_runtime.go:27-41 | source |
| added | function | claimDataDirLock | cmd/legato/runtime_lock.go:16-30 |  | source |
| added | function | TestRuntimeRolesLockTheirOwnFiles | cmd/legato/runtime_lock_test.go:13-45 |  | test |
| modified | function | claimUsagePoller | cmd/legato/usage_runtime.go:30-32 | cmd/legato/usage_runtime.go:32-46 | source |
| added | function | TestConfigImportsNoLegatoPackage | config/architecture_test.go:17-46 |  | test |
| modified | type | Config | config/config.go:17-87 | config/config.go:18-88 | source |
| modified | function | extractMacrosRaw | config/config.go:824-832 | config/config.go:825-833 | source |
| renamed | type | Macro | config/macros.go:22-38 | internal/engine/macros/macros.go:21-37 | source |
| renamed | type | macroFields | config/macros.go:41-50 | internal/engine/macros/macros.go:40-49 | source |
| renamed | method | macroFields.macro | config/macros.go:52-61 | internal/engine/macros/macros.go:51-60 | source |
| renamed | method | Macro.UnmarshalYAML | config/macros.go:64-75 | internal/engine/macros/macros.go:63-74 | source |
| renamed | method | Macro.UnmarshalJSON | config/macros.go:78-89 | internal/engine/macros/macros.go:77-88 | source |
| added | method | Macro.Validate | config/macros.go:92-99 |  | source |
| renamed | function | AnyFavorite → AnyFavoriteMacro | config/macros.go:102-109 | internal/engine/macros/macros.go:106-113 | source |
| added | function | SplitFavoriteMacros | config/macros.go:114-126 |  | source |
| renamed | function | TestMacroJSONRoundTrip | config/macros_test.go:10-23 | internal/engine/macros/macros_test.go:10-23 | test |
| renamed | function | TestMessageIsAnotherNameForKeys | config/macros_test.go:25-46 | internal/engine/macros/macros_test.go:51-72 | test |
| renamed | function | TestMessageAndKeysThatDifferAreRefused | config/macros_test.go:48-56 | internal/engine/macros/macros_test.go:74-82 | test |
| renamed | function | TestMacroJSONStillWritesKeys | config/macros_test.go:58-66 | internal/engine/macros/macros_test.go:84-92 | test |
| modified | method | Client.Read | internal/engine/claudeusage/claudeusage.go:119-173 | internal/engine/claudeusage/claudeusage.go:118-172 | source |
| formatting_only | function | transportError | internal/engine/claudeusage/claudeusage.go:186-196 | internal/engine/claudeusage/claudeusage.go:185-195 | source |
| modified | method | FileCredentials.Login | internal/engine/claudeusage/credentials.go:34-44 | internal/engine/claudeusage/credentials.go:32-42 | source |
| deleted | function | openRegular → openRegular |  | internal/engine/claudeusage/open.go:17-28 | source |
| modified | function | lastRateLimits | internal/engine/codexusage/logs.go:129-169 | internal/engine/codexusage/logs.go:127-167 | source |
| deleted | function | openRegular → openRegular |  | internal/engine/codexusage/open.go:17-28 | source |
| deleted | method | Macro.Validate → Macro.Validate |  | internal/engine/macros/macros.go:91-98 | source |
| deleted | type | ListResult → ListResult |  | internal/engine/macros/macros.go:101-103 | source |
| deleted | function | SplitFavorites → SplitFavorites |  | internal/engine/macros/macros.go:118-130 | source |
| deleted | function | TestListResultJSON → TestListResultJSON |  | internal/engine/macros/macros_test.go:25-49 | test |
| modified | method | Client.Read | internal/engine/opencodeusage/opencodeusage.go:111-159 | internal/engine/opencodeusage/opencodeusage.go:110-158 | source |
| formatting_only | function | transportError | internal/engine/opencodeusage/opencodeusage.go:174-184 | internal/engine/opencodeusage/opencodeusage.go:173-183 | source |
| deleted | function | retryAfter → retryAfter |  | internal/engine/opencodeusage/opencodeusage.go:185-197 | source |
| modified | function | TestEnginePackagesDoNotImportSiblingEnginePackages | internal/engine/store/architecture_test.go:48-77 | internal/engine/store/architecture_test.go:40-67 | test |
| formatting_only | function | TestProductionCommandsUseSharedStoreOpener | internal/engine/store/architecture_test.go:79-112 | internal/engine/store/architecture_test.go:69-102 | test |
| added | function | TestLeafPackagesImportOnlyTheStandardLibrary | internal/engine/store/architecture_test.go:117-149 |  | test |
| modified | function | checkReviewLoopStep | internal/engine/workflow/review_loop.go:149-195 | internal/engine/workflow/review_loop.go:149-209 | source |
| added | function | validReviewLoopPRs | internal/engine/workflow/review_loop.go:199-218 |  | source |
| added | function | TestAReviewLoopOpenStepsPRsMustBeAListOfNumbers | internal/engine/workflow/review_loop_test.go:197-227 |  | test |
| modified | method | upgrader.mapText | internal/engine/workflow/upgrade.go:528-548 | internal/engine/workflow/upgrade.go:528-558 | source |
| added | function | wholeReferences | internal/engine/workflow/upgrade.go:553-569 |  | source |
| added | function | TestAStepRewriteMapsOnlyWholeReferences | internal/engine/workflow/upgrade_test.go:239-279 |  | test |
| modified | function | TestReadyRefusesAnApprovalWithheldFromAnOutsider | internal/handofftest/review_loop_ready_test.go:137-150 | internal/handofftest/review_loop_ready_test.go:137-148 | test |
| formatting_only | function | TestASecondReviewLoopOnAPullRequestIsRefusedWithALinkToTheFirst | internal/handofftest/review_run_one_test.go:10-23 | internal/handofftest/review_run_one_test.go:8-21 | test |
| added | function | TestAStepWaitingToOpenHoldsItsPullRequest | internal/handofftest/review_run_one_test.go:28-57 |  | test |
| added | function | TestAnOpenStepsLoopKeepsItsClones | internal/handofftest/review_run_sweep_test.go:52-90 |  | test |
| added | method | FakeReviewGitHub.editPR | internal/handofftest/reviews.go:141-150 |  | source |
| modified | method | FakeReviewGitHub.Draft | internal/handofftest/reviews.go:153-155 | internal/handofftest/reviews.go:141-150 | source |
| modified | method | FakeReviewGitHub.SetAuthor | internal/handofftest/reviews.go:365-367 | internal/handofftest/reviews.go:360-369 | source |
| modified | method | FakeReviewGitHub.Describe | internal/handofftest/reviews.go:371-373 | internal/handofftest/reviews.go:373-382 | source |
| modified | method | FakeReviewGitHub.Push | internal/handofftest/reviews.go:821-823 | internal/handofftest/reviews.go:830-839 | source |
| added | function | Open | internal/leaf/regularfile/regularfile.go:25-36 |  | source |
| added | function | TestOpenReadsARegularFile | internal/leaf/regularfile/regularfile_test.go:12-26 |  | test |
| added | function | TestOpenKeepsNotExist | internal/leaf/regularfile/regularfile_test.go:28-33 |  | test |
| added | function | TestOpenRefusesAFolder | internal/leaf/regularfile/regularfile_test.go:35-40 |  | test |
| added | function | TestOpenFollowsALinkToARegularFile | internal/leaf/regularfile/regularfile_test.go:42-57 |  | test |
| added | function | TestOpenDoesNotHangOnAFIFO | internal/leaf/regularfile/regularfile_unix_test.go:14-35 |  | test |
| renamed | function | retryAfter → Parse | internal/leaf/retryafter/retryafter.go:19-31 | internal/engine/claudeusage/claudeusage.go:197-209 | source |
| added | function | TestParse | internal/leaf/retryafter/retryafter_test.go:9-34 |  | test |
| modified | method | Server.SetMacros | internal/server/macros.go:23-25 | internal/server/macros.go:24-26 | source |
| modified | method | Server.macrosHandler | internal/server/macros.go:33-63 | internal/server/macros.go:34-64 | source |
| added | type | macroListResult | internal/server/macros.go:66-68 |  | source |
| modified | type | macroWriter | internal/server/macros.go:77-83 | internal/server/macros.go:73-79 | source |
| modified | method | Server.macroEntryHandler | internal/server/macros.go:113-152 | internal/server/macros.go:109-148 | source |
| modified | function | TestMacrosHandlerEmpty | internal/server/macros_test.go:14-32 | internal/server/macros_test.go:14-32 | test |
| modified | function | TestMacrosHandlerUsesLowercaseJSONFields | internal/server/macros_test.go:34-53 | internal/server/macros_test.go:34-53 | test |
| modified | function | TestMacrosHandlerPopulated | internal/server/macros_test.go:55-82 | internal/server/macros_test.go:55-82 | test |
| modified | function | TestMacrosHandlerReadsTheMacroServiceAtEachRequest | internal/server/macros_test.go:124-153 | internal/server/macros_test.go:124-153 | test |
| added | function | TestMacroListResultJSON | internal/server/macros_test.go:156-180 |  | test |
| modified | type | Server | internal/server/server.go:28-151 | internal/server/server.go:28-151 | source |
| added | type | accountResolveFixture | internal/service/agent_account_resolve_test.go:21-26 |  | test |
| added | function | newAccountResolveFixture | internal/service/agent_account_resolve_test.go:28-56 |  | test |
| added | method | accountResolveFixture.setAccountDir | internal/service/agent_account_resolve_test.go:59-65 |  | test |
| added | method | accountResolveFixture.moveWork | internal/service/agent_account_resolve_test.go:68-74 |  | test |
| added | method | accountResolveFixture.workIdentity | internal/service/agent_account_resolve_test.go:76-78 |  | test |
| added | method | accountResolveFixture.defaultIdentity | internal/service/agent_account_resolve_test.go:80-83 |  | test |
| added | function | mustMkdir | internal/service/agent_account_resolve_test.go:85-90 |  | test |
| added | function | mustWrite | internal/service/agent_account_resolve_test.go:92-98 |  | test |
| added | function | mustSymlink | internal/service/agent_account_resolve_test.go:100-106 |  | test |
| added | function | writeAccountLogin | internal/service/agent_account_resolve_test.go:108-111 |  | test |
| added | function | errWorkAuthVar | internal/service/agent_account_resolve_test.go:129-131 |  | test |
| added | function | errWorkNoLogin | internal/service/agent_account_resolve_test.go:133-135 |  | test |
| added | type | accountResolveCase | internal/service/agent_account_resolve_test.go:137-145 |  | test |
| added | function | fixed | internal/service/agent_account_resolve_test.go:147-149 |  | test |
| added | function | named | internal/service/agent_account_resolve_test.go:151-153 |  | test |
| added | function | accountResolveCases | internal/service/agent_account_resolve_test.go:157-592 |  | test |
| added | function | TestResolveClaudeAccountCharacterisation | internal/service/agent_account_resolve_test.go:594-630 |  | test |
| modified | method | agentService.resolveClaudeAccount | internal/service/agent_accounts.go:67-134 | internal/service/agent_accounts.go:67-197 | source |
| added | function | namedAccountEnvConflict | internal/service/agent_accounts.go:139-149 |  | source |
| added | method | agentService.checkNamedAccountIsolation | internal/service/agent_accounts.go:156-215 |  | source |
| added | function | checkNamedAccountLogin | internal/service/agent_accounts.go:219-228 |  | source |
| added | function | TestHistoryConversationTakesTheNewestRunsAccount | internal/service/agent_history_account_test.go:14-53 |  | test |
| modified | method | agentService.prepareTaskResume | internal/service/agent_resume.go:60-140 | internal/service/agent_resume.go:60-186 | source |
| added | method | agentService.inferResumeKind | internal/service/agent_resume.go:146-159 |  | source |
| added | function | resumeStorage | internal/service/agent_resume.go:164-172 |  | source |
| added | method | agentService.resumeConversation | internal/service/agent_resume.go:180-206 |  | source |
| added | method | agentService.bindResumeAccount | internal/service/agent_resume.go:213-229 |  | source |
| modified | method | agentService.historyConversation | internal/service/agent_resume.go:233-246 | internal/service/agent_resume.go:190-213 | source |
| added | method | agentService.latestRunAccount | internal/service/agent_resume.go:251-263 |  | source |
| added | type | prepFixture | internal/service/agent_resume_prepare_test.go:25-31 |  | test |
| added | function | newPrepFixture | internal/service/agent_resume_prepare_test.go:35-57 |  | test |
| added | method | prepFixture.save | internal/service/agent_resume_prepare_test.go:61-85 |  | test |
| added | method | prepFixture.prepare | internal/service/agent_resume_prepare_test.go:87-90 |  | test |
| added | method | prepFixture.defaultIdentity | internal/service/agent_resume_prepare_test.go:93-100 |  | test |
| added | type | resumeFields | internal/service/agent_resume_prepare_test.go:103-106 |  | test |
| added | function | fieldsOf | internal/service/agent_resume_prepare_test.go:108-110 |  | test |
| added | function | wantFields | internal/service/agent_resume_prepare_test.go:112-120 |  | test |
| added | function | wantRefusal | internal/service/agent_resume_prepare_test.go:122-130 |  | test |
| added | function | TestPrepareResumeLeavesOptionsAloneWithoutAResumePolicy | internal/service/agent_resume_prepare_test.go:134-143 |  | test |
| added | function | TestPrepareResumeNewPolicyChecksOnlyTheAdapter | internal/service/agent_resume_prepare_test.go:145-160 |  | test |
| added | function | TestPrepareResumeRefusesAnUnsupportedPolicy | internal/service/agent_resume_prepare_test.go:162-166 |  | test |
| added | function | TestPrepareResumeRefusesAToolThatIsNotConfigured | internal/service/agent_resume_prepare_test.go:168-180 |  | test |
| added | function | TestPrepareResumeInfersTheToolFromTheNewestSavedSession | internal/service/agent_resume_prepare_test.go:182-213 |  | test |
| added | function | TestPrepareResumeSelectsTheSavedOrAnEarlierConversation | internal/service/agent_resume_prepare_test.go:215-289 |  | test |
| added | function | TestPrepareResumeChecksTheSavedDirectory | internal/service/agent_resume_prepare_test.go:291-362 |  | test |
| added | function | TestPrepareResumeBindsTheClaudeAccountTheConversationRanWith | internal/service/agent_resume_prepare_test.go:364-436 |  | test |
| modified | method | agentService.sandboxPrefix | internal/service/agent_sandbox.go:49-240 | internal/service/agent_sandbox.go:48-263 | source |
| formatting_only | function | sessionEnvValues | internal/service/agent_sandbox.go:426-434 | internal/service/agent_sandbox.go:449-457 | source |
| added | function | withClaudeAccounts | internal/service/agent_sandbox.go:443-470 |  | source |
| added | function | claudeMounts | internal/service/agent_sandbox_claude_test.go:18-30 |  | test |
| added | type | claudeFixture | internal/service/agent_sandbox_claude_test.go:35-39 |  | test |
| added | function | newClaudeFixture | internal/service/agent_sandbox_claude_test.go:41-69 |  | test |
| added | function | TestASandboxSeesOnlyItsOwnClaudeAccount | internal/service/agent_sandbox_claude_test.go:77-138 |  | test |
| added | method | claudeFixture.rebase | internal/service/agent_sandbox_claude_test.go:141-148 |  | test |
| added | function | TestWithClaudeAccounts | internal/service/agent_sandbox_claude_test.go:153-188 |  | test |
| added | type | EventStream | internal/service/event_stream.go:11-14 |  | source |
| added | type | PlanVerdict | internal/service/event_stream.go:20-25 |  | source |
| added | function | SendPlanVerdict | internal/service/event_stream.go:29-37 |  | source |
| formatting_only | function | newMacroStartFixture | internal/service/macro_start_test.go:50-67 | internal/service/macro_start_test.go:50-65 | test |
| modified | function | TestAMacroStartsTheTasksAgentWithItsToolAndTier | internal/service/macro_start_test.go:69-93 | internal/service/macro_start_test.go:67-91 | test |
| modified | function | TestMacroAccountIsUsedForFreshWorkAndSavedIdentityForResume | internal/service/macro_start_test.go:182-213 | internal/service/macro_start_test.go:180-211 | test |
| modified | function | TestMacroAccountSelectsClaudeWithoutInheritingOtherToolProfile | internal/service/macro_start_test.go:224-236 | internal/service/macro_start_test.go:222-234 | test |
| modified | function | BuiltinMacros | internal/service/macros.go:70-74 | internal/service/macros.go:71-75 | source |
| modified | function | builtinMacro | internal/service/macros.go:76-83 | internal/service/macros.go:77-84 | source |
| modified | type | MacroEntry | internal/service/macros.go:86-95 | internal/service/macros.go:87-96 | source |
| modified | method | MacroService.Macros | internal/service/macros.go:217-227 | internal/service/macros.go:218-228 | source |
| modified | method | MacroService.Save | internal/service/macros.go:247-267 | internal/service/macros.go:248-268 | source |
| modified | function | macroValue | internal/service/macros.go:280-282 | internal/service/macros.go:281-283 | source |
| modified | method | MacroService.SetFavorite | internal/service/macros.go:286-314 | internal/service/macros.go:287-315 | source |
| modified | method | MacroService.Reset | internal/service/macros.go:388-415 | internal/service/macros.go:389-416 | source |
| modified | function | TestMacrosListBuiltInsThenTheFilesThenTheDatabases | internal/service/macros_test.go:31-67 | internal/service/macros_test.go:32-68 | test |
| modified | function | TestAnEditedBuiltInShadowsItAndResetRestoresIt | internal/service/macros_test.go:69-99 | internal/service/macros_test.go:70-100 | test |
| modified | function | TestABuiltInCannotBeDeleted | internal/service/macros_test.go:101-120 | internal/service/macros_test.go:102-121 | test |
| modified | function | TestAFileEntryIsPinnedEvenOverABuiltIn | internal/service/macros_test.go:122-140 | internal/service/macros_test.go:123-141 | test |
| modified | function | TestAMacroRefusesAnUnknownAppliesTo | internal/service/macros_test.go:142-151 | internal/service/macros_test.go:143-152 | test |
| modified | function | TestSettingsListsMacrosWithTheBuiltIns | internal/service/macros_test.go:153-178 | internal/service/macros_test.go:154-179 | test |
| modified | function | TestAMacroOfYoursCanBeRenamed | internal/service/macros_test.go:180-206 | internal/service/macros_test.go:181-207 | test |
| modified | function | TestARenameIsRefusedForBuiltInsFileMacrosAndNamesInUse | internal/service/macros_test.go:208-237 | internal/service/macros_test.go:209-238 | test |
| modified | function | TestAMacroEntryReadsBackWhole | internal/service/macros_test.go:261-274 | internal/service/macros_test.go:262-275 | test |
| modified | function | TestAStarSurvivesEditsAndAReset | internal/service/macros_test.go:309-359 | internal/service/macros_test.go:310-360 | test |
| modified | method | PullRequestModes.PullRequestMode | internal/service/pull_request_modes.go:23-29 | internal/service/pull_request_modes.go:23-33 | source |
| added | method | PullRequestModes.rules | internal/service/pull_request_modes.go:34-44 |  | source |
| modified | method | PullRequestModes.ReadyAfterReview | internal/service/pull_request_modes.go:56-59 | internal/service/pull_request_modes.go:45-55 | source |
| added | function | TestPullRequestModesAnswerFromTheRulesOrNothing | internal/service/pull_request_modes_test.go:105-140 |  | test |
| modified | function | OpenReviewLoopClones | internal/service/review_loop_open.go:20-90 | internal/service/review_loop_open.go:20-91 | source |
| added | function | stepOwnedLoopHosts | internal/service/review_loop_open.go:94-106 |  | source |
| modified | method | ReviewRunCaller.openSteps | internal/service/review_loop_steps.go:502-526 | internal/service/review_loop_steps.go:502-530 | source |
| added | method | ReviewRunCaller.failWaitingLogged | internal/service/review_loop_steps.go:531-535 |  | source |
| added | method | ReviewRunCaller.completeWaitingLogged | internal/service/review_loop_steps.go:539-543 |  | source |
| modified | method | ReviewRunCaller.SettleSteps | internal/service/review_loop_steps.go:635-665 | internal/service/review_loop_steps.go:622-656 | source |
| modified | method | ReviewRunCaller.closeOutputs | internal/service/review_loop_steps.go:716-759 | internal/service/review_loop_steps.go:707-764 | source |
| added | method | ReviewRunCaller.postedPasses | internal/service/review_loop_steps.go:764-774 |  | source |
| modified | method | ReviewRunCaller.settleObservers | internal/service/review_loop_steps.go:941-966 | internal/service/review_loop_steps.go:931-960 | source |
| modified | method | ReviewRunCaller.settleReady | internal/service/review_loop_steps.go:1205-1232 | internal/service/review_loop_steps.go:1199-1230 | source |
| modified | method | ReviewRunCaller.markStackReady | internal/service/review_loop_steps.go:1242-1275 | internal/service/review_loop_steps.go:1240-1305 | source |
| added | type | readyLayer | internal/service/review_loop_steps.go:1279-1282 |  | source |
| added | method | ReviewRunCaller.approvedOpenLayers | internal/service/review_loop_steps.go:1289-1315 |  | source |
| added | function | notApprovedAtHeads | internal/service/review_loop_steps.go:1319-1321 |  | source |
| added | method | ReviewRunCaller.cleanPassEvents | internal/service/review_loop_steps.go:1326-1334 |  | source |
| added | type | settleWorkflows | internal/service/review_loop_steps_settle_test.go:16-20 |  | test |
| added | method | settleWorkflows.WaitingAttempts | internal/service/review_loop_steps_settle_test.go:22-24 |  | test |
| added | method | settleWorkflows.CompleteWaiting | internal/service/review_loop_steps_settle_test.go:26-29 |  | test |
| added | method | settleWorkflows.FailWaiting | internal/service/review_loop_steps_settle_test.go:31-33 |  | test |
| added | method | settleWorkflows.FailWaitingError | internal/service/review_loop_steps_settle_test.go:35-38 |  | test |
| added | function | captureLog | internal/service/review_loop_steps_settle_test.go:40-48 |  | test |
| added | function | TestSettlingAWaitingStepLogsOnlyARealRefusal | internal/service/review_loop_steps_settle_test.go:52-96 |  | test |
| modified | method | ReviewRunCaller.waitingToOpenOn | internal/service/review_pass_run.go:522-537 | internal/service/review_pass_run.go:522-567 | source |
| added | method | ReviewRunCaller.passesStepWaitingOn | internal/service/review_pass_run.go:542-560 |  | source |
| added | method | ReviewRunCaller.openStepWaitingOn | internal/service/review_pass_run.go:565-587 |  | source |
| modified | method | ReviewRunCaller.AuthorDone | internal/service/review_pass_run.go:1093-1131 | internal/service/review_pass_run.go:1073-1122 | source |
| added | method | ReviewRunCaller.roundPRs | internal/service/review_pass_run.go:1137-1147 |  | source |
| added | method | ReviewRunCaller.closedOutcome | internal/service/review_pass_run.go:1152-1159 |  | source |
| modified | method | ReviewRunCaller.agentDone | internal/service/review_pass_run.go:1195-1211 | internal/service/review_pass_run.go:1158-1198 | source |
| added | function | stewardRoundDone | internal/service/review_pass_run.go:1216-1228 |  | source |
| added | function | reviewRoundDone | internal/service/review_pass_run.go:1233-1249 |  | source |
| added | function | roundItem | internal/service/review_pass_run_round_test.go:9-11 |  | test |
| added | function | TestAStewardedRoundWaitsOnlyForBlockingFindings | internal/service/review_pass_run_round_test.go:16-35 |  | test |
| added | function | TestARoundWaitsForEveryItemButNits | internal/service/review_pass_run_round_test.go:40-62 |  | test |
| modified | type | watchedPull | internal/service/review_publisher_watch.go:137-168 | internal/service/review_publisher_watch.go:137-179 | source |
| modified | method | watchedPull.watch | internal/service/review_publisher_watch.go:206-232 | internal/service/review_publisher_watch.go:217-260 | source |
| added | type | watchedCheck | internal/service/review_publisher_watch.go:236-247 |  | source |
| added | method | watchedCheck.toWatchCheck | internal/service/review_publisher_watch.go:253-274 |  | source |
| formatting_only | function | TestTheWatchReadsEveryLayerInOneGraphQLQuery | internal/service/review_publisher_watch_test.go:28-76 | internal/service/review_publisher_watch_test.go:27-75 | test |
| added | function | TestTheWatchMapsEachCheckContext | internal/service/review_publisher_watch_test.go:82-112 |  | test |
| modified | method | ReviewPassService.addStacked | internal/service/review_steward.go:344-388 | internal/service/review_steward.go:344-411 | source |
| added | function | stackedLayersToAdd | internal/service/review_steward.go:395-421 |  | source |
| added | function | TestStackedLayersToAddAreTheAuthorsOwn | internal/service/review_steward_stacked_test.go:14-47 |  | test |
| modified | method | ReviewPassService.watchLoop | internal/service/review_watch.go:305-370 | internal/service/review_watch.go:305-398 | source |
| added | type | watchRound | internal/service/review_watch.go:375-386 |  | source |
| added | method | ReviewPassService.readWatchRound | internal/service/review_watch.go:392-421 |  | source |

… 60 more in JSON

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
