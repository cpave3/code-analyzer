# PRD: Codebase-Aware Workflow Evaluator

**Version:** v0 MVP  
**Status:** Draft  
**Working name:** CodeEvaluator

## 1. Overview

CodeEval evaluates the engineering quality of a code change in the context of its target repository.

The system is intended to help compare coding-agent workflows and determine which workflows produce better code. It evaluates the resulting patch rather than relying on the workflow's own claims, execution trace, or self-assessment.

The initial implementation will be a local CLI tool written in Go. It will analyze a Git diff, inspect the surrounding codebase, calculate a set of measurable quality metrics, and produce a structured evaluation report with scores, reasons, and evidence.

The MVP will prioritize deterministic analysis. An LLM-based evaluator may be introduced as an optional component for judgments that cannot be reliably made through static analysis alone.

## 2. Problem Statement

Coding-agent workflows can produce functionally correct changes that nevertheless degrade the codebase.

Examples include:

- Reimplementing functionality that already exists.
- Introducing unnecessary complexity.
- Creating poorly factored abstractions.
- Increasing coupling between modules.
- Making changes across an unnecessarily broad area of the repository.
- Failing to remove obsolete implementations after refactoring.

Existing test results and conventional CI checks do not adequately measure these outcomes.

We need a mechanism to evaluate the resulting code change against the repository's existing structure and conventions, producing evidence-backed measurements that can be compared across workflows.

## 3. Goals

The MVP must:

1. Evaluate a Git change in the context of the repository.
2. Identify identical and structurally similar code.
3. Measure cyclomatic complexity before and after the change.
4. Measure additional complexity and maintainability signals.
5. Detect changes to dependency relationships and architecture.
6. Identify changes that appear unnecessarily broad or unrelated.
7. Produce scores with clear explanations and source-code evidence.
8. Run locally without requiring a hosted service.
9. Produce machine-readable output suitable for CI and later workflow comparison.

## 4. Non-goals

The MVP will not:

- Automatically modify or refactor code.
- Guarantee semantic correctness.
- Replace unit tests, integration tests, or existing static analysis.
- Support every programming language equally.
- Build a complete semantic model of every repository.
- Automatically decide whether every abstraction is architecturally appropriate.
- Require an LLM for basic evaluation.
- Optimize prompts or agent instructions directly.
- Produce a single opaque overall code-quality score.

## 5. Target User and Workflow

The primary user is a developer experimenting with coding-agent workflows and orchestration strategies.

Typical workflow:

1. Run a coding-agent workflow against a repository.
2. Obtain a Git diff representing the resulting change.
3. Run CodeEval against the diff and its base revision.
4. Inspect the metrics, findings, and evidence.
5. Compare the evaluation against changes produced by other workflows.

The evaluator should be independent of the agent or workflow that produced the patch.

## 6. Technical Stack

| Component               | Technology                                 |
| ----------------------- | ------------------------------------------ |
| Main implementation     | Go                                         |
| CLI                     | Go                                         |
| Git integration         | Git CLI                                    |
| Syntax parsing          | Tree-sitter                                |
| Local index and history | SQLite                                     |
| Analyzer output         | JSON                                       |
| Initial interface       | CLI and JSON report                        |
| Optional judgment layer | LLM adapter, not required for core metrics |

Language-specific adapters may be added where Tree-sitter alone is insufficient for reliable symbol resolution or call-graph analysis.

## 7. Functional Requirements

### 7.1 Git Change Input

The evaluator must accept:

- Repository path.
- Base Git revision.
- Target Git revision or working-tree changes.
- Optional evaluation configuration.

The evaluator must:

- Generate the relevant diff.
- Identify added, modified, deleted, and moved files where possible.
- Map changed ranges to functions, classes, methods, and other supported symbols.
- Preserve source locations for findings.
- Support evaluating a patch against its base revision.

### 7.2 Repository Index

