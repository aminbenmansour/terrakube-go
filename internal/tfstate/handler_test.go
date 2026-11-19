package tfstate

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"

	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	testOrg = "11111111-1111-1111-1111-111111111111"
	testWS1 = "22222222-2222-2222-2222-222222222222"
	testWS2 = "33333333-3333-3333-3333-333333333333"
	testWS3 = "44444444-4444-4444-4444-444444444444"
	testWS4 = "55555555-5555-5555-5555-555555555555"
	testWS5 = "66666666-6666-6666-6666-666666666666"
)

type memObject struct {
	data        []byte
	contentType string
	modifiedAt  time.Time
}

type memStorage struct {
	mu      sync.Mutex
	objects map[string]memObject
	failPut func(key string) error
}

func newMemStorage() *memStorage {
	return &memStorage{objects: make(map[string]memObject)}
}

func (m *memStorage) Put(_ context.Context, key string, r io.Reader, _ int64, contentType string) error {
	if m.failPut != nil {
		if err := m.failPut(key); err != nil {
			return err
		}
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.objects[key] = memObject{data: data, contentType: contentType, modifiedAt: time.Now()}
	return nil
}

func (m *memStorage) Get(_ context.Context, key string) (io.ReadCloser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	obj, ok := m.objects[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(obj.data)), nil
}

func (m *memStorage) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.objects, key)
	return nil
}

func (m *memStorage) Stat(_ context.Context, key string) (storage.Object, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	obj, ok := m.objects[key]
	if !ok {
		return storage.Object{}, storage.ErrNotFound
	}
	return storage.Object{Key: key, Size: int64(len(obj.data)), ContentType: obj.contentType, ModifiedAt: obj.modifiedAt}, nil
}

func (m *memStorage) Health(context.Context) error { return nil }

type stubStore struct {
	store.Querier
	workspaces map[string]store.Workspace
	steps      map[string]store.Step
	jobs       map[int64]store.Job
	archives   map[string]store.GetArchiveRow
	history    map[string]store.History
}

func (s *stubStore) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	ws, ok := s.workspaces[id]
	if !ok {
		return store.Workspace{}, pgx.ErrNoRows
	}
	return ws, nil
}

func (s *stubStore) GetStep(_ context.Context, id string) (store.Step, error) {
	step, ok := s.steps[id]
	if !ok {
		return store.Step{}, pgx.ErrNoRows
	}
	return step, nil
}

func (s *stubStore) GetJob(_ context.Context, id int64) (store.Job, error) {
	job, ok := s.jobs[id]
	if !ok {
		return store.Job{}, pgx.ErrNoRows
	}
	return job, nil
}

func (s *stubStore) GetArchive(_ context.Context, id string) (store.GetArchiveRow, error) {
	archive, ok := s.archives[id]
	if !ok {
		return store.GetArchiveRow{}, pgx.ErrNoRows
	}
	return archive, nil
}

func (s *stubStore) DeleteArchive(_ context.Context, id string) error {
	delete(s.archives, id)
	return nil
}

func (s *stubStore) CreateHistory(_ context.Context, arg store.CreateHistoryParams) (store.History, error) {
	history := store.History{
		ID:           arg.ID,
		WorkspaceID:  arg.WorkspaceID,
		JobReference: arg.JobReference,
		Output:       arg.Output,
		Serial:       arg.Serial,
		Md5:          arg.Md5,
		Lineage:      arg.Lineage,
	}
	s.history[arg.ID] = history
	return history, nil
}

func (s *stubStore) UpdateHistoryOutput(_ context.Context, arg store.UpdateHistoryOutputParams) error {
	history := s.history[arg.ID]
	history.Output = arg.Output
	s.history[arg.ID] = history
	return nil
}

type fixture struct {
	handler http.Handler
	storage *memStorage
	store   *stubStore
	secret  string
	token   string
}

