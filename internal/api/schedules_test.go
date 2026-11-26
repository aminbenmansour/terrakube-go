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

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type scheduleStub struct {
	store.Querier
	workspace store.Workspace
	templates []store.Template
	schedules []store.Schedule
	created   *store.CreateScheduleParams
	updated   *store.UpdateScheduleParams
	deleted   []string
}

func (s *scheduleStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	if id == s.workspace.ID {
		return s.workspace, nil
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *scheduleStub) GetTemplate(_ context.Context, id string) (store.Template, error) {
	for _, template := range s.templates {
		if template.ID == id {
			return template, nil
		}
	}
	return store.Template{}, pgx.ErrNoRows
}

func (s *scheduleStub) ListSchedulesByWorkspace(_ context.Context, workspaceID string) ([]store.Schedule, error) {
	items := []store.Schedule{}
	for _, schedule := range s.schedules {
		if schedule.WorkspaceID == workspaceID {
			items = append(items, schedule)
		}
	}
	return items, nil
}

func (s *scheduleStub) GetSchedule(_ context.Context, id string) (store.Schedule, error) {
	for _, schedule := range s.schedules {
		if schedule.ID == id {
			return schedule, nil
		}
	}
	return store.Schedule{}, pgx.ErrNoRows
}

func (s *scheduleStub) CreateSchedule(_ context.Context, arg store.CreateScheduleParams) (store.Schedule, error) {
	s.created = &arg
	return store.Schedule{
		ID:                arg.ID,
		Cron:              arg.Cron,
		Description:       arg.Description,
		TemplateReference: arg.TemplateReference,
		Enabled:           arg.Enabled,
		WorkspaceID:       arg.WorkspaceID,
	}, nil
}

func (s *scheduleStub) UpdateSchedule(_ context.Context, arg store.UpdateScheduleParams) (store.Schedule, error) {
	s.updated = &arg
	for _, schedule := range s.schedules {
		if schedule.ID == arg.ID {
			schedule.Cron = arg.Cron
			schedule.Description = arg.Description
			schedule.TemplateReference = arg.TemplateReference
			schedule.Enabled = arg.Enabled
			return schedule, nil
		}
	}
	return store.Schedule{}, pgx.ErrNoRows
}

func (s *scheduleStub) DeleteSchedule(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func scheduleStr(v string) *string { return &v }

const (
	scheduleOrg        = "70000000-0000-4000-8000-000000000001"
	scheduleOtherOrg   = "70000000-0000-4000-8000-000000000009"
	scheduleWorkspace  = "70000000-0000-4000-8000-000000000002"
	scheduleTemplateID = "70000000-0000-4000-8000-000000000003"
)

func newScheduleAPI(q *scheduleStub) http.Handler {
	return New(Deps{Store: q, Logger: slog.New(slog.NewTextHandler(io.Discard, nil))})
}

func newScheduleStub() *scheduleStub {
	return &scheduleStub{
		workspace: store.Workspace{ID: scheduleWorkspace, Name: "ws", OrganizationID: scheduleOrg},
		templates: []store.Template{
			{ID: scheduleTemplateID, OrganizationID: scheduleOrg},
			{ID: "70000000-0000-4000-8000-000000000004", OrganizationID: scheduleOtherOrg},
		},
	}
}

func scheduleBody(attributes map[string]any) *strings.Reader {
	payload, _ := json.Marshal(map[string]any{
		"data": map[string]any{
			"type":       "schedule",
			"attributes": attributes,
		},
	})
	return strings.NewReader(string(payload))
}

func TestCreateScheduleCopiesAttributesAndDefaultsEnabled(t *testing.T) {
	stub := newScheduleStub()
	handler := newScheduleAPI(stub)

	req := httptest.NewRequest(http.MethodPost,
		"/organization/"+scheduleOrg+"/workspace/"+scheduleWorkspace+"/schedule",
		scheduleBody(map[string]any{"cron": "0 0/5 * * * ?", "templateReference": scheduleTemplateID}))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if stub.created == nil {
		t.Fatal("schedule was not created")
	}
	if got := stub.created.TemplateReference; got == nil || *got != scheduleTemplateID {
		t.Errorf("templateReference = %v, want %s", got, scheduleTemplateID)
	}
	if !stub.created.Enabled {
		t.Error("enabled did not default to true")
	}

	var doc struct {
		Data struct {
			Type          string         `json:"type"`
			Attributes    map[string]any `json:"attributes"`
			Relationships map[string]struct {
				Data struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"relationships"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Type != "schedule" {
		t.Errorf("type = %q, want schedule", doc.Data.Type)
	}
	if doc.Data.Attributes["cron"] != "0 0/5 * * * ?" {
		t.Errorf("cron = %v", doc.Data.Attributes["cron"])
	}
	if got := doc.Data.Relationships["workspace"].Data.ID; got != scheduleWorkspace {
		t.Errorf("workspace relationship = %q", got)
	}
}

func TestCreateScheduleRejectsForeignTemplate(t *testing.T) {
	stub := newScheduleStub()
	handler := newScheduleAPI(stub)

	for name, template := range map[string]string{
		"another organization": "70000000-0000-4000-8000-000000000004",
		"unknown template":     "70000000-0000-4000-8000-000000000005",
		"not a uuid":           "nope",
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost,
				"/organization/"+scheduleOrg+"/workspace/"+scheduleWorkspace+"/schedule",
				scheduleBody(map[string]any{"cron": "0 0 * * * ?", "templateReference": template}))
			req.Header.Set("Content-Type", "application/vnd.api+json")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Fatalf("status = %d, want 403; body = %s", rec.Code, rec.Body.String())
			}
			if stub.created != nil {
				t.Fatal("schedule was created despite the rejected template")
			}
		})
	}
}

func TestListAndUpdateSchedule(t *testing.T) {
	existing := store.Schedule{
		ID:          "70000000-0000-4000-8000-00000000000a",
		Cron:        scheduleStr("0 0 * * * ?"),
		Enabled:     true,
		WorkspaceID: scheduleWorkspace,
	}
	stub := newScheduleStub()
	stub.schedules = []store.Schedule{existing}
	handler := newScheduleAPI(stub)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/organization/"+scheduleOrg+"/workspace/"+scheduleWorkspace+"/schedule", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), existing.ID) {
		t.Errorf("list body lacks the schedule: %s", rec.Body.String())
	}

	req := httptest.NewRequest(http.MethodPatch,
		"/organization/"+scheduleOrg+"/workspace/"+scheduleWorkspace+"/schedule/"+existing.ID,
		scheduleBody(map[string]any{"cron": "0 30 * * * ?", "enabled": false}))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("patch status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if stub.updated == nil {
		t.Fatal("schedule was not updated")
	}
	if got := stub.updated.Cron; got == nil || *got != "0 30 * * * ?" {
		t.Errorf("cron = %v, want the new expression", got)
	}
	if stub.updated.Enabled {
		t.Error("enabled was not turned off")
	}
	if got := stub.updated.TemplateReference; got != nil {
		t.Errorf("untouched templateReference = %v", got)
	}
}

