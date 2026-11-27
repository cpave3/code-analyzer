package tags

import "example.com/reuse/text"

// Normalize turns user-entered tags into canonical, de-duplicated slugs.
func Normalize(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, t := range in {
		s := text.Slugify(t)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
