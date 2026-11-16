package tfe

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// stubQuerier embeds the generated interface so tests only implement the
// methods they exercise.
type stubQuerier struct {
	store.Querier

	org    store.Organization
	orgErr error

	workspaces    []store.Workspace
	workspacesErr error
	workspace     store.Workspace
	workspaceErr  error
	createdWS     *store.CreateWorkspaceParams
	lockParams    *store.SetWorkspaceLockParams

	tags                     []store.Tag
	workspaceTags            []store.Workspacetag
	workspaceTagsByWorkspace map[string][]store.Workspacetag

	projects    []store.Project
	projectsErr error
	project     store.Project
	projectErr  error
	createdProj *store.CreateProjectParams
}

func (s *stubQuerier) GetOrganizationByName(context.Context, string) (store.Organization, error) {
	return s.org, s.orgErr
}

func (s *stubQuerier) ListWorkspacesByOrganization(context.Context, string) ([]store.Workspace, error) {
	return s.workspaces, s.workspacesErr
}

func (s *stubQuerier) GetWorkspaceByName(context.Context, store.GetWorkspaceByNameParams) (store.Workspace, error) {
	return s.workspace, s.workspaceErr
}

func (s *stubQuerier) GetWorkspace(context.Context, string) (store.Workspace, error) {
	return s.workspace, s.workspaceErr
}

func (s *stubQuerier) CreateWorkspace(_ context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
	s.createdWS = &arg
	return s.workspace, nil
}

func (s *stubQuerier) UpdateWorkspace(context.Context, store.UpdateWorkspaceParams) (store.Workspace, error) {
	return s.workspace, nil
}

func (s *stubQuerier) SetWorkspaceLock(_ context.Context, arg store.SetWorkspaceLockParams) (store.Workspace, error) {
	s.lockParams = &arg
	ws := s.workspace
	ws.Locked = arg.Locked
	ws.LockDescription = arg.LockDescription
	return ws, nil
}

func (s *stubQuerier) ListTagsByOrganization(context.Context, string) ([]store.Tag, error) {
	return s.tags, nil
}

func (s *stubQuerier) ListWorkspaceTagsByWorkspace(_ context.Context, workspaceID string) ([]store.Workspacetag, error) {
	if s.workspaceTagsByWorkspace != nil {
		return s.workspaceTagsByWorkspace[workspaceID], nil
	}
	return s.workspaceTags, nil
}

func (s *stubQuerier) ListProjectsByOrganization(context.Context, string) ([]store.Project, error) {
	return s.projects, s.projectsErr
}

func (s *stubQuerier) GetProject(context.Context, string) (store.Project, error) {
	return s.project, s.projectErr
}

func (s *stubQuerier) CreateProject(_ context.Context, arg store.CreateProjectParams) (store.Project, error) {
	s.createdProj = &arg
	return s.project, nil
}

