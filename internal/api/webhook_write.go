package api

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	vcsprovider "github.com/aminbenmansour/terrakube-go/internal/vcs"
)

// webhookAttrs is what a webhook PATCH can change. remoteHookId is owned by
// the provider registration, so a caller cannot set it.
type webhookAttrs struct {
	MigratedV2 *bool `json:"migratedV2"`
}

// webhookEventAttrs mirrors upstream's WebhookEvent fields the UI edits.
type webhookEventAttrs struct {
	Branch            *string `json:"branch"`
	Path              *string `json:"path"`
	PathType          *string `json:"pathType"`
	Event             *string `json:"event"`
	TemplateID        *string `json:"templateId"`
	Priority          *int32  `json:"priority"`
	PrWorkflowEnabled *bool   `json:"prWorkflowEnabled"`
	PrApplyEnabled    *bool   `json:"prApplyEnabled"`
}

var webhookEventTypes = map[string]bool{
	"PUSH": true, "PULL_REQUEST": true, "PR_COMMENT": true, "PING": true, "RELEASE": true, "UNKNOWN": true,
}

var webhookPathTypes = map[string]bool{"PATTERN": true, "REGEX": true}

// workspaceForWebhook loads the workspace named in the route; the access
// middleware already checked the caller is a member of its organization.
func (a *API) workspaceForWebhook(w http.ResponseWriter, r *http.Request) (store.Workspace, bool) {
	ws, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return store.Workspace{}, false
	}
	return ws, true
}

func (a *API) createWebhook(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceForWebhook(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "webhook") {
		return
	}
	webhook, ok := a.createWebhookRow(w, r, ws, data, uuid.NewString())
	if !ok {
		return
	}
	// A direct POST has no triggers yet; the event calls that follow refresh
	// the registration, as they do for the UI's atomic save.
	jsonapi.WriteResource(w, http.StatusCreated, webhookResource(webhook), jsonapi.ParseFieldsets(r))
}

// createWebhookRow stores the webhook and checks the one invariant upstream's
// one-to-one relationship implies: a workspace has a single webhook.
func (a *API) createWebhookRow(w http.ResponseWriter, r *http.Request, ws store.Workspace, data *jsonapi.RequestData, fallbackID string) (store.Webhook, bool) {
	id := data.ID
	if id == "" {
		id = fallbackID
	}
	if _, err := uuid.Parse(id); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
		return store.Webhook{}, false
	}
	if _, err := a.deps.Store.GetWebhook(r.Context(), id); err == nil {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "A webhook with that id already exists.")
		return store.Webhook{}, false
	} else if !errors.Is(err, pgx.ErrNoRows) {
		a.internalError(w, "get webhook", err)
		return store.Webhook{}, false
	}
	if _, err := a.deps.Store.GetWebhookByWorkspace(r.Context(), ws.ID); err == nil {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This workspace already has a webhook.")
		return store.Webhook{}, false
	} else if !errors.Is(err, pgx.ErrNoRows) {
		a.internalError(w, "get workspace webhook", err)
		return store.Webhook{}, false
	}

	attrs, err := jsonapi.DecodeAttributes[webhookAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return store.Webhook{}, false
	}
	// Upstream refuses to save a webhook for a workspace without a VCS
	// connection, because there is no repository to register it on.
	connection, err := a.webhookConnection(r.Context(), ws)
	if err != nil {
		a.writeWebhookError(w, err)
		return store.Webhook{}, false
	}
	if err := vcsprovider.ValidateConnection(connection); err != nil {
		a.writeWebhookError(w, err)
		return store.Webhook{}, false
	}

	webhook, err := a.deps.Store.CreateWebhook(r.Context(), store.CreateWebhookParams{
		ID:          id,
		WorkspaceID: ws.ID,
		MigratedV2:  boolValue(attrs.MigratedV2, false),
	})
	if err != nil {
		if !a.writeStoreError(w, "create webhook", err) {
			a.internalError(w, "create webhook", err)
		}
		return store.Webhook{}, false
	}
	return webhook, true
}

