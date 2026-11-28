// Package gitsrc reads a Git repository strictly read-only through the Git
// CLI. It never checks out, resets, stashes, fetches or writes refs, and it
// neutralises repository configuration that could run programs (fsmonitor,
// external diff, textconv, hooks).
package gitsrc

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// WorktreeRev is the pseudo revision used for working-tree evaluation.
const WorktreeRev = "WORKTREE"

// Repo is a read-only handle on a repository.
type Repo struct {
	Path string // absolute top-level path

	mu    sync.Mutex
	batch *catFile
}

// safeArgs are prepended to every git invocation.
var safeArgs = []string{
	"--no-optional-locks",
	"-c", "core.fsmonitor=false",
	"-c", "core.hooksPath=/dev/null",
	"-c", "core.quotePath=false",
	"-c", "core.untrackedCache=false",
	"-c", "protocol.allow=never", // no network transports, ever
	"-c", "credential.helper=",
}

func gitEnv() []string {
	env := []string{}
	for _, kv := range os.Environ() {
		k, _, _ := strings.Cut(kv, "=")
		switch {
		case strings.HasPrefix(k, "GIT_"):
			continue // do not inherit GIT_DIR, GIT_INDEX_FILE, ...
		}
		env = append(env, kv)
	}
	return append(env,
		"GIT_OPTIONAL_LOCKS=0",
		"GIT_TERMINAL_PROMPT=0",
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_NO_REPLACE_OBJECTS=1",
		"LC_ALL=C",
	)
}

func (r *Repo) command(args ...string) *exec.Cmd {
	full := append(append([]string{}, safeArgs...), args...)
	cmd := exec.Command("git", full...)
	cmd.Dir = r.Path
	cmd.Env = gitEnv()
	return cmd
}

func (r *Repo) run(args ...string) ([]byte, error) {
	cmd := r.command(args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

// Open resolves the repository top level.
func Open(path string) (*Repo, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	r := &Repo{Path: abs}
	out, err := r.run("rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %s", path)
	}
	r.Path = strings.TrimSpace(string(out))
	return r, nil
}

// Close stops the background cat-file process.
func (r *Repo) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.batch != nil {
		r.batch.close()
		r.batch = nil
	}
}

// ResolveCommit resolves a revision expression to a full commit SHA. Option
// injection is refused.
func (r *Repo) ResolveCommit(rev string) (string, error) {
	if strings.HasPrefix(rev, "-") || rev == "" {
		return "", fmt.Errorf("invalid revision %q", rev)
	}
	out, err := r.run("rev-parse", "--verify", "--end-of-options", rev+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("cannot resolve revision %q", rev)
	}
	return strings.TrimSpace(string(out)), nil
}

