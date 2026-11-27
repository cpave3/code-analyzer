#!/usr/bin/env bash
# Claude Code "Stop" hook: before Claude finishes, check the working tree for
# complexity/duplication/dependency debt introduced against the merge base.
# On failure it exits 2 and prints the fix list on stderr, which Claude Code
# feeds back to Claude so it can address the factors (or justify them).
#
# Environment (all optional):
#   CODEEVAL_BIN      path to codeeval (default: codeeval on PATH)
#   CODEEVAL_BASE     base ref to diff against (default: merge-base with origin/main, else HEAD)
#   CODEEVAL_CONFIG   codeeval config file
#   CODEEVAL_FAIL_ON  new-debt | net-debt | gates | none (default: net-debt)
#   CODEEVAL_PATHS    comma-separated scope (default: whole repository)
#   CODEEVAL_TOP      factors to report (default: 10)
set -u
input=$(cat)
# Avoid loops: if this Stop was already triggered by a previous block, let it finish.
if printf '%s' "$input" | grep -q '"stop_hook_active"[[:space:]]*:[[:space:]]*true'; then exit 0; fi

bin=${CODEEVAL_BIN:-codeeval}
command -v "$bin" >/dev/null 2>&1 || exit 0   # not installed: do nothing
repo=${CLAUDE_PROJECT_DIR:-$(pwd)}
cd "$repo" 2>/dev/null || exit 0
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || exit 0
base=${CODEEVAL_BASE:-$(git merge-base HEAD origin/main 2>/dev/null || git rev-parse HEAD)}
args=(check --repo "$repo" --base "$base" --working-tree --fail-on "${CODEEVAL_FAIL_ON:-net-debt}" --top "${CODEEVAL_TOP:-10}" --format md)
[ -n "${CODEEVAL_CONFIG:-}" ] && args+=(--config "$CODEEVAL_CONFIG")
[ -n "${CODEEVAL_PATHS:-}" ] && args+=(--path "$CODEEVAL_PATHS")

out=$("$bin" "${args[@]}" 2>/dev/null); code=$?
if [ "$code" -eq 3 ]; then
  {
    echo "codeeval found complexity/duplication/dependency debt added by your changes."
    echo "Address these factors (see the codeeval-complexity skill), or explain in your reply why they are justified:"
    echo
    printf '%s\n' "$out" | sed -n '/^## Fix these/,/^## Improved by this change/p' | sed '$d' | head -c 6000
  } >&2
  exit 2
fi
exit 0   # pass, advisory, or tool error: never block on our own failure
