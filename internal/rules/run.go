package rules

import (
	"fmt"
	"sync"

	sitter "github.com/tree-sitter/go-tree-sitter"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/memo"
	"github.com/cameron/code-evaluator/internal/model"
)

// Options select what a run checks.
type Options struct {
	// Only keeps rules whose ID equals or starts with one of these
	// (e.g. "dead-code" or "dead-code/unused-import"). Empty = all.
	Only []string
	// InScope filters findings by file path after checking (rules always
	// see the whole snapshot, e.g. to know whether a function is used).
	InScope func(path string) bool
	// Cache keeps per-file results across runs (e.g. the commits of a
	// series); nil caches within this run only.
	Cache *memo.Cache
}

// Coverage records what one language implementation of a rule checked.
type Coverage struct {
	Rule     string `json:"rule"`
	Language string `json:"language"`
	Status   string `json:"status"` // checked | not_implemented | disabled
	Files    int    `json:"files"`
	Note     string `json:"note,omitempty"`
}

// LanguageFiles counts files per language.
type LanguageFiles struct {
	Language string `json:"language"`
	Checked  int    `json:"checked"`
	// Checked counts source and test files. Generated files are read as
	// context (a function used only by generated code is used) but never
	// reported on.
	// SkippedSyntaxErrors files are not checked: rules on a broken tree
	// would report noise.
	SkippedSyntaxErrors int `json:"skippedSyntaxErrors"`
}

// Result is the outcome of running rules on one snapshot.
type Result struct {
	Findings []Finding `json:"findings"`
	// Suppressed findings matched a codeeval-ignore comment; they are kept
	// for audit and never counted or compared.
	Suppressed  []Finding          `json:"suppressed,omitempty"`
	Coverage    []Coverage         `json:"coverage"`
	Files       []LanguageFiles    `json:"files"`
	Diagnostics []model.Diagnostic `json:"diagnostics,omitempty"`

	mu sync.Mutex
}

func selected(id string, only, disabled []string) (bool, bool) {
	for _, d := range disabled {
		if id == d || strings.HasPrefix(id, d+"/") {
			return false, true
		}
	}
	if len(only) == 0 {
		return true, false
	}
	for _, o := range only {
		if id == o || strings.HasPrefix(id, o+"/") {
			return true, false
		}
	}
	return false, false
}

type langFiles struct {
	files []*File
	count LanguageFiles
}

// checkable reports whether a snapshot file is given to rules: parsed
// source, test or generated (context only) files.
func checkable(f *index.File, adapter string) bool {
	if f.Adapter == nil || f.Adapter.ID() != adapter || f.Parsed == nil {
		return false
	}
	return f.Class == model.ClassSource || f.Class == model.ClassTest || f.Class == model.ClassGenerated
}

// loadLanguage lists one language's checkable files. Content and syntax
// trees load lazily, on first use.
func (res *Result) loadLanguage(snap *index.Snapshot, ad lang.Adapter, cache *memo.Cache) *langFiles {
	lf := &langFiles{count: LanguageFiles{Language: ad.ID()}}
	sp, ok := ad.(lang.SyntaxParser)
	if !ok {
		res.Diagnostics = append(res.Diagnostics, model.Diagnostic{Severity: model.SevInfo, Code: "rules_no_syntax", Message: "adapter exposes no syntax tree; syntax rules cannot run", Language: ad.ID()})
		return lf
	}
	for _, p := range snap.Paths {
		f := snap.Files[p]
		if !checkable(f, ad.ID()) {
			continue
		}
		if len(f.Parsed.ParseErrors) > 0 {
			lf.count.SkippedSyntaxErrors++
			continue
		}
		rf := &File{Path: f.Path, Blob: f.Blob, Class: f.Class, Parsed: f.Parsed, Language: ad.ID(), cache: cache}
		path, language := f.Path, ad.ID()
		rf.load = func() ([]byte, *sitter.Tree, error) {
			src, err := snap.Content(path)
			if err != nil {
				res.diag(model.Diagnostic{Severity: model.SevWarning, Code: "rules_unreadable", Message: err.Error(), File: path, Language: language})
				return nil, nil, err
			}
			tree, err := sp.ParseSyntax(src)
			if err != nil {
				res.diag(model.Diagnostic{Severity: model.SevWarning, Code: "rules_parse_failed", Message: err.Error(), File: path, Language: language})
				return src, nil, err
			}
			return src, tree, nil
		}
		lf.files = append(lf.files, rf)
		if f.Class != model.ClassGenerated {
			lf.count.Checked++
		}
	}
	return lf
}

