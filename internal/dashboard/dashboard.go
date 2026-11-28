// Package dashboard serves a local, read-only web dashboard of quality over
// a repository's first-parent history: per-commit scoring/v2 debt, score,
// hotspots, defect findings and near-duplication, the module graph of each
// commit, and commit provenance (agent, harness, workflow ...).
//
// Metrics are computed in the background, oldest commit first, and the page
// polls for them; graphs are built on demand. Everything is derived from
// Git objects through the engine; nothing in the repository is written.
package dashboard

import (
	"embed"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cameron/code-evaluator/internal/clones"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/gitsrc"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/memo"
	"github.com/cameron/code-evaluator/internal/provenance"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// web/dist is the built page (Vite + React, in web/); it is committed, so
// building the binary needs no Node.
//
//go:embed web/dist
var web embed.FS

// Options select the history and the initially measured window.
type Options struct {
	Base     string // first commit listed (the history after it, plus itself); "" = back to the root
	Target   string
	Max      int // most commits listed (0 = the whole first-parent history)
	Last     int // the window measured at start: the newest Last commits (default 50)
	Since    string
	Until    string // or a date range (RFC 3339 date or time); overrides Last
	Defects  bool
	Clones   bool
	NotesRef string
	Infer    bool // infer models from Co-Authored-By lines
}

// Commit is one first-parent commit with its metadata.
type Commit struct {
	Index      int                 `json:"index"`
	SHA        string              `json:"sha"`
	Parent     string              `json:"parent"`
	Subject    string              `json:"subject"` // untrusted text
	Author     string              `json:"author"`
	Time       string              `json:"time"`
	Merge      bool                `json:"merge"`
	Provenance []provenance.Record `json:"provenance"` // commits the step brings in (a merge: the merged branch too)
}

// Metrics are one commit's absolute measurements, plus step changes
// against the previous commit in the history.
type Metrics struct {
	Done         bool         `json:"done"`
	Error        string       `json:"error,omitempty"`
	Debt         scorev2.Debt `json:"debt"`
	DebtPerKLOC  *float64     `json:"debtPerKLOC"`
	Score        *float64     `json:"score"`
	Lines        int          `json:"lines"`
	Files        int          `json:"files"`
	Callables    int          `json:"callables"`
	Hotspots     int          `json:"hotspots"`
	Decisions    int          `json:"decisionPoints"`
	Cloned       int          `json:"clonedCallables"`
	Violations   int          `json:"ruleViolations"`
	CycleUnits   int          `json:"cycleUnits"`
	ChangedLines int          `json:"changedLines"` // vs previous commit
	// StepDone is set once the changes against the previous commit
	// (changed lines, defects introduced/resolved) are known.
	StepDone bool `json:"stepDone"`
	// Detail is set once defects and near-duplication are measured (the
	// slower phase); DetailStep once defects introduced/resolved are known.
	Detail     bool `json:"detail"`
	DetailStep bool `json:"detailStep"`
	Gate       *int `json:"gate,omitempty"`
	Advisory   *int `json:"advisory,omitempty"`
	IntroGate  *int `json:"introducedGate,omitempty"`
	IntroAdv   *int `json:"introducedAdvisory,omitempty"`
	Resolved   *int `json:"resolved,omitempty"`
	NearMisses *int `json:"nearMisses,omitempty"`
	Fragments  *int `json:"fragments,omitempty"`
}

// Server holds the evaluator, the listed history and what has been measured.
type Server struct {
	opts  Options
	ev    *engine.Evaluator
	cfg   *config.Config
	langs *lang.Registry
	reg   *rules.Registry

	evMu    sync.Mutex // the evaluator is not safe for concurrent use
	commits []Commit
	meta    map[string]any

	mu       sync.Mutex
	work     *sync.Cond
	metrics  []Metrics
	findings map[int][]rules.Finding // measured commits' defect findings, for step comparison
	cache    *memo.Cache             // per-file results by content, shared across commits
	// defectRuns keeps whole-commit rule results checked on demand.
	defectRuns  map[string]*rules.Result
	defectOrder []string
	lo, hi      int // measured window (inclusive indexes)
	graphs      map[string][]byte
	order       []string
}

// New lists the history and its provenance (cheap: a few git calls); call
// Run to measure the window.
func New(ev *engine.Evaluator, cfg *config.Config, langs *lang.Registry, reg *rules.Registry, opts Options) (*Server, error) {
	s := &Server{opts: opts, ev: ev, cfg: cfg, langs: langs, reg: reg, graphs: map[string][]byte{}, findings: map[int][]rules.Finding{}, cache: memo.New(), defectRuns: map[string]*rules.Result{}}
	s.work = sync.NewCond(&s.mu)
	log, err := s.history()
	if err != nil {
		return nil, err
	}
	if len(log) == 0 {
		return nil, fmt.Errorf("no commits in range")
	}
	records, err := s.provenance(log)
	if err != nil {
		return nil, err
	}
	for i, e := range log {
		c := Commit{Index: i, SHA: e.SHA, Subject: e.Subject, Author: e.Author, Time: e.Time, Merge: len(e.Parents) > 1}
		if len(e.Parents) > 0 {
			c.Parent = e.Parents[0]
		}
		if c.Provenance, err = s.stepRecords(c, records); err != nil {
			return nil, err
		}
		s.commits = append(s.commits, c)
	}
	s.metrics = make([]Metrics, len(s.commits))
	if s.lo, s.hi, err = s.initialWindow(); err != nil {
		return nil, err
	}
	s.meta = s.describe()
	return s, nil
}

func (s *Server) history() ([]gitsrc.LogEntry, error) {
	tsha, err := s.ev.ResolveCommit(s.opts.Target)
	if err != nil {
		return nil, err
	}
	base := ""
	if s.opts.Base != "" {
		if base, err = s.ev.ResolveCommit(s.opts.Base); err != nil {
			return nil, err
		}
		if _, err := s.ev.FirstParentSteps(base, tsha); err != nil {
			return nil, err
		}
	}
	log, err := s.ev.FirstParentLog(base, tsha, s.opts.Max)
	if err != nil || base == "" || (s.opts.Max > 0 && len(log) >= s.opts.Max) {
		return log, err
	}
	// Show the base itself first, so the first listed step has a parent.
	b, err := s.ev.FirstParentLog("", base, 1)
	if err != nil {
		return nil, err
	}
	return append(b, log...), nil
}

// provenance parses messages and notes for every commit the history brings
// in, in one git call.
func (s *Server) provenance(log []gitsrc.LogEntry) (map[string]provenance.Record, error) {
	base := ""
	if p := log[0].Parents; len(p) > 0 {
		base = p[0]
	}
	list, err := s.ev.CommitsBetween(base, log[len(log)-1].SHA, s.opts.NotesRef)
	if err != nil {
		return nil, err
	}
	out := map[string]provenance.Record{}
	for _, m := range list {
		rec := provenance.Parse(m.SHA, m.Message, m.Note, s.opts.Infer)
		rec.Merge = len(m.Parents) > 1
		out[m.SHA] = rec
	}
	return out, nil
}

// stepRecords lists the provenance of the commits a step brings in: the
// commit itself, plus the merged branch for a merge (one git call per merge).
func (s *Server) stepRecords(c Commit, records map[string]provenance.Record) ([]provenance.Record, error) {
	shas := []string{c.SHA}
	if c.Merge {
		var err error
		if shas, err = s.ev.RevList(c.Parent, c.SHA); err != nil {
			return nil, err
		}
	}
	var out []provenance.Record
	for _, sha := range shas {
		rec, ok := records[sha]
		if !ok {
			rec = provenance.Record{Commit: sha, Source: provenance.SourceNone}
		}
		out = append(out, rec)
	}
	return out, nil
}

// initialWindow picks the measured window from --since/--until or --last.
func (s *Server) initialWindow() (int, int, error) {
	n := len(s.commits)
	if s.opts.Since != "" || s.opts.Until != "" {
		since, until, err := parseRange(s.opts.Since, s.opts.Until)
		if err != nil {
			return 0, 0, err
		}
		lo, hi := -1, -1
		for i, c := range s.commits {
			t, err := time.Parse(time.RFC3339, c.Time)
			if err != nil || (!since.IsZero() && t.Before(since)) || (!until.IsZero() && t.After(until)) {
				continue
			}
			if lo < 0 {
				lo = i
			}
			hi = i
		}
		if lo < 0 {
			return 0, 0, fmt.Errorf("no commits between %q and %q", s.opts.Since, s.opts.Until)
		}
		return lo, hi, nil
	}
	last := s.opts.Last
	if last <= 0 || last > n {
		last = n
	}
	return n - last, n - 1, nil
}

func parseRange(since, until string) (time.Time, time.Time, error) {
	parse := func(v string, end bool) (time.Time, error) {
		if v == "" {
			return time.Time{}, nil
		}
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t, nil
		}
		t, err := time.ParseInLocation("2006-01-02", v, time.Local)
		if err != nil {
			return time.Time{}, fmt.Errorf("date %q: want YYYY-MM-DD or RFC 3339", v)
		}
		if end {
			t = t.Add(24*time.Hour - time.Nanosecond)
		}
		return t, nil
	}
	a, err := parse(since, false)
	if err != nil {
		return a, a, err
	}
	b, err := parse(until, true)
	return a, b, err
}

