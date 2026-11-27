# CodeEval (POC)

A local Go CLI that measures whether a Git change improves or degrades a
codebase **in the context of that repository**: duplication against existing
code, complexity before and after, dependency structure, scope, and code left
obsolete. It judges the patch itself, not what the agent or workflow says
about it. See [PRD.md](PRD.md) for the full product brief and
[TECH.md](TECH.md) for dependencies and trade-offs.

This is a **proof of concept**. It measures structure. It does not prove a
design is good, and it never combines dimensions into an overall score.

## Build and run

Requires Go 1.26+, git and a C compiler (Tree-sitter uses cgo).

### Build environment (where this POC was built and tested)

| Tool | Location on the dev host | Notes |
|---|---|---|
| Go | `/home/cameron/.local/share/mise/installs/go/1.26.0/bin/go` (go1.26.0 linux/amd64) | Installed by **mise**, per user. It is on `PATH` only in shells where mise is activated, such as Cameron's normal terminal. It is not in `/usr/bin`. |
| C compiler | `/usr/bin/gcc` (GCC 16.1.1) | Needed for cgo (Tree-sitter) |
| git | `/usr/bin/git` (2.55.0) | |
| OS | Bazzite 44 (Fedora Kinoite), on the host directly | No container or toolbox needed |
| Go modules | `~/go/pkg/mod` (already populated) | The build works offline with `GOPROXY=off` |

**`go: command not found` (exit 127) is a toolchain `PATH` setup issue, not a code failure.**
A minimal `PATH` (`/usr/bin:/bin`) has no `go`, so the command fails before
any code is compiled. Put the existing toolchain directory first; nothing
needs installing globally:

```bash
cd /home/cameron/Projects/code-evaluator
export PATH=/home/cameron/.local/share/mise/installs/go/1.26.0/bin:$PATH
export GOTOOLCHAIN=local GOPROXY=off GOFLAGS=-mod=readonly   # offline, pinned toolchain, go.mod not modified
go version                                   # go1.26.0 linux/amd64
go build -o bin/codeeval ./cmd/codeeval
go vet ./... && go test -count=1 ./...
```

This exact sequence was verified under `env -i HOME=$HOME PATH=<go dir>:/usr/bin:/bin`:
build, vet and every test passed. The e2e tests themselves call `go` (to
build the binary and to run `go test` in fixture variants), so the Go
directory must be on `PATH`, not just invoked by absolute path. Without
mise, any Go ≥ 1.26 plus gcc works. Point `PATH` at its `bin` instead.
The scripts in `scripts/` need the same `PATH` and stop with a clear message
if `go` is missing.

```bash
go build -o bin/codeeval ./cmd/codeeval

# Compare two revisions (JSON + Markdown summary)
bin/codeeval evaluate --repo ../some-repo --base main --target HEAD \
  --output /tmp/eval.json --summary /tmp/eval.md

# Tracked working-tree changes against a base (untracked files are never read)
bin/codeeval evaluate --repo ../some-repo --base HEAD --working-tree

# Net change plus one report per first-parent commit or merge in a window
bin/codeeval history --repo ../some-repo --base <sha> --target <sha> \
  --config my-rules.json --output /tmp/history.json --summary /tmp/history.md --timing

bin/codeeval languages   # adapters, capabilities, limits, unsupported languages
bin/codeeval analyzers   # analyzers, dimensions, required capabilities
```

Flags: `--config FILE` (JSON, `configVersion: 1`), `--format md|json|none`
(stdout), `--cache-dir DIR` (SQLite index, default `$XDG_CACHE_HOME/codeeval`),
`--no-cache`, `--timing`. Output is deterministic unless you pass `--timing`.

### Read-only guarantees

The evaluated repository is only read through Git plumbing (`ls-tree`,
`cat-file --batch`, `diff --name-status -z`). The tool never checks out,
resets, stashes, fetches, creates worktrees or runs the repository's code,
scripts or hooks. It also:

- runs git with `--no-optional-locks`, `core.fsmonitor=false`,
  `core.hooksPath=/dev/null` and `protocol.allow=never`, without inherited
  `GIT_*` environment variables;
- disables external diff and textconv drivers and computes line diffs in
  process;
- reads NUL-separated paths, so unusual filenames are handled verbatim;
- refuses `--output`, `--summary` or `--cache-dir` paths inside the evaluated
  repository;
- treats repository content, including commit subjects, as data.

