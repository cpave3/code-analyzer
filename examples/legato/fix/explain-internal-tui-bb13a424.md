# Debt factors: legato-stack @ `bb13a424dc73` — internal/tui

- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.2.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 1555 of 1848 production source files measured (90.0% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2; 4 file(s) with parse errors

**Scope:** 72 source files, 20003 lines, 802 callables (114 hotspots) in 13 module(s).

| Debt | complexity | nesting | duplication | architecture | Debt/kLOC | Score |
|---:|---:|---:|---:|---:|---:|---:|
| **1296.13** | 1041.00 | 98.00 | 47.13 | 110.00 | 64.80 | 27.84 |

**Concentration:** 166 factor(s) in scope; the largest 19 carry 50% of the debt and the largest 67 carry 80%.

## Factors (20 listed, 146 smaller omitted; largest first)

### 1. complexity `App.update` — 299.00 debt · 23.1% of scope (cumulative 23.1%)

- **Where:** `internal/tui/app.go:696-1827` (module `internal/tui`, go)
- **Why:** cyclomatic 308 is 299 point(s) at or above the hotspot threshold 10 (307 decision points, 1132 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 308) removes this factor; each point below 308 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `App.update` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 2. complexity `Model.HandleAction` — 41.00 debt · 3.2% of scope (cumulative 26.2%)

- **Where:** `internal/tui/agents/model.go:537-673` (module `internal/tui/agents`, go)
- **Why:** cyclomatic 50 is 41 point(s) at or above the hotspot threshold 10 (49 decision points, 137 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 50) removes this factor; each point below 50 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `Model.HandleAction` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 3. complexity `Model.HandleAction` — 27.00 debt · 2.1% of scope (cumulative 28.3%)

- **Where:** `internal/tui/attention/model.go:123-212` (module `internal/tui/attention`, go)
- **Why:** cyclomatic 36 is 27 point(s) at or above the hotspot threshold 10 (35 decision points, 90 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 36) removes this factor; each point below 36 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `Model.HandleAction` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 4. complexity `CreateOverlay.Update` — 25.00 debt · 1.9% of scope (cumulative 30.2%)

- **Where:** `internal/tui/overlay/create.go:98-184` (module `internal/tui/overlay`, go)
- **Why:** cyclomatic 34 is 25 point(s) at or above the hotspot threshold 10 (33 decision points, 87 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 34) removes this factor; each point below 34 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `CreateOverlay.Update` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 5. complexity `App.dispatchHotkey` — 24.00 debt · 1.8% of scope (cumulative 32.1%)

- **Where:** `internal/tui/app.go:3384-3467` (module `internal/tui`, go)
- **Why:** cyclomatic 33 is 24 point(s) at or above the hotspot threshold 10 (32 decision points, 84 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 33) removes this factor; each point below 33 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `App.dispatchHotkey` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

### 6. complexity `Model.handleAction` — 22.00 debt · 1.7% of scope (cumulative 33.8%)

- **Where:** `internal/tui/board/model.go:248-333` (module `internal/tui/board`, go)
- **Why:** cyclomatic 31 is 22 point(s) at or above the hotspot threshold 10 (30 decision points, 86 lines)

### 7. complexity `validate` — 22.00 debt · 1.7% of scope (cumulative 35.5%)

- **Where:** `internal/tui/hotkey/registry.go:359-429` (module `internal/tui/hotkey`, go)
- **Why:** cyclomatic 31 is 22 point(s) at or above the hotspot threshold 10 (30 decision points, 71 lines)

### 8. complexity `Model.Update` — 21.00 debt · 1.6% of scope (cumulative 37.1%)

- **Where:** `internal/tui/agents/model.go:397-511` (module `internal/tui/agents`, go)
- **Why:** cyclomatic 30 is 21 point(s) at or above the hotspot threshold 10 (29 decision points, 115 lines)

### 9. complexity `Model.HandleAction` — 20.00 debt · 1.5% of scope (cumulative 38.6%)

- **Where:** `internal/tui/detail/model.go:279-349` (module `internal/tui/detail`, go)
- **Why:** cyclomatic 29 is 20 point(s) at or above the hotspot threshold 10 (28 decision points, 71 lines)

### 10. complexity `CardContextOverlay.View` — 18.00 debt · 1.4% of scope (cumulative 40.0%)

- **Where:** `internal/tui/overlay/card_context.go:308-380` (module `internal/tui/overlay`, go)
- **Why:** cyclomatic 27 is 18 point(s) at or above the hotspot threshold 10 (26 decision points, 73 lines)

### 11. complexity `CardContextOverlay.updateKey` — 18.00 debt · 1.4% of scope (cumulative 41.4%)

- **Where:** `internal/tui/overlay/card_context.go:152-219` (module `internal/tui/overlay`, go)
- **Why:** cyclomatic 27 is 18 point(s) at or above the hotspot threshold 10 (26 decision points, 68 lines)

### 12. complexity `Model.renderSidebarEntry` — 17.00 debt · 1.3% of scope (cumulative 42.7%)

- **Where:** `internal/tui/agents/list.go:372-459` (module `internal/tui/agents`, go)
- **Why:** cyclomatic 26 is 17 point(s) at or above the hotspot threshold 10 (25 decision points, 88 lines)

### 13. complexity `AgentSpawnOverlay.handleKey` — 17.00 debt · 1.3% of scope (cumulative 44.0%)

- **Where:** `internal/tui/overlay/agent_spawn.go:506-570` (module `internal/tui/overlay`, go)
- **Why:** cyclomatic 26 is 17 point(s) at or above the hotspot threshold 10 (25 decision points, 65 lines)

### 14. complexity `Model.renderHeader` — 16.00 debt · 1.2% of scope (cumulative 45.3%)

- **Where:** `internal/tui/detail/model.go:663-762` (module `internal/tui/detail`, go)
- **Why:** cyclomatic 25 is 16 point(s) at or above the hotspot threshold 10 (24 decision points, 100 lines)

### 15. complexity `WorkflowTriggerOverlay.Update` — 15.00 debt · 1.2% of scope (cumulative 46.5%)

- **Where:** `internal/tui/overlay/workflow_trigger.go:46-108` (module `internal/tui/overlay`, go)
- **Why:** cyclomatic 24 is 15 point(s) at or above the hotspot threshold 10 (23 decision points, 63 lines)

### 16. complexity `Model.HandleAction` — 15.00 debt · 1.2% of scope (cumulative 47.6%)

- **Where:** `internal/tui/plan/model.go:148-204` (module `internal/tui/plan`, go)
- **Why:** cyclomatic 24 is 15 point(s) at or above the hotspot threshold 10 (23 decision points, 57 lines)

### 17. complexity `Model.panelTask` — 12.00 debt · 0.9% of scope (cumulative 48.5%)

- **Where:** `internal/tui/board/view.go:352-431` (module `internal/tui/board`, go)
- **Why:** cyclomatic 21 is 12 point(s) at or above the hotspot threshold 10 (20 decision points, 80 lines)

### 18. complexity `GroupOverlay.Update` — 12.00 debt · 0.9% of scope (cumulative 49.5%)

- **Where:** `internal/tui/overlay/group.go:51-98` (module `internal/tui/overlay`, go)
- **Why:** cyclomatic 21 is 12 point(s) at or above the hotspot threshold 10 (20 decision points, 48 lines)

### 19. duplication `buildMWShortcuts` — 10.45 debt · 0.8% of scope (cumulative 50.3%)

- **Where:** `internal/tui/overlay/move_workspace.go:72-109` (module `internal/tui/overlay`, go)
- **Why:** structural clone group of 2 callables (2 in scope); 209 redundant tokens attributed to this scope
- **Related:** `buildMWShortcuts` `internal/tui/overlay/move_workspace.go:72-109`; `buildWSShortcuts` `internal/tui/overlay/workspace.go:72-110`

### 20. complexity `Model.conversationWithRuns` — 10.00 debt · 0.8% of scope (cumulative 51.0%)

- **Where:** `internal/tui/agents/conversation.go:202-278` (module `internal/tui/agents`, go)
- **Why:** cyclomatic 19 is 10 point(s) at or above the hotspot threshold 10 (18 decision points, 77 lines)

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
