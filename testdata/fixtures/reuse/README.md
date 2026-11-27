Both variants add tag normalisation with identical behaviour (same test).
- good: reuses the existing text.Slugify.
- bad: re-implements Slugify as a private copy with renamed identifiers.
Expected: duplication 0 vs -1 (structural clone of pre-existing text.Slugify).
