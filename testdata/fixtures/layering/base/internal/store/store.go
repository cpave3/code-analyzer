package store

import "sort"

// Store is an in-memory key/value store.
type Store struct{ items map[string]int }

// New returns an empty store.
func New() *Store { return &Store{items: map[string]int{}} }

// Put sets a value.
func (s *Store) Put(k string, v int) { s.items[k] = v }

// Get reads a value.
func (s *Store) Get(k string) (int, bool) { v, ok := s.items[k]; return v, ok }

// Keys lists keys in order.
func (s *Store) Keys() []string {
	keys := make([]string, 0, len(s.items))
	for k := range s.items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
