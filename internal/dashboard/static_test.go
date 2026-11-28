package dashboard

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestStaticServesAssetsAndFallsBackToTheApp(t *testing.T) {
	page := fstest.MapFS{
		"index.html":         {Data: []byte(`<div id="root"></div>`)},
		"assets/app-1a2b.js": {Data: []byte(`console.log(1)`)},
	}
	h := static(page)
	get := func(path string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		return w
	}
	body := func(w *httptest.ResponseRecorder) string { b, _ := io.ReadAll(w.Result().Body); return string(b) }

	if w := get("/"); w.Code != 200 || !strings.Contains(body(w), "root") || w.Header().Get("Cache-Control") != "no-cache" {
		t.Errorf("/: %d %q %q", w.Code, body(w), w.Header().Get("Cache-Control"))
	}
	// Client-side routes load the app.
	for _, p := range []string{"/modules/internal/service", "/changes", "/graph"} {
		if w := get(p); w.Code != 200 || !strings.Contains(body(w), "root") {
			t.Errorf("%s: %d %q", p, w.Code, body(w))
		}
	}
	// Hashed assets are served as files and cached for good.
	if w := get("/assets/app-1a2b.js"); w.Code != 200 || body(w) != "console.log(1)" || !strings.Contains(w.Header().Get("Cache-Control"), "immutable") {
		t.Errorf("asset: %d %q %q", w.Code, body(w), w.Header().Get("Cache-Control"))
	}
	// A missing asset or API path is a 404, not the app.
	for _, p := range []string{"/assets/missing.js", "/assets/", "/api/nope", "/favicon.ico"} {
		if w := get(p); w.Code != 404 {
			t.Errorf("%s: want 404, got %d", p, w.Code)
		}
	}
}
