// Package golang implements the defect rules for Go, on Tree-sitter syntax
// trees. Nothing is type-checked or built.
package golang

import (
	"path"
	"regexp"
	"strings"
	"unicode"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
)

// Language is the adapter ID these implementations target.
const Language = "go"

// Register adds the Go implementations.
func Register(r *rules.Registry) error {
	impls := []rules.Impl{
		{Rule: catalog.UnusedImport, Check: unusedImport,
			Note: "named imports are checked exactly; for unnamed imports the package name is inferred from the path (findings for non-standard-library paths are advisory)"},
		{Rule: catalog.UnusedVariable, Check: unusedVariable,
			Note: ":= and var declarations inside function and method bodies; the Go compiler rejects these, so a finding means the code does not build"},
		{Rule: catalog.UnreachableCode, Check: unreachable,
			Note: "statements after return, panic(), goto, break, continue in the same statement list; labelled statements are treated as reachable"},
		{Rule: catalog.UnusedPrivateFunction, Check: unusedPrivateFunction,
			Note: "unexported top-level functions (not methods) with no identifier of the same name in their package directory, tests and generated files included; //go:linkname and cgo //export are honoured"},
		{Rule: catalog.UnusedExport, Check: unusedExport,
			Note: "exported top-level functions in packages under an internal/ directory only (only the module can import those), never referenced by name anywhere"},
		{Rule: catalog.EmptyHandler, Check: emptyHandler,
			Note: "empty `if err != nil {}` branches, recover() results discarded, and `_ = err` (advisory)"},
		{Rule: catalog.ErrorToSuccess, Check: errorToSuccess,
			Note: "in functions whose last result is error: an `if err != nil` branch whose only statement returns nil as the error"},
		{Rule: catalog.UnresolvedImport, Check: unresolvedImport,
			Note: "imports inside a module present in the snapshot must name a directory with non-test .go files"},
		{Rule: catalog.UndeclaredDependency, Check: undeclaredDependency,
			Note: "non-standard-library imports must be covered by a require of the nearest go.mod or be a module in the snapshot"},
		{Rule: catalog.UnusedDependency, Check: unusedDependency,
			Note: "direct (non-// indirect) requires that no .go file of the module imports and no tool directive names"},
		{Rule: catalog.TautologicalAssertion, Check: tautological,
			Note: "testify assert/require comparisons (and method forms) and reflect.DeepEqual with identical call-free operands or two literals; True(true), False(false), Nil(nil)"},
		{Rule: catalog.LostAsyncAssertion, Check: lostAsync,
			Note: "t.Fatal*/FailNow/Skip* on the test's *testing.T/B/TB and require.* called inside a go statement (FailNow must run on the test goroutine)"},
	}
	for _, impl := range impls {
		impl.Language = Language
		if err := r.Implement(impl); err != nil {
			return err
		}
	}
	return nil
}

// ---- helpers ----

// funcName names the enclosing top-level declaration of n.
func funcName(f *rules.File, n *sitter.Node) string {
	d := n
	if d.Kind() != "function_declaration" && d.Kind() != "method_declaration" {
		d = rules.Ancestor(n, "function_declaration", "method_declaration")
	}
	if d == nil {
		return ""
	}
	name := f.Text(d.ChildByFieldName("name"))
	if d.Kind() == "method_declaration" {
		if recv := d.ChildByFieldName("receiver"); recv != nil {
			var t string
			rules.Walk(recv, func(c *sitter.Node) bool {
				if t == "" && c.Kind() == "type_identifier" {
					t = f.Text(c)
				}
				return t == ""
			})
			if t != "" {
				name = t + "." + name
			}
		}
	}
	return name
}

// idents counts identifier-like leaves by text, skipping subtrees of the
// given kinds.
func idents(f *rules.File, n *sitter.Node, kinds map[string]bool, skip ...string) map[string]int {
	out := map[string]int{}
	rules.Walk(n, func(c *sitter.Node) bool {
		for _, s := range skip {
			if c.Kind() == s {
				return false
			}
		}
		if kinds[c.Kind()] {
			out[f.Text(c)]++
		}
		return true
	})
	return out
}

var valueIdents = map[string]bool{"identifier": true}

func isExported(name string) bool {
	for _, r := range name {
		return unicode.IsUpper(r)
	}
	return false
}

func isErrName(s string) bool {
	l := strings.ToLower(s)
	return l == "err" || strings.HasSuffix(l, "err") || strings.HasSuffix(l, "error")
}

