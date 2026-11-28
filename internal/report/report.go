// Package report defines the versioned JSON report and renders the
// human-readable Markdown summary.
package report

import (
	"github.com/cameron/code-evaluator/internal/analyzer"
	"github.com/cameron/code-evaluator/internal/change"
	"github.com/cameron/code-evaluator/internal/model"
)

// Schema identifies the report format.
const Schema = "codeeval.report/v1"

// HistorySchema identifies the history (series) format.
const HistorySchema = "codeeval.history/v1"

// ToolVersion is the evaluator version.
const ToolVersion = "0.4.0-poc"

// Tool identifies the producer.
type Tool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Revision pins one side of the evaluation.
type Revision struct {
	Ref        string `json:"ref"`
	Commit     string `json:"commit"`
	Tree       string `json:"tree"`
	CommitTime string `json:"commitTime,omitempty"`
}

// LanguageCoverage summarises one language across the target snapshot and
// the change.
type LanguageCoverage struct {
	ID           string            `json:"id"`
	Supported    bool              `json:"supported"`
	AdapterVer   string            `json:"adapterVersion,omitempty"`
	Capabilities map[string]string `json:"capabilities,omitempty"`
	Limits       []string          `json:"limits,omitempty"`
	Note         string            `json:"note,omitempty"`
	FilesTotal   int               `json:"filesInTarget"`
	FilesByClass map[string]int    `json:"filesInTargetByClass"`
	FilesChanged int               `json:"filesChanged"`
	LinesChanged int               `json:"linesChanged"`
}

// SymbolRow is a changed symbol in the report.
type SymbolRow struct {
	Status           string `json:"status"`
	Kind             string `json:"kind"`
	Name             string `json:"name"`
	Exported         bool   `json:"exported"`
	File             string `json:"file,omitempty"`
	StartLine        int    `json:"startLine,omitempty"`
	EndLine          int    `json:"endLine,omitempty"`
	BaseName         string `json:"baseName,omitempty"`
	BaseFile         string `json:"baseFile,omitempty"`
	BaseStartLine    int    `json:"baseStartLine,omitempty"`
	BaseEndLine      int    `json:"baseEndLine,omitempty"`
	SignatureChanged bool   `json:"signatureChanged,omitempty"`
	Class            string `json:"class"`
}

// ChangeSummary lists changed files and symbols.
type ChangeSummary struct {
	FilesChanged int                 `json:"filesChanged"`
	LinesAdded   int                 `json:"linesAdded"`
	LinesDeleted int                 `json:"linesDeleted"`
	Files        []change.FileChange `json:"files"`
	Symbols      []SymbolRow         `json:"symbols"`
}

// AnalyzerScore is one analyzer's headline inside a dimension.
type AnalyzerScore struct {
	Analyzer string `json:"analyzer"`
	Status   string `json:"status"`
	Score    *int   `json:"score"`
	Summary  string `json:"summary"`
}

// Dimension groups analyzer scores. There is deliberately no combined score.
type Dimension struct {
	Dimension string          `json:"dimension"`
	Analyzers []AnalyzerScore `json:"analyzers"`
}

// Diagnostics holds diagnostics relevant to the change plus snapshot-wide
// counts.
type Diagnostics struct {
	Change              []model.Diagnostic `json:"change"`
	SnapshotCountByCode map[string]int     `json:"snapshotCountByCode"`
}

// Timing is optional run metadata (omitted by default for determinism).
type Timing struct {
	StartedAt     string `json:"startedAt"`
	ElapsedMillis int64  `json:"elapsedMillis"`
	FilesParsed   int    `json:"filesParsed"`
	CacheHits     int    `json:"cacheHits"`
	CacheMisses   int    `json:"cacheMisses"`
}

// Report is the evaluation report.
type Report struct {
	Schema         string             `json:"schema"`
	Tool           Tool               `json:"tool"`
	ScoringVersion string             `json:"scoringVersion"`
	ConfigVersion  int                `json:"configVersion"`
	ConfigHash     string             `json:"configHash"`
	Mode           string             `json:"mode"` // commits | worktree
	Repository     string             `json:"repository"`
	Base           Revision           `json:"base"`
	Target         Revision           `json:"target"`
	ScoreScale     string             `json:"scoreScale"`
	Languages      []LanguageCoverage `json:"languages"`
	Change         ChangeSummary      `json:"change"`
	Dimensions     []Dimension        `json:"dimensions"`
	Results        []*analyzer.Result `json:"results"`
	Diagnostics    Diagnostics        `json:"diagnostics"`
	Limitations    []string           `json:"limitations"`
	Timing         *Timing            `json:"timing,omitempty"`
}

// HistoryStep is one first-parent step.
type HistoryStep struct {
	Index   int     `json:"index"`
	Commit  string  `json:"commit"`
	Parent  string  `json:"parent"`
	Subject string  `json:"subject"` // untrusted commit text, data only
	Time    string  `json:"commitTime"`
	Report  *Report `json:"report"`
}

// History is a first-parent series plus the net evaluation.
type History struct {
	Schema string        `json:"schema"`
	Tool   Tool          `json:"tool"`
	Base   Revision      `json:"base"`
	Target Revision      `json:"target"`
	Net    *Report       `json:"net"`
	Steps  []HistoryStep `json:"steps"`
	Notes  []string      `json:"notes"`
	Timing *Timing       `json:"timing,omitempty"`
}
