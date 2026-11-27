# CodeEval v2 delta: legato-stack `471bfd24a696` → `bb13a424dc73`

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `bb13a424dc73f3eeb98f8bce43db932a9d19f7e6` tree `25b725950d6a` 2026-11-27T20:29:04+11:00
- **Scoring:** `scoring/v2-experimental` (experimental) · config `3cef63633359769f` · tool 0.1.0-poc
- **Weights (uncalibrated):** cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC
- **Thresholds:** hotspot cyclomatic 10, nesting 4, clone min tokens 40

**Coverage:** 990 of 1860 production source files measured (61.9% of source bytes); not measured: javascript 23, kotlin 54, python 5, rust 34, shell 15, sql 160, swift 2, typescript 577

**Change size:** 37 production source file(s), 37 measured (1698 changed lines)

## Ledger (after − before; telescopes across commits)

| | Before | After | Δ | Δ per changed kLOC |
|---|---:|---:|---:|---:|
| **Total debt** | 11611.30 | 11387.20 | **-224.10** | -131.98 |
| complexity debt | 10027.00 | 9888.00 | -139.00 | -81.86 |
| nesting debt | 618.00 | 612.00 | -6.00 | -3.53 |
| duplication debt | 816.30 | 777.20 | -39.10 | -23.03 |
| architecture debt | 150.00 | 110.00 | -40.00 | -23.56 |
| score (0–100) | 33.93 | 34.40 | +0.47 | |
| debt / kLOC | 48.68 | 47.67 | -1.01 | |

| Raw component | Before | After | Δ |
|---|---:|---:|---:|
| source lines | 238533 | 238865 | +332 |
| callables | 9298 | 9336 | +38 |
| decision points | 38387 | 38363 | -24 |
| cyclomatic total | 47685 | 47699 | +14 |
| hotspots | 1288 | 1285 | -3 |
| cyclomatic excess | 10027 | 9888 | -139 |
| nesting excess | 309 | 306 | -3 |
| exact redundant callables | 10 | 8 | -2 |
| structural redundant callables | 260 | 253 | -7 |
| rule-violating edges | 15 | 11 | -4 |
| units in cycles | 0 | 0 | +0 |
| internal edges | 200 | 199 | -1 |
| redundant tokens | 16326.00 | 15544.00 | -782.00 |

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

## First-parent steps

