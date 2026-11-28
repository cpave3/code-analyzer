// Package clones finds duplication that the exact and structural clone
// groups (scoring/v2) miss:
//
//   - near-miss clones: whole callables whose abstracted token sequences
//     (identifiers → ID, literals → LIT) are similar but not identical,
//     e.g. one statement added or a condition changed. Candidates come from
//     MinHash/LSH over token shingles; each pair is verified with a
//     longest-common-subsequence similarity 2·LCS/(|a|+|b|).
//   - fragment clones: a run of at least FragmentTokens identical
//     abstracted tokens shared by two different callables (or two
//     non-overlapping places in one), i.e. a copied block inside otherwise
//     different code.
//
// Results are findings, not debt: scoring/v2 counts only exact and
// structural clone groups, and this package does not change it.
package clones

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/cameron/code-evaluator/internal/lang/fingerprint"
	"github.com/cameron/code-evaluator/internal/model"
)

// Version identifies the detection semantics.
const Version = "clones/v1"

// Options tune detection.
type Options struct {
	MinTokens      int     // smallest callable compared (as for structural clones)
	Similarity     float64 // near-miss threshold on 2·LCS/(|a|+|b|)
	FragmentTokens int     // smallest shared run reported as a fragment clone
	MaxPairsPer    int     // near-miss partners kept per callable
}

// DefaultOptions mirror the duplication config defaults.
func DefaultOptions(minTokens int, near float64) Options {
	if minTokens <= 0 {
		minTokens = 40
	}
	if near <= 0 {
		near = 0.85
	}
	frag := 50
	if minTokens > frag {
		frag = minTokens
	}
	return Options{MinTokens: minTokens, Similarity: near, FragmentTokens: frag, MaxPairsPer: 3}
}

// Callable is one callable's abstracted token stream.
type Callable struct {
	File      string
	Language  string
	Name      string
	Kind      string
	StartLine int
	EndLine   int
	StructFP  string
	Tokens    []string // abstracted kinds
	Lines     []int    // source line per token
}

var literalKinds = map[string]bool{
	"string": true, "template_string": true, "number": true, "regex": true, "jsx_text": true, "true": true, "false": true,
	"null": true, "undefined": true, "nil": true, "int_literal": true, "float_literal": true, "imaginary_literal": true,
	"interpreted_string_literal": true, "raw_string_literal": true, "rune_literal": true,
}

// Extract abstracts the callables of one file from its syntax tree.
func Extract(path, language string, src []byte, root *sitter.Node, syms []*model.Symbol, minTokens int) []Callable {
	var out []Callable
	for _, s := range syms {
		if !s.Callable() || s.InErrorRegion || s.Tokens < minTokens {
			continue
		}
		c := Callable{File: path, Language: language, Name: s.DisplayName(), Kind: s.Kind, StartLine: s.StartLine, EndLine: s.EndLine, StructFP: s.StructFP}
		lo, hi := uint(s.StartLine-1), uint(s.EndLine-1)
		var visit func(n *sitter.Node)
		visit = func(n *sitter.Node) {
			if n.EndPosition().Row < lo || n.StartPosition().Row > hi {
				return
			}
			k := n.Kind()
			if k == "comment" {
				return
			}
			inside := n.StartPosition().Row >= lo && n.EndPosition().Row <= hi
			if inside && literalKinds[k] {
				c.Tokens, c.Lines = append(c.Tokens, "LIT"), append(c.Lines, int(n.StartPosition().Row)+1)
				return
			}
			if n.ChildCount() == 0 {
				if !inside {
					return
				}
				tok := k
				if strings.Contains(k, "identifier") {
					tok = "ID"
				}
				c.Tokens, c.Lines = append(c.Tokens, tok), append(c.Lines, int(n.StartPosition().Row)+1)
				return
			}
			for i := uint(0); i < n.ChildCount(); i++ {
				visit(n.Child(i))
			}
		}
		visit(root)
		if len(c.Tokens) >= minTokens {
			out = append(out, c)
		}
	}
	return out
}

