// Package csharp implements the defect rules for C# on Tree-sitter syntax
// trees. Nothing is compiled or type-checked.
//
// Not implemented for C# (recorded in internal/rules/all): unused usings
// (needs type and namespace resolution), unused public API (C# public members
// are largely consumed through DI, reflection and framework conventions),
// and the dependency rules (need .csproj PackageReference and a
// namespace-to-package mapping).
package csharp

import (
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
)

// Language is the adapter ID these implementations target.
const Language = "csharp"

// Register adds the C# implementations.
func Register(r *rules.Registry) error {
	impls := []rules.Impl{
		{Rule: catalog.UnusedVariable, FileLocal: true, Check: unusedVariable,
			Note: "locals declared with a plain name inside a member and never referenced there; `using var` declarations (disposal is the point), out and pattern variables are skipped"},
		{Rule: catalog.UnreachableCode, FileLocal: true, Check: unreachable,
			Note: "statements after return, throw, break, continue or goto in the same block or switch section; local functions are ignored"},
		{Rule: catalog.UnusedPrivateFunction, Check: unusedPrivate,
			Note: "private (or unmodified) methods without attributes, override, partial, extern or explicit interface names, whose name appears nowhere else in the repository's C# (name-based, so partial classes across files are covered)"},
		{Rule: catalog.EmptyHandler, FileLocal: true, Check: emptyCatch,
			Note: "empty catch blocks (advisory when they hold only a comment)"},
		{Rule: catalog.ErrorToSuccess, FileLocal: true, Check: catchConstant,
			Note: "catch blocks whose only statement returns a constant (null, default, a literal, an empty collection) or nothing"},
		{Rule: catalog.TautologicalAssertion, FileLocal: true, Check: tautological,
			Note: "xUnit/NUnit/MSTest Assert.Equal/AreEqual/Same/NotEqual (and Assert.That(x, Is.EqualTo(x))) and FluentAssertions Should().Be/BeEquivalentTo comparing call-free identical expressions or two literals; Assert.True(true)/IsTrue(true)/False(false)"},
		{Rule: catalog.LostAsyncAssertion, FileLocal: true, Check: lostAsync,
			Note: "Assert.*Async(...) and FluentAssertions .Should()...Async() assertions used as a bare statement (not awaited)"},
	}
	for _, impl := range impls {
		impl.Language = Language
		if err := r.Implement(impl); err != nil {
			return err
		}
	}
	r.Prepare(Language, func(f *rules.File) { fileNames(f); filePrivate(f) })
	return nil
}

var memberKinds = map[string]bool{
	"method_declaration": true, "constructor_declaration": true, "destructor_declaration": true, "operator_declaration": true,
	"conversion_operator_declaration": true, "property_declaration": true, "indexer_declaration": true,
}

// symbol names the enclosing member of n as Type.Member.
func symbol(f *rules.File, n *sitter.Node) string {
	member, typ := "", ""
	for p := n; p != nil; p = p.Parent() {
		k := p.Kind()
		if member == "" && memberKinds[k] {
			member = f.Text(p.ChildByFieldName("name"))
			if member == "" {
				member = strings.TrimSuffix(k, "_declaration")
			}
		}
		if strings.HasSuffix(k, "class_declaration") || k == "struct_declaration" || k == "record_declaration" || k == "interface_declaration" {
			if typ = f.Text(p.ChildByFieldName("name")); typ != "" {
				break
			}
		}
	}
	switch {
	case typ != "" && member != "":
		return typ + "." + member
	case member != "":
		return member
	}
	return typ
}

func at(c *rules.Context, f *rules.File, n *sitter.Node, key, msg string, sev rules.Severity) {
	c.At(f, n, symbol(f, n), key, msg, sev)
}

// ---- dead code ----

func unreachable(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "block" && n.Kind() != "switch_section" {
				return true
			}
			stmts := rules.Named(n)
			for i, s := range stmts {
				switch s.Kind() {
				case "return_statement", "throw_statement", "break_statement", "continue_statement", "goto_statement":
				default:
					continue
				}
				for _, next := range stmts[i+1:] {
					if !strings.HasSuffix(next.Kind(), "_statement") || next.Kind() == "local_function_statement" || next.Kind() == "empty_statement" || next.Kind() == "labeled_statement" {
						continue
					}
					at(c, f, next, "after:"+rules.Norm(f.Text(s))+"|"+rules.Norm(f.Text(next)),
						"unreachable: `"+rules.Norm(f.Text(next))+"` follows `"+rules.Norm(f.Text(s))+"`", "")
					break
				}
				break
			}
			return true
		})
	}
}

