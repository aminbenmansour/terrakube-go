package executor

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

func TestHandleAcksTerminalJob(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "completed"}}
	q := &fakeQueue{}
	e := newTestExecutor(st, q)

	e.handle(context.Background(), jobMessage("1-1", 7, 0))

	if got := len(q.acks); got != 1 {
		t.Fatalf("acks = %d, want 1", got)
	}
	if got := len(q.published); got != 0 {
		t.Fatalf("published = %d, want 0", got)
	}
}

func TestHandleRetriesTransientFailure(t *testing.T) {
	st := &reliabilityStore{getJobErrs: []error{errors.New("database unavailable")}}
	q := &fakeQueue{}
	e := newTestExecutor(st, q)

	e.handle(context.Background(), jobMessage("1-1", 7, 0))

	if got := len(q.acks); got != 1 {
		t.Fatalf("acks = %d, want 1", got)
	}
	if got := len(q.published); got != 1 {
		t.Fatalf("published = %d, want 1", got)
	}
	if got := q.published[0]["attempt"]; got != 1 {
		t.Errorf("retry attempt = %v, want 1", got)
	}
	if got := q.published[0]["job_id"]; got != "7" {
		t.Errorf("retry job_id = %v, want 7", got)
	}
}

func TestHandleGivesUpAfterMaxAttempts(t *testing.T) {
	st := &reliabilityStore{
		job:        store.Job{ID: 7, Status: "queue"},
		getJobErrs: []error{errors.New("database unavailable")},
	}
	q := &fakeQueue{}
	e := newTestExecutor(st, q)

	e.handle(context.Background(), jobMessage("1-1", 7, maxJobAttempts-1))

	if got := len(q.published); got != 0 {
		t.Fatalf("published = %d, want 0 after the last attempt", got)
	}
	if got := len(q.acks); got != 1 {
		t.Fatalf("acks = %d, want 1", got)
	}
	if got := st.status(); got != "failed" {
		t.Fatalf("job status = %q, want failed", got)
	}
}

func TestHandleLeavesMessagePendingOnShutdown(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "running"}}
	q := &fakeQueue{}
	e := newTestExecutor(st, q)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	e.handle(ctx, jobMessage("1-1", 7, 0))

	if got := len(q.acks); got != 0 {
		t.Fatalf("acks = %d, want 0 during shutdown", got)
	}
	if got := len(q.published); got != 0 {
		t.Fatalf("published = %d, want 0 during shutdown", got)
	}
}

func TestWatchJobCancelsWhenStatusBecomesCancelled(t *testing.T) {
	st := &reliabilityStore{
		job:              store.Job{ID: 7, Status: "running"},
		cancelAfterCalls: 2,
	}
	e := New(Deps{Store: st, Logger: testLogger()})
	e.cancelPollInterval = 5 * time.Millisecond

	ctx, stop := e.watchJob(context.Background(), 7)
	defer stop()

	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("job context was not cancelled after the job was cancelled")
	}
}

func TestWatchJobStopsWhenJobCompletes(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "running"}}
	e := New(Deps{Store: st, Logger: testLogger()})
	e.cancelPollInterval = 5 * time.Millisecond

	ctx, stop := e.watchJob(context.Background(), 7)
	st.setStatus("completed")
	time.Sleep(25 * time.Millisecond)

	if ctx.Err() != nil {
		t.Fatalf("job context was cancelled for a completed job: %v", ctx.Err())
	}
	stop()
}

func TestArchiveHistorySkipsExistingRow(t *testing.T) {
	ctx := context.Background()
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	state := []byte(`{"version":4,"serial":7,"lineage":"lineage-1","terraform_version":"1.9.0"}`)
	if err := backend.Put(ctx, layout.CurrentState("org-1", "ws-1"), bytes.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("seed state: %v", err)
	}

	st := &reliabilityStore{existingHistory: &store.History{ID: "history-1"}}
	e := New(Deps{
		Store:   st,
		Storage: backend,
		Layout:  layout,
		Logger:  testLogger(),
	})
	job := store.Job{ID: 9, OrganizationID: "org-1", WorkspaceID: "ws-1"}
	if err := e.archiveHistory(ctx, job, store.Workspace{ID: "ws-1"}); err != nil {
		t.Fatalf("archive history: %v", err)
	}

	if got := len(st.createHistoryCalls); got != 0 {
		t.Fatalf("created %d history rows, want 0 for a retried job", got)
	}
}

