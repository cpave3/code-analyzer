package store

import (
	"fmt"
	"strings"
)

// Dump renders every entry as key=value lines.
func (s *Store) Dump() string {
	var lines []string
	for _, k := range s.Keys() {
		lines = append(lines, fmt.Sprintf("%s=%d", k, s.items[k]))
	}
	return strings.Join(lines, "\n")
}
