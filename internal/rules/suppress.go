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
	return suppressionAt(ignoreLines(src), line, rule)
}

// ignoreLines maps 1-based line numbers to the codeeval-ignore comments on
// them (nil when the file has none).
func ignoreLines(src []byte) map[int][][]byte {
	if !bytes.Contains(src, []byte("codeeval-ignore")) {
		return nil
	}
	out := map[int][][]byte{}
	for i, l := range bytes.Split(src, []byte{'\n'}) {
		if m := ignoreRe.FindSubmatch(l); m != nil {
			out[i+1] = m
		}
	}
	return out
}

// suppressionLines returns a file's ignore comments, cached by content.
func suppressionLines(f *File) map[int][][]byte {
	return f.Memo("codeeval-ignore", func() any { return ignoreLines(f.Source()) }).(map[int][][]byte)
}

func suppressionAt(lines map[int][][]byte, line int, rule string) (string, bool) {
	for _, l := range []int{line, line - 1} {
		m := lines[l]
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
