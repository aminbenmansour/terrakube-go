package tfe

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type runAttributes struct {
	IsDestroy *bool `json:"is-destroy"`
	AutoApply *bool `json:"auto-apply"`
	Refresh   *bool `json:"refresh"`
}

func (a *API) createRun(w http.ResponseWriter, r *http.Request) {
	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}

	workspaceID, present, err := jsonapi.RelationshipID(data, "workspace")
	if err != nil || !present || workspaceID == "" {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid Relationship", "runs require a workspace relationship.")
		return
	}
	configurationID, present, err := jsonapi.RelationshipID(data, "configuration-version")
	if err != nil || !present || configurationID == "" {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid Relationship", "runs require a configuration-version relationship.")
		return
	}

	ws, err := a.deps.Store.GetWorkspace(r.Context(), workspaceID)
	if err != nil {
		if !a.writeStoreError(w, "get run workspace", err) {
			a.internalError(w, "get run workspace", err)
		}
		return
	}
	content, err := a.deps.Store.GetContent(r.Context(), configurationID)
	if err != nil {
		if !a.writeStoreError(w, "get run configuration", err) {
			a.internalError(w, "get run configuration", err)
		}
		return
	}
	if content.WorkspaceID == nil || *content.WorkspaceID != ws.ID {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid Configuration", "The configuration version belongs to another workspace.")
		return
	}
	if content.Status == nil || *content.Status != "uploaded" {
		jsonapi.WriteError(w, http.StatusUnprocessableEntity, "Invalid Configuration", "The configuration version has not been uploaded.")
		return
	}

	attrs, err := jsonapi.DecodeAttributes[runAttributes](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	isDestroy := attrs.IsDestroy != nil && *attrs.IsDestroy
	autoApply := attrs.AutoApply != nil && *attrs.AutoApply
	refresh := true
	if attrs.Refresh != nil {
		refresh = *attrs.Refresh
	}
	command := "plan"
	if isDestroy {
		command = "destroy"
	}

	job, err := a.deps.Store.CreateJob(r.Context(), store.CreateJobParams{
		OrganizationID: ws.OrganizationID,
		WorkspaceID:    ws.ID,
		Command:        stringPointer(command),
		Status:         "pending",
		Via:            "CLI",
		AutoApply:      autoApply,
		Refresh:        refresh,
		PlanChanges:    true,
		CreatedBy:      stringPointer("tfe-api"),
	})
	if err != nil {
		a.internalError(w, "create run", err)
		return
	}
	if job.Comments == nil {
		job.Comments = stringPointer("terraform-cli")
	}
	if a.deps.Queue == nil {
		_, _ = a.deps.Store.UpdateJobStatus(r.Context(), store.UpdateJobStatusParams{ID: job.ID, Status: "failed"})
		a.internalError(w, "enqueue run", errors.New("run queue is not configured"))
		return
	}
	if _, err := a.deps.Queue.EnqueueJob(r.Context(), job.ID); err != nil {
		_, _ = a.deps.Store.UpdateJobStatus(r.Context(), store.UpdateJobStatusParams{ID: job.ID, Status: "failed"})
		a.internalError(w, "enqueue run", err)
		return
	}

	jsonapi.WriteResource(w, http.StatusCreated, runResource(job, ws))
}

func (a *API) getRun(w http.ResponseWriter, r *http.Request) {
	job, ws, ok := a.jobAndWorkspace(w, r, "runID")
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, runResource(job, ws))
}

func (a *API) listWorkspaceRuns(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}
	jobs, err := a.deps.Store.ListActiveJobsByWorkspace(r.Context(), ws.ID)
	if err != nil {
		a.internalError(w, "list workspace runs", err)
		return
	}
	resources := make([]jsonapi.Resource, 0, len(jobs))
	for _, job := range jobs {
		resources = append(resources, runResource(job, ws))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), false)
}

func (a *API) getRunQueue(w http.ResponseWriter, r *http.Request) {
	org, ok := a.organizationByName(w, r)
	if !ok {
		return
	}
	jobs, err := a.deps.Store.ListActiveJobsByOrganization(r.Context(), org.Name)
	if err != nil {
		a.internalError(w, "list run queue", err)
		return
	}
	resources := make([]jsonapi.Resource, 0, len(jobs))
	for _, job := range jobs {
		ws, err := a.deps.Store.GetWorkspace(r.Context(), job.WorkspaceID)
		if err != nil {
			if !a.writeStoreError(w, "get queued run workspace", err) {
				a.internalError(w, "get queued run workspace", err)
			}
			return
		}
		resources = append(resources, runResource(job, ws))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), false)
}

func (a *API) getRunEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "{\n  \"data\": []\n}")
}

func (a *API) applyRun(w http.ResponseWriter, r *http.Request) {
	job, ws, ok := a.jobAndWorkspace(w, r, "runID")
	if !ok {
		return
	}
	if job.Status == "cancelled" || job.Status == "failed" {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "The run cannot be applied in its current state.")
		return
	}
	updated, err := a.deps.Store.UpdateJobStatus(r.Context(), store.UpdateJobStatusParams{
		ID:     job.ID,
		Status: "pending",
	})
	if err != nil {
		a.internalError(w, "apply run", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, runResource(updated, ws))
}

func (a *API) discardRun(w http.ResponseWriter, r *http.Request) {
	job, ws, ok := a.jobAndWorkspace(w, r, "runID")
	if !ok {
		return
	}
	if !activeRunStatus(job.Status) {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "The run cannot be discarded in its current state.")
		return
	}
	updated, err := a.deps.Store.UpdateJobStatus(r.Context(), store.UpdateJobStatusParams{
		ID:     job.ID,
		Status: "cancelled",
	})
	if err != nil {
		a.internalError(w, "discard run", err)
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, runResource(updated, ws))
}

