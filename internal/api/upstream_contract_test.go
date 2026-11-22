package api

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// upstreamCapturesDir holds responses recorded from a live upstream
// (azbuilder/api-server:2.34.3) by scripts/upstream-capture. They are the
// contract, so they are replayed here rather than trusted to a reading of the
// OpenAPI file.
const upstreamCapturesDir = "../../docs/upstream/rest"

type upstreamCapture struct {
	name   string
	status int
	path   string
	body   []byte
}

func (c upstreamCapture) resourceID(t *testing.T) string {
	t.Helper()
	var doc struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(c.body, &doc); err != nil || doc.Data.ID == "" {
		t.Fatalf("%s: no data.id to replay against: %v", c.name, err)
	}
	return doc.Data.ID
}

func loadUpstreamCaptures(t *testing.T) map[string]upstreamCapture {
	t.Helper()

	index, err := os.ReadFile(filepath.Join(upstreamCapturesDir, "index.tsv"))
	if err != nil {
		t.Fatalf("read captured index: %v", err)
	}

	captures := make(map[string]upstreamCapture)
	for _, line := range strings.Split(strings.TrimSpace(string(index)), "\n") {
		name, rest, ok := strings.Cut(line, "\t")
		rawStatus, path, ok2 := strings.Cut(rest, "\t")
		if !ok || !ok2 {
			t.Fatalf("malformed index line %q, want name<TAB>status<TAB>path", line)
		}
		status, err := strconv.Atoi(rawStatus)
		if err != nil {
			t.Fatalf("%s: status %q is not a number", name, rawStatus)
		}
		body, err := os.ReadFile(filepath.Join(upstreamCapturesDir, name+".json"))
		if err != nil {
			t.Fatalf("%s: read capture: %v", name, err)
		}
		captures[name] = upstreamCapture{name: name, status: status, path: path, body: body}
	}
	return captures
}

// replay is how this API answers one captured request, with stub data standing
// in for the database. target "resource" means "/organization/{captured id}".
type replay struct {
	target string
	// status overrides the captured status when we deliberately answer
	// differently; 0 means it must match.
	status  int
	querier *stubQuerier
	// authToken enables authentication, for the unauthenticated capture.
	authToken string
	// idFrom names another capture that carries the resource id, for paths
	// built against an organization described by a different request.
	idFrom string
	// method defaults to GET; body is a template in which {id} is replaced with
	// the captured resource id (PATCH cases).
	method string
	body   string
	// accept overrides both media-type headers, for the atomic extension.
	accept string
}

