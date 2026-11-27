# CodeEval report: edge

- **Base:** `558cf05221d55f55ce1588532914292c021130ee` (558cf05221d55f55ce1588532914292c021130ee) tree `81e6e56ea3ff` 2026-01-01T00:00:00Z
- **Target:** `fcb21f565aaeb3f3e13a0831ecbdc8552c4d38e0` (fcb21f565aaeb3f3e13a0831ecbdc8552c4d38e0) tree `97fc22c1ac8d` 2026-01-01T00:00:00Z
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `d5b3653d4ad9faa0` · tool 0.1.0-poc
- **Change:** 1 files, +1/-1 lines, 2 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | partial | 0 | structural redundancy 0→0 (Δ+0), exact 0→0; 0 near-structural candidate(s) |
| complexity | complexity | partial | 0 | modified callables net cyclomatic Δ+1; 0 added (cc 0), 1 deleted (cc 2); hotspots 0→0 |
| architecture_coupling | dependencies | partial | 0 | internal edges 0→0 (+0/-0); cycles 0→0; rule violations 0→0 |
| maintainability_scope | obsolete | partial | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | partial | 0 | 1 files changed (1 source); symbols +0/-1/~1 (renamed 0, formatting-only 0); exported API +0/-1/~0; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 0 test files changed, +0/-0 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 1 | 1 | 2 |  |
| typescript | no | 1 | 0 | 0 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
source 1 (+1/-1); 

## duplication — Callable duplication delta (`duplication`)

**Status:** partial · **Score:** 0 · **Rule:** net structural redundancy Δ+0: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 0→0 (Δ+0), exact 0→0; 0 near-structural candidate(s)

Coverage: 1 changed code file(s) in scope, 0 analysed, 1 with parse errors, unsupported {}, excluded by class {}, non-code 0

- callables with ≥40 tokens in production source are compared repository-wide (base 0, target 0)
- introduced 0 redundant instance(s), removed 0; 0 pre-existing group(s) touched but unchanged in size
- partial coverage: score reflects 0 analysed file(s); 0 unsupported (none) and 1 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 0 | 0 | 0 | – |
| structural_redundant_callables | repository | 0 | 0 | 0 | – |
| callables_considered | repository | 0 | 0 | 0 | – |
| redundant_instances_introduced | change | – | – | – | 0 |
| redundant_instances_removed | change | – | – | – | 0 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

## complexity — Callable complexity deltas (`complexity`)

**Status:** partial · **Score:** 0 · **Rule:** no hotspot change and \|modified net cyclomatic\| below threshold → 0

modified callables net cyclomatic Δ+1; 0 added (cc 0), 1 deleted (cc 2); hotspots 0→0

Coverage: 1 changed code file(s) in scope, 0 analysed, 1 with parse errors, unsupported {}, excluded by class {}, non-code 0

- modified callables: 1 increased, 0 decreased; net cyclomatic Δ+1
- added callables: 0 (cyclomatic 0); deleted callables: 1 (cyclomatic 2)
- changed-scope cyclomatic 3→2 (Δ-1), hotspots 0→0
- partial coverage: score reflects 0 analysed file(s); 0 unsupported (none) and 1 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 1 | 2 | +1 | – |
| decision_points_sum | modified_callables | 0 | 1 | +1 | – |
| lines_sum | modified_callables | 1 | 10 | +9 | – |
| max_nesting | modified_callables | 0 | 2 | +2 | – |
| cyclomatic_sum | changed_scope | 3 | 2 | -1 | – |
| max_cyclomatic | changed_scope | 2 | 2 | 0 | – |
| hotspot_callables | changed_scope | 0 | 0 | 0 | – |
| callables | changed_scope | 2 | 1 | -1 | – |
| callables_modified |  | – | – | – | 1 |
| callables_added |  | – | – | – | 0 |
| callables_deleted |  | – | – | – | 1 |
| cyclomatic_sum_added | added_callables | – | – | – | 0 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 2 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 0 |
| hotspots_added | added_callables | – | – | – | 0 |

Findings (1):

- **complexity.modified** [degradation; measured] Greet: cyclomatic 1→2, nesting 0→2, lines 1→10
  - base `558cf05221d5` util/util.go:8-8 `Greet` (before)
  - target `fcb21f565aae` util/util.go:8-17 `Greet` (after)

