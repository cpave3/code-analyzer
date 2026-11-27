package invoices

import "example.com/consolidate/money"

// Invoice is a bill.
type Invoice struct{ DueCents int64 }

// Due renders the amount due.
func (i Invoice) Due() string { return money.Format(i.DueCents) }
