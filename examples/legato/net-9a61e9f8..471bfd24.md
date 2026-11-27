# CodeEval report: legato

- **Base:** `9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc` (9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc) tree `a55db6103bf4` 2026-11-25T21:49:11+11:00
- **Target:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `3cef63633359769f` · tool 0.1.0-poc
- **Run:** started 2026-11-27T17:11:36+11:00, 0.8s, 0 files parsed, cache 4278 hits / 0 misses
- **Change:** 680 files, +60618/-4046 lines, 1836 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | partial | -2 | structural redundancy 243→260 (Δ+17), exact 8→10; 3 near-structural candidate(s) |
| complexity | complexity | partial | -2 | modified callables net cyclomatic Δ+194; 569 added (cc 2976), 6 deleted (cc 21); hotspots 104→195 |
| architecture_coupling | dependencies | partial | 0 | internal edges 193→200 (+7/-0); cycles 0→0; rule violations 15→15 |
| maintainability_scope | obsolete | partial | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | partial | 0 | 680 files changed (296 source); symbols +707/-6/~326 (renamed 2, formatting-only 13); exported API +280/-4/~1; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 250 test files changed, +19910/-545 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 2207 | 363 | 32703 |  |
| javascript | no | 33 | 2 | 119 | unsupported: no adapter |
| kotlin | no | 57 | 2 | 104 | unsupported: no adapter |
| python | no | 10 | 0 | 0 | unsupported: no adapter |
| rust | no | 34 | 0 | 0 | unsupported: no adapter |
| shell | no | 15 | 0 | 0 | unsupported: no adapter |
| sql | no | 160 | 7 | 36 | unsupported: no adapter |
| swift | no | 2 | 0 | 0 | unsupported: no adapter |
| typescript | no | 999 | 171 | 9556 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
docs 96 (+1564/-774); generated 2 (+9249/-23); other 35 (+9421/-945); source 296 (+20462/-1759); test 250 (+19910/-545); testdata 1 (+12/-0); 

## duplication — Callable duplication delta (`duplication`)

**Status:** partial · **Score:** -2 · **Rule:** net structural redundancy Δ+17: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 243→260 (Δ+17), exact 8→10; 3 near-structural candidate(s)

Coverage: 296 changed code file(s) in scope, 201 analysed, 0 with parse errors, unsupported {javascript:2, kotlin:2, sql:7, typescript:84}, excluded by class {test:249}, non-code 135

- callables with ≥40 tokens in production source are compared repository-wide (base 6362, target 6762)
- introduced 20 redundant instance(s), removed 3; 1 pre-existing group(s) touched but unchanged in size
- partial coverage: score reflects 201 analysed file(s); 95 unsupported (javascript:2, kotlin:2, sql:7, typescript:84) and 0 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 8 | 10 | +2 | – |
| structural_redundant_callables | repository | 243 | 260 | +17 | – |
| callables_considered | repository | 6362 | 6762 | +400 | – |
| redundant_instances_introduced | change | – | – | – | 20 |
| redundant_instances_removed | change | – | – | – | 3 |
| preexisting_groups_touched | change | – | – | – | 1 |
| near_structural_candidates | change | – | – | – | 3 |

Findings (26):

- **duplication.introduced_exact** [degradation; measured] identical callable body introduced: [retryAfter retryAfter]
  - target `471bfd24a696` internal/engine/claudeusage/claudeusage.go:197-209 `retryAfter` (new or changed (added))
  - target `471bfd24a696` internal/engine/opencodeusage/opencodeusage.go:185-197 `retryAfter` (new or changed (added))
