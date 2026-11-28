// Package hookinstall installs and removes CodeEval's agent hooks in a
// project. It is the one part of CodeEval that writes into a repository,
// and only when the user runs `codeeval hook install` there.
package hookinstall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cameron/code-evaluator/integrations"
)

// Tools that can be installed.
const ToolClaudeCode = "claude-code"

// Options configure an installation.
type Options struct {
	Root   string // project root (the git top level)
	Bin    string // codeeval binary for the hook ("" = codeeval on PATH)
	Config string // codeeval config file for the hook (absolute), optional
	DryRun bool
}

type script struct {
	name, event, embedded string
	timeout               int
}

var claudeScripts = []script{
	{"codeeval-session-base.sh", "SessionStart", "claude-code/hooks/codeeval-session-base.sh", 10},
	{"codeeval-stop-analyze.sh", "Stop", "claude-code/hooks/codeeval-stop-analyze.sh", 180},
}

// Change is one file the installer wrote, would write or removed.
type Change struct {
	Path   string `json:"path"`
	Action string `json:"action"` // created | updated | unchanged | removed
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

func command(sc script, o Options) string {
	var env []string
	if o.Bin != "" {
		env = append(env, "CODEEVAL_BIN="+shellQuote(o.Bin))
	}
	if o.Config != "" {
		env = append(env, "CODEEVAL_CONFIG="+shellQuote(o.Config))
	}
	cmd := `"$CLAUDE_PROJECT_DIR"/.claude/hooks/` + sc.name
	if len(env) > 0 {
		cmd = strings.Join(env, " ") + " " + cmd
	}
	return cmd
}

// settingsFile is .claude/settings.json with its top-level key order kept,
// so only the "hooks" value changes when the file is written back.
type settingsFile struct {
	order []string
	raw   map[string]json.RawMessage
	hooks map[string]any
}

// readSettings loads .claude/settings.json, refusing to rewrite a file it
// cannot parse.
func readSettings(path string) (*settingsFile, error) {
	sf := &settingsFile{raw: map[string]json.RawMessage{}, hooks: map[string]any{}}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) || (err == nil && len(bytes.TrimSpace(data)) == 0) {
		return sf, nil
	}
	if err != nil {
		return nil, err
	}
	bad := func(e error) error {
		return fmt.Errorf("%s is not valid JSON (%v); fix it or add the hooks by hand", path, e)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	if t, err := dec.Token(); err != nil || t != json.Delim('{') {
		return nil, bad(fmt.Errorf("want an object"))
	}
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, bad(err)
		}
		key, _ := t.(string)
		var v json.RawMessage
		if err := dec.Decode(&v); err != nil {
			return nil, bad(err)
		}
		if _, dup := sf.raw[key]; !dup {
			sf.order = append(sf.order, key)
		}
		sf.raw[key] = v
	}
	if h, ok := sf.raw["hooks"]; ok {
		if err := json.Unmarshal(h, &sf.hooks); err != nil {
			return nil, bad(fmt.Errorf("hooks: %v", err))
		}
	}
	return sf, nil
}

func ours(entry any) bool {
	h, ok := entry.(map[string]any)
	if !ok {
		return false
	}
	c, _ := h["command"].(string)
	for _, sc := range claudeScripts {
		if strings.Contains(c, ".claude/hooks/"+sc.name) {
			return true
		}
	}
	return false
}

// setHook adds or updates our hook for an event, leaving other hooks alone.
func setHook(hooks map[string]any, event string, hook map[string]any) {
	groups, _ := hooks[event].([]any)
	for _, g := range groups {
		gm, ok := g.(map[string]any)
		if !ok {
			continue
		}
		inner, _ := gm["hooks"].([]any)
		for i, h := range inner {
			if ours(h) {
				inner[i] = hook
				return
			}
		}
	}
	hooks[event] = append(groups, map[string]any{"hooks": []any{hook}})
}

// removeHooks drops our hooks, and any group or event left empty.
func removeHooks(hooks map[string]any) bool {
	removed := false
	for event, v := range hooks {
		groups, _ := v.([]any)
		var keptGroups []any
		for _, g := range groups {
			gm, ok := g.(map[string]any)
			if !ok {
				keptGroups = append(keptGroups, g)
				continue
			}
			inner, _ := gm["hooks"].([]any)
			var kept []any
			for _, h := range inner {
				if ours(h) {
					removed = true
					continue
				}
				kept = append(kept, h)
			}
			if len(kept) > 0 {
				gm["hooks"] = kept
				keptGroups = append(keptGroups, gm)
			}
		}
		if len(keptGroups) == 0 {
			delete(hooks, event)
		} else {
			hooks[event] = keptGroups
		}
	}
	return removed
}