// TreeOf returns the tree SHA of a commit.
func (r *Repo) TreeOf(commit string) (string, error) {
	out, err := r.run("rev-parse", "--verify", "--end-of-options", commit+"^{tree}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// FirstParentChain lists commits on the first-parent path base..target,
// oldest first.
func (r *Repo) FirstParentChain(base, target string) ([]string, error) {
	out, err := r.run("rev-list", "--first-parent", "--reverse", base+".."+target)
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(out)), nil
}

// FirstParent returns the first parent of a commit ("" for a root commit).
func (r *Repo) FirstParent(commit string) (string, error) {
	out, err := r.run("rev-list", "--parents", "-n", "1", commit)
	if err != nil {
		return "", err
	}
	f := strings.Fields(string(out))
	if len(f) < 2 {
		return "", nil
	}
	return f[1], nil
}

// IsAncestor reports whether a is an ancestor of b.
func (r *Repo) IsAncestor(a, b string) bool {
	cmd := r.command("merge-base", "--is-ancestor", a, b)
	return cmd.Run() == nil
}

// Subject returns the one-line subject of a commit (untrusted data).
func (r *Repo) Subject(commit string) string {
	out, err := r.run("log", "-1", "--format=%s", commit)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// CommitTime returns the committer time of a commit in ISO-8601 strict form.
func (r *Repo) CommitTime(commit string) string {
	out, err := r.run("log", "-1", "--format=%cI", commit)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Entry is a file in a tree.
type Entry struct {
	Path string
	Mode string
	Blob string // object id (or computed blob id for worktree)
	Size int64
}

// ListTree lists all blobs in a commit's tree (recursively). Submodules
// (commit entries) and symlinks are reported with their mode so callers can
// exclude them.
func (r *Repo) ListTree(commit string) ([]Entry, error) {
	out, err := r.run("ls-tree", "-r", "-z", "--full-tree", "-l", commit)
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for _, rec := range bytes.Split(out, []byte{0}) {
		if len(rec) == 0 {
			continue
		}
		meta, path, ok := bytes.Cut(rec, []byte{'\t'})
		if !ok {
			continue
		}
		f := strings.Fields(string(meta))
		if len(f) < 4 {
			continue
		}
		size, _ := strconv.ParseInt(f[3], 10, 64)
		entries = append(entries, Entry{Path: string(path), Mode: f[0], Blob: f[2], Size: size})
	}
	return entries, nil
}

// MaxUntrackedBytes bounds untracked files read in worktree mode.
const MaxUntrackedBytes = 8 << 20

// ListWorktree lists tracked files present in the working tree with content
// hashed like Git blobs. With untracked, files Git does not track but does
// not ignore (.gitignore, .git/info/exclude, core.excludesFile) are included
// too, up to MaxUntrackedBytes each; ignored files are never read. Symlinks
// and non-regular files are skipped.
func (r *Repo) ListWorktree(untracked bool) ([]Entry, map[string][]byte, error) {
	out, err := r.run("ls-files", "-z", "--cached", "--full-name")
	if err != nil {
		return nil, nil, err
	}
	tracked := len(bytes.Split(out, []byte{0}))
	if untracked {
		more, err := r.run("ls-files", "-z", "--others", "--exclude-standard", "--full-name")
		if err != nil {
			return nil, nil, err
		}
		out = append(out, more...)
	}
	seen := map[string]bool{}
	var entries []Entry
	contents := map[string][]byte{}
	for i, p := range bytes.Split(out, []byte{0}) {
		if len(p) == 0 || seen[string(p)] {
			continue
		}
		isUntracked := i >= tracked
		rel := string(p)
		seen[rel] = true
		full := filepath.Join(r.Path, filepath.FromSlash(rel))
		fi, err := os.Lstat(full)
		if errors.Is(err, fs.ErrNotExist) {
			continue // deleted in worktree
		}
		if err != nil || !fi.Mode().IsRegular() || (isUntracked && fi.Size() > MaxUntrackedBytes) {
			continue
		}
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		mode := "100644"
		if fi.Mode()&0o111 != 0 {
			mode = "100755"
		}
		id := BlobID(data)
		contents[id] = data
		entries = append(entries, Entry{Path: rel, Mode: mode, Blob: id, Size: int64(len(data))})
	}
	return entries, contents, nil
}

// BlobID computes the Git SHA-1 blob id of content.
func BlobID(data []byte) string {
	h := sha1.New()
	fmt.Fprintf(h, "blob %d\x00", len(data))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// NameStatus is one entry of a name-status diff.
type NameStatus struct {
	Status     string // A, M, D, R, C, T
	Similarity int
	OldPath    string
	NewPath    string
}

// DiffNameStatus returns changed paths between base and target commits, or
// between base and the working tree (tracked files only) when target is
// WorktreeRev. Rename detection is enabled; external diff and textconv are
// disabled.
func (r *Repo) DiffNameStatus(base, target string) ([]NameStatus, error) {
	args := []string{"diff", "--no-ext-diff", "--no-textconv", "--no-color", "-M", "-z", "--name-status"}
	if target == WorktreeRev {
		args = append(args, base, "--")
	} else {
		args = append(args, base, target, "--")
	}
	out, err := r.run(args...)
	if err != nil {
		return nil, err
	}
	parts := bytes.Split(out, []byte{0})
	var res []NameStatus
	for i := 0; i < len(parts); i++ {
		code := string(parts[i])
		if code == "" {
			continue
		}
		ns := NameStatus{Status: code[:1]}
		if len(code) > 1 {
			ns.Similarity, _ = strconv.Atoi(code[1:])
		}
		switch ns.Status {
		case "R", "C":
			if i+2 >= len(parts) {
				return nil, fmt.Errorf("malformed name-status output")
			}
			ns.OldPath, ns.NewPath = string(parts[i+1]), string(parts[i+2])
			i += 2
		default:
			if i+1 >= len(parts) {
				return nil, fmt.Errorf("malformed name-status output")
			}
			ns.OldPath, ns.NewPath = string(parts[i+1]), string(parts[i+1])
			i++
		}
		res = append(res, ns)
	}
	return res, nil
}

// ReadBlob returns the content of a blob object.
func (r *Repo) ReadBlob(id string) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.batch == nil {
		b, err := r.startCatFile()
		if err != nil {
			return nil, err
		}
		r.batch = b
	}
	return r.batch.read(id)
}

// ReadTreeFile reads a path from a commit (used for build files like go.mod).
func (r *Repo) ReadTreeFile(commit, path string) ([]byte, error) {
	return r.run("cat-file", "blob", commit+":"+path)
}

type catFile struct {
	cmd *exec.Cmd
	in  io.WriteCloser
	out *bufio.Reader
}

func (r *Repo) startCatFile() (*catFile, error) {
	cmd := r.command("cat-file", "--batch")
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &catFile{cmd: cmd, in: in, out: bufio.NewReaderSize(out, 1<<16)}, nil
}

func (c *catFile) read(id string) ([]byte, error) {
	if _, err := fmt.Fprintln(c.in, id); err != nil {
		return nil, err
	}
	header, err := c.out.ReadString('\n')
	if err != nil {
		return nil, err
	}
	f := strings.Fields(header)
	if len(f) == 2 && f[1] == "missing" {
		return nil, fmt.Errorf("object %s missing", id)
	}
	if len(f) != 3 {
		return nil, fmt.Errorf("unexpected cat-file header %q", header)
	}
	size, err := strconv.Atoi(f[2])
	if err != nil {
		return nil, err
	}
	buf := make([]byte, size+1) // content + trailing newline
	if _, err := io.ReadFull(c.out, buf); err != nil {
		return nil, err
	}
	return buf[:size], nil
}

func (c *catFile) close() {
	c.in.Close()
	_ = c.cmd.Wait()
}

// CommitMessage is a commit's message and parents.
type CommitMessage struct {
	SHA     string
	Parents []string
	Message string // full message (untrusted data)
	Note    string // note text from the requested notes ref, if any
}

// HasRef reports whether a fully qualified ref exists.
func (r *Repo) HasRef(ref string) bool {
	_, err := r.run("rev-parse", "--verify", "-q", "--end-of-options", ref)
	return err == nil
}

// CommitsBetween lists commits reachable from target but not from base
// (all parents, not only the first-parent path; base "" = every commit
// reachable from target), newest first, with their
// messages and, when notesRef exists, their notes.
func (r *Repo) CommitsBetween(base, target, notesRef string) ([]CommitMessage, error) {
	args := []string{"log", "--no-show-signature", "--format=%H%x1f%P%x1f%B%x1f%N%x1e"}
	if notesRef != "" && r.HasRef(notesRef) {
		args = append(args, "--notes="+notesRef)
	} else {
		args = append(args, "--no-notes")
	}
	spec := target
	if base != "" {
		spec = base + ".." + target
	}
	args = append(args, "--end-of-options", spec)
	out, err := r.run(args...)
	if err != nil {
		return nil, err
	}
	var cs []CommitMessage
	for _, rec := range strings.Split(string(out), "\x1e") {
		rec = strings.TrimLeft(rec, "\n")
		if rec == "" {
			continue
		}
		f := strings.SplitN(rec, "\x1f", 4)
		if len(f) < 4 {
			continue
		}
		cs = append(cs, CommitMessage{SHA: f[0], Parents: strings.Fields(f[1]), Message: f[2], Note: strings.TrimSpace(f[3])})
	}
	return cs, nil
}

// RevList lists commits reachable from target but not from base.
func (r *Repo) RevList(base, target string) ([]string, error) {
	spec := target
	if base != "" {
		spec = base + ".." + target
	}
	out, err := r.run("rev-list", "--end-of-options", spec)
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(out)), nil
}

// CommitInfo is a commit's author, author time and subject (untrusted text).
type CommitInfo struct {
	Author, Time, Subject string
	Parents               []string
}

// Info returns a commit's author, author time (RFC 3339) and subject.
func (r *Repo) Info(commit string) (CommitInfo, error) {
	out, err := r.run("log", "-1", "--no-show-signature", "--format=%an%x1f%aI%x1f%P%x1f%s", "--end-of-options", commit)
	if err != nil {
		return CommitInfo{}, err
	}
	f := strings.SplitN(strings.TrimRight(string(out), "\n"), "\x1f", 4)
	if len(f) < 4 {
		return CommitInfo{}, fmt.Errorf("unexpected log output for %s", commit)
	}
	return CommitInfo{Author: f[0], Time: f[1], Parents: strings.Fields(f[2]), Subject: f[3]}, nil
}

// LastFirstParents lists up to n commits (n <= 0: all) on target's
// first-parent chain, oldest first, ending with target.
func (r *Repo) LastFirstParents(target string, n int) ([]string, error) {
	args := []string{"rev-list", "--first-parent"}
	if n > 0 {
		args = append(args, "-n", strconv.Itoa(n))
	}
	out, err := r.run(append(args, "--end-of-options", target)...)
	if err != nil {
		return nil, err
	}
	shas := strings.Fields(string(out))
	for i, j := 0, len(shas)-1; i < j; i, j = i+1, j-1 {
		shas[i], shas[j] = shas[j], shas[i]
	}
	return shas, nil
}
