# CodeEval v2 delta: legato-stack `9a61e9f8a021` → `471bfd24a696`

- **Base:** `9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc` tree `a55db6103bf4` 2026-11-25T21:49:11+11:00
- **Target:** `471bfd24a69628d985cdd7767dcb340da6a569b7` tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.2.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 1553 of 1846 production source files measured (90.0% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2; 4 file(s) with parse errors

**Change size:** 293 production source file(s), 282 measured (21906 changed lines); not measured: javascript 2, kotlin 2, sql 7

## Ledger (after − before; telescopes across commits)

| | Before | After | Δ | Δ per changed kLOC |
|---|---:|---:|---:|---:|
| **Total debt** | 18387.20 | 19814.75 | **+1427.55** | +65.17 |
| complexity debt | 16287.00 | 17626.00 | +1339.00 | +61.12 |
| nesting debt | 830.00 | 860.00 | +30.00 | +1.37 |
| duplication debt | 860.20 | 918.75 | +58.55 | +2.67 |
| architecture debt | 410.00 | 410.00 | +0.00 | +0.00 |
| score (0–100) | 28.84 | 28.54 | -0.30 | |
| debt / kLOC | 61.69 | 62.59 | +0.90 | |

| Language | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |
|---|---|---:|---:|---:|---:|---:|---|
| go | 10783.30 → 11611.30 | +828.00 | +752.00 | +20.00 | +56.00 | +0.00 | 34.19 → 33.93 |
| typescript | 7603.90 → 8203.45 | +599.55 | +587.00 | +10.00 | +2.55 | +0.00 | 19.57 → 19.21 |

| Raw component | Before | After | Δ |
|---|---:|---:|---:|
| source lines | 298071 | 316568 | +18497 |
| callables | 11238 | 11968 | +730 |
| decision points | 50926 | 54640 | +3714 |
| cyclomatic total | 62164 | 66608 | +4444 |
| hotspots | 1656 | 1776 | +120 |
| cyclomatic excess | 16287 | 17626 | +1339 |
| nesting excess | 415 | 430 | +15 |
| exact redundant callables | 17 | 20 | +3 |
| structural redundant callables | 277 | 295 | +18 |
| rule-violating edges | 15 | 15 | +0 |
| units in cycles | 13 | 13 | +0 |
| internal edges | 262 | 279 | +17 |
| redundant tokens | 17204.00 | 18375.00 | +1171.00 |

## Changed callables: decision points versus function base points

- modified pre-existing: 338, decision points 5882 → 6541 (+659), cyclomatic 6220 → 6879
- added: 741 (decision points 3094); deleted: 11 (decision points 39)
- **changed scope: decision points +3714, function base points +730** (cyclomatic = decisions + 1 per function)

## Advisory gates

| Gate | Result | Detail |
|---|---|---|
| no_new_exact_duplicates | **fail** | +3 exact redundant callables |
| no_new_structural_duplicates | **fail** | +18 structural redundant callables |
| no_new_rule_violations | pass | 0 rule-violating edges |
| no_new_cycles | pass | 0 units in cycles |
| no_branch_growth_in_existing_callables | **fail** | +659 decision points in modified pre-existing callables |
| no_hotspot_excess_growth | **fail** | +1339 cyclomatic excess, +15 nesting excess |

## Module changes (28 changed; the 13 worst and 12 best)

| Module | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |
|---|---|---:|---:|---:|---:|---:|---|
| `internal/service` | 4523.95 → 5070.35 | +546.40 | +522.00 | +6.00 | +18.40 | +0.00 | 32.69 → 32.14 |
| `packages/even-hub/src` | 221.00 → 355.28 | +134.27 | +133.00 | +0.00 | +1.27 | +0.00 | 20.36 → 18.03 |
| `web/src/pages` | 1490.00 → 1595.00 | +105.00 | +105.00 | +0.00 | +0.00 | +0.00 | 10.93 → 10.90 |
| `internal/engine/workflow` | 500.00 → 597.00 | +97.00 | +85.00 | +12.00 | +0.00 | +0.00 | 27.47 → 28.78 |
| `web/src/components` | 2094.25 → 2178.25 | +84.00 | +84.00 | +0.00 | +0.00 | +0.00 | 18.71 → 18.62 |
| `web/src/lib` | 1157.35 → 1217.63 | +60.28 | +55.00 | +4.00 | +1.28 | +0.00 | 26.38 → 26.34 |
| `web/src/components/workflow-graph` | 42.00 → 96.00 | +54.00 | +54.00 | +0.00 | +0.00 | +0.00 | 19.35 → 13.40 |
| `cmd/legato` | 844.03 → 888.13 | +44.10 | +36.00 | +0.00 | +8.10 | +0.00 | 34.39 → 34.35 |
| `web/src/components/workflows` | 0.00 → 43.00 | +43.00 | +43.00 | +0.00 | +0.00 | +0.00 | n/a → 21.71 |
| `web/src/components/workflow-editor` | 465.00 → 504.00 | +39.00 | +39.00 | +0.00 | +0.00 | +0.00 | 13.83 → 13.30 |
| `web/src/components/workflow-canvas` | 0.00 → 33.00 | +33.00 | +31.00 | +2.00 | +0.00 | +0.00 | n/a → 15.98 |
| `config` | 224.80 → 256.10 | +31.30 | +28.00 | +0.00 | +3.30 | +0.00 | 34.69 → 34.05 |
| `web/src/components/settings` | 347.00 → 372.00 | +25.00 | +25.00 | +0.00 | +0.00 | +0.00 | 18.45 → 17.60 |
| `packages/even-hub/dev` | 42.00 → 53.00 | +11.00 | +9.00 | +2.00 | +0.00 | +0.00 | 12.86 → 11.04 |
| `internal/tui/overlay` | 361.20 → 371.20 | +10.00 | +10.00 | +0.00 | +0.00 | +0.00 | 28.18 → 27.96 |
| `internal/handofftest` | 41.25 → 51.05 | +9.80 | +7.00 | +0.00 | +2.80 | +0.00 | 54.74 → 53.20 |
| `internal/engine/store` | 824.90 → 833.90 | +9.00 | +10.00 | +0.00 | -1.00 | +0.00 | 37.29 → 37.49 |
| `internal/server` | 514.05 → 522.80 | +8.75 | +7.00 | +0.00 | +1.75 | +0.00 | 48.70 → 48.76 |
| `web/src/components/board` | 380.00 → 387.00 | +7.00 | +7.00 | +0.00 | +0.00 | +0.00 | 21.77 → 21.62 |
| `internal/tui` | 441.40 → 448.20 | +6.80 | +4.00 | +0.00 | +2.80 | +0.00 | 20.98 → 21.17 |
| `internal/engine/burnlite` | 0.00 → 4.80 | +4.80 | +0.00 | +0.00 | +4.80 | +0.00 | 100.00 → 77.70 |
| `internal/engine/jira` | 70.00 → 74.70 | +4.70 | +1.00 | +0.00 | +3.70 | +0.00 | 32.71 → 32.43 |
| `internal/tui/board` | 69.00 → 73.00 | +4.00 | +4.00 | +0.00 | +0.00 | +0.00 | 35.51 → 34.46 |
| `browser-extension/src` | 163.00 → 167.00 | +4.00 | +2.00 | +2.00 | +0.00 | +0.00 | 33.79 → 33.73 |
| `internal/tui/agents` | 146.15 → 148.15 | +2.00 | +2.00 | +0.00 | +0.00 | +0.00 | 32.94 → 33.16 |

## First-parent steps

| # | Commit | Subject (untrusted) | Changed lines | Δ debt | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Δ score | Gates failed |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | `bc203a430415` | Merge pull request #389 from cpave3/legato/leg-649 | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 2 | `4d6d56c2cea0` | Merge pull request #395 from cpave3/legato/leg-559-588-pr-ready | 4633 | +220.80 | +202.00 | +16.00 | +2.80 | +0.00 | +0.02 | no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 3 | `fee8984ccdeb` | Merge pull request #401 from cpave3/legato/LEG-650-narrow-prs-output | 225 | +12.00 | +12.00 | +0.00 | +0.00 | +0.00 | +0.00 | no_hotspot_excess_growth |
| 4 | `3a223ed5243a` | Merge pull request #403 from cpave3/legato/leg-651 | 574 | +70.00 | +68.00 | +2.00 | +0.00 | +0.00 | -0.05 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 5 | `d380cd0d7444` | Merge pull request #402 from cpave3/legato/LEG-653-nuget-packages-pin | 17 | +4.00 | +4.00 | +0.00 | +0.00 | +0.00 | +0.00 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 6 | `f740befd0cf2` | Merge pull request #404 from cpave3/legato/leg-652-task-id-as-text | 44 | +1.00 | +1.00 | +0.00 | +0.00 | +0.00 | +0.00 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 7 | `e2c149366ed5` | Merge pull request #406 from cpave3/legato/LEG-656-ticket-update | 743 | +48.50 | +40.00 | +0.00 | +8.50 | +0.00 | -0.01 | no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 8 | `ae8500ddb577` | Merge pull request #407 from cpave3/legato/LEG-655-definition-canvas | 1044 | +55.00 | +55.00 | +0.00 | +0.00 | +0.00 | -0.03 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 9 | `2ee8cd6adbd5` | Merge pull request #405 from cpave3/legato/LEG-654-store-close-waits | 32 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 10 | `3b3d3eea4937` | Merge pull request #408 from cpave3/legato/leg-657 | 1331 | +94.05 | +92.00 | +0.00 | +2.05 | +0.00 | -0.03 | no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 11 | `f4fca83ceb64` | Merge pull request #417 from cpave3/legato/LEG-662-related-work-runs | 52 | -0.30 | +2.00 | +0.00 | -2.30 | +0.00 | +0.01 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 12 | `e756ef170ec1` | Merge pull request #419 from cpave3/legato/LEG-659-subscription-usage | 2553 | +87.50 | +60.00 | +2.00 | +25.50 | +0.00 | +0.07 | no_new_exact_duplicates, no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 13 | `a33dd8ca3906` | Merge pull request #418 from cpave3/legato/leg-559-599-deprecate-inli… | 5878 | +257.20 | +248.00 | +2.00 | +7.20 | +0.00 | +0.03 | no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 14 | `60b19e5a32a5` | Merge pull request #425 from cpave3/legato/LEG-660-usage-profiles | 1829 | +142.95 | +144.00 | +2.00 | -3.05 | +0.00 | -0.06 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 15 | `bac2e695ffd4` | Merge pull request #424 from cpave3/legato/fix-config-show-sources-go… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 16 | `21336cd88c35` | Merge pull request #426 from cpave3/legato/LEG-660-startup-switch | 19 | +2.00 | +2.00 | +0.00 | +0.00 | +0.00 | +0.00 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 17 | `9678029a5506` | Merge pull request #420 from cpave3/legato/LEG-663-live-glasses-pages | 152 | +4.00 | +4.00 | +0.00 | +0.00 | +0.00 | +0.00 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 18 | `86e6d89843d9` | Merge pull request #423 from cpave3/legato/leg-661 | 861 | +142.55 | +138.00 | +2.00 | +2.55 | +0.00 | -0.10 | no_new_exact_duplicates, no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 19 | `4df5dadeb911` | Merge pull request #427 from cpave3/legato/6aeum288-delivery-docs | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 20 | `7d7de93024f9` | Merge pull request #428 from cpave3/legato/qykhugdr-watch-batches | 17 | +4.00 | +4.00 | +0.00 | +0.00 | +0.00 | +0.00 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 21 | `bf94ed1de8f4` | Merge pull request #429 from cpave3/legato/cajixhp6-opencode-usage | 766 | +26.75 | +22.00 | +0.00 | +4.75 | +0.00 | +0.00 | no_new_exact_duplicates, no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 22 | `5f7a847077b0` | Merge pull request #430 from cpave3/legato/8z3913fg-claude-accounts | 1860 | +253.55 | +239.00 | +4.00 | +10.55 | +0.00 | -0.17 | no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 23 | `471bfd24a696` | Merge pull request #431 from cpave3/legato/leg-666 | 576 | +2.00 | +2.00 | +0.00 | +0.00 | +0.00 | +0.02 | no_hotspot_excess_growth |

**Telescoping check:** Σ step Δdebt +1427.55 vs net +1427.55 → true (chain starts at base: true)

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
