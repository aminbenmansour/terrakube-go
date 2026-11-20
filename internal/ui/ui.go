// Package ui serves the upstream React UI from the API process. The UI is an
// unmodified upstream build: a single-page app that reads its runtime settings
// from /env-config.js and calls the API on its own origin.
package ui

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aminbenmansour/terrakube-go/internal/version"
)

// Deps are the collaborators the UI handler needs.
type Deps struct {
	// Dir is the directory holding the built UI (index.html and assets/).
	Dir string
	// PublicURL is this API's external URL; the UI derives every API URL from
	// it, including the WebSocket origin for subscriptions.
	PublicURL string
	// DexClientID and DexAuthority configure the OIDC login redirect.
	DexClientID  string
	DexAuthority string
	Logger       *slog.Logger
}

// New builds the UI handler. It returns nil when no UI directory is configured
// or the directory does not exist, so callers can mount it conditionally.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.Dir == "" {
		return nil
	}
	index := filepath.Join(deps.Dir, "index.html")
	if _, err := os.Stat(index); err != nil {
		deps.Logger.Warn("ui directory has no index.html; the UI is not served", "dir", deps.Dir, "error", err)
		return nil
	}

	h := &handler{
		deps:      deps,
		index:     index,
		envConfig: envConfig(deps),
	}
	files := http.FileServer(http.Dir(deps.Dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/env-config.js":
			// Regenerated per deployment; a cached copy points at a stale API.
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			_, _ = w.Write(h.envConfig)
		case strings.HasPrefix(r.URL.Path, "/assets/"):
			// Vite fingerprints build assets, so they never change in place.
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			files.ServeHTTP(w, r)
		case r.URL.Path == "/" || r.URL.Path == "/index.html":
			h.serveIndex(w, r)
		case exists(filepath.Join(deps.Dir, filepath.FromSlash(path.Clean(r.URL.Path)))):
			files.ServeHTTP(w, r)
		case path.Ext(r.URL.Path) != "":
			// A missing asset must 404 rather than return HTML, or the browser
			// reports a syntax error instead of a missing file.
			http.NotFound(w, r)
		default:
			// Deep links and the router's client-side paths.
			h.serveIndex(w, r)
		}
	})
}

type handler struct {
	deps      Deps
	index     string
	envConfig []byte
}

func (h *handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, h.index)
}

func exists(name string) bool {
	info, err := os.Stat(name)
	return err == nil && !info.IsDir()
}

// envConfig renders the window._env_ object the UI reads at startup. The keys
// are the upstream ones, so an upstream build works unchanged.
func envConfig(deps Deps) []byte {
	values := map[string]string{
		"REACT_APP_TERRAKUBE_API_URL": deps.PublicURL,
		"REACT_APP_CLIENT_ID":         deps.DexClientID,
		"REACT_APP_AUTHORITY":         deps.DexAuthority,
		"REACT_APP_REDIRECT_URI":      deps.PublicURL,
		"REACT_APP_SCOPE":             "openid profile email",
		"REACT_APP_BASE_PATH":         "/",
		"REACT_APP_REGISTRY_URI":      deps.PublicURL + "/registry",
		"REACT_APP_TERRAKUBE_VERSION": version.Get().Version,
	}
	// json.Marshal escapes < and >, which is what keeps a configured URL from
	// closing the script tag.
	encoded, err := json.Marshal(values)
	if err != nil {
		return []byte("window._env_ = {};\n")
	}
	return append([]byte("window._env_ = "), append(encoded, ";\n"...)...)
}
