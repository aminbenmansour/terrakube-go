package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	vcsOrgID     = "21212121-2121-4121-8121-212121212121"
	vcsOtherID   = "22222222-2222-4222-8222-222222222222"
	vcsID        = "23232323-2323-4323-8323-232323232323"
	vcsForeignID = "24242424-2424-4424-8424-242424242424"
	vcsWSID      = "25252525-2525-4525-8525-252525252525"
)

// vcsStub answers only the provider methods; the embedded interface panics for
// anything else, which would mean the handler grew a dependency.
type vcsStub struct {
	store.Querier
	org        store.Organization
	providers  []store.Vc
	workspaces []store.Workspace
	created    *store.CreateVCSParams
	updated    *store.UpdateVCSParams
	deleted    string
	cleared    *string
}

func (s *vcsStub) GetOrganization(context.Context, string) (store.Organization, error) {
	return s.org, nil
}

func (s *vcsStub) ListVCSByOrganization(context.Context, string) ([]store.Vc, error) {
	return s.providers, nil
}

func (s *vcsStub) GetVCS(_ context.Context, id string) (store.Vc, error) {
	for _, provider := range s.providers {
		if provider.ID == id {
			return provider, nil
		}
	}
	return store.Vc{}, pgx.ErrNoRows
}

func (s *vcsStub) CreateVCS(_ context.Context, arg store.CreateVCSParams) (store.Vc, error) {
	s.created = &arg
	return store.Vc{
		ID: arg.ID, Name: arg.Name, Description: arg.Description, VcsType: arg.VcsType,
		ClientID: arg.ClientID, ClientSecret: arg.ClientSecret, AccessToken: arg.AccessToken,
		Status: arg.Status, Callback: arg.Callback, RedirectUrl: arg.RedirectUrl,
		Endpoint: arg.Endpoint, ApiUrl: arg.ApiUrl, ConnectionType: arg.ConnectionType,
		PrivateKey: arg.PrivateKey, OrganizationID: arg.OrganizationID,
	}, nil
}

func (s *vcsStub) UpdateVCS(_ context.Context, arg store.UpdateVCSParams) (store.Vc, error) {
	s.updated = &arg
	return store.Vc{ID: arg.ID, Name: arg.Name, Description: arg.Description}, nil
}

func (s *vcsStub) ClearWorkspaceVCS(_ context.Context, id *string) error {
	s.cleared = id
	return nil
}

func (s *vcsStub) DeleteVCS(_ context.Context, id string) error {
	s.deleted = id
	return nil
}

func (s *vcsStub) ListWorkspaceIDsByVCS(context.Context, *string) ([]string, error) {
	return []string{vcsWSID}, nil
}

func (s *vcsStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *vcsStub) UpdateWorkspace(_ context.Context, arg store.UpdateWorkspaceParams) (store.Workspace, error) {
	return store.Workspace{ID: arg.ID, VcsID: arg.VcsID}, nil
}

func newVCSStub() *vcsStub {
	return &vcsStub{
		org: store.Organization{ID: vcsOrgID},
		providers: []store.Vc{{
			ID: vcsID, Name: "github", Description: "main provider", VcsType: "GITHUB",
			ClientID: strPtr("client"), ClientSecret: strPtr("secret"), AccessToken: strPtr("token"),
			PrivateKey: strPtr("key"), Status: strPtr("COMPLETED"), ConnectionType: strPtr("STANDALONE"),
			OrganizationID: vcsOrgID,
		}},
		workspaces: []store.Workspace{
			{ID: vcsWSID, OrganizationID: vcsOrgID, Name: "ws"},
			{ID: vcsForeignID, OrganizationID: vcsOtherID, Name: "other"},
		},
	}
}

func vcsTestAPI(q store.Querier) http.Handler {
	return New(Deps{Store: q})
}

