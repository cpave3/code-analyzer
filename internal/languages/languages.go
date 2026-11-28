// Package languages assembles the default language registry. This is the one
// place a new adapter is wired in.
package languages

import (
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/lang/csharp"
	"github.com/cameron/code-evaluator/internal/lang/golang"
	"github.com/cameron/code-evaluator/internal/lang/typescript"
)

// Default returns the registry used by the CLI: the Go and TypeScript adapters plus declared,
// explicitly unsupported languages so their files are counted as coverage
// gaps rather than ignored.
func Default() *lang.Registry {
	r := lang.NewRegistry()
	must(r.Register(golang.New()))
	must(r.Register(typescript.New()))
	must(r.Register(csharp.New()))

	// Other recognised languages, reported as unsupported coverage.
	r.Declare(lang.Declared{ID: "javascript", Extensions: []string{".js", ".jsx", ".mjs", ".cjs"}, Note: "no adapter"})
	r.Declare(lang.Declared{ID: "kotlin", Extensions: []string{".kt", ".kts"}, Note: "no adapter"})
	r.Declare(lang.Declared{ID: "rust", Extensions: []string{".rs"}, Note: "no adapter"})
	r.Declare(lang.Declared{ID: "python", Extensions: []string{".py"}, Note: "no adapter"})
	r.Declare(lang.Declared{ID: "swift", Extensions: []string{".swift"}, Note: "no adapter"})
	r.Declare(lang.Declared{ID: "java", Extensions: []string{".java"}, Note: "no adapter"})
	r.Declare(lang.Declared{ID: "sql", Extensions: []string{".sql"}, Note: "no adapter"})
	r.Declare(lang.Declared{ID: "shell", Extensions: []string{".sh", ".bash"}, Note: "no adapter"})
	return r
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
