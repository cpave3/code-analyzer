package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cameron/code-evaluator/internal/hookinstall"
)

// hookCmd installs or removes CodeEval's agent hooks in the current project.
func hookCmd(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || (args[0] != "install" && args[0] != "uninstall") {
		fmt.Fprintln(stderr, "usage: codeeval hook install|uninstall --tool claude-code [--dir PATH] [--bin PATH] [--config FILE] [--dry-run]")
		return 2
	}
	action := args[0]
	fs := flag.NewFlagSet("hook "+action, flag.ContinueOnError)
	fs.SetOutput(stderr)
	tool := fs.String("tool", "", "agent tool: claude-code")
	dir := fs.String("dir", ".", "a directory inside the project (hooks go in its git top level)")
	bin := fs.String("bin", "", "codeeval binary the hook runs (default: codeeval on PATH, else this binary)")
	cfg := fs.String("config", "", "codeeval config file the hook passes to analyze")
	dry := fs.Bool("dry-run", false, "show what would change without writing")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}
	if *tool == "" || fs.NArg() > 0 {
		fmt.Fprintln(stderr, "hook "+action+" needs --tool claude-code")
		return 2
	}
	root, err := gitTop(*dir)
	if err != nil {
		return fail(stderr, err)
	}
	o := hookinstall.Options{Root: root, DryRun: *dry}
	if action == "install" {
		if o.Bin, err = hookBinary(*bin); err != nil {
			return fail(stderr, err)
		}
		if *cfg != "" {
			if o.Config, err = filepath.Abs(*cfg); err != nil {
				return fail(stderr, err)
			}
			if _, err := os.Stat(o.Config); err != nil {
				return fail(stderr, err)
			}
		}
	}
	var changes []hookinstall.Change
	if action == "install" {
		changes, err = hookinstall.Install(*tool, o)
	} else {
		changes, err = hookinstall.Uninstall(*tool, o)
	}
	if err != nil {
		return fail(stderr, err)
	}
	verb := ""
	if *dry {
		verb = "would be "
	}
	for _, c := range changes {
		rel, _ := filepath.Rel(root, c.Path)
		fmt.Fprintf(stdout, "%s%-9s %s\n", verb, c.Action, rel)
	}
	if action == "install" && !*dry {
		runs := "codeeval on PATH"
		if o.Bin != "" {
			runs = o.Bin
		}
		fmt.Fprintf(stdout, "\nInstalled the CodeEval hooks for Claude Code in %s (running %s).\n", root, runs)
		fmt.Fprintln(stdout, "- SessionStart records the commit a session starts from and prints the check to run before pushing; Stop runs `codeeval analyze` from there (or from its last passing stop) and blocks once on introduced gate defects or a debt rise of more than CODEEVAL_DEBT_BUDGET (default 10) points, unless the pull request is already approved.")
		fmt.Fprintln(stdout, "- Escape hatches: CODEEVAL_HOOK=off, CODEEVAL_HOOK_MAX_BLOCKS (default 3), `codeeval-ignore <rule>: <reason>` comments.")
		fmt.Fprintln(stdout, "- New Claude Code sessions pick it up; commit .claude/ to share it with the team, or run `codeeval hook uninstall --tool claude-code` to remove it.")
	}
	return 0
}

// gitTop returns the git top level containing dir.
func gitTop(dir string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s is not inside a git repository", dir)
	}
	return strings.TrimSpace(string(out)), nil
}

// hookBinary picks the binary the hook runs: the --bin given, nothing when
// codeeval is on PATH, else this binary's absolute path.
func hookBinary(bin string) (string, error) {
	if bin != "" {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		if _, err := os.Stat(abs); err != nil {
			return "", err
		}
		return abs, nil
	}
	if _, err := exec.LookPath("codeeval"); err == nil {
		return "", nil
	}
	self, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(self)
}