// errNotNil returns the identifier X of a condition `X != nil` with an
// error-like name, or "".
func errNotNil(f *rules.File, cond *sitter.Node) string {
	if cond == nil || cond.Kind() != "binary_expression" {
		return ""
	}
	op, l, r := cond.ChildByFieldName("operator"), cond.ChildByFieldName("left"), cond.ChildByFieldName("right")
	if op == nil || f.Text(op) != "!=" || l == nil || r == nil || r.Kind() != "nil" || l.Kind() != "identifier" || !isErrName(f.Text(l)) {
		return ""
	}
	return f.Text(l)
}

var versionElem = regexp.MustCompile(`^v[0-9]+$`)

// packageNames lists plausible package names for an import path: Go does
// not require the package name to match the path.
func packageNames(p string) []string {
	elems := strings.Split(p, "/")
	last := elems[len(elems)-1]
	if versionElem.MatchString(last) && len(elems) > 1 {
		last = elems[len(elems)-2]
	}
	if i := strings.Index(last, ".v"); i > 0 { // gopkg.in/yaml.v3
		last = last[:i]
	}
	cands := []string{last}
	trimmed := strings.TrimSuffix(strings.TrimSuffix(strings.TrimPrefix(last, "go-"), "-go"), ".go")
	cands = append(cands, trimmed, strings.ReplaceAll(trimmed, "-", ""), strings.ReplaceAll(trimmed, "-", "_"), strings.ReplaceAll(trimmed, ".", ""))
	return cands
}

func stdlib(p string) bool { return !strings.Contains(strings.Split(p, "/")[0], ".") }

// ---- dead code ----

func unusedImport(c *rules.Context) {
	for _, f := range c.Files {
		used := idents(f, f.Root(), map[string]bool{"identifier": true, "package_identifier": true}, "import_declaration")
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "import_spec" {
				return true
			}
			p := strings.Trim(f.Text(n.ChildByFieldName("path")), "\"`")
			if p == "C" || p == "" {
				return false
			}
			if name := n.ChildByFieldName("name"); name != nil {
				alias := f.Text(name)
				if alias == "_" || alias == "." || used[alias] > 0 {
					return false
				}
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Key: "import:" + p,
					Message: "import " + alias + " \"" + p + "\" is never used"})
				return false
			}
			for _, cand := range packageNames(p) {
				if used[cand] > 0 {
					return false
				}
			}
			sev, msg := rules.Gate, "import \""+p+"\" is never used"
			if !stdlib(p) || !regexp.MustCompile(`^[a-z][a-z0-9]*$`).MatchString(path.Base(p)) {
				sev, msg = rules.Advisory, msg+" (package name inferred from the path; it may differ)"
			}
			c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Key: "import:" + p, Severity: sev, Message: msg})
			return false
		})
	}
}

func unusedVariable(c *rules.Context) {
	for _, f := range c.Files {
		for _, fn := range rules.Named(f.Root()) {
			if fn.Kind() != "function_declaration" && fn.Kind() != "method_declaration" {
				continue
			}
			body := fn.ChildByFieldName("body")
			if body == nil {
				continue
			}
			var order []string
			first := map[string]*sitter.Node{}
			declCount := map[string]int{}
			addDecl := func(id *sitter.Node) {
				name := f.Text(id)
				if name == "_" || strings.HasPrefix(name, "_") {
					return
				}
				if first[name] == nil {
					first[name] = id
					order = append(order, name)
				}
				declCount[name]++
			}
			rules.Walk(body, func(n *sitter.Node) bool {
				switch n.Kind() {
				case "short_var_declaration":
					if left := n.ChildByFieldName("left"); left != nil {
						for _, id := range rules.Named(left) {
							if id.Kind() == "identifier" {
								addDecl(id)
							}
						}
					}
				case "var_spec":
					for i := uint(0); i < n.ChildCount(); i++ {
						if n.FieldNameForChild(uint32(i)) == "name" {
							addDecl(n.Child(i))
						}
					}
				}
				return true
			})
			if len(order) == 0 {
				continue
			}
			total := idents(f, body, valueIdents)
			for _, name := range order {
				if total[name]-declCount[name] > 0 {
					continue
				}
				n := first[name]
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: funcName(f, fn), Key: funcName(f, fn) + "|var:" + name,
					Message: "local variable " + name + " is declared and never used"})
			}
		}
	}
}

