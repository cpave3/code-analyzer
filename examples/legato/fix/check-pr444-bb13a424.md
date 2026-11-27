# codeeval check: PASS (`4cf88e484251` → `bb13a424dc73`, whole repository)

- **Fail on:** `net-debt` · scoring `scoring/v2-experimental` (experimental, uncalibrated)
- **Change:** 1 measured source file(s), 63 changed lines
- **Debt in scope:** 19598.65 → 19590.65 (-8.00)
- **Repository debt:** -8.00 (complexity -8.00, nesting +0.00, duplication +0.00, architecture +0.00)
- **Changed callables:** decision points +0, function base points +1; modified pre-existing 1 (decisions 66 → 57), added 1

| Gate | Result |
|---|---|
| no_new_exact_duplicates | pass |
| no_new_structural_duplicates | pass |
| no_new_rule_violations | pass |
| no_new_cycles | pass |
| no_branch_growth_in_existing_callables | pass |
| no_hotspot_excess_growth | pass |

## Fix these (1 introduced or worsened; largest first)

### 1. complexity `withClaudeAccounts` — 1.00 debt · **introduced** (Δ +1.00)

- **Where:** `internal/service/agent_sandbox.go:443-470` (module `internal/service`, go)
- **Why:** cyclomatic 10 is 1 point(s) at or above the hotspot threshold 10 (9 decision points, 28 lines)
- **Fix:** reduce branching in this callable. Target: cyclomatic ≤ 9 (now 10) removes this factor; each point below 10 removes 1 debt point.
  - identify independent responsibilities (validation, mapping, side effects) and extract each into a cohesive, well-named function
  - replace repeated condition chains with early returns / guard clauses, or a lookup only where the mapping is genuinely data
  - remove duplicated or dead branches; merge conditions that lead to the same outcome
  - keep each extracted helper below the hotspot threshold
- **Done when:** `withClaudeAccounts` cyclomatic ≤ 9, or a justified won't-fix; no new hotspot callables introduced in the same change; decision points across the changed callables do not increase; existing and characterisation tests pass
- **Caution:** large roots usually need several focused changes; split work per responsibility and verify between steps

## Improved by this change (1)

- complexity `agentService.sandboxPrefix` at `internal/service/agent_sandbox.go:48-263`: reduced (Δ -9.00)

## How to fix

- preserve behaviour: add or confirm characterisation tests for a factor's code before changing it; the evaluator does not run code
- reduce decision points, do not only move them: a split that leaves the same branching in helpers removes hotspot debt but not logic (check the decision-point ledger with `codeeval check`)
- keep helpers cohesive and named for one purpose; avoid creating new hotspot helpers (cc ≥ threshold) while removing a root's debt
- do not delete features, inline configuration into data tables, or reshape syntax solely to lower the number; debt is an uncalibrated structural proxy
- verify each fix by re-running `codeeval explain` on the same scope (the factor should disappear or shrink) and `codeeval check` on the change (gates pass, no new factors)

This check measures structure only. It does not run code or establish correctness.
