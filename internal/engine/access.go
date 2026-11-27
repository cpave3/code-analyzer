package engine

import (
	"fmt"
	"path/filepath"

	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/gitsrc"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
)

// Accessors used by scorers that work directly on snapshots (scoring/v2).
// They expose the same read-only pipeline the v1 analyzers use and do not
// change v1 behaviour.

// ResolveSnapshot resolves a revision and returns its commit and snapshot.
func (e *Evaluator) ResolveSnapshot(ref string) (string, *index.Snapshot, error) {
	sha, err := e.repo.ResolveCommit(ref)
	if err != nil {
		return "", nil, err
	}
	s, err := e.snapshot(sha)
	if err != nil {
		return "", nil, fmt.Errorf("index %s: %w", short(sha), err)
	}
	return sha, s, nil
}

// WorktreeSnapshot indexes the tracked files of the working tree.
func (e *Evaluator) WorktreeSnapshot() (*index.Snapshot, error) {
	return e.snapshot(gitsrc.WorktreeRev)
}

// ChangeSet computes the change context between two snapshots.
func (e *Evaluator) ChangeSet(base, target *index.Snapshot) (*change.Set, error) {
	cs, _, err := change.Compute(e.repo, base, target)
	return cs, err
}

// FirstParentSteps returns (parent, commit) pairs along the first-parent chain
// from base (exclusive) to target (inclusive). base must be an ancestor.
func (e *Evaluator) FirstParentSteps(baseSHA, targetSHA string) ([][2]string, error) {
	if !e.repo.IsAncestor(baseSHA, targetSHA) {
		return nil, fmt.Errorf("base %s is not an ancestor of target %s", short(baseSHA), short(targetSHA))
	}
	chain, err := e.repo.FirstParentChain(baseSHA, targetSHA)
	if err != nil {
		return nil, err
	}
	out := make([][2]string, 0, len(chain))
	for _, c := range chain {
		p, err := e.repo.FirstParent(c)
		if err != nil {
			return nil, err
		}
		out = append(out, [2]string{p, c})
	}
	return out, nil
}

// RepoName is the base name of the repository directory.
func (e *Evaluator) RepoName() string { return filepath.Base(e.repo.Path) }

// CommitTime returns a commit's committer time (RFC 3339).
func (e *Evaluator) CommitTime(sha string) string { return e.repo.CommitTime(sha) }

// Subject returns a commit's subject line (untrusted text).
func (e *Evaluator) Subject(sha string) string { return e.repo.Subject(sha) }

// Config returns the effective configuration.
func (e *Evaluator) Config() *config.Config { return e.opts.Config }

// Langs returns the language registry.
func (e *Evaluator) Langs() *lang.Registry { return e.opts.Langs }
