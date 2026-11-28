// Package typescript implements the defect rules for TypeScript (.ts, .tsx,
// .mts, .cts) on Tree-sitter syntax trees. Nothing is type-checked, built or
// installed; node_modules are never read.
package typescript

import (
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
)

// Language is the adapter ID these implementations target.
const Language = "typescript"

// Register adds the TypeScript implementations.
func Register(r *rules.Registry) error {
	impls := []rules.Impl{
		{Rule: catalog.UnusedImport, FileLocal: true, Check: unusedImport,
			Note: "default, namespace and named import bindings (type-only included) never referenced in the file; a default React import counts as used when the file contains JSX"},
		{Rule: catalog.UnusedVariable, FileLocal: true, Check: unusedVariable,
			Note: "const/let/var bindings with a plain identifier name inside functions, and non-exported module-level bindings in module files; destructuring and loop variables are skipped"},
		{Rule: catalog.UnreachableCode, FileLocal: true, Check: unreachable,
			Note: "statements after return, throw, break or continue in the same block or case; hoisted function declarations and type-only declarations are ignored"},
		{Rule: catalog.UnusedPrivateFunction, FileLocal: true, Check: unusedPrivateFunction,
			Note: "non-exported module-level functions (declarations and const arrow/function expressions) unused in their module, and private/#private class methods unused in their class"},
		{Rule: catalog.UnusedExport, Check: unusedExport,
			Note: "exported functions of a module that other files import only by name (no namespace, wildcard re-export or dynamic import), where no file imports that name and the module does not use it itself"},
		{Rule: catalog.EmptyHandler, FileLocal: true, Check: emptyHandler,
			Note: "empty catch blocks and empty .catch()/.then(_, onRejected) callbacks"},
		{Rule: catalog.ErrorToSuccess, FileLocal: true, Check: errorToSuccess,
			Note: "catch blocks whose only statement returns a constant (null, undefined, [], {}, a literal) and .catch(() => constant) callbacks"},
		{Rule: catalog.UnresolvedImport, Check: unresolvedImport,
			Note: "relative imports must resolve to a file in the snapshot (TypeScript, JavaScript, JSON or the named asset); unresolved tsconfig-alias and workspace imports are advisory"},
		{Rule: catalog.UndeclaredDependency, Check: undeclaredDependency,
			Note: "bare imports must be declared (dependencies, devDependencies, peer or optional, or an @types package) in a package.json between the file and the repository root; node builtins, virtual modules and ambient `declare module` names are exempt"},
		{Rule: catalog.UnusedDependency, Check: unusedDependency,
			Note: "package.json \"dependencies\" (not @types) that no file under the package imports or names in quotes (TypeScript imports, JS/CSS/HTML/config text, package scripts)"},
		{Rule: catalog.TautologicalAssertion, FileLocal: true, Check: tautological,
			Note: "expect(A).toBe/toEqual/toStrictEqual(A) and chai equal/eql forms, node:assert equal/deepEqual families, and constant truthiness assertions; .not chains are skipped"},
		{Rule: catalog.LostAsyncAssertion, FileLocal: true, Check: lostAsync,
			Note: "expect(...).resolves/.rejects chains used as a bare statement (neither awaited nor returned)"},
	}
	for _, impl := range impls {
		impl.Language = Language
		if err := r.Implement(impl); err != nil {
			return err
		}
	}
	r.Prepare(Language, prepare)
	return nil
}

// ---- helpers ----

var functionKinds = map[string]bool{"function_declaration": true, "generator_function_declaration": true, "function_expression": true,
	"function": true, "generator_function": true, "arrow_function": true, "method_definition": true}

