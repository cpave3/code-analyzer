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
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/cameron/code-evaluator/internal/clones"
	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/engine"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/provenance"
	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/scorev2"
	"github.com/cameron/code-evaluator/internal/viz"
)

//go:embed web
var web embed.FS

// Options select the history and what is measured.
type Options struct {
	Base     string // exclusive start; "" = the last Max first-parent commits
	Target   string
	Max      int
	Defects  bool
	Clones   bool
	NotesRef string
	Infer    bool // infer agents from Co-Authored-By lines
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
// against the previous commit in the series.
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
	ChangedLines int          `json:"changedLines"` // vs previous commit in the series
	Gate         *int         `json:"gate,omitempty"`
	Advisory     *int         `json:"advisory,omitempty"`
	IntroGate    *int         `json:"introducedGate,omitempty"` // vs previous commit
	IntroAdv     *int         `json:"introducedAdvisory,omitempty"`
	Resolved     *int         `json:"resolved,omitempty"`
	NearMisses   *int         `json:"nearMisses,omitempty"`
	Fragments    *int         `json:"fragments,omitempty"`
}

// Server holds the evaluator and the computed series.
type Server struct {
	opts  Options
	ev    *engine.Evaluator
	cfg   *config.Config
	langs *lang.Registry
	reg   *rules.Registry

	evMu    sync.Mutex // the evaluator is not safe for concurrent use
	commits []Commit
	meta    map[string]any

	mu      sync.Mutex
	metrics []Metrics
	graphs  map[string][]byte
	order   []string
}

// New resolves the history and its provenance; call Run to compute metrics.
func New(ev *engine.Evaluator, cfg *config.Config, langs *lang.Registry, reg *rules.Registry, opts Options) (*Server, error) {
	s := &Server{opts: opts, ev: ev, cfg: cfg, langs: langs, reg: reg, graphs: map[string][]byte{}}
	shas, err := s.history()
	if err != nil {
		return nil, err
	}
	if len(shas) == 0 {
		return nil, fmt.Errorf("no commits in range")
	}
	records, err := s.provenance(shas)
	if err != nil {
		return nil, err
	}
	for i, sha := range shas {
		info, err := ev.Info(sha)
		if err != nil {
			return nil, err
		}
		c := Commit{Index: i, SHA: sha, Subject: info.Subject, Author: info.Author, Time: info.Time, Merge: len(info.Parents) > 1}
		if len(info.Parents) > 0 {
			c.Parent = info.Parents[0]
		}
		if c.Provenance, err = s.stepRecords(c, records); err != nil {
			return nil, err
		}
		s.commits = append(s.commits, c)
	}
	s.metrics = make([]Metrics, len(s.commits))
	s.meta = s.describe()
	return s, nil
}

func (s *Server) history() ([]string, error) {
	tsha, err := s.ev.ResolveCommit(s.opts.Target)
	if err != nil {
		return nil, err
	}
	if s.opts.Base == "" {
		return s.ev.LastFirstParents(tsha, s.opts.Max)
	}
	bsha, err := s.ev.ResolveCommit(s.opts.Base)
	if err != nil {
		return nil, err
	}
	steps, err := s.ev.FirstParentSteps(bsha, tsha)
	if err != nil {
		return nil, err
	}
	shas := []string{bsha}
	for _, p := range steps {
		shas = append(shas, p[1])
	}
	if s.opts.Max > 0 && len(shas) > s.opts.Max {
		shas = shas[len(shas)-s.opts.Max:]
	}
	return shas, nil
}

