package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/teamtoken"
)

// fakeAccessTokenStore records the parameters it receives and returns rows by
// name only. The org/project filters live in the SQL, so this fake deliberately
// ignores them: a handler that queries the wrong parent leaks a foreign row and
// the tests catch it.
type fakeAccessTokenStore struct {
	teams            []store.Team
	access           []store.Access
	projectAccess    []store.ProjectAccess
	workspaces       map[string]store.Workspace
	lastTeamOrgID    string
	lastTeamNames    []string
	lastAccessID     string
	lastAccessNames  []string
	lastProjectID    string
	lastProjectNames []string
}

func (f *fakeAccessTokenStore) ListTeamsByOrganizationAndNames(_ context.Context, arg store.ListTeamsByOrganizationAndNamesParams) ([]store.Team, error) {
	f.lastTeamOrgID = arg.OrganizationID
	f.lastTeamNames = arg.Names
	var out []store.Team
	for _, team := range f.teams {
		if contains(arg.Names, team.Name) {
			out = append(out, team)
		}
	}
	return out, nil
}

func (f *fakeAccessTokenStore) ListAccessByWorkspaceAndNames(_ context.Context, arg store.ListAccessByWorkspaceAndNamesParams) ([]store.Access, error) {
	f.lastAccessID = arg.WorkspaceID
	f.lastAccessNames = arg.Names
	var out []store.Access
	for _, access := range f.access {
		if contains(arg.Names, access.Name) {
			out = append(out, access)
		}
	}
	return out, nil
}

