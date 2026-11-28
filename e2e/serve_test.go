package e2e

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestServeDashboard(t *testing.T) {
	if r := cli(t, "serve", "--repo", repos["defects"], "--addr", "0.0.0.0:0"); r.code != 2 || !strings.Contains(r.stderr, "loopback") {
		t.Fatalf("non-loopback address accepted: %d %s", r.code, r.stderr)
	}
	cmd := exec.Command(binary, "serve", "--repo", repos["defects"], "--config", scenarios["defects"].Config, "--base", "base", "--target", "bad", "--addr", "127.0.0.1:0", "--no-cache")
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
	if code, page := get(""); code != 200 || !strings.Contains(string(page), "Quality over time") {
		t.Fatalf("page: %d", code)
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
	var series struct {
		Done, Total int
		Metrics     []struct {
			Done      bool     `json:"done"`
			Error     string   `json:"error"`
			Gate      *int     `json:"gate"`
			IntroGate *int     `json:"introducedGate"`
			Near      *int     `json:"nearMisses"`
			Score     *float64 `json:"score"`
		} `json:"metrics"`
	}
	deadline := time.Now().Add(60 * time.Second)
	for {
		_, b = get("api/metrics")
		if err := json.Unmarshal(b, &series); err != nil {
			t.Fatal(err)
		}
		if series.Done == series.Total || time.Now().After(deadline) {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if series.Done != 2 {
		t.Fatalf("metrics not computed: %+v", series)
	}
	base, bad := series.Metrics[0], series.Metrics[1]
	if base.Error != "" || bad.Error != "" || base.Gate == nil || *base.Gate != 0 || bad.IntroGate == nil || *bad.IntroGate == 0 || bad.Near == nil || bad.Score == nil {
		t.Errorf("series: base %+v bad %+v", base, bad)
	}
	if code, g := get("api/graph?commit=" + meta.Commits[1].SHA); code != 200 || !strings.Contains(string(g), `"modules"`) {
		t.Errorf("graph: %d", code)
	}
	if code, _ := get("api/graph?commit=0000000"); code != 404 {
		t.Errorf("unknown commit: %d", code)
	}
}
