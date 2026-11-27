package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	webhookID           = "41414141-4141-4141-8141-414141414141"
	webhookEventID      = "42424242-4242-4242-8242-424242424242"
	webhookOtherID      = "43434343-4343-4343-8343-434343434343"
	webhookTemplateID   = "44444444-4444-4444-8444-444444444444"
	webhookOtherFetchID = "45454545-4545-4545-8545-454545454545"
)

func (s *stubQuerier) GetWebhook(_ context.Context, id string) (store.Webhook, error) {
	for _, webhook := range s.webhooks {
		if webhook.ID == id {
			return webhook, nil
		}
	}
	return store.Webhook{}, pgx.ErrNoRows
}

func (s *stubQuerier) GetWebhookByWorkspace(_ context.Context, workspaceID string) (store.Webhook, error) {
	for _, webhook := range s.webhooks {
		if webhook.WorkspaceID == workspaceID {
			return webhook, nil
		}
	}
	return store.Webhook{}, pgx.ErrNoRows
}

func (s *stubQuerier) ListWebhookEvents(_ context.Context, webhookID string) ([]store.WebhookEvent, error) {
	events := []store.WebhookEvent{}
	for _, event := range s.webhookEvents {
		if event.WebhookID == webhookID {
			events = append(events, event)
		}
	}
	return events, nil
}

func newWebhookStub() *stubQuerier {
	remoteHookID := "987654"
	return &stubQuerier{
		org:        store.Organization{ID: vcsOrgID},
		workspaces: []store.Workspace{{ID: vcsWSID, OrganizationID: vcsOrgID}},
		webhooks: []store.Webhook{{
			ID: webhookID, WorkspaceID: vcsWSID, RemoteHookID: &remoteHookID,
		}},
		webhookEvents: []store.WebhookEvent{
			{
				ID: webhookEventID, WebhookID: webhookID, Branch: "main", Path: "modules/**",
				PathType: "PATTERN", Event: "PULL_REQUEST", TemplateID: webhookTemplateID,
				Priority: 1, PrWorkflowEnabled: true,
			},
			{
				ID: webhookOtherFetchID, WebhookID: webhookID, Branch: ".*", Path: "**",
				PathType: "PATTERN", Event: "PUSH", TemplateID: webhookTemplateID, Priority: 2,
			},
		},
	}
}

func TestWebhookReadAndEvents(t *testing.T) {
	h := New(Deps{Store: newWebhookStub()})
	base := "/organization/" + vcsOrgID + "/workspace/" + vcsWSID + "/webhook/" + webhookID

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, base, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("webhook = %d; body = %s", rec.Code, rec.Body.String())
	}
	var document struct {
		Data struct {
			Type       string         `json:"type"`
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &document); err != nil {
		t.Fatalf("decode: %v; body = %s", err, rec.Body.String())
	}
	if document.Data.Type != "webhook" || document.Data.ID != webhookID {
		t.Fatalf("data = %s/%s", document.Data.Type, document.Data.ID)
	}
	if document.Data.Attributes["remoteHookId"] != "987654" || document.Data.Attributes["migratedV2"] != false {
		t.Fatalf("attributes = %+v", document.Data.Attributes)
	}
	if document.Data.Attributes["syncError"] != nil {
		t.Fatalf("syncError = %v, want null", document.Data.Attributes["syncError"])
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, base+"/events", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("events = %d; body = %s", rec.Code, rec.Body.String())
	}
	var events struct {
		Data []struct {
			Type       string         `json:"type"`
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &events); err != nil {
		t.Fatalf("decode events: %v; body = %s", err, rec.Body.String())
	}
	if len(events.Data) != 2 || events.Data[0].ID != webhookEventID || events.Data[1].ID != webhookOtherFetchID {
		t.Fatalf("events = %+v", events.Data)
	}
	first := events.Data[0]
	if first.Type != "webhook_event" || first.Attributes["event"] != "PULL_REQUEST" ||
		first.Attributes["pathType"] != "PATTERN" || first.Attributes["templateId"] != webhookTemplateID ||
		first.Attributes["priority"] != float64(1) || first.Attributes["prWorkflowEnabled"] != true ||
		first.Attributes["prApplyEnabled"] != false {
		t.Fatalf("first event = %+v", first)
	}
}

func TestWebhookInAnotherWorkspace(t *testing.T) {
	q := newWebhookStub()
	q.webhooks = append(q.webhooks, store.Webhook{ID: webhookOtherID, WorkspaceID: webhookOtherFetchID})
	h := New(Deps{Store: q})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/organization/"+vcsOrgID+"/workspace/"+vcsWSID+"/webhook/"+webhookOtherID, nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("foreign webhook = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/organization/"+vcsOrgID+"/workspace/"+vcsWSID+"/webhook/00000000-0000-4000-8000-000000000001", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown webhook = %d, want 404", rec.Code)
	}
}

func TestWorkspaceWebhookRelationship(t *testing.T) {
	h := New(Deps{Store: newWebhookStub()})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/organization/"+vcsOrgID+"/workspace/"+vcsWSID+"?include=webhook", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("workspace = %d; body = %s", rec.Code, rec.Body.String())
	}
	var document struct {
		Data struct {
			Relationships map[string]struct {
				Data *struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"relationships"`
		} `json:"data"`
		Included []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"included"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &document); err != nil {
		t.Fatalf("decode: %v; body = %s", err, rec.Body.String())
	}
	relationship := document.Data.Relationships["webhook"].Data
	if relationship == nil || relationship.ID != webhookID || relationship.Type != "webhook" {
		t.Fatalf("webhook relationship = %+v", relationship)
	}
	if len(document.Included) != 1 || document.Included[0].Type != "webhook" || document.Included[0].ID != webhookID {
		t.Fatalf("included = %+v", document.Included)
	}

	// A workspace without a webhook reports null, which the UI reads as "off".
	empty := &stubQuerier{org: store.Organization{ID: vcsOrgID},
		workspaces: []store.Workspace{{ID: vcsWSID, OrganizationID: vcsOrgID}}}
	rec = httptest.NewRecorder()
	New(Deps{Store: empty}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/organization/"+vcsOrgID+"/workspace/"+vcsWSID+"?include=webhook", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("empty workspace = %d; body = %s", rec.Code, rec.Body.String())
	}
	document = struct {
		Data struct {
			Relationships map[string]struct {
				Data *struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"relationships"`
		} `json:"data"`
		Included []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"included"`
	}{}
	if err := json.Unmarshal(rec.Body.Bytes(), &document); err != nil {
		t.Fatalf("decode empty: %v; body = %s", err, rec.Body.String())
	}
	if document.Data.Relationships["webhook"].Data != nil {
		t.Fatalf("webhook relationship = %+v, want null", document.Data.Relationships["webhook"].Data)
	}
	if len(document.Included) != 0 {
		t.Fatalf("included = %+v, want none", document.Included)
	}
}