func unusedVariable(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(m *sitter.Node) bool {
			if !memberKinds[m.Kind()] && m.Kind() != "accessor_declaration" {
				return true
			}
			if m.Kind() == "accessor_declaration" {
				return false // handled with its property
			}
			names := map[string]int{}
			rules.Walk(m, func(n *sitter.Node) bool {
				if n.Kind() == "identifier" {
					names[f.Text(n)]++
				}
				return true
			})
			rules.Walk(m, func(n *sitter.Node) bool {
				if n.Kind() != "local_declaration_statement" {
					return true
				}
				for i := uint(0); i < n.ChildCount(); i++ {
					if n.Child(i).Kind() == "using" {
						return false // `using var x = ...`: disposal is the point
					}
				}
				rules.Walk(n, func(d *sitter.Node) bool {
					if d.Kind() != "variable_declarator" {
						return true
					}
					nm := d.ChildByFieldName("name")
					if nm == nil || nm.Kind() != "identifier" {
						return false
					}
					name := f.Text(nm)
					if strings.HasPrefix(name, "_") || names[name] > 1 {
						return false
					}
					at(c, f, nm, "var:"+name, "local variable "+name+" is declared and never used", "")
					return false
				})
				return false
			})
			return false
		})
	}
}

// fileNames counts identifiers in a file (cached by content).
func fileNames(f *rules.File) map[string]int {
	return f.Memo("cs:names", func() any {
		out := map[string]int{}
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() == "identifier" {
				out[f.Text(n)]++
			}
			return true
		})
		return out
	}).(map[string]int)
}

type privateMethod struct {
	name, sym string
	line      int
}

// filePrivate lists private methods that are candidates for the unused
// check (cached by content).
func filePrivate(f *rules.File) []privateMethod {
	return f.Memo("cs:private", func() any {
		var out []privateMethod
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "method_declaration" {
				return true
			}
			name := f.Text(n.ChildByFieldName("name"))
			if name == "" || strings.Contains(name, ".") || n.ChildByFieldName("body") == nil {
				return false
			}
			private, skip := true, false
			for i := uint(0); i < n.NamedChildCount(); i++ {
				ch := n.NamedChild(i)
				switch ch.Kind() {
				case "attribute_list":
					skip = true // framework entry points ([Fact], [HttpGet], ...)
				case "explicit_interface_specifier":
					skip = true
				case "modifier":
					switch f.Text(ch) {
					case "public", "protected", "internal":
						private = false
					case "override", "partial", "extern", "virtual", "abstract":
						skip = true
					}
				}
			}
			if private && !skip {
				out = append(out, privateMethod{name, symbol(f, n), rules.Line(n)})
			}
			return false
		})
		return out
	}).([]privateMethod)
}

func unusedPrivate(c *rules.Context) {
	total := c.Memo("cs:allNames", func() any {
		out := map[string]int{}
		for _, f := range c.Files {
			for k, v := range fileNames(f) {
				out[k] += v
			}
		}
		return out
	}).(map[string]int)
	for _, f := range c.Files {
		for _, m := range filePrivate(f) {
			if total[m.name] > 1 {
				continue
			}
			c.Report(rules.Finding{File: f.Path, Line: m.line, Symbol: m.sym, Key: "method:" + m.sym,
				Message: "private method " + m.sym + " is never referenced (name-based, across the repository's C#)"})
		}
	}
}

// ---- lost errors ----

func emptyCatch(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "catch_clause" {
				return true
			}
			body := n.ChildByFieldName("body")
			if body == nil {
				for _, ch := range rules.Named(n) {
					if ch.Kind() == "block" {
						body = ch
					}
				}
			}
			if body == nil || len(rules.Named(body)) > 0 {
				return true
			}
			sev, msg := rules.Gate, "empty catch block discards the exception"
			if rules.HasComment(body) {
				sev, msg = rules.Advisory, msg+" (the block holds only a comment)"
			}
			at(c, f, n, "empty-catch", msg, sev)
			return true
		})
	}
}

var constKinds = map[string]bool{"null_literal": true, "boolean_literal": true, "integer_literal": true, "real_literal": true,
	"string_literal": true, "character_literal": true, "default_expression": true, "verbatim_string_literal": true}

func constant(f *rules.File, n *sitter.Node) bool {
	if n == nil {
		return true
	}
	if constKinds[n.Kind()] || f.Text(n) == "default" {
		return true
	}
	return n.Kind() == "collection_expression" && len(rules.Named(n)) == 0
}

