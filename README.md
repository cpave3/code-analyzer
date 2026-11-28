# CodeEval (POC)

A local, read-only Go CLI that measures code quality **in the context of the
repository**: duplication against existing code, complexity, nesting,
dependency structure and rules, scope, and code left obsolete. It judges a
change (or a snapshot) itself, not what an agent or workflow says about it.
See [PRD.md](PRD.md) for the product brief, [TECH.md](TECH.md) for design and
trade-offs, and [CHANGELOG.md](CHANGELOG.md) for what changed in each version.

This is a **proof of concept**. It measures structure. It never runs the
evaluated code and makes no claim about correctness or design quality.

## Commands at a glance

| Command | What it answers | Scoring |
|---|---|---|
| `codeeval evaluate --base A --target B` | Did this change improve or degrade each dimension? Gives a −2…+2 score per analyzer with evidence; there is no overall score | `scoring/v1` (frozen) |
| `codeeval history --base A --target B` | The same, for the net change and for each first-parent merge | `scoring/v1` |
| `codeeval score --rev R` | Absolute debt and 0–100 score for a codebase, each module and each language | `scoring/v2-experimental` |
| `codeeval score --base A --target B [--steps]` | Debt ledger for a change, where per-commit steps sum exactly to the net | v2 |
| `codeeval explain --rev R --path P` | Which items make up a scope's debt, where they are, how much each contributes, and how to fix them | v2 |
| `codeeval check --base A --target B\|--working-tree` | PR check: items a change introduced or worsened. Exit 3 on failure | v2 |
| `codeeval graph --rev R --output g.html` | Interactive offline debt graph (POC) | v2 |
| `codeeval defects --rev R` or `--base A --target B\|--working-tree` | Concrete defects (dead code, lost errors, dependency integrity, architecture policy, test-oracle defects), or those a change introduced. Exit 3 on introduced gate findings | `rules/v1` (no score) |
| `codeeval analyze --rev R` or `--base A --target B\|--working-tree` | Combined view in one report: v2 debt factors, defect findings, and near-miss and fragment duplication. A per-file table ranks files by gate defects, then debt, then advisory defects, then near-duplicate pairs. No combined score | v2 + `rules/v1` + `clones/v1` |
| `codeeval serve [--base A] [--target B] [--max N]` | Local web dashboard on 127.0.0.1: per-commit debt, score, hotspots, defects and near-duplication charts; a commit stepper that shows the module graph change; provenance grouping (agent, harness, workflow) | v2 + `rules/v1` + `clones/v1` |
| `codeeval rules` | The defect rules and which languages implement each, with per-language coverage notes | – |
| `codeeval attribution --base A --target B --group-by model,harness,workflow [--defects]` | Per-step v2 debt deltas (and defects) grouped by commit provenance: which model, harness and workflow made each change | v2 (+ `rules/v1`) |
| `codeeval languages` / `analyzers` / `version` | Adapters and their limits, analyzers, version | – |

Languages: **Go** and **TypeScript** (Tree-sitter adapters). Other languages
are counted as unsupported coverage, never scored as neutral. Claude Code
templates (a skill and a Stop hook) are in
[integrations/claude-code/](integrations/claude-code/README.md). Worked
examples on a real repository are in [examples/legato/](examples/legato/).

**scoring/v1 vs v2:** v1 is the original change-attributed −2…+2 model. Its
code and outputs are unchanged, and the Legato v1 reports were made with a
frozen build of it. v2 is an experimental, uncalibrated debt model added
later. Neither v2's debt nor its score should be optimised on its own.

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

### Pinned, reproducible build

For a binary that another workflow pins by checksum, build from a clean,
detached clone with `-trimpath` (no local paths in the binary, and
`vcs.modified=false` in its build info). `W` is any scratch directory
outside the repository:

```bash
W=$(mktemp -d)
git -c core.hooksPath=/dev/null clone -q --no-local /home/cameron/Projects/code-evaluator "$W/src"
git -C "$W/src" checkout -q --detach <commit>
cd "$W/src" && env PATH=/home/cameron/.local/share/mise/installs/go/1.26.0/bin:/usr/bin:/bin \
  GOTOOLCHAIN=local GOPROXY=off GOFLAGS=-mod=readonly CGO_ENABLED=1 \
  go build -trimpath -o "$W/bin/codeeval" ./cmd/codeeval
sha256sum "$W/bin/codeeval"; go version -m "$W/bin/codeeval" | grep vcs
```

