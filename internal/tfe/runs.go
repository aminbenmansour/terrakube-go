package tfe

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Step numbers mirror the upstream Terraform-Plan/Apply-Cli and
// Terraform-Plan/Destroy-Cli templates.
const (
	planStepNumber     int32 = 100
	approvalStepNumber int32 = 150
	applyStepNumber    int32 = 200
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
	jobID := job.ID

	// The CLI uploads its configuration through a configuration version, so the
	// job carries the tarball URL in override_source and the reserved
	// remote-content branch, exactly like upstream.
	source := a.tfeURL(fmt.Sprintf("/configuration-versions/%s/terraformContent.tar.gz", configurationID))
	job, err = a.deps.Store.SetJobOverrideSource(r.Context(), store.SetJobOverrideSourceParams{
		ID:             jobID,
		OverrideSource: stringPointer(source),
		OverrideBranch: stringPointer("remote-content"),
	})
	if err != nil {
		_, _ = a.deps.Store.UpdateJobStatus(r.Context(), store.UpdateJobStatusParams{ID: jobID, Status: "failed"})
		a.internalError(w, "set run configuration source", err)
		return
	}
	// Speculative configurations use upstream's "Plan" template, which has a
	// single step. Destroy takes precedence over speculative, as upstream's
	// getTemplateName does.
	speculative := content.Speculative != nil && *content.Speculative && !isDestroy
	if err := a.createRunSteps(r.Context(), jobID, isDestroy, speculative); err != nil {
		_, _ = a.deps.Store.UpdateJobStatus(r.Context(), store.UpdateJobStatusParams{ID: jobID, Status: "failed"})
		a.internalError(w, "create run steps", err)
		return
	}
	_ = a.deps.Store.SetWorkspaceLastJob(r.Context(), store.SetWorkspaceLastJobParams{
		ID:            ws.ID,
		LastJobStatus: "pending",
	})

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

	a.writeRun(r.Context(), w, http.StatusCreated, job, ws, r.URL.Query().Get("include"))
}

func (a *API) getRun(w http.ResponseWriter, r *http.Request) {
	job, ws, ok := a.jobAndWorkspace(w, r, "runID")
	if !ok {
		return
	}
	a.writeRun(r.Context(), w, http.StatusOK, job, ws, r.URL.Query().Get("include"))
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
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{})
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
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{})
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

	// Upstream only confirms a run whose plan finished. The approval step is
	// completed and the job is re-queued so the executor can run the saved plan.
	if job.Status != "completed" && job.Status != "waitingApproval" {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "The run cannot be applied in its current state.")
		return
	}

	approval, err := a.deps.Store.GetStepByJobAndNumber(r.Context(), store.GetStepByJobAndNumberParams{
		JobID:      job.ID,
		StepNumber: approvalStepNumber,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		approval, err = a.deps.Store.CreateStep(r.Context(), store.CreateStepParams{
			ID:         uuid.NewString(),
			JobID:      job.ID,
			StepNumber: approvalStepNumber,
			Name:       stringPointer("Approve Plan from Terraform CLI"),
			Status:     stringPointer("pending"),
			LogStatus:  stringPointer("pending"),
		})
	}
	if err != nil {
		a.internalError(w, "load approval step", err)
		return
	}
	if approval.Status == nil || *approval.Status != "completed" {
		if _, err := a.deps.Store.UpdateStepStatus(r.Context(), store.UpdateStepStatusParams{
			ID:     approval.ID,
			Status: stringPointer("completed"),
		}); err != nil {
			a.internalError(w, "complete approval step", err)
			return
		}
	}

	updated, err := a.deps.Store.ApproveJob(r.Context(), store.ApproveJobParams{
		ID:         job.ID,
		ApprovedBy: stringPointer("terraform-cli"),
	})
	if err != nil {
		a.internalError(w, "apply run", err)
		return
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

	a.writeRun(r.Context(), w, http.StatusOK, updated, ws, r.URL.Query().Get("include"))
}

func (a *API) createRunSteps(ctx context.Context, jobID int64, destroy, speculative bool) error {
	planName := "Terraform Plan from Terraform CLI"
	if destroy {
		planName = "Terraform Plan Destroy from Terraform CLI"
	}
	type stepSpec struct {
		number int32
		name   string
	}
	steps := []stepSpec{{planStepNumber, planName}}
	if !speculative {
		steps = append(steps,
			stepSpec{approvalStepNumber, "Approve Plan from Terraform CLI"},
			stepSpec{applyStepNumber, "Terraform Apply from Terraform CLI"},
		)
	}
	for _, step := range steps {
		if _, err := a.deps.Store.CreateStep(ctx, store.CreateStepParams{
			ID:         uuid.NewString(),
			JobID:      jobID,
			StepNumber: step.number,
			Name:       stringPointer(step.name),
			Status:     stringPointer("pending"),
			LogStatus:  stringPointer("pending"),
		}); err != nil {
			return err
		}
	}
	return nil
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
	_ = a.deps.Store.SetWorkspaceLastJob(r.Context(), store.SetWorkspaceLastJobParams{
		ID:            ws.ID,
		LastJobStatus: "cancelled",
	})
	a.writeRun(r.Context(), w, http.StatusOK, updated, ws, r.URL.Query().Get("include"))
}

func (a *API) getPlan(w http.ResponseWriter, r *http.Request) {
	job, _, ok := a.jobAndWorkspace(w, r, "planID")
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, a.planResource(r.Context(), job))
}