// Ref locates one side of a clone.
type Ref struct {
	File      string `json:"file"`
	Name      string `json:"symbol"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Tokens    int    `json:"tokens"`
}

// Pair is a near-miss or fragment clone between two places.
type Pair struct {
	Kind       string  `json:"kind"` // near_miss | fragment
	Language   string  `json:"language"`
	A          Ref     `json:"a"`
	B          Ref     `json:"b"`
	Similarity float64 `json:"similarity,omitempty"` // near_miss: 2·LCS/(|a|+|b|)
	Shared     int     `json:"sharedTokens"`         // near_miss: LCS length; fragment: run length
	Key        string  `json:"key"`                  // line-independent identity
}

func ref(c *Callable) Ref {
	return Ref{File: c.File, Name: c.Name, StartLine: c.StartLine, EndLine: c.EndLine, Tokens: len(c.Tokens)}
}

func pairKey(kind string, a, b Ref) string {
	x, y := a.File+"#"+a.Name, b.File+"#"+b.Name
	if y < x {
		x, y = y, x
	}
	return kind + "|" + x + "|" + y
}

// lcs returns the longest common subsequence length of two token streams.
func lcs(a, b []int) int {
	if len(a) < len(b) {
		a, b = b, a
	}
	prev, cur := make([]int, len(b)+1), make([]int, len(b)+1)
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			switch {
			case a[i-1] == b[j-1]:
				cur[j] = prev[j-1] + 1
			case prev[j] >= cur[j-1]:
				cur[j] = prev[j]
			default:
				cur[j] = cur[j-1]
			}
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

type coded struct {
	c    *Callable
	ids  []int
	sig  []uint64
	size int
}

func encode(cs []Callable) []coded {
	dict := map[string]int{}
	out := make([]coded, len(cs))
	for i := range cs {
		c := &cs[i]
		ids := make([]int, len(c.Tokens))
		var b fingerprint.Builder
		for j, t := range c.Tokens {
			id, ok := dict[t]
			if !ok {
				id = len(dict)
				dict[t] = id
			}
			ids[j] = id
			b.Add(t, t)
		}
		out[i] = coded{c: c, ids: ids, sig: b.MinHash(), size: len(ids)}
	}
	return out
}

// maxLCSTokens bounds the quadratic verification.
const maxLCSTokens = 4000

// NearMisses finds similar, structurally non-identical callable pairs in
// the same language.
func NearMisses(cs []Callable, opt Options) []Pair {
	enc := encode(cs)
	buckets := map[string][]int{}
	for i, e := range enc {
		for _, band := range fingerprint.Bands(e.sig, 8) {
			k := e.c.Language + fmt.Sprint(band)
			buckets[k] = append(buckets[k], i)
		}
	}
	seen := map[[2]int]bool{}
	partners := map[int][]Pair{}
	for _, idx := range buckets {
		for x := 0; x < len(idx); x++ {
			for y := x + 1; y < len(idx); y++ {
				i, j := idx[x], idx[y]
				if i > j {
					i, j = j, i
				}
				if i == j || seen[[2]int{i, j}] {
					continue
				}
				seen[[2]int{i, j}] = true
				a, b := enc[i], enc[j]
				if a.c.StructFP != "" && a.c.StructFP == b.c.StructFP {
					continue // identical structure: already a structural clone group
				}
				lo, hi := a.size, b.size
				if lo > hi {
					lo, hi = hi, lo
				}
				if float64(lo)/float64(hi) < 0.67 || hi > maxLCSTokens {
					continue
				}
				if float64(2*lo)/float64(lo+hi) < opt.Similarity {
					continue // even a perfect subsequence cannot reach the threshold
				}
				l := lcs(a.ids, b.ids)
				sim := float64(2*l) / float64(a.size+b.size)
				if sim < opt.Similarity {
					continue
				}
				p := Pair{Kind: "near_miss", Language: a.c.Language, A: ref(a.c), B: ref(b.c), Similarity: float64(int(sim*1000)) / 1000, Shared: l}
				if p.B.File < p.A.File || (p.B.File == p.A.File && p.B.StartLine < p.A.StartLine) {
					p.A, p.B = p.B, p.A
				}
				p.Key = pairKey(p.Kind, p.A, p.B)
				partners[i] = append(partners[i], p)
				partners[j] = append(partners[j], p)
			}
		}
	}
	keep := map[string]Pair{}
	for _, ps := range partners {
		sort.Slice(ps, func(x, y int) bool { return ps[x].Similarity > ps[y].Similarity })
		for n, p := range ps {
			if opt.MaxPairsPer > 0 && n >= opt.MaxPairsPer {
				break
			}
			keep[p.Key] = p
		}
	}
	return sortPairs(keep)
}

func sortPairs(m map[string]Pair) []Pair {
	out := make([]Pair, 0, len(m))
	for _, p := range m {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Shared != out[j].Shared {
			return out[i].Shared > out[j].Shared
		}
		return out[i].Key < out[j].Key
	})
	return out
}

type occ struct{ c, pos int }

// Fragments finds shared runs of at least FragmentTokens abstracted tokens
// between different callables, or non-overlapping places in one callable.
// Pairs already reported as near misses (or identical structure) are
// skipped. One fragment pair is kept per callable pair: its longest run.
func Fragments(cs []Callable, opt Options, near []Pair) []Pair {
	w := opt.FragmentTokens
	enc := encode(cs)
	skip := map[string]bool{}
	for _, p := range near {
		skip[pairKey("fragment", p.A, p.B)] = true
	}
	windows := map[uint64][]occ{}
	for ci, e := range enc {
		if e.size < w {
			continue
		}
		for pos := 0; pos+w <= e.size; pos++ {
			h := fnv.New64a()
			fmt.Fprint(h, e.c.Language)
			for _, t := range e.ids[pos : pos+w] {
				fmt.Fprintf(h, ",%d", t)
			}
			windows[h.Sum64()] = append(windows[h.Sum64()], occ{ci, pos})
		}
	}
	// diagonals: (callable a, callable b, offset) -> positions in a
	type diag struct{ a, b, off int }
	runs := map[diag][]int{}
	for _, os := range windows {
		if len(os) > 50 {
			continue // boilerplate shared by many places; not a copy
		}
		for x := 0; x < len(os); x++ {
			for y := x + 1; y < len(os); y++ {
				a, b := os[x], os[y]
				if a.c > b.c || (a.c == b.c && a.pos > b.pos) {
					a, b = b, a
				}
				if a.c == b.c && b.pos-a.pos < w {
					continue // overlapping windows in one callable
				}
				runs[diag{a.c, b.c, b.pos - a.pos}] = append(runs[diag{a.c, b.c, b.pos - a.pos}], a.pos)
			}
		}
	}
	best := map[string]Pair{}
	for d, ps := range runs {
		ca, cb := enc[d.a].c, enc[d.b].c
		if d.a != d.b && ca.StructFP != "" && ca.StructFP == cb.StructFP {
			continue
		}
		sort.Ints(ps)
		start, prev := ps[0], ps[0]
		flush := func(s, e int) {
			length := e - s + w
			if d.a == d.b && length > d.off {
				return // the two regions overlap: repetition, not a copy
			}
			if periodic(enc[d.a].ids[s:s+length], w) {
				return // a short unit repeated (table rows, a list of calls), not a copied block
			}
			ra, rb := ref(ca), ref(cb)
			ra.StartLine, ra.EndLine = ca.Lines[s], ca.Lines[e+w-1]
			rb.StartLine, rb.EndLine = cb.Lines[s+d.off], cb.Lines[e+d.off+w-1]
			ra.Tokens, rb.Tokens = length, length
			k := pairKey("fragment", ref(ca), ref(cb))
			if skip[k] || ra.EndLine-ra.StartLine < 2 {
				return
			}
			if rb.File < ra.File || (rb.File == ra.File && rb.StartLine < ra.StartLine) {
				ra, rb = rb, ra // stable order, as for near misses
			}
			if old, ok := best[k]; !ok || length > old.Shared {
				best[k] = Pair{Kind: "fragment", Language: ca.Language, A: ra, B: rb, Shared: length, Key: k}
			}
		}
		for _, p := range ps[1:] {
			if p == prev+1 {
				prev = p
				continue
			}
			flush(start, prev)
			start, prev = p, p
		}
		flush(start, prev)
	}
	return sortPairs(best)
}

// periodic reports whether a run is a repetition of a unit shorter than
// maxPeriod tokens.
func periodic(run []int, maxPeriod int) bool {
	for p := 1; p < maxPeriod && p < len(run); p++ {
		ok := true
		for i := p; i < len(run); i++ {
			if run[i] != run[i-p] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
