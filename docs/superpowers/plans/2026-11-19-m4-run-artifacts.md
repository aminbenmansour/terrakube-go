# M4 run logs and artifacts Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Archive final step logs to object storage, serve them through `/tfoutput/v1`, record the state URL on history rows, and fail a run when a log, plan or history upload fails.

**Architecture:** The executor keeps writing inline `step.output`/`job.output` text but also puts the same bounded text under the upstream `tfoutput/{org}/{job}/{step}.tfoutput` key. The log handler prefers that object for finished steps and falls back to the column. History rows gain the upstream state JSON URL, and archive failures propagate to the job status.

**Tech Stack:** Go 1.27, pgx/sqlc, Redis Streams, chi, the existing `storage.Backend` abstraction.

**Spec:** `docs/superpowers/specs/2026-11-19-m4-run-artifacts-design.md`

## Global Constraints

- Upstream is the specification; no HTTP API or schema changes in this plan.
- `internal/store/*.sql.go` is generated; do not edit it. No migration is needed.
- TDD: write the failing test, run it, implement, run it again.
- Run `export PATH=$PATH:/usr/local/go/bin` before Go commands.
- Run `make check` and `go test -race ./...` before the final commit.
- Never commit to `main`; work stays on `feat/m4-run-artifacts`.
- Logs are bounded at `maxLogBytes` (10 MiB); do not remove the bound.

## Review Focus

- A step whose log upload fails but whose Terraform command succeeded must end `failed`, not `completed`.
- A cancelled step must stay `cancelled` even if its log upload fails.
- A missing archived object must fall back to `step.output`, so rows written before this change still load.
- An empty `TERRAKUBE_PUBLIC_URL` must leave `history.output` empty rather than produce a malformed URL.
- A history database/storage failure after a successful apply must fail the run, not warn.

---

### Task 1: Archive the final step log

**Files:**
- Modify: `internal/tfstate/layout.go`
- Modify: `internal/tfstate/layout_test.go`
- Modify: `internal/executor/run.go:164-185` (`completeStep`), call sites at `run.go:71,78,82,138,216,220,223`
- Create: `internal/executor/artifacts_test.go`

**Interfaces:**
- Produces: `func (l Layout) Output(organizationID string, jobID int64, stepID string) string` → `tfoutput/{org}/{job}/{step}.tfoutput`.
- Produces: `func (e *Executor) archiveStepLog(ctx context.Context, job store.Job, step store.Step, output string) error`.
- Changes: `func (e *Executor) completeStep(ctx context.Context, job store.Job, step store.Step, sink *logSink, status string) error` (was no return).

- [ ] **Step 1: Write the failing layout test**

Add to both key maps in `internal/tfstate/layout_test.go`:

```go
layout.Output("org", 7, "step"): "tfoutput/org/7/step.tfoutput",
```

- [ ] **Step 2: Write the failing executor tests**

Create `internal/executor/artifacts_test.go`:

