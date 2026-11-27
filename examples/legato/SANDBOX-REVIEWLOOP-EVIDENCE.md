# Batch deltas: #439 updated head (6af114f8) and #443 review loop (f3ddb233)

These are two independent frozen evaluations, both from pinned main
`471bfd24`. **Do not add the deltas together.** #439 and #443 are separate
branches, and the repository-wide measurements interact. The earlier #439
measurement at `1b322d6c` ([BATCHES-439-441-EVIDENCE.md](BATCHES-439-441-EVIDENCE.md))
stays as it was, with its own pins. It is not relabelled.

## Provenance (shared)

| | |
|---|---|
| Base | `471bfd24a69628d985cdd7767dcb340da6a569b7`, tree `ae91f4192b23…`. This is the pinned benchmark main, not a claim about the current remote main |
| #439 updated | `6af114f8aeb04529fee83616d3bbd542c1329413`, tree `d6fa0d75ef0f…`, ref `legato/i95qqzwb-sandbox-plan`. The previously measured `1b322d6c` **is an ancestor**: the update is one commit on top (`ffd5e720 → 1b322d6c → 6af114f8`), and `1b322d6c..6af114f8` touches only `internal/service/sandbox_plan.go` (+2/−10) |
| #443 review loop | `f3ddb233148b7792387b3ba9e0f4abe35b02aa13`, tree `978f92a58bf0…`, ref `legato/i95qqzwb-review-loop`. 9 linear first-parent commits on top of base |
| Evaluator | Built from `git archive 5e0c679`, tool 0.1.0-poc, `scoring/v1`. Config sha256 `768b2fc6…85f6`, configHash `3cef63633359769f`. Unchanged, nothing tuned |
| Object source | Offline bundle `.tmp/i95qqzwb-sandbox6af-reviewf3dd.bundle`. `bundle verify`: complete, and the 2 heads equal the pins. `merge-base --is-ancestor` verified. Cloned with `--no-checkout --no-local` into `.tmp/legato-b5`, which has no worktree files and is not committed; it was unchanged after the runs. No network. The shared Legato checkout was untouched: HEAD `471bfd24`, fingerprint `3ab8895f206d466e` before and after |
| Run | 2026-11-27T20:07:44+11:00. Evaluate took 0.65–0.74 s and history 0.66–3.7 s |
| Per-callable detail | Helper cc and nesting come from the evaluator's own index cache (adapter `go-adapter/3`) for the target blobs. These are the same numbers the analyzers use. "Decisions" is the evaluator's `decisions` metric (cc = decisions + 1 per callable) |

| Comparison | Net | History |
|---|---|---|
| #439 full: `471bfd24 → 6af114f8` | `sandbox-471bfd24..6af114f8.{json,md}` | `sandbox-history-471bfd24..6af114f8.{json,md}` |
| #439 increment: `1b322d6c → 6af114f8` | `sandbox-1b322d6c..6af114f8.{json,md}` | `sandbox-history-1b322d6c..6af114f8.{json,md}` |
| #443: `471bfd24 → f3ddb233` | `reviewloop-471bfd24..f3ddb233.{json,md}` | `reviewloop-history-471bfd24..f3ddb233.{json,md}` |

## Scores and raw changed-scope totals

| | #439 full (base→6af1) | #439 increment (1b32→6af1) | #443 (base→f3dd) |
|---|---:|---:|---:|
| duplication | 0 | 0 | 0 |
| complexity | **+1** | 0 | **0 (mixed)** |
| dependencies / obsolete / scope | 0 / 0 / 0 | 0 / 0 / 0 | 0 / 0 / 0 |
| correctness | n/a | n/a | n/a |
| Root (modified) cc | 107 → 96 (−11) | 98 → 96 (−2) | 237 → 145 (−92), 13 callables, all down |
| Added helpers | 6, cc 18 | 0 | 16, cc 101 |
| Deleted | 0 | 0 | 0 |
| **Changed-scope cc** | **107 → 114 (+7)** | **98 → 96 (−2)** | **237 → 246 (+9)** |
| of which: base points (+1 per new callable) | +6 | 0 | +16 |
| of which: decision points (net) | **+1** | **−2** | **−7** |
| Max cc in scope | 107 → 96 | 98 → 96 | 33 → 20 |
| Max nesting, modified | 4 → 4 | 4 → 4 | 4 → 3 |
| Hotspots in scope | 1 → 1 | 1 → 1 | 13 → 8 (−7 pre-existing, **+2 new**) |