func terminates(f *rules.File, n *sitter.Node) bool {
	switch n.Kind() {
	case "return_statement", "goto_statement", "break_statement", "continue_statement", "fallthrough_statement":
		return true
	case "expression_statement":
		call := rules.Named(n)
		if len(call) == 1 && call[0].Kind() == "call_expression" {
			fn := call[0].ChildByFieldName("function")
			return fn != nil && fn.Kind() == "identifier" && f.Text(fn) == "panic"
		}
	}
	return false
}

func unreachable(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "statement_list" {
				return true
			}
			stmts := rules.Named(n)
			for i := 0; i+1 < len(stmts); i++ {
				if !terminates(f, stmts[i]) {
					continue
				}
				next := stmts[i+1]
				if next.Kind() == "labeled_statement" || next.Kind() == "empty_statement" {
					continue
				}
				sym := funcName(f, n)
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(next), Symbol: sym,
					Key:     sym + "|after:" + rules.Norm(f.Text(stmts[i])) + "|" + rules.Norm(f.Text(next)),
					Message: "unreachable: `" + rules.Norm(f.Text(next)) + "` follows `" + rules.Norm(f.Text(stmts[i])) + "`"})
				break
			}
			return true
		})
	}
}

// dirIdents counts identifiers per package directory across all files.
func dirIdents(c *rules.Context) map[string]map[string]int {
	return c.Memo("go:dirIdents", func() any {
		out := map[string]map[string]int{}
		for _, f := range c.Files {
			d := path.Dir(f.Path)
			if out[d] == nil {
				out[d] = map[string]int{}
			}
			for k, v := range idents(f, f.Root(), valueIdents) {
				out[d][k] += v
			}
		}
		return out
	}).(map[string]map[string]int)
}

// allNames counts identifiers and selector fields across every file.
func allNames(c *rules.Context) map[string]int {
	return c.Memo("go:allNames", func() any {
		out := map[string]int{}
		for _, f := range c.Files {
			for k, v := range idents(f, f.Root(), map[string]bool{"identifier": true, "field_identifier": true}) {
				out[k] += v
			}
		}
		return out
	}).(map[string]int)
}

func directives(c *rules.Context, dir string) string {
	var b strings.Builder
	for _, f := range c.Files {
		if path.Dir(f.Path) == dir && (strings.Contains(string(f.Src), "go:linkname") || strings.Contains(string(f.Src), "//export ")) {
			b.Write(f.Src)
		}
	}
	return b.String()
}

func unusedPrivateFunction(c *rules.Context) {
	counts := dirIdents(c)
	dirDirectives := map[string]string{}
	for _, f := range c.Files {
		dir := path.Dir(f.Path)
		for _, fn := range rules.Named(f.Root()) {
			if fn.Kind() != "function_declaration" {
				continue
			}
			name := f.Text(fn.ChildByFieldName("name"))
			if name == "" || isExported(name) || name == "init" || name == "main" || name == "_" {
				continue
			}
			if counts[dir][name] > 1 { // the declaration itself is one
				continue
			}
			if _, ok := dirDirectives[dir]; !ok {
				dirDirectives[dir] = directives(c, dir)
			}
			if d := dirDirectives[dir]; strings.Contains(d, "go:linkname "+name) || strings.Contains(d, "//export "+name) {
				continue
			}
			c.Report(rules.Finding{File: f.Path, Line: rules.Line(fn), Symbol: name, Key: "func:" + name,
				Message: "unexported function " + name + " is never referenced in package " + dir})
		}
	}
}

func underInternal(p string) bool {
	for _, e := range strings.Split(path.Dir(p), "/") {
		if e == "internal" {
			return true
		}
	}
	return false
}

func unusedExport(c *rules.Context) {
	names := allNames(c)
	for _, f := range c.Files {
		if f.Test() || !underInternal(f.Path) || (f.Parsed != nil && f.Parsed.PackageName == "main") {
			continue
		}
		for _, fn := range rules.Named(f.Root()) {
			if fn.Kind() != "function_declaration" {
				continue
			}
			name := f.Text(fn.ChildByFieldName("name"))
			if !isExported(name) || names[name] > 1 {
				continue
			}
			c.Report(rules.Finding{File: f.Path, Line: rules.Line(fn), Symbol: name, Key: "func:" + name,
				Message: "exported function " + name + " is not referenced anywhere in the repository (name-based)"})
		}
	}
}

