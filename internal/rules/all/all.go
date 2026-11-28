// Package all assembles the default defect-rule registry. This is the one
// place rule packages are wired in.
package all

import (
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
	"github.com/cameron/code-evaluator/internal/rules/generic"
	"github.com/cameron/code-evaluator/internal/rules/golang"
	"github.com/cameron/code-evaluator/internal/rules/typescript"
)

// Default returns the catalog with every language implementation.
func Default(langs *lang.Registry) *rules.Registry {
	r := rules.NewRegistry()
	must(catalog.Define(r))
	must(golang.Register(r))
	must(typescript.Register(r))
	must(generic.Register(r, langs))
	return r
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
