package cli

import "fmt"

// FormatKV renders a key/value pair for terminals.
func FormatKV(k string, v int) string { return fmt.Sprintf("%s=%d", k, v) }

// Render prints a lookup result.
func Render(key string, value int, found bool) string {
	if !found {
		return key + " (missing)"
	}
	return FormatKV(key, value)
}
