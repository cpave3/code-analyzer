# Agent provenance on commits

`codeeval attribution` groups per-commit quality deltas by who or what made
the commit, e.g. "Claude Code + Opus 5.5 with workflow v1" versus "workflow
v2". For that, commits must say what made them.

## Format

Plain git trailers in the commit message. No extension or custom object
format is needed, and they survive push, rebase and pull requests:

```
Fix pagination off-by-one

Agent-Model: claude-opus-5-5
Agent-Harness: claude-code/2.1.0
Agent-Workflow: proposal-stack@2
```

| Trailer | Meaning |
|---|---|
| `Agent-Model:` | **Standard.** The model, ideally its ID (`claude-opus-5-5`). It may repeat when several models contributed. A bare `Agent:` line (the 0.4–0.6 format) is still read as the model |
| `Agent-Harness:` | **Standard.** The tool that ran it, ideally `name/version` |
| `Agent-<Key>:` | **Any other metadata**: workflow, workflow version, session, effort, prompt pack, … Keys are case-insensitive and read as lower-case (`Agent-Workflow` → `workflow`) |

Rules the reader applies (`internal/provenance`):
- **Where it looks:** the lines are read anywhere in the message body, but
  never in the subject. A squash merge that concatenates commit messages
  therefore keeps every contributor's lines.
- **Notes:** commits that are already made can be annotated without
  rewriting them, with a note in `refs/notes/agent-provenance` holding the
  same lines:
  `git notes --ref=agent-provenance add -m "Agent-Model: x" <sha>`. Notes are not
  pushed by default. Share them with
  `git push origin refs/notes/agent-provenance`. Message values win over a
  note's, key by key.
- **Inference:** a commit with neither may be inferred from a
  `Co-Authored-By` line naming a known coding agent (Claude, Codex,
  Copilot, Gemini, …). That gives the agent only, labelled
  `co-author (inferred)`. The name is ID-styled ("Claude Opus 5.5" →
  `claude-opus-5-5`) so it groups with declared model IDs. Turn it off with
  `--infer-co-authors=false`.
- **Not verified:** all of it is self-reported text. It records what a commit
  claims, not what produced it.

## Writing it without touching messages: the `post-commit` hook (notes)

[`post-commit`](post-commit) reads the same variables (below) and attaches
the lines as a note in `refs/notes/agent-provenance` to each new commit.
Commit messages stay exactly as written, and the reader treats notes and
trailers the same way.

**Trade-offs** (these are why trailers are still offered):
- **Rewrites:** notes are tied to the commit SHA. Set
  `git config notes.rewriteRef refs/notes/agent-provenance` so amend and
  rebase carry them over. A squash or rebase merge on a forge creates new
  SHAs without notes, so the merged commit has no provenance unless
  something copies it.
- **Sharing:** notes are not pushed or fetched by default. Use
  `git push origin refs/notes/agent-provenance` and add a fetch refspec
  (`+refs/notes/agent-provenance:refs/notes/agent-provenance`).
- **Visibility:** forges do not display them.

Install it like the other hook. This repository dogfoods it: `.git/hooks/post-commit`
links to the template, and `.claude/settings.json` sets `AGENT_MODEL` and
`AGENT_HARNESS` for Claude Code sessions.

Other carriers were considered and rejected:
- extra commit-object headers: porcelain cannot write them, and rebase drops
  them;
- a ledger file in the repository: it adds churn and conflicts;
- an external database keyed by SHA: it loses provenance on every rewrite.

## Writing it: the `commit-msg` hook (trailers)

[`commit-msg`](commit-msg) appends trailers from environment variables, so a
harness or wrapper sets them once and every commit carries them:

| Variable | Trailer |
|---|---|
| `AGENT_MODEL` (or `AGENT_ID`) | `Agent-Model` |
| `AGENT_HARNESS` | `Agent-Harness` |
| `AGENT_META_<KEY>` | `Agent-<Key>` (`AGENT_META_WORKFLOW` → `Agent-Workflow`, `AGENT_META_SESSION_ID` → `Agent-Session-Id`) |

When none are set the hook does nothing, so human commits are unchanged. It
does not add a trailer that is already present. Install it in each
repository that should record provenance, or call it from your hook
manager:

```bash
cp integrations/git/commit-msg <repo>/.git/hooks/commit-msg && chmod +x <repo>/.git/hooks/commit-msg
```

For Claude Code, the variables can be set in the project's
`.claude/settings.json` `env` block. The model can change during a session
(`/model`), but the variables stay as set, so keep `AGENT_MODEL` in step with
the model in use. A harness that writes the trailers itself does not need
the hook. CodeEval only reads; it never installs hooks or writes to
evaluated repositories.

## Reading it

```bash
codeeval attribution --repo PATH --base A --target B --group-by model,harness,workflow [--defects]
```

The report has one row per group: first-parent steps, commits, changed
lines, v2 debt delta and its rate per changed kLOC, steps whose debt went up
or down, the decision-point delta, and with `--defects` the defect findings
introduced and resolved. A merge step counts every commit it brings in. A
merge commit without provenance is ignored for grouping, and a step whose
commits disagree is `(mixed)`.

The comparison is **observational**. Groups differ in tasks, size and
difficulty, so a difference between workflow v1 and v2 is a lead to
investigate, not a measured effect. Compare like with like (similar tasks,
enough steps), and use a controlled comparison (the same tasks run under
each workflow) before concluding that one degrades quality.
