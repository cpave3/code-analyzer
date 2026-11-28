#!/usr/bin/env bash
# Claude Code "Stop" hook: before Claude finishes, run `codeeval analyze` on
# everything changed since the session base (v2 debt, defect rules,
# near-duplication). When the change introduces gate defects or (by default)
# raises debt in scope, it exits 2 with the report head on stderr, which
# Claude Code feeds back to Claude.
#
# Base, first that exists: $CODEEVAL_BASE, <git-dir>/codeeval-session-base
# (written by codeeval-session-base.sh), merge-base with origin/main, HEAD.
#
# Environment (optional): CODEEVAL_BIN, CODEEVAL_BASE, CODEEVAL_CONFIG,
# CODEEVAL_DEBT_FAIL_ON (new-debt|net-debt|gates|none, default net-debt),
# CODEEVAL_DEFECTS_FAIL_ON (gate|any|none, default gate), CODEEVAL_PATHS.
set -u
input=$(cat)
if printf '%s' "$input" | grep -q '"stop_hook_active"[[:space:]]*:[[:space:]]*true'; then exit 0; fi
bin=${CODEEVAL_BIN:-codeeval}
command -v "$bin" >/dev/null 2>&1 || [ -x "$bin" ] || exit 0
repo=${CLAUDE_PROJECT_DIR:-$(pwd)}
cd "$repo" 2>/dev/null || exit 0
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || exit 0
gitdir=$(git rev-parse --git-dir)
base=${CODEEVAL_BASE:-}
if [ -z "$base" ] && [ -s "$gitdir/codeeval-session-base" ]; then base=$(cat "$gitdir/codeeval-session-base"); fi
git rev-parse -q --verify "${base:-x}^{commit}" >/dev/null 2>&1 || base=$(git merge-base HEAD origin/main 2>/dev/null || git rev-parse HEAD)
args=(analyze --repo "$repo" --base "$base" --working-tree --debt-fail-on "${CODEEVAL_DEBT_FAIL_ON:-net-debt}"
  --defects-fail-on "${CODEEVAL_DEFECTS_FAIL_ON:-gate}" --files 15 --top 10 --format md)
[ -n "${CODEEVAL_CONFIG:-}" ] && args+=(--config "$CODEEVAL_CONFIG")
[ -n "${CODEEVAL_PATHS:-}" ] && args+=(--path "$CODEEVAL_PATHS")
out=$("$bin" "${args[@]}" 2>/dev/null); code=$?
if [ "$code" -eq 3 ]; then
  {
    echo "codeeval analyze (base ${base:0:12} → working tree) failed: fix what this change introduced, or explain in your reply why it is justified."
    untracked=$(git ls-files --others --exclude-standard | head -20)
    [ -n "$untracked" ] && echo "Note: untracked files are not checked until added: $(echo $untracked)"
    echo
    printf '%s\n' "$out" | head -c 7000
  } >&2
  exit 2
fi
exit 0   # pass, or a tool error: never block on our own failure
