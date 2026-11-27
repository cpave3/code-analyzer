// Package obsolete flags non-exported callables that were referenced before
// the change and have no remaining references after it, but were left in
// place: a candidate for code made obsolete by the patch. It also credits
// deletion of callables that were already unreferenced. References are
// name-based, so findings are heuristics.
package obsolete

import (
	"fmt"
	"sort"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

// Analyzer is the obsolete-code analyzer.
type Analyzer struct{}

// New returns the analyzer.
func New() *Analyzer { return &Analyzer{} }

func (Analyzer) Descriptor() analyzer.Descriptor {
	return analyzer.Descriptor{
		ID: "obsolete", Dimension: analyzer.DimMaintainability, Title: "Code left obsolete", Version: "1",
		Requires:    []lang.Capability{lang.CapSymbols, lang.CapReferences},
		Description: "Non-exported callables in touched units that lose all in-unit references but remain; deletion of already-unreferenced callables.",
		Classes:     []model.FileClass{model.ClassSource},
	}
}

// refCounts sums name references in a unit, split by production vs test files.
type refCounts struct{ prod, test map[string]int }

func unitRefs(snap *index.Snapshot, unit, language string) refCounts {
	rc := refCounts{prod: map[string]int{}, test: map[string]int{}}
	for _, p := range snap.Paths {
		f := snap.Files[p]
		if f.Parsed == nil || f.Language != language || f.Adapter.PackageOf(p) != unit {
			continue
		}
		dst := rc.prod
		if f.Class == model.ClassTest {
			dst = rc.test
		}
		for k, v := range f.Parsed.Refs {
			dst[k] += v
		}
	}
	return rc
}

func (an Analyzer) Analyze(ctx *analyzer.Context) (*analyzer.Result, error) {
	d := an.Descriptor()
	res := &analyzer.Result{}
	type unitKey struct{ unit, lang string }
	units := map[unitKey]bool{}
	for _, fc := range ctx.Change.Files {
		ad := ctx.Langs.Adapter(fc.Language)
		if ad == nil || !lang.Has(ad, lang.CapReferences) {
			continue
		}
		units[unitKey{ad.PackageOf(fc.Path), fc.Language}] = true
	}
	var keys []unitKey
	for k := range units {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].unit < keys[j].unit })

	deletedKeys := map[string]bool{}
	for _, sc := range ctx.Change.Symbols {
		if sc.Status == change.SymDeleted {
			deletedKeys[sc.Base.Key] = true
		}
	}
	targetByKey := map[string]*model.Symbol{}
	for _, s := range ctx.Target.Symbols {
		targetByKey[s.Key] = s
	}

	newlyUnref, removedDead := 0, 0
	for _, k := range keys {
		before, after := unitRefs(ctx.Base, k.unit, k.lang), unitRefs(ctx.Target, k.unit, k.lang)
		for _, b := range ctx.Base.Symbols {
			if b.Package != k.unit || !b.Callable() || b.Exported || b.Name == "main" || b.Name == "init" || b.Name == "_" {
				continue
			}
			if !ctx.SymbolAnalyzable(d, "base", b) {
				continue
			}
			wasRef := before.prod[b.Name] > 0
			if deletedKeys[b.Key] {
				if !wasRef && before.test[b.Name] == 0 {
					removedDead++
					res.Findings = append(res.Findings, analyzer.Finding{
						Kind: "obsolete.unreferenced_removed", Basis: analyzer.BasisHeuristic, Direction: analyzer.DirImprovement, Confidence: "medium",
						Message:  fmt.Sprintf("deleted %s had no name references in its unit before the change", b.DisplayName()),
						Caveat:   "name-based references; reflection, build tags and interface satisfaction are invisible",
						Evidence: []analyzer.Evidence{ctx.Ev("base", b, "deleted")},
					})
				}
				continue
			}
			t := targetByKey[b.Key]
			if t == nil || !wasRef || after.prod[t.Name] > 0 {
				continue
			}
			newlyUnref++
			conf := "medium"
			if b.Kind == model.KindMethod {
				conf = "low" // may satisfy an interface without explicit references
			}
			msg := fmt.Sprintf("%s was referenced %d time(s) in its unit before the change and has no production references after it, but was not removed", t.DisplayName(), before.prod[b.Name])
			if after.test[t.Name] > 0 {
				msg += fmt.Sprintf(" (still referenced %d time(s) from tests)", after.test[t.Name])
			}
			res.Findings = append(res.Findings, analyzer.Finding{
				Kind: "obsolete.left_unreferenced", Basis: analyzer.BasisHeuristic, Direction: analyzer.DirDegradation, Confidence: conf,
				Message:  msg,
				Caveat:   "name-based references; reflection, build tags, linkname and interface satisfaction are invisible",
				Evidence: []analyzer.Evidence{ctx.Ev("base", b, "referenced before"), ctx.Ev("target", t, "unreferenced after")},
			})
		}
	}
	res.Measurements = append(res.Measurements,
		analyzer.Count("callables_left_unreferenced", "touched_units", newlyUnref, "callables"),
		analyzer.Count("unreferenced_callables_removed", "touched_units", removedDead, "callables"),
	)
	medium := 0
	for _, f := range res.Findings {
		if f.Kind == "obsolete.left_unreferenced" && f.Confidence == "medium" {
			medium++
		}
	}
	score := 0
	switch {
	case medium >= 3:
		score = -2
		res.ScoreRule = "≥3 non-exported functions newly unreferenced but kept → -2"
	case medium >= 1:
		score = -1
		res.ScoreRule = "≥1 non-exported function newly unreferenced but kept → -1"
	case removedDead > 0 && newlyUnref == 0:
		score = 1
		res.ScoreRule = "already-unreferenced callables deleted, none newly orphaned → +1"
	default:
		res.ScoreRule = "no newly orphaned functions (methods alone are low confidence and not scored) → 0"
	}
	res.Score = &score
	res.Summary = fmt.Sprintf("%d callable(s) newly unreferenced but kept; %d unreferenced callable(s) removed", newlyUnref, removedDead)
	return res, nil
}
