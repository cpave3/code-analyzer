#!/usr/bin/env bash
# Dogfood: rebuild codeeval from this checkout (falling back to the last good
# build when the working tree does not compile), then run the analyze Stop
# hook template on this repository.
set -u
repo=${CLAUDE_PROJECT_DIR:-$(pwd)}
input=$(cat)
export PATH="$PATH:/home/cameron/.local/share/mise/installs/go/1.26.0/bin"
(cd "$repo" && GOTOOLCHAIN=local GOPROXY=off GOFLAGS=-mod=readonly go build -o bin/codeeval.new ./cmd/codeeval) >/dev/null 2>&1 \
  && mv -f "$repo/bin/codeeval.new" "$repo/bin/codeeval"
[ -x "$repo/bin/codeeval" ] || exit 0
printf '%s' "$input" | CODEEVAL_BIN="$repo/bin/codeeval" "$repo/integrations/claude-code/hooks/codeeval-stop-analyze.sh"
