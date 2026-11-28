package scorev2

import (
	"cmp"
	"slices"
	"strings"

	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/lang/fingerprint"
	"github.com/cameron/code-evaluator/internal/model"
)

// HeirSimilarity is the least structural similarity (MinHash estimate of
// Jaccard over token shingles) at which a new callable is taken to be an old
// one renamed or moved, possibly with edits.
const HeirSimilarity = 0.5

// lineage maps each changed target callable's key to the base callable it
// descends from, so a check compares a callable with its own past rather
// than reporting a rename as new debt.
//
// A callable kept under its key descends from its base self, unless its body
// left: a new callable whose body is at least HeirSimilarity like an old one
// that was deleted, or whose old name now holds something else (typically a
// wrapper that calls the new name), takes over that old callable's lineage.
// Each old callable has at most one heir; the most similar pair wins.
// change.Set keeps its own exact-body rename pairing for scoring/v1.
func lineage(cs *change.Set) map[string]*model.Symbol {
	l := collect(cs)
	usedT, usedB := map[*model.Symbol]bool{}, map[*model.Symbol]bool{}
	for _, p := range l.pairs() {
		if usedT[p.t] || usedB[p.b] {
			continue
		}
		usedT[p.t], usedB[p.b] = true, true
		l.pred[p.t.Key] = p.b
		if kt := l.keyedTarget[p.b]; kt != nil {
			delete(l.pred, kt.Key) // the old name now holds something new
		}
	}
	return l.pred
}

// lines holds what lineage starts from: predecessors by key (an exact-body
// rename pairs as similarity 1 already), new callables, and old callables
// whose body may have gone to one of them.
type lines struct {
	pred        map[string]*model.Symbol
	keyedTarget map[*model.Symbol]*model.Symbol // base → its same-key target, when the body left
	added       []*model.Symbol
	candidates  []*model.Symbol
}

func collect(cs *change.Set) *lines {
	l := &lines{pred: map[string]*model.Symbol{}, keyedTarget: map[*model.Symbol]*model.Symbol{}}
	for _, sc := range cs.Symbols {
		b, t := sc.Base, sc.Target
		switch {
		case b != nil && t != nil:
			l.pred[t.Key] = b
			if b.Callable() && fingerprint.Similarity(b.MinHash, t.MinHash) < HeirSimilarity {
				l.keyedTarget[b] = t
				l.candidates = append(l.candidates, b)
			}
		case t != nil && t.Callable():
			l.added = append(l.added, t)
		case b != nil && b.Callable():
			l.candidates = append(l.candidates, b)
		}
	}
	return l
}

type heirPair struct {
	t, b *model.Symbol
	sim  float64
	same int // 1 when both are in the same file
}

// pairs lists every similar enough (new, old) pair, most similar first, then
// same file first, then by key for determinism.
func (l *lines) pairs() []heirPair {
	var out []heirPair
	for _, t := range l.added {
		for _, b := range l.candidates {
			sim := fingerprint.Similarity(b.MinHash, t.MinHash)
			if b.Kind == t.Kind && sim >= HeirSimilarity && comparableSize(b.Tokens, t.Tokens) {
				out = append(out, heirPair{t, b, sim, sameFile(b, t)})
			}
		}
	}
	slices.SortStableFunc(out, func(x, y heirPair) int {
		return cmp.Or(cmp.Compare(y.sim, x.sim), cmp.Compare(y.same, x.same), strings.Compare(x.t.Key, y.t.Key), strings.Compare(x.b.Key, y.b.Key))
	})
	return out
}

// comparableSize keeps a small helper from inheriting a large callable's
// lineage (or the reverse) on structure alone.
func comparableSize(a, b int) bool {
	if a > b {
		a, b = b, a
	}
	return b > 0 && 2*a >= b
}

func sameFile(a, b *model.Symbol) int {
	if a.File == b.File {
		return 1
	}
	return 0
}
