// Package languages assembles the default language registry. This is the one
// place a new adapter is wired in.
package languages

import (
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/lang/golang"
)

// Default returns the registry used by the CLI: the Go adapter plus declared,
// explicitly unsupported languages so their files are counted as coverage
// gaps rather than ignored.
func Default() *lang.Registry {
	r := lang.NewRegistry()
	must(r.Register(golang.New()))

	// Planned next adapters (Tree-sitter grammars exist for both). Declared
	// only: files are reported as unsupported, never scored.
	r.Declare(lang.Declared{ID: "typescript", Extensions: []string{".ts", ".tsx", ".mts", ".cts"}, Note: "planned adapter: tree-sitter-typescript; graph from import/export + tsconfig paths"})
	r.Declare(lang.Declared{ID: "csharp", Extensions: []string{".cs"}, Note: "planned adapter: tree-sitter-c-sharp; graph from using/namespace + .csproj ProjectReference"})

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
