Both variants add the same new discount rules (member spend tier, members-only
VIP coupon, bulk tiers, higher cap) and pass the same test.
- good: extracts each rule into a small helper (table-driven coupons).
- bad: inlines the rules as nested conditionals.
Expected: complexity +1 (with an "extraction, not net removal" note) vs -1.
