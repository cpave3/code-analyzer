# CodeEval report: legato-b6

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `8a2e03d7d9373b5e2295b9bcc019601dab138b87` (8a2e03d7d9373b5e2295b9bcc019601dab138b87) tree `9d18e5608c18` 2026-11-27T20:12:56+11:00
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `3cef63633359769f` · tool 0.1.0-poc
- **Change:** 2 files, +529/-50 lines, 22 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | ok | 0 | structural redundancy 260→260 (Δ+0), exact 10→10; 0 near-structural candidate(s) |
| complexity | complexity | ok | 0 | modified callables net cyclomatic Δ-20; 4 added (cc 26), 0 deleted (cc 0); hotspots 1→2 |
| architecture_coupling | dependencies | ok | 0 | internal edges 200→200 (+0/-0); cycles 0→0; rule violations 15→15 |
| maintainability_scope | obsolete | ok | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | ok | 0 | 2 files changed (1 source); symbols +4/-0/~1 (renamed 0, formatting-only 0); exported API +0/-0/~0; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 1 test files changed, +436/-0 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 2208 | 2 | 579 |  |
| javascript | no | 33 | 0 | 0 | unsupported: no adapter |
| kotlin | no | 57 | 0 | 0 | unsupported: no adapter |
| python | no | 10 | 0 | 0 | unsupported: no adapter |
| rust | no | 34 | 0 | 0 | unsupported: no adapter |
| shell | no | 15 | 0 | 0 | unsupported: no adapter |
| sql | no | 160 | 0 | 0 | unsupported: no adapter |
| swift | no | 2 | 0 | 0 | unsupported: no adapter |
| typescript | no | 999 | 0 | 0 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
source 1 (+93/-50); test 1 (+436/-0); 

## duplication — Callable duplication delta (`duplication`)

**Status:** ok · **Score:** 0 · **Rule:** net structural redundancy Δ+0: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 260→260 (Δ+0), exact 10→10; 0 near-structural candidate(s)

Coverage: 1 changed code file(s) in scope, 1 analysed, 0 with parse errors, unsupported {}, excluded by class {test:1}, non-code 0

- callables with ≥40 tokens in production source are compared repository-wide (base 6762, target 6765)
- introduced 0 redundant instance(s), removed 0; 0 pre-existing group(s) touched but unchanged in size

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 10 | 10 | 0 | – |
| structural_redundant_callables | repository | 260 | 260 | 0 | – |
| callables_considered | repository | 6762 | 6765 | +3 | – |
| redundant_instances_introduced | change | – | – | – | 0 |
| redundant_instances_removed | change | – | – | – | 0 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

## complexity — Callable complexity deltas (`complexity`)

**Status:** ok · **Score:** 0 · **Rule:** mixed: degradation and improvement signals both present → 0

modified callables net cyclomatic Δ-20; 4 added (cc 26), 0 deleted (cc 0); hotspots 1→2

Coverage: 1 changed code file(s) in scope, 1 analysed, 0 with parse errors, unsupported {}, excluded by class {test:1}, non-code 0

- modified callables: 0 increased, 1 decreased; net cyclomatic Δ-20
- added callables: 4 (cyclomatic 26); deleted callables: 0 (cyclomatic 0)
- changed-scope cyclomatic 48→54 (Δ+6), hotspots 1→2
- complexity removed from modified callables (Δ-20) is matched by 4 added callables (cyclomatic 26): this is extraction/redistribution, not net removal

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 48 | 28 | -20 | – |
| decision_points_sum | modified_callables | 47 | 27 | -20 | – |
| lines_sum | modified_callables | 127 | 81 | -46 | – |
| max_nesting | modified_callables | 3 | 2 | -1 | – |
| cyclomatic_sum | changed_scope | 48 | 54 | +6 | – |
| max_cyclomatic | changed_scope | 48 | 28 | -20 | – |
| hotspot_callables | changed_scope | 1 | 2 | +1 | – |
| callables | changed_scope | 1 | 5 | +4 | – |
| callables_modified |  | – | – | – | 1 |
| callables_added |  | – | – | – | 4 |
| callables_deleted |  | – | – | – | 0 |
| cyclomatic_sum_added | added_callables | – | – | – | 26 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 0 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 0 |
| hotspots_added | added_callables | – | – | – | 1 |

