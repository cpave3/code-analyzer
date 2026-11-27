# Final combined stack: Legato cleanup #436 → #444 (`471bfd24 → bb13a424`)

This is the **observed combined result** of the 11-layer cleanup stack,
measured with the frozen evaluator against the retained baseline
`471bfd24`. It is **not** a sum of the earlier independent batch deltas.
The stack layers were rebased (new SHAs) and some carry integration
changes, and repository-wide measurements interact across layers. All
earlier reports are unchanged.

**#434 is not in this stack.** Its two removed clone groups (server
invalidation events, jira `doPost`/`doPut`) are **not** included in any
number here.

## Provenance

| | |
|---|---|
| Baseline (retained) | `471bfd24a69628d985cdd7767dcb340da6a569b7`, tree `ae91f4192b23…`. This is the pinned benchmark main. No assumption about a moving HEAD or the current remote main |
| Stack tip | `bb13a424dc73f3eeb98f8bce43db932a9d19f7e6`, tree `25b725950d6a…` |
| Manifest | `.tmp/i95qqzwb-stack-bb13a424.json` (sha256 `05755c9a8365…`), copied verbatim to [final-stack-manifest-bb13a424.json](final-stack-manifest-bb13a424.json) |
| Layer verification | For each of the 11 layers: the ref in the bundle equals the manifest head; `merge-base --is-ancestor parent head`; the first-parent chain from head reaches **exactly** the manifest parent (actual Git `^1`, not GitHub base metadata); 0 merge commits. All passed. The stack is linear: 40 commits, 83 files, +4497/−774. #437's GitHub base reportedly still names the ledger, but its actual Git parent is the #435 head `679e0ed2`, and that is what was measured |
| Evaluator | Built from `git archive 5e0c679`, tool 0.1.0-poc, `scoring/v1`. Config sha256 `768b2fc6…85f6`, configHash `3cef63633359769f`. Unchanged, nothing tuned |
| Object source | Offline bundle `.tmp/i95qqzwb-stack-bb13a424.bundle`. `bundle verify`: complete history, 11 heads. Cloned with `--no-checkout --no-local` into `.tmp/legato-stack`, which has no worktree files and is not committed; it was unchanged after the runs. The historical base `9a61e9f8` is present and is an ancestor of `471bfd24`. No network. The shared Legato checkout was untouched: HEAD `471bfd24`, fingerprint `3ab8895f206d466e` before and after |
| Run | 2026-11-27T20:37:02+11:00. Combined evaluate took 0.93 s, combined history (40 steps) 12.7 s, each layer about 0.6 s, and drift 1.0 s |

| Report | Files |
|---|---|
| **Primary**: combined `471bfd24 → bb13a424` | `final-stack-471bfd24..bb13a424.{json,md}` |
| Combined first-parent history (40 commits) | `final-stack-history-471bfd24..bb13a424.{json,md}` |
| Per layer, manifest parent → head (11) | `final-stack-layers/layerNN-prPPP-<parent8>..<head8>.{json,md}` |
| Secondary, residual drift: `9a61e9f8 → bb13a424` | `final-drift-9a61e9f8..bb13a424.{json,md}` |

## Primary result: combined `471bfd24 → bb13a424`

Coverage: 83 files. 75 are Go: 37 source (+1014/−684) and 38 test
(+2988/−79). The other 8 are docs (+495/−11): ledger, ADRs 0001–0003,
AGENTS.md, docs/ai and the site page. Every analyzer status is `ok`, with
no unsupported changed files.

