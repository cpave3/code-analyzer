package b

import "example.com/cyclerm/a"

// Name of b.
const Name = "beta"

// Label labels b.
func Label() string { return Name + " of " + a.Name }
