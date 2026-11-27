// Package fingerprint provides language-neutral fingerprint helpers that any
// adapter can feed with its own token stream: an exact fingerprint over token
// text, a structural fingerprint over abstracted token kinds, and a MinHash
// sketch of structural shingles for near-structural candidate search.
package fingerprint

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"hash/fnv"
)

// ShingleSize is the number of structural tokens per shingle.
const ShingleSize = 5

// NumHashes is the MinHash sketch length.
const NumHashes = 32

// Builder accumulates tokens of one declaration.
type Builder struct {
	exact  []byte
	kinds  []string
	tokens int
}

// Add appends a token. text is the literal token text (used for the exact
// fingerprint); kind is the abstracted structural token (e.g. "ID", "LIT",
// "if", "{").
func (b *Builder) Add(text, kind string) {
	b.exact = append(b.exact, text...)
	b.exact = append(b.exact, 0)
	b.kinds = append(b.kinds, kind)
	b.tokens++
}

// Tokens returns the token count.
func (b *Builder) Tokens() int { return b.tokens }

// Exact returns the exact fingerprint.
func (b *Builder) Exact() string {
	h := sha256.Sum256(b.exact)
	return hex.EncodeToString(h[:12])
}

// Struct returns the structural fingerprint.
func (b *Builder) Struct() string {
	h := sha256.New()
	for _, k := range b.kinds {
		h.Write([]byte(k))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)[:12])
}

// MinHash returns a MinHash sketch over structural shingles. Declarations
// shorter than one shingle return nil.
func (b *Builder) MinHash() []uint64 {
	if len(b.kinds) < ShingleSize {
		return nil
	}
	sig := make([]uint64, NumHashes)
	for i := range sig {
		sig[i] = ^uint64(0)
	}
	for i := 0; i+ShingleSize <= len(b.kinds); i++ {
		h := fnv.New64a()
		for _, k := range b.kinds[i : i+ShingleSize] {
			h.Write([]byte(k))
			h.Write([]byte{0})
		}
		base := h.Sum64()
		for j := range sig {
			v := mix(base, uint64(j))
			if v < sig[j] {
				sig[j] = v
			}
		}
	}
	return sig
}

func mix(x, seed uint64) uint64 {
	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:8], x)
	binary.LittleEndian.PutUint64(buf[8:], seed*0x9E3779B97F4A7C15+1)
	h := fnv.New64a()
	h.Write(buf[:])
	z := h.Sum64()
	z ^= z >> 31
	z *= 0xBF58476D1CE4E5B9
	z ^= z >> 29
	return z
}

// Similarity estimates Jaccard similarity of two sketches (0 when either is
// missing).
func Similarity(a, b []uint64) float64 {
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}
	eq := 0
	for i := range a {
		if a[i] == b[i] {
			eq++
		}
	}
	return float64(eq) / float64(len(a))
}

// Bands splits a sketch into LSH band keys (rows per band = len/bands).
func Bands(sig []uint64, bands int) []uint64 {
	if len(sig) == 0 {
		return nil
	}
	rows := len(sig) / bands
	out := make([]uint64, 0, bands)
	for i := 0; i < bands; i++ {
		h := fnv.New64a()
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(i))
		h.Write(buf[:])
		for _, v := range sig[i*rows : (i+1)*rows] {
			binary.LittleEndian.PutUint64(buf[:], v)
			h.Write(buf[:])
		}
		out = append(out, h.Sum64())
	}
	return out
}
