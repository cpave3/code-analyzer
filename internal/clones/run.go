package clones

import (
	"fmt"

	"sort"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/memo"
	"github.com/cameron/code-evaluator/internal/model"
)

// Result is the clone analysis of one snapshot (optionally scoped).
type Result struct {
	Version    string         `json:"version"`
	Options    Options        `json:"options"`
	Callables  map[string]int `json:"callablesCompared"` // by language
	NearMisses []Pair         `json:"nearMisses"`
	Fragments  []Pair         `json:"fragments"`
	Notes      []string       `json:"notes"`
}

// Notes on what the result means.
func Notes() []string {
	return []string{
		"unscored: scoring/v2 duplication debt counts only exact and structural clone groups; near misses and fragments are listed for review",
		"tokens are abstracted (identifiers → ID, literals → LIT, comments dropped), so renamed copies match; similarity is 2·LCS/(|a|+|b|) over those tokens",
		"similar structure is not always a copy: small handlers, table-driven tests and parallel per-language implementations can look alike",
	}
}

// Analyze runs both detectors on a snapshot's source callables. inScope
// keeps pairs with at least one side in scope (nil = all).
func Analyze(snap *index.Snapshot, cfg *config.Config, langs *lang.Registry, inScope func(string) bool) *Result {
	return AnalyzeCached(snap, cfg, langs, inScope, nil)
}

// AnalyzeCached is Analyze with per-file extraction cached by content, so
// a series of commits only re-extracts the files each commit changed.
func AnalyzeCached(snap *index.Snapshot, cfg *config.Config, langs *lang.Registry, inScope func(string) bool, cache *memo.Cache) *Result {
	opt := DefaultOptions(cfg.Duplication.MinTokens, cfg.Duplication.NearThreshold)
	r := &Result{Version: Version, Options: opt, Callables: map[string]int{}, NearMisses: []Pair{}, Fragments: []Pair{}, Notes: Notes()}
	var all []Callable
	for _, ad := range langs.Adapters() {
		sp, ok := ad.(lang.SyntaxParser)
		if !ok {
			continue
		}
		var cs []Callable
		for _, p := range snap.Paths {
			f := snap.Files[p]
			if f.Adapter == nil || f.Adapter.ID() != ad.ID() || f.Parsed == nil || len(f.Parsed.ParseErrors) > 0 {
				continue
			}
			if f.Class != model.ClassSource && !(cfg.Duplication.IncludeTests && f.Class == model.ClassTest) {
				continue
			}
			key := fmt.Sprintf("clones:%s:%s:%s:%d", Version, ad.ID(), f.Blob, opt.MinTokens)
			extracted := cache.Get(key, func() any {
				src, err := snap.Content(p)
				if err != nil {
					return []Callable(nil)
				}
				tree, err := sp.ParseSyntax(src)
				if err != nil {
					return []Callable(nil)
				}
				defer tree.Close()
				return Extract(p, ad.ID(), src, tree.RootNode(), snap.SymbolsIn(p), opt.MinTokens)
			}).([]Callable)
			for _, c := range extracted {
				c.File = p // the same content may sit at another path
				cs = append(cs, c)
			}
		}
		r.Callables[ad.ID()] = len(cs)
		all = append(all, cs...)
	}
	near := NearMisses(all, opt)
	frags := Fragments(all, opt, near)
	keep := func(ps []Pair) []Pair {
		out := []Pair{}
		for _, p := range ps {
			if inScope == nil || inScope(p.A.File) || inScope(p.B.File) {
				out = append(out, p)
			}
		}
		return out
	}
	r.NearMisses, r.Fragments = keep(near), keep(frags)
	return r
}

// Introduced returns the pairs in after whose key is not in before.
func Introduced(before, after []Pair) []Pair {
	had := map[string]bool{}
	for _, p := range before {
		had[p.Key] = true
	}
	out := []Pair{}
	for _, p := range after {
		if !had[p.Key] {
			out = append(out, p)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Shared > out[j].Shared })
	return out
}