| Analyzer | Score | Raw (combined, observed) |
|---|---:|---|
| duplication | **+2** | **exact 10→8 (−2)**, **structural 260→253 (−7)**. 8 redundant instances removed, 1 introduced; 0 near-structural candidates |
| complexity | **0 (mixed)** | 55 modified callables, all decreased: cc 796→585 (−211), lines 2443→1920 (−523). 43 added (cc 245), 5 deleted (cc 20). **Changed-scope cc 816→830 (+14)**. Hotspots 29→26 (**−8 pre-existing cleared, +5 new**) |
| dependencies | **+1** | units 86→87; internal edges 200→199 (+5/−6); cycles 0→0; frozen-rule violations **15→11** (4 resolved, 0 new) |
| obsolete | 0 | nothing newly orphaned |
| scope | 0 | symbols +50/−6/~63 (renamed 7, formatting-only 3); exported API +6/−3/~12; 1 unit group |
| correctness | n/a | not assessed. 38 test files changed (+2988/−79), which is context only |

### Complexity: decision points versus function base points

| Combined, changed scope | Before | After | Δ |
|---|---:|---:|---:|
| **Total cyclomatic complexity** | **816** | **830** | **+14** |
| Function base points (1 per callable; 60 → 98 callables) | 60 | 98 | **+38** (43 added − 5 deleted) |
| **Decision points** (`if`/`for`/`case`/`&&`/`\|\|`) | **756** | **732** | **−24** |
| Roots (modified), cc | 796 | 585 | −211 |
| Roots plus helpers, cc | 816 | 830 | +14 |
| Max cc in scope | 107 | 96 | −11 (`BuildSandboxPlan`) |
| Max nesting, modified | 4 | 4 | 0 |
| Hotspots in scope | 29 | 26 | −3 |

Read this carefully. Root complexity fell by 211, but **total complexity in
the changed scope rose by 14**, because 43 new functions each add a base
point. Net **decision points fell by 24**. That is the only part of the
combined change that is a reduction in branching rather than a
redistribution of it. Most of the change is extraction into named helpers.

**Hotspots cleared (8 pre-existing):**
- 7 from #443: `agentDone` 24→7, `waitingToOpenOn` 17→6, `watchedPull.watch`
  13→9 (nesting 4→2), `SettleSteps` 12→8, `settleReady` 12→8,
  `settleObservers` 11→7 and `openSteps` 10→6. The other 6 #443 roots
  (`OpenReviewLoopClones` 20, `watchLoop` 19, `markStackReady` 15,
  `closeOutputs` 15, `addStacked` 14 and `AuthorDone` 11) shrank but are
  still hotspots.
- 1 from #441: `upgrader.mapText` 13→8 (nesting 4→3).

**New hotspots (5, all extracted helpers):** `checkNamedAccountIsolation`
cc 22 (#445), `readWatchRound` 15 (#443), `reviewRoundDone` 13 (#443),
`resumeConversation` 12 (#446) and `withClaudeAccounts` 10 (#444, exactly
at the threshold).

**Still hotspots:** `BuildSandboxPlan` (96), `sandboxPrefix` (58),
`resolveClaudeAccount` (28), `prepareTaskResume` (28), `checkReviewLoopStep`
(22), `OpenReviewLoopClones` (20), `watchLoop` (19), `TicketCreateExecutor.Start`
(18), `Recover` (16, nesting 4) and others, 26 in scope in total.

### Duplication

| Group | Base | Tip |
|---|---|---|
| `retryAfter` (exact, cross-engine) | 2 | 1, now `leaf/retryafter.Parse` |
| `openRegular` (exact, cross-engine) | 2 | 0; `leaf/regularfile.Open` |
| `claim*` runtime trio | 3 | 0, now `claimDataDirLock` |
| `FakeReviewGitHub` Draft/Push and SetAuthor/Describe (test fake, classified as source) | 2 + 2 | 0 + 0, now `editPR` |
| `PullRequestMode` / `ReadyAfterReview` | 2 | 0, now `rules` |
| `claudeUsageError` / `opencodeUsageError` | 2 (22 lines each) | **still a structural pair** (7-line wrappers over the reason tables): reduced, not eliminated |

The groups kept as won't-fix remain counted in the 253: B3
`CredentialsPath`/`SessionsDir`, C2, C4 and D1–D7. So do the #434 groups,
which are not in this stack. Reuse inside a function is invisible to the
duplication analyzer, because loops and closures are not indexed symbols.
That covers #439's `realPaths` loops (and #434's `plan.go` closure, which
is not in this stack).

