Both variants add Store.Dump() producing "key=value" lines (same test).
- good: formats locally with fmt.
- bad: imports the presentation package internal/cli to reuse its formatter,
  violating the configured rule "store must not depend on presentation".
Expected: architecture 0 vs -1 (rule violation with import-site evidence).