// symbol names the enclosing declaration of n: a named function, a method,
// a variable bound to a function, or a test case (it/test("title")).
func symbol(f *rules.File, n *sitter.Node) string {
	for p := n; p != nil; p = p.Parent() {
		switch p.Kind() {
		case "function_declaration", "generator_function_declaration":
			return f.Text(p.ChildByFieldName("name"))
		case "method_definition":
			name := f.Text(p.ChildByFieldName("name"))
			if cls := rules.Ancestor(p, "class_declaration", "abstract_class_declaration", "class"); cls != nil {
				if cn := cls.ChildByFieldName("name"); cn != nil {
					return f.Text(cn) + "." + name
				}
			}
			return name
		case "variable_declarator":
			if v := p.ChildByFieldName("value"); v != nil && functionKinds[v.Kind()] {
				return f.Text(p.ChildByFieldName("name"))
			}
		case "call_expression":
			fn := p.ChildByFieldName("function")
			if fn != nil && fn.Kind() == "identifier" {
				switch f.Text(fn) {
				case "it", "test", "describe":
					if args := rules.Named(p.ChildByFieldName("arguments")); len(args) > 0 && args[0].Kind() == "string" {
						return f.Text(fn) + "(" + rules.Norm(f.Text(args[0])) + ")"
					}
				}
			}
		}
	}
	return ""
}

var useKinds = map[string]bool{"identifier": true, "type_identifier": true, "shorthand_property_identifier": true}

func countIdents(f *rules.File, n *sitter.Node, kinds map[string]bool, skip ...string) map[string]int {
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

// isModule reports whether a file is an ES module (has import or export);
// top-level bindings of scripts are global and may be used elsewhere.
func isModule(f *rules.File) bool {
	for _, n := range rules.Named(f.Root()) {
		if n.Kind() == "import_statement" || n.Kind() == "export_statement" {
			return true
		}
	}
	return false
}

func isDeclarationFile(p string) bool {
	return strings.HasSuffix(p, ".d.ts") || strings.HasSuffix(p, ".d.mts") || strings.HasSuffix(p, ".d.cts")
}

// functionValue unwraps parentheses/as/satisfies and wrapper calls
// (memo(), forwardRef()) to find a function-valued initialiser.
func functionValue(v *sitter.Node) bool {
	for v != nil && (v.Kind() == "parenthesized_expression" || v.Kind() == "as_expression" || v.Kind() == "satisfies_expression") {
		v = v.NamedChild(0)
	}
	if v == nil {
		return false
	}
	if functionKinds[v.Kind()] {
		return true
	}
	if v.Kind() == "call_expression" {
		if args := rules.Named(v.ChildByFieldName("arguments")); len(args) > 0 {
			return functionValue(args[0])
		}
	}
	return false
}

var literalKinds = map[string]bool{"number": true, "string": true, "true": true, "false": true, "null": true, "undefined": true, "regex": true}

func isLiteral(f *rules.File, n *sitter.Node) bool {
	if n == nil {
		return false
	}
	if literalKinds[n.Kind()] {
		return true
	}
	if n.Kind() == "identifier" && f.Text(n) == "undefined" {
		return true
	}
	if n.Kind() == "template_string" && !rules.Contains(n, "template_substitution") {
		return true
	}
	if n.Kind() == "unary_expression" {
		if a := n.ChildByFieldName("argument"); a != nil && a.Kind() == "number" {
			return true
		}
	}
	return false
}

// constantValue: a literal or an empty array/object.
func constantValue(f *rules.File, n *sitter.Node) bool {
	if isLiteral(f, n) {
		return true
	}
	return n != nil && (n.Kind() == "array" || n.Kind() == "object") && len(rules.Named(n)) == 0
}

func hasJSX(f *rules.File) bool {
	return rules.Contains(f.Root(), "jsx_element", "jsx_self_closing_element")
}

// ---- dead code ----

func unusedImport(c *rules.Context) {
	for _, f := range c.Files {
		uses := countIdents(f, f.Root(), useKinds, "import_statement")
		jsx := -1
		for _, st := range rules.Named(f.Root()) {
			if st.Kind() != "import_statement" {
				continue
			}
			src := f.Text(st.ChildByFieldName("source"))
			spec := strings.Trim(src, "\"'`")
			var clause *sitter.Node
			for _, ch := range rules.Named(st) {
				if ch.Kind() == "import_clause" {
					clause = ch
				}
			}
			if clause == nil {
				continue // side-effect import
			}
			var bindings []*sitter.Node
			rules.Walk(clause, func(n *sitter.Node) bool {
				switch n.Kind() {
				case "import_specifier":
					if a := n.ChildByFieldName("alias"); a != nil {
						bindings = append(bindings, a)
					} else if nm := n.ChildByFieldName("name"); nm != nil {
						bindings = append(bindings, nm)
					}
					return false
				case "identifier":
					bindings = append(bindings, n) // default or namespace binding
				}
				return true
			})
			for _, b := range bindings {
				name := f.Text(b)
				if uses[name] > 0 {
					continue
				}
				if name == "React" {
					if jsx < 0 {
						jsx = 0
						if hasJSX(f) {
							jsx = 1
						}
					}
					if jsx == 1 {
						continue
					}
				}
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(b), Key: "import:" + spec + "#" + name,
					Message: "imported " + name + " from \"" + spec + "\" is never used"})
			}
		}
	}
}

