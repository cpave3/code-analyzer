package blog

import "example.com/reuse/text"

// Post is a blog post.
type Post struct{ Title string }

// URL is the post's canonical path.
func (p Post) URL() string { return "/posts/" + text.Slugify(p.Title) }