The Legato cleanup pilot pins the build of `6639e38`:
sha256 `e2ed35adaa66e85664fabf9a1dd71cb9a60b93a69df5c3d5e62ec284fba90229`
(same host toolchain; reproduced independently). The Legato v1 reports use
the frozen `scoring/v1` build of `5e0c679` (see
[examples/legato/EVIDENCE.md](examples/legato/EVIDENCE.md)). Don't rebuild a
pinned binary from a different commit under the same name.

```bash
go build -o bin/codeeval ./cmd/codeeval

# Compare two revisions (JSON + Markdown summary)
bin/codeeval evaluate --repo ../some-repo --base main --target HEAD \
  --output /tmp/eval.json --summary /tmp/eval.md

# Working-tree changes against a base: tracked files plus untracked files that
# are not ignored (since 0.5.0; --untracked=false reads tracked files only)
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
- treats repository content, including commit subjects, as data;
- in `--working-tree` mode, reads tracked files and untracked files that Git
  does not ignore (`.gitignore`, `.git/info/exclude`, `core.excludesFile`),
  each up to 8 MB. Ignored files are never read. `--untracked=false` limits
  it to tracked files. Contents are only parsed or searched and never
  copied into reports; paths of untracked files can appear as added files.

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

## Experimental scoring/v2: debt, absolute score and telescoping deltas

`codeeval score` runs a second, separately versioned scoring model
(`scoring/v2-experimental`). It does not change `evaluate`/`history` or any
`scoring/v1` output; v1 JSON from this binary is byte-identical to earlier
runs.

```bash
codeeval score --repo PATH --rev REV [--config FILE] [--output F.json] [--summary F.md]          # absolute snapshot
codeeval score --repo PATH --base REV --target REV [--steps] [--config FILE] ...                # delta (+ per first-parent step)
```

**Debt points** are measured on a single snapshot, for the codebase and for
each module (Go package):

| Component | Raw | Default weight |
|---|---|---|
| complexity | Σ max(0, cyclomatic − H + 1): debt exactly when a callable is a hotspot (H = 10) | 1 per point |
| nesting | Σ max(0, depth − N + 1), N = 4 | 2 per level |
| duplication | redundant tokens in structural clone groups (each member carries (n−1)/n of its tokens) | 1 per 20 tokens |
| architecture | configured-rule-violating edges, units inside dependency cycles | 10 per edge, 20 per unit |

Properties: module debts **sum** to codebase debt; deltas are after − before
of one snapshot function, so per-step deltas **telescope** to the net delta
(`--steps` checks this). The **score** is `100 × h ÷ (h + debt per 1,000
supported source lines)` with `h = 25`: a size-normalised view for comparing
modules or codebases. **Neither scalar should be optimised on its own.**
The score can be diluted by adding simple code (shown on Legato in
[examples/legato/v2/V2-EVIDENCE.md](examples/legato/v2/V2-EVIDENCE.md)).
Debt is not diluted that way, but it is an uncalibrated, size-sensitive
structural proxy: deleting features or code, or reshaping syntax (for
example splitting functions below the thresholds), can lower it without
improving engineering. Use debt together with its components, the raw
measurements and evidence, and behaviour/test, coverage and scope checks.
Telescoping is mathematical consistency under a fixed model and config, not
empirical calibration.
Deltas also report the changed-callable split into decision points versus
function base points, debt per changed kLOC, per-module deltas, and advisory
gates (no new exact/structural duplicates, rule violations or cycles; no
branch growth in existing callables; no hotspot-excess growth).

Limits: weights are uncalibrated defaults recorded in every report; only
languages with adapters are measured (coverage is reported, unsupported is
not neutral); complexity growth below the hotspot thresholds carries no debt
(it shows as decision points); architecture debt reflects only configured
import-level rules; correctness is not measured.

## Debt factor reports and PR checks (for fixing, by people or agents)

```bash
# Why does this module, application, path or glob carry debt? Largest factors first.
codeeval explain --repo PATH --rev REV [--path internal/service] [--path 'web/**/*.tsx'] [--top N] [--detailed N]
# What did this change add? Exit 3 when the check fails.
codeeval check --repo PATH --base REV (--target REV | --working-tree) [--path P] \
  [--fail-on new-debt|net-debt|gates|none] [--top N]
