package util

import "strings"

func helper(s string) string { return strings.TrimSpace(strings.ToLower(s)) }

// Greet greets a person.
func Greet(name string) string { return "hello " + helper(name) }

// Classify buckets a number.
func Classify(n int) string {
	switch {
	case n < 0:
		return "negative"
	case n == 0:
		return "zero"
	case n < 10:
		return "small"
	}
	return "large"
}
