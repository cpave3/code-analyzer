# Claude Code integration (templates)

These files are templates. Nothing here is installed automatically.

## Skill

Copy `skills/codeeval-complexity/` into a project's `.claude/skills/` (or your
user skills directory). Claude then uses it when asked to check a PR for
complexity, reduce debt in a module, or respond to the hook below.

## Stop hook

`hooks/codeeval-stop-check.sh` runs `codeeval check --working-tree` against the
merge base when Claude is about to finish. If debt in scope rose (default
`--fail-on net-debt`), it blocks the stop and hands Claude the fix list.
- **No loops:** it exits 0 when `stop_hook_active` is set, so it blocks at
  most once per turn.
- **Never blocks on itself:** it also exits 0 when codeeval is not installed
  or fails.

Example `.claude/settings.json` entry:

```json
{
  "hooks": {
    "Stop": [
      { "hooks": [ { "type": "command",
        "command": "\"$CLAUDE_PROJECT_DIR\"/.claude/hooks/codeeval-stop-check.sh",
        "timeout": 120 } ] }
    ]
  }
}
```

Configure it with environment variables: `CODEEVAL_BIN`, `CODEEVAL_BASE`,
`CODEEVAL_CONFIG`, `CODEEVAL_FAIL_ON` (`new-debt`, `net-debt`, `gates` or
`none`), `CODEEVAL_PATHS` and `CODEEVAL_TOP`.

## Analyze Stop hook (session-wide)

`hooks/codeeval-stop-analyze.sh` runs `codeeval analyze` instead of `check`.
That covers v2 debt, the defect rules and near-duplication. It compares the
working tree with the commit the session started from, which
`hooks/codeeval-session-base.sh` records at `SessionStart` in
`.git/codeeval-session-base`, so work the session already committed is
checked too.
- **When it blocks:** on introduced gate defects (`CODEEVAL_DEFECTS_FAIL_ON`,
  default `gate`) or a debt rise (`CODEEVAL_DEBT_FAIL_ON`, default
  `net-debt`).
- **Untracked files:** they are not read until they are added, and the block
  message lists them.

This repository uses both hooks on itself: see `.claude/settings.json` and
`.claude/hooks/dogfood-stop.sh`. The dogfood script uses the existing
`bin/codeeval` and does not rebuild it, so a change is never judged by its own
evaluator code. Rebuild it on purpose to update it.

## CI / PR

The same check works in CI:

```bash
codeeval check --repo . --base "$BASE_SHA" --target "$HEAD_SHA" --fail-on net-debt \
  --summary /tmp/codeeval-check.md --output /tmp/codeeval-check.json
```

Exit 3 means the check failed. Post the Markdown as a PR comment, or feed the
JSON to an agent: each factor carries its location, explanation and
remediation. All of this is experimental scoring/v2. The numbers are
uncalibrated structural proxies, and the check never establishes correctness.
