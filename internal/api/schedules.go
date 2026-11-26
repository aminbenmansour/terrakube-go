package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type scheduleAttrs struct {
	Cron              *string `json:"cron"`
	Description       *string `json:"description"`
	TemplateReference *string `json:"templateReference"`
	Enabled           *bool   `json:"enabled"`
}

func scheduleResource(schedule store.Schedule) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "schedule",
		ID:   schedule.ID,
		Attributes: map[string]any{
			"cron":              schedule.Cron,
			"description":       schedule.Description,
			"enabled":           schedule.Enabled,
			"templateReference": schedule.TemplateReference,
			"createdBy":         schedule.CreatedBy,
			"createdDate":       timestamp(schedule.CreatedDate),
			"updatedBy":         schedule.UpdatedBy,
			"updatedDate":       timestamp(schedule.UpdatedDate),
		},
		Relationships: map[string]jsonapi.Relationship{
			"workspace": jsonapi.Link("workspace", schedule.WorkspaceID),
		},
	}
}

func (a *API) listSchedules(w http.ResponseWriter, r *http.Request) {
	wsID := param(r, "workspaceID")
	if _, err := a.deps.Store.GetWorkspace(r.Context(), wsID); err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}

	schedules, err := a.deps.Store.ListSchedulesByWorkspace(r.Context(), wsID)
	if err != nil {
		if !a.writeStoreError(w, "list schedules", err) {
			a.internalError(w, "list schedules", err)
		}
		return
	}

	resources := make([]jsonapi.Resource, 0, len(schedules))
	for _, schedule := range schedules {
		resources = append(resources, scheduleResource(schedule))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{}, jsonapi.ParseFieldsets(r))
}

func (a *API) getSchedule(w http.ResponseWriter, r *http.Request) {
	schedule, ok := a.scheduleInWorkspace(w, r)
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, scheduleResource(schedule), jsonapi.ParseFieldsets(r))
}

func (a *API) createSchedule(w http.ResponseWriter, r *http.Request) {
	ws, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "schedule") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[scheduleAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	ok, err := a.scheduleTemplateAllowed(r, ws.OrganizationID, attrs.TemplateReference)
	if err != nil {
		a.internalError(w, "check schedule template", err)
		return
	}
	if !ok {
		writePermissionDenied(w, r.Method)
		return
	}

	schedule, err := a.deps.Store.CreateSchedule(r.Context(), store.CreateScheduleParams{
		ID:                uuid.NewString(),
		Cron:              attrs.Cron,
		Description:       attrs.Description,
		TemplateReference: attrs.TemplateReference,
		Enabled:           boolValue(attrs.Enabled, true),
		WorkspaceID:       ws.ID,
	})
	if err != nil {
		if !a.writeStoreError(w, "create schedule", err) {
			a.internalError(w, "create schedule", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusCreated, scheduleResource(schedule), jsonapi.ParseFieldsets(r))
}

func (a *API) updateSchedule(w http.ResponseWriter, r *http.Request) {
	schedule, ok := a.scheduleInWorkspace(w, r)
	if !ok {
		return
	}
	ws, err := a.deps.Store.GetWorkspace(r.Context(), schedule.WorkspaceID)
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "schedule") {
		return
	}
	attrs, err := jsonapi.DecodeAttributes[scheduleAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}

	templateReference := pickString(data, "templateReference", attrs.TemplateReference, schedule.TemplateReference)
	allowed, err := a.scheduleTemplateAllowed(r, ws.OrganizationID, templateReference)
	if err != nil {
		a.internalError(w, "check schedule template", err)
		return
	}
	if !allowed {
		writePermissionDenied(w, r.Method)
		return
	}

	enabled := schedule.Enabled
	if jsonapi.HasAttribute(data, "enabled") && attrs.Enabled != nil {
		enabled = *attrs.Enabled
	}
	updated, err := a.deps.Store.UpdateSchedule(r.Context(), store.UpdateScheduleParams{
		ID:                schedule.ID,
		Cron:              pickString(data, "cron", attrs.Cron, schedule.Cron),
		Description:       pickString(data, "description", attrs.Description, schedule.Description),
		TemplateReference: templateReference,
		Enabled:           enabled,
	})
	if err != nil {
		if !a.writeStoreError(w, "update schedule", err) {
			a.internalError(w, "update schedule", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, scheduleResource(updated), jsonapi.ParseFieldsets(r))
}

func (a *API) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	schedule, ok := a.scheduleInWorkspace(w, r)
	if !ok {
		return
	}
	if err := a.deps.Store.DeleteSchedule(r.Context(), schedule.ID); err != nil {
		if !a.writeStoreError(w, "delete schedule", err) {
			a.internalError(w, "delete schedule", err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// scheduleInWorkspace loads the schedule named in the route and answers 404
// when it belongs to another workspace, as the relationship route would.
func (a *API) scheduleInWorkspace(w http.ResponseWriter, r *http.Request) (store.Schedule, bool) {
	schedule, err := a.deps.Store.GetSchedule(r.Context(), param(r, "scheduleID"))
	if err != nil {
		if !a.writeStoreError(w, "get schedule", err) {
			a.internalError(w, "get schedule", err)
		}
		return store.Schedule{}, false
	}
	if schedule.WorkspaceID != param(r, "workspaceID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Schedule{}, false
	}
	return schedule, true
}

// scheduleTemplateAllowed mirrors upstream's TeamManageSchedule check: an empty
// reference is fine, otherwise the template must belong to the workspace's
// organization. A malformed or unknown id fails the same check.
func (a *API) scheduleTemplateAllowed(r *http.Request, organizationID string, templateReference *string) (bool, error) {
	if templateReference == nil || *templateReference == "" {
		return true, nil
	}
	template, err := a.deps.Store.GetTemplate(r.Context(), *templateReference)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return template.OrganizationID == organizationID, nil
}