### Dependencies and intentional surface

- **Edges added (5, intentional):** `engine/claudeusage → leaf/retryafter`,
  `engine/claudeusage → leaf/regularfile`, `engine/codexusage → leaf/regularfile`,
  `engine/opencodeusage → leaf/retryafter` (ADR 0001) and
  `tui/overlay → config` (ADR 0002).
- **Edges removed (6):** `config`, `server`, `service`, `tui` and
  `tui/overlay` → `engine/macros` (the package was deleted, ADR 0002), and
  `tui → engine/ipc` (ADR 0003).
- **Frozen violations resolved (4):** `config → engine/macros`,
  `tui → engine/macros`, `tui/overlay → engine/macros` and
  `tui → engine/ipc`.
- **Remaining 11 (all pre-existing `tui → engine`):** `events`, `store` ×5
  units, `workflow` ×2, `github`, `analytics` and `git`. Under ADR 0003
  these may be **permitted value-type imports**. The frozen strict rule
  still counts them, so **11 is not a count of ADR 0003 violations**.
- **Exported API +6 / −3 / ~12:**
  - New: `regularfile.Open`, `service.EventStream`, `service.PlanVerdict`,
    `service.SendPlanVerdict`, plus `config.Macro.Validate` and
    `config.SplitFavoriteMacros`, which are moves.
  - Removed: `macros.Macro.Validate` and `macros.SplitFavorites` (moved),
    and `macros.ListResult` (now server-private: a real narrowing).
  - Renamed or moved: `retryAfter → retryafter.Parse`, plus `Macro` and its
    `UnmarshalYAML`/`UnmarshalJSON`, and `AnyFavorite → AnyFavoriteMacro`.
  - Signature changes: 7 are only the type path `macros.Macro →
    config.Macro`. `NewApp` changed both its bus parameter (`*events.Bus →
    service.EventStream`) and its macro type.

## Rule and version mismatch (frozen config versus new ADRs)

| New rule | Frozen evaluator can attest? | Evidence at tip |
|---|---|---|
| ADR 0001: leaf packages import only the standard library | **No.** There is no stdlib-only rule type | Measured: both leaves have internal fan-out 0. Manual: 0 Legato imports in non-test `internal/leaf/**`. Enforced by Legato's `TestLeafPackagesImportOnlyTheStandardLibrary` (not run here) |
| ADR 0001: engine imports only leaf packages among Legato's non-engine packages | **No.** There are no allow-list rules | Manual: the only non-engine Legato imports from non-test engine code are the 2 leaves (2 imports each) |
| Engine-sibling prohibition | Yes (frozen rule) | 0 violations at base and tip |
| ADR 0002: `config` imports no Legato package | **Partly.** The frozen rule covers `config → internal/**` | Measured: `config` fan-out 1→0 across the whole graph (non-test files). Manual: 0 Legato imports in non-test `config/*.go` |
| ADR 0002: engine must not import `config` | **No.** It is not in the frozen config | Manual: 0 non-test engine files import `config`. Enforced by Legato's architecture test (not run here) |
| ADR 0003: TUI names only audited engine values and pure helpers; behaviour goes through service | **No.** The frozen rule is import-level and strict | Frozen counts 11 remaining `tui → engine` imports regardless of intent. Symbol-level allowlist enforced by `TestTUIUsesOnlyEngineValueTypesAndPureHelpers` (not run here). The integration change in #442 drops the obsolete `macros` row |

Expressing these rules would need new rule types and a new config version.
Neither was applied, so measurements stay comparable.

## Per-layer comparisons (manifest parent → head; actual trees)

Each row is its own report. **Rows are not independent effects.** Each
layer is measured on top of the layers below it in this stack. For this
stack the changed scopes are disjoint, so the changed-scope cc and
structural deltas happen to sum to the combined figures (Σ cc +14, Σ
structural −7). That is an observation about this stack, not a property to
rely on.