Findings (2):

- **complexity.added_hotspot** [degradation; measured] new callable agentService.resumeConversation has cyclomatic 12, nesting 2 (hotspot thresholds 10/4)
  - target `8a2e03d7d937` internal/service/agent_resume.go:180-206 `agentService.resumeConversation` (added)
- **complexity.modified** [improvement; measured] agentService.prepareTaskResume: cyclomatic 48→28, nesting 3→2, lines 127→81
  - base `471bfd24a696` internal/service/agent_resume.go:60-186 `agentService.prepareTaskResume` (before)
  - target `8a2e03d7d937` internal/service/agent_resume.go:60-140 `agentService.prepareTaskResume` (after)

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** ok · **Score:** 0 · **Rule:** no cycles or rule violations changed → 0 (new edges alone are reported, not judged: features legitimately add dependencies)

internal edges 200→200 (+0/-0); cycles 0→0; rule violations 15→15

Coverage: 1 changed code file(s) in scope, 1 analysed, 0 with parse errors, unsupported {}, excluded by class {test:1}, non-code 0

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| internal_units | repository | 86 | 86 | 0 | – |
| internal_edges | repository | 200 | 200 | 0 | – |
| external_dependency_refs | repository | 1028 | 1028 | 0 | – |
| cycles | repository | 0 | 0 | 0 | – |
| rule_violations | repository | 15 | 15 | 0 | – |
| edges_added | change | – | – | – | 0 |
| edges_removed | change | – | – | – | 0 |

Findings (15):

- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "config is standalone, no internal imports (CLAUDE.md prose)" already violated before the change: config → internal/engine/macros (not scored)
  - target `8a2e03d7d937` config/config.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/events (not scored)
  - target `8a2e03d7d937` internal/tui/app.go:14-14 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/ipc (not scored)
  - target `8a2e03d7d937` internal/tui/app.go:15-15 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/macros (not scored)
  - target `8a2e03d7d937` internal/tui/app.go:16-16 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/store (not scored)
  - target `8a2e03d7d937` internal/tui/app.go:17-17 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/workflow (not scored)
  - target `8a2e03d7d937` internal/tui/app.go:18-18 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/board → internal/engine/store (not scored)
  - target `8a2e03d7d937` internal/tui/board/model.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/github (not scored)
  - target `8a2e03d7d937` internal/tui/overlay/link_pr.go:10-10 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/macros (not scored)
  - target `8a2e03d7d937` internal/tui/overlay/macro_picker.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/workflow (not scored)
  - target `8a2e03d7d937` internal/tui/overlay/workflow_trigger.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/store (not scored)
  - target `8a2e03d7d937` internal/tui/overlay/workspace.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/plan → internal/engine/store (not scored)
  - target `8a2e03d7d937` internal/tui/plan/model.go:13-13 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/report → internal/engine/analytics (not scored)
  - target `8a2e03d7d937` internal/tui/report/model.go:12-12 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/git (not scored)
  - target `8a2e03d7d937` internal/tui/review/messages.go:4-4 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/store (not scored)
  - target `8a2e03d7d937` internal/tui/review/model.go:17-17 (import)

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** ok · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 1 changed code file(s) in scope, 1 analysed, 0 with parse errors, unsupported {}, excluded by class {test:1}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** ok · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

2 files changed (1 source); symbols +4/-0/~1 (renamed 0, formatting-only 0); exported API +0/-0/~0; 1 unit group(s)

