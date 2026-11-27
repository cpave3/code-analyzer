package b

import "example.com/cycle/a"

// Name is b's display name.
const Name = "beta"

// Label returns b's label.
func Label() string { return "[" + Name + "] of " + a.Name }
