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
		if o.Subtotal > 500 {
			d += 0.05
		}
	}
	if o.Coupon != "" {
		if o.Coupon == "SAVE10" {
			d += 0.10
		} else if o.Coupon == "VIP" {
			if o.Member {
				d += 0.15
			}
		}
	}
	if o.Items >= 10 {
		if o.Items >= 50 {
			d += 0.10
		} else {
			d += 0.05
		}
	}
	if d > 0.30 {
		d = 0.30
	}
	return d
}