- **duplication.introduced_exact** [degradation; measured] identical callable body introduced: [openRegular openRegular]
  - target `471bfd24a696` internal/engine/claudeusage/open.go:17-28 `openRegular` (new or changed (added))
  - target `471bfd24a696` internal/engine/codexusage/open.go:17-28 `openRegular` (new or changed (added))
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [claimUsagePoller]; matches pre-existing [claimChatBots claimMQTTPublisher] — possible reuse candidate
  - target `471bfd24a696` cmd/legato/chatbots_runtime.go:68-82 `claimChatBots` (existing match)
  - target `471bfd24a696` cmd/legato/mqtt_runtime.go:27-41 `claimMQTTPublisher` (existing match)
  - target `471bfd24a696` cmd/legato/usage_runtime.go:32-46 `claimUsagePoller` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [accountPathContains]; matches pre-existing [pathWithin] — possible reuse candidate
  - target `471bfd24a696` config/claude_accounts.go:76-79 `accountPathContains` (new or changed (added))
  - target `471bfd24a696` internal/service/learning_drift.go:153-156 `pathWithin` (existing match)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [UsageProfile.UnmarshalYAML UsageProfileReviewer.UnmarshalYAML]
  - target `471bfd24a696` config/usage_profiles.go:120-126 `UsageProfile.UnmarshalYAML` (new or changed (added))
  - target `471bfd24a696` config/usage_profiles.go:128-134 `UsageProfileReviewer.UnmarshalYAML` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [Client.ClaimTask Client.ReleaseClaim Client.SetStage]
  - target `471bfd24a696` internal/engine/burnlite/assign.go:19-22 `Client.ClaimTask` (new or changed (added))
  - target `471bfd24a696` internal/engine/burnlite/assign.go:26-29 `Client.ReleaseClaim` (new or changed (added))
  - target `471bfd24a696` internal/engine/burnlite/assign.go:33-36 `Client.SetStage` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [CredentialsPath SessionsDir]
  - target `471bfd24a696` internal/engine/claudeusage/credentials.go:17-22 `CredentialsPath` (new or changed (added))
  - target `471bfd24a696` internal/engine/codexusage/logs.go:24-29 `SessionsDir` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [Client.doPut]; matches pre-existing [Client.doPost] — possible reuse candidate
  - target `471bfd24a696` internal/engine/jira/client.go:157-164 `Client.doPost` (existing match)
  - target `471bfd24a696` internal/engine/jira/users.go:57-64 `Client.doPut` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [Store.LatestWorkflowRuns]; matches pre-existing [Store.ListAgentSessions Store.ListInterruptedAgentSessions Store.ListAwayTasks] — possible reuse candidate
  - target `471bfd24a696` internal/engine/store/agent_session.go:104-109 `Store.ListAgentSessions` (existing match)
  - target `471bfd24a696` internal/engine/store/agent_session.go:396-401 `Store.ListInterruptedAgentSessions` (existing match)
  - target `471bfd24a696` internal/engine/store/handoffs.go:287-291 `Store.ListAwayTasks` (existing match)
  - target `471bfd24a696` internal/engine/store/workflow_stats.go:17-23 `Store.LatestWorkflowRuns` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [Store.ListSubscriptionUsage]; matches pre-existing [Store.ListSecrets Store.ListLaunchMessages Store.GitHubAppLogins Store.ListReevePeers Store.ListReeveInbox Store.ListRepositoryGuidance Store.ListRepositoryPresets Store.ListSettings Store.ListSettingEntries Store.WorkflowTriggerSwitches] — possible reuse candidate
  - target `471bfd24a696` internal/engine/store/credential.go:57-61 `Store.ListSecrets` (existing match)
  - target `471bfd24a696` internal/engine/store/launch_messages.go:37-42 `Store.ListLaunchMessages` (existing match)
  - target `471bfd24a696` internal/engine/store/pr_reviews_seen.go:64-68 `Store.GitHubAppLogins` (existing match)
  - target `471bfd24a696` internal/engine/store/reeve_peers.go:45-49 `Store.ListReevePeers` (existing match)
  - target `471bfd24a696` internal/engine/store/reeve_peers.go:105-109 `Store.ListReeveInbox` (existing match)
  - target `471bfd24a696` internal/engine/store/repository_guidance.go:21-25 `Store.ListRepositoryGuidance` (existing match)
  - target `471bfd24a696` internal/engine/store/repository_preset.go:32-36 `Store.ListRepositoryPresets` (existing match)
  - target `471bfd24a696` internal/engine/store/settings.go:19-23 `Store.ListSettings` (existing match)
  - target `471bfd24a696` internal/engine/store/settings.go:75-79 `Store.ListSettingEntries` (existing match)
  - target `471bfd24a696` internal/engine/store/subscription_usage.go:14-18 `Store.ListSubscriptionUsage` (new or changed (added))
  - target `471bfd24a696` internal/engine/store/workflow_trigger_switches.go:14-18 `Store.WorkflowTriggerSwitches` (existing match)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [Store.ListReviewLoopsClosedSince]; matches pre-existing [Store.OAuthClientRegistrations Store.ListSubtasksByParentAndStatus] — possible reuse candidate
  - target `471bfd24a696` internal/engine/store/oauth_client_registrations.go:61-65 `Store.OAuthClientRegistrations` (existing match)
  - target `471bfd24a696` internal/engine/store/review_loops.go:574-579 `Store.ListReviewLoopsClosedSince` (new or changed (added))
  - target `471bfd24a696` internal/engine/store/swarm.go:78-84 `Store.ListSubtasksByParentAndStatus` (existing match)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [FakeReviewGitHub.Draft]; matches pre-existing [FakeReviewGitHub.Push] — possible reuse candidate
  - target `471bfd24a696` internal/handofftest/reviews.go:141-150 `FakeReviewGitHub.Draft` (new or changed (added))
  - target `471bfd24a696` internal/handofftest/reviews.go:830-839 `FakeReviewGitHub.Push` (existing match)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [Server.StartSubscriptionUsageEvents]; matches pre-existing [Server.StartServerThemeEvents] — possible reuse candidate
  - target `471bfd24a696` internal/server/subscription_usage.go:70-87 `Server.StartSubscriptionUsageEvents` (new or changed (added))
  - target `471bfd24a696` internal/server/theme_preference.go:23-40 `Server.StartServerThemeEvents` (existing match)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [PullRequestModes.ReadyAfterReview]; matches pre-existing [PullRequestModes.PullRequestMode] — possible reuse candidate
  - target `471bfd24a696` internal/service/pull_request_modes.go:23-33 `PullRequestModes.PullRequestMode` (existing match)
  - target `471bfd24a696` internal/service/pull_request_modes.go:45-55 `PullRequestModes.ReadyAfterReview` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [ReviewLoopSteering.StopWatching ReviewLoopSteering.WatchAgain]
  - target `471bfd24a696` internal/service/review_loop_steering.go:102-107 `ReviewLoopSteering.StopWatching` (new or changed (added))
  - target `471bfd24a696` internal/service/review_loop_steering.go:111-116 `ReviewLoopSteering.WatchAgain` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [ReviewPassService.StopWatchingOn ReviewPassService.WatchAgainOn]
  - target `471bfd24a696` internal/service/review_steward_hold.go:102-108 `ReviewPassService.StopWatchingOn` (new or changed (added))
  - target `471bfd24a696` internal/service/review_steward_hold.go:114-120 `ReviewPassService.WatchAgainOn` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [ReviewRunCaller.stopWatchingOn ReviewRunCaller.watchAgainOn]
  - target `471bfd24a696` internal/service/review_steward_hold.go:203-208 `ReviewRunCaller.stopWatchingOn` (new or changed (added))
  - target `471bfd24a696` internal/service/review_steward_hold.go:210-215 `ReviewRunCaller.watchAgainOn` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [claudeUsageError opencodeUsageError]
  - target `471bfd24a696` internal/service/subscription_usage_readers.go:52-73 `claudeUsageError` (new or changed (added))
  - target `471bfd24a696` internal/service/subscription_usage_readers.go:176-197 `opencodeUsageError` (new or changed (added))
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.introduced_structural** [degradation; measured] structurally identical callable introduced (names/literals differ): [App.loadSubscriptionUsage]; matches pre-existing [App.loadUsage] — possible reuse candidate
  - target `471bfd24a696` internal/tui/app_subscription_usage.go:35-44 `App.loadSubscriptionUsage` (new or changed (added))
  - target `471bfd24a696` internal/tui/app_usage.go:31-40 `App.loadUsage` (existing match)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.near_structural_candidate** [info; heuristic, confidence low] Store.ListWorkflowRunsForTask (added) is structurally similar to Store.ListAgentUIItems (estimated shingle Jaccard 0.88)
  - target `471bfd24a696` internal/engine/store/workflow.go:1041-1053 `Store.ListWorkflowRunsForTask` (changed callable)
  - target `471bfd24a696` internal/engine/store/agent_ui.go:144-152 `Store.ListAgentUIItems` (existing callable)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.near_structural_candidate** [info; heuristic, confidence low] Store.ListWorkflowRunsForTask (added) is structurally similar to Store.ListWorkflowTriggerRequestsForTask (estimated shingle Jaccard 0.88)
  - target `471bfd24a696` internal/engine/store/workflow.go:1041-1053 `Store.ListWorkflowRunsForTask` (changed callable)
  - target `471bfd24a696` internal/engine/store/workflow_trigger.go:103-107 `Store.ListWorkflowTriggerRequestsForTask` (existing callable)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.near_structural_candidate** [info; heuristic, confidence low] Store.ListWorkflowRunsForTask (added) is structurally similar to Store.ListWorkflowTriggerRequestsForTransition (estimated shingle Jaccard 0.88)
  - target `471bfd24a696` internal/engine/store/workflow.go:1041-1053 `Store.ListWorkflowRunsForTask` (changed callable)
  - target `471bfd24a696` internal/engine/store/workflow_trigger.go:97-101 `Store.ListWorkflowTriggerRequestsForTransition` (existing callable)
  - _caveat: structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse_
- **duplication.preexisting_touched** [neutral; measured] pre-existing clone group of 3 still has 3 members; not counted as introduced
  - target `471bfd24a696` internal/engine/store/review_loops.go:783-797 `ReviewLoopTx.SavePass` (changed member (modified))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 1 (removed/changed: [Store.ListWorkflowRunsByRelatedTask])
  - base `9a61e9f8a021` internal/engine/store/agent_session.go:158-163 `Store.ListTaskAgentSessions` (remaining copy)
  - base `9a61e9f8a021` internal/engine/store/workflow.go:1030-1034 `Store.ListWorkflowRunsByRelatedTask` (removed or changed (deleted))
- **duplication.removed** [improvement; measured] clone group shrank from 2 to 1 (removed/changed: [passOf])
  - base `9a61e9f8a021` internal/engine/store/handoff_reviews.go:565-569 `passOf` (removed or changed (modified))
  - base `9a61e9f8a021` internal/engine/store/handoff_reviews.go:985-989 `findingOf` (remaining copy)
- … 1 more in JSON

## complexity — Callable complexity deltas (`complexity`)

**Status:** partial · **Score:** -2 · **Rule:** pre-existing hotspotΔ ≥ 2 or modified net cyclomatic ≥ 10 → -2