```go
package executor

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

func TestCompleteStepArchivesLog(t *testing.T) {
	ctx := context.Background()
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	st := &artifactStore{}
	e := New(Deps{Store: st, Storage: backend, Layout: layout, Logger: testLogger()})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	sink.WriteString("line one\nline two\n")
	err := e.completeStep(ctx, store.Job{ID: 9, OrganizationID: "org-1"}, store.Step{ID: "step-1"}, sink, "completed")
	if err != nil {
		t.Fatalf("completeStep: %v", err)
	}

	got := string(backend.objects[layout.Output("org-1", 9, "step-1")])
	if got != "line one\nline two\n" {
		t.Fatalf("archived log = %q", got)
	}
	if st.stepStatus != "completed" {
		t.Errorf("step status = %q, want completed", st.stepStatus)
	}
	if st.stepOutput == nil || *st.stepOutput != "line one\nline two\n" {
		t.Errorf("inline step output = %v", st.stepOutput)
	}
}

func TestCompleteStepLogUploadFailureFailsStep(t *testing.T) {
	ctx := context.Background()
	e := New(Deps{
		Store:   &artifactStore{},
		Storage: failingBackend{err: errors.New("storage down")},
		Layout:  tfstate.NewLayout("local"),
		Logger:  testLogger(),
	})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	sink.WriteString("output\n")
	err := e.completeStep(ctx, store.Job{ID: 9}, store.Step{ID: "step-1"}, sink, "completed")
	if err == nil || !strings.Contains(err.Error(), "archive step log") {
		t.Fatalf("error = %v, want archive failure", err)
	}
	if st := e.deps.Store.(*artifactStore); st.stepStatus != "failed" {
		t.Errorf("step status = %q, want failed", st.stepStatus)
	}
}

func TestCompleteStepCancelledKeepsStatusWhenLogUploadFails(t *testing.T) {
	ctx := context.Background()
	st := &artifactStore{}
	e := New(Deps{
		Store:   st,
		Storage: failingBackend{err: errors.New("storage down")},
		Layout:  tfstate.NewLayout("local"),
		Logger:  testLogger(),
	})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	if err := e.completeStep(ctx, store.Job{ID: 9}, store.Step{ID: "step-1"}, sink, "cancelled"); err != nil {
		t.Fatalf("completeStep: %v", err)
	}
	if st.stepStatus != "cancelled" {
		t.Errorf("step status = %q, want cancelled", st.stepStatus)
	}
}

type artifactStore struct {
	store.Querier
	stepStatus string
	stepOutput *string
	jobOutput  *string
	historyErr error
	created    *store.CreateHistoryParams
}

func (s *artifactStore) UpdateStepStatus(_ context.Context, arg store.UpdateStepStatusParams) (store.Step, error) {
	if arg.Status != nil {
		s.stepStatus = *arg.Status
	}
	return store.Step{ID: arg.ID}, nil
}

func (s *artifactStore) UpdateStepOutput(_ context.Context, arg store.UpdateStepOutputParams) error {
	s.stepOutput = arg.Output
	return nil
}

func (s *artifactStore) UpdateJobOutput(_ context.Context, arg store.UpdateJobOutputParams) error {
	s.jobOutput = arg.Output
	return nil
}

func (s *artifactStore) GetHistoryByJobReference(context.Context, store.GetHistoryByJobReferenceParams) (store.History, error) {
	if s.historyErr != nil {
		return store.History{}, s.historyErr
	}
	return store.History{}, pgx.ErrNoRows
}

func (s *artifactStore) CreateHistory(_ context.Context, arg store.CreateHistoryParams) (store.History, error) {
	s.created = &arg
	return store.History{ID: arg.ID}, nil
}

type failingBackend struct{ err error }

func (f failingBackend) Put(context.Context, string, io.Reader, int64, string) error { return f.err }
func (f failingBackend) Get(context.Context, string) (io.ReadCloser, error)          { return nil, f.err }
func (f failingBackend) Delete(context.Context, string) error                        { return f.err }
func (f failingBackend) Stat(context.Context, string) (storage.Object, error) {
	return storage.Object{}, f.err
}
func (f failingBackend) Health(context.Context) error { return nil }

var _ storage.Backend = failingBackend{}
```

- [ ] **Step 3: Run the tests to verify RED**

Run: `go test ./internal/tfstate ./internal/executor -run 'TestLayout|TestCompleteStep'`
Expected: FAIL — `layout.Output undefined` and `completeStep` returns no value.

- [ ] **Step 4: Add `Layout.Output`**

In `internal/tfstate/layout.go`:

```go
// Output is the archived console log for a job step. Upstream uses the same
// S3 key; its local backend writes to a home-directory path instead.
func (l Layout) Output(organizationID string, jobID int64, stepID string) string {
	return fmt.Sprintf("tfoutput/%s/%d/%s.tfoutput", organizationID, jobID, stepID)
}
```

- [ ] **Step 5: Implement `archiveStepLog` and change `completeStep`**

In `internal/executor/run.go`, add:

```go
func (e *Executor) archiveStepLog(ctx context.Context, job store.Job, step store.Step, output string) error {
	if e.deps.Storage == nil {
		return nil
	}
	key := e.deps.Layout.Output(job.OrganizationID, job.ID, step.ID)
	return e.deps.Storage.Put(ctx, key, strings.NewReader(output), int64(len(output)), "text/plain")
}
```

Rewrite `completeStep` to return `error`:

