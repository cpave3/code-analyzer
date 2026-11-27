#!/usr/bin/env bash
# Rebuild the fixture repositories (outside this repo) and regenerate the
# paired fixture reports in examples/fixtures using the real CLI.
set -euo pipefail
cd "$(dirname "$0")/.."
REPOS="${1:-$(mktemp -d -t codeeval-fixture-repos-XXXXXX)}"
go build -o bin/codeeval ./cmd/codeeval
go run ./cmd/codeeval-fixtures --codeeval ./bin/codeeval --repos "$REPOS" --out examples/fixtures
echo "fixture repos: $REPOS"
