package billing

// Item is a line item.
type Item struct {
	Cents int
	Qty   int
}

// Total sums line items.
func Total(items []Item) int {
	sum := 0
	for _, it := range items {
		sum += it.Cents
	}
	return sum
}
