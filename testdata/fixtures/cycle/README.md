Base: package a imports b.
- introduce: b imports a back, creating a package cycle. Go rejects import
  cycles at compile time, so this variant intentionally does not build; it
  exercises the language-neutral cycle detector.
- decouple: a stops importing b (the constant is inlined), removing an edge.
Expected: architecture -2 (cycle introduced) vs +1 (edge removed, none added).
