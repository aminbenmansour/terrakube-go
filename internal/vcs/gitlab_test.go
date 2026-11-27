package vcs

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func gitlabRoutes(t *testing.T) map[string]func(http.ResponseWriter, *http.Request) {
	return map[string]func(http.ResponseWriter, *http.Request){
		"/user": func(w http.ResponseWriter, r *http.Request) {
			requireBearer(t, r, "gl-token")
			writeJSON(w, `{"username":"jdoe"}`)
		},
		"/groups": func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("membership"); got != "true" {
				t.Errorf("membership = %q, want true", got)
			}
			writeJSON(w, `[{"id":7,"full_path":"acme/sub"}]`)
		},
		"/projects": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, `[{"name":"app","path_with_namespace":"jdoe/app",
				"http_url_to_repo":"https://gitlab.com/jdoe/app.git","visibility":"public",
				"default_branch":"main","namespace":{"full_path":"jdoe"}}]`)
		},
		"/groups/7/projects": func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("include_subgroups"); got != "true" {
				t.Errorf("include_subgroups = %q, want true", got)
			}
			w.Header().Set("X-Next-Page", "2")
			writeJSON(w, `[{"name":"infra","path_with_namespace":"acme/sub/infra",
				"http_url_to_repo":"https://gitlab.com/acme/sub/infra.git","visibility":"private",
				"default_branch":null,"namespace":{"full_path":"acme/sub"}}]`)
		},
		"/projects/acme/sub/infra/repository/branches": func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, `[{"name":"main"}]`)
		},
	}
}

func TestGitLabDiscovery(t *testing.T) {
	server := newFakeProvider(t, gitlabRoutes(t))
	connection := Connection{AllowLocal: true, Type: TypeGitLab, APIURL: server.URL, AccessToken: "gl-token"}

	groups, err := Groups(context.Background(), connection)
	if err != nil {
		t.Fatalf("Groups: %v", err)
	}
	if len(groups) != 2 || groups[0].ID != "@me" || groups[1] != (Group{ID: "7", Name: "acme/sub"}) {
		t.Fatalf("groups = %+v", groups)
	}

	personal, err := Repositories(context.Background(), connection, "@me", "", 1)
	if err != nil {
		t.Fatalf("Repositories: %v", err)
	}
	if len(personal.Items) != 1 || personal.Items[0].PrivateRepo || personal.Items[0].DefaultBranch == nil {
		t.Fatalf("personal repository = %+v", personal.Items)
	}

	group, err := Repositories(context.Background(), connection, "7", "", 2)
	if err != nil {
		t.Fatalf("group Repositories: %v", err)
	}
	if !group.HasMore || group.Page != 2 {
		t.Fatalf("group page = %+v, want hasMore from X-Next-Page and page 2", group)
	}
	if group.Items[0].DefaultBranch != nil {
		t.Fatalf("default branch = %v, want null for an empty provider value", *group.Items[0].DefaultBranch)
	}
	if !group.Items[0].PrivateRepo || group.Items[0].Group != "acme/sub" {
		t.Fatalf("group repository = %+v", group.Items[0])
	}

	branches, err := Branches(context.Background(), connection, "acme/sub/infra", 1)
	if err != nil {
		t.Fatalf("Branches: %v", err)
	}
	if len(branches.Items) != 1 || branches.Items[0].Name != "main" {
		t.Fatalf("branches = %+v", branches.Items)
	}

	message, err := Test(context.Background(), connection)
	if err != nil || message != "Connected as jdoe" {
		t.Fatalf("Test = %q/%v", message, err)
	}
}

func TestGitLabConnectionWithoutAToken(t *testing.T) {
	_, err := Groups(context.Background(), Connection{AllowLocal: true, Type: TypeGitLab, APIURL: "https://gitlab.com"})
	var discoveryErr *Error
	if !errors.As(err, &discoveryErr) || discoveryErr.Status != http.StatusBadRequest {
		t.Fatalf("error = %v, want 400", err)
	}
}
