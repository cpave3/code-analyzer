package util

import "strings"

// Trim trims spaces.
func Trim(s string) string { return strings.TrimSpace(s) }

// Shout upper-cases.
func Shout(s string) string { return strings.ToUpper(s) }
