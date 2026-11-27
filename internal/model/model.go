// Package model holds the language-neutral data model shared by language
// adapters, the repository index and analyzers. Nothing here knows about a
// specific language or metric.
package model

// FileClass says what role a file plays in the repository. Analyzers use it to
// keep tests, generated code, vendored code and documentation apart from
// production source.
type FileClass string

const (
	ClassSource    FileClass = "source"
	ClassTest      FileClass = "test"
	ClassTestdata  FileClass = "testdata"
	ClassGenerated FileClass = "generated"
	ClassVendor    FileClass = "vendor"
	ClassBuild     FileClass = "build"
	ClassDocs      FileClass = "docs"
	ClassExcluded  FileClass = "excluded" // matched a configured exclude glob
	ClassOther     FileClass = "other"
)

// Severity levels for diagnostics.
const (
	SevInfo    = "info"
	SevWarning = "warning"
	SevError   = "error"
)

// Diagnostic records an analysis limitation, unsupported input or failure. It
// is always surfaced in the report; nothing is silently dropped.
type Diagnostic struct {
	Severity string `json:"severity"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Side     string `json:"side,omitempty"` // base | target
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Analyzer string `json:"analyzer,omitempty"`
	Language string `json:"language,omitempty"`
}

// Metrics are per-callable measurements produced by a language adapter.
type Metrics struct {
	Cyclomatic int `json:"cyclomatic"`
	Decisions  int `json:"decisions"`
	MaxNesting int `json:"maxNesting"`
	Lines      int `json:"lines"`
	Statements int `json:"statements"`
}

// Symbol kinds understood across languages.
const (
	KindFunction = "function"
	KindMethod   = "method"
	KindType     = "type"
)

// Symbol is a declaration produced by an adapter. Path-dependent fields (File,
// Package, Key) are filled in by the index so that parse results can be cached
// by content alone.
type Symbol struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Container string `json:"container,omitempty"` // receiver / enclosing type
	Exported  bool   `json:"exported"`
	Signature string `json:"signature,omitempty"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`

	Metrics *Metrics `json:"metrics,omitempty"` // nil for non-callables

	// ExactFP hashes the declaration body token text (comments and layout
	// ignored, symbol name excluded). StructFP hashes the token kinds with
	// identifiers and literals abstracted. MinHash is a structural shingle
	// sketch used to find near-structural candidates.
	ExactFP  string   `json:"exactFP,omitempty"`
	StructFP string   `json:"structFP,omitempty"`
	Tokens   int      `json:"tokens"`
	MinHash  []uint64 `json:"minHash,omitempty"`

	InErrorRegion bool `json:"inErrorRegion,omitempty"`

	// Filled by the index.
	File    string `json:"-"`
	Package string `json:"-"`
	Key     string `json:"-"`
}

// DisplayName is the human name, e.g. "(*Server).Start" style is avoided in
// favour of "Server.Start".
func (s *Symbol) DisplayName() string {
	if s.Container != "" {
		return s.Container + "." + s.Name
	}
	return s.Name
}

// Callable reports whether the symbol has complexity metrics.
func (s *Symbol) Callable() bool { return s.Metrics != nil }

// Import is a dependency declaration in a source file.
type Import struct {
	Path string `json:"path"`
	Name string `json:"name,omitempty"` // explicit alias, if any
	Line int    `json:"line"`
}

// ParsedFile is the path-independent parse result of one file. It is cached by
// (adapter, adapter version, content hash).
type ParsedFile struct {
	PackageName string         `json:"packageName,omitempty"`
	Imports     []Import       `json:"imports,omitempty"`
	Symbols     []Symbol       `json:"symbols,omitempty"`
	Refs        map[string]int `json:"refs,omitempty"` // name-based identifier uses, "q.Name" for qualified
	Generated   bool           `json:"generated,omitempty"`
	ParseErrors []Diagnostic   `json:"parseErrors,omitempty"`
	Lines       int            `json:"lines"`
}

// Location is a file position used for evidence of graph edges.
type Location struct {
	File string `json:"file"`
	Line int    `json:"line"`
}
