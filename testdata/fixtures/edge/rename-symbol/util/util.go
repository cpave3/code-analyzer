package util

import "strings"

func normalize(s string) string { return strings.TrimSpace(strings.ToLower(s)) }

// Greet greets a person.
func Greet(name string) string { return "hello " + normalize(name) }

// Legacy is an old entry point.
func Legacy(n int) int {
	if n > 10 {
		return 10
	}
	return n
}
