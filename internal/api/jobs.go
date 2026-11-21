package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// jobStatuses are the values upstream's Job.status enum accepts. It is
// case-sensitive: "CANCELLED" is rejected with 400 by a live 2.34.3.
var jobStatuses = map[string]bool{
	"pending": true, "queue": true, "running": true, "waitingApproval": true,
	"completed": true, "failed": true, "cancelled": true, "rejected": true,
	"approved": true, "noChanges": true, "notExecuted": true, "NeverExecuted": true,
}

type jobAttrs struct {
	Status *string `json:"status"`
}

func (a *API) listJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := a.deps.Store.ListJobsByOrganization(r.Context(), param(r, "organizationID"))
	if err != nil {
		if !a.writeStoreError(w, "list jobs", err) {
			a.internalError(w, "list jobs", err)
		}
		return
	}

	resources := make([]jsonapi.Resource, 0, len(jobs))
	for _, job := range jobs {
		resources = append(resources, jobResource(job))
	}
	jsonapi.WriteCollection(w, http.StatusOK, resources, len(resources), jsonapi.Page{})
}

func (a *API) getJob(w http.ResponseWriter, r *http.Request) {
	job, ok := a.jobInOrganization(w, r)
	if !ok {
		return
	}

	resource := jobResource(job)
	var included []jsonapi.Resource
	want := parseInclude(r)

	if want["step"] {
		steps, err := a.deps.Store.ListStepsByJob(r.Context(), job.ID)
		if err != nil {
			if !a.writeStoreError(w, "list steps", err) {
				a.internalError(w, "list steps", err)
			}
			return
		}
		ids := make([]string, 0, len(steps))
		for _, step := range steps {
			ids = append(ids, step.ID)
			included = append(included, stepResource(step))
		}
		resource.Relationships["step"] = jsonapi.LinkList("step", ids...)
	}

	if want["workspace"] {
		ws, err := a.deps.Store.GetWorkspace(r.Context(), job.WorkspaceID)
		if err != nil {
			if !a.writeStoreError(w, "get workspace", err) {
				a.internalError(w, "get workspace", err)
			}
			return
		}
		included = append(included, workspaceResource(ws))
	}

	jsonapi.Write(w, http.StatusOK, jsonapi.Document{Data: resource, Included: included})
}

// updateJob serves the run page's cancel, approve and reject buttons. Upstream
// answers 204 with no body, so this cannot return the updated resource.
func (a *API) updateJob(w http.ResponseWriter, r *http.Request) {
	job, ok := a.jobInOrganization(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	if !requireType(w, data, "job") {
		return
	}
	if !jsonapi.HasAttribute(data, "status") {
		// Upstream accepts an empty attributes member as a no-op.
		w.WriteHeader(http.StatusNoContent)
		return
	}

	attrs, err := jsonapi.DecodeAttributes[jobAttrs](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Attributes", err.Error())
		return
	}
	if attrs.Status == nil {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "This operation conflicts with data that already exists.")
		return
	}
	status := *attrs.Status
	if !jobStatuses[status] {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+status)
		return
	}
	if status == "approved" && job.Status != "waitingApproval" {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict",
			"Job "+strconv.FormatInt(job.ID, 10)+" cannot be approved from status '"+job.Status+"': only a job that is waiting for approval can be approved.")
		return
	}

	// The same call the TFE discard action makes, so jobstatus.Wrap publishes the
	// change to the GraphQL subscriptions.
	if _, err := a.deps.Store.UpdateJobStatus(r.Context(), store.UpdateJobStatusParams{ID: job.ID, Status: status}); err != nil {
		if !a.writeStoreError(w, "update job status", err) {
			a.internalError(w, "update job status", err)
		}
		return
	}
	_ = a.deps.Store.SetWorkspaceLastJob(r.Context(), store.SetWorkspaceLastJobParams{
		ID:            job.WorkspaceID,
		LastJobStatus: status,
	})
	w.WriteHeader(http.StatusNoContent)
}

// jobInOrganization loads the job and proves it belongs to the organization in
// the path, the same way the nested workspace routes do.
func (a *API) jobInOrganization(w http.ResponseWriter, r *http.Request) (store.Job, bool) {
	id, err := strconv.ParseInt(param(r, "jobID"), 10, 64)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+param(r, "jobID"))
		return store.Job{}, false
	}
	job, err := a.deps.Store.GetJob(r.Context(), id)
	if err != nil {
		if !a.writeStoreError(w, "get job", err) {
			a.internalError(w, "get job", err)
		}
		return store.Job{}, false
	}
	if job.OrganizationID != param(r, "organizationID") {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return store.Job{}, false
	}
	return job, true
}

// jobResource carries every attribute upstream exposes, including the ones whose
// columns we do not have yet: they are null (or the column default) rather than
// missing, so a client that reads them sees null instead of undefined.
func jobResource(job store.Job) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "job",
		ID:   strconv.FormatInt(job.ID, 10),
		Attributes: map[string]any{
			"approvalTeam":                   nil,
			"approvedAt":                     timestamp(job.ApprovedAt),
			"approvedBy":                     job.ApprovedBy,
			"cascadeDepth":                   0,
			"comments":                       job.Comments,
			"commitId":                       job.CommitID,
			"createdBy":                      job.CreatedBy,
			"createdDate":                    timestamp(job.CreatedDate),
			"dispatchFailureCount":           0,
			"dispatchFirstFailureAt":         nil,
			"dispatchLastAttemptAt":          nil,
			"dispatchLastError":              nil,
			"dispatchNextRetryAt":            nil,
			"output":                         job.Output,
			"overrideBranch":                 job.OverrideBranch,
			"overrideSource":                 job.OverrideSource,
			"planChanges":                    job.PlanChanges,
			"prCommentError":                 nil,
			"prNumber":                       nil,
			"refresh":                        job.Refresh,
			"refreshOnly":                    false,
			"replaceAddrs":                   nil,
			"resolvedTerraformVersion":       nil,
			"status":                         job.Status,
			"targetAddrs":                    nil,
			"templateReference":              job.TemplateReference,
			"terraformPlan":                  nil,
			"terraformPlanArtifacts":         nil,
			"terraformPlanArtifactsChecksum": nil,
			"triggeredByJobId":               nil,
			"updatedBy":                      job.UpdatedBy,
			"updatedDate":                    timestamp(job.UpdatedDate),
			"variableNames":                  nil,
			"variablesAsArguments":           false,
			"via":                            job.Via,
		},
		Relationships: map[string]jsonapi.Relationship{
			"organization": jsonapi.Link("organization", job.OrganizationID),
			"workspace":    jsonapi.Link("workspace", job.WorkspaceID),
		},
	}
}

// stepResource is the shape upstream returns inside include=step. Our step rows
// carry a single log status rather than the delivery bookkeeping, so the
// delivery fields are the durable default upstream uses.
func stepResource(step store.Step) jsonapi.Resource {
	return jsonapi.Resource{
		Type: "step",
		ID:   step.ID,
		Attributes: map[string]any{
			"deliveryRevision":  nil,
			"deliveryState":     "NOT_DURABLE",
			"deliveryTimestamp": nil,
			"name":              step.Name,
			"output":            step.Output,
			"status":            step.Status,
			"stepNumber":        step.StepNumber,
		},
		Relationships: map[string]jsonapi.Relationship{
			"job": jsonapi.Link("job", strconv.FormatInt(step.JobID, 10)),
		},
	}
}

func timestamp(t pgtype.Timestamptz) any {
	if !t.Valid {
		return nil
	}
	return t.Time.UTC().Format(time.RFC3339)
}
