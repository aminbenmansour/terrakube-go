package api

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	authzOrg     = "70000000-0000-0000-0000-000000000001"
	authzOther   = "70000000-0000-0000-0000-000000000002"
	authzWS      = "70000000-0000-0000-0000-000000000067"
	authzProject = "70000000-0000-0000-0000-000000000010"
)

// groupVerifier maps "groups:a,b" to an identity with those groups, so a test
// can act as any team without minting a token.
type groupVerifier struct{}

func (groupVerifier) Verify(_ context.Context, raw string) (*auth.Identity, error) {
	groups, ok := strings.CutPrefix(raw, "groups:")
	if !ok {
		return nil, errors.New("unknown token")
	}
	return &auth.Identity{Issuer: "dex", Email: "user@example.com", Groups: strings.Split(groups, ",")}, nil
}

func authzAPI(t *testing.T, q *stubQuerier) http.Handler {
	t.Helper()
	return New(Deps{
		Store:      q,
		Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
		Verifiers:  []auth.Verifier{groupVerifier{}},
		AdminGroup: "TERRAKUBE_ADMIN",
	})
}

func authzRequest(t *testing.T, handler http.Handler, groups, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer groups:"+groups)
	if body != "" {
		request.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, request)
	return rec
}

func TestOrganizationPermissionMatrix(t *testing.T) {
	const vcsBody = `{"data":{"type":"vcs","attributes":{"name":"probe","description":"probe","vcsType":"GITLAB","connectionType":"STANDALONE","clientId":"c","clientSecret":"s","endpoint":"https://gitlab.example.com"}}}`
	const teamBody = `{"data":{"type":"team","attributes":{"name":"new-team","role":"read"}}}`

	tests := []struct {
		name   string
		groups string
		method string
		path   string
		body   string
		want   int
		detail string
	}{
		{"a non-member cannot read the organization", "strangers", http.MethodGet, "/organization/" + authzOrg, "", http.StatusForbidden, "ReadPermission Denied"},
		{"a member can read the organization", "readers", http.MethodGet, "/organization/" + authzOrg, "", http.StatusOK, ""},
		{"a member can list VCS providers", "readers", http.MethodGet, "/organization/" + authzOrg + "/vcs", "", http.StatusOK, ""},
		{"a read role cannot create a VCS provider", "readers", http.MethodPost, "/organization/" + authzOrg + "/vcs", vcsBody, http.StatusForbidden, "CreatePermission Denied"},
		{"an admin role can create a VCS provider", "admins", http.MethodPost, "/organization/" + authzOrg + "/vcs", vcsBody, http.StatusCreated, ""},
		{"a read role cannot create a team", "readers", http.MethodPost, "/organization/" + authzOrg + "/team", teamBody, http.StatusForbidden, "CreatePermission Denied"},
		{"the instance owner can create a team", "TERRAKUBE_ADMIN", http.MethodPost, "/organization/" + authzOrg + "/team", teamBody, http.StatusCreated, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := &stubQuerier{
				org: store.Organization{ID: authzOrg, Name: "mine"},
				teams: []store.Team{
					{Name: "readers", OrganizationID: authzOrg, Role: strPtr("read")},
					{Name: "admins", OrganizationID: authzOrg, Role: strPtr("admin")},
				},
			}
			rec := authzRequest(t, authzAPI(t, q), test.groups, test.method, test.path, test.body)
			if rec.Code != test.want {
				t.Fatalf("status = %d, want %d: %s", rec.Code, test.want, rec.Body.String())
			}
			if test.detail != "" && !strings.Contains(rec.Body.String(), test.detail) {
				t.Errorf("body = %s, want detail %q", rec.Body.String(), test.detail)
			}
		})
	}
}

func TestOrganizationListIsScopedToTheCaller(t *testing.T) {
	q := &stubQuerier{
		orgs: []store.Organization{
			{ID: authzOrg, Name: "mine"},
			{ID: authzOther, Name: "other"},
		},
		teams: []store.Team{{Name: "readers", OrganizationID: authzOrg, Role: strPtr("read")}},
	}

	rec := authzRequest(t, authzAPI(t, q), "readers", http.MethodGet, "/organization", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), authzOther) {
		t.Errorf("body = %s, want only the caller's organization", rec.Body.String())
	}

	rec = authzRequest(t, authzAPI(t, q), "strangers", http.MethodGet, "/organization", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"data":[]`) {
		t.Fatalf("a token without groups status = %d body = %s, want an empty list", rec.Code, rec.Body.String())
	}
}

func TestWorkspaceGrantAllowsTheWrite(t *testing.T) {
	const body = `{"data":{"type":"workspace","id":"` + authzWS + `","attributes":{"description":"x"}}}`

	withoutGrant := &stubQuerier{
		workspaces: []store.Workspace{{ID: authzWS, OrganizationID: authzOrg}},
		teams:      []store.Team{{Name: "readers", OrganizationID: authzOrg, Role: strPtr("read")}},
	}
	rec := authzRequest(t, authzAPI(t, withoutGrant), "readers", http.MethodPatch,
		"/organization/"+authzOrg+"/workspace/"+authzWS, body)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("without a grant status = %d, want 403: %s", rec.Code, rec.Body.String())
	}

	withGrant := &stubQuerier{
		workspaces: []store.Workspace{{ID: authzWS, OrganizationID: authzOrg}},
		teams:      []store.Team{{Name: "readers", OrganizationID: authzOrg, Role: strPtr("read")}},
		access: []store.Access{{
			WorkspaceID: authzWS, Name: "readers", ManageWorkspace: true,
		}},
	}
	rec = authzRequest(t, authzAPI(t, withGrant), "readers", http.MethodPatch,
		"/organization/"+authzOrg+"/workspace/"+authzWS, body)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("with a workspace grant status = %d, want 204: %s", rec.Code, rec.Body.String())
	}
}

func TestProjectGrantOfAnotherOrganizationIsNotFound(t *testing.T) {
	q := &stubQuerier{
		projects: []store.Project{{ID: authzProject, OrganizationID: authzOther}},
		teams:    []store.Team{{Name: "readers", OrganizationID: authzOrg, Role: strPtr("read")}},
	}

	rec := authzRequest(t, authzAPI(t, q), "readers", http.MethodGet,
		"/organization/"+authzOrg+"/project/"+authzProject+"/projectAccess", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Unknown identifier ["+authzProject+"] for project") {
		t.Errorf("body = %s, want upstream's identifier message", rec.Body.String())
	}
}
