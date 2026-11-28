package store

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// ErrMissing is returned for unknown keys.
var ErrMissing = errors.New("missing")

// Store is an in-memory key/value store.
type Store struct{ data map[string]string }

// New returns an empty store.
func New() *Store { return &Store{data: map[string]string{}} }

// Get returns a value.
func (s *Store) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("get %q: %w", key, ErrMissing)
	}
	return v, nil
}

// Put stores a value.
func (s *Store) Put(key, value string) { s.data[normalize(key)] = value }

// Dump renders the store as YAML.
func (s *Store) Dump() (string, error) {
	out, err := yaml.Marshal(s.data)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Load reads YAML.
func (s *Store) Load(raw string) error {
	if err := yaml.Unmarshal([]byte(raw), &s.data); err != nil {
		return fmt.Errorf("load: %w", err)
	}
	return nil
}

// ID returns a new identifier.
func ID() string { return uuid.NewString() }

// Each visits entries until visit returns false.
func (s *Store) Each(visit func(k, v string) bool) {
	count := 0
	for k, v := range s.data {
		count++
		if !visit(k, v) {
			break
		}
	}
	_ = count
}

func normalize(k string) string {
	if k == "" {
		return "_"
	}
	return k
}