func (a *API) updateWebhook(w http.ResponseWriter, r *http.Request) {
	webhook, ok := a.webhookInWorkspace(w, r)
	if !ok {
		return
	}
	webhook, ok = a.applyWebhookUpdate(w, r, webhook)
	if !ok {
		return
	}
	if webhook.MigratedV2 {
		// The shared repository hook owns the registration; drop the old
		// per-workspace hook so the two never deliver twice.
		if webhook, ok = a.dropWorkspaceHook(w, r, webhook); !ok {
			return
		}
	} else if webhook, ok = a.refreshWebhook(w, r, webhook); !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, webhookResource(webhook), jsonapi.ParseFieldsets(r))
}

// applyWebhookUpdate changes the stored attributes; migratedV2 is the only
// writable one. A PATCH without it is a no-op that still refreshes the
// registration, matching upstream's update hook.
func (a *API) applyWebhookUpdate(w http.ResponseWriter, r *http.Request, webhook store.Webhook) (store.Webhook, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return store.Webhook{}, false
	}
	if !requireType(w, data, "webhook") {
		return store.Webhook{}, false
	}
	return a.applyWebhookData(w, r, webhook, data)
}

func (a *API) applyWebhookData(w http.ResponseWriter, r *http.Request, webhook store.Webhook, data *jsonapi.RequestData) (store.Webhook, bool) {
	attrs, err := jsonapi.DecodeAttributes[webhookAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return store.Webhook{}, false
	}
	if attrs.MigratedV2 == nil || *attrs.MigratedV2 == webhook.MigratedV2 {
		return webhook, true
	}
	updated, err := a.deps.Store.UpdateWebhook(r.Context(), store.UpdateWebhookParams{
		ID:           webhook.ID,
		RemoteHookID: webhook.RemoteHookID,
		MigratedV2:   *attrs.MigratedV2,
		SyncError:    webhook.SyncError,
	})
	if err != nil {
		if !a.writeStoreError(w, "update webhook", err) {
			a.internalError(w, "update webhook", err)
		}
		return store.Webhook{}, false
	}
	return updated, true
}

func (a *API) deleteWebhook(w http.ResponseWriter, r *http.Request) {
	webhook, ok := a.webhookInWorkspace(w, r)
	if !ok {
		return
	}
	if !webhook.MigratedV2 {
		if _, ok := a.dropWorkspaceHook(w, r, webhook); !ok {
			return
		}
	}
	if err := a.deps.Store.DeleteWebhook(r.Context(), webhook.ID); err != nil {
		if !a.writeStoreError(w, "delete webhook", err) {
			a.internalError(w, "delete webhook", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) createWebhookEvent(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceForWebhook(w, r)
	if !ok {
		return
	}
	webhook, ok := a.webhookInWorkspace(w, r)
	if !ok {
		return
	}
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "webhook_event") {
		return
	}
	event, ok := a.createWebhookEventRow(w, r, ws, webhook, data, uuid.NewString())
	if !ok {
		return
	}
	if _, ok := a.refreshWebhook(w, r, webhook); !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, webhookEventResource(event), jsonapi.ParseFieldsets(r))
}

func (a *API) createWebhookEventRow(w http.ResponseWriter, r *http.Request, ws store.Workspace, webhook store.Webhook, data *jsonapi.RequestData, fallbackID string) (store.WebhookEvent, bool) {
	id := data.ID
	if id == "" {
		id = fallbackID
	}
	attrs, err := jsonapi.DecodeAttributes[webhookEventAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return store.WebhookEvent{}, false
	}
	values, ok := a.validateWebhookEvent(w, r, ws, webhook, attrs, store.WebhookEvent{})
	if !ok {
		return store.WebhookEvent{}, false
	}
	event, err := a.deps.Store.CreateWebhookEvent(r.Context(), store.CreateWebhookEventParams{
		ID:                id,
		WebhookID:         webhook.ID,
		Branch:            values.Branch,
		Path:              values.Path,
		PathType:          values.PathType,
		Event:             values.Event,
		TemplateID:        values.TemplateID,
		Priority:          values.Priority,
		PrWorkflowEnabled: values.PrWorkflowEnabled,
		PrApplyEnabled:    values.PrApplyEnabled,
	})
	if err != nil {
		if !a.writeStoreError(w, "create webhook event", err) {
			a.internalError(w, "create webhook event", err)
		}
		return store.WebhookEvent{}, false
	}
	return event, true
}

func (a *API) updateWebhookEvent(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceForWebhook(w, r)
	if !ok {
		return
	}
	webhook, event, ok := a.webhookEventInWorkspace(w, r)
	if !ok {
		return
	}
	event, ok = a.applyWebhookEventUpdate(w, r, ws, webhook, event)
	if !ok {
		return
	}
	if _, ok := a.refreshWebhook(w, r, webhook); !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, webhookEventResource(event), jsonapi.ParseFieldsets(r))
}

