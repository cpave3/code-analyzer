package rules

import (
	"bytes"
	"regexp"
	"strings"
)

// ignoreRe matches a suppression comment:
//
//	// codeeval-ignore dead-code/unused-export: kept for the plugin API
//	/* codeeval-ignore lost-errors: best-effort cleanup */
//
// The rule list is comma-separated IDs or categories; the text after ":" is
// the reason (recommended; reports show it).
var ignoreRe = regexp.MustCompile(`codeeval-ignore\s+([a-z0-9/_-]+(?:\s*,\s*[a-z0-9/_-]+)*)(?:\s*:\s*([^*\n]*))?`)

// Suppression reports whether a finding of rule at line is suppressed by a
// codeeval-ignore comment on that line or the line above, with its reason.
func Suppression(src []byte, line int, rule string) (string, bool) {
	if len(src) == 0 || line <= 0 || !bytes.Contains(src, []byte("codeeval-ignore")) {
		return "", false
	}
	lines := bytes.Split(src, []byte{'\n'})
	for _, l := range []int{line - 1, line - 2} {
		if l < 0 || l >= len(lines) {
			continue
		}
		m := ignoreRe.FindSubmatch(lines[l])
		if m == nil {
			continue
		}
		for _, id := range strings.Split(string(m[1]), ",") {
			if id = strings.TrimSpace(id); id != "" && (rule == id || strings.HasPrefix(rule, id+"/")) {
				reason := strings.TrimSpace(string(m[2]))
				if reason == "" {
					reason = "no reason given"
				}
				return reason, true
			}
		}
	}
	return "", false
}
