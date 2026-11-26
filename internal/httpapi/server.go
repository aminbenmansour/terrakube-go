// Package httpapi contains the HTTP surface shared by the Terrakube services:
// middleware, health endpoints and (from M1 onwards) the JSON:API routers.
package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/version"
)

// Pinger is implemented by the database pool and the Redis queue.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Deps are the collaborators the HTTP layer needs.
type Deps struct {
	Logger  *slog.Logger
	Service string
	Started time.Time
	DB      Pinger
	Redis   Pinger
	Storage storage.Backend
	// RateLimit caps requests per second per caller (a rate of zero disables
	// it) and RateBurst is the bucket size (zero uses the rate).
	RateLimit int
	RateBurst int
}

// New builds the base router. Domain routers are mounted on top of it.
func New(deps Deps) *chi.Mux {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.Started.IsZero() {
		deps.Started = time.Now()
	}

	r := chi.NewRouter()
	r.Use(SecurityHeaders())
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(RateLimit(deps.RateLimit, deps.RateBurst))
	r.Use(middleware.Recoverer)
	r.Use(requestLogger(deps.Logger))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "ok",
			"service": deps.Service,
		})
	})

	r.Get("/readyz", func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
		defer cancel()

		checks := map[string]string{}
		ready := true

		if deps.DB != nil {
			if err := deps.DB.Ping(ctx); err != nil {
				checks["database"] = err.Error()
				ready = false
			} else {
				checks["database"] = "ok"
			}
		}
		if deps.Redis != nil {
			if err := deps.Redis.Ping(ctx); err != nil {
				checks["redis"] = err.Error()
				ready = false
			} else {
				checks["redis"] = "ok"
			}
		}
		if deps.Storage != nil {
			if err := deps.Storage.Health(ctx); err != nil {
				checks["storage"] = err.Error()
				ready = false
			} else {
				checks["storage"] = "ok"
			}
		}

		status := http.StatusOK
		state := "ok"
		if !ready {
			status = http.StatusServiceUnavailable
			state = "not ready"
		}
		writeJSON(w, status, map[string]any{
			"status":  state,
			"service": deps.Service,
			"checks":  checks,
			"uptime":  time.Since(deps.Started).Round(time.Second).String(),
		})
	})

	r.Get("/version", func(w http.ResponseWriter, _ *http.Request) {
		info := version.Get()
		writeJSON(w, http.StatusOK, map[string]string{
			"service": deps.Service,
			"version": info.Version,
			"commit":  info.Commit,
			"date":    info.Date,
		})
	})

	return r
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration", time.Since(start).String(),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}
