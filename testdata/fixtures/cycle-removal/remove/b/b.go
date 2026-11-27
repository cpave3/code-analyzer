package b

import "example.com/cyclerm/names"

// Name of b.
const Name = names.Beta

// Label labels b.
func Label() string { return Name + " of " + names.Alpha }