var replayedCaptures = map[string]replay{
	"organization-list": {
		target:  "/organization",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-list-sorted": {
		target:  "/organization?sort=-name",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-list-filtered": {
		target:  "/organization?filter[organization]=name==azure",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-get": {
		target: "resource",
	},
	"organization-get-fields": {
		target:  "resource-query",
		querier: &stubQuerier{org: store.Organization{ID: capturedOrgID, Name: "azure", Description: strPtr("azure sample organization")}},
	},
	"organization-get-fields-empty": {
		target:  "resource-query",
		querier: &stubQuerier{org: store.Organization{ID: capturedOrgID, Name: "azure"}},
	},
	"organization-get-fields-other": {
		target:  "resource-query",
		querier: &stubQuerier{org: store.Organization{ID: capturedOrgID, Name: "azure"}},
	},
	"organization-list-fields": {
		target:  "/organization?fields[organization]=name",
		querier: &stubQuerier{orgs: demoOrganizations()},
	},
	"organization-projects": {
		target:  "resource-projects",
		idFrom:  "organization-get",
		querier: &stubQuerier{},
	},
	"organization-missing": {
		target:  "/organization/00000000-0000-0000-0000-000000000000",
		querier: &stubQuerier{err: pgx.ErrNoRows},
	},
	"organization-bad-id": {
		target: "/organization/not-a-uuid",
	},
	"organization-no-token": {
		target:    "/organization",
		authToken: "secret",
	},
	"organization-job-list": {
		target:  "/organization/" + capturedOrgID + "/job",
		querier: &stubQuerier{jobs: captureJobs()},
	},
	"organization-job-get": {
		target:  "job",
		querier: &stubQuerier{jobs: captureJobs()},
	},
	"organization-job-include": {
		target:  "job-include",
		querier: &stubQuerier{jobs: captureJobs(), steps: captureSteps(), workspaces: captureWorkspaces()},
	},
	"organization-job-patch": {
		target:  "job",
		idFrom:  "organization-job-get",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"job","id":"{id}","attributes":{"status":"cancelled"}}}`,
		querier: &stubQuerier{jobs: captureJobs()},
	},
	"organization-job-patch-invalid": {
		target:  "job",
		idFrom:  "organization-job-get",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"job","id":"{id}","attributes":{"status":"nonsense"}}}`,
		querier: &stubQuerier{jobs: captureJobs()},
	},
	"workspace-patch-attribute": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"workspace","id":"` + capturedRelWS + `","attributes":{"description":"captured"}}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    workspaceProjectBody(capturedRelP1),
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project-moved": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    workspaceProjectBody(capturedRelP2),
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project-foreign": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    workspaceProjectBody(capturedRelForeign),
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project-unknown": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    workspaceProjectBody(capturedRelUnknown),
		querier: captureWorkspaceStore(),
	},
	"workspace-patch-project-clear": {
		target:  "workspace",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"workspace","id":"` + capturedRelWS + `","relationships":{"project":{"data":null}}}}`,
		querier: captureWorkspaceStore(),
	},
	// The recording was made against a project of another organization; which
	// id it was is not part of the contract, so the replay names its own.
	"workspace-create-project-foreign": {
		target:  "workspace-collection",
		method:  http.MethodPost,
		body:    `{"data":{"type":"workspace","attributes":{"name":"from-vcs","source":"https://github.com/terrakube-io/terrakube.git","branch":"main","terraformVersion":"1.9.0"},"relationships":{"project":{"data":{"type":"project","id":"` + capturedRelForeign + `"}}}}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-project-set": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"project","id":"` + capturedRelP1 + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-project-clear": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-project-foreign": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"project","id":"` + capturedRelForeign + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-project-unknown": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"project","id":"` + capturedRelUnknown + `"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-vcs-set": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"vcs","id":"aaaaaaaa-0000-4000-8000-000000000001"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-vcs-clear": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-ssh-set": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"ssh","id":"aaaaaaaa-0000-4000-8000-000000000002"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-ssh-clear": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-agent-set": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"agent","id":"aaaaaaaa-0000-4000-8000-000000000003"}}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-agent-clear": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"workspace-relationship-unknown": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":null}`,
		querier: captureWorkspaceStore(),
	},
	"organization-atomic-create": {
		target:  "resource-query",
		method:  http.MethodPost,
		body:    `{"atomic:operations":[{"op":"add","href":"/organization/` + capturedRelOrg + `/workspace","data":{"type":"workspace","lid":"11111111-1111-4111-8111-111111111111","attributes":{"source":"https://github.com/terrakube-io/terrakube.git","folder":"/","name":"capture-atomic","terraformVersion":"1.9.0","branch":"main","iacType":"terraform","executionMode":"remote"}}}]}`,
		accept:  atomicMedia,
		querier: captureWorkspaceStore(),
	},
	"organization-atomic-update": {
		target:  "resource-query",
		method:  http.MethodPost,
		body:    `{"atomic:operations":[{"op":"update","href":"/organization/` + capturedRelOrg + `/workspace/` + capturedRelWS + `","data":{"type":"workspace","id":"` + capturedRelWS + `","attributes":{"description":"atomic-after"}}}]}`,
		accept:  atomicMedia,
		querier: captureWorkspaceStore(),
	},
	"organization-globalvar-list": {
		target:  "globalvar-collection",
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-get": {
		target:  "globalvar-item",
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-get-sensitive": {
		target:  "globalvar-secret",
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-create": {
		target:  "globalvar-collection",
		method:  http.MethodPost,
		body:    `{"data":{"type":"globalvar","attributes":{"key":"CAPTURE_NEW","value":"new","description":"capture","sensitive":false,"hcl":false,"category":"ENV"}}}`,
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-patch": {
		target:  "globalvar-item",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"globalvar","id":"` + capturedGvID + `","attributes":{"key":"CAPTURE_ONE","value":"after","description":"after","hcl":true,"category":"TERRAFORM"}}}`,
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-delete": {
		target:  "globalvar-secret",
		method:  http.MethodDelete,
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-invalid-category": {
		target:  "globalvar-collection",
		method:  http.MethodPost,
		body:    `{"data":{"type":"globalvar","attributes":{"key":"CAPTURE_BAD","value":"x","description":"d","sensitive":false,"hcl":false,"category":"env"}}}`,
		querier: captureGlobalVarStore(),
	},
	"organization-globalvar-missing-key": {
		target:  "globalvar-collection",
		method:  http.MethodPost,
		body:    `{"data":{"type":"globalvar","attributes":{"value":"x","description":"d","sensitive":false,"hcl":false,"category":"ENV"}}}`,
		querier: captureGlobalVarStore(),
	},
	"action-list": {
		target:  "resource-query",
		querier: captureActionStore(),
	},
	"action-list-filtered": {
		target:  "resource-query",
		querier: captureActionStore(),
	},
	"action-get": {
		target:  "resource-query",
		querier: captureActionStore(),
	},
	"action-fields": {
		target:  "resource-query",
		querier: captureActionStore(),
	},
	"action-missing": {
		target:  "resource-query",
		querier: &stubQuerier{},
	},
	"action-create": {
		target:  "resource-query",
		method:  http.MethodPost,
		body:    `{"data":{"type":"action","id":"capture.temp-action","attributes":{"name":"Capture Temp","label":"Temp","category":"General","description":"capture","type":"Workspace/Action","action":"dGVzdA==","displayCriteria":"[{\"filter\":\"true\"}]","version":"1.0.0","active":true}}}`,
		querier: captureActionStore(),
	},
	"action-patch": {
		target:  "resource-query",
		method:  http.MethodPatch,
		body:    `{"data":{"type":"action","id":"capture.temp-action","attributes":{"label":"Temp after","active":false}}}`,
		querier: captureActionStore(),
	},
	"action-delete": {
		target:  "resource-query",
		method:  http.MethodDelete,
		querier: captureActionStore(),
	},
}

