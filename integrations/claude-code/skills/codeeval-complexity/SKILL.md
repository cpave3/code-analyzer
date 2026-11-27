---
name: codeeval-complexity
description: Check a change or a module for complexity, nesting, duplication and dependency-rule debt with codeeval, and fix the reported factors without gaming the metrics. Use when asked to check a PR or branch for added complexity, to reduce complexity/duplication in a module or path, or when a codeeval hook reports new debt.
---

# Checking and fixing complexity with codeeval

codeeval measures structure read-only (Go and TypeScript). It never runs code,
so behaviour must be protected by tests. Its numbers are an uncalibrated
structural proxy: use them to find and verify work, not as a target to game.

## Commands

Run from the repository root (outputs must go outside the repository, e.g. `/tmp`):

```bash
# PR / branch check: factors this change introduced or worsened (exit 3 = fails)
codeeval check --repo . --base "$(git merge-base HEAD origin/main)" --working-tree \
  --fail-on net-debt --top 15 --format md
# (use --target <sha> instead of --working-tree for a committed range;
#  --fail-on new-debt is stricter: any new or worsened factor fails)

# Debt factors in a module, application folder, path or glob, largest first
codeeval explain --repo . --rev HEAD --path internal/service --top 20 --format md
codeeval explain --repo . --rev HEAD --path 'web/src/**/*.tsx' --format json   # machine-readable
```

Add `--config <file>` when the project has a codeeval config (architecture rules,
exclusions). `--format json` gives a stable schema (`codeeval.check/v2`,
`codeeval.explain/v2`): each factor has `kind`, `location` (file, lines, symbol),
`metrics`, `debt`, `shareOfScopeDebt`, `explanation` and `remediation`
(`strategy`, `target`, `steps`, `acceptance`, `caution`).

## Workflow: fixing a PR check

1. Run `codeeval check`. Read the gate table and the "Fix these" list.
2. For each factor, largest `debtDelta` first:
   - Open the location. Understand what the code does and why it grew.
   - Make sure tests cover the behaviour you are about to touch; add
     characterisation tests first if they do not.
   - Apply the factor's remediation (see "How to fix" below).
3. Re-run `codeeval check`. The factor should be gone or reduced, no new
   factors should appear, and `decision points` for changed callables should
   not rise. Run the project's tests.
4. Stop when the check passes, or when the remaining factors are justified. In
   that case, write a short won't-fix note (factor, reason) in the PR
   description instead of forcing a change.

## Workflow: reducing debt in a module or path

1. `codeeval explain --path <scope>`. Note the concentration line: usually a
   few factors carry most of the debt.
2. Pick one factor (or one responsibility inside a large root). Keep each
   change small and separately reviewable.
3. Follow steps 2–4 above, using `explain` on the same scope to confirm the
   factor shrank and `check --working-tree` to confirm nothing new was added.

## How to fix each kind

- **complexity** (cyclomatic ≥ threshold): separate independent
  responsibilities into cohesive, well-named functions. Use guard clauses and
  early returns. Merge conditions with the same outcome. Remove dead branches.
  Keep every new helper below the threshold.
- **nesting**: invert conditions into guard clauses, and extract the
  innermost loop or branch body into a named function.
- **duplication**: confirm the copies really mean the same thing, then move
  one implementation to a module all callers may import (respecting
  architecture rules) and delete the copies. Structural matches may differ in
  ways that matter; if they do, or if layering forbids a shared home, record a
  won't-fix.
- **rule_violation**: check the rule still reflects intent, then route the
  dependency through a permitted layer (an interface owned by the importer,
  or a service) and remove the import.
- **cycle**: break the edge with the fewest import sites by moving or
  inverting the shared piece. TypeScript directory cycles are legal: treat
  them as coupling to review.

## Do not

- Do not split functions mechanically just to get under the threshold. If the
  branching only moves into helpers, `decision points` stay the same; the
  report shows this.
- Do not delete features, move logic into opaque data tables, or reshape
  syntax only to lower the number.
- Do not weaken or delete tests, or change behaviour, to make a refactor pass.
- Do not edit the codeeval config (thresholds, exclusions, rules) to make a
  check pass unless the user asks for that.
- Do not claim correctness from a passing check. It measures structure only.