func writeIfChanged(path string, data []byte, mode os.FileMode, dry bool) (Change, error) {
	ch := Change{Path: path, Action: "created"}
	if old, err := os.ReadFile(path); err == nil {
		if bytes.Equal(old, data) {
			if fi, err := os.Stat(path); err == nil && fi.Mode().Perm() == mode {
				return Change{Path: path, Action: "unchanged"}, nil
			}
		}
		ch.Action = "updated"
	}
	if dry {
		return ch, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return ch, err
	}
	if err := os.WriteFile(path, data, mode); err != nil {
		return ch, err
	}
	return ch, os.Chmod(path, mode)
}

// encode writes the settings back in their original key order; "hooks"
// is appended when new and dropped when empty.
func (sf *settingsFile) encode() ([]byte, error) {
	var hb bytes.Buffer
	enc := json.NewEncoder(&hb)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(sf.hooks); err != nil {
		return nil, err
	}
	order := sf.order
	if _, ok := sf.raw["hooks"]; !ok {
		order = append(order, "hooks")
	}
	var b bytes.Buffer
	b.WriteString("{")
	first := true
	for _, k := range order {
		v := sf.raw[k]
		if k == "hooks" {
			if len(sf.hooks) == 0 {
				continue
			}
			v = bytes.TrimSpace(hb.Bytes())
		}
		var ind bytes.Buffer
		if err := json.Indent(&ind, v, "  ", "  "); err != nil {
			return nil, err
		}
		key, _ := json.Marshal(k)
		if !first {
			b.WriteString(",")
		}
		first = false
		fmt.Fprintf(&b, "\n  %s: %s", key, ind.Bytes())
	}
	b.WriteString("\n}\n")
	return b.Bytes(), nil
}

// Install writes the Claude Code hooks into o.Root/.claude and registers
// them in .claude/settings.json. Running it again updates in place.
func Install(tool string, o Options) ([]Change, error) {
	if tool != ToolClaudeCode {
		return nil, fmt.Errorf("unknown tool %q (supported: %s)", tool, ToolClaudeCode)
	}
	settingsPath := filepath.Join(o.Root, ".claude", "settings.json")
	settings, err := readSettings(settingsPath)
	if err != nil {
		return nil, err
	}
	var changes []Change
	for _, sc := range claudeScripts {
		data, err := integrations.Files.ReadFile(sc.embedded)
		if err != nil {
			return nil, err
		}
		ch, err := writeIfChanged(filepath.Join(o.Root, ".claude", "hooks", sc.name), data, 0o755, o.DryRun)
		if err != nil {
			return nil, err
		}
		changes = append(changes, ch)
		setHook(settings.hooks, sc.event, map[string]any{"type": "command", "command": command(sc, o), "timeout": sc.timeout})
	}
	data, err := settings.encode()
	if err != nil {
		return nil, err
	}
	ch, err := writeIfChanged(settingsPath, data, 0o644, o.DryRun)
	if err != nil {
		return nil, err
	}
	return append(changes, ch), nil
}

// Uninstall removes the hook scripts and their settings entries; other
// settings and hooks are kept.
func Uninstall(tool string, o Options) ([]Change, error) {
	if tool != ToolClaudeCode {
		return nil, fmt.Errorf("unknown tool %q (supported: %s)", tool, ToolClaudeCode)
	}
	settingsPath := filepath.Join(o.Root, ".claude", "settings.json")
	settings, err := readSettings(settingsPath)
	if err != nil {
		return nil, err
	}
	var changes []Change
	if removeHooks(settings.hooks) {
		data, err := settings.encode()
		if err != nil {
			return nil, err
		}
		if !o.DryRun {
			if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
				return nil, err
			}
		}
		changes = append(changes, Change{Path: settingsPath, Action: "updated"})
	}
	for _, sc := range claudeScripts {
		p := filepath.Join(o.Root, ".claude", "hooks", sc.name)
		if _, err := os.Stat(p); err != nil {
			continue
		}
		if !o.DryRun {
			if err := os.Remove(p); err != nil {
				return nil, err
			}
		}
		changes = append(changes, Change{Path: p, Action: "removed"})
	}
	return changes, nil
}
