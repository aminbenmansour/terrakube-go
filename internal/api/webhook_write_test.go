package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const webhookVcsID = "46464646-4646-4646-8646-464646464646"

// fakeHooks records the provider calls a webhook save makes.
type fakeHooks struct {
	method string
	path   string
	body   map[string]any
	calls  int
}

func (f *fakeHooks) server(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.calls++
		f.method = r.Method
		f.path = r.URL.EscapedPath()
		raw, _ := io.ReadAll(r.Body)
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &f.body)
		}
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id": 55}`))
	}))
	t.Cleanup(server.Close)
	return server
}

func webhookWriteStub(apiURL string) *stubQuerier {
	source := "https://github.com/acme/infra.git"
	token := "gh-token"
	return &stubQuerier{
		org: store.Organization{ID: vcsOrgID},
		workspaces: []store.Workspace{{
			ID: vcsWSID, OrganizationID: vcsOrgID, VcsID: strPtr(webhookVcsID),
			Source: &source, IacType: "terraform",
		}},
		vcs: []store.Vc{{
			ID: webhookVcsID, OrganizationID: vcsOrgID, VcsType: "GITHUB",
			ApiUrl: &apiURL, AccessToken: &token,
		}},
		templates: []store.Template{{ID: webhookTemplateID, OrganizationID: vcsOrgID}},
	}
}

func (s *stubQuerier) CreateWebhook(_ context.Context, arg store.CreateWebhookParams) (store.Webhook, error) {
	webhook := store.Webhook{
		ID: arg.ID, WorkspaceID: arg.WorkspaceID, RemoteHookID: arg.RemoteHookID, MigratedV2: arg.MigratedV2,
	}
	s.webhooks = append(s.webhooks, webhook)
	return webhook, nil
}

func (s *stubQuerier) UpdateWebhook(_ context.Context, arg store.UpdateWebhookParams) (store.Webhook, error) {
	for i := range s.webhooks {
		if s.webhooks[i].ID == arg.ID {
			s.webhooks[i].RemoteHookID = arg.RemoteHookID
			s.webhooks[i].MigratedV2 = arg.MigratedV2
			s.webhooks[i].SyncError = arg.SyncError
			return s.webhooks[i], nil
		}
	}
	return store.Webhook{}, pgx.ErrNoRows
}

func (s *stubQuerier) DeleteWebhook(_ context.Context, id string) error {
	kept := s.webhooks[:0]
	for _, webhook := range s.webhooks {
		if webhook.ID != id {
			kept = append(kept, webhook)
		}
	}
	s.webhooks = kept
	keptEvents := s.webhookEvents[:0]
	for _, event := range s.webhookEvents {
		if event.WebhookID != id {
			keptEvents = append(keptEvents, event)
		}
	}
	s.webhookEvents = keptEvents
	return nil
}

func (s *stubQuerier) CreateWebhookEvent(_ context.Context, arg store.CreateWebhookEventParams) (store.WebhookEvent, error) {
	event := store.WebhookEvent{
		ID: arg.ID, WebhookID: arg.WebhookID, Branch: arg.Branch, Path: arg.Path, PathType: arg.PathType,
		Event: arg.Event, TemplateID: arg.TemplateID, Priority: arg.Priority,
		PrWorkflowEnabled: arg.PrWorkflowEnabled, PrApplyEnabled: arg.PrApplyEnabled,
	}
	s.webhookEvents = append(s.webhookEvents, event)
	return event, nil
}

func (s *stubQuerier) UpdateWebhookEvent(_ context.Context, arg store.UpdateWebhookEventParams) (store.WebhookEvent, error) {
	for i := range s.webhookEvents {
		if s.webhookEvents[i].ID == arg.ID {
			s.webhookEvents[i].Branch = arg.Branch
			s.webhookEvents[i].Path = arg.Path
			s.webhookEvents[i].PathType = arg.PathType
			s.webhookEvents[i].Event = arg.Event
			s.webhookEvents[i].TemplateID = arg.TemplateID
			s.webhookEvents[i].Priority = arg.Priority
			s.webhookEvents[i].PrWorkflowEnabled = arg.PrWorkflowEnabled
			s.webhookEvents[i].PrApplyEnabled = arg.PrApplyEnabled
			return s.webhookEvents[i], nil
		}
	}
	return store.WebhookEvent{}, pgx.ErrNoRows
}

