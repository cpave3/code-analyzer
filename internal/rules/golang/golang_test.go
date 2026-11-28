package golang

import (
	"strings"
	"testing"

	lgo "github.com/cameron/code-evaluator/internal/lang/golang"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
	"github.com/cameron/code-evaluator/internal/rules/rulestest"
)

type tc struct {
	name  string
	rule  string
	files map[string]string
	want  []string // substrings, one per expected finding, "gate:" or "advisory:" prefix
}

func check(t *testing.T, cases []tc) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := rulestest.Run(t, lgo.New(), Register, c.rule, c.files)
			if len(got) != len(c.want) {
				t.Fatalf("got %d findings, want %d:\n%s", len(got), len(c.want), strings.Join(rulestest.Summary(got), "\n"))
			}
			for i, w := range c.want {
				sev, sub, _ := strings.Cut(w, ":")
				if string(got[i].Severity) != sev || !strings.Contains(got[i].Message, sub) {
					t.Errorf("finding %d = %s %q, want %s containing %q", i, got[i].Severity, got[i].Message, sev, sub)
				}
				if got[i].Key == "" {
					t.Errorf("finding %d has no key", i)
				}
			}
		})
	}
}

const mod = "module example.com/m\n\ngo 1.22\n\nrequire (\n\tgithub.com/stretchr/testify v1.9.0\n\tgopkg.in/yaml.v3 v3.0.1\n\tgithub.com/x/tool v1.0.0\n\tgithub.com/x/ind v1.0.0 // indirect\n)\n\ntool github.com/x/tool/cmd/gen\n"

func TestDeadCode(t *testing.T) {
	check(t, []tc{
		{"imports", catalog.UnusedImport, map[string]string{"a/a.go": `package a
import (
	"fmt"
	str "strings"
	_ "embed"
	. "math"
	"gopkg.in/yaml.v3"
	"github.com/x/go-thing"
	"github.com/x/unused-lib"
	"os"
)
var _ = yaml.Marshal
var _ = thing.X
func f() { fmt.Println(Pi) }
`}, []string{"gate:import str \"strings\"", "advisory:\"github.com/x/unused-lib\" is never used (package name inferred", "gate:import \"os\" is never used"}},
		{"variables", catalog.UnusedVariable, map[string]string{"a/a.go": `package a
var pkgLevel = 1
func f() (int, error) {
	a, err := g()
	b, err := g()
	c := 1
	go func() { _ = c }()
	var d, _e int
	_ = d
	var unused string
	return a + b, nil
}
func g() (int, error) { return 0, nil }
`}, []string{"gate:err is declared", "gate:unused is declared"}},
		{"unreachable", catalog.UnreachableCode, map[string]string{"a/a.go": `package a
func f(x int) int {
	if x > 0 {
		return 1
		x++
	}
	switch x {
	case 1:
		panic("no")
		x--
	}
	goto end
end:
	for {
		break
	}
	return 0
	// trailing comment
}
`}, []string{"gate:`x++` follows `return 1`", "gate:`x--` follows `panic(\"no\")`"}},
		{"private functions", catalog.UnusedPrivateFunction, map[string]string{
			"a/a.go":      "package a\nfunc used() {}\nfunc onlyTests() {}\nfunc dead() {}\nfunc linked() {}\nfunc init() {}\nfunc (T) method() {}\ntype T struct{}\nfunc F() { used() }\n",
			"a/b.go":      "package a\nimport _ \"unsafe\"\n//go:linkname linked runtime.linked\n",
			"a/a_test.go": "package a\nfunc helper() { onlyTests() }\n",
			"b/b.go":      "package b\nfunc dead() {}\nfunc x() { dead() }\n",
		}, []string{"gate:unexported function dead", "gate:unexported function helper", "gate:unexported function x"}},
		{"exports only under internal", catalog.UnusedExport, map[string]string{
			"internal/a/a.go":      "package a\nfunc Used() {}\nfunc Dead() {}\ntype T struct{}\nfunc (T) Method() {}\n",
			"internal/b/b.go":      "package b\nimport \"example.com/m/internal/a\"\nfunc g() { a.Used() }\n",
			"pub/p.go":             "package pub\nfunc Public() {}\n",
			"internal/cmd/main.go": "package main\nfunc Exported() {}\nfunc main() {}\n",
		}, []string{"advisory:exported function Dead"}},
	})
}

