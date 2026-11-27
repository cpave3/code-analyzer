// Package duplication finds exact and structural callable clones against the
// whole repository on both sides of the change, and attributes changes in
// redundancy to the patch: introduced, removed/consolidated, or pre-existing
// duplication that was merely touched. Near-structural matches are reported
// as reuse candidates, never as proof of equivalence or missed reuse.
package duplication

import (
	"fmt"
	"sort"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/lang/fingerprint"
	"github.com/cameron/code-evaluator/internal/model"
)

// Analyzer is the duplication analyzer.
type Analyzer struct{}

// New returns the analyzer.
func New() *Analyzer { return &Analyzer{} }

func (Analyzer) Descriptor() analyzer.Descriptor {
	return analyzer.Descriptor{
		ID: "duplication", Dimension: analyzer.DimDuplication, Title: "Callable duplication delta", Version: "1",
		Requires:    []lang.Capability{lang.CapSymbols, lang.CapFingerprints},
		Description: "Exact (token text) and structural (identifiers/literals abstracted) clone groups of callables across the repository, before and after; near-structural reuse candidates via MinHash.",
		Classes:     []model.FileClass{model.ClassSource},
	}
}

const caveat = "structural similarity is a candidate signal, not proof of semantic equivalence or of missed reuse"

func population(ctx *analyzer.Context, d analyzer.Descriptor, snap *index.Snapshot, minTokens int, includeTests bool) []*model.Symbol {
	var out []*model.Symbol
	for _, s := range snap.Symbols {
		if !s.Callable() || s.Tokens < minTokens || s.InErrorRegion {
			continue
		}
		f := snap.Files[s.File]
		if f == nil || !f.Supported() {
			continue
		}
		if !(d.InClasses(f.Class) || (includeTests && f.Class == model.ClassTest)) {
			continue
		}
		if !lang.Has(f.Adapter, lang.CapFingerprints) {
			continue
		}
		out = append(out, s)
	}
	return out
}

func groups(pop []*model.Symbol, key func(*model.Symbol) string) map[string][]*model.Symbol {
	g := map[string][]*model.Symbol{}
	for _, s := range pop {
		k := key(s)
		if k == "" {
			continue
		}
		g[k] = append(g[k], s)
	}
	return g
}

func redundancy(g map[string][]*model.Symbol) int {
	n := 0
	for _, v := range g {
		if len(v) > 1 {
			n += len(v) - 1
		}
	}
	return n
}

func exactKey(s *model.Symbol) string  { return s.Kind + ":" + s.ExactFP }
func structKey(s *model.Symbol) string { return s.Kind + ":" + s.StructFP }

