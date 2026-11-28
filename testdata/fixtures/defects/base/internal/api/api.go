package api

import (
	"example.com/defects/internal/store"
	"example.com/defects/internal/util"
)

// Handler serves values.
type Handler struct{ S *store.Store }

// Lookup returns a trimmed value.
func (h Handler) Lookup(key string) (string, error) {
	v, err := h.S.Get(key)
	if err != nil {
		return "", err
	}
	return util.Trim(v), nil
}
