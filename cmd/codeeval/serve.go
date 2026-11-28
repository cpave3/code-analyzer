package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/cameron/code-evaluator/internal/dashboard"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/provenance"
	"github.com/cameron/code-evaluator/internal/rules/all"
)

// serveCmd serves the quality-over-time dashboard on a loopback address.
// With no arguments it shows the whole first-parent history of the current
// checkout's HEAD; `codeeval serve PATH` does the same for another
// checkout, and --target/--base/--max narrow it.
func serveCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var k conn
	k.register(fs)
	addr := fs.String("addr", "127.0.0.1:8765", "listen address (loopback only)")
	base := fs.String("base", "", "start after this commit (exclusive); default: the whole first-parent history of --target")
	target := fs.String("target", "HEAD", "last commit (default: the checkout's HEAD)")
	max := fs.Int("max", 0, "most recent commits shown (0 = all)")
	defects := fs.Bool("defects", true, "run the defect rules on every commit")
	near := fs.Bool("near-duplication", true, "run near-miss and fragment clone detection on every commit")
	infer := fs.Bool("infer-co-authors", true, "infer agents from Co-Authored-By lines when a commit has no provenance")
	notesRef := fs.String("notes-ref", provenance.NotesRef, "notes ref with provenance notes")
	// PATH may come before the flags (`serve ../repo --max 50`): Go's flag
	// parser stops at the first non-flag argument, so take it off first.
	var path string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		path, args = args[0], args[1:]
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if path != "" && fs.NArg() > 0 {
		fmt.Fprintln(stderr, "usage: codeeval serve [PATH] [flags]: give one PATH")
		return 2
	}
	if path == "" && fs.NArg() == 1 {
		path = fs.Arg(0)
	} else if fs.NArg() > 1 {
		fmt.Fprintln(stderr, "usage: codeeval serve [PATH] [flags]: give one PATH")
		return 2
	}
	if *max < 0 {
		fmt.Fprintln(stderr, "usage: codeeval serve [PATH] [flags]  (PATH defaults to the current directory; --max must not be negative)")
		return 2
	}
	if path != "" {
		if repoSet(fs) {
			fmt.Fprintln(stderr, "give the repository either as PATH or with --repo, not both")
			return 2
		}
		k.repo = path
	}
	if err := loopback(*addr); err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 2
	}
	ev, cfg, code := k.open(stderr)
	if ev == nil {
		return code
	}
	defer ev.Close()
	langs := languages.Default()
	srv, err := dashboard.New(ev, cfg, langs, all.Default(langs), dashboard.Options{Base: *base, Target: *target, Max: *max,
		Defects: *defects, Clones: *near, NotesRef: *notesRef, Infer: *infer})
	if err != nil {
		return fail(stderr, err)
	}
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		return fail(stderr, err)
	}
	go srv.Run()
	fmt.Fprintf(stdout, "codeeval dashboard for %s: http://%s/  (%d commits, measured newest first; Ctrl-C to stop)\n", ev.RepoName(), ln.Addr(), srv.Commits())
	hs := &http.Server{Handler: srv.Handler(), ReadHeaderTimeout: 10 * time.Second}
	if err := hs.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fail(stderr, err)
	}
	return 0
}

// loopback refuses listen addresses other people could reach: the dashboard
// shows source structure and commit metadata.
func loopback(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if host == "localhost" {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return nil
	}
	return fmt.Errorf("--addr %s is not a loopback address; the dashboard only listens on 127.0.0.1, ::1 or localhost", addr)
}

// repoSet reports whether --repo was given explicitly.
func repoSet(fs *flag.FlagSet) bool {
	set := false
	fs.Visit(func(f *flag.Flag) { set = set || f.Name == "repo" })
	return set
}
