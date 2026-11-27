TypeScript paired scenarios (TSX grammar, relative imports). Behaviour is not
executed by the Go test runner (`_nocompile.txt`); variants are written to be
behaviour-equivalent pairs.
- reuse-good: tags reuse the existing slugify via a relative import.
- reuse-bad: tags re-implement slugify as a private copy with renamed identifiers → structural clone.
- nest-bad: discount gains nested conditionals in place → complexity growth in an existing function.
- extract-good: the same discount rules expressed through small helpers → root shrinks, helpers added.
- cycle: text imports from blog while blog imports text → directory cycle introduced.
