// Package vcs lists repositories and branches from a VCS provider using the
// connection's stored credentials, mirroring upstream's discovery services
// (api/.../plugin/vcs/discovery).
package vcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Upstream's VcsType values this package can discover.
const (
	TypeGitHub = "GITHUB"
	TypeGitLab = "GITLAB"
)

const (
	// pageSize is the repository page size upstream requests.
	pageSize = 50
	// branchesPerPage is GitHub's and GitLab's default maximum.
	branchesPerPage = 100
	// maxBody keeps a hostile or broken provider response bounded.
	maxBody = 4 << 20
)

// Group is a provider organization or the caller's personal namespace.
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Repository is the summary the workspace picker reads; the field names are
// upstream's VcsRepositorySummary.
type Repository struct {
	Name          string  `json:"name"`
	FullName      string  `json:"fullName"`
	Group         string  `json:"group"`
	URL           string  `json:"url"`
	PrivateRepo   bool    `json:"privateRepo"`
	DefaultBranch *string `json:"defaultBranch"`
}

// RepositoryPage mirrors upstream's VcsRepositoryPage.
type RepositoryPage struct {
	Items   []Repository `json:"items"`
	HasMore bool         `json:"hasMore"`
	Page    int          `json:"page"`
}

// Branch is one branch name; the branch picker only needs the name.
type Branch struct {
	Name string `json:"name"`
}

// BranchPage is the branch-listing response.
type BranchPage struct {
	Items   []Branch `json:"items"`
	HasMore bool     `json:"hasMore"`
	Page    int      `json:"page"`
}

// Connection is what discovery needs from a stored VCS record. The store
// decrypts Token and PrivateKey before they reach this package.
type Connection struct {
	Type           string
	APIURL         string
	AccessToken    string
	PrivateKey     string
	ConnectionType string
	// AllowLocal permits loopback and link-local API URLs. Development stacks
	// and tests set it; production leaves it false so a connection cannot be
	// used to reach cloud metadata or a loopback service.
	AllowLocal bool
}

// Error carries the status the API should answer with, so a rejected
// credential never surfaces as a 500.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return e.Message }

func failed(status int, format string, args ...any) *Error {
	return &Error{Status: status, Message: fmt.Sprintf(format, args...)}
}

// Groups lists the organizations (and the personal namespace) the connection
// can see.
func Groups(ctx context.Context, connection Connection) ([]Group, error) {
	provider, err := providerFor(connection)
	if err != nil {
		return nil, err
	}
	return provider.groups(ctx)
}

// Repositories lists one group's repositories.
func Repositories(ctx context.Context, connection Connection, group, search string, page int) (RepositoryPage, error) {
	provider, err := providerFor(connection)
	if err != nil {
		return RepositoryPage{Items: []Repository{}, Page: page}, err
	}
	return provider.repositories(ctx, group, search, page)
}

// Branches lists a repository's branches. repository is the provider's full
// name (owner/name for GitHub, path_with_namespace for GitLab).
func Branches(ctx context.Context, connection Connection, repository string, page int) (BranchPage, error) {
	provider, err := providerFor(connection)
	if err != nil {
		return BranchPage{Items: []Branch{}, Page: page}, err
	}
	if strings.TrimSpace(repository) == "" {
		return BranchPage{Items: []Branch{}, Page: page}, failed(http.StatusBadRequest, "A repository is required")
	}
	return provider.branches(ctx, repository, page)
}

// Test checks the stored credential and returns a short description of the
// authenticated account.
func Test(ctx context.Context, connection Connection) (string, error) {
	provider, err := providerFor(connection)
	if err != nil {
		return "", err
	}
	return provider.test(ctx)
}

type provider interface {
	groups(context.Context) ([]Group, error)
	repositories(context.Context, string, string, int) (RepositoryPage, error)
	branches(context.Context, string, int) (BranchPage, error)
	test(context.Context) (string, error)
	syncWebhook(context.Context, WebhookSpec, string) (string, error)
	removeWebhook(context.Context, string, string) error
}

