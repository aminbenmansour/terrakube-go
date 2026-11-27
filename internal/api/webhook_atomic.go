package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// webhookSyncRef names a webhook whose provider registration the atomic save
// must reconcile once every trigger in the batch has been written.
type webhookSyncRef struct {
	id      string
	created bool
}

// parseWebhookHref reads the four paths the Webhook settings screen sends:
// the webhook collection and item, and the event collection and item.
func parseWebhookHref(href string) (orgID, wsID, webhookID, eventID string, ok bool) {
	segments := strings.Split(strings.Trim(strings.TrimPrefix(href, "/organization/"), "/"), "/")
	if len(segments) < 4 || len(segments) > 7 || segments[1] != "workspace" || segments[3] != "webhook" {
		return "", "", "", "", false
	}
	for _, index := range []int{0, 2} {
		if _, err := uuid.Parse(segments[index]); err != nil {
			return "", "", "", "", false
		}
	}
	orgID, wsID = segments[0], segments[2]
	if len(segments) >= 5 {
		if _, err := uuid.Parse(segments[4]); err != nil {
			return "", "", "", "", false
		}
		webhookID = segments[4]
	}
	switch len(segments) {
	case 4, 5:
		return orgID, wsID, webhookID, "", true
	case 7:
		if segments[5] != "events" {
			return "", "", "", "", false
		}
		if _, err := uuid.Parse(segments[6]); err != nil {
			return "", "", "", "", false
		}
		return orgID, wsID, webhookID, segments[6], true
	default: // 6 is the event collection.
		if segments[5] != "events" {
			return "", "", "", "", false
		}
		return orgID, wsID, webhookID, "", true
	}
}

func isWebhookAtomicHref(href string) bool {
	_, _, _, _, ok := parseWebhookHref(href)
	return ok
}

// atomicWebhook runs one webhook or webhook_event operation from the UI's
// atomic save. The provider call is deferred to syncAtomicWebhooks so a batch
// that creates a webhook and its triggers registers the finished trigger set
// once.
func (a *API) atomicWebhook(w http.ResponseWriter, r *http.Request, operation atomicOperation) (any, webhookSyncRef, bool) {
	orgID, wsID, webhookID, eventID, ok := parseWebhookHref(operation.Href)
	if !ok {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "Unsupported href: "+operation.Href)
		return nil, webhookSyncRef{}, false
	}
	ws, err := a.deps.Store.GetWorkspace(r.Context(), wsID)
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return nil, webhookSyncRef{}, false
	}
	if ws.OrganizationID != orgID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier "+wsID+" for workspace")
		return nil, webhookSyncRef{}, false
	}
	if !a.authorizeWorkspaceWrite(w, r, orgID, &ws) {
		return nil, webhookSyncRef{}, false
	}

	switch operation.Op {
	case "add":
		if operation.Data == nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data is required for an add")
			return nil, webhookSyncRef{}, false
		}
		if eventID == "" && webhookID == "" {
			if operation.Data.Type != "webhook" {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data.type must be webhook")
				return nil, webhookSyncRef{}, false
			}
			webhook, ok := a.createWebhookRow(w, r, ws, operation.Data, uuid.NewString())
			if !ok {
				return nil, webhookSyncRef{}, false
			}
			return webhookResource(webhook), webhookSyncRef{id: webhook.ID, created: true}, true
		}
		if eventID == "" {
			webhook, ok := a.atomicWebhookInWorkspace(w, r, wsID, webhookID)
			if !ok {
				return nil, webhookSyncRef{}, false
			}
			if operation.Data.Type != "webhook_event" {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data.type must be webhook_event")
				return nil, webhookSyncRef{}, false
			}
			event, ok := a.createWebhookEventRow(w, r, ws, webhook, operation.Data, uuid.NewString())
			if !ok {
				return nil, webhookSyncRef{}, false
			}
			return webhookEventResource(event), webhookSyncRef{id: webhook.ID}, true
		}
	case "update":
		if operation.Data == nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data is required for an update")
			return nil, webhookSyncRef{}, false
		}
		// The UI addresses an existing webhook by its collection href and
		// carries the id in data.
		if webhookID == "" {
			webhookID = operation.Data.ID
		}
		if webhookID == "" {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data.id is required for a webhook update")
			return nil, webhookSyncRef{}, false
		}
		if eventID == "" {
			if operation.Data.Type != "webhook" {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data.type must be webhook")
				return nil, webhookSyncRef{}, false
			}
			webhook, ok := a.atomicWebhookInWorkspace(w, r, wsID, webhookID)
			if !ok {
				return nil, webhookSyncRef{}, false
			}
			webhook, ok = a.applyWebhookData(w, r, webhook, operation.Data)
			if !ok {
				return nil, webhookSyncRef{}, false
			}
			return webhookResource(webhook), webhookSyncRef{id: webhook.ID}, true
		}
		webhook, event, ok := a.atomicWebhookEventInWorkspace(w, r, wsID, webhookID, eventID)
		if !ok {
			return nil, webhookSyncRef{}, false
		}
		if operation.Data.Type != "webhook_event" {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "data.type must be webhook_event")
			return nil, webhookSyncRef{}, false
		}
		event, ok = a.applyWebhookEventData(w, r, ws, webhook, event, operation.Data)
		if !ok {
			return nil, webhookSyncRef{}, false
		}
		return webhookEventResource(event), webhookSyncRef{id: webhook.ID}, true
	case "remove":
		if eventID != "" {
			webhook, event, ok := a.atomicWebhookEventInWorkspace(w, r, wsID, webhookID, eventID)
			if !ok {
				return nil, webhookSyncRef{}, false
			}
			if err := a.deps.Store.DeleteWebhookEvent(r.Context(), event.ID); err != nil {
				if !a.writeStoreError(w, "delete webhook event", err) {
					a.internalError(w, "delete webhook event", err)
				}
				return nil, webhookSyncRef{}, false
			}
			return nil, webhookSyncRef{id: webhook.ID}, true
		}
		if webhookID == "" {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "a webhook id is required for a remove")
			return nil, webhookSyncRef{}, false
		}
		webhook, ok := a.atomicWebhookInWorkspace(w, r, wsID, webhookID)
		if !ok {
			return nil, webhookSyncRef{}, false
		}
		if !webhook.MigratedV2 {
			if _, ok := a.dropWorkspaceHook(w, r, webhook); !ok {
				return nil, webhookSyncRef{}, false
			}
		}
		if err := a.deps.Store.DeleteWebhook(r.Context(), webhook.ID); err != nil {
			if !a.writeStoreError(w, "delete webhook", err) {
				a.internalError(w, "delete webhook", err)
			}
			return nil, webhookSyncRef{}, false
		}
		return nil, webhookSyncRef{}, true
	}
	jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Operation", "Unsupported op: "+operation.Op)
	return nil, webhookSyncRef{}, false
}

