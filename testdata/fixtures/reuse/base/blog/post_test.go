package blog

import "testing"

func TestURL(t *testing.T) {
	if got := (Post{Title: "Hello, World!"}).URL(); got != "/posts/hello-world" {
		t.Fatal(got)
	}
}