```

Every **factor** is one item of v2 debt in scope. It has a `kind`
(`complexity`, `nesting`, `duplication`, `rule_violation` or `cycle`), a
location (file, lines, symbol), metrics, debt points, its share and
cumulative share of the scope's debt, an explanation, related locations
(clone members, cycle edges), and deterministic **remediation**: strategy,
target, steps, acceptance criteria and caution.
- **explain:** factors add up exactly to the scope's debt. A concentration
  line shows how many factors carry 50% and 80% of it.
- **check:** lists the factors a change *introduced* or *worsened* in scope.
  For a worsened factor the target is relative to its pre-change value. It
  also lists what the change improved, the decision-point and base-point
  split, and the advisory gates.
- **Fail modes:** `new-debt` (any new or worsened factor), `net-debt` (debt
  in scope rises), `gates`, or `none`.

JSON schemas are `codeeval.explain/v2` and `codeeval.check/v2`. Claude Code
templates (a skill and a Stop hook) are in
[integrations/claude-code/](integrations/claude-code/README.md), with Legato
samples in `examples/legato/fix/`. The remediation text is rule-based guidance,
not a guarantee. The numbers are uncalibrated structural proxies, and nothing
here checks behaviour.

## Debt graph view (POC)

```bash
codeeval graph --repo PATH --rev REV --output debt-graph.html [--json data.json] [--max-callables N] [--config FILE]
```

This writes a self-contained, offline HTML page: a force-directed graph in the
style of Obsidian's graph view, showing v2 debt.
- **Nodes:** modules, plus the callables that carry complexity or nesting
  debt or belong to a structural clone group.
- **Links:** cyan links join clone-group members. Optional import edges
  are shown in red where they violate a configured rule or form a cycle.
- **Controls:** group by top-level folder, language or dominant debt
  component; colour by debt density, component, language or folder; size by
  debt, lines or complexity; filter by language and minimum debt; search;
  hover and click for details.

Sizes and colours are uncalibrated structural proxies, not quality verdicts.
See `examples/legato/viz/`.

## Defect rules (`codeeval defects`, rules/v1)

```bash
# Every finding in a revision (optionally one category or rule, one scope)
codeeval defects --repo PATH --rev REV [--rule dead-code] [--path web/src]
# What a change introduced and resolved. Exit 3 when it introduces a gate finding.
codeeval defects --repo PATH --base origin/main --target HEAD
codeeval defects --repo PATH --base HEAD --working-tree
codeeval rules    # rule × language matrix with coverage notes
```

These checks report concrete defects, not whether code "looks AI-written".
They are findings, not a score, and they are not part of v1 or v2 debt.
Nothing is type-checked, built or run, so each rule covers what syntax and
repository files can prove. Each language's implementation states what it
covers and misses (`codeeval rules`, and `rules[].implementations` in the
JSON report).

| Rule | Severity | What it catches (Go / TypeScript) |
|---|---|---|
| `dead-code/unused-import` | gate | Unused import or binding. Go: unnamed imports whose package name is only inferred from a non-stdlib path are advisory |
| `dead-code/unused-variable` | gate | Local declared and never read in its function (TS: also non-exported module-level bindings) |
| `dead-code/unreachable-code` | gate | Statement after return/throw/panic/break/continue/goto in the same block |
| `dead-code/unused-private-function` | gate | Unexported Go function unused in its package; module-private TS function or private method unused |
| `dead-code/unused-export` | advisory | Go: exported function under `internal/` never referenced. TS: export no file imports, in a module others import by name |
| `lost-errors/empty-handler` | gate | Empty `if err != nil {}`, discarded `recover()`, `_ = err` (advisory); empty `catch {}`, empty `.catch()`/`.then(_, fn)` handler. Comment-only handlers are advisory |
| `lost-errors/error-to-success` | advisory | `if err != nil { return …, nil }`; `catch { return null/[]/{} }`, `.catch(() => constant)` |
| `dependencies/unresolved-import` | gate | Import of repository code that resolves to nothing in the snapshot (TS alias/workspace misses are advisory) |
| `dependencies/undeclared-dependency` | gate | External import not covered by `go.mod` requires / any `package.json` from the file up to the root |
| `dependencies/unused-dependency` | advisory | Direct `go.mod` require or `package.json` dependency nothing imports or names |
| `architecture/import-policy` | gate | An edge that breaks a configured rule: `to` (forbidden) or the new `only` allow-list |
| `test-oracle/tautological-assertion` | gate | An assertion comparing an expression with itself or two literals, or asserting a constant |
| `test-oracle/lost-async-assertion` | gate | `expect(…).resolves/.rejects` neither awaited nor returned; `t.Fatal`/`FailNow`/`require.*` inside a goroutine |

- **Change mode:** findings are matched across revisions by rule, file and a
  line-independent key, so pre-existing defects never fail a change. A
  finding in a moved file counts as resolved plus introduced.
- **Defaults:** `--fail-on gate` in change mode. In snapshot mode the
  default is `none` (report only).
- **Not covered:** redundant conditions need types and control flow. Mutation
  testing needs the tests to run, which the read-only rules forbid. Both are
  better left to `tsc --strict`, staticcheck and a mutation tool.
- **Suppressing a false positive:** put a comment on the line or the line
  above: `// codeeval-ignore dead-code/unused-export: kept for the plugin API`.
  It accepts a comma-separated list of rule IDs or categories, and the text
  after `:` is the reason. Suppressed findings are left out of counts and
  comparisons but listed in reports with their reason, so a suppression is
  visible in review.
