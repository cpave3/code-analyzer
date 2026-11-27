# CodeEval v2 delta: legato-stack `9a61e9f8a021` → `471bfd24a696`

- **Base:** `9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc` tree `a55db6103bf4` 2026-11-25T21:49:11+11:00
- **Target:** `471bfd24a69628d985cdd7767dcb340da6a569b7` tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.1.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 988 of 1858 production source files measured (61.8% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2, typescript 577

**Change size:** 296 production source file(s), 201 measured (16581 changed lines); not measured: javascript 2, kotlin 2, sql 7, typescript 84

## Ledger (after − before; telescopes across commits)

| | Before | After | Δ | Δ per changed kLOC |
|---|---:|---:|---:|---:|
| **Total debt** | 10783.30 | 11611.30 | **+828.00** | +49.94 |
| complexity debt | 9275.00 | 10027.00 | +752.00 | +45.35 |
| nesting debt | 598.00 | 618.00 | +20.00 | +1.21 |
| duplication debt | 760.30 | 816.30 | +56.00 | +3.38 |
| architecture debt | 150.00 | 150.00 | +0.00 | +0.00 |
| score (0–100) | 34.19 | 33.93 | -0.26 | |
| debt / kLOC | 48.12 | 48.68 | +0.56 | |

| Raw component | Before | After | Δ |
|---|---:|---:|---:|
| source lines | 224082 | 238533 | +14451 |
| callables | 8734 | 9298 | +564 |
| decision points | 35801 | 38387 | +2586 |
| cyclomatic total | 44535 | 47685 | +3150 |
| hotspots | 1197 | 1288 | +91 |
| cyclomatic excess | 9275 | 10027 | +752 |
| nesting excess | 299 | 309 | +10 |
| exact redundant callables | 8 | 10 | +2 |
| structural redundant callables | 243 | 260 | +17 |
| rule-violating edges | 15 | 15 | +0 |
| units in cycles | 0 | 0 | +0 |
| internal edges | 193 | 200 | +7 |
| redundant tokens | 15206.00 | 16326.00 | +1120.00 |

## Changed callables: decision points versus function base points

- modified pre-existing: 242, decision points 3257 → 3451 (+194), cyclomatic 3499 → 3693
- added: 570 (decision points 2407); deleted: 6 (decision points 15)
- **changed scope: decision points +2586, function base points +564** (cyclomatic = decisions + 1 per function)

## Advisory gates

| Gate | Result | Detail |
|---|---|---|
| no_new_exact_duplicates | **fail** | +2 exact redundant callables |
| no_new_structural_duplicates | **fail** | +17 structural redundant callables |
| no_new_rule_violations | pass | 0 rule-violating edges |
| no_new_cycles | pass | 0 units in cycles |
| no_branch_growth_in_existing_callables | **fail** | +194 decision points in modified pre-existing callables |
| no_hotspot_excess_growth | **fail** | +752 cyclomatic excess, +10 nesting excess |

## Module changes (16 changed; worst first)

| Module | Debt before → after | Δ total | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Score before → after |
|---|---|---:|---:|---:|---:|---:|---|
| `internal/service` | 4523.95 → 5070.35 | +546.40 | +522.00 | +6.00 | +18.40 | +0.00 | 32.69 → 32.14 |
| `internal/engine/workflow` | 500.00 → 597.00 | +97.00 | +85.00 | +12.00 | +0.00 | +0.00 | 27.47 → 28.78 |
| `cmd/legato` | 844.03 → 888.13 | +44.10 | +36.00 | +0.00 | +8.10 | +0.00 | 34.39 → 34.35 |
| `config` | 224.80 → 256.10 | +31.30 | +28.00 | +0.00 | +3.30 | +0.00 | 34.69 → 34.05 |
| `internal/engine/codexusage` | 0.00 → 18.30 | +18.30 | +13.00 | +2.00 | +3.30 | +0.00 | n/a → 42.95 |
| `internal/engine/claudeusage` | 0.00 → 16.68 | +16.68 | +11.00 | +0.00 | +5.68 | +0.00 | n/a → 40.45 |
| `internal/engine/opencodeusage` | 0.00 → 14.38 | +14.38 | +12.00 | +0.00 | +2.38 | +0.00 | n/a → 31.71 |
| `internal/tui/overlay` | 361.20 → 371.20 | +10.00 | +10.00 | +0.00 | +0.00 | +0.00 | 28.18 → 27.96 |
| `internal/handofftest` | 41.25 → 51.05 | +9.80 | +7.00 | +0.00 | +2.80 | +0.00 | 54.74 → 53.20 |
| `internal/engine/store` | 824.90 → 833.90 | +9.00 | +10.00 | +0.00 | -1.00 | +0.00 | 37.29 → 37.49 |
| `internal/server` | 514.05 → 522.80 | +8.75 | +7.00 | +0.00 | +1.75 | +0.00 | 48.70 → 48.76 |
| `internal/tui` | 441.40 → 448.20 | +6.80 | +4.00 | +0.00 | +2.80 | +0.00 | 20.98 → 21.17 |
| `internal/engine/burnlite` | 0.00 → 4.80 | +4.80 | +0.00 | +0.00 | +4.80 | +0.00 | 100.00 → 77.70 |
| `internal/engine/jira` | 70.00 → 74.70 | +4.70 | +1.00 | +0.00 | +3.70 | +0.00 | 32.71 → 32.43 |
| `internal/tui/board` | 69.00 → 73.00 | +4.00 | +4.00 | +0.00 | +0.00 | +0.00 | 35.51 → 34.46 |
| `internal/tui/agents` | 146.15 → 148.15 | +2.00 | +2.00 | +0.00 | +0.00 | +0.00 | 32.94 → 33.16 |

## First-parent steps

| # | Commit | Subject (untrusted) | Changed lines | Δ debt | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Δ score | Gates failed |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | `bc203a430415` | Merge pull request #389 from cpave3/legato/leg-649 | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 2 | `4d6d56c2cea0` | Merge pull request #395 from cpave3/legato/leg-559-588-pr-ready | 4229 | +157.80 | +141.00 | +14.00 | +2.80 | +0.00 | +0.04 | no_new_structural_duplicates, no_hotspot_excess_growth |
| 3 | `fee8984ccdeb` | Merge pull request #401 from cpave3/legato/LEG-650-narrow-prs-output | 225 | +12.00 | +12.00 | +0.00 | +0.00 | +0.00 | -0.01 | no_hotspot_excess_growth |
| 4 | `3a223ed5243a` | Merge pull request #403 from cpave3/legato/leg-651 | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 5 | `d380cd0d7444` | Merge pull request #402 from cpave3/legato/LEG-653-nuget-packages-pin | 17 | +4.00 | +4.00 | +0.00 | +0.00 | +0.00 | +0.00 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 6 | `f740befd0cf2` | Merge pull request #404 from cpave3/legato/leg-652-task-id-as-text | 42 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 7 | `e2c149366ed5` | Merge pull request #406 from cpave3/legato/LEG-656-ticket-update | 739 | +48.50 | +40.00 | +0.00 | +8.50 | +0.00 | -0.03 | no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 8 | `ae8500ddb577` | Merge pull request #407 from cpave3/legato/LEG-655-definition-canvas | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 9 | `2ee8cd6adbd5` | Merge pull request #405 from cpave3/legato/LEG-654-store-close-waits | 32 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 10 | `3b3d3eea4937` | Merge pull request #408 from cpave3/legato/leg-657 | 432 | +17.05 | +15.00 | +0.00 | +2.05 | +0.00 | +0.00 | no_new_structural_duplicates, no_hotspot_excess_growth |
| 11 | `f4fca83ceb64` | Merge pull request #417 from cpave3/legato/LEG-662-related-work-runs | 19 | -2.30 | +0.00 | +0.00 | -2.30 | +0.00 | +0.01 |  |
| 12 | `e756ef170ec1` | Merge pull request #419 from cpave3/legato/LEG-659-subscription-usage | 2239 | +75.50 | +48.00 | +2.00 | +25.50 | +0.00 | +0.07 | no_new_exact_duplicates, no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 13 | `a33dd8ca3906` | Merge pull request #418 from cpave3/legato/leg-559-599-deprecate-inli… | 5197 | +179.20 | +174.00 | -2.00 | +7.20 | +0.00 | +0.03 | no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 14 | `60b19e5a32a5` | Merge pull request #425 from cpave3/legato/LEG-660-usage-profiles | 1577 | +102.95 | +104.00 | +2.00 | -3.05 | +0.00 | -0.09 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 15 | `bac2e695ffd4` | Merge pull request #424 from cpave3/legato/fix-config-show-sources-go… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 16 | `21336cd88c35` | Merge pull request #426 from cpave3/legato/LEG-660-startup-switch | 19 | +2.00 | +2.00 | +0.00 | +0.00 | +0.00 | +0.00 | no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 17 | `9678029a5506` | Merge pull request #420 from cpave3/legato/LEG-663-live-glasses-pages | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 18 | `86e6d89843d9` | Merge pull request #423 from cpave3/legato/leg-661 | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 19 | `4df5dadeb911` | Merge pull request #427 from cpave3/legato/6aeum288-delivery-docs | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 20 | `7d7de93024f9` | Merge pull request #428 from cpave3/legato/qykhugdr-watch-batches | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 21 | `bf94ed1de8f4` | Merge pull request #429 from cpave3/legato/cajixhp6-opencode-usage | 611 | +18.75 | +14.00 | +0.00 | +4.75 | +0.00 | +0.01 | no_new_exact_duplicates, no_new_structural_duplicates, no_hotspot_excess_growth |
| 22 | `5f7a847077b0` | Merge pull request #430 from cpave3/legato/8z3913fg-claude-accounts | 1755 | +212.55 | +198.00 | +4.00 | +10.55 | +0.00 | -0.29 | no_new_structural_duplicates, no_branch_growth_in_existing_callables, no_hotspot_excess_growth |
| 23 | `471bfd24a696` | Merge pull request #431 from cpave3/legato/leg-666 | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |

**Telescoping check:** Σ step Δdebt +828.00 vs net +828.00 → true (chain starts at base: true)

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
