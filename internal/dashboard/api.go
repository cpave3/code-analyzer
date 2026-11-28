package dashboard

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/cameron/code-evaluator/internal/report"
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
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	page, _ := fs.Sub(web, "web")
	mux.Handle("GET /", http.FileServerFS(page))
	mux.HandleFunc("GET /api/meta", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, s.meta) })
	mux.HandleFunc("GET /api/metrics", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, s.windowMetrics()) })
	mux.HandleFunc("POST /api/window", s.handleWindow)
	mux.HandleFunc("GET /api/graph", s.handleGraph)
	return securityHeaders(mux)
}

// windowMetrics returns the window, progress and per-index metrics.
func (s *Server) windowMetrics() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	from := max(0, s.lo-1)
	ms := map[string]Metrics{}
	done, steps := 0, 0
	for i := from; i <= s.hi; i++ {
		m := s.metrics[i]
		ms[strconv.Itoa(i)] = m
		if i >= s.lo && m.Done {
			done++
		}
		if i >= s.lo && m.StepDone {
			steps++
		}
	}
	measured := 0
	for _, m := range s.metrics {
		if m.Done {
			measured++
		}
	}
	return map[string]any{"window": map[string]int{"from": s.lo, "to": s.hi}, "done": done, "stepsDone": steps,
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
