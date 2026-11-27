# CodeEval report: legato-b3

- **Base:** `471bfd24a69628d985cdd7767dcb340da6a569b7` (471bfd24a69628d985cdd7767dcb340da6a569b7) tree `ae91f4192b23` 2026-11-27T15:46:34+11:00
- **Target:** `7b6de49760dfd4a71ed4110374612348e47a643a` (7b6de49760dfd4a71ed4110374612348e47a643a) tree `9f13299d5822` 2026-11-27T19:54:47+11:00
- **Mode:** commits · schema `codeeval.report/v1` · scoring/v1 · config v1 `3cef63633359769f` · tool 0.1.0-poc
- **Change:** 23 files, +227/-140 lines, 72 changed symbols

Scale: -2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.

## Scores by dimension (independent; no overall score)

| Dimension | Analyzer | Status | Score | Summary |
|---|---|---|---:|---|
| duplication | duplication | ok | 0 | structural redundancy 260→260 (Δ+0), exact 10→10; 0 near-structural candidate(s) |
| complexity | complexity | ok | 0 | modified callables net cyclomatic Δ+0; 2 added (cc 6), 2 deleted (cc 6); hotspots 5→5 |
| architecture_coupling | dependencies | ok | +1 | internal edges 200→196 (+1/-5); cycles 0→0; rule violations 15→12 |
| maintainability_scope | obsolete | ok | 0 | 0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed |
| maintainability_scope | scope | ok | 0 | 23 files changed (7 source); symbols +3/-3/~23 (renamed 6, formatting-only 0); exported API +2/-3/~11; 1 unit group(s) |
| correctness_confidence | correctness | not_assessed | n/a | not assessed (context: 11 test files changed, +122/-75 lines; test churn is not evidence of correctness) |

## Language coverage

| Language | Supported | Files in target | Changed files | Changed lines | Note |
|---|---|---:|---:|---:|---|
| go | yes | 2208 | 18 | 315 |  |
| javascript | no | 33 | 0 | 0 | unsupported: no adapter |
| kotlin | no | 57 | 0 | 0 | unsupported: no adapter |
| python | no | 10 | 0 | 0 | unsupported: no adapter |
| rust | no | 34 | 0 | 0 | unsupported: no adapter |
| shell | no | 15 | 0 | 0 | unsupported: no adapter |
| sql | no | 160 | 0 | 0 | unsupported: no adapter |
| swift | no | 2 | 0 | 0 | unsupported: no adapter |
| typescript | no | 999 | 0 | 0 | unsupported: planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths |

Changed files by class: 
docs 5 (+47/-5); source 7 (+58/-60); test 11 (+122/-75); 

## duplication — Callable duplication delta (`duplication`)

**Status:** ok · **Score:** 0 · **Rule:** net structural redundancy Δ+0: ≥1 → -1, ≥3 → -2, ≤-1 → +1, ≤-3 → +2; near-structural candidates are not scored

structural redundancy 260→260 (Δ+0), exact 10→10; 0 near-structural candidate(s)

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:11}, non-code 5

- callables with ≥40 tokens in production source are compared repository-wide (base 6762, target 6762)
- introduced 0 redundant instance(s), removed 0; 0 pre-existing group(s) touched but unchanged in size

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| exact_redundant_callables | repository | 10 | 10 | 0 | – |
| structural_redundant_callables | repository | 260 | 260 | 0 | – |
| callables_considered | repository | 6762 | 6762 | 0 | – |
| redundant_instances_introduced | change | – | – | – | 0 |
| redundant_instances_removed | change | – | – | – | 0 |
| preexisting_groups_touched | change | – | – | – | 0 |
| near_structural_candidates | change | – | – | – | 0 |

## complexity — Callable complexity deltas (`complexity`)

**Status:** ok · **Score:** 0 · **Rule:** no hotspot change and \|modified net cyclomatic\| below threshold → 0

modified callables net cyclomatic Δ+0; 2 added (cc 6), 2 deleted (cc 6); hotspots 5→5

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:11}, non-code 5

