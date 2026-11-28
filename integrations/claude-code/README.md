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

## Install into a project

From anywhere inside the project's git repository:

```bash
codeeval hook install --tool claude-code [--config codeeval.json] [--bin /path/to/codeeval] [--dry-run]
codeeval hook uninstall --tool claude-code
```

- **What it writes:** `.claude/hooks/codeeval-session-base.sh` and
  `.claude/hooks/codeeval-stop-analyze.sh`, both embedded in the binary.
- **Settings:** it registers them in `.claude/settings.json` under
  `SessionStart` and `Stop`. Other settings and hooks are kept, with
  top-level key order unchanged.
- **Running it again:** a reinstall updates in place and changes nothing
  when already current.
- **Refusals:** it never rewrites a `settings.json` that is not valid JSON.
- **Which binary:** the hook runs `codeeval` from `PATH` when it is there.
  Otherwise it pins this binary's path. `--bin` overrides either.
- **Uninstall:** removes only CodeEval's entries and scripts.

This is the only CodeEval command that writes into a repository, and only
the one you run it in.

## Analyze Stop hook (session-wide)

`hooks/codeeval-stop-analyze.sh` runs `codeeval analyze` instead of `check`.
That covers v2 debt, the defect rules and near-duplication. It compares the
working tree with the commit the session started from, which
`hooks/codeeval-session-base.sh` records at `SessionStart` in
`.git/codeeval-session-base`, so work the session already committed is
checked too.
- **Two tiers, so a stop stays fast** (`CODEEVAL_STOP_CHECKS`, default
  `tiered`):
  - *Quick, at every stop:* v2 debt and the single-file defect rules, on the
    changed files only (`analyze --checks quick`). On Legato (4,561 files)
    a stop takes about 1 s instead of about 72 s for the full analysis.
  - *Broad, in the background:* the cross-file defect rules (unused
    exports and private functions, unresolved imports, dependency and
    import-policy rules) and near-duplication (`--checks broad`, about 22 s
    on Legato). Each stop starts one, unless one is still running, and the
    next stop reports its result. That result describes the working tree at
    the earlier stop, and the message says so.
  - `quick` drops the background tier; `all` runs everything in the
    foreground, as before 0.8.5.
  - A new session stops a background run left from the previous one and
    clears its result. A result left when a session ends is never reported.
- **When it blocks:** on introduced gate defects (`CODEEVAL_DEFECTS_FAIL_ON`,
  default `gate`) or a debt rise (`CODEEVAL_DEBT_FAIL_ON`, default
  `net-debt`).
- **Escape hatches:** these keep the hook from trapping an agent in a fixing
  loop or forcing it to "fix" a false positive.
  - It blocks at most once per turn.
  - It never blocks twice in a session on the same failure. `analyze` prints
    a failure fingerprint of the failing items, and fingerprints already
    shown pass.
  - It blocks at most `CODEEVAL_HOOK_MAX_BLOCKS` times per session
    (default 3). A stop that reports both tiers counts as one block.
  - `CODEEVAL_HOOK=off` disables it.
  - False positives can be suppressed with a reasoned
    `codeeval-ignore <rule>: <reason>` comment, which reports list.
- **Untracked files:** they are checked by default (`--untracked`), except
  those Git ignores.

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
