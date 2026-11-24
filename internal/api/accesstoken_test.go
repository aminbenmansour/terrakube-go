package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// fakeAccessTokenStore answers the four lookups from maps, so the handlers can
// be tested without a database.
type fakeAccessTokenStore struct {
	teams         []store.Team
	access        []store.Access
	projectAccess []store.ProjectAccess
	workspaces    map[string]store.Workspace
	lastTeamNames []string
	lastAccessID  string
	lastProjectID string
}

func (f *fakeAccessTokenStore) ListTeamsByOrganizationAndNames(_ context.Context, arg store.ListTeamsByOrganizationAndNamesParams) ([]store.Team, error) {
	f.lastTeamNames = arg.Names
	var out []store.Team
	for _, team := range f.teams {
		if team.OrganizationID == arg.OrganizationID && contains(arg.Names, team.Name) {
			out = append(out, team)
		}
	}
	return out, nil
}

func (f *fakeAccessTokenStore) ListAccessByWorkspaceAndNames(_ context.Context, arg store.ListAccessByWorkspaceAndNamesParams) ([]store.Access, error) {
	f.lastAccessID = arg.WorkspaceID
	var out []store.Access
	for _, access := range f.access {
		if access.WorkspaceID == arg.WorkspaceID && contains(arg.Names, access.Name) {
			out = append(out, access)
		}
	}
	return out, nil
}

func (f *fakeAccessTokenStore) ListProjectAccessByProjectAndNames(_ context.Context, arg store.ListProjectAccessByProjectAndNamesParams) ([]store.ProjectAccess, error) {
	f.lastProjectID = arg.ProjectID
	var out []store.ProjectAccess
	for _, access := range f.projectAccess {
		if access.ProjectID == arg.ProjectID && access.Name != nil && contains(arg.Names, *access.Name) {
			out = append(out, access)
		}
	}
	return out, nil
}

func (f *fakeAccessTokenStore) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	workspace, ok := f.workspaces[id]
	if !ok {
		return store.Workspace{}, pgx.ErrNoRows
	}
	return workspace, nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

type stubIdentityVerifier struct {
	token    string
	identity auth.Identity
}

func (v stubIdentityVerifier) Verify(_ context.Context, raw string) (*auth.Identity, error) {
	if raw != v.token {
		return nil, errors.New("unknown token")
	}
	identity := v.identity
	return &identity, nil
}

const (
	testOrgID       = "11111111-1111-4111-8111-111111111111"
	testOtherOrgID  = "22222222-2222-4222-8222-222222222222"
	testWorkspaceID = "33333333-3333-4333-8333-333333333333"
	testProjectID   = "44444444-4444-4444-8444-444444444444"
)

func serveAccessToken(t *testing.T, fake *fakeAccessTokenStore, verifier stubIdentityVerifier, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	router := NewAccessToken(Deps{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken:    "dev-token",
		Verifiers:    []auth.Verifier{verifier},
		AdminGroup:   "TERRAKUBE_ADMIN",
		AccessTokens: fake,
	})
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Authorization", "Bearer "+verifier.token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodePermissions(t *testing.T, body []byte) map[string]bool {
	t.Helper()
	var permissions map[string]bool
	if err := json.Unmarshal(body, &permissions); err != nil {
		t.Fatalf("decode permissions: %v (%s)", err, body)
	}
	return permissions
}