func (a *API) atomicWebhookInWorkspace(w http.ResponseWriter, r *http.Request, wsID, webhookID string) (store.Webhook, bool) {
	webhook, err := a.deps.Store.GetWebhook(r.Context(), webhookID)
	if err != nil {
		if !a.writeStoreError(w, "get webhook", err) {
			a.internalError(w, "get webhook", err)
		}
		return store.Webhook{}, false
	}
	if webhook.WorkspaceID != wsID {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Webhook{}, false
	}
	return webhook, true
}

func (a *API) atomicWebhookEventInWorkspace(w http.ResponseWriter, r *http.Request, wsID, webhookID, eventID string) (store.Webhook, store.WebhookEvent, bool) {
	webhook, ok := a.atomicWebhookInWorkspace(w, r, wsID, webhookID)
	if !ok {
		return store.Webhook{}, store.WebhookEvent{}, false
	}
	event, err := a.deps.Store.GetWebhookEvent(r.Context(), eventID)
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

// rollbackWebhooks removes webhooks created earlier in a failed batch. The
// cascade drops their triggers too.
func (a *API) rollbackWebhooks(r *http.Request, ids []string) {
	for _, id := range ids {
		_ = a.deps.Store.DeleteWebhook(r.Context(), id)
	}
}

// syncAtomicWebhooks registers or refreshes each touched webhook once the
// batch is stored. A webhook created by this batch is removed again when the
// provider rejects the registration, so a failed save leaves nothing behind.
func (a *API) syncAtomicWebhooks(w http.ResponseWriter, r *http.Request, refs []webhookSyncRef) bool {
	seen := map[string]bool{}
	for _, ref := range refs {
		if seen[ref.id] {
			continue
		}
		seen[ref.id] = true
		webhook, err := a.deps.Store.GetWebhook(r.Context(), ref.id)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			a.internalError(w, "get webhook", err)
			return false
		}
		if _, err := a.syncWebhook(r.Context(), webhook); err != nil {
			if ref.created {
				_ = a.deps.Store.DeleteWebhook(r.Context(), webhook.ID)
			}
			a.writeWebhookError(w, err)
			return false
		}
	}
	return true
}