func (a *API) getApply(w http.ResponseWriter, r *http.Request) {
	job, _, ok := a.jobAndWorkspace(w, r, "applyID")
	if !ok {
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, a.applyResource(r.Context(), job))
}

func (a *API) getPlanLogs(w http.ResponseWriter, r *http.Request) {
	a.writeRunLogs(w, r, "planID")
}

func (a *API) getApplyLogs(w http.ResponseWriter, r *http.Request) {
	a.writeRunLogs(w, r, "applyID")
}

// logToken is the capability embedded in log-read-url. The Terraform CLI
// fetches that URL without an Authorization header, so the token itself must
// authorize the read.
func (a *API) logToken(id string) string {
	mac := hmac.New(sha256.New, []byte(a.deps.LogTokenSecret))
	_, _ = mac.Write([]byte("log:" + id))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (a *API) validLogToken(id, token string) bool {
	expected := a.logToken(id)
	return subtle.ConstantTimeCompare([]byte(expected), []byte(token)) == 1
}

func (a *API) logURL(kind, id string) string {
	return a.tfeURL(fmt.Sprintf("/%s/logs/%s/%s", kind, id, a.logToken(id)))
}

func (a *API) getRunPolicyChecks(w http.ResponseWriter, r *http.Request) {
	if _, _, ok := a.jobAndWorkspace(w, r, "runID"); !ok {
		return
	}
	jsonapi.WriteCollection(w, http.StatusOK, []jsonapi.Resource{}, 0, jsonapi.Page{})
}

func (a *API) getPolicyCheck(w http.ResponseWriter, r *http.Request) {
	jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The policy check does not exist.")
}

func (a *API) writeRunLogs(w http.ResponseWriter, r *http.Request, param string) {
	if !a.validLogToken(chi.URLParam(r, param), chi.URLParam(r, "logToken")) {
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "The log read token is invalid.")
		return
	}
	job, _, ok := a.jobAndWorkspace(w, r, param)
	if !ok {
		return
	}

	// Plan logs come from the plan step, apply logs from the apply step. Older
	// jobs only have job.output, so fall back to it.
	stepNumber := planStepNumber
	if param == "applyID" {
		stepNumber = applyStepNumber
	}
	output := ""
	step, err := a.deps.Store.GetStepByJobAndNumber(r.Context(), store.GetStepByJobAndNumberParams{
		JobID:      job.ID,
		StepNumber: stepNumber,
	})
	if err == nil && step.Output != nil {
		output = *step.Output
	} else if job.Output != nil {
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

func (a *API) writeRun(ctx context.Context, w http.ResponseWriter, status int, job store.Job, ws store.Workspace, include string) {
	resource := runResource(job, ws)
	included := make([]jsonapi.Resource, 0, 2)
	for _, part := range strings.Split(include, ",") {
		switch strings.TrimSpace(part) {
		case "plan":
			included = append(included, a.planResource(ctx, job))
		case "apply":
			included = append(included, a.applyResource(ctx, job))
		}
	}
	if len(included) == 0 {
		jsonapi.WriteResource(w, status, resource)
		return
	}
	jsonapi.Write(w, status, jsonapi.Document{Data: resource, Included: included})
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
				"is-confirmable": job.Status == "waitingApproval" || job.Status == "completed",
				"is-discardable": activeRunStatus(job.Status),
			},
			"permissions": map[string]any{
				"can-apply": job.Status == "waitingApproval" || job.Status == "completed",
			},
		},
		Relationships: map[string]jsonapi.Relationship{
			"plan":       jsonapi.Link("plans", id),
			"apply":      jsonapi.Link("applies", id),
			"workspace":  jsonapi.Link("workspaces", ws.ID),
			"run-events": jsonapi.LinkList("run-events"),
		},
	}
}