The evaluator must build a local repository index containing:

- Files and language identifiers.
- Symbol definitions.
- Symbol locations.
- Basic symbol relationships.
- Dependency edges where supported.
- Relevant source fingerprints.
- Existing metric measurements.

The index must support querying symbols and related code without repeatedly parsing the entire repository for every analyzer.

The initial implementation may rebuild the index when necessary. Incremental indexing is desirable but not required for the first MVP.

### 7.3 Duplication Analysis

The evaluator must detect:

**Exact duplication**

- Identical code blocks or functions.
- Repeated code introduced by the patch.

**Structural duplication**

- Similar syntax or control-flow structures with differences in names, literals, or formatting.

**Duplication delta**

- Duplication introduced by the patch.
- Duplication removed by the patch.
- Existing duplication modified by the patch.

**Reuse candidates**

- Existing functions or abstractions that appear to implement overlapping functionality.
- New implementations that appear to bypass existing abstractions.

Each finding must identify the relevant code locations and provide a similarity measurement or classification.

Semantic similarity is a candidate-finding mechanism, not proof that two implementations are equivalent.

### 7.4 Complexity Analysis

The evaluator must calculate before-and-after complexity measurements for changed symbols.

Required metrics:

- Cyclomatic complexity.
- Function or method length.
- Maximum nesting depth.
- Branch or decision-point count.

Cognitive complexity may be included where a reliable implementation is available.

The evaluator must report:

- Before value.
- After value.
- Delta.
- Relevant symbol location.
- Aggregate delta across changed symbols.

The system should distinguish a local complexity reduction from an increase elsewhere in the patch.

### 7.5 Dependency and Architecture Analysis

The evaluator must measure supported dependency changes, including:

- New and removed dependency edges.
- Fan-in and fan-out changes.
- Cycles introduced or removed.
- Dependency direction changes.
- Configured architectural boundary violations.

Architecture rules should be configurable and should not be hardcoded to one repository structure.

### 7.6 Change Scope Analysis

The evaluator must report:

- Files and symbols changed.
- Lines added, deleted, and moved where detectable.
- Public interface changes.
- Existing callers affected.
- Changes outside the direct dependency neighborhood of modified symbols, where detectable.

The MVP should identify potentially unrelated changes but must not automatically classify all broad changes as bad.

### 7.7 Refactoring and Abstraction Signals

The evaluator should identify candidate signals such as:

- Duplicate implementations consolidated into one.
- New helper functions or abstractions.
- Existing abstractions bypassed.
- New wrappers with limited apparent value.
- Obsolete code left behind after replacement.
- Changes that increase fragmentation across symbols.

These findings should be treated as heuristics and should include supporting evidence.

### 7.8 Evaluation Report

The evaluator must produce:

1. A human-readable summary.
2. A machine-readable JSON report.
3. Metric-level scores.
4. Before-and-after measurements.
5. Findings with explanations.
6. Evidence references to source locations.
7. Confidence values where analysis is uncertain.

The report must distinguish measured facts from heuristic interpretations.

## 8. Scoring Model

The MVP will use a directional score rather than a universal code-quality rating.

Suggested scale:

| Score | Meaning                      |
| ----: | ---------------------------- |
|    -2 | Significant degradation      |
|    -1 | Degradation                  |
|     0 | Neutral, unchanged, or mixed |
|    +1 | Improvement                  |
|    +2 | Significant improvement      |

Each metric should be scored independently.

The initial dimensions are:

- Duplication.
- Complexity.
- Architecture and coupling.
- Maintainability and change scope.
- Correctness confidence.

Raw measurements must always be retained alongside scores.

The MVP will not combine these dimensions into a single weighted score. A composite score may be considered after the metrics have been validated against real evaluations.

## 9. Analyzer Result Schema

Each analyzer must return a result using a common schema.

Example:

