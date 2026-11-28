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

type job struct {
	i       int
	measure bool // else fill commit i's step changes
}

// next picks the next job in the window, newest first: measure an
// unmeasured commit (the window plus the parent of its first commit, so the
// first step has a delta), else fill a step whose parent is measured.
// Called with s.mu held.
func (s *Server) next() (job, bool) {
	lo := max(0, s.lo-1)
	for i := s.hi; i >= lo; i-- {
		if !s.metrics[i].Done {
			return job{i, true}, true
		}
	}
	for i := s.hi; i >= s.lo; i-- {
		m := s.metrics[i]
		if m.StepDone || m.Error != "" {
			continue
		}
		if i == 0 {
			s.metrics[i].StepDone = true // no parent in the history
			continue
		}
		if p := s.metrics[i-1]; p.Done && p.Error == "" {
			return job{i, false}, true
		}
	}
	return job{}, false
}

// Run measures the window and fills step changes, waiting for new windows.
// It never returns.
func (s *Server) Run() {
	for {
		s.mu.Lock()
		j, ok := s.next()
		for !ok {
			s.work.Wait()
			j, ok = s.next()
		}
		s.mu.Unlock()
		if j.measure {
			m, findings := s.measure(s.commits[j.i].SHA)
			s.mu.Lock()
			s.metrics[j.i] = m
			if findings != nil {
				s.findings[j.i] = findings
			}
			s.mu.Unlock()
		} else {
			s.fillStep(j.i)
		}
	}
}

// measure computes one commit's absolute metrics and keeps its defect
// findings (slimmed) for step comparisons.
func (s *Server) measure(sha string) (m Metrics, findings []rules.Finding) {
	s.evMu.Lock()
	defer s.evMu.Unlock()
	m.Done = true
	_, snap, err := s.ev.ResolveSnapshot(sha)
	if err != nil {
		m.Error = err.Error()
		return m, nil
	}
	a := scorev2.Analyze(snap, s.cfg, scorev2.DefaultWeights())
	c := a.Codebase.Components
	m.Debt, m.DebtPerKLOC, m.Score = a.Codebase.Debt, a.Codebase.DebtPerKLOC, a.Codebase.Score
	m.Lines, m.Files, m.Callables, m.Hotspots, m.Decisions = c.SourceLines, c.SourceFiles, c.Callables, c.Hotspots, c.DecisionPoints
	m.Cloned, m.Violations, m.CycleUnits = c.ClonedCallables, c.RuleViolations, c.CycleUnits
	if s.opts.Defects {
		res := rules.Run(s.reg, snap, s.cfg, s.langs, rules.Options{})
		gate, adv := 0, 0
		findings = make([]rules.Finding, 0, len(res.Findings))
		for _, f := range res.Findings {
			if f.Severity == rules.Gate {
				gate++
			} else {
				adv++
			}
			findings = append(findings, rules.Finding{Rule: f.Rule, File: f.File, Key: f.Key, Severity: f.Severity})
		}
		m.Gate, m.Advisory = &gate, &adv
	}
	if s.opts.Clones {
		cr := clones.Analyze(snap, s.cfg, s.langs, nil)
		n, f := len(cr.NearMisses), len(cr.Fragments)
		m.NearMisses, m.Fragments = &n, &f
	}
	return m, findings
}

// fillStep records commit i's changes against its parent (i-1): changed
// measured-source lines, and defects introduced and resolved.
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
	defer s.mu.Unlock()
	m := &s.metrics[i]
	m.ChangedLines, m.StepDone = changed, true
	if s.opts.Defects {
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
}
