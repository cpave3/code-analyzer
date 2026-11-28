package dashboard

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/cameron/code-evaluator/internal/report"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/scorev2"
	"github.com/cameron/code-evaluator/internal/viz"
)

// Handler serves the page and the JSON API:
//
//	GET  /api/meta     every commit (metadata, provenance) and settings
//	GET  /api/metrics  the window and the metrics of its commits (+ the
//	                   parent of the first)
//	POST /api/window   {"from": i, "to": j} moves the measured window
//	GET  /api/graph    ?commit=SHA: one commit's module graph
//	GET  /api/explain  ?commit=SHA&path=DIR: v2 debt factors for a path
//	                   (codeeval explain --path), all factors
//	GET  /api/defects  ?commit=SHA&path=DIR: defect findings for a path
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	page, _ := fs.Sub(web, "web")
	mux.Handle("GET /", http.FileServerFS(page))
	mux.HandleFunc("GET /api/meta", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, s.meta) })
	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, s.windowMetrics()) })
	mux.HandleFunc("POST /api/window", s.handleWindow)
	mux.HandleFunc("GET /api/graph", s.handleGraph)
	mux.HandleFunc("GET /api/explain", s.handleExplain)
	mux.HandleFunc("GET /api/defects", s.handleDefects)
	return securityHeaders(mux)
}

// windowMetrics returns the window, progress and per-index metrics.
func (s *Server) windowMetrics() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	from := max(0, s.lo-1)
	ms := map[string]Metrics{}
	done, steps, detail, detailSteps := 0, 0, 0, 0
	for i := from; i <= s.hi; i++ {
		m := s.metrics[i]
		ms[strconv.Itoa(i)] = m
		if i < s.lo {
			continue
		}
		for _, x := range []struct {
			ok bool
			n  *int
		}{{m.Done, &done}, {m.StepDone, &steps}, {m.Detail, &detail}, {m.DetailStep, &detailSteps}} {
			if x.ok {
				*x.n++
			}
		}
	}
	measured := 0
	for _, m := range s.metrics {
		if m.Done {
			measured++
		}
	}
	slow := s.opts.Defects || s.opts.Clones
	if !slow {
		detail, detailSteps = s.hi-s.lo+1, s.hi-s.lo+1
	}
	entries, hits, misses := s.cache.Stats()
	return map[string]any{"window": map[string]int{"from": s.lo, "to": s.hi}, "done": done, "stepsDone": steps,
		"detailDone": detail, "detailStepsDone": detailSteps, "cache": map[string]int{"entries": entries, "hits": hits, "misses": misses},
		"total": s.hi - s.lo + 1, "measuredOverall": measured, "metrics": ms}
}

// handleWindow moves the window. Only same-origin JSON requests are
// accepted, so another site cannot make the local server do work.
func (s *Server) handleWindow(w http.ResponseWriter, r *http.Request) {
	if o := r.Header.Get("Origin"); o != "" && o != "http://"+r.Host {
		http.Error(w, "cross-origin request refused", http.StatusForbidden)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "want application/json", http.StatusUnsupportedMediaType)
		return
	}
	var req struct{ From, To int }
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<10)).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lo, hi := s.SetWindow(req.From, req.To)
	writeJSON(w, map[string]int{"from": lo, "to": hi})
}

func (s *Server) handleGraph(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) known(sha string) bool {
	for _, c := range s.commits {
		if c.SHA == sha {
			return true
		}
	}
	return false
}

// cached returns a cached JSON body for key, or builds and stores it.
func (s *Server) cached(key string, build func() (any, error)) ([]byte, error) {
	s.mu.Lock()
	if b, ok := s.graphs[key]; ok {
		s.mu.Unlock()
		return b, nil
	}
	s.mu.Unlock()
	v, err := build()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.graphs[key] = b
	s.order = append(s.order, key)
	if len(s.order) > 96 {
		delete(s.graphs, s.order[0])
		s.order = s.order[1:]
	}
	s.mu.Unlock()
	return b, nil
}

func (s *Server) serveJSON(w http.ResponseWriter, b []byte, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func pathScope(p string) scorev2.Scope {
	p = strings.Trim(strings.TrimPrefix(p, "./"), "/")
	if p == "" || p == "." {
		return scorev2.Scope{Paths: []string{}}
	}
	return scorev2.Scope{Paths: []string{p}}
}

// handleExplain returns the scoped explain report (all factors), as
// `codeeval explain --rev SHA --path DIR --top 0` would.
func (s *Server) handleExplain(w http.ResponseWriter, r *http.Request) {
	sha, path := r.URL.Query().Get("commit"), r.URL.Query().Get("path")
	if !s.known(sha) {
		http.Error(w, "unknown commit", http.StatusNotFound)
		return
	}
	b, err := s.cached("explain|"+sha+"|"+path, func() (any, error) {
		s.evMu.Lock()
		defer s.evMu.Unlock()
		_, snap, err := s.ev.ResolveSnapshot(sha)
		if err != nil {
			return nil, err
		}
		rev := report.Revision{Ref: sha, Commit: sha, Tree: snap.Tree, CommitTime: s.ev.CommitTime(sha)}
		a := scorev2.Analyze(snap, s.cfg, scorev2.DefaultWeights())
		return scorev2.NewExplainReport(s.ev.RepoName(), s.cfg.Hash(), rev, snap, a, pathScope(path), 0), nil
	})
	s.serveJSON(w, b, err)
}

// handleDefects returns the defect findings within a path. The rules run
// once per commit (whole repository: cross-file rules need it) and the
// result is kept, so every folder of that commit filters it instantly.
// With cached=1 it answers 204 when the commit has not been checked yet,
// so the page can offer the check instead of starting it.
func (s *Server) handleDefects(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sha, path := q.Get("commit"), q.Get("path")
	if !s.known(sha) {
		http.Error(w, "unknown commit", http.StatusNotFound)
		return
	}
	s.mu.Lock()
	res := s.defectRuns[sha]
	s.mu.Unlock()
	if res == nil && q.Get("cached") == "1" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if res == nil {
		s.evMu.Lock()
		_, snap, err := s.ev.ResolveSnapshot(sha)
		s.evMu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res = rules.Run(s.reg, snap, s.cfg, s.langs, rules.Options{Cache: s.cache})
		s.mu.Lock()
		s.defectRuns[sha] = res
		s.defectOrder = append(s.defectOrder, sha)
		if len(s.defectOrder) > 16 {
			delete(s.defectRuns, s.defectOrder[0])
			s.defectOrder = s.defectOrder[1:]
		}
		s.mu.Unlock()
	}
	scope := pathScope(path)
	scoped := &rules.Result{Coverage: res.Coverage, Files: res.Files, Diagnostics: res.Diagnostics}
	for _, f := range res.Findings {
		if scope.Match(f.File) {
			scoped.Findings = append(scoped.Findings, f)
		}
	}
	for _, f := range res.Suppressed {
		if scope.Match(f.File) {
			scoped.Suppressed = append(scoped.Suppressed, f)
		}
	}
	rep := rules.NewReport(s.reg, nil, scoped, rules.FailNone)
	rep.Repository, rep.ConfigHash, rep.Scope = s.ev.RepoName(), s.cfg.Hash(), scope.Paths
	rep.Revision = &report.Revision{Ref: sha, Commit: sha}
	writeJSON(w, rep)
}
