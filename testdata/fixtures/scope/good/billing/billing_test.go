package billing

import "testing"

func TestTotal(t *testing.T) {
	if got := Total([]Item{{Cents: 250, Qty: 2}, {Cents: 100, Qty: 1}}); got != 600 {
		t.Fatal(got)
	}
}
