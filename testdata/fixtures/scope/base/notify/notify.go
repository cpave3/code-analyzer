package notify

import "fmt"

// Subject builds an email subject.
func Subject(user string, n int) string { return fmt.Sprintf("%s, you have %d new messages", user, n) }

// Body builds an email body.
func Body(lines []string) string {
	out := ""
	for _, l := range lines {
		out += "* " + l + "\n"
	}
	return out
}
