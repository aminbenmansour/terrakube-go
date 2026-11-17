package tfe

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// These tests pin the JSON:API envelope and the status codes the Terraform CLI
// depends on. They intentionally assert shapes, not implementation details, so
// refactors do not silently change the compatibility surface.

func assertContractJSONAPI(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, jsonapi.MediaType) {
		t.Fatalf("Content-Type = %q, want %s", ct, jsonapi.MediaType)
	}
	var doc struct {
		Data   json.RawMessage   `json:"data"`
		Errors []json.RawMessage `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("response is not JSON: %v; body = %s", err, rec.Body.String())
	}
	if len(doc.Data) == 0 && len(doc.Errors) == 0 {
		t.Fatalf("response has neither data nor errors: %s", rec.Body.String())
	}
}

func assertContractError(t *testing.T, rec *httptest.ResponseRecorder, status int) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, status, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, jsonapi.MediaType) {
		t.Fatalf("Content-Type = %q, want %s", ct, jsonapi.MediaType)
	}
	var doc struct {
		Errors []struct {
			Status string `json:"status"`
			Title  string `json:"title"`
			Detail string `json:"detail"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("error response is not JSON: %v; body = %s", err, rec.Body.String())
	}
	if len(doc.Errors) == 0 {
		t.Fatalf("missing errors array: %s", rec.Body.String())
	}
	if doc.Errors[0].Status != strconv.Itoa(status) {
		t.Errorf("error status = %q, want %q", doc.Errors[0].Status, strconv.Itoa(status))
	}
	if doc.Errors[0].Title == "" || doc.Errors[0].Detail == "" {
		t.Errorf("error object = %+v, want title and detail", doc.Errors[0])
	}
}

func TestContractDiscoveryAndPing(t *testing.T) {
	discovery := Discovery(Deps{PublicURL: "http://localhost:8080"})
	rec := httptest.NewRecorder()
	discovery.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/terraform.json", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("discovery status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("discovery Content-Type = %q", ct)
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode discovery: %v", err)
	}
	for _, key := range []string{"tfe.v2", "tfe.v2.1", "state.v2", "login.v1"} {
		if len(body[key]) == 0 {
			t.Errorf("discovery missing %s: %s", key, rec.Body.String())
		}
	}

	handler := newTestAPI(t, &stubQuerier{})
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("ping status = %d, want 404", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("ping body = %q, want empty", rec.Body.String())
	}
	if got := rec.Header().Get("TFP-API-Version"); !strings.HasPrefix(got, "2.5") {
		t.Errorf("TFP-API-Version = %q", got)
	}
	if got := rec.Header().Get("TFP-AppName"); got != "Terrakube" {
		t.Errorf("TFP-AppName = %q", got)
	}
}

func TestContractOrganizationEnvelope(t *testing.T) {
	q := &stubQuerier{org: store.Organization{ID: "org-1", Name: "acme"}}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	assertContractJSONAPI(t, rec)
	var doc struct {
		Data struct {
			Type       string         `json:"type"`
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Type != "organizations" || doc.Data.ID != "acme" {
		t.Errorf("organization = %+v", doc.Data)
	}
	if doc.Data.Attributes["name"] != "acme" {
		t.Errorf("attributes = %+v", doc.Data.Attributes)
	}
}

func TestContractWorkspaceAndRunEnvelopes(t *testing.T) {
	q := &stubQuerier{
		org:       store.Organization{ID: "org-1", Name: "acme"},
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
		content: store.Content{
			ID:          "cv-1",
			WorkspaceID: stringPtr("ws-1"),
			Status:      stringPtr("uploaded"),
		},
	}
	handler, _ := newTestAPIWithQueue(t, q, &stubJobQueue{})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/workspaces/dev", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("workspace status = %d; body = %s", rec.Code, rec.Body.String())
	}
	assertContractJSONAPI(t, rec)
	if !strings.Contains(rec.Body.String(), `"type":"workspaces"`) {
		t.Errorf("workspace type missing: %s", rec.Body.String())
	}

	body := `{"data":{"type":"runs","relationships":{"workspace":{"data":{"type":"workspaces","id":"ws-1"}},"configuration-version":{"data":{"type":"configuration-versions","id":"cv-1"}}}}}`
	req := httptest.NewRequest(http.MethodPost, "/runs", strings.NewReader(body))
	req.Header.Set("Content-Type", jsonapi.MediaType)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("run status = %d; body = %s", rec.Code, rec.Body.String())
	}
	assertContractJSONAPI(t, rec)
	if !strings.Contains(rec.Body.String(), `"type":"runs"`) || !strings.Contains(rec.Body.String(), `"id":"run-42"`) {
		t.Errorf("run envelope = %s", rec.Body.String())
	}
}

func TestContractErrorShapes(t *testing.T) {
	q := &stubQuerier{orgErr: pgx.ErrNoRows, workspaceErr: pgx.ErrNoRows}
	handler := newTestAPI(t, q)

	cases := []struct {
		name   string
		method string
		path   string
		body   string
		status int
	}{
		{name: "missing organization", method: http.MethodGet, path: "/organizations/missing", status: http.StatusNotFound},
		{name: "missing workspace", method: http.MethodGet, path: "/workspaces/missing/runs", status: http.StatusNotFound},
		{name: "invalid run relationships", method: http.MethodPost, path: "/runs", body: `{"data":{"type":"runs"}}`, status: http.StatusUnprocessableEntity},
		{name: "missing upload ticket", method: http.MethodPut, path: "/configuration-versions/missing", body: "x", status: http.StatusUnauthorized},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", jsonapi.MediaType)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assertContractError(t, rec, tc.status)
		})
	}
}

func TestContractAuthentication(t *testing.T) {
	handler := New(Deps{
		Store:     &stubQuerier{},
		AuthToken: "secret",
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme", nil))
	assertContractError(t, rec, http.StatusUnauthorized)
}
