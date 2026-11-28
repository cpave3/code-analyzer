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
	assert.Equal(t, got, got)
	done := make(chan struct{})
	go func() {
		if _, err := s.Get("missing"); err == nil {
			t.Fatal("want error")
		}
		close(done)
	}()
	<-done
}