func (an Analyzer) Analyze(ctx *analyzer.Context) (*analyzer.Result, error) {
	d := an.Descriptor()
	cfg := ctx.Config.Duplication
	res := &analyzer.Result{}
	basePop := population(ctx, d, ctx.Base, cfg.MinTokens, cfg.IncludeTests)
	targetPop := population(ctx, d, ctx.Target, cfg.MinTokens, cfg.IncludeTests)
	bExact, tExact := groups(basePop, exactKey), groups(targetPop, exactKey)
	bStruct, tStruct := groups(basePop, structKey), groups(targetPop, structKey)

	changedTarget := map[*model.Symbol]string{}
	changedBase := map[*model.Symbol]string{}
	for _, sc := range ctx.Change.Symbols {
		if sc.Target != nil {
			changedTarget[sc.Target] = sc.Status
		}
		if sc.Base != nil {
			changedBase[sc.Base] = sc.Status
		}
	}

	exB, exA := redundancy(bExact), redundancy(tExact)
	stB, stA := redundancy(bStruct), redundancy(tStruct)
	res.Measurements = append(res.Measurements,
		analyzer.Measure("exact_redundant_callables", "repository", exB, exA, "count"),
		analyzer.Measure("structural_redundant_callables", "repository", stB, stA, "count"),
		analyzer.Measure("callables_considered", "repository", len(basePop), len(targetPop), "count"),
	)

	introduced, removed, touched := 0, 0, 0
	keys := map[string]bool{}
	for k := range bStruct {
		keys[k] = true
	}
	for k := range tStruct {
		keys[k] = true
	}
	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	inGroupFinding := map[*model.Symbol]bool{}
	for _, k := range sorted {
		bs, ts := bStruct[k], tStruct[k]
		before, after := max0(len(bs)-1), max0(len(ts)-1)
		switch {
		case after > before:
			introduced += after - before
			exact := false
			var ev []analyzer.Evidence
			var newNames, existing []string
			for _, s := range ts {
				role := "existing match"
				if st, ok := changedTarget[s]; ok && (st == change.SymAdded || st == change.SymModified || st == change.SymRenamed) {
					role = "new or changed (" + st + ")"
					newNames = append(newNames, s.DisplayName())
					inGroupFinding[s] = true
					if len(bExact[exactKey(s)]) < len(tExact[exactKey(s)]) && len(tExact[exactKey(s)]) > 1 {
						exact = true
					}
				} else {
					existing = append(existing, s.DisplayName())
				}
				ev = append(ev, ctx.Ev("target", s, role))
			}
			kind, basis, msg := "duplication.introduced_structural", analyzer.BasisMeasured, ""
			f := analyzer.Finding{Basis: basis, Direction: analyzer.DirDegradation, Evidence: ev,
				Measurements: []analyzer.Measurement{analyzer.Measure("group_size", "", len(bs), len(ts), "callables")}}
			if exact {
				kind = "duplication.introduced_exact"
				msg = fmt.Sprintf("identical callable body introduced: %v", newNames)
			} else {
				msg = fmt.Sprintf("structurally identical callable introduced (names/literals differ): %v", newNames)
				f.Caveat = caveat
			}
			if len(existing) > 0 {
				msg += fmt.Sprintf("; matches pre-existing %v — possible reuse candidate", existing)
			}
			f.Kind, f.Message = kind, msg
			res.Findings = append(res.Findings, f)
		case after < before:
			removed += before - after
			var ev []analyzer.Evidence
			var gone []string
			for _, s := range bs {
				role := "remaining copy"
				if st, ok := changedBase[s]; ok && (st == change.SymDeleted || st == change.SymModified) {
					role = "removed or changed (" + st + ")"
					gone = append(gone, s.DisplayName())
				}
				ev = append(ev, ctx.Ev("base", s, role))
			}
			res.Findings = append(res.Findings, analyzer.Finding{
				Kind: "duplication.removed", Basis: analyzer.BasisMeasured, Direction: analyzer.DirImprovement,
				Message:      fmt.Sprintf("clone group shrank from %d to %d (removed/changed: %v)", len(bs), len(ts), gone),
				Measurements: []analyzer.Measurement{analyzer.Measure("group_size", "", len(bs), len(ts), "callables")},
				Evidence:     ev,
			})
		case after > 0:
			var ev []analyzer.Evidence
			hit := false
			for _, s := range ts {
				if st, ok := changedTarget[s]; ok && st != change.SymFormatting {
					hit = true
					ev = append(ev, ctx.Ev("target", s, "changed member ("+st+")"))
				}
			}
			if hit {
				touched++
				res.Findings = append(res.Findings, analyzer.Finding{
					Kind: "duplication.preexisting_touched", Basis: analyzer.BasisMeasured, Direction: analyzer.DirNeutral,
					Message:  fmt.Sprintf("pre-existing clone group of %d still has %d members; not counted as introduced", len(bs), len(ts)),
					Evidence: ev,
				})
			}
		}
	}

	near := an.nearCandidates(ctx, targetPop, basePop, changedTarget, inGroupFinding, cfg.NearThreshold)
	res.Findings = append(res.Findings, near...)
	res.Measurements = append(res.Measurements,
		analyzer.Count("redundant_instances_introduced", "change", introduced, "count"),
		analyzer.Count("redundant_instances_removed", "change", removed, "count"),
		analyzer.Count("preexisting_groups_touched", "change", touched, "count"),
		analyzer.Count("near_structural_candidates", "change", len(near), "count"),
	)
	net := stA - stB
	score := analyzer.ScoreFromThresholds(net, cfg.DegradeAt, cfg.SignificantDegradeAt)
	res.Score = &score
	res.ScoreRule = fmt.Sprintf("net structural redundancy Δ%+d: ≥%d → -1, ≥%d → -2, ≤-%d → +1, ≤-%d → +2; near-structural candidates are not scored",
		net, cfg.DegradeAt, cfg.SignificantDegradeAt, cfg.DegradeAt, cfg.SignificantDegradeAt)
	res.Rationale = append(res.Rationale,
		fmt.Sprintf("callables with ≥%d tokens in production source are compared repository-wide (base %d, target %d)", cfg.MinTokens, len(basePop), len(targetPop)),
		fmt.Sprintf("introduced %d redundant instance(s), removed %d; %d pre-existing group(s) touched but unchanged in size", introduced, removed, touched),
	)
	res.Summary = fmt.Sprintf("structural redundancy %d→%d (Δ%+d), exact %d→%d; %d near-structural candidate(s)", stB, stA, net, exB, exA, len(near))
	return res, nil
}

