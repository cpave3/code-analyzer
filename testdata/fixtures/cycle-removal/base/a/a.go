package a

import "example.com/cyclerm/b"

// Name of a.
const Name = "alpha"

// Describe joins names.
func Describe() string { return Name + "/" + b.Name }
