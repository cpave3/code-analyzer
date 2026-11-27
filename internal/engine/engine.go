// Package engine runs the evaluation pipeline: resolve revisions, index both
// snapshots, compute the change context, run every registered analyzer with
// engine-computed coverage, and assemble the report.
package engine

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/gitsrc"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
	"github.com/cameron/code-evaluator/internal/report"
)

// Options configures an Evaluator.
type Options struct {
	RepoPath  string
	Config    *config.Config
	CacheDir  string // "" = in-memory cache
	Langs     *lang.Registry
	Analyzers *analyzer.Registry
	Timing    bool // include timing metadata (non-deterministic)
}

// Evaluator evaluates changes in one repository.
type Evaluator struct {
	opts    Options
	repo    *gitsrc.Repo
	cache   *index.Cache
	builder *index.Builder
	snaps   []*index.Snapshot // small LRU by commit
}

// New opens the repository and cache. The cache directory must not be inside
// the evaluated repository.
func New(opts Options) (*Evaluator, error) {
	repo, err := gitsrc.Open(opts.RepoPath)
	if err != nil {
		return nil, err
	}
	if opts.CacheDir != "" {
		abs, _ := filepath.Abs(opts.CacheDir)
		if rel, err := filepath.Rel(repo.Path, abs); err == nil && !strings.HasPrefix(rel, "..") {
			repo.Close()
			return nil, fmt.Errorf("cache directory %s is inside the evaluated repository; refusing to write there", opts.CacheDir)
		}
	}
	cache, err := index.OpenCache(opts.CacheDir)
	if err != nil {
		repo.Close()
		return nil, err
	}
	if opts.Config == nil {
		opts.Config = config.Default()
	}
	return &Evaluator{opts: opts, repo: repo, cache: cache,
		builder: &index.Builder{Repo: repo, Langs: opts.Langs, Config: opts.Config, Cache: cache}}, nil
}

// Close releases resources.
func (e *Evaluator) Close() {
	e.repo.Close()
	e.cache.Close()
}

// RepoPath is the resolved repository top level.
func (e *Evaluator) RepoPath() string { return e.repo.Path }

func (e *Evaluator) snapshot(commit string) (*index.Snapshot, error) {
	for _, s := range e.snaps {
		if s.Commit == commit {
			return s, nil
		}
	}
	var s *index.Snapshot
	var err error
	if commit == gitsrc.WorktreeRev {
		s, err = e.builder.Worktree()
	} else {
		s, err = e.builder.Commit(commit)
	}
	if err != nil {
		return nil, err
	}
	e.snaps = append(e.snaps, s)
	if len(e.snaps) > 3 {
		e.snaps = e.snaps[1:]
	}
	return s, nil
}

// Evaluate compares baseRef with targetRef (or the working tree).
func (e *Evaluator) Evaluate(baseRef, targetRef string, worktree bool) (*report.Report, error) {
	start := time.Now()
	parsed0, hits0, miss0 := e.builder.Parsed, e.cache.Hits, e.cache.Misses
	baseSHA, err := e.repo.ResolveCommit(baseRef)
	if err != nil {
		return nil, err
	}
	targetSHA := gitsrc.WorktreeRev
	if !worktree {
		if targetSHA, err = e.repo.ResolveCommit(targetRef); err != nil {
			return nil, err
		}
	} else {
		targetRef = gitsrc.WorktreeRev
	}
	rep, err := e.evaluateResolved(baseRef, baseSHA, targetRef, targetSHA)
	if err != nil {
		return nil, err
	}
	if e.opts.Timing {
		rep.Timing = &report.Timing{StartedAt: start.Format(time.RFC3339), ElapsedMillis: time.Since(start).Milliseconds(),
			FilesParsed: e.builder.Parsed - parsed0, CacheHits: e.cache.Hits - hits0, CacheMisses: e.cache.Misses - miss0}
	}
	return rep, nil
}