How to read the complexity row: neither +1 nor 0 means logic was removed.
Total cyclomatic complexity in the changed scope **rose** in #439 (+7) and in
#443 (+9). Almost all of that rise is the +1 baseline of each new function.
Net decision points barely move: +1 in #439 and −7 in #443.

---

## #439 at updated head `6af114f8`

**Increment (`1b322d6c → 6af114f8`).** One modified callable:
`BuildSandboxPlan` went from cc 98 to 96, decisions 97 to 95, and 337 to 329
lines. Nesting stayed at 4. Nothing was added or deleted. It scores 0
because a −2 change is below the ±3 threshold. The commit replaces two
in-function loops (for ToolState and GH.Overlays) with calls to the
existing `realPaths` helper (cc 2, unchanged). **The duplication analyzer
cannot see this reuse:** loops inside a function are not indexed symbols,
so redundancy stays at 260→260. It shows up only as −2 cc and −8 lines in
the root. This is the same blind spot as #434's `plan.go` closure.

**Full (`471bfd24 → 6af114f8`).**
- `BuildSandboxPlan`: cc 107→96, decisions 106→95, lines 356→329,
  nesting 4→4. It is still a hotspot.
- The 6 helpers are unchanged from `1b322d6c`, with cc 18 between them:
  `realPaths` 2, `inProtectedClaudeAccount` 3, `usesClaudeAccount` 2,
  `claudeAccountStateMounts` 4, `nugetPackagesEnv` 5 and
  `coveredInsideExemptMounts` 2, with maximum nesting 2.
- Changed-scope cc went 107→114 (+7): +6 base points and +1 decision point.
  At `1b322d6c` it was +9 (+6 base, +3 decisions). The update removed 2 net
  decision points by reusing a helper.
- The score is +1 because a modified net of −11 is at most −3. The analyzer
  labels it "extraction/redistribution, not net removal". No duplicates,
  coupling or API change.
- Per commit: `ffd5e720` (tests) 0 everywhere; `1b322d6c` complexity +1;
  `6af114f8` 0.

---

## #443 review loop (9 commits)

Coverage: 14 Go files (6 source, 8 test), status `ok` everywhere. Exported
API +0/−0/~0, symbols +19/−0/~14. The +19 is 16 callables and 3 types
(`readyLayer`, `watchedCheck` and `watchRound`, all unexported). No new
edges, no duplicates and no orphans.

**Root reductions** (all 13 modified callables decreased):

| Root | cc | nesting | lines |
|---|---|---|---|
| `ReviewPassService.watchLoop` | 33→19 | 3→2 | 94→66 |
| `ReviewRunCaller.agentDone` | 24→7 | 3→1 | 41→17 |
| `ReviewRunCaller.markStackReady` | 24→15 | 2→2 | 66→34 |
| `OpenReviewLoopClones` | 22→20 | 3→3 | 72→71 |
| `ReviewPassService.addStacked` | 22→14 | 3→3 | 68→45 |
| `ReviewRunCaller.closeOutputs` | 20→15 | 3→3 | 58→44 |
| `ReviewRunCaller.waitingToOpenOn` | 17→6 | 2→2 | 46→16 |
| `ReviewRunCaller.AuthorDone` | 17→11 | 3→2 | 50→39 |
| `watchedPull.watch` | 13→9 | **4→2** | 44→27 |
| `ReviewRunCaller.SettleSteps` | 12→8 | 3→2 | 35→31 |
| `ReviewRunCaller.settleReady` | 12→8 | 3→2 | 32→28 |
| `ReviewRunCaller.settleObservers` | 11→7 | 3→2 | 30→26 |
| `ReviewRunCaller.openSteps` | 10→6 | 3→2 | 29→25 |

