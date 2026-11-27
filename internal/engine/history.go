package engine

import (
	"fmt"
	"time"

	"github.com/cameron/code-evaluator/internal/report"
)

// History evaluates the net change base..target and each first-parent step
// within it (parent^1 → commit), so merges are attributed once without
// double counting their branch commits.
func (e *Evaluator) History(baseRef, targetRef string, progress func(string)) (*report.History, error) {
	start := time.Now()
	parsed0, hits0, miss0 := e.builder.Parsed, e.cache.Hits, e.cache.Misses
	baseSHA, err := e.repo.ResolveCommit(baseRef)
	if err != nil {
		return nil, err
	}
	targetSHA, err := e.repo.ResolveCommit(targetRef)
	if err != nil {
		return nil, err
	}
	if !e.repo.IsAncestor(baseSHA, targetSHA) {
		return nil, fmt.Errorf("base %s is not an ancestor of target %s", baseSHA, targetSHA)
	}
	chain, err := e.repo.FirstParentChain(baseSHA, targetSHA)
	if err != nil {
		return nil, err
	}
	h := &report.History{
		Schema: report.HistorySchema,
		Tool:   report.Tool{Name: "codeeval", Version: report.ToolVersion},
		Base:   report.Revision{Ref: baseRef, Commit: baseSHA, CommitTime: e.repo.CommitTime(baseSHA)},
		Target: report.Revision{Ref: targetRef, Commit: targetSHA, CommitTime: e.repo.CommitTime(targetSHA)},
		Steps:  []report.HistoryStep{},
		Notes: []string{
			"each step compares a first-parent commit with its first parent; for merge commits this is the whole merged branch, counted once",
			"step results are observed contributions to metrics, not proof that a workflow or person caused a quality change",
			"step reports do not sum to the net report: repository-wide measurements (duplication, graph) interact across steps",
		},
	}
	if progress != nil {
		progress(fmt.Sprintf("net %s..%s", short(baseSHA), short(targetSHA)))
	}
	net, err := e.evaluateResolved(baseRef, baseSHA, targetRef, targetSHA)
	if err != nil {
		return nil, err
	}
	h.Net = net
	h.Base.Tree, h.Target.Tree = net.Base.Tree, net.Target.Tree
	for i, c := range chain {
		parent, err := e.repo.FirstParent(c)
		if err != nil {
			return nil, err
		}
		if i == 0 && parent != baseSHA {
			h.Notes = append(h.Notes, fmt.Sprintf("first step's parent %s is not the base (base is not on the first-parent chain)", short(parent)))
		}
		if progress != nil {
			progress(fmt.Sprintf("step %d/%d %s", i+1, len(chain), short(c)))
		}
		r, err := e.evaluateResolved(parent, parent, c, c)
		if err != nil {
			return nil, fmt.Errorf("step %s: %w", c, err)
		}
		h.Steps = append(h.Steps, report.HistoryStep{Index: i + 1, Commit: c, Parent: parent, Subject: e.repo.Subject(c), Time: e.repo.CommitTime(c), Report: r})
	}
	if e.opts.Timing {
		h.Timing = &report.Timing{StartedAt: start.Format(time.RFC3339), ElapsedMillis: time.Since(start).Milliseconds(),
			FilesParsed: e.builder.Parsed - parsed0, CacheHits: e.cache.Hits - hits0, CacheMisses: e.cache.Misses - miss0}
	}
	return h, nil
}

func short(s string) string {
	if len(s) > 12 {
		return s[:12]
	}
	return s
}