func (e *Evaluator) evaluateResolved(baseRef, baseSHA, targetRef, targetSHA string) (*report.Report, error) {
	base, err := e.snapshot(baseSHA)
	if err != nil {
		return nil, fmt.Errorf("index base: %w", err)
	}
	target, err := e.snapshot(targetSHA)
	if err != nil {
		return nil, fmt.Errorf("index target: %w", err)
	}
	cs, cdiags, err := change.Compute(e.repo, base, target)
	if err != nil {
		return nil, fmt.Errorf("diff: %w", err)
	}
	rep := &report.Report{
		Schema:         report.Schema,
		Tool:           report.Tool{Name: "codeeval", Version: report.ToolVersion},
		ScoringVersion: config.ScoringVersion,
		ConfigVersion:  e.opts.Config.ConfigVersion,
		ConfigHash:     e.opts.Config.Hash(),
		Mode:           "commits",
		Repository:     filepath.Base(e.repo.Path),
		Base:           report.Revision{Ref: baseRef, Commit: baseSHA, Tree: base.Tree, CommitTime: e.repo.CommitTime(baseSHA)},
		Target:         report.Revision{Ref: targetRef, Commit: targetSHA, Tree: target.Tree},
		ScoreScale:     "-2 significant degradation, -1 degradation, 0 neutral/mixed, +1 improvement, +2 significant improvement; null = not scored (unsupported, error or not assessed). No combined score.",
	}
	if targetSHA == gitsrc.WorktreeRev {
		rep.Mode = "worktree"
	} else {
		rep.Target.CommitTime = e.repo.CommitTime(targetSHA)
	}
	rep.Languages = languageCoverage(e.opts.Langs, target, cs)
	rep.Change = changeSummary(cs, base, target)
	rep.Diagnostics = diagnostics(cs, base, target, cdiags)

	ctx := &analyzer.Context{Base: base, Target: target, Change: cs, Config: e.opts.Config, Langs: e.opts.Langs}
	for _, a := range e.opts.Analyzers.All() {
		rep.Results = append(rep.Results, runAnalyzer(a, ctx))
	}
	rep.Dimensions = dimensions(rep.Results)
	rep.Limitations = limitations(e.opts.Langs, rep.Mode)
	return rep, nil
}

// Coverage computes what an analyzer can see of the change from its declared
// classes and capability requirements.
func Coverage(d analyzer.Descriptor, ctx *analyzer.Context) analyzer.Coverage {
	cov := analyzer.Coverage{UnsupportedByLanguage: map[string]int{}, ExcludedByClass: map[string]int{}}
	for _, fc := range ctx.Change.Files {
		if fc.Language == "" {
			cov.NonCodeFiles++
			continue
		}
		if !d.InClasses(fc.Class) {
			cov.ExcludedByClass[string(fc.Class)]++
			continue
		}
		cov.ChangedCodeFiles++
		ad := ctx.Langs.Adapter(fc.Language)
		ok := ad != nil
		for _, c := range d.Requires {
			ok = ok && lang.Has(ad, c)
		}
		if !ok {
			cov.UnsupportedByLanguage[fc.Language]++
			continue
		}
		hasErr := false
		for _, snapFile := range []*index.File{ctx.Base.Files[fc.OldPath], ctx.Base.Files[fc.Path], ctx.Target.Files[fc.Path]} {
			if snapFile != nil && (snapFile.Parsed == nil || len(snapFile.Parsed.ParseErrors) > 0) {
				hasErr = true
			}
		}
		if hasErr {
			cov.WithParseErrors++
		} else {
			cov.Analyzed++
		}
	}
	return cov
}

func runAnalyzer(a analyzer.Analyzer, ctx *analyzer.Context) (res *analyzer.Result) {
	d := a.Descriptor()
	cov := Coverage(d, ctx)
	ctx.Coverage = cov
	defer func() {
		if r := recover(); r != nil {
			res = &analyzer.Result{Status: analyzer.StatusError, ScoreRule: "analyzer failed; not scored",
				Summary:     fmt.Sprintf("analyzer panicked: %v", r),
				Diagnostics: []model.Diagnostic{{Severity: model.SevError, Code: "analyzer_panic", Message: fmt.Sprint(r), Analyzer: d.ID}}}
		}
		finalize(res, d, cov)
	}()
	out, err := a.Analyze(ctx)
	if err != nil {
		return &analyzer.Result{Status: analyzer.StatusError, ScoreRule: "analyzer failed; not scored", Summary: err.Error(),
			Diagnostics: []model.Diagnostic{{Severity: model.SevError, Code: "analyzer_error", Message: err.Error(), Analyzer: d.ID}}}
	}
	if out == nil {
		return &analyzer.Result{Status: analyzer.StatusError, ScoreRule: "analyzer returned no result", Summary: "no result"}
	}
	return out
}