modified callables net cyclomatic Δ+194; 569 added (cc 2976), 6 deleted (cc 21); hotspots 104→195

Coverage: 296 changed code file(s) in scope, 201 analysed, 0 with parse errors, unsupported {javascript:2, kotlin:2, sql:7, typescript:84}, excluded by class {test:249}, non-code 135

- modified callables: 129 increased, 24 decreased; net cyclomatic Δ+194
- added callables: 569 (cyclomatic 2976); deleted callables: 6 (cyclomatic 21)
- changed-scope cyclomatic 3517→6666 (Δ+3149), hotspots 104→195
- partial coverage: score reflects 201 analysed file(s); 95 unsupported (javascript:2, kotlin:2, sql:7, typescript:84) and 0 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 3496 | 3690 | +194 | – |
| decision_points_sum | modified_callables | 3255 | 3449 | +194 | – |
| lines_sum | modified_callables | 13266 | 13861 | +595 | – |
| max_nesting | modified_callables | 6 | 6 | 0 | – |
| cyclomatic_sum | changed_scope | 3517 | 6666 | +3149 | – |
| max_cyclomatic | changed_scope | 304 | 308 | +4 | – |
| hotspot_callables | changed_scope | 104 | 195 | +91 | – |
| callables | changed_scope | 247 | 810 | +563 | – |
| callables_modified |  | – | – | – | 241 |
| callables_added |  | – | – | – | 569 |
| callables_deleted |  | – | – | – | 6 |
| cyclomatic_sum_added | added_callables | – | – | – | 2976 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 21 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 7 |
| hotspots_added | added_callables | – | – | – | 84 |

Findings (278):

- **complexity.added_hotspot** [degradation; measured] new callable formatSubscriptionUsage has cyclomatic 19, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` cmd/legato/usage_cmd.go:80-142 `formatSubscriptionUsage` (added)
- **complexity.added_hotspot** [degradation; measured] new callable runUsagePinCmd has cyclomatic 13, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` cmd/legato/usage_cmd.go:170-216 `runUsagePinCmd` (added)
- **complexity.added_hotspot** [degradation; measured] new callable runWorkflowUpgrade has cyclomatic 17, nesting 2 (hotspot thresholds 10/4)
  - target `471bfd24a696` cmd/legato/workflow_upgrade.go:19-71 `runWorkflowUpgrade` (added)
