Base has the same money formatter copy-pasted into orders and invoices.
Both variants introduce a shared money.Format and route callers through it.
- good: deletes both copies.
- bad: leaves both now-unused copies behind.
Expected: duplication +1 vs -1; obsolete 0 vs -1 (two orphaned helpers).
