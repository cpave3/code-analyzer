export interface Order {
  total: number;
  member: boolean;
  coupon?: string;
}

export function discount(order: Order): number {
  let rate = 0;
  if (order.member) {
    rate = 0.05;
  }
  return Math.round(order.total * rate * 100) / 100;
}