- **Configuration:** `"defects": {"disable": ["dead-code/unused-export"]}`
  turns off rules or whole categories. An architecture rule with
  `"only": ["web/src/shared/**"]` (and `"to": ""`) is an allow-list:
  `from` units may depend only on matching units or on units matching
  `from` itself. Allow-lists are enforced by `defects` only; v1 and v2 still
  read only `to`.

## Near-miss and fragment duplication (`clones/v1`, unscored)

v2 duplication debt counts only whole callables with identical structure
(exact, or with identifiers and literals renamed). `codeeval analyze` also
reports two kinds of duplication that this misses:
- **Near-miss clones:** callables whose abstracted token sequences
  (identifiers → ID, literals → LIT, no comments) are similar but not
  identical. The threshold is `duplication.nearThreshold` (default 0.85) on
  2·LCS/(|a|+|b|), and both callables must be ≥ `minTokens`. Candidates come
  from MinHash/LSH and every pair is verified exactly. This catches renamed
  copies where a statement was added, removed or changed.
- **Fragment clones:** a run of ≥ max(50, `minTokens`) identical abstracted
  tokens shared by two different callables, or by two non-overlapping places
  in one. This catches a copied block inside otherwise different code. Runs
  made of a short repeated unit (table rows, a list of calls) are dropped.

In change mode only the pairs the change introduced are listed (matched by
callable names, not lines). These pairs are not part of v2 debt.
Folding them into debt would be a new scoring version, after calibration.
Similar structure is not always a copy: parallel per-language
implementations and small boilerplate functions also match.

## Dashboard: quality over time (`codeeval serve`)

```bash
codeeval serve --repo PATH                      # the last 60 first-parent commits of HEAD
codeeval serve --repo PATH --base v1.0 --target main --addr 127.0.0.1:8765
```

It serves a self-contained page (no external scripts) on a loopback
address; other addresses are refused. Metrics are computed in the
background, oldest commit first, about 4 s per commit on this repository,
and appear as they arrive.
- **Charts:** v2 debt by component, score, debt/kLOC, hotspots and size,
  defect findings and near-duplication pairs, per commit. Points are
  coloured by the provenance group chosen above them.
- **Stepper:** ◀ ▶, a slider, ▶ Play and the arrow keys step through
  commits. The module graph keeps node positions between commits, and rings
  mark nodes whose debt rose (red) or fell (green), and new nodes (dashed).
  The panel also shows the commit's subject, author, provenance and deltas
  against the previous commit.