```json
{
  "metric": "cyclomatic_complexity",
  "score": -1,
  "confidence": 1.0,
  "summary": "Complexity increased in a modified function",
  "before": 4,
  "after": 7,
  "delta": 3,
  "rationale": ["Three additional decision paths were introduced"],
  "evidence": [
    {
      "file": "src/parser.go",
      "symbol": "parseExpression",
      "startLine": 42,
      "endLine": 81
    }
  ]
}
```

The schema must allow metrics to include multiple evidence locations and additional raw measurements.

## 10. Evaluation Pipeline

The MVP pipeline will be:

1. Resolve the base and target revisions.
2. Generate the Git diff.
3. Parse changed files.
4. Map diff ranges to symbols.
5. Load or build the repository index.
6. Run independent metric analyzers.
7. Aggregate analyzer results.
8. Generate the human-readable and JSON reports.

Analyzers should be independent and should not need to know which coding-agent workflow produced the patch.

## 11. CLI Interface

Initial command:

```bash
codeeval evaluate \
  --repo ./my-repository \
  --base main \
  --target HEAD \
  --output evaluation.json
```

The CLI should also support evaluating working-tree changes:

```bash
codeeval evaluate \
  --repo ./my-repository \
  --base HEAD \
  --working-tree
```

The exact CLI syntax may evolve during implementation.

## 12. Performance and Operational Requirements

- Must run locally without a hosted backend.
- Must not modify the repository being evaluated.
- Must support repositories of practical medium-to-large size.
- Must avoid reparsing unchanged files when an index can be reused.
- Must provide deterministic results for deterministic analyzers.
- Must report unsupported languages or analysis limitations explicitly.
- Must not silently omit failed analyzers.

The MVP should prioritize correctness and evidence quality over aggressive optimization.

## 13. MVP Acceptance Criteria

The MVP is considered complete when it can:

1. Evaluate a Git patch against a base revision.
2. Parse supported source files and identify changed symbols.
3. Detect exact and structural duplication candidates.
4. Calculate cyclomatic complexity before and after.
5. Report complexity and duplication deltas.
6. Measure basic dependency graph changes.
7. Report changed files, symbols, and scope.
8. Produce metric-level scores with explanations.
9. Link findings to precise source locations.
10. Emit a stable JSON report.
11. Run locally through a Go CLI.
12. Evaluate the same patch consistently regardless of which workflow produced it.

## 14. Implementation Phases

### Phase 1: Diff and Symbol Model

- Git integration.
- Tree-sitter parsing.
- Changed symbol identification.
- Source range tracking.

### Phase 2: Core Deterministic Metrics

- Exact duplication.
- Structural duplication.
- Cyclomatic complexity.
- Function length and nesting.
- Before-and-after reporting.

### Phase 3: Repository Context

- SQLite repository index.
- Symbol lookup.
- Basic references and dependency graph.
- Reuse candidate detection.

### Phase 4: Architecture and Scope

- Dependency graph deltas.
- Cycle detection.
- Configurable architecture rules.
- Change locality measurements.

### Phase 5: Evaluation Output

- Metric scoring.
- JSON schema.
- Human-readable report.
- Confidence and evidence handling.

### Phase 6: Optional Evaluator

- Constrained LLM interpretation of ambiguous findings.
- No unrestricted repository access by default.
- Evaluator must cite repository evidence.
- LLM judgments must remain separate from deterministic measurements.

## 15. Success Criteria

The MVP should demonstrate that it can distinguish patches that:

- Reduce duplication from patches that introduce it.
- Improve complexity from patches that increase it.
- Reuse existing abstractions from patches that reimplement them.
- Improve dependency structure from patches that increase coupling.
- Make focused changes from patches that introduce unnecessary scope.

The longer-term success criterion is that CodeEval can compare coding-agent workflows using repeatable, repository-aware measurements, allowing workflow changes to be evaluated empirically rather than by subjective review alone.
