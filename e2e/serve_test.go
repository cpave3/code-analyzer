package e2e

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestServeDashboard(t *testing.T) {
	if r := cli(t, "serve", "--repo", repos["defects"], "--addr", "0.0.0.0:0"); r.code != 2 || !strings.Contains(r.stderr, "loopback") {
		t.Fatalf("non-loopback address accepted: %d %s", r.code, r.stderr)
	}
	if r := cli(t, "serve", repos["defects"], "--repo", repos["defects"]); r.code != 2 {
		t.Errorf("PATH together with --repo accepted: %d", r.code)
	}
	// The repository as a positional PATH; --base/--target narrow the history.
	cmd := exec.Command(binary, "serve", repos["defects"], "--config", scenarios["defects"].Config, "--base", "base", "--target", "bad", "--addr", "127.0.0.1:0", "--no-cache", "--defects", "--near-duplication")
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cmd.Process.Kill(); cmd.Wait() })
	line, err := bufio.NewReader(out).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	i := strings.Index(line, "http://")
	if i < 0 {
		t.Fatalf("no URL in %q", line)
	}
	url := strings.Fields(line[i:])[0]
	get := func(path string) (int, []byte) {
		resp, err := http.Get(url + path)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, b
	}
	code, page := get("")
	if code != 200 || !strings.Contains(string(page), "Quality over time") || !strings.Contains(string(page), `id="root"`) {
		t.Fatalf("page: %d", code)
	}
	// The page's script is embedded in the binary, and the page's own routes
	// load the page.
	m := regexp.MustCompile(`src="/(assets/[^"]+\.js)"`).FindStringSubmatch(string(page))
	if m == nil {
		t.Fatalf("no script in the page: %s", page)
	}
	if code, js := get(m[1]); code != 200 || len(js) < 1000 {
		t.Errorf("script %s: %d (%d bytes)", m[1], code, len(js))
	}
	for _, route := range []string{"modules/internal", "changes", "graph"} {
		if code, b := get(route); code != 200 || !strings.Contains(string(b), `id="root"`) {
			t.Errorf("route %s: %d", route, code)
		}
	}
	if code, _ := get("api/unknown"); code != 404 {
		t.Errorf("unknown API path: %d, want 404", code)
	}
	var meta struct {
		Commits []struct {
			SHA     string `json:"sha"`
			Subject string `json:"subject"`
		} `json:"commits"`
		GroupFields []string `json:"groupFields"`
	}
	_, b := get("api/meta")
	if err := json.Unmarshal(b, &meta); err != nil || len(meta.Commits) != 2 || len(meta.GroupFields) < 2 {
		t.Fatalf("meta: %v %s", err, b)
	}
	type metric struct {
		Done      bool     `json:"done"`
		StepDone  bool     `json:"stepDone"`
		Error     string   `json:"error"`
		Gate      *int     `json:"gate"`
		IntroGate *int     `json:"introducedGate"`
		Near      *int     `json:"nearMisses"`
		Score     *float64 `json:"score"`
	}
	var series struct {
		Done      int               `json:"done"`
		StepsDone int               `json:"stepsDone"`
		DetailSD  int               `json:"detailStepsDone"`
		Total     int               `json:"total"`
		Window    map[string]int    `json:"window"`
		Metrics   map[string]metric `json:"metrics"`
	}
	deadline := time.Now().Add(60 * time.Second)
	for {
		_, b = get("api/metrics")
		if err := json.Unmarshal(b, &series); err != nil {
			t.Fatal(err)
		}
		if (series.Done == series.Total && series.StepsDone == series.Total && series.DetailSD == series.Total) || time.Now().After(deadline) {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if series.Done != 2 {
		t.Fatalf("metrics not computed: %+v", series)
	}
	if series.Window["from"] != 0 || series.Window["to"] != 1 {
		t.Errorf("default window (--last 50) should cover both commits: %v", series.Window)
	}
	base, bad := series.Metrics["0"], series.Metrics["1"]
	if base.Error != "" || bad.Error != "" || base.Gate == nil || *base.Gate != 0 || bad.IntroGate == nil || *bad.IntroGate == 0 || bad.Near == nil || bad.Score == nil {
		t.Errorf("series: base %+v bad %+v", base, bad)
	}
	if code, g := get("api/graph?commit=" + meta.Commits[1].SHA); code != 200 || !strings.Contains(string(g), `"modules"`) {
		t.Errorf("graph: %d", code)
	}
	post := func(body, origin, ctype string) int {
		req, _ := http.NewRequest("POST", url+"api/window", strings.NewReader(body))
		if ctype != "" {
			req.Header.Set("Content-Type", ctype)
		}
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		return resp.StatusCode
	}
	host := strings.TrimSuffix(strings.TrimPrefix(url, "http://"), "/")
	if code := post(`{"from":1,"to":1}`, "http://evil.example", "application/json"); code != 403 {
		t.Errorf("cross-origin window change: %d", code)
	}
	if code := post(`{"from":1,"to":1}`, "", "text/plain"); code != 415 {
		t.Errorf("non-JSON window change: %d", code)
	}
	if code := post(`{"from":1,"to":1}`, "http://"+host, "application/json"); code != 200 {
		t.Errorf("same-origin window change: %d", code)
	}
	_, b = get("api/metrics")
	json.Unmarshal(b, &series)
	if series.Window["from"] != 1 || series.Window["to"] != 1 || series.Total != 1 {
		t.Errorf("window after POST: %v total %d", series.Window, series.Total)
	}
	if code, _ := get("api/graph?commit=0000000"); code != 404 {
		t.Errorf("unknown commit: %d", code)
	}

	// Folder report: explain for a path, as `codeeval explain --path`.
	badSHA := meta.Commits[1].SHA
	var ex struct {
		Scope   any   `json:"scope"`
		Factors []any `json:"factors"`
	}
	if code, b := get("api/explain?commit=" + badSHA + "&path="); code != 200 || json.Unmarshal(b, &ex) != nil {
		t.Fatalf("explain: %d %s", code, b)
	}
	// Defects run only on demand: cached=1 answers 204 until a check ran.
	if code, _ := get("api/defects?commit=" + badSHA + "&path=&cached=1"); code != 204 {
		t.Errorf("defects before a check: %d, want 204", code)
	}
	type defects struct {
		Findings []struct {
			File string `json:"file"`
		} `json:"findings"`
		Totals struct{ Gate, Advisory int } `json:"totals"`
	}
	var all, none defects
	if code, b := get("api/defects?commit=" + badSHA + "&path="); code != 200 || json.Unmarshal(b, &all) != nil || all.Totals.Gate == 0 {
		t.Fatalf("defects: %d %s", code, b)
	}
	if code, b := get("api/defects?commit=" + badSHA + "&path=&cached=1"); code != 200 || json.Unmarshal(b, &all) != nil {
		t.Errorf("defects after a check: %d", code)
	}
	if code, b := get("api/defects?commit=" + badSHA + "&path=no/such/dir&cached=1"); code != 200 || json.Unmarshal(b, &none) != nil || len(none.Findings) != 0 {
		t.Errorf("defects outside the path should be filtered out: %d %s", code, b)
	}
	if code, _ := get("api/defects?commit=0000000&cached=1"); code != 404 {
		t.Errorf("defects, unknown commit: %d", code)
	}
}
