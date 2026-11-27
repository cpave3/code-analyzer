# Legato findings worklist: everything CodeEval flagged in 9a61e9f8..471bfd24

Requested by Cameron: address **all** findings from the two-day window, not
only the subset PR #434 fixed. Locations are at `471bfd24`. Line numbers may
have shifted slightly in PR #434's three files, but nothing else moved.

Source data: `net-9a61e9f8..471bfd24.json` (`results[].findings`). The full
lists of the 127 modified callables that got more complex and the 84 new
hotspots are in that JSON under the `complexity` result.

**Ground rules for the fixing agent**

- Every item gets an outcome: **fixed**, or **won't fix** with a written
  reason in the PR (for example "idiomatic one-liner; abstraction would add
  indirection").
- CodeEval matches code shape. It does not prove two pieces of code are
  equivalent. Confirm semantics before merging code.
- Preserve behaviour. Add tests where behaviour isn't pinned. Legato's
  architecture tests must stay green, including the test-enforced rule
  "engine packages must not import sibling engine packages"
  (`internal/engine/store/architecture_test.go`). No new cycles.
- Keep commits focused, one concern each, so the per-commit delta can be
  attributed.
- Base the work on the current head of PR #434 (`41c80614`) or on main once
  it's merged.

---

## A. Duplication: done in PR #434 (for reference)

| Group | Outcome |
|---|---|
| `Server.StartSubscriptionUsageEvents` `internal/server/subscription_usage.go:70-87` ↔ `Server.StartServerThemeEvents` `internal/server/theme_preference.go:23-40` | Fixed: both now delegate to `broadcastInvalidations` |
| `Client.doPut` `internal/engine/jira/users.go:57-64` ↔ `Client.doPost` `internal/engine/jira/client.go:157-164` | Fixed: replaced by `doJSON(method, …)` |
| `internal/service/plan.go` inline closure ↔ `pathWithin` | Fixed (not visible to CodeEval: closures aren't indexed) |

## B. Exact copies across sibling engine packages (blocked on a layering decision)

These are **byte-identical** (exact fingerprint match) or structurally
identical copies across `claudeusage`, `codexusage` and `opencodeusage`.
PR #434 left them alone because the engine-sibling rule leaves no permitted
shared home.

| # | Copies | Kind |
|---|---|---|
| B1 | `retryAfter` `internal/engine/claudeusage/claudeusage.go:197-209` and `internal/engine/opencodeusage/opencodeusage.go:185-197` | exact (the opencode copy duplicates the pre-existing claude one) |
| B2 | `openRegular` `internal/engine/claudeusage/open.go:17-28` and `internal/engine/codexusage/open.go:17-28` | exact |
| B3 | `CredentialsPath` `internal/engine/claudeusage/credentials.go:17-22` ↔ `SessionsDir` `internal/engine/codexusage/logs.go:24-29` | structural (same home-dir path resolution) |

**Decision needed (Cameron or Tessitura) before the fix:** where may
provider-agnostic helpers live? Options:

1. a new leaf package **outside** `internal/engine` (for example
   `internal/fsutil` or `internal/httputil`) that imports only stdlib and that
   engine packages may depend on;
2. an explicit, documented exception to the sibling rule for one shared
   `internal/engine/usagecommon` package;
3. keep the copies and document why.

The agent should propose option 1 with an architecture-test update if needed,
and must not silently relax the test.

## C. Remaining structural duplicates: likely real, fix or justify

| # | Copies | Notes |
|---|---|---|
| C1 | `claudeUsageError` `internal/service/subscription_usage_readers.go:52-73` ↔ `opencodeUsageError` same file `:176-197` | 22-line identical shape. PR #434 kept them as "provider-specific". Extract a shared mapper that takes the provider-specific parts (error types and messages) as parameters, or justify. |
| C2 | `App.loadSubscriptionUsage` `internal/tui/app_subscription_usage.go:35-44` ↔ pre-existing `App.loadUsage` `internal/tui/app_usage.go:31-40` | Same load-command shape. Consider a generic loader helper. |
| C3 | `claimUsagePoller` `cmd/legato/usage_runtime.go:32-46` ↔ pre-existing `claimChatBots` `cmd/legato/chatbots_runtime.go:68-82` and `claimMQTTPublisher` `cmd/legato/mqtt_runtime.go:27-41` | Three copies of the same "claim a singleton runtime role" pattern. A shared `claimRuntime(name, …)` would remove two. |
| C4 | `accountPathContains` `config/claude_accounts.go:76-79` ↔ `pathWithin` `internal/service/learning_drift.go:153-156` | `config` must stay standalone (CLAUDE.md), so the shared home has the same question as B (a leaf path helper both can import). |
| C5 | `PullRequestModes.ReadyAfterReview` `internal/service/pull_request_modes.go:45-55` ↔ `PullRequestModes.PullRequestMode` same file `:23-33` | Same lookup shape. Consider one private lookup used by both. |

## D. Structural matches that are probably idiom: verify, fix only if it helps

These are flagged but look like normal Go idiom. Record "won't fix" with a
reason unless a clean shared form exists.

| # | Copies |
|---|---|
| D1 | `Store.ListSubscriptionUsage` `internal/engine/store/subscription_usage.go:14-18`, matching 10 existing thin `Store.List*` wrappers (`ListSecrets` `credential.go:57-61`, `ListLaunchMessages` `launch_messages.go:37-42`, `GitHubAppLogins` `pr_reviews_seen.go:64-68`, `ListReevePeers`/`ListReeveInbox` `reeve_peers.go:45-49,105-109`, `ListRepositoryGuidance` `repository_guidance.go:21-25`, `ListRepositoryPresets` `repository_preset.go:32-36`, `ListSettings`/`ListSettingEntries` `settings.go:19-23,75-79`, `WorkflowTriggerSwitches` `workflow_trigger_switches.go:14-18`) |
| D2 | `Store.LatestWorkflowRuns` `internal/engine/store/workflow_stats.go:17-23` ↔ `ListAgentSessions` `agent_session.go:104-109`, `ListInterruptedAgentSessions` `:396-401`, `ListAwayTasks` `handoffs.go:287-291` |
| D3 | `Store.ListReviewLoopsClosedSince` `internal/engine/store/review_loops.go:574-579` ↔ `OAuthClientRegistrations` `oauth_client_registrations.go:61-65`, `ListSubtasksByParentAndStatus` `swarm.go:78-84` |
| D4 | `Store.ListWorkflowRunsForTask` `internal/engine/store/workflow.go:1041-1053` ≈ `ListAgentUIItems` `agent_ui.go:144-152`, `ListWorkflowTriggerRequestsForTask`/`ForTransition` `workflow_trigger.go:103-107,97-101` (near-structural, similarity 0.88, low confidence) |
| D5 | `UsageProfile.UnmarshalYAML` ↔ `UsageProfileReviewer.UnmarshalYAML` `config/usage_profiles.go:120-126,128-134` |
| D6 | burnlite `Client.ClaimTask`/`ReleaseClaim`/`SetStage` `internal/engine/burnlite/assign.go:19-22,26-29,33-36` |
| D7 | Stop/watch pairs: `ReviewLoopSteering.StopWatching`/`WatchAgain` `internal/service/review_loop_steering.go:102-107,111-116`; `ReviewPassService.StopWatchingOn`/`WatchAgainOn` `internal/service/review_steward_hold.go:102-108,114-120`; `ReviewRunCaller.stopWatchingOn`/`watchAgainOn` same file `:203-208,210-215`. Three layers each have a stop/again pair; consider whether one `setWatching(on bool)` per layer is clearer. |
| D8 | Test fake `FakeReviewGitHub.Draft` `internal/handofftest/reviews.go:141-150` ↔ `.Push` `:830-839` |

## E. Complexity growth in pre-existing functions (127 functions got more complex; net +194)

Top increases, before → after cyclomatic complexity (cc):

| # | Function | Location (target) | cc | Nesting | Lines |
|---|---|---|---|---|---|
| E1 | `agentService.SpawnAgent` | `internal/service/agent.go:837-1693` | 266→287 | 5 | 810→857 |
| E2 | `agentService.prepareTaskResume` | `internal/service/agent_resume.go:60-186` | 28→48 | 2→3 | 82→127 |
| E3 | `BuildSandboxPlan` | `internal/service/sandbox_plan.go:414-769` | 94→107 | 4 | 320→356 |
| E4 | `MacroService.Start` | `internal/service/macro_start.go:72-152` | 14→26 | 2→4 | 48→81 |
| E5 | `HandoffWorkerService.Accept` | `internal/service/task_handoff_worker.go:190-315` | 25→36 | 3 | 91→126 |
| E6 | `agentService.sandboxPrefix` | `internal/service/agent_sandbox.go:48-263` | 58→67 | 3 | 191→216 |
| E7 | `ReviewRunCaller.AuthorDone` | `internal/service/review_pass_run.go:1073-1122` | 9→17 | 2→3 | 32→50 |
| E8 | `stepAgent` | `internal/service/workflow_macro.go:39-60` | 3→11 | 1 | 8→22 |
| E9 | `ReviewRunCaller.waitingToOpenOn` | `internal/service/review_pass_run.go:522-567` | 10→17 | 2 | 22→46 |
| E10 | `stripFor` | `internal/tui/board/card.go:126-153` | 6→13 | 1 | 18→28 |
| E11 | `ReviewPassService.confirm` | `internal/service/review_loop_ops.go:234-312` | 18→24 | 4 | 57→79 |
| E12 | `ReviewPassService.feedback` | `internal/service/review_pass_feedback.go:331-381` | 15→21 | 2→3 | 34→51 |
| E13 | `ReviewRunCaller.agentDone` | `internal/service/review_pass_run.go:1158-1198` | 18→24 | 2→3 | 28→41 |
| E14 | `HandoffWorkerService.resumeAuthor` | `internal/service/task_handoff_resume.go:52-105` | 12→18 | 2→3 | 43→54 |
| E15 | `schemaReflector.object` | `internal/engine/workflow/schema.go:444-488` | 12→17 | 3 | 43→45 |
| E16 | `ReviewPromptAttentionAdapter.Execute` | `internal/service/attention_review_prompt.go:192-236` | 7→12 | 1→3 | 29→45 |
| E17 | `SessionHistoryService.Recover` | `internal/service/session_recovery.go:28-102` | 16→21 | 4 | 58→75 |
| E18 | `TicketCreateExecutor.Start` | `internal/service/workflow_ticket.go:171-246` | 16→21 | 3 | 56→76 |
| E19 | `Store.RecordTaskNativeSession` | `internal/engine/store/task_native_session.go:183-252` | 20→24 | 1 | 63→70 |
| E20 | `Server.taskAgentHandler` | `internal/server/agents.go:656-844` | 48→52 | 4 | 184→189 |

The other 107 increases are mostly +1 to +4 and are listed in the JSON.
Guidance: extract the **code added in the window** (the new branches) into
named helpers, so these functions return to roughly their base complexity.
Rewriting the whole of `SpawnAgent` (E1) or `BuildSandboxPlan` (E3) is a
separate, higher-risk job; split those into their own PRs with
characterisation tests first.

## F. New hotspot functions (84 added callables at cc ≥ 10 or nesting ≥ 4)

Top 20 by complexity (locations at target):

| # | Function | Location | cc / nesting |
|---|---|---|---|
| F1 | `agentService.resolveClaudeAccount` | `internal/service/agent_accounts.go:67-197` | 55 / 3 |
| F2 | `UsagePolicyService.evaluate` | `internal/service/usage_profile.go:111-227` | 38 / 4 |
| F3 | `TicketUpdateExecutor.Start` | `internal/service/workflow_ticket_update.go:115-229` | 37 / 3 |
| F4 | `ReviewPassService.watchLoop` | `internal/service/review_watch.go:305-398` | 33 / 3 |
| F5 | `ReviewPassService.routeItems` | `internal/service/review_steward.go:238-335` | 32 / 3 |
| F6 | `checkReviewLoopStep` | `internal/engine/workflow/review_loop.go:149-209` | 29 / 4 |
| F7 | `EvaluateUsageRules` | `internal/service/usage_policy.go:55-122` | 25 / 3 |
| F8 | `prTrackingService.RefreshBranchStates` | `internal/service/pr_branches.go:96-155` | 24 / 3 |
| F9 | `parseReviewOpen` | `internal/service/review_loop_steps.go:196-263` | 24 / 2 |
| F10 | `ReviewRunCaller.markStackReady` | `internal/service/review_loop_steps.go:1240-1305` | 24 / 2 |
| F11 | `ReviewRunCaller.details` | `internal/service/review_pass_run.go:1835-1897` | 24 / 3 |
| F12 | `upgrade` | `internal/engine/workflow/upgrade.go:138-207` | 23 / 3 |
| F13 | `upgrader.upgradeStep` | `internal/engine/workflow/upgrade.go:383-448` | 23 / 3 |
| F14 | `ReviewPassService.addStacked` | `internal/service/review_steward.go:344-411` | 22 / 3 |
| F15 | `ValidateClaudeAccounts` | `config/claude_accounts.go:20-74` | 21 / 3 |
| F16 | `SetTags` | `internal/engine/workflow/tags.go:60-110` | 20 / 2 |
| F17 | `ReviewRunCaller.closeOutputs` | `internal/service/review_loop_steps.go:707-764` | 20 / 3 |
| F18 | `HandoffReviewCaller.sweepLoops` | `internal/service/review_pass_handoff.go:263-297` | 20 / 3 |
| F19 | `ReviewRunCaller.recordDetails` | `internal/service/review_run_page.go:19-68` | 20 / 2 |
| F20 | `formatSubscriptionUsage` | `cmd/legato/usage_cmd.go:80-142` | 19 / 3 |

The remaining 64 (cc 10 to 19) are in the JSON (`complexity.added_hotspot`).
Note on F11: `ReviewRunCaller.Details` went from 26 to 3 in the window by
moving its logic into the new private `details` (cc 24). That was
redistribution, not reduction. Prioritise by complexity and how often the
code changes (F1 to F11 first). Validation functions (F15, `ValidateUsageProfiles`)
are often legitimately branchy; table-driven validation is an option.

## G. Architecture: documented rules vs code (decision, not a mechanical fix)

- 14 pre-existing imports from TUI into engine packages contradict CLAUDE.md
  ("tui imports service via interfaces, never engine/"), but Legato's tests
  don't enforce the rule:
  - `internal/tui/app.go` → `engine/events`, `engine/ipc`, `engine/macros`, `engine/store`, `engine/workflow`
  - `internal/tui/board/model.go` → `engine/store`
  - `internal/tui/overlay/link_pr.go` → `engine/github`; `macro_picker.go` → `engine/macros`; `workflow_trigger.go` → `engine/workflow`; `workspace.go` → `engine/store`
  - `internal/tui/plan/model.go` → `engine/store`
  - `internal/tui/report/model.go` → `engine/analytics`
  - `internal/tui/review/messages.go` → `engine/git`; `model.go` → `engine/store`
- `config/config.go:14` → `internal/engine/macros` contradicts "config is
  standalone".

These predate the window. Either update CLAUDE.md to describe the real rule
(for example "tui may use engine types, not engine behaviour"), or route these
imports through service interfaces and add an architecture test. This is
Cameron's call; agents should not start this without a decision.

---

## How the result will be measured

The evaluator stays frozen (code-evaluator `5e0c679`, `scoring/v1`, config
`legato-architecture.json`, configHash `3cef63633359769f`). The primary
comparison is the fix branch's pre-work base SHA → fix SHA; drift from
`471bfd24` is reported separately. A complete fix should show, against
`471bfd24`:

- exact redundancy 10 → 8 or lower (B1, B2)
- structural redundancy falling well below 258
- modified-callable cc deltas negative for E1 to E20
- fewer new hotspots
- no new rule violations or cycles

Behavioural correctness comes from Legato's tests and CI, not from these
numbers.
