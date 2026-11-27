# Batch delta: Legato PR #442 (TUI uses engine values, not behaviour; worklist G)

This is an independent batch from pinned main `471bfd24`, measured with the
frozen evaluator. It does **not** stack with #440 or any other batch.
#440's `engine/macros` move is absent here, so the TUI → `engine/macros`
imports are still present in this branch. Earlier reports are unchanged.

## Provenance

| | |
|---|---|
| PR | cpave3/legato#442, task i95qqzwb, ref `legato/i95qqzwb-tui-boundary` |
| Base | `471bfd24a69628d985cdd7767dcb340da6a569b7`, tree `ae91f4192b23…`. This is the pinned benchmark main, not a claim about the current remote main |
| Fix | `c05e6bdb0e338410fcc12a5b07748e2d05382a7c`, tree `1c6eab010c50…`. 2 commits: `27265986` (characterisation test), then `c05e6bdb` (refactor). `merge-base --is-ancestor` verified |
| Evaluator | Built from `git archive 5e0c679`, tool 0.1.0-poc, `scoring/v1`. Config sha256 `768b2fc6…85f6`, configHash `3cef63633359769f`. Unchanged, nothing tuned |
| Object source | Offline bundle `.tmp/i95qqzwb-tui-c05e6bdb.bundle`. `bundle verify`: complete, single head = FIX. Cloned with `--no-checkout --no-local` into `.tmp/legato-tui`, which has no worktree files and is not committed; it was unchanged after the run. No network. The shared Legato checkout was untouched: HEAD `471bfd24`, fingerprint `3ab8895f206d466e` before and after |
| Run | 2026-11-27T20:03:52+11:00. Evaluate took 0.79 s and history 0.96 s |
| Reports | `tui-471bfd24..c05e6bdb.{json,md}` and `tui-history-471bfd24..c05e6bdb.{json,md}` |

## Result

Coverage: 8 files. 5 are Go: 2 source (`internal/service/event_stream.go`,
which is new, and `internal/tui/app.go`) and 3 test. The other 3 are docs
(ADR 0003, AGENTS.md, site architecture page). Status `ok` everywhere.

| Analyzer | Score | Raw |
|---|---:|---|
| duplication | 0 | structural 260→260, exact 10→10; nothing introduced or removed |
| complexity | 0 | 3 modified callables at cc 12→12 (decisions 9→9, lines 105→103). 1 added callable, `SendPlanVerdict` (cc 1). **Changed-scope cc 12→13 (+1, the base point of the new forwarding function)**. Hotspots 0→0 |
| dependencies | **+1** | internal edges 200→199 (+0/−1): `internal/tui → internal/engine/ipc` removed. Units 86→86, cycles 0→0, frozen rule violations **15→14** (that edge was one of the 14 tui → engine violations), tui fan-out 19→18 |
| obsolete | 0 | nothing newly orphaned |
| scope | 0 | symbols +3/−0/~4; **exported API +3/−0/~1**; 1 unit group |
| correctness | n/a | not assessed. 3 test files (+193/−1), which is context only. CI and review are pending |

Per commit: `27265986` adds `app_plan_verdict_test.go` (+82). It scores 0
everywhere and changes no production code, so the test verifiably precedes
the refactor. `c05e6bdb` carries all of the measured deltas above.

## Actual edge, type and API deltas

**Edges.** Only one edge changed: `tui → engine/ipc` was removed. No edge was
added. `internal/service` already imported both `engine/events` and
`engine/ipc`, so moving the call into service adds no service edge. The IPC
dependency moved inward, from tui to service, rather than disappearing
from the program.

**New exported service surface** (`internal/service/event_stream.go`, new file):

| Symbol | Kind | What it is |
|---|---|---|
| `EventStream` | interface (`Subscribe`, `Publish`) | Abstracts `*events.Bus`. It still exposes `events.EventType` and `events.Event` in its method signatures, so it abstracts the bus's *behaviour*, not its *types*. `var _ EventStream = (*events.Bus)(nil)` asserts conformance at compile time |
| `PlanVerdict` | struct (4 string fields) | The parameter object for a verdict. It mirrors a subset of `ipc.Message` |
| `SendPlanVerdict` | function, cc 1 | A pure forwarder: it builds an `ipc.Message{Type: "plan_verdict", …}` and calls `ipc.Send`. It has no logic of its own |

