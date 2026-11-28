// Package provenance reads which agent, harness and workflow produced a
// commit, so quality deltas can be attributed to them.
//
// Carrier (in priority order):
//
//  1. Message lines (normally the trailer block):
//     Agent: claude-opus-5-5              standard: the model or agent
//     Agent-Harness: claude-code/2.1      standard: the tool that ran it
//     Agent-<Key>: <value>                any other metadata, e.g.
//     Agent-Workflow: proposal-stack@2
//     Lines are read anywhere in the body (not the subject), so trailers
//     survive squash merges that concatenate commit messages.
//  2. A note on the commit in refs/notes/agent-provenance with the same
//     lines, for commits that are already made. Message values win per key.
//  3. Inference from Co-Authored-By lines that name a known coding agent
//     (Claude, Codex, Copilot, Gemini, ...), labelled "co-author (inferred)":
//     agent only, no harness or metadata.
//
// All of it is self-reported, untrusted text: provenance says what a commit
// claims, not what produced it.
package provenance

import (
	"regexp"
	"sort"
	"strings"
)

// NotesRef holds provenance notes.
const NotesRef = "refs/notes/agent-provenance"

// Sources.
const (
	SourceMessage  = "message"
	SourceNote     = "note"
	SourceBoth     = "message+note"
	SourceInferred = "co-author (inferred)"
	SourceNone     = "none"
)

// Standard field names.
const (
	FieldAgent   = "agent"
	FieldHarness = "harness"
)

// Record is the provenance of one commit.
type Record struct {
	Commit  string            `json:"commit"`
	Agent   []string          `json:"agent,omitempty"`
	Harness []string          `json:"harness,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"` // Agent-<Key> → value; keys lower-case
	Source  string            `json:"source"`
	Merge   bool              `json:"merge,omitempty"` // a merge commit (set by the caller)
}

var lineRe = regexp.MustCompile(`(?i)^\s*(?:[-*]\s+)?agent(?:-([a-z0-9][a-z0-9_-]*))?\s*:\s*(\S.*?)\s*$`)
var coAuthorRe = regexp.MustCompile(`(?i)^\s*co-authored-by\s*:\s*([^<]+?)\s*<`)
var knownAgent = regexp.MustCompile(`(?i)\b(claude|codex|copilot|gemini|cursor|devin|aider|gpt|openhands|jules|amp)\b`)

type fields struct {
	agent, harness []string
	meta           map[string]string
}

func (f *fields) empty() bool { return len(f.agent) == 0 && len(f.harness) == 0 && len(f.meta) == 0 }

func addUnique(list []string, v string) []string {
	for _, x := range list {
		if x == v {
			return list
		}
	}
	return append(list, v)
}

func parseLines(text string, skipFirst bool) fields {
	f := fields{meta: map[string]string{}}
	for i, line := range strings.Split(text, "\n") {
		if i == 0 && skipFirst {
			continue // the subject is prose, not metadata
		}
		m := lineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		key, val := strings.ToLower(m[1]), m[2]
		switch key {
		case "":
			f.agent = addUnique(f.agent, val)
		case "harness":
			f.harness = addUnique(f.harness, val)
		default:
			key = strings.ReplaceAll(key, "_", "-")
			if old, ok := f.meta[key]; ok && old != val && !strings.Contains(","+old+",", ","+val+",") {
				val = old + "," + val
			}
			f.meta[key] = val
		}
	}
	return f
}

// Parse reads provenance from a commit message and an optional note.
func Parse(commit, message, note string, infer bool) Record {
	r := Record{Commit: commit, Source: SourceNone}
	msg := parseLines(message, true)
	nt := parseLines(note, false)
	switch {
	case !msg.empty() && !nt.empty():
		r.Source = SourceBoth
	case !msg.empty():
		r.Source = SourceMessage
	case !nt.empty():
		r.Source = SourceNote
	}
	pick := func(a, b []string) []string {
		if len(a) > 0 {
			return a
		}
		return b
	}
	r.Agent, r.Harness = pick(msg.agent, nt.agent), pick(msg.harness, nt.harness)
	if len(msg.meta)+len(nt.meta) > 0 {
		r.Meta = map[string]string{}
		for k, v := range nt.meta {
			r.Meta[k] = v
		}
		for k, v := range msg.meta {
			r.Meta[k] = v
		}
	}
	if r.Source == SourceNone && infer {
		for _, line := range strings.Split(message, "\n") {
			if m := coAuthorRe.FindStringSubmatch(line); m != nil && knownAgent.MatchString(m[1]) {
				r.Agent = addUnique(r.Agent, strings.ToLower(strings.Join(strings.Fields(m[1]), "-")))
			}
		}
		if len(r.Agent) > 0 {
			r.Source = SourceInferred
		}
	}
	sort.Strings(r.Agent)
	sort.Strings(r.Harness)
	return r
}

// Value returns a field for grouping: "agent", "harness" or a metadata key.
// Multiple values are joined with "+"; "" when absent.
func (r Record) Value(field string) string {
	switch field {
	case FieldAgent:
		return strings.Join(r.Agent, "+")
	case FieldHarness:
		return strings.Join(r.Harness, "+")
	}
	return r.Meta[strings.ToLower(field)]
}

// None is the group value for a missing field.
const None = "(none)"

// Mixed is the group key of a step whose commits disagree.
const Mixed = "(mixed)"

// Key renders the grouping key of a record, e.g. "agent=x, workflow=v1".
func (r Record) Key(groupBy []string) string {
	parts := make([]string, len(groupBy))
	for i, g := range groupBy {
		v := r.Value(g)
		if v == "" {
			v = None
		}
		parts[i] = g + "=" + v
	}
	return strings.Join(parts, ", ")
}