func (a *API) applyWebhookEventUpdate(w http.ResponseWriter, r *http.Request, ws store.Workspace, webhook store.Webhook, event store.WebhookEvent) (store.WebhookEvent, bool) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return store.WebhookEvent{}, false
	}
	if !requireType(w, data, "webhook_event") {
		return store.WebhookEvent{}, false
	}
	return a.applyWebhookEventData(w, r, ws, webhook, event, data)
}

func (a *API) applyWebhookEventData(w http.ResponseWriter, r *http.Request, ws store.Workspace, webhook store.Webhook, event store.WebhookEvent, data *jsonapi.RequestData) (store.WebhookEvent, bool) {
	attrs, err := jsonapi.DecodeAttributes[webhookEventAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return store.WebhookEvent{}, false
	}
	values, ok := a.validateWebhookEvent(w, r, ws, webhook, attrs, event)
	if !ok {
		return store.WebhookEvent{}, false
	}
	event, err = a.deps.Store.UpdateWebhookEvent(r.Context(), store.UpdateWebhookEventParams{
		ID:                event.ID,
		Branch:            values.Branch,
		Path:              values.Path,
		PathType:          values.PathType,
		Event:             values.Event,
		TemplateID:        values.TemplateID,
		Priority:          values.Priority,
		PrWorkflowEnabled: values.PrWorkflowEnabled,
		PrApplyEnabled:    values.PrApplyEnabled,
	})
	if err != nil {
		if !a.writeStoreError(w, "update webhook event", err) {
			a.internalError(w, "update webhook event", err)
		}
		return store.WebhookEvent{}, false
	}
	return event, true
}