func catchConstant(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "catch_clause" {
				return true
			}
			var body *sitter.Node
			for _, ch := range rules.Named(n) {
				if ch.Kind() == "block" {
					body = ch
				}
			}
			stmts := rules.Named(body)
			if len(stmts) != 1 || stmts[0].Kind() != "return_statement" {
				return true
			}
			vals := rules.Named(stmts[0])
			var v *sitter.Node
			if len(vals) > 0 {
				v = vals[0]
			}
			if !constant(f, v) {
				return true
			}
			at(c, f, n, "catch:"+rules.Norm(f.Text(stmts[0])),
				"catch block only does `"+rules.Norm(f.Text(stmts[0]))+"`: the failure becomes an ordinary value and the exception is not used", "")
			return true
		})
	}
}

// ---- test oracle ----

var assertEq = map[string]bool{"Equal": true, "AreEqual": true, "StrictEqual": true, "Same": true, "AreSame": true,
	"NotEqual": true, "AreNotEqual": true, "Equivalent": true}
var fluentEq = map[string]bool{"Be": true, "BeEquivalentTo": true, "BeSameAs": true, "NotBe": true}

// method returns the name of an invocation's member access and its receiver.
func method(f *rules.File, call *sitter.Node) (name string, recv *sitter.Node) {
	fn := call.ChildByFieldName("function")
	if fn == nil || fn.Kind() != "member_access_expression" {
		return "", nil
	}
	nm := fn.ChildByFieldName("name")
	if nm != nil && nm.Kind() == "generic_name" {
		nm = nm.NamedChild(0)
	}
	return f.Text(nm), fn.ChildByFieldName("expression")
}

func args(call *sitter.Node) []*sitter.Node {
	var out []*sitter.Node
	for _, a := range rules.Named(call.ChildByFieldName("arguments")) {
		if a.Kind() == "argument" {
			if v := rules.Named(a); len(v) > 0 {
				out = append(out, v[len(v)-1])
			}
		}
	}
	return out
}

func literal(n *sitter.Node) bool { return constKinds[n.Kind()] }

// assertion returns compared operands, or a reason for a constant assertion.
func assertion(f *rules.File, call *sitter.Node) (a, b *sitter.Node, why string) {
	name, recv := method(f, call)
	as := args(call)
	if recv != nil && f.Text(recv) == "Assert" {
		switch {
		case assertEq[name] && len(as) >= 2:
			return as[0], as[1], ""
		case name == "That" && len(as) >= 2 && as[1].Kind() == "invocation_expression":
			if inner, r := method(f, as[1]); inner == "EqualTo" && r != nil && f.Text(r) == "Is" {
				if ia := args(as[1]); len(ia) == 1 {
					return as[0], ia[0], ""
				}
			}
		case (name == "True" || name == "IsTrue") && len(as) >= 1 && f.Text(as[0]) == "true",
			(name == "False" || name == "IsFalse") && len(as) >= 1 && f.Text(as[0]) == "false":
			return nil, nil, "asserts a constant"
		}
	}
	if fluentEq[name] && len(as) >= 1 && recv != nil && recv.Kind() == "invocation_expression" {
		if should, subject := method(f, recv); should == "Should" && subject != nil {
			return subject, as[0], ""
		}
	}
	return nil, nil, ""
}

func tautological(c *rules.Context) {
	for _, f := range c.Files {
		if !f.Test() {
			continue
		}
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "invocation_expression" {
				return true
			}
			a, b, why := assertion(f, n)
			if a != nil {
				why = rules.Tautology(a, b, f.Source(), literal, "invocation_expression", "object_creation_expression", "await_expression")
			}
			if why == "" {
				return true
			}
			text := rules.Norm(f.Text(n))
			at(c, f, n, "assert:"+text, "`"+text+"` "+why+" and cannot fail", "")
			return false
		})
	}
}

// asyncAssertion reports whether a call is an asynchronous assertion:
// Assert.XxxAsync(...) or a FluentAssertions chain through .Should() ending
// in an ...Async method.
func asyncAssertion(f *rules.File, call *sitter.Node) bool {
	name, recv := method(f, call)
	if !strings.HasSuffix(name, "Async") || recv == nil {
		return false
	}
	if f.Text(recv) == "Assert" {
		return true
	}
	for r := recv; r != nil; {
		if r.Kind() != "invocation_expression" {
			return false
		}
		n, next := method(f, r)
		if n == "Should" {
			return true
		}
		r = next
	}
	return false
}

func lostAsync(c *rules.Context) {
	for _, f := range c.Files {
		if !f.Test() {
			continue
		}
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "expression_statement" {
				return true
			}
			e := rules.Named(n)
			if len(e) == 1 && e[0].Kind() == "invocation_expression" && asyncAssertion(f, e[0]) {
				text := rules.Norm(f.Text(e[0]))
				at(c, f, n, "unawaited:"+text, "`"+text+"` is not awaited: the test can finish before it completes and a failure is lost", "")
			}
			return true
		})
	}
}
