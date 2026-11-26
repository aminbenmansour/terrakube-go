package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimiterAllowsBurstThenRefills(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	l := newLimiter(2, 3, clock) // two per second, three at once

	for i := range 3 {
		if ok, _ := l.allow("a", now); !ok {
			t.Fatalf("request %d refused inside the burst", i+1)
		}
	}
	ok, retry := l.allow("a", now)
	if ok {
		t.Fatal("fourth request allowed past the burst")
	}
	if retry != 500*time.Millisecond {
		t.Errorf("retry = %v, want 500ms", retry)
	}

	now = now.Add(500 * time.Millisecond)
	if ok, _ := l.allow("a", now); !ok {
		t.Error("request refused after half a refill")
	}
	if ok, _ := l.allow("a", now); ok {
		t.Error("second request allowed without another refill")
	}

	// Another caller has its own bucket.
	if ok, _ := l.allow("b", now); !ok {
		t.Error("independent caller refused")
	}
}

func TestLimiterDropsIdleBuckets(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	l := newLimiter(1, 1, clock)

	l.allow("gone", now)
	now = now.Add(11 * time.Minute)
	l.allow("fresh", now)

	if _, ok := l.buckets["gone"]; ok {
		t.Error("idle bucket was kept")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	handler := RateLimit(1, 2)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	send := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer same")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec
	}

	for i := range 2 {
		if rec := send(); rec.Code != http.StatusNoContent {
			t.Fatalf("request %d status = %d, want 204", i+1, rec.Code)
		}
	}
	rec := send()
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", rec.Code)
	}
	if got := rec.Header().Get("Retry-After"); got == "" || got == "0" {
		t.Errorf("Retry-After = %q, want a positive delay", got)
	}

	// Preflights and probes are never throttled.
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	preflight := httptest.NewRecorder()
	handler.ServeHTTP(preflight, req)
	if preflight.Code != http.StatusNoContent {
		t.Errorf("preflight status = %d, want 204", preflight.Code)
	}
}

func TestRateLimitDisabledByDefault(t *testing.T) {
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	if got := RateLimit(0, 0); got(next) == nil {
		t.Fatal("disabled limiter returned nil")
	}
}
