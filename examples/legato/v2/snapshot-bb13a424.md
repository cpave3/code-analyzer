# CodeEval v2 snapshot: legato-stack @ `bb13a424dc73`

- **Revision:** `bb13a424dc73f3eeb98f8bce43db932a9d19f7e6` tree `25b725950d6a` 2026-11-27T20:29:04+11:00
- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.1.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 990 of 1860 production source files measured (61.9% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2, typescript 577

## Codebase

| Score (0–100) | Debt | Debt / kLOC | Source lines | Callables |
|---:|---:|---:|---:|---:|
| **34.40** | **11387.20** | 47.67 | 238865 | 9336 |

| Component | Raw | Debt |
|---|---|---:|
| complexity | cyclomatic excess 9888 (hotspots 1285; decision points 38363, cyclomatic total 47699) | 9888.00 |
| nesting | nesting excess 306 | 612.00 |
| duplication | redundant tokens 15544.00 (structural redundant 253, exact 8, cloned callables 409) | 777.20 |
| architecture | rule-violating edges 11, units in cycles 0 (internal edges 199) | 110.00 |

## Modules (top 25 of 87 by debt)

| Module | Score | Debt | Debt/kLOC | Lines | Complexity | Nesting | Duplication | Architecture |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `internal/service` | 32.82 | 4932.45 | 51.17 | 96395 | 4618.00 | 242.00 | 72.45 | 0.00 |
| `cmd/legato` | 34.75 | 871.93 | 46.94 | 18575 | 809.00 | 26.00 | 36.93 | 0.00 |
| `internal/engine/store` | 37.49 | 833.90 | 41.69 | 20004 | 403.00 | 14.00 | 416.90 | 0.00 |
| `internal/engine/workflow` | 29.35 | 582.00 | 60.19 | 9670 | 538.00 | 42.00 | 2.00 | 0.00 |
| `internal/server` | 48.77 | 522.80 | 26.26 | 19905 | 389.00 | 66.00 | 67.80 | 0.00 |
| `internal/tui` | 21.94 | 428.20 | 88.93 | 4815 | 372.00 | 18.00 | 8.20 | 30.00 |
| `internal/cli` | 33.04 | 426.75 | 50.67 | 8422 | 376.00 | 16.00 | 34.75 | 0.00 |
| `internal/tui/overlay` | 28.51 | 361.20 | 62.68 | 5763 | 265.00 | 50.00 | 16.20 | 30.00 |
| `config` | 35.48 | 246.10 | 45.46 | 5414 | 232.00 | 8.00 | 6.10 | 0.00 |
| `internal/engine/dockerproxy` | 17.25 | 242.30 | 119.95 | 2020 | 230.00 | 8.00 | 4.30 | 0.00 |
| `internal/engine/credential` | 30.44 | 225.00 | 57.14 | 3938 | 223.00 | 2.00 | 0.00 | 0.00 |
| `internal/tui/agents` | 33.16 | 148.15 | 50.39 | 2940 | 142.00 | 4.00 | 2.15 | 0.00 |
| `internal/engine/codexbridge` | 34.48 | 140.48 | 47.51 | 2957 | 117.00 | 6.00 | 17.48 | 0.00 |
| `internal/engine/claudebridge` | 40.81 | 116.83 | 36.26 | 3222 | 91.00 | 12.00 | 13.83 | 0.00 |
| `internal/setup` | 27.68 | 77.00 | 65.31 | 1179 | 73.00 | 4.00 | 0.00 | 0.00 |
| `internal/engine/jira` | 32.43 | 74.70 | 52.09 | 1434 | 59.00 | 12.00 | 3.70 | 0.00 |
| `internal/tui/board` | 34.46 | 73.00 | 47.56 | 1535 | 54.00 | 4.00 | 5.00 | 10.00 |
| `internal/engine/hooks` | 50.32 | 58.95 | 24.69 | 2388 | 41.00 | 10.00 | 7.95 | 0.00 |
| `internal/tui/detail` | 27.88 | 58.00 | 64.66 | 897 | 54.00 | 4.00 | 0.00 | 0.00 |
| `internal/tui/review` | 28.39 | 56.25 | 63.06 | 892 | 30.00 | 4.00 | 2.25 | 20.00 |
| `internal/tui/plan` | 15.62 | 52.25 | 135.01 | 387 | 36.00 | 4.00 | 2.25 | 10.00 |
| `internal/engine/oauthclient` | 38.12 | 47.00 | 40.59 | 1158 | 45.00 | 2.00 | 0.00 | 0.00 |
| `internal/handofftest` | 56.12 | 45.05 | 19.54 | 2305 | 33.00 | 2.00 | 10.05 | 0.00 |
| `internal/engine/discordbot` | 41.33 | 45.00 | 35.49 | 1268 | 43.00 | 2.00 | 0.00 | 0.00 |
| `internal/tui/attention` | 14.57 | 40.03 | 146.61 | 273 | 31.00 | 4.00 | 5.03 | 0.00 |

## Largest callable debts

| Callable | Location | cc | nesting | lines | debt |
|---|---|---:|---:|---:|---:|
| `App.update` | internal/tui/app.go:696-1827 | 308 | 6 | 1132 | 305.00 |
| `agentService.SpawnAgent` | internal/service/agent.go:837-1693 | 287 | 5 | 857 | 282.00 |
| `BuildSandboxPlan` | internal/service/sandbox_plan.go:414-742 | 96 | 4 | 329 | 89.00 |
| `runReviewCmd` | cmd/legato/main.go:347-703 | 93 | 4 | 357 | 86.00 |
| `prepareMigration` | internal/engine/store/store.go:364-683 | 93 | 4 | 320 | 86.00 |
| `runLearnCmd` | cmd/legato/learning_cmd.go:69-349 | 81 | 3 | 281 | 72.00 |
| `runTUI` | cmd/legato/main.go:1232-1806 | 67 | 4 | 575 | 60.00 |
| `TaskHandoffService.HandOff` | internal/service/task_handoff.go:203-396 | 69 | 2 | 194 | 60.00 |
| `Proxy.checkHostConfig` | internal/engine/dockerproxy/policy.go:146-318 | 62 | 3 | 173 | 53.00 |
| `agentService.sandboxPrefix` | internal/service/agent_sandbox.go:49-240 | 58 | 3 | 192 | 49.00 |
| `runPlanCmd` | cmd/legato/main.go:134-345 | 57 | 3 | 212 | 48.00 |
| `applyClaudeHookAt` | internal/cli/claude_hook.go:303-450 | 55 | 3 | 148 | 46.00 |
| `Server.taskAgentHandler` | internal/server/agents.go:656-844 | 52 | 4 | 189 | 45.00 |
| `validateCuratorOperation` | internal/service/learning_curation.go:414-498 | 54 | 2 | 85 | 45.00 |
| `Proxy.check` | internal/engine/dockerproxy/proxy.go:185-263 | 52 | 3 | 79 | 43.00 |
| `Model.HandleAction` | internal/tui/agents/model.go:537-673 | 50 | 4 | 137 | 43.00 |
| `renderHandoffMarkdown` | internal/service/handoff_export.go:12-178 | 47 | 5 | 167 | 42.00 |
| `RunChimeraLink` | internal/cli/chimera_link.go:38-223 | 46 | 5 | 186 | 41.00 |
| `Server.taskHandler` | internal/server/board.go:260-463 | 46 | 5 | 204 | 41.00 |
| `prTrackingService.syncTask` | internal/service/pr_branches.go:367-488 | 45 | 5 | 122 | 40.00 |
| `Bridge.handle` | internal/engine/codexbridge/bridge.go:311-494 | 44 | 5 | 184 | 39.00 |
| `validateHandoffInput` | internal/service/handoff.go:383-487 | 48 | 3 | 105 | 39.00 |
| `HTTPProxy.ServeHTTP` | internal/engine/credential/httpproxy.go:158-314 | 47 | 3 | 157 | 38.00 |
| `AgentGraphService.workflowPlaces` | internal/service/agent_graph_workflow.go:17-161 | 46 | 3 | 145 | 37.00 |
| `HTTPProxy.serveGit` | internal/engine/credential/httpproxy_git.go:71-236 | 45 | 3 | 166 | 36.00 |

## Largest clone groups

| Kind | Redundant tokens | Debt | Members |
|---|---:|---:|---|
| structural | 1276.00 | 63.80 | `Store.ClosedLoopReviewers` internal/engine/store/review_loops.go:1115; `Store.ConversationStorageDirs` internal/engine/store/agent_conversation_recovery.go:41; `Store.ListChosenExecutionResources` internal/engine/store/resource.go:148; `Store.ListExecutionResources` internal/engine/store/resource.go:157; `Store.ListOAuthAttempts` internal/engine/store/oauth_attempts.go:188; `Store.ListOpenIntegrationAttention` internal/engine/store/integration_attention.go:76; `Store.ListPlanComments` internal/engine/store/plan.go:183; `Store.ListPlanMessages` internal/engine/store/plan.go:331; `Store.ListPlanQuestions` internal/engine/store/plan.go:106; `Store.ListPlanResponses` internal/engine/store/plan.go:168; `Store.ListPlanRevisions` internal/engine/store/plan.go:306; `Store.ListPlansByTask` internal/engine/store/plan.go:21; `Store.ListReeveWatchEvents` internal/engine/store/reeve_watches.go:55; `Store.ListReviewChapterHunks` internal/engine/store/review.go:453; `Store.ListReviewFindingsByPass` internal/engine/store/review.go:123; `Store.ListReviewHunkNotesByPass` internal/engine/store/review.go:108; `Store.ListReviewLoopFindings` internal/engine/store/review_loops.go:993; `Store.ListReviewLoopItems` internal/engine/store/review_loops.go:873; `Store.ListReviewLoopPRs` internal/engine/store/review_loops.go:937; `Store.ListReviewLoopPasses` internal/engine/store/review_loops.go:986; `Store.ListReviewMessagesByPass` internal/engine/store/review.go:102; `Store.ListReviewPasses` internal/engine/store/review.go:79; `Store.ListReviewStepsByPass` internal/engine/store/review.go:94; `Store.ListReviewToursByTask` internal/engine/store/review.go:64; `Store.ListRunReviewLoopsAuthoredBy` internal/engine/store/handoff_reviews.go:1297; `Store.ListSubtasksByParent` internal/engine/store/swarm.go:70; `Store.ListSyncLogs` internal/engine/store/store.go:1073; `Store.ListTaskNotes` internal/engine/store/task_notes.go:10; `Store.ListUnackedSwarmEvents` internal/engine/store/swarm.go:245; `Store.OpenSubagents` internal/engine/store/subagent_activity.go:123 |
| structural | 520.00 | 26.00 | `Store.GetExecution` internal/engine/store/execution.go:70; `Store.GetLearningCurationGeneration` internal/engine/store/learning_curation.go:163; `Store.GetLearningCurationGenerationByAttempt` internal/engine/store/learning_curation.go:54; `Store.GetLearningIncorporation` internal/engine/store/learning.go:415; `Store.GetLearningPromptManifest` internal/engine/store/learning_prompt.go:74; `Store.GetLearningReflectionRequest` internal/engine/store/learning.go:38; `Store.GetLearningReflectionRequestByAttempt` internal/engine/store/learning.go:49; `Store.GetTaskTransition` internal/engine/store/workflow_trigger.go:86; `Store.LatestTaskTransition` internal/engine/store/workflow_trigger.go:66 |
| structural | 413.00 | 20.65 | `Store.GetActiveReviewPass` internal/engine/store/review.go:85; `Store.GetDefaultReviewTour` internal/engine/store/review.go:54; `Store.GetReeveChangeSet` internal/engine/store/reeve_change_sets.go:47; `Store.GetReevePeer` internal/engine/store/reeve_peers.go:35; `Store.GetReviewPassPlan` internal/engine/store/review.go:184; `Store.GetReviewStep` internal/engine/store/review.go:286; `Store.GetReviewTour` internal/engine/store/review.go:44; `Store.HandoffObservation` internal/engine/store/handoff.go:200 |
| structural | 410.00 | 20.50 | `Store.GitHubAppLogins` internal/engine/store/pr_reviews_seen.go:64; `Store.ListLaunchMessages` internal/engine/store/launch_messages.go:37; `Store.ListReeveInbox` internal/engine/store/reeve_peers.go:105; `Store.ListReevePeers` internal/engine/store/reeve_peers.go:45; `Store.ListRepositoryGuidance` internal/engine/store/repository_guidance.go:21; `Store.ListRepositoryPresets` internal/engine/store/repository_preset.go:32; `Store.ListSecrets` internal/engine/store/credential.go:57; `Store.ListSettingEntries` internal/engine/store/settings.go:75; `Store.ListSettings` internal/engine/store/settings.go:19; `Store.ListSubscriptionUsage` internal/engine/store/subscription_usage.go:14; `Store.WorkflowTriggerSwitches` internal/engine/store/workflow_trigger_switches.go:14 |
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
| structural | 159.00 | 7.95 | `InstallClaudeMod` internal/engine/hooks/claude_mod.go:70; `InstallClaudeStateMod` internal/engine/hooks/claude_mod.go:103; `InstallClaudeToolsMod` internal/engine/hooks/claude_mod.go:92; `InstallClaudeUsageMod` internal/engine/hooks/claude_mod.go:81 |
| structural | 150.00 | 7.50 | `Store.VerifyLearningIncorporation` internal/engine/store/learning.go:475; `Store.VerifyLearningIncorporationTarget` internal/engine/store/learning.go:444 |
| structural | 148.00 | 7.40 | `Store.GetAttentionItemBySource` internal/engine/store/attention.go:76; `Store.GetPlanByTaskName` internal/engine/store/plan.go:27; `Store.GetPlanQuestionByKey` internal/engine/store/plan.go:148 |
| structural | 143.00 | 7.15 | `agentService.discoverClaudeBridge` internal/service/claude_bridge.go:24; `agentService.discoverCodexBridge` internal/service/codex_bridge.go:27 |
| structural | 138.00 | 6.90 | `Store.ListAutoReeveChangeSets` internal/engine/store/reeve_change_sets.go:105; `Store.ListDelegationsByState` internal/engine/store/delegations.go:108; `Store.ListDelegationsForTask` internal/engine/store/delegations.go:115; `Store.ListRestartRecords` internal/engine/store/restart_record.go:78 |
| structural | 138.00 | 6.90 | `Store.ListReviewHunkNotes` internal/engine/store/review.go:600; `Store.ListReviewMessages` internal/engine/store/review.go:699; `Store.ListReviewSteps` internal/engine/store/review.go:276; `Store.ListTasksForRepository` internal/engine/store/integration_tokens.go:179 |
| structural | 136.00 | 6.80 | `Store.ApprovePlan` internal/engine/store/plan.go:213; `Store.BindLearningPromptManifestSession` internal/engine/store/learning_prompt.go:59; `Store.SetParentActiveStep` internal/engine/store/swarm.go:111 |
| structural | 136.00 | 6.80 | `Store.GetScopeBatch` internal/engine/store/scope.go:173; `Store.GetScopePin` internal/engine/store/scope.go:125; `Store.GetScopeSubscription` internal/engine/store/scope.go:260 |
| structural | 135.00 | 6.75 | `relativeAge` cmd/legato/attention_cmd.go:111; `formatRelativeTime` internal/tui/attention/model.go:253 |

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

## Notes

- experimental scoring/v2: separate from and not comparable to scoring/v1 −2..+2 scores; v1 reports are unchanged
- debt points = complexity (Σ max(0, cyclomatic − H + 1): debt exactly when a callable is a hotspot) + nesting (Σ max(0, depth − N + 1) × weight) + duplication (structural-clone redundant tokens ÷ tokens-per-point) + architecture (configured-rule-violating edges and units in cycles × weights)
- debt is additive across modules and telescopes across commits (Δ = after − before of one snapshot function); it is the quantity to optimise within one codebase because adding code cannot dilute it
- score = 100 × h ÷ (h + debt per 1,000 supported source lines), h = halfScoreDensity: a size-normalised view for comparing modules or codebases; it can be diluted by adding simple code, so do not optimise it alone
- weights are uncalibrated defaults; the score is a structural index, not a measure of correctness, design quality or behaviour
- only languages with adapters are measured; see coverage — unsupported source is not scored and not neutral
- architecture debt reflects only the configured import-level rules; rules the configuration cannot express are not measured
- complexity growth below the hotspot thresholds carries no debt; it is reported as decision points and checked by the existing-callable branch gate
- gates are advisory checks, not merge decisions
