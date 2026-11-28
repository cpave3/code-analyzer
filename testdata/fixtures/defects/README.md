Defect rules (`codeeval defects`, rules/v1) in Go and TypeScript in one repository.
Not built by the Go test runner (`_nocompile.txt`): the bad variant does not compile by design, and the
fixtures have no go.sum or node_modules.
- base: clean code with idioms the rules must not flag (renamed package import, hoisted function after
  return, JSX with a default React import, closures, a used private method, an err branch that returns
  the error, an allowed import inside the lib allow-list).
- bad: one instance of every rule in each language (see the per-file comments). Every rule × language
  pair must report at least one finding, and the change must fail the default gate.
- clean: the same features as bad, written without the defects. No gate finding is introduced.
codeeval.json: store must not depend on api (forbid), web/src/lib may depend only on web/src/shared (allow-list),
svc/Domain (C#) must not depend on svc/Api (forbid).
C#: base is clean (using var, attributed test, awaited ThrowsAsync); bad fires every rule implemented for C#
(unused variable, unreachable, unused private method, empty catch, catch returning null, tautological
Assert.Equal, un-awaited Assert.ThrowsAsync, Domain -> Api dependency: Order calls OrdersController.Notify;
a bare unused `using Shop.Api;` would not count, since an unused using is not a dependency).
