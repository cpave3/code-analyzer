package notify

import "fmt"

// Subject builds the subject line of a notification email.
func Subject(user string, n int) string {
	return fmt.Sprintf("%s, you have %d new messages", user, n)
}

// Body builds a bulleted notification body,
// one bullet per line.
func Body(lines []string) string {
	out := ""
	for _, l := range lines {
		out += "* " + l + "\n" // bullet
	}
	return out
}
