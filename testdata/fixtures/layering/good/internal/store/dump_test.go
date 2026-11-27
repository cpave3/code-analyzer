package store

import "testing"

func TestDump(t *testing.T) {
	s := New()
	s.Put("b", 2)
	s.Put("a", 1)
	if got := s.Dump(); got != "a=1\nb=2" {
		t.Fatalf("%q", got)
	}
}