The e2e tests fingerprint refs, status, stash, worktrees and every file
(content and mtime) before and after each run.

## What it measures

| Dimension | Analyzer | Measures | Scored on |
|---|---|---|---|
| duplication | `duplication` | Exact clones (same token text) and structural clones (identifiers and literals abstracted) of callables with ≥40 tokens, compared across the whole repository before and after. Near-structural reuse candidates come from MinHash. | Net change in structural redundancy. Pre-existing duplication is not counted, and neither are renames or moves. Near-structural candidates are never scored. |
| complexity | `complexity` | Per-callable cyclomatic complexity, decision points, max nesting, lines and statements for modified, added and deleted callables. Hotspot is cc ≥ 10 or nesting ≥ 4. | Net cyclomatic delta of modified (pre-existing) callables, plus hotspot transitions in existing code. New hotspot callables cap at -1. Moving complexity into new helpers is labelled as redistribution, not removal. |
| architecture_coupling | `dependencies` | Internal package edges added or removed, fan-in and fan-out of touched units, cycles (Tarjan SCC), and violations of configured forbidden-dependency rules, both new and pre-existing. | Cycles and new rule violations. New edges alone are reported but not judged. |
| maintainability_scope | `scope` | Files, lines and symbols by class and status. Exported API changes. Callers affected (matched by name). Locality of the changed units in the dependency graph. | Only formatting- or rename-only churn in units unrelated to the main change is scored. A focused change scores 0. This analyzer never scores positive. |
| maintainability_scope | `obsolete` | Non-exported functions that were referenced before the change, have no production references after it, and were left in place. Also credits deleting functions that were already unreferenced. | Heuristic and name-based. Methods are low confidence and not scored. |
| correctness_confidence | `correctness` | Always `not_assessed`. Test churn is shown as context only. | Never scored. |

Scores use the PRD scale -2…+2. `null` means not scored (unsupported, error
or not assessed). Unsupported is **not** neutral: if none of the changed code
in an analyzer's scope can be analysed, its score is `null` and its status is
`unsupported`. If only part of it can be analysed, the status is `partial` and
the coverage is reported. Every score carries the rule that produced it, the
raw before/after measurements and evidence. Each piece of evidence gives the
side (base or target), the pinned commit SHA, the file, the symbol and the
line range. Removed code is cited at the base revision. Findings are labelled
`measured` or `heuristic`, and heuristics also carry a confidence label and a
caveat. File classes (source, test, testdata, generated, vendor, build, docs,
excluded, other) are reported. The duplication, complexity and obsolete
analyzers measure production source only.

## Languages

| Language | Status |
|---|---|
| Go | Supported: Tree-sitter adapter with symbols, complexity, fingerprints, package import graph (via the nearest `go.mod`) and name-based references |
| TypeScript, C# | Recognised and counted as unsupported coverage. Next adapters (see below) |
| JS, Kotlin, Rust, Python, Swift, Java, SQL, shell | Recognised and counted as unsupported coverage |

The Go adapter does not type-check, resolve interface dispatch or build a call
graph, and it indexes all build-tag variants together. These limits are also
written into every report.

## Extending

The CLI has no `switch` over languages or metrics. It uses two registries,
each wired in one file. The registration contract is tested in
`internal/engine/contract_test.go`: a toy language and a toy analyzer are
added through the public registries without touching existing code.

### Add a language (e.g. TypeScript, C#)

