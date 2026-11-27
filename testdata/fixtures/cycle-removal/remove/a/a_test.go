package a

import "testing"

func TestDescribe(t *testing.T) {
	if Describe() != "alpha/beta" {
		t.Fatal(Describe())
	}
}