// deferredCaptures are recorded but cannot be asserted yet. Each needs a reason,
// and the test below fails when a capture is added without one, so a new
// recording cannot be silently ignored.
var deferredCaptures = map[string]string{
	"organization-list-totals":         "we accept paging on organization, upstream answers 400; deliberate superset (docs/upstream/rest/README.md)",
	"organization-list-paged":          "same as organization-list-totals",
	"organization-list-include":        "we ignore include= and emit no `included` member; the UI sends no include on /api/v1",
	"organization-relationships":       "no /relationships/{relation} routes; the UI patches relationships on workspace only (#116)",
	"workspace-list-paged":             "no global /workspace collection; the capture documents meta.page, not a route we serve",
	"workspace-list-totals":            "same as workspace-list-paged",
	"delete-method-not-allowed":        "upstream answers 500 with a servlet error page, which is not a contract worth copying",
	"organization-globalvar-duplicate": "upstream's duplicate-key lifecycle hook throws and Spring answers 500; our unique constraint answers 409",
	"organization-atomic-unknown":      "upstream answers 400 with an array of per-operation results; we answer one JSON:API error object",
}

// The ids upstream's demo seed gave the captures; the job captures were taken
// against these rows, so the replay uses the same ids.
const (
	capturedOrgID = "1c839e0a-ca5d-49a9-9c78-7a569df534b7"
	capturedWSID  = "0beb7db4-b073-4bd4-9a98-94fe7e2d9939"
	// The workspace relationship captures were seeded against this org, which
	// the upstream demo gives two projects.
	capturedRelOrg     = "70000000-0000-0000-0000-000000000001"
	capturedRelWS      = "77777777-7777-4777-8777-777777777777"
	capturedRelP1      = "70000000-0000-0000-0000-000000000010"
	capturedRelP2      = "70000000-0000-0000-0000-000000000020"
	capturedRelForeign = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	capturedRelUnknown = "00000000-0000-0000-0000-000000000000"
	capturedForeignOrg = "d9b58bd3-f3fc-4056-a026-1163297e80a8"
	capturedGvID       = "99999999-9999-4999-8999-999999999990"
	capturedGvSecret   = "99999999-9999-4999-8999-999999999993"
	capturedStepID     = "66666666-6666-4666-8666-666666666666"
)

