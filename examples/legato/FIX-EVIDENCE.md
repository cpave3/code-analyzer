# Fix delta: Legato PR #434 (duplication cleanup)

This measures the cleanup that followed the findings in [EVIDENCE.md](EVIDENCE.md).
The evaluator was frozen and nothing was tuned for this patch.

## Provenance

| | |
|---|---|
| PR | cpave3/legato#434, Claude Work task ah7qmrsa (branch `legato/ah7qmrsa-consolidate-duplication`) |
| Base | `471bfd24a69628d985cdd7767dcb340da6a569b7`, tree `ae91f4192b23…` (the same tree as the original benchmark target) |
| Fix | `41c8061476b6592b5ea4faa192d68a2a803b6c59`, tree `555c4b3e8792…`, 3 linear commits on top of base |
| Comparisons | Main had not moved, so the cleanup-only delta (base→fix) and the 471bfd24→fix delta are the **same comparison** |
| Evaluator | code-evaluator `5e0c679`, tool 0.1.0-poc, `scoring/v1`, config `legato-architecture.json`, configHash `3cef63633359769f`. Thresholds, exclusions and scorer unchanged |
| Object source | Offline bundle `.tmp/ah7qmrsa-fix.bundle` provided by the Reeve (verified complete, head = FIX_SHA). It was cloned locally into `.tmp/legato-fix.git` (bare) and then `.tmp/legato-fix` (`--no-checkout`). Neither is committed. No network, and the shared Legato checkout was not touched (status fingerprint `3ab8895f206d466e…` and HEAD unchanged). |
| Run | 2026-11-27T17:35:16+11:00: net in 1.6 s, history in about 1 s |

**Input-handling mismatch:** the frozen evaluator cannot open a **bare**
repository (`rev-parse --show-toplevel` fails), so the first attempt exited
with "not a git repository". Rather than change the evaluator, I evaluated a
`--no-checkout` clone of the same objects. No files were checked out, and
commit mode only reads objects. Refs and files in the clone were unchanged by
the run. The report's `repository` field therefore reads `legato-fix`.
Supporting bare repos is a small follow-up to input handling; scoring is
unaffected.

Rerun:
```bash
export PATH=/home/cameron/.local/share/mise/installs/go/1.26.0/bin:$PATH GOTOOLCHAIN=local GOPROXY=off GOFLAGS=-mod=readonly
git -c core.hooksPath=/dev/null clone -q --bare --no-local .tmp/ah7qmrsa-fix.bundle .tmp/legato-fix.git
git -c core.hooksPath=/dev/null clone -q --no-checkout --no-local .tmp/legato-fix.git .tmp/legato-fix
go build -o bin/codeeval ./cmd/codeeval
bin/codeeval evaluate --repo .tmp/legato-fix --base 471bfd24a69628d985cdd7767dcb340da6a569b7 --target 41c8061476b6592b5ea4faa192d68a2a803b6c59 \
  --config examples/legato/legato-architecture.json --output examples/legato/fix-delta-471bfd24..41c80614.json --summary examples/legato/fix-delta-471bfd24..41c80614.md
bin/codeeval history  --repo .tmp/legato-fix --base 471bfd24a69628d985cdd7767dcb340da6a569b7 --target 41c8061476b6592b5ea4faa192d68a2a803b6c59 \
  --config examples/legato/legato-architecture.json --output examples/legato/fix-history-471bfd24..41c80614.json --summary examples/legato/fix-history-471bfd24..41c80614.md
```

## Result (base → fix)

Coverage: 11 files changed, all Go (7 source, 4 test). Every analyzer status
is `ok`, with no unsupported files and no diagnostics.

| Analyzer | Score | Raw |
|---|---:|---|
| duplication | **+1** | structural redundancy 260→258 (−2); exact 10→10; 2 groups removed, 0 introduced, 0 near-structural candidates |
| complexity | **+1** | modified callables net cc 37→29 (−8). Changed-scope total 41→35 (−6): a **real net reduction**, not just redistribution. 2 added (cc 6), 2 deleted (cc 4). Hotspots 1→1 |
| dependencies | 0 | edges 200→200, cycles 0→0, rule violations 15→15 (all pre-existing, none new) |
| obsolete | 0 | nothing newly orphaned |
| scope | 0 | 3 packages in one connected group; exported API unchanged |
| correctness | n/a | not assessed. 4 test files changed (+129/−3), which is context only |

