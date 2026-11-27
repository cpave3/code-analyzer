package api

import "example.com/layering/internal/store"

// Response is a lookup result.
type Response struct {
	Key   string
	Value int
	Found bool
}

// Lookup reads a key from the store.
func Lookup(s *store.Store, key string) Response {
	v, ok := s.Get(key)
	return Response{Key: key, Value: v, Found: ok}
}
