# CodeEval report: edge

- **Base:** `558cf05221d55f55ce1588532914292c021130ee` (558cf05221d55f55ce1588532914292c021130ee) tree `81e6e56ea3ff` 2026-01-01T00:00:00Z
- **Target:** `281bb2dbd746f91c785c61744b28abc58cc36412` (281bb2dbd746f91c785c61744b28abc58cc36412) tree `c109c778bea3` 2026-01-01T00:00:00Z
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `d5b3653d4ad9faa0` · tool 0.1.0-poc
- **Change:** 1 files, +2/-1 lines, 0 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | unsupported | n/a | unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1) |
| complexity | complexity | unsupported | n/a | unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1) |
| architecture_coupling | dependencies | unsupported | n/a | unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1) |
| maintainability_scope | obsolete | unsupported | n/a | unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1) |
| maintainability_scope | scope | unsupported | n/a | unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 0 test files changed, +0/-0 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 1 | 0 | 0 |  |
| typescript | no | 1 | 1 | 3 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
source 1 (+2/-1); 

## duplication — Callable duplication delta (`duplication`)

**Status:** unsupported · **Score:** n/a · **Rule:** no changed code in scope could be analysed → not scored (unsupported is not neutral)

unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1)

Coverage: 1 changed code file(s) in scope, 0 analysed, 0 with parse errors, unsupported {typescript:1}, excluded by class {}, non-code 0

- callables with ≥40 tokens in production source are compared repository-wide (base 0, target 0)
- introduced 0 redundant instance(s), removed 0; 0 pre-existing group(s) touched but unchanged in size

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

**Status:** unsupported · **Score:** n/a · **Rule:** no changed code in scope could be analysed → not scored (unsupported is not neutral)

unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1)

Coverage: 1 changed code file(s) in scope, 0 analysed, 0 with parse errors, unsupported {typescript:1}, excluded by class {}, non-code 0

- modified callables: 0 increased, 0 decreased; net cyclomatic Δ+0
- added callables: 0 (cyclomatic 0); deleted callables: 0 (cyclomatic 0)
- changed-scope cyclomatic 0→0 (Δ+0), hotspots 0→0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 0 | 0 | 0 | – |
| decision_points_sum | modified_callables | 0 | 0 | 0 | – |
| lines_sum | modified_callables | 0 | 0 | 0 | – |
| max_nesting | modified_callables | 0 | 0 | 0 | – |
| cyclomatic_sum | changed_scope | 0 | 0 | 0 | – |
| max_cyclomatic | changed_scope | 0 | 0 | 0 | – |
| hotspot_callables | changed_scope | 0 | 0 | 0 | – |
| callables | changed_scope | 0 | 0 | 0 | – |
| callables_modified |  | – | – | – | 0 |
| callables_added |  | – | – | – | 0 |
| callables_deleted |  | – | – | – | 0 |
| cyclomatic_sum_added | added_callables | – | – | – | 0 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 0 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 0 |
| hotspots_added | added_callables | – | – | – | 0 |

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** unsupported · **Score:** n/a · **Rule:** no changed code in scope could be analysed → not scored (unsupported is not neutral)

unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1)

Coverage: 1 changed code file(s) in scope, 0 analysed, 0 with parse errors, unsupported {typescript:1}, excluded by class {}, non-code 0

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod
- no architecture rules configured: boundary violations not assessed

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

**Status:** unsupported · **Score:** n/a · **Rule:** no changed code in scope could be analysed → not scored (unsupported is not neutral)

unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1)

Coverage: 1 changed code file(s) in scope, 0 analysed, 0 with parse errors, unsupported {typescript:1}, excluded by class {}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** unsupported · **Score:** n/a · **Rule:** no changed code in scope could be analysed → not scored (unsupported is not neutral)

unsupported: 1 changed file(s) in languages without the required capabilities (typescript:1)

Coverage: 1 changed code file(s) in scope, 0 analysed, 0 with parse errors, unsupported {typescript:1}, excluded by class {}, non-code 0

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | source | – | – | – | 1 |
| lines_added | source | – | – | – | 2 |
| lines_deleted | source | – | – | – | 1 |
| symbols_added | source | – | – | – | 0 |
| symbols_deleted | source | – | – | – | 0 |
| symbols_modified | source | – | – | – | 0 |
| symbols_renamed | source | – | – | – | 0 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 0 |
| exported_symbols_added | source | – | – | – | 0 |
| exported_symbols_removed | source | – | – | – | 0 |
| exported_symbols_changed | source | – | – | – | 0 |
| units_changed | source | – | – | – | 0 |
| disconnected_unit_groups | source | – | – | – | 0 |
| unrelated_churn_symbols | source | – | – | – | 0 |

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 0 test files changed, +0/-0 lines; test churn is not evidence of correctness)

Coverage: 0 changed code file(s) in scope, 0 analysed, 0 with parse errors, unsupported {}, excluded by class {source:1}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 0 |
| test_lines_added | context_only | – | – | – | 0 |
| test_lines_deleted | context_only | – | – | – | 0 |

## Diagnostics

- **info** `unsupported_language`  web/app.ts changed file in a language without an adapter is not analysed

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
