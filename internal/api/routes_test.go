package api

import (
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

var updateGolden = flag.Bool("update-golden", false, "rewrite testdata/routes.txt")

// TestManagementAPIRoutes pins the routes this API answers. The upstream UI
// needs far more than this (docs/upstream/ui-rest-operations.md); the list is
// here so a refactor cannot quietly drop one, and so the gap can be measured
// instead of guessed.
func TestManagementAPIRoutes(t *testing.T) {
	handler := New(Deps{
		Store:  &stubQuerier{},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	mux, ok := handler.(*chi.Mux)
	if !ok {
		t.Fatalf("New returned %T, want *chi.Mux", handler)
	}

	var routes []string
	if err := chi.Walk(mux, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		routes = append(routes, method+" "+route)
		return nil
	}); err != nil {
		t.Fatalf("walk routes: %v", err)
	}
	sort.Strings(routes)
	got := strings.Join(routes, "\n") + "\n"

	golden := filepath.Join("testdata", "routes.txt")
	if *updateGolden {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}

	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden (run with -update-golden to create it): %v", err)
	}
	if got != string(want) {
		t.Fatalf("routes changed; run `go test ./internal/api -update-golden` if that is intended\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
