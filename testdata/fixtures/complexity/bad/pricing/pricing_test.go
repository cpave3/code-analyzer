package pricing

import (
	"math"
	"testing"
)

func TestDiscount(t *testing.T) {
	cases := []struct {
		o    Order
		want float64
	}{
		{Order{}, 0},
		{Order{Member: true}, 0.05},
		{Order{Member: true, Subtotal: 600}, 0.10},
		{Order{Coupon: "SAVE10"}, 0.10},
		{Order{Coupon: "VIP"}, 0},
		{Order{Coupon: "VIP", Member: true}, 0.20},
		{Order{Items: 10}, 0.05},
		{Order{Items: 60}, 0.10},
		{Order{Items: 60, Member: true, Subtotal: 900, Coupon: "VIP"}, 0.30},
	}
	for _, c := range cases {
		if got := Discount(c.o); math.Abs(got-c.want) > 1e-9 {
			t.Errorf("%+v: got %v want %v", c.o, got, c.want)
		}
	}
}