func workspaceProjectBody(projectID string) string {
	return `{"data":{"type":"workspace","id":"` + capturedRelWS +
		`","relationships":{"project":{"data":{"type":"project","id":"` + projectID + `"}}}}}`
}

func captureWorkspaceStore() *stubQuerier {
	return &stubQuerier{
		workspaces: []store.Workspace{{
			ID: capturedRelWS, OrganizationID: capturedRelOrg, Name: "capture-rel-workspace",
			ProjectID: strPtr(capturedRelP1),
		}},
		projects: []store.Project{
			{ID: capturedRelP1, Name: "governance-terraform", OrganizationID: capturedRelOrg},
			{ID: capturedRelP2, Name: "governance-tofu", OrganizationID: capturedRelOrg},
			{ID: capturedRelForeign, Name: "demo-project", OrganizationID: capturedForeignOrg},
		},
	}
}

func captureGlobalVarStore() *stubQuerier {
	return &stubQuerier{
		globals: []store.Globalvar{
			{ID: capturedGvID, VariableKey: "CAPTURE_ONE", VariableValue: "one", VariableDescription: "capture", VariableCategory: "ENV", OrganizationID: capturedRelOrg},
			{ID: capturedGvSecret, VariableKey: "CAPTURE_SECRET", VariableValue: "hidden", VariableDescription: "capture", VariableCategory: "ENV", Sensitive: true, OrganizationID: capturedRelOrg},
			{ID: "99999999-9999-4999-8999-999999999992", VariableKey: "CAPTURE_TWO", VariableValue: "two", VariableDescription: "capture", VariableCategory: "TERRAFORM", Hcl: true, OrganizationID: capturedRelOrg},
		},
	}
}

func captureActionStore() *stubQuerier {
	return &stubQuerier{actions: []store.Action{
		{ID: "terrakube.open-documentation", Name: "Open Documentation", Label: "Open Documentation", Category: "General", Type: "Workspace/ResourceDrawer/Action", Action: "dGVzdA==", DisplayCriteria: `[{"filter": "true"}]`, Version: "1.0.0", Active: true},
		{ID: "terrakube.resource-details", Name: "Resource Details", Label: "Details", Category: "General", Type: "Workspace/ResourceDrawer/Tab", Action: "dGVzdA==", DisplayCriteria: `[{"filter": "true"}]`, Version: "1.0.0", Active: true},
		{ID: "terrakube.open-ai", Name: "Open AI", Label: "Open AI", Category: "AI", Type: "Workspace/Action", Action: "dGVzdA==", DisplayCriteria: `[{"filter": "true"}]`, Version: "0.1.0", Active: false},
		{ID: "capture.temp-action", Name: "Capture Temp", Label: "Temp", Category: "General", Type: "Workspace/Action", Action: "dGVzdA==", DisplayCriteria: `[{"filter":"true"}]`, Version: "1.0.0", Active: true},
	}}
}

func captureJobs() []store.Job {
	return []store.Job{{
		ID:             2,
		OrganizationID: capturedOrgID,
		WorkspaceID:    capturedWSID,
		Status:         "completed",
		Via:            "UI",
	}}
}

// captureJobStore answers for whichever job id the job capture was taken
// against; the serial changes on every capture run.
func captureJobStore(id string) *stubQuerier {
	jobID, _ := strconv.ParseInt(id, 10, 64)
	return &stubQuerier{
		jobs: []store.Job{{
			ID:             jobID,
			OrganizationID: capturedOrgID,
			WorkspaceID:    capturedWSID,
			Status:         "completed",
			Via:            "UI",
		}},
		steps: []store.Step{{
			ID:         capturedStepID,
			JobID:      jobID,
			StepNumber: 1,
			Name:       strPtr("Plan"),
			Status:     strPtr("completed"),
			Output:     strPtr("ok"),
		}},
		workspaces: captureWorkspaces(),
	}
}

