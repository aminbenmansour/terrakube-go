package httpapi

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	handler := New(Deps{Service: "api"})

	cases := []struct {
		name     string
		path     string
		tls      bool
		forward  string
		wantHSTS bool
	}{
		{"plain http", "/healthz", false, "", false},
		{"terminator reports https", "/healthz", false, "https", true},
		{"direct https", "/healthz", true, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.tls {
				req.TLS = &tls.ConnectionState{}
			}
			if tc.forward != "" {
				req.Header.Set("X-Forwarded-Proto", tc.forward)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			want := map[string]string{
				"X-Content-Type-Options": "nosniff",
				"X-Frame-Options":        "DENY",
				"Referrer-Policy":        "strict-origin-when-cross-origin",
			}
			for header, value := range want {
				if got := rec.Header().Get(header); got != value {
					t.Errorf("%s = %q, want %q", header, got, value)
				}
			}
			if got := rec.Header().Get("Strict-Transport-Security"); (got != "") != tc.wantHSTS {
				t.Errorf("Strict-Transport-Security = %q, want set = %v", got, tc.wantHSTS)
			}
		})
	}
}
