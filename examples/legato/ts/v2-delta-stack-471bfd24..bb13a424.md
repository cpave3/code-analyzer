# CodeEval v2 delta: legato-stack `471bfd24a696` → `bb13a424dc73`

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `bb13a424dc73f3eeb98f8bce43db932a9d19f7e6` tree `25b725950d6a` 2026-11-27T20:29:04+11:00
- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.2.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 1555 of 1848 production source files measured (90.0% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2; 4 file(s) with parse errors

**Change size:** 37 production source file(s), 37 measured (1698 changed lines)

## Ledger (after − before; telescopes across commits)

| | Before | After | Δ | Δ per changed kLOC |
|---|---:|---:|---:|---:|
| **Total debt** | 19814.75 | 19590.65 | **-224.10** | -131.98 |
| complexity debt | 17626.00 | 17487.00 | -139.00 | -81.86 |
| nesting debt | 860.00 | 854.00 | -6.00 | -3.53 |
| duplication debt | 918.75 | 879.65 | -39.10 | -23.03 |
| architecture debt | 410.00 | 370.00 | -40.00 | -23.56 |
| score (0–100) | 28.54 | 28.80 | +0.26 | |
| debt / kLOC | 62.59 | 61.82 | -0.77 | |

| Language | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |
|---|---|---:|---:|---:|---:|---:|---|
| go | 11611.30 → 11387.20 | -224.10 | -139.00 | -6.00 | -39.10 | -40.00 | 33.93 → 34.40 |
| typescript | 8203.45 → 8203.45 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | 19.21 → 19.21 |

| Raw component | Before | After | Δ |
|---|---:|---:|---:|
| source lines | 316568 | 316900 | +332 |
| callables | 11968 | 12006 | +38 |
| decision points | 54640 | 54616 | -24 |
| cyclomatic total | 66608 | 66622 | +14 |
| hotspots | 1776 | 1773 | -3 |
| cyclomatic excess | 17626 | 17487 | -139 |
| nesting excess | 430 | 427 | -3 |
| exact redundant callables | 20 | 18 | -2 |
| structural redundant callables | 295 | 288 | -7 |
| rule-violating edges | 15 | 11 | -4 |
| units in cycles | 13 | 13 | +0 |
| internal edges | 279 | 278 | -1 |
| redundant tokens | 18375.00 | 17593.00 | -782.00 |

## Changed callables: decision points versus function base points

- modified pre-existing: 60, decision points 755 → 544 (-211), cyclomatic 815 → 604
- added: 43 (decision points 202); deleted: 5 (decision points 15)
- **changed scope: decision points -24, function base points +38** (cyclomatic = decisions + 1 per function)

## Advisory gates

| Gate | Result | Detail |
|---|---|---|
| no_new_exact_duplicates | pass | -2 exact redundant callables |
| no_new_structural_duplicates | pass | -7 structural redundant callables |
| no_new_rule_violations | pass | -4 rule-violating edges |
| no_new_cycles | pass | 0 units in cycles |
| no_branch_growth_in_existing_callables | pass | -211 decision points in modified pre-existing callables |
| no_hotspot_excess_growth | pass | -139 cyclomatic excess, -3 nesting excess |

## Module changes (10 changed; worst first)

| Module | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |
|---|---|---:|---:|---:|---:|---:|---|
| `internal/engine/codexusage` | 18.30 → 16.18 | -2.13 | +0.00 | +0.00 | -2.13 | +0.00 | 42.95 → 44.75 |
| `internal/engine/opencodeusage` | 14.38 → 12.00 | -2.38 | +0.00 | +0.00 | -2.38 | +0.00 | 31.71 → 34.60 |
| `internal/engine/claudeusage` | 16.68 → 12.18 | -4.50 | +0.00 | +0.00 | -4.50 | +0.00 | 40.45 → 45.89 |
| `internal/handofftest` | 51.05 → 45.05 | -6.00 | +0.00 | +0.00 | -6.00 | +0.00 | 53.20 → 56.12 |
| `config` | 256.10 → 246.10 | -10.00 | +0.00 | +0.00 | +0.00 | -10.00 | 34.05 → 35.48 |
| `internal/tui/overlay` | 371.20 → 361.20 | -10.00 | +0.00 | +0.00 | +0.00 | -10.00 | 27.96 → 28.51 |
| `internal/engine/workflow` | 597.00 → 582.00 | -15.00 | -11.00 | -4.00 | +0.00 | +0.00 | 28.78 → 29.35 |
| `cmd/legato` | 888.13 → 871.93 | -16.20 | +0.00 | +0.00 | -16.20 | +0.00 | 34.35 → 34.75 |
| `internal/tui` | 448.20 → 428.20 | -20.00 | +0.00 | +0.00 | +0.00 | -20.00 | 21.17 → 21.94 |
| `internal/service` | 5070.35 → 4932.45 | -137.90 | -128.00 | -2.00 | -7.90 | +0.00 | 32.14 → 32.82 |

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
