package vcs

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newFakeProvider serves the routes discovery calls, and fails the test when a
// request does not look like the one upstream sends.
func newFakeProvider(t *testing.T, routes map[string]func(http.ResponseWriter, *http.Request)) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, ok := routes[r.URL.Path]
		if !ok {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.String())
			http.NotFound(w, r)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	return server
}

func requireBearer(t *testing.T, r *http.Request, token string) {
	t.Helper()
	if got := r.Header.Get("Authorization"); got != "Bearer "+token {
		t.Errorf("Authorization = %q, want bearer token", got)
	}
}

func writeJSON(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}

func TestRejectedCredentialIsAUseful4xx(t *testing.T) {
	server := newFakeProvider(t, map[string]func(http.ResponseWriter, *http.Request){
		"/expired/user": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			writeJSON(w, `{"message":"Bad credentials"}`)
		},
	})
	_, err := Test(context.Background(),
		Connection{AllowLocal: true, Type: TypeGitHub, APIURL: server.URL + "/expired", AccessToken: "gh-token"})
	var discoveryErr *Error
	if !errors.As(err, &discoveryErr) {
		t.Fatalf("error = %v, want *Error", err)
	}
	if discoveryErr.Status != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", discoveryErr.Status)
	}
	if !strings.Contains(discoveryErr.Message, "Bad credentials") {
		t.Fatalf("message = %q, want the provider detail", discoveryErr.Message)
	}
}

func TestUnsupportedConnections(t *testing.T) {
	cases := []struct {
		name       string
		connection Connection
		status     int
	}{
		{"unknown type", Connection{AllowLocal: true, Type: "BITBUCKET", AccessToken: "x"}, http.StatusUnprocessableEntity},
		{"public", Connection{AllowLocal: true, Type: "PUBLIC", AccessToken: "x"}, http.StatusUnprocessableEntity},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Groups(context.Background(), testCase.connection)
			var discoveryErr *Error
			if !errors.As(err, &discoveryErr) || discoveryErr.Status != testCase.status {
				t.Fatalf("error = %v, want status %d", err, testCase.status)
			}
		})
	}
}