func newFixture(t *testing.T, maxSize int64) *fixture {
	t.Helper()
	secret := base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{0x17}, 64))
	validator, err := internalauth.NewValidator(secret)
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	statusRunning := "running"
	statusDone := "completed"
	st := &stubStore{
		workspaces: map[string]store.Workspace{
			testWS1: {ID: testWS1, OrganizationID: testOrg},
			testWS2: {ID: testWS2, OrganizationID: testOrg, GlobalRemoteState: true},
			testWS3: {ID: testWS3, OrganizationID: testOrg, SharedIds: strPtr(testWS1)},
			testWS4: {ID: testWS4, OrganizationID: testOrg, SharedIds: strPtr("someone-else")},
			testWS5: {ID: testWS5, OrganizationID: "other-org"},
		},
		jobs: map[int64]store.Job{
			42: {ID: 42, OrganizationID: testOrg, WorkspaceID: testWS1},
		},
		steps: map[string]store.Step{
			"step-1":    {ID: "step-1", JobID: 42, Status: &statusRunning},
			"step-done": {ID: "step-done", JobID: 42, Status: &statusDone},
		},
		archives: make(map[string]store.GetArchiveRow),
		history:  make(map[string]store.History),
	}
	token, err := internalauth.Sign(secret, internalauth.Claims{
		OrganizationID: testOrg,
		WorkspaceID:    testWS1,
		JobID:          "42",
		StepID:         "step-1",
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	ms := newMemStorage()
	handler := New(Deps{
		Store:             st,
		Storage:           ms,
		Logger:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		Layout:            NewLayout("local"),
		Auth:              validator,
		MaxStateSizeBytes: maxSize,
		PublicURL:         "https://terrakube.test",
	})
	return &fixture{handler: handler, storage: ms, store: st, secret: secret, token: token}
}

func strPtr(value string) *string { return &value }

func (f *fixture) request(t *testing.T, method, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, body)
	if headers == nil {
		headers = map[string]string{}
	}
	if _, ok := headers["Authorization"]; !ok {
		headers["Authorization"] = "Bearer " + f.token
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	f.handler.ServeHTTP(rec, req)
	return rec
}

func (f *fixture) requestAs(t *testing.T, token, method, path string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()
	return f.request(t, method, path, body, map[string]string{"Authorization": "Bearer " + token})
}

func backendPath(ws string) string {
	return "/organization/" + testOrg + "/workspace/" + ws + "/backend/terraform.tfstate"
}

func statePath(ws string) string {
	return "/organization/" + testOrg + "/workspace/" + ws + "/state/terraform.tfstate"
}

func TestBackendStateLifecycle(t *testing.T) {
	f := newFixture(t, 1024)
	state := []byte(`{"version":4,"serial":1}`)

	if rec := f.request(t, http.MethodGet, backendPath(testWS1), nil, map[string]string{"Authorization": ""}); rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET without auth = %d, want 401", rec.Code)
	}
	if rec := f.request(t, http.MethodGet, backendPath(testWS1), nil, nil); rec.Code != http.StatusNoContent {
		t.Fatalf("GET missing = %d, want 204", rec.Code)
	}
	if rec := f.request(t, http.MethodPut, backendPath(testWS1), bytes.NewReader(state), nil); rec.Code != http.StatusOK {
		t.Fatalf("PUT = %d, want 200", rec.Code)
	}
	rec := f.request(t, http.MethodGet, backendPath(testWS1), nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET = %d, want 200", rec.Code)
	}
	if got := rec.Body.String(); got != string(state) {
		t.Errorf("GET body = %q, want %q", got, state)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("GET Content-Type = %q", got)
	}

	sum := md5.Sum(state)
	md5Header := base64.StdEncoding.EncodeToString(sum[:])
	if rec := f.request(t, http.MethodPut, backendPath(testWS1), bytes.NewReader(state), map[string]string{"Content-MD5": md5Header}); rec.Code != http.StatusOK {
		t.Fatalf("PUT with matching MD5 = %d, want 200", rec.Code)
	}
	if rec := f.request(t, http.MethodPut, backendPath(testWS1), bytes.NewReader(state), map[string]string{"Content-MD5": "not-the-md5"}); rec.Code != http.StatusBadRequest {
		t.Fatalf("PUT with bad MD5 = %d, want 400", rec.Code)
	}
	if rec := f.request(t, http.MethodPut, backendPath(testWS1), bytes.NewReader(nil), nil); rec.Code != http.StatusBadRequest {
		t.Fatalf("PUT empty = %d, want 400", rec.Code)
	}
	if rec := f.request(t, http.MethodPut, backendPath(testWS1), bytes.NewReader(bytes.Repeat([]byte{'x'}, 1025)), nil); rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("PUT too large = %d, want 413", rec.Code)
	}

	if rec := f.request(t, http.MethodDelete, backendPath(testWS1), nil, nil); rec.Code != http.StatusOK {
		t.Fatalf("DELETE = %d, want 200", rec.Code)
	}
	if rec := f.request(t, http.MethodGet, backendPath(testWS1), nil, nil); rec.Code != http.StatusNoContent {
		t.Fatalf("GET after DELETE = %d, want 204", rec.Code)
	}
}

