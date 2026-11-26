package httpapi

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMetricsDisabledWithoutToken(t *testing.T) {
	handler := New(Deps{Metrics: NewMetrics()})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 when no token is configured", rec.Code)
	}
}

func TestMetricsRequiresToken(t *testing.T) {
	handler := New(Deps{Metrics: NewMetrics(), MetricsToken: "scrape"})

	cases := []struct {
		name   string
		header string
		want   int
	}{
		{"no token", "", http.StatusUnauthorized},
		{"wrong token", "Bearer nope", http.StatusUnauthorized},
		{"right token", "Bearer scrape", http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d", rec.Code, tc.want)
			}
			if tc.want == http.StatusOK && !strings.Contains(rec.Body.String(), "go_goroutines") {
				t.Errorf("exposition lacks the runtime collector:\n%s", rec.Body.String())
			}
		})
	}
}

func TestMetricsRecordsRequests(t *testing.T) {
	handler := New(Deps{Metrics: NewMetrics(), MetricsToken: "scrape"})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("healthz status = %d", rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	req.Header.Set("Authorization", "Bearer scrape")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	for _, want := range []string{
		`terrakube_http_requests_total{method="GET",route="/healthz",status="200"} 1`,
		"terrakube_http_request_duration_seconds_bucket",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("exposition lacks %q", want)
		}
	}
}

type fakeStats struct {
	jobs  map[string]int64
	depth int64
	hit   chan struct{}
	err   error
}

func (f *fakeStats) JobStatusCounts(context.Context) (map[string]int64, error) {
	defer func() { f.hit <- struct{}{} }()
	return f.jobs, f.err
}

func (f *fakeStats) QueueDepth(context.Context) (int64, error) {
	return f.depth, f.err
}

func TestSampleRefreshesGauges(t *testing.T) {
	m := NewMetrics()
	src := &fakeStats{jobs: map[string]int64{"completed": 4, "failed": 1}, depth: 7, hit: make(chan struct{}, 4)}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		m.Sample(ctx, slog.New(slog.NewTextHandler(io.Discard, nil)), src, time.Millisecond)
	}()

	<-src.hit
	cancel()
	<-done

	handler := m.Handler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := rec.Body.String()

	for _, want := range []string{
		`terrakube_jobs{status="completed"} 4`,
		`terrakube_jobs{status="failed"} 1`,
		`terrakube_queue_depth{stream="jobs"} 7`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("exposition lacks %q", want)
		}
	}

	// A status that disappears from the database disappears from the gauge.
	m.TrackJobStatuses(map[string]int64{"completed": 5})
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	if strings.Contains(rec.Body.String(), `terrakube_jobs{status="failed"}`) {
		t.Error("stale status series kept")
	}
}
