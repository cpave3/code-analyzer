# Debt factors: legato-stack @ `bb13a424dc73` — internal/engine/store

- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.2.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 1555 of 1848 production source files measured (90.0% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2; 4 file(s) with parse errors

**Scope:** 89 source files, 20004 lines, 953 callables (62 hotspots) in 1 module(s).

| Debt | complexity | nesting | duplication | architecture | Debt/kLOC | Score |
|---:|---:|---:|---:|---:|---:|---:|
| **833.90** | 403.00 | 14.00 | 416.90 | 0.00 | 41.69 | 37.49 |

**Concentration:** 134 factor(s) in scope; the largest 21 carry 50% of the debt and the largest 63 carry 80%.

## Factors (20 listed, 114 smaller omitted; largest first)

### 1. complexity `prepareMigration` — 84.00 debt · 10.1% of scope (cumulative 10.1%)

- **Where:** `internal/engine/store/store.go:364-683` (module `internal/engine/store`, go)
- **Why:** cyclomatic 93 is 84 point(s) at or above the hotspot threshold 10 (92 decision points, 320 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 93) removes this factor; each point below 93 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `prepareMigration` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 2. duplication `Store.ClosedLoopReviewers` — 63.80 debt · 7.6% of scope (cumulative 17.7%)

- **Where:** `internal/engine/store/review_loops.go:1115-1123` (module `internal/engine/store`, go)
- **Why:** structural clone group of 30 callables (30 in scope); 1276 redundant tokens attributed to this scope
- **Related:** `Store.ClosedLoopReviewers` `internal/engine/store/review_loops.go:1115-1123`; `Store.ConversationStorageDirs` `internal/engine/store/agent_conversation_recovery.go:41-46`; `Store.ListChosenExecutionResources` `internal/engine/store/resource.go:148-153`; `Store.ListExecutionResources` `internal/engine/store/resource.go:157-162`; `Store.ListOAuthAttempts` `internal/engine/store/oauth_attempts.go:188-192`; `Store.ListOpenIntegrationAttention` `internal/engine/store/integration_attention.go:76-81`; `Store.ListPlanComments` `internal/engine/store/plan.go:183-187`; `Store.ListPlanMessages` `internal/engine/store/plan.go:331-335`; `Store.ListPlanQuestions` `internal/engine/store/plan.go:106-110`; `Store.ListPlanResponses` `internal/engine/store/plan.go:168-172`; `Store.ListPlanRevisions` `internal/engine/store/plan.go:306-310`; `Store.ListPlansByTask` `internal/engine/store/plan.go:21-25`; `Store.ListReeveWatchEvents` `internal/engine/store/reeve_watches.go:55-59`; `Store.ListReviewChapterHunks` `internal/engine/store/review.go:453-458`; `Store.ListReviewFindingsByPass` `internal/engine/store/review.go:123-127`; `Store.ListReviewHunkNotesByPass` `internal/engine/store/review.go:108-112`; `Store.ListReviewLoopFindings` `internal/engine/store/review_loops.go:993-997`; `Store.ListReviewLoopItems` `internal/engine/store/review_loops.go:873-877`; `Store.ListReviewLoopPRs` `internal/engine/store/review_loops.go:937-941`; `Store.ListReviewLoopPasses` `internal/engine/store/review_loops.go:986-990`; `Store.ListReviewMessagesByPass` `internal/engine/store/review.go:102-106`; `Store.ListReviewPasses` `internal/engine/store/review.go:79-83`; `Store.ListReviewStepsByPass` `internal/engine/store/review.go:94-100`; `Store.ListReviewToursByTask` `internal/engine/store/review.go:64-69`; `Store.ListRunReviewLoopsAuthoredBy` `internal/engine/store/handoff_reviews.go:1297-1302`; `Store.ListSubtasksByParent` `internal/engine/store/swarm.go:70-75`; `Store.ListSyncLogs` `internal/engine/store/store.go:1073-1078`; `Store.ListTaskNotes` `internal/engine/store/task_notes.go:10-14`; `Store.ListUnackedSwarmEvents` `internal/engine/store/swarm.go:245-252`; `Store.OpenSubagents` `internal/engine/store/subagent_activity.go:123-128`
- **Fix:** consolidate duplicated callables. Target: one implementation instead of 30 copies.
  - confirm the copies have the same semantics (structural clones may differ in literals or names that matter)
  - move one implementation to a module every caller may import under the configured architecture rules
  - replace the other copies with calls; delete the copies
- **Done when:** the clone group shrinks to one member (or is removed); no new architecture rule violations; tests for every caller pass
- **Caution:** structural similarity is a candidate signal, not proof of equivalence; if semantics differ or layering forbids a shared home, record a won't-fix with the reason

### 3. duplication `Store.GetExecution` — 26.00 debt · 3.1% of scope (cumulative 20.8%)

- **Where:** `internal/engine/store/execution.go:70-88` (module `internal/engine/store`, go)
- **Why:** structural clone group of 9 callables (9 in scope); 520 redundant tokens attributed to this scope
- **Related:** `Store.GetExecution` `internal/engine/store/execution.go:70-88`; `Store.GetLearningCurationGeneration` `internal/engine/store/learning_curation.go:163-172`; `Store.GetLearningCurationGenerationByAttempt` `internal/engine/store/learning_curation.go:54-63`; `Store.GetLearningIncorporation` `internal/engine/store/learning.go:415-424`; `Store.GetLearningPromptManifest` `internal/engine/store/learning_prompt.go:74-83`; `Store.GetLearningReflectionRequest` `internal/engine/store/learning.go:38-47`; `Store.GetLearningReflectionRequestByAttempt` `internal/engine/store/learning.go:49-58`; `Store.GetTaskTransition` `internal/engine/store/workflow_trigger.go:86-95`; `Store.LatestTaskTransition` `internal/engine/store/workflow_trigger.go:66-75`
- **Fix:** consolidate duplicated callables. Target: one implementation instead of 9 copies.
  - confirm the copies have the same semantics (structural clones may differ in literals or names that matter)
  - move one implementation to a module every caller may import under the configured architecture rules
  - replace the other copies with calls; delete the copies
- **Done when:** the clone group shrinks to one member (or is removed); no new architecture rule violations; tests for every caller pass
- **Caution:** structural similarity is a candidate signal, not proof of equivalence; if semantics differ or layering forbids a shared home, record a won't-fix with the reason

### 4. duplication `Store.GetActiveReviewPass` — 20.65 debt · 2.5% of scope (cumulative 23.3%)

- **Where:** `internal/engine/store/review.go:85-92` (module `internal/engine/store`, go)
- **Why:** structural clone group of 8 callables (8 in scope); 413 redundant tokens attributed to this scope
- **Related:** `Store.GetActiveReviewPass` `internal/engine/store/review.go:85-92`; `Store.GetDefaultReviewTour` `internal/engine/store/review.go:54-61`; `Store.GetReeveChangeSet` `internal/engine/store/reeve_change_sets.go:47-54`; `Store.GetReevePeer` `internal/engine/store/reeve_peers.go:35-42`; `Store.GetReviewPassPlan` `internal/engine/store/review.go:184-196`; `Store.GetReviewStep` `internal/engine/store/review.go:286-293`; `Store.GetReviewTour` `internal/engine/store/review.go:44-51`; `Store.HandoffObservation` `internal/engine/store/handoff.go:200-207`
- **Fix:** consolidate duplicated callables. Target: one implementation instead of 8 copies.
  - confirm the copies have the same semantics (structural clones may differ in literals or names that matter)
  - move one implementation to a module every caller may import under the configured architecture rules
  - replace the other copies with calls; delete the copies
- **Done when:** the clone group shrinks to one member (or is removed); no new architecture rule violations; tests for every caller pass
- **Caution:** structural similarity is a candidate signal, not proof of equivalence; if semantics differ or layering forbids a shared home, record a won't-fix with the reason

### 5. duplication `Store.GitHubAppLogins` — 20.50 debt · 2.5% of scope (cumulative 25.8%)

- **Where:** `internal/engine/store/pr_reviews_seen.go:64-68` (module `internal/engine/store`, go)
- **Why:** structural clone group of 11 callables (11 in scope); 410 redundant tokens attributed to this scope
- **Related:** `Store.GitHubAppLogins` `internal/engine/store/pr_reviews_seen.go:64-68`; `Store.ListLaunchMessages` `internal/engine/store/launch_messages.go:37-42`; `Store.ListReeveInbox` `internal/engine/store/reeve_peers.go:105-109`; `Store.ListReevePeers` `internal/engine/store/reeve_peers.go:45-49`; `Store.ListRepositoryGuidance` `internal/engine/store/repository_guidance.go:21-25`; `Store.ListRepositoryPresets` `internal/engine/store/repository_preset.go:32-36`; `Store.ListSecrets` `internal/engine/store/credential.go:57-61`; `Store.ListSettingEntries` `internal/engine/store/settings.go:75-79`; `Store.ListSettings` `internal/engine/store/settings.go:19-23`; `Store.ListSubscriptionUsage` `internal/engine/store/subscription_usage.go:14-18`; `Store.WorkflowTriggerSwitches` `internal/engine/store/workflow_trigger_switches.go:14-18`
- **Fix:** consolidate duplicated callables. Target: one implementation instead of 11 copies.
  - confirm the copies have the same semantics (structural clones may differ in literals or names that matter)
  - move one implementation to a module every caller may import under the configured architecture rules
  - replace the other copies with calls; delete the copies
- **Done when:** the clone group shrinks to one member (or is removed); no new architecture rule violations; tests for every caller pass
- **Caution:** structural similarity is a candidate signal, not proof of equivalence; if semantics differ or layering forbids a shared home, record a won't-fix with the reason

### 6. complexity `Store.migrate` — 19.00 debt · 2.3% of scope (cumulative 28.1%)

- **Where:** `internal/engine/store/store.go:225-317` (module `internal/engine/store`, go)
- **Why:** cyclomatic 28 is 19 point(s) at or above the hotspot threshold 10 (27 decision points, 93 lines)

### 7. complexity `openReviewLoopTx` — 18.00 debt · 2.2% of scope (cumulative 30.2%)

- **Where:** `internal/engine/store/review_loops.go:438-507` (module `internal/engine/store`, go)
- **Why:** cyclomatic 27 is 18 point(s) at or above the hotspot threshold 10 (26 decision points, 70 lines)

### 8. complexity `Store.SyncColumnMappings` — 15.00 debt · 1.8% of scope (cumulative 32.0%)

- **Where:** `internal/engine/store/columns.go:23-127` (module `internal/engine/store`, go)
- **Why:** cyclomatic 24 is 15 point(s) at or above the hotspot threshold 10 (23 decision points, 105 lines)

### 9. complexity `Store.ApplyLearningCuration` — 15.00 debt · 1.8% of scope (cumulative 33.8%)

- **Where:** `internal/engine/store/learning_curation.go:177-252` (module `internal/engine/store`, go)
- **Why:** cyclomatic 24 is 15 point(s) at or above the hotspot threshold 10 (23 decision points, 76 lines)

### 10. complexity `Store.RecordTaskNativeSession` — 15.00 debt · 1.8% of scope (cumulative 35.6%)

- **Where:** `internal/engine/store/task_native_session.go:183-252` (module `internal/engine/store`, go)
- **Why:** cyclomatic 24 is 15 point(s) at or above the hotspot threshold 10 (23 decision points, 70 lines)

### 11. complexity `Store.ApplyLearningProposalDecision` — 14.00 debt · 1.7% of scope (cumulative 37.3%)

- **Where:** `internal/engine/store/learning_curation.go:77-161` (module `internal/engine/store`, go)
- **Why:** cyclomatic 23 is 14 point(s) at or above the hotspot threshold 10 (22 decision points, 85 lines)

### 12. complexity `Store.FailWorkflowStepWith` — 14.00 debt · 1.7% of scope (cumulative 39.0%)

- **Where:** `internal/engine/store/workflow.go:759-827` (module `internal/engine/store`, go)
- **Why:** cyclomatic 23 is 14 point(s) at or above the hotspot threshold 10 (22 decision points, 69 lines)

### 13. complexity `Store.UpdateMCPConnection` — 13.00 debt · 1.6% of scope (cumulative 40.5%)

- **Where:** `internal/engine/store/mcp_connections.go:114-161` (module `internal/engine/store`, go)
- **Why:** cyclomatic 22 is 13 point(s) at or above the hotspot threshold 10 (21 decision points, 48 lines)

### 14. duplication `Store.DeleteRepositoryGuidance` — 13.00 debt · 1.6% of scope (cumulative 42.1%)

- **Where:** `internal/engine/store/repository_guidance.go:28-39` (module `internal/engine/store`, go)
- **Why:** structural clone group of 5 callables (5 in scope); 260 redundant tokens attributed to this scope
- **Related:** `Store.DeleteRepositoryGuidance` `internal/engine/store/repository_guidance.go:28-39`; `Store.DeleteRepositoryPreset` `internal/engine/store/repository_preset.go:39-50`; `Store.RevokeIntegrationToken` `internal/engine/store/integration_tokens.go:156-167`; `Store.RevokeMainIntegrationToken` `internal/engine/store/integration_tokens.go:113-124`; `Store.SetExecutionReadOnly` `internal/engine/store/execution_broker.go:38-49`

### 15. complexity `Store.CheckReviewLoopWrite` — 11.00 debt · 1.3% of scope (cumulative 43.4%)

- **Where:** `internal/engine/store/handoff_reviews.go:540-604` (module `internal/engine/store`, go)
- **Why:** cyclomatic 20 is 11 point(s) at or above the hotspot threshold 10 (19 decision points, 65 lines)

### 16. complexity `Store.HoldHandoffReviewForRetry` — 11.00 debt · 1.3% of scope (cumulative 44.7%)

- **Where:** `internal/engine/store/handoff_reviews.go:978-1026` (module `internal/engine/store`, go)
- **Why:** cyclomatic 20 is 11 point(s) at or above the hotspot threshold 10 (19 decision points, 49 lines)

### 17. complexity `Store.FailRecordedWorkflowStep` — 11.00 debt · 1.3% of scope (cumulative 46.0%)

- **Where:** `internal/engine/store/workflow_wait.go:206-285` (module `internal/engine/store`, go)
- **Why:** cyclomatic 20 is 11 point(s) at or above the hotspot threshold 10 (19 decision points, 80 lines)

### 18. duplication `Store.GetApprovalDevice` — 10.50 debt · 1.3% of scope (cumulative 47.3%)

- **Where:** `internal/engine/store/approval_devices.go:53-62` (module `internal/engine/store`, go)
- **Why:** structural clone group of 4 callables (4 in scope); 210 redundant tokens attributed to this scope
- **Related:** `Store.GetApprovalDevice` `internal/engine/store/approval_devices.go:53-62`; `Store.GetAttentionItem` `internal/engine/store/attention.go:65-74`; `Store.GetPlan` `internal/engine/store/plan.go:10-19`; `Store.GetSecret` `internal/engine/store/credential.go:34-43`

### 19. duplication `Store.CountApprovals` — 10.00 debt · 1.2% of scope (cumulative 48.5%)

- **Where:** `internal/engine/store/approval.go:178-182` (module `internal/engine/store`, go)
- **Why:** structural clone group of 6 callables (6 in scope); 200 redundant tokens attributed to this scope
- **Related:** `Store.CountApprovals` `internal/engine/store/approval.go:178-182`; `Store.CountUnansweredRequiredPlanQuestions` `internal/engine/store/plan.go:139-146`; `Store.HasHandoffPackets` `internal/engine/store/handoff.go:238-242`; `Store.HasReeveWatchEventKind` `internal/engine/store/reeve_watches.go:106-110`; `Store.TaskActivityMark` `internal/engine/store/reeve_watches.go:96-102`; `Store.WorkflowAttemptInUse` `internal/engine/store/workflow.go:1114-1121`

### 20. duplication `Store.ClaimHandoffNotice` — 9.15 debt · 1.1% of scope (cumulative 49.6%)

- **Where:** `internal/engine/store/handoff_notices.go:10-17` (module `internal/engine/store`, go)
- **Why:** structural clone group of 4 callables (4 in scope); 183 redundant tokens attributed to this scope
- **Related:** `Store.ClaimHandoffNotice` `internal/engine/store/handoff_notices.go:10-17`; `Store.ClaimInstanceValue` `internal/engine/store/approval_devices.go:221-228`; `Store.SetWorkerRefused` `internal/engine/store/workers.go:103-111`; `Store.SetWorkerUnreachable` `internal/engine/store/workers.go:91-98`

## Guidance for fixing

- preserve behaviour: add or confirm characterisation tests for a factor's code before changing it; the evaluator does not run code
- reduce decision points, do not only move them: a split that leaves the same branching in helpers removes hotspot debt but not logic (check the decision-point ledger with `codeeval check`)
- keep helpers cohesive and named for one purpose; avoid creating new hotspot helpers (cc ≥ threshold) while removing a root's debt
- do not delete features, inline configuration into data tables, or reshape syntax solely to lower the number; debt is an uncalibrated structural proxy
- verify each fix by re-running `codeeval explain` on the same scope (the factor should disappear or shrink) and `codeeval check` on the change (gates pass, no new factors)

## Notes

- experimental scoring/v2: separate from and not comparable to scoring/v1 −2..+2 scores; v1 reports are unchanged
- debt points = complexity (Σ max(0, cyclomatic − H + 1): debt exactly when a callable is a hotspot) + nesting (Σ max(0, depth − N + 1) × weight) + duplication (structural-clone redundant tokens ÷ tokens-per-point) + architecture (configured-rule-violating edges and units in cycles × weights)
- debt is additive across modules and telescopes across commits (Δ = after − before of one snapshot function); telescoping is mathematical consistency under a fixed model and config, not empirical calibration
- debt is an uncalibrated, size-sensitive structural proxy: deleting features or code, or reshaping syntax (e.g. splitting functions below the thresholds), can lower it without improving engineering; do not optimise it alone — keep the component and raw evidence, behaviour/test, coverage and scope checks
- score = 100 × h ÷ (h + debt per 1,000 supported source lines), h = halfScoreDensity: a size-normalised view for comparing modules or codebases; it can be diluted by adding simple code and is never a target on its own
- weights are uncalibrated defaults; the score is a structural index, not a measure of correctness, design quality or behaviour
- only languages with adapters are measured; see coverage — unsupported source is not scored and not neutral
- architecture debt reflects only the configured import-level rules; rules the configuration cannot express are not measured
- complexity growth below the hotspot thresholds carries no debt; it is reported as decision points and checked by the existing-callable branch gate
- gates are advisory checks, not merge decisions
