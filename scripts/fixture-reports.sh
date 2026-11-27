#!/usr/bin/env bash
# Rebuild the fixture repositories (outside this repo) and regenerate the
# paired fixture reports in examples/fixtures using the real CLI.
set -euo pipefail
cd "$(dirname "$0")/.."
command -v go >/dev/null || { echo "go not on PATH (toolchain setup, not a code failure). Try: export PATH=/home/cameron/.local/share/mise/installs/go/1.26.0/bin:\$PATH" >&2; exit 127; }
REPOS="${1:-$(mktemp -d -t codeeval-fixture-repos-XXXXXX)}"
go build -o bin/codeeval ./cmd/codeeval
go run ./cmd/codeeval-fixtures --codeeval ./bin/codeeval --repos "$REPOS" --out examples/fixtures
echo "fixture repos: $REPOS"