func (res *Result) diag(d model.Diagnostic) {
	res.mu.Lock()
	res.Diagnostics = append(res.Diagnostics, d)
	res.mu.Unlock()
}

// Run checks a snapshot with every selected rule × language implementation.
//
// Per language it makes one pass over the files: file-local checks run file
// by file (their findings cached by content when opt.Cache is set), the
// language's Prepare steps warm the per-file facts cross-file rules read,
// and the file's syntax tree is released. Cross-file checks then run on
// those facts.
func Run(reg *Registry, snap *index.Snapshot, cfg *config.Config, langs *lang.Registry, opt Options) *Result {
	res := &Result{}
	var disabled []string
	if cfg != nil {
		disabled = cfg.Defects.Disable
	}
	type plan struct{ local, global []Impl }
	plans := map[string]*plan{}
	byLang := map[string]*langFiles{}
	for _, ad := range langs.Adapters() {
		lf := res.loadLanguage(snap, ad, opt.Cache)
		byLang[ad.ID()] = lf
		res.Files = append(res.Files, lf.count)
		plans[ad.ID()] = &plan{}
	}
	for _, rule := range reg.Rules() {
		run, off := selected(rule.ID, opt.Only, disabled)
		for _, ad := range langs.Adapters() {
			cov := Coverage{Rule: rule.ID, Language: ad.ID()}
			impl, ok := reg.Impl(rule.ID, ad.ID())
			switch {
			case !ok:
				cov.Status = "not_implemented"
			case off:
				cov.Status = "disabled"
			case !run:
				continue
			default:
				p := plans[ad.ID()]
				if impl.FileLocal {
					p.local = append(p.local, impl)
				} else {
					p.global = append(p.global, impl)
				}
				cov.Status, cov.Files, cov.Note = "checked", byLang[ad.ID()].count.Checked, impl.Note
			}
			res.Coverage = append(res.Coverage, cov)
		}
	}
	for _, ad := range langs.Adapters() {
		lf, p := byLang[ad.ID()], plans[ad.ID()]
		for _, f := range lf.files {
			if f.Class != model.ClassGenerated {
				for _, impl := range p.local {
					res.checkFile(reg, f, impl, cfg)
				}
			}
			if len(p.global) > 0 {
				suppressionLines(f)
				for _, prep := range reg.prep[ad.ID()] {
					prep(f)
				}
			}
			f.release()
		}
		memos := map[string]any{}
		for _, impl := range p.global {
			rule, _ := reg.Rule(impl.Rule)
			res.check(snap, cfg, rule, impl, lf, memos)
		}
		for _, f := range lf.files {
			f.release()
		}
	}
	res.Findings = inScope(res.Findings, opt.InScope)
	res.Suppressed = inScope(res.Suppressed, opt.InScope)
	SortFindings(res.Findings)
	SortFindings(res.Suppressed)
	return res
}

// fileResult is a file-local check's outcome for one file content.
type fileResult struct{ kept, suppressed []Finding }

