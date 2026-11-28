// Package rulestest runs one rule implementation on in-memory sources, for
// the per-language rule tests.
package rulestest

import (
	"sort"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/model"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
)

// Adapter is what the harness needs from a language adapter.
type Adapter interface {
	lang.Adapter
	lang.SyntaxParser
}

// Run checks files (path -> content) with one rule's implementation for the
// adapter's language. Files whose path the adapter claims are parsed;
// every path is readable (manifests, go.mod ...).
func Run(t *testing.T, ad Adapter, register func(*rules.Registry) error, ruleID string, files map[string]string) []rules.Finding {
	t.Helper()
	reg := rules.NewRegistry()
	if err := catalog.Define(reg); err != nil {
		t.Fatal(err)
	}
	if err := register(reg); err != nil {
		t.Fatal(err)
	}
	rule, ok := reg.Rule(ruleID)
	if !ok {
		t.Fatalf("unknown rule %s", ruleID)
	}
	impl, ok := reg.Impl(ruleID, ad.ID())
	if !ok {
		t.Fatalf("%s not implemented for %s", ruleID, ad.ID())
	}
	var paths []string
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	var rf []*rules.File
	for _, p := range paths {
		claimed := false
		for _, e := range ad.Extensions() {
			if strings.HasSuffix(p, e) {
				claimed = true
			}
		}
		if !claimed {
			continue
		}
		src := []byte(files[p])
		pf, err := ad.Parse(src)
		if err != nil {
			t.Fatal(err)
		}
		if len(pf.ParseErrors) > 0 {
			t.Fatalf("%s: syntax errors %+v", p, pf.ParseErrors)
		}
		tree, err := ad.ParseSyntax(src)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(tree.Close)
		class := model.ClassSource
		if ad.IsTestPath(p) {
			class = model.ClassTest
		}
		rf = append(rf, rules.NewFile(p, class, pf, src, tree))
	}
	read := func(p string) ([]byte, error) {
		if s, ok := files[p]; ok {
			return []byte(s), nil
		}
		return nil, &missing{p}
	}
	ctx := rules.NewContext(rule, ad.ID(), rf, paths, read, config.Default())
	impl.Check(ctx)
	return ctx.Findings()
}

type missing struct{ p string }

func (m *missing) Error() string { return m.p + ": not found" }

// Summary renders findings as "severity line message" lines for assertions.
func Summary(fs []rules.Finding) []string {
	var out []string
	for _, f := range fs {
		out = append(out, string(f.Severity)+" "+f.File+":"+itoa(f.Line)+" "+f.Message)
	}
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