// provenance parses messages and notes for every commit the series brings in.
func (s *Server) provenance(shas []string) (map[string]provenance.Record, error) {
	first, err := s.ev.Info(shas[0])
	if err != nil {
		return nil, err
	}
	base := ""
	if len(first.Parents) > 0 {
		base = first.Parents[0]
	}
	list, err := s.ev.CommitsBetween(base, shas[len(shas)-1], s.opts.NotesRef)
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

// stepRecords lists the provenance of the commits a step brings in.
func (s *Server) stepRecords(c Commit, records map[string]provenance.Record) ([]provenance.Record, error) {
	shas := []string{c.SHA}
	if c.Parent != "" {
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
	fields := []string{provenance.FieldAgent, provenance.FieldHarness}
	var extra []string
	for k := range keys {
		extra = append(extra, k)
	}
	sort.Strings(extra)
	notes := []string{
		"scoring/v2 debt and score are uncalibrated structural proxies; defects (rules/v1) and near-duplication (clones/v1) are findings, not scores",
		"provenance is self-reported text in commit messages and notes (or inferred from Co-Authored-By lines); comparisons between agents or workflows are observational, not causal",
		"each point is a first-parent commit; a merge step includes every commit it brings in",
	}
	return map[string]any{
		"repository": s.ev.RepoName(), "tool": report.ToolVersion, "scoring": scorev2.Version, "rules": rules.Version, "clones": clones.Version,
		"defects": s.opts.Defects, "nearDuplication": s.opts.Clones, "groupFields": append(fields, extra...), "notes": notes,
		"commits": s.commits,
	}
}

// Run computes every commit's metrics, oldest first.
func (s *Server) Run() {
	var prevSnap *index.Snapshot
	var prevRules *rules.Result
	for i := range s.commits {
		m, snap, res := s.measure(s.commits[i].SHA, prevSnap, prevRules)
		s.mu.Lock()
		s.metrics[i] = m
		s.mu.Unlock()
		if m.Error == "" {
			prevSnap, prevRules = snap, res
		}
	}
}

func (s *Server) measure(sha string, prev *index.Snapshot, prevRules *rules.Result) (m Metrics, snap *index.Snapshot, res *rules.Result) {
	s.evMu.Lock()
	defer s.evMu.Unlock()
	m.Done = true
	_, snap, err := s.ev.ResolveSnapshot(sha)
	if err != nil {
		m.Error = err.Error()
		return m, nil, nil
	}
	a := scorev2.Analyze(snap, s.cfg, scorev2.DefaultWeights())
	c := a.Codebase.Components
	m.Debt, m.DebtPerKLOC, m.Score = a.Codebase.Debt, a.Codebase.DebtPerKLOC, a.Codebase.Score
	m.Lines, m.Files, m.Callables, m.Hotspots, m.Decisions = c.SourceLines, c.SourceFiles, c.Callables, c.Hotspots, c.DecisionPoints
	m.Cloned, m.Violations, m.CycleUnits = c.ClonedCallables, c.RuleViolations, c.CycleUnits
	if prev != nil {
		if cs, err := s.ev.ChangeSet(prev, snap); err == nil {
			for _, f := range cs.Files {
				if f.Supported {
					m.ChangedLines += f.Added + f.Deleted
				}
			}
		}
	}
	if s.opts.Defects {
		res = rules.Run(s.reg, snap, s.cfg, s.langs, rules.Options{})
		gate, adv := 0, 0
		for _, f := range res.Findings {
			if f.Severity == rules.Gate {
				gate++
			} else {
				adv++
			}
		}
		m.Gate, m.Advisory = &gate, &adv
		if prevRules != nil {
			in, out, _ := rules.Compare(prevRules.Findings, res.Findings)
			ig, ia, r := 0, 0, len(out)
			for _, f := range in {
				if f.Severity == rules.Gate {
					ig++
				} else {
					ia++
				}
			}
			m.IntroGate, m.IntroAdv, m.Resolved = &ig, &ia, &r
		}
	}
	if s.opts.Clones {
		cr := clones.Analyze(snap, s.cfg, s.langs, nil)
		n, f := len(cr.NearMisses), len(cr.Fragments)
		m.NearMisses, m.Fragments = &n, &f
	}
	return m, snap, res
}

// graph builds (and caches) one commit's module graph.
func (s *Server) graph(sha string, maxCallables int) ([]byte, error) {
	key := sha + "/" + strconv.Itoa(maxCallables)
	s.mu.Lock()
	if b, ok := s.graphs[key]; ok {
		s.mu.Unlock()
		return b, nil
	}
	s.mu.Unlock()
	s.evMu.Lock()
	_, snap, err := s.ev.ResolveSnapshot(sha)
	if err != nil {
		s.evMu.Unlock()
		return nil, err
	}
	d := viz.Build(s.ev.RepoName(), report.Revision{Ref: sha, Commit: sha, Tree: snap.Tree, CommitTime: s.ev.CommitTime(sha)}, snap, s.cfg, scorev2.Analyze(snap, s.cfg, scorev2.DefaultWeights()))
	s.evMu.Unlock()
	d.Trim(maxCallables)
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.graphs[key] = b
	s.order = append(s.order, key)
	if len(s.order) > 64 {
		delete(s.graphs, s.order[0])
		s.order = s.order[1:]
	}
	s.mu.Unlock()
	return b, nil
}

// Handler serves the page and the JSON API.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	page, _ := fs.Sub(web, "web")
	mux.Handle("GET /", http.FileServerFS(page))
	mux.HandleFunc("GET /api/meta", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, s.meta) })
	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		done := 0
		for _, m := range s.metrics {
			if m.Done {
				done++
			}
		}
		body := map[string]any{"done": done, "total": len(s.metrics), "metrics": append([]Metrics(nil), s.metrics...)}
		s.mu.Unlock()
		writeJSON(w, body)
	})
	mux.HandleFunc("GET /api/graph", func(w http.ResponseWriter, r *http.Request) {
		sha := r.URL.Query().Get("commit")
		known := false
		for _, c := range s.commits {
			known = known || c.SHA == sha
		}
		if !known {
			http.Error(w, "unknown commit", http.StatusNotFound)
			return
		}
		max, _ := strconv.Atoi(r.URL.Query().Get("callables"))
		if max <= 0 {
			max = 120
		}
		b, err := s.graph(sha, max)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	return securityHeaders(mux)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// securityHeaders keeps the page self-contained and same-origin.
func securityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; connect-src 'self'; img-src 'self' data:")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		h.ServeHTTP(w, r)
	})
}
