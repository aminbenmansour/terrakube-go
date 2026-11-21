package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// jobStub is an in-memory job store. The job handlers must scope reads to the
// organization in the path and go through UpdateJobStatus (the path the
// jobstatus wrapper publishes) - a stub that answers blindly would hide both.
type jobStub struct {
	store.Querier
	jobs       []store.Job
	steps      []store.Step
	workspaces []store.Workspace
	lastJob    string
	err        error
}

func (s *jobStub) GetJob(_ context.Context, id int64) (store.Job, error) {
	for _, job := range s.jobs {
		if job.ID == id {
			return job, nil
		}
	}
	return store.Job{}, pgx.ErrNoRows
}

func (s *jobStub) ListJobsByOrganization(_ context.Context, orgID string) ([]store.Job, error) {
	var out []store.Job
	for _, job := range s.jobs {
		if job.OrganizationID == orgID {
			out = append(out, job)
		}
	}
	return out, s.err
}

func (s *jobStub) ListStepsByJob(_ context.Context, jobID int64) ([]store.Step, error) {
	var out []store.Step
	for _, step := range s.steps {
		if step.JobID == jobID {
			out = append(out, step)
		}
	}
	return out, nil
}

func (s *jobStub) UpdateJobStatus(_ context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	for i, job := range s.jobs {
		if job.ID != arg.ID {
			continue
		}
		s.jobs[i].Status = arg.Status
		return s.jobs[i], nil
	}
	return store.Job{}, pgx.ErrNoRows
}

func (s *jobStub) SetWorkspaceLastJob(_ context.Context, arg store.SetWorkspaceLastJobParams) error {
	s.lastJob = arg.LastJobStatus
	return nil
}

func (s *jobStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

const (
	jobOrg  = "11111111-1111-4111-8111-111111111111"
	jobID   = int64(7)
	jobWSID = "22222222-2222-4222-8222-222222222222"
)

func newJobStub() *jobStub {
	return &jobStub{
		jobs: []store.Job{{
			ID:             jobID,
			OrganizationID: jobOrg,
			WorkspaceID:    jobWSID,
			Status:         "completed",
			Via:            "UI",
		}},
		steps: []store.Step{{
			ID:         "66666666-6666-4666-8666-666666666666",
			JobID:      jobID,
			StepNumber: 1,
			Name:       strPtr("Plan"),
			Status:     strPtr("completed"),
			Output:     strPtr("ok"),
		}},
		workspaces: []store.Workspace{{ID: jobWSID, OrganizationID: jobOrg, Name: "demo"}},
	}
}

func serveJob(t *testing.T, s *jobStub, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	newTestAPI(t, s).ServeHTTP(rec, req)
	return rec
}

func TestJobGetMatchesTheOrganizationInThePath(t *testing.T) {
	rec := serveJob(t, newJobStub(), http.MethodGet, "/organization/"+jobOrg+"/job/7", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Type != "job" || doc.Data.ID != "7" {
		t.Errorf("data = %s/%s, want job/7", doc.Data.Type, doc.Data.ID)
	}
	if doc.Data.Attributes["status"] != "completed" {
		t.Errorf("status = %v, want completed", doc.Data.Attributes["status"])
	}
}

func TestJobGetFromAnotherOrganizationIsNotFound(t *testing.T) {
	rec := serveJob(t, newJobStub(), http.MethodGet,
		"/organization/99999999-9999-4999-8999-999999999999/job/7", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestJobListReturnsTheOrganizationsJobs(t *testing.T) {
	rec := serveJob(t, newJobStub(), http.MethodGet, "/organization/"+jobOrg+"/job", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 1 || doc.Data[0].ID != "7" {
		t.Fatalf("data = %v, want job 7", doc.Data)
	}
}

func TestJobIncludeReturnsStepsAndWorkspace(t *testing.T) {
	rec := serveJob(t, newJobStub(), http.MethodGet,
		"/organization/"+jobOrg+"/job/7?include=step,workspace", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var step, workspace bool
	for _, inc := range doc.Included {
		switch inc.Type {
		case "step":
			step = inc.Attributes["name"] == "Plan" && inc.Attributes["stepNumber"] == float64(1)
		case "workspace":
			workspace = inc.ID == jobWSID
		}
	}
	if !step || !workspace {
		t.Errorf("included = %v, want the job's step and workspace", doc.Included)
	}
}

func TestJobCancelIsNoContentAndReachesTheStatusPath(t *testing.T) {
	s := newJobStub()
	rec := serveJob(t, s, http.MethodPatch, "/organization/"+jobOrg+"/job/7",
		`{"data":{"type":"job","id":"7","attributes":{"status":"cancelled"}}}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
	if s.jobs[0].Status != "cancelled" {
		t.Errorf("job status = %q, want cancelled", s.jobs[0].Status)
	}
	if s.lastJob != "cancelled" {
		t.Errorf("workspace last job = %q, want cancelled", s.lastJob)
	}
}

func TestJobPatchRejectsAnUnknownStatus(t *testing.T) {
	rec := serveJob(t, newJobStub(), http.MethodPatch, "/organization/"+jobOrg+"/job/7",
		`{"data":{"type":"job","id":"7","attributes":{"status":"nonsense"}}}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Errors) == 0 || !strings.Contains(doc.Errors[0].Detail, "Invalid value: nonsense") {
		t.Errorf("errors = %v, want Invalid value: nonsense", doc.Errors)
	}
}

func TestJobApprovalRequiresWaitingApproval(t *testing.T) {
	rec := serveJob(t, newJobStub(), http.MethodPatch, "/organization/"+jobOrg+"/job/7",
		`{"data":{"type":"job","id":"7","attributes":{"status":"approved"}}}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body = %s", rec.Code, rec.Body.String())
	}

	s := newJobStub()
	s.jobs[0].Status = "waitingApproval"
	rec = serveJob(t, s, http.MethodPatch, "/organization/"+jobOrg+"/job/7",
		`{"data":{"type":"job","id":"7","attributes":{"status":"approved"}}}`)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("waitingApproval approve: status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if s.jobs[0].Status != "approved" {
		t.Errorf("job status = %q, want approved", s.jobs[0].Status)
	}
}

// Upstream treats an empty attributes member as a no-op, not an error.
func TestJobPatchWithoutStatusChangesNothing(t *testing.T) {
	s := newJobStub()
	rec := serveJob(t, s, http.MethodPatch, "/organization/"+jobOrg+"/job/7",
		`{"data":{"type":"job","id":"7","attributes":{}}}`)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if s.jobs[0].Status != "completed" {
		t.Errorf("status = %q, want untouched", s.jobs[0].Status)
	}
}