// checkFile runs a file-local check on one file, through the content cache.
func (res *Result) checkFile(reg *Registry, f *File, impl Impl, cfg *config.Config) {
	rule, _ := reg.Rule(impl.Rule)
	r := f.Memo("findings:"+Version+":"+impl.Rule, func() any {
		ctx := &Context{Rule: rule, Language: impl.Language, Files: []*File{f}, Config: cfg, cache: map[string]any{}}
		if err := safeCheck(impl, ctx); err != nil {
			res.diag(model.Diagnostic{Severity: model.SevError, Code: "rule_failed", Message: err.Error(), Analyzer: "rules:" + rule.ID, Language: impl.Language, File: f.Path})
		}
		var out fileResult
		lines := suppressionLines(f)
		for _, x := range ctx.findings {
			if reason, ok := suppressionAt(lines, x.Line, x.Rule); ok {
				x.Message += " [suppressed: " + reason + "]"
				out.suppressed = append(out.suppressed, x)
				continue
			}
			out.kept = append(out.kept, x)
		}
		return out
	}).(fileResult)
	for _, x := range r.kept {
		x.File = f.Path
		res.Findings = append(res.Findings, x)
	}
	for _, x := range r.suppressed {
		x.File = f.Path
		res.Suppressed = append(res.Suppressed, x)
	}
}

// check runs one cross-file implementation and keeps its findings outside
// generated code.
func (res *Result) check(snap *index.Snapshot, cfg *config.Config, rule Rule, impl Impl, lf *langFiles, memo map[string]any) {
	ctx := &Context{Rule: rule, Language: impl.Language, Files: lf.files, AllPaths: snap.Paths, ReadFile: snap.Content,
		Graph: snap.Graphs[impl.Language], Config: cfg, cache: memo}
	if err := safeCheck(impl, ctx); err != nil {
		res.diag(model.Diagnostic{Severity: model.SevError, Code: "rule_failed", Message: err.Error(), Analyzer: "rules:" + rule.ID, Language: impl.Language})
	}
	byPath := map[string]*File{}
	for _, f := range lf.files {
		byPath[f.Path] = f
	}
	for _, f := range ctx.findings {
		if g := snap.Files[f.File]; g != nil && g.Class == model.ClassGenerated {
			continue // generated code is context, not checked
		}
		if file := byPath[f.File]; file != nil {
			if reason, ok := suppressionAt(suppressionLines(file), f.Line, f.Rule); ok {
				f.Message += " [suppressed: " + reason + "]"
				res.Suppressed = append(res.Suppressed, f)
				continue
			}
		}
		res.Findings = append(res.Findings, f)
	}
	res.Diagnostics = append(res.Diagnostics, ctx.diags...)
}

func inScope(fs []Finding, keep func(string) bool) []Finding {
	if keep == nil {
		return fs
	}
	out := fs[:0]
	for _, f := range fs {
		if keep(f.File) {
			out = append(out, f)
		}
	}
	return out
}

func safeCheck(impl Impl, ctx *Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("rule %s (%s) panicked: %v", impl.Rule, impl.Language, r)
		}
	}()
	impl.Check(ctx)
	return nil
}

// SortFindings orders findings by severity (gate first), file, line, rule.
func SortFindings(fs []Finding) {
	sort.SliceStable(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		if a.Severity != b.Severity {
			return a.Severity == Gate
		}
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Rule != b.Rule {
			return a.Rule < b.Rule
		}
		return a.Key < b.Key
	})
}

func identity(f Finding) string { return f.Rule + "\x00" + f.File + "\x00" + f.Key }

// Compare matches findings across two revisions by rule, file and key (not
// line), as multisets. Introduced findings are in after but not before;
// resolved ones the reverse. A finding in a renamed or moved file counts as
// resolved in the old path and introduced in the new one.
func Compare(before, after []Finding) (introduced, resolved []Finding, persisting int) {
	count := map[string]int{}
	for _, f := range before {
		count[identity(f)]++
	}
	for _, f := range after {
		id := identity(f)
		if count[id] > 0 {
			count[id]--
			persisting++
			continue
		}
		f.Change = "introduced"
		introduced = append(introduced, f)
	}
	left := map[string]int{}
	for _, f := range after {
		left[identity(f)]++
	}
	for _, f := range before {
		id := identity(f)
		if left[id] > 0 {
			left[id]--
			continue
		}
		f.Change = "resolved"
		resolved = append(resolved, f)
	}
	SortFindings(introduced)
	SortFindings(resolved)
	return introduced, resolved, persisting
}