func TestVCSJSONAPI(t *testing.T) {
	q := newVCSStub()
	h := vcsTestAPI(q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+vcsOrgID+"/vcs", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d; body = %s", rec.Code, rec.Body.String())
	}
	var list struct {
		Data []struct {
			Type       string         `json:"type"`
			Attributes map[string]any `json:"attributes"`
			Rel        map[string]any `json:"relationships"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 1 || list.Data[0].Type != "vcs" {
		t.Fatalf("list = %+v", list.Data)
	}
	for _, secret := range []string{"clientSecret", "accessToken", "refreshToken", "privateKey"} {
		if _, ok := list.Data[0].Attributes[secret]; ok {
			t.Errorf("list echoes %s", secret)
		}
	}
	if list.Data[0].Attributes["vcsType"] != "GITHUB" || list.Data[0].Attributes["status"] != "COMPLETED" {
		t.Errorf("attributes = %+v", list.Data[0].Attributes)
	}
	workspaceRel, _ := list.Data[0].Rel["workspace"].(map[string]any)
	if workspaceRel == nil {
		t.Errorf("workspace relationship missing: %+v", list.Data[0].Rel)
	}

	// The UI sends name, description and vcsType; the rest is optional.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+vcsOrgID+"/vcs",
		strings.NewReader(`{"data":{"type":"vcs","attributes":{"name":"gitlab","description":"gitlab","vcsType":"GITLAB"}}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.created == nil || q.created.Status == nil || *q.created.Status != "PENDING" ||
		q.created.ConnectionType == nil || *q.created.ConnectionType != "OAUTH" {
		t.Errorf("create defaults = %+v", q.created)
	}

	// Upstream's NOT NULL columns surface as Elide's data-integrity 409.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+vcsOrgID+"/vcs",
		strings.NewReader(`{"data":{"type":"vcs","attributes":{"name":"missing"}}}`)))
	if rec.Code != http.StatusConflict {
		t.Errorf("missing required = %d, want 409", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+vcsOrgID+"/vcs",
		strings.NewReader(`{"data":{"type":"vcs","attributes":{"name":"bad","description":"bad","vcsType":"NOPE"}}}`)))
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: NOPE") {
		t.Errorf("invalid type = %d; body = %s", rec.Code, rec.Body.String())
	}

	// A patch that leaves the secrets out keeps the stored ones.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/organization/"+vcsOrgID+"/vcs/"+vcsID,
		strings.NewReader(`{"data":{"type":"vcs","id":"`+vcsID+`","attributes":{"description":"renamed"}}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.updated == nil || q.updated.Description != "renamed" ||
		q.updated.ClientSecret == nil || *q.updated.ClientSecret != "secret" ||
		q.updated.PrivateKey == nil || *q.updated.PrivateKey != "key" {
		t.Errorf("patch lost secrets: %+v", q.updated)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/organization/"+vcsOrgID+"/vcs/"+vcsID,
		strings.NewReader(`{"data":{"type":"vcs","id":"`+vcsID+`","attributes":{"name":null}}}`)))
	if rec.Code != http.StatusConflict {
		t.Errorf("null name = %d, want 409", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/organization/"+vcsOrgID+"/vcs/"+vcsID,
		strings.NewReader(`{"data":{"type":"vcs","id":"`+vcsID+`","attributes":{"connectionType":"NOPE"}}}`)))
	if rec.Code != http.StatusBadRequest {
		t.Errorf("invalid connection type = %d, want 400", rec.Code)
	}

	// Deleting unlinks its workspaces first, so no dangling vcs id survives.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+vcsOrgID+"/vcs/"+vcsID, nil))
	if rec.Code != http.StatusNoContent || q.deleted != vcsID {
		t.Fatalf("delete = %d, deleted %q", rec.Code, q.deleted)
	}
	if q.cleared == nil || *q.cleared != vcsID {
		t.Errorf("delete did not clear workspace links: %v", q.cleared)
	}

	other := newVCSStub()
	other.org = store.Organization{ID: vcsOtherID}
	rec = httptest.NewRecorder()
	vcsTestAPI(other).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+vcsOtherID+"/vcs/"+vcsID, nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("cross-org get = %d, want 404", rec.Code)
	}
}

func TestVCSWorkspaceLink(t *testing.T) {
	q := newVCSStub()
	q.providers = append(q.providers, store.Vc{ID: vcsForeignID, Name: "foreign", OrganizationID: vcsOtherID})
	h := vcsTestAPI(q)
	base := "/organization/" + vcsOrgID + "/workspace/" + vcsWSID

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, base+"/relationships/vcs",
		strings.NewReader(`{"data":{"type":"vcs","id":"`+vcsID+`"}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("link = %d; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, base+"/relationships/vcs",
		strings.NewReader(`{"data":{"type":"vcs","id":"00000000-0000-4000-8000-000000000095"}}`)))
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown link = %d, want 404", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, base+"/relationships/vcs",
		strings.NewReader(`{"data":{"type":"vcs","id":"`+vcsForeignID+`"}}`)))
	if rec.Code != http.StatusForbidden {
		t.Errorf("foreign link = %d, want 403", rec.Code)
	}

	// The same check guards a vcs relationship sent inside a workspace PATCH.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, base,
		strings.NewReader(`{"data":{"type":"workspace","id":"`+vcsWSID+`","relationships":{"vcs":{"data":{"type":"vcs","id":"00000000-0000-4000-8000-000000000095"}}}}}`)))
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown patch link = %d, want 404", rec.Code)
	}

	// The workspace page reads the provider out of include=vcs.
	q.workspaces[0].VcsID = strPtr(vcsID)
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, base+"?include=vcs", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"type":"vcs"`) {
		t.Errorf("include = %d; body = %s", rec.Code, rec.Body.String())
	}
}
