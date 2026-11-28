// Package index builds an immutable snapshot of a repository revision: every
// tracked file classified by role and language, parse results for supported
// languages (via the SQLite cache), and per-language dependency graphs.
package index

import (
	"fmt"
	"path"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/gitsrc"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
)

// MaxParseBytes bounds files given to adapters.
const MaxParseBytes = 2 << 20

// File is one tracked file in a snapshot.
type File struct {
	Path     string
	Language string       // "" when unknown
	Adapter  lang.Adapter // nil when unsupported
	Class    model.FileClass
	Blob     string
	Mode     string
	Size     int64
	Parsed   *model.ParsedFile
	Skipped  string // reason a supported file was not parsed
}

// Supported reports whether the file has a parse result.
func (f *File) Supported() bool { return f.Adapter != nil && f.Parsed != nil }

// Snapshot is an indexed revision.
type Snapshot struct {
	Commit      string // full SHA or gitsrc.WorktreeRev
	Tree        string // tree SHA, or content fingerprint for the worktree
	Files       map[string]*File
	Paths       []string // sorted
	Symbols     []*model.Symbol
	Graphs      map[string]*lang.Graph
	Diagnostics []model.Diagnostic

	repo     *gitsrc.Repo
	worktree map[string][]byte
	byFile   map[string][]*model.Symbol
}

// Content returns the bytes of a file in the snapshot.
func (s *Snapshot) Content(p string) ([]byte, error) {
	f := s.Files[p]
	if f == nil {
		return nil, fmt.Errorf("%s not in snapshot", p)
	}
	if s.worktree != nil {
		if b, ok := s.worktree[f.Blob]; ok {
			return b, nil
		}
	}
	return s.repo.ReadBlob(f.Blob)
}

// SymbolsIn returns the symbols of one file.
func (s *Snapshot) SymbolsIn(p string) []*model.Symbol {
	return s.byFile[p]
}

// Builder builds snapshots.
type Builder struct {
	Repo   *gitsrc.Repo
	Langs  *lang.Registry
	Config *config.Config
	Cache  *Cache
	Parsed int // files parsed (cache misses) across all builds
}

// Commit indexes a commit.
func (b *Builder) Commit(sha string) (*Snapshot, error) {
	entries, err := b.Repo.ListTree(sha)
	if err != nil {
		return nil, err
	}
	tree, err := b.Repo.TreeOf(sha)
	if err != nil {
		return nil, err
	}
	s := &Snapshot{Commit: sha, Tree: tree, repo: b.Repo}
	return s, b.build(s, entries)
}

// Worktree indexes the working tree: tracked files, plus untracked files
// that are not ignored when untracked is set.
func (b *Builder) Worktree(untracked bool) (*Snapshot, error) {
	entries, contents, err := b.Repo.ListWorktree(untracked)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		ids = append(ids, e.Path+"\x00"+e.Blob)
	}
	sort.Strings(ids)
	s := &Snapshot{Commit: gitsrc.WorktreeRev, Tree: "worktree:" + gitsrc.BlobID([]byte(strings.Join(ids, "\n")))[:16], repo: b.Repo, worktree: contents}
	return s, b.build(s, entries)
}

