package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

func TestHealthz(t *testing.T) {
	handler := New(Deps{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), Service: "api"})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %v, want ok", body["status"])
	}
	if body["service"] != "api" {
		t.Errorf("service = %v, want api", body["service"])
	}
}

func TestReadyz(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("ready", func(t *testing.T) {
		handler := New(Deps{
			Logger:  logger,
			Service: "api",
			DB:      fakePinger{},
			Redis:   fakePinger{},
		})
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("database down", func(t *testing.T) {
		handler := New(Deps{
			Logger:  logger,
			Service: "api",
			DB:      fakePinger{err: context.DeadlineExceeded},
		})
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want 503", rec.Code)
		}
	})
}

func TestVersion(t *testing.T) {
	handler := New(Deps{
		Logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
		Service: "executor",
		Started: time.Now(),
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["service"] != "executor" {
		t.Errorf("service = %q, want executor", body["service"])
	}
	if body["version"] == "" {
		t.Error("version is empty")
	}
}
