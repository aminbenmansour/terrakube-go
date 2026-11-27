package vcs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// github talks to GitHub's REST API, including GitHub Enterprise Server (whose
// api_url ends in /api/v3).
type github struct {
	base  string
	token string
}

func (g *github) headers() http.Header {
	header := http.Header{}
	header.Set("Accept", "application/vnd.github+json")
	header.Set("X-GitHub-Api-Version", "2022-11-28")
	return header
}

func (g *github) groups(ctx context.Context) ([]Group, error) {
	groups := []Group{}
	if g.token != "" {
		var me struct {
			Login string `json:"login"`
		}
		if _, err := getJSON(ctx, g.base+"/user", g.token, g.headers(), &me); err != nil {
			return nil, err
		}
		groups = append(groups, Group{ID: "@me", Name: me.Login + " (personal)"})
	}

	var organizations []struct {
		Login string `json:"login"`
	}
	if _, err := getJSON(ctx, g.base+"/user/orgs?per_page=100", g.token, g.headers(), &organizations); err != nil {
		return nil, err
	}
	for _, organization := range organizations {
		groups = append(groups, Group{ID: organization.Login, Name: organization.Login})
	}
	return groups, nil
}

func (g *github) repositories(ctx context.Context, group, search string, page int) (RepositoryPage, error) {
	if strings.TrimSpace(search) != "" {
		query := url.Values{}
		scope := search + " in:name"
		if group != "" && group != "@me" {
			scope += " user:" + group
		}
		query.Set("q", scope)
		query.Set("per_page", strconv.Itoa(pageSize))
		query.Set("page", strconv.Itoa(page))
		var response struct {
			TotalCount int `json:"total_count"`
			Items      []githubRepositoryJSON
		}
		if _, err := getJSON(ctx, g.base+"/search/repositories?"+query.Encode(), g.token, g.headers(), &response); err != nil {
			return RepositoryPage{Items: []Repository{}, Page: page}, err
		}
		items := make([]Repository, 0, len(response.Items))
		for _, repo := range response.Items {
			items = append(items, repo.summary())
		}
		return RepositoryPage{Items: items, HasMore: response.TotalCount > page*pageSize, Page: page}, nil
	}

	rawURL := fmt.Sprintf("%s/user/repos?affiliation=owner&per_page=%d&page=%d", g.base, pageSize, page)
	if group != "" && group != "@me" {
		rawURL = fmt.Sprintf("%s/orgs/%s/repos?per_page=%d&page=%d", g.base, url.PathEscape(group), pageSize, page)
	}
	var response []githubRepositoryJSON
	if _, err := getJSON(ctx, rawURL, g.token, g.headers(), &response); err != nil {
		return RepositoryPage{Items: []Repository{}, Page: page}, err
	}
	items := make([]Repository, 0, len(response))
	for _, repo := range response {
		items = append(items, repo.summary())
	}
	return RepositoryPage{Items: items, HasMore: len(items) == pageSize, Page: page}, nil
}

func (g *github) branches(ctx context.Context, repository string, page int) (BranchPage, error) {
	path := strings.Trim(strings.TrimSpace(repository), "/")
	rawURL := fmt.Sprintf("%s/repos/%s/branches?per_page=%d&page=%d", g.base, path, branchesPerPage, page)
	var response []struct {
		Name string `json:"name"`
	}
	if _, err := getJSON(ctx, rawURL, g.token, g.headers(), &response); err != nil {
		return BranchPage{Items: []Branch{}, Page: page}, err
	}
	items := make([]Branch, 0, len(response))
	for _, branch := range response {
		items = append(items, Branch{Name: branch.Name})
	}
	return BranchPage{Items: items, HasMore: len(items) == branchesPerPage, Page: page}, nil
}

func (g *github) test(ctx context.Context) (string, error) {
	if g.token == "" {
		if _, err := getJSON(ctx, g.base+"/repositories?per_page=1", "", g.headers(), &json.RawMessage{}); err != nil {
			return "", err
		}
		return "The public GitHub API is reachable", nil
	}
	var me struct {
		Login string `json:"login"`
	}
	if _, err := getJSON(ctx, g.base+"/user", g.token, g.headers(), &me); err != nil {
		return "", err
	}
	return "Connected as " + me.Login, nil
}

type githubRepositoryJSON struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	CloneURL      string `json:"clone_url"`
	HTMLURL       string `json:"html_url"`
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func (r githubRepositoryJSON) summary() Repository {
	repositoryURL := r.CloneURL
	if repositoryURL == "" && r.HTMLURL != "" {
		repositoryURL = r.HTMLURL + ".git"
	}
	return Repository{
		Name: r.Name, FullName: r.FullName, Group: r.Owner.Login, URL: repositoryURL,
		PrivateRepo: r.Private, DefaultBranch: optional(r.DefaultBranch),
	}
}
