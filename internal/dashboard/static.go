package dashboard

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// static serves the built page (web/dist, from `pnpm build`): files as they
// are, with hashed assets cached for good, and index.html for any other path
// without a file extension, so the page's own routes (/modules/..., /changes)
// load the app. API paths and missing files are 404s.
func static(page fs.FS) http.Handler {
	files := http.FileServerFS(page)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if p == "" || p == "." || p == "index.html" {
			serveIndex(w, r, page)
			return
		}
		if st, err := fs.Stat(page, p); err == nil && !st.IsDir() {
			if strings.HasPrefix(p, "assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			files.ServeHTTP(w, r)
			return
		}
		if !route(p) && (p == "api" || p == "assets" || strings.HasPrefix(p, "api/") || strings.HasPrefix(p, "assets/") || path.Ext(p) != "") {
			http.NotFound(w, r)
			return
		}
		serveIndex(w, r, page)
	})
}

// route reports whether p is one of the page's own routes. Folder paths
// under /modules/ may contain dots (a C# project such as My.App), so they
// are not mistaken for missing files.
func route(p string) bool {
	return p == "modules" || strings.HasPrefix(p, "modules/")
}

func serveIndex(w http.ResponseWriter, r *http.Request, page fs.FS) {
	b, err := fs.ReadFile(page, "index.html")
	if err != nil {
		http.Error(w, "the dashboard page is not built: run `task web:build`", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(b)
}
