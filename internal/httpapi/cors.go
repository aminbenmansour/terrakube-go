package httpapi

import (
	"net/http"
	"strings"
)

// CORS allows browser access to a surface from the given origins. An empty list
// keeps the historical permissive default, which is what a UI served by this
// process relies on; anything else is an allowlist.
//
// It is shared rather than applied per router because the React UI calls three
// surfaces on this origin — /api/v1, /graphql/api/v1 and /subscriptions — and a
// split-origin deployment needs the same answer from all of them.
func CORS(origins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		allowed[strings.TrimSpace(origin)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			switch {
			case len(allowed) == 0:
				w.Header().Set("Access-Control-Allow-Origin", "*")
			case allowed[origin]:
				// Echoing the origin keeps credentials and any future
				// Vary-sensitive caching correct.
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
