package httpapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLimitBody(t *testing.T) {
	reached := false
	echo := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		w.Write(body)
	})
	handler := LimitBody(8)(echo)

	cases := []struct {
		name   string
		body   string
		length int64 // -1 leaves the length undeclared, as a chunked request does
		want   int
	}{
		{"under the limit", "12345678", 8, http.StatusOK},
		{"over the limit", "123456789", 9, http.StatusRequestEntityTooLarge},
		{"chunked under the limit", "12345678", -1, http.StatusOK},
		{"chunked over the limit", "123456789", -1, http.StatusRequestEntityTooLarge},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reached = false
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			req.ContentLength = tc.length
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, tc.want, rec.Body.String())
			}
			if tc.want == http.StatusOK && (!reached || rec.Body.String() != tc.body) {
				t.Errorf("handler got %q, want %q", rec.Body.String(), tc.body)
			}
			if tc.want == http.StatusRequestEntityTooLarge {
				if reached {
					t.Error("handler ran for an oversized body")
				}
				if !strings.Contains(rec.Body.String(), "exceeds the") {
					t.Errorf("body = %s, want a JSON:API 413 error", rec.Body.String())
				}
			}
		})
	}
}
