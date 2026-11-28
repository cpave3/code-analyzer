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
	errs := make(chan error, 1)
	go func() {
		_, err := s.Get("missing")
		errs <- err
	}()
	if err := <-errs; err == nil {
		t.Fatal("want error")
	}
}
