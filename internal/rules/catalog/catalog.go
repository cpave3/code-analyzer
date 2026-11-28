// Package catalog defines the language-neutral defect rules. Language
// implementations live in internal/rules/{golang,typescript,generic}.
package catalog

import "github.com/cameron/code-evaluator/internal/rules"

// Rule IDs.
const (
	UnusedImport          = "dead-code/unused-import"
	UnusedVariable        = "dead-code/unused-variable"
	UnreachableCode       = "dead-code/unreachable-code"
	UnusedPrivateFunction = "dead-code/unused-private-function"
	UnusedExport          = "dead-code/unused-export"

	EmptyHandler   = "lost-errors/empty-handler"
	ErrorToSuccess = "lost-errors/error-to-success"

	UnresolvedImport     = "dependencies/unresolved-import"
	UndeclaredDependency = "dependencies/undeclared-dependency"
	UnusedDependency     = "dependencies/unused-dependency"

	ImportPolicy = "architecture/import-policy"

	TautologicalAssertion = "test-oracle/tautological-assertion"
	LostAsyncAssertion    = "test-oracle/lost-async-assertion"
)

// Rules is the initial rule set.
var Rules = []rules.Rule{
	{ID: UnusedImport, Category: rules.CatDeadCode, Severity: rules.Gate,
		Summary: "An imported module or binding is never referenced in the file.",
		Limits:  []string{"name-based: a local declaration that shadows the import name counts as a use"}},
	{ID: UnusedVariable, Category: rules.CatDeadCode, Severity: rules.Gate,
		Summary: "A local variable is declared and never referenced in its enclosing function (or, for module-level TypeScript bindings, in its module).",
		Limits:  []string{"name-based within the enclosing function: a shadowing declaration with the same name that is used hides the defect", "destructuring patterns and loop variables are not checked", "names starting with _ are ignored"}},
	{ID: UnreachableCode, Category: rules.CatDeadCode, Severity: rules.Gate,
		Summary: "A statement follows return, throw, panic, break, continue or goto in the same block.",
		Limits:  []string{"only the block-local case: no constant-condition or whole-function reachability analysis", "process-exit calls (os.Exit, process.exit, log.Fatal) are not treated as terminators"}},
	{ID: UnusedPrivateFunction, Category: rules.CatDeadCode, Severity: rules.Gate,
		Summary: "A non-exported function (or private method) is never referenced where it is visible.",
		Limits:  []string{"name-based: any identifier with the same name in the visible scope counts as a use", "reflection, code generation and linker directives are not modelled"}},
	{ID: UnusedExport, Category: rules.CatDeadCode, Severity: rules.Advisory,
		Summary: "An exported function is not referenced from anywhere else in the repository.",
		Limits:  []string{"advisory: consumers outside the repository, frameworks and entry points are invisible", "name-based"}},

	{ID: EmptyHandler, Category: rules.CatLostErrors, Severity: rules.Gate,
		Summary: "A failure is caught and discarded with an empty handler (empty catch, empty .catch callback, empty err != nil branch, discarded recover()).",
		Limits:  []string{"a handler holding only a comment is reported as advisory (intent is documented but the failure is still dropped)", "no type information: dropped error results and floating promises are not detected"}},
	{ID: ErrorToSuccess, Category: rules.CatLostErrors, Severity: rules.Advisory,
		Summary: "A failure is converted into a success value without using the error (return nil / a default value from an error branch or catch).",
		Limits:  []string{"advisory: some fallbacks are intended; review whether the caller can tell the operation failed"}},

	{ID: UnresolvedImport, Category: rules.CatDependencies, Severity: rules.Gate,
		Summary: "An import of repository code resolves to no file or package in the snapshot.",
		Limits:  []string{"files generated at build time and not committed are reported as unresolved", "TypeScript: package exports maps and node_modules are not read"}},
	{ID: UndeclaredDependency, Category: rules.CatDependencies, Severity: rules.Gate,
		Summary: "An external package is imported but not declared in the governing manifest (go.mod, package.json).",
		Limits:  []string{"TypeScript: declarations in any package.json between the file and the repository root count (hoisted workspaces)", "Go: replace directives and go.work files are not modelled beyond modules present in the snapshot"}},
	{ID: UnusedDependency, Category: rules.CatDependencies, Severity: rules.Advisory,
		Summary: "A declared direct dependency is not imported by any code it governs.",
		Limits:  []string{"advisory: tools, plugins, CLIs and runtime-loaded packages are used without an import", "TypeScript: only \"dependencies\" are checked, not devDependencies"}},

	{ID: ImportPolicy, Category: rules.CatArchitecture, Severity: rules.Gate,
		Summary: "An internal dependency breaks a configured architecture rule (forbidden target, or a target outside the rule's allow-list).",
		Limits:  []string{"only explicit rules from the configuration are enforced; nothing is inferred", "unit granularity is the language's graph unit (directory); test files are not in the graph"}},

	{ID: TautologicalAssertion, Category: rules.CatTestOracle, Severity: rules.Gate,
		Summary: "A test assertion compares an expression with itself, or one literal with another, so it cannot fail.",
		Limits:  []string{"expressions containing calls are not compared (they may be deliberate determinism checks)", "only recognised assertion forms (testify, expect, node:assert)"}},
	{ID: LostAsyncAssertion, Category: rules.CatTestOracle, Severity: rules.Gate,
		Summary: "An asynchronous assertion can be lost: a promise assertion that is neither awaited nor returned, or a fatal test assertion called from a goroutine.",
		Limits:  []string{"syntactic: assertions passed through helper functions are not followed"}},
}

// Define adds the catalog to a registry.
func Define(r *rules.Registry) error {
	for _, rule := range Rules {
		if err := r.Define(rule); err != nil {
			return err
		}
	}
	return nil
}