- **complexity.added_hotspot** [degradation; measured] new callable ValidateClaudeAccounts has cyclomatic 21, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` config/claude_accounts.go:20-74 `ValidateClaudeAccounts` (added)
- **complexity.added_hotspot** [degradation; measured] new callable ParseUsageMeter has cyclomatic 10, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` config/claude_accounts.go:107-131 `ParseUsageMeter` (added)
- **complexity.added_hotspot** [degradation; measured] new callable ValidateUsageProfiles has cyclomatic 19, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` config/usage_profiles.go:28-83 `ValidateUsageProfiles` (added)
- **complexity.added_hotspot** [degradation; measured] new callable Client.Read has cyclomatic 16, nesting 1 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/claudeusage/claudeusage.go:118-172 `Client.Read` (added)
- **complexity.added_hotspot** [degradation; measured] new callable decode has cyclomatic 13, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/claudeusage/claudeusage.go:216-250 `decode` (added)
- **complexity.added_hotspot** [degradation; measured] new callable AppServer.Read has cyclomatic 11, nesting 2 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/codexusage/appserver.go:47-105 `AppServer.Read` (added)
- **complexity.added_hotspot** [degradation; measured] new callable converse has cyclomatic 15, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/codexusage/appserver.go:109-163 `converse` (added)
- **complexity.added_hotspot** [degradation; measured] new callable lookPath has cyclomatic 8, nesting 4 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/codexusage/appserver.go:178-195 `lookPath` (added)
- **complexity.added_hotspot** [degradation; measured] new callable newestRollouts has cyclomatic 10, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/codexusage/logs.go:60-87 `newestRollouts` (added)
- **complexity.added_hotspot** [degradation; measured] new callable lastRateLimits has cyclomatic 13, nesting 2 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/codexusage/logs.go:127-167 `lastRateLimits` (added)
- **complexity.added_hotspot** [degradation; measured] new callable Client.Read has cyclomatic 15, nesting 1 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/opencodeusage/opencodeusage.go:110-158 `Client.Read` (added)
- **complexity.added_hotspot** [degradation; measured] new callable decodeMeter has cyclomatic 15, nesting 1 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/opencodeusage/opencodeusage.go:237-266 `decodeMeter` (added)
- **complexity.added_hotspot** [degradation; measured] new callable Store.AdoptHandoffReviewLoop has cyclomatic 11, nesting 2 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/store/handoff_reviews.go:345-377 `Store.AdoptHandoffReviewLoop` (added)
- **complexity.added_hotspot** [degradation; measured] new callable graphBuilder.scope has cyclomatic 10, nesting 4 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/graph.go:223-260 `graphBuilder.scope` (added)
- **complexity.added_hotspot** [degradation; measured] new callable conditionText has cyclomatic 19, nesting 4 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/graph.go:392-447 `conditionText` (added)
- **complexity.added_hotspot** [degradation; measured] new callable checkReviewLoopStep has cyclomatic 29, nesting 4 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/review_loop.go:149-209 `checkReviewLoopStep` (added)
- **complexity.added_hotspot** [degradation; measured] new callable SetTags has cyclomatic 20, nesting 2 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/tags.go:60-110 `SetTags` (added)
- **complexity.added_hotspot** [degradation; measured] new callable upgrade has cyclomatic 23, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/upgrade.go:138-207 `upgrade` (added)
- **complexity.added_hotspot** [degradation; measured] new callable upgrader.rewriteSteps has cyclomatic 11, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/upgrade.go:250-288 `upgrader.rewriteSteps` (added)
- **complexity.added_hotspot** [degradation; measured] new callable upgrader.editStep has cyclomatic 13, nesting 2 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/upgrade.go:292-321 `upgrader.editStep` (added)
- **complexity.added_hotspot** [degradation; measured] new callable upgrader.wrapWith has cyclomatic 14, nesting 2 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/upgrade.go:325-380 `upgrader.wrapWith` (added)
- **complexity.added_hotspot** [degradation; measured] new callable upgrader.upgradeStep has cyclomatic 23, nesting 3 (hotspot thresholds 10/4)
  - target `471bfd24a696` internal/engine/workflow/upgrade.go:383-448 `upgrader.upgradeStep` (added)
- … 253 more in JSON

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** partial · **Score:** 0 · **Rule:** no cycles or rule violations changed → 0 (new edges alone are reported, not judged: features legitimately add dependencies)

internal edges 193→200 (+7/-0); cycles 0→0; rule violations 15→15

Coverage: 296 changed code file(s) in scope, 201 analysed, 0 with parse errors, unsupported {javascript:2, kotlin:2, sql:7, typescript:84}, excluded by class {test:249}, non-code 135

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod
- partial coverage: score reflects 201 analysed file(s); 95 unsupported (javascript:2, kotlin:2, sql:7, typescript:84) and 0 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| internal_units | repository | 83 | 86 | +3 | – |
| internal_edges | repository | 193 | 200 | +7 | – |
| external_dependency_refs | repository | 982 | 1028 | +46 | – |
| cycles | repository | 0 | 0 | 0 | – |
| rule_violations | repository | 15 | 15 | 0 | – |
| edges_added | change | – | – | – | 7 |
| edges_removed | change | – | – | – | 0 |

Findings (29):

- **dependencies.edge_added** [info; measured] [go] new dependency cmd/legato → internal/engine/claudeusage
  - target `471bfd24a696` cmd/legato/usage_runtime.go:12-12 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency cmd/legato → internal/engine/codexusage
  - target `471bfd24a696` cmd/legato/usage_runtime.go:13-13 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency cmd/legato → internal/engine/opencodeusage
  - target `471bfd24a696` cmd/legato/usage_runtime.go:17-17 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency internal/service → internal/engine/claudeusage
  - target `471bfd24a696` internal/service/agent_accounts.go:13-13 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency internal/service → internal/engine/codexusage
  - target `471bfd24a696` internal/service/subscription_usage_readers.go:9-9 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency internal/service → internal/engine/opencodeusage
  - target `471bfd24a696` internal/service/subscription_usage_readers.go:10-10 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency internal/tui → config
  - target `471bfd24a696` internal/tui/app_spawn.go:7-7 (import)
- **dependencies.fan_change** [info; measured] [go] cmd/legato fan-out 42→45, fan-in 0→0
- **dependencies.fan_change** [info; measured] [go] config fan-out 1→1, fan-in 5→6
- **dependencies.fan_change** [info; measured] [go] internal/engine/claudeusage fan-out 0→0, fan-in 0→2
- **dependencies.fan_change** [info; measured] [go] internal/engine/codexusage fan-out 0→0, fan-in 0→2
- **dependencies.fan_change** [info; measured] [go] internal/engine/opencodeusage fan-out 0→0, fan-in 0→2
- **dependencies.fan_change** [info; measured] [go] internal/service fan-out 45→48, fan-in 15→15
- **dependencies.fan_change** [info; measured] [go] internal/tui fan-out 18→19, fan-in 1→1
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "config is standalone, no internal imports (CLAUDE.md prose)" already violated before the change: config → internal/engine/macros (not scored)
  - target `471bfd24a696` config/config.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/events (not scored)
  - target `471bfd24a696` internal/tui/app.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/ipc (not scored)
  - target `471bfd24a696` internal/tui/app.go:15-15 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/macros (not scored)
  - target `471bfd24a696` internal/tui/app.go:16-16 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/store (not scored)
  - target `471bfd24a696` internal/tui/app.go:17-17 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/workflow (not scored)
  - target `471bfd24a696` internal/tui/app.go:18-18 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/board → internal/engine/store (not scored)
  - target `471bfd24a696` internal/tui/board/model.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/github (not scored)
  - target `471bfd24a696` internal/tui/overlay/link_pr.go:10-10 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/macros (not scored)
  - target `471bfd24a696` internal/tui/overlay/macro_picker.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/workflow (not scored)
  - target `471bfd24a696` internal/tui/overlay/workflow_trigger.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/store (not scored)
  - target `471bfd24a696` internal/tui/overlay/workspace.go:9-9 (import)
- … 4 more in JSON

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** partial · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 296 changed code file(s) in scope, 201 analysed, 0 with parse errors, unsupported {javascript:2, kotlin:2, sql:7, typescript:84}, excluded by class {test:249}, non-code 135

- partial coverage: score reflects 201 analysed file(s); 95 unsupported (javascript:2, kotlin:2, sql:7, typescript:84) and 0 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** partial · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

680 files changed (296 source); symbols +707/-6/~326 (renamed 2, formatting-only 13); exported API +280/-4/~1; 1 unit group(s)

Coverage: 296 changed code file(s) in scope, 201 analysed, 0 with parse errors, unsupported {javascript:2, kotlin:2, sql:7, typescript:84}, excluded by class {test:249}, non-code 135

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction
- partial coverage: score reflects 201 analysed file(s); 95 unsupported (javascript:2, kotlin:2, sql:7, typescript:84) and 0 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | docs | – | – | – | 96 |
| lines_added | docs | – | – | – | 1564 |
| lines_deleted | docs | – | – | – | 774 |
| files_changed | generated | – | – | – | 2 |
| lines_added | generated | – | – | – | 9249 |
| lines_deleted | generated | – | – | – | 23 |
| files_changed | other | – | – | – | 35 |
| lines_added | other | – | – | – | 9421 |
| lines_deleted | other | – | – | – | 945 |
| files_changed | source | – | – | – | 296 |
| lines_added | source | – | – | – | 20462 |
| lines_deleted | source | – | – | – | 1759 |
| files_changed | test | – | – | – | 250 |
| lines_added | test | – | – | – | 19910 |
| lines_deleted | test | – | – | – | 545 |
| files_changed | testdata | – | – | – | 1 |
| lines_added | testdata | – | – | – | 12 |
| lines_deleted | testdata | – | – | – | 0 |
| symbols_added | source | – | – | – | 707 |
| symbols_deleted | source | – | – | – | 6 |
| symbols_modified | source | – | – | – | 326 |
| symbols_renamed | source | – | – | – | 2 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 13 |
| exported_symbols_added | source | – | – | – | 280 |
| exported_symbols_removed | source | – | – | – | 4 |
| exported_symbols_changed | source | – | – | – | 1 |
| units_changed | source | – | – | – | 21 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

Findings (285):

- **scope.api_added** [info; measured] exported symbol added: ClaudeAccount
  - target `471bfd24a696` config/claude_accounts.go:12-16 `ClaudeAccount` (after)
- **scope.api_added** [info; measured] exported symbol added: ValidateClaudeAccounts
  - target `471bfd24a696` config/claude_accounts.go:20-74 `ValidateClaudeAccounts` (after)
- **scope.api_added** [info; measured] exported symbol added: Config.ValidateClaudeSelection
  - target `471bfd24a696` config/claude_accounts.go:81-94 `Config.ValidateClaudeSelection` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageReaderID
  - target `471bfd24a696` config/claude_accounts.go:98-103 `UsageReaderID` (after)
- **scope.api_added** [info; measured] exported symbol added: ParseUsageMeter
  - target `471bfd24a696` config/claude_accounts.go:107-131 `ParseUsageMeter` (after)
- **scope.api_added** [info; measured] exported symbol added: Config.IsUsageMeter
  - target `471bfd24a696` config/claude_accounts.go:133-136 `Config.IsUsageMeter` (after)
- **scope.api_added** [info; measured] exported symbol added: NotificationsConfig.UsageEnabled
  - target `471bfd24a696` config/config.go:303-305 `NotificationsConfig.UsageEnabled` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageConfig
  - target `471bfd24a696` config/config.go:1231-1241 `UsageConfig` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageProviderConfig
  - target `471bfd24a696` config/config.go:1244-1247 `UsageProviderConfig` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageProviderConfig.IsEnabled
  - target `471bfd24a696` config/config.go:1250-1250 `UsageProviderConfig.IsEnabled` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageOpenCodeConfig
  - target `471bfd24a696` config/config.go:1254-1262 `UsageOpenCodeConfig` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageOpenCodeConfig.IsEnabled
  - target `471bfd24a696` config/config.go:1265-1265 `UsageOpenCodeConfig.IsEnabled` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageConfig.PollIntervalDuration
  - target `471bfd24a696` config/config.go:1275-1280 `UsageConfig.PollIntervalDuration` (after)
- **scope.api_added** [info; measured] exported symbol added: PullRequestsConfig.ReadyAfterReview
  - target `471bfd24a696` config/pull_requests.go:95-98 `PullRequestsConfig.ReadyAfterReview` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageProfile
  - target `471bfd24a696` config/usage_profiles.go:12-18 `UsageProfile` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageProfileReviewer
  - target `471bfd24a696` config/usage_profiles.go:21-25 `UsageProfileReviewer` (after)
- **scope.api_added** [info; measured] exported symbol added: ValidateUsageProfiles
  - target `471bfd24a696` config/usage_profiles.go:28-83 `ValidateUsageProfiles` (after)
- **scope.api_added** [info; measured] exported symbol added: IsUsageMeter
  - target `471bfd24a696` config/usage_profiles.go:87-90 `IsUsageMeter` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageProfile.UnmarshalYAML
  - target `471bfd24a696` config/usage_profiles.go:120-126 `UsageProfile.UnmarshalYAML` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageProfileReviewer.UnmarshalYAML
  - target `471bfd24a696` config/usage_profiles.go:128-134 `UsageProfileReviewer.UnmarshalYAML` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageRule
  - target `471bfd24a696` config/usage_profiles.go:149-155 `UsageRule` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageMeterInfo
  - target `471bfd24a696` config/usage_providers.go:14-21 `UsageMeterInfo` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageProviderInfo
  - target `471bfd24a696` config/usage_providers.go:24-33 `UsageProviderInfo` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageProviders
  - target `471bfd24a696` config/usage_providers.go:53-53 `UsageProviders` (after)
- **scope.api_added** [info; measured] exported symbol added: UsageMeterInfoOf
  - target `471bfd24a696` config/usage_providers.go:56-63 `UsageMeterInfoOf` (after)
- … 260 more in JSON

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 250 test files changed, +19910/-545 lines; test churn is not evidence of correctness)

Coverage: 249 changed code file(s) in scope, 162 analysed, 0 with parse errors, unsupported {typescript:87}, excluded by class {source:296}, non-code 135

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 250 |
| test_lines_added | context_only | – | – | – | 19910 |
| test_lines_deleted | context_only | – | – | – | 545 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| modified | method | approvalRuntime.trackAgentBranches | cmd/legato/agent_branches_runtime.go:13-23 | cmd/legato/agent_branches_runtime.go:13-20 | source |
| modified | function | wireAttention | cmd/legato/attention_runtime.go:16-77 | cmd/legato/attention_runtime.go:16-76 | source |
| added | function | claudeManagedSettingsDir | cmd/legato/claude_account_runtime.go:5-10 |  | source |
| added | function | TestCommandAgentsRefuseManagedAuthForNamedSubscription | cmd/legato/claude_account_runtime_linux_test.go:20-101 |  | test |
| modified | function | commandTree | cmd/legato/command_tree.go:94-465 | cmd/legato/command_tree.go:94-451 | source |
| modified | function | configureWorkflowCommands | cmd/legato/command_tree.go:629-660 | cmd/legato/command_tree.go:615-645 | source |
| modified | function | handleServeIPCMessage | cmd/legato/ipc.go:58-84 | cmd/legato/ipc.go:58-84 | source |
| modified | function | handleIPCMessage | cmd/legato/ipc.go:86-233 | cmd/legato/ipc.go:86-224 | source |
| modified | function | runMacroRunContract | cmd/legato/macro_contract.go:31-84 | cmd/legato/macro_contract.go:31-84 | source |
| added | method | macroRunSpy.NewWorkDefaults | cmd/legato/macro_contract_test.go:38-40 |  | test |
| deleted | method | macroRunSpy.DefaultAdapter → macroRunSpy.DefaultAdapter |  | cmd/legato/macro_contract_test.go:38-38 | test |
| formatting_only | function | captureStderr | cmd/legato/macro_contract_test.go:118-132 | cmd/legato/macro_contract_test.go:116-130 | test |
| added | function | TestMacroRunExplicitClaudeAccountReachesSpawn | cmd/legato/macro_contract_test.go:134-140 |  | test |
| modified | function | runServeCmd | cmd/legato/main.go:882-1230 | cmd/legato/main.go:882-1221 | source |
| modified | function | runTUI | cmd/legato/main.go:1232-1806 | cmd/legato/main.go:1223-1775 | source |
| modified | function | loadSwarmServiceForCLI | cmd/legato/main.go:2080-2147 | cmd/legato/main.go:2049-2115 | source |
| modified | function | classifyLegatoCommand | cmd/legato/reeve_commands.go:28-56 | cmd/legato/reeve_commands.go:28-50 | source |
| modified | method | legatoCommands.Commands | cmd/legato/reeve_commands.go:139-156 | cmd/legato/reeve_commands.go:133-151 | source |
| modified | function | TestEveryLegatoSubcommandHasAClass | cmd/legato/reeve_commands_test.go:9-47 | cmd/legato/reeve_commands_test.go:9-43 | test |
| modified | function | reeveDefaults | cmd/legato/reeve_runtime.go:89-91 | cmd/legato/reeve_runtime.go:89-91 | source |
| modified | function | TestTheReviewEngineBuiltWithoutAWebServerKnowsHandedOffReviews | cmd/legato/review_engine_handoff_test.go:21-94 | cmd/legato/review_engine_handoff_test.go:20-87 | test |
| modified | type | reviewLoopSteerer | cmd/legato/review_loop_contract.go:42-48 | cmd/legato/review_loop_contract.go:40-43 | source |
| formatting_only | function | runReviewLoopAnotherPassContract | cmd/legato/review_loop_contract.go:79-105 | cmd/legato/review_loop_contract.go:74-100 | source |
| added | function | runReviewLoopStewardContract | cmd/legato/review_loop_contract.go:110-135 |  | source |
| added | function | runReviewLoopStopWatchingContract | cmd/legato/review_loop_contract.go:138-141 |  | source |
| added | function | runReviewLoopWatchAgainContract | cmd/legato/review_loop_contract.go:144-147 |  | source |
| added | function | runReviewLoopEndContract | cmd/legato/review_loop_contract.go:150-153 |  | source |
| modified | function | reviewPassError | cmd/legato/review_pass_contract.go:120-136 | cmd/legato/review_pass_contract.go:120-130 | source |
| modified | function | TestInsideASessionLoopIsRefused | cmd/legato/review_pass_loop_test.go:97-125 | cmd/legato/review_pass_loop_test.go:97-122 | test |
| modified | function | TestOnTheHostLoopAndTaskAreExclusive | cmd/legato/review_pass_loop_test.go:128-145 | cmd/legato/review_pass_loop_test.go:125-139 | test |
| added | function | stewardedLoop | cmd/legato/review_pass_loop_test.go:149-174 |  | test |
| added | function | TestOnTheHostLoopStopsWatchingWatchesAgainAndEndsASteward | cmd/legato/review_pass_loop_test.go:179-219 |  | test |
| added | function | TestOnTheHostLoopRefusesStewardActionsToALoopWithNoSteward | cmd/legato/review_pass_loop_test.go:223-240 |  | test |
| modified | function | TestAHandOffCallsReviewOverrideRunsOnAWorkerWiredAsTheRuntimeWiresIt | cmd/legato/review_pass_override_test.go:23-134 | cmd/legato/review_pass_override_test.go:22-128 | test |
| modified | function | reviewPassConfig | cmd/legato/review_pass_runtime.go:53-64 | cmd/legato/review_pass_runtime.go:53-64 | source |
| modified | function | workerReviewPassOptions | cmd/legato/review_pass_runtime.go:70-76 | cmd/legato/review_pass_runtime.go:70-73 | source |
| modified | function | startReviewEngine | cmd/legato/review_pass_runtime.go:98-161 | cmd/legato/review_pass_runtime.go:95-155 | source |
| modified | function | wireHandoffReviewCaller | cmd/legato/review_pass_runtime.go:168-191 | cmd/legato/review_pass_runtime.go:159-175 | source |
| modified | function | reviewPassesAdvert | cmd/legato/review_passes_runtime.go:19-70 | cmd/legato/review_passes_runtime.go:18-61 | source |
| added | function | TestUsageProfileReviewerAppearsInTheLiveAdvert | cmd/legato/review_passes_runtime_test.go:13-46 |  | test |
| modified | function | attachReviewRuns | cmd/legato/review_run_runtime.go:16-41 | cmd/legato/review_run_runtime.go:14-36 | source |
| formatting_only | function | wireReviewRuns | cmd/legato/review_run_runtime.go:45-50 | cmd/legato/review_run_runtime.go:40-45 | source |
| added | function | reviewAuthors | cmd/legato/review_run_runtime.go:55-57 |  | source |
| added | function | runUsageCmd | cmd/legato/usage_cmd.go:19-52 |  | source |
| added | function | cliSubscriptionUsage | cmd/legato/usage_cmd.go:56-64 |  | source |
| added | function | refreshUsageFromCLI | cmd/legato/usage_cmd.go:68-76 |  | source |
| added | function | formatSubscriptionUsage | cmd/legato/usage_cmd.go:80-142 |  | source |
| added | function | agoText | cmd/legato/usage_cmd.go:145-150 |  | source |
| added | function | spanText | cmd/legato/usage_cmd.go:153-168 |  | source |
| added | function | runUsagePinCmd | cmd/legato/usage_cmd.go:170-216 |  | source |
| added | function | TestUsageIsPrintedPerProviderWithResetsAndFreshness | cmd/legato/usage_cmd_test.go:10-48 |  | test |
| added | function | TestUsagePrintsOpenCodeGoWithItsMonthlyMeter | cmd/legato/usage_cmd_test.go:52-78 |  | test |
| added | function | usageNewWorkDefaults | cmd/legato/usage_defaults_runtime.go:9-14 |  | source |
| added | function | usageReviewerDefaults | cmd/legato/usage_defaults_runtime.go:16-28 |  | source |
| added | type | usagePolicyHost | cmd/legato/usage_policy_runtime.go:24-28 |  | source |
| added | method | usagePolicyHost.setOwner | cmd/legato/usage_policy_runtime.go:31-35 |  | source |
| added | method | usagePolicyHost.pin | cmd/legato/usage_policy_runtime.go:36-43 |  | source |
| added | function | usageDirectory | cmd/legato/usage_policy_runtime.go:47-50 |  | source |
| added | function | attachUsagePolicy | cmd/legato/usage_policy_runtime.go:52-94 |  | source |
| added | function | reviewerPair | cmd/legato/usage_policy_runtime.go:96-101 |  | source |
| added | function | usageOSNotifier | cmd/legato/usage_policy_runtime.go:102-107 |  | source |
| added | function | requestUsagePin | cmd/legato/usage_policy_runtime.go:109-131 |  | source |
| added | function | handleUsagePinIPC | cmd/legato/usage_policy_runtime.go:133-157 |  | source |
| added | function | registerUsageAttention | cmd/legato/usage_policy_runtime.go:160-162 |  | source |
| added | function | TestUsagePinIsAcknowledgedByTheDirectoryOwner | cmd/legato/usage_policy_runtime_test.go:20-79 |  | test |
| added | function | TestUsagePinCLIControlsPolicyWithReadersOff | cmd/legato/usage_policy_runtime_test.go:81-120 |  | test |
| added | function | TestUsagePolicyRefusesAnUninstalledTierInsteadOfLaunchingBaseArgs | cmd/legato/usage_policy_runtime_test.go:122-165 |  | test |
| added | function | TestUsageCommittedDefaultsReachTheSpawnLaunchCommand | cmd/legato/usage_policy_runtime_test.go:167-210 |  | test |
| added | type | usageLaunchTmux | cmd/legato/usage_policy_runtime_test.go:212-215 |  | test |
| added | method | usageLaunchTmux.Spawn | cmd/legato/usage_policy_runtime_test.go:217-220 |  | test |
| added | method | usageLaunchTmux.SendKeysLine | cmd/legato/usage_policy_runtime_test.go:221-221 |  | test |
| added | method | usageLaunchTmux.SendKeysShellCommand | cmd/legato/usage_policy_runtime_test.go:222-222 |  | test |
| added | function | TestUsageRootAndPinCommandsKeepSeparatePermissionClasses | cmd/legato/usage_policy_runtime_test.go:224-233 |  | test |
| added | function | claimUsagePoller | cmd/legato/usage_runtime.go:32-46 |  | source |
| added | function | usageReaders | cmd/legato/usage_runtime.go:52-77 |  | source |
| added | function | usageAPIKey | cmd/legato/usage_runtime.go:84-102 |  | source |
| added | function | subscriptionUsageFor | cmd/legato/usage_runtime.go:107-125 |  | source |
| added | function | startSubscriptionUsage | cmd/legato/usage_runtime.go:135-158 |  | source |
| added | function | awaitUsagePoller | cmd/legato/usage_runtime.go:162-181 |  | source |
| added | function | pollSubscriptionUsage | cmd/legato/usage_runtime.go:185-200 |  | source |
| added | function | TestOneProcessHoldsTheUsagePoller | cmd/legato/usage_runtime_test.go:24-40 |  | test |
| added | function | TestUsageReadersFollowTheSettings | cmd/legato/usage_runtime_test.go:42-61 |  | test |
| added | function | TestOpenCodeUsageReaderIsWiredWithTheResolvedKey | cmd/legato/usage_runtime_test.go:66-131 |  | test |
| added | type | countingUsageReader | cmd/legato/usage_runtime_test.go:133-133 |  | test |
| added | method | countingUsageReader.Provider | cmd/legato/usage_runtime_test.go:135-135 |  | test |
| added | method | countingUsageReader.Read | cmd/legato/usage_runtime_test.go:136-139 |  | test |
| added | function | TestASurvivingProcessTakesOverThePoller | cmd/legato/usage_runtime_test.go:143-188 |  | test |
| added | function | TestUsageRuntimeReadsEachNamedClaudeDirectory | cmd/legato/usage_runtime_test.go:190-231 |  | test |
| modified | function | addWorkflowExample | cmd/legato/workflow_add.go:45-72 | cmd/legato/workflow_add.go:45-69 | source |
| modified | function | TestWorkflowAddWritesAnExampleIntoTheChosenFolder | cmd/legato/workflow_add_test.go:12-30 | cmd/legato/workflow_add_test.go:12-30 | test |
| modified | function | TestWorkflowAddRefusesToOverwriteWithoutForce | cmd/legato/workflow_add_test.go:32-50 | cmd/legato/workflow_add_test.go:32-50 | test |
| modified | function | TestWorkflowAddListsItsExamplesAndRefusesAnUnknownOne | cmd/legato/workflow_add_test.go:52-64 | cmd/legato/workflow_add_test.go:52-64 | test |
| added | function | TestWorkflowAddSaysReviewLoopReplacedReviewPRPasses | cmd/legato/workflow_add_test.go:68-80 |  | test |
| modified | function | TestWorkflowAddIsReachableFromTheCommandLine | cmd/legato/workflow_add_test.go:82-100 | cmd/legato/workflow_add_test.go:66-84 | test |
| modified | function | parseWorkflowArgs | cmd/legato/workflow_args.go:13-64 | cmd/legato/workflow_args.go:13-64 | source |
| modified | function | runWorkflowCmd | cmd/legato/workflow_cmd.go:50-172 | cmd/legato/workflow_cmd.go:50-169 | source |
| modified | function | resolveWorkflowDefinition | cmd/legato/workflow_cmd.go:236-247 | cmd/legato/workflow_cmd.go:233-244 | source |
| modified | function | buildWorkflowService | cmd/legato/workflow_runtime.go:70-197 | cmd/legato/workflow_runtime.go:70-194 | source |
| modified | function | commandAgentService | cmd/legato/workflow_runtime.go:259-284 | cmd/legato/workflow_runtime.go:256-279 | source |
| formatting_only | method | approvalRuntime.compactWaits | cmd/legato/workflow_runtime.go:296-305 | cmd/legato/workflow_runtime.go:291-300 | source |
| added | function | reviewLoopSteps | cmd/legato/workflow_runtime.go:310-315 |  | source |
| modified | function | workflowDiscoveryProvider | cmd/legato/workflow_trigger_runtime.go:17-25 | cmd/legato/workflow_trigger_runtime.go:17-25 | source |
| added | function | runWorkflowUpgrade | cmd/legato/workflow_upgrade.go:19-71 |  | source |
| added | function | wasVersion2 | cmd/legato/workflow_upgrade.go:74-77 |  | source |
| added | function | TestWorkflowUpgradeRewritesTheFileInPlace | cmd/legato/workflow_upgrade_test.go:45-64 |  | test |
| added | function | TestWorkflowUpgradeStdoutWritesNothing | cmd/legato/workflow_upgrade_test.go:66-79 |  | test |
| added | function | TestWorkflowUpgradeJSONSaysWhatChanged | cmd/legato/workflow_upgrade_test.go:81-99 |  | test |
| added | function | TestWorkflowUpgradeRefusesAFileThatDoesNotLoad | cmd/legato/workflow_upgrade_test.go:101-112 |  | test |
| added | function | TestWorkflowUpgradeRewritesTheInlineReviewStepOfAVersion2File | cmd/legato/workflow_upgrade_test.go:116-140 |  | test |
| added | type | ClaudeAccount | config/claude_accounts.go:12-16 |  | source |
| added | function | ValidateClaudeAccounts | config/claude_accounts.go:20-74 |  | source |
| added | function | accountPathContains | config/claude_accounts.go:76-79 |  | source |
| added | method | Config.ValidateClaudeSelection | config/claude_accounts.go:81-94 |  | source |
| added | function | UsageReaderID | config/claude_accounts.go:98-103 |  | source |
| added | function | ParseUsageMeter | config/claude_accounts.go:107-131 |  | source |
| added | method | Config.IsUsageMeter | config/claude_accounts.go:133-136 |  | source |
| added | function | TestNamedClaudeAccountProfileAndMeterLoad | config/claude_accounts_test.go:9-25 |  | test |
| added | function | TestClaudeAccountConfigurationRejectsAmbiguousSelections | config/claude_accounts_test.go:27-49 |  | test |
| modified | type | Config | config/config.go:18-88 | config/config.go:18-84 | source |
| modified | type | SwarmParticipantConfig | config/config.go:223-227 | config/config.go:219-222 | source |
| modified | type | NotificationsConfig | config/config.go:277-289 | config/config.go:272-282 | source |
| added | method | NotificationsConfig.UsageEnabled | config/config.go:303-305 |  | source |
| modified | function | validate | config/config.go:668-746 | config/config.go:656-728 | source |
| added | type | UsageConfig | config/config.go:1231-1241 |  | source |
| added | type | UsageProviderConfig | config/config.go:1244-1247 |  | source |
| added | method | UsageProviderConfig.IsEnabled | config/config.go:1250-1250 |  | source |
| added | type | UsageOpenCodeConfig | config/config.go:1254-1262 |  | source |
| added | method | UsageOpenCodeConfig.IsEnabled | config/config.go:1265-1265 |  | source |
| added | method | UsageConfig.PollIntervalDuration | config/config.go:1275-1280 |  | source |
| added | function | checkUsagePollInterval | config/config.go:1283-1296 |  | source |
| modified | type | ReeveConfig | config/config.go:1300-1337 | config/config.go:1214-1250 | source |
| modified | type | PullRequestRule | config/pull_requests.go:31-41 | config/pull_requests.go:31-37 | source |
| modified | function | normalizePullRequests | config/pull_requests.go:48-80 | config/pull_requests.go:44-73 | source |
| modified | method | PullRequestsConfig.ModeFor | config/pull_requests.go:85-90 | config/pull_requests.go:78-104 | source |
| added | method | PullRequestsConfig.ReadyAfterReview | config/pull_requests.go:95-98 |  | source |
| added | method | PullRequestsConfig.ruleFor | config/pull_requests.go:101-128 |  | source |
| added | function | TestReadyAfterReviewIsParsedOnADraftRule | config/pull_requests_test.go:90-104 |  | test |
| added | function | TestReadyAfterReviewOnARuleThatIsNotDraftIsAnError | config/pull_requests_test.go:106-113 |  | test |
| added | function | TestReadyAfterReviewOfARepository | config/pull_requests_test.go:117-134 |  | test |
| modified | type | ReviewPassesConfig | config/review_passes.go:13-35 | config/review_passes.go:13-31 | source |
| modified | type | ReviewerAgentConfig | config/review_passes.go:39-45 | config/review_passes.go:35-40 | source |
| modified | function | defaultReviewPasses | config/review_passes.go:47-52 | config/review_passes.go:42-47 | source |
| modified | function | ValidateReviewPasses | config/review_passes.go:70-144 | config/review_passes.go:65-135 | source |
| modified | function | TestReviewPassesAreOffWithTheDecidedDefaults | config/review_passes_test.go:26-43 | config/review_passes_test.go:26-43 | test |
| modified | function | TestReviewPassesParseAWorkersSettings | config/review_passes_test.go:45-73 | config/review_passes_test.go:45-72 | test |
| modified | function | TestInvalidReviewPassSettingsFailToLoad | config/review_passes_test.go:75-99 | config/review_passes_test.go:74-97 | test |
| formatting_only | function | TestReviewPassesNeedTheAgentsMCPTools | config/review_passes_test.go:110-122 | config/review_passes_test.go:108-120 | test |
| added | function | TestReviewerAccountConfigurationLoadsAndRejectsInvalidSelections | config/review_passes_test.go:124-149 |  | test |
| modified | function | pointerDefault | config/settings.go:630-640 | config/settings.go:616-626 | source |
| modified | function | TestANilSettingReportsTheDefaultTheRuntimeUses | config/settings_defaults_test.go:12-59 | config/settings_defaults_test.go:12-54 | test |
| modified | function | TestADefaultLabelIsWhatTheRuntimeUses | config/settings_layout_test.go:130-155 | config/settings_layout_test.go:130-154 | test |
| added | type | UsageProfile | config/usage_profiles.go:12-18 |  | source |
| added | type | UsageProfileReviewer | config/usage_profiles.go:21-25 |  | source |
| added | function | ValidateUsageProfiles | config/usage_profiles.go:28-83 |  | source |
| added | function | IsUsageMeter | config/usage_profiles.go:87-90 |  | source |
| added | function | usageName | config/usage_profiles.go:92-102 |  | source |
| added | function | usageAgent | config/usage_profiles.go:104-116 |  | source |
| added | method | UsageProfile.UnmarshalYAML | config/usage_profiles.go:120-126 |  | source |
| added | method | UsageProfileReviewer.UnmarshalYAML | config/usage_profiles.go:128-134 |  | source |
| added | function | rejectUsageModel | config/usage_profiles.go:136-145 |  | source |
| added | type | UsageRule | config/usage_profiles.go:149-155 |  | source |
| added | function | TestUsageProfilesLoadWithTiersAndReviewer | config/usage_profiles_test.go:9-47 |  | test |
| added | function | TestUsageProfilesRejectInvalidDefaults | config/usage_profiles_test.go:49-67 |  | test |
| added | function | TestUsageProfilesRejectRawModels | config/usage_profiles_test.go:69-78 |  | test |
| added | function | TestUsageRulesLoadQualifiedMetersAndExplicitOrder | config/usage_profiles_test.go:80-106 |  | test |
| added | function | TestUsageRulesRejectInvalidPolicy | config/usage_profiles_test.go:108-129 |  | test |
| added | function | TestUsagePolicyOverlaysProfilesRulesAndFallbackTogether | config/usage_profiles_test.go:131-159 |  | test |
| added | function | TestUsagePolicyOverlaysRejectInvalidEntriesAndReferences | config/usage_profiles_test.go:161-189 |  | test |
| added | function | TestUsageNotificationsAndSettingsAreLive | config/usage_profiles_test.go:191-233 |  | test |
| added | type | UsageMeterInfo | config/usage_providers.go:14-21 |  | source |
| added | type | UsageProviderInfo | config/usage_providers.go:24-33 |  | source |
| added | function | UsageProviders | config/usage_providers.go:53-53 |  | source |
| added | function | UsageMeterInfoOf | config/usage_providers.go:56-63 |  | source |
| added | function | UsageQualifiedMeters | config/usage_providers.go:67-75 |  | source |
| added | function | TestUsageProvidersDeclareIdentityAndMeters | config/usage_providers_test.go:12-50 |  | test |
| added | function | TestUsageMetersCarryTheirOwnLabels | config/usage_providers_test.go:52-72 |  | test |
| added | function | TestOpenCodeUsageIsOffUntilItIsTurnedOnAndNamesAKeySecret | config/usage_providers_test.go:76-117 |  | test |
| added | function | TestSubscriptionUsageReadsBothProvidersEveryFiveMinutesByDefault | config/usage_test.go:9-22 |  | test |
| added | function | TestSubscriptionUsageSettingsAreReadFromTheFile | config/usage_test.go:24-32 |  | test |
| added | function | TestAPollIntervalUnderTwoMinutesIsRefused | config/usage_test.go:34-45 |  | test |
| modified | type | draftWorkflowInput | internal/cli/mcp_reeve_workflow.go:14-17 | internal/cli/mcp_reeve_workflow.go:14-17 | source |
| modified | function | TestTheReeveDraftsAWorkflowAndProposesItPaused | internal/cli/mcp_reeve_workflow_test.go:21-56 | internal/cli/mcp_reeve_workflow_test.go:21-52 | test |
| modified | function | WorkflowDefinitions | internal/cli/workflow.go:16-18 | internal/cli/workflow.go:16-18 | source |
| added | method | Client.ClaimTask | internal/engine/burnlite/assign.go:19-22 |  | source |
| added | method | Client.ReleaseClaim | internal/engine/burnlite/assign.go:26-29 |  | source |
| added | method | Client.SetStage | internal/engine/burnlite/assign.go:33-36 |  | source |
| added | function | TestClaimTaskClaimsUnderAHolder | internal/engine/burnlite/assign_test.go:12-25 |  | test |
| added | function | TestClaimTaskHeldByAnotherIsItsError | internal/engine/burnlite/assign_test.go:27-39 |  | test |
| added | function | TestReleaseClaimReleasesTheHoldersClaim | internal/engine/burnlite/assign_test.go:41-52 |  | test |
| added | function | TestSetStageUpdatesTheTasksStage | internal/engine/burnlite/assign_test.go:54-67 |  | test |
| added | function | TestGetTaskReadsTheStage | internal/engine/burnlite/assign_test.go:69-83 |  | test |
| modified | type | Task | internal/engine/burnlite/client.go:131-154 | internal/engine/burnlite/client.go:131-152 | source |
| modified | method | Client.GetTask | internal/engine/burnlite/client.go:377-429 | internal/engine/burnlite/client.go:375-426 | source |
| modified | method | fakeBurnlite.server | internal/engine/burnlite/client_test.go:28-62 | internal/engine/burnlite/client_test.go:28-62 | test |
| added | type | RateLimitedError | internal/engine/claudeusage/claudeusage.go:67-69 |  | source |
| added | method | RateLimitedError.Error | internal/engine/claudeusage/claudeusage.go:71-76 |  | source |
| added | type | Window | internal/engine/claudeusage/claudeusage.go:79-86 |  | source |
| added | type | Reading | internal/engine/claudeusage/claudeusage.go:89-91 |  | source |
| added | type | Login | internal/engine/claudeusage/claudeusage.go:94-98 |  | source |
| added | type | CredentialSource | internal/engine/claudeusage/claudeusage.go:102-104 |  | source |

… 1636 more in JSON

## Diagnostics

- **info** `unsupported_language`  browser-extension/src/actions.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  browser-extension/src/background.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  browser-extension/src/background.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  browser-extension/src/config.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  browser-extension/src/legato.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  browser-extension/src/legato.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  internal/engine/store/migrations/155_workflow_runs_by_name.sql changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  internal/engine/store/migrations/156_subscription_usage.sql changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  internal/engine/store/migrations/157_review_loop_approval_dismissed.sql changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  internal/engine/store/migrations/158_usage_profiles.sql changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  internal/engine/store/migrations/159_usage_work_identity.sql changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  internal/engine/store/migrations/160_execution_accounts.sql changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  internal/engine/store/migrations/161_handoff_accounts.sql changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  mobile/app/src/main/java/dev/legato/mobile/GlassesHttp.kt changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  mobile/app/src/test/java/dev/legato/mobile/GlassesHttpTest.kt changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/dev/legato-dev.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/dev/legato-dev.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/scripts/screenshots.mjs changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/agent.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/agent.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/app.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/app.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/approvals.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/glance.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/glance.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/layout.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/legato.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/legato.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/main.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/questions.fixtures.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/questions.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/questions.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/settings.test.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  packages/even-hub/src/settings.ts changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  scripts/verify-watch-disclosures.mjs changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/AgentChatHeader.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/AgentInfoDialog.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/AgentKindBadge.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/AgentSidebar.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/MacroForm.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/ReviewLoopFacts.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/RunWorkflowDialog.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/SpawnAgentDialog.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/SpawnEphemeralModal.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/UsageStrip.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/WorkflowTriggerSwitch.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/__tests__/AgentInfoDialog.test.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/__tests__/AgentSidebar.test.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/__tests__/ListSettings.test.tsx changed file in a language without an adapter is not analysed
- **info** `unsupported_language`  web/src/components/__tests__/MacroFormAccount.test.tsx changed file in a language without an adapter is not analysed
- … 132 more in JSON

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
