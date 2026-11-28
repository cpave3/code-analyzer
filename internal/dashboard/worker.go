package dashboard

import (
	"github.com/cameron/code-evaluator/internal/clones"
	"github.com/cameron/code-evaluator/internal/index"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/scorev2"
)

// SetWindow changes the measured window (inclusive indexes, clamped) and
// wakes the worker. Commits already measured stay cached.
func (s *Server) SetWindow(lo, hi int) (int, int) {
	n := len(s.commits)
	if lo > hi {
		lo, hi = hi, lo
	}
	lo, hi = max(0, min(lo, n-1)), max(0, min(hi, n-1))
	s.mu.Lock()
	s.lo, s.hi = lo, hi
	s.mu.Unlock()
	s.work.Broadcast()
	return lo, hi
}

// Job kinds, in priority order. The fast phase (v2 metrics and changed
// lines, ~0.4 s per commit on a 3,000-file repository) finishes for the
// whole window before the detail phase (defect rules and near-duplication,
// seconds per commit, cached by file content) starts, so charts and graphs
// appear quickly.
const (
	jobMetrics = iota
	jobStep
	jobDetail
	jobDetailStep
)

type job struct{ i, kind int }

// next picks the next job in the window, newest first. The window includes
// the parent of its first commit, so the first step has a delta. Called
// with s.mu held.
func (s *Server) next() (job, bool) {
	lo := max(0, s.lo-1)
	for i := s.hi; i >= lo; i-- {
		if !s.metrics[i].Done {
			return job{i, jobMetrics}, true
		}
	}
	ready := func(i int) bool { m := s.metrics[i]; return m.Done && m.Error == "" }
	for i := s.hi; i >= s.lo; i-- {
		if m := s.metrics[i]; !m.StepDone && m.Error == "" {
			if i == 0 {
				s.metrics[i].StepDone = true // no parent in the history
				continue
			}
			if ready(i - 1) {
				return job{i, jobStep}, true
			}
		}
	}
	if !s.opts.Defects && !s.opts.Clones {
		return job{}, false
	}
	for i := s.hi; i >= lo; i-- {
		if m := s.metrics[i]; !m.Detail && m.Error == "" {
			return job{i, jobDetail}, true
		}
	}
	for i := s.hi; i >= s.lo; i-- {
		if m := s.metrics[i]; m.Detail && !m.DetailStep {
			if i == 0 {
				s.metrics[i].DetailStep = true
				continue
			}
			if s.metrics[i-1].Detail {
				return job{i, jobDetailStep}, true
			}
		}
	}
	return job{}, false
}

// Run works through the window's jobs, waiting for new windows. It never
// returns.
func (s *Server) Run() {
	for {
		s.mu.Lock()
		j, ok := s.next()
		for !ok {
			s.work.Wait()
			j, ok = s.next()
		}
		s.mu.Unlock()
		switch j.kind {
		case jobMetrics:
			m := s.measure(s.commits[j.i].SHA)
			s.mu.Lock()
			s.metrics[j.i] = m
			s.mu.Unlock()
		case jobStep:
			s.fillStep(j.i)
		case jobDetail:
			s.detail(j.i)
		case jobDetailStep:
			s.fillDetailStep(j.i)
		}
	}
}

// measure computes one commit's v2 metrics (fast: parse results come from
// the index cache).
func (s *Server) measure(sha string) (m Metrics) {
	s.evMu.Lock()
	defer s.evMu.Unlock()
	m.Done = true
	_, snap, err := s.ev.ResolveSnapshot(sha)
	if err != nil {
		m.Error = err.Error()
		return m
	}
	a := scorev2.Analyze(snap, s.cfg, scorev2.DefaultWeights())
	c := a.Codebase.Components
	m.Debt, m.DebtPerKLOC, m.Score = a.Codebase.Debt, a.Codebase.DebtPerKLOC, a.Codebase.Score
	m.Lines, m.Files, m.Callables, m.Hotspots, m.Decisions = c.SourceLines, c.SourceFiles, c.Callables, c.Hotspots, c.DecisionPoints
	m.Cloned, m.Violations, m.CycleUnits = c.ClonedCallables, c.RuleViolations, c.CycleUnits
	return m
}

// fillStep records commit i's changed measured-source lines against i-1.
func (s *Server) fillStep(i int) {
	s.evMu.Lock()
	changed := 0
	var older, newer *index.Snapshot
	_, older, err1 := s.ev.ResolveSnapshot(s.commits[i-1].SHA)
	_, newer, err2 := s.ev.ResolveSnapshot(s.commits[i].SHA)
	if err1 == nil && err2 == nil {
		if cs, err := s.ev.ChangeSet(older, newer); err == nil {
			for _, f := range cs.Files {
				if f.Supported {
					changed += f.Added + f.Deleted
				}
			}
		}
	}
	s.evMu.Unlock()
	s.mu.Lock()
	s.metrics[i].ChangedLines, s.metrics[i].StepDone = changed, true
	s.mu.Unlock()
}

// detail runs the defect rules and near-duplication on commit i. Only
// resolving the snapshot holds the evaluator lock, so graph requests are
// not held up; per-file work is cached by content across commits.
func (s *Server) detail(i int) {
	s.evMu.Lock()
	_, snap, err := s.ev.ResolveSnapshot(s.commits[i].SHA)
	s.evMu.Unlock()
	var gate, adv, near, frag *int
	var findings []rules.Finding
	if err == nil && s.opts.Defects {
		res := rules.Run(s.reg, snap, s.cfg, s.langs, rules.Options{Cache: s.cache})
		g, a := 0, 0
		findings = make([]rules.Finding, 0, len(res.Findings))
		for _, f := range res.Findings {
			if f.Severity == rules.Gate {
				g++
			} else {
				a++
			}
			findings = append(findings, rules.Finding{Rule: f.Rule, File: f.File, Key: f.Key, Severity: f.Severity})
		}
		gate, adv = &g, &a
	}
	if err == nil && s.opts.Clones {
		cr := clones.AnalyzeCached(snap, s.cfg, s.langs, nil, s.cache)
		n, f := len(cr.NearMisses), len(cr.Fragments)
		near, frag = &n, &f
	}
	s.mu.Lock()
	m := &s.metrics[i]
	m.Detail, m.Gate, m.Advisory, m.NearMisses, m.Fragments = true, gate, adv, near, frag
	if findings != nil {
		s.findings[i] = findings
	}
	s.mu.Unlock()
}

// fillDetailStep records defects introduced and resolved by commit i.
func (s *Server) fillDetailStep(i int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := &s.metrics[i]
	m.DetailStep = true
	if !s.opts.Defects {
		return
	}
	in, out, _ := rules.Compare(s.findings[i-1], s.findings[i])
	g, a, r := 0, 0, len(out)
	for _, f := range in {
		if f.Severity == rules.Gate {
			g++
		} else {
			a++
		}
	}
	m.IntroGate, m.IntroAdv, m.Resolved = &g, &a, &r
}
