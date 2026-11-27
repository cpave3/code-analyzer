# Debt factors: legato-stack @ `bb13a424dc73` — internal/service

- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.2.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 1555 of 1848 production source files measured (90.0% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2; 4 file(s) with parse errors

**Scope:** 331 source files, 96395 lines, 3778 callables (606 hotspots) in 1 module(s).

| Debt | complexity | nesting | duplication | architecture | Debt/kLOC | Score |
|---:|---:|---:|---:|---:|---:|---:|
| **4932.45** | 4618.00 | 242.00 | 72.45 | 0.00 | 51.17 | 32.82 |

**Concentration:** 711 factor(s) in scope; the largest 103 carry 50% of the debt and the largest 292 carry 80%.

## Factors (40 listed, 671 smaller omitted; largest first)

### 1. complexity `agentService.SpawnAgent` — 278.00 debt · 5.6% of scope (cumulative 5.6%)

- **Where:** `internal/service/agent.go:837-1693` (module `internal/service`, go)
- **Why:** cyclomatic 287 is 278 point(s) at or above the hotspot threshold 10 (286 decision points, 857 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 287) removes this factor; each point below 287 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `agentService.SpawnAgent` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 2. complexity `BuildSandboxPlan` — 87.00 debt · 1.8% of scope (cumulative 7.4%)

- **Where:** `internal/service/sandbox_plan.go:414-742` (module `internal/service`, go)
- **Why:** cyclomatic 96 is 87 point(s) at or above the hotspot threshold 10 (95 decision points, 329 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 96) removes this factor; each point below 96 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `BuildSandboxPlan` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 3. complexity `TaskHandoffService.HandOff` — 60.00 debt · 1.2% of scope (cumulative 8.6%)

- **Where:** `internal/service/task_handoff.go:203-396` (module `internal/service`, go)
- **Why:** cyclomatic 69 is 60 point(s) at or above the hotspot threshold 10 (68 decision points, 194 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 69) removes this factor; each point below 69 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `TaskHandoffService.HandOff` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 4. complexity `agentService.sandboxPrefix` — 49.00 debt · 1.0% of scope (cumulative 9.6%)

- **Where:** `internal/service/agent_sandbox.go:49-240` (module `internal/service`, go)
- **Why:** cyclomatic 58 is 49 point(s) at or above the hotspot threshold 10 (57 decision points, 192 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 58) removes this factor; each point below 58 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `agentService.sandboxPrefix` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 5. complexity `validateCuratorOperation` — 45.00 debt · 0.9% of scope (cumulative 10.5%)

- **Where:** `internal/service/learning_curation.go:414-498` (module `internal/service`, go)
- **Why:** cyclomatic 54 is 45 point(s) at or above the hotspot threshold 10 (53 decision points, 85 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 54) removes this factor; each point below 54 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `validateCuratorOperation` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 6. complexity `validateHandoffInput` — 39.00 debt · 0.8% of scope (cumulative 11.3%)

- **Where:** `internal/service/handoff.go:383-487` (module `internal/service`, go)
- **Why:** cyclomatic 48 is 39 point(s) at or above the hotspot threshold 10 (47 decision points, 105 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 48) removes this factor; each point below 48 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `validateHandoffInput` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 7. complexity `renderHandoffMarkdown` — 38.00 debt · 0.8% of scope (cumulative 12.1%)

- **Where:** `internal/service/handoff_export.go:12-178` (module `internal/service`, go)
- **Why:** cyclomatic 47 is 38 point(s) at or above the hotspot threshold 10 (46 decision points, 167 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 47) removes this factor; each point below 47 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `renderHandoffMarkdown` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 8. complexity `AgentGraphService.workflowPlaces` — 37.00 debt · 0.8% of scope (cumulative 12.8%)

- **Where:** `internal/service/agent_graph_workflow.go:17-161` (module `internal/service`, go)
- **Why:** cyclomatic 46 is 37 point(s) at or above the hotspot threshold 10 (45 decision points, 145 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 46) removes this factor; each point below 46 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `AgentGraphService.workflowPlaces` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 9. complexity `MCPConnectionService.update` — 36.00 debt · 0.7% of scope (cumulative 13.6%)

- **Where:** `internal/service/mcp_connections.go:424-546` (module `internal/service`, go)
- **Why:** cyclomatic 45 is 36 point(s) at or above the hotspot threshold 10 (44 decision points, 123 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 45) removes this factor; each point below 45 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `MCPConnectionService.update` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 10. complexity `prTrackingService.syncTask` — 36.00 debt · 0.7% of scope (cumulative 14.3%)

- **Where:** `internal/service/pr_branches.go:367-488` (module `internal/service`, go)
- **Why:** cyclomatic 45 is 36 point(s) at or above the hotspot threshold 10 (44 decision points, 122 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 45) removes this factor; each point below 45 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `prTrackingService.syncTask` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 11. complexity `WorkflowService.runStep` — 36.00 debt · 0.7% of scope (cumulative 15.0%)

- **Where:** `internal/service/workflow_loop.go:187-326` (module `internal/service`, go)
- **Why:** cyclomatic 45 is 36 point(s) at or above the hotspot threshold 10 (44 decision points, 140 lines)

### 12. complexity `ReviewPassService.post` — 35.00 debt · 0.7% of scope (cumulative 15.7%)

- **Where:** `internal/service/review_pass_submit.go:455-623` (module `internal/service`, go)
- **Why:** cyclomatic 44 is 35 point(s) at or above the hotspot threshold 10 (43 decision points, 169 lines)

### 13. complexity `ReviewService.EditAnnotation` — 32.00 debt · 0.7% of scope (cumulative 16.4%)

- **Where:** `internal/service/review.go:696-804` (module `internal/service`, go)
- **Why:** cyclomatic 41 is 32 point(s) at or above the hotspot threshold 10 (40 decision points, 109 lines)

### 14. complexity `NeedsYouByTask` — 31.00 debt · 0.6% of scope (cumulative 17.0%)

- **Where:** `internal/service/card_needs_you.go:85-199` (module `internal/service`, go)
- **Why:** cyclomatic 40 is 31 point(s) at or above the hotspot threshold 10 (39 decision points, 115 lines)

### 15. complexity `LearningService.PrepareLearningPacket` — 30.00 debt · 0.6% of scope (cumulative 17.6%)

- **Where:** `internal/service/learning_curation.go:78-217` (module `internal/service`, go)
- **Why:** cyclomatic 39 is 30 point(s) at or above the hotspot threshold 10 (38 decision points, 140 lines)

### 16. complexity `OAuthFlowService.Begin` — 30.00 debt · 0.6% of scope (cumulative 18.2%)

- **Where:** `internal/service/oauth_flow.go:153-277` (module `internal/service`, go)
- **Why:** cyclomatic 39 is 30 point(s) at or above the hotspot threshold 10 (38 decision points, 125 lines)

### 17. complexity `ReeveSnapshot` — 29.00 debt · 0.6% of scope (cumulative 18.8%)

- **Where:** `internal/service/reeve_prompt.go:97-208` (module `internal/service`, go)
- **Why:** cyclomatic 38 is 29 point(s) at or above the hotspot threshold 10 (37 decision points, 112 lines)

### 18. complexity `UsagePolicyService.evaluate` — 29.00 debt · 0.6% of scope (cumulative 19.4%)

- **Where:** `internal/service/usage_profile.go:111-227` (module `internal/service`, go)
- **Why:** cyclomatic 38 is 29 point(s) at or above the hotspot threshold 10 (37 decision points, 117 lines)

### 19. complexity `ReviewRunCaller.open` — 28.00 debt · 0.6% of scope (cumulative 20.0%)

- **Where:** `internal/service/review_pass_run.go:622-750` (module `internal/service`, go)
- **Why:** cyclomatic 37 is 28 point(s) at or above the hotspot threshold 10 (36 decision points, 129 lines)

### 20. complexity `TicketUpdateExecutor.Start` — 28.00 debt · 0.6% of scope (cumulative 20.5%)

- **Where:** `internal/service/workflow_ticket_update.go:115-229` (module `internal/service`, go)
- **Why:** cyclomatic 37 is 28 point(s) at or above the hotspot threshold 10 (36 decision points, 115 lines)

### 21. complexity `ReviewRunCaller.disputed` — 27.00 debt · 0.5% of scope (cumulative 21.1%)

- **Where:** `internal/service/review_pass_run.go:1260-1346` (module `internal/service`, go)
- **Why:** cyclomatic 36 is 27 point(s) at or above the hotspot threshold 10 (35 decision points, 87 lines)

### 22. complexity `HandoffWorkerService.Accept` — 27.00 debt · 0.5% of scope (cumulative 21.6%)

- **Where:** `internal/service/task_handoff_worker.go:190-315` (module `internal/service`, go)
- **Why:** cyclomatic 36 is 27 point(s) at or above the hotspot threshold 10 (35 decision points, 126 lines)

### 23. complexity `BurnliteConnections.Migrate` — 25.00 debt · 0.5% of scope (cumulative 22.1%)

- **Where:** `internal/service/burnlite_migration.go:66-198` (module `internal/service`, go)
- **Why:** cyclomatic 34 is 25 point(s) at or above the hotspot threshold 10 (33 decision points, 133 lines)

### 24. complexity `validateHandoffCall` — 25.00 debt · 0.5% of scope (cumulative 22.7%)

- **Where:** `internal/service/task_handoff_worker.go:141-181` (module `internal/service`, go)
- **Why:** cyclomatic 34 is 25 point(s) at or above the hotspot threshold 10 (33 decision points, 41 lines)

### 25. complexity `prTrackingService.pollInternal` — 24.00 debt · 0.5% of scope (cumulative 23.1%)

- **Where:** `internal/service/pr_tracking.go:360-470` (module `internal/service`, go)
- **Why:** cyclomatic 33 is 24 point(s) at or above the hotspot threshold 10 (32 decision points, 111 lines)

### 26. complexity `ReviewPassService.validate` — 24.00 debt · 0.5% of scope (cumulative 23.6%)

- **Where:** `internal/service/review_pass_submit.go:183-266` (module `internal/service`, go)
- **Why:** cyclomatic 33 is 24 point(s) at or above the hotspot threshold 10 (32 decision points, 84 lines)

### 27. complexity `swarmService.CreateAdhocSwarm` — 24.00 debt · 0.5% of scope (cumulative 24.1%)

- **Where:** `internal/service/swarm.go:633-757` (module `internal/service`, go)
- **Why:** cyclomatic 33 is 24 point(s) at or above the hotspot threshold 10 (32 decision points, 125 lines)

### 28. complexity `StudioNote` — 23.00 debt · 0.5% of scope (cumulative 24.6%)

- **Where:** `internal/service/agent_ui_studio.go:65-126` (module `internal/service`, go)
- **Why:** cyclomatic 32 is 23 point(s) at or above the hotspot threshold 10 (31 decision points, 62 lines)

### 29. complexity `NewCredentialBroker` — 23.00 debt · 0.5% of scope (cumulative 25.0%)

- **Where:** `internal/service/credential_broker.go:245-330` (module `internal/service`, go)
- **Why:** cyclomatic 32 is 23 point(s) at or above the hotspot threshold 10 (31 decision points, 86 lines)

### 30. complexity `ReviewService.SubmitHandoff` — 23.00 debt · 0.5% of scope (cumulative 25.5%)

- **Where:** `internal/service/handoff.go:57-170` (module `internal/service`, go)
- **Why:** cyclomatic 32 is 23 point(s) at or above the hotspot threshold 10 (31 decision points, 114 lines)

### 31. complexity `ReeveWatcher.news` — 23.00 debt · 0.5% of scope (cumulative 26.0%)

- **Where:** `internal/service/reeve_watch.go:485-586` (module `internal/service`, go)
- **Why:** cyclomatic 32 is 23 point(s) at or above the hotspot threshold 10 (31 decision points, 102 lines)

### 32. complexity `reportService.GenerateReport` — 23.00 debt · 0.5% of scope (cumulative 26.4%)

- **Where:** `internal/service/report.go:35-243` (module `internal/service`, go)
- **Why:** cyclomatic 32 is 23 point(s) at or above the hotspot threshold 10 (31 decision points, 209 lines)

### 33. complexity `ReviewPassService.routeItems` — 23.00 debt · 0.5% of scope (cumulative 26.9%)

- **Where:** `internal/service/review_steward.go:238-335` (module `internal/service`, go)
- **Why:** cyclomatic 32 is 23 point(s) at or above the hotspot threshold 10 (31 decision points, 98 lines)

### 34. complexity `conversationTail` — 21.00 debt · 0.4% of scope (cumulative 27.3%)

- **Where:** `internal/service/agent_conversation.go:174-303` (module `internal/service`, go)
- **Why:** cyclomatic 30 is 21 point(s) at or above the hotspot threshold 10 (29 decision points, 130 lines)

### 35. complexity `validateUIComponent` — 21.00 debt · 0.4% of scope (cumulative 27.8%)

- **Where:** `internal/service/agent_ui_components.go:125-203` (module `internal/service`, go)
- **Why:** cyclomatic 30 is 21 point(s) at or above the hotspot threshold 10 (29 decision points, 79 lines)

### 36. complexity `ChangeSetService.run` — 21.00 debt · 0.4% of scope (cumulative 28.2%)

- **Where:** `internal/service/reeve_changes.go:698-795` (module `internal/service`, go)
- **Why:** cyclomatic 30 is 21 point(s) at or above the hotspot threshold 10 (29 decision points, 98 lines)

### 37. complexity `WorkflowService.agentOutputs` — 21.00 debt · 0.4% of scope (cumulative 28.6%)

- **Where:** `internal/service/workflow_agent_wait.go:344-439` (module `internal/service`, go)
- **Why:** cyclomatic 30 is 21 point(s) at or above the hotspot threshold 10 (29 decision points, 96 lines)

### 38. complexity `agentService.endGoneSessionFor` — 20.00 debt · 0.4% of scope (cumulative 29.0%)

- **Where:** `internal/service/agent.go:2361-2434` (module `internal/service`, go)
- **Why:** cyclomatic 29 is 20 point(s) at or above the hotspot threshold 10 (28 decision points, 74 lines)

### 39. complexity `ReeveVoice.Status` — 20.00 debt · 0.4% of scope (cumulative 29.4%)

- **Where:** `internal/service/reeve_voice.go:160-241` (module `internal/service`, go)
- **Why:** cyclomatic 29 is 20 point(s) at or above the hotspot threshold 10 (28 decision points, 82 lines)

### 40. complexity `ExportReportMarkdown` — 20.00 debt · 0.4% of scope (cumulative 29.8%)

- **Where:** `internal/service/report_export.go:10-145` (module `internal/service`, go)
- **Why:** cyclomatic 29 is 20 point(s) at or above the hotspot threshold 10 (28 decision points, 136 lines)

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
