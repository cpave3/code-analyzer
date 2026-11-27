export interface Order {
  total: number;
  member: boolean;
  coupon?: string;
}

const premium = new Set(["VIP", "GOLD"]);

function memberRate(order: Order): number {
  if (order.total <= 100) return 0.05;
  if (!order.coupon) return 0.08;
  return premium.has(order.coupon) ? 0.2 : 0.1;
}

function guestRate(order: Order): number {
  if (!order.coupon || order.total <= 50) return 0;
  return order.coupon === "VIP" ? 0.1 : 0.03;
}

export function discount(order: Order): number {
  const rate = order.member ? memberRate(order) : guestRate(order);
  return Math.round(order.total * rate * 100) / 100;
}
