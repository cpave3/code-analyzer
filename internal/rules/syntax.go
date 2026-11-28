package rules

import (
	"regexp"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// Helpers shared by language implementations.

// Walk visits n and its descendants in source order. Returning false from
// visit skips the node's children.
func Walk(n *sitter.Node, visit func(*sitter.Node) bool) {
	if n == nil || !visit(n) {
		return
	}
	for i := uint(0); i < n.ChildCount(); i++ {
		Walk(n.Child(i), visit)
	}
}

// Line is a node's 1-based start line.
func Line(n *sitter.Node) int { return int(n.StartPosition().Row) + 1 }

// Named returns a node's named children, without comments.
func Named(n *sitter.Node) []*sitter.Node {
	var out []*sitter.Node
	if n == nil {
		return nil
	}
	for i := uint(0); i < n.NamedChildCount(); i++ {
		c := n.NamedChild(i)
		if c.Kind() != "comment" {
			out = append(out, c)
		}
	}
	return out
}

// HasComment reports whether a node directly contains a comment.
func HasComment(n *sitter.Node) bool {
	for i := uint(0); i < n.NamedChildCount(); i++ {
		if n.NamedChild(i).Kind() == "comment" {
			return true
		}
	}
	return false
}

// Ancestor returns the nearest ancestor whose kind is in kinds.
func Ancestor(n *sitter.Node, kinds ...string) *sitter.Node {
	for p := n.Parent(); p != nil; p = p.Parent() {
		for _, k := range kinds {
			if p.Kind() == k {
				return p
			}
		}
	}
	return nil
}

// Contains reports whether any descendant of n (or n) has one of kinds.
func Contains(n *sitter.Node, kinds ...string) bool {
	found := false
	Walk(n, func(c *sitter.Node) bool {
		if found {
			return false
		}
		for _, k := range kinds {
			if c.Kind() == k {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

var space = regexp.MustCompile(`\s+`)

// Norm collapses whitespace and shortens text for keys and messages.
func Norm(s string) string {
	s = strings.TrimSpace(space.ReplaceAllString(s, " "))
	if len(s) > 80 {
		s = s[:77] + "..."
	}
	return s
}

// Same reports whether two nodes have identical text up to whitespace.
func Same(a, b *sitter.Node, src []byte) bool {
	if a == nil || b == nil {
		return false
	}
	return space.ReplaceAllString(a.Utf8Text(src), "") == space.ReplaceAllString(b.Utf8Text(src), "")
}

// Tautology explains why an assertion comparing a with b cannot fail: two
// literals, or the same expression on both sides (unless it contains one of
// effectKinds, e.g. calls, which may be a deliberate determinism check).
// It returns "" when the comparison can fail.
func Tautology(a, b *sitter.Node, src []byte, literal func(*sitter.Node) bool, effectKinds ...string) string {
	switch {
	case a == nil || b == nil:
		return ""
	case literal(a) && literal(b):
		return "compares two literals"
	case Same(a, b, src) && !Contains(a, effectKinds...):
		return "compares " + Norm(a.Utf8Text(src)) + " with itself"
	}
	return ""
}
