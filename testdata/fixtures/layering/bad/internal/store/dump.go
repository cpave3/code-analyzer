package store

import (
	"strings"

	"example.com/layering/internal/cli"
)

// Dump renders every entry as key=value lines.
func (s *Store) Dump() string {
	var lines []string
	for _, k := range s.Keys() {
		lines = append(lines, cli.FormatKV(k, s.items[k]))
	}
	return strings.Join(lines, "\n")
}