func captureSteps() []store.Step {
	return []store.Step{{
		ID:         capturedStepID,
		JobID:      2,
		StepNumber: 1,
		Name:       strPtr("Plan"),
		Status:     strPtr("completed"),
		Output:     strPtr("ok"),
	}}
}

func captureWorkspaces() []store.Workspace {
	return []store.Workspace{{ID: capturedWSID, OrganizationID: capturedOrgID, Name: "sample_cosmosdb"}}
}

func (s *stubQuerier) ListJobsByOrganization(context.Context, string) ([]store.Job, error) {
	return s.jobs, s.err
}

func (s *stubQuerier) GetJob(_ context.Context, id int64) (store.Job, error) {
	for _, job := range s.jobs {
		if job.ID == id {
			return job, nil
		}
	}
	return store.Job{}, pgx.ErrNoRows
}

func (s *stubQuerier) ListStepsByJob(_ context.Context, jobID int64) ([]store.Step, error) {
	var out []store.Step
	for _, step := range s.steps {
		if step.JobID == jobID {
			out = append(out, step)
		}
	}
	return out, nil
}

func (s *stubQuerier) UpdateJobStatus(_ context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	for i, job := range s.jobs {
		if job.ID == arg.ID {
			s.jobs[i].Status = arg.Status
			return s.jobs[i], nil
		}
	}
	return store.Job{}, pgx.ErrNoRows
}

func (s *stubQuerier) SetWorkspaceLastJob(context.Context, store.SetWorkspaceLastJobParams) error {
	return s.err
}

func (s *stubQuerier) ListGlobalVariablesByOrganization(_ context.Context, orgID string) ([]store.Globalvar, error) {
	var out []store.Globalvar
	for _, variable := range s.globals {
		if variable.OrganizationID == orgID {
			out = append(out, variable)
		}
	}
	return out, s.err
}

