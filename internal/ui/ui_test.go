package ui

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newUI(t *testing.T) http.Handler {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "assets"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>app</html>"), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("console.log(1)"), 0o644); err != nil {
		t.Fatalf("write asset: %v", err)
	}
	handler := New(Deps{
		Dir:          dir,
		PublicURL:    "https://terrakube.example.com",
		DexClientID:  "terrakube-ui",
		DexAuthority: "https://auth.example.com/dex",
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	if handler == nil {
		t.Fatal("New returned nil for a valid directory")
	}
	return handler
}

func get(t *testing.T, handler http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestServesTheIndexAndAssets(t *testing.T) {
	handler := newUI(t)

	rec := get(t, handler, "/")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "app") {
		t.Fatalf("index = %d %q", rec.Code, rec.Body.String())
	}

	rec = get(t, handler, "/assets/app.js")
	if rec.Code != http.StatusOK || rec.Body.String() != "console.log(1)" {
		t.Fatalf("asset = %d %q", rec.Code, rec.Body.String())
	}
	if cache := rec.Header().Get("Cache-Control"); !strings.Contains(cache, "immutable") {
		t.Fatalf("asset Cache-Control = %q", cache)
	}
}

// A deep link is a client-side route: the app must load and route itself.
func TestDeepLinksServeTheIndex(t *testing.T) {
	handler := newUI(t)
	rec := get(t, handler, "/organizations/acme/workspaces")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "app") {
		t.Fatalf("deep link = %d %q", rec.Code, rec.Body.String())
	}
	if cache := rec.Header().Get("Cache-Control"); cache != "no-cache" {
		t.Fatalf("index Cache-Control = %q, want no-cache", cache)
	}
}

// Returning HTML for a missing asset makes the browser report a syntax error
// instead of a missing file, which hides deployment mistakes.
func TestMissingAssetsAreNotFound(t *testing.T) {
	handler := newUI(t)
	if rec := get(t, handler, "/assets/gone.js"); rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if rec := get(t, handler, "/logo.png"); rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestEnvConfigIsNeverCached(t *testing.T) {
	handler := newUI(t)
	rec := get(t, handler, "/env-config.js")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if cache := rec.Header().Get("Cache-Control"); cache != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", cache)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/javascript") {
		t.Fatalf("Content-Type = %q", ct)
	}

	body := strings.TrimSpace(rec.Body.String())
	const prefix, suffix = "window._env_ = ", ";"
	if !strings.HasPrefix(body, prefix) || !strings.HasSuffix(body, suffix) {
		t.Fatalf("env-config.js = %q", body)
	}
	var values map[string]string
	if err := json.Unmarshal([]byte(strings.TrimSuffix(strings.TrimPrefix(body, prefix), suffix)), &values); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if values["REACT_APP_TERRAKUBE_API_URL"] != "https://terrakube.example.com" {
		t.Fatalf("API URL = %q", values["REACT_APP_TERRAKUBE_API_URL"])
	}
	if values["REACT_APP_CLIENT_ID"] != "terrakube-ui" || values["REACT_APP_AUTHORITY"] != "https://auth.example.com/dex" {
		t.Fatalf("login settings = %#v", values)
	}
}

func TestNoDirectoryMeansNoHandler(t *testing.T) {
	if New(Deps{}) != nil {
		t.Fatal("New returned a handler without a directory")
	}
	if New(Deps{Dir: t.TempDir(), Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}) != nil {
		t.Fatal("New returned a handler for a directory without index.html")
	}
}
