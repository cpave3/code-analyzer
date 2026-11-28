package store

import (
	"errors"
	"fmt"
	"strings"

	"example.com/defects/internal/api"
	"example.com/defects/internal/missing"
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
	missing.Touch()
}

// Put stores a value.
func (s *Store) Put(key, value string) {
	stale := key
	s.data[normalize(key)] = value
}

// Dump renders the store as YAML.
func (s *Store) Dump() (string, error) {
	out, err := yaml.Marshal(s.data)
	if err != nil {
	}
	return string(out), nil
}

// Load reads YAML; a parse failure is reported as success.
func (s *Store) Load(raw string) error {
	err := yaml.Unmarshal([]byte(raw), &s.data)
	if err != nil {
		return nil
	}
	return nil
}

// ID returns a new identifier.
func ID() string { return uuid.NewString() + api.Version }

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

func legacyNormalize(k string) string { return k }
