// Package change computes the change context between two snapshots: changed
// files with line ranges and counts, and symbols mapped to added, deleted,
// modified, renamed, moved or formatting-only status.
package change

import (
	"bytes"
	"sort"

	"github.com/cameron/code-evaluator/internal/gitsrc"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/model"
	"github.com/cameron/code-evaluator/internal/textdiff"
)

// File statuses (Git name-status letters).
const (
	Added    = "added"
	Deleted  = "deleted"
	Modified = "modified"
	Renamed  = "renamed"
	Copied   = "copied"
	TypeChg  = "type_changed"
)

// FileChange is one changed file.
type FileChange struct {
	Status      string          `json:"status"`
	OldPath     string          `json:"oldPath,omitempty"`
	Path        string          `json:"path"`
	Similarity  int             `json:"similarity,omitempty"`
	Language    string          `json:"language,omitempty"`
	Supported   bool            `json:"supported"`
	Class       model.FileClass `json:"class"`
	Added       int             `json:"linesAdded"`
	Deleted     int             `json:"linesDeleted"`
	Binary      bool            `json:"binary,omitempty"`
	Approximate bool            `json:"approximateDiff,omitempty"`

	OldChanged []textdiff.Range `json:"-"`
	NewChanged []textdiff.Range `json:"-"`
}

// Symbol statuses.
const (
	SymAdded      = "added"
	SymDeleted    = "deleted"
	SymModified   = "modified"
	SymRenamed    = "renamed"         // same body, different name (or package)
	SymMoved      = "moved"           // same identity and body, different file
	SymFormatting = "formatting_only" // touched, but only layout/comments changed
)

// SymbolChange maps a symbol across the change.
type SymbolChange struct {
	Status string        `json:"status"`
	Base   *model.Symbol `json:"-"`
	Target *model.Symbol `json:"-"`
	// SignatureChanged is set for modified symbols whose signature differs.
	SignatureChanged bool `json:"signatureChanged,omitempty"`
}

// Primary returns the target symbol, or the base one for deletions.
func (c *SymbolChange) Primary() *model.Symbol {
	if c.Target != nil {
		return c.Target
	}
	return c.Base
}

// Set is the full change context.
type Set struct {
	Files   []FileChange
	Symbols []SymbolChange
}

// FileByPath finds a change by new path.
func (s *Set) FileByPath(p string) *FileChange {
	for i := range s.Files {
		if s.Files[i].Path == p {
			return &s.Files[i]
		}
	}
	return nil
}

// Compute builds the change set between two snapshots.
func Compute(repo *gitsrc.Repo, base, target *index.Snapshot) (*Set, []model.Diagnostic, error) {
	ns, err := repo.DiffNameStatus(base.Commit, target.Commit)
	if err != nil {
		return nil, nil, err
	}
	var diags []model.Diagnostic
	set := &Set{}
	for _, n := range ns {
		fc := FileChange{OldPath: n.OldPath, Path: n.NewPath, Similarity: n.Similarity}
		switch n.Status {
		case "A":
			fc.Status, fc.OldPath = Added, ""
		case "D":
			fc.Status = Deleted
		case "R":
			fc.Status = Renamed
		case "C":
			fc.Status = Copied
		case "T":
			fc.Status = TypeChg
		default:
			fc.Status = Modified
		}
		if fc.OldPath == fc.Path && fc.Status != Deleted {
			fc.OldPath = ""
		}
		var bf, tf *index.File
		if fc.Status != Added {
			bf = base.Files[n.OldPath]
		}
		if fc.Status != Deleted {
			tf = target.Files[n.NewPath]
		}
		ref := tf
		if ref == nil {
			ref = bf
		}
		if ref == nil {
			continue // symlink/submodule
		}
		fc.Language, fc.Class = ref.Language, ref.Class
		fc.Supported = (bf == nil || bf.Supported()) && (tf == nil || tf.Supported())
		var oldB, newB []byte
		if bf != nil {
			if oldB, err = base.Content(bf.Path); err != nil {
				diags = append(diags, model.Diagnostic{Severity: model.SevError, Code: "blob_unreadable", Message: err.Error(), File: bf.Path, Side: "base"})
			}
		}
		if tf != nil {
			if newB, err = target.Content(tf.Path); err != nil {
				diags = append(diags, model.Diagnostic{Severity: model.SevError, Code: "blob_unreadable", Message: err.Error(), File: tf.Path, Side: "target"})
			}
		}
		if bytes.IndexByte(oldB, 0) >= 0 || bytes.IndexByte(newB, 0) >= 0 {
			fc.Binary = true
		} else {
			d := textdiff.Lines(oldB, newB)
			fc.Added, fc.Deleted, fc.OldChanged, fc.NewChanged, fc.Approximate = d.Added, d.Deleted, d.OldChanged, d.NewChanged, d.Approximate
		}
		set.Files = append(set.Files, fc)
	}
	sort.Slice(set.Files, func(i, j int) bool { return set.Files[i].Path < set.Files[j].Path })
	set.Symbols = mapSymbols(set, base, target)
	return set, diags, nil
}

