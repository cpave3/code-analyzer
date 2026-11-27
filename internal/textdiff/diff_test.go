package textdiff

import (
	"reflect"
	"testing"
)

func TestLines(t *testing.T) {
	old := []byte("a\nb\nc\nd\ne\n")
	nw := []byte("a\nB\nc\nd\nx\ny\ne\n")
	r := Lines(old, nw)
	if r.Added != 3 || r.Deleted != 1 {
		t.Fatalf("added=%d deleted=%d", r.Added, r.Deleted)
	}
	if !reflect.DeepEqual(r.OldChanged, []Range{{2, 2}}) {
		t.Fatalf("old %v", r.OldChanged)
	}
	if !reflect.DeepEqual(r.NewChanged, []Range{{2, 2}, {5, 6}}) {
		t.Fatalf("new %v", r.NewChanged)
	}
	if r := Lines(old, old); r.Added+r.Deleted != 0 {
		t.Fatalf("identical: %+v", r)
	}
	if r := Lines(nil, old); r.Added != 5 || r.NewChanged[0] != (Range{1, 5}) {
		t.Fatalf("add all: %+v", r)
	}
}
