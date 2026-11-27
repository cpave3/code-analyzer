package billing

// Item is a line item.
type Item struct {
	Cents int
	Qty   int
}

// Total sums line items, honouring quantity.
func Total(items []Item) int {
	sum := 0
	for _, it := range items {
		sum += it.Cents * it.Qty
	}
	return sum
}