func (f *fakeAccessTokenStore) ListProjectAccessByProjectAndNames(_ context.Context, arg store.ListProjectAccessByProjectAndNamesParams) ([]store.ProjectAccess, error) {
	f.lastProjectID = arg.ProjectID
	f.lastProjectNames = arg.Names
	var out []store.ProjectAccess
	for _, access := range f.projectAccess {
		if access.Name != nil && contains(arg.Names, *access.Name) {
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

func expectPermissions(t *testing.T, rec *httptest.ResponseRecorder, want map[string]bool) {
	t.Helper()
	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202: %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
	var permissions map[string]bool
	if err := json.Unmarshal(rec.Body.Bytes(), &permissions); err != nil {
		t.Fatalf("decode permissions: %v (%s)", err, rec.Body.String())
	}
	if len(permissions) != len(want) {
		t.Fatalf("field count = %d, want %d: %s", len(permissions), len(want), rec.Body.String())
	}
	for name, value := range want {
		if permissions[name] != value {
			t.Errorf("%s = %v, want %v", name, permissions[name], value)
		}
	}
}

// none is the all-false permission set, so a test only names the granted ones.
var none = map[string]bool{
	"manageState": false, "manageWorkspace": false, "manageModule": false,
	"manageProvider": false, "manageVcs": false, "manageTemplate": false,
	"manageCollection": false, "manageJob": false, "planJob": false,
	"approveJob": false, "managePermission": false, "managePolicies": false,
}

func with(grants map[string]bool) map[string]bool {
	permissions := map[string]bool{}
	for name, value := range none {
		permissions[name] = value
	}
	for name, value := range grants {
		permissions[name] = value
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
	expectPermissions(t, rec, none)
	if fake.lastTeamOrgID != testOrgID {
		t.Errorf("queried organization = %q, want %q", fake.lastTeamOrgID, testOrgID)
	}
	if len(fake.lastTeamNames) != 1 || fake.lastTeamNames[0] != "readers" {
		t.Errorf("queried names = %v, want only the caller's groups", fake.lastTeamNames)
	}
}

func TestAccessTokenOrganizationPermissionsWithoutGroups(t *testing.T) {
	fake := &fakeAccessTokenStore{teams: []store.Team{
		{Name: "readers", OrganizationID: testOrgID, Role: strPtr("admin")},
	}}
	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Email: "user@example.com"}}
	rec := serveAccessToken(t, fake, verifier, http.MethodGet, "/teams/permissions/organization/"+testOrgID)

	expectPermissions(t, rec, none)
	if fake.lastTeamOrgID != "" {
		t.Errorf("queried organization = %q, want no team lookup for a token without groups", fake.lastTeamOrgID)
	}
}

func TestAccessTokenOrganizationPermissionsForTheOwner(t *testing.T) {
	owner := stubIdentityVerifier{token: "owner-token", identity: auth.Identity{Groups: []string{"TERRAKUBE_ADMIN"}}}
	rec := serveAccessToken(t, &fakeAccessTokenStore{}, owner,
		http.MethodGet, "/teams/permissions/organization/"+testOtherOrgID)

	expectPermissions(t, rec, with(map[string]bool{"managePermission": true, "managePolicies": true}))
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
	expectPermissions(t, rec, with(map[string]bool{
		"manageWorkspace": true, "managePermission": true,
		"planJob": true, "manageJob": true, "manageState": true,
	}))
	if fake.lastAccessID != testWorkspaceID || fake.lastProjectID != testProjectID {
		t.Errorf("lookups = %q/%q, want the workspace and its project", fake.lastAccessID, fake.lastProjectID)
	}
	if len(fake.lastAccessNames) != 2 || len(fake.lastProjectNames) != 2 {
		t.Errorf("names = %v/%v, want the caller's groups", fake.lastAccessNames, fake.lastProjectNames)
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

	rec = serveAccessToken(t, fake, verifier, http.MethodGet,
		"/teams/permissions/organization/"+testOrgID+"/workspace/not-a-uuid")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("bad workspace status = %d, want 400 (upstream answers 500)", rec.Code)
	}
}

func TestAccessTokenWithoutStore(t *testing.T) {
	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Groups: []string{"team-a"}}}
	router := NewAccessToken(Deps{
		Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken:  "dev-token",
		Verifiers:  []auth.Verifier{verifier},
		AdminGroup: "TERRAKUBE_ADMIN",
	})

	current := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/teams/current-teams", nil)
	request.Header.Set("Authorization", "Bearer "+verifier.token)
	router.ServeHTTP(current, request)
	if current.Code != http.StatusAccepted {
		t.Errorf("current-teams without a store = %d, want 202", current.Code)
	}

	permissions := httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/teams/permissions/organization/"+testOrgID, nil)
	request.Header.Set("Authorization", "Bearer "+verifier.token)
	router.ServeHTTP(permissions, request)
	if permissions.Code != http.StatusServiceUnavailable {
		t.Errorf("permissions without a store = %d, want 503", permissions.Code)
	}
}

type stubTeamTokens struct {
	createToken     string
	createErr       error
	lastGroup       string
	lastDescription string
	lastDays        int
	lastHours       int
	lastMinutes     int
	tokens          []teamtoken.Token
	listErr         error
	deleteErr       error
	lastDeleteID    string
}

func (s *stubTeamTokens) Verify(context.Context, string) (*auth.Identity, error) {
	return nil, errors.New("not used")
}

func (s *stubTeamTokens) Create(_ context.Context, _ auth.Identity, group string, days, hours, minutes int, description string) (string, error) {
	s.lastGroup, s.lastDays, s.lastHours, s.lastMinutes, s.lastDescription = group, days, hours, minutes, description
	return s.createToken, s.createErr
}

func (s *stubTeamTokens) List(context.Context, auth.Identity) ([]teamtoken.Token, error) {
	return s.tokens, s.listErr
}

func (s *stubTeamTokens) Delete(_ context.Context, id string, _ auth.Identity) error {
	s.lastDeleteID = id
	return s.deleteErr
}

func serveTeamTokens(t *testing.T, teamTokens TeamTokenService, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	verifier := stubIdentityVerifier{token: "user-token", identity: auth.Identity{Email: "user@example.com", Groups: []string{"platform"}}}
	router := NewAccessToken(Deps{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken:    "dev-token",
		Verifiers:    []auth.Verifier{verifier},
		AdminGroup:   "TERRAKUBE_ADMIN",
		AccessTokens: &fakeAccessTokenStore{},
		TeamTokens:   teamTokens,
	})
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+verifier.token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	return rec
}

