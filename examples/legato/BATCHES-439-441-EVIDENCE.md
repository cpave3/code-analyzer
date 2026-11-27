# Batch deltas: Legato PRs #439, #440 and #441 (independent, each from `471bfd24`)

These are three independent batches of the all-findings follow-up. Each is
measured **separately** against pinned main `471bfd24` with the frozen
evaluator. They do not stack with each other or with #434, #435 or #437.
**Do not add their deltas together**: the clone groups and the import graph
are repository-wide, and the batches touch overlapping code. Earlier reports
are unchanged.

## Provenance (shared)

| | |
|---|---|
| Base | `471bfd24a69628d985cdd7767dcb340da6a569b7`, tree `ae91f4192b23…`. This is the pinned benchmark main, not a claim about the current remote main |
| #439 sandbox-plan | `1b322d6c9a045f31135d02cfd13c2f140f3a1561`, tree `a9145021f5d8…`, ref `legato/i95qqzwb-sandbox-plan`. 2 commits: `ffd5e720` (tests), then `1b322d6c` (refactor) |
| #440 config-macros | `7b6de49760dfd4a71ed4110374612348e47a643a`, tree `9f13299d5822…`, ref `legato/i95qqzwb-config-macros`. 1 commit |
| #441 small-extracts | `b9696a9c71129006afb881645501d29520dbbcfd`, tree `29c8697082bd…`, ref `legato/i95qqzwb-small-extracts`. 4 commits: `1dad97dd`, `6e0d8e4c`, `bb78ad57`, `b9696a9c` |
| Evaluator | Built from `git archive 5e0c679`, tool 0.1.0-poc, `scoring/v1`. Config `legato-architecture.json` sha256 `768b2fc6…85f6`, configHash `3cef63633359769f`. Unchanged, nothing tuned |
| Object source | Offline bundle `.tmp/i95qqzwb-batches-439-440-441.bundle`. `bundle verify`: complete history, and the 3 heads equal the 3 pins. `merge-base --is-ancestor 471bfd24 <pin>` verified for each. Cloned with `--no-checkout --no-local` into `.tmp/legato-b3`, which has no worktree files and is not committed; it was unchanged after the runs. No network. The shared Legato checkout was untouched: HEAD `471bfd24`, status fingerprint `3ab8895f206d466e` before and after |
| Run | 2026-11-27T19:58:28+11:00. Each evaluate took 0.7–0.8 s and each history 0.7–1.9 s |
| Per-callable detail | The nesting and cc of each helper below were read from the evaluator's own index cache (`symbol` and `parsed_file` rows for the target blobs, adapter `go-adapter/3`). These are the same numbers the analyzers use, not a separate tool |

| Batch | Net report | History |
|---|---|---|
| #439 | `sandbox-471bfd24..1b322d6c.{json,md}` | `sandbox-history-471bfd24..1b322d6c.{json,md}` |
| #440 | `macros-471bfd24..7b6de497.{json,md}` | `macros-history-471bfd24..7b6de497.{json,md}` |
| #441 | `extracts-471bfd24..b9696a9c.{json,md}` | `extracts-history-471bfd24..b9696a9c.{json,md}` |

**#435 docs-only head.** The Reeve reports that the #435 head moved from
measured `96298449` to `3fc6d42cc43f7606f6a9afd22e3cd8c277ac75de`, with the
only difference being `site/src/content/docs/contributing/architecture.md`
and no Go change. That claim was **not re-measured** here. The #435 report
stays pinned to `96298449` and was not relabelled.

## Scores at a glance

| Analyzer | #439 sandbox | #440 macros | #441 extracts |
|---|---:|---:|---:|
| duplication | 0 | 0 | 0 |
| complexity | +1 | 0 | +1 |
| dependencies | 0 | **+1** | 0 |
| obsolete | 0 | 0 | 0 |
| scope | 0 | 0 | 0 |
| correctness | n/a | n/a | n/a |

---

## #439 BuildSandboxPlan (Claude-account and NuGet parts)

Coverage: 2 Go files (1 source, 1 test), status `ok` everywhere.

| | Before | After | Δ |
|---|---:|---:|---:|
| `BuildSandboxPlan` cc | 107 | 98 | −9 |
| `BuildSandboxPlan` decision points | 106 | 97 | −9 |
| `BuildSandboxPlan` lines | 356 | 337 | −19 |
| `BuildSandboxPlan` max nesting | 4 | 4 | 0 |
| Added callables | – | 6, cc 18 (decisions 12) | +18 |
| Deleted callables | – | 0 | 0 |
| **Changed-scope cc** (root plus helpers) | **107** | **116** | **+9** |
| Changed-scope decision points | 106 | 109 | +3 |
| Hotspots in scope | 1 | 1 | 0 (`BuildSandboxPlan` is still one); hotspots added 0 |

