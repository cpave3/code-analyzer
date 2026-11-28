package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHookInstallUninstall(t *testing.T) {
	dir := t.TempDir()
	git(t, dir, "init", "-q")
	os.MkdirAll(filepath.Join(dir, ".claude"), 0o755)
	os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	original := "{\n  \"model\": \"opus\",\n  \"hooks\": {\"Stop\": [{\"hooks\": [{\"type\": \"command\", \"command\": \"echo mine\"}]}]},\n  \"env\": {\"A\": \"1\"}\n}\n"
	settings := filepath.Join(dir, ".claude", "settings.json")
	os.WriteFile(settings, []byte(original), 0o644)

	run := func(wd string, args ...string) run {
		t.Helper()
		cmd := exec.Command(binary, args...)
		cmd.Dir = wd
		var out, errb strings.Builder
		cmd.Stdout, cmd.Stderr = &out, &errb
		err := cmd.Run()
		code := 0
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		}
		return struct {
			code           int
			stdout, stderr string
		}{code, out.String(), errb.String()}
	}
	if r := run(filepath.Join(dir, "sub"), "hook", "install", "--tool", "claude-code"); r.code != 0 {
		t.Fatalf("install: %d %s", r.code, r.stderr)
	}
	b, _ := os.ReadFile(settings)
	s := string(b)
	for _, want := range []string{"codeeval-stop-analyze.sh", "codeeval-session-base.sh", "echo mine", `"A": "1"`} {
		if !strings.Contains(s, want) {
			t.Errorf("settings lack %q:\n%s", want, s)
		}
	}
	if strings.Index(s, `"model"`) > strings.Index(s, `"hooks"`) || strings.Index(s, `"hooks"`) > strings.Index(s, `"env"`) {
		t.Errorf("top-level key order not kept:\n%s", s)
	}
	for _, f := range []string{"codeeval-stop-analyze.sh", "codeeval-session-base.sh"} {
		if fi, err := os.Stat(filepath.Join(dir, ".claude", "hooks", f)); err != nil || fi.Mode().Perm()&0o100 == 0 {
			t.Errorf("%s missing or not executable: %v", f, err)
		}
	}
	if r := run(dir, "hook", "install", "--tool", "claude-code"); r.code != 0 || strings.Count(r.stdout, "unchanged") != 3 {
		t.Errorf("reinstall should change nothing:\n%s", r.stdout)
	}
	if r := run(dir, "hook", "uninstall", "--tool", "claude-code"); r.code != 0 {
		t.Fatalf("uninstall: %s", r.stderr)
	}
	b, _ = os.ReadFile(settings)
	if strings.Contains(string(b), "codeeval") || !strings.Contains(string(b), "echo mine") || !strings.Contains(string(b), `"model"`) {
		t.Errorf("uninstall left codeeval or dropped other settings:\n%s", b)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude", "hooks", "codeeval-stop-analyze.sh")); err == nil {
		t.Error("hook script not removed")
	}
	os.WriteFile(settings, []byte("{not json"), 0o644)
	if r := run(dir, "hook", "install", "--tool", "claude-code"); r.code == 0 || !strings.Contains(r.stderr, "not valid JSON") {
		t.Errorf("invalid settings.json must be refused: %d %s", r.code, r.stderr)
	}
	if b, _ := os.ReadFile(settings); string(b) != "{not json" {
		t.Error("invalid settings.json was rewritten")
	}
	if r := run(dir, "hook", "install", "--tool", "cursor"); r.code == 0 {
		t.Error("unknown tool accepted")
	}
	if r := run(t.TempDir(), "hook", "install", "--tool", "claude-code"); r.code == 0 || !strings.Contains(r.stderr, "not inside a git repository") {
		t.Errorf("outside a repository: %d %s", r.code, r.stderr)
	}
}