```go
func (e *Executor) completeStep(ctx context.Context, job store.Job, step store.Step, sink *logSink, status string) error {
	sink.Flush()
	output := sink.String()

	effective := status
	archiveErr := e.archiveStepLog(ctx, job, step, output)
	if archiveErr != nil {
		e.deps.Logger.Warn("could not archive step log", "step_id", step.ID, "error", archiveErr)
		if status == "completed" {
			effective = "failed"
		}
	}

	if err := e.deps.Store.UpdateStepOutput(ctx, store.UpdateStepOutputParams{
		ID:     step.ID,
		Output: stringPointer(output),
	}); err != nil {
		e.deps.Logger.Warn("could not save step output", "step_id", step.ID, "error", err)
	}
	if err := e.deps.Store.UpdateJobOutput(ctx, store.UpdateJobOutputParams{
		ID:     job.ID,
		Output: stringPointer(output),
	}); err != nil {
		e.deps.Logger.Warn("could not save job output", "job_id", job.ID, "error", err)
	}
	if _, err := e.deps.Store.UpdateStepStatus(ctx, store.UpdateStepStatusParams{
		ID:     step.ID,
		Status: stringPointer(effective),
	}); err != nil {
		e.deps.Logger.Warn("could not save step status", "step_id", step.ID, "error", err)
	}

	if archiveErr != nil && status == "completed" {
		return fmt.Errorf("archive step log: %w", archiveErr)
	}
	return nil
}
```

Update call sites:

- `run.go:71` and `run.go:78`: `_ = e.completeStep(ctx, job, step, sink, "failed")`
- `run.go:82`: `if err := e.completeStep(ctx, job, step, sink, "completed"); err != nil { return err }`
- `run.go:138`: same as line 82.
- `run.go:216`, `220`, `223`: `_ = e.completeStep(...)` (the original error is already being returned).

- [ ] **Step 6: Run the tests to verify GREEN**

Run: `go test ./internal/tfstate ./internal/executor`
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/tfstate/layout.go internal/tfstate/layout_test.go internal/executor/run.go internal/executor/artifacts_test.go
git commit -m "feat(executor): archive final step logs to object storage"
```

---

### Task 2: Serve archived logs from `/tfoutput/v1`

**Files:**
- Modify: `internal/tfoutput/handler.go:36-43` (`Deps`), `getStepLog` around line 91
- Modify: `cmd/api/main.go:120-126`
- Modify: `internal/tfoutput/handler_test.go`

**Interfaces:**
- Consumes: `Layout.Output(organizationID string, jobID int64, stepID string) string` from Task 1.
- Produces: `Deps` fields `Storage storage.Backend` and `Layout tfstate.Layout`.

- [ ] **Step 1: Write the failing handler tests**

Append to `internal/tfoutput/handler_test.go`:

```go
func TestGetStepLogPrefersArchivedObject(t *testing.T) {
	inline := "inline text"
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7, Output: &inline, Status: stringPtr("completed")},
	}
	backend := &logBackend{objects: map[string][]byte{}}
	layout := tfstate.NewLayout("local")
	backend.objects[layout.Output("org-1", 7, "step-1")] = []byte("archived text")

	rec := httptest.NewRecorder()
	handler := New(Deps{Store: q, Storage: backend, Layout: layout})
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "archived text" {
		t.Fatalf("body = %q, want archived text", rec.Body.String())
	}
}

func TestGetStepLogFallsBackToInlineOutput(t *testing.T) {
	inline := "inline text"
	q := &stubQuerier{
		job:  store.Job{ID: 7, OrganizationID: "org-1"},
		step: store.Step{ID: "step-1", JobID: 7, Output: &inline, Status: stringPtr("completed")},
	}
	handler := New(Deps{Store: q, Storage: &logBackend{objects: map[string][]byte{}}, Layout: tfstate.NewLayout("local")})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/org-1/job/7/step/step-1", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != inline {
		t.Fatalf("body = %q, want inline text", rec.Body.String())
	}
}

type logBackend struct{ objects map[string][]byte }

func (b *logBackend) Put(_ context.Context, key string, r io.Reader, _ int64, _ string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	b.objects[key] = data
	return nil
}

