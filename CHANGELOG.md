# Changelog

All versions are local proof-of-concept builds; there is no remote or release.

## 0.8.5-poc

- **Tiered Stop hook.** The analyze Stop hook no longer runs the full
  analysis at every stop (about 72 s on Legato's 4,561 files):
  - *Quick, every stop:* v2 debt and the single-file defect rules on the
    changed files, about 1 s on Legato.
  - *Broad, in the background:* cross-file defect rules and
    near-duplication, about 22 s on Legato. It is reported at the next
    stop, labelled with the time of the tree it checked.
  - `CODEEVAL_STOP_CHECKS=quick|all` selects one tier or the old behaviour.
    A new session stops a leftover background run.
- **`analyze --checks quick|broad|all`** (change mode). Rules report the
  skipped tier as `deferred` in coverage. Quick plus broad equals all: this
  is tested on the fixtures and checked on Legato.
- **Faster change checks.** `analyze` and `defects` share one content cache
  between the base and target runs, so an unchanged file is checked once:
  `defects` on a Legato commit drops from 59 s to 30 s, with identical
  output.
- **Fix: near-duplication was nondeterministic.** Each callable keeps its
  best few near-miss partners, and ties were broken by map order, so
  repeated runs of the same change reported 97, 79 and 50 "introduced"
  pairs on Legato (the real number was 0). Ties now break on size, key and
  position. This also affects `serve --near-duplication` and `analyze`
  history.

## 0.8.4-poc

- **Dashboard folder reports:** clicking a folder in the Structure view
  opens a drawer with that folder's `codeeval explain --path` analysis:
  - totals, density against the repository, score and hotspots;
  - debt by component and by subfolder (subfolders are clickable);
  - every factor ranked into severity bands, with collapsed metrics,
    related code and remediation;
  - the folder's defect findings, grouped by rule.
  New endpoints: `GET /api/explain` and `GET /api/defects` (`cached=1`
  answers 204 when the commit has not been checked).
- **Defects on demand:** `serve` no longer runs defect rules or
  near-duplication for every commit by default (`--defects` and
  `--near-duplication` turn them back on). The commit panel and folder
  reports offer a one-commit check instead; its result is kept (last 16
  commits) and shared by every folder report of that commit.

## 0.8.3-poc

- **Dashboard Dependencies view** (alongside Structure and Network):
  - Imports aggregate to project or package groups (`.csproj`,
    `package.json`, `go.mod`, else the top folder), with arrow thickness by
    import count.
  - Layout is layered on the condensation of the strongly connected
    components, so users sit above what they use. Cycle members share a row
    with red edges.
  - Barycentre ordering reduces crossings, and wide layers wrap onto extra
    rows.
  - Clicking a group gives its uses and used-by lists with counts and
    sample import sites; double-click or *Expand* opens a group into
    subfolders.
  - On autopilot the top level is 74 groups in 5 layers with no cycles
    (`Agx.Core` at the bottom). On Legato, expanding `internal` and `web`
    shows `cmd` → `server`/`cli`/`tui` → `service` → `engine` → `leaf`,
    beside the `web/src` cycle.
- **Graph data** (`codeeval graph --json`, `/api/graph`): imports carry
  `count` and up to 3 sample `sites`; `roots` lists manifest directories.

## 0.8.2-poc

- **Dashboard structure view**, now the default graph: a squarified
  treemap of the directory tree.
  - Modules sit inside their folders. Single-child chains collapse, and
    one-module folders become the tile.
  - Area is lines or debt (*Size by*); colour is debt density.
  - Modules in an import cycle get a red corner marker instead of a mesh of
    arrows. On a 1,483-module repository the old view drew 1,705 cycle
    arrows.
  - A hovered or selected module shows its imports and importers.
  - The force layout stays as the "Network" view, which leaves out modules
    with no debt and under 150 lines.

## 0.8.1-poc

- **C# dependency graph fix.** The 0.8.0 graph linked a `using N` to every
  directory declaring namespace `N`. On a real repository this produced an
  impossible cycle across two projects, because one project declared
  namespaces named after the other. It also turned unused usings into
  edges. A `using` now links only to the following:
  - directories of `N` that declare a type the file references (`[Foo]`
    matches `FooAttribute`);
  - failing that, directories declaring an extension method (first parameter
    `this`) the file calls by name.

  Edges are also limited to the project itself and to projects it
  references (`.csproj` `ProjectReference`, transitively).

  On the same repository: cross-project cycles went from 1 to 0, directories
  in cycles from 142 to 128, and edges from 5,293 to 4,438. The 20 cycles
  that remain each lie within one project, and the sample checked against
  source is real. C# graphs, and the v2 architecture debt from them, differ
  from 0.8.0.

## 0.8.0-poc

- **C# adapter** (`internal/lang/csharp`, tree-sitter-c-sharp v0.23.5, the
  only new dependency):
  - Symbols: types (nested as Outer.Inner) and callables (methods,
    constructors, destructors, operators, properties and indexers with
    bodies).
  - Complexity: if, loops, switch sections, switch-expression arms (not `_`),
    catch, ternary, `&&`, `||` and `??`. Lambdas and local functions count
    toward their member.
  - Fingerprints and name-based references.
  - A directory graph from `using` directives resolved to the directories
    that declare each namespace. Namespaces not declared in the repository
    are external.
  - Limits, stated in reports: no type resolution; implicit
    enclosing-namespace access and fully qualified names are not followed;
    `.csproj` is not read; top-level statements are not measured.
  - On a 5,221-file C# repository, all but 40 files parse (34 of the 40 are
    generated with `#if` branches), and C# plus TypeScript now cover 98% of
    its source bytes.
