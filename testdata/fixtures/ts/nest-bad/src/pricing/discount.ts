export interface Order {
  total: number;
  member: boolean;
  coupon?: string;
}

export function discount(order: Order): number {
  let rate = 0;
  if (order.member) {
    rate = 0.05;
    if (order.total > 100) {
      if (order.coupon) {
        if (order.coupon === "VIP" || order.coupon === "GOLD") {
          rate = 0.2;
        } else {
          rate = 0.1;
        }
      } else {
        rate = 0.08;
      }
    }
  } else if (order.coupon && order.total > 50) {
    rate = order.coupon === "VIP" ? 0.1 : 0.03;
  }
  return Math.round(order.total * rate * 100) / 100;
}
