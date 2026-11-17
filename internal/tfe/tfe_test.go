package tfe

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

// stubJobQueue records enqueued run IDs for handler tests.
type stubJobQueue struct {
	enqueued []int64
	err      error
}

func (s *stubJobQueue) EnqueueJob(_ context.Context, jobID int64) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	s.enqueued = append(s.enqueued, jobID)
	return "1-0", nil
}

// stubQuerier embeds the generated interface so tests only implement the
// methods they exercise.
type stubQuerier struct {
	store.Querier

	org    store.Organization
	orgErr error

	workspaces         []store.Workspace
	workspacesErr      error
	workspace          store.Workspace
	workspaceErr       error
	workspaceByNameErr error
	createdWS          *store.CreateWorkspaceParams
	lockParams         *store.SetWorkspaceLockParams

	tags                     []store.Tag
	workspaceTags            []store.Workspacetag
	workspaceTagsByWorkspace map[string][]store.Workspacetag
	createdTag               *store.CreateTagParams
	createdWorkspaceTag      *store.CreateWorkspaceTagParams
	updatedWorkspaceTag      *store.UpdateWorkspaceTagValueParams

	projects    []store.Project
	projectsErr error
	project     store.Project
	projectErr  error
	createdProj *store.CreateProjectParams

	job            store.Job
	jobErr         error
	jobs           []store.Job
	updatedJob     *store.UpdateJobStatusParams
	step           store.Step
	history        store.History
	historyErr     error
	archives       map[string][]store.TempArchive
	createdArchive *store.CreateArchiveParams

	content              store.Content
	contentErr           error
	createdContent       *store.CreateContentParams
	updatedContentStatus *store.UpdateContentStatusParams
}

func (s *stubQuerier) GetOrganizationByName(context.Context, string) (store.Organization, error) {
	return s.org, s.orgErr
}

func (s *stubQuerier) ListWorkspacesByOrganization(context.Context, string) ([]store.Workspace, error) {
	return s.workspaces, s.workspacesErr
}

func (s *stubQuerier) GetWorkspaceByName(context.Context, store.GetWorkspaceByNameParams) (store.Workspace, error) {
	return s.workspace, s.workspaceByNameErr
}

func (s *stubQuerier) GetWorkspace(context.Context, string) (store.Workspace, error) {
	return s.workspace, s.workspaceErr
}

func (s *stubQuerier) CreateWorkspace(_ context.Context, arg store.CreateWorkspaceParams) (store.Workspace, error) {
	s.createdWS = &arg
	return s.workspace, nil
}

func (s *stubQuerier) UpdateWorkspace(context.Context, store.UpdateWorkspaceParams) (store.Workspace, error) {
	return s.workspace, nil
}

func (s *stubQuerier) SetWorkspaceLock(_ context.Context, arg store.SetWorkspaceLockParams) (store.Workspace, error) {
	s.lockParams = &arg
	ws := s.workspace
	ws.Locked = arg.Locked
	ws.LockDescription = arg.LockDescription
	return ws, nil
}

func (s *stubQuerier) ListTagsByOrganization(context.Context, string) ([]store.Tag, error) {
	return s.tags, nil
}

func (s *stubQuerier) ListWorkspaceTagsByWorkspace(_ context.Context, workspaceID string) ([]store.Workspacetag, error) {
	if s.workspaceTagsByWorkspace != nil {
		return s.workspaceTagsByWorkspace[workspaceID], nil
	}
	return s.workspaceTags, nil
}

