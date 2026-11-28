// Package delivery decides whether keys may be typed into an agent's pane.
package delivery

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

// Session is one agent session.
type Session struct {
	ID          string
	AgentKind   string
	Native      bool
	Endpoint    *string
	StartupOver bool
}

// Pane reads the terminal an agent runs in.
type Pane interface {
	ShellReturned(name string) (bool, error)
}

// Service guards input to agent panes.
type Service struct {
	pane       Pane
	sessions   map[string]Session
	ended      map[string]string
	foreground map[string]bool
}

// New returns a Service over the given sessions.
func New(p Pane, sessions map[string]Session) *Service {
	return &Service{pane: p, sessions: sessions, ended: map[string]string{}, foreground: map[string]bool{}}
}

// Ended reports how a session was ended, if it was.
func (s *Service) Ended(id string) string { return s.ended[id] }

// VerifyPaneCommand checks that the agent itself is in the pane's foreground
// before keys are typed into it.
func (s *Service) VerifyPaneCommand(name, command string) (err error) {
	defer func() {
		if err != nil {
			log.Printf("agent input refused session=%s: %v", name, err)
		}
	}()
	session, ok := s.sessions[name]
	if !ok || session.ID == "" {
		return errors.New("message not delivered: " + name + " has no running agent")
	}
	if !agentForeground(session.AgentKind, command) {
		// Native tools may have a bridge in the foreground, but they must
		// never receive raw keys, regardless of that process.
		if !isShell(command) && (session.Endpoint != nil || session.Native) {
			return fmt.Errorf("message not delivered: %s requires its native route", name)
		}
		if command != "" && session.StartupOver {
			status := "interrupted"
			if isShell(command) {
				returned, err := s.pane.ShellReturned(name)
				if err != nil {
					return fmt.Errorf("message not delivered: verifying shell return: %w", err)
				}
				if !returned {
					return fmt.Errorf("message not delivered: %s foreground shell may still be running the agent", name)
				}
				status = "dead"
			}
			s.ended[session.ID] = status
		}
		return fmt.Errorf("message not delivered: %s foreground is %q, not agent %q", name, command, session.AgentKind)
	}
	s.foreground[session.ID] = true
	if session.Endpoint != nil || session.Native {
		return fmt.Errorf("message not delivered: %s requires its native route", name)
	}
	return nil
}

func agentForeground(kind, command string) bool {
	command = filepath.Base(command)
	switch kind {
	case "claude":
		return command == "claude" || strings.HasPrefix(command, "claude-")
	case "codex":
		return command == "codex"
	}
	return command == kind
}

func isShell(command string) bool {
	switch filepath.Base(command) {
	case "bash", "zsh", "sh", "fish":
		return true
	}
	return false
}