func (b *logBackend) Get(_ context.Context, key string) (io.ReadCloser, error) {
	data, ok := b.objects[key]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (b *logBackend) Delete(_ context.Context, key string) error { delete(b.objects, key); return nil }
func (b *logBackend) Stat(_ context.Context, key string) (storage.Object, error) {
	data, ok := b.objects[key]
	if !ok {
		return storage.Object{}, storage.ErrNotFound
	}
	return storage.Object{Key: key, Size: int64(len(data))}, nil
}
func (b *logBackend) Health(context.Context) error { return nil }
```

Add imports: `bytes`, `io`, `internal/storage`, `internal/tfstate`.

- [ ] **Step 2: Run the tests to verify RED**

Run: `go test ./internal/tfoutput -run 'TestGetStepLogPrefers|TestGetStepLogFallsBack'`
Expected: FAIL — `Deps` has no `Storage`/`Layout`.

- [ ] **Step 3: Add the storage fields and read path**

In `internal/tfoutput/handler.go`:

```go
type Deps struct {
	Store        store.Querier
	Queue        *queue.Queue
	Storage      storage.Backend
	Layout       tfstate.Layout
	Logger       *slog.Logger
	AuthToken    string
	InternalAuth *internalauth.Validator
}
```

In `getStepLog`, after the live-stream branch and before the `step.Output` check:

```go
	if h.deps.Storage != nil {
		key := h.deps.Layout.Output(job.OrganizationID, job.ID, step.ID)
		if rc, err := h.deps.Storage.Get(r.Context(), key); err == nil {
			defer rc.Close()
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			if _, err := io.Copy(w, rc); err != nil {
				h.deps.Logger.Warn("stream archived step log failed", "job_id", job.ID, "step_id", step.ID, "error", err)
			}
			return
		} else if !errors.Is(err, storage.ErrNotFound) {
			h.deps.Logger.Warn("read archived step log failed", "job_id", job.ID, "step_id", step.ID, "error", err)
		}
	}
```

Add imports `io`, `internal/storage`, `internal/tfstate`.

- [ ] **Step 4: Wire the API**

In `cmd/api/main.go`, add to the `tfoutput.New` call:

```go
		Storage:      storeBackend,
		Layout:       tfstate.NewLayout(cfg.Storage.Type),
```

- [ ] **Step 5: Run the tests to verify GREEN**

Run: `go test ./internal/tfoutput ./cmd/...`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/tfoutput/handler.go internal/tfoutput/handler_test.go cmd/api/main.go
git commit -m "feat(api): serve archived step logs"
```

---

### Task 3: Record the state URL and fail on history errors

**Files:**
- Modify: `internal/executor/history.go`
- Modify: `internal/executor/run.go` (`runApply` tail)
- Modify: `internal/executor/executor_test.go` (`TestArchiveHistory`)
- Modify: `internal/executor/artifacts_test.go`

**Interfaces:**
- Produces: `func (e *Executor) stateJSONURL(job store.Job, ws store.Workspace, historyID string) string`.
- Produces: `func (e *Executor) finishApply(ctx context.Context, job store.Job, ws store.Workspace, step store.Step, sink *logSink) error`.

- [ ] **Step 1: Write the failing tests**

In `internal/executor/executor_test.go`, inside `TestArchiveHistory`, after the job reference assertion:

```go
	wantURL := "https://terrakube.example.com/tfstate/v1/organization/org-1/workspace/ws-1/state/" + q.created.ID + ".json"
	if q.created.Output == nil || *q.created.Output != wantURL {
		t.Errorf("history output = %v, want %q", q.created.Output, wantURL)
	}
```

The test already constructs `e := New(Deps{...})` without `PublicURL`; add `PublicURL: "https://terrakube.example.com/"` there.

Append to `internal/executor/artifacts_test.go`:

```go
func TestFinishApplyFailsWhenHistoryArchiveFails(t *testing.T) {
	ctx := context.Background()
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	state := []byte(`{"version":4,"serial":3,"lineage":"lineage-1"}`)
	if err := backend.Put(ctx, layout.CurrentState("org-1", "ws-1"), bytes.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("seed state: %v", err)
	}
	st := &artifactStore{historyErr: errors.New("history database down")}
	e := New(Deps{Store: st, Storage: backend, Layout: layout, Logger: testLogger()})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	err := e.finishApply(ctx, store.Job{ID: 9, OrganizationID: "org-1"}, store.Workspace{ID: "ws-1"}, store.Step{ID: "step-1"}, sink)
	if err == nil || !strings.Contains(err.Error(), "archive state history") {
		t.Fatalf("error = %v, want archive state history", err)
	}
	if st.stepStatus != "failed" {
		t.Errorf("step status = %q, want failed", st.stepStatus)
	}
}

func TestFinishApplyCompletesStep(t *testing.T) {
	ctx := context.Background()
	st := &artifactStore{}
	e := New(Deps{Store: st, Logger: testLogger()})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	if err := e.finishApply(ctx, store.Job{ID: 9, OrganizationID: "org-1"}, store.Workspace{ID: "ws-1"}, store.Step{ID: "step-1"}, sink); err != nil {
		t.Fatalf("finishApply: %v", err)
	}
	if st.stepStatus != "completed" {
		t.Errorf("step status = %q, want completed", st.stepStatus)
	}
}
```

Add `bytes` to the imports of `artifacts_test.go`.

- [ ] **Step 2: Run the tests to verify RED**

Run: `go test ./internal/executor -run 'TestArchiveHistory|TestFinishApply'`
Expected: FAIL — `finishApply undefined` and the history output assertion fails.

- [ ] **Step 3: Add the URL helper and set it on history rows**

In `internal/executor/history.go`, before `CreateHistory`:

```go
	output := e.stateJSONURL(job, ws, historyID)
```

and pass `Output: &output` instead of `nil`.

Add:

```go
// stateJSONURL is the upstream state URL stored on a history row.
func (e *Executor) stateJSONURL(job store.Job, ws store.Workspace, historyID string) string {
	base := strings.TrimRight(e.deps.PublicURL, "/")
	if base == "" {
		return ""
	}
	return fmt.Sprintf("%s/tfstate/v1/organization/%s/workspace/%s/state/%s.json",
		base, job.OrganizationID, ws.ID, historyID)
}
```

Add `strings` to the imports.

- [ ] **Step 4: Add `finishApply` and use it in `runApply`**

In `internal/executor/run.go`:

```go
// finishApply archives the state produced by apply and completes the step. A
// history upload failure is a run failure, not a warning.
func (e *Executor) finishApply(ctx context.Context, job store.Job, ws store.Workspace, step store.Step, sink *logSink) error {
	if err := e.archiveHistory(ctx, job, ws); err != nil {
		_ = e.completeStep(ctx, job, step, sink, "failed")
		return fmt.Errorf("archive state history: %w", err)
	}
	return e.completeStep(ctx, job, step, sink, "completed")
}
```

Replace the tail of `runApply`:

```go
	if err := e.finishApply(ctx, job, ws, step, sink); err != nil {
		return err
	}
	return nil
```

- [ ] **Step 5: Run the tests to verify GREEN**

Run: `go test ./internal/executor -run 'TestArchiveHistory|TestFinishApply'`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/executor/history.go internal/executor/run.go internal/executor/executor_test.go internal/executor/artifacts_test.go
git commit -m "feat(executor): record state URL and fail on archive errors"
```

---

### Task 4: End-to-end coverage, docs and graph

**Files:**
- Modify: `internal/executor/e2e_test.go`
- Modify: `docs/compat.md` (§ Executor contract, Logs/History/Divergences)
- Regenerate: `graphify-out/`

**Interfaces:**
- Consumes: `Layout.Output` (Task 1), archived-log reads (Task 2), `history.output` (Task 3).

- [ ] **Step 1: Seed a state object and extend the e2e assertions**

In `internal/executor/e2e_test.go`, after the workspace is created and before the job, seed state:

```go
	state := []byte(`{"version":4,"serial":3,"lineage":"lineage-e2e"}`)
	if err := backend.Put(ctx, layout.CurrentState(orgID, wsID), bytes.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("seed state: %v", err)
	}
```

After the step loop, assert each step has an archived log:

```go
	for _, step := range steps {
		rc, err := backend.Get(ctx, layout.Output(orgID, job.ID, step.ID))
		if err != nil {
			t.Fatalf("read archived log for step %d: %v", step.StepNumber, err)
		}
		data, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			t.Fatalf("read archived log body: %v", err)
		}
		if !bytes.Contains(data, []byte("fake tofu")) {
			t.Errorf("archived log for step %d = %q", step.StepNumber, data)
		}
	}
