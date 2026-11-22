package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	notificationOrgID       = "51515151-5151-4151-8151-515151515151"
	notificationOtherOrgID  = "52525252-5252-4252-8252-525252525252"
	notificationWorkspaceID = "53535353-5353-4353-8353-535353535353"
	notificationTemplateID  = "54545454-5454-4454-8454-545454545454"
	notificationConfigID    = "55555555-5555-4555-8555-555555555555"
	notificationTriggerID   = "56565656-5656-4656-8656-565656565656"
	notificationMissingID   = "57575757-5757-4757-8757-575757575757"
)

func newNotificationStub() *stubQuerier {
	return &stubQuerier{
		org:        store.Organization{ID: notificationOrgID},
		workspaces: []store.Workspace{{ID: notificationWorkspaceID, OrganizationID: notificationOrgID}},
		templates:  []store.Template{{ID: notificationTemplateID, OrganizationID: notificationOrgID}},
		notifications: []store.NotificationConfiguration{{
			ID: notificationConfigID, Name: "alerts", OrganizationID: notificationOrgID,
			ChannelType: "WEBHOOK", DestinationUrl: "https://example.com/hook",
			Active: true, MessageStyle: "DETAILED",
		}},
		triggers: []store.NotificationTrigger{{
			ID: notificationTriggerID, JobStatus: "failed", ConfigurationID: notificationConfigID,
		}},
	}
}

func notificationTestAPI(q store.Querier) http.Handler {
	return New(Deps{Store: q})
}

func TestNotificationConfigurationJSONAPI(t *testing.T) {
	q := newNotificationStub()
	h := notificationTestAPI(q)
	root := "/notification_configuration/" + notificationConfigID

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+notificationOrgID+"/notificationConfiguration", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list = %d; body = %s", rec.Code, rec.Body.String())
	}
	var list struct {
		Data []struct {
			Type       string         `json:"type"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 1 || list.Data[0].Type != "notification_configuration" {
		t.Fatalf("list = %+v", list.Data)
	}
	if list.Data[0].Attributes["name"] != "alerts" || list.Data[0].Attributes["active"] != true {
		t.Errorf("attributes = %+v", list.Data[0].Attributes)
	}
	if _, leaked := list.Data[0].Attributes["signingSecret"]; leaked {
		t.Error("signingSecret must never be returned")
	}

	// Missing required columns answer upstream's data-integrity 409.
	for _, body := range []string{
		`{"data":{"type":"notification_configuration","attributes":{"channelType":"WEBHOOK","destinationUrl":"https://x"}}}`,
		`{"data":{"type":"notification_configuration","attributes":{"name":"n","destinationUrl":"https://x"}}}`,
		`{"data":{"type":"notification_configuration","attributes":{"name":"n","channelType":"WEBHOOK"}}}`,
	} {
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+notificationOrgID+"/notificationConfiguration", strings.NewReader(body)))
		if rec.Code != http.StatusConflict {
			t.Errorf("missing attribute = %d, want 409; body = %s", rec.Code, rec.Body.String())
		}
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+notificationOrgID+"/notificationConfiguration",
		strings.NewReader(`{"data":{"type":"notification_configuration","attributes":{"name":"slack","channelType":"SLACK","destinationUrl":"https://hooks.slack.com/x","signingSecret":"s3cret"}}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create = %d; body = %s", rec.Code, rec.Body.String())
	}
	created := q.notifications[len(q.notifications)-1]
	if !created.Active || created.MessageStyle != "DETAILED" || created.WorkspaceID != nil {
		t.Errorf("defaults = %+v", created)
	}
	if strings.Contains(rec.Body.String(), "s3cret") || strings.Contains(rec.Body.String(), "signingSecret") {
		t.Errorf("create leaked the signing secret: %s", rec.Body.String())
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+notificationOrgID+"/notificationConfiguration",
		strings.NewReader(`{"data":{"type":"notification_configuration","attributes":{"name":"bad","channelType":"NOPE","destinationUrl":"https://x"}}}`)))
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: NOPE") {
		t.Errorf("invalid channel = %d; body = %s", rec.Code, rec.Body.String())
	}

	// A workspace-scoped configuration inherits both ids from the nested URL.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+notificationOrgID+"/workspace/"+notificationWorkspaceID+"/notificationConfiguration",
		strings.NewReader(`{"data":{"type":"notification_configuration","attributes":{"name":"ws","channelType":"TEAMS","destinationUrl":"https://teams/x"}}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("workspace create = %d; body = %s", rec.Code, rec.Body.String())
	}
	scoped := q.notifications[len(q.notifications)-1]
	if scoped.WorkspaceID == nil || *scoped.WorkspaceID != notificationWorkspaceID || scoped.OrganizationID != notificationOrgID {
		t.Errorf("workspace create = %+v", scoped)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, root,
		strings.NewReader(`{"data":{"type":"notification_configuration","id":"`+notificationConfigID+`","attributes":{"description":null,"active":false}}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch = %d; body = %s", rec.Code, rec.Body.String())
	}
	updated := q.notifications[0]
	if updated.Name != "alerts" || updated.Active || updated.Description != nil {
		t.Errorf("patch params = %+v", updated)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, root,
		strings.NewReader(`{"data":{"type":"notification_configuration","id":"`+notificationConfigID+`","attributes":{"name":null}}}`)))
	if rec.Code != http.StatusConflict {
		t.Errorf("null name = %d, want 409", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/notification_configuration/"+notificationMissingID, nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("missing get = %d, want 404", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/notification_configuration/not-a-uuid", nil))
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: not-a-uuid") {
		t.Errorf("bad id = %d; body = %s", rec.Code, rec.Body.String())
	}
}

