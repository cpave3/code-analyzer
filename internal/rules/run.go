package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
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
	Findings    []Finding          `json:"findings"`
	Coverage    []Coverage         `json:"coverage"`
	Files       []LanguageFiles    `json:"files"`
	Diagnostics []model.Diagnostic `json:"diagnostics,omitempty"`
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

// loadLanguage parses the syntax trees of one language's checkable files.
func (res *Result) loadLanguage(snap *index.Snapshot, ad lang.Adapter) *langFiles {
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
		if rf := res.parseFile(snap, f, sp, ad.ID()); rf != nil {
			lf.files = append(lf.files, rf)
			if f.Class != model.ClassGenerated {
				lf.count.Checked++
			}
		}
	}
	return lf
}

func (res *Result) parseFile(snap *index.Snapshot, f *index.File, sp lang.SyntaxParser, language string) *File {
	src, err := snap.Content(f.Path)
	if err != nil {
		res.Diagnostics = append(res.Diagnostics, model.Diagnostic{Severity: model.SevWarning, Code: "rules_unreadable", Message: err.Error(), File: f.Path, Language: language})
		return nil
	}
	tree, err := sp.ParseSyntax(src)
	if err != nil {
		res.Diagnostics = append(res.Diagnostics, model.Diagnostic{Severity: model.SevWarning, Code: "rules_parse_failed", Message: err.Error(), File: f.Path, Language: language})
		return nil
	}
	return &File{Path: f.Path, Class: f.Class, Parsed: f.Parsed, Src: src, Tree: tree}
}

// Run checks a snapshot with every selected rule × language implementation.
func Run(reg *Registry, snap *index.Snapshot, cfg *config.Config, langs *lang.Registry, opt Options) *Result {
	res := &Result{}
	byLang := map[string]*langFiles{}
	for _, ad := range langs.Adapters() {
		lf := res.loadLanguage(snap, ad)
		byLang[ad.ID()] = lf
		res.Files = append(res.Files, lf.count)
	}
	defer func() {
		for _, lf := range byLang {
			for _, f := range lf.files {
				f.Tree.Close()
			}
		}
	}()
	var disabled []string
	if cfg != nil {
		disabled = cfg.Defects.Disable
	}
	memo := map[string]map[string]any{}
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
				if memo[ad.ID()] == nil {
					memo[ad.ID()] = map[string]any{}
				}
				res.check(snap, cfg, rule, impl, byLang[ad.ID()], memo[ad.ID()])
				cov.Status, cov.Files, cov.Note = "checked", byLang[ad.ID()].count.Checked, impl.Note
			}
			res.Coverage = append(res.Coverage, cov)
		}
	}
	res.Findings = inScope(res.Findings, opt.InScope)
	SortFindings(res.Findings)
	return res
}

// check runs one implementation and keeps its findings outside generated code.
func (res *Result) check(snap *index.Snapshot, cfg *config.Config, rule Rule, impl Impl, lf *langFiles, memo map[string]any) {
	ctx := &Context{Rule: rule, Language: impl.Language, Files: lf.files, AllPaths: snap.Paths, ReadFile: snap.Content,
		Graph: snap.Graphs[impl.Language], Config: cfg, cache: memo}
	if err := safeCheck(impl, ctx); err != nil {
		res.Diagnostics = append(res.Diagnostics, model.Diagnostic{Severity: model.SevError, Code: "rule_failed", Message: err.Error(), Analyzer: "rules:" + rule.ID, Language: impl.Language})
	}
	for _, f := range ctx.findings {
		if g := snap.Files[f.File]; g != nil && g.Class == model.ClassGenerated {
			continue // generated code is context, not checked
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
