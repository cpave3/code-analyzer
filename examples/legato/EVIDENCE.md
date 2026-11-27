# Legato evaluation: evidence notes

This is the headline real-world evaluation. It measures the quality delta in
Legato's Go code from the start of 2026-11-26 to local HEAD on 2026-11-27.
It shows what the measurements can see. It is not a scientific validation,
not a merge gate, and not a verdict on Legato's overall quality.

## Provenance

| | |
|---|---|
| Source | Read-only local clone of Legato (the operator-provided resource) |
| Window pinned by | the operator, 2026-11-27 16:44 AEDT (+11:00) |
| Base | `9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc`, the last first-parent ancestor before 2026-11-26T00:00:00+11:00 (committed 2026-11-25T21:49:11+11:00), tree `a55db6103bf4…` |
| Target | `471bfd24a69628d985cdd7767dcb340da6a569b7` (merge #431), the local HEAD at capture time (committed 2026-11-27T15:46:34+11:00), tree `ae91f4192b23…`. This was the local HEAD; it does not prove this is the newest remote main. Nothing was fetched. |
| First-parent steps | 23 merges (#389 … #431). 172 commits in total including branch commits. Branch commits are not double counted. |
| Evaluator | codeeval 0.1.0-poc, `scoring/v1`, config `legato-architecture.json` (hash in reports) |
| Run | 2026-11-27T17:11:25+11:00. History (net plus 23 steps) took 10.0 s with a cold cache: 2,602 files parsed, 50,796 cache hits. The net evaluation alone took 0.8 s with a warm cache, or about 1.9 s cold. |

The architecture rules in `legato-architecture.json` come from Legato's own
`CLAUDE.md` at the target commit. Each rule's name says whether Legato's own
architecture tests enforce it (only the engine-sibling rule is test-enforced)
or whether it appears only in the prose.

## Reproduce

```bash
# from the code-evaluator repo
scripts/legato-eval.sh /home/cameron/Projects/legato /tmp/codeeval-cache
# equivalent explicit commands:
go build -o bin/codeeval ./cmd/codeeval
bin/codeeval history  --repo /home/cameron/Projects/legato \
  --base 9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc --target 471bfd24a69628d985cdd7767dcb340da6a569b7 \
  --config examples/legato/legato-architecture.json --cache-dir /tmp/codeeval-cache --timing \
  --output examples/legato/history-9a61e9f8..471bfd24.json --summary examples/legato/history-9a61e9f8..471bfd24.md --format none
bin/codeeval evaluate --repo /home/cameron/Projects/legato \
  --base 9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc --target 471bfd24a69628d985cdd7767dcb340da6a569b7 \
  --config examples/legato/legato-architecture.json --cache-dir /tmp/codeeval-cache --timing \
  --output examples/legato/net-9a61e9f8..471bfd24.json --summary examples/legato/net-9a61e9f8..471bfd24.md --format none
```

Outputs:
- `net-9a61e9f8..471bfd24.{json,md}`: the net base→target report
- `history-9a61e9f8..471bfd24.{json,md}`: a per-merge table and the full report for each step, followed by the net report

Rerunning reproduces the same results. Only the `timing` block differs.

## Source state checks

- The evaluator only ran Git plumbing against Legato. No checkout, fetch,
  hooks, scripts, builds or tests were run. Only tracked content at the two
  pinned commits was read; untracked files were never read or listed.
- The status fingerprint was `sha256(git status --porcelain=v1 --untracked-files=all)`,
  and only the hash was printed, never filenames. It was
  `3ab8895f206d466e…` before the first run and after every run.
  `HEAD = 471bfd24…` throughout.
- Disclosure: the very first fingerprint (16:44) was taken with plain
  `git status`, which may refresh the stat cache in `.git/index`. It does
  not touch tracked content, refs or the worktree. Every later check, and
  the evaluator itself, used `--no-optional-locks`.

## Coverage: Go only, so this is not Legato as a whole

| Language | Files at target | Changed files | Changed lines | Analysed? |
|---|---:|---:|---:|---|
| Go | 2,207 (988 source / 1,219 test) | 363 (201 source / 162 test) | 32,703 | yes |
| TypeScript | 999 | 171 (84 source / 87 test) | 9,556 | **no** (unsupported) |
| SQL migrations | 160 | 7 | 36 | no |
| JavaScript / Kotlin | 33 / 57 | 2 / 2 | 119 / 104 | no |

Non-code changes were 96 docs files, 35 "other" files (images, JSON, YAML)
and 2 generated files. Every code analyzer reports `partial` coverage:
201 of 296 changed source files were analysed, and 95 were unsupported
(TS 84, SQL 7, JS 2, Kotlin 2). Six of the 23 merges (#403, #407, #420,
#423, #428, #431) changed no Go source. For those, every code analyzer
returns `unsupported` with a null score, not a neutral 0. **None of the
findings below say anything about the TypeScript, web or mobile code.**

## Net result (base → target)

| Analyzer | Score | Raw measurement |
|---|---:|---|
| duplication | -2 | structural redundancy 243→260 (Δ+17), exact 8→10 |
| complexity | -2 | modified pre-existing callables net cyclomatic Δ+194 (129 up, 24 down). 569 callables added (cc 2,976), 6 deleted. Pre-existing hotspot Δ+7, new hotspot callables 84 |
| dependencies | 0 | internal edges 193→200 (+7/−0), cycles 0→0, rule violations 15→15 (all pre-existing) |
| obsolete | 0 | no newly orphaned functions |
| scope | 0 | 680 files; symbols +707/−6/~326; exported API +280/−4/~1; one connected group of changed units |
| correctness | n/a | not assessed |

How to read this: for a two-day window with this much feature work
(+60.6k/−4.0k lines), the -2 scores mostly reflect **volume**. The
thresholds are absolute, so a large window saturates (see limits). The
useful part is the evidence behind each score, not the score itself.

## Per-merge attribution (first-parent)

These are observed contributions to the metrics. They do not show that a
workflow or a person caused a quality change. Full table:
`history-9a61e9f8..471bfd24.md`.

- **#419 subscription-usage** (dup -2, cx -2): added `openRegular`
  identically in `internal/engine/claudeusage/open.go:17` and
  `internal/engine/codexusage/open.go:17`. Struct redundancy +6.
- **#429 opencode-usage** (dup -1): added `retryAfter` in
  `internal/engine/opencodeusage/opencodeusage.go:185`, an exact copy of the
  pre-existing `claudeusage.retryAfter` (`claudeusage.go:197`). The finding
  marks it as a reuse candidate.
- **#430 claude-accounts** (cx -2): `agentService.SpawnAgent` went from
  cc 273 to 287 and from 820 to 857 lines; `prepareTaskResume` went from
  cc 34 to 48. Modified net cc +148.
- **#418 deprecate-inline-step** (cx -2, dup -2): `ReviewPassService.watchLoop`
  went from cc 17 to 33 and `ReviewRunCaller.AuthorDone` from 9 to 17. It
  also reduced `pollComments` from 38 to 17 and removed a duplicate
  `Server.writeAttentionError`. The step is mixed, and the report shows both
  sides.
- **#395 pr-ready** (dup -1, cx 0 mixed): `ReviewRunCaller.Details` went from
  cc 26 to 3 (extraction) and pre-existing hotspots fell by 2, against 24 new hotspot callables among 167 additions. The score is 0 because the signals are mixed.
- **#417 related-work-runs** (dup +1): removed the duplicate
  `Store.ListWorkflowRunsByRelatedTask`, with base evidence.
- **#425 usage-profiles** (obsolete +1): deleted a function that was already
  unreferenced. Complexity -2: `SpawnAgent` cc 266→273 and
  `formatSubscriptionUsage` cc 11→19.

## Useful signals

1. **Copy-paste across sibling packages.** `retryAfter`, `openRegular`,
   `CredentialsPath`/`SessionsDir` and `claudeUsageError`/`opencodeUsageError`
   were added as identical or structurally identical code across the
   `claudeusage`, `codexusage` and `opencodeusage` packages. This matches
   Legato's own rule that engine packages must not import their siblings: the
   rule pushes toward duplication unless there is a shared lower-level
   package. That is a real design tension the tool surfaced with exact
   evidence.
2. **Growth of existing god functions.** `SpawnAgent` (cc 266→287, 857
   lines), `App.update` (cc 304→306, 1,125 lines) and `BuildSandboxPlan`
   (94→107) are pinned with before and after locations. Tests and
   self-assessment would not show this.
3. **Extraction is credited correctly.** Examples are `ReviewRunCaller.Details`
   (26→3), `ReviewPublisher.PollComments` (18→2) and `BuildGraph` (17→3).
   Where helpers absorb the complexity, the report says so.
4. **The dependency picture is precise.** Seven new edges are listed with
   their import lines: `cmd/legato` and `internal/service` gained three new
   engine usage-reader packages, and `internal/tui` gained a dependency on
   `config`. There were no cycles and no new rule violations.
5. **Documentation drift.** Encoding the CLAUDE.md prose rule "tui … never
   engine/" literally yields 14 pre-existing violations (for example
   `internal/tui/app.go` imports `engine/store`, `engine/workflow` and
   `engine/ipc`), plus one in `config` (`config → internal/engine/macros`).
   Legato's architecture tests enforce the engine-sibling rule (0 violations)
   but not this one. The violations are pre-existing and not scored.
6. **API surface.** +280 exported symbols, 4 removed and 1 signature change
   (`Store.FindTaskConversation` gained a `storageDir` parameter), with
   name-based caller counts.

## Misleading or limited signals (false positives and limits)

1. **Absolute thresholds saturate.** Any large feature window scores -2 on
   duplication and complexity. Compare raw deltas and per-symbol evidence,
   not the scores of windows of different sizes.
2. **Idiomatic boilerplate counts as structural duplication.** Thin `Store.List*`
   wrappers, paired `StopWatching`/`WatchAgain` methods and `UnmarshalYAML`
   pairs are flagged as structural clones. They are idioms, not missed reuse.
   Of the 19 "introduced" groups at the net level (2 exact, 17
   structural), my manual read is that about half look like real copy-paste
   (`retryAfter`, `openRegular`, `accountPathContains`/`pathWithin`,
   `doPut`/`doPost`, `*UsageError`, `load*Usage`, `Start*Events`) and the
   rest like idiom. Raising `minTokens`, or treating
   short thin wrappers separately, would cut the noise. The defaults were not
   tuned on Legato.
3. **New hotspots in feature code.** 84 new callables at cc ≥ 10 or nesting
   ≥ 4 is a real count, but they cap the score at -1 and are shown
   separately from growth in pre-existing code.
4. **Go only.** 84 changed TypeScript source files (+87 test files) are
   invisible. Merges #403, #407, #420, #423, #428 and #431 are unscored.
5. **References are name-based.** `obsolete` found nothing at the net level.
   Low-confidence method candidates are reported and not scored. Reflection
   and interface satisfaction are invisible.
6. **Correctness is not assessed.** 162 Go test files changed (+19.9k
   lines). That is context, not evidence.
7. **Step reports are not additive.** Repository-wide clone groups and the
   graph interact across merges, so the per-merge redundancy deltas sum to
   more than the net +17.

## Follow-up: fix delta

The cleanup PR #434 (`471bfd24 → 41c80614`) was measured with the frozen evaluator; see [FIX-EVIDENCE.md](FIX-EVIDENCE.md).
The independent leaf-helper batch PR #435 (`471bfd24 → 96298449`, worklist B1–B3) is in [LEAF-EVIDENCE.md](LEAF-EVIDENCE.md).
The service, cmd and test-fake duplicate batch PR #437 (batch only `8661dc42 → e0cf802e`, plus `471bfd24 → e0cf802e` including the #436 ledger, worklist C1, C3, C5 and D8) is in [SVC-EVIDENCE.md](SVC-EVIDENCE.md).
Three independent batches, each from `471bfd24`, are in [BATCHES-439-441-EVIDENCE.md](BATCHES-439-441-EVIDENCE.md): #439 sandbox-plan (`1b322d6c`), #440 config-macros (`7b6de497`) and #441 small-extracts (`b9696a9c`).
The independent TUI-boundary batch PR #442 (`471bfd24 → c05e6bdb`, worklist G) is in [TUI-EVIDENCE.md](TUI-EVIDENCE.md).
The #439 updated head (`6af114f8`, with the `1b322d6c → 6af114f8` increment) and the independent review-loop batch #443 (`471bfd24 → f3ddb233`) are in [SANDBOX-REVIEWLOOP-EVIDENCE.md](SANDBOX-REVIEWLOOP-EVIDENCE.md).
Three independent batches, each from `471bfd24`, are in [BATCHES-444-446-EVIDENCE.md](BATCHES-444-446-EVIDENCE.md): #444 sandbox-prefix (`94ce0d0d`), #445 Claude account resolution (`3faa8213`) and #446 resume preparation (`8a2e03d7`).
**The final combined stack** (11 layers, `471bfd24 → bb13a424`, excluding #434), with per-layer comparisons and the residual drift from `9a61e9f8`, is in [FINAL-STACK-EVIDENCE.md](FINAL-STACK-EVIDENCE.md).

## Bottom line

The POC can tell apart, with precise and reproducible evidence: copy-paste
versus reuse, extraction versus nesting, consolidation versus leaving copies
behind, new coupling versus decoupling, and focused versus churny changes
(shown on the fixtures in both directions). On Legato it surfaced concrete,
reviewable items: cross-package copies, growth of existing god functions, and
drift between documented and actual layering. Those items are what is worth
acting on. The directional scores are a coarse index into that evidence. They
are not a quality verdict, especially for large windows and while the tool
covers Go only.
