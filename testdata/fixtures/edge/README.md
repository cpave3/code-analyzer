Edge cases evaluated against the same base:
- rename-only: util/util.go moved to util/strings.go unchanged → symbols "moved", all scores 0.
- rename-symbol: helper renamed to normalize, caller updated → renamed + modified, no duplication/complexity change.
- added-deleted: Legacy deleted, Classify added → added/deleted callables reported separately.
- parse-error: util/util.go has a syntax error → parse_error diagnostics, partial status.
- unsupported-only: only a Kotlin file changes (Kotlin has no adapter) → code analyzers "unsupported", score null (not neutral).
- mixed: Go and Kotlin change together → partial coverage, Go scored.
- odd-filename: adds "util/odd name ü.go" → path handled verbatim (NUL-separated Git output).
- (base..base): unchanged baseline → no changes, scores 0.