func unusedVariable(c *rules.Context) {
	for _, f := range c.Files {
		if isDeclarationFile(f.Path) {
			continue
		}
		module := isModule(f)
		scopeCounts := map[*sitter.Node]map[string]int{}
		declCounts := map[*sitter.Node]map[string]int{}
		type decl struct {
			name  *sitter.Node
			scope *sitter.Node
		}
		var decls []decl
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() == "ambient_declaration" {
				return false
			}
			if n.Kind() != "variable_declarator" {
				return true
			}
			decl0 := n.Parent()
			if decl0 == nil || (decl0.Kind() != "lexical_declaration" && decl0.Kind() != "variable_declaration") {
				return true
			}
			name := n.ChildByFieldName("name")
			if name == nil || name.Kind() != "identifier" || strings.HasPrefix(f.Text(name), "_") {
				return true
			}
			var scope *sitter.Node
			for p := decl0.Parent(); p != nil; p = p.Parent() {
				if functionKinds[p.Kind()] {
					scope = p
					break
				}
			}
			if scope == nil {
				if !module || functionValue(n.ChildByFieldName("value")) || (decl0.Parent() != nil && decl0.Parent().Kind() == "export_statement") {
					return true // scripts are global; module-level functions belong to unused-private-function
				}
				scope = f.Root()
			}
			decls = append(decls, decl{name, scope})
			return true
		})
		for _, d := range decls {
			if scopeCounts[d.scope] == nil {
				scopeCounts[d.scope] = countIdents(f, d.scope, useKinds)
				dc := map[string]int{}
				rules.Walk(d.scope, func(n *sitter.Node) bool {
					if n.Kind() == "variable_declarator" {
						if nm := n.ChildByFieldName("name"); nm != nil && nm.Kind() == "identifier" {
							dc[f.Text(nm)]++
						}
					}
					return true
				})
				declCounts[d.scope] = dc
			}
			name := f.Text(d.name)
			if scopeCounts[d.scope][name]-declCounts[d.scope][name] > 0 {
				continue
			}
			sym := symbol(f, d.name.Parent().Parent())
			c.Report(rules.Finding{File: f.Path, Line: rules.Line(d.name), Symbol: sym, Key: sym + "|var:" + name,
				Message: "variable " + name + " is declared and never used"})
		}
	}
}

var hoisted = map[string]bool{"function_declaration": true, "generator_function_declaration": true, "interface_declaration": true,
	"type_alias_declaration": true, "empty_statement": true, "ambient_declaration": true}