func (b *Builder) build(s *Snapshot, entries []gitsrc.Entry) error {
	s.Files = map[string]*File{}
	s.Graphs = map[string]*lang.Graph{}
	for _, e := range entries {
		if e.Mode == "160000" || e.Mode == "120000" {
			continue // submodule / symlink: not content we analyse
		}
		id, ad := b.Langs.Detect(e.Path)
		f := &File{Path: e.Path, Language: id, Adapter: ad, Blob: e.Blob, Mode: e.Mode, Size: e.Size}
		s.Files[e.Path] = f
		s.Paths = append(s.Paths, e.Path)
	}
	sort.Strings(s.Paths)

	// Parse supported files: cache first, then parse misses in parallel.
	type job struct {
		f   *File
		src []byte
	}
	var jobs []job
	for _, p := range s.Paths {
		f := s.Files[p]
		if f.Adapter == nil || preClass(b.Config, f.Path) == model.ClassVendor || preClass(b.Config, f.Path) == model.ClassExcluded {
			continue
		}
		if f.Size > MaxParseBytes {
			f.Skipped = "file larger than parse limit"
			s.Diagnostics = append(s.Diagnostics, model.Diagnostic{Severity: model.SevWarning, Code: "file_too_large", Message: f.Skipped, File: f.Path, Language: f.Language})
			continue
		}
		if pf, ok := b.Cache.Get(f.Adapter.ID(), f.Adapter.Version(), f.Blob); ok {
			f.Parsed = pf
			continue
		}
		src, err := s.Content(f.Path)
		if err != nil {
			f.Skipped = "unreadable: " + err.Error()
			s.Diagnostics = append(s.Diagnostics, model.Diagnostic{Severity: model.SevError, Code: "file_unreadable", Message: err.Error(), File: f.Path, Language: f.Language})
			continue
		}
		jobs = append(jobs, job{f, src})
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var puts []PutEntry
	ch := make(chan job)
	for w := 0; w < runtime.GOMAXPROCS(0); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range ch {
				pf, err := safeParse(j.f.Adapter, j.src)
				mu.Lock()
				if err != nil {
					j.f.Skipped = "adapter failure: " + err.Error()
					s.Diagnostics = append(s.Diagnostics, model.Diagnostic{Severity: model.SevError, Code: "adapter_failure", Message: err.Error(), File: j.f.Path, Language: j.f.Language})
				} else {
					j.f.Parsed = pf
					puts = append(puts, PutEntry{Adapter: j.f.Adapter.ID(), Version: j.f.Adapter.Version(), Blob: j.f.Blob, Parsed: pf})
				}
				mu.Unlock()
			}
		}()
	}
	for _, j := range jobs {
		ch <- j
	}
	close(ch)
	wg.Wait()
	b.Parsed += len(jobs)
	sort.Slice(puts, func(i, j int) bool { return puts[i].Blob < puts[j].Blob })
	if err := b.Cache.PutAll(puts); err != nil {
		s.Diagnostics = append(s.Diagnostics, model.Diagnostic{Severity: model.SevWarning, Code: "cache_write_failed", Message: err.Error()})
	}

	// Classify, attach symbols, record parse errors.
	for _, p := range s.Paths {
		f := s.Files[p]
		f.Class = classify(b.Config, f)
		if f.Parsed == nil {
			continue
		}
		for _, d := range f.Parsed.ParseErrors {
			d.File = f.Path
			s.Diagnostics = append(s.Diagnostics, d)
		}
		pkg := f.Adapter.PackageOf(f.Path)
		for i := range f.Parsed.Symbols {
			sym := &f.Parsed.Symbols[i]
			sym.File, sym.Package = f.Path, pkg
			sym.Key = f.Language + ":" + pkg + "|" + f.Parsed.PackageName + "|" + sym.DisplayName()
			s.Symbols = append(s.Symbols, sym)
		}
	}
	// Disambiguate duplicate keys (e.g. build-tag variants) deterministically.
	count := map[string]int{}
	for _, sym := range s.Symbols {
		count[sym.Key]++
	}
	for _, sym := range s.Symbols {
		if count[sym.Key] > 1 {
			sym.Key += "@" + path.Base(sym.File)
		}
	}

	s.byFile = map[string][]*model.Symbol{}
	for _, sym := range s.Symbols {
		s.byFile[sym.File] = append(s.byFile[sym.File], sym)
	}

	// Dependency graphs per adapter with the capability.
	for _, ad := range b.Langs.Adapters() {
		if !lang.Has(ad, lang.CapDependencies) {
			continue
		}
		var files []lang.GraphFile
		for _, p := range s.Paths {
			f := s.Files[p]
			if f.Adapter == ad && f.Parsed != nil {
				files = append(files, lang.GraphFile{Path: f.Path, Class: f.Class, Parsed: f.Parsed})
			}
		}
		if len(files) == 0 {
			continue
		}
		g, diags := ad.BuildGraph(lang.GraphInput{Files: files, AllPaths: s.Paths, ReadFile: s.Content})
		s.Graphs[ad.ID()] = g
		s.Diagnostics = append(s.Diagnostics, diags...)
	}
	return nil
}

func safeParse(a lang.Adapter, src []byte) (pf *model.ParsedFile, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return a.Parse(src)
}

// preClass classifies by path alone (before parsing).
func preClass(c *config.Config, p string) model.FileClass {
	for _, g := range c.Exclude {
		if config.Match(g, p) {
			return model.ClassExcluded
		}
	}
	seg := "/" + p
	for _, v := range []string{"/vendor/", "/node_modules/", "/third_party/"} {
		if strings.Contains(seg, v) {
			return model.ClassVendor
		}
	}
	return ""
}

var buildFiles = map[string]bool{
	"go.mod": true, "go.sum": true, "go.work": true, "go.work.sum": true, "package.json": true,
	"package-lock.json": true, "pnpm-lock.yaml": true, "yarn.lock": true, "Makefile": true,
	"Taskfile.yml": true, "Dockerfile": true, "tsconfig.json": true, "Cargo.toml": true, "Cargo.lock": true,
}

var buildExts = map[string]bool{".csproj": true, ".sln": true, ".props": true, ".targets": true, ".gradle": true}

var docExts = map[string]bool{".md": true, ".mdx": true, ".rst": true, ".txt": true, ".adoc": true}

func classify(c *config.Config, f *File) model.FileClass {
	if pc := preClass(c, f.Path); pc != "" {
		return pc
	}
	seg := "/" + f.Path
	if strings.Contains(seg, "/testdata/") {
		return model.ClassTestdata
	}
	base := path.Base(f.Path)
	if f.Adapter != nil && f.Adapter.IsTestPath(f.Path) {
		return model.ClassTest
	}
	if f.Adapter == nil && (strings.Contains(base, ".test.") || strings.Contains(base, ".spec.") ||
		strings.Contains(seg, "/__tests__/") || strings.HasSuffix(strings.TrimSuffix(base, path.Ext(base)), "Tests")) {
		return model.ClassTest
	}
	if f.Parsed != nil && f.Parsed.Generated {
		return model.ClassGenerated
	}
	if strings.HasSuffix(base, ".pb.go") || strings.Contains(base, ".generated.") || strings.Contains(base, "_generated.") || strings.Contains(seg, "/generated/") {
		return model.ClassGenerated
	}
	if buildFiles[base] || buildExts[path.Ext(base)] || strings.HasSuffix(base, ".gradle.kts") {
		return model.ClassBuild
	}
	if docExts[strings.ToLower(path.Ext(base))] || strings.HasPrefix(f.Path, "docs/") {
		return model.ClassDocs
	}
	if f.Language != "" {
		return model.ClassSource
	}
	return model.ClassOther
}