Per commit (first-parent):
1. `6f4e4616` server broadcast loop: duplication +1, complexity +1 (redundancy −1, modified cc −6)
2. `b8ebb812` jira doJSON: duplication +1, complexity 0 (redundancy −1; doPost and doPut deleted, cc 4; doJSON added, cc 2)
3. `41c80614` plan uses pathWithin: duplication 0, complexity 0 (modified cc −2, below the ±3 threshold)

## Which groups disappeared

| Group | Before (base `471bfd24`) | After (fix `41c80614`) |
|---|---|---|
| Server invalidation events | `Server.StartSubscriptionUsageEvents` `internal/server/subscription_usage.go:70-87` (cc 4, nesting 3) and `Server.StartServerThemeEvents` `internal/server/theme_preference.go:23-40` (cc 4, nesting 3) | Each is a 3-line delegate (`subscription_usage.go:69-71`, `theme_preference.go:22-24`, cc 1) to the new `Server.broadcastInvalidations` `internal/server/invalidation_events.go:12-29` (cc 4). Group 2→0 |
| Jira JSON writes | `Client.doPost` `internal/engine/jira/client.go:157-164` and `Client.doPut` `internal/engine/jira/users.go:57-64` | Both deleted. The new `Client.doJSON` `client.go:158-165` takes the method as a parameter. All four callers were updated: `DoTransition`, `Assign`, `CreateIssue`, `AddComment`. Group 2→0 |

## What the tool did not see, or that needs judgement

- **The `plan.go` → `pathWithin` change is invisible to the duplication
  analyzer.** It replaced an inline closure in `planBundleImage`
  (`internal/service/plan.go:72-97` → `72-93`) with the existing service
  helper. Closures are not indexed symbols, so this showed up only as cc 8→6
  in that function, below the scoring threshold. It does not affect the
  separately flagged `accountPathContains` (config) ↔ `pathWithin` group,
  which persists by design, so no metric credits this commit.
- **The intentionally kept groups are still present and still count:**
  `retryAfter` and `openRegular` (exact: exact redundancy stays at 10),
  `CredentialsPath`/`SessionsDir`, the `*UsageError` pair,
  `loadSubscriptionUsage`/`loadUsage` and
  `accountPathContains`/`pathWithin`. None of these symbols changed. The
  author's reason is layering: there is no permitted shared home under the
  engine-sibling and config-standalone rules. CodeEval cannot judge that
  trade-off. It only shows the cost still present.
- **New wrappers:** the two `Start*Events` methods are now cc-1 delegates
  that keep the public entry points. They fall under `minTokens`, so they do
  not count as duplication, and the tool takes no view on whether they add
  value. `doJSON` adds a method parameter (string), a small loss of
  type-level specificity that the tool does not measure.
- **No coupling change:** no imports were added or removed, so the
  dependency picture is identical.
- **Correctness is not established here.** The author reports focused and
  architecture tests green, with full CI and review pending. The new tests
  (`TestWritesSendJSONWithBasicAuth`,
  `TestAssignRateLimitedIsRetriedWithItsBody`,
  `TestInvalidationBroadcastersStopCleanly`,
  `TestPlanImageInASiblingFolderSharingTheBundlesPrefixIsOutside`) are listed
  as context, not evidence.

## Bottom line

The frozen measurements match what the fix claims to do: two clone groups
removed (redundancy −2), a real net drop in complexity (changed-scope cc −6)
with no new hotspots, no coupling change, no orphans, and no API change.
The third commit (closure → `pathWithin`) is a real reuse that this
evaluator under-counts. The deliberately kept cross-package copies remain
measurable debt. The +1 scores reflect a small, focused cleanup; they say
nothing about behaviour.