func (s *stubQuerier) DeleteWebhookEvent(_ context.Context, id string) error {
	kept := s.webhookEvents[:0]
	for _, event := range s.webhookEvents {
		if event.ID != id {
			kept = append(kept, event)
		}
	}
	s.webhookEvents = kept
	return nil
}

func (s *stubQuerier) GetWebhookEvent(_ context.Context, id string) (store.WebhookEvent, error) {
	for _, event := range s.webhookEvents {
		if event.ID == id {
			return event, nil
		}
	}
	return store.WebhookEvent{}, pgx.ErrNoRows
}

func TestCreateWebhookEventRegistersProvider(t *testing.T) {
	hooks := &fakeHooks{}
	server := hooks.server(t)
	stub := webhookWriteStub(server.URL)
	h := New(Deps{Store: stub, PublicURL: "https://terrakube.example.com"})
	base := "/organization/" + vcsOrgID + "/workspace/" + vcsWSID + "/webhook"

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, base,
		strings.NewReader(`{"data":{"type":"webhook","id":"`+webhookID+`"}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create webhook = %d; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, base+"/"+webhookID+"/events",
		strings.NewReader(`{"data":{"type":"webhook_event","id":"`+webhookEventID+`","attributes":{`+
			`"event":"PUSH","branch":"main","path":"modules/**","pathType":"PATTERN",`+
			`"templateId":"`+webhookTemplateID+`","priority":1,"prWorkflowEnabled":false,"prApplyEnabled":false}}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create event = %d; body = %s", rec.Code, rec.Body.String())
	}
	if hooks.calls != 1 || hooks.method != http.MethodPost || hooks.path != "/repos/acme/infra/hooks" {
		t.Fatalf("provider call = %d %s %s", hooks.calls, hooks.method, hooks.path)
	}
	events, _ := hooks.body["events"].([]any)
	if len(events) != 1 || events[0] != "push" {
		t.Fatalf("subscribed events = %+v", events)
	}
	config, _ := hooks.body["config"].(map[string]any)
	// The secret is derived from the workspace id and only travels to the
	// provider; the API never returns it.
	if config["secret"] != "MjUyNTI1MjUtMjUyNS00NTI1LTg1MjUtMjUyNTI1MjUyNTI1" {
		t.Fatalf("secret = %v", config["secret"])
	}
	if config["url"] != "https://terrakube.example.com/webhook/v1/"+webhookID {
		t.Fatalf("callback = %v", config["url"])
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, base+"/"+webhookID, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("read webhook = %d; body = %s", rec.Code, rec.Body.String())
	}
	var document struct {
		Data struct {
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &document); err != nil {
		t.Fatalf("decode webhook: %v", err)
	}
	if document.Data.Attributes["remoteHookId"] != "55" {
		t.Fatalf("remoteHookId = %v", document.Data.Attributes["remoteHookId"])
	}
	for name := range document.Data.Attributes {
		if strings.Contains(strings.ToLower(name), "secret") || strings.Contains(strings.ToLower(name), "token") {
			t.Fatalf("response exposes %q", name)
		}
	}
}

func TestDeleteWebhookRemovesProviderHook(t *testing.T) {
	hooks := &fakeHooks{}
	server := hooks.server(t)
	remote := "55"
	stub := webhookWriteStub(server.URL)
	stub.webhooks = []store.Webhook{{ID: webhookID, WorkspaceID: vcsWSID, RemoteHookID: &remote}}
	h := New(Deps{Store: stub, PublicURL: "https://terrakube.example.com"})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete,
		"/organization/"+vcsOrgID+"/workspace/"+vcsWSID+"/webhook/"+webhookID, nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete webhook = %d; body = %s", rec.Code, rec.Body.String())
	}
	if hooks.method != http.MethodDelete || hooks.path != "/repos/acme/infra/hooks/55" {
		t.Fatalf("provider call = %s %s", hooks.method, hooks.path)
	}
	if len(stub.webhooks) != 0 {
		t.Fatalf("webhooks after delete = %+v", stub.webhooks)
	}
}

func TestAtomicWebhookSaveRegistersFinalTriggers(t *testing.T) {
	hooks := &fakeHooks{}
	server := hooks.server(t)
	stub := webhookWriteStub(server.URL)
	h := New(Deps{Store: stub, PublicURL: "https://terrakube.example.com"})

	body := `{"atomic:operations":[
		{"op":"add","href":"/organization/` + vcsOrgID + `/workspace/` + vcsWSID + `/webhook",
		 "data":{"type":"webhook","id":"` + webhookID + `"}},
		{"op":"add","href":"/organization/` + vcsOrgID + `/workspace/` + vcsWSID + `/webhook/` + webhookID + `/events",
		 "data":{"type":"webhook_event","id":"` + webhookEventID + `","attributes":{
			"event":"PULL_REQUEST","branch":".*","path":"**","pathType":"PATTERN",
			"templateId":"` + webhookTemplateID + `","priority":1,"prWorkflowEnabled":true,"prApplyEnabled":false}}}
	]}`
	request := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/vnd.api+json; ext=jsonapi.org/ext/atomic")
	request.Header.Set("Accept", "application/vnd.api+json; ext=jsonapi.org/ext/atomic")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, request)
	if rec.Code != http.StatusOK {
		t.Fatalf("atomic save = %d; body = %s", rec.Code, rec.Body.String())
	}
	// One provider call for the finished batch, not one per operation.
	if hooks.calls != 1 || hooks.method != http.MethodPost {
		t.Fatalf("provider calls = %d, last = %s", hooks.calls, hooks.method)
	}
	events, _ := hooks.body["events"].([]any)
	if len(events) != 2 || events[0] != "pull_request" || events[1] != "issue_comment" {
		t.Fatalf("subscribed events = %+v", events)
	}
	if len(stub.webhooks) != 1 || stub.webhooks[0].RemoteHookID == nil || *stub.webhooks[0].RemoteHookID != "55" {
		t.Fatalf("stored webhook = %+v", stub.webhooks)
	}
	if len(stub.webhookEvents) != 1 {
		t.Fatalf("stored events = %+v", stub.webhookEvents)
	}
}