func TestNotificationTriggersAndTemplates(t *testing.T) {
	q := newNotificationStub()
	h := notificationTestAPI(q)
	root := "/notification_configuration/" + notificationConfigID

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, root+"/triggers",
		strings.NewReader(`{"data":{"type":"notification_trigger","attributes":{"jobStatus":"completed"}}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("trigger create = %d; body = %s", rec.Code, rec.Body.String())
	}
	if len(q.triggers) != 2 || q.triggers[1].JobStatus != "completed" || q.triggers[1].ConfigurationID != notificationConfigID {
		t.Errorf("triggers = %+v", q.triggers)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, root+"/triggers",
		strings.NewReader(`{"data":{"type":"notification_trigger","attributes":{}}}`)))
	if rec.Code != http.StatusConflict {
		t.Errorf("trigger without status = %d, want 409", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, root+"/triggers",
		strings.NewReader(`{"data":{"type":"notification_trigger","attributes":{"jobStatus":"nope"}}}`)))
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: nope") {
		t.Errorf("trigger invalid status = %d; body = %s", rec.Code, rec.Body.String())
	}

	// A trigger that belongs to another configuration is not addressable here.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/notification_configuration/"+notificationMissingID+"/triggers/"+notificationTriggerID, nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("cross-config trigger delete = %d, want 404", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, root+"/triggers/"+notificationTriggerID, nil))
	if rec.Code != http.StatusNoContent || len(q.triggers) != 1 {
		t.Errorf("trigger delete = %d, triggers = %+v", rec.Code, q.triggers)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, root+"/relationships/templates",
		strings.NewReader(`{"data":[{"type":"template","id":"`+notificationTemplateID+`"}]}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("templates patch = %d; body = %s", rec.Code, rec.Body.String())
	}
	if ids := q.notificationTemplates[notificationConfigID]; len(ids) != 1 || ids[0] != notificationTemplateID {
		t.Errorf("templates = %v", ids)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, root+"/relationships/templates",
		strings.NewReader(`{"data":[{"type":"template","id":"`+notificationMissingID+`"}]}`)))
	if rec.Code != http.StatusNotFound {
		t.Errorf("unknown template = %d, want 404", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, root+"/relationships/templates",
		strings.NewReader(`{"data":[{"type":"organization","id":"`+notificationOrgID+`"}]}`)))
	if rec.Code != http.StatusConflict {
		t.Errorf("wrong relationship type = %d, want 409", rec.Code)
	}

	// An empty list clears the filter again (every template applies).
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, root+"/relationships/templates", strings.NewReader(`{"data":[]}`)))
	if rec.Code != http.StatusNoContent || len(q.notificationTemplates[notificationConfigID]) != 0 {
		t.Errorf("clear templates = %d, %v", rec.Code, q.notificationTemplates[notificationConfigID])
	}
}

func TestDeleteNotificationConfigurationCascadesTriggers(t *testing.T) {
	q := newNotificationStub()
	h := notificationTestAPI(q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/notification_configuration/"+notificationConfigID, nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete = %d; body = %s", rec.Code, rec.Body.String())
	}
	if len(q.notifications) != 0 || len(q.triggers) != 0 {
		t.Errorf("after delete: configs = %+v, triggers = %+v", q.notifications, q.triggers)
	}
}
