package vcs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// WebhookSpec is the repository webhook upstream's WebhookService registers for
// a workspace: the repository to attach it to, our callback URL, the derived
// secret and the trigger events the user configured.
type WebhookSpec struct {
	SourceURL string
	// CallbackURL is where the provider posts deliveries.
	CallbackURL string
	// Secret is sent to the provider and is never returned to a caller.
	Secret string
	// Events are upstream's WebhookEventType names (PUSH, PULL_REQUEST, ...).
	Events []string
	// PullRequestWorkflow subscribes to pull-request comments as well.
	PullRequestWorkflow bool
}

// ValidateConnection reports whether this package can talk to the connection's
// provider at all (supported type and usable credential).
func ValidateConnection(connection Connection) error {
	_, err := providerFor(connection)
	return err
}

// SyncWebhook creates the repository webhook, or updates the existing one when
// remoteID is set, and returns the provider's id. GitHub's update only changes
// the subscription; GitLab replaces the whole hook.
func SyncWebhook(ctx context.Context, connection Connection, spec WebhookSpec, remoteID string) (string, error) {
	provider, err := providerFor(connection)
	if err != nil {
		return "", err
	}
	return provider.syncWebhook(ctx, spec, remoteID)
}

// DeleteWebhook removes a repository webhook. A provider that no longer knows
// the id is treated as success, since the desired state already holds.
func DeleteWebhook(ctx context.Context, connection Connection, sourceURL, remoteID string) error {
	provider, err := providerFor(connection)
	if err != nil {
		return err
	}
	return provider.removeWebhook(ctx, sourceURL, remoteID)
}

func (g *github) syncWebhook(ctx context.Context, spec WebhookSpec, remoteID string) (string, error) {
	owner, repo, err := githubOwnerAndRepo(spec.SourceURL)
	if err != nil {
		return "", err
	}
	events := make([]string, 0, len(spec.Events)+1)
	seen := map[string]bool{}
	for _, event := range spec.Events {
		name := strings.ToLower(event)
		if !seen[name] {
			seen[name] = true
			events = append(events, name)
		}
	}
	// Upstream subscribes to issue_comment when a trigger posts pull-request
	// plans, which GitHub does not deliver through pull_request events.
	if spec.PullRequestWorkflow && !seen["issue_comment"] {
		events = append(events, "issue_comment")
	}

	rawURL := fmt.Sprintf("%s/repos/%s/%s/hooks", g.base, url.PathEscape(owner), url.PathEscape(repo))
	method := http.MethodPost
	body := map[string]any{"active": true, "events": events}
	if remoteID == "" {
		body["name"] = "web"
		body["config"] = map[string]any{
			"url":          spec.CallbackURL,
			"secret":       spec.Secret,
			"content_type": "json",
			// Upstream registers hooks this way so a self-signed Terrakube instance works.
			"insecure_ssl": "1",
		}
	} else {
		method = http.MethodPatch
		rawURL += "/" + url.PathEscape(remoteID)
	}

	var response struct {
		ID json.Number `json:"id"`
	}
	if _, err := doJSON(ctx, method, rawURL, g.token, g.headers(), body, &response); err != nil {
		return "", err
	}
	if remoteID != "" {
		return remoteID, nil
	}
	if response.ID == "" {
		return "", failed(http.StatusBadGateway, "The VCS provider did not return the new webhook id")
	}
	return response.ID.String(), nil
}

func (g *github) removeWebhook(ctx context.Context, sourceURL, remoteID string) error {
	owner, repo, err := githubOwnerAndRepo(sourceURL)
	if err != nil {
		return err
	}
	rawURL := remoteID
	if !strings.HasPrefix(remoteID, "http") {
		rawURL = fmt.Sprintf("%s/repos/%s/%s/hooks/%s", g.base, url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(remoteID))
	}
	if _, err := doJSON(ctx, http.MethodDelete, rawURL, g.token, g.headers(), nil, nil); err != nil {
		if providerErr, ok := err.(*Error); ok && providerErr.Status == http.StatusNotFound {
			return nil
		}
		return err
	}
	return nil
}

func githubOwnerAndRepo(sourceURL string) (string, string, error) {
	parsed, err := url.Parse(strings.TrimSpace(sourceURL))
	if err != nil {
		return "", "", failed(http.StatusUnprocessableEntity, "The workspace source URL is invalid: %v", err)
	}
	path := strings.TrimSuffix(strings.Trim(parsed.Path, "/"), ".git")
	segments := strings.Split(path, "/")
	if len(segments) < 2 || segments[0] == "" || segments[1] == "" {
		return "", "", failed(http.StatusUnprocessableEntity,
			"The workspace source URL does not name a GitHub repository (expected https://host/owner/repository)")
	}
	return segments[0], segments[1], nil
}

func (g *gitlab) syncWebhook(ctx context.Context, spec WebhookSpec, remoteID string) (string, error) {
	projectID, err := g.projectID(ctx, spec.SourceURL)
	if err != nil {
		return "", err
	}
	rawURL := fmt.Sprintf("%s/projects/%s/hooks", g.base, url.PathEscape(projectID))
	method := http.MethodPost
	if remoteID != "" {
		method = http.MethodPut
		rawURL += "/" + url.PathEscape(remoteID)
	}
	body := map[string]any{
		"url":                     spec.CallbackURL,
		"push_events":             true,
		"merge_requests_events":   true,
		"releases_events":         true,
		"enable_ssl_verification": false,
		"token":                   spec.Secret,
	}
	if spec.PullRequestWorkflow {
		body["note_events"] = true
	}

	var response struct {
		ID json.Number `json:"id"`
	}
	if _, err := doJSON(ctx, method, rawURL, g.token, g.headers(), body, &response); err != nil {
		return "", err
	}
	if remoteID != "" {
		return remoteID, nil
	}
	if response.ID == "" {
		return "", failed(http.StatusBadGateway, "The VCS provider did not return the new webhook id")
	}
	return response.ID.String(), nil
}

func (g *gitlab) removeWebhook(ctx context.Context, sourceURL, remoteID string) error {
	projectID, err := g.projectID(ctx, sourceURL)
	if err != nil {
		return err
	}
	rawURL := fmt.Sprintf("%s/projects/%s/hooks/%s", g.base, url.PathEscape(projectID), url.PathEscape(remoteID))
	if _, err := doJSON(ctx, http.MethodDelete, rawURL, g.token, g.headers(), nil, nil); err != nil {
		if providerErr, ok := err.(*Error); ok && providerErr.Status == http.StatusNotFound {
			return nil
		}
		return err
	}
	return nil
}

// projectID resolves a repository URL to GitLab's numeric project id. The
// encoded full path is the API's own addressing, so one lookup is enough;
// shortcut: rename-tolerant search is not implemented, add it if a project is
// only reachable by search.
func (g *gitlab) projectID(ctx context.Context, sourceURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(sourceURL))
	if err != nil {
		return "", failed(http.StatusUnprocessableEntity, "The workspace source URL is invalid: %v", err)
	}
	path := strings.TrimSuffix(strings.Trim(parsed.Path, "/"), ".git")
	if path == "" {
		return "", failed(http.StatusUnprocessableEntity,
			"The workspace source URL does not name a GitLab project")
	}
	var project struct {
		ID json.Number `json:"id"`
	}
	if _, err := doJSON(ctx, http.MethodGet, g.base+"/projects/"+url.PathEscape(path), g.token, g.headers(), nil, &project); err != nil {
		return "", err
	}
	if project.ID == "" {
		return "", failed(http.StatusBadGateway, "The VCS provider did not return the project id")
	}
	return project.ID.String(), nil
}