func TestRunPlanStopsWhenJobIsCancelled(t *testing.T) {
	ctx := context.Background()
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	archive := tarGz(t, map[string]string{"main.tf": `resource "terraform_data" "x" {}`})
	if err := backend.Put(ctx, layout.Content("cv-1"), bytes.NewReader(archive), int64(len(archive)), "application/gzip"); err != nil {
		t.Fatalf("seed configuration: %v", err)
	}

	script := filepath.Join(t.TempDir(), "tofu")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nexec sleep 5\n"), 0o755); err != nil {
		t.Fatalf("write fake tofu: %v", err)
	}

	st := &cancelRunStore{
		job: store.Job{
			ID:             9,
			OrganizationID: "org-1",
			WorkspaceID:    "ws-1",
			Status:         "running",
			OverrideSource: ptr("https://host/remote/tfe/v2/configuration-versions/cv-1/terraformContent.tar.gz"),
		},
		workspace:        store.Workspace{ID: "ws-1"},
		cancelAfterCalls: 2,
	}
	e := New(Deps{
		Store:   st,
		Storage: backend,
		Layout:  layout,
		Logger:  testLogger(),
		Binary:  script,
		WorkDir: t.TempDir(),
	})
	e.cancelPollInterval = 5 * time.Millisecond

	start := time.Now()
	err := e.runPlan(ctx, st.job, st.workspace, store.Step{ID: "step-1", JobID: 9, StepNumber: 100, Status: ptr("pending")})
	elapsed := time.Since(start)

	if !errors.Is(err, errJobCancelled) {
		t.Fatalf("runPlan error = %v, want errJobCancelled", err)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("runPlan took %s; the terraform command was not stopped", elapsed)
	}
	if got := st.stepStatus(); got != "cancelled" {
		t.Fatalf("step status = %q, want cancelled", got)
	}
	if got := st.status(); got != "cancelled" {
		t.Fatalf("job status = %q, want cancelled", got)
	}
}

func jobMessage(id string, jobID int64, attempt int) queue.Message {
	return queue.Message{ID: id, Values: map[string]any{
		"job_id":  jobID,
		"attempt": attempt,
	}}
}