func TestAtomicWebhookInvalidTriggerChangesNothing(t *testing.T) {
	hooks := &fakeHooks{}
	server := hooks.server(t)
	stub := webhookWriteStub(server.URL)
	h := New(Deps{Store: stub, PublicURL: "https://terrakube.example.com"})

	body := `{"atomic:operations":[
		{"op":"add","href":"/organization/` + vcsOrgID + `/workspace/` + vcsWSID + `/webhook",
		 "data":{"type":"webhook","id":"` + webhookID + `"}},
		{"op":"add","href":"/organization/` + vcsOrgID + `/workspace/` + vcsWSID + `/webhook/` + webhookID + `/events",
		 "data":{"type":"webhook_event","id":"` + webhookEventID + `","attributes":{
			"event":"NOT_AN_EVENT","branch":"main","path":"**","pathType":"PATTERN",
			"templateId":"` + webhookTemplateID + `","priority":1}}}
	]}`
	request := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/vnd.api+json; ext=jsonapi.org/ext/atomic")
	request.Header.Set("Accept", "application/vnd.api+json; ext=jsonapi.org/ext/atomic")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, request)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("atomic save = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if len(stub.webhooks) != 0 || len(stub.webhookEvents) != 0 {
		t.Fatalf("stored rows after failure = %+v / %+v", stub.webhooks, stub.webhookEvents)
	}
	if hooks.calls != 0 {
		t.Fatalf("provider was called %d times for an invalid trigger", hooks.calls)
	}
}

func TestWebhookEventTemplateOutsideOrganizationIsRejected(t *testing.T) {
	hooks := &fakeHooks{}
	server := hooks.server(t)
	stub := webhookWriteStub(server.URL)
	stub.templates = []store.Template{{ID: webhookTemplateID, OrganizationID: "other-org"}}
	stub.webhooks = []store.Webhook{{ID: webhookID, WorkspaceID: vcsWSID}}
	h := New(Deps{Store: stub, PublicURL: "https://terrakube.example.com"})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost,
		"/organization/"+vcsOrgID+"/workspace/"+vcsWSID+"/webhook/"+webhookID+"/events",
		strings.NewReader(`{"data":{"type":"webhook_event","id":"`+webhookEventID+`","attributes":{`+
			`"event":"PUSH","branch":"main","path":"**","pathType":"PATTERN","templateId":"`+webhookTemplateID+`"}}}`)))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("create event = %d, want 403; body = %s", rec.Code, rec.Body.String())
	}
	if len(stub.webhookEvents) != 0 || hooks.calls != 0 {
		t.Fatalf("stored events = %+v, provider calls = %d", stub.webhookEvents, hooks.calls)
	}
}
