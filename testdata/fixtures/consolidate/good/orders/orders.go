package orders

import "example.com/consolidate/money"

// Order is a customer order.
type Order struct{ Cents int64 }

// Total renders the order total.
func (o Order) Total() string { return money.Format(o.Cents) }