- **C# defect rules** (`internal/rules/csharp`): 8 of 13.
  - Implemented: unused variable, unreachable code, unused private method
    (name-based across the repository, for partial classes), empty catch,
    catch returning a constant, tautological xUnit/NUnit/MSTest/
    FluentAssertions assertions, un-awaited async assertions, and
    architecture rules through the graph.
  - Recorded gaps, with reasons: unused usings, unused public API, and the
    three dependency rules.
  - The defects fixture has C# in its clean base and bad variant.
- Reports on repositories with C# are not comparable to 0.7 reports, which
  counted C# as unsupported coverage.
- **`codeeval hook install|uninstall --tool claude-code`:**
  - Installs the analyze Stop hook and session-base hook (embedded in the
    binary) into the git repository you run it from, and registers them in
    `.claude/settings.json`.
  - Keeps other settings, hooks and key order, and refuses invalid JSON.
  - Idempotent.
  - `--bin`, `--config`, `--dry-run`.

## 0.7.0-poc

- **Provenance standard fields are now model and harness:**
  - Trailers and notes use `Agent-Model:` and `Agent-Harness:`, and any other
    `Agent-<Key>:` stays free-form. A bare `Agent:` line (0.4–0.6) is still
    read as the model.
  - Records and JSON use `model` in place of `agent`. This changes
    `codeeval.attribution/v1` and the dashboard API.
  - `attribution` and the dashboard group by model and harness by default.
    `agent` is still accepted as an alias of `model` in `--group-by`.
  - The hooks read `AGENT_MODEL`, with `AGENT_ID` still accepted.
- **`Taskfile.yml`:** build, install, link, test (plus short, race and cover
  variants), vet, fmt, check, run, serve, analyze and fixtures:reports. It
  uses the offline pinned environment and `CGO_ENABLED=1`, because the
  Tree-sitter grammars are C.
- **Incremental defect rules** (`internal/memo`, `rules.Options.Cache`,
  `Impl.FileLocal`, `Registry.Prepare`):
  - Files load lazily, and their syntax trees are freed after the file's
    checks.
  - Single-file rules run file by file, with findings cached by content
    (blob id).
  - Cross-file rules read cached per-file facts (identifiers, declarations,
    imports) instead of syntax trees.
  - Clone extraction is cached by content the same way
    (`clones.AnalyzeCached`).
  - Results are byte-identical: the golden outputs and Legato defect reports
    match. On Legato, one run's peak memory went from 909 MB to 224 MB.
  - Across commits, `serve`'s detail phase went from about 36 s to about 5 s
    per commit after the first.
- **Dashboard in two phases:** charts and graphs first (v2 metrics, about
  0.4 s per commit), then defects and near-duplication. The detail phase
  does not hold the evaluator lock, so graph requests no longer wait on it.
- **Dashboard measured window:**
  - The whole first-parent history is listed from a few batched git calls,
    and only a window is measured: `--last 50` by default, or
    `--since`/`--until`.
  - On the page, a history strip shows every commit; dragging across it
    or using the last-N or date controls moves the window
    (`POST /api/window`, same-origin JSON only).
  - Charts, the stepper and the group table cover the window, and measured
    commits stay cached.
  - Graphs no longer prefetch neighbours except while playing.
  - `--base` lists the base commit as the first point.
