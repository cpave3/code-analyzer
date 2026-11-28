// Package config holds the versioned evaluation configuration: exclusions,
// analyzer thresholds, directional scoring thresholds and architecture rules.
// Configuration lives outside the evaluated repository and is passed with
// --config; the evaluator never reads configuration from the target repo.
package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

// Version of the configuration format.
const Version = 1

// ScoringVersion identifies the scoring rules; bump when thresholds or rule
// semantics change so reports stay comparable.
const ScoringVersion = "scoring/v1"

// Config is the evaluation configuration.
type Config struct {
	ConfigVersion int `json:"configVersion"`
	// Exclude globs (repository-relative, ** supported) classify files as
	// "excluded"; they are counted but not analysed.
	Exclude      []string     `json:"exclude"`
	Duplication  Duplication  `json:"duplication"`
	Complexity   Complexity   `json:"complexity"`
	Architecture Architecture `json:"architecture"`
	Scope        Scope        `json:"scope"`
	// Defects configures `codeeval defects` (internal/rules). Omitted when
	// empty so existing configuration hashes are unchanged.
	Defects Defects `json:"defects,omitzero"`
}

// Defects configures the defect rules.
type Defects struct {
	// Disable rule IDs or category prefixes (e.g. "dead-code/unused-export").
	Disable []string `json:"disable,omitempty"`
}

// Duplication thresholds.
type Duplication struct {
	MinTokens     int     `json:"minTokens"`     // smallest callable considered for clones
	NearThreshold float64 `json:"nearThreshold"` // MinHash similarity for near-structural candidates
	IncludeTests  bool    `json:"includeTests"`  // analyse test files as clone sources
	// Scoring: introduced redundant instances (exact+structural, net of removed)
	DegradeAt            int `json:"degradeAt"`
	SignificantDegradeAt int `json:"significantDegradeAt"`
}

// Complexity thresholds.
type Complexity struct {
	HotspotCyclomatic int `json:"hotspotCyclomatic"` // a callable at/above this is a hotspot
	HotspotNesting    int `json:"hotspotNesting"`
	// Scoring on modified (pre-existing) callables' net cyclomatic delta.
	DegradeAt            int `json:"degradeAt"`
	SignificantDegradeAt int `json:"significantDegradeAt"`
}

// Architecture rules.
type Architecture struct {
	Rules []Rule `json:"rules"`
}

// Rule restricts dependencies of units matching From. With To, it forbids
// dependencies on units matching To (enforced by v1, v2 and `codeeval
// defects`). With Only, units matching From may depend only on units matching
// one of the Only globs or From itself (an allow-list, enforced by `codeeval
// defects` only).
type Rule struct {
	Name string   `json:"name"`
	From string   `json:"from"`
	To   string   `json:"to"`
	Only []string `json:"only,omitempty"`
}

// Scope thresholds.
type Scope struct {
	// FormattingOnlyChurnAt: number of symbols changed only in layout/comments
	// that triggers a -1 when unrelated to the rest of the change.
	FormattingOnlyChurnAt int `json:"formattingOnlyChurnAt"`
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		ConfigVersion: Version,
		Exclude:       []string{},
		Duplication:   Duplication{MinTokens: 40, NearThreshold: 0.85, DegradeAt: 1, SignificantDegradeAt: 3},
		Complexity:    Complexity{HotspotCyclomatic: 10, HotspotNesting: 4, DegradeAt: 3, SignificantDegradeAt: 10},
		Architecture:  Architecture{Rules: []Rule{}},
		Scope:         Scope{FormattingOnlyChurnAt: 3},
	}
}

// Load reads a config file over the defaults.
func Load(path string) (*Config, error) {
	c := Default()
	if path == "" {
		return c, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(c); err != nil {
		return nil, fmt.Errorf("config %s: %w", path, err)
	}
	if c.ConfigVersion != Version {
		return nil, fmt.Errorf("config %s: configVersion %d unsupported (want %d)", path, c.ConfigVersion, Version)
	}
	for _, r := range c.Architecture.Rules {
		if r.From == "" || (r.To == "" && len(r.Only) == 0) {
			return nil, fmt.Errorf("config %s: architecture rule %q needs from and either to or only", path, r.Name)
		}
	}
	for _, g := range append(append([]string{}, c.Exclude...), ruleGlobs(c)...) {
		if _, err := compile(g); err != nil {
			return nil, fmt.Errorf("config %s: bad glob %q: %w", path, g, err)
		}
	}
	return c, nil
}

func ruleGlobs(c *Config) []string {
	var out []string
	for _, r := range c.Architecture.Rules {
		out = append(out, r.From, r.To)
		out = append(out, r.Only...)
	}
	return out
}

// Hash returns a stable hash of the effective configuration.
func (c *Config) Hash() string {
	data, _ := json.Marshal(c)
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:8])
}

var globCache sync.Map // glob -> *regexp.Regexp

func compile(glob string) (*regexp.Regexp, error) {
	if re, ok := globCache.Load(glob); ok {
		return re.(*regexp.Regexp), nil
	}
	var b strings.Builder
	b.WriteString("^")
	for i := 0; i < len(glob); i++ {
		c := glob[i]
		switch c {
		case '*':
			if i+1 < len(glob) && glob[i+1] == '*' {
				i++
				if i+1 < len(glob) && glob[i+1] == '/' {
					i++
					b.WriteString("(?:.*/)?")
				} else {
					b.WriteString(".*")
				}
			} else {
				b.WriteString("[^/]*")
			}
		case '?':
			b.WriteString("[^/]")
		default:
			b.WriteString(regexp.QuoteMeta(string(c)))
		}
	}
	b.WriteString("$")
	re, err := regexp.Compile(b.String())
	if err == nil {
		globCache.Store(glob, re)
	}
	return re, err
}

// Match reports whether a repository-relative path matches a glob.
func Match(glob, p string) bool {
	// "dir/**" also matches "dir" itself (graph units are directories).
	if strings.HasSuffix(glob, "/**") && p == strings.TrimSuffix(glob, "/**") {
		return true
	}
	re, err := compile(glob)
	return err == nil && re.MatchString(p)
}
