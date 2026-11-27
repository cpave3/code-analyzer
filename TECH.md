# TECH: dependencies, design and trade-offs

## Dependencies (pinned in go.mod)

| Purpose | Choice | Why / trade-off |
|---|---|---|
| Language | Go 1.26 | Required by the PRD. Single static CLI. |
| Parsing | `github.com/tree-sitter/go-tree-sitter` v0.25.0 (official bindings) + `github.com/tree-sitter/tree-sitter-go` v0.25.0 | Required by the PRD. These are the official, maintained bindings, and grammars for TypeScript and C# exist for the same runtime. Trade-off: cgo, so a C compiler is needed to build. There was no toolchain blocker, so the prescribed parser was kept. `go/ast` would have been simpler for Go alone but would not generalise to TS/C#. |
| Index | `modernc.org/sqlite` v1.63.0 | Pure-Go SQLite, so no second cgo dependency. Stores parse results and a queryable `symbol` table keyed by (adapter, adapter version, blob id). Unchanged blobs are never reparsed across snapshots, commits or runs. |
| Git | `git` CLI | Required by the PRD. Uses plumbing only: `rev-parse`, `ls-tree -z`, `cat-file --batch`, `diff --name-status -M -z`, `rev-list --first-parent`, `ls-files -z`. |

Toolchain used for development and verification: go1.26.0 installed by
mise at `/home/cameron/.local/share/mise/installs/go/1.26.0/bin`, with
`/usr/bin/gcc` 16.1.1 and `/usr/bin/git` 2.55.0, on the Bazzite 44 host (no
container). Go is not on a minimal system `PATH`. Prefix `PATH` with that
directory and set `GOTOOLCHAIN=local GOPROXY=off GOFLAGS=-mod=readonly` for
an offline, reproducible build (see README → Build environment).

There are no other third-party dependencies. Line diffs (bounded Myers),
glob matching, MinHash and Tarjan SCC are small in-process implementations.

## Pipeline

```
resolve revs ─▶ index base snapshot ─┐
               index target snapshot ┴▶ change set (files, line ranges, symbol mapping)
                                     ─▶ analyzers (engine computes coverage per analyzer)
                                     ─▶ report (JSON schema codeeval.report/v1 + Markdown)
```

- **Snapshot** (`internal/index`): every tracked file is classified by role
  and language. Supported files are parsed through the cache, in parallel,
  and symbols get identity keys (`lang:unit|package|Container.Name`). Graphs
  are built per adapter. A snapshot is immutable, and the evaluator keeps a
  three-entry LRU so a history run reuses each commit's snapshot for the next
  step.
- **Change set** (`internal/change`): name-status with rename detection, and
  in-process line diffs for line counts and changed ranges. Symbols are
  matched by key across all changed files. Unmatched symbols with an
  identical exact fingerprint (which excludes the name) are paired as
  `renamed`. The other statuses are `modified`, `moved` (same identity and
  body, different file), `formatting_only` (touched but token-identical),
  `added` and `deleted`.
- **Analyzers** (`internal/analyzers/*`) are independent and stateless. They
  depend only on `analyzer`, `model`, `lang` and the snapshots.
- **Engine** (`internal/engine`): coverage from the descriptor, the
  unsupported/partial rules, panic recovery, deterministic finding order and
  report assembly. `history` evaluates the net change, then parent^1→commit
  for each first-parent commit.

## Key decisions

- **Scores are per analyzer, with no rollup**, even within a dimension
  (maintainability has two analyzers). This follows the PRD's "no single
  opaque score".
- **Directional scores come from transparent, absolute thresholds**
  (`scoring/v1`). Each result carries the rule string that fired. The
  absolute thresholds saturate on large windows; a size-normalised variant
  is deferred until there is data to calibrate it.
- **New functionality is not degradation by itself.** In complexity, -2 is
  reserved for complexity added to pre-existing callables. New hotspot
  callables cap at -1. In dependencies, new edges are info, and only cycles
  and configured rule violations score. Scope never scores positive and only
  penalises unrelated formatting or rename churn.
- **Duplication is attributed by delta**, by comparing the size of each
  repository-wide clone group before and after. Pre-existing duplication,
  renames and moves leave group sizes unchanged and are not counted.
  Consolidation shrinks groups and is credited.
- **Structural similarity is a candidate signal.** Exact-structure matches
  are measured facts, but their meaning is labelled with a caveat.
  Near-structural MinHash matches are heuristics and never scored.
- **Unsupported is not neutral.** The engine sets status `unsupported` and a
  null score when no changed code in an analyzer's scope could be analysed,
  and `partial` (with per-language counts) when some could not.
- **Deterministic output.** No timestamps unless `--timing` is passed, sorted
  findings, sorted maps, and deterministic Git identities in fixtures. A test
  asserts byte-identical JSON across cold-cache, warm-cache and no-cache runs.
- **Safety.** See the README's read-only guarantees. Configuration is never
  read from the evaluated repository.

## Known limitations

- Go references are matched by name. `obsolete` and `callers` can miss
  reflection, linkname, interface satisfaction and build tags. They are
  labelled heuristics.
- Go import cycles do not compile, so cycle fixtures are intentionally
  non-building. The detector is language-neutral for future TS/C# graphs.
- Locality uses a two-hop neighbourhood in the import graph only. Runtime
  wiring, config and tests are invisible to it.
- Lines moved within a file are not detected, only symbol moves.
- Very large diffs (more than 2000 line edits in a file) fall back to treating
  the whole file as changed (`approximateDiff: true` on that file).
- Per-step history reports are not additive. Repository-wide measurements
  interact across steps.