const bands = 8

// nearCandidates finds changed callables similar (MinHash) to other callables,
// excluding pairs that were already similar before the change.
func (an Analyzer) nearCandidates(ctx *analyzer.Context, targetPop, basePop []*model.Symbol, changed map[*model.Symbol]string, skip map[*model.Symbol]bool, threshold float64) []analyzer.Finding {
	buckets := map[uint64][]*model.Symbol{}
	for _, s := range targetPop {
		for _, b := range fingerprint.Bands(s.MinHash, bands) {
			buckets[b] = append(buckets[b], s)
		}
	}
	baseByKey := map[string]*model.Symbol{}
	for _, s := range basePop {
		baseByKey[s.Key] = s
	}
	var out []analyzer.Finding
	seen := map[[2]*model.Symbol]bool{}
	var subjects []*model.Symbol
	for _, s := range targetPop {
		if st, ok := changed[s]; ok && (st == change.SymAdded || st == change.SymModified) && !skip[s] {
			subjects = append(subjects, s)
		}
	}
	for _, s := range subjects {
		cands := map[*model.Symbol]bool{}
		for _, b := range fingerprint.Bands(s.MinHash, bands) {
			for _, c := range buckets[b] {
				if c != s {
					cands[c] = true
				}
			}
		}
		var best []*model.Symbol
		for c := range cands {
			if c.StructFP == s.StructFP && c.Kind == s.Kind {
				continue // exact structural match handled by group analysis
			}
			if ratio(c.Tokens, s.Tokens) < 0.67 {
				continue
			}
			sim := fingerprint.Similarity(s.MinHash, c.MinHash)
			if sim < threshold {
				continue
			}
			// Pre-existing similarity between the same pair is not new.
			if bs, ok := baseByKey[s.Key]; ok && changed[s] == change.SymModified {
				if bc, ok := baseByKey[c.Key]; ok && fingerprint.Similarity(bs.MinHash, bc.MinHash) >= threshold {
					continue
				}
			}
			pair := [2]*model.Symbol{s, c}
			if seen[[2]*model.Symbol{c, s}] || seen[pair] {
				continue
			}
			seen[pair] = true
			best = append(best, c)
		}
		sort.Slice(best, func(i, j int) bool {
			si, sj := fingerprint.Similarity(s.MinHash, best[i].MinHash), fingerprint.Similarity(s.MinHash, best[j].MinHash)
			if si != sj {
				return si > sj
			}
			if best[i].File != best[j].File {
				return best[i].File < best[j].File
			}
			return best[i].StartLine < best[j].StartLine
		})
		if len(best) > 3 {
			best = best[:3]
		}
		for _, c := range best {
			sim := fingerprint.Similarity(s.MinHash, c.MinHash)
			role := "existing callable"
			if _, ok := changed[c]; ok {
				role = "also changed in this patch"
			}
			conf := "low"
			if sim >= 0.95 {
				conf = "medium"
			}
			out = append(out, analyzer.Finding{
				Kind: "duplication.near_structural_candidate", Basis: analyzer.BasisHeuristic, Direction: analyzer.DirInfo,
				Message:    fmt.Sprintf("%s (%s) is structurally similar to %s (estimated shingle Jaccard %.2f)", s.DisplayName(), changed[s], c.DisplayName(), sim),
				Similarity: &sim, Confidence: conf, Caveat: caveat,
				Evidence: []analyzer.Evidence{ctx.Ev("target", s, "changed callable"), ctx.Ev("target", c, role)},
			})
		}
	}
	return out
}

func ratio(a, b int) float64 {
	if a > b {
		a, b = b, a
	}
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func max0(v int) int {
	if v < 0 {
		return 0
	}
	return v
}
