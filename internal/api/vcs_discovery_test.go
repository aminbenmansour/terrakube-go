package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const vcsDiscoveryID = "31313131-3131-4131-8131-313131313131"

// discoveryStub adds the team lookup the non-owner path needs; with no
// organization in the list the caller is not a member.
type discoveryStub struct {
	*stubQuerier
	organizationIDs []string
}

func (s *discoveryStub) ListOrganizationIDsByGroupNames(context.Context, []string) ([]string, error) {
	return s.organizationIDs, nil
}

// fakeGitHub serves the endpoints the discovery handlers call.
func fakeGitHub(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/user":
			_, _ = w.Write([]byte(`{"login":"octocat"}`))
		case "/user/orgs":
			_, _ = w.Write([]byte(`[{"login":"acme"}]`))
		case "/orgs/acme/repos":
			_, _ = w.Write([]byte(`[{"name":"infra","full_name":"acme/infra",
				"clone_url":"https://github.com/acme/infra.git","private":true,
				"default_branch":"main","owner":{"login":"acme"}}]`))
		case "/search/repositories":
			_, _ = w.Write([]byte(`{"total_count":51,"items":[{"name":"infra","full_name":"acme/infra",
				"clone_url":"https://github.com/acme/infra.git","private":true,
				"default_branch":"main","owner":{"login":"acme"}}]}`))
		case "/repos/acme/infra/branches":
			_, _ = w.Write([]byte(`[{"name":"main"},{"name":"release"}]`))
		default:
			t.Errorf("unexpected request %s", r.URL.String())
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

// vcsDiscoveryAPI serves the root /vcs/v1 router against a fake provider.
func vcsDiscoveryAPI(t *testing.T, providerURL string) http.Handler {
	t.Helper()
	q := &discoveryStub{stubQuerier: &stubQuerier{
		org: store.Organization{ID: vcsOrgID},
		vcs: []store.Vc{{
			ID: vcsDiscoveryID, OrganizationID: vcsOrgID, Name: "github", VcsType: "GITHUB",
			ApiUrl: strPtr(providerURL), AccessToken: strPtr("gh-token"),
			ConnectionType: strPtr("OAUTH"),
		}},
	}}
	return NewVCS(Deps{Store: q})
}

func decodeBody[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var value T
	if err := json.Unmarshal(rec.Body.Bytes(), &value); err != nil {
		t.Fatalf("decode: %v; body = %s", err, rec.Body.String())
	}
	return value
}

type discoveryGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type discoveryRepositoryPage struct {
	Items   []map[string]any `json:"items"`
	HasMore bool             `json:"hasMore"`
	Page    int              `json:"page"`
}

func TestVCSDiscoveryGroups(t *testing.T) {
	server := fakeGitHub(t)
	h := vcsDiscoveryAPI(t, server.URL)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/"+vcsDiscoveryID+"/groups", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	groups := decodeBody[[]discoveryGroup](t, rec)
	if len(groups) != 2 || groups[0].ID != "@me" || groups[1].Name != "acme" {
		t.Fatalf("groups = %+v", groups)
	}
}

func TestVCSDiscoveryRepositoriesAndBranches(t *testing.T) {
	server := fakeGitHub(t)
	h := vcsDiscoveryAPI(t, server.URL)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/"+vcsDiscoveryID+"/repositories?group=acme&search=infra&page=1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	page := decodeBody[discoveryRepositoryPage](t, rec)
	if len(page.Items) != 1 || page.Items[0]["fullName"] != "acme/infra" || !page.HasMore || page.Page != 1 {
		t.Fatalf("page = %+v", page)
	}
	if page.Items[0]["privateRepo"] != true {
		t.Fatalf("privateRepo = %v", page.Items[0]["privateRepo"])
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/"+vcsDiscoveryID+"/branches?repository=acme/infra", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("branches status = %d; body = %s", rec.Code, rec.Body.String())
	}
	branches := decodeBody[struct {
		Items []struct {
			Name string `json:"name"`
		} `json:"items"`
	}](t, rec)
	if len(branches.Items) != 2 || branches.Items[1].Name != "release" {
		t.Fatalf("branches = %+v", branches.Items)
	}
}

func TestVCSDiscoveryTestConnection(t *testing.T) {
	server := fakeGitHub(t)
	h := vcsDiscoveryAPI(t, server.URL)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/"+vcsDiscoveryID+"/test", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	connection := decodeBody[struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
	}](t, rec)
	if !connection.OK || connection.Message != "Connected as octocat" {
		t.Fatalf("test = %+v", connection)
	}
}

func TestVCSDiscoveryRejectsBadRequests(t *testing.T) {
	server := fakeGitHub(t)
	h := vcsDiscoveryAPI(t, server.URL)

	cases := []struct {
		name   string
		method string
		target string
		status int
	}{
		{"malformed id", http.MethodGet, "/not-a-uuid/groups", http.StatusBadRequest},
		{"unknown connection", http.MethodGet, "/99999999-9999-4999-8999-999999999999/groups", http.StatusNotFound},
		{"missing group", http.MethodGet, "/" + vcsDiscoveryID + "/repositories", http.StatusBadRequest},
		{"bad page", http.MethodGet, "/" + vcsDiscoveryID + "/repositories?group=acme&page=zero", http.StatusBadRequest},
		{"missing repository", http.MethodGet, "/" + vcsDiscoveryID + "/branches", http.StatusBadRequest},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(testCase.method, testCase.target, nil))
			if rec.Code != testCase.status {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, testCase.status, rec.Body.String())
			}
		})
	}
}

func TestVCSDiscoveryDeniedWithoutPermission(t *testing.T) {
	server := fakeGitHub(t)
	h := vcsDiscoveryAPI(t, server.URL)

	// An authenticated caller whose group maps to no team in the connection's
	// organization gets the 403 upstream's @PreAuthorize answers.
	rec := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/"+vcsDiscoveryID+"/groups", nil)
	request = request.WithContext(auth.WithIdentity(request.Context(),
		&auth.Identity{Issuer: "dex", Groups: []string{"other-team"}}))
	h.ServeHTTP(rec, request)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body.String())
	}
}
