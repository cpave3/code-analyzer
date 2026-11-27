# CodeEval v2 delta: legato-stack `a3ee21f209d4` → `4cf88e484251`

- **Base:** `a3ee21f209d477ae0e574c4375aa3aaf09dc89f0` tree `593122bfffc2` 2026-11-27T20:29:04+11:00
- **Target:** `4cf88e48425163499193263ffec92c72adbe6c82` tree `4fbb5b2b4feb` 2026-11-27T20:29:04+11:00
- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.1.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 990 of 1860 production source files measured (61.9% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2, typescript 577

**Change size:** 1 production source file(s), 1 measured (114 changed lines)

## Ledger (after − before; telescopes across commits)

| | Before | After | Δ | Δ per changed kLOC |
|---|---:|---:|---:|---:|
| **Total debt** | 11406.20 | 11395.20 | **-11.00** | -96.49 |
| complexity debt | 9907.00 | 9896.00 | -11.00 | -96.49 |
| nesting debt | 612.00 | 612.00 | +0.00 | +0.00 |
| duplication debt | 777.20 | 777.20 | +0.00 | +0.00 |
| architecture debt | 110.00 | 110.00 | +0.00 | +0.00 |
| score (0–100) | 34.36 | 34.38 | +0.02 | |
| debt / kLOC | 47.76 | 47.71 | -0.05 | |

| Raw component | Before | After | Δ |
|---|---:|---:|---:|
| source lines | 238812 | 238852 | +40 |
| callables | 9329 | 9335 | +6 |
| decision points | 38362 | 38363 | +1 |
| cyclomatic total | 47691 | 47698 | +7 |
| hotspots | 1284 | 1284 | +0 |
| cyclomatic excess | 9907 | 9896 | -11 |
| nesting excess | 306 | 306 | +0 |
| exact redundant callables | 8 | 8 | +0 |
| structural redundant callables | 253 | 253 | +0 |
| rule-violating edges | 11 | 11 | +0 |
| units in cycles | 0 | 0 | +0 |
| internal edges | 199 | 199 | +0 |
| redundant tokens | 15544.00 | 15544.00 | +0.00 |

## Changed callables: decision points versus function base points

- modified pre-existing: 1, decision points 106 → 95 (-11), cyclomatic 107 → 96
- added: 6 (decision points 12); deleted: 0 (decision points 0)
- **changed scope: decision points +1, function base points +6** (cyclomatic = decisions + 1 per function)

## Advisory gates

| Gate | Result | Detail |
|---|---|---|
| no_new_exact_duplicates | pass | 0 exact redundant callables |
| no_new_structural_duplicates | pass | 0 structural redundant callables |
| no_new_rule_violations | pass | 0 rule-violating edges |
| no_new_cycles | pass | 0 units in cycles |
| no_branch_growth_in_existing_callables | pass | -11 decision points in modified pre-existing callables |
| no_hotspot_excess_growth | pass | -11 cyclomatic excess, 0 nesting excess |

## Module changes (1 changed; worst first)

| Module | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |
|---|---|---:|---:|---:|---:|---:|---|
| `internal/service` | 4951.45 → 4940.45 | -11.00 | -11.00 | +0.00 | +0.00 | +0.00 | 32.72 → 32.78 |

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