func providerFor(connection Connection) (provider, error) {
	base := baseURL(connection)
	// Only a stored URL can point somewhere unexpected; an empty api_url uses
	// the fixed public provider host.
	if strings.TrimSpace(connection.APIURL) != "" {
		if err := validateEndpoint(base, connection.AllowLocal); err != nil {
			return nil, err
		}
	}
	switch connection.Type {
	case TypeGitHub:
		if strings.TrimSpace(connection.AccessToken) == "" {
			if strings.TrimSpace(connection.PrivateKey) != "" {
				return nil, failed(http.StatusNotImplemented,
					"GitHub App connections need an installation token, which is not implemented yet")
			}
			return nil, failed(http.StatusBadRequest,
				"This GitHub connection has no access token; connect it with OAuth or store a token")
		}
		return &github{base: base, token: connection.AccessToken}, nil
	case TypeGitLab:
		if strings.TrimSpace(connection.AccessToken) == "" {
			return nil, failed(http.StatusBadRequest,
				"This GitLab connection has no access token; connect it with OAuth or store a token")
		}
		return &gitlab{base: base, token: connection.AccessToken}, nil
	default:
		return nil, failed(http.StatusUnprocessableEntity,
			"Repository discovery is not supported for this VCS connection type")
	}
}

func baseURL(connection Connection) string {
	base := strings.TrimRight(strings.TrimSpace(connection.APIURL), "/")
	if base != "" {
		return base
	}
	switch connection.Type {
	case TypeGitLab:
		return "https://gitlab.com/api/v4"
	default:
		return "https://api.github.com"
	}
}

// getJSON performs one authenticated GET and decodes a 2xx body into target.
func getJSON(ctx context.Context, rawURL, token string, header http.Header, target any) (http.Header, error) {
	return doJSON(ctx, http.MethodGet, rawURL, token, header, nil, target)
}

// doJSON performs one authenticated request and decodes a 2xx body into target.
// A nil target ignores the body. Provider failures carry the status the API
// should answer with, so a rejected credential never surfaces as a 500.
func doJSON(ctx context.Context, method, rawURL, token string, header http.Header, payload, target any) (http.Header, error) {
	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, failed(http.StatusInternalServerError, "Could not encode the request: %v", err)
		}
		body = bytes.NewReader(encoded)
	}
	request, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, failed(http.StatusBadGateway, "The provider URL is invalid: %v", err)
	}
	request.Header.Set("Accept", "application/json")
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	for key, values := range header {
		request.Header[key] = values
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		// A url.Error carries the request URL only; the token travels in a
		// header, so it never appears in this message.
		if urlErr, ok := err.(*url.Error); ok && urlErr.Err != nil {
			err = urlErr.Err
		}
		return nil, failed(http.StatusBadGateway, "Could not reach the VCS provider: %v", err)
	}
	defer response.Body.Close()
	responseBody, _ := io.ReadAll(io.LimitReader(response.Body, maxBody))

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, providerError(response.StatusCode, responseBody)
	}
	if target != nil && len(responseBody) > 0 {
		if err := json.Unmarshal(responseBody, target); err != nil {
			return nil, failed(http.StatusBadGateway, "The VCS provider returned an unexpected response: %v", err)
		}
	}
	return response.Header, nil
}

// providerError maps a provider status and body to the API's answer.
func providerError(status int, body []byte) error {
	message := providerMessage(body)
	switch status {
	case http.StatusUnauthorized:
		return failed(http.StatusUnauthorized,
			"The VCS provider rejected the credential (HTTP 401). Reconnect it or check the token. %s", message)
	case http.StatusForbidden:
		return failed(http.StatusForbidden,
			"The VCS provider refused the request (HTTP 403). Check the token's scopes. %s", message)
	case http.StatusNotFound:
		return failed(http.StatusNotFound, "The VCS provider could not find that resource. %s", message)
	case http.StatusTooManyRequests:
		return failed(http.StatusTooManyRequests, "The VCS provider rate limit was reached. Try again later.")
	default:
		if status >= 500 {
			return failed(http.StatusBadGateway,
				"The VCS provider failed the request (HTTP %d). %s", status, message)
		}
		return failed(http.StatusBadRequest,
			"The VCS provider rejected the request (HTTP %d). %s", status, message)
	}
}

// providerMessage extracts the short error text GitHub and GitLab send.
func providerMessage(body []byte) string {
	var payload struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if json.Unmarshal(body, &payload) == nil {
		if payload.Message != "" {
			return "(" + payload.Message + ")"
		}
		if payload.Error != "" {
			return "(" + payload.Error + ")"
		}
	}
	return ""
}

func optional(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