func TestBackendScopeAndStepValidation(t *testing.T) {
	f := newFixture(t, 1024)
	otherWSToken, err := internalauth.Sign(f.secret, internalauth.Claims{
		OrganizationID: testOrg,
		WorkspaceID:    testWS2,
		JobID:          "42",
		StepID:         "step-1",
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	if rec := f.requestAs(t, otherWSToken, http.MethodGet, backendPath(testWS1), nil); rec.Code != http.StatusForbidden {
		t.Fatalf("wrong workspace token = %d, want 403", rec.Code)
	}

	doneToken, _ := internalauth.Sign(f.secret, internalauth.Claims{
		OrganizationID: testOrg, WorkspaceID: testWS1, JobID: "42", StepID: "step-done",
	})
	if rec := f.requestAs(t, doneToken, http.MethodGet, backendPath(testWS1), nil); rec.Code != http.StatusUnauthorized {
		t.Fatalf("completed step token = %d, want 401", rec.Code)
	}

	missingStepToken, _ := internalauth.Sign(f.secret, internalauth.Claims{
		OrganizationID: testOrg, WorkspaceID: testWS1, JobID: "42", StepID: "missing",
	})
	if rec := f.requestAs(t, missingStepToken, http.MethodGet, backendPath(testWS1), nil); rec.Code != http.StatusUnauthorized {
		t.Fatalf("missing step token = %d, want 401", rec.Code)
	}

	serviceToken, _ := internalauth.Sign(f.secret, internalauth.Claims{
		Subject: internalauth.SubjectExecutor,
	})
	if rec := f.requestAs(t, serviceToken, http.MethodGet, backendPath(testWS1), nil); rec.Code != http.StatusForbidden {
		t.Fatalf("executor service token on backend = %d, want 403", rec.Code)
	}

	foreignToken, _ := internalauth.Sign(f.secret, internalauth.Claims{
		RegisteredClaims: jwtRegisteredClaims("not-terrakube"),
		WorkspaceID:      testWS1,
	})
	if rec := f.requestAs(t, foreignToken, http.MethodGet, backendPath(testWS1), nil); rec.Code != http.StatusForbidden {
		t.Fatalf("foreign issuer token = %d, want 403", rec.Code)
	}

	if rec := f.request(t, http.MethodGet, backendPath(testWS1), nil, map[string]string{"Authorization": "Bearer garbage"}); rec.Code != http.StatusUnauthorized {
		t.Fatalf("malformed token = %d, want 401", rec.Code)
	}
}

func TestStateReads(t *testing.T) {
	f := newFixture(t, 1024)

	rec := f.request(t, http.MethodGet, statePath(testWS1), nil, nil)
	if rec.Code != http.StatusOK || rec.Body.Len() != 0 {
		t.Fatalf("missing current state = %d len=%d, want 200 empty", rec.Code, rec.Body.Len())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("current state Content-Type = %q", got)
	}

	current := []byte(`{"version":4}`)
	if err := f.storage.Put(context.Background(), "backend/"+testOrg+"/"+testWS1+"/terraform.tfstate", bytes.NewReader(current), int64(len(current)), "application/json"); err != nil {
		t.Fatalf("Put current state: %v", err)
	}
	rec = f.request(t, http.MethodGet, statePath(testWS1), nil, nil)
	if rec.Code != http.StatusOK || rec.Body.String() != string(current) {
		t.Fatalf("current state = %d body=%q, want 200 %q", rec.Code, rec.Body.String(), current)
	}

	archived := []byte(`{"version":4,"serial":7}`)
	if err := f.storage.Put(context.Background(), "state/"+testOrg+"/"+testWS1+"/state/history-1.json", bytes.NewReader(archived), int64(len(archived)), "application/json"); err != nil {
		t.Fatalf("Put archived state: %v", err)
	}
	rec = f.request(t, http.MethodGet, "/organization/"+testOrg+"/workspace/"+testWS1+"/state/history-1.json", nil, nil)
	if rec.Code != http.StatusOK || rec.Body.String() != string(archived) {
		t.Fatalf("archived state = %d body=%q", rec.Code, rec.Body.String())
	}
	if rec := f.request(t, http.MethodGet, "/organization/"+testOrg+"/workspace/"+testWS1+"/state/history-1", nil, nil); rec.Code != http.StatusNotFound {
		t.Fatalf("archived state without .json = %d, want 404", rec.Code)
	}

	plan := []byte("plan-bytes")
	planKey := "state/" + testOrg + "/" + testWS1 + "/42/step-1/terraformLibrary.tfPlan"
	if err := f.storage.Put(context.Background(), planKey, bytes.NewReader(plan), int64(len(plan)), "application/octet-stream"); err != nil {
		t.Fatalf("Put plan: %v", err)
	}
	planPath := "/organization/" + testOrg + "/workspace/" + testWS1 + "/jobId/42/step/step-1/terraform.tfstate"
	rec = f.request(t, http.MethodGet, planPath, nil, nil)
	if rec.Code != http.StatusOK || rec.Body.String() != string(plan) {
		t.Fatalf("plan = %d body=%q", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/octet-stream" {
		t.Errorf("plan Content-Type = %q", got)
	}

	rec = f.request(t, http.MethodGet, "/organization/"+testOrg+"/workspace/"+testWS1+"/jobId/99/step/missing/terraform.tfstate", nil, nil)
	if rec.Code != http.StatusOK || rec.Body.Len() != 0 {
		t.Fatalf("missing plan = %d len=%d, want 200 empty", rec.Code, rec.Body.Len())
	}
}

func TestCrossWorkspaceReads(t *testing.T) {
	f := newFixture(t, 1024)

	if rec := f.request(t, http.MethodGet, statePath(testWS2), nil, nil); rec.Code != http.StatusOK {
		t.Fatalf("global remote state read = %d, want 200", rec.Code)
	}
	if rec := f.request(t, http.MethodGet, statePath(testWS3), nil, nil); rec.Code != http.StatusOK {
		t.Fatalf("shared remote state read = %d, want 200", rec.Code)
	}
	if rec := f.request(t, http.MethodGet, statePath(testWS4), nil, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("unshared remote state read = %d, want 403", rec.Code)
	}
	if rec := f.request(t, http.MethodGet, statePath(testWS5), nil, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("cross-org remote state read = %d, want 403", rec.Code)
	}
	planPath := "/organization/" + testOrg + "/workspace/" + testWS2 + "/jobId/42/step/step-1/terraform.tfstate"
	if rec := f.request(t, http.MethodGet, planPath, nil, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("cross-workspace plan read = %d, want 403", rec.Code)
	}

	serviceToken, _ := internalauth.Sign(f.secret, internalauth.Claims{Subject: internalauth.SubjectExecutor})
	if rec := f.requestAs(t, serviceToken, http.MethodGet, statePath(testWS4), nil); rec.Code != http.StatusOK {
		t.Fatalf("executor service state read = %d, want 200", rec.Code)
	}
	if rec := f.requestAs(t, serviceToken, http.MethodGet, planPath, nil); rec.Code != http.StatusOK {
		t.Fatalf("executor service plan read = %d, want 200", rec.Code)
	}
}

func TestArchiveUploads(t *testing.T) {
	f := newFixture(t, 1024)
	f.store.archives["archive-raw"] = store.GetArchiveRow{
		ID: "archive-raw", HistoryID: "history-raw", WorkspaceID: testWS1, OrganizationID: testOrg,
	}
	raw := []byte(`{"version":4,"serial":1}`)
	rec := f.request(t, http.MethodPut, "/archive/archive-raw/terraform.tfstate", bytes.NewReader(raw), map[string]string{"Authorization": ""})
	if rec.Code != http.StatusCreated {
		t.Fatalf("raw archive upload = %d, want 201", rec.Code)
	}
	assertObject(t, f.storage, "backend/"+testOrg+"/"+testWS1+"/terraform.tfstate", raw)
	assertObject(t, f.storage, "state/"+testOrg+"/"+testWS1+"/state/history-raw.raw.json", raw)
	if _, ok := f.store.archives["archive-raw"]; ok {
		t.Fatal("archive row was not deleted")
	}

	f.store.archives["archive-json"] = store.GetArchiveRow{
		ID: "archive-json", HistoryID: "history-json", WorkspaceID: testWS1, OrganizationID: testOrg,
	}
	jsonState := []byte(`{"version":4,"serial":2}`)
	rec = f.request(t, http.MethodPut, "/archive/archive-json/terraform.json.tfstate", bytes.NewReader(jsonState), map[string]string{"Authorization": ""})
	if rec.Code != http.StatusCreated {
		t.Fatalf("JSON archive upload = %d, want 201", rec.Code)
	}
	assertObject(t, f.storage, "state/"+testOrg+"/"+testWS1+"/state/history-json.json", jsonState)

	rec = f.request(t, http.MethodPut, "/archive/missing/terraform.tfstate", bytes.NewReader(raw), map[string]string{"Authorization": ""})
	if rec.Code != http.StatusForbidden {
		t.Fatalf("unknown archive = %d, want 403", rec.Code)
	}
}

func TestRollback(t *testing.T) {
	f := newFixture(t, 1024)
	jsonState := []byte(`{"version":4,"serial":2}`)
	rawState := []byte(`{"version":4,"serial":2,"raw":true}`)
	if err := f.storage.Put(context.Background(), "state/"+testOrg+"/"+testWS1+"/state/history-1.json", bytes.NewReader(jsonState), int64(len(jsonState)), "application/json"); err != nil {
		t.Fatalf("Put JSON state: %v", err)
	}
	if err := f.storage.Put(context.Background(), "state/"+testOrg+"/"+testWS1+"/state/history-1.raw.json", bytes.NewReader(rawState), int64(len(rawState)), "application/json"); err != nil {
		t.Fatalf("Put raw state: %v", err)
	}

	path := "/organization/" + testOrg + "/workspace/" + testWS1 + "/rollback/history-1.json"
	rec := f.request(t, http.MethodPut, path, nil, nil)
	if rec.Code != http.StatusCreated || rec.Body.String() != "Rollback successful" {
		t.Fatalf("rollback = %d body=%q", rec.Code, rec.Body.String())
	}

	if len(f.store.history) != 1 {
		t.Fatalf("history rows = %d, want 1", len(f.store.history))
	}
	var history store.History
	for _, value := range f.store.history {
		history = value
	}
	if history.JobReference == nil || *history.JobReference != "history-1" {
		t.Errorf("job reference = %v, want history-1", history.JobReference)
	}
	if history.Serial != 1 || history.Md5 != "0" || history.Lineage != "0" {
		t.Errorf("history serial/md5/lineage = %d/%s/%s, want 1/0/0", history.Serial, history.Md5, history.Lineage)
	}
	wantOutput := "https://terrakube.test/tfstate/v1/organization/" + testOrg + "/workspace/" + testWS1 + "/state/" + history.ID + ".json"
	if history.Output == nil || *history.Output != wantOutput {
		t.Errorf("history output = %v, want %s", history.Output, wantOutput)
	}
	assertObject(t, f.storage, "backend/"+testOrg+"/"+testWS1+"/terraform.tfstate", rawState)
	assertObject(t, f.storage, "state/"+testOrg+"/"+testWS1+"/state/"+history.ID+".raw.json", rawState)
	assertObject(t, f.storage, "state/"+testOrg+"/"+testWS1+"/state/"+history.ID+".json", jsonState)

	rec = f.request(t, http.MethodPut, "/organization/"+testOrg+"/workspace/"+testWS1+"/rollback/missing.json", nil, nil)
	if rec.Code != http.StatusNotFound || rec.Body.String() != "JSON state not found" {
		t.Fatalf("missing JSON rollback = %d body=%q", rec.Code, rec.Body.String())
	}
	if err := f.storage.Put(context.Background(), "state/"+testOrg+"/"+testWS1+"/state/no-raw.json", bytes.NewReader(jsonState), int64(len(jsonState)), "application/json"); err != nil {
		t.Fatalf("Put JSON-only state: %v", err)
	}
	rec = f.request(t, http.MethodPut, "/organization/"+testOrg+"/workspace/"+testWS1+"/rollback/no-raw.json", nil, nil)
	if rec.Code != http.StatusNotFound || rec.Body.String() != "Raw Terraform state not found" {
		t.Fatalf("missing raw rollback = %d body=%q", rec.Code, rec.Body.String())
	}

	otherWSToken, _ := internalauth.Sign(f.secret, internalauth.Claims{
		OrganizationID: testOrg, WorkspaceID: testWS2, JobID: "42", StepID: "step-1",
	})
	if rec := f.requestAs(t, otherWSToken, http.MethodPut, path, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("cross-workspace rollback = %d, want 403", rec.Code)
	}
}

func assertObject(t *testing.T, storage *memStorage, key string, want []byte) {
	t.Helper()
	storage.mu.Lock()
	defer storage.mu.Unlock()
	obj, ok := storage.objects[key]
	if !ok {
		t.Fatalf("object %q not found", key)
	}
	if !bytes.Equal(obj.data, want) {
		t.Errorf("object %q = %q, want %q", key, obj.data, want)
	}
}

func jwtRegisteredClaims(issuer string) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{Issuer: issuer}
}