- modified callables: 0 increased, 0 decreased; net cyclomatic Δ+0
- added callables: 2 (cyclomatic 6); deleted callables: 2 (cyclomatic 6)
- changed-scope cyclomatic 94→94 (Δ+0), hotspots 5→5

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| cyclomatic_sum | modified_callables | 88 | 88 | 0 | – |
| decision_points_sum | modified_callables | 72 | 72 | 0 | – |
| lines_sum | modified_callables | 322 | 322 | 0 | – |
| max_nesting | modified_callables | 4 | 4 | 0 | – |
| cyclomatic_sum | changed_scope | 94 | 94 | 0 | – |
| max_cyclomatic | changed_scope | 15 | 15 | 0 | – |
| hotspot_callables | changed_scope | 5 | 5 | 0 | – |
| callables | changed_scope | 18 | 18 | 0 | – |
| callables_modified |  | – | – | – | 16 |
| callables_added |  | – | – | – | 2 |
| callables_deleted |  | – | – | – | 2 |
| cyclomatic_sum_added | added_callables | – | – | – | 6 |
| cyclomatic_sum_deleted | deleted_callables | – | – | – | 6 |
| hotspot_delta_preexisting | modified_and_deleted | – | – | – | 0 |
| hotspots_added | added_callables | – | – | – | 0 |

## architecture_coupling — Dependency graph delta (`dependencies`)

**Status:** ok · **Score:** +1 · **Rule:** rule violation resolved, or edges removed with none added → +1

internal edges 200→196 (+1/-5); cycles 0→0; rule violations 15→12

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:11}, non-code 5

- go: Go package imports (non-test files), internal packages resolved via nearest go.mod

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| internal_units | repository | 86 | 85 | -1 | – |
| internal_edges | repository | 200 | 196 | -4 | – |
| external_dependency_refs | repository | 1028 | 1025 | -3 | – |
| cycles | repository | 0 | 0 | 0 | – |
| rule_violations | repository | 15 | 12 | -3 | – |
| edges_added | change | – | – | – | 1 |
| edges_removed | change | – | – | – | 5 |

Findings (26):

