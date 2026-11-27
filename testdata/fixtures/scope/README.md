Both variants fix the same billing bug (quantity ignored) and pass the same test.
- good: changes billing only.
- bad: also reformats and re-comments five functions in the unrelated report
  and notify packages (no behavioural change there).
Expected: scope 0 vs -1 (formatting-only churn in units unrelated to the fix).
Breadth by itself is never scored negative.
