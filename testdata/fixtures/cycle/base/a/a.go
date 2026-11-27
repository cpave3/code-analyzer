package a

import "example.com/cycle/b"

// Name is a's display name.
const Name = "alpha"

// Describe joins both names.
func Describe() string { return Name + "/" + b.Name }