// finalize applies coverage-driven status. Unsupported is never neutral: a
// score is withheld when no changed code in scope could be analysed.
func finalize(res *analyzer.Result, d analyzer.Descriptor, cov analyzer.Coverage) {
	res.Analyzer, res.Dimension, res.Title, res.Version = d.ID, d.Dimension, d.Title, d.Version
	res.Coverage = cov
	if res.Measurements == nil {
		res.Measurements = []analyzer.Measurement{}
	}
	if res.Findings == nil {
		res.Findings = []analyzer.Finding{}
	}
	for i := range res.Findings {
		if res.Findings[i].Evidence == nil {
			res.Findings[i].Evidence = []analyzer.Evidence{}
		}
	}
	analyzer.SortFindings(res.Findings)
	for i := range res.Diagnostics {
		res.Diagnostics[i].Analyzer = d.ID
	}
	if res.Status == analyzer.StatusError || res.Status == analyzer.StatusNotAssessed {
		res.Score = nil
		return
	}
	if len(d.Requires) == 0 {
		if res.Status == "" {
			res.Status = analyzer.StatusOK
		}
		return
	}
	unsupported := cov.Unsupported()
	switch {
	case cov.ChangedCodeFiles > 0 && cov.Analyzed+cov.WithParseErrors == 0:
		res.Status = analyzer.StatusUnsupported
		res.Score = nil
		res.ScoreRule = "no changed code in scope could be analysed → not scored (unsupported is not neutral)"
		res.Summary = fmt.Sprintf("unsupported: %d changed file(s) in languages without the required capabilities (%s)", unsupported, langList(cov.UnsupportedByLanguage))
	case unsupported > 0 || cov.WithParseErrors > 0:
		res.Status = analyzer.StatusPartial
		res.Rationale = append(res.Rationale, fmt.Sprintf("partial coverage: score reflects %d analysed file(s); %d unsupported (%s) and %d with parse errors are not reflected",
			cov.Analyzed, unsupported, langList(cov.UnsupportedByLanguage), cov.WithParseErrors))
	default:
		res.Status = analyzer.StatusOK
	}
	if cov.ChangedCodeFiles == 0 && res.Score != nil && *res.Score == 0 {
		res.Rationale = append(res.Rationale, "no changed code files in this analyzer's scope")
	}
}

func langList(m map[string]int) string {
	var parts []string
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("%s:%d", k, v))
	}
	sort.Strings(parts)
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ", ")
}

func dimensions(results []*analyzer.Result) []report.Dimension {
	var out []report.Dimension
	for _, dim := range analyzer.DimensionOrder {
		d := report.Dimension{Dimension: dim, Analyzers: []report.AnalyzerScore{}}
		for _, r := range results {
			if r.Dimension == dim {
				d.Analyzers = append(d.Analyzers, report.AnalyzerScore{Analyzer: r.Analyzer, Status: r.Status, Score: r.Score, Summary: r.Summary})
			}
		}
		out = append(out, d)
	}
	return out
}

