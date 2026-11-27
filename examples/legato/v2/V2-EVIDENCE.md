# scoring/v2 (experimental) on Legato

This note shows the experimental v2 model on the same Legato revisions as
the frozen v1 evaluation. It is a **separate, new measurement**. The v1
reports, the frozen `5e0c679` build and the original results are unchanged
and remain the provenance of record. v2 was designed after seeing the v1
results. Its weights are uncalibrated defaults. They were **not** tuned to
favour the cleanup.

## Setup

| | |
|---|---|
| Tool | codeeval 0.1.0-poc with `score` (built from this repo's working tree; `scoring/v2-experimental`) |
| Config | `legato-architecture.json`, configHash `3cef63633359769f`: the same rules and thresholds as v1 |
| Weights | cyclomatic excess ×1, nesting excess ×2, 20 redundant tokens = 1 point, rule-violating edge ×10, unit in cycle ×20; score half-point at 25 debt/kLOC |
| Objects | The existing private `--no-checkout` clones `.tmp/legato-stack` (stack bundle, contains `9a61e9f8`…`bb13a424`) and `.tmp/legato-fix` (#434). No network; the shared Legato checkout was untouched |
| Run | 2026-11-27T21:00+11:00. A snapshot takes about 0.3 s warm; the 40-step stack delta about 13 s |
| v1 regression check | `evaluate` from the same binary reproduces `final-stack-471bfd24..bb13a424.json` byte for byte, and `fix-delta-471bfd24..41c80614.json` except for its `--timing` block |

Reports in this folder: `snapshot-<rev>.{json,md}`,
`delta-window-9a61e9f8..471bfd24`, `delta-stack-471bfd24..bb13a424`,
`delta-drift-9a61e9f8..bb13a424`, `delta-434-471bfd24..41c80614` (each
`.json` and `.md`), and `layers/` (one delta per stack layer).

## Absolute snapshots (Go source only: 61.8% of production source bytes)

| Revision | Debt | complexity | nesting | duplication | architecture | Source lines | Debt/kLOC | Score |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| `9a61e9f8` (before the window) | 10,783.30 | 9,275 | 598 | 760.30 | 150 | 224,082 | 48.12 | 34.19 |
| `471bfd24` (after the window) | 11,611.30 | 10,027 | 618 | 816.30 | 150 | 238,533 | 48.68 | 33.93 |
| `41c80614` (#434 only) | 11,601.65 | 10,027 | 618 | 806.65 | 150 | 238,519 | 48.64 | 33.95 |
| `bb13a424` (cleanup stack tip) | 11,387.20 | 9,888 | 612 | 777.20 | 110 | 238,865 | 47.67 | 34.40 |

Not measured: TypeScript 577 files, SQL 160, Kotlin 54, Rust 34,
JavaScript 23, shell 15, Python 5, Swift 2. Each snapshot report lists the
modules by debt, the largest callable debts (for example `App.update` with
cc 308 is 305 points and `SpawnAgent` with cc 287 is 282 points) and the
largest clone groups.

## Deltas: v1 versus v2

| Comparison | v1 scores (dup / cx / dep) | v2 Δ debt | complexity | nesting | duplication | architecture | Decision points / base points | Gates failed |
|---|---|---:|---:|---:|---:|---:|---|---|
| Two-day window `9a61→471b` | −2 / −2 / 0 | **+828.00** | +752 | +20 | +56.00 | 0 | +2,586 / +564 | exact dups, structural dups, branch growth, hotspot excess |
| Cleanup stack `471b→bb13` | +2 / 0 / +1 | **−224.10** | −139 | −6 | −39.10 | −40 | −24 / +38 | none |
| #434 `471b→41c8` | +1 / +1 / 0 | **−9.65** | 0 | 0 | −9.65 | 0 | −6 / 0 | none |
| Window and stack `9a61→bb13` | −2 / −2 / +1 | **+603.90** | +613 | +14 | +16.90 | −40 | +2,562 / +602 | structural dups, branch growth, hotspot excess |

What v2 shows that v1 could not:
- **Magnitudes on one scale.** The window added 828 debt points and the
  stack removed 224, so the stack recovered **27%** of what the window
  added. Under v1 both were capped at −2 and +2.
- **Additivity.** Window plus stack is +828 − 224 = +604, and the drift
  delta measures exactly +603.90. The 11 per-layer deltas sum to −224.1,
  which equals the stack delta. In both step series (23 window merges and
  40 stack commits) Σ step Δdebt equals the net Δdebt (`stepCheck.telescopes:
  true`).
- **Refactors get credit without being called removals.** The stack's
  complexity debt fell by 139, because the large functions shrank and only
  a few new helpers cross the hotspot threshold. The decision-point ledger
  (−24 decisions, +38 base points) still shows that most of the logic moved
  rather than disappeared.

Per-layer v2 deltas (manifest parent → head):

| Layer | Δ debt | complexity | nesting | duplication | architecture | Decisions / base |
|---|---:|---:|---:|---:|---:|---|
| #436 ledger | 0.00 | 0 | 0 | 0 | 0 | 0 / 0 |
| #435 leaf | −9.00 | 0 | 0 | −9.00 | 0 | −8 / −2 |
| #437 service | −30.10 | 0 | 0 | −30.10 | 0 | −19 / +4 |
| #440 macros | −30.00 | 0 | 0 | 0 | −30 | 0 / 0 |
| #442 TUI | −10.00 | 0 | 0 | 0 | −10 | 0 / +1 |
| #441 small | −23.00 | −19 | −4 | 0 | 0 | +4 / +5 |
| #443 review | −72.00 | −70 | −2 | 0 | 0 | −7 / +16 |
| #446 resume | −17.00 | −17 | 0 | 0 | 0 | +2 / +4 |
| #445 accounts | −14.00 | −14 | 0 | 0 | 0 | +3 / +3 |
| #439 sandbox plan | −11.00 | −11 | 0 | 0 | 0 | +1 / +6 |
| #444 sandbox prefix | −8.00 | −8 | 0 | 0 | 0 | 0 / +1 |
| **Σ = stack** | **−224.10** | −139 | −6 | −39.10 | −40 | −24 / +38 |

## Important finding: the normalised score can be diluted

Between `9a61e9f8` and `bb13a424`, debt **rose by 603.90**, yet the
codebase score **rose** from 34.19 to 34.40. About 14,800 new lines of
lower-density feature code pulled the average down. The same happens at
module level: in the window, `internal/engine/workflow` gained 97 debt
points while its score went *up* from 27.47 to 28.78. Any per-line score
has this property. So:

- **Use debt (total and by module) as the metric to optimise.** Adding code
  can never lower it, and splitting it by module and component shows where
  it lives.
- Use the 0–100 score only to **compare** modules or codebases of different
  sizes, and never as a target on its own.
- For changes, use **Δdebt and debt per changed kLOC** (the window was
  +49.9 per changed kLOC, the stack −132.0) plus the advisory gates.

Also note how little the score moves. At 238k measured lines, the whole
cleanup moved it by 0.47 points. The codebase score is a slow-moving index,
while module scores and debt respond to individual changes.

## Limits

- **Weights are uncalibrated.** Complexity dominates total debt (about 87%
  at `bb13a424`), so the total mostly tracks hotspot complexity. Any
  composite inherits its weights, so the component breakdown is reported
  alongside it.
- **Hotspot excess rewards splitting.** Breaking a large function into
  helpers below cc 10 removes its debt even if no logic is removed. The
  decision-point ledger and the existing-callable branch gate are there to
  expose that.
- Complexity growth below the thresholds carries no debt.
- Architecture debt covers only the configured import-level rules. As with
  v1, the 11 remaining `tui → engine` edges may be permitted under ADR 3.
  Stdlib-only leaves, engine ↛ config and symbol-level rules are not
  measured.
- **Go only** (61.8% of production source bytes). Names are not resolved,
  and nothing is executed, so there is no correctness claim.
- Calibrated bands (percentiles of this repository's own step rates) are
  not implemented yet. The 23 window merges and 40 stack commits here are a
  start for that data.