Coverage: 1 changed code file(s) in scope, 1 analysed, 0 with parse errors, unsupported {}, excluded by class {test:1}, non-code 0

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | source | – | – | – | 1 |
| lines_added | source | – | – | – | 93 |
| lines_deleted | source | – | – | – | 50 |
| files_changed | test | – | – | – | 1 |
| lines_added | test | – | – | – | 436 |
| lines_deleted | test | – | – | – | 0 |
| symbols_added | source | – | – | – | 4 |
| symbols_deleted | source | – | – | – | 0 |
| symbols_modified | source | – | – | – | 1 |
| symbols_renamed | source | – | – | – | 0 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 0 |
| exported_symbols_added | source | – | – | – | 0 |
| exported_symbols_removed | source | – | – | – | 0 |
| exported_symbols_changed | source | – | – | – | 0 |
| units_changed | source | – | – | – | 1 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 1 test files changed, +436/-0 lines; test churn is not evidence of correctness)

Coverage: 1 changed code file(s) in scope, 1 analysed, 0 with parse errors, unsupported {}, excluded by class {source:1}, non-code 0

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 1 |
| test_lines_added | context_only | – | – | – | 436 |
| test_lines_deleted | context_only | – | – | – | 0 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| modified | method | agentService.prepareTaskResume | internal/service/agent_resume.go:60-140 | internal/service/agent_resume.go:60-186 | source |
| added | method | agentService.inferResumeKind | internal/service/agent_resume.go:146-159 |  | source |
| added | function | resumeStorage | internal/service/agent_resume.go:164-172 |  | source |
| added | method | agentService.resumeConversation | internal/service/agent_resume.go:180-206 |  | source |
| added | method | agentService.bindResumeAccount | internal/service/agent_resume.go:213-229 |  | source |
| added | type | prepFixture | internal/service/agent_resume_prepare_test.go:25-31 |  | test |
| added | function | newPrepFixture | internal/service/agent_resume_prepare_test.go:35-57 |  | test |
| added | method | prepFixture.save | internal/service/agent_resume_prepare_test.go:61-85 |  | test |
| added | method | prepFixture.prepare | internal/service/agent_resume_prepare_test.go:87-90 |  | test |
| added | method | prepFixture.defaultIdentity | internal/service/agent_resume_prepare_test.go:93-100 |  | test |
| added | type | resumeFields | internal/service/agent_resume_prepare_test.go:103-106 |  | test |
| added | function | fieldsOf | internal/service/agent_resume_prepare_test.go:108-110 |  | test |
| added | function | wantFields | internal/service/agent_resume_prepare_test.go:112-120 |  | test |
| added | function | wantRefusal | internal/service/agent_resume_prepare_test.go:122-130 |  | test |
| added | function | TestPrepareResumeLeavesOptionsAloneWithoutAResumePolicy | internal/service/agent_resume_prepare_test.go:134-143 |  | test |
| added | function | TestPrepareResumeNewPolicyChecksOnlyTheAdapter | internal/service/agent_resume_prepare_test.go:145-160 |  | test |
| added | function | TestPrepareResumeRefusesAnUnsupportedPolicy | internal/service/agent_resume_prepare_test.go:162-166 |  | test |
| added | function | TestPrepareResumeRefusesAToolThatIsNotConfigured | internal/service/agent_resume_prepare_test.go:168-180 |  | test |
| added | function | TestPrepareResumeInfersTheToolFromTheNewestSavedSession | internal/service/agent_resume_prepare_test.go:182-213 |  | test |
| added | function | TestPrepareResumeSelectsTheSavedOrAnEarlierConversation | internal/service/agent_resume_prepare_test.go:215-289 |  | test |
| added | function | TestPrepareResumeChecksTheSavedDirectory | internal/service/agent_resume_prepare_test.go:291-362 |  | test |
| added | function | TestPrepareResumeBindsTheClaudeAccountTheConversationRanWith | internal/service/agent_resume_prepare_test.go:364-436 |  | test |

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
