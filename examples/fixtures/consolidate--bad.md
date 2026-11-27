# CodeEval report: consolidate

- **Base:** `354c9dc79dd525f1e39f17ed42d92455311577f5` (354c9dc79dd525f1e39f17ed42d92455311577f5) tree `f92bd4bf7ce9` 2026-01-01T00:00:00Z
- **Target:** `35b796f9def8f7797857c3bda3263d7d4e4e3160` (35b796f9def8f7797857c3bda3263d7d4e4e3160) tree `90bfcdbd0f86` 2026-01-01T00:00:00Z
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `d5b3653d4ad9faa0` · tool 0.1.0-poc
- **Change:** 3 files, +30/-4 lines, 3 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | ok | -1 | structural redundancy 1→2 (Δ+1), exact 1→2; 0 near-structural candidate(s) |
| complexity | complexity | ok | 0 | modified callables net cyclomatic Δ+0; 1 added (cc 3), 0 deleted (cc 0); hotspots 0→0 |
| architecture_coupling | dependencies | ok | 0 | internal edges 0→2 (+2/-0); cycles 0→0; rule violations 0→0 |
| maintainability_scope | obsolete | ok | -1 | 2 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | ok | 0 | 3 files changed (3 source); symbols +1/-0/~2 (renamed 0, formatting-only 0); exported API +1/-0/~0; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 0 test files changed, +0/-0 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 5 | 3 | 34 |  |

Changed files by class: 
source 3 (+30/-4); 

## duplication — Callable duplication delta (`duplication`)

**Status:** ok · **Score:** -1 · **Rule:** net structural redundancy Δ+1: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 1→2 (Δ+1), exact 1→2; 0 near-structural candidate(s)

Coverage: 3 changed code file(s) in scope, 3 analysed, 0 with parse errors, unsupported {}, excluded by class {}, non-code 0

- callables with ≥40 tokens in production source are compared repository-wide (base 2, target 3)
- introduced 1 redundant instance(s), removed 0; 0 pre-existing group(s) touched but unchanged in size

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 1 | 2 | +1 | – |
| structural_redundant_callables | repository | 1 | 2 | +1 | – |
| callables_considered | repository | 2 | 3 | +1 | – |
| redundant_instances_introduced | change | – | – | – | 1 |
| redundant_instances_removed | change | – | – | – | 0 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

Findings (1):

- **duplication.introduced_exact** [degradation; measured] identical callable body introduced: [Format]; matches pre-existing [formatAmount formatMoney] — possible reuse candidate
  - target `35b796f9def8` invoices/invoices.go:15-27 `formatAmount` (existing match)
  - target `35b796f9def8` money/money.go:6-18 `Format` (new or changed (added))
  - target `35b796f9def8` orders/orders.go:15-27 `formatMoney` (existing match)

## complexity — Callable complexity deltas (`complexity`)

**Status:** ok · **Score:** 0 · **Rule:** no hotspot change and \|modified net cyclomatic\| below threshold → 0

modified callables net cyclomatic Δ+0; 1 added (cc 3), 0 deleted (cc 0); hotspots 0→0

Coverage: 3 changed code file(s) in scope, 3 analysed, 0 with parse errors, unsupported {}, excluded by class {}, non-code 0

- modified callables: 0 increased, 0 decreased; net cyclomatic Δ+0
- added callables: 1 (cyclomatic 3); deleted callables: 0 (cyclomatic 0)
- changed-scope cyclomatic 2→5 (Δ+3), hotspots 0→0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 2 | 2 | 0 | – |
| decision_points_sum | modified_callables | 0 | 0 | 0 | – |
| lines_sum | modified_callables | 2 | 2 | 0 | – |
| max_nesting | modified_callables | 0 | 0 | 0 | – |
| cyclomatic_sum | changed_scope | 2 | 5 | +3 | – |
| max_cyclomatic | changed_scope | 1 | 3 | +2 | – |
| hotspot_callables | changed_scope | 0 | 0 | 0 | – |
| callables | changed_scope | 2 | 3 | +1 | – |
| callables_modified |  | – | – | – | 2 |
| callables_added |  | – | – | – | 1 |
| callables_deleted |  | – | – | – | 0 |
| cyclomatic_sum_added | added_callables | – | – | – | 3 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 0 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 0 |
| hotspots_added | added_callables | – | – | – | 0 |

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** ok · **Score:** 0 · **Rule:** no cycles or rule violations changed → 0 (new edges alone are reported, not judged: features legitimately add dependencies)