func languageCoverage(reg *lang.Registry, target *index.Snapshot, cs *change.Set) []report.LanguageCoverage {
	byID := map[string]*report.LanguageCoverage{}
	get := func(id string) *report.LanguageCoverage {
		if c, ok := byID[id]; ok {
			return c
		}
		c := &report.LanguageCoverage{ID: id, FilesByClass: map[string]int{}}
		if a := reg.Adapter(id); a != nil {
			c.Supported, c.AdapterVer, c.Limits = true, a.Version(), a.Limits()
			c.Capabilities = map[string]string{}
			for k, v := range a.Capabilities() {
				c.Capabilities[string(k)] = v
			}
		} else {
			for _, d := range reg.DeclaredLanguages() {
				if d.ID == id {
					c.Note = "unsupported: " + d.Note
				}
			}
		}
		byID[id] = c
		return c
	}
	for _, a := range reg.Adapters() {
		get(a.ID())
	}
	for _, p := range target.Paths {
		f := target.Files[p]
		if f.Language == "" {
			continue
		}
		c := get(f.Language)
		c.FilesTotal++
		c.FilesByClass[string(f.Class)]++
	}
	for _, fc := range cs.Files {
		if fc.Language == "" {
			continue
		}
		c := get(fc.Language)
		c.FilesChanged++
		c.LinesChanged += fc.Added + fc.Deleted
	}
	var out []report.LanguageCoverage
	for _, c := range byID {
		out = append(out, *c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Supported != out[j].Supported {
			return out[i].Supported
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func changeSummary(cs *change.Set, base, target *index.Snapshot) report.ChangeSummary {
	sum := report.ChangeSummary{Files: cs.Files, Symbols: []report.SymbolRow{}}
	if sum.Files == nil {
		sum.Files = []change.FileChange{}
	}
	for _, f := range cs.Files {
		sum.FilesChanged++
		sum.LinesAdded += f.Added
		sum.LinesDeleted += f.Deleted
	}
	for _, sc := range cs.Symbols {
		p := sc.Primary()
		row := report.SymbolRow{Status: sc.Status, Kind: p.Kind, Name: p.DisplayName(), Exported: p.Exported, SignatureChanged: sc.SignatureChanged}
		if sc.Target != nil {
			row.File, row.StartLine, row.EndLine = sc.Target.File, sc.Target.StartLine, sc.Target.EndLine
			if f := target.Files[sc.Target.File]; f != nil {
				row.Class = string(f.Class)
			}
		}
		if sc.Base != nil {
			row.BaseFile, row.BaseStartLine, row.BaseEndLine = sc.Base.File, sc.Base.StartLine, sc.Base.EndLine
			if sc.Target == nil || sc.Base.DisplayName() != sc.Target.DisplayName() {
				row.BaseName = sc.Base.DisplayName()
			}
			if row.Class == "" {
				if f := base.Files[sc.Base.File]; f != nil {
					row.Class = string(f.Class)
				}
			}
		}
		sum.Symbols = append(sum.Symbols, row)
	}
	return sum
}

func diagnostics(cs *change.Set, base, target *index.Snapshot, extra []model.Diagnostic) report.Diagnostics {
	out := report.Diagnostics{Change: append([]model.Diagnostic{}, extra...), SnapshotCountByCode: map[string]int{}}
	changedOld, changedNew := map[string]bool{}, map[string]bool{}
	for _, fc := range cs.Files {
		changedNew[fc.Path] = true
		if fc.OldPath != "" {
			changedOld[fc.OldPath] = true
		} else {
			changedOld[fc.Path] = true
		}
	}
	for _, side := range []struct {
		name    string
		snap    *index.Snapshot
		changed map[string]bool
	}{{"base", base, changedOld}, {"target", target, changedNew}} {
		for _, d := range side.snap.Diagnostics {
			out.SnapshotCountByCode[side.name+":"+d.Code]++
			if d.File == "" || side.changed[d.File] {
				d.Side = side.name
				out.Change = append(out.Change, d)
			}
		}
	}
	for _, fc := range cs.Files {
		if fc.Language != "" && !fc.Supported {
			if fc.Class == model.ClassSource || fc.Class == model.ClassTest || fc.Class == model.ClassGenerated {
				out.Change = append(out.Change, model.Diagnostic{Severity: model.SevInfo, Code: "unsupported_language", Message: "changed file in a language without an adapter is not analysed", File: fc.Path, Language: fc.Language})
			}
		}
	}
	sort.SliceStable(out.Change, func(i, j int) bool {
		a, b := out.Change[i], out.Change[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Side != b.Side {
			return a.Side < b.Side
		}
		return a.Code < b.Code
	})
	return out
}

func limitations(reg *lang.Registry, mode string) []string {
	out := []string{
		"POC: measurements of structure, not proof of design quality; lower metrics do not by themselves prove better architecture",
		"scores are per analyzer and directional; there is deliberately no weighted overall score",
		"correctness confidence is not assessed; the evaluated repository's code is never executed",
		"duplication, complexity and obsolete-code analyzers measure production source only; tests, generated, vendor and docs are counted in scope and coverage",
		"structural/near-structural similarity marks candidates, not semantic equivalence or proven missed reuse",
	}
	if mode == "worktree" {
		out = append(out, "worktree mode reads tracked files only; untracked files are never read")
	}
	for _, a := range reg.Adapters() {
		for _, l := range a.Limits() {
			out = append(out, a.ID()+": "+l)
		}
	}
	for _, d := range reg.DeclaredLanguages() {
		out = append(out, d.ID+": no adapter — files are counted as unsupported coverage, never scored")
	}
	return out
}