- **Group table:** step changes summed by provenance group (`--group-by` as
  checkboxes: model, harness and any `Agent-<Key>` metadata; model and harness by default). As with
  `attribution`, this is observational, not causal.

`--defects=false` and `--near-duplication=false` skip the slower per-commit
checks.

## Attribution by model, harness and workflow

Commits can declare what made them with git trailers. `Agent-Model:` and
`Agent-Harness:` are the standard fields (a bare `Agent:` from 0.4–0.6 is
read as the model); any `Agent-<Key>:` is free-form
metadata, such as `Agent-Workflow: proposal-stack@2`. Commits that are
already made can be annotated with notes in `refs/notes/agent-provenance`.
`codeeval attribution` groups first-parent step deltas by those fields. The
format, the `commit-msg` hook template and the reading rules are in
[integrations/git/](integrations/git/README.md). The comparison is
observational: groups differ in tasks, so differences are leads, not
effects.

## Languages

| Language | Status |
|---|---|
| Go | Supported: Tree-sitter adapter with symbols, complexity, fingerprints, package import graph (via the nearest `go.mod`) and name-based references |
| TypeScript (`.ts`, `.tsx`, `.mts`, `.cts`) | Supported (since 0.2.0): Tree-sitter adapter (TSX grammar with TypeScript fallback) with functions, classes/methods, interfaces/types/enums and top-level arrow-function components/hooks; complexity; fingerprints; a directory-level import graph resolved via relative specifiers, `tsconfig` `paths` and workspace package names; name-based references |
| C# | Recognised and counted as unsupported coverage. Next adapter (see below) |
| JS, Kotlin, Rust, Python, Swift, Java, SQL, shell | Recognised and counted as unsupported coverage |

The Go adapter does not type-check, resolve interface dispatch or build a call
graph, and it indexes all build-tag variants together. The TypeScript adapter
does not type-check either. It measures inner callbacks as part of the
enclosing callable (as Go func literals are), so React components with many
inline handlers and JSX conditionals score high. Its graph units are
directories, where cycles are legal. It does not follow node_modules, package
`exports` maps or re-export chains. These limits are also written into every
report. Cyclomatic counts are not directly comparable across languages.

Adding TypeScript changes what a report covers. Reports from 0.2.0 on
repositories with TypeScript are not comparable to 0.1.0 Go-only reports;
`tool.version` records which one produced a report.

## Extending

The CLI has no `switch` over languages or metrics. It uses two registries,
each wired in one file. The registration contract is tested in
`internal/engine/contract_test.go`: a toy language and a toy analyzer are
added through the public registries without touching existing code.

### Add a language (e.g. C#)

`internal/lang/typescript/` is the most complete reference adapter (two
grammars, tsconfig/workspace resolution, callback-inclusive complexity).

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
   - `BuildGraph(GraphInput)`: resolve only what you can. The TS adapter
     reads `tsconfig.json` paths and workspace `package.json` names; C#
     would read `.csproj` `ProjectReference` and namespaces
2. Register it in `internal/languages/languages.go`. Its `Register` call
   supersedes the existing `Declare` entry for the same extensions.
3. Add a fixture scenario under `testdata/fixtures/` and assertions in
   `e2e/`.

Analyzers need no change. If your adapter lacks a capability, the analyzers
that need it report your files as unsupported coverage.

### Add a defect rule

Rules are defined once and implemented per language, so a rule can start in
one language and gain others later.

1. Define it in `internal/rules/catalog/catalog.go`: ID `<category>/<name>`,
   severity (`gate` = provable, `advisory` = needs a human), summary and
   limits.
2. Implement it per language as a `rules.Impl{Rule, Note, Check}`. Put it in
   `internal/rules/golang/` and `internal/rules/typescript/` and add it to
   that package's `Register` list. A check that needs only the dependency
   graph goes in `internal/rules/generic/`. `Check(ctx)` gets the language's
   files with their Tree-sitter trees, all snapshot paths, a file reader, the
   graph and the config. It calls `ctx.Report(Finding{File, Line, Key,
   Message})`. The `Key` must not contain line numbers. Shared values (a
   resolver, parsed manifests) go through `ctx.Memo`.
