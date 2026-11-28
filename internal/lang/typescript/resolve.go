package typescript

import (
	"encoding/json"
	"path"
	"sort"
	"strings"
)

// Resolution kinds for an import specifier.
const (
	KindRelative  = "relative"  // ./ or ../
	KindAlias     = "alias"     // tsconfig "paths"
	KindWorkspace = "workspace" // package.json "name" of a package in the repository
	KindExternal  = "external"  // a package outside the repository (node_modules)
	KindAsset     = "asset"     // stylesheet, image, JSON, query-suffixed import ...
)

// Resolution is where an import specifier points. Target is the resolved
// TypeScript file, or "" when an internal specifier resolves to no indexed
// TypeScript file. Package is the package name for external specifiers.
type Resolution struct {
	Kind    string
	Target  string
	Package string
}

// Resolver resolves import specifiers the way the dependency graph does:
// relative paths, tsconfig "paths" aliases (applied to files under the
// config's directory) and workspace package names. It never reads
// node_modules or package "exports" maps.
type Resolver struct {
	exists  map[string]bool
	configs []tsconfig
	pkgs    map[string]string // workspace package name -> directory
}

// NewResolver indexes tsconfig files and workspace packages of a snapshot.
func NewResolver(allPaths []string, read func(string) ([]byte, error)) *Resolver {
	r := &Resolver{exists: map[string]bool{}, pkgs: map[string]string{}}
	for _, p := range allPaths {
		r.exists[p] = true
	}
	for _, p := range allPaths {
		if strings.Contains(p, "node_modules/") {
			continue
		}
		base := path.Base(p)
		if strings.HasPrefix(base, "tsconfig") && strings.HasSuffix(base, ".json") {
			baseDir, paths := loadTSConfig(p, read, 0)
			if len(paths) == 0 {
				continue
			}
			c := tsconfig{dir: path.Dir(p)}
			keys := make([]string, 0, len(paths))
			for k := range paths {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				al := alias{}
				if i := strings.Index(k, "*"); i >= 0 {
					al.prefix, al.suffix, al.wild = k[:i], k[i+1:], true
				} else {
					al.prefix = k
				}
				for _, t := range paths[k] {
					al.targets = append(al.targets, path.Clean(path.Join(baseDir, t)))
				}
				c.aliases = append(c.aliases, al)
			}
			r.configs = append(r.configs, c)
		}
		if base == "package.json" {
			data, err := read(p)
			if err != nil {
				continue
			}
			var pj struct {
				Name string `json:"name"`
			}
			if json.Unmarshal(data, &pj) == nil && pj.Name != "" {
				r.pkgs[pj.Name] = path.Dir(p)
			}
		}
	}
	sort.Slice(r.configs, func(i, j int) bool { return len(r.configs[i].dir) > len(r.configs[j].dir) })
	return r
}

// Exists reports whether a path is in the snapshot.
func (r *Resolver) Exists(p string) bool { return r.exists[p] }

// WorkspacePackage reports whether name is a package defined in the repository.
func (r *Resolver) WorkspacePackage(name string) bool { _, ok := r.pkgs[name]; return ok }

func (r *Resolver) resolveFile(cand string) string {
	cand = path.Clean(cand)
	if r.exists[cand] && isTS(cand) {
		return cand
	}
	stem := cand
	for _, js := range []string{".js", ".jsx", ".mjs", ".cjs"} {
		if strings.HasSuffix(cand, js) {
			stem = strings.TrimSuffix(cand, js)
		}
	}
	for _, e := range tsExts {
		if r.exists[stem+e] {
			return stem + e
		}
	}
	for _, e := range tsExts {
		if r.exists[cand+"/index"+e] {
			return cand + "/index" + e
		}
	}
	return ""
}

// IsAsset reports whether a specifier names a non-code asset.
func IsAsset(spec string) bool {
	switch path.Ext(spec) {
	case ".css", ".scss", ".sass", ".less", ".json", ".svg", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".woff", ".woff2", ".ttf", ".md", ".html", ".wasm", ".txt":
		return true
	}
	return strings.Contains(spec, "?")
}

// IsRelative reports whether a specifier is a relative path.
func IsRelative(spec string) bool {
	return strings.HasPrefix(spec, "./") || strings.HasPrefix(spec, "../") || spec == "." || spec == ".."
}

// Resolve resolves spec imported from file.
func (r *Resolver) Resolve(file, spec string) Resolution {
	from := path.Dir(file)
	if IsRelative(spec) {
		if IsAsset(spec) {
			return Resolution{Kind: KindAsset}
		}
		return Resolution{Kind: KindRelative, Target: r.resolveFile(path.Join(from, spec))}
	}
	for _, c := range r.configs {
		if c.dir != "." && !strings.HasPrefix(file, c.dir+"/") {
			continue
		}
		for _, al := range c.aliases {
			var star string
			if al.wild {
				if !strings.HasPrefix(spec, al.prefix) || !strings.HasSuffix(spec, al.suffix) || len(spec) < len(al.prefix)+len(al.suffix) {
					continue
				}
				star = spec[len(al.prefix) : len(spec)-len(al.suffix)]
			} else if spec != al.prefix {
				continue
			}
			external := true
			for _, t := range al.targets {
				if !strings.Contains(t, "node_modules/") {
					external = false
				}
			}
			if external {
				// An alias into node_modules (e.g. react → node_modules/@types/react)
				// names an external package, not internal code.
				break
			}
			res := Resolution{Kind: KindAlias}
			for _, t := range al.targets {
				if res.Target = r.resolveFile(strings.Replace(t, "*", star, 1)); res.Target != "" {
					break
				}
			}
			return res
		}
	}
	name, sub := packageName(spec)
	if dir, ok := r.pkgs[name]; ok {
		res := Resolution{Kind: KindWorkspace, Package: name}
		for _, cand := range []string{path.Join(dir, sub), path.Join(dir, "src", sub)} {
			if res.Target = r.resolveFile(cand); res.Target != "" {
				break
			}
		}
		return res
	}
	if IsAsset(spec) {
		return Resolution{Kind: KindAsset}
	}
	return Resolution{Kind: KindExternal, Package: name}
}

// PackageName splits a bare specifier into package name and subpath.
func PackageName(spec string) (string, string) { return packageName(spec) }