```

After the job-output assertions, assert the history row:

```go
	jobReference := strconv.FormatInt(job.ID, 10)
	history, err := st.GetHistoryByJobReference(ctx, store.GetHistoryByJobReferenceParams{
		WorkspaceID:  wsID,
		JobReference: &jobReference,
	})
	if err != nil {
		t.Fatalf("load history: %v", err)
	}
	wantURL := "https://example.invalid/tfstate/v1/organization/" + orgID + "/workspace/" + wsID + "/state/" + history.ID + ".json"
	if history.Output == nil || *history.Output != wantURL {
		t.Errorf("history output = %v, want %q", history.Output, wantURL)
	}
	if history.Serial != 3 || history.Lineage != "lineage-e2e" {
		t.Errorf("history serial/lineage = %d/%q", history.Serial, history.Lineage)
	}
	sum := md5.Sum(state)
	if history.Md5 != hex.EncodeToString(sum[:]) {
		t.Errorf("history md5 = %q", history.Md5)
	}
```

Add `crypto/md5`, `encoding/hex`, `io` to the imports.

- [ ] **Step 2: Run the gated e2e test**

```bash
TERRAKUBE_TEST_DATABASE_URL='postgres://terrakube@localhost:5432/terrakube?sslmode=disable' \
TERRAKUBE_TEST_REDIS_URL='redis://localhost:6379/9' \
go test -race -run TestExecutorEndToEnd -v ./internal/executor
```

Expected: PASS.

- [ ] **Step 3: Update `docs/compat.md`**

Replace the **Logs** bullet with:

```markdown
- **Logs**: each command line is published to the Redis stream named after the
  job with `stepId`/`output` fields. The stream is capped at 200,000 records;
  ordering and step boundaries come from the stream order and `stepId`, and the
  stream ID carries the timestamp. When a step finishes, the final output is
  uploaded to `tfoutput/{orgId}/{jobId}/{stepId}.tfoutput` (the upstream S3
  key). `/tfoutput/v1` serves that object for finished steps and falls back to
  `step.output`, which still holds the same text inline (truncated at 10 MiB)
  for the TFE log endpoints. A log upload failure marks the step `failed` and
  fails the run, unless the step was already cancelled or failed.
