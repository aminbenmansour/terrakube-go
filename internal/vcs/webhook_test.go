package vcs

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fakeHookServer serves one provider request and records it.
type fakeHookServer struct {
	*httptest.Server
	method string
	path   string
	body   map[string]any
	auth   string
	status int
}

func newHookServer(t *testing.T, response string, status int) *fakeHookServer {
	t.Helper()
	fake := &fakeHookServer{status: status}
	fake.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fake.method = r.Method
		fake.path = r.URL.EscapedPath()
		fake.auth = r.Header.Get("Authorization")
		raw, _ := io.ReadAll(r.Body)
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &fake.body); err != nil {
				t.Errorf("request body is not JSON: %v (%s)", err, raw)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(fake.status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(fake.Close)
	return fake
}

func TestGitHubWebhookCreate(t *testing.T) {
	server := newHookServer(t, `{"id": 4242}`, http.StatusCreated)

	remote, err := SyncWebhook(context.Background(),
		Connection{AllowLocal: true, Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token"},
		WebhookSpec{
			SourceURL: "https://github.com/acme/infra.git", CallbackURL: "https://terrakube.example.com/webhook/v1/hook-id",
			Secret: "c2VjcmV0", Events: []string{"PUSH", "PULL_REQUEST", "PUSH"}, PullRequestWorkflow: true,
		}, "")
	if err != nil {
		t.Fatalf("sync webhook: %v", err)
	}
	if remote != "4242" {
		t.Fatalf("remote id = %q", remote)
	}
	if server.method != http.MethodPost || server.path != "/repos/acme/infra/hooks" {
		t.Fatalf("request = %s %s", server.method, server.path)
	}
	if server.auth != "Bearer gh-token" {
		t.Fatalf("Authorization = %q", server.auth)
	}
	if server.body["name"] != "web" || server.body["active"] != true {
		t.Fatalf("body = %+v", server.body)
	}
	events, _ := server.body["events"].([]any)
	if len(events) != 3 || events[0] != "push" || events[1] != "pull_request" || events[2] != "issue_comment" {
		t.Fatalf("events = %+v, want deduplicated lower-case events plus issue_comment", events)
	}
	config, _ := server.body["config"].(map[string]any)
	if config["url"] != "https://terrakube.example.com/webhook/v1/hook-id" || config["secret"] != "c2VjcmV0" ||
		config["content_type"] != "json" || config["insecure_ssl"] != "1" {
		t.Fatalf("config = %+v", config)
	}
}

func TestGitHubWebhookUpdateAndDelete(t *testing.T) {
	server := newHookServer(t, `{"id": 4242}`, http.StatusOK)
	connection := Connection{AllowLocal: true, Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token"}
	spec := WebhookSpec{SourceURL: "https://github.com/acme/infra", Events: []string{"RELEASE"}}

	remote, err := SyncWebhook(context.Background(), connection, spec, "4242")
	if err != nil {
		t.Fatalf("update webhook: %v", err)
	}
	if remote != "4242" || server.method != http.MethodPatch || server.path != "/repos/acme/infra/hooks/4242" {
		t.Fatalf("update = %q %s %s", remote, server.method, server.path)
	}
	if _, ok := server.body["config"]; ok {
		t.Fatalf("update body = %+v, want only active and events", server.body)
	}

	server.status = http.StatusNoContent
	if err := DeleteWebhook(context.Background(), connection, spec.SourceURL, "4242"); err != nil {
		t.Fatalf("delete webhook: %v", err)
	}
	if server.method != http.MethodDelete || server.path != "/repos/acme/infra/hooks/4242" {
		t.Fatalf("delete = %s %s", server.method, server.path)
	}
}

func TestDeleteWebhookTreats404AsRemoved(t *testing.T) {
	server := newHookServer(t, `{"message":"Not Found"}`, http.StatusNotFound)
	if err := DeleteWebhook(context.Background(),
		Connection{AllowLocal: true, Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token"},
		"https://github.com/acme/infra", "9"); err != nil {
		t.Fatalf("delete gone webhook = %v", err)
	}
}

func TestGitLabWebhookCreate(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/projects/acme%2Finfra":
			_, _ = w.Write([]byte(`{"id": 7}`))
		case r.Method == http.MethodPost && r.URL.EscapedPath() == "/projects/7/hooks":
			raw, _ := io.ReadAll(r.Body)
			var body map[string]any
			if err := json.Unmarshal(raw, &body); err != nil {
				t.Errorf("body is not JSON: %v", err)
			}
			if body["token"] != "c2VjcmV0" || body["url"] != "https://terrakube.example.com/webhook/v1/hook-id" ||
				body["push_events"] != true || body["merge_requests_events"] != true ||
				body["releases_events"] != true || body["note_events"] != true || body["enable_ssl_verification"] != false {
				t.Errorf("body = %+v", body)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": 99}`))
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.String())
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	remote, err := SyncWebhook(context.Background(),
		Connection{AllowLocal: true, Type: TypeGitLab, APIURL: server.URL, AccessToken: "gl-token"},
		WebhookSpec{
			SourceURL: "https://gitlab.com/acme/infra.git", CallbackURL: "https://terrakube.example.com/webhook/v1/hook-id",
			Secret: "c2VjcmV0", Events: []string{"PUSH"}, PullRequestWorkflow: true,
		}, "")
	if err != nil {
		t.Fatalf("sync webhook: %v", err)
	}
	if remote != "99" || calls != 2 {
		t.Fatalf("remote = %q after %d calls", remote, calls)
	}
}

func TestGitLabWebhookDelete(t *testing.T) {
	// Route both the project lookup and the delete through one server.
	handler := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.EscapedPath() == "/projects/acme%2Finfra":
			_, _ = w.Write([]byte(`{"id": 7}`))
		case r.Method == http.MethodDelete && r.URL.EscapedPath() == "/projects/7/hooks/99":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.String())
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(handler.Close)

	if err := DeleteWebhook(context.Background(),
		Connection{AllowLocal: true, Type: TypeGitLab, APIURL: handler.URL, AccessToken: "gl-token"},
		"https://gitlab.com/acme/infra", "99"); err != nil {
		t.Fatalf("delete webhook: %v", err)
	}
}

func TestWebhookProviderErrorCarriesStatus(t *testing.T) {
	server := newHookServer(t, `{"message":"Resource not accessible by integration"}`, http.StatusForbidden)
	_, err := SyncWebhook(context.Background(),
		Connection{AllowLocal: true, Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token"},
		WebhookSpec{SourceURL: "https://github.com/acme/infra"}, "")
	var providerErr *Error
	if !errors.As(err, &providerErr) || providerErr.Status != http.StatusForbidden {
		t.Fatalf("error = %v, want a 403 provider error", err)
	}
	if !strings.Contains(providerErr.Message, "scopes") {
		t.Fatalf("message = %q, want the scope hint", providerErr.Message)
	}
}
