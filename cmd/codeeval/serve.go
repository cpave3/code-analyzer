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
	base := fs.String("base", "", "list history from this commit (it is the first point); default: the whole first-parent history of --target")
	target := fs.String("target", "HEAD", "last commit (default: the checkout's HEAD)")
	max := fs.Int("max", 0, "most recent commits listed (0 = all); listing is cheap, measuring is not")
	last := fs.Int("last", 50, "measure the newest N commits at start (the page can move the window)")
	since := fs.String("since", "", "measure commits from this date at start (YYYY-MM-DD or RFC 3339); with --until, overrides --last")
	until := fs.String("until", "", "measure commits up to this date at start")
	defects := fs.Bool("defects", false, "also chart defect findings for every commit in the window (slow; by default defects are checked on demand per commit)")
	near := fs.Bool("near-duplication", false, "also chart near-miss and fragment clones for every commit in the window (slow)")
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
	if *max < 0 || *last < 0 {
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
		Last: *last, Since: *since, Until: *until, Defects: *defects, Clones: *near, NotesRef: *notesRef, Infer: *infer})
	if err != nil {
		return fail(stderr, err)
	}
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		return fail(stderr, err)
	}
	go srv.Run()
	lo, hi := srv.Window()
	fmt.Fprintf(stdout, "codeeval dashboard for %s: http://%s/  (%d commits listed; measuring %d of them, newest first; Ctrl-C to stop)\n", ev.RepoName(), ln.Addr(), srv.Commits(), hi-lo+1)
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