3. Test it with `internal/rules/rulestest.Run` on in-memory sources, and add
   the defect to `testdata/fixtures/defects/bad`.
   `TestDefectsEveryRuleAndLanguageDetects` requires every rule × language
   pair to fire there. `internal/rules/all` fails when a rule lacks a
   language without a recorded reason (`knownGaps`).

A new language adapter gets rules by implementing `lang.SyntaxParser` and
adding a `internal/rules/<lang>` package to `internal/rules/all`.

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
  "architecture": {"rules": [
    {"name": "store must not depend on presentation", "from": "internal/store/**", "to": "internal/cli/**"},
    {"name": "lib allow-list", "from": "web/src/lib/**", "to": "", "only": ["web/src/shared/**"]}
  ]},
  "scope": {"formattingOnlyChurnAt": 3},
  "defects": {"disable": []}
}
```

Each report records `configVersion`, a hash of the effective config and the
`scoringVersion` (`scoring/v1`). The `defects` section and `only` lists are
omitted from the hash when empty or absent, so hashes of earlier
configurations are unchanged.

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
`partial`), unsupported-only (a Kotlin file: `unsupported`, score null), mixed
Go and Kotlin (`partial`), and an odd filename (`util/odd name ü.go`).
TypeScript pairs (`testdata/fixtures/ts`): reuse versus a renamed copy
(duplication 0 vs -1), extraction versus in-place nesting (complexity 0 vs
-1), and an introduced directory cycle (dependencies -2).

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

The same folder holds, each linked from `EVIDENCE.md`:

- frozen-evaluator measurements of the Legato cleanup fix batches
  (#434–#446) and of the final combined stack (`*-EVIDENCE.md`);
- the experimental v2 debt ledgers (`v2/`) and the TypeScript evidence
  (`ts/`), which records partial coverage and parser gaps;
- the debt graph (`viz/`) and sample factor reports and PR checks (`fix/`).

## POC scope and deferred PRD work

Done here:
- **Inputs and index:** Git change input (refs and working tree), Tree-sitter
  Go and TypeScript adapters, and a SQLite parse/symbol index keyed by blob.
- **v1 analyzers:** exact and structural duplication with change
  attribution, complexity deltas, dependency graph deltas and cycles,
  configurable forbidden rules, scope and locality, and obsolete-code
  heuristics.
- **v1 output:** versioned JSON and Markdown, first-parent history, and
  explicit coverage and diagnostics.
- **v2 (experimental):** snapshot debt per module and language, telescoping
  deltas, advisory gates, scoped factor reports with remediation, a PR check,
  and the graph view.
- **Integration:** Claude Code skill and hook templates.

Deferred:

- C# adapter
- calibrating v2 weights and thresholds (per language), and percentile bands
  from a repository's own history
- rule types the config cannot express (allow-lists, stdlib-only,
  symbol-level), and TypeScript architecture rules for real projects
- type-resolved references and a call graph (references are name-based)
- detecting lines moved within a file (only symbol moves are detected)
- cognitive complexity
- an optional LLM judgment layer (PRD phase 6); remediation text is
  rule-based today
- a correctness signal from independently run tests
- calibrating any score against human review

## Dogfooding in this repository

`.claude/settings.json` wires two Claude Code hooks:
- **SessionStart:** records the session's base commit.
- **Stop:** runs the existing `bin/codeeval analyze` (it does not rebuild)
  from the session base to the working tree. It blocks the stop on new gate defects or
  a debt rise.

The same file sets `AGENT_MODEL` and `AGENT_HARNESS`, and `.git/hooks/post-commit`
(a link to `integrations/git/post-commit`, local and not tracked) records
them as provenance notes. Run `codeeval attribution` to read them.

## For agents working on this repository

`main` is the integration branch (the POC was built on
`poc/go-first-evaluator` and merged into it). There is no remote.

This project's learnings, pitfalls and open threads are kept in the engram
identity **`assayer`**. Run `engram --as assayer boot` at the start of a
session and `engram recall "<query>" --as assayer` before assuming. Build and
test with the offline environment described under
[Build environment](#build-environment-where-this-poc-was-built-and-tested).
Never modify a repository you evaluate.