The 6 helpers (cc / max nesting / lines):

| Helper | cc | nesting | lines |
|---|---:|---:|---:|
| `realPaths` | 2 | 1 | 7 |
| `inProtectedClaudeAccount` | 3 | 2 | 8 |
| `usesClaudeAccount` | 2 | 0 | 3 |
| `claudeAccountStateMounts` | 4 | 2 | 13 |
| `nugetPackagesEnv` | 5 | 2 | 9 |
| `coveredInsideExemptMounts` | 2 | 0 | 3 |

The analyzer labels this outright: *"complexity removed from modified
callables (Δ-9) is matched by 6 added callables (cyclomatic 18): this is
extraction/redistribution, not net removal."* Total cc in the changed scope
rose by 9. Six of that is the baseline of 1 per new function, and 3 is
added decision points. The root function's cc agrees with the author's
gocyclo figure (107→98). That agreement is a coincidence of definitions,
not proof that the two metrics are equivalent. `BuildSandboxPlan` is still
a cc-98, 337-line hotspot at maximum nesting 4. This is a first slice, not
a resolution. The score is +1 only because `scoring/v1` credits a modified
net of −9 or less. It does not credit the redistribution.

Characterisation tests: the first-parent order is verified. `ffd5e720`
adds `sandbox_plan_claude_test.go` (+483) on top of base. It scores 0 on
every code analyzer and changes no production code. `1b322d6c` then does
the refactor. The tool sees 2 test functions and 5 test helpers. It does
not count table cases, so "21 characterisation cases" is the author's
claim and is not verified here. Duplication 0, no new duplicates. No
coupling or API change. **No correctness claim.**

---

## #440 config owns the macro model; `internal/engine/macros` deleted

Coverage: 23 files: 7 source, 11 test, 5 docs (ADR 0002, AGENTS.md,
docs/ai ×2, site architecture page). Status `ok` everywhere.

**Dependencies (+1).** This is the only batch with a coupling change. The
operator approved the removal.

| | Before | After |
|---|---:|---:|
| internal units | 86 | 85 (`engine/macros` gone) |
| internal edges | 200 | 196 (+1 / −5) |
| cycles | 0 | 0 |
| rule violations (frozen rules) | 15 | **12** |
| `config` fan-out | 1 | **0** |

- Edges removed: `config`, `internal/server`, `internal/service`,
  `internal/tui` and `internal/tui/overlay` → `internal/engine/macros`.
- Edge added: `internal/tui/overlay → config`. `overlay` already reached
  `config` indirectly. No frozen rule forbids tui → config.
