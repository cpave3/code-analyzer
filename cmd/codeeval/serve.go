package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/cameron/code-evaluator/internal/dashboard"
	"github.com/cameron/code-evaluator/internal/languages"
	"github.com/cameron/code-evaluator/internal/provenance"
	"github.com/cameron/code-evaluator/internal/rules/all"
)

// serveCmd serves the quality-over-time dashboard on a loopback address.
func serveCmd(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var k conn
	k.register(fs)
	addr := fs.String("addr", "127.0.0.1:8765", "listen address (loopback only)")
	base := fs.String("base", "", "first commit's parent (exclusive); default: the last --max first-parent commits of --target")
	target := fs.String("target", "HEAD", "last commit")
	max := fs.Int("max", 60, "most commits shown")
	defects := fs.Bool("defects", true, "run the defect rules on every commit")
	near := fs.Bool("near-duplication", true, "run near-miss and fragment clone detection on every commit")
	infer := fs.Bool("infer-co-authors", true, "infer agents from Co-Authored-By lines when a commit has no provenance")
	notesRef := fs.String("notes-ref", provenance.NotesRef, "notes ref with provenance notes")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 || *max < 1 {
		fmt.Fprintln(stderr, "serve takes flags only; --max must be at least 1")
		return 2
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
	fmt.Fprintf(stdout, "codeeval dashboard for %s: http://%s/  (Ctrl-C to stop)\n", ev.RepoName(), ln.Addr())
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