func (s *stubQuerier) GetGlobalVariable(_ context.Context, id string) (store.Globalvar, error) {
	for _, variable := range s.globals {
		if variable.ID == id {
			return variable, nil
		}
	}
	return store.Globalvar{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateGlobalVariable(_ context.Context, arg store.CreateGlobalVariableParams) (store.Globalvar, error) {
	return store.Globalvar{
		ID: arg.ID, VariableKey: arg.VariableKey, VariableValue: arg.VariableValue,
		VariableDescription: arg.VariableDescription, VariableCategory: arg.VariableCategory,
		Sensitive: arg.Sensitive, Hcl: arg.Hcl, OrganizationID: arg.OrganizationID,
	}, s.err
}

func (s *stubQuerier) UpdateGlobalVariable(_ context.Context, arg store.UpdateGlobalVariableParams) (store.Globalvar, error) {
	return store.Globalvar{
		ID: arg.ID, VariableKey: arg.VariableKey, VariableValue: arg.VariableValue,
		VariableDescription: arg.VariableDescription, VariableCategory: arg.VariableCategory,
		Sensitive: arg.Sensitive, Hcl: arg.Hcl, OrganizationID: capturedRelOrg,
	}, s.err
}

func (s *stubQuerier) DeleteGlobalVariable(context.Context, string) error {
	return s.err
}

func (s *stubQuerier) GetProject(_ context.Context, id string) (store.Project, error) {
	for _, project := range s.projects {
		if project.ID == id {
			return project, nil
		}
	}
	return store.Project{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateWorkspace(_ context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
	return store.Workspace{
		ID: arg.ID, Name: arg.Name, OrganizationID: arg.OrganizationID, ProjectID: arg.ProjectID,
	}, s.err
}

func (s *stubQuerier) UpdateWorkspace(_ context.Context, arg store.UpdateWorkspaceParams) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == arg.ID {
			ws.ProjectID = arg.ProjectID
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *stubQuerier) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

// ListProjectsByOrganization is stubbed here because only the capture replay
// needs it; the organization handlers above it do not.
func (s *stubQuerier) ListProjectsByOrganization(context.Context, string) ([]store.Project, error) {
	return nil, s.err
}

// demoOrganizations stands in for upstream's demo seed; only the names matter,
// and no assertion here reads them.
func demoOrganizations() []store.Organization {
	return []store.Organization{
		{ID: "1", Name: "azure", ExecutionMode: "remote"},
		{ID: "2", Name: "aws", ExecutionMode: "local"},
	}
}

func TestEveryUpstreamCaptureIsAccountedFor(t *testing.T) {
	for name := range loadUpstreamCaptures(t) {
		_, replayed := replayedCaptures[name]
		reason, deferred := deferredCaptures[name]
		if !replayed && !deferred {
			t.Errorf("%s: captured from upstream but neither replayed nor deferred; add it to replayedCaptures or explain it in deferredCaptures", name)
		}
		if deferred && strings.TrimSpace(reason) == "" {
			t.Errorf("%s: deferred without a reason", name)
		}
	}
}

func TestCapturedShapesAreWhatWeReturn(t *testing.T) {
	captures := loadUpstreamCaptures(t)

	for name, plan := range replayedCaptures {
		t.Run(name, func(t *testing.T) {
			capture, ok := captures[name]
			if !ok {
				t.Fatalf("no capture named %s", name)
			}

			target, querier := plan.target, plan.querier
			idCapture := capture
			if plan.idFrom != "" {
				idCapture = captures[plan.idFrom]
			}
			switch target {
			case "resource":
				id := idCapture.resourceID(t)
				target = "/organization/" + id
				querier = &stubQuerier{org: store.Organization{ID: id, Name: "azure", ExecutionMode: "remote"}}
			case "resource-query":
				target = strings.TrimPrefix(capture.path, "/api/v1")
			case "resource-projects":
				target = "/organization/" + idCapture.resourceID(t) + "/project"
			case "job":
				target = "/organization/" + capturedOrgID + "/job/" + idCapture.resourceID(t)
				querier = captureJobStore(idCapture.resourceID(t))
			case "job-include":
				target = "/organization/" + capturedOrgID + "/job/" + idCapture.resourceID(t) + "?include=step,workspace"
				querier = captureJobStore(idCapture.resourceID(t))
			case "workspace":
				target = "/organization/" + capturedRelOrg + "/workspace/" + capturedRelWS
			case "workspace-collection":
				target = "/organization/" + capturedRelOrg + "/workspace"
			case "globalvar-collection":
				target = "/organization/" + capturedRelOrg + "/globalvar"
			case "globalvar-item":
				target = "/organization/" + capturedRelOrg + "/globalvar/" + capturedGvID
			case "globalvar-secret":
				target = "/organization/" + capturedRelOrg + "/globalvar/" + capturedGvSecret
			}
			if querier == nil {
				querier = &stubQuerier{}
			}

			handler := New(Deps{
				Store:     querier,
				Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
				AuthToken: plan.authToken,
			})
			method := plan.method
			if method == "" {
				method = http.MethodGet
			}
			body := plan.body
			if strings.Contains(body, "{id}") {
				body = strings.ReplaceAll(body, "{id}", idCapture.resourceID(t))
			}
			req := httptest.NewRequest(method, target, strings.NewReader(body))
			if body != "" {
				req.Header.Set("Content-Type", "application/vnd.api+json")
			}
			if plan.accept != "" {
				req.Header.Set("Content-Type", plan.accept)
				req.Header.Set("Accept", plan.accept)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			wantStatus := capture.status
			if plan.status != 0 {
				wantStatus = plan.status
			}
			if rec.Code != wantStatus {
				t.Fatalf("%s: status = %d, want %d (captured); body = %s", capture.path, rec.Code, wantStatus, rec.Body.String())
			}
			assertSameDocumentShape(t, capture, rec.Body.Bytes())
		})
	}
}

// assertSameDocumentShape compares the members a client can see, not values: an
// attribute rename or a missing member must fail, a different id or name must
// not.
func assertSameDocumentShape(t *testing.T, capture upstreamCapture, ours []byte) {
	t.Helper()

	if len(strings.TrimSpace(string(capture.body))) == 0 {
		// Upstream's 401 has an empty body; only the status is comparable.
		return
	}

	var theirs, mine map[string]any
	if err := json.Unmarshal(capture.body, &theirs); err != nil {
		t.Fatalf("%s: captured body is not JSON: %v", capture.name, err)
	}
	if err := json.Unmarshal(ours, &mine); err != nil {
		t.Fatalf("%s: our body is not JSON (%v): %s", capture.name, err, ours)
	}

	if theirErrors, ok := theirs["errors"].([]any); ok {
		// We add status and title; JSON:API allows it and the UI ignores them.
		myErrors, ok := mine["errors"].([]any)
		if !ok || len(myErrors) != len(theirErrors) {
			t.Fatalf("%s: errors = %v, want %d error object(s)", capture.name, mine["errors"], len(theirErrors))
		}
		for i, raw := range theirErrors {
			theirsErr, _ := raw.(map[string]any)
			mineErr, _ := myErrors[i].(map[string]any)
			if theirsErr["detail"] != nil && mineErr["detail"] == nil {
				t.Errorf("%s: errors[%d] has no detail; upstream says %v", capture.name, i, theirsErr["detail"])
			}
		}
		return
	}

	if got, want := sortedKeys(mine), sortedKeys(theirs); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("%s: top-level members = %v, want %v", capture.name, got, want)
	}

	switch their := theirs["data"].(type) {
	case []any:
		my, ok := mine["data"].([]any)
		if !ok {
			t.Fatalf("%s: data = %T, want array", capture.name, mine["data"])
		}
		if len(their) == 0 {
			if len(my) != 0 {
				t.Fatalf("%s: upstream returned an empty collection, we returned %d", capture.name, len(my))
			}
			return
		}
		if len(my) == 0 {
			t.Fatalf("%s: we returned an empty collection, upstream returned %d", capture.name, len(their))
		}
		assertSameResourceShape(t, capture.name, their[0], my[0])
	case map[string]any:
		assertSameResourceShape(t, capture.name, their, mine["data"])
	}
}

func assertSameResourceShape(t *testing.T, name string, theirs, mine any) {
	t.Helper()

	theirResource, ok := theirs.(map[string]any)
	if !ok {
		t.Fatalf("%s: captured data is %T, want a resource object", name, theirs)
	}
	myResource, ok := mine.(map[string]any)
	if !ok {
		t.Fatalf("%s: our data is %T, want a resource object", name, mine)
	}

	if myResource["type"] != theirResource["type"] {
		t.Errorf("%s: type = %v, want %v", name, myResource["type"], theirResource["type"])
	}
	if _, ok := myResource["id"].(string); !ok {
		t.Errorf("%s: id = %v, want a string", name, myResource["id"])
	}

	theirAttrs, _ := theirResource["attributes"].(map[string]any)
	myAttrs, _ := myResource["attributes"].(map[string]any)
	if got, want := sortedKeys(myAttrs), sortedKeys(theirAttrs); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("%s: attributes = %v, want %v", name, got, want)
	}
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *stubQuerier) ListActions(context.Context) ([]store.Action, error) {
	return s.actions, s.err
}

func (s *stubQuerier) GetAction(_ context.Context, id string) (store.Action, error) {
	for _, action := range s.actions {
		if action.ID == id {
			return action, nil
		}
	}
	return store.Action{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateAction(_ context.Context, arg store.CreateActionParams) (store.Action, error) {
	return store.Action{
		ID: arg.ID, Name: arg.Name, Label: arg.Label, Category: arg.Category,
		Description: arg.Description, Type: arg.Type, Action: arg.Action,
		DisplayCriteria: arg.DisplayCriteria, Version: arg.Version, Active: arg.Active,
	}, s.err
}

func (s *stubQuerier) UpdateAction(_ context.Context, arg store.UpdateActionParams) (store.Action, error) {
	for _, action := range s.actions {
		if action.ID == arg.ID {
			return action, s.err
		}
	}
	return store.Action{}, pgx.ErrNoRows
}

func (s *stubQuerier) DeleteAction(context.Context, string) error {
	return s.err
}
