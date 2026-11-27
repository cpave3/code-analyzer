package invoices

import (
	"fmt"

	"example.com/consolidate/money"
)

// Invoice is a bill.
type Invoice struct{ DueCents int64 }

// Due renders the amount due.
func (i Invoice) Due() string { return money.Format(i.DueCents) }

func formatAmount(cents int64) string {
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
