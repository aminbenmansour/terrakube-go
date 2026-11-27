package vcs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// gitlab talks to GitLab's v4 API.
type gitlab struct {
	base  string
	token string
}

func (g *gitlab) headers() http.Header {
	return http.Header{}
}

func (g *gitlab) groups(ctx context.Context) ([]Group, error) {
	groups := []Group{{ID: "@me", Name: "My Projects"}}
	var response []struct {
		ID       int    `json:"id"`
		FullPath string `json:"full_path"`
	}
	if _, err := getJSON(ctx, g.base+"/groups?membership=true&per_page=100", g.token, g.headers(), &response); err != nil {
		return nil, err
	}
	for _, group := range response {
		groups = append(groups, Group{ID: strconv.Itoa(group.ID), Name: group.FullPath})
	}
	return groups, nil
}

func (g *gitlab) repositories(ctx context.Context, group, search string, page int) (RepositoryPage, error) {
	var rawURL string
	if group == "" || group == "@me" {
		rawURL = fmt.Sprintf("%s/projects?owned=true&order_by=name", g.base)
	} else {
		rawURL = fmt.Sprintf("%s/groups/%s/projects?include_subgroups=true&order_by=name", g.base, url.PathEscape(group))
	}
	query := url.Values{}
	query.Set("per_page", strconv.Itoa(pageSize))
	query.Set("page", strconv.Itoa(page))
	if strings.TrimSpace(search) != "" {
		query.Set("search", search)
	}
	rawURL += "&" + query.Encode()

	var response []gitlabProjectJSON
	header, err := getJSON(ctx, rawURL, g.token, g.headers(), &response)
	if err != nil {
		return RepositoryPage{Items: []Repository{}, Page: page}, err
	}
	items := make([]Repository, 0, len(response))
	for _, project := range response {
		items = append(items, project.summary())
	}
	return RepositoryPage{Items: items, HasMore: header.Get("X-Next-Page") != "", Page: page}, nil
}

func (g *gitlab) branches(ctx context.Context, repository string, page int) (BranchPage, error) {
	path := url.PathEscape(strings.Trim(strings.TrimSpace(repository), "/"))
	rawURL := fmt.Sprintf("%s/projects/%s/repository/branches?per_page=%d&page=%d", g.base, path, branchesPerPage, page)
	var response []struct {
		Name string `json:"name"`
	}
	header, err := getJSON(ctx, rawURL, g.token, g.headers(), &response)
	if err != nil {
		return BranchPage{Items: []Branch{}, Page: page}, err
	}
	items := make([]Branch, 0, len(response))
	for _, branch := range response {
		items = append(items, Branch{Name: branch.Name})
	}
	return BranchPage{Items: items, HasMore: header.Get("X-Next-Page") != "", Page: page}, nil
}

func (g *gitlab) test(ctx context.Context) (string, error) {
	var me struct {
		Username string `json:"username"`
	}
	if _, err := getJSON(ctx, g.base+"/user", g.token, g.headers(), &me); err != nil {
		return "", err
	}
	return "Connected as " + me.Username, nil
}

type gitlabProjectJSON struct {
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	HTTPURLToRepo     string `json:"http_url_to_repo"`
	Visibility        string `json:"visibility"`
	DefaultBranch     string `json:"default_branch"`
	Namespace         struct {
		FullPath string `json:"full_path"`
	} `json:"namespace"`
}

func (p gitlabProjectJSON) summary() Repository {
	return Repository{
		Name: p.Name, FullName: p.PathWithNamespace, Group: p.Namespace.FullPath,
		URL: p.HTTPURLToRepo, PrivateRepo: p.Visibility != "public",
		DefaultBranch: optional(p.DefaultBranch),
	}
}