func unreachable(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			var stmts []*sitter.Node
			switch n.Kind() {
			case "statement_block":
				stmts = rules.Named(n)
			case "switch_case", "switch_default":
				for i := uint(0); i < n.ChildCount(); i++ {
					ch := n.Child(i)
					if ch.IsNamed() && ch.Kind() != "comment" && n.FieldNameForChild(uint32(i)) != "value" {
						stmts = append(stmts, ch)
					}
				}
			default:
				return true
			}
			for i, s := range stmts {
				switch s.Kind() {
				case "return_statement", "throw_statement", "break_statement", "continue_statement":
				default:
					continue
				}
				for _, next := range stmts[i+1:] {
					if hoisted[next.Kind()] {
						continue
					}
					sym := symbol(f, n)
					c.Report(rules.Finding{File: f.Path, Line: rules.Line(next), Symbol: sym,
						Key:     sym + "|after:" + rules.Norm(f.Text(s)) + "|" + rules.Norm(f.Text(next)),
						Message: "unreachable: `" + rules.Norm(f.Text(next)) + "` follows `" + rules.Norm(f.Text(s)) + "`"})
					break
				}
				break
			}
			return true
		})
	}
}

func unusedPrivateFunction(c *rules.Context) {
	for _, f := range c.Files {
		if isDeclarationFile(f.Path) {
			continue
		}
		var uses map[string]int
		if isModule(f) {
			uses = countIdents(f, f.Root(), useKinds)
			for _, st := range rules.Named(f.Root()) {
				var names []*sitter.Node
				switch st.Kind() {
				case "function_declaration", "generator_function_declaration":
					names = append(names, st.ChildByFieldName("name"))
				case "lexical_declaration", "variable_declaration":
					for _, d := range rules.Named(st) {
						if d.Kind() == "variable_declarator" && functionValue(d.ChildByFieldName("value")) {
							if nm := d.ChildByFieldName("name"); nm != nil && nm.Kind() == "identifier" {
								names = append(names, nm)
							}
						}
					}
				}
				for _, nm := range names {
					if nm == nil || uses[f.Text(nm)] > 1 || strings.HasPrefix(f.Text(nm), "_") {
						continue
					}
					c.Report(rules.Finding{File: f.Path, Line: rules.Line(nm), Symbol: f.Text(nm), Key: "func:" + f.Text(nm),
						Message: "module-private function " + f.Text(nm) + " is never used"})
				}
			}
		}
		// private class methods
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			if n.Kind() != "class_body" {
				return true
			}
			props := countIdents(f, n, map[string]bool{"property_identifier": true, "private_property_identifier": true})
			for _, m := range rules.Named(n) {
				if m.Kind() != "method_definition" {
					continue
				}
				nm := m.ChildByFieldName("name")
				if nm == nil {
					continue
				}
				private := nm.Kind() == "private_property_identifier"
				for i := uint(0); i < m.ChildCount(); i++ {
					if ch := m.Child(i); ch.Kind() == "accessibility_modifier" && f.Text(ch) == "private" {
						private = true
					}
				}
				if !private || props[f.Text(nm)] > 1 {
					continue
				}
				sym := symbol(f, nm)
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(m), Symbol: sym, Key: "method:" + sym,
					Message: "private method " + sym + " is never used in its class"})
			}
			return true
		})
	}
}

// ---- lost errors ----

func emptyFn(n *sitter.Node) bool {
	if n == nil || !functionKinds[n.Kind()] {
		return false
	}
	b := n.ChildByFieldName("body")
	return b != nil && b.Kind() == "statement_block" && len(rules.Named(b)) == 0
}

