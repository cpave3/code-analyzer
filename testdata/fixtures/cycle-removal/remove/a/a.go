package a

import "example.com/cyclerm/names"

// Name of a.
const Name = names.Alpha

// Describe joins names.
func Describe() string { return Name + "/" + names.Beta }
