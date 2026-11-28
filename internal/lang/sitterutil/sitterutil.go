// Package sitterutil holds Tree-sitter plumbing shared by language adapters:
// parsing, syntax-error collection and the complexity walk. Language
// differences enter through MeasureSpec.
package sitterutil

import (
	"fmt"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/cameron/code-evaluator/internal/model"
)

// Parse parses src with a grammar. The caller closes the tree.
func Parse(language *sitter.Language, src []byte) (*sitter.Tree, error) {
	parser := sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(language); err != nil {
		return nil, err
	}
	tree := parser.Parse(src, nil)
	if tree == nil {
		return nil, fmt.Errorf("tree-sitter returned no tree")
	}
	return tree, nil
}

// Line is a node's 1-based start line; EndLine its end line.
func Line(n *sitter.Node) int    { return int(n.StartPosition().Row) + 1 }
func EndLine(n *sitter.Node) int { return int(n.EndPosition().Row) + 1 }

// MaxErrors bounds the syntax errors recorded per file.
const MaxErrors = 20

// Errors collects ERROR and MISSING nodes as diagnostics.
type Errors struct {
	Language string
	Diags    []model.Diagnostic
	Lines    []int
}

// Collect walks n (only into subtrees that contain errors).
func (e *Errors) Collect(n *sitter.Node) {
	if len(e.Diags) >= MaxErrors {
		return
	}
	if n.IsError() || n.IsMissing() {
		msg := "syntax error"
		if n.IsMissing() {
			msg = "missing " + n.Kind()
		}
		e.Diags = append(e.Diags, model.Diagnostic{Severity: model.SevError, Code: "parse_error", Message: msg, Line: Line(n), Language: e.Language})
		e.Lines = append(e.Lines, Line(n))
		return
	}
	if !n.HasError() {
		return
	}
	for i := uint(0); i < n.ChildCount(); i++ {
		e.Collect(n.Child(i))
	}
}

// In reports whether any error line falls within [start, end].
func (e *Errors) In(start, end int) bool {
	for _, l := range e.Lines {
		if l >= start && l <= end {
			return true
		}
	}
	return false
}

// MeasureSpec describes a language's complexity rules.
type MeasureSpec struct {
	Decision  func(n *sitter.Node) bool // adds one decision point
	Statement func(kind string) bool    // counts as a statement
	Nests     func(n *sitter.Node) bool // opens a nesting level (else-if excluded by the language)
}

// Measure walks a callable body accumulating decisions, statements and the
// maximum nesting depth.
func Measure(n *sitter.Node, depth int, m *model.Metrics, spec *MeasureSpec) {
	if spec.Decision(n) {
		m.Decisions++
	}
	if spec.Statement(n.Kind()) {
		m.Statements++
	}
	childDepth := depth
	if spec.Nests(n) {
		childDepth = depth + 1
		if childDepth > m.MaxNesting {
			m.MaxNesting = childDepth
		}
	}
	for i := uint(0); i < n.ChildCount(); i++ {
		Measure(n.Child(i), childDepth, m, spec)
	}
}

// LogicalOp reports whether a binary expression's operator is one of ops.
func LogicalOp(n *sitter.Node, ops ...string) bool {
	if n.Kind() != "binary_expression" {
		return false
	}
	op := n.ChildByFieldName("operator")
	if op == nil {
		return false
	}
	for _, o := range ops {
		if op.Kind() == o {
			return true
		}
	}
	return false
}

// Walk visits n and its descendants; returning false skips a node's children.
func Walk(n *sitter.Node, visit func(*sitter.Node) bool) {
	if n == nil || !visit(n) {
		return
	}
	for i := uint(0); i < n.ChildCount(); i++ {
		Walk(n.Child(i), visit)
	}
}
