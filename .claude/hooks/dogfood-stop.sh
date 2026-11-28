#!/usr/bin/env bash
# Dogfood: run the analyze Stop hook template on this repository with the
# already-built bin/codeeval (falling back to codeeval on PATH). It does not
# rebuild, so the evaluator never judges a change with that change's own code;
# rebuild deliberately (go build -o bin/codeeval ./cmd/codeeval) to update it.
set -u
repo=${CLAUDE_PROJECT_DIR:-$(pwd)}
bin="$repo/bin/codeeval"
[ -x "$bin" ] || bin=codeeval
CODEEVAL_BIN="$bin" exec "$repo/integrations/claude-code/hooks/codeeval-stop-analyze.sh"
