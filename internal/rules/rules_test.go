package rules

import (
	"strings"
	"testing"
)

func TestRegistryValidation(t *testing.T) {
	r := NewRegistry()
	ok := Rule{ID: "dead-code/x", Category: CatDeadCode, Severity: Gate, Summary: "x"}
	if err := r.Define(ok); err != nil {
		t.Fatal(err)
	}
	for name, bad := range map[string]Rule{
		"duplicate":       ok,
		"category prefix": {ID: "other/x", Category: CatDeadCode, Severity: Gate, Summary: "x"},
		"severity":        {ID: "dead-code/y", Category: CatDeadCode, Severity: "warn", Summary: "x"},
		"summary":         {ID: "dead-code/z", Category: CatDeadCode, Severity: Gate},
	} {
		if err := r.Define(bad); err == nil {
			t.Errorf("%s: accepted %+v", name, bad)
		}
	}
	check := func(*Context) {}
	if err := r.Implement(Impl{Rule: "dead-code/missing", Language: "go", Check: check}); err == nil {
		t.Error("implementation of an undefined rule accepted")
	}
	if err := r.Implement(Impl{Rule: ok.ID, Language: "go", Check: check}); err != nil {
		t.Fatal(err)
	}
	if err := r.Implement(Impl{Rule: ok.ID, Language: "go", Check: check}); err == nil {
		t.Error("second go implementation accepted")
	}
	if err := r.Implement(Impl{Rule: ok.ID, Language: "typescript", Check: check}); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(r.Languages(ok.ID), ","); got != "go,typescript" {
		t.Errorf("languages %s", got)
	}
}

func TestCompareIsLineIndependentMultiset(t *testing.T) {
	f := func(key string, line int) Finding {
		return Finding{Rule: "dead-code/x", File: "a.go", Key: key, Line: line, Severity: Gate}
	}
	before := []Finding{f("a", 10), f("b", 20), f("b", 30)}
	after := []Finding{f("a", 15), f("b", 25), f("c", 40)}
	in, out, same := Compare(before, after)
	if same != 2 || len(in) != 1 || in[0].Key != "c" || in[0].Change != "introduced" || len(out) != 1 || out[0].Key != "b" || out[0].Change != "resolved" {
		t.Fatalf("introduced %+v resolved %+v persisting %d", in, out, same)
	}
}

func TestReportFailOn(t *testing.T) {
	reg := NewRegistry()
	_ = reg.Define(Rule{ID: "dead-code/x", Category: CatDeadCode, Severity: Gate, Summary: "x"})
	cov := []Coverage{{Rule: "dead-code/x", Language: "go", Status: "checked"}}
	base := &Result{Coverage: cov, Findings: []Finding{{Rule: "dead-code/x", File: "a", Key: "k", Severity: Gate}}}
	same := &Result{Coverage: cov, Findings: base.Findings}
	if r := NewReport(reg, base, same, FailGate); r.Failed || r.Persisting != 1 {
		t.Errorf("unchanged pre-existing finding failed the change: %+v", r)
	}
	adv := &Result{Coverage: cov, Findings: append(append([]Finding{}, base.Findings...), Finding{Rule: "dead-code/x", File: "a", Key: "k2", Severity: Advisory})}
	if r := NewReport(reg, base, adv, FailGate); r.Failed {
		t.Error("advisory finding failed the gate")
	}
	if r := NewReport(reg, base, adv, FailAny); !r.Failed {
		t.Error("fail-on any passed with an introduced finding")
	}
	if r := NewReport(reg, nil, base, FailGate); !r.Failed || r.Mode != "snapshot" {
		t.Error("snapshot gate finding did not fail fail-on gate")
	}
}