internal edges 0→2 (+2/-0); cycles 0→0; rule violations 0→0

Coverage: 3 changed code file(s) in scope, 3 analysed, 0 with parse errors, unsupported {}, excluded by class {}, non-code 0

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod
- no architecture rules configured: boundary violations not assessed

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| internal_units | repository | 2 | 3 | +1 | – |
| internal_edges | repository | 0 | 2 | +2 | – |
| external_dependency_refs | repository | 2 | 3 | +1 | – |
| cycles | repository | 0 | 0 | 0 | – |
| rule_violations | repository | 0 | 0 | 0 | – |
| edges_added | change | – | – | – | 2 |
| edges_removed | change | – | – | – | 0 |

Findings (5):

- **dependencies.edge_added** [info; measured] [go] new dependency invoices → money
  - target `35b796f9def8` invoices/invoices.go:6-6 (import)
- **dependencies.edge_added** [info; measured] [go] new dependency orders → money
  - target `35b796f9def8` orders/orders.go:6-6 (import)
- **dependencies.fan_change** [info; measured] [go] invoices fan-out 0→1, fan-in 0→0
- **dependencies.fan_change** [info; measured] [go] money fan-out 0→0, fan-in 0→2
- **dependencies.fan_change** [info; measured] [go] orders fan-out 0→1, fan-in 0→0

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** ok · **Score:** -1 · **Rule:** ≥1 non-exported function newly unreferenced but kept → -1

2 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 3 changed code file(s) in scope, 3 analysed, 0 with parse errors, unsupported {}, excluded by class {}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 2 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

Findings (2):

- **obsolete.left_unreferenced** [degradation; heuristic, confidence medium] formatAmount was referenced 1 time(s) in its unit before the change and has no production references after it, but was not removed
  - base `354c9dc79dd5` invoices/invoices.go:11-23 `formatAmount` (referenced before)
  - target `35b796f9def8` invoices/invoices.go:15-27 `formatAmount` (unreferenced after)
  - _caveat: name-based references; reflection, build tags, linkname and interface satisfaction are invisible_
- **obsolete.left_unreferenced** [degradation; heuristic, confidence medium] formatMoney was referenced 1 time(s) in its unit before the change and has no production references after it, but was not removed
  - base `354c9dc79dd5` orders/orders.go:11-23 `formatMoney` (referenced before)
  - target `35b796f9def8` orders/orders.go:15-27 `formatMoney` (unreferenced after)
  - _caveat: name-based references; reflection, build tags, linkname and interface satisfaction are invisible_

## maintainability_scope — Change scope and locality (`scope`)

**Status:** ok · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

3 files changed (3 source); symbols +1/-0/~2 (renamed 0, formatting-only 0); exported API +1/-0/~0; 1 unit group(s)

Coverage: 3 changed code file(s) in scope, 3 analysed, 0 with parse errors, unsupported {}, excluded by class {}, non-code 0

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | source | – | – | – | 3 |
| lines_added | source | – | – | – | 30 |
| lines_deleted | source | – | – | – | 4 |
| symbols_added | source | – | – | – | 1 |
| symbols_deleted | source | – | – | – | 0 |
| symbols_modified | source | – | – | – | 2 |
| symbols_renamed | source | – | – | – | 0 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 0 |
| exported_symbols_added | source | – | – | – | 1 |
| exported_symbols_removed | source | – | – | – | 0 |
| exported_symbols_changed | source | – | – | – | 0 |
| units_changed | source | – | – | – | 3 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

Findings (1):

- **scope.api_added** [info; measured] exported symbol added: Format
  - target `35b796f9def8` money/money.go:6-18 `Format` (after)

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 0 test files changed, +0/-0 lines; test churn is not evidence of correctness)

Coverage: 0 changed code file(s) in scope, 0 analysed, 0 with parse errors, unsupported {}, excluded by class {source:3}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 0 |
| test_lines_added | context_only | – | – | – | 0 |
| test_lines_deleted | context_only | – | – | – | 0 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| modified | method | Invoice.Due | invoices/invoices.go:13-13 | invoices/invoices.go:9-9 | source |
| added | function | Format | money/money.go:6-18 |  | source |
| modified | method | Order.Total | orders/orders.go:13-13 | orders/orders.go:9-9 | source |

## Diagnostics

No diagnostics for changed files.

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
