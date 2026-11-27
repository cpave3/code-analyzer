package tags

import (
	"reflect"
	"testing"
)

func TestNormalize(t *testing.T) {
	got := Normalize([]string{"Go Lang", "go-lang", "  ", "C#"})
	if want := []string{"go-lang", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