- diagnostic **warning** `metrics_in_error_region`: 1 changed callables overlap syntax errors; their metrics are unreliable

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** partial · **Score:** 0 · **Rule:** no cycles or rule violations changed → 0 (new edges alone are reported, not judged: features legitimately add dependencies)

internal edges 0→0 (+0/-0); cycles 0→0; rule violations 0→0

Coverage: 1 changed code file(s) in scope, 0 analysed, 1 with parse errors, unsupported {}, excluded by class {}, non-code 0

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod
- no architecture rules configured: boundary violations not assessed
- partial coverage: score reflects 0 analysed file(s); 0 unsupported (none) and 1 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| internal_units | repository | 1 | 1 | 0 | – |
| internal_edges | repository | 0 | 0 | 0 | – |
| external_dependency_refs | repository | 1 | 1 | 0 | – |
| cycles | repository | 0 | 0 | 0 | – |
| rule_violations | repository | 0 | 0 | 0 | – |
| edges_added | change | – | – | – | 0 |
| edges_removed | change | – | – | – | 0 |

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** partial · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 1 changed code file(s) in scope, 0 analysed, 1 with parse errors, unsupported {}, excluded by class {}, non-code 0

- partial coverage: score reflects 0 analysed file(s); 0 unsupported (none) and 1 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** partial · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

1 files changed (1 source); symbols +0/-1/~1 (renamed 0, formatting-only 0); exported API +0/-1/~0; 1 unit group(s)

Coverage: 1 changed code file(s) in scope, 0 analysed, 1 with parse errors, unsupported {}, excluded by class {}, non-code 0

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction
- partial coverage: score reflects 0 analysed file(s); 0 unsupported (none) and 1 with parse errors are not reflected

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | source | – | – | – | 1 |
| lines_added | source | – | – | – | 1 |
| lines_deleted | source | – | – | – | 1 |
| symbols_added | source | – | – | – | 0 |
| symbols_deleted | source | – | – | – | 1 |
| symbols_modified | source | – | – | – | 1 |
| symbols_renamed | source | – | – | – | 0 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 0 |
| exported_symbols_added | source | – | – | – | 0 |
| exported_symbols_removed | source | – | – | – | 1 |
| exported_symbols_changed | source | – | – | – | 0 |
| units_changed | source | – | – | – | 1 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

Findings (1):

- **scope.api_removed** [info; measured] exported symbol removed: Legacy
  - base `558cf05221d5` util/util.go:11-16 `Legacy` (before)

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 0 test files changed, +0/-0 lines; test churn is not evidence of correctness)

Coverage: 0 changed code file(s) in scope, 0 analysed, 0 with parse errors, unsupported {}, excluded by class {source:1}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 0 |
| test_lines_added | context_only | – | – | – | 0 |
| test_lines_deleted | context_only | – | – | – | 0 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| modified | function | Greet | util/util.go:8-17 | util/util.go:8-8 | source |
| deleted | function | Legacy → Legacy |  | util/util.go:11-16 | source |

## Diagnostics

- **error** `parse_error` target util/util.go:11 syntax error
- **error** `parse_error` target util/util.go:17 missing }

Snapshot-wide diagnostic counts: `target:parse_error`=2 

## Limitations

- POC: measurements of structure, not proof of design quality; lower metrics do not by themselves prove better architecture
- scores are per analyzer and directional; there is deliberately no weighted overall score
- correctness confidence is not assessed; the evaluated repository's code is never executed
- duplication, complexity and obsolete-code analyzers measure production source only; tests, generated, vendor and docs are counted in scope and coverage
- structural/near-structural similarity marks candidates, not semantic equivalence or proven missed reuse
- go: no type checking: references are matched by name, not resolved
- go: no call graph or interface dispatch resolution
- go: build tags and GOOS/GOARCH file variants are all indexed together
- go: import cycles are compile errors in Go; detected cycles indicate non-building code
- csharp: no adapter — files are counted as unsupported coverage, never scored
- java: no adapter — files are counted as unsupported coverage, never scored
- javascript: no adapter — files are counted as unsupported coverage, never scored
- kotlin: no adapter — files are counted as unsupported coverage, never scored
- python: no adapter — files are counted as unsupported coverage, never scored
- rust: no adapter — files are counted as unsupported coverage, never scored
- shell: no adapter — files are counted as unsupported coverage, never scored
- sql: no adapter — files are counted as unsupported coverage, never scored
- swift: no adapter — files are counted as unsupported coverage, never scored
- typescript: no adapter — files are counted as unsupported coverage, never scored