**Helpers added (16 callables, cc 101, decisions 85):**

| Helper | cc | nesting | lines | Note |
|---|---:|---:|---:|---|
| `ReviewPassService.readWatchRound` | **15** | 2 | 30 | **new hotspot** (cc ≥ 10), from `watchLoop` |
| `reviewRoundDone` | **13** | 2 | 17 | **new hotspot**, from `agentDone` |
| `ReviewRunCaller.passesStepWaitingOn` | 9 | 2 | 19 | from `waitingToOpenOn` |
| `ReviewRunCaller.approvedOpenLayers` | 9 | 2 | 27 | from `markStackReady` |
| `stackedLayersToAdd` | 9 | 2 | 27 | from `addStacked` |
| `ReviewRunCaller.openStepWaitingOn` | 8 | 2 | 23 | from `waitingToOpenOn` |
| `stewardRoundDone` | 6 | 2 | 13 | |
| `watchedCheck.toWatchCheck` | 6 | 3 | 22 | |
| `stepOwnedLoopHosts` | 5 | 2 | 13 | |
| `ReviewRunCaller.postedPasses` | 4 | **3** | 11 | |
| `ReviewRunCaller.closedOutcome` | 4 | 2 | 8 | |
| `ReviewRunCaller.failWaitingLogged` | 3 | 1 | 5 | |
| `ReviewRunCaller.completeWaitingLogged` | 3 | 1 | 5 | |
| `ReviewRunCaller.roundPRs` | 3 | 1 | 11 | |
| `ReviewRunCaller.cleanPassEvents` | 3 | 2 | 9 | |
| `notApprovedAtHeads` | 1 | 0 | 3 | |

**Root versus helper-inclusive.** The roots lost 92 cc and 196 lines.
Including the helpers, changed-scope cc rose from 237 to 246 (+9). That
breaks down as +16 base points and −7 net decision points. The real
structural change is in **shape, not amount**:
- The largest function went from 33 to 20.
- 7 pre-existing hotspots dropped below the thresholds.
- Maximum nesting fell from 4 to 3. `watchedPull.watch` was the only
  callable at nesting 4.
- Two new cc-13 and cc-15 helpers are themselves hotspots.

`scoring/v1` sees both an improvement (modified net −92, pre-existing
hotspots −7) and a degradation (2 added hotspots, capped at −1), so it
scores complexity **0, "mixed"**. That is the frozen rule working as
designed, not a judgement that nothing improved.

Still hotspots after the batch (8 in scope): the 6 modified roots at
cc ≥ 10 (`OpenReviewLoopClones` 20, `watchLoop` 19, `markStackReady` 15,
`closeOutputs` 15, `addStacked` 14 and `AuthorDone` 11) plus the 2 new
helpers above.

**Per commit** (first-parent; each touches one source file and usually one
test file):

| Commit | cx | Modified cc Δ | Added (cc) | Hotspots |
|---|---:|---:|---|---|
| `55e43bd6` watchedCheck.toWatchCheck | +1 | −4 | 1 (6) | 1→0 |
| `9fc6d90c` readWatchRound | 0 | −14 | 1 (15) | 1→2 (adds a hotspot) |
| `98f9aa08` stackedLayersToAdd | +1 | −8 | 1 (9) | 1→1 |
| `e378fe6c` fail/complete waiting | +2 | −16 | 2 (6) | 4→0 |
| `15ecf47d` ready head check, clean-pass events | +1 | −14 | 4 (17) | 2→2 |
| `2fe55ce7` round done | 0 | −17 | 2 (19) | 1→1 (adds a hotspot) |
| `aaa2c3e3` roundPRs, closedOutcome | +1 | −6 | 2 (7) | 1→1 |
| `58d57839` two waiting-to-open scans | +1 | −11 | 2 (17) | 1→0 |
| `f3ddb233` stepOwnedLoopHosts | 0 | −2 | 1 (5) | 1→1 |

