// Package fixturerepo materialises the paired fixture scenarios under
// testdata/fixtures into real Git repositories with deterministic commits.
//
// Layout of a scenario directory:
//
//	<scenario>/base/...            full base tree
//	<scenario>/<variant>/...       overlay: files replace or add to base
//	<variant>/_delete.txt          paths to delete (one per line)
//	<variant>/_rename.txt          "old new" per line: move content unchanged
//	<variant or base>/_nocompile.txt  reason the tree intentionally does not build
//	<scenario>/codeeval.json       optional evaluation config for the scenario
//	<scenario>/README.md           what the pair demonstrates (not copied)
//
// The repository gets tag "base" and one tag per variant, each a child of base.
package fixturerepo

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Scenario is one fixture scenario.
type Scenario struct {
	Name     string
	Dir      string
	Variants []string
	Config   string // path to codeeval.json, or ""
}

var special = map[string]bool{"_delete.txt": true, "_rename.txt": true, "_nocompile.txt": true}

// Load lists scenarios under root.
func Load(root string) ([]Scenario, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var out []Scenario
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		s := Scenario{Name: e.Name(), Dir: filepath.Join(root, e.Name())}
		sub, err := os.ReadDir(s.Dir)
		if err != nil {
			return nil, err
		}
		hasBase := false
		for _, v := range sub {
			switch {
			case v.IsDir() && v.Name() == "base":
				hasBase = true
			case v.IsDir():
				s.Variants = append(s.Variants, v.Name())
			case v.Name() == "codeeval.json":
				s.Config = filepath.Join(s.Dir, v.Name())
			}
		}
		if !hasBase {
			return nil, fmt.Errorf("scenario %s has no base/", s.Name)
		}
		sort.Strings(s.Variants)
		out = append(out, s)
	}
	return out, nil
}

// NoCompile returns the documented reason a variant ("base" for the base)
// intentionally does not build, or "".
func (s Scenario) NoCompile(variant string) string {
	data, err := os.ReadFile(filepath.Join(s.Dir, variant, "_nocompile.txt"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// Materialize writes the base tree plus a variant overlay into dst.
// variant "base" writes the base only.
func (s Scenario) Materialize(variant, dst string) error {
	if err := copyTree(filepath.Join(s.Dir, "base"), dst); err != nil {
		return err
	}
	if variant == "base" {
		return nil
	}
	vdir := filepath.Join(s.Dir, variant)
	for _, line := range readLines(filepath.Join(vdir, "_rename.txt")) {
		f := strings.Fields(line)
		if len(f) != 2 {
			return fmt.Errorf("%s: bad rename line %q", s.Name, line)
		}
		to := filepath.Join(dst, filepath.FromSlash(f[1]))
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			return err
		}
		if err := os.Rename(filepath.Join(dst, filepath.FromSlash(f[0])), to); err != nil {
			return err
		}
	}
	for _, line := range readLines(filepath.Join(vdir, "_delete.txt")) {
		if err := os.Remove(filepath.Join(dst, filepath.FromSlash(line))); err != nil {
			return err
		}
	}
	return copyTree(vdir, dst)
}

func readLines(p string) []string {
	f, err := os.Open(p)
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if l := strings.TrimSpace(sc.Text()); l != "" && !strings.HasPrefix(l, "#") {
			out = append(out, l)
		}
	}
	return out
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		if special[d.Name()] {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), data, 0o644)
	})
}

// git runs git with a deterministic identity and dates.
func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "-c", "commit.gpgsign=false", "-c", "tag.gpgsign=false", "-c", "core.hooksPath=/dev/null", "-c", "init.defaultBranch=main"}, args...)...)
	cmd.Dir = dir
	env := []string{}
	for _, kv := range os.Environ() {
		if !strings.HasPrefix(kv, "GIT_") {
			env = append(env, kv)
		}
	}
	cmd.Env = append(env, "GIT_AUTHOR_DATE=2026-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2026-01-01T00:00:00Z", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %v: %v: %s", args, err, stderr.String())
	}
	return strings.TrimSpace(string(out)), nil
}

// BuildRepo creates a Git repository at dst with tag "base" and one tag per
// variant. It returns tag -> commit SHA.
func (s Scenario) BuildRepo(dst string) (map[string]string, error) {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return nil, err
	}
	if _, err := git(dst, "init", "-q"); err != nil {
		return nil, err
	}
	refs := map[string]string{}
	commitTree := func(variant string) error {
		// Replace the working tree content with the variant.
		entries, _ := os.ReadDir(dst)
		for _, e := range entries {
			if e.Name() != ".git" {
				if err := os.RemoveAll(filepath.Join(dst, e.Name())); err != nil {
					return err
				}
			}
		}
		if err := s.Materialize(variant, dst); err != nil {
			return err
		}
		if _, err := git(dst, "add", "-A"); err != nil {
			return err
		}
		if _, err := git(dst, "commit", "-q", "--allow-empty", "-m", s.Name+": "+variant); err != nil {
			return err
		}
		if _, err := git(dst, "tag", variant); err != nil {
			return err
		}
		sha, err := git(dst, "rev-parse", "HEAD")
		refs[variant] = sha
		return err
	}
	if err := commitTree("base"); err != nil {
		return nil, err
	}
	for _, v := range s.Variants {
		if _, err := git(dst, "checkout", "-q", "--detach", "base"); err != nil {
			return nil, err
		}
		if err := commitTree(v); err != nil {
			return nil, err
		}
	}
	if _, err := git(dst, "checkout", "-q", "--detach", "base"); err != nil {
		return nil, err
	}
	return refs, nil
}