// describe lists the fields the page can group by and the run's settings.
func (s *Server) describe() map[string]any {
	keys := map[string]bool{}
	for _, c := range s.commits {
		for _, r := range c.Provenance {
			for k := range r.Meta {
				keys[k] = true
			}
		}
	}
	fields := []string{provenance.FieldModel, provenance.FieldHarness}
	var extra []string
	for k := range keys {
		extra = append(extra, k)
	}
	sort.Strings(extra)
	notes := []string{
		"scoring/v2 debt and score are uncalibrated structural proxies; defects (rules/v1) and near-duplication (clones/v1) are findings, not scores",
		"provenance is self-reported text in commit messages and notes (or inferred from Co-Authored-By lines); comparisons between models or workflows are observational, not causal",
		"each point is a first-parent commit; a merge step includes every commit it brings in",
		"only the selected window is measured (each commit once, then cached); the history strip lists every commit",
	}
	return map[string]any{
		"repository": s.ev.RepoName(), "tool": report.ToolVersion, "scoring": scorev2.Version, "rules": rules.Version, "clones": clones.Version,
		"defects": s.opts.Defects, "nearDuplication": s.opts.Clones, "groupFields": append(fields, extra...), "notes": notes,
		"thresholds": map[string]int{"cyclomatic": s.cfg.Complexity.HotspotCyclomatic, "nesting": s.cfg.Complexity.HotspotNesting, "minCloneTokens": s.cfg.Duplication.MinTokens},
		"weights":    scorev2.DefaultWeights(),
		"commits":    s.commits,
	}
}

// Commits is the number of commits listed.
func (s *Server) Commits() int { return len(s.commits) }

// Window returns the measured window (inclusive indexes).
func (s *Server) Window() (int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lo, s.hi
}