func (a *API) deleteWebhookEvent(w http.ResponseWriter, r *http.Request) {
	webhook, event, ok := a.webhookEventInWorkspace(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteWebhookEvent(r.Context(), event.ID); err != nil {
		if !a.writeStoreError(w, "delete webhook event", err) {
			a.internalError(w, "delete webhook event", err)
		}
		return
	}
	if _, ok := a.refreshWebhook(w, r, webhook); !ok {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) webhookEventInWorkspace(w http.ResponseWriter, r *http.Request) (store.Webhook, store.WebhookEvent, bool) {
	webhook, ok := a.webhookInWorkspace(w, r)
	if !ok {
		return store.Webhook{}, store.WebhookEvent{}, false
	}
	event, err := a.deps.Store.GetWebhookEvent(r.Context(), param(r, "eventID"))
	if err != nil {
		if !a.writeStoreError(w, "get webhook event", err) {
			a.internalError(w, "get webhook event", err)
		}
		return store.Webhook{}, store.WebhookEvent{}, false
	}
	if event.WebhookID != webhook.ID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Webhook{}, store.WebhookEvent{}, false
	}
	return webhook, event, true
}

// validateWebhookEvent merges a PATCH over the stored trigger (empty on create)
// and checks the values upstream validates: a known event and path type, the
// required fields, and a template from the workspace's own organization.
func (a *API) validateWebhookEvent(w http.ResponseWriter, r *http.Request, ws store.Workspace, webhook store.Webhook, attrs webhookEventAttrs, existing store.WebhookEvent) (store.WebhookEvent, bool) {
	values := existing
	values.WebhookID = webhook.ID
	if attrs.Branch != nil {
		values.Branch = *attrs.Branch
	}
	if attrs.Path != nil {
		values.Path = *attrs.Path
	}
	if attrs.PathType != nil {
		values.PathType = *attrs.PathType
	}
	if attrs.Event != nil {
		values.Event = *attrs.Event
	}
	if attrs.TemplateID != nil {
		values.TemplateID = *attrs.TemplateID
	}
	if attrs.Priority != nil {
		values.Priority = *attrs.Priority
	}
	if attrs.PrWorkflowEnabled != nil {
		values.PrWorkflowEnabled = *attrs.PrWorkflowEnabled
	}
	if attrs.PrApplyEnabled != nil {
		values.PrApplyEnabled = *attrs.PrApplyEnabled
	}
	if values.PathType == "" {
		values.PathType = "REGEX"
	}

	switch {
	case !webhookEventTypes[values.Event]:
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "event must be one of PUSH, PULL_REQUEST, PR_COMMENT, PING, RELEASE or UNKNOWN.")
		return store.WebhookEvent{}, false
	case !webhookPathTypes[values.PathType]:
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "pathType must be PATTERN or REGEX.")
		return store.WebhookEvent{}, false
	case strings.TrimSpace(values.Branch) == "":
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "branch is required.")
		return store.WebhookEvent{}, false
	case strings.TrimSpace(values.Path) == "":
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "path is required.")
		return store.WebhookEvent{}, false
	case strings.TrimSpace(values.TemplateID) == "":
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "templateId is required.")
		return store.WebhookEvent{}, false
	}
	allowed, err := a.webhookTemplateAllowed(r.Context(), ws.OrganizationID, values.TemplateID)
	if err != nil {
		a.internalError(w, "check webhook template", err)
		return store.WebhookEvent{}, false
	}
	if !allowed {
		writePermissionDenied(w, r.Method)
		return store.WebhookEvent{}, false
	}
	return values, true
}

func (a *API) webhookTemplateAllowed(ctx context.Context, organizationID, templateID string) (bool, error) {
	template, err := a.deps.Store.GetTemplate(ctx, templateID)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return template.OrganizationID == organizationID, nil
}

// refreshWebhook registers the workspace's per-workspace hook with the final
// trigger set, or updates the existing one. A brand-new webhook without
// triggers has nothing to subscribe to yet.
func (a *API) refreshWebhook(w http.ResponseWriter, r *http.Request, webhook store.Webhook) (store.Webhook, bool) {
	updated, err := a.syncWebhook(r.Context(), webhook)
	if err != nil {
		a.writeWebhookError(w, err)
		return store.Webhook{}, false
	}
	return updated, true
}

