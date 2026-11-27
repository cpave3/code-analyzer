package tags

import (
	"strings"
	"unicode"
)

// Normalize turns user-entered tags into canonical, de-duplicated slugs.
func Normalize(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, t := range in {
		s := slugTag(t)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

func slugTag(tag string) string {
	var sb strings.Builder
	sep := false
	for _, ch := range strings.ToLower(tag) {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			sb.WriteRune(ch)
			sep = false
			continue
		}
		if !sep && sb.Len() > 0 {
			sb.WriteByte('-')
			sep = true
		}
	}
	return strings.TrimSuffix(sb.String(), "-")
}
