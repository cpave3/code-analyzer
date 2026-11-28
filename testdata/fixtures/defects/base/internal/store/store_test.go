package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutGet(t *testing.T) {
	s := New()
	s.Put("a", "1")
	got, err := s.Get("a")
	assert.NoError(t, err)
	assert.Equal(t, "1", got)
	done := make(chan error)
	go func() {
		_, err := s.Get("missing")
		done <- err
	}()
	if err := <-done; err == nil {
		t.Fatal("want error")
	}
}
