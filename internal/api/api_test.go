package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// stubQuerier embeds the generated interface so tests only implement the
// methods they exercise. Unimplemented methods panic, which is what we want.
type stubQuerier struct {
	store.Querier
	orgs            []store.Organization
	org             store.Organization
	jobs            []store.Job
	steps           []store.Step
	workspaces      []store.Workspace
	projects        []store.Project
	globals         []store.Globalvar
	actions         []store.Action
	collections     []store.Collection
	collectionItems []store.Item
	collectionRefs  []store.Reference
	sshKeys         []store.Ssh
	agents          []store.Agent
	sshKeyUses      map[string]int64
	agentUses       map[string]int64
	tags            []store.Tag
	workspaceTags   []store.Workspacetag
	err             error
	created         *store.CreateOrganizationParams
	createdAction   *store.CreateActionParams
	updatedAction   *store.UpdateActionParams
	deletedAction   string
}

func (s *stubQuerier) ListOrganizations(context.Context) ([]store.Organization, error) {
	return s.orgs, s.err
}

func (s *stubQuerier) GetOrganization(context.Context, string) (store.Organization, error) {
	return s.org, s.err
}

func (s *stubQuerier) CreateOrganization(_ context.Context, arg store.CreateOrganizationParams) (store.Organization, error) {
	s.created = &arg
	return s.org, s.err
}