| # | Layer | dup | cx | dep | Struct | Exact | Root cc Δ | Added n (cc) | Deleted n (cc) | Changed-scope cc | Hotspots (new) | Edges +/− | Violations | Files src/test/docs |
|---:|---|---:|---:|---:|---|---|---:|---|---|---|---|---|---|---|
| 1 | #436 ledger `471bfd24..95295362` | 0 | 0 | 0 | 260→260 | 10→10 | 0 | 0 | 0 | 0→0 | 0→0 | +0/−0 | 15→15 | 0/0/1 |
| 2 | #435 leaf `95295362..679e0ed2` | +1 | 0 | 0 | 260→258 | 10→8 | 0 | 1 (4) | 3 (14) | 61→51 | 3→3 | +4/−0 | 15→15 | 8/4/4 |
| 3 | #437 service `679e0ed2..907025f4` | +2 | +2 | 0 | 258→253 | 8→8 | −29 | 4 (14) | 0 | 46→31 | 0→0 | +0/−0 | 15→15 | 7/4/0 |
| 4 | #440 macros `907025f4..ff6ba3d0` | 0 | 0 | +1 | 253→253 | 8→8 | 0 | 2 (6) | 2 (6) | 94→94 | 5→5 | +1/−5 | 15→12 | 7/11/5 |
| 5 | #442 TUI `ff6ba3d0..3d112f74` | 0 | 0 | +1 | 253→253 | 8→8 | 0 | 1 (1) | 0 | 12→13 | 0→0 | +0/−1 | 12→11 | 2/3/3 |
| 6 | #441 small `3d112f74..f506d37e` | 0 | +1 | 0 | 253→253 | 8→8 | −23 | 5 (32) | 0 | 93→102 | 4→3 | +0/−0 | 11→11 | 5/4/0 |
| 7 | #443 review `f506d37e..41f83616` | 0 | 0 | 0 | 253→253 | 8→8 | −92 | 16 (101) | 0 | 237→246 | 13→8 (+2) | +0/−0 | 11→11 | 6/8/0 |
| 8 | #446 resume `41f83616..6ac9257d` | 0 | 0 | 0 | 253→253 | 8→8 | −20 | 4 (26) | 0 | 48→54 | 1→2 (+1) | +0/−0 | 11→11 | 1/1/0 |
| 9 | #445 accounts `6ac9257d..a3ee21f2` | 0 | 0 | 0 | 253→253 | 8→8 | −27 | 3 (33) | 0 | 55→61 | 1→2 (+1) | +0/−0 | 11→11 | 1/1/0 |
| 10 | #439 sandbox plan `a3ee21f2..4cf88e48` | 0 | +1 | 0 | 253→253 | 8→8 | −11 | 6 (18) | 0 | 107→114 | 1→1 | +0/−0 | 11→11 | 1/1/0 |
| 11 | #444 sandbox prefix `4cf88e48..bb13a424` | 0 | 0 | 0 | 253→253 | 8→8 | −9 | 1 (10) | 0 | 67→68 | 1→2 (+1) | +0/−0 | 11→11 | 1/1/0 |
| **Σ** | **combined (observed)** | **+2** | **0** | **+1** | **260→253** | **10→8** | **−211** | **43 (245)** | **5 (20)** | **816→830** | **29→26 (+5)** | **+5/−6** | **15→11** | **37/38/8** |

**Differences from the earlier independent measurements** (the actual
layer trees are ground truth):
- **#441's layer is larger** than the independently measured `b9696a9c`.
  It has 6 commits instead of 4. It adds `821efd67` (test) and `f506d37e`,
  which extracts `validReviewLoopPRs` (cc 9, nesting 3) from
  `checkReviewLoopStep` (cc 29→22, nesting 4→3; still a hotspot). So this
  layer's root change is −23 (not −16) and its helpers total cc 32 (not 23).
- **#437's layer** adds a test-only commit (`907025f4`, usage-error
  precedence, +53 test lines). Code metrics are the same as the
  independent run.
