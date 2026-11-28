package rules

import "testing"

func TestSuppression(t *testing.T) {
	src := []byte("package a\n\n// codeeval-ignore dead-code/unused-export, lost-errors: plugin API\nfunc A() {}\nfunc B() {} /* codeeval-ignore dead-code */\nfunc C() {}\n// codeeval-ignore test-oracle\nfunc D() {}\n")
	cases := []struct {
		line   int
		rule   string
		ok     bool
		reason string
	}{
		{4, "dead-code/unused-export", true, "plugin API"},                // comment on the line above
		{4, "lost-errors/empty-handler", true, "plugin API"},              // category in a list
		{4, "dead-code/unused-import", false, ""},                         // other rule in the category
		{5, "dead-code/unused-private-function", true, "no reason given"}, // same line, whole category
		{6, "dead-code/unused-private-function", true, "no reason given"}, // line above is line 5
		{8, "dead-code/unused-import", false, ""},                         // different category
		{8, "test-oracle/tautological-assertion", true, "no reason given"},
		{2, "dead-code/x", false, ""},
	}
	for _, c := range cases {
		reason, ok := Suppression(src, c.line, c.rule)
		if ok != c.ok || reason != c.reason {
			t.Errorf("line %d %s: got %v %q, want %v %q", c.line, c.rule, ok, reason, c.ok, c.reason)
		}
	}
	if _, ok := Suppression(nil, 1, "dead-code/x"); ok {
		t.Error("nil source suppressed")
	}
}