func newTestAPI(t *testing.T, q store.Querier) http.Handler {
	t.Helper()
	return New(Deps{
		Store:  q,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
}

func TestListOrganizationsJSONAPI(t *testing.T) {
	q := &stubQuerier{orgs: []store.Organization{
		{ID: "1", Name: "alpha", ExecutionMode: "remote"},
		{ID: "2", Name: "beta", ExecutionMode: "remote"},
	}}
	handler := newTestAPI(t, q)

	req := httptest.NewRequest(http.MethodGet, "/organization?sort=-name&page[size]=1&page[number]=1&page[totals]", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	var doc struct {
		Data []struct {
			Type       string         `json:"type"`
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
		Meta struct {
			Page struct {
				Number       int `json:"number"`
				Limit        int `json:"limit"`
				TotalPages   int `json:"totalPages"`
				TotalRecords int `json:"totalRecords"`
			} `json:"page"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 1 {
		t.Fatalf("data length = %d, want 1", len(doc.Data))
	}
	if doc.Data[0].Attributes["name"] != "beta" {
		t.Errorf("first name = %v, want beta (descending sort)", doc.Data[0].Attributes["name"])
	}
	if doc.Meta.Page.TotalRecords != 2 || doc.Meta.Page.TotalPages != 2 || doc.Meta.Page.Limit != 1 || doc.Meta.Page.Number != 1 {
		t.Errorf("meta.page = %+v, want number 1, limit 1, totalPages 2, totalRecords 2", doc.Meta.Page)
	}
}

func TestGetOrganizationNotFound(t *testing.T) {
	handler := newTestAPI(t, &stubQuerier{err: pgx.ErrNoRows})

	req := httptest.NewRequest(http.MethodGet, "/organization/00000000-0000-0000-0000-000000000000", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	var doc struct {
		Errors []struct {
			Status string `json:"status"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Errors) != 1 || doc.Errors[0].Status != "404" {
		t.Errorf("errors = %+v, want one 404", doc.Errors)
	}
}

func TestCreateOrganizationRequiresName(t *testing.T) {
	q := &stubQuerier{}
	handler := newTestAPI(t, q)

	body := `{"data":{"type":"organization","attributes":{"description":"no name"}}}`
	req := httptest.NewRequest(http.MethodPost, "/organization", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body = %s", rec.Code, rec.Body.String())
	}
	if q.created != nil {
		t.Error("CreateOrganization was called for an invalid request")
	}
}

func TestCreateOrganization(t *testing.T) {
	q := &stubQuerier{org: store.Organization{ID: "new-id", Name: "simple", ExecutionMode: "remote"}}
	handler := New(Deps{
		Store:  q,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	body := `{"data":{"type":"organization","attributes":{"name":"simple"}}}`
	req := httptest.NewRequest(http.MethodPost, "/organization", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if q.created == nil {
		t.Fatal("CreateOrganization was not called")
	}
	if q.created.ExecutionMode != "remote" {
		t.Errorf("ExecutionMode = %q, want remote", q.created.ExecutionMode)
	}
	if q.created.ID == "" {
		t.Error("ID is empty")
	}
}

func TestAuthMiddleware(t *testing.T) {
	q := &stubQuerier{}
	handler := New(Deps{
		Store:     q,
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken: "secret",
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/organization", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("with token status = %d, want 200", rec.Code)
	}
}

// Upstream rejects page parameters on organization (it cannot paginate that
// collection) and answers an unknown path with a JSON:API document. Both were
// captured from upstream 2.34.3; we are permissive about the first and match the
// second, because a client that cannot parse a 404 body shows the wrong message.
func TestUnknownCollectionIsJSONAPI(t *testing.T) {
	handler := newTestAPI(t, &stubQuerier{})

	req := httptest.NewRequest(http.MethodGet, "/nope", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	var doc struct {
		Errors []struct {
			Detail string `json:"detail"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v; body = %s", err, rec.Body.String())
	}
	if len(doc.Errors) != 1 || doc.Errors[0].Detail != "Unknown collection nope" {
		t.Fatalf("errors = %+v, want one detail 'Unknown collection nope'", doc.Errors)
	}
}

func TestPageNumberWithoutSizeUsesElideDefaultLimit(t *testing.T) {
	q := &stubQuerier{orgs: []store.Organization{
		{ID: "1", Name: "alpha"},
		{ID: "2", Name: "beta"},
	}}
	handler := newTestAPI(t, q)

	req := httptest.NewRequest(http.MethodGet, "/organization?page[number]=2", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var doc struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
		Meta struct {
			Page struct {
				Number int `json:"number"`
				Limit  int `json:"limit"`
			} `json:"page"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Meta.Page.Number != 2 || doc.Meta.Page.Limit != 10000 {
		t.Errorf("meta.page = %+v, want number 2, limit 10000", doc.Meta.Page)
	}
	if len(doc.Data) != 0 {
		t.Errorf("data = %v, want empty second page at the default limit", doc.Data)
	}
}

func TestMalformedIDIsBadRequest(t *testing.T) {
	handler := newTestAPI(t, &stubQuerier{})

	req := httptest.NewRequest(http.MethodGet, "/organization/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Errors []struct {
			Detail string `json:"detail"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Errors) != 1 || doc.Errors[0].Detail != "Invalid value: not-a-uuid" {
		t.Fatalf("errors = %+v, want upstream's 'Invalid value: not-a-uuid'", doc.Errors)
	}
}

// The router is mounted under MountPrefix in production, and chi leaves the
// prefix in URL.Path while routing against the stripped path. Calling New()
// directly in the other tests never exercised that.
func TestMountPrefixDoesNotConfuseTheIDCheck(t *testing.T) {
	mux := chi.NewMux()
	mux.Mount(MountPrefix, newTestAPI(t, &stubQuerier{orgs: []store.Organization{{ID: "1", Name: "azure"}}}))

	for _, tc := range []struct {
		path string
		want int
	}{
		{MountPrefix + "/organization", http.StatusOK},
		{MountPrefix + "/organization?filter[organization]=name==azure", http.StatusOK},
		{MountPrefix + "/organization/not-a-uuid", http.StatusBadRequest},
		{MountPrefix + "/nope", http.StatusNotFound},
	} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.want {
			t.Errorf("%s: status = %d, want %d; body = %s", tc.path, rec.Code, tc.want, rec.Body.String())
		}
		if tc.want == http.StatusNotFound && !strings.Contains(rec.Body.String(), "Unknown collection nope") {
			t.Errorf("%s: body = %s, want the collection blamed, not the mount prefix", tc.path, rec.Body.String())
		}
	}
}