// ---- lost errors ----

func emptyHandler(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			switch n.Kind() {
			case "if_statement":
				name := errNotNil(f, n.ChildByFieldName("condition"))
				cons := n.ChildByFieldName("consequence")
				if name == "" || cons == nil || n.ChildByFieldName("alternative") != nil || len(rules.Named(cons)) > 0 {
					return true
				}
				sev, msg := rules.Gate, "empty `if "+name+" != nil {}` branch drops the error"
				if rules.HasComment(cons) || rules.Contains(cons, "comment") {
					sev, msg = rules.Advisory, msg+" (the branch holds only a comment)"
				}
				sym := funcName(f, n)
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Severity: sev, Key: sym + "|empty-if:" + name, Message: msg})
			case "expression_statement":
				e := rules.Named(n)
				if len(e) == 1 && e[0].Kind() == "call_expression" {
					fn, args := e[0].ChildByFieldName("function"), e[0].ChildByFieldName("arguments")
					if fn != nil && f.Text(fn) == "recover" && len(rules.Named(args)) == 0 {
						sym := funcName(f, n)
						c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Key: sym + "|recover",
							Message: "recover() result is discarded: a panic is silently turned into normal return"})
					}
				}
			case "assignment_statement":
				l, r := n.ChildByFieldName("left"), n.ChildByFieldName("right")
				if l != nil && r != nil && strings.TrimSpace(f.Text(l)) == "_" {
					if rv := rules.Named(r); len(rv) == 1 && rv[0].Kind() == "identifier" && isErrName(f.Text(rv[0])) {
						sym := funcName(f, n)
						c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Severity: rules.Advisory, Key: sym + "|discard:" + f.Text(rv[0]),
							Message: "`_ = " + f.Text(rv[0]) + "` explicitly discards an error"})
					}
				}
			}
			return true
		})
	}
}

// returnsError reports whether a function's last result type is error.
func returnsError(f *rules.File, fn *sitter.Node) bool {
	res := fn.ChildByFieldName("result")
	if res == nil {
		return false
	}
	if res.Kind() == "type_identifier" {
		return f.Text(res) == "error"
	}
	if res.Kind() == "parameter_list" {
		ps := rules.Named(res)
		if len(ps) == 0 {
			return false
		}
		t := ps[len(ps)-1].ChildByFieldName("type")
		return t != nil && f.Text(t) == "error"
	}
	return false
}

func errorToSuccess(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(fn *sitter.Node) bool {
			k := fn.Kind()
			if k != "function_declaration" && k != "method_declaration" && k != "func_literal" {
				return true
			}
			body := fn.ChildByFieldName("body")
			if body == nil || !returnsError(f, fn) {
				return true
			}
			var visit func(n *sitter.Node)
			visit = func(n *sitter.Node) {
				if n.Kind() == "func_literal" {
					return // checked on its own
				}
				if n.Kind() == "if_statement" {
					name := errNotNil(f, n.ChildByFieldName("condition"))
					cons := n.ChildByFieldName("consequence")
					if name != "" && cons != nil && n.ChildByFieldName("alternative") == nil {
						var stmts []*sitter.Node
						for _, sl := range rules.Named(cons) {
							stmts = append(stmts, rules.Named(sl)...)
						}
						if len(stmts) == 1 && stmts[0].Kind() == "return_statement" {
							if vals := rules.Named(stmts[0]); len(vals) == 1 {
								exprs := rules.Named(vals[0])
								if vals[0].Kind() != "expression_list" {
									exprs = vals
								}
								if len(exprs) > 0 && exprs[len(exprs)-1].Kind() == "nil" {
									sym := funcName(f, n)
									c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Key: sym + "|if:" + name + "|" + rules.Norm(f.Text(stmts[0])),
										Message: "`if " + name + " != nil { " + rules.Norm(f.Text(stmts[0])) + " }` reports success for a failure without using " + name})
								}
							}
						}
					}
				}
				for i := uint(0); i < n.ChildCount(); i++ {
					visit(n.Child(i))
				}
			}
			visit(body)
			return true
		})
	}
}

// ---- test oracle ----

var comparators = map[string]bool{"Equal": true, "Equalf": true, "EqualValues": true, "EqualValuesf": true, "Exactly": true, "Exactlyf": true,
	"NotEqual": true, "NotEqualf": true, "Same": true, "Samef": true, "ElementsMatch": true, "ElementsMatchf": true, "JSONEq": true, "YAMLEq": true}