```

Replace the **History** bullet with:

```markdown
- **History**: after a successful apply the executor reads the current state,
  creates a `history` row with the state `serial`, `lineage` and MD5, sets
  `history.output` to
  `{publicUrl}/tfstate/v1/organization/{orgId}/workspace/{wsId}/state/{historyId}.json`,
  and archives the raw and JSON state objects under the upstream `state/...`
  keys. Archiving is idempotent per job reference, so a retried job cannot
  create a second row. A history upload failure marks the step `failed` and
  fails the run.
```

Update the divergence bullet to:

```markdown
- `step.output` stores log text inline; upstream stores a `/tfoutput/v1/...`
  URL there. The executor also archives the text under the upstream
  `tfoutput/...` object key, which `/tfoutput/v1` prefers; the column remains
  for the TFE log endpoints and rows written before archiving.
```

- [ ] **Step 4: Run the full checks**

```bash
make check
go test -race ./...
```

Expected: PASS.

- [ ] **Step 5: Rebuild the knowledge graph**

```bash
graphify update .
```

Expected: `graph.json`, `GRAPH_REPORT.md` and `manifest.json` updated.

- [ ] **Step 6: Commit**

```bash
git add internal/executor/e2e_test.go docs/compat.md graphify-out/graph.json graphify-out/GRAPH_REPORT.md graphify-out/manifest.json
git commit -m "test(executor): cover archived logs and history end to end"
```

---

## Self-review notes

- Spec coverage: log archiving (Task 1), serving (Task 2), history URL and
  failure semantics (Task 3), end-to-end proof and docs (Task 4).
- Plan upload failure already marks the step failed in `runPlan`; no code
  change is needed there, and Task 3's `finishApply` covers the history case.
- No schema or HTTP API changes; `step.output` stays inline as agreed.
