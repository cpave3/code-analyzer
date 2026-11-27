# Batch delta: Legato PR #437 (service, cmd and test-fake duplicates, worklist C1, C3, C5, D8)

This is the second batch of the all-findings follow-up, measured with the
frozen evaluator. It is stacked on #436, a docs-only ledger. It is
**independent of #434 and #435**: neither is an ancestor, and both have
separate reports ([FIX-EVIDENCE.md](FIX-EVIDENCE.md) and
[LEAF-EVIDENCE.md](LEAF-EVIDENCE.md)).

## Provenance

| | |
|---|---|
| PR | cpave3/legato#437, task i95qqzwb, ref `legato/i95qqzwb-service-dups`, stacked on #436 `legato/i95qqzwb-ledger` |
| Pinned main | `471bfd24a69628d985cdd7767dcb340da6a569b7`, tree `ae91f4192b23…`. This is the pinned benchmark commit. It is **not** claimed to be the current remote main |
| Batch base (#436 head) | `8661dc42a368d137f807798ae52625253029b059`, tree `330af3951166…`. The older pre-restack `7d9862b6` named in the PR body is **not** used as the batch base. It appears only as the first ledger commit on the first-parent path from main |
| Fix | `e0cf802e15cf3c2007285cbdc43e759073b1ebbb`, tree `fad541164613…` |
| Ancestry | `471bfd24 → 7d9862b6 → 8661dc42` (2 docs commits, #436) `→ 0242ef2a → 16d79149 → a3e7f445 → e0cf802e` (4 commits, #437). `merge-base --is-ancestor` verified for both links |
| Evaluator | Built from `git archive 5e0c679`, tool 0.1.0-poc, `scoring/v1`. Config `legato-architecture.json` sha256 `768b2fc6…85f6`, configHash `3cef63633359769f`. Unchanged, nothing tuned |
| Object source | Offline bundle `.tmp/i95qqzwb-service-dups-e0cf802e.bundle`. `bundle verify`: complete history, heads `service-dups@e0cf802e` and `ledger@8661dc42`. Cloned with `--no-checkout --no-local` into `.tmp/legato-svc`, which has no worktree files and is not committed; its status and refs were unchanged after the runs. No network. The shared Legato checkout was untouched: HEAD `471bfd24`, status fingerprint `3ab8895f206d466e` before and after. The measured pins are recorded above whatever the branches do later |
| Run | 2026-11-27T19:54:47+11:00. Evaluate took 0.8 and 0.9 s; history took 2.1 and 2.6 s |

## Reports

| Comparison | Net | History (first-parent) |
|---|---|---|
| (2) **Primary, batch only:** `8661dc42 → e0cf802e` | `svc-batch-8661dc42..e0cf802e.{json,md}` | `svc-history-batch-8661dc42..e0cf802e.{json,md}` (4 steps) |
| (1) Operator-requested, pinned main to fix, **includes the #436 ledger**: `471bfd24 → e0cf802e` | `svc-main-471bfd24..e0cf802e.{json,md}` | `svc-history-main-471bfd24..e0cf802e.{json,md}` (6 steps) |

## Result

| Analyzer | (2) batch `8661→e0cf` | (1) main `471b→e0cf` | Raw (identical in both) |
|---|---:|---:|---|
| duplication | **+2** | **+2** | structural redundancy 260→255 (−5): 6 redundant instances removed, 1 introduced; exact 10→10; callables considered 6762→6757; 0 near-structural candidates |
| complexity | **+2** | **+2** | 11 modified callables, all decreased: net cc 46→17 (−29), lines 151→46, max nesting 2→1. 4 added (cc 14), 0 deleted. Changed-scope cc 46→31 (−15), max cc 8→4, hotspots 0→0 |
| dependencies | 0 | 0 | edges 200→200 (+0/−0), cycles 0→0, rule violations 15→15 (pre-existing) |
| obsolete | 0 | 0 | nothing newly orphaned |
| scope | 0 | 0 | symbols +5/−0/~11; exported API +0/−0/~0; 1 unit group; unrelated churn 0 |
| correctness | n/a | n/a | not assessed. 3 test files changed (+156/−0), which is context only. CI and review are pending |

**Why the two comparisons differ only in scope counts.** #436 changes one
docs file (`docs/analysis/codeeval-findings-ledger.md`, +347/−0) and no
code. Comparison (1) therefore reports 11 files and +615/−121 lines,
against 10 files and +268/−121 for comparison (2). It also records the
extra file as a non-code file in each analyzer's coverage. Every
duplication, complexity, dependency and obsolete measurement and finding is
identical after normalising the base SHA (verified by a JSON diff). In the
main history the two ledger steps score `ok`/0 on every code analyzer, with
coverage showing `changedCodeFiles: 0`. That 0 means *no code in scope*,
not a neutral quality judgement. The frozen engine marks a step
`unsupported` only when it has unsupported code, so docs-only steps appear
as `ok`. Read their coverage, not their score.

Per commit (first-parent, not additive):

| Commit | dup | cx | Struct Δ | Modified cc Δ | Added (cc) |
|---|---:|---:|---:|---:|---|
| `7d9862b6` ledger (docs) | 0 | 0 | 0 | 0 | none; no code |
| `8661dc42` ledger update (docs) | 0 | 0 | 0 | 0 | none; no code |
| `0242ef2a` usage errors become a reason table (C1) | 0 | +2 | 0 | −10 | `usageError` (4) |
| `16d79149` PullRequestModes read their rules in one place (C5) | +1 | 0 | −1 | −2 | `rules` (3) |
| `a3e7f445` runtime roles share `claimDataDirLock` (C3) | +1 | +1 | −2 | −9 | `claimDataDirLock` (4) |
| `e0cf802e` review fake uses `editPR` (D8) | +1 | +1 | −2 | −8 | `editPR` (3) |

## Group by group

| Item | Before (`8661dc42`) | After (`e0cf802e`) | Measured |
|---|---|---|---|
| C3 `claim*` runtime trio | `claimChatBots` `chatbots_runtime.go:68-82`, `claimMQTTPublisher` `mqtt_runtime.go:27-41`, `claimUsagePoller` `usage_runtime.go:32-46`: each cc 4, nesting 2, 15 lines | Each is a cc-1, 3-line delegate to the new `claimDataDirLock` `runtime_lock.go:16-30` (cc 4) | Group 3→0 |
| D8 `FakeReviewGitHub` | `Draft`/`Push` and `SetAuthor`/`Describe` pairs, cc 3, 10 lines each | Each is a cc-1, 3-line delegate to the new `editPR` `reviews.go:141-150` (cc 3) | 2 groups 2→0 |
| C5 PR rules | `PullRequestMode` `pull_request_modes.go:23-33` and `ReadyAfterReview` `:45-55`, cc 3 each | cc 2 (7 lines) and cc 2 (4 lines). The nil check, read, log and fallback move into the new `rules` `:34-44` (cc 3) | Group 2→0 |
| C1 `*UsageError` | `claudeUsageError` `subscription_usage_readers.go:52-73` and `opencodeUsageError` `:176-197`, cc 8, 22 lines each | cc 3, 7 lines each, plus the new `usageError` `:82-92` (cc 4) and the type `usageErrorReason` `:74-77`. Each provider keeps its own reason table | Old group 2→0, but the **new 7-line wrappers are still structurally identical**, so 1 redundant instance is introduced. Net: the pair shrank and was reduced, not eliminated |

The worklist items left out of this batch are C2, C4 and D1–D7. The ledger
records them as won't-fix, and they are unchanged and still counted in the
remaining redundancy of 255.

## Trade-offs the scores do not show

- **Wrappers introduced and kept.** The three `claim*` functions and four
  `FakeReviewGitHub` methods are now one-line delegates. The author kept
  the names because callers and #432's tests use them. They are under
  `minTokens`, so they do not count as duplication, and the tool takes no
  view on whether the extra layer adds value. The call depth is one level
  deeper at each site.
- **Error mapping becomes data plus a lookup.** C1 replaces two explicit
  `switch` statements with per-provider `[]usageErrorReason` tables and a
  shared first-match loop. Part of the drop in cyclomatic complexity (8→3
  per wrapper) is complexity **moved into data**, which cc does not count:
  table entries are not decision points. Correctness now depends on the
  order of table entries ("first matching reason"). The tool cannot verify
  that the order and reasons match the old switch. The author's
  characterisation test `TestProviderUsageErrorsGiveReasons` is context,
  not evaluator evidence.
- **Rule lookup centralised.** `rules(ctx)` now owns the nil, read, log and
  fallback path for both questions. Behaviour such as "each question reads
  the rules once" and the logged message is claimed by the author and
  tested in `TestPullRequestModesAnswerFromTheRulesOrNothing`. The tool
  does not measure it.
- **Runtime lock sharing.** All three roles now go through one flock helper.
  The lock file names and release-once semantics are said to be preserved,
  but the tool does not execute code. A defect in `claimDataDirLock` would
  now affect all three roles at once, which is the usual trade-off of
  consolidation. `TestRuntimeRolesLockTheirOwnFiles` is context.
- **Test-support code counts as production source.** `internal/handofftest/reviews.go`
  is not a `_test.go` file, so the evaluator classifies it as source. Two of
  the 6 removed redundant instances and −8 of the modified cc come from a
  test fake. That is real cleanup, but it is not product code.
- **No coupling or API change.** There are no new imports and no exported
  symbols added, removed or changed. Added callables: 4 unexported, total
  cc 14.
- **Coverage.** Go only. All 10 batch files are Go, so this batch has no
  unsupported-language gap.

## Bottom line

On the primary batch-only comparison, #437 removed 6 redundant instances
and introduced 1 (net structural redundancy −5; the exact count is
unchanged at 10). The 11 callables it touched all lost complexity (net cc
−29, changed-scope cc −15) with no new hotspots, coupling or API change.
The costs: 4 new helpers and a type, 7 delegate wrappers kept by design,
complexity partly moved into data tables in C1, a C1 pair that still reads
as structurally identical, and a third of the gain coming from a test fake.
The operator-requested comparison from pinned main has the same code
measurements. It adds only the #436 ledger docs file, which is shown in
scope and coverage. Correctness is not assessed: behaviour preservation
rests on the author's tests and on CI and review, which are pending.
