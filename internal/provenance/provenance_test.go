package provenance

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	msg := "Agent: subject is not metadata\n\nBody text.\n\nAgent-Model: claude-opus-5-5\nAgent-Harness: claude-code/2.1\nAgent-Workflow: proposal-stack@2\nagent-effort: high\nCo-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>\n"
	r := Parse("c1", msg, "Agent-Workflow: from-note\nAgent-Session: s1", true)
	want := Record{Commit: "c1", Model: []string{"claude-opus-5-5"}, Harness: []string{"claude-code/2.1"},
		Meta: map[string]string{"workflow": "proposal-stack@2", "effort": "high", "session": "s1"}, Source: SourceBoth}
	if !reflect.DeepEqual(r, want) {
		t.Fatalf("got %+v\nwant %+v", r, want)
	}
	if k := r.Key([]string{"model", "workflow", "missing"}); k != "model=claude-opus-5-5, workflow=proposal-stack@2, missing=(none)" {
		t.Errorf("key %q", k)
	}
}

func TestSquashedMessageAndNotesOnly(t *testing.T) {
	squash := "Feature (#12)\n\n* first\n\nAgent: model-a\nAgent-Workflow: v1\n\n* second\n\nAgent: model-b\nAgent-Workflow: v1\n"
	r := Parse("c", squash, "", true)
	if r.Value("model") != "model-a+model-b" || r.Value("workflow") != "v1" || r.Source != SourceMessage {
		t.Errorf("squash: %+v", r)
	}
	n := Parse("c", "Fix\n", "Agent: model-c\nAgent-Harness: h", true)
	if n.Value("agent") != "model-c" || n.Value("harness") != "h" || n.Source != SourceNote {
		t.Errorf("note: %+v", n)
	}
}

func TestInference(t *testing.T) {
	msg := "Fix\n\nCo-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>\nCo-authored-by: Jane Doe <jane@example.com>\n"
	if r := Parse("c", msg, "", true); r.Value("model") != "claude-opus-5-5" || r.Source != SourceInferred || len(r.Harness) != 0 {
		t.Errorf("inferred: %+v", r)
	}
	if r := Parse("c", msg, "", false); r.Source != SourceNone || len(r.Model) != 0 {
		t.Errorf("inference disabled: %+v", r)
	}
	if r := Parse("c", "Fix\n\nCo-authored-by: Jane Doe <jane@example.com>\n", "", true); r.Source != SourceNone {
		t.Errorf("human co-author inferred as agent: %+v", r)
	}
}

func TestProseIsNotProvenance(t *testing.T) {
	// From this repository's history: a wrapped sentence put "Agent:" at
	// the start of a body line.
	msg := "Read agent provenance\n\nAgent: and Agent-Harness: are standard fields; any Agent-<Key>: line is\nfree-form metadata (e.g. Agent-Workflow).\n\nCo-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>\n"
	r := Parse("c", msg, "Agent-Model: m\nAgent-Harness: h", false)
	if r.Value("model") != "m" || r.Source != SourceNote {
		t.Errorf("prose parsed as provenance: %+v", r)
	}
	mixed := "Fix\n\nSome prose.\n\nAgent-Model: x\nAgent-Harness: y\nCo-Authored-By: A <a@b>\n"
	if r := Parse("c", mixed, "", false); r.Value("model") != "x" || r.Value("harness") != "y" {
		t.Errorf("trailer block next to other trailers: %+v", r)
	}
}
