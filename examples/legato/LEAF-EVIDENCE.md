# Batch delta: Legato PR #435 (stdlib-only leaf helpers, worklist B1–B3)

This is the first batch from the all-findings follow-up in
[FINDINGS-WORKLIST.md](FINDINGS-WORKLIST.md), measured with the frozen
evaluator. It is **independent of #434**: both branch from `471bfd24`, so
they are not cumulative. The #434 reports ([FIX-EVIDENCE.md](FIX-EVIDENCE.md))
are untouched.

## Provenance

| | |
|---|---|
| PR | cpave3/legato#435, Claude Work coordinator task i95qqzwb, ref `legato/i95qqzwb-leaf-helpers` |
| Base | `471bfd24a69628d985cdd7767dcb340da6a569b7`, tree `ae91f4192b23…` |
| Fix | `9629844901a383e1c2556c3fe955d650b28911c4`, tree `35dc08abb5ff…`, 3 linear commits on top of base (`git merge-base --is-ancestor` verified) |
| Comparisons | The base is the pre-work commit and also `471bfd24`, so the batch-only delta and the drift from `471bfd24` are **the same comparison** |
| Evaluator | Built from `git archive 5e0c679`, tool 0.1.0-poc, `scoring/v1`. The code is identical to the #434 measurement: the only non-example difference to HEAD is `.gitignore`. Config `legato-architecture.json` sha256 `768b2fc6…85f6`, configHash `3cef63633359769f`. Thresholds, exclusions and rules unchanged |
| Object source | Offline bundle `.tmp/i95qqzwb-leaf-helpers.bundle` (`git bundle verify`: ok, complete history, single head = FIX_SHA). It was cloned with `clone --no-checkout --no-local` into `.tmp/legato-leaf`. The clone was not committed and had no worktree files; its status and refs were unchanged after the runs. No network was used and the shared Legato checkout was not touched (HEAD `471bfd24`, status fingerprint `3ab8895f206d466e` before and after) |
| Run | 2026-11-27T19:43:56+11:00: net in 0.85 s, history (net plus 3 steps) in 1.25 s |
| Reports | `leaf-delta-471bfd24..96298449.{json,md}` and `leaf-history-471bfd24..96298449.{json,md}` |

## Result (base → fix)

Coverage: 15 files changed. 12 are Go (8 source, 4 test) and 3 are docs (ADR, AGENTS.md, packages.md). Every analyzer status is `ok`, with no unsupported files and no diagnostics.

| Analyzer | Score | Raw |
|---|---:|---|
| duplication | **+1** | structural redundancy 260→258 (−2); **exact 10→8 (−2)**; 2 groups shrank 2→1, 0 introduced, 0 near-structural candidates; callables considered 6762→6760 |
| complexity | 0 | modified callables net cc 47→47 (0, 4 modified); 1 added (cc 4), 3 deleted (cc 14); changed-scope cc 61→51 (−10); hotspots 3→3; hotspots added 0 |
| dependencies | 0 | units 86→88; internal edges 200→204 (+4/−0); cycles 0→0; rule violations 15→15 (all pre-existing, none new); external refs 1028→1033 |
| obsolete | 0 | nothing newly orphaned |
| scope | 0 | 1 connected unit group; exported API +1 (`regularfile.Open`), ~1 (`retryAfter` → `retryafter.Parse`, matched as a rename); unrelated churn 0 |
| correctness | n/a | not assessed. 4 test files changed (+173/−0), which is context only. CI and review are pending |

Per commit (first-parent; steps are not additive):

| Step | Commit | dup | cx | dep | Struct Δ | Notes |
|---|---|---:|---:|---:|---:|---|
| 1 | `dab61830` add leaf packages, ADR, architecture test | −1 | 0 | 0 | +2 | Transient: the leaf copies exist alongside the originals |
| 2 | `a1f1c8d9` claude and opencode use `retryafter.Parse` | +1 | 0 | 0 | −2 | 2 callables deleted (cc 12), +2 edges |
| 3 | `96298449` claude and codex use `regularfile.Open` | +1 | 0 | 0 | −2 | 2 callables deleted (cc 8), +2 edges |

Step 1's −1 shows the intermediate state of an extract-then-switch
sequence. The net comparison is the one that counts.

## Groups resolved and kept