func TestTeamTokenCreatePassesTheRequestAndReturnsTheToken(t *testing.T) {
	service := &stubTeamTokens{createToken: "signed-token"}
	rec := serveTeamTokens(t, service, http.MethodPost, "/teams",
		`{"group":"platform","description":"ci","days":1,"hours":2,"minutes":30}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", rec.Code, rec.Body.String())
	}
	if service.lastGroup != "platform" || service.lastDescription != "ci" ||
		service.lastDays != 1 || service.lastHours != 2 || service.lastMinutes != 30 {
		t.Errorf("service got %+v", service)
	}
	var body struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Token != "signed-token" {
		t.Fatalf("body = %s/%v", rec.Body.String(), err)
	}
}

func TestTeamTokenCreateRejectsNonMembersAndBadInput(t *testing.T) {
	for name, testCase := range map[string]struct {
		body   string
		err    error
		status int
	}{
		"non-member":        {`{"group":"platform"}`, teamtoken.ErrForbidden, http.StatusForbidden},
		"token token":       {`{"group":"platform"}`, teamtoken.ErrNonInteractive, http.StatusForbidden},
		"negative lifetime": {`{"group":"platform","days":-1}`, teamtoken.ErrInvalidLifetime, http.StatusBadRequest},
		"invalid json":      {"{", nil, http.StatusBadRequest},
	} {
		t.Run(name, func(t *testing.T) {
			rec := serveTeamTokens(t, &stubTeamTokens{createErr: testCase.err}, http.MethodPost, "/teams", testCase.body)
			if rec.Code != testCase.status {
				t.Fatalf("status = %d, want %d: %s", rec.Code, testCase.status, rec.Body.String())
			}
		})
	}
}

func TestTeamTokenListIsAccepted(t *testing.T) {
	description := "ci"
	service := &stubTeamTokens{tokens: []teamtoken.Token{{ID: "token-1", Description: &description}}}
	rec := serveTeamTokens(t, service, http.MethodGet, "/teams", "")

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202: %s", rec.Code, rec.Body.String())
	}
	var body []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body) != 1 || body[0]["id"] != "token-1" || body[0]["description"] != "ci" {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestTeamTokenDeleteMapsErrors(t *testing.T) {
	for name, testCase := range map[string]struct {
		err    error
		status int
	}{
		"unknown token": {teamtoken.ErrNotFound, http.StatusBadRequest},
		"foreign team":  {teamtoken.ErrForbidden, http.StatusForbidden},
		"revoked":       {nil, http.StatusAccepted},
	} {
		t.Run(name, func(t *testing.T) {
			service := &stubTeamTokens{deleteErr: testCase.err}
			rec := serveTeamTokens(t, service, http.MethodDelete, "/teams/token-1", "")
			if rec.Code != testCase.status {
				t.Fatalf("status = %d, want %d: %s", rec.Code, testCase.status, rec.Body.String())
			}
			if testCase.err == nil && service.lastDeleteID != "token-1" {
				t.Errorf("deleted %q, want token-1", service.lastDeleteID)
			}
		})
	}
}

func TestTeamTokensAnswer503WithoutAService(t *testing.T) {
	rec := serveTeamTokens(t, nil, http.MethodGet, "/teams", "")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
}

func TestTeamTokenRoutesRefuseOversizedBodies(t *testing.T) {
	router := NewAccessToken(Deps{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken:    "dev-token",
		AccessTokens: &fakeAccessTokenStore{},
		TeamTokens:   &stubTeamTokens{},
	})
	request := httptest.NewRequest(http.MethodPost, "/teams", strings.NewReader("{}"))
	request.Header.Set("Authorization", "Bearer dev-token")
	request.ContentLength = jsonapi.MaxBodyBytes + 1
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, request)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413: %s", rec.Code, rec.Body.String())
	}
}