1. Create `internal/lang/<name>/` and implement `lang.Adapter`
   (`internal/lang/lang.go`):
   - `ID`, `Version` (part of the parse-cache key) and `Extensions`
   - `Capabilities()`: declare only what you really extract, with a note on
     fidelity (`symbols`, `complexity`, `fingerprints`, `dependencies`,
     `references`)
   - `Limits()`: what you deliberately do not resolve
   - `IsTestPath`, `PackageOf`: the unit used for symbol identity and graph
     nodes (a directory for Go; a module file or a project for TS/C#)
   - `Parse(src) → model.ParsedFile`: symbols with ranges, metrics,
     fingerprints (feed tokens into `lang/fingerprint.Builder`), imports and
     refs. Report syntax errors in `ParseErrors`; don't return an error for
     them
   - `BuildGraph(GraphInput)`: resolve only what you can. TS would read
     `tsconfig.json` paths; C# would read `.csproj` `ProjectReference` and
     namespaces
2. Register it in `internal/languages/languages.go`. Its `Register` call
   supersedes the existing `Declare` entry for the same extensions.
3. Add a fixture scenario under `testdata/fixtures/` and assertions in
   `e2e/`.

Analyzers need no change. If your adapter lacks a capability, the analyzers
that need it report your files as unsupported coverage.

### Add an analyzer

1. Implement `analyzer.Analyzer` (`internal/analyzer/analyzer.go`). The
   descriptor declares the ID, dimension, version, required capabilities and
   the file classes it measures.
2. Read `ctx.Base`, `ctx.Target` (indexed snapshots) and `ctx.Change` (files
   and mapped symbols). Return measurements, findings with evidence, and a
   score with its rule.
3. Register it in `internal/analyzers/analyzers.go`.

The engine computes coverage from the descriptor, sets
`partial`/`unsupported` status, recovers panics as `error` results, sorts
findings deterministically and adds the result to the report schema.

### Configuration

```json
{
  "configVersion": 1,
  "exclude": ["gen/**"],
  "duplication": {"minTokens": 40, "nearThreshold": 0.85, "includeTests": false, "degradeAt": 1, "significantDegradeAt": 3},
  "complexity": {"hotspotCyclomatic": 10, "hotspotNesting": 4, "degradeAt": 3, "significantDegradeAt": 10},
  "architecture": {"rules": [{"name": "store must not depend on presentation", "from": "internal/store/**", "to": "internal/cli/**"}]},
  "scope": {"formattingOnlyChurnAt": 3}
}
```

Each report records `configVersion`, a hash of the effective config and the
`scoringVersion` (`scoring/v1`).

## Proof of value

### Fixtures

`testdata/fixtures/` holds paired scenarios. Both variants of a pair
implement the same behaviour and pass the same tests (the e2e suite runs
`go test` in each variant), but they differ in quality.
`examples/fixtures/README.md` has the score table, with a JSON and Markdown
report for every pair.

| Pair | Better variant | Worse variant |
|---|---|---|
| reuse | reuses `text.Slugify` → duplication 0 | private structural copy → duplication **-1**, evidence cites both |
| consolidate | deletes both copies → duplication **+1** | leaves both copies → duplication **-1**, obsolete **-1** |
| complexity | extracts helpers → complexity **+1**, labelled redistribution (changed-scope total rises) | nested conditionals → complexity **-1** (Discount cc 4→10, new hotspot) |
| layering (with config) | local formatting → architecture 0 | imports the presentation package → **-1**, cites the import line |
| cycle / cycle-removal | decouple **+1**, remove cycle **+2** | introduce cycle **-2** (does not compile in Go, so it only exercises the detector) |
| scope | focused fix → 0 | the same fix plus formatting churn in unrelated packages → **-1** |

Edge cases: unchanged baseline (everything 0), rename-only (symbols marked
"moved", all 0), symbol rename (renamed, not penalised), added and deleted
symbols (with base locations), parse error (`parse_error` diagnostic,
`partial`), TypeScript-only (`unsupported`, score null), mixed Go and TS
(`partial`), and an odd filename (`util/odd name ü.go`).

```bash
go test ./...                    # unit, contract and e2e (real CLI subprocess) tests
scripts/fixture-reports.sh       # regenerate examples/fixtures
```

### Real repository: Legato

`examples/legato/` holds the headline evaluation of a pinned Legato window,
`9a61e9f8..471bfd24` (start of 2026-11-26 to local HEAD on 2026-11-27, 23
first-parent merges). It has the net report, a report per merge, rerun
commands, coverage, timings and a list of findings separating useful ones
from misleading ones: [examples/legato/EVIDENCE.md](examples/legato/EVIDENCE.md).

## POC scope and deferred PRD work

Done here: Git change input (refs and working tree), a Tree-sitter Go
adapter, a SQLite parse/symbol index keyed by blob (rebuilt or reused per
run), exact and structural duplication with change attribution, complexity
deltas, dependency graph deltas and cycles, configurable forbidden rules,
scope and locality, obsolete-code heuristics, a versioned JSON schema, a
Markdown summary, first-parent history, and explicit coverage and
diagnostics.

Deferred:

- TypeScript and C# adapters
- type-resolved references and a call graph (references are name-based)
- detecting lines moved within a file (only symbol moves are detected)
- incremental index across runs beyond the blob cache
- scoring thresholds relative to change size (they are absolute today, so
  large windows saturate at ±2)
- cognitive complexity
- an optional LLM judgment layer (PRD phase 6)
- a correctness signal from independently run tests
- calibrating any score against human review
