package orders

import "testing"

func TestTotal(t *testing.T) {
	if got := (Order{Cents: 123450}).Total(); got != "$1,234.50" {
		t.Fatal(got)
	}
	if got := (Order{Cents: -5}).Total(); got != "-$0.05" {
		t.Fatal(got)
	}
}