**Changed exported signature.** In `NewApp`, the parameter `bus *events.Bus`
became `bus service.EventStream`. `App.eventBus` changed to match.
`handlePlanApprove` and `handlePlanReject` keep cc 4. The inline
`ipc.Message` literal moved behind `SendPlanVerdict` with the same fields.
Status strings stay at the call site.

**Shifted complexity.** Essentially none: +1 cc for the forwarder's base
point. Nothing moved out of the handlers except a struct literal.

## Rule and version mismatch (reported, not tuned)

The frozen config encodes the CLAUDE.md prose rule literally at **import
level**: `internal/tui/** → internal/engine/**` is forbidden. ADR 0003
replaces that with a **symbol-level** rule: the TUI may name audited engine
value types, constants and pure helpers (the `engineValues` allowlist), and
engine *behaviour* reaches it only through service interfaces. The two do
not line up:

1. **13 of the 14 frozen tui → engine "violations" remain**, and under ADR
   0003 they may be **intentional**: value-type imports such as `events`
   (types and constants), `store` records, `workflow`, `github`, `git`,
   `analytics` and `macros`. The frozen evaluator still flags all 13 as
   pre-existing violations. It cannot tell an allowed value import from a
   forbidden behaviour import. The remaining count of 14 is therefore
   **not** a count of ADR 0003 violations.
2. **Symbol-level purity is not attested.** That "the TUI uses only
   `Event`, `EventType`, `Event…` constants and `…Payload` from `events`,
   never `Bus` or `New`", and the whole `engineValues` list, are checked
   only by Legato's new `TestTUIUsesOnlyEngineValueTypesAndPureHelpers`.
   The evaluator lists it as an added test and runs nothing.
3. **Interface boundaries are not attested.** The evaluator sees that the
   type of `NewApp` changed. It does not verify that no TUI code reaches the
   concrete bus another way (type assertion, reflection), that
   `EventStream`'s behaviour matches the bus, or that the ADR's
   "no allowed type has methods doing I/O" audit is true. The ADR's own
   Limits section says the symbol test also cannot see method calls on
   values of allowed types.
4. **Frozen scoring rewards the edge removal (+1)** because a frozen-rule
   violation was resolved. That reflects the old import-level rule. Fewer
   imports is not the same as a correct boundary.

Expressing ADR 0003 would need a symbol-level rule type, which the
evaluator does not have, plus a new config version. Neither was applied.

## Judgement items the tool does not measure

- **Nil-interface semantics.** The new `NewApp` comment says *"a nil bus
  (the interface itself, not a nil `*events.Bus` inside it) subscribes to
  nothing"*. A caller passing a typed-nil `*events.Bus` would now get a
  non-nil interface and pass the `bus != nil` and `a.eventBus != nil`
  guards (`app.go:421`, `:2487`, `:2513`). By manual read, the one
  production caller (`cmd/legato/main.go:1649`) passes `bus :=
  events.New()`, which is not nil. This is a latent hazard for future
  callers or tests, not a demonstrated defect.
- **API growth for a forwarder.** There are 3 new exported service symbols
  to remove 1 import. `PlanVerdict` duplicates a subset of `ipc.Message`'s
  fields. Whether that indirection pays off is a design judgement.
- **Correctness.** `TestPlanVerdictsReachTheConductorAndTheBus` (committed
  first) is context, not evaluator evidence.

## Bottom line

Measured under the frozen import-level rules: one forbidden edge
(`tui → engine/ipc`) was removed, so frozen violations went from 15 to 14.
No edges were added, no cycles formed, complexity is flat (+1 for a cc-1
forwarder) and there are no duplicates. The cost is 3 new exported service
symbols and one changed constructor signature. The IPC dependency moved
from tui into service; it did not go away. The new symbol-level boundary
(ADR 0003) and the intent behind the 13 remaining value imports are
**outside what the frozen evaluator can attest**. They rest on Legato's
architecture test and on review. This is directional structural evidence,
not proof that the boundary is correct.
