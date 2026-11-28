// Package memo caches values derived from file content across runs, keyed
// by the content's blob id. A file that is unchanged between two commits
// has the same blob id, so analysing a series of commits only recomputes
// what each commit changed.
package memo

import "sync"

// Cache is a concurrency-safe map of derived values. Keys must include
// everything the value depends on besides the content itself (language,
// rule ID, analysis version ...). A nil *Cache caches nothing.
type Cache struct {
	mu           sync.Mutex
	m            map[string]any
	hits, misses int
}

// New returns an empty cache.
func New() *Cache { return &Cache{m: map[string]any{}} }

// Get returns build()'s value for key, computing and storing it once.
func (c *Cache) Get(key string, build func() any) any {
	if c == nil {
		return build()
	}
	c.mu.Lock()
	if v, ok := c.m[key]; ok {
		c.hits++
		c.mu.Unlock()
		return v
	}
	c.misses++
	c.mu.Unlock()
	v := build()
	c.mu.Lock()
	c.m[key] = v
	c.mu.Unlock()
	return v
}

// Stats reports entries, hits and misses.
func (c *Cache) Stats() (entries, hits, misses int) {
	if c == nil {
		return 0, 0, 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.m), c.hits, c.misses
}
