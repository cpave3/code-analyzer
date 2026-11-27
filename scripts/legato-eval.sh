#!/usr/bin/env bash
# Read-only headline evaluation of a pinned Legato window.
# Usage: scripts/legato-eval.sh [LEGATO_REPO] [CACHE_DIR]
# The source repo is only read through `git` plumbing (no checkout, fetch,
# hooks or scripts). Status is fingerprinted before and after without
# printing any filenames.
set -euo pipefail
cd "$(dirname "$0")/.."
command -v go >/dev/null || { echo "go not on PATH (toolchain setup, not a code failure). Try: export PATH=/home/cameron/.local/share/mise/installs/go/1.26.0/bin:\$PATH" >&2; exit 127; }
LEGATO="${1:-/home/cameron/Projects/legato}"
CACHE="${2:-$(mktemp -d -t codeeval-cache-XXXXXX)}"
BASE=9a61e9f8a021fa95a616966fe3977d4a5f6e0ccc   # last first-parent ancestor before 2026-11-26T00:00:00+11:00
TARGET=471bfd24a69628d985cdd7767dcb340da6a569b7 # local HEAD at 2026-11-27 16:44 AEDT (merge #431)
OUT=examples/legato
fingerprint() {
  git -C "$LEGATO" --no-optional-locks status --porcelain=v1 --untracked-files=all | sha256sum | cut -c1-16
}
go build -o bin/codeeval ./cmd/codeeval
before=$(fingerprint); head_before=$(git -C "$LEGATO" --no-optional-locks rev-parse HEAD)
./bin/codeeval history --repo "$LEGATO" --base $BASE --target $TARGET \
  --config $OUT/legato-architecture.json --cache-dir "$CACHE" --timing \
  --output $OUT/history-9a61e9f8..471bfd24.json --summary $OUT/history-9a61e9f8..471bfd24.md --format none
./bin/codeeval evaluate --repo "$LEGATO" --base $BASE --target $TARGET \
  --config $OUT/legato-architecture.json --cache-dir "$CACHE" --timing \
  --output $OUT/net-9a61e9f8..471bfd24.json --summary $OUT/net-9a61e9f8..471bfd24.md --format none
after=$(fingerprint); head_after=$(git -C "$LEGATO" --no-optional-locks rev-parse HEAD)
echo "status fingerprint before=$before after=$after; HEAD before=$head_before after=$head_after"
[ "$before" = "$after" ] && [ "$head_before" = "$head_after" ] || { echo "SOURCE STATE CHANGED" >&2; exit 1; }