| Item | Before (`471bfd24`) | After (`96298449`) |
|---|---|---|
| B1 `retryAfter` | `claudeusage.go:197-209` and `opencodeusage.go:185-197`, an exact pair | Both deleted. A single copy is in `internal/leaf/retryafter/retryafter.go:19-31` `Parse`. Group 2→1, so it is no longer redundant. The tool pairs the old claudeusage copy with `Parse` as a cross-package rename, so the finding labels claudeusage as the "remaining copy" |
| B2 `openRegular` | `claudeusage/open.go:17-28` and `codexusage/open.go:17-28`, an exact pair | Both files deleted. The callers `FileCredentials.Login` (claudeusage) and `lastRateLimits` (codexusage) call `regularfile.Open` (`regularfile.go:25-36`, cc 4) |
| B3 `CredentialsPath` / `SessionsDir` | structural pair | **Kept, a justified won't-fix** (4-line provider path resolvers). Still present at `claudeusage/credentials.go:19` and `codexusage/logs.go:26` and still counted in structural redundancy (258). The tool does not judge the justification |

The hypothesis was "exact redundancy −2", and that is what was measured
(10→8). It was not treated as an acceptance target.

## Dependency graph trade-off and rule check

**Graph cost.** The batch adds 2 units and 4 edges:
`claudeusage → retryafter`, `claudeusage → regularfile`,
`codexusage → regularfile` and `opencodeusage → retryafter`. Fan-out goes
from 0 to 2, 1 and 1 for those three engine packages. Both leaves have
fan-in 2 and **fan-out 0** (no internal imports). There are no cycles before
or after. This is the expected price of extraction: 4 edges replace 2 exact
copies. The dependencies analyzer reports edges as info and does not score
them.

**Retained rules (frozen config, 6 distinct rules):**

| Rule (frozen) | Provenance | Violations before→after |
|---|---|---:|
| engine must not import sibling engine packages | test-enforced | 0→0 |
| engine must not import service | CLAUDE.md | 0→0 |
| engine must not import tui | CLAUDE.md | 0→0 |
| service must not import tui | CLAUDE.md | 0→0 |
| tui must not import engine | CLAUDE.md prose, not test-enforced | 14→14 (pre-existing drift, G) |
| config is standalone | CLAUDE.md prose | 1→1 (`config → engine/macros`, pre-existing, G) |

The new engine→leaf edges violate none of these. The sibling rule's `to`
pattern is `internal/engine/**`, and `internal/leaf/**` is outside it.

**Rule/version mismatch (reported, not tuned).** The frozen config predates
ADR 0001 and does not encode two parts of the authorised rule:

1. *A leaf imports only the standard library.* There is no frozen rule
   `internal/leaf/** → internal/**`, and the evaluator has no
   "stdlib-only" rule type. As a measured fact, both leaves have internal
   fan-out 0. As a manual read of the tracked blobs at FIX, the imports are
   `net/http`, `strconv`, `strings` and `time` (retryafter), and `errors`,
   `os` and `syscall` (regularfile): standard library only.
2. *Of Legato's own packages, an engine package imports only leaf
   packages.* The frozen config can express only deny rules, not an allow
   list. As a manual check, `git grep` at FIX finds no non-test engine
   import of any non-engine Legato package except the two leaves (2 imports
   each).

So the batch conforms to the ADR rule according to these manual checks, but
the frozen evaluator **cannot attest to that**. Expressing it would need a
new config version (a deny rule `internal/leaf/** → internal/**`, plus an
allow-list or stdlib-only rule type). This is deferred and not applied,
because changing config mid-measurement would break comparability. Legato's
own `TestLeafPackagesImportOnlyTheStandardLibrary` (an added test function,
listed as context) is the enforcement of record.

## What the tool did not see, or that needs judgement

- **Semantic equivalence is not checked.** The tool shows that the copies
  were byte-identical before (exact fingerprint) and that the callers now
  route through one function. It does not show that behaviour is preserved.
  The ADR's distinction that Slack, Discord and Jira Retry-After readers are
  look-alikes with different semantics is a judgement the tool cannot make.
  Those readers were not in the flagged groups.
- **The exported API widened by 2 symbols** in new packages
  (`retryafter.Parse` and `regularfile.Open`), replacing 2 unexported
  copies. This is reported, not judged.
- **Docs and test changes** (ADR, AGENTS.md, packages.md, architecture test,
  leaf tests) are counted in scope and coverage only.

## Bottom line

The frozen measurements match the batch's claims. Both exact cross-engine
copies (B1 and B2) are gone (exact 10→8, structural 260→258). Changed-scope
complexity fell by 10 with no new hotspots. There is no orphaned code and no
unrelated churn. The graph grew by 4 edges into two stdlib-only leaves with
zero fan-out, and no rule was newly violated and no cycle formed. B3 remains
as a stated won't-fix. Under `scoring/v1` this scores duplication +1 and
everything else 0. Correctness is not assessed.