func newTestExecutor(st store.Querier, q JobQueue) *Executor {
	e := New(Deps{Store: st, Queue: q, Logger: testLogger()})
	e.retryDelay = func(int) time.Duration { return 0 }
	e.cancelPollInterval = 5 * time.Millisecond
	return e
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// fakeQueue records acknowledgements and republished retries.
type fakeQueue struct {
	acks      []string
	deleted   []string
	dropped   []string
	published []map[string]any
	stale     []queue.Message
}

func (q *fakeQueue) EnsureGroup(context.Context, string, string) error { return nil }

func (q *fakeQueue) Delete(_ context.Context, _, id string) error {
	q.deleted = append(q.deleted, id)
	return nil
}

func (q *fakeQueue) DeleteStream(_ context.Context, stream string) error {
	q.dropped = append(q.dropped, stream)
	return nil
}

func (q *fakeQueue) ClaimStale(context.Context, string, string, string, time.Duration, int64) ([]queue.Message, error) {
	return q.stale, nil
}

func (q *fakeQueue) ReadGroup(context.Context, string, string, string, int64, time.Duration) ([]queue.Message, error) {
	return nil, nil
}

func (q *fakeQueue) Ack(_ context.Context, _, _, id string) error {
	q.acks = append(q.acks, id)
	return nil
}

func (q *fakeQueue) Publish(_ context.Context, _ string, values map[string]any) (string, error) {
	q.published = append(q.published, values)
	return "2-2", nil
}

func (q *fakeQueue) PublishCapped(context.Context, string, map[string]any, int64) (string, error) {
	return "2-2", nil
}

// reliabilityStore implements only the queries the reliability paths touch.
// Every other method panics through the embedded nil interface, which makes an
// unexpected call visible in the test rather than silently succeeding.
type reliabilityStore struct {
	store.Querier

	mu               sync.Mutex
	job              store.Job
	getJobErrs       []error
	getJobCalls      int
	cancelAfterCalls int
	workspace        store.Workspace
	existingHistory  *store.History

	updateJobStatusCalls  []store.UpdateJobStatusParams
	updateStepStatusCalls []store.UpdateStepStatusParams
	createHistoryCalls    []store.CreateHistoryParams
}

func (s *reliabilityStore) GetJob(ctx context.Context, _ int64) (store.Job, error) {
	if err := ctx.Err(); err != nil {
		return store.Job{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.getJobCalls++
	if s.cancelAfterCalls > 0 && s.getJobCalls >= s.cancelAfterCalls {
		s.job.Status = "cancelled"
	}
	if len(s.getJobErrs) > 0 {
		err := s.getJobErrs[0]
		s.getJobErrs = s.getJobErrs[1:]
		return store.Job{}, err
	}
	return s.job, nil
}

func (s *reliabilityStore) GetWorkspace(context.Context, string) (store.Workspace, error) {
	return s.workspace, nil
}

func (s *reliabilityStore) GetPendingStepByJob(context.Context, int64) (store.Step, error) {
	return store.Step{}, pgx.ErrNoRows
}

func (s *reliabilityStore) UpdateJobStatus(_ context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.updateJobStatusCalls = append(s.updateJobStatusCalls, arg)
	s.job.Status = arg.Status
	return s.job, nil
}

func (s *reliabilityStore) UpdateStepStatus(_ context.Context, arg store.UpdateStepStatusParams) (store.Step, error) {
	s.updateStepStatusCalls = append(s.updateStepStatusCalls, arg)
	return store.Step{ID: arg.ID}, nil
}

func (s *reliabilityStore) GetHistoryByJobReference(context.Context, store.GetHistoryByJobReferenceParams) (store.History, error) {
	if s.existingHistory == nil {
		return store.History{}, pgx.ErrNoRows
	}
	return *s.existingHistory, nil
}

func (s *reliabilityStore) CreateHistory(_ context.Context, arg store.CreateHistoryParams) (store.History, error) {
	s.createHistoryCalls = append(s.createHistoryCalls, arg)
	return store.History{ID: arg.ID}, nil
}

func (s *reliabilityStore) setStatus(status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.job.Status = status
}

func (s *reliabilityStore) status() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.job.Status
}

// cancelRunStore flips the job to cancelled after a few status reads, which is
// what the API does when a run is discarded while the executor is working.
type cancelRunStore struct {
	store.Querier

	mu               sync.Mutex
	job              store.Job
	workspace        store.Workspace
	calls            int
	cancelAfterCalls int
	lastStepStatus   string
}

func (s *cancelRunStore) GetJob(context.Context, int64) (store.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	if s.cancelAfterCalls > 0 && s.calls >= s.cancelAfterCalls {
		s.job.Status = "cancelled"
	}
	return s.job, nil
}

func (s *cancelRunStore) GetWorkspace(context.Context, string) (store.Workspace, error) {
	return s.workspace, nil
}

func (s *cancelRunStore) UpdateJobStatus(_ context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.job.Status = arg.Status
	return s.job, nil
}

func (s *cancelRunStore) UpdateStepStatus(_ context.Context, arg store.UpdateStepStatusParams) (store.Step, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if arg.Status != nil {
		s.lastStepStatus = *arg.Status
	}
	return store.Step{ID: arg.ID}, nil
}

func (s *cancelRunStore) status() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.job.Status
}

func (s *cancelRunStore) stepStatus() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastStepStatus
}

func (s *cancelRunStore) UpdateStepOutput(context.Context, store.UpdateStepOutputParams) error {
	return nil
}

func (s *cancelRunStore) UpdateJobOutput(context.Context, store.UpdateJobOutputParams) error {
	return nil
}

func (s *cancelRunStore) ListVariablesByWorkspace(context.Context, string) ([]store.Variable, error) {
	return nil, nil
}

var (
	_ JobQueue = (*fakeQueue)(nil)
	_          = time.Second
)
