# CodeEval v2 delta: legato-stack `9a61e9f8a021` → `bb13a424dc73`

- **Base:** `9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc` tree `a55db6103bf4` 2026-11-25T21:49:11+11:00
- **Target:** `bb13a424dc73f3eeb98f8bce43db932a9d19f7e6` tree `25b725950d6a` 2026-11-27T20:29:04+11:00
- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.1.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 990 of 1860 production source files measured (61.9% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2, typescript 577

**Change size:** 301 production source file(s), 206 measured (17229 changed lines); not measured: javascript 2, kotlin 2, sql 7, typescript 84

## Ledger (after − before; telescopes across commits)

| | Before | After | Δ | Δ per changed kLOC |
|---|---:|---:|---:|---:|
| **Total debt** | 10783.30 | 11387.20 | **+603.90** | +35.05 |
| complexity debt | 9275.00 | 9888.00 | +613.00 | +35.58 |
| nesting debt | 598.00 | 612.00 | +14.00 | +0.81 |
| duplication debt | 760.30 | 777.20 | +16.90 | +0.98 |
| architecture debt | 150.00 | 110.00 | -40.00 | -2.32 |
| score (0–100) | 34.19 | 34.40 | +0.21 | |
| debt / kLOC | 48.12 | 47.67 | -0.45 | |

| Raw component | Before | After | Δ |
|---|---:|---:|---:|
| source lines | 224082 | 238865 | +14783 |
| callables | 8734 | 9336 | +602 |
| decision points | 35801 | 38363 | +2562 |
| cyclomatic total | 44535 | 47699 | +3164 |
| hotspots | 1197 | 1285 | +88 |
| cyclomatic excess | 9275 | 9888 | +613 |
| nesting excess | 299 | 306 | +7 |
| exact redundant callables | 8 | 8 | +0 |
| structural redundant callables | 243 | 253 | +10 |
| rule-violating edges | 15 | 11 | -4 |
| units in cycles | 0 | 0 | +0 |
| internal edges | 193 | 199 | +6 |
| redundant tokens | 15206.00 | 15544.00 | +338.00 |

## Changed callables: decision points versus function base points

- modified pre-existing: 267, decision points 3352 → 3446 (+94), cyclomatic 3619 → 3713
- added: 611 (decision points 2490); deleted: 9 (decision points 22)
- **changed scope: decision points +2562, function base points +602** (cyclomatic = decisions + 1 per function)

## Advisory gates

| Gate | Result | Detail |
|---|---|---|
| no_new_exact_duplicates | pass | 0 exact redundant callables |
| no_new_structural_duplicates | **fail** | +10 structural redundant callables |
| no_new_rule_violations | pass | -4 rule-violating edges |
| no_new_cycles | pass | 0 units in cycles |
| no_branch_growth_in_existing_callables | **fail** | +94 decision points in modified pre-existing callables |
| no_hotspot_excess_growth | **fail** | +613 cyclomatic excess, +7 nesting excess |

## Module changes (16 changed; worst first)

| Module | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |
|---|---|---:|---:|---:|---:|---:|---|
| `internal/service` | 4523.95 → 4932.45 | +408.50 | +394.00 | +4.00 | +10.50 | +0.00 | 32.69 → 32.82 |
| `internal/engine/workflow` | 500.00 → 582.00 | +82.00 | +74.00 | +8.00 | +0.00 | +0.00 | 27.47 → 29.35 |
| `cmd/legato` | 844.03 → 871.93 | +27.90 | +36.00 | +0.00 | -8.10 | +0.00 | 34.39 → 34.75 |
| `config` | 224.80 → 246.10 | +21.30 | +28.00 | +0.00 | +3.30 | -10.00 | 34.69 → 35.48 |
| `internal/engine/codexusage` | 0.00 → 16.18 | +16.18 | +13.00 | +2.00 | +1.18 | +0.00 | n/a → 44.75 |
| `internal/engine/claudeusage` | 0.00 → 12.18 | +12.18 | +11.00 | +0.00 | +1.18 | +0.00 | n/a → 45.89 |
| `internal/engine/opencodeusage` | 0.00 → 12.00 | +12.00 | +12.00 | +0.00 | +0.00 | +0.00 | n/a → 34.60 |
| `internal/engine/store` | 824.90 → 833.90 | +9.00 | +10.00 | +0.00 | -1.00 | +0.00 | 37.29 → 37.49 |
| `internal/server` | 514.05 → 522.80 | +8.75 | +7.00 | +0.00 | +1.75 | +0.00 | 48.70 → 48.77 |
| `internal/engine/burnlite` | 0.00 → 4.80 | +4.80 | +0.00 | +0.00 | +4.80 | +0.00 | 100.00 → 77.70 |
| `internal/engine/jira` | 70.00 → 74.70 | +4.70 | +1.00 | +0.00 | +3.70 | +0.00 | 32.71 → 32.43 |
| `internal/tui/board` | 69.00 → 73.00 | +4.00 | +4.00 | +0.00 | +0.00 | +0.00 | 35.51 → 34.46 |
| `internal/handofftest` | 41.25 → 45.05 | +3.80 | +7.00 | +0.00 | -3.20 | +0.00 | 54.74 → 56.12 |
| `internal/tui/agents` | 146.15 → 148.15 | +2.00 | +2.00 | +0.00 | +0.00 | +0.00 | 32.94 → 33.16 |
| `internal/tui/overlay` | 361.20 → 361.20 | +0.00 | +10.00 | +0.00 | +0.00 | -10.00 | 28.18 → 28.51 |
| `internal/tui` | 441.40 → 428.20 | -13.20 | +4.00 | +0.00 | +2.80 | -20.00 | 20.98 → 21.94 |

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