func TestDeleteScheduleAndCrossWorkspaceLookup(t *testing.T) {
	schedule := store.Schedule{
		ID:          "70000000-0000-4000-8000-00000000000b",
		WorkspaceID: scheduleWorkspace,
		Enabled:     true,
	}
	stub := newScheduleStub()
	stub.schedules = []store.Schedule{schedule}
	handler := newScheduleAPI(stub)

	// The schedule exists, but not under this workspace id.
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/organization/"+scheduleOrg+"/workspace/70000000-0000-4000-8000-0000000000ff/schedule/"+schedule.ID, nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("cross-workspace status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete,
		"/organization/"+scheduleOrg+"/workspace/"+scheduleWorkspace+"/schedule/"+schedule.ID, nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if len(stub.deleted) != 1 || stub.deleted[0] != schedule.ID {
		t.Errorf("deleted = %v, want [%s]", stub.deleted, schedule.ID)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("delete returned a body: %s", rec.Body.String())
	}
}

func TestWorkspaceDetailIncludesSchedules(t *testing.T) {
	stub := newScheduleStub()
	stub.schedules = []store.Schedule{{
		ID:          "70000000-0000-4000-8000-00000000000c",
		WorkspaceID: scheduleWorkspace,
		Cron:        scheduleStr("0 0 * * * ?"),
		Enabled:     true,
	}}
	handler := newScheduleAPI(stub)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
		"/organization/"+scheduleOrg+"/workspace/"+scheduleWorkspace+"?include=schedule", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"schedule"`) {
		t.Errorf("workspace relationships lack schedule: %s", body)
	}
	if !strings.Contains(body, `"type":"schedule"`) {
		t.Errorf("included lacks the schedule resource: %s", body)
	}
}