**State and lifecycle trade-offs the tool cannot judge:**
- **Indirection over long-lived loops.** `watchLoop`/`readWatchRound`
  (with a new `watchRound` value type) and `watchedPull.watch` are polling
  loops whose per-round state now crosses a function boundary. Whether
  context cancellation, the "prior read" carried between rounds, and error
  paths behave exactly as before cannot be seen from structure. The author
  reports characterisation tests (watch context/prior read) green on base.
  That is **context, not evaluator evidence**.
- **Settle and close lifecycles.** `failWaitingLogged` and
  `completeWaitingLogged` centralise how waiting review steps end. A defect
  there now affects every caller, which is the usual cost of consolidation.
  The order of "log, then fail or complete" is a behavioural property the
  tool does not see.
- **Clone retention** (`stepOwnedLoopHosts` in `OpenReviewLoopClones`)
  decides which hosts' clones are kept. The author reports mutation checks
  on the waiting-to-open and clone-retention tests. Those are author
  claims, not verified here.
- **AuthorDone** relies on existing handoff tests, with no new
  characterisation test in its commit (`2fe55ce7` adds
  `review_pass_run_round_test.go`). The tool cannot assess coverage of the
  changed paths.
- **Test-first ordering cannot be verified** from history for #443. Each
  test was committed in the same commit as its refactor, unlike #439 and
  #442, which have separate test commits. The claim "green on base" was
  not executed (target code is never run).
- **Rejected extraction.** The author reports declining a renderer-only
  `writeSummary` extraction. Nothing about it appears in the measured
  delta, which is consistent with that report.

---

## #442 follow-up: typed-nil judgement note

The author reports checking every `NewApp` caller. One production caller
passes `events.New()`; the test callers pass an untyped `nil` (a nil
interface) or `events.New()`. No caller passes a typed-nil `*events.Bus`,
and a typed nil is documented as unsupported. No reflection was added.

My own partial check at `c05e6bdb`: by a textual grep, 47 calls pass a
literal `nil`, 4 pass a variable named `bus`, and 1 multi-line call was not
parsed. The production `bus` is `events.New()`. That fits the author's
report but is not a full re-verification. **Status:** this stays a
judgement note on a latent hazard. It was never a measured failure or a
demonstrated bug, and according to the author's caller check it does not
occur today. [TUI-EVIDENCE.md](TUI-EVIDENCE.md) is unchanged apart from a
pointer here.

## Bottom line

- **#439 at `6af114f8`:** the increment reuses `realPaths` for two loops
  (root −2 cc, −8 lines, decision points −2). That is real reuse which the
  duplication analyzer cannot see. In full, the root went 107→96, but
  changed-scope cc still rose by 7 (+6 base, +1 decision), and
  `BuildSandboxPlan` remains a cc-96, nesting-4 hotspot.
- **#443:** a broad change in shape. 13 roots all got simpler (−92 cc,
  −196 lines), the largest function fell from 33 to 20, maximum nesting
  fell from 4 to 3, and 7 pre-existing hotspots cleared. But 16 helpers
  carry cc 101, including 2 new hotspots. Changed-scope cc rose by 9 (+16
  base, −7 decisions), so the frozen score is 0 (mixed). Lifecycle and
  state behaviour across the new function boundaries is not attested.

This is directional structural evidence with no correctness claim. CI and
review are pending.
