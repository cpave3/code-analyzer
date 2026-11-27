package orders

import (
	"fmt"

	"example.com/consolidate/money"
)

// Order is a customer order.
type Order struct{ Cents int64 }

// Total renders the order total.
func (o Order) Total() string { return money.Format(o.Cents) }

func formatMoney(cents int64) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	whole, frac := cents/100, cents%100
	digits := fmt.Sprintf("%d", whole)
	for i := len(digits) - 3; i > 0; i -= 3 {
		digits = digits[:i] + "," + digits[i:]
	}
	return fmt.Sprintf("%s$%s.%02d", sign, digits, frac)
}
