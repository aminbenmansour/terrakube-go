package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// workspaceStub records what the handlers ask the store to write. The project
// relationship is the point: a stub that ignored ProjectID would let the tests
// pass while nothing moves.
type workspaceStub struct {
	store.Querier
	orgs       []store.Organization
	workspaces []store.Workspace
	projects   []store.Project
	created    *store.CreateWorkspaceParams
	updated    *store.UpdateWorkspaceParams
}

func (s *workspaceStub) GetOrganization(_ context.Context, id string) (store.Organization, error) {
	for _, org := range s.orgs {
		if org.ID == id {
			return org, nil
		}
	}
	return store.Organization{}, pgx.ErrNoRows
}

func (s *workspaceStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *workspaceStub) GetProject(_ context.Context, id string) (store.Project, error) {
	for _, project := range s.projects {
		if project.ID == id {
			return project, nil
		}
	}
	return store.Project{}, pgx.ErrNoRows
}

func (s *workspaceStub) CreateWorkspace(_ context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
	s.created = &arg
	return store.Workspace{
		ID:             arg.ID,
		Name:           arg.Name,
		OrganizationID: arg.OrganizationID,
		ProjectID:      arg.ProjectID,
	}, nil
}

func (s *workspaceStub) UpdateWorkspace(_ context.Context, arg store.UpdateWorkspaceParams) (store.Workspace, error) {
	s.updated = &arg
	for _, ws := range s.workspaces {
		if ws.ID == arg.ID {
			ws.ProjectID = arg.ProjectID
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

const (
	relOrgID     = "70000000-0000-4000-8000-000000000001"
	relWorkspace = "70000000-0000-4000-8000-000000000002"
	relProjectA  = "70000000-0000-4000-8000-000000000010"
	relProjectB  = "70000000-0000-4000-8000-000000000020"
	foreignOrg   = "70000000-0000-4000-8000-000000000099"
	foreignProj  = "70000000-0000-4000-8000-000000000098"
	unknownProj  = "00000000-0000-0000-0000-000000000000"
)

func newWorkspaceStub() *workspaceStub {
	return &workspaceStub{
		orgs: []store.Organization{{ID: relOrgID, Name: "simple-governance"}},
		workspaces: []store.Workspace{{
			ID: relWorkspace, OrganizationID: relOrgID, Name: "capture-rel-workspace",
			ProjectID: strPtr(relProjectA),
		}},
		projects: []store.Project{
			{ID: relProjectA, Name: "governance-terraform", OrganizationID: relOrgID},
			{ID: relProjectB, Name: "governance-tofu", OrganizationID: relOrgID},
			{ID: foreignProj, Name: "demo-project", OrganizationID: foreignOrg},
		},
	}
}

func serveWorkspace(t *testing.T, s *workspaceStub, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	newTestAPI(t, s).ServeHTTP(rec, req)
	return rec
}

func TestCreateWorkspaceStoresTheProjectRelationship(t *testing.T) {
	s := newWorkspaceStub()
	rec := serveWorkspace(t, s, http.MethodPost, "/organization/"+relOrgID+"/workspace",
		`{"data":{"type":"workspace","attributes":{"name":"from-vcs"},
		  "relationships":{"project":{"data":{"type":"project","id":"`+relProjectA+`"}}}}}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if s.created == nil || s.created.ProjectID == nil || *s.created.ProjectID != relProjectA {
		t.Fatalf("created project = %v, want %s", s.created, relProjectA)
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got := string(doc.Data.Relationships["project"].Data); !strings.Contains(got, relProjectA) {
		t.Errorf("relationships.project.data = %s, want %s", got, relProjectA)
	}
}

func TestCreateWorkspaceRefusesAForeignProject(t *testing.T) {
	s := newWorkspaceStub()
	rec := serveWorkspace(t, s, http.MethodPost, "/organization/"+relOrgID+"/workspace",
		`{"data":{"type":"workspace","attributes":{"name":"from-vcs"},
		  "relationships":{"project":{"data":{"type":"project","id":"`+foreignProj+`"}}}}}`)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body.String())
	}
	if s.created != nil {
		t.Errorf("workspace was created with project %v, want refusal", s.created.ProjectID)
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Errors) == 0 || doc.Errors[0].Detail != "Project does not belong to the workspace organization" {
		t.Errorf("errors = %v, want upstream's foreign-project detail", doc.Errors)
	}
}

func TestCreateWorkspaceReportsAnUnknownProject(t *testing.T) {
	rec := serveWorkspace(t, newWorkspaceStub(), http.MethodPost, "/organization/"+relOrgID+"/workspace",
		`{"data":{"type":"workspace","attributes":{"name":"from-vcs"},
		  "relationships":{"project":{"data":{"type":"project","id":"`+unknownProj+`"}}}}}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	want := "Unknown identifier " + unknownProj + " for project"
	if len(doc.Errors) == 0 || doc.Errors[0].Detail != want {
		t.Errorf("errors = %v, want %q", doc.Errors, want)
	}
}

func TestPatchWorkspaceMovesTheProject(t *testing.T) {
	s := newWorkspaceStub()
	rec := serveWorkspace(t, s, http.MethodPatch, "/organization/"+relOrgID+"/workspace/"+relWorkspace,
		`{"data":{"type":"workspace","id":"`+relWorkspace+`",
		  "relationships":{"project":{"data":{"type":"project","id":"`+relProjectB+`"}}}}}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
	if s.updated == nil || s.updated.ProjectID == nil || *s.updated.ProjectID != relProjectB {
		t.Fatalf("updated project = %v, want %s", s.updated, relProjectB)
	}
}

func TestPatchWorkspaceClearsTheProject(t *testing.T) {
	s := newWorkspaceStub()
	rec := serveWorkspace(t, s, http.MethodPatch, "/organization/"+relOrgID+"/workspace/"+relWorkspace,
		`{"data":{"type":"workspace","id":"`+relWorkspace+`",
		  "relationships":{"project":{"data":null}}}}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if s.updated == nil || s.updated.ProjectID != nil {
		t.Errorf("updated project = %v, want nil", s.updated)
	}
}

func TestPatchWorkspaceRefusesAForeignProject(t *testing.T) {
	s := newWorkspaceStub()
	rec := serveWorkspace(t, s, http.MethodPatch, "/organization/"+relOrgID+"/workspace/"+relWorkspace,
		`{"data":{"type":"workspace","id":"`+relWorkspace+`",
		  "relationships":{"project":{"data":{"type":"project","id":"`+foreignProj+`"}}}}}`)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body.String())
	}
	if s.updated != nil {
		t.Errorf("workspace was updated with project %v, want refusal", s.updated.ProjectID)
	}
}

// Upstream answers attribute-only patches with 204 too; the UI refetches.
func TestPatchWorkspaceAnswersNoContent(t *testing.T) {
	rec := serveWorkspace(t, newWorkspaceStub(), http.MethodPatch,
		"/organization/"+relOrgID+"/workspace/"+relWorkspace,
		`{"data":{"type":"workspace","id":"`+relWorkspace+`","attributes":{"description":"after"}}}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

// The UI's create form sends vcs, ssh and agent relationships through the
// atomic endpoint; the same fields must survive a direct POST.
func TestCreateWorkspaceStoresLinksAndSSHSettings(t *testing.T) {
	s := newWorkspaceStub()
	const vcs, ssh, agent = "70000000-0000-4000-8000-000000000031", "70000000-0000-4000-8000-000000000032", "70000000-0000-4000-8000-000000000033"
	rec := serveWorkspace(t, s, http.MethodPost, "/organization/"+relOrgID+"/workspace",
		`{"data":{"type":"workspace","attributes":{"name":"from-vcs","moduleSshKey":"`+ssh+`","sharedIds":"`+relWorkspace+`"},
		  "relationships":{"vcs":{"data":{"type":"vcs","id":"`+vcs+`"}},
		                   "ssh":{"data":{"type":"ssh","id":"`+ssh+`"}},
		                   "agent":{"data":{"type":"agent","id":"`+agent+`"}}}}}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if s.created == nil {
		t.Fatal("no create recorded")
	}
	if s.created.VcsID == nil || *s.created.VcsID != vcs || s.created.SshID == nil || *s.created.SshID != ssh ||
		s.created.AgentID == nil || *s.created.AgentID != agent {
		t.Errorf("links = %v/%v/%v, want %s/%s/%s", s.created.VcsID, s.created.SshID, s.created.AgentID, vcs, ssh, agent)
	}
	if s.created.ModuleSshKey == nil || *s.created.ModuleSshKey != ssh {
		t.Errorf("moduleSshKey = %v, want %s", s.created.ModuleSshKey, ssh)
	}
	if s.created.SharedIds == nil || *s.created.SharedIds != relWorkspace {
		t.Errorf("sharedIds = %v, want %s", s.created.SharedIds, relWorkspace)
	}
}

func TestPatchWorkspaceRelationshipSetsAndClearsLinks(t *testing.T) {
	const vcs = "70000000-0000-4000-8000-000000000031"

	s := newWorkspaceStub()
	rec := serveWorkspace(t, s, http.MethodPatch,
		"/organization/"+relOrgID+"/workspace/"+relWorkspace+"/relationships/vcs",
		`{"data":{"type":"vcs","id":"`+vcs+`"}}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if s.updated == nil || s.updated.VcsID == nil || *s.updated.VcsID != vcs {
		t.Fatalf("updated vcs = %v, want %s", s.updated, vcs)
	}

	rec = serveWorkspace(t, s, http.MethodPatch,
		"/organization/"+relOrgID+"/workspace/"+relWorkspace+"/relationships/vcs",
		`{"data":null}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if s.updated.VcsID != nil {
		t.Errorf("updated vcs = %v, want cleared", *s.updated.VcsID)
	}
}

func TestPatchWorkspaceRelationshipValidatesProject(t *testing.T) {
	s := newWorkspaceStub()
	rec := serveWorkspace(t, s, http.MethodPatch,
		"/organization/"+relOrgID+"/workspace/"+relWorkspace+"/relationships/project",
		`{"data":{"type":"project","id":"`+foreignProj+`"}}`)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body.String())
	}

	rec = serveWorkspace(t, s, http.MethodPatch,
		"/organization/"+relOrgID+"/workspace/"+relWorkspace+"/relationships/project",
		`{"data":{"type":"project","id":"`+unknownProj+`"}}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestPatchWorkspaceRelationshipRejectsUnknownAndMismatched(t *testing.T) {
	rec := serveWorkspace(t, newWorkspaceStub(), http.MethodPatch,
		"/organization/"+relOrgID+"/workspace/"+relWorkspace+"/relationships/nope",
		`{"data":null}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}

	rec = serveWorkspace(t, newWorkspaceStub(), http.MethodPatch,
		"/organization/"+relOrgID+"/workspace/"+relWorkspace+"/relationships/vcs",
		`{"data":{"type":"ssh","id":"`+relProjectA+`"}}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
}