func TestAccessTokenCurrentTeams(t *testing.T) {
	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Groups: []string{"team-a", "team-b"}}}
	rec := serveAccessToken(t, &fakeAccessTokenStore{}, verifier, http.MethodGet, "/teams/current-teams")

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", rec.Code)
	}
	var body struct {
		Groups []string `json:"groups"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Groups) != 2 || body.Groups[0] != "team-a" || body.Groups[1] != "team-b" {
		t.Errorf("groups = %v, want the token's groups in order", body.Groups)
	}
}

func TestAccessTokenCurrentTeamsWithoutGroups(t *testing.T) {
	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Email: "user@example.com"}}
	rec := serveAccessToken(t, &fakeAccessTokenStore{}, verifier, http.MethodGet, "/teams/current-teams")

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	groups, ok := body["groups"].([]any)
	if !ok || len(groups) != 0 {
		t.Errorf("groups = %#v, want an empty list, not null", body["groups"])
	}
}

func TestAccessTokenOrganizationPermissions(t *testing.T) {
	fake := &fakeAccessTokenStore{teams: []store.Team{
		{Name: "readers", OrganizationID: testOrgID, Role: strPtr("read")},
		{Name: "writers", OrganizationID: testOrgID, Role: strPtr("write")},
		{Name: "foreign", OrganizationID: testOtherOrgID, Role: strPtr("admin")},
	}}
	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Groups: []string{"readers"}}}
	rec := serveAccessToken(t, fake, verifier, http.MethodGet, "/teams/permissions/organization/"+testOrgID)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", rec.Code)
	}
	permissions := decodePermissions(t, rec.Body.Bytes())
	if permissions["manageWorkspace"] || permissions["manageState"] || permissions["planJob"] ||
		permissions["manageModule"] || permissions["managePermission"] || permissions["managePolicies"] {
		t.Errorf("permissions = %v, want a read team to grant nothing", permissions)
	}
	if len(fake.lastTeamNames) != 1 || fake.lastTeamNames[0] != "readers" {
		t.Errorf("queried names = %v, want only the caller's groups", fake.lastTeamNames)
	}
}

func TestAccessTokenOrganizationPermissionsForTheOwner(t *testing.T) {
	owner := stubIdentityVerifier{token: "owner-token", identity: auth.Identity{Groups: []string{"TERRAKUBE_ADMIN"}}}
	rec := serveAccessToken(t, &fakeAccessTokenStore{}, owner,
		http.MethodGet, "/teams/permissions/organization/"+testOtherOrgID)

	permissions := decodePermissions(t, rec.Body.Bytes())
	if !permissions["managePermission"] || !permissions["managePolicies"] {
		t.Errorf("permissions = %v, want the owner flags even for an unknown organization", permissions)
	}
	if permissions["manageWorkspace"] || permissions["planJob"] {
		t.Errorf("permissions = %v, want the owner to add only the two flags", permissions)
	}
}

func TestAccessTokenWorkspacePermissions(t *testing.T) {
	projectID := testProjectID
	fake := &fakeAccessTokenStore{
		teams: []store.Team{{Name: "platform", OrganizationID: testOrgID, Role: strPtr("custom"), ManageWorkspace: true}},
		access: []store.Access{
			{Name: "deployers", WorkspaceID: testWorkspaceID, Role: strPtr("plan")},
			{Name: "strangers", WorkspaceID: testWorkspaceID, Role: strPtr("admin")},
		},
		projectAccess: []store.ProjectAccess{
			{Name: strPtr("deployers"), ProjectID: projectID, Role: strPtr("custom"), ManageState: true},
			{Name: strPtr("strangers"), ProjectID: projectID, Role: strPtr("admin")},
		},
		workspaces: map[string]store.Workspace{
			testWorkspaceID: {ID: testWorkspaceID, OrganizationID: testOrgID, ProjectID: &projectID},
		},
	}
	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Groups: []string{"platform", "deployers"}}}
	rec := serveAccessToken(t, fake, verifier, http.MethodGet,
		"/teams/permissions/organization/"+testOrgID+"/workspace/"+testWorkspaceID)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202: %s", rec.Code, rec.Body.String())
	}
	permissions := decodePermissions(t, rec.Body.Bytes())
	if !permissions["manageWorkspace"] || !permissions["planJob"] || !permissions["manageState"] {
		t.Errorf("permissions = %v, want the team plus the deployers rows", permissions)
	}
	if permissions["approveJob"] {
		t.Errorf("permissions = %v, want the plan access not to approve", permissions)
	}
	if fake.lastAccessID != testWorkspaceID || fake.lastProjectID != testProjectID {
		t.Errorf("lookups = %q/%q, want the workspace and its project", fake.lastAccessID, fake.lastProjectID)
	}
}

func TestAccessTokenWorkspacePermissionsForAMissingWorkspace(t *testing.T) {
	fake := &fakeAccessTokenStore{
		teams:      []store.Team{{Name: "platform", OrganizationID: testOrgID, Role: strPtr("read")}},
		workspaces: map[string]store.Workspace{},
	}
	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Groups: []string{"platform"}}}
	rec := serveAccessToken(t, fake, verifier, http.MethodGet,
		"/teams/permissions/organization/"+testOrgID+"/workspace/"+testWorkspaceID)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202 like upstream's optional lookup", rec.Code)
	}
	if fake.lastAccessID != "" {
		t.Error("a missing workspace must not trigger the access lookup")
	}
}

func TestAccessTokenRejectsBadInput(t *testing.T) {
	fake := &fakeAccessTokenStore{}

	unauthenticated := httptest.NewRecorder()
	NewAccessToken(Deps{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken:    "dev-token",
		AdminGroup:   "TERRAKUBE_ADMIN",
		AccessTokens: fake,
	}).ServeHTTP(unauthenticated, httptest.NewRequest(http.MethodGet, "/teams/current-teams", nil))
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Errorf("unauthenticated status = %d, want 401", unauthenticated.Code)
	}

	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Groups: []string{"team-a"}}}
	rec := serveAccessToken(t, fake, verifier, http.MethodGet, "/teams/permissions/organization/not-a-uuid")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("bad organization status = %d, want 400 (upstream answers 500)", rec.Code)
	}
}