- **`codeeval serve` defaults:**
  - `codeeval serve` (or `codeeval serve PATH`) shows the whole first-parent
    history of the checkout's HEAD. `--max` (now default 0 = all),
    `--target` and `--base` narrow it.
  - Metrics are computed newest-first.
  - Step changes are filled in when the parent is measured (`stepDone`).
- **Dashboard graph:**
  - It fills its card.
  - A "How to read this graph" section explains nodes, sizes, colours, lines,
    rings and the debt formula, using the configured thresholds.
  - Colour scales are explicit and fixed.
  - Clicking a node shows its detail (module debt by component; for a
    function, where its points come from and its clone partners).
  - A clickable list shows the biggest changes in each commit.
  - Functions are matched across commits by file and name, not line.
- **Provenance parsing:** message lines are read only from trailer-style
  paragraphs, in which every line is `Key: value`. Before, a prose line
  that wrapped to start with `Agent:` was parsed as the model; this
  happened in this repository's commit 18c442e.

## 0.6.0-poc

- **`codeeval serve`** (`internal/dashboard`): a local, read-only web
  dashboard of quality over the first-parent history.
  - Per-commit charts: v2 debt by component, score, debt/kLOC, hotspots,
    defects, near-duplication.
  - A commit stepper with the module graph changing in place: stable
    positions, rings for debt up, down and new.
  - Commit details with provenance and deltas.
  - A table of step changes grouped by agent, harness or any provenance
    metadata.
  - Loopback addresses only; a self-contained page with a strict CSP; metrics
    are computed in the background.
- **Provenance:** inferred co-author agent names are ID-styled
  ("Claude Opus 5.5" → `claude-opus-5-5`), so they group with declared IDs.

## 0.5.0-poc

- **Escape hatches for the Stop hook and the rules:**
  - `codeeval-ignore <rule>[, …]: <reason>` comments suppress defect findings
    on the same or the next line. Suppressed findings are listed in reports
    and never counted.
  - `analyze` reports a `failureFingerprint`.
  - The analyze Stop hook never blocks twice on the same fingerprint in a
    session, has a block budget (`CODEEVAL_HOOK_MAX_BLOCKS`, default 3) and
    can be switched off (`CODEEVAL_HOOK=off`).

- **`--working-tree` reads untracked files by default:** files Git does not
  ignore are read, each up to 8 MB, and appear as added files in the change
  set. `--untracked=false` restores the earlier tracked-only behaviour.
  Ignored files are never read. This applies to `evaluate`, `check`,
  `defects` and `analyze`, so their working-tree results differ from 0.4.0.
  Commit-to-commit results are unchanged.
- **Refactors guided by `codeeval analyze` on this repository.** Golden
  outputs of every command on all fixtures and on this repository are
  byte-identical before and after.
  - Shared CLI plumbing (`cmd/codeeval/cli.go`): v2 debt −124.
  - Shared Tree-sitter helpers for both adapters
    (`internal/lang/sitterutil`): −18.
- **Dogfooding:**
  - an analyze Stop hook;
  - a session-base SessionStart hook;
  - provenance written as git notes (`integrations/git/post-commit`).

## 0.4.0-poc

- **`codeeval analyze`** (schema `codeeval.analysis/v1`): one report holding
  the v2 debt factors (explain, or check in change mode) and the defect
  findings, with a summary and a per-file table ranked by gate defects, then
  debt, then advisory defects. Debt and defects are not summed into a score.
  `--debt-fail-on` (change mode) and `--defects-fail-on` control exit 3.
- **Near-miss and fragment duplication** (`internal/clones`, `clones/v1`,
  unscored), reported by `analyze`:
  - near misses are found with MinHash/LSH candidates, then verified exactly
    with LCS similarity on abstracted tokens;
  - fragments are shared runs of ≥ 50 abstracted tokens, with periodic
    repetition dropped;
  - in change mode only introduced pairs are listed.

  v2 duplication debt is unchanged.

- **Agent provenance:** `internal/provenance` reads it from commit-message
  lines (`Agent:` and `Agent-Harness:` are standard, `Agent-<Key>:` is
  arbitrary metadata), anywhere in the body except the subject, so it
  survives squash merges. Notes in `refs/notes/agent-provenance` fill in
  missing keys. Optionally (on by default), the agent is inferred from
  `Co-Authored-By` lines naming a known coding agent; that source is always
  labelled.