func newTestAPI(t *testing.T, q store.Querier) http.Handler {
	t.Helper()
	return New(Deps{
		Store:  q,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
}

func TestDiscovery(t *testing.T) {
	handler := Discovery(Deps{
		PublicURL:    "http://localhost:8080",
		DexClientID:  "terrakube",
		DexIssuerURL: "http://localhost:5556/dex",
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/terraform.json", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var doc map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc["state.v2"] != "/remote/state/v2/" || doc["tfe.v2"] != "/remote/tfe/v2/" {
		t.Errorf("discovery paths = %#v", doc)
	}
	login, ok := doc["login.v1"].(map[string]any)
	if !ok {
		t.Fatalf("login.v1 = %#v", doc["login.v1"])
	}
	if login["client"] != "terrakube" {
		t.Errorf("login client = %v", login["client"])
	}
	if login["token"] != "http://localhost:5556/dex/token" {
		t.Errorf("login token = %v", login["token"])
	}
}

func TestPing(t *testing.T) {
	handler := newTestAPI(t, &stubQuerier{})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if got := rec.Header().Get("TFP-API-Version"); got != "2.5" {
		t.Errorf("TFP-API-Version = %q, want 2.5", got)
	}
	if got := rec.Header().Get("TFP-AppName"); got != "Terrakube" {
		t.Errorf("TFP-AppName = %q, want Terrakube", got)
	}
}

func TestAuth(t *testing.T) {
	q := &stubQuerier{org: store.Organization{ID: "org-1", Name: "acme", ExecutionMode: "remote"}}
	handler := New(Deps{
		Store:     q,
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken: "secret",
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/organizations/acme", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("with token status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
}

func TestGetOrganization(t *testing.T) {
	q := &stubQuerier{org: store.Organization{ID: "org-1", Name: "acme", ExecutionMode: "remote"}}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

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
		t.Errorf("data = %+v", doc.Data)
	}
	if doc.Data.Attributes["name"] != "acme" {
		t.Errorf("name = %v", doc.Data.Attributes["name"])
	}
}

func TestEntitlementSet(t *testing.T) {
	q := &stubQuerier{org: store.Organization{ID: "org-1", Name: "acme"}}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/entitlement-set", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

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
	if doc.Data.Type != "entitlement-sets" || doc.Data.ID != "org-acme" {
		t.Errorf("data = %+v", doc.Data)
	}
	if doc.Data.Attributes["state-storage"] != true {
		t.Errorf("state-storage = %v", doc.Data.Attributes["state-storage"])
	}
}

func TestListWorkspacesRequiresFilter(t *testing.T) {
	q := &stubQuerier{
		org: store.Organization{ID: "org-1", Name: "acme"},
		workspaces: []store.Workspace{
			{ID: "ws-1", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
		},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/workspaces", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var empty struct {
		Data []any          `json:"data"`
		Meta map[string]any `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &empty); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(empty.Data) != 0 {
		t.Errorf("data = %#v, want empty", empty.Data)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/workspaces?search[name]=dev", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("filtered status = %d, want 200", rec.Code)
	}
	var filtered struct {
		Data []struct {
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &filtered); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(filtered.Data) != 1 || filtered.Data[0].ID != "ws-1" {
		t.Fatalf("filtered data = %+v", filtered.Data)
	}
	if filtered.Data[0].Attributes["name"] != "dev" {
		t.Errorf("name = %v", filtered.Data[0].Attributes["name"])
	}
}

func TestListWorkspacesTagFilter(t *testing.T) {
	q := &stubQuerier{
		org: store.Organization{ID: "org-1", Name: "acme"},
		workspaces: []store.Workspace{
			{ID: "ws-1", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
			{ID: "ws-2", Name: "prod", OrganizationID: "org-1", GlobalRemoteState: true},
		},
		tags: []store.Tag{{ID: "tag-1", Name: "env", OrganizationID: "org-1"}},
		workspaceTagsByWorkspace: map[string][]store.Workspacetag{
			"ws-1": {{ID: "wt-1", TagID: "tag-1", WorkspaceID: "ws-1", TagValue: stringPtr("dev")}},
			"ws-2": {{ID: "wt-2", TagID: "tag-1", WorkspaceID: "ws-2", TagValue: stringPtr("prod")}},
		},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/workspaces?search%5Btags%5D=env&filter%5Btagged%5D%5B0%5D%5Bkey%5D=env&filter%5Btagged%5D%5B0%5D%5Bvalue%5D=dev", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data []struct {
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 1 || doc.Data[0].ID != "ws-1" {
		t.Fatalf("data = %+v, want ws-1", doc.Data)
	}
	names, ok := doc.Data[0].Attributes["tag-names"].([]any)
	if !ok || len(names) != 1 || names[0] != "env" {
		t.Errorf("tag-names = %#v", doc.Data[0].Attributes["tag-names"])
	}
}

func TestCreateWorkspace(t *testing.T) {
	q := &stubQuerier{
		org:          store.Organization{ID: "org-1", Name: "acme", ExecutionMode: "remote"},
		workspace:    store.Workspace{ID: "ws-new", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
		workspaceErr: pgx.ErrNoRows,
	}
	handler := newTestAPI(t, q)

	body := `{"data":{"type":"workspaces","attributes":{"name":"dev","terraform-version":"1.9.0","execution-mode":"remote"}}}`
	req := httptest.NewRequest(http.MethodPost, "/organizations/acme/workspaces", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if q.createdWS == nil {
		t.Fatal("CreateWorkspace was not called")
	}
	if q.createdWS.Name != "dev" || q.createdWS.OrganizationID != "org-1" {
		t.Errorf("created = %+v", q.createdWS)
	}
	if q.createdWS.TerraformVersion == nil || *q.createdWS.TerraformVersion != "1.9.0" {
		t.Errorf("terraform version = %v", q.createdWS.TerraformVersion)
	}
	if q.createdWS.Branch == nil || *q.createdWS.Branch != "remote-content" {
		t.Errorf("branch = %v", q.createdWS.Branch)
	}
}

func TestLockWorkspace(t *testing.T) {
	q := &stubQuerier{
		org:       store.Organization{ID: "org-1", Name: "acme"},
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/actions/lock", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if q.lockParams == nil || !q.lockParams.Locked {
		t.Fatalf("lock params = %+v", q.lockParams)
	}

	q.workspace.Locked = true
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/actions/lock", nil))
	if rec.Code != http.StatusConflict {
		t.Fatalf("already locked status = %d, want 409", rec.Code)
	}
}

func TestProjects(t *testing.T) {
	q := &stubQuerier{
		org:      store.Organization{ID: "org-1", Name: "acme"},
		projects: []store.Project{{ID: "proj-1", Name: "platform", OrganizationID: "org-1"}},
		project:  store.Project{ID: "proj-new", Name: "app", OrganizationID: "org-1"},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/projects", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200", rec.Code)
	}
	var list struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != "proj-1" {
		t.Errorf("projects = %+v", list.Data)
	}

	body := `{"data":{"type":"projects","attributes":{"name":"app"}}}`
	req := httptest.NewRequest(http.MethodPost, "/organizations/acme/projects", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if q.createdProj == nil || q.createdProj.Name != "app" {
		t.Errorf("created project = %+v", q.createdProj)
	}
}

func stringPtr(value string) *string {
	return &value
}
