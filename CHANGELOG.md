# Changelog

All versions are local proof-of-concept builds; there is no remote or release.

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
