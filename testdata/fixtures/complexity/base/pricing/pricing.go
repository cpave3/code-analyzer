package pricing

// Order is a priced order.
type Order struct {
	Items    int
	Subtotal float64
	Member   bool
	Coupon   string
}

// Discount returns the discount fraction for an order.
func Discount(o Order) float64 {
	d := 0.0
	if o.Member {
		d += 0.05
	}
	if o.Coupon == "SAVE10" {
		d += 0.10
	}
	if d > 0.25 {
		d = 0.25
	}
	return d
}