| # | Commit | Subject (untrusted) | Changed lines | Δ debt | Δ complexity | Δ nesting | Δ duplication | Δ architecture | Δ score | Gates failed |
|---|---|---|---:|---:|---:|---:|---:|---:|---:|---|
| 1 | `7d9862b658af` | docs(analysis): CodeEval findings ledger for 9a61e9f8..471bfd24 (i95q… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 2 | `8661dc42a368` | docs(analysis): ledger records Cameron's G decisions, #435 measuremen… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 3 | `952953624e3d` | docs(analysis): ledger maps every item to its PR (#435–#446) and the … | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 4 | `fb7d1f8f08bb` | refactor(arch): stdlib-only leaf packages under internal/leaf (i95qqz… | 67 | +9.00 | +0.00 | +0.00 | +9.00 | +0.00 | -0.01 | no_new_exact_duplicates, no_new_structural_duplicates |
| 5 | `d40cc841b696` | refactor(usage): claude and opencode readers share retryafter.Parse (… | 38 | -9.50 | +0.00 | +0.00 | -9.50 | +0.00 | +0.02 |  |
| 6 | `6d99f8bc622f` | refactor(usage): claude and codex readers share regularfile.Open (i95… | 64 | -8.50 | +0.00 | +0.00 | -8.50 | +0.00 | +0.01 |  |
| 7 | `679e0ed2d35b` | docs(site): the architecture page states the leaf-package rule (i95qq… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 8 | `adc5c3aa7b79` | refactor(usage): provider usage errors map through a reason table (i9… | 84 | -5.35 | +0.00 | +0.00 | -5.35 | +0.00 | +0.01 |  |
| 9 | `306d76debd00` | refactor(service): PullRequestModes reads its rules in one place (i95… | 28 | -2.55 | +0.00 | +0.00 | -2.55 | +0.00 | +0.00 |  |
| 10 | `5eb0119b5982` | refactor(cmd): the singleton runtime roles share claimDataDirLock (i9… | 83 | -16.20 | +0.00 | +0.00 | -16.20 | +0.00 | +0.04 |  |
| 11 | `459afaa30d18` | refactor(handofftest): the review fake edits a pull request through e… | 38 | -6.00 | +0.00 | +0.00 | -6.00 | +0.00 | +0.01 |  |
| 12 | `907025f419bf` | test(usage): pin the precedence of overlapping provider usage errors … | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 13 | `ff6ba3d09493` | refactor(config): config owns the macro model; delete engine/macros (… | 118 | -30.00 | +0.00 | +0.00 | +0.00 | -30.00 | +0.05 |  |
| 14 | `db51f9c6e63e` | test(tui): pin how plan verdicts reach the conductor and the bus (i95… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 15 | `684dc2dad39b` | refactor(tui): engine behaviour reaches the TUI through service (ADR … | 64 | -10.00 | +0.00 | +0.00 | +0.00 | -10.00 | +0.03 |  |
| 16 | `3d112f744e19` | test(tui): drop the macros row from the engine allowlist (i95qqzwb) | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 17 | `705c80a44f64` | refactor(agents): a conversation's newest recorded account is one loo… | 25 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 18 | `992e301fe032` | refactor(history): recovery asks storageDirsFor where a run's convers… | 46 | -5.00 | -5.00 | +0.00 | +0.00 | +0.00 | +0.01 |  |
| 19 | `6b0d83c96343` | refactor(workflow): upgrade finds whole references in one function (i… | 35 | -6.00 | -4.00 | -2.00 | +0.00 | +0.00 | +0.01 |  |
| 20 | `a8e735372651` | refactor(workflow): ticket.create resolves its assignee in one place … | 40 | -3.00 | -3.00 | +0.00 | +0.00 | +0.00 | +0.01 |  |
| 21 | `821efd6774c0` | test(workflow): pin the shape a review loop step's prs must have (i95… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 22 | `f506d37e1b0f` | refactor(workflow): a review loop step's prs shape is validReviewLoop… | 39 | -9.00 | -7.00 | -2.00 | +0.00 | +0.00 | +0.02 |  |
| 23 | `3c3283390708` | refactor(service): map a watched check context in watchedCheck.toWatc… | 76 | -6.00 | -4.00 | -2.00 | +0.00 | +0.00 | +0.01 |  |
| 24 | `035f776204f5` | refactor(service): read one watch round in readWatchRound (i95qqzwb) | 89 | -8.00 | -8.00 | +0.00 | +0.00 | +0.00 | +0.02 |  |
| 25 | `8ff58506707d` | refactor(service): choose the layers a steward adds in stackedLayersT… | 58 | -8.00 | -8.00 | +0.00 | +0.00 | +0.00 | +0.01 |  |
| 26 | `2b9404c4d011` | refactor(service): settle waiting review steps through one fail and c… | 49 | -9.00 | -9.00 | +0.00 | +0.00 | +0.00 | +0.02 |  |
| 27 | `d32ce9ed8d14` | refactor(service): name the ready step's head check and the clean pas… | 116 | -14.00 | -14.00 | +0.00 | +0.00 | +0.00 | +0.03 |  |
| 28 | `24c9601c95c7` | refactor(service): name when an agent's review round is done (i95qqzw… | 34 | -11.00 | -11.00 | +0.00 | +0.00 | +0.00 | +0.02 |  |
| 29 | `f6758acfc9f8` | refactor(service): read an author's round from roundPRs, end with clo… | 47 | -6.00 | -6.00 | +0.00 | +0.00 | +0.00 | +0.02 |  |
| 30 | `766845ae5773` | refactor(service): name the two scans of steps waiting to open (i95qq… | 40 | -8.00 | -8.00 | +0.00 | +0.00 | +0.00 | +0.01 |  |
| 31 | `41f836167310` | refactor(service): list the step-owned loops' hosts in stepOwnedLoopH… | 35 | -2.00 | -2.00 | +0.00 | +0.00 | +0.00 | +0.01 |  |
| 32 | `af697f4305d9` | test(service): pin every path of prepareTaskResume (i95qqzwb) | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 33 | `6ac9257d1067` | refactor(service): split prepareTaskResume by responsibility (i95qqzw… | 143 | -17.00 | -17.00 | +0.00 | +0.00 | +0.00 | +0.04 |  |
| 34 | `67be03c20fd8` | test(service): characterise resolveClaudeAccount's every return (i95q… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 35 | `a3ee21f209d4` | refactor(service): split resolveClaudeAccount's named-account guards … | 67 | -14.00 | -14.00 | +0.00 | +0.00 | +0.00 | +0.03 |  |
| 36 | `52bfa43e2600` | test(service): pin sandbox plans around Claude accounts and NuGet (i9… | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 37 | `17c33dcc481c` | refactor(service): move the Claude-account and NuGet parts out of Bui… | 102 | -9.00 | -9.00 | +0.00 | +0.00 | +0.00 | +0.02 |  |
| 38 | `4cf88e484251` | refactor(service): resolve ToolState and GH overlays through realPath… | 12 | -2.00 | -2.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 39 | `1c70f46760e2` | test(sandbox): pin which Claude accounts a sandboxed agent can reach … | 0 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 | +0.00 |  |
| 40 | `bb13a424dc73` | refactor(sandbox): sandboxPrefix applies the Claude account rule in w… | 63 | -8.00 | -8.00 | +0.00 | +0.00 | +0.00 | +0.02 |  |

**Telescoping check:** Σ step Δdebt -224.10 vs net -224.10 → true (chain starts at base: true)

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