// planResource mirrors upstream getPlan: the plan status comes from the plan
// step, and the CLI keeps polling until the scheduler reaches the approval
// gate when the plan has changes.
func (a *API) planResource(ctx context.Context, job store.Job) jsonapi.Resource {
	id := strconv.FormatInt(job.ID, 10)
	status := "pending"
	if step, err := a.deps.Store.GetStepByJobAndNumber(ctx, store.GetStepByJobAndNumberParams{
		JobID:      job.ID,
		StepNumber: planStepNumber,
	}); err == nil {
		status = stepStatusTFE(step.Status)
	}
	if status == "finished" && !approvalGateStatus(job.Status) && job.PlanChanges {
		if approval, err := a.deps.Store.GetStepByJobAndNumber(ctx, store.GetStepByJobAndNumberParams{
			JobID:      job.ID,
			StepNumber: approvalStepNumber,
		}); err == nil && stepStatusValue(approval.Status) != "completed" {
			if failedJobStatus(job.Status) {
				status = "errored"
			} else {
				status = "running"
			}
		}
	}
	return jsonapi.Resource{
		Type: "plans",
		ID:   id,
		Attributes: map[string]any{
			"status": status,
			"actions": map[string]any{
				"is-confirmable": true,
				"is-discardable": false,
			},
			"log-read-url": a.logURL("plans", id),
		},
	}
}

// applyResource mirrors upstream getApply: the apply status comes from the
// apply step rather than the job status.
func (a *API) applyResource(ctx context.Context, job store.Job) jsonapi.Resource {
	id := strconv.FormatInt(job.ID, 10)
	status := "pending"
	if step, err := a.deps.Store.GetStepByJobAndNumber(ctx, store.GetStepByJobAndNumberParams{
		JobID:      job.ID,
		StepNumber: applyStepNumber,
	}); err == nil {
		status = stepStatusTFE(step.Status)
	}
	return jsonapi.Resource{
		Type: "applies",
		ID:   id,
		Attributes: map[string]any{
			"status":       status,
			"log-read-url": a.logURL("applies", id),
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

func stepStatusValue(status *string) string {
	if status == nil {
		return "pending"
	}
	return strings.ToLower(strings.TrimSpace(*status))
}

// stepStatusTFE maps a step status to the plan/apply status the Terraform CLI
// expects.
func stepStatusTFE(status *string) string {
	switch stepStatusValue(status) {
	case "completed":
		return "finished"
	case "queue", "running":
		return "running"
	case "failed":
		return "errored"
	default:
		return "pending"
	}
}

func approvalGateStatus(status string) bool {
	switch status {
	case "waitingApproval", "completed", "noChanges":
		return true
	default:
		return false
	}
}

func failedJobStatus(status string) bool {
	switch status {
	case "failed", "cancelled", "canceled", "rejected":
		return true
	default:
		return false
	}
}

func activeRunStatus(status string) bool {
	switch status {
	case "pending", "approved", "running", "queue", "waitingApproval":
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