func emptyHandler(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			switch n.Kind() {
			case "catch_clause":
				body := n.ChildByFieldName("body")
				if body == nil || len(rules.Named(body)) > 0 {
					return true
				}
				sev, msg := rules.Gate, "empty catch block discards the error"
				if rules.HasComment(body) {
					sev, msg = rules.Advisory, msg+" (the block holds only a comment)"
				}
				sym := symbol(f, n)
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Severity: sev, Key: sym + "|empty-catch", Message: msg})
			case "call_expression":
				fn := n.ChildByFieldName("function")
				if fn == nil || fn.Kind() != "member_expression" {
					return true
				}
				prop := f.Text(fn.ChildByFieldName("property"))
				args := rules.Named(n.ChildByFieldName("arguments"))
				var handler *sitter.Node
				switch {
				case prop == "catch" && len(args) >= 1:
					handler = args[0]
				case prop == "then" && len(args) >= 2:
					handler = args[1]
				}
				if emptyFn(handler) {
					sym := symbol(f, n)
					c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Key: sym + "|empty-" + prop + ":" + rules.Norm(f.Text(fn.ChildByFieldName("object"))),
						Message: "." + prop + "() with an empty rejection handler discards the error"})
				}
			}
			return true
		})
	}
}

func errorToSuccess(c *rules.Context) {
	for _, f := range c.Files {
		rules.Walk(f.Root(), func(n *sitter.Node) bool {
			switch n.Kind() {
			case "catch_clause":
				body := n.ChildByFieldName("body")
				stmts := rules.Named(body)
				if len(stmts) != 1 || stmts[0].Kind() != "return_statement" {
					return true
				}
				vals := rules.Named(stmts[0])
				if len(vals) == 1 && !constantValue(f, vals[0]) {
					return true
				}
				sym := symbol(f, n)
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Key: sym + "|catch:" + rules.Norm(f.Text(stmts[0])),
					Message: "catch block only does `" + rules.Norm(f.Text(stmts[0])) + "`: the failure becomes an ordinary value and the error is not used"})
			case "call_expression":
				fn := n.ChildByFieldName("function")
				if fn == nil || fn.Kind() != "member_expression" || f.Text(fn.ChildByFieldName("property")) != "catch" {
					return true
				}
				args := rules.Named(n.ChildByFieldName("arguments"))
				if len(args) == 0 || args[0].Kind() != "arrow_function" {
					return true
				}
				if body := args[0].ChildByFieldName("body"); body != nil && body.Kind() != "statement_block" && constantValue(f, body) {
					sym := symbol(f, n)
					c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Key: sym + "|catch-value:" + rules.Norm(f.Text(args[0])),
						Message: ".catch(" + rules.Norm(f.Text(args[0])) + ") turns a rejection into an ordinary value"})
				}
			}
			return true
		})
	}
}

// ---- test oracle ----

var expectMatchers = map[string]bool{"toBe": true, "toEqual": true, "toStrictEqual": true}
var chaiMatchers = map[string]bool{"equal": true, "equals": true, "eql": true}
var chaiChain = map[string]bool{"to": true, "be": true, "deep": true, "been": true, "is": true, "that": true, "which": true, "and": true,
	"has": true, "have": true, "with": true, "at": true, "of": true, "same": true, "strict": true}
var assertEq = map[string]bool{"equal": true, "strictEqual": true, "deepEqual": true, "deepStrictEqual": true, "equals": true, "eq": true}

// expectArg returns A for a chain rooted at expect(A), passing through the
// allowed chain properties, or nil (also when the chain contains .not,
// .resolves or .rejects).
func expectArg(f *rules.File, n *sitter.Node, through map[string]bool) *sitter.Node {
	for n != nil && n.Kind() == "member_expression" {
		if !through[f.Text(n.ChildByFieldName("property"))] {
			return nil
		}
		n = n.ChildByFieldName("object")
	}
	if n == nil || n.Kind() != "call_expression" {
		return nil
	}
	fn := n.ChildByFieldName("function")
	if fn == nil || fn.Kind() != "identifier" || f.Text(fn) != "expect" {
		return nil
	}
	args := rules.Named(n.ChildByFieldName("arguments"))
	if len(args) != 1 {
		return nil
	}
	return args[0]
}

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
				why = rules.Tautology(a, b, f.Source(), func(x *sitter.Node) bool { return isLiteral(f, x) }, "call_expression", "new_expression", "await_expression")
			}
			if why == "" {
				return true
			}
			text := rules.Norm(f.Text(n))
			c.At(f, n, symbol(f, n), "assert:"+text, "`"+text+"` "+why+" and cannot fail", "")
			return false
		})
	}
}