func mapSymbols(set *Set, base, target *index.Snapshot) []SymbolChange {
	oldRanges := map[string][]textdiff.Range{}
	newRanges := map[string][]textdiff.Range{}
	var baseSyms, targetSyms []*model.Symbol
	for _, fc := range set.Files {
		if fc.Status != Added {
			old := fc.Path
			if fc.OldPath != "" {
				old = fc.OldPath
			}
			oldRanges[old] = fc.OldChanged
			baseSyms = append(baseSyms, base.SymbolsIn(old)...)
		}
		if fc.Status != Deleted {
			newRanges[fc.Path] = fc.NewChanged
			targetSyms = append(targetSyms, target.SymbolsIn(fc.Path)...)
		}
	}
	byKey := map[string]*model.Symbol{}
	for _, s := range baseSyms {
		byKey[s.Key] = s
	}
	matched := map[*model.Symbol]bool{}
	var out []SymbolChange
	var addedC []*model.Symbol
	for _, t := range targetSyms {
		b, ok := byKey[t.Key]
		if !ok {
			addedC = append(addedC, t)
			continue
		}
		matched[b] = true
		touched := textdiff.Overlaps(newRanges[t.File], t.StartLine, t.EndLine) ||
			textdiff.Overlaps(oldRanges[b.File], b.StartLine, b.EndLine)
		switch {
		case b.ExactFP != t.ExactFP || b.Signature != t.Signature:
			out = append(out, SymbolChange{Status: SymModified, Base: b, Target: t, SignatureChanged: b.Signature != t.Signature})
		case b.File != t.File:
			out = append(out, SymbolChange{Status: SymMoved, Base: b, Target: t})
		case touched:
			out = append(out, SymbolChange{Status: SymFormatting, Base: b, Target: t})
		}
	}
	var deletedC []*model.Symbol
	for _, b := range baseSyms {
		if !matched[b] {
			deletedC = append(deletedC, b)
		}
	}
	// Pair renames: same kind and identical exact fingerprint (name excluded).
	usedDel := map[*model.Symbol]bool{}
	for _, a := range addedC {
		var pair *model.Symbol
		for _, d := range deletedC {
			if !usedDel[d] && d.Kind == a.Kind && d.ExactFP == a.ExactFP && a.ExactFP != "" {
				pair = d
				break
			}
		}
		if pair != nil {
			usedDel[pair] = true
			out = append(out, SymbolChange{Status: SymRenamed, Base: pair, Target: a})
			continue
		}
		out = append(out, SymbolChange{Status: SymAdded, Target: a})
	}
	for _, d := range deletedC {
		if !usedDel[d] {
			out = append(out, SymbolChange{Status: SymDeleted, Base: d})
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i].Primary(), out[j].Primary()
		if a.File != b.File {
			return a.File < b.File
		}
		if a.StartLine != b.StartLine {
			return a.StartLine < b.StartLine
		}
		return out[i].Status < out[j].Status
	})
	return out
}
