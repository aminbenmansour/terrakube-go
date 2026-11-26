package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	triggerOrg         = "70000000-0000-4000-8000-000000000001"
	triggerOtherOrg    = "70000000-0000-4000-8000-000000000009"
	triggerSource      = "70000000-0000-4000-8000-000000000002"
	triggerDestination = "70000000-0000-4000-8000-000000000003"
	triggerThird       = "70000000-0000-4000-8000-000000000004"
	triggerTemplate    = "70000000-0000-4000-8000-000000000005"
	triggerForeign     = "70000000-0000-4000-8000-000000000006"
	triggerEdgeA       = "70000000-0000-4000-8000-0000000000a1"
	triggerEdgeB       = "70000000-0000-4000-8000-0000000000a2"
)

// runTriggerStub borrows the management API's RBAC fakes and overrides the
// workspace and template lookups the trigger handlers need.
type runTriggerStub struct {
	*stubQuerier
	workspaces []store.Workspace
	templates  []store.Template
}

func (s *runTriggerStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, workspace := range s.workspaces {
		if workspace.ID == id {
			return workspace, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *runTriggerStub) GetTemplate(_ context.Context, id string) (store.Template, error) {
	for _, template := range s.templates {
		if template.ID == id {
			return template, nil
		}
	}
	return store.Template{}, pgx.ErrNoRows
}

func newRunTriggerStub() *runTriggerStub {
	return &runTriggerStub{
		stubQuerier: &stubQuerier{},
		workspaces: []store.Workspace{
			{ID: triggerSource, Name: "source", OrganizationID: triggerOrg},
			{ID: triggerDestination, Name: "destination", OrganizationID: triggerOrg},
			{ID: triggerThird, Name: "third", OrganizationID: triggerOrg},
			{ID: triggerForeign, Name: "foreign", OrganizationID: triggerOtherOrg},
		},
		templates: []store.Template{
			{ID: triggerTemplate, Name: strPtr("template"), OrganizationID: triggerOrg},
			{ID: triggerForeign, Name: strPtr("foreign-template"), OrganizationID: triggerOtherOrg},
		},
	}
}

func runTriggerAPI(stub *runTriggerStub) http.Handler {
	return New(Deps{
		Store:  stub,
		Lists:  stub,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
}

func runTriggerRBACAPI(stub *runTriggerStub) http.Handler {
	return New(Deps{
		Store:      stub,
		Lists:      stub,
		Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
		Verifiers:  []auth.Verifier{groupVerifier{}},
		AdminGroup: "TERRAKUBE_ADMIN",
	})
}

func runTriggerBody(source, destination string, extra map[string]any) *strings.Reader {
	relationships := map[string]any{
		"sourceWorkspace":      map[string]any{"data": map[string]any{"type": "workspace", "id": source}},
		"destinationWorkspace": map[string]any{"data": map[string]any{"type": "workspace", "id": destination}},
	}
	attributes := map[string]any{}
	for key, value := range extra {
		switch key {
		case "template":
			relationships["template"] = map[string]any{"data": map[string]any{"type": "template", "id": value}}
		default:
			attributes[key] = value
		}
	}
	payload, _ := json.Marshal(map[string]any{"data": map[string]any{
		"type": "runTrigger", "attributes": attributes, "relationships": relationships,
	}})
	return strings.NewReader(string(payload))
}

func TestCreateRunTriggerDerivesOrganizationAndRelationships(t *testing.T) {
	stub := newRunTriggerStub()
	handler := runTriggerAPI(stub)

	req := httptest.NewRequest(http.MethodPost, "/runTrigger",
		runTriggerBody(triggerSource, triggerDestination, map[string]any{"template": triggerTemplate}))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if stub.createdRunTrigger == nil {
		t.Fatal("run trigger was not created")
	}
	created := stub.createdRunTrigger
	if created.OrganizationID != triggerOrg {
		t.Errorf("organization = %q, want the destination's %q", created.OrganizationID, triggerOrg)
	}
	if !created.Enabled {
		t.Error("enabled did not default to true")
	}
	if created.TemplateID == nil || *created.TemplateID != triggerTemplate {
		t.Errorf("template = %v, want %s", created.TemplateID, triggerTemplate)
	}

	var doc struct {
		Data struct {
			Type          string         `json:"type"`
			Attributes    map[string]any `json:"attributes"`
			Relationships map[string]struct {
				Data *struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"relationships"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Type != "runTrigger" || doc.Data.Attributes["enabled"] != true {
		t.Errorf("resource = %+v, want a runTrigger with enabled=true", doc.Data)
	}
	if got := doc.Data.Relationships["sourceWorkspace"].Data; got == nil || got.Type != "workspace" || got.ID != triggerSource {
		t.Errorf("sourceWorkspace = %+v", got)
	}
	if got := doc.Data.Relationships["destinationWorkspace"].Data; got == nil || got.ID != triggerDestination {
		t.Errorf("destinationWorkspace = %+v", got)
	}
	if got := doc.Data.Relationships["template"].Data; got == nil || got.ID != triggerTemplate {
		t.Errorf("template = %+v", got)
	}
}

func TestCreateRunTriggerRejectsStructuralViolations(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		destination string
		extra       map[string]any
	}{
		{"self link", triggerSource, triggerSource, nil},
		{"cross organization", triggerSource, triggerForeign, nil},
		{"unknown workspace", triggerSource, "70000000-0000-4000-8000-0000000000ff", nil},
		{"foreign template", triggerSource, triggerDestination, map[string]any{"template": triggerForeign}},
		{"unknown template", triggerSource, triggerDestination, map[string]any{"template": "70000000-0000-4000-8000-0000000000fe"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stub := newRunTriggerStub()
			handler := runTriggerAPI(stub)

			req := httptest.NewRequest(http.MethodPost, "/runTrigger",
				runTriggerBody(test.source, test.destination, test.extra))
			req.Header.Set("Content-Type", "application/vnd.api+json")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body.String())
			}
			if stub.createdRunTrigger != nil {
				t.Fatal("run trigger was created despite the rejected shape")
			}
		})
	}
}

func TestCreateRunTriggerRejectsCycle(t *testing.T) {
	stub := newRunTriggerStub()
	stub.runTriggers = []store.WorkspaceRunTrigger{
		{ID: triggerEdgeA, SourceWorkspaceID: triggerSource, DestinationWorkspaceID: triggerDestination,
			OrganizationID: triggerOrg, Enabled: true},
		{ID: triggerEdgeB, SourceWorkspaceID: triggerDestination, DestinationWorkspaceID: triggerThird,
			OrganizationID: triggerOrg, Enabled: true},
	}
	handler := runTriggerAPI(stub)

	req := httptest.NewRequest(http.MethodPost, "/runTrigger",
		runTriggerBody(triggerThird, triggerSource, nil))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if stub.createdRunTrigger != nil {
		t.Fatal("a cyclic run trigger was created")
	}
	if !strings.Contains(rec.Body.String(), "circular dependency") {
		t.Errorf("detail does not name the cycle: %s", rec.Body.String())
	}
}

func TestCreateRunTriggerEnforcesFanOutLimit(t *testing.T) {
	stub := newRunTriggerStub()
	for i := 0; i < maxOutboundRunTriggers; i++ {
		stub.runTriggers = append(stub.runTriggers, store.WorkspaceRunTrigger{
			ID:                     string(rune('a' + i)),
			SourceWorkspaceID:      triggerSource,
			DestinationWorkspaceID: triggerThird,
			OrganizationID:         triggerOrg,
			Enabled:                true,
		})
	}
	handler := runTriggerAPI(stub)

	req := httptest.NewRequest(http.MethodPost, "/runTrigger",
		runTriggerBody(triggerSource, triggerDestination, nil))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if stub.createdRunTrigger != nil {
		t.Fatal("a run trigger past the fan-out limit was created")
	}
}

func TestListRunTriggersFiltersBothDirectionsAndIncludesNames(t *testing.T) {
	stub := newRunTriggerStub()
	stub.runTriggers = []store.WorkspaceRunTrigger{
		{ID: triggerEdgeA, SourceWorkspaceID: triggerSource, DestinationWorkspaceID: triggerDestination,
			TemplateID: strPtr(triggerTemplate), OrganizationID: triggerOrg, Enabled: true},
		{ID: triggerEdgeB, SourceWorkspaceID: triggerDestination, DestinationWorkspaceID: triggerThird,
			OrganizationID: triggerOrg, Enabled: false},
	}
	handler := runTriggerAPI(stub)

	req := httptest.NewRequest(http.MethodGet,
		"/runTrigger?include=sourceWorkspace,destinationWorkspace,template&filter[runTrigger]="+
			"sourceWorkspace.id=="+triggerDestination+",destinationWorkspace.id=="+triggerDestination, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if stub.runTriggerFilter == nil {
		t.Error("filter[runTrigger] was not parsed and passed to the list query")
	}
	body := rec.Body.String()
	for _, want := range []string{`"type":"runTrigger"`, `"type":"workspace"`, `"type":"template"`} {
		if !strings.Contains(body, want) {
			t.Errorf("body lacks %s: %s", want, body)
		}
	}
	if strings.Count(body, `"type":"workspace"`) < 2 {
		t.Errorf("include did not resolve both ends: %s", body)
	}
}

func TestUpdateRunTriggerTogglesEnabledWithoutTouchingTopology(t *testing.T) {
	stub := newRunTriggerStub()
	stub.runTriggers = []store.WorkspaceRunTrigger{{
		ID: triggerEdgeA, SourceWorkspaceID: triggerSource, DestinationWorkspaceID: triggerDestination,
		OrganizationID: triggerOrg, Enabled: true,
	}}
	handler := runTriggerAPI(stub)

	body := `{"data":{"type":"runTrigger","id":"` + triggerEdgeA + `","attributes":{"enabled":false}}}`
	req := httptest.NewRequest(http.MethodPatch, "/runTrigger/"+triggerEdgeA, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("patch returned a body: %s", rec.Body.String())
	}
	if stub.updatedRunTrigger == nil {
		t.Fatal("run trigger was not updated")
	}
	if stub.updatedRunTrigger.Enabled {
		t.Error("enabled was not turned off")
	}
	if stub.updatedRunTrigger.SourceWorkspaceID != triggerSource ||
		stub.updatedRunTrigger.DestinationWorkspaceID != triggerDestination {
		t.Error("topology changed during a benign update")
	}
}

func TestNonOwnerCannotRepointRunTrigger(t *testing.T) {
	stub := newRunTriggerStub()
	stub.teams = []store.Team{{
		ID: "team", Name: "TERRAKUBE_DEVELOPERS", OrganizationID: triggerOrg,
		ManageWorkspace: true,
	}}
	stub.runTriggers = []store.WorkspaceRunTrigger{{
		ID: triggerEdgeA, SourceWorkspaceID: triggerSource, DestinationWorkspaceID: triggerDestination,
		OrganizationID: triggerOrg, Enabled: true,
	}}
	handler := runTriggerRBACAPI(stub)

	body := `{"data":{"type":"runTrigger","id":"` + triggerEdgeA + `","relationships":{"sourceWorkspace":{"data":{"type":"workspace","id":"` +
		triggerThird + `"}}}}}`
	req := httptest.NewRequest(http.MethodPatch, "/runTrigger/"+triggerEdgeA, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", "Bearer groups:TERRAKUBE_DEVELOPERS")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body.String())
	}
	if stub.updatedRunTrigger != nil {
		t.Fatal("a non-owner repointed the run trigger")
	}
}

func TestDeleteRunTriggerAndContentTypeQuirk(t *testing.T) {
	stub := newRunTriggerStub()
	stub.runTriggers = []store.WorkspaceRunTrigger{{
		ID: triggerEdgeA, SourceWorkspaceID: triggerSource, DestinationWorkspaceID: triggerDestination,
		OrganizationID: triggerOrg, Enabled: true,
	}}
	handler := runTriggerAPI(stub)

	req := httptest.NewRequest(http.MethodDelete, "/runTrigger/"+triggerEdgeA, nil)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("delete with a content type = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if len(stub.deletedRunTriggers) != 0 {
		t.Fatal("the edge was deleted despite the rejected request")
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/runTrigger/"+triggerEdgeA, nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if len(stub.deletedRunTriggers) != 1 || stub.deletedRunTriggers[0] != triggerEdgeA {
		t.Errorf("deleted = %v, want [%s]", stub.deletedRunTriggers, triggerEdgeA)
	}
}