func (a *API) syncWebhook(ctx context.Context, webhook store.Webhook) (store.Webhook, error) {
	if webhook.MigratedV2 {
		return a.dropWorkspaceHookContext(ctx, webhook)
	}
	events, err := a.deps.Store.ListWebhookEvents(ctx, webhook.ID)
	if err != nil {
		return store.Webhook{}, err
	}
	remoteID := ""
	if webhook.RemoteHookID != nil {
		remoteID = *webhook.RemoteHookID
	}
	if remoteID == "" && len(events) == 0 {
		return webhook, nil
	}
	ws, err := a.deps.Store.GetWorkspace(ctx, webhook.WorkspaceID)
	if err != nil {
		return store.Webhook{}, err
	}
	connection, err := a.webhookConnection(ctx, ws)
	if err != nil {
		return store.Webhook{}, err
	}
	spec := vcsprovider.WebhookSpec{
		SourceURL:   strValue(ws.Source),
		CallbackURL: strings.TrimRight(a.deps.PublicURL, "/") + "/webhook/v1/" + webhook.ID,
		Secret:      base64.StdEncoding.EncodeToString([]byte(ws.ID)),
	}
	for _, event := range events {
		spec.Events = append(spec.Events, event.Event)
		spec.PullRequestWorkflow = spec.PullRequestWorkflow || event.PrWorkflowEnabled
	}
	remote, err := vcsprovider.SyncWebhook(ctx, connection, spec, remoteID)
	if err != nil {
		return store.Webhook{}, err
	}
	if remote == remoteID {
		return webhook, nil
	}
	return a.deps.Store.UpdateWebhook(ctx, store.UpdateWebhookParams{
		ID: webhook.ID, RemoteHookID: strPtr(remote), MigratedV2: false, SyncError: webhook.SyncError,
	})
}

// dropWorkspaceHook removes the provider hook and clears the stored id. A hook
// that is already gone (the provider answers 404) counts as removed.
func (a *API) dropWorkspaceHook(w http.ResponseWriter, r *http.Request, webhook store.Webhook) (store.Webhook, bool) {
	updated, err := a.dropWorkspaceHookContext(r.Context(), webhook)
	if err != nil {
		a.writeWebhookError(w, err)
		return store.Webhook{}, false
	}
	return updated, true
}

func (a *API) dropWorkspaceHookContext(ctx context.Context, webhook store.Webhook) (store.Webhook, error) {
	if webhook.RemoteHookID == nil || *webhook.RemoteHookID == "" {
		return webhook, nil
	}
	ws, err := a.deps.Store.GetWorkspace(ctx, webhook.WorkspaceID)
	if err != nil {
		return store.Webhook{}, err
	}
	connection, err := a.webhookConnection(ctx, ws)
	if err != nil {
		return store.Webhook{}, err
	}
	if err := vcsprovider.DeleteWebhook(ctx, connection, strValue(ws.Source), *webhook.RemoteHookID); err != nil {
		return store.Webhook{}, err
	}
	return a.deps.Store.UpdateWebhook(ctx, store.UpdateWebhookParams{
		ID: webhook.ID, RemoteHookID: nil, MigratedV2: webhook.MigratedV2, SyncError: webhook.SyncError,
	})
}

// webhookConnection builds the provider client the registration calls need.
func (a *API) webhookConnection(ctx context.Context, ws store.Workspace) (vcsprovider.Connection, error) {
	if ws.VcsID == nil || *ws.VcsID == "" {
		return vcsprovider.Connection{}, &vcsprovider.Error{
			Status: http.StatusFailedDependency, Message: "No VCS defined for workspace",
		}
	}
	provider, err := a.deps.Store.GetVCS(ctx, *ws.VcsID)
	if err != nil {
		return vcsprovider.Connection{}, err
	}
	return vcsprovider.Connection{
		Type:           provider.VcsType,
		APIURL:         strValue(provider.ApiUrl),
		AccessToken:    strValue(provider.AccessToken),
		PrivateKey:     strValue(provider.PrivateKey),
		ConnectionType: strValue(provider.ConnectionType),
	}, nil
}

// writeWebhookError answers a provider failure the way upstream's lifecycle
// hook does: 424, except for a connection type this API cannot register on.
func (a *API) writeWebhookError(w http.ResponseWriter, err error) {
	var providerErr *vcsprovider.Error
	if errors.As(err, &providerErr) {
		status := http.StatusFailedDependency
		if providerErr.Status == http.StatusUnprocessableEntity || providerErr.Status == http.StatusNotImplemented {
			status = providerErr.Status
		}
		jsonapi.WriteError(w, status, http.StatusText(status), providerErr.Message)
		return
	}
	a.internalError(w, "webhook provider", err)
}
