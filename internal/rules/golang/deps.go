package golang

import (
	"path"
	"sort"
	"strings"

	"github.com/cameron/code-evaluator/internal/model"
	"github.com/cameron/code-evaluator/internal/rules"
)

type require struct {
	path     string
	line     int
	indirect bool
}

type module struct {
	dir, path, file string
	requires        []require
	tools           []string
}

type modules struct {
	list   []module        // longest dir first
	goDirs map[string]bool // directories with non-test .go files
}

func parseGoMod(file string, data []byte) module {
	m := module{dir: path.Dir(file), file: file}
	block := ""
	for i, raw := range strings.Split(string(data), "\n") {
		line, comment := strings.TrimSpace(raw), ""
		if j := strings.Index(line, "//"); j >= 0 {
			comment, line = line[j+2:], strings.TrimSpace(line[:j])
		}
		if line == "" {
			continue
		}
		if block != "" {
			if line == ")" {
				block = ""
			} else {
				m.add(block, line, comment, i+1)
			}
			continue
		}
		kw := strings.Fields(line)[0]
		rest := strings.TrimSpace(line[len(kw):])
		switch {
		case kw == "module":
			m.path = strings.Trim(rest, "\"")
		case rest == "(":
			block = kw
		default:
			m.add(kw, rest, comment, i+1)
		}
	}
	return m
}

func (m *module) add(kind, line, comment string, ln int) {
	f := strings.Fields(line)
	if len(f) == 0 {
		return
	}
	p := strings.Trim(f[0], "\"")
	switch kind {
	case "require":
		m.requires = append(m.requires, require{path: p, line: ln, indirect: strings.Contains(comment, "indirect")})
	case "tool":
		m.tools = append(m.tools, p)
	}
}

func loadModules(c *rules.Context) *modules {
	return c.Memo("go:modules", func() any {
		ms := &modules{goDirs: map[string]bool{}}
		for _, p := range c.AllPaths {
			if strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go") {
				ms.goDirs[path.Dir(p)] = true
			}
			if path.Base(p) != "go.mod" || strings.Contains(p, "vendor/") || strings.Contains(p, "testdata/") {
				continue
			}
			data, err := c.ReadFile(p)
			if err != nil {
				c.Diagnose(model.Diagnostic{Severity: model.SevWarning, Code: "gomod_unreadable", Message: err.Error(), File: p})
				continue
			}
			if m := parseGoMod(p, data); m.path != "" {
				ms.list = append(ms.list, m)
			}
		}
		sort.Slice(ms.list, func(i, j int) bool { return len(ms.list[i].dir) > len(ms.list[j].dir) })
		return ms
	}).(*modules)
}

// governing returns the module whose directory contains file.
func (ms *modules) governing(file string) *module {
	dir := path.Dir(file)
	for i := range ms.list {
		m := &ms.list[i]
		if m.dir == "." || dir == m.dir || strings.HasPrefix(dir, m.dir+"/") {
			return m
		}
	}
	return nil
}

// internal returns the snapshot module an import path belongs to.
func (ms *modules) internal(imp string) *module {
	var best *module
	for i := range ms.list {
		m := &ms.list[i]
		if (imp == m.path || strings.HasPrefix(imp, m.path+"/")) && (best == nil || len(m.path) > len(best.path)) {
			best = m
		}
	}
	return best
}

func covers(mod, imp string) bool { return imp == mod || strings.HasPrefix(imp, mod+"/") }

func eachImport(c *rules.Context, visit func(f *rules.File, imp model.Import)) {
	for _, f := range c.Files {
		if f.Parsed == nil {
			continue
		}
		for _, imp := range f.Parsed.Imports {
			visit(f, imp)
		}
	}
}

func unresolvedImport(c *rules.Context) {
	ms := loadModules(c)
	eachImport(c, func(f *rules.File, imp model.Import) {
		m := ms.internal(imp.Path)
		if m == nil {
			return
		}
		dir := m.dir
		if rel := strings.TrimPrefix(strings.TrimPrefix(imp.Path, m.path), "/"); rel != "" {
			dir = path.Join(m.dir, rel)
		}
		if !ms.goDirs[path.Clean(dir)] {
			c.Report(rules.Finding{File: f.Path, Line: imp.Line, Key: "import:" + imp.Path,
				Message: "import \"" + imp.Path + "\" names no package in module " + m.path + " (no non-test .go files in " + path.Clean(dir) + ")"})
		}
	})
}

func undeclaredDependency(c *rules.Context) {
	ms := loadModules(c)
	eachImport(c, func(f *rules.File, imp model.Import) {
		if stdlib(imp.Path) || ms.internal(imp.Path) != nil {
			return
		}
		m := ms.governing(f.Path)
		if m == nil {
			return
		}
		for _, r := range m.requires {
			if covers(r.path, imp.Path) {
				return
			}
		}
		c.Report(rules.Finding{File: f.Path, Line: imp.Line, Key: "import:" + imp.Path,
			Message: "import \"" + imp.Path + "\" is not required by " + m.file})
	})
}

func unusedDependency(c *rules.Context) {
	ms := loadModules(c)
	imported := map[string][]string{} // go.mod -> imports
	eachImport(c, func(f *rules.File, imp model.Import) {
		if m := ms.governing(f.Path); m != nil {
			imported[m.file] = append(imported[m.file], imp.Path)
		}
	})
	for _, m := range ms.list {
		if len(imported[m.file]) == 0 {
			continue // no checked Go files: nothing to compare against
		}
	next:
		for _, r := range m.requires {
			if r.indirect {
				continue
			}
			for _, imp := range imported[m.file] {
				if covers(r.path, imp) {
					continue next
				}
			}
			for _, t := range m.tools {
				if covers(r.path, t) {
					continue next
				}
			}
			c.Report(rules.Finding{File: m.file, Line: r.line, Key: "require:" + r.path,
				Message: "go.mod requires " + r.path + " directly, but no .go file in the module imports it"})
		}
	}
}
