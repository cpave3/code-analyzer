# Debt factors: legato-stack @ `bb13a424dc73` — web

- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.2.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 1555 of 1848 production source files measured (90.0% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2; 4 file(s) with parse errors

**Scope:** 459 source files, 63542 lines, 2152 callables (405 hotspots) in 17 module(s).

| Debt | complexity | nesting | duplication | architecture | Debt/kLOC | Score |
|---:|---:|---:|---:|---:|---:|---:|
| **6966.17** | 6493.00 | 192.00 | 101.17 | 180.00 | 109.63 | 18.57 |

**Concentration:** 502 factor(s) in scope; the largest 45 carry 50% of the debt and the largest 155 carry 80%.

## Factors (40 listed, 462 smaller omitted; largest first)

### 1. complexity `NativeAgentPanel` — 322.00 debt · 4.6% of scope (cumulative 4.6%)

- **Where:** `web/src/components/NativeAgentPanel.tsx:270-1453` (module `web/src/components`, typescript)
- **Why:** cyclomatic 331 is 322 point(s) at or above the hotspot threshold 10 (330 decision points, 1184 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 331) removes this factor; each point below 331 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `NativeAgentPanel` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 2. complexity `BoardPage` — 288.00 debt · 4.1% of scope (cumulative 8.8%)

- **Where:** `web/src/pages/Board.tsx:190-1424` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 297 is 288 point(s) at or above the hotspot threshold 10 (296 decision points, 1235 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 297) removes this factor; each point below 297 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `BoardPage` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 3. complexity `SpawnAgentDialog` — 241.00 debt · 3.5% of scope (cumulative 12.2%)

- **Where:** `web/src/components/SpawnAgentDialog.tsx:104-1050` (module `web/src/components`, typescript)
- **Why:** cyclomatic 250 is 241 point(s) at or above the hotspot threshold 10 (249 decision points, 947 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 250) removes this factor; each point below 250 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `SpawnAgentDialog` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 4. complexity `AgentsPage` — 238.00 debt · 3.4% of scope (cumulative 15.6%)

- **Where:** `web/src/pages/Agents.tsx:113-1128` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 247 is 238 point(s) at or above the hotspot threshold 10 (246 decision points, 1016 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 247) removes this factor; each point below 247 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `AgentsPage` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 5. complexity `Outline` — 195.00 debt · 2.8% of scope (cumulative 18.4%)

- **Where:** `web/src/components/workflow-editor/WorkflowOutline.tsx:194-647` (module `web/src/components/workflow-editor`, typescript)
- **Why:** cyclomatic 204 is 195 point(s) at or above the hotspot threshold 10 (203 decision points, 454 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 204) removes this factor; each point below 204 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `Outline` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 6. complexity `WorkflowEditPage` — 184.00 debt · 2.6% of scope (cumulative 21.1%)

- **Where:** `web/src/pages/WorkflowEditPage.tsx:166-595` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 193 is 184 point(s) at or above the hotspot threshold 10 (192 decision points, 430 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 193) removes this factor; each point below 193 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `WorkflowEditPage` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 7. complexity `ReviewTourPage` — 100.00 debt · 1.4% of scope (cumulative 22.5%)

- **Where:** `web/src/pages/ReviewTourPage.tsx:53-477` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 109 is 100 point(s) at or above the hotspot threshold 10 (108 decision points, 425 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 109) removes this factor; each point below 109 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `ReviewTourPage` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 8. complexity `reduceNativeEvent` — 87.00 debt · 1.2% of scope (cumulative 23.8%)

- **Where:** `web/src/lib/nativeAgentSync.ts:573-797` (module `web/src/lib`, typescript)
- **Why:** cyclomatic 96 is 87 point(s) at or above the hotspot threshold 10 (95 decision points, 225 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 96) removes this factor; each point below 96 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `reduceNativeEvent` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 9. complexity `SettingsPage` — 81.00 debt · 1.2% of scope (cumulative 24.9%)

- **Where:** `web/src/pages/Settings.tsx:102-343` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 90 is 81 point(s) at or above the hotspot threshold 10 (89 decision points, 242 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 90) removes this factor; each point below 90 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `SettingsPage` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 10. complexity `runNeedsYou` — 78.00 debt · 1.1% of scope (cumulative 26.0%)

- **Where:** `web/src/lib/runNeedsYou.ts:73-165` (module `web/src/lib`, typescript)
- **Why:** cyclomatic 87 is 78 point(s) at or above the hotspot threshold 10 (86 decision points, 93 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 87) removes this factor; each point below 87 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `runNeedsYou` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 11. complexity `AgentArtifactPanel` — 70.00 debt · 1.0% of scope (cumulative 27.0%)

- **Where:** `web/src/components/AgentArtifactPanel.tsx:69-430` (module `web/src/components`, typescript)
- **Why:** cyclomatic 79 is 70 point(s) at or above the hotspot threshold 10 (78 decision points, 362 lines)

### 12. complexity `QuestionFocus` — 68.00 debt · 1.0% of scope (cumulative 28.0%)

- **Where:** `web/src/components/QuestionFocus.tsx:25-292` (module `web/src/components`, typescript)
- **Why:** cyclomatic 77 is 68 point(s) at or above the hotspot threshold 10 (76 decision points, 268 lines)

### 13. complexity `LearningDetail` — 68.00 debt · 1.0% of scope (cumulative 29.0%)

- **Where:** `web/src/pages/LearningDetailPage.tsx:18-146` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 77 is 68 point(s) at or above the hotspot threshold 10 (76 decision points, 129 lines)

### 14. complexity `WorkflowRunPage` — 68.00 debt · 1.0% of scope (cumulative 30.0%)

- **Where:** `web/src/pages/WorkflowRunPage.tsx:694-839` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 77 is 68 point(s) at or above the hotspot threshold 10 (76 decision points, 146 lines)

### 15. complexity `stateFromNativeSnapshot` — 66.00 debt · 0.9% of scope (cumulative 30.9%)

- **Where:** `web/src/lib/nativeAgentSync.ts:327-456` (module `web/src/lib`, typescript)
- **Why:** cyclomatic 75 is 66 point(s) at or above the hotspot threshold 10 (74 decision points, 130 lines)

### 16. complexity `useMenuLayer` — 63.00 debt · 0.9% of scope (cumulative 31.8%)

- **Where:** `web/src/hooks/useMenuLayer.ts:19-152` (module `web/src/hooks`, typescript)
- **Why:** cyclomatic 72 is 63 point(s) at or above the hotspot threshold 10 (71 decision points, 134 lines)

### 17. complexity `ToolCallCard` — 56.00 debt · 0.8% of scope (cumulative 32.6%)

- **Where:** `web/src/components/NativeAgentPanel.tsx:1686-1914` (module `web/src/components`, typescript)
- **Why:** cyclomatic 65 is 56 point(s) at or above the hotspot threshold 10 (64 decision points, 229 lines)

### 18. complexity `WithField` — 55.00 debt · 0.8% of scope (cumulative 33.4%)

- **Where:** `web/src/components/workflow-editor/StepsEditor.tsx:186-258` (module `web/src/components/workflow-editor`, typescript)
- **Why:** cyclomatic 64 is 55 point(s) at or above the hotspot threshold 10 (63 decision points, 73 lines)

### 19. complexity `Connection` — 53.00 debt · 0.8% of scope (cumulative 34.2%)

- **Where:** `web/src/components/settings/MCPSettings.tsx:159-424` (module `web/src/components/settings`, typescript)
- **Why:** cyclomatic 62 is 53 point(s) at or above the hotspot threshold 10 (61 decision points, 266 lines)

### 20. complexity `Now` — 53.00 debt · 0.8% of scope (cumulative 34.9%)

- **Where:** `web/src/pages/WorkflowRunPage.tsx:90-202` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 62 is 53 point(s) at or above the hotspot threshold 10 (61 decision points, 113 lines)

### 21. complexity `pathType` — 52.00 debt · 0.8% of scope (cumulative 35.7%)

- **Where:** `web/src/lib/workflowExpressions.ts:202-264` (module `web/src/lib`, typescript)
- **Why:** cyclomatic 61 is 52 point(s) at or above the hotspot threshold 10 (60 decision points, 63 lines)

### 22. complexity `AgentGraphView` — 50.00 debt · 0.7% of scope (cumulative 36.4%)

- **Where:** `web/src/components/AgentGraphView.tsx:58-269` (module `web/src/components`, typescript)
- **Why:** cyclomatic 59 is 50 point(s) at or above the hotspot threshold 10 (58 decision points, 212 lines)

### 23. complexity `AgentSelect` — 50.00 debt · 0.7% of scope (cumulative 37.1%)

- **Where:** `web/src/components/AgentSelect.tsx:62-265` (module `web/src/components`, typescript)
- **Why:** cyclomatic 59 is 50 point(s) at or above the hotspot threshold 10 (58 decision points, 204 lines)

### 24. complexity `MacroForm` — 50.00 debt · 0.7% of scope (cumulative 37.8%)

- **Where:** `web/src/components/MacroForm.tsx:38-153` (module `web/src/components`, typescript)
- **Why:** cyclomatic 59 is 50 point(s) at or above the hotspot threshold 10 (58 decision points, 116 lines)

### 25. complexity `HandOffModal` — 50.00 debt · 0.7% of scope (cumulative 38.6%)

- **Where:** `web/src/components/board/HandOffModal.tsx:78-276` (module `web/src/components/board`, typescript)
- **Why:** cyclomatic 59 is 50 point(s) at or above the hotspot threshold 10 (58 decision points, 199 lines)

### 26. complexity `layoutAgentGraph` — 49.00 debt · 0.7% of scope (cumulative 39.3%)

- **Where:** `web/src/lib/agentGraph.ts:265-430` (module `web/src/lib`, typescript)
- **Why:** cyclomatic 58 is 49 point(s) at or above the hotspot threshold 10 (57 decision points, 166 lines)

### 27. complexity `NativeModelSettingsForm` — 48.00 debt · 0.7% of scope (cumulative 40.0%)

- **Where:** `web/src/components/NativeModelSettings.tsx:26-122` (module `web/src/components`, typescript)
- **Why:** cyclomatic 57 is 48 point(s) at or above the hotspot threshold 10 (56 decision points, 97 lines)

### 28. complexity `EntryForm` — 48.00 debt · 0.7% of scope (cumulative 40.6%)

- **Where:** `web/src/components/settings/ListSettings.tsx:162-276` (module `web/src/components/settings`, typescript)
- **Why:** cyclomatic 57 is 48 point(s) at or above the hotspot threshold 10 (56 decision points, 115 lines)

### 29. complexity `Inspector` — 45.00 debt · 0.7% of scope (cumulative 41.3%)

- **Where:** `web/src/components/workflow-editor/WorkflowOutline.tsx:676-782` (module `web/src/components/workflow-editor`, typescript)
- **Why:** cyclomatic 54 is 45 point(s) at or above the hotspot threshold 10 (53 decision points, 107 lines)

### 30. complexity `PromptBar` — 43.00 debt · 0.6% of scope (cumulative 41.9%)

- **Where:** `web/src/components/PromptBar.tsx:167-419` (module `web/src/components`, typescript)
- **Why:** cyclomatic 52 is 43 point(s) at or above the hotspot threshold 10 (51 decision points, 253 lines)

### 31. complexity `PlanPage` — 43.00 debt · 0.6% of scope (cumulative 42.5%)

- **Where:** `web/src/pages/PlanPage.tsx:210-367` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 52 is 43 point(s) at or above the hotspot threshold 10 (51 decision points, 158 lines)

### 32. complexity `RunWorkflowDialog` — 42.00 debt · 0.6% of scope (cumulative 43.1%)

- **Where:** `web/src/components/RunWorkflowDialog.tsx:8-76` (module `web/src/components`, typescript)
- **Why:** cyclomatic 51 is 42 point(s) at or above the hotspot threshold 10 (50 decision points, 69 lines)

### 33. complexity `WorkflowPage` — 42.00 debt · 0.6% of scope (cumulative 43.7%)

- **Where:** `web/src/pages/WorkflowPage.tsx:164-286` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 51 is 42 point(s) at or above the hotspot threshold 10 (50 decision points, 123 lines)

### 34. complexity `ApprovalCard` — 41.00 debt · 0.6% of scope (cumulative 44.3%)

- **Where:** `web/src/components/ApprovalTray.tsx:122-288` (module `web/src/components`, typescript)
- **Why:** cyclomatic 50 is 41 point(s) at or above the hotspot threshold 10 (49 decision points, 167 lines)

### 35. complexity `AgentPickerControl` — 41.00 debt · 0.6% of scope (cumulative 44.9%)

- **Where:** `web/src/components/settings/Pickers.tsx:44-102` (module `web/src/components/settings`, typescript)
- **Why:** cyclomatic 50 is 41 point(s) at or above the hotspot threshold 10 (49 decision points, 59 lines)

### 36. complexity `route` — 41.00 debt · 0.6% of scope (cumulative 45.5%)

- **Where:** `web/src/lib/agentGraph.ts:588-738` (module `web/src/lib`, typescript)
- **Why:** cyclomatic 50 is 41 point(s) at or above the hotspot threshold 10 (49 decision points, 151 lines)

### 37. complexity `check` — 41.00 debt · 0.6% of scope (cumulative 46.1%)

- **Where:** `web/src/lib/workflowSchema.ts:369-440` (module `web/src/lib`, typescript)
- **Why:** cyclomatic 50 is 41 point(s) at or above the hotspot threshold 10 (49 decision points, 72 lines)

### 38. complexity `DetailPanel` — 40.00 debt · 0.6% of scope (cumulative 46.7%)

- **Where:** `web/src/components/board/DetailPanel.tsx:84-243` (module `web/src/components/board`, typescript)
- **Why:** cyclomatic 49 is 40 point(s) at or above the hotspot threshold 10 (48 decision points, 160 lines)

### 39. complexity `CardContextMenu` — 39.00 debt · 0.6% of scope (cumulative 47.2%)

- **Where:** `web/src/components/board/CardContextMenu.tsx:51-198` (module `web/src/components/board`, typescript)
- **Why:** cyclomatic 48 is 39 point(s) at or above the hotspot threshold 10 (47 decision points, 148 lines)

### 40. complexity `SessionDetailPanel` — 39.00 debt · 0.6% of scope (cumulative 47.8%)

- **Where:** `web/src/pages/SessionHistoryPage.tsx:148-263` (module `web/src/pages`, typescript)
- **Why:** cyclomatic 48 is 39 point(s) at or above the hotspot threshold 10 (47 decision points, 116 lines)

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