var goLiterals = map[string]bool{"int_literal": true, "float_literal": true, "interpreted_string_literal": true, "raw_string_literal": true,
	"rune_literal": true, "true": true, "false": true, "nil": true, "imaginary_literal": true}

func tautological(c *rules.Context) {
	for _, f := range c.Files {
		if !f.Test() {
			continue
		}
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "call_expression" {
				return true
			}
			a, b, why := assertion(f, n)
			if a != nil {
				why = rules.Tautology(a, b, f.Src, func(x *sitter.Node) bool { return goLiterals[x.Kind()] }, "call_expression")
			}
			if why != "" {
				text := rules.Norm(f.Text(n))
				c.At(f, n, funcName(f, n), "assert:"+text, "`"+text+"` "+why+" and cannot fail", "")
			}
			return true
		})
	}
}

// assertion returns the compared operands of a recognised assertion call,
// or a reason when it asserts a constant.
func assertion(f *rules.File, n *sitter.Node) (a, b *sitter.Node, why string) {
	fn := n.ChildByFieldName("function")
	if fn == nil || fn.Kind() != "selector_expression" {
		return nil, nil, ""
	}
	recv, method := f.Text(fn.ChildByFieldName("operand")), f.Text(fn.ChildByFieldName("field"))
	args := rules.Named(n.ChildByFieldName("arguments"))
	pkgForm := recv == "assert" || recv == "require"
	switch {
	case recv == "reflect" && method == "DeepEqual" && len(args) == 2:
		return args[0], args[1], ""
	case comparators[method] && pkgForm && len(args) >= 3:
		return args[1], args[2], ""
	case comparators[method] && !pkgForm && len(args) >= 2 && fn.ChildByFieldName("operand").Kind() == "identifier" && f.Text(args[0]) != "t":
		return args[0], args[1], ""
	case pkgForm && len(args) >= 2 && constantAssert[method] == args[1].Kind():
		return nil, nil, "asserts a constant"
	}
	return nil, nil, ""
}

var constantAssert = map[string]string{"True": "true", "False": "false", "Nil": "nil"}

var fatalMethods = map[string]bool{"Fatal": true, "Fatalf": true, "FailNow": true, "Skip": true, "Skipf": true, "SkipNow": true}

// testingParams returns the names of *testing.T / *testing.B / testing.TB
// parameters of a function.
func testingParams(f *rules.File, fn *sitter.Node) map[string]bool {
	out := map[string]bool{}
	params := fn.ChildByFieldName("parameters")
	for _, p := range rules.Named(params) {
		t := strings.ReplaceAll(f.Text(p.ChildByFieldName("type")), " ", "")
		if t == "*testing.T" || t == "*testing.B" || t == "testing.TB" || t == "*testing.F" {
			for i := uint(0); i < p.ChildCount(); i++ {
				if p.FieldNameForChild(uint32(i)) == "name" {
					out[f.Text(p.Child(i))] = true
				}
			}
		}
	}
	return out
}

func lostAsync(c *rules.Context) {
	for _, f := range c.Files {
		if !f.Test() {
			continue
		}
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "go_statement" {
				return true
			}
			tnames := map[string]bool{}
			for p := n.Parent(); p != nil; p = p.Parent() {
				if k := p.Kind(); k == "function_declaration" || k == "method_declaration" || k == "func_literal" {
					for name := range testingParams(f, p) {
						tnames[name] = true
					}
				}
			}
			rules.Walk(n, func(call *sitter.Node) bool {
				if call.Kind() != "call_expression" {
					return true
				}
				fn := call.ChildByFieldName("function")
				if fn == nil || fn.Kind() != "selector_expression" {
					return true
				}
				recv, method := f.Text(fn.ChildByFieldName("operand")), f.Text(fn.ChildByFieldName("field"))
				if (tnames[recv] && fatalMethods[method]) || recv == "require" {
					sym := funcName(f, n)
					c.Report(rules.Finding{File: f.Path, Line: rules.Line(call), Symbol: sym, Key: sym + "|go:" + rules.Norm(f.Text(call)),
						Message: "`" + rules.Norm(f.Text(call)) + "` runs inside a goroutine: FailNow there does not stop the test and the failure can be lost"})
				}
				return true
			})
			return false
		})
	}
}
