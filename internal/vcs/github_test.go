package vcs

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func githubRoutes(t *testing.T) map[string]func(http.ResponseWriter, *http.Request) {
	return map[string]func(http.ResponseWriter, *http.Request){
		"/user": func(w http.ResponseWriter, r *http.Request) {
			requireBearer(t, r, "gh-token")
			writeJSON(w, `{"login":"octocat"}`)
		},
		"/user/orgs": func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("per_page"); got != "100" {
				t.Errorf("per_page = %q, want 100", got)
			}
			writeJSON(w, `[{"login":"acme"}]`)
		},
		"/user/repos": func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("affiliation"); got != "owner" {
				t.Errorf("affiliation = %q, want owner", got)
			}
			if got := r.URL.Query().Get("page"); got != "1" {
				t.Errorf("page = %q, want 1", got)
			}
			writeJSON(w, `[{"name":"personal","full_name":"octocat/personal",
				"clone_url":"https://github.com/octocat/personal.git",
				"private":false,"default_branch":"main","owner":{"login":"octocat"}}]`)
		},
		"/orgs/acme/repos": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, `[{"name":"infra","full_name":"acme/infra",
				"clone_url":"https://github.com/acme/infra.git","html_url":"https://github.com/acme/infra",
				"private":true,"default_branch":"trunk","owner":{"login":"acme"}}]`)
		},
		"/search/repositories": func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("q"); got != "infra in:name user:acme" {
				t.Errorf("q = %q", got)
			}
			writeJSON(w, `{"total_count":51,"items":[{"name":"infra","full_name":"acme/infra",
				"clone_url":"https://github.com/acme/infra.git","private":true,
				"default_branch":"main","owner":{"login":"acme"}}]}`)
		},
		"/repos/acme/infra/branches": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, `[{"name":"main"},{"name":"release"}]`)
		},
	}
}

func TestGroupsGitHub(t *testing.T) {
	server := newFakeProvider(t, githubRoutes(t))
	groups, err := Groups(context.Background(), Connection{
		Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token",
	})
	if err != nil {
		t.Fatalf("Groups: %v", err)
	}
	want := []Group{{ID: "@me", Name: "octocat (personal)"}, {ID: "acme", Name: "acme"}}
	if len(groups) != len(want) || groups[0] != want[0] || groups[1] != want[1] {
		t.Fatalf("groups = %+v, want %+v", groups, want)
	}
}

func TestRepositoriesGitHub(t *testing.T) {
	server := newFakeProvider(t, githubRoutes(t))
	connection := Connection{Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token"}

	page, err := Repositories(context.Background(), connection, "acme", "", 1)
	if err != nil {
		t.Fatalf("Repositories: %v", err)
	}
	if len(page.Items) != 1 {
		t.Fatalf("items = %+v", page.Items)
	}
	repo := page.Items[0]
	if repo.Name != "infra" || repo.FullName != "acme/infra" || repo.Group != "acme" ||
		repo.URL != "https://github.com/acme/infra.git" || !repo.PrivateRepo ||
		repo.DefaultBranch == nil || *repo.DefaultBranch != "trunk" {
		t.Fatalf("repository = %+v", repo)
	}
	if page.HasMore || page.Page != 1 {
		t.Fatalf("page = %+v", page)
	}

	searched, err := Repositories(context.Background(), connection, "acme", "infra", 1)
	if err != nil {
		t.Fatalf("search Repositories: %v", err)
	}
	if !searched.HasMore {
		t.Fatalf("hasMore = false, want true for total_count 51 on page 1")
	}

	personal, err := Repositories(context.Background(), connection, "@me", "", 1)
	if err != nil {
		t.Fatalf("personal Repositories: %v", err)
	}
	if len(personal.Items) != 1 || personal.Items[0].FullName != "octocat/personal" {
		t.Fatalf("personal items = %+v", personal.Items)
	}
	if personal.Items[0].DefaultBranch == nil || *personal.Items[0].DefaultBranch != "main" {
		t.Fatalf("default branch = %v", personal.Items[0].DefaultBranch)
	}
}

func TestBranchesGitHub(t *testing.T) {
	server := newFakeProvider(t, githubRoutes(t))
	page, err := Branches(context.Background(),
		Connection{Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token"}, "acme/infra", 1)
	if err != nil {
		t.Fatalf("Branches: %v", err)
	}
	if len(page.Items) != 2 || page.Items[0].Name != "main" || page.Items[1].Name != "release" {
		t.Fatalf("branches = %+v", page.Items)
	}

	_, err = Branches(context.Background(),
		Connection{Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token"}, "  ", 1)
	var discoveryErr *Error
	if !errors.As(err, &discoveryErr) || discoveryErr.Status != http.StatusBadRequest {
		t.Fatalf("empty repository error = %v", err)
	}
}

func TestTestConnection(t *testing.T) {
	server := newFakeProvider(t, githubRoutes(t))
	message, err := Test(context.Background(),
		Connection{Type: TypeGitHub, APIURL: server.URL, AccessToken: "gh-token"})
	if err != nil {
		t.Fatalf("Test: %v", err)
	}
	if message != "Connected as octocat" {
		t.Fatalf("message = %q", message)
	}
}

func TestGitHubConnectionsWithoutAToken(t *testing.T) {
	cases := []struct {
		name       string
		connection Connection
		status     int
	}{
		{"no token", Connection{Type: TypeGitHub, APIURL: "https://api.github.com"}, http.StatusBadRequest},
		{"app", Connection{Type: TypeGitHub, APIURL: "https://api.github.com", PrivateKey: "pem"}, http.StatusNotImplemented},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Groups(context.Background(), testCase.connection)
			var discoveryErr *Error
			if !errors.As(err, &discoveryErr) || discoveryErr.Status != testCase.status {
				t.Fatalf("error = %v, want status %d", err, testCase.status)
			}
		})
	}
}
