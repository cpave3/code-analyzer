// Package integrations embeds the hook templates so `codeeval hook install`
// works from the binary alone.
package integrations

import "embed"

// Files holds the Claude Code and git hook templates.
//
//go:embed claude-code/hooks/codeeval-stop-analyze.sh claude-code/hooks/codeeval-session-base.sh git/post-commit git/commit-msg
var Files embed.FS