func (a *API) getPlan(w http.ResponseWriter, r *http.Request) {
	job, _, ok := a.jobAndWorkspace(w, r, "planID")
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, planResource(a, job))
}

func (a *API) getApply(w http.ResponseWriter, r *http.Request) {
	job, _, ok := a.jobAndWorkspace(w, r, "applyID")
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, applyResource(a, job))
}

func (a *API) getPlanLogs(w http.ResponseWriter, r *http.Request) {
	a.writeRunLogs(w, r, "planID")
}

func (a *API) getApplyLogs(w http.ResponseWriter, r *http.Request) {
	a.writeRunLogs(w, r, "applyID")
}

func (a *API) getRunPolicyChecks(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := a.jobAndWorkspace(w, r, "runID"); !ok {
		return
	}
	jsonapi.WriteCollection(w, http.StatusOK, []jsonapi.Resource{}, 0, false)
}

func (a *API) getPolicyCheck(w http.ResponseWriter, r *http.Request) {
	jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The policy check does not exist.")
}

func (a *API) writeRunLogs(w http.ResponseWriter, r *http.Request, param string) {
	job, _, ok := a.jobAndWorkspace(w, r, param)
	if !ok {
		return
	}
	output := ""
	if job.Output != nil {
		output = *job.Output
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if offset < 0 || limit <= 0 || offset >= len(output) {
		output = ""
	} else {
		end := offset + limit
		if end > len(output) {
			end = len(output)
		}
		output = output[offset:end]
	}
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, output)
}

func (a *API) jobAndWorkspace(w http.ResponseWriter, r *http.Request, param string) (store.Job, store.Workspace, bool) {
	runID, err := parseRunID(chi.URLParam(r, param))
	if err != nil {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The run does not exist.")
		return store.Job{}, store.Workspace{}, false
	}
	job, err := a.deps.Store.GetJob(r.Context(), runID)
	if err != nil {
		if !a.writeStoreError(w, "get run", err) {
			a.internalError(w, "get run", err)
		}
		return store.Job{}, store.Workspace{}, false
	}
	ws, err := a.deps.Store.GetWorkspace(r.Context(), job.WorkspaceID)
	if err != nil {
		if !a.writeStoreError(w, "get run workspace", err) {
			a.internalError(w, "get run workspace", err)
		}
		return store.Job{}, store.Workspace{}, false
	}
	return job, ws, true
}

func runResource(job store.Job, ws store.Workspace) jsonapi.Resource {
	id := strconv.FormatInt(job.ID, 10)
	return jsonapi.Resource{
		Type: "runs",
		ID:   "run-" + id,
		Attributes: map[string]any{
			"status":                runStatus(job.Status),
			"has-changes":           job.PlanChanges,
			"resource-additions":    1,
			"resource-changes":      1,
			"resource-destructions": 0,
			"actions": map[string]any{
				"is-confirmable": true,
				"is-discardable": true,
			},
			"permissions": map[string]any{
				"can-apply": true,
			},
		},
		Relationships: map[string]jsonapi.Relationship{
			"plan":      jsonapi.Link("plans", id),
			"apply":     jsonapi.Link("applies", id),
			"workspace": jsonapi.Link("workspaces", ws.ID),
		},
	}
}

func planResource(a *API, job store.Job) jsonapi.Resource {
	id := strconv.FormatInt(job.ID, 10)
	return jsonapi.Resource{
		Type: "plans",
		ID:   id,
		Attributes: map[string]any{
			"status": planStatus(job.Status),
			"actions": map[string]any{
				"is-confirmable": true,
				"is-discardable": false,
			},
			"log-read-url": a.tfeURL(fmt.Sprintf("/plans/logs/%s/0", id)),
		},
	}
}

func applyResource(a *API, job store.Job) jsonapi.Resource {
	id := strconv.FormatInt(job.ID, 10)
	return jsonapi.Resource{
		Type: "applies",
		ID:   id,
		Attributes: map[string]any{
			"status":       applyStatus(job.Status),
			"log-read-url": a.tfeURL(fmt.Sprintf("/applies/logs/%s/0", id)),
		},
	}
}

func runStatus(status string) string {
	switch status {
	case "completed":
		return "finished"
	case "pending":
		return "pending"
	case "running", "queue", "waitingApproval":
		return "running"
	case "failed":
		return "errored"
	default:
		return "unknown"
	}
}

func planStatus(status string) string {
	switch status {
	case "completed":
		return "finished"
	case "running", "queue", "waitingApproval":
		return "running"
	case "failed":
		return "errored"
	default:
		return "pending"
	}
}

func applyStatus(status string) string {
	switch status {
	case "completed":
		return "finished"
	case "running", "queue", "waitingApproval":
		return "running"
	case "failed":
		return "errored"
	default:
		return "pending"
	}
}

func activeRunStatus(status string) bool {
	switch status {
	case "pending", "running", "queue", "waitingApproval":
		return true
	default:
		return false
	}
}

func parseRunID(raw string) (int64, error) {
	raw = strings.TrimPrefix(raw, "run-")
	if raw == "" {
		return 0, errors.New("run id is empty")
	}
	return strconv.ParseInt(raw, 10, 64)
}

func (a *API) tfeURL(path string) string {
	return strings.TrimRight(a.deps.PublicURL, "/") + "/remote/tfe/v2" + path
}
