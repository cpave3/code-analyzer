// Package correctness reports the correctness-confidence dimension as not
// assessed. Static structure cannot establish behavioural correctness, and
// the evaluator does not execute the evaluated repository's tests. It records
// test-file churn as context only, explicitly not as evidence.
package correctness

import (
	"fmt"

	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/model"
)

// Analyzer is the correctness-confidence analyzer.
type Analyzer struct{}

// New returns the analyzer.
func New() *Analyzer { return &Analyzer{} }

func (Analyzer) Descriptor() analyzer.Descriptor {
	return analyzer.Descriptor{
		ID: "correctness", Dimension: analyzer.DimCorrectness, Title: "Correctness confidence", Version: "1",
		Description: "Not assessed: requires independent evidence (e.g. trusted test results supplied by the caller). Test churn is reported as context only.",
		Classes:     []model.FileClass{model.ClassTest},
	}
}

func (Analyzer) Analyze(ctx *analyzer.Context) (*analyzer.Result, error) {
	files, added, deleted := 0, 0, 0
	for _, fc := range ctx.Change.Files {
		if fc.Class == model.ClassTest {
			files++
			added += fc.Added
			deleted += fc.Deleted
		}
	}
	return &analyzer.Result{
		Status:    analyzer.StatusNotAssessed,
		ScoreRule: "never scored in this POC: no independent correctness evidence is collected",
		Summary:   fmt.Sprintf("not assessed (context: %d test files changed, +%d/-%d lines; test churn is not evidence of correctness)", files, added, deleted),
		Measurements: []analyzer.Measurement{
			analyzer.Count("test_files_changed", "context_only", files, "files"),
			analyzer.Count("test_lines_added", "context_only", added, "lines"),
			analyzer.Count("test_lines_deleted", "context_only", deleted, "lines"),
		},
	}, nil
}
