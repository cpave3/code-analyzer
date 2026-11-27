# Batch deltas: Legato PRs #444, #445 and #446 (independent, each from `471bfd24`)

These are three independent frozen evaluations from pinned main `471bfd24`.
**Do not add them together**, or to any earlier batch. Earlier reports and
pins are unchanged.

## Provenance (shared)

| | |
|---|---|
| Base | `471bfd24a69628d985cdd7767dcb340da6a569b7`, tree `ae91f4192b23…`. This is the pinned benchmark main, not a claim about the current remote main |
| #444 sandbox-prefix | `94ce0d0ddc8035343592ed1e92bd812ac5954cd0`, tree `fae5e81375fe…`, ref `legato/i95qqzwb-sandbox-prefix`. `f6468797` adds the test file only; `94ce0d0d` changes `agent_sandbox.go` **and** the test file (+40 test lines) |
| #445 claude-account-resolve | `3faa8213318e0d28ce792b365b48a75477d2cb2d`, tree `b6a37cdcfbe8…`, ref `legato/i95qqzwb-claude-account-resolve`. `e37bf8f3` adds the test file only; `3faa8213` changes only `agent_accounts.go` |
| #446 prepare-resume | `8a2e03d7d9373b5e2295b9bcc019601dab138b87`, tree `9d18e5608c18…`, ref `legato/i95qqzwb-prepare-resume`. `b516c13c` adds the test file only; `8a2e03d7` changes only `agent_resume.go` |
| Evaluator | Built from `git archive 5e0c679`, tool 0.1.0-poc, `scoring/v1`. Config sha256 `768b2fc6…85f6`, configHash `3cef63633359769f`. Unchanged, nothing tuned |
| Object source | Offline bundle `.tmp/i95qqzwb-batches444-445-446.bundle`. `bundle verify`: complete, and the 3 heads equal the pins. `merge-base --is-ancestor` verified for each. Cloned with `--no-checkout --no-local` into `.tmp/legato-b6`, which has no worktree files and is not committed; it was unchanged after the runs. No network. The shared Legato checkout was untouched: HEAD `471bfd24`, fingerprint `3ab8895f206d466e` before and after |
| Run | 2026-11-27T20:19:10+11:00. Each evaluate took 0.65–0.68 s and each history 1.0 s |
| Per-callable detail | Helper cc, decisions and nesting come from the evaluator's own index cache (adapter `go-adapter/3`). cc = decisions + 1 per callable |

| Batch | Net | History |
|---|---|---|
| #444 | `sandboxprefix-471bfd24..94ce0d0d.{json,md}` | `sandboxprefix-history-471bfd24..94ce0d0d.{json,md}` |
| #445 | `claudeaccount-471bfd24..3faa8213.{json,md}` | `claudeaccount-history-471bfd24..3faa8213.{json,md}` |
| #446 | `prepresume-471bfd24..8a2e03d7.{json,md}` | `prepresume-history-471bfd24..8a2e03d7.{json,md}` |

**Test-first ordering:** in all three batches, a commit that adds only a test
file comes before the refactor on the first-parent path. Each of those test
commits scores 0 everywhere and changes no production code. #444's refactor
commit also edits its test file, adding 40 lines.

## Scores and raw complexity

| | #444 sandboxPrefix | #445 resolveClaudeAccount | #446 prepareTaskResume |
|---|---:|---:|---:|
| duplication | 0 | 0 | 0 |
| **complexity** | **0 (mixed)** | **0 (mixed)** | **0 (mixed)** |
| dependencies / obsolete / scope | 0 / 0 / 0 | 0 / 0 / 0 | 0 / 0 / 0 |
| correctness | n/a | n/a | n/a |
| Root cc | 67 → 58 (−9) | 55 → 28 (−27) | 48 → 28 (−20) |
| Root lines | 216 → 192 | 131 → 68 | 127 → 81 |
| Root max nesting | 3 → 3 | 3 → 3 | **3 → 2** |
| Added helpers | 1, cc 10 | 3, cc 33 | 4, cc 26 |
| **Changed-scope cc** | **67 → 68 (+1)** | **55 → 61 (+6)** | **48 → 54 (+6)** |
| of which: base points (+1 per new function) | +1 | +3 | +4 |
| of which: branch/decision points (net) | **0** | **+3** | **+2** |
| Hotspots in scope | 1 → 2 (**+1 new**) | 1 → 2 (**+1 new**) | 1 → 2 (**+1 new**) |
| Exported API | +0/−0/~0 | +0/−0/~0 | +0/−0/~0 |
| Edges / cycles / rule violations | unchanged | unchanged | unchanged |
| Duplicates introduced / removed | 0 / 0 | 0 / 0 | 0 / 0 |

In every batch, complexity was **redistributed, not removed**. The root
lost decision points, and the helpers took on at least as many: exactly as
many in #444, and 3 and 2 more in #445 and #446. The rest of the rise is the
+1 baseline of each new function. Each batch also creates one new hotspot
helper. Under `scoring/v1` that offsets the root improvement, so all three
score complexity **0 ("mixed")**. That is the frozen rule working as
designed, not a verdict that nothing improved.

---

## #444 sandboxPrefix → withClaudeAccounts