func (s *stubQuerier) GetTagByName(_ context.Context, arg store.GetTagByNameParams) (store.Tag, error) {
	for _, tag := range s.tags {
		if tag.OrganizationID == arg.OrganizationID && tag.Name == arg.Name {
			return tag, nil
		}
	}
	return store.Tag{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateTag(_ context.Context, arg store.CreateTagParams) (store.Tag, error) {
	s.createdTag = &arg
	tag := store.Tag{ID: arg.ID, Name: arg.Name, OrganizationID: arg.OrganizationID}
	s.tags = append(s.tags, tag)
	return tag, nil
}

func (s *stubQuerier) CreateWorkspaceTag(_ context.Context, arg store.CreateWorkspaceTagParams) (store.Workspacetag, error) {
	s.createdWorkspaceTag = &arg
	if s.workspaceTagsByWorkspace == nil {
		s.workspaceTagsByWorkspace = map[string][]store.Workspacetag{}
	}
	workspaceTag := store.Workspacetag{
		ID:          arg.ID,
		TagID:       arg.TagID,
		TagValue:    arg.TagValue,
		WorkspaceID: arg.WorkspaceID,
	}
	s.workspaceTagsByWorkspace[arg.WorkspaceID] = append(s.workspaceTagsByWorkspace[arg.WorkspaceID], workspaceTag)
	return workspaceTag, nil
}

func (s *stubQuerier) UpdateWorkspaceTagValue(_ context.Context, arg store.UpdateWorkspaceTagValueParams) (store.Workspacetag, error) {
	s.updatedWorkspaceTag = &arg
	for workspaceID, list := range s.workspaceTagsByWorkspace {
		for i := range list {
			if list[i].ID == arg.ID {
				list[i].TagValue = arg.TagValue
				s.workspaceTagsByWorkspace[workspaceID] = list
				return list[i], nil
			}
		}
	}
	return store.Workspacetag{}, pgx.ErrNoRows
}

func (s *stubQuerier) CreateJob(_ context.Context, arg store.CreateJobParams) (store.Job, error) {
	s.job = store.Job{
		ID:             42,
		OrganizationID: arg.OrganizationID,
		WorkspaceID:    arg.WorkspaceID,
		Command:        arg.Command,
		Status:         arg.Status,
		Via:            arg.Via,
		AutoApply:      arg.AutoApply,
		Refresh:        arg.Refresh,
		PlanChanges:    arg.PlanChanges,
		CreatedBy:      arg.CreatedBy,
	}
	return s.job, nil
}

func (s *stubQuerier) GetJob(context.Context, int64) (store.Job, error) {
	if s.jobErr != nil {
		return store.Job{}, s.jobErr
	}
	if s.job.ID == 0 {
		return store.Job{}, pgx.ErrNoRows
	}
	return s.job, nil
}

func (s *stubQuerier) UpdateJobStatus(_ context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	s.updatedJob = &arg
	s.job.Status = arg.Status
	return s.job, nil
}

func (s *stubQuerier) ListActiveJobsByWorkspace(context.Context, string) ([]store.Job, error) {
	return s.jobs, nil
}

func (s *stubQuerier) ListActiveJobsByOrganization(context.Context, string) ([]store.Job, error) {
	return s.jobs, nil
}

func (s *stubQuerier) CreateStep(_ context.Context, arg store.CreateStepParams) (store.Step, error) {
	s.step = store.Step{
		ID:         arg.ID,
		JobID:      arg.JobID,
		StepNumber: arg.StepNumber,
		Name:       arg.Name,
		Status:     arg.Status,
		LogStatus:  arg.LogStatus,
		Output:     arg.Output,
	}
	return s.step, nil
}

func (s *stubQuerier) CreateHistory(_ context.Context, arg store.CreateHistoryParams) (store.History, error) {
	s.history = store.History{
		ID:           arg.ID,
		WorkspaceID:  arg.WorkspaceID,
		JobReference: arg.JobReference,
		Output:       arg.Output,
		Serial:       arg.Serial,
		Md5:          arg.Md5,
		Lineage:      arg.Lineage,
		CreatedBy:    arg.CreatedBy,
	}
	return s.history, nil
}

func (s *stubQuerier) GetHistory(context.Context, string) (store.History, error) {
	if s.historyErr != nil {
		return store.History{}, s.historyErr
	}
	if s.history.ID == "" {
		return store.History{}, pgx.ErrNoRows
	}
	return s.history, nil
}

func (s *stubQuerier) CreateArchive(_ context.Context, arg store.CreateArchiveParams) (store.TempArchive, error) {
	s.createdArchive = &arg
	archive := store.TempArchive{ID: arg.ID, Type: arg.Type, HistoryID: arg.HistoryID}
	if s.archives == nil {
		s.archives = map[string][]store.TempArchive{}
	}
	s.archives[arg.HistoryID] = append(s.archives[arg.HistoryID], archive)
	return archive, nil
}

func (s *stubQuerier) ListArchivesByHistory(_ context.Context, historyID string) ([]store.TempArchive, error) {
	return s.archives[historyID], nil
}

func (s *stubQuerier) CreateContent(_ context.Context, arg store.CreateContentParams) (store.Content, error) {
	s.createdContent = &arg
	s.content = store.Content{
		ID:            arg.ID,
		AutoQueueRuns: arg.AutoQueueRuns,
		Source:        arg.Source,
		Speculative:   arg.Speculative,
		Status:        arg.Status,
		WorkspaceID:   arg.WorkspaceID,
	}
	return s.content, nil
}

func (s *stubQuerier) GetContent(context.Context, string) (store.Content, error) {
	if s.contentErr != nil {
		return store.Content{}, s.contentErr
	}
	if s.content.ID == "" {
		return store.Content{}, pgx.ErrNoRows
	}
	return s.content, nil
}

func (s *stubQuerier) UpdateContentStatus(_ context.Context, arg store.UpdateContentStatusParams) (store.Content, error) {
	s.updatedContentStatus = &arg
	s.content.Status = arg.Status
	return s.content, nil
}

func (s *stubQuerier) ListProjectsByOrganization(context.Context, string) ([]store.Project, error) {
	return s.projects, s.projectsErr
}

func (s *stubQuerier) GetProject(context.Context, string) (store.Project, error) {
	return s.project, s.projectErr
}

func (s *stubQuerier) CreateProject(_ context.Context, arg store.CreateProjectParams) (store.Project, error) {
	s.createdProj = &arg
	return s.project, nil
}

func newTestAPI(t *testing.T, q store.Querier) http.Handler {
	t.Helper()
	handler, _ := newTestAPIWithStorage(t, q)
	return handler
}

func newTestAPIWithStorage(t *testing.T, q store.Querier) (http.Handler, storage.Backend) {
	t.Helper()
	return newTestAPIWithQueue(t, q, &stubJobQueue{})
}

func newTestAPIWithQueue(t *testing.T, q store.Querier, jq JobQueue) (http.Handler, storage.Backend) {
	t.Helper()
	backend, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("create local storage: %v", err)
	}
	return New(Deps{
		Store:              q,
		Storage:            backend,
		Queue:              jq,
		Layout:             tfstate.NewLayout("local"),
		PublicURL:          "https://terrakube.local",
		UploadTicketSecret: "test-secret",
		Logger:             slog.New(slog.NewTextHandler(io.Discard, nil)),
	}), backend
}

func TestDiscovery(t *testing.T) {
	handler := Discovery(Deps{
		PublicURL:    "http://localhost:8080",
		DexClientID:  "terrakube",
		DexIssuerURL: "http://localhost:5556/dex",
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/terraform.json", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var doc map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc["state.v2"] != "/remote/state/v2/" || doc["tfe.v2"] != "/remote/tfe/v2/" {
		t.Errorf("discovery paths = %#v", doc)
	}
	login, ok := doc["login.v1"].(map[string]any)
	if !ok {
		t.Fatalf("login.v1 = %#v", doc["login.v1"])
	}
	if login["client"] != "terrakube" {
		t.Errorf("login client = %v", login["client"])
	}
	if login["token"] != "http://localhost:5556/dex/token" {
		t.Errorf("login token = %v", login["token"])
	}
}

func TestPing(t *testing.T) {
	handler := newTestAPI(t, &stubQuerier{})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if got := rec.Header().Get("TFP-API-Version"); got != "2.5" {
		t.Errorf("TFP-API-Version = %q, want 2.5", got)
	}
	if got := rec.Header().Get("TFP-AppName"); got != "Terrakube" {
		t.Errorf("TFP-AppName = %q, want Terrakube", got)
	}
}

func TestAuth(t *testing.T) {
	q := &stubQuerier{org: store.Organization{ID: "org-1", Name: "acme", ExecutionMode: "remote"}}
	handler := New(Deps{
		Store:     q,
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken: "secret",
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/organizations/acme", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("with token status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
}

func TestGetOrganization(t *testing.T) {
	q := &stubQuerier{org: store.Organization{ID: "org-1", Name: "acme", ExecutionMode: "remote"}}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	var doc struct {
		Data struct {
			Type       string         `json:"type"`
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Type != "organizations" || doc.Data.ID != "acme" {
		t.Errorf("data = %+v", doc.Data)
	}
	if doc.Data.Attributes["name"] != "acme" {
		t.Errorf("name = %v", doc.Data.Attributes["name"])
	}
}

func TestEntitlementSet(t *testing.T) {
	q := &stubQuerier{org: store.Organization{ID: "org-1", Name: "acme"}}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/entitlement-set", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	var doc struct {
		Data struct {
			Type       string         `json:"type"`
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Type != "entitlement-sets" || doc.Data.ID != "org-acme" {
		t.Errorf("data = %+v", doc.Data)
	}
	if doc.Data.Attributes["state-storage"] != true {
		t.Errorf("state-storage = %v", doc.Data.Attributes["state-storage"])
	}
}

func TestListWorkspacesRequiresFilter(t *testing.T) {
	q := &stubQuerier{
		org: store.Organization{ID: "org-1", Name: "acme"},
		workspaces: []store.Workspace{
			{ID: "ws-1", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
		},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/workspaces", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var empty struct {
		Data []any          `json:"data"`
		Meta map[string]any `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &empty); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(empty.Data) != 0 {
		t.Errorf("data = %#v, want empty", empty.Data)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/workspaces?search[name]=dev", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("filtered status = %d, want 200", rec.Code)
	}
	var filtered struct {
		Data []struct {
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &filtered); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(filtered.Data) != 1 || filtered.Data[0].ID != "ws-1" {
		t.Fatalf("filtered data = %+v", filtered.Data)
	}
	if filtered.Data[0].Attributes["name"] != "dev" {
		t.Errorf("name = %v", filtered.Data[0].Attributes["name"])
	}
}

func TestListWorkspacesTagFilter(t *testing.T) {
	q := &stubQuerier{
		org: store.Organization{ID: "org-1", Name: "acme"},
		workspaces: []store.Workspace{
			{ID: "ws-1", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
			{ID: "ws-2", Name: "prod", OrganizationID: "org-1", GlobalRemoteState: true},
		},
		tags: []store.Tag{{ID: "tag-1", Name: "env", OrganizationID: "org-1"}},
		workspaceTagsByWorkspace: map[string][]store.Workspacetag{
			"ws-1": {{ID: "wt-1", TagID: "tag-1", WorkspaceID: "ws-1", TagValue: stringPtr("dev")}},
			"ws-2": {{ID: "wt-2", TagID: "tag-1", WorkspaceID: "ws-2", TagValue: stringPtr("prod")}},
		},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/workspaces?search%5Btags%5D=env&filter%5Btagged%5D%5B0%5D%5Bkey%5D=env&filter%5Btagged%5D%5B0%5D%5Bvalue%5D=dev", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data []struct {
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 1 || doc.Data[0].ID != "ws-1" {
		t.Fatalf("data = %+v, want ws-1", doc.Data)
	}
	names, ok := doc.Data[0].Attributes["tag-names"].([]any)
	if !ok || len(names) != 1 || names[0] != "env" {
		t.Errorf("tag-names = %#v", doc.Data[0].Attributes["tag-names"])
	}
}

func TestCreateWorkspace(t *testing.T) {
	q := &stubQuerier{
		org:                store.Organization{ID: "org-1", Name: "acme", ExecutionMode: "remote"},
		workspace:          store.Workspace{ID: "ws-new", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
		workspaceByNameErr: pgx.ErrNoRows,
	}
	handler := newTestAPI(t, q)

	body := `{"data":{"type":"workspaces","attributes":{"name":"dev","terraform-version":"1.9.0","execution-mode":"remote"}}}`
	req := httptest.NewRequest(http.MethodPost, "/organizations/acme/workspaces", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if q.createdWS == nil {
		t.Fatal("CreateWorkspace was not called")
	}
	if q.createdWS.Name != "dev" || q.createdWS.OrganizationID != "org-1" {
		t.Errorf("created = %+v", q.createdWS)
	}
	if q.createdWS.TerraformVersion == nil || *q.createdWS.TerraformVersion != "1.9.0" {
		t.Errorf("terraform version = %v", q.createdWS.TerraformVersion)
	}
	if q.createdWS.Branch == nil || *q.createdWS.Branch != "remote-content" {
		t.Errorf("branch = %v", q.createdWS.Branch)
	}
}

func TestLockWorkspace(t *testing.T) {
	q := &stubQuerier{
		org:       store.Organization{ID: "org-1", Name: "acme"},
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/actions/lock", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if q.lockParams == nil || !q.lockParams.Locked {
		t.Fatalf("lock params = %+v", q.lockParams)
	}

	q.workspace.Locked = true
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/actions/lock", nil))
	if rec.Code != http.StatusConflict {
		t.Fatalf("already locked status = %d, want 409", rec.Code)
	}
}

func TestLockPreservesDescriptionAndUnlockClears(t *testing.T) {
	description := "PR apply lock"
	q := &stubQuerier{
		workspace: store.Workspace{
			ID:              "ws-1",
			Name:            "dev",
			OrganizationID:  "org-1",
			LockDescription: &description,
		},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/actions/lock", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("lock status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.lockParams == nil || q.lockParams.LockDescription == nil || *q.lockParams.LockDescription != description {
		t.Fatalf("lock description = %+v, want %q", q.lockParams, description)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/actions/unlock", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("unlock status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.lockParams == nil || q.lockParams.LockDescription != nil {
		t.Fatalf("unlock description = %+v, want nil", q.lockParams)
	}
}

func TestProjects(t *testing.T) {
	q := &stubQuerier{
		org:      store.Organization{ID: "org-1", Name: "acme"},
		projects: []store.Project{{ID: "proj-1", Name: "platform", OrganizationID: "org-1"}},
		project:  store.Project{ID: "proj-new", Name: "app", OrganizationID: "org-1"},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/projects", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200", rec.Code)
	}
	var list struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(list.Data) != 1 || list.Data[0].ID != "proj-1" {
		t.Errorf("projects = %+v", list.Data)
	}

	body := `{"data":{"type":"projects","attributes":{"name":"app"}}}`
	req := httptest.NewRequest(http.MethodPost, "/organizations/acme/projects", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if q.createdProj == nil || q.createdProj.Name != "app" {
		t.Errorf("created project = %+v", q.createdProj)
	}
}

func TestTagBindingsRoundTrip(t *testing.T) {
	q := &stubQuerier{
		org:       store.Organization{ID: "org-1", Name: "acme"},
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/tag-bindings", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("initial list status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	body := `{"data":[{"type":"tag-bindings","attributes":{"key":"env","value":"dev"}}]}`
	req := httptest.NewRequest(http.MethodPatch, "/workspaces/ws-1/tag-bindings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if q.createdTag == nil || q.createdTag.Name != "env" {
		t.Fatalf("created tag = %+v", q.createdTag)
	}
	if q.createdWorkspaceTag == nil || q.createdWorkspaceTag.TagValue == nil || *q.createdWorkspaceTag.TagValue != "dev" {
		t.Fatalf("created workspace tag = %+v", q.createdWorkspaceTag)
	}

	body = `{"data":[{"type":"tag-bindings","attributes":{"key":"env","value":"prod"}}]}`
	req = httptest.NewRequest(http.MethodPatch, "/workspaces/ws-1/tag-bindings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if q.updatedWorkspaceTag == nil || q.updatedWorkspaceTag.TagValue == nil || *q.updatedWorkspaceTag.TagValue != "prod" {
		t.Fatalf("updated workspace tag = %+v", q.updatedWorkspaceTag)
	}

	body = `{"data":[{"type":"tags","attributes":{"name":"team"}}]}`
	req = httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/relationships/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("tags status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/tag-bindings", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("final list status = %d, want 200", rec.Code)
	}
	var doc struct {
		Data []struct {
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 2 {
		t.Fatalf("bindings = %+v, want 2", doc.Data)
	}
	values := map[string]string{}
	for _, item := range doc.Data {
		values[item.Attributes["key"].(string)] = item.Attributes["value"].(string)
	}
	if values["env"] != "prod" || values["team"] != "" {
		t.Errorf("bindings = %#v", values)
	}
}

func TestTagBindingsRejectEmptyKey(t *testing.T) {
	q := &stubQuerier{
		org:       store.Organization{ID: "org-1", Name: "acme"},
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1", GlobalRemoteState: true},
	}
	handler := newTestAPI(t, q)

	body := `{"data":[{"type":"tag-bindings","attributes":{"key":"","value":"x"}}]}`
	req := httptest.NewRequest(http.MethodPatch, "/workspaces/ws-1/tag-bindings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body = %s", rec.Code, rec.Body.String())
	}
}

func TestCreateWorkspaceWithTags(t *testing.T) {
	q := &stubQuerier{
		org:                store.Organization{ID: "org-1", Name: "acme", ExecutionMode: "remote"},
		workspace:          store.Workspace{ID: "ws-new", Name: "tagged", OrganizationID: "org-1", GlobalRemoteState: true},
		workspaceByNameErr: pgx.ErrNoRows,
	}
	handler := newTestAPI(t, q)

	body := `{"data":{"type":"workspaces","attributes":{"name":"tagged"},"relationships":{"tags":{"data":[{"type":"tags","attributes":{"name":"team"}}]},"tag-bindings":{"data":[{"type":"tag-bindings","attributes":{"key":"env","value":"dev"}}]}}}}`
	req := httptest.NewRequest(http.MethodPost, "/organizations/acme/workspaces", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-new/tag-bindings", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data []struct {
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 2 {
		t.Fatalf("bindings = %+v, want 2", doc.Data)
	}
}

func TestCreateAndGetStateVersion(t *testing.T) {
	state := `{"version":4,"terraform_version":"1.9.0","serial":7,"lineage":"lineage-7","outputs":{},"resources":[]}`
	jsonState := `{"format_version":"1.0","terraform_version":"1.9.0","values":{"root_module":{}}}`
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1", TerraformVersion: stringPtr("1.9.0")},
	}
	handler, backend := newTestAPIWithStorage(t, q)

	body := fmt.Sprintf(`{"data":{"type":"state-versions","attributes":{"state":%q,"json-state":%q}}}`,
		base64.StdEncoding.EncodeToString([]byte(state)),
		base64.StdEncoding.EncodeToString([]byte(jsonState)))
	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/state-versions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data struct {
			ID         string         `json:"id"`
			Type       string         `json:"type"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Type != "state-versions" || doc.Data.ID == "" {
		t.Fatalf("state version = %+v", doc.Data)
	}
	if doc.Data.Attributes["status"] != "finalized" {
		t.Errorf("status = %v, want finalized", doc.Data.Attributes["status"])
	}
	if q.history.Serial != 7 || q.history.Lineage != "lineage-7" {
		t.Errorf("history = %+v", q.history)
	}
	if q.job.Status != "completed" || q.step.Name == nil || *q.step.Name != "Dummy State Uploaded" {
		t.Errorf("dummy job/step = %+v / %+v", q.job, q.step)
	}

	raw, err := backend.Get(context.Background(), tfstate.NewLayout("local").RawState("org-1", "ws-1", doc.Data.ID))
	if err != nil {
		t.Fatalf("read raw state: %v", err)
	}
	rawBody, _ := io.ReadAll(raw)
	raw.Close()
	if string(rawBody) != state {
		t.Errorf("raw state = %s", rawBody)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/state-versions/"+doc.Data.ID, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"hosted-state-download-url"`) {
		t.Errorf("get body = %s", rec.Body.String())
	}
}

func TestCreateStateVersionPendingUpload(t *testing.T) {
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
	}
	handler := newTestAPI(t, q)

	body := `{"data":{"type":"state-versions","attributes":{}}}`
	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/state-versions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data struct {
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Attributes["status"] != "pending" {
		t.Errorf("status = %v, want pending", doc.Data.Attributes["status"])
	}
	if doc.Data.Attributes["hosted-state-upload-url"] == nil || doc.Data.Attributes["hosted-json-state-upload-url"] == nil {
		t.Errorf("upload URLs = %#v", doc.Data.Attributes)
	}
	if q.createdArchive == nil {
		t.Fatal("no archive created")
	}
}

func TestCreateStateVersionUpdatesCurrentState(t *testing.T) {
	state := `{"version":4,"terraform_version":"1.9.0","serial":7,"lineage":"lineage-7","outputs":{"name":{"value":"demo","type":"string","sensitive":false}}}`
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
	}
	handler, _ := newTestAPIWithStorage(t, q)

	body := fmt.Sprintf(`{"data":{"type":"state-versions","attributes":{"state":"%s"}}}`, base64.StdEncoding.EncodeToString([]byte(state)))
	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/state-versions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/current-state-version", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("current state status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var current struct {
		Data struct {
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &current); err != nil {
		t.Fatalf("decode current state: %v", err)
	}
	if current.Data.Attributes["serial"] != float64(7) {
		t.Errorf("serial = %v, want 7", current.Data.Attributes["serial"])
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/current-state-version-outputs", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"name":"name"`) || !strings.Contains(rec.Body.String(), `"value":"demo"`) {
		t.Fatalf("outputs status = %d; body = %s", rec.Code, rec.Body.String())
	}
}

func TestCurrentStateVersion(t *testing.T) {
	state := `{"version":4,"terraform_version":"1.10.1","serial":12,"lineage":"lineage-12"}`
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1", TerraformVersion: stringPtr("1.9.0")},
	}
	handler, backend := newTestAPIWithStorage(t, q)
	layout := tfstate.NewLayout("local")
	if err := backend.Put(context.Background(), layout.CurrentState("org-1", "ws-1"), strings.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("put current state: %v", err)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/current-state-version", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data struct {
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Attributes["serial"] != float64(12) {
		t.Errorf("serial = %v, want 12", doc.Data.Attributes["serial"])
	}
	if doc.Data.Attributes["terraform-version"] != "1.10.1" {
		t.Errorf("terraform-version = %v", doc.Data.Attributes["terraform-version"])
	}
	if doc.Data.Attributes["status"] != "finalized" {
		t.Errorf("status = %v", doc.Data.Attributes["status"])
	}
}

func TestCurrentStateVersionMissing(t *testing.T) {
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/current-state-version", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestCurrentStateVersionOutputs(t *testing.T) {
	state := `{"version":4,"outputs":{"z_count":{"value":3,"type":"number","sensitive":false},"a_name":{"value":"demo","type":"string","sensitive":true}}}`
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
	}
	handler, backend := newTestAPIWithStorage(t, q)
	layout := tfstate.NewLayout("local")
	if err := backend.Put(context.Background(), layout.CurrentState("org-1", "ws-1"), strings.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("put current state: %v", err)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/current-state-version-outputs", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data []struct {
			Type       string         `json:"type"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 2 {
		t.Fatalf("outputs = %d, want 2; body = %s", len(doc.Data), rec.Body.String())
	}
	if doc.Data[0].Type != "state-version-outputs" || doc.Data[0].Attributes["name"] != "a_name" {
		t.Errorf("first output = %+v", doc.Data[0])
	}
	if doc.Data[0].Attributes["value"] != "demo" || doc.Data[0].Attributes["sensitive"] != true {
		t.Errorf("string output = %+v", doc.Data[0].Attributes)
	}
	if doc.Data[0].Attributes["type"] != "string" || doc.Data[0].Attributes["detailed-type"] != "string" {
		t.Errorf("string types = %+v", doc.Data[0].Attributes)
	}
	if doc.Data[1].Attributes["type"] != "string" || doc.Data[1].Attributes["detailed-type"] != "string" {
		t.Errorf("number types = %+v", doc.Data[1].Attributes)
	}
}

func TestCurrentStateVersionOutputsMissing(t *testing.T) {
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/current-state-version-outputs", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"data":[]`) {
		t.Errorf("body = %s, want empty data", rec.Body.String())
	}
}

func TestCreateGetUploadConfigurationVersion(t *testing.T) {
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
	}
	handler, backend := newTestAPIWithStorage(t, q)

	body := `{"data":{"type":"configuration-versions","attributes":{"speculative":true,"auto-queue-runs":true}}}`
	req := httptest.NewRequest(http.MethodPost, "/workspaces/ws-1/configuration-versions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data struct {
			ID         string         `json:"id"`
			Type       string         `json:"type"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.Type != "configuration-versions" || doc.Data.ID == "" {
		t.Fatalf("configuration version = %+v", doc.Data)
	}
	if doc.Data.Attributes["status"] != "pending" || doc.Data.Attributes["speculative"] != true || doc.Data.Attributes["auto-queue-runs"] != true {
		t.Errorf("attributes = %#v", doc.Data.Attributes)
	}
	uploadURL, _ := doc.Data.Attributes["upload-url"].(string)
	if uploadURL == "" {
		t.Fatal("missing upload-url")
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/configuration-versions/"+doc.Data.ID, nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"status":"pending"`) {
		t.Fatalf("get status = %d; body = %s", rec.Code, rec.Body.String())
	}

	parsed, err := url.Parse(uploadURL)
	if err != nil {
		t.Fatalf("parse upload url: %v", err)
	}
	uploadPath := strings.TrimPrefix(parsed.Path, "/remote/tfe/v2")
	if parsed.RawQuery != "" {
		uploadPath += "?" + parsed.RawQuery
	}
	req = httptest.NewRequest(http.MethodPut, uploadPath, strings.NewReader("configuration-tarball"))
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("upload status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if q.updatedContentStatus == nil || q.updatedContentStatus.Status == nil || *q.updatedContentStatus.Status != "uploaded" {
		t.Fatalf("updated status = %+v", q.updatedContentStatus)
	}

	rc, err := backend.Get(context.Background(), tfstate.NewLayout("local").Content(doc.Data.ID))
	if err != nil {
		t.Fatalf("read stored content: %v", err)
	}
	stored, _ := io.ReadAll(rc)
	rc.Close()
	if string(stored) != "configuration-tarball" {
		t.Errorf("stored content = %q", stored)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/configuration-versions/"+doc.Data.ID+"/terraformContent.tar.gz", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "configuration-tarball" {
		t.Fatalf("download status = %d; body = %q", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPut, uploadPath, strings.NewReader("again"))
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("second upload status = %d, want 409; body = %s", rec.Code, rec.Body.String())
	}
}

func TestConfigurationUploadTickets(t *testing.T) {
	handler := newTestAPI(t, &stubQuerier{})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/configuration-versions/c1", strings.NewReader("x")))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("missing ticket status = %d, want 401", rec.Code)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/configuration-versions/c1?ticket=bad", strings.NewReader("x")))
	if rec.Code != http.StatusForbidden {
		t.Fatalf("invalid ticket status = %d, want 403", rec.Code)
	}

	api := &API{deps: Deps{UploadTicketSecret: "test-secret"}}
	expired := api.signUploadTicket("c1", time.Now().Add(-time.Minute))
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/configuration-versions/c1?ticket="+url.QueryEscape(expired), strings.NewReader("x")))
	if rec.Code != http.StatusGone {
		t.Fatalf("expired ticket status = %d, want 410", rec.Code)
	}
}

func TestGetConfigurationVersionNotFound(t *testing.T) {
	handler := newTestAPI(t, &stubQuerier{})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/configuration-versions/missing", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestCreateGetAndDiscardRun(t *testing.T) {
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
		content: store.Content{
			ID:          "cv-1",
			WorkspaceID: stringPtr("ws-1"),
			Status:      stringPtr("uploaded"),
		},
	}
	jq := &stubJobQueue{}
	handler, _ := newTestAPIWithQueue(t, q, jq)

	body := `{"data":{"type":"runs","attributes":{"is-destroy":false,"auto-apply":false,"refresh":true},"relationships":{"workspace":{"data":{"type":"workspaces","id":"ws-1"}},"configuration-version":{"data":{"type":"configuration-versions","id":"cv-1"}}}}}`
	req := httptest.NewRequest(http.MethodPost, "/runs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data struct {
			ID            string `json:"id"`
			Type          string `json:"type"`
			Attributes    map[string]any
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
	if doc.Data.Type != "runs" || doc.Data.ID != "run-42" {
		t.Fatalf("run = %+v", doc.Data)
	}
	if doc.Data.Attributes["status"] != "pending" {
		t.Errorf("status = %v", doc.Data.Attributes["status"])
	}
	if doc.Data.Relationships["workspace"].Data.ID != "ws-1" {
		t.Errorf("workspace relationship = %+v", doc.Data.Relationships["workspace"])
	}
	if q.job.Command == nil || *q.job.Command != "plan" {
		t.Errorf("command = %v", q.job.Command)
	}
	if len(jq.enqueued) != 1 || jq.enqueued[0] != 42 {
		t.Errorf("enqueued = %v, want [42]", jq.enqueued)
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/runs/run-42", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"run-42"`) {
		t.Fatalf("get status = %d; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/runs/run-42/actions/discard", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("discard status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.updatedJob == nil || q.updatedJob.Status != "cancelled" {
		t.Fatalf("updated job = %+v", q.updatedJob)
	}
}

func TestCreateRunEnqueueFailure(t *testing.T) {
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
		content: store.Content{
			ID:          "cv-1",
			WorkspaceID: stringPtr("ws-1"),
			Status:      stringPtr("uploaded"),
		},
	}
	jq := &stubJobQueue{err: errors.New("redis down")}
	handler, _ := newTestAPIWithQueue(t, q, jq)

	body := `{"data":{"type":"runs","relationships":{"workspace":{"data":{"type":"workspaces","id":"ws-1"}},"configuration-version":{"data":{"type":"configuration-versions","id":"cv-1"}}}}}`
	req := httptest.NewRequest(http.MethodPost, "/runs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/vnd.api+json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body = %s", rec.Code, rec.Body.String())
	}
	if q.job.Status != "failed" {
		t.Errorf("job status = %q, want failed", q.job.Status)
	}
}

func TestListWorkspaceRunsAndQueue(t *testing.T) {
	q := &stubQuerier{
		org:       store.Organization{ID: "org-1", Name: "acme"},
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
		jobs: []store.Job{
			{ID: 42, OrganizationID: "org-1", WorkspaceID: "ws-1", Status: "pending"},
		},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspaces/ws-1/runs", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"run-42"`) {
		t.Fatalf("workspace runs status = %d; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organizations/acme/runs/queue", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"run-42"`) {
		t.Fatalf("queue status = %d; body = %s", rec.Code, rec.Body.String())
	}
}

func TestRunEventsPlanAndApply(t *testing.T) {
	q := &stubQuerier{
		workspace: store.Workspace{ID: "ws-1", Name: "dev", OrganizationID: "org-1"},
		job:       store.Job{ID: 42, OrganizationID: "org-1", WorkspaceID: "ws-1", Status: "pending", PlanChanges: true},
	}
	handler := newTestAPI(t, q)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/runs/run-42/run-events", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"data": []`) {
		t.Fatalf("run events status = %d; body = %q", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/plans/42", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"type":"plans"`) {
		t.Fatalf("plan status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "/remote/tfe/v2/plans/logs/42/0") {
		t.Errorf("missing log-read-url: %s", rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/applies/42", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"type":"applies"`) {
		t.Fatalf("apply status = %d; body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/plans/logs/42/0?offset=0&limit=10", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("plan logs status = %d", rec.Code)
	}
}

func stringPtr(value string) *string {
	return &value
}
