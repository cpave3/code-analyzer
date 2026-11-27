// Package analyzers assembles the default analyzer registry. This is the one
// place a new analyzer is wired in.
package analyzers

import (
	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/analyzers/complexity"
	"github.com/cameron/code-evaluator/internal/analyzers/correctness"
	"github.com/cameron/code-evaluator/internal/analyzers/dependencies"
	"github.com/cameron/code-evaluator/internal/analyzers/duplication"
	"github.com/cameron/code-evaluator/internal/analyzers/obsolete"
	"github.com/cameron/code-evaluator/internal/analyzers/scope"
)

// Default returns the registry used by the CLI.
func Default() *analyzer.Registry {
	r := analyzer.NewRegistry()
	for _, a := range []analyzer.Analyzer{
		duplication.New(), complexity.New(), dependencies.New(), scope.New(), obsolete.New(), correctness.New(),
	} {
		if err := r.Register(a); err != nil {
			panic(err)
		}
	}
	return r
}
