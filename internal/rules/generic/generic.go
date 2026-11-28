// Package generic holds rule implementations that work for any language
// whose adapter builds a dependency graph.
package generic

import (
	"fmt"
	"sort"

	"github.com/cameron/code-evaluator/internal/config"
	"github.com/cameron/code-evaluator/internal/lang"
	"github.com/cameron/code-evaluator/internal/rules"
	"github.com/cameron/code-evaluator/internal/rules/catalog"
)

// Register implements the graph-based rules for each adapter with a
// dependency graph.
func Register(r *rules.Registry, langs *lang.Registry) error {
	for _, ad := range langs.Adapters() {
		if !lang.Has(ad, lang.CapDependencies) {
			continue
		}
		if err := r.Implement(rules.Impl{Rule: catalog.ImportPolicy, Language: ad.ID(), Check: importPolicy,
			Note: "graph edges between " + ad.ID() + " units (" + ad.Capabilities()[lang.CapDependencies] + ")"}); err != nil {
			return err
		}
	}
	return nil
}

// violation returns the reason an edge breaks a rule, or "".
func violation(rule config.Rule, from, to string) string {
	if !config.Match(rule.From, from) {
		return ""
	}
	if rule.To != "" && config.Match(rule.To, to) {
		return fmt.Sprintf("forbidden by rule %q (%s must not depend on %s)", rule.Name, rule.From, rule.To)
	}
	if len(rule.Only) > 0 && !config.Match(rule.From, to) {
		for _, g := range rule.Only {
			if config.Match(g, to) {
				return ""
			}
		}
		return fmt.Sprintf("outside the allow-list of rule %q (%s may depend only on %v)", rule.Name, rule.From, rule.Only)
	}
	return ""
}

func importPolicy(c *rules.Context) {
	if c.Graph == nil || c.Config == nil || len(c.Config.Architecture.Rules) == 0 {
		return
	}
	froms := make([]string, 0, len(c.Graph.Edges))
	for f := range c.Graph.Edges {
		froms = append(froms, f)
	}
	sort.Strings(froms)
	for _, from := range froms {
		tos := make([]string, 0, len(c.Graph.Edges[from]))
		for t := range c.Graph.Edges[from] {
			tos = append(tos, t)
		}
		sort.Strings(tos)
		for _, to := range tos {
			for _, rule := range c.Config.Architecture.Rules {
				why := violation(rule, from, to)
				if why == "" {
					continue
				}
				for _, loc := range c.Graph.Edges[from][to] {
					c.Report(rules.Finding{File: loc.File, Line: loc.Line,
						Message: fmt.Sprintf("%s imports %s: %s", from, to, why),
						Key:     "rule:" + rule.Name + "|" + from + "->" + to})
				}
			}
		}
	}
}
