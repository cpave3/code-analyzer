package config

import "testing"

func TestMatch(t *testing.T) {
	cases := []struct {
		glob, path string
		want       bool
	}{
		{"internal/store/**", "internal/store", true},
		{"internal/store/**", "internal/store/sql", true},
		{"internal/store/**", "internal/storefront", false},
		{"**/gen/*.go", "a/b/gen/x.go", true},
		{"**/gen/*.go", "gen/x.go", true},
		{"*.md", "docs/a.md", false},
	}
	for _, c := range cases {
		if got := Match(c.glob, c.path); got != c.want {
			t.Errorf("Match(%q,%q)=%v want %v", c.glob, c.path, got, c.want)
		}
	}
}
