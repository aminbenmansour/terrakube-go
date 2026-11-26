package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