| Callable | cc | decisions | nesting | lines |
|---|---:|---:|---:|---:|
| `agentService.sandboxPrefix` (root) | 67 → 58 | 66 → 57 | 3 → 3 | 216 → 192 |
| `withClaudeAccounts` (new) | **10** | 9 | 2 | 28 |

- **Pure extraction, as the metric reads it.** Exactly the 9 decision
  points that left the root reappear in the helper (net 0). The only
  increase is the helper's +1 base point.
- `withClaudeAccounts` sits **exactly at** the cc-10 hotspot threshold, so
  it counts as a new hotspot. That is a threshold-edge effect, but it is
  the frozen rule.
- `sessionEnvValues` was touched without changing its tokens (moved within
  the file): formatting-only, with no unrelated churn found.
- The root is still a cc-58, 192-line hotspot.

## #445 resolveClaudeAccount → named-account guards

| Callable | cc | decisions | nesting | lines |
|---|---:|---:|---:|---:|
| `agentService.resolveClaudeAccount` (root) | 55 → 28 | 54 → 27 | 3 → 3 | 131 → 68 |
| `agentService.checkNamedAccountIsolation` (new) | **22** | 21 | 2 | 60 |
| `checkNamedAccountLogin` (new) | 7 | 6 | 1 | 10 |
| `namedAccountEnvConflict` (new) | 4 | 3 | 2 | 11 |

- The root's cc halves (55→28; the author's gocyclo figure matches). Of the
  27 decision points it lost, 30 reappear in the three helpers: net **+3**.
- `checkNamedAccountIsolation` is a **new cc-22 hotspot**, 60 lines long.
  It takes most of the moved logic. The two hotspots after the batch are
  the root (28) and this helper.
- **This is security-relevant code**: account isolation, refusals, login
  checks, and environment conflicts for named accounts. The evaluator
  measures structure only. Refusal text and order, account identity, path
  handling and isolation semantics are **not** attested by any metric. They
  rest on tests and review.

## #446 prepareTaskResume → split by responsibility

| Callable | cc | decisions | nesting | lines |
|---|---:|---:|---:|---:|
| `agentService.prepareTaskResume` (root) | 48 → 28 | 47 → 27 | **3 → 2** | 127 → 81 |
| `agentService.resumeConversation` (new) | **12** | 11 | 2 | 27 |
| `agentService.bindResumeAccount` (new) | 6 | 5 | 2 | 17 |
| `agentService.inferResumeKind` (new) | 5 | 4 | 1 | 14 |
| `resumeStorage` (new) | 3 | 2 | 1 | 9 |

- The root went 48→28 and lost one level of nesting. The helpers match the
  author's 12/6/5/3. Of the 20 decision points the root lost, 22 reappear
  in the helpers: net **+2**.
- `resumeConversation` is a **new cc-12 hotspot**. After the batch the
  hotspots are the root (28) and this helper.
- `prepareTaskResume` was one of the worklist's top growth items (cc
  28→48 inside the original window). This batch brings the root back to 28,
  but the whole resume path in scope is 54, compared with 48 before the
  batch.

---

## Precise limits

- **Correctness and behaviour are not assessed.** The evaluator does not
  run code. Exact refusal text and order, account selection, isolation,
  resume semantics and sandbox mount rules are outside every metric.
- **Branch points versus logic.** The decision count is a syntactic count
  (`if`, `for`, `case`, `&&`, `||`). A net +2 or +3 may be checks repeated
  at helper boundaries rather than new behaviour. The evaluator cannot tell
  which.
- **The hotspot threshold is absolute** (cc ≥ 10 or nesting ≥ 4). #444's
  helper is a hotspot at exactly 10.
- **Test volume is context only.** The tests added are +188, +630 and +436
  lines respectively. The evaluator counts files and lines, not test
  functions, table cases or mutations.

## Test and mutation evidence: who produced it

| Claim | Source | Status here |
|---|---|---|
| #444: tests-first commit `f6468797`, `withClaudeAccounts` is a pure extraction | Author | Commit order **verified from history**. "Pure" is consistent with net 0 decision points, but behaviour is not verified |
| #445: 99 exact text/order cases and 27 mutation failures in `e37bf8f3` | Author | **Not verified.** The evaluator does not count cases or run mutations. Commit order verified |
| #446: 8 test functions and 5 mutations in `b516c13c` | Author | **Not verified.** Commit order verified |
| #437, #441, #443: the branch characterisation test files pass on the exact base `471bfd24`; #443's lifecycle set passes on base and fix; mutation checks reported earlier | **Coordinator-reported** (reproduced by the coordinator, not historical commit-order proof) | **Not independently executed by me.** I never ran Legato code. This is recorded as reported context and does not change the earlier reports. For #437, #441 and #443, history still cannot show test-first order, because their tests share commits with the refactors |

## Bottom line

All three batches roughly halve or substantially shrink a large root
function: 67→58, 55→28 and 48→28. #446 also removes a nesting level. In
each case the logic moved into helpers, not out of the code. Changed-scope
cc went +1, +6 and +6, with net decision points 0, +3 and +2. Each batch
creates one new hotspot helper (cc 10, 22 and 12), so the frozen
complexity score is 0 (mixed) for all three. There is no duplication,
dependency, API or obsolete-code change. Correctness, especially #445's
account-isolation and refusal behaviour, is not attested. CI and review
are pending.