- **dependencies.edge_added** [info; measured] [go] new dependency internal/tui/overlay → config
  - target `7b6de49760df` internal/tui/overlay/macro_picker.go:5-5 (import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed config → internal/engine/macros
  - base `471bfd24a696` config/config.go:14-14 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/server → internal/engine/macros
  - base `471bfd24a696` internal/server/macros.go:12-12 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/service → internal/engine/macros
  - base `471bfd24a696` internal/service/macros.go:12-12 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/tui → internal/engine/macros
  - base `471bfd24a696` internal/tui/app.go:16-16 (removed import)
- **dependencies.edge_removed** [info; measured] [go] dependency removed internal/tui/overlay → internal/engine/macros
  - base `471bfd24a696` internal/tui/overlay/macro_picker.go:9-9 (removed import)
- **dependencies.fan_change** [info; measured] [go] config fan-out 1→0, fan-in 6→7
- **dependencies.fan_change** [info; measured] [go] internal/engine/macros fan-out 0→0, fan-in 5→0
- **dependencies.fan_change** [info; measured] [go] internal/server fan-out 18→17, fan-in 2→2
- **dependencies.fan_change** [info; measured] [go] internal/service fan-out 48→47, fan-in 15→15
- **dependencies.fan_change** [info; measured] [go] internal/tui fan-out 19→18, fan-in 1→1
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/events (not scored)
  - target `7b6de49760df` internal/tui/app.go:15-15 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/ipc (not scored)
  - target `7b6de49760df` internal/tui/app.go:16-16 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/store (not scored)
  - target `7b6de49760df` internal/tui/app.go:17-17 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui → internal/engine/workflow (not scored)
  - target `7b6de49760df` internal/tui/app.go:18-18 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/board → internal/engine/store (not scored)
  - target `7b6de49760df` internal/tui/board/model.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/github (not scored)
  - target `7b6de49760df` internal/tui/overlay/link_pr.go:10-10 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/workflow (not scored)
  - target `7b6de49760df` internal/tui/overlay/workflow_trigger.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/overlay → internal/engine/store (not scored)
  - target `7b6de49760df` internal/tui/overlay/workspace.go:9-9 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/plan → internal/engine/store (not scored)
  - target `7b6de49760df` internal/tui/plan/model.go:13-13 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/report → internal/engine/analytics (not scored)
  - target `7b6de49760df` internal/tui/report/model.go:12-12 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/git (not scored)
  - target `7b6de49760df` internal/tui/review/messages.go:4-4 (import)
- **dependencies.rule_violation_preexisting** [info; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" already violated before the change: internal/tui/review → internal/engine/store (not scored)
  - target `7b6de49760df` internal/tui/review/model.go:17-17 (import)
- **dependencies.rule_violation_resolved** [improvement; measured] [go] rule "config is standalone, no internal imports (CLAUDE.md prose)" no longer violated: config → internal/engine/macros
  - base `471bfd24a696` config/config.go:14-14 (removed import)
- **dependencies.rule_violation_resolved** [improvement; measured] [go] rule "tui must not import engine (CLAUDE.md prose; NOT enforced by Legato architecture tests)" no longer violated: internal/tui → internal/engine/macros
  - base `471bfd24a696` internal/tui/app.go:16-16 (removed import)
- … 1 more in JSON

## maintainability_scope — Code left obsolete (`obsolete`)

**Status:** ok · **Score:** 0 · **Rule:** no newly orphaned functions (methods alone are low confidence and not scored) → 0

0 callable(s) newly unreferenced but kept; 0 unreferenced callable(s) removed

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:11}, non-code 5

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| callables_left_unreferenced | touched_units | – | – | – | 0 |
| unreferenced_callables_removed | touched_units | – | – | – | 0 |

## maintainability_scope — Change scope and locality (`scope`)

**Status:** ok · **Score:** 0 · **Rule:** no unrelated churn detected → 0 (breadth and API changes are reported, not judged; this analyzer never scores positive)

23 files changed (7 source); symbols +3/-3/~23 (renamed 6, formatting-only 0); exported API +2/-3/~11; 1 unit group(s)

Coverage: 7 changed code file(s) in scope, 7 analysed, 0 with parse errors, unsupported {}, excluded by class {test:11}, non-code 5

- affected callers are counted by identifier name (qualified by package name across packages), without type resolution
- a focused change scores 0: scope has no improvement direction

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| files_changed | docs | – | – | – | 5 |
| lines_added | docs | – | – | – | 47 |
| lines_deleted | docs | – | – | – | 5 |
| files_changed | source | – | – | – | 7 |
| lines_added | source | – | – | – | 58 |
| lines_deleted | source | – | – | – | 60 |
| files_changed | test | – | – | – | 11 |
| lines_added | test | – | – | – | 122 |
| lines_deleted | test | – | – | – | 75 |
| symbols_added | source | – | – | – | 3 |
| symbols_deleted | source | – | – | – | 3 |
| symbols_modified | source | – | – | – | 23 |
| symbols_renamed | source | – | – | – | 6 |
| symbols_moved | source | – | – | – | 0 |
| symbols_formatting_only | source | – | – | – | 0 |
| exported_symbols_added | source | – | – | – | 2 |
| exported_symbols_removed | source | – | – | – | 3 |
| exported_symbols_changed | source | – | – | – | 11 |
| units_changed | source | – | – | – | 5 |
| disconnected_unit_groups | source | – | – | – | 1 |
| unrelated_churn_symbols | source | – | – | – | 0 |

Findings (16):

- **scope.api_added** [info; measured] exported symbol added: Macro.Validate
  - target `7b6de49760df` config/macros.go:92-99 `Macro.Validate` (after)
- **scope.api_added** [info; measured] exported symbol added: SplitFavoriteMacros
  - target `7b6de49760df` config/macros.go:114-126 `SplitFavoriteMacros` (after)
- **scope.api_removed** [info; measured] exported symbol removed: Macro.Validate
  - base `471bfd24a696` internal/engine/macros/macros.go:91-98 `Macro.Validate` (before)
- **scope.api_removed** [info; measured] exported symbol removed: ListResult
  - base `471bfd24a696` internal/engine/macros/macros.go:101-103 `ListResult` (before)
- **scope.api_removed** [info; measured] exported symbol removed: SplitFavorites
  - base `471bfd24a696` internal/engine/macros/macros.go:118-130 `SplitFavorites` (before)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: Macro → Macro
  - base `471bfd24a696` internal/engine/macros/macros.go:21-37 `Macro` (before)
  - target `7b6de49760df` config/macros.go:22-38 `Macro` (after)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: Macro.UnmarshalYAML → Macro.UnmarshalYAML
  - base `471bfd24a696` internal/engine/macros/macros.go:63-74 `Macro.UnmarshalYAML` (before)
  - target `7b6de49760df` config/macros.go:64-75 `Macro.UnmarshalYAML` (after)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: Macro.UnmarshalJSON → Macro.UnmarshalJSON
  - base `471bfd24a696` internal/engine/macros/macros.go:77-88 `Macro.UnmarshalJSON` (before)
  - target `7b6de49760df` config/macros.go:78-89 `Macro.UnmarshalJSON` (after)
- **scope.api_renamed** [info; measured] exported symbol renamed/moved: AnyFavorite → AnyFavoriteMacro
  - base `471bfd24a696` internal/engine/macros/macros.go:106-113 `AnyFavorite` (before)
  - target `7b6de49760df` config/macros.go:102-109 `AnyFavoriteMacro` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: Server.SetMacros (m []macros.Macro) → (m []config.Macro)
  - base `471bfd24a696` internal/server/macros.go:24-26 `Server.SetMacros` (before)
  - target `7b6de49760df` internal/server/macros.go:23-25 `Server.SetMacros` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: BuiltinMacros () []macros.Macro → () []config.Macro
  - base `471bfd24a696` internal/service/macros.go:71-75 `BuiltinMacros` (before)
  - target `7b6de49760df` internal/service/macros.go:70-74 `BuiltinMacros` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: MacroService.Macros (ctx context.Context) ([]macros.Macro, error) → (ctx context.Context) ([]config.Macro, error)
  - base `471bfd24a696` internal/service/macros.go:218-228 `MacroService.Macros` (before)
  - target `7b6de49760df` internal/service/macros.go:217-227 `MacroService.Macros` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: MacroService.Save (ctx context.Context, name string, m macros.Macro, by string) error → (ctx context.Context, name string, m config.Macro, by string) error
  - base `471bfd24a696` internal/service/macros.go:248-268 `MacroService.Save` (before)
  - target `7b6de49760df` internal/service/macros.go:247-267 `MacroService.Save` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: NewApp (svc service.BoardService, syncSvc service.SyncService, agentSvc service.AgentService, prSvc service.PRTrackingService, reportSvc service.ReportService, icons theme.Icons, bus *events.Bus, editor string, workspaces []service.Workspace, tmux service.TmuxManager, workDir string, swarmSvc service.SwarmService, macrosList []macros.Macro, reviewServices ...review.Service) App → (svc service.BoardService, syncSvc service.SyncService, agentSvc service.AgentService, prSvc service.PRTrackingService, reportSvc service.ReportService, icons theme.Icons, bus *events.Bus, editor string, workspaces []service.Workspace, tmux service.TmuxManager, workDir string, swarmSvc service.SwarmService, macrosList []config.Macro, reviewServices ...review.Service) App
  - base `471bfd24a696` internal/tui/app.go:392-444 `NewApp` (before)
  - target `7b6de49760df` internal/tui/app.go:392-444 `NewApp` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: App.SetMacroSource (source func(context.Context) ([]macros.Macro, error)) → (source func(context.Context) ([]config.Macro, error))
  - base `471bfd24a696` internal/tui/app.go:2881-2883 `App.SetMacroSource` (before)
  - target `7b6de49760df` internal/tui/app.go:2881-2883 `App.SetMacroSource` (after)
- **scope.api_signature_changed** [info; measured] exported signature changed: NewMacroPicker (macrosList []macros.Macro) MacroPickerOverlay → (macrosList []config.Macro) MacroPickerOverlay
  - base `471bfd24a696` internal/tui/overlay/macro_picker.go:33-36 `NewMacroPicker` (before)
  - target `7b6de49760df` internal/tui/overlay/macro_picker.go:33-36 `NewMacroPicker` (after)

## correctness_confidence — Correctness confidence (`correctness`)

**Status:** not_assessed · **Score:** n/a · **Rule:** never scored in this POC: no independent correctness evidence is collected

not assessed (context: 11 test files changed, +122/-75 lines; test churn is not evidence of correctness)

Coverage: 11 changed code file(s) in scope, 11 analysed, 0 with parse errors, unsupported {}, excluded by class {source:7}, non-code 5

| Metric | Scope | Before | After | Δ | Value |
|---|---|---:|---:|---:|---:|
| test_files_changed | context_only | – | – | – | 11 |
| test_lines_added | context_only | – | – | – | 122 |
| test_lines_deleted | context_only | – | – | – | 75 |

## Changed symbols

| Status | Kind | Symbol | Location | Base location | Class |
|---|---|---|---|---|---|
| added | function | TestConfigImportsNoLegatoPackage | config/architecture_test.go:17-46 |  | test |
| modified | type | Config | config/config.go:17-87 | config/config.go:18-88 | source |
| modified | function | extractMacrosRaw | config/config.go:824-832 | config/config.go:825-833 | source |
| renamed | type | Macro | config/macros.go:22-38 | internal/engine/macros/macros.go:21-37 | source |
| renamed | type | macroFields | config/macros.go:41-50 | internal/engine/macros/macros.go:40-49 | source |
| renamed | method | macroFields.macro | config/macros.go:52-61 | internal/engine/macros/macros.go:51-60 | source |
| renamed | method | Macro.UnmarshalYAML | config/macros.go:64-75 | internal/engine/macros/macros.go:63-74 | source |
| renamed | method | Macro.UnmarshalJSON | config/macros.go:78-89 | internal/engine/macros/macros.go:77-88 | source |
| added | method | Macro.Validate | config/macros.go:92-99 |  | source |
| renamed | function | AnyFavorite → AnyFavoriteMacro | config/macros.go:102-109 | internal/engine/macros/macros.go:106-113 | source |
| added | function | SplitFavoriteMacros | config/macros.go:114-126 |  | source |
| renamed | function | TestMacroJSONRoundTrip | config/macros_test.go:10-23 | internal/engine/macros/macros_test.go:10-23 | test |
| renamed | function | TestMessageIsAnotherNameForKeys | config/macros_test.go:25-46 | internal/engine/macros/macros_test.go:51-72 | test |
| renamed | function | TestMessageAndKeysThatDifferAreRefused | config/macros_test.go:48-56 | internal/engine/macros/macros_test.go:74-82 | test |
| renamed | function | TestMacroJSONStillWritesKeys | config/macros_test.go:58-66 | internal/engine/macros/macros_test.go:84-92 | test |
| deleted | method | Macro.Validate → Macro.Validate |  | internal/engine/macros/macros.go:91-98 | source |
| deleted | type | ListResult → ListResult |  | internal/engine/macros/macros.go:101-103 | source |
| deleted | function | SplitFavorites → SplitFavorites |  | internal/engine/macros/macros.go:118-130 | source |
| deleted | function | TestListResultJSON → TestListResultJSON |  | internal/engine/macros/macros_test.go:25-49 | test |
| modified | method | Server.SetMacros | internal/server/macros.go:23-25 | internal/server/macros.go:24-26 | source |
| modified | method | Server.macrosHandler | internal/server/macros.go:33-63 | internal/server/macros.go:34-64 | source |
| added | type | macroListResult | internal/server/macros.go:66-68 |  | source |
| modified | type | macroWriter | internal/server/macros.go:77-83 | internal/server/macros.go:73-79 | source |
| modified | method | Server.macroEntryHandler | internal/server/macros.go:113-152 | internal/server/macros.go:109-148 | source |
| modified | function | TestMacrosHandlerEmpty | internal/server/macros_test.go:14-32 | internal/server/macros_test.go:14-32 | test |
| modified | function | TestMacrosHandlerUsesLowercaseJSONFields | internal/server/macros_test.go:34-53 | internal/server/macros_test.go:34-53 | test |
| modified | function | TestMacrosHandlerPopulated | internal/server/macros_test.go:55-82 | internal/server/macros_test.go:55-82 | test |
| modified | function | TestMacrosHandlerReadsTheMacroServiceAtEachRequest | internal/server/macros_test.go:124-153 | internal/server/macros_test.go:124-153 | test |
| added | function | TestMacroListResultJSON | internal/server/macros_test.go:156-180 |  | test |
| modified | type | Server | internal/server/server.go:28-151 | internal/server/server.go:28-151 | source |
| formatting_only | function | newMacroStartFixture | internal/service/macro_start_test.go:50-67 | internal/service/macro_start_test.go:50-65 | test |
| modified | function | TestAMacroStartsTheTasksAgentWithItsToolAndTier | internal/service/macro_start_test.go:69-93 | internal/service/macro_start_test.go:67-91 | test |
| modified | function | TestMacroAccountIsUsedForFreshWorkAndSavedIdentityForResume | internal/service/macro_start_test.go:182-213 | internal/service/macro_start_test.go:180-211 | test |
| modified | function | TestMacroAccountSelectsClaudeWithoutInheritingOtherToolProfile | internal/service/macro_start_test.go:224-236 | internal/service/macro_start_test.go:222-234 | test |
| modified | function | BuiltinMacros | internal/service/macros.go:70-74 | internal/service/macros.go:71-75 | source |
| modified | function | builtinMacro | internal/service/macros.go:76-83 | internal/service/macros.go:77-84 | source |
| modified | type | MacroEntry | internal/service/macros.go:86-95 | internal/service/macros.go:87-96 | source |
| modified | method | MacroService.Macros | internal/service/macros.go:217-227 | internal/service/macros.go:218-228 | source |
| modified | method | MacroService.Save | internal/service/macros.go:247-267 | internal/service/macros.go:248-268 | source |
| modified | function | macroValue | internal/service/macros.go:280-282 | internal/service/macros.go:281-283 | source |
| modified | method | MacroService.SetFavorite | internal/service/macros.go:286-314 | internal/service/macros.go:287-315 | source |
| modified | method | MacroService.Reset | internal/service/macros.go:388-415 | internal/service/macros.go:389-416 | source |
| modified | function | TestMacrosListBuiltInsThenTheFilesThenTheDatabases | internal/service/macros_test.go:31-67 | internal/service/macros_test.go:32-68 | test |
| modified | function | TestAnEditedBuiltInShadowsItAndResetRestoresIt | internal/service/macros_test.go:69-99 | internal/service/macros_test.go:70-100 | test |
| modified | function | TestABuiltInCannotBeDeleted | internal/service/macros_test.go:101-120 | internal/service/macros_test.go:102-121 | test |
| modified | function | TestAFileEntryIsPinnedEvenOverABuiltIn | internal/service/macros_test.go:122-140 | internal/service/macros_test.go:123-141 | test |
| modified | function | TestAMacroRefusesAnUnknownAppliesTo | internal/service/macros_test.go:142-151 | internal/service/macros_test.go:143-152 | test |
| modified | function | TestSettingsListsMacrosWithTheBuiltIns | internal/service/macros_test.go:153-178 | internal/service/macros_test.go:154-179 | test |
| modified | function | TestAMacroOfYoursCanBeRenamed | internal/service/macros_test.go:180-206 | internal/service/macros_test.go:181-207 | test |
| modified | function | TestARenameIsRefusedForBuiltInsFileMacrosAndNamesInUse | internal/service/macros_test.go:208-237 | internal/service/macros_test.go:209-238 | test |
| modified | function | TestAMacroEntryReadsBackWhole | internal/service/macros_test.go:261-274 | internal/service/macros_test.go:262-275 | test |
| modified | function | TestAStarSurvivesEditsAndAReset | internal/service/macros_test.go:309-359 | internal/service/macros_test.go:310-360 | test |
| modified | function | TestWorkflowMacroAccountReachesBothLaunchPaths | internal/service/workflow_agent_account_test.go:82-115 | internal/service/workflow_agent_account_test.go:83-116 | test |
| modified | function | TestAWorkflowAgentStepResolvesTheMacrosCurrentText | internal/service/workflow_macro_test.go:12-55 | internal/service/workflow_macro_test.go:12-55 | test |
| modified | function | TestAMacroOnAGoneFolderStartsInAFreshCloneOfTheBranch | internal/service/workflow_merge_wait_test.go:490-517 | internal/service/workflow_merge_wait_test.go:490-517 | test |
| modified | function | TestAMacroResumesInItsLiveFolder | internal/service/workflow_merge_wait_test.go:519-534 | internal/service/workflow_merge_wait_test.go:519-534 | test |
| modified | function | TestAMacroOnAGoneFolderWithNoRepositoryDoesNotStart | internal/service/workflow_merge_wait_test.go:536-554 | internal/service/workflow_merge_wait_test.go:536-554 | test |
| modified | type | App | internal/tui/app.go:168-271 | internal/tui/app.go:168-271 | source |
| modified | function | NewApp | internal/tui/app.go:392-444 | internal/tui/app.go:392-444 | source |
| modified | function | cardMacroItems | internal/tui/app.go:2837-2877 | internal/tui/app.go:2837-2877 | source |
| modified | method | App.SetMacroSource | internal/tui/app.go:2881-2883 | internal/tui/app.go:2881-2883 | source |
| modified | function | TestMacroSelectedClosesOverlay | internal/tui/app_test.go:2076-2089 | internal/tui/app_test.go:2076-2089 | test |
| modified | function | TestAPickedMacroIsSentAsAMessage | internal/tui/app_test.go:2094-2112 | internal/tui/app_test.go:2094-2112 | test |
| modified | function | macroEntry | internal/tui/card_context_macro_test.go:27-29 | internal/tui/card_context_macro_test.go:27-29 | test |
| modified | type | MacroSelectedMsg | internal/tui/overlay/macro_picker.go:14-16 | internal/tui/overlay/macro_picker.go:14-16 | source |
| modified | type | MacroPickerOverlay | internal/tui/overlay/macro_picker.go:23-30 | internal/tui/overlay/macro_picker.go:23-30 | source |
| modified | function | NewMacroPicker | internal/tui/overlay/macro_picker.go:33-36 | internal/tui/overlay/macro_picker.go:33-36 | source |
| modified | method | MacroPickerOverlay.Update | internal/tui/overlay/macro_picker.go:52-84 | internal/tui/overlay/macro_picker.go:52-84 | source |
| modified | function | TestMacroPickerNavigateAndSelect | internal/tui/overlay/macro_picker_test.go:19-52 | internal/tui/overlay/macro_picker_test.go:19-52 | test |
| modified | function | TestMacroPickerCancel | internal/tui/overlay/macro_picker_test.go:54-67 | internal/tui/overlay/macro_picker_test.go:54-67 | test |
| modified | function | TestMacroPickerShowsFavouritesAndTheRestBehindShowAll | internal/tui/overlay/macro_picker_test.go:69-100 | internal/tui/overlay/macro_picker_test.go:69-100 | test |
| modified | function | TestMacroPickerWithNoStarsListsEveryMacro | internal/tui/overlay/macro_picker_test.go:102-108 | internal/tui/overlay/macro_picker_test.go:102-108 | test |

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
