package api

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// webhookResource mirrors upstream's Webhook entity. syncError is why the last
// background sync of a shared repository webhook failed; the Webhook settings
// screen shows it as a warning.
func webhookResource(webhook store.Webhook) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "webhook",
		ID:   webhook.ID,
		Attributes: map[string]any{
			"migratedV2":   webhook.MigratedV2,
			"remoteHookId": webhook.RemoteHookID,
			"syncError":    webhook.SyncError,
			"createdBy":    webhook.CreatedBy,
			"createdDate":  timestamp(webhook.CreatedDate),
			"updatedBy":    webhook.UpdatedBy,
			"updatedDate":  timestamp(webhook.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"workspace": jsonapi.Link("workspace", webhook.WorkspaceID),
		},
	}
}

// webhookEventResource mirrors upstream's WebhookEvent entity: one trigger,
// checked in ascending priority order.
func webhookEventResource(event store.WebhookEvent) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "webhook_event",
		ID:   event.ID,
		Attributes: map[string]any{
			"branch":            event.Branch,
			"event":             event.Event,
			"path":              event.Path,
			"pathType":          event.PathType,
			"prApplyEnabled":    event.PrApplyEnabled,
			"prWorkflowEnabled": event.PrWorkflowEnabled,
			"priority":          event.Priority,
			"templateId":        event.TemplateID,
			"createdBy":         event.CreatedBy,
			"createdDate":       timestamp(event.CreatedDate),
			"updatedBy":         event.UpdatedBy,
			"updatedDate":       timestamp(event.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"webhook": jsonapi.Link("webhook", event.WebhookID),
		},
	}
}

func (a *API) getWebhook(w http.ResponseWriter, r *http.Request) {
	webhook, ok := a.webhookInWorkspace(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, webhookResource(webhook), jsonapi.ParseFieldsets(r))
}

func (a *API) listWebhookEvents(w http.ResponseWriter, r *http.Request) {
	webhook, ok := a.webhookInWorkspace(w, r)
	if !ok {
		return
	}
	events, err := a.deps.Store.ListWebhookEvents(r.Context(), webhook.ID)
	if err != nil {
		if !a.writeStoreError(w, "list webhook events", err) {
			a.internalError(w, "list webhook events", err)
		}
		return
	}
	resources := make([]jsonapi.Resource, 0, len(events))
	for _, event := range events {
		resources = append(resources, webhookEventResource(event))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

// webhookInWorkspace answers a webhook that belongs to another workspace the
// way Elide does for an id outside the relationship: 404, not 403.
func (a *API) webhookInWorkspace(w http.ResponseWriter, r *http.Request) (store.Webhook, bool) {
	webhook, err := a.deps.Store.GetWebhook(r.Context(), param(r, "webhookID"))
	if err != nil {
		if !a.writeStoreError(w, "get webhook", err) {
			a.internalError(w, "get webhook", err)
		}
		return store.Webhook{}, false
	}
	if webhook.WorkspaceID != param(r, "workspaceID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Webhook{}, false
	}
	return webhook, true
}

// workspaceWebhook loads the workspace's webhook for the `webhook`
// relationship; a workspace without one is not an error.
func (a *API) workspaceWebhook(r *http.Request, workspaceID string) (*store.Webhook, error) {
	webhook, err := a.deps.Store.GetWebhookByWorkspace(r.Context(), workspaceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}
