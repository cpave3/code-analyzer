# Changelog

All versions are local proof-of-concept builds; there is no remote or release.

## 0.4.0-poc

- **`codeeval analyze`** (schema `codeeval.analysis/v1`): one report holding
  the v2 debt factors (explain, or check in change mode) and the defect
  findings, with a summary and a per-file table ranked by gate defects, then
  debt, then advisory defects. Debt and defects are not summed into a score.
  `--debt-fail-on` (change mode) and `--defects-fail-on` control exit 3.

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
