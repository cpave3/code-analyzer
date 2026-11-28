package scorev2

import (
	"sort"

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
	pred := map[string]*model.Symbol{}
	keyedTarget := map[*model.Symbol]*model.Symbol{} // base → its same-key target
	var added, candidates []*model.Symbol
	for _, sc := range cs.Symbols {
		b, t := sc.Base, sc.Target
		switch {
		case b != nil && t != nil && sc.Status == change.SymRenamed:
			pred[t.Key] = b // exact body, new name
		case b != nil && t != nil:
			pred[t.Key] = b
			if b.Callable() && fingerprint.Similarity(b.MinHash, t.MinHash) < HeirSimilarity {
				keyedTarget[b] = t
				candidates = append(candidates, b)
			}
		case t != nil && t.Callable():
			added = append(added, t)
		case b != nil && b.Callable():
			candidates = append(candidates, b)
		}
	}
	type pair struct {
		t, b *model.Symbol
		sim  float64
		same bool
	}
	var pairs []pair
	for _, t := range added {
		for _, b := range candidates {
			if b.Kind != t.Kind || !comparableSize(b.Tokens, t.Tokens) {
				continue
			}
			if sim := fingerprint.Similarity(b.MinHash, t.MinHash); sim >= HeirSimilarity {
				pairs = append(pairs, pair{t, b, sim, b.File == t.File})
			}
		}
	}
	sort.SliceStable(pairs, func(i, j int) bool {
		if pairs[i].sim != pairs[j].sim {
			return pairs[i].sim > pairs[j].sim
		}
		if pairs[i].same != pairs[j].same {
			return pairs[i].same
		}
		if pairs[i].t.Key != pairs[j].t.Key {
			return pairs[i].t.Key < pairs[j].t.Key
		}
		return pairs[i].b.Key < pairs[j].b.Key
	})
	usedT, usedB := map[*model.Symbol]bool{}, map[*model.Symbol]bool{}
	for _, p := range pairs {
		if usedT[p.t] || usedB[p.b] {
			continue
		}
		usedT[p.t], usedB[p.b] = true, true
		pred[p.t.Key] = p.b
		if kt := keyedTarget[p.b]; kt != nil {
			delete(pred, kt.Key) // the old name now holds something new
		}
	}
	return pred
}

// comparableSize keeps a small helper from inheriting a large callable's
// lineage (or the reverse) on structure alone.
func comparableSize(a, b int) bool {
	if a > b {
		a, b = b, a
	}
	return b > 0 && 2*a >= b
}