- Violations resolved (3): `config → engine/macros` ("config is
  standalone") and `tui → engine/macros` and `tui/overlay → engine/macros`
  ("tui must not import engine", a prose rule). The other 12 tui → engine
  violations are pre-existing and unchanged. The operator's G decision
  ("audit intent first") still applies to them.

**Move versus deletion, and API count churn.** The scope analyzer reports
exported API +2 / −3 / ~11 and symbols +3 / −3 / ~23 (renamed 6). Almost
all of that is relocation:

| Kind | Symbols | Reading |
|---|---|---|
| Matched as moved/renamed (6) | `Macro`, `macroFields`, `Macro.UnmarshalYAML`, `Macro.UnmarshalJSON`, `macroFields.macro`, and `AnyFavorite → AnyFavoriteMacro` | Same body, new package (`config/macros.go`) |
| Deleted plus added (not matched; the body or identity changed) | `Macro.Validate` (deleted from engine, added to config), `SplitFavorites → SplitFavoriteMacros`, and `ListResult` (exported, engine) → `macroListResult` (unexported, server) | **Moves with edits**, not lost behaviour. The 2 added callables (cc 6) exactly offset the 2 deleted (cc 6). The tool cannot confirm the edits are behaviour-preserving |
| Signature changes (8) | `Server.SetMacros`, `BuiltinMacros`, `MacroService.Macros`, `MacroService.Save`, `NewApp`, `App.SetMacroSource`, `NewMacroPicker` and so on | Only the type path changed: `macros.Macro → config.Macro` |

So "−3 exported" counts one real narrowing (`ListResult` becomes
server-private) and two moves with renames or edits. Duplication: 0
introduced, 0 removed. No copy was left behind in engine:
`internal/engine/macros` has 0 files at FIX. Complexity: 16 modified
callables at cc 88→88, the 2 added/deleted pairs cancel, and hotspots stay
5→5. Those are pre-existing hotspots touched only by the type rename. That
is pure churn with no structural change. Obsolete: 0.

**Rule and version limits.**
- The frozen rule `config/** → internal/**` sees the removal of
  `config → engine/macros` (resolved). The **full** ADR 0002 rule is
  "config imports no Legato package". Here the evaluator shows `config`
  fan-out 1→0 across all module packages in its graph, so it does measure
  that at FIX. The graph excludes `_test.go` files and does no build-tag
  resolution (all variants are indexed together).
- "Engine packages must not import config" is cited in the ADR and
  enforced by Legato's `internal/engine/store/architecture_test.go`. It
  **is not in the frozen config**, so the evaluator cannot attest to it.
  A manual `git grep` at FIX finds 0 non-test engine imports of `config`.
  That is separate evidence.
- The new `TestConfigImportsNoLegatoPackage` is Legato's enforcement of
  record. The tool lists it only as an added test.
- The 3 resolved violations score +1 under `scoring/v1`. Two belong to a
  rule that Legato's tests do not enforce. The score reflects the frozen
  rule set, not a judgement on that rule.

**No correctness claim.** 11 test files changed (+122/−75), which is
context only.

---

## #441 small extracts (historyConversation, Recover, mapText, TicketCreate)

Coverage: 7 Go files (4 source, 3 test), status `ok` everywhere.

| Root (modified) | cc | nesting | lines | Helper added | cc | nesting | lines |
|---|---|---|---|---|---:|---:|---:|
| `agentService.historyConversation` | 9→6 | 2→1 | 24→14 | `latestRunAccount` | 5 | 2 | 13 |
| `SessionHistoryService.Recover` | 21→16 | 4→4 | 75→59 | `storageDirsFor` | 6 | 2 | 20 |
| `upgrader.mapText` (engine/workflow) | 13→8 | 4→3 | 31→21 | `wholeReferences` | 7 | 2 | 17 |
| `TicketCreateExecutor.Start` | 21→18 | 3→3 | 76→66 | `resolveCreateAssignee` | 5 | 1 | 17 |

| | Before | After | Δ |
|---|---:|---:|---:|
| Modified callables cc | 64 | 48 | −16 |
| Modified decision points | 60 | 44 | −16 |
| Modified lines | 206 | 160 | −46 |
| Added callables | – | 4, cc 23 (decisions 19) | +23 |
| **Changed-scope cc** | **64** | **71** | **+7** |
| Changed-scope decision points | 60 | 63 | +3 |
| Max cc in scope | 21 | 18 | −3 |
| Max nesting, modified | 4 | 4 | 0 (`Recover`) |
| Hotspots (cc ≥ 10 or nesting ≥ 4) | 3 | 2 | −1 |

- **Genuine reduction versus redistribution.** The analyzer again says
  *"extraction/redistribution, not net removal"*. The roots lost 16 cc and
  46 lines, but total cc in the changed scope rose by 7 (4 from new
  function baselines, 3 from added decisions). The gains are in shape: the
  largest function fell from 21 to 18, `mapText` dropped out of hotspot
  status (13/4 → 8/3, so pre-existing hotspots fell by 1), and
  `historyConversation` lost a nesting level. They are not in total
  logic.
- **Still hotspots:** `Recover` (cc 16, nesting 4) and
  `TicketCreateExecutor.Start` (cc 18).
- **New hotspots: 0.** The largest helper is cc 7 with nesting 2.
  **New duplicates: 0.** Structural and exact redundancy are unchanged.
- **Per commit:** each of the 4 commits scores complexity +1 and
  everything else 0. The modified cc changes are −3, −5, −5 and −3, with
  helpers of cc 5, 6, 7 and 5.
- **Tests first:** each refactor commit carries its own test (for example
  `TestHistoryConversationTakesTheNewestRunsAccount` and
  `TestAStepRewriteMapsOnlyWholeReferences`). The tool cannot tell whether
  a test was written before the code inside a single commit. Only
  #439's separate test commit makes that order visible.
- No coupling or API change. **No correctness claim.**

---

## Bottom line

- **#439:** the root function is smaller (cc 107→98), but this is
  redistribution: changed-scope cc +9 across 6 helpers, and
  `BuildSandboxPlan` remains a large hotspot. The characterisation tests
  verifiably precede the refactor.
- **#440:** a real decoupling. 4 net edges removed, `engine/macros`
  deleted with no leftover copy, 3 frozen-rule violations resolved, and
  `config` fan-out 0. The API and symbol churn is mostly moves and a type
  path change, plus one real narrowing (`ListResult`). Part of the new
  rule (engine ↛ config) is outside the frozen config and attested only
  by Legato's tests and a manual check.
- **#441:** the shape improved (one hotspot fewer, the largest function
  smaller, no new hotspots or duplicates). Total complexity in scope rose
  slightly (+7) because of extraction.

All of this is directional structural evidence. None of it establishes
correctness or that the code is clearer. CI and review are pending.
