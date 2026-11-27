package invoices

import "testing"

func TestDue(t *testing.T) {
	if got := (Invoice{DueCents: 100000000}).Due(); got != "$1,000,000.00" {
		t.Fatal(got)
	}
}