- **`codeeval attribution`** (schema `codeeval.attribution/v1`): for each
  first-parent step it reports the v2 debt delta, changed lines and
  decision points, and with `--defects` the defects introduced and resolved.
  Steps are grouped by `--group-by` provenance fields. Merge steps include
  the merged commits; a merge commit without provenance is ignored for
  grouping, and a step whose commits disagree is `(mixed)`. Observational
  only.
- **`integrations/git/commit-msg`:** a hook template that writes the
  trailers from `AGENT_ID`, `AGENT_HARNESS` and `AGENT_META_<KEY>`.

## 0.3.0-poc

- **`codeeval defects`** (`rules/v1`, schema `codeeval.defects/v1`): concrete
  defect findings for a revision, or the findings a change introduced and
  resolved. Findings are matched by line-independent keys. Exit 3 on
  introduced gate findings (`--fail-on gate|any|none`); also `--rule` and
  `--path`. There are 13 rules in five categories (dead code, lost errors,
  dependency integrity, architecture policy, test-oracle defects), and every
  rule is implemented for both Go and TypeScript. Findings are separate from
  v1 and v2 scoring, and that scoring is unchanged.
- **`codeeval rules`:** lists the rule × language matrix with
  per-implementation coverage notes.
- **Rule framework** (`internal/rules`): a rule is defined once in
  `catalog` and implemented per language (`rules/golang`,
  `rules/typescript`, or `rules/generic` for graph-only checks). Includes a
  `rulestest` harness and a contract test that every rule covers every
  language.
- **Adapters:** they implement `lang.SyntaxParser`. The TypeScript import
  resolution has moved into an exported `typescript.Resolver`, which the
  graph and the dependency rules share. Graph behaviour is unchanged.
- **Config:** architecture rules accept `only` allow-lists (enforced by
  `defects` only). A new `defects.disable` setting turns rules off. Neither
  changes existing configuration hashes. A rule now needs `from` plus `to`
  or `only`.
- **Fixtures:** `testdata/fixtures/defects`. Its base is clean, `bad` has
  every rule in both languages, and `clean` has the same features without
  defects.

## 0.2.0-poc

- **TypeScript adapter** (`.ts`, `.tsx`, `.mts`, `.cts`), built on
  tree-sitter-typescript v0.23.2 (the only new dependency).
  - Parsing: uses the TSX grammar, falling back to the TypeScript grammar.
  - Symbols: functions, classes and methods, interfaces, types and enums,
    and top-level arrow-function components and hooks.
  - Complexity: counts if, loops, case, catch, ternary, `&&`, `||` and
    `??`. Inline callbacks count toward the enclosing callable.
  - Graph: a directory-level import graph resolved through relative paths,
    tsconfig `paths` (aliases into node_modules are external) and workspace
    packages.
  - Fixtures: a paired TypeScript fixture scenario. The edge-case fixtures
    now use Kotlin as the unsupported language.
- **`codeeval score`** (experimental `scoring/v2`):
  - snapshot debt points (hotspot complexity and nesting excess,
    structural-clone redundant tokens, configured rule violations, cycle
    units) for the codebase, each module and each language;
  - a 0–100 size-normalised score, for comparison only;
  - deltas that telescope across first-parent steps;
  - the decision-point versus function-base-point split, and advisory gates.
- **`codeeval explain`:** a scoped factor report (module, application, path
  or glob). Each factor has a location, debt, share and rule-based
  remediation.
- **`codeeval check`:** a PR check listing factors introduced or worsened,
  with targets relative to the pre-change value. Fail modes are `new-debt`,
  `net-debt`, `gates` and `none`; exit 3 on failure. It also supports
  `--working-tree`.
- **`codeeval graph`:** an offline interactive force-graph view of v2 debt
  (POC).
- **Claude Code templates:** the `codeeval-complexity` skill and a Stop hook
  (`integrations/claude-code/`).
- Engine accessors for snapshot-level scorers. `scoring/v1` output is
  unchanged for Go-only repositories; repositories with TypeScript now get
  TypeScript coverage.

## 0.1.0-poc

- Go adapter, a SQLite index, and the v1 analyzers (duplication,
  complexity, dependencies, scope, obsolete, correctness = not assessed)
  with `scoring/v1` −2…+2 per analyzer and no overall score.
- `evaluate` and `history` commands, paired fixtures, end-to-end CLI tests,
  and the Legato evaluation. The frozen build `5e0c679` is the provenance for
  the Legato v1 reports and fix-batch measurements.