- **#435's layer** adds a docs commit (`679e0ed2`, site architecture page).
  Code metrics are the same.
- **#442's layer** sits on #440, so it uses `config.Macro`, and adds
  `3d112f74`, which drops the `macros` row from the TUI allowlist test. The
  violations go 12→11 here because #440 had already resolved the two
  `→ engine/macros` violations.
- **#436's ledger** has a third docs commit (`95295362`).
- **#440**: AGENTS.md combines the leaf and config documentation (an
  integration change the author reported). Code metrics are the same as
  the independent run.

## Secondary: residual drift `9a61e9f8 → bb13a424` (separate from the cleanup)

This compares the **original pre-window baseline** with the stack tip.
That covers the whole two-day feature window **plus** the cleanup stack
(without #434). It is **not** the cleanup delta and must not be read as one.

| | Original window `9a61e9f8 → 471bfd24` ([EVIDENCE.md](EVIDENCE.md)) | Window plus stack `9a61e9f8 → bb13a424` |
|---|---|---|
| Scores (dup / cx / dep / obs / scope) | −2 / −2 / 0 / 0 / 0 | −2 / −2 / **+1** / 0 / 0 |
| Coverage | partial (Go only) | partial (Go only) |
| Structural redundancy | 243 → 260 (+17) | 243 → 253 (+10) |
| Exact redundancy | 8 → 10 | **8 → 8** |
| Modified pre-existing callables, net cc | +194 | +94 |
| Added callables | 569 | 610 (cc 3100) |
| New hotspot callables | 84 | 83 |
| Internal edges | 193 → 200 | 193 → 199 |
| Frozen-rule violations | 15 → 15 | 15 → 11 |

The scores still saturate at −2 for duplication and complexity, because the
window's feature volume dominates (+64.6k/−4.3k lines, 610 new callables).
On these metrics, the stack recovers roughly 40% of the structural
redundancy the window added (17 → 10 net). It brings exact redundancy back
to its pre-window level, roughly halves the window's net growth in
pre-existing callables (+194 → +94), removes 4 frozen-rule violations and
adds no cycles. The 84 → 83 new hotspots barely move, because the cleanup
replaced several hotspot roots with hotspot helpers. TypeScript and other
non-Go changes remain unmeasured.

## Conclusion: useful versus limited

**Useful (measured, reproducible):**
- **Cross-package copy-paste is down.** Exact 10→8 and structural 260→253,
  with each removed group pinned to before and after locations. The kept
  groups and the one remaining C1 pair are named.
- **Coupling moved in the intended direction.** −6 edges for +5 intentional
  leaf and config edges, `engine/macros` deleted, `config` fan-out 0, 4
  frozen-rule violations resolved, no cycles.
- **Large functions got smaller and flatter.** 55 roots, all decreased
  (−211 cc, −523 lines). The largest function in scope went 107→96. 8
  pre-existing hotspots cleared.
- **The decision-point budget fell slightly** (−24) once base points are
  separated out.

**Limited (do not over-read):**
- **Complexity was mostly redistributed.** Changed-scope cc is +14 overall,
  and 5 new hotspot helpers appeared. The frozen score is 0 (mixed), and the
  positive duplication and dependency scores do not mean the code is
  simpler overall.
- **The new architecture rules are only partly attestable.** Stdlib-only
  leaves, engine-only-leaf, engine ↛ config and ADR 0003's symbol-level
  boundary rest on Legato's architecture tests and manual checks. 11
  strict-rule "violations" may be permitted value imports.
- **The API surface grew** (+6 exported, 3 of them a service forwarding
  layer), with one real narrowing.
- **No correctness or behaviour claim.** 38 test files were added or changed
  and author or coordinator test and mutation evidence was reported. None
  was executed here. Review and CI are pending, and main has known baseline
  failures.
- **Go only.** Reuse inside functions is invisible, references are matched
  by name, and the absolute thresholds saturate on large windows.
- **#434 is excluded.** Its 2 clone-group removals and changed-scope −6 are
  not in these numbers.
