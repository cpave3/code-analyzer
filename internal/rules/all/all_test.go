package all

import (
	"testing"

	"github.com/cameron/code-evaluator/internal/languages"
)

// The initial rule set must be implemented for every supported language.
// A later rule may cover a subset; it then needs an entry here with the
// reason, so a gap is a decision and not an accident.
var knownGaps = map[string]map[string]string{}

func TestEveryRuleImplementedForEveryLanguage(t *testing.T) {
	langs := languages.Default()
	reg := Default(langs)
	if len(reg.Rules()) == 0 {
		t.Fatal("no rules")
	}
	for _, rule := range reg.Rules() {
		for _, ad := range langs.Adapters() {
			if _, ok := reg.Impl(rule.ID, ad.ID()); !ok && knownGaps[rule.ID][ad.ID()] == "" {
				t.Errorf("%s has no %s implementation", rule.ID, ad.ID())
			}
		}
	}
}
