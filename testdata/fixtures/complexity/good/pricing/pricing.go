package pricing

import "math"

// Order is a priced order.
type Order struct {
	Items    int
	Subtotal float64
	Member   bool
	Coupon   string
}

const maxDiscount = 0.30

// Discount returns the discount fraction for an order.
func Discount(o Order) float64 {
	return math.Min(memberDiscount(o)+couponDiscount(o)+bulkDiscount(o.Items), maxDiscount)
}

func memberDiscount(o Order) float64 {
	if !o.Member {
		return 0
	}
	if o.Subtotal > 500 {
		return 0.10
	}
	return 0.05
}

type coupon struct {
	rate        float64
	membersOnly bool
}

var coupons = map[string]coupon{"SAVE10": {0.10, false}, "VIP": {0.15, true}}

func couponDiscount(o Order) float64 {
	c, ok := coupons[o.Coupon]
	if !ok || (c.membersOnly && !o.Member) {
		return 0
	}
	return c.rate
}

func bulkDiscount(items int) float64 {
	switch {
	case items >= 50:
		return 0.10
	case items >= 10:
		return 0.05
	}
	return 0
}
