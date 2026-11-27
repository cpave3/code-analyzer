# CodeEval v2 snapshot: legato-stack @ `bb13a424dc73`

- **Revision:** `bb13a424dc73f3eeb98f8bce43db932a9d19f7e6` tree `25b725950d6a` 2026-11-27T20:29:04+11:00
- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.2.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 1555 of 1848 production source files measured (90.0% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2; 4 file(s) with parse errors

## Codebase

| Score (0–100) | Debt | Debt / kLOC | Source lines | Callables |
|---:|---:|---:|---:|---:|
| **28.80** | **19590.65** | 61.82 | 316900 | 12006 |

| Component | Raw | Debt |
|---|---|---:|
| complexity | cyclomatic excess 17487 (hotspots 1773; decision points 54616, cyclomatic total 66622) | 17487.00 |
| nesting | nesting excess 427 | 854.00 |
| duplication | redundant tokens 17593.00 (structural redundant 288, exact 18, cloned callables 470) | 879.65 |
| architecture | rule-violating edges 11, units in cycles 13 (internal edges 278) | 370.00 |

## By language

| Language | Score | Debt | Debt/kLOC | Lines | Callables | Hotspots | Complexity | Nesting | Duplication | Architecture |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| go | 34.40 | 11387.20 | 47.67 | 238865 | 9336 | 1285 | 9888.00 | 612.00 | 777.20 | 110.00 |
| typescript | 19.21 | 8203.45 | 105.13 | 78035 | 2670 | 488 | 7599.00 | 242.00 | 102.45 | 260.00 |

Complexity is measured per language adapter; cyclomatic counts are not directly comparable across languages (see adapter capability notes).

## Modules (top 25 of 120 by debt)

| Module | Score | Debt | Debt/kLOC | Lines | Complexity | Nesting | Duplication | Architecture |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `internal/service` | 32.82 | 4932.45 | 51.17 | 96395 | 4618.00 | 242.00 | 72.45 | 0.00 |
| `web/src/components` | 18.62 | 2178.25 | 109.28 | 19933 | 2105.00 | 50.00 | 3.25 | 20.00 |
| `web/src/pages` | 10.90 | 1595.00 | 204.43 | 7802 | 1577.00 | 18.00 | 0.00 | 0.00 |
| `web/src/lib` | 26.34 | 1217.63 | 69.91 | 17416 | 1068.00 | 46.00 | 83.63 | 20.00 |
| `cmd/legato` | 34.75 | 871.93 | 46.94 | 18575 | 809.00 | 26.00 | 36.93 | 0.00 |
| `internal/engine/store` | 37.49 | 833.90 | 41.69 | 20004 | 403.00 | 14.00 | 416.90 | 0.00 |
| `internal/engine/workflow` | 29.35 | 582.00 | 60.19 | 9670 | 538.00 | 42.00 | 2.00 | 0.00 |
| `internal/server` | 48.77 | 522.80 | 26.26 | 19905 | 389.00 | 66.00 | 67.80 | 0.00 |
| `web/src/components/workflow-editor` | 13.30 | 504.00 | 163.00 | 3092 | 476.00 | 24.00 | 4.00 | 0.00 |
| `internal/tui` | 21.94 | 428.20 | 88.93 | 4815 | 372.00 | 18.00 | 8.20 | 30.00 |
| `internal/cli` | 33.04 | 426.75 | 50.67 | 8422 | 376.00 | 16.00 | 34.75 | 0.00 |
| `web/src/components/board` | 21.62 | 387.00 | 90.63 | 4270 | 381.00 | 6.00 | 0.00 | 0.00 |
| `web/src/components/settings` | 17.60 | 372.00 | 117.02 | 3179 | 344.00 | 8.00 | 0.00 | 20.00 |
| `internal/tui/overlay` | 28.51 | 361.20 | 62.68 | 5763 | 265.00 | 50.00 | 16.20 | 30.00 |
| `packages/even-hub/src` | 18.03 | 355.28 | 113.65 | 3126 | 352.00 | 2.00 | 1.27 | 0.00 |
| `web/src/hooks` | 23.43 | 288.38 | 81.69 | 3530 | 226.00 | 34.00 | 8.38 | 20.00 |
| `browser-extension/src/content` | 18.35 | 251.00 | 111.21 | 2257 | 241.00 | 10.00 | 0.00 | 0.00 |
| `config` | 35.48 | 246.10 | 45.46 | 5414 | 232.00 | 8.00 | 6.10 | 0.00 |
| `internal/engine/dockerproxy` | 17.25 | 242.30 | 119.95 | 2020 | 230.00 | 8.00 | 4.30 | 0.00 |
| `internal/engine/credential` | 30.44 | 225.00 | 57.14 | 3938 | 223.00 | 2.00 | 0.00 | 0.00 |
| `browser-extension/src` | 33.73 | 167.00 | 49.12 | 3400 | 137.00 | 10.00 | 0.00 | 20.00 |
| `internal/tui/agents` | 33.16 | 148.15 | 50.39 | 2940 | 142.00 | 4.00 | 2.15 | 0.00 |
| `internal/engine/codexbridge` | 34.48 | 140.48 | 47.51 | 2957 | 117.00 | 6.00 | 17.48 | 0.00 |
| `desktop/src/components` | 22.11 | 126.00 | 88.05 | 1431 | 98.00 | 8.00 | 0.00 | 20.00 |
| `internal/engine/claudebridge` | 40.81 | 116.83 | 36.26 | 3222 | 91.00 | 12.00 | 13.83 | 0.00 |

## Largest callable debts

| Callable | Location | cc | nesting | lines | debt |
|---|---|---:|---:|---:|---:|
| `NativeAgentPanel` | web/src/components/NativeAgentPanel.tsx:270-1453 | 331 | 5 | 1184 | 326.00 |
| `startGlance` | packages/even-hub/src/app.ts:538-1507 | 320 | 4 | 970 | 313.00 |
| `App.update` | internal/tui/app.go:696-1827 | 308 | 6 | 1132 | 305.00 |
| `BoardPage` | web/src/pages/Board.tsx:190-1424 | 297 | 5 | 1235 | 292.00 |
| `agentService.SpawnAgent` | internal/service/agent.go:837-1693 | 287 | 5 | 857 | 282.00 |
| `SpawnAgentDialog` | web/src/components/SpawnAgentDialog.tsx:104-1050 | 250 | 5 | 947 | 245.00 |
| `AgentsPage` | web/src/pages/Agents.tsx:113-1128 | 247 | 5 | 1016 | 242.00 |
| `Outline` | web/src/components/workflow-editor/WorkflowOutline.tsx:194-647 | 204 | 6 | 454 | 201.00 |
| `WorkflowEditPage` | web/src/pages/WorkflowEditPage.tsx:166-595 | 193 | 4 | 430 | 186.00 |
| `createScopeOverlay` | browser-extension/src/content/scope-overlay.ts:135-665 | 116 | 5 | 531 | 111.00 |
| `ReviewTourPage` | web/src/pages/ReviewTourPage.tsx:53-477 | 109 | 4 | 425 | 102.00 |
| `BuildSandboxPlan` | internal/service/sandbox_plan.go:414-742 | 96 | 4 | 329 | 89.00 |
| `reduceNativeEvent` | web/src/lib/nativeAgentSync.ts:573-797 | 96 | 4 | 225 | 89.00 |
| `runReviewCmd` | cmd/legato/main.go:347-703 | 93 | 4 | 357 | 86.00 |
| `prepareMigration` | internal/engine/store/store.go:364-683 | 93 | 4 | 320 | 86.00 |
| `SettingsPage` | web/src/pages/Settings.tsx:102-343 | 90 | 4 | 242 | 83.00 |
| `runNeedsYou` | web/src/lib/runNeedsYou.ts:73-165 | 87 | 3 | 93 | 78.00 |
| `runLearnCmd` | cmd/legato/learning_cmd.go:69-349 | 81 | 3 | 281 | 72.00 |
| `AgentArtifactPanel` | web/src/components/AgentArtifactPanel.tsx:69-430 | 79 | 3 | 362 | 70.00 |
| `QuestionFocus` | web/src/components/QuestionFocus.tsx:25-292 | 77 | 3 | 268 | 68.00 |
| `stateFromNativeSnapshot` | web/src/lib/nativeAgentSync.ts:327-456 | 75 | 4 | 130 | 68.00 |
| `LearningDetail` | web/src/pages/LearningDetailPage.tsx:18-146 | 77 | 3 | 129 | 68.00 |
| `WorkflowRunPage` | web/src/pages/WorkflowRunPage.tsx:694-839 | 77 | 2 | 146 | 68.00 |
| `useMenuLayer` | web/src/hooks/useMenuLayer.ts:19-152 | 72 | 4 | 134 | 65.00 |
| `runTUI` | cmd/legato/main.go:1232-1806 | 67 | 4 | 575 | 60.00 |
| `TaskHandoffService.HandOff` | internal/service/task_handoff.go:203-396 | 69 | 2 | 194 | 60.00 |
| `App` | desktop/src/App.tsx:27-301 | 65 | 4 | 275 | 58.00 |
| `ToolCallCard` | web/src/components/NativeAgentPanel.tsx:1686-1914 | 65 | 3 | 229 | 56.00 |
| `Connection` | web/src/components/settings/MCPSettings.tsx:159-424 | 62 | 4 | 266 | 55.00 |
| `WithField` | web/src/components/workflow-editor/StepsEditor.tsx:186-258 | 64 | 3 | 73 | 55.00 |

## Largest clone groups

| Kind | Redundant tokens | Debt | Members |
|---|---:|---:|---|
| structural | 1276.00 | 63.80 | `Store.ClosedLoopReviewers` internal/engine/store/review_loops.go:1115; `Store.ConversationStorageDirs` internal/engine/store/agent_conversation_recovery.go:41; `Store.ListChosenExecutionResources` internal/engine/store/resource.go:148; `Store.ListExecutionResources` internal/engine/store/resource.go:157; `Store.ListOAuthAttempts` internal/engine/store/oauth_attempts.go:188; `Store.ListOpenIntegrationAttention` internal/engine/store/integration_attention.go:76; `Store.ListPlanComments` internal/engine/store/plan.go:183; `Store.ListPlanMessages` internal/engine/store/plan.go:331; `Store.ListPlanQuestions` internal/engine/store/plan.go:106; `Store.ListPlanResponses` internal/engine/store/plan.go:168; `Store.ListPlanRevisions` internal/engine/store/plan.go:306; `Store.ListPlansByTask` internal/engine/store/plan.go:21; `Store.ListReeveWatchEvents` internal/engine/store/reeve_watches.go:55; `Store.ListReviewChapterHunks` internal/engine/store/review.go:453; `Store.ListReviewFindingsByPass` internal/engine/store/review.go:123; `Store.ListReviewHunkNotesByPass` internal/engine/store/review.go:108; `Store.ListReviewLoopFindings` internal/engine/store/review_loops.go:993; `Store.ListReviewLoopItems` internal/engine/store/review_loops.go:873; `Store.ListReviewLoopPRs` internal/engine/store/review_loops.go:937; `Store.ListReviewLoopPasses` internal/engine/store/review_loops.go:986; `Store.ListReviewMessagesByPass` internal/engine/store/review.go:102; `Store.ListReviewPasses` internal/engine/store/review.go:79; `Store.ListReviewStepsByPass` internal/engine/store/review.go:94; `Store.ListReviewToursByTask` internal/engine/store/review.go:64; `Store.ListRunReviewLoopsAuthoredBy` internal/engine/store/handoff_reviews.go:1297; `Store.ListSubtasksByParent` internal/engine/store/swarm.go:70; `Store.ListSyncLogs` internal/engine/store/store.go:1073; `Store.ListTaskNotes` internal/engine/store/task_notes.go:10; `Store.ListUnackedSwarmEvents` internal/engine/store/swarm.go:245; `Store.OpenSubagents` internal/engine/store/subagent_activity.go:123 |
| structural | 520.00 | 26.00 | `Store.GetExecution` internal/engine/store/execution.go:70; `Store.GetLearningCurationGeneration` internal/engine/store/learning_curation.go:163; `Store.GetLearningCurationGenerationByAttempt` internal/engine/store/learning_curation.go:54; `Store.GetLearningIncorporation` internal/engine/store/learning.go:415; `Store.GetLearningPromptManifest` internal/engine/store/learning_prompt.go:74; `Store.GetLearningReflectionRequest` internal/engine/store/learning.go:38; `Store.GetLearningReflectionRequestByAttempt` internal/engine/store/learning.go:49; `Store.GetTaskTransition` internal/engine/store/workflow_trigger.go:86; `Store.LatestTaskTransition` internal/engine/store/workflow_trigger.go:66 |
| structural | 413.00 | 20.65 | `Store.GetActiveReviewPass` internal/engine/store/review.go:85; `Store.GetDefaultReviewTour` internal/engine/store/review.go:54; `Store.GetReeveChangeSet` internal/engine/store/reeve_change_sets.go:47; `Store.GetReevePeer` internal/engine/store/reeve_peers.go:35; `Store.GetReviewPassPlan` internal/engine/store/review.go:184; `Store.GetReviewStep` internal/engine/store/review.go:286; `Store.GetReviewTour` internal/engine/store/review.go:44; `Store.HandoffObservation` internal/engine/store/handoff.go:200 |
| structural | 410.00 | 20.50 | `Store.GitHubAppLogins` internal/engine/store/pr_reviews_seen.go:64; `Store.ListLaunchMessages` internal/engine/store/launch_messages.go:37; `Store.ListReeveInbox` internal/engine/store/reeve_peers.go:105; `Store.ListReevePeers` internal/engine/store/reeve_peers.go:45; `Store.ListRepositoryGuidance` internal/engine/store/repository_guidance.go:21; `Store.ListRepositoryPresets` internal/engine/store/repository_preset.go:32; `Store.ListSecrets` internal/engine/store/credential.go:57; `Store.ListSettingEntries` internal/engine/store/settings.go:75; `Store.ListSettings` internal/engine/store/settings.go:19; `Store.ListSubscriptionUsage` internal/engine/store/subscription_usage.go:14; `Store.WorkflowTriggerSwitches` internal/engine/store/workflow_trigger_switches.go:14 |
| exact | 325.00 | 16.25 | `ok` web/src/lib/jobClones.ts:31; `ok` web/src/lib/mcpConnections.ts:132; `ok` web/src/lib/repositories.ts:30; `ok` web/src/lib/repositoryGuidance.ts:15; `ok` web/src/lib/resources.ts:20; `ok` web/src/lib/skills.ts:22 |
| structural | 260.00 | 13.00 | `Store.DeleteRepositoryGuidance` internal/engine/store/repository_guidance.go:28; `Store.DeleteRepositoryPreset` internal/engine/store/repository_preset.go:39; `Store.RevokeIntegrationToken` internal/engine/store/integration_tokens.go:156; `Store.RevokeMainIntegrationToken` internal/engine/store/integration_tokens.go:113; `Store.SetExecutionReadOnly` internal/engine/store/execution_broker.go:38 |
| structural | 240.00 | 12.00 | `Client.ListCommands` internal/engine/chimera/client.go:455; `Client.ListModels` internal/engine/chimera/client.go:473; `Client.ListModes` internal/engine/chimera/client.go:431; `Client.ListSkills` internal/engine/chimera/client.go:464; `Client.ListSubagents` internal/engine/chimera/client.go:217 |
| structural | 214.00 | 10.70 | `Server.handoffHandler` internal/server/handoff.go:19; `Server.planHandler` internal/server/plan.go:67; `Server.planRevisionsHandler` internal/server/plan.go:85 |
| structural | 210.00 | 10.50 | `Store.GetApprovalDevice` internal/engine/store/approval_devices.go:53; `Store.GetAttentionItem` internal/engine/store/attention.go:65; `Store.GetPlan` internal/engine/store/plan.go:10; `Store.GetSecret` internal/engine/store/credential.go:34 |
| structural | 209.00 | 10.45 | `buildMWShortcuts` internal/tui/overlay/move_workspace.go:72; `buildWSShortcuts` internal/tui/overlay/workspace.go:72 |
| structural | 206.00 | 10.30 | `Server.swarmCancelHandler` internal/server/swarm.go:293; `Server.swarmCloseHandler` internal/server/swarm.go:196 |
| structural | 200.00 | 10.00 | `Store.CountApprovals` internal/engine/store/approval.go:178; `Store.CountUnansweredRequiredPlanQuestions` internal/engine/store/plan.go:139; `Store.HasHandoffPackets` internal/engine/store/handoff.go:238; `Store.HasReeveWatchEventKind` internal/engine/store/reeve_watches.go:106; `Store.TaskActivityMark` internal/engine/store/reeve_watches.go:96; `Store.WorkflowAttemptInUse` internal/engine/store/workflow.go:1114 |
| structural | 190.00 | 9.50 | `Server.StartAttentionEvents` internal/server/server.go:990; `Server.StartPlanEvents` internal/server/server.go:974; `Server.StartReeveSpokeEvents` internal/server/reeve_voice.go:76 |
| structural | 183.00 | 9.15 | `Store.ClaimHandoffNotice` internal/engine/store/handoff_notices.go:10; `Store.ClaimInstanceValue` internal/engine/store/approval_devices.go:221; `Store.SetWorkerRefused` internal/engine/store/workers.go:103; `Store.SetWorkerUnreachable` internal/engine/store/workers.go:91 |
| structural | 178.00 | 8.90 | `runSwarmCancel` cmd/legato/main.go:2246; `runSwarmInbox` cmd/legato/main.go:2442; `runSwarmStatus` cmd/legato/main.go:2425 |
| exact | 177.00 | 8.85 | `Bridge.publish` internal/engine/claudebridge/stream.go:27; `Bridge.publish` internal/engine/codexbridge/stream.go:30 |
| structural | 177.00 | 8.85 | `Store.ActiveOAuthAttempt` internal/engine/store/oauth_attempts.go:78; `Store.CurrentOAuthAttempt` internal/engine/store/oauth_lifecycle.go:27; `Store.GetMCPConnection` internal/engine/store/mcp_connections.go:93; `Store.GetOAuthAttempt` internal/engine/store/oauth_attempts.go:69 |
| structural | 164.00 | 8.20 | `fetchMovable` web/src/lib/configSettings.ts:122; `fetchReeve` web/src/lib/reeve.ts:42; `fetchReeveEngram` web/src/lib/reeve.ts:194; `listDevices` web/src/lib/devices.ts:24; `listWorkers` web/src/lib/workers.ts:111 |
| structural | 159.00 | 7.95 | `InstallClaudeMod` internal/engine/hooks/claude_mod.go:70; `InstallClaudeStateMod` internal/engine/hooks/claude_mod.go:103; `InstallClaudeToolsMod` internal/engine/hooks/claude_mod.go:92; `InstallClaudeUsageMod` internal/engine/hooks/claude_mod.go:81 |
| structural | 150.00 | 7.50 | `Store.VerifyLearningIncorporation` internal/engine/store/learning.go:475; `Store.VerifyLearningIncorporationTarget` internal/engine/store/learning.go:444 |
| structural | 148.00 | 7.40 | `Store.GetAttentionItemBySource` internal/engine/store/attention.go:76; `Store.GetPlanByTaskName` internal/engine/store/plan.go:27; `Store.GetPlanQuestionByKey` internal/engine/store/plan.go:148 |
| structural | 146.00 | 7.30 | `cancelSwarm` web/src/lib/swarm.ts:190; `closeWorker` web/src/lib/swarm.ts:135; `dispatchWorker` web/src/lib/swarm.ts:100 |
| structural | 143.00 | 7.15 | `agentService.discoverClaudeBridge` internal/service/claude_bridge.go:24; `agentService.discoverCodexBridge` internal/service/codex_bridge.go:27 |
| structural | 138.00 | 6.90 | `Store.ListAutoReeveChangeSets` internal/engine/store/reeve_change_sets.go:105; `Store.ListDelegationsByState` internal/engine/store/delegations.go:108; `Store.ListDelegationsForTask` internal/engine/store/delegations.go:115; `Store.ListRestartRecords` internal/engine/store/restart_record.go:78 |
| structural | 138.00 | 6.90 | `Store.ListReviewHunkNotes` internal/engine/store/review.go:600; `Store.ListReviewMessages` internal/engine/store/review.go:699; `Store.ListReviewSteps` internal/engine/store/review.go:276; `Store.ListTasksForRepository` internal/engine/store/integration_tokens.go:179 |
| structural | 136.00 | 6.80 | `Store.ApprovePlan` internal/engine/store/plan.go:213; `Store.BindLearningPromptManifestSession` internal/engine/store/learning_prompt.go:59; `Store.SetParentActiveStep` internal/engine/store/swarm.go:111 |
| structural | 136.00 | 6.80 | `Store.GetScopeBatch` internal/engine/store/scope.go:173; `Store.GetScopePin` internal/engine/store/scope.go:125; `Store.GetScopeSubscription` internal/engine/store/scope.go:260 |
| structural | 135.00 | 6.75 | `relativeAge` cmd/legato/attention_cmd.go:111; `formatRelativeTime` internal/tui/attention/model.go:253 |
| structural | 135.00 | 6.75 | `Store.ResetWorkflowStepForRetry` internal/engine/store/workflow.go:968; `Store.TriggerWorkflowStep` internal/engine/store/workflow.go:411 |
| structural | 132.00 | 6.60 | `conceptIDs` internal/service/handoff.go:639; `evidenceIDs` internal/service/handoff.go:646; `scenarioIDs` internal/service/handoff.go:653; `sourceInputIDs` internal/service/handoff.go:632 |

## Rule-violating edges

| Rule | Edge | Import |
|---|---|---|
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui` → `internal/engine/events` | internal/tui/app.go:15 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui` → `internal/engine/store` | internal/tui/app.go:16 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui` → `internal/engine/workflow` | internal/tui/app.go:17 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui/board` → `internal/engine/store` | internal/tui/board/model.go:9 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui/overlay` → `internal/engine/github` | internal/tui/overlay/link_pr.go:10 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui/overlay` → `internal/engine/store` | internal/tui/overlay/workspace.go:9 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui/overlay` → `internal/engine/workflow` | internal/tui/overlay/workflow_trigger.go:9 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui/plan` → `internal/engine/store` | internal/tui/plan/model.go:13 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui/report` → `internal/engine/analytics` | internal/tui/report/model.go:12 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui/review` → `internal/engine/git` | internal/tui/review/messages.go:4 |
| tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests) | `internal/tui/review` → `internal/engine/store` | internal/tui/review/model.go:17 |

## Cycles

- browser-extension/src ↔ browser-extension/src/scope-source
- desktop/src ↔ desktop/src/components
- web/src/components ↔ web/src/components/agent-graph ↔ web/src/components/needsYou ↔ web/src/components/settings ↔ web/src/components/ui ↔ web/src/components/workflow-graph ↔ web/src/hooks ↔ web/src/lib ↔ web/src/voice

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