func TestLostErrors(t *testing.T) {
	check(t, []tc{
		{"empty handlers", catalog.EmptyHandler, map[string]string{"a/a.go": `package a
func f() error {
	err := g()
	if err != nil {
	}
	if err != nil {
		// expected when the file is missing
	}
	if err != nil {
	} else {
		return nil
	}
	if x := 1; x != 0 {
	}
	defer func() { recover() }()
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	_ = err
	_ = x
	return err
}
var x int
func g() error { return nil }
`}, []string{"gate:empty `if err != nil {}`", "advisory:holds only a comment", "gate:recover() result is discarded", "advisory:`_ = err`"}},
		{"error to success", catalog.ErrorToSuccess, map[string]string{"a/a.go": `package a
func f() (int, error) {
	n, err := g()
	if err != nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	run := func() error {
		if loadErr := h(); loadErr != nil {
			return nil
		}
		return nil
	}
	_ = run
	return n, nil
}
func noError() int {
	if err := h(); err != nil {
		return 0
	}
	return 1
}
func g() (int, error) { return 0, nil }
func h() error { return nil }
`}, []string{"advisory:`if err != nil { return 0, nil }`", "advisory:`if loadErr != nil { return nil }`"}},
	})
}

func TestDependencies(t *testing.T) {
	files := map[string]string{
		"go.mod":              mod,
		"a/a.go":              "package a\nimport (\n\t\"fmt\"\n\t\"example.com/m/b\"\n\t\"example.com/m/nope\"\n\t\"example.com/m/onlytests\"\n\t\"gopkg.in/yaml.v3\"\n\t\"github.com/x/undeclared/sub\"\n\t\"example.com/other/pkg\"\n)\n",
		"a/a_test.go":         "package a\nimport \"github.com/stretchr/testify/assert\"\n",
		"b/b.go":              "package b\n",
		"onlytests/x_test.go": "package onlytests\n",
		"other/go.mod":        "module example.com/other\n\ngo 1.22\n",
		"other/pkg/p.go":      "package pkg\n",
	}
	check(t, []tc{
		{"unresolved", catalog.UnresolvedImport, files, []string{"gate:\"example.com/m/nope\" names no package", "gate:\"example.com/m/onlytests\" names no package"}},
		{"undeclared", catalog.UndeclaredDependency, files, []string{"gate:\"github.com/x/undeclared/sub\" is not required by go.mod"}},
		{"unused", catalog.UnusedDependency, files, nil},
	})
	files["go.mod"] = mod + "require github.com/x/never v1.0.0\n"
	check(t, []tc{{"unused direct require", catalog.UnusedDependency, files, []string{"advisory:requires github.com/x/never directly"}}})
}

func TestTestOracle(t *testing.T) {
	check(t, []tc{
		{"tautologies", catalog.TautologicalAssertion, map[string]string{"a/a_test.go": `package a
import (
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
func TestX(t *testing.T) {
	got, want := 1, 1
	assert.Equal(t, got, got)
	require.Equal(t, 1, 1, "msg")
	assert.Equal(t, want, got)
	assert.Equal(t, f(), f())
	assert.True(t, true)
	a := assert.New(t)
	a.Equal(got.Field, got.Field)
	_ = reflect.DeepEqual(got, got)
}
func f() int { return 1 }
`, "a/a.go": "package a\nimport \"reflect\"\nvar _ = reflect.DeepEqual(1, 1)\n"},
			[]string{"gate:compares got with itself", "gate:compares two literals", "gate:asserts a constant", "gate:compares got.Field with itself", "gate:compares got with itself"}},
		{"goroutines", catalog.LostAsyncAssertion, map[string]string{"a/a_test.go": `package a
import (
	"log"
	"testing"
	"github.com/stretchr/testify/require"
)
func TestX(t *testing.T) {
	go func() { t.Fatal("x") }()
	go func() { t.Error("x") }()
	go func() { log.Fatal("x") }()
	go func() { require.NoError(t, nil) }()
	t.Run("sub", func(st *testing.T) {
		go func() { st.FailNow() }()
	})
	t.Fatal("on the test goroutine")
}
`}, []string{"gate:`t.Fatal(\"x\")` runs inside a goroutine", "gate:`require.NoError(t, nil)`", "gate:`st.FailNow()`"}},
	})
}

func TestEveryFindingHasRuleMetadata(t *testing.T) {
	fs := rulestest.Run(t, lgo.New(), Register, catalog.UnusedImport, map[string]string{"a.go": "package a\nimport \"os\"\n"})
	if len(fs) != 1 || fs[0].Rule != catalog.UnusedImport || fs[0].Category != rules.CatDeadCode || fs[0].Language != "go" {
		t.Fatalf("%+v", fs)
	}
}
