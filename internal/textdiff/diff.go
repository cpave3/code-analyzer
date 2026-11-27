// Package textdiff computes line-level differences in process, so the
// evaluator never depends on parsing Git patch headers (which are awkward for
// unusual filenames) and never invokes external diff drivers.
package textdiff

import (
	"bytes"
	"hash/fnv"
)

// Range is an inclusive 1-based line range.
type Range struct{ Start, End int }

// Result summarises the difference between two texts.
type Result struct {
	Added, Deleted int
	OldChanged     []Range // lines deleted or replaced in the old text
	NewChanged     []Range // lines added or replaced in the new text
	Approximate    bool    // edit distance exceeded the budget; whole files treated as changed
}

// MaxEdits bounds the Myers search. Larger diffs fall back to treating both
// sides as fully changed, which is reported as Approximate.
const MaxEdits = 2000

// SplitLines splits text into lines without trailing newline characters.
func SplitLines(b []byte) [][]byte {
	if len(b) == 0 {
		return nil
	}
	lines := bytes.Split(b, []byte{'\n'})
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func hashLines(lines [][]byte) []uint64 {
	out := make([]uint64, len(lines))
	for i, l := range lines {
		h := fnv.New64a()
		h.Write(bytes.TrimRight(l, "\r"))
		out[i] = h.Sum64()
	}
	return out
}

// Lines diffs two texts line by line.
func Lines(oldText, newText []byte) Result {
	a, b := hashLines(SplitLines(oldText)), hashLines(SplitLines(newText))
	// Trim common prefix/suffix.
	pre := 0
	for pre < len(a) && pre < len(b) && a[pre] == b[pre] {
		pre++
	}
	suf := 0
	for suf < len(a)-pre && suf < len(b)-pre && a[len(a)-1-suf] == b[len(b)-1-suf] {
		suf++
	}
	am, bm := a[pre:len(a)-suf], b[pre:len(b)-suf]
	ops, ok := myers(am, bm)
	var res Result
	if !ok {
		res.Approximate = true
		res.Deleted, res.Added = len(am), len(bm)
		if len(am) > 0 {
			res.OldChanged = []Range{{pre + 1, pre + len(am)}}
		}
		if len(bm) > 0 {
			res.NewChanged = []Range{{pre + 1, pre + len(bm)}}
		}
		return res
	}
	// ops: sequence of 'e' (equal), 'd' (delete from a), 'i' (insert from b)
	ai, bi := pre, pre
	for _, op := range ops {
		switch op {
		case 'e':
			ai++
			bi++
		case 'd':
			ai++
			res.Deleted++
			res.OldChanged = addLine(res.OldChanged, ai)
		case 'i':
			bi++
			res.Added++
			res.NewChanged = addLine(res.NewChanged, bi)
		}
	}
	return res
}

func addLine(rs []Range, line int) []Range {
	if n := len(rs); n > 0 && rs[n-1].End == line-1 {
		rs[n-1].End = line
		return rs
	}
	return append(rs, Range{line, line})
}

// myers returns the edit script for a -> b, or ok=false when the edit
// distance exceeds MaxEdits.
func myers(a, b []uint64) ([]byte, bool) {
	n, m := len(a), len(b)
	if n == 0 && m == 0 {
		return nil, true
	}
	max := n + m
	if max > MaxEdits {
		max = MaxEdits
	}
	offset := max + 1
	v := make([]int, 2*max+3)
	// trace[d] holds v[-(d+1) .. d+1] as it was at the start of round d.
	var trace [][]int
	for d := 0; d <= max; d++ {
		snapshot := make([]int, 2*d+3)
		copy(snapshot, v[offset-d-1:offset+d+2])
		trace = append(trace, snapshot)
		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[offset+k-1] < v[offset+k+1]) {
				x = v[offset+k+1]
			} else {
				x = v[offset+k-1] + 1
			}
			y := x - k
			for x < n && y < m && a[x] == b[y] {
				x++
				y++
			}
			v[offset+k] = x
			if x >= n && y >= m {
				return backtrack(trace, n, m, d), true
			}
		}
	}
	return nil, false
}

func backtrack(trace [][]int, n, m, dEnd int) []byte {
	var ops []byte
	x, y := n, m
	for d := dEnd; d > 0; d-- {
		w := trace[d]
		at := func(k int) int { return w[k+d+1] }
		k := x - y
		var prevK int
		if k == -d || (k != d && at(k-1) < at(k+1)) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}
		prevX := at(prevK)
		prevY := prevX - prevK
		for x > prevX && y > prevY {
			ops = append(ops, 'e')
			x--
			y--
		}
		if x == prevX {
			ops = append(ops, 'i')
		} else {
			ops = append(ops, 'd')
		}
		x, y = prevX, prevY
	}
	for x > 0 && y > 0 {
		ops = append(ops, 'e')
		x--
		y--
	}
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}
	return ops
}

// Overlaps reports whether [start,end] intersects any range.
func Overlaps(rs []Range, start, end int) bool {
	for _, r := range rs {
		if r.Start <= end && start <= r.End {
			return true
		}
	}
	return false
}