// assertion returns the compared operands of a recognised assertion call,
// or a reason when it asserts a constant.
func assertion(f *rules.File, n *sitter.Node) (a, b *sitter.Node, why string) {
	fn := n.ChildByFieldName("function")
	args := rules.Named(n.ChildByFieldName("arguments"))
	switch {
	case fn == nil:
	case fn.Kind() == "member_expression":
		return memberAssertion(f, f.Text(fn.ChildByFieldName("property")), fn.ChildByFieldName("object"), args)
	case fn.Kind() == "identifier" && f.Text(fn) == "assert" && len(args) >= 1 && args[0].Kind() == "true":
		return nil, nil, "asserts a constant"
	}
	return nil, nil, ""
}

// memberAssertion handles obj.prop(args): expect/chai matchers and
// node:assert methods.
func memberAssertion(f *rules.File, prop string, obj *sitter.Node, args []*sitter.Node) (a, b *sitter.Node, why string) {
	switch {
	case expectMatchers[prop] && len(args) == 1:
		return expectArg(f, obj, nil), args[0], ""
	case chaiMatchers[prop] && len(args) == 1:
		return expectArg(f, obj, chaiChain), args[0], ""
	case assertEq[prop] && len(args) >= 2 && obj.Kind() == "identifier" && f.Text(obj) == "assert":
		return args[0], args[1], ""
	case prop == "ok" && len(args) >= 1 && f.Text(obj) == "assert" && args[0].Kind() == "true":
		return nil, nil, "asserts a constant"
	case len(args) == 0:
		if x := expectArg(f, obj, nil); x != nil && constantMatcher(f, prop, x) {
			return nil, nil, "asserts a constant"
		}
	}
	return nil, nil, ""
}

// constantMatcher: expect(<constant>).toBeTruthy() and similar.
func constantMatcher(f *rules.File, prop string, x *sitter.Node) bool {
	switch prop {
	case "toBeTruthy":
		return x.Kind() == "true"
	case "toBeFalsy":
		return x.Kind() == "false"
	case "toBeNull":
		return x.Kind() == "null"
	case "toBeUndefined":
		return f.Text(x) == "undefined"
	case "toBeDefined":
		return isLiteral(f, x) && f.Text(x) != "undefined"
	}
	return false
}

// promiseAssertion reports whether an expression is a call chain through
// expect(...).resolves or .rejects.
func promiseAssertion(f *rules.File, n *sitter.Node) bool {
	found := false
	for n != nil {
		switch n.Kind() {
		case "call_expression":
			fn := n.ChildByFieldName("function")
			if fn != nil && fn.Kind() == "identifier" && f.Text(fn) == "expect" {
				return found
			}
			n = fn
		case "member_expression":
			if p := f.Text(n.ChildByFieldName("property")); p == "resolves" || p == "rejects" {
				found = true
			}
			n = n.ChildByFieldName("object")
		default:
			return false
		}
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
			if len(e) == 1 && e[0].Kind() == "call_expression" && promiseAssertion(f, e[0]) {
				sym := symbol(f, n)
				c.Report(rules.Finding{File: f.Path, Line: rules.Line(n), Symbol: sym, Key: sym + "|unawaited:" + rules.Norm(f.Text(e[0])),
					Message: "`" + rules.Norm(f.Text(e[0])) + "` is neither awaited nor returned: the test can finish before it settles and a failure is lost"})
			}
			return true
		})
	}
}
