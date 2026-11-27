Base has an a <-> b import cycle (does not build; Go forbids it).
- remove: both packages depend on a new leaf package names instead.
Expected: architecture +2 (cycle removed). Paired with cycle/introduce.
