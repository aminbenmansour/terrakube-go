package runtrigger

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

var testNow = time.Date(2026, 11, 26, 12, 0, 0, 0, time.UTC)

type fakeStore struct {
	store.Querier

	jobs         map[int64]store.Job
	steps        []store.Step
	triggers     []store.WorkspaceRunTrigger
	workspaces   map[string]store.Workspace
	templates    map[string]store.Template
	jobHistory   map[string]store.History
	prevHistory  map[string]store.History
	events       []store.RunTriggerEvent
	records      []store.RecordRunTriggerEventParams
	createdJobs  []store.CreateJobParams
	createdSteps []store.CreateStepParams
	deletedJobs  []int64
	lastJobs     []string
	enqueueErr   error
	nextID       int64
	now          func() time.Time
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		jobs:        map[int64]store.Job{},
		workspaces:  map[string]store.Workspace{},
		templates:   map[string]store.Template{},
		jobHistory:  map[string]store.History{},
		prevHistory: map[string]store.History{},
		nextID:      100,
		now:         func() time.Time { return testNow },
	}
}

func (f *fakeStore) GetJob(_ context.Context, id int64) (store.Job, error) {
	job, ok := f.jobs[id]
	if !ok {
		return store.Job{}, pgx.ErrNoRows
	}
	return job, nil
}

func (f *fakeStore) ListStepsByJob(_ context.Context, jobID int64) ([]store.Step, error) {
	var steps []store.Step
	for _, step := range f.steps {
		if step.JobID == jobID {
			steps = append(steps, step)
		}
	}
	return steps, nil
}

func (f *fakeStore) GetHistoryByJobReference(_ context.Context, arg store.GetHistoryByJobReferenceParams) (store.History, error) {
	history, ok := f.jobHistory[arg.WorkspaceID+"|"+*arg.JobReference]
	if !ok {
		return store.History{}, pgx.ErrNoRows
	}
	return history, nil
}

func (f *fakeStore) GetHistoryBefore(_ context.Context, arg store.GetHistoryBeforeParams) (store.History, error) {
	history, ok := f.prevHistory[arg.WorkspaceID]
	if !ok {
		return store.History{}, pgx.ErrNoRows
	}
	return history, nil
}

func (f *fakeStore) ListEnabledRunTriggersBySource(_ context.Context, sourceWorkspaceID string) ([]store.WorkspaceRunTrigger, error) {
	var triggers []store.WorkspaceRunTrigger
	for _, trigger := range f.triggers {
		if trigger.SourceWorkspaceID == sourceWorkspaceID && trigger.Enabled {
			triggers = append(triggers, trigger)
		}
	}
	return triggers, nil
}

func (f *fakeStore) GetJobByWorkspaceAndTrigger(_ context.Context, arg store.GetJobByWorkspaceAndTriggerParams) (store.Job, error) {
	for _, job := range f.jobs {
		if job.WorkspaceID == arg.WorkspaceID && job.TriggeredByJobID != nil && *job.TriggeredByJobID == *arg.TriggeredByJobID {
			return job, nil
		}
	}
	return store.Job{}, pgx.ErrNoRows
}

func (f *fakeStore) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	workspace, ok := f.workspaces[id]
	if !ok {
		return store.Workspace{}, pgx.ErrNoRows
	}
	return workspace, nil
}

func (f *fakeStore) GetTemplate(_ context.Context, id string) (store.Template, error) {
	template, ok := f.templates[id]
	if !ok {
		return store.Template{}, pgx.ErrNoRows
	}
	return template, nil
}

func (f *fakeStore) CreateJob(_ context.Context, arg store.CreateJobParams) (store.Job, error) {
	f.createdJobs = append(f.createdJobs, arg)
	f.nextID++
	job := store.Job{
		ID:                f.nextID,
		OrganizationID:    arg.OrganizationID,
		WorkspaceID:       arg.WorkspaceID,
		Status:            arg.Status,
		Via:               arg.Via,
		TemplateReference: arg.TemplateReference,
		Tcl:               arg.Tcl,
		TriggeredByJobID:  arg.TriggeredByJobID,
		CascadeDepth:      arg.CascadeDepth,
	}
	f.jobs[job.ID] = job
	return job, nil
}

func (f *fakeStore) CreateStep(_ context.Context, arg store.CreateStepParams) (store.Step, error) {
	f.createdSteps = append(f.createdSteps, arg)
	return store.Step{ID: arg.ID, JobID: arg.JobID, StepNumber: arg.StepNumber}, nil
}

func (f *fakeStore) DeleteJob(_ context.Context, id int64) error {
	f.deletedJobs = append(f.deletedJobs, id)
	delete(f.jobs, id)
	return nil
}

func (f *fakeStore) SetWorkspaceLastJob(_ context.Context, arg store.SetWorkspaceLastJobParams) error {
	f.lastJobs = append(f.lastJobs, arg.LastJobStatus)
	return nil
}

func (f *fakeStore) ClaimRunTriggerEvent(context.Context) (store.RunTriggerEvent, error) {
	for i, event := range f.events {
		if event.Status != "PENDING" {
			continue
		}
		if event.NextAttemptAt.Valid && event.NextAttemptAt.Time.After(f.now()) {
			continue
		}
		event.Status = "PROCESSING"
		event.AttemptCount++
		event.LastAttemptAt = pgtype.Timestamptz{Time: f.now(), Valid: true}
		f.events[i] = event
		return event, nil
	}
	return store.RunTriggerEvent{}, pgx.ErrNoRows
}

func (f *fakeStore) RecordRunTriggerEvent(_ context.Context, arg store.RecordRunTriggerEventParams) (int64, error) {
	f.records = append(f.records, arg)
	for i, event := range f.events {
		if event.ID == arg.ID {
			event.Status = arg.Status
			event.LastError = arg.LastError
			event.NextAttemptAt = arg.NextAttemptAt
			f.events[i] = event
		}
	}
	return 1, nil
}

func (f *fakeStore) ReclaimStuckRunTriggerEvents(context.Context, store.ReclaimStuckRunTriggerEventsParams) (int64, error) {
	return 0, nil
}

func (f *fakeStore) FailStuckRunTriggerEvents(context.Context, store.FailStuckRunTriggerEventsParams) (int64, error) {
	return 0, nil
}

type fakeQueue struct {
	enqueued []int64
	err      error
}

func (f *fakeQueue) EnqueueJob(_ context.Context, jobID int64) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	f.enqueued = append(f.enqueued, jobID)
	return "1-0", nil
}

func newDispatcherForTest(st *fakeStore, q *fakeQueue) *Dispatcher {
	return &Dispatcher{
		Store:    st,
		Queue:    q,
		Logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
		Now:      func() time.Time { return testNow },
		Interval: time.Second,
	}
}

// seedCompletedJob installs a completed job with a completed apply step and an
// outbox row the claim will pick up. Callers add triggers and workspaces.
func seedCompletedJob(st *fakeStore) store.Job {
	job := store.Job{
		ID: 1, OrganizationID: "org", WorkspaceID: "source", Status: "completed",
		Via: "UI", Refresh: true, PlanChanges: true,
	}
	st.jobs[job.ID] = job
	st.steps = []store.Step{
		{ID: "plan", JobID: job.ID, StepNumber: planStepNumber, Status: strptr("completed")},
		{ID: "apply", JobID: job.ID, StepNumber: applyStepNumber, Status: strptr("completed")},
	}
	st.events = []store.RunTriggerEvent{{
		ID: "event", JobID: job.ID, Status: "PENDING",
	}}
	return job
}

func addDestination(st *fakeStore, workspaceID, templateID string) store.WorkspaceRunTrigger {
	st.workspaces[workspaceID] = store.Workspace{
		ID: workspaceID, OrganizationID: "org", IacType: "terraform",
	}
	var template *string
	if templateID != "" {
		st.templates[templateID] = store.Template{
			ID: templateID, OrganizationID: "org", Tcl: strptr("flow"),
		}
		template = &templateID
	}
	trigger := store.WorkspaceRunTrigger{
		ID: "trigger-" + workspaceID, SourceWorkspaceID: "source",
		DestinationWorkspaceID: workspaceID, TemplateID: template,
		OrganizationID: "org", Enabled: true,
	}
	st.triggers = append(st.triggers, trigger)
	return trigger
}

func TestTickQueuesOneDependentPerEnabledTrigger(t *testing.T) {
	st := newFakeStore()
	seedCompletedJob(st)
	addDestination(st, "dest-a", "template-a")
	addDestination(st, "dest-b", "template-b")
	q := &fakeQueue{}
	d := newDispatcherForTest(st, q)

	d.Tick(context.Background())

	if len(st.createdJobs) != 2 {
		t.Fatalf("created jobs = %d, want 2: %+v", len(st.createdJobs), st.createdJobs)
	}
	for _, created := range st.createdJobs {
		if created.Via != "RUN_TRIGGER" {
			t.Errorf("via = %q, want RUN_TRIGGER", created.Via)
		}
		if created.Status != "pending" {
			t.Errorf("status = %q, want pending", created.Status)
		}
		if created.TriggeredByJobID == nil || *created.TriggeredByJobID != 1 {
			t.Errorf("triggered_by_job_id = %v, want 1", created.TriggeredByJobID)
		}
		if created.CascadeDepth != 1 {
			t.Errorf("cascade_depth = %d, want 1", created.CascadeDepth)
		}
		if created.TemplateReference == nil || *created.TemplateReference == "" {
			t.Error("template reference was not carried")
		}
	}
	if len(st.createdSteps) != 6 {
		t.Errorf("steps = %d, want 3 per dependent", len(st.createdSteps))
	}
	if len(q.enqueued) != 2 {
		t.Errorf("enqueued = %v, want both dependents", q.enqueued)
	}
	if len(st.records) != 1 || st.records[0].Status != "PROCESSED" {
		t.Fatalf("records = %+v, want one PROCESSED", st.records)
	}
	if len(st.lastJobs) != 2 || st.lastJobs[0] != "pending" {
		t.Errorf("last job updates = %v, want pending per dependent", st.lastJobs)
	}
}

func TestTickIgnoresPlanOnlyRun(t *testing.T) {
	st := newFakeStore()
	job := seedCompletedJob(st)
	st.steps = st.steps[:1] // plan completed, no apply step
	addDestination(st, "dest", "template")
	d := newDispatcherForTest(st, &fakeQueue{})

	d.Tick(context.Background())

	if len(st.createdJobs) != 0 {
		t.Errorf("created jobs = %d, want 0 for a plan-only run", len(st.createdJobs))
	}
	if len(st.records) != 1 || st.records[0].Status != "PROCESSED" {
		t.Errorf("records = %+v, want the event processed", st.records)
	}
	if _, err := st.GetJob(context.Background(), job.ID); err != nil {
		t.Fatalf("job disappeared: %v", err)
	}
}

func TestTickIgnoresUnchangedState(t *testing.T) {
	st := newFakeStore()
	seedCompletedJob(st)
	addDestination(st, "dest", "template")
	st.jobHistory["source|1"] = store.History{WorkspaceID: "source", Serial: 3, Md5: "abc"}
	st.prevHistory["source"] = store.History{WorkspaceID: "source", Serial: 3, Md5: "abc"}
	d := newDispatcherForTest(st, &fakeQueue{})

	d.Tick(context.Background())

	if len(st.createdJobs) != 0 {
		t.Errorf("created jobs = %d, want 0 when the state did not change", len(st.createdJobs))
	}
}

func TestTickDispatchesWhenSerialIncreased(t *testing.T) {
	st := newFakeStore()
	seedCompletedJob(st)
	addDestination(st, "dest", "template")
	st.jobHistory["source|1"] = store.History{WorkspaceID: "source", Serial: 4, Md5: "abc"}
	st.prevHistory["source"] = store.History{WorkspaceID: "source", Serial: 3, Md5: "def"}
	d := newDispatcherForTest(st, &fakeQueue{})

	d.Tick(context.Background())

	if len(st.createdJobs) != 1 {
		t.Fatalf("created jobs = %d, want 1", len(st.createdJobs))
	}
}

func TestTickStopsAtCascadeLimit(t *testing.T) {
	st := newFakeStore()
	job := seedCompletedJob(st)
	job.CascadeDepth = maxCascadeDepth
	st.jobs[job.ID] = job
	addDestination(st, "dest", "template")
	d := newDispatcherForTest(st, &fakeQueue{})

	d.Tick(context.Background())

	if len(st.createdJobs) != 0 {
		t.Errorf("created jobs = %d, want 0 past the depth limit", len(st.createdJobs))
	}
	if len(st.records) != 1 || st.records[0].Status != "PROCESSED" {
		t.Errorf("records = %+v, want the event processed", st.records)
	}
}

func TestTickRequeuesAPendingDependentInsteadOfDuplicatingIt(t *testing.T) {
	st := newFakeStore()
	seedCompletedJob(st)
	addDestination(st, "dest", "template")
	existing := store.Job{
		ID: 50, WorkspaceID: "dest", OrganizationID: "org", Status: "pending",
		TriggeredByJobID: int64ptr(1),
	}
	st.jobs[existing.ID] = existing
	q := &fakeQueue{}
	d := newDispatcherForTest(st, q)

	d.Tick(context.Background())

	if len(st.createdJobs) != 0 {
		t.Errorf("created jobs = %d, want no duplicate", len(st.createdJobs))
	}
	if len(q.enqueued) != 1 || q.enqueued[0] != existing.ID {
		t.Errorf("enqueued = %v, want the pending %d", q.enqueued, existing.ID)
	}
}

func TestTickSkipsADestinationWithoutATemplate(t *testing.T) {
	st := newFakeStore()
	seedCompletedJob(st)
	addDestination(st, "dest", "")
	d := newDispatcherForTest(st, &fakeQueue{})

	d.Tick(context.Background())

	if len(st.createdJobs) != 0 {
		t.Errorf("created jobs = %d, want 0 without a template", len(st.createdJobs))
	}
	if len(st.records) != 1 || st.records[0].Status != "PROCESSED" {
		t.Errorf("records = %+v, want the event processed", st.records)
	}
}

func TestTickRetriesAFailedEnqueueWithBackoff(t *testing.T) {
	st := newFakeStore()
	seedCompletedJob(st)
	addDestination(st, "dest", "template")
	q := &fakeQueue{err: errors.New("redis down")}
	d := newDispatcherForTest(st, q)

	d.Tick(context.Background())

	if len(st.records) != 1 {
		t.Fatalf("records = %+v, want one", st.records)
	}
	record := st.records[0]
	if record.Status != "PENDING" {
		t.Errorf("status = %q, want PENDING for a retry", record.Status)
	}
	if !record.NextAttemptAt.Valid || !record.NextAttemptAt.Time.After(testNow) {
		t.Errorf("next attempt = %+v, want a future retry", record.NextAttemptAt)
	}
	if record.LastError == nil {
		t.Error("last_error was not recorded")
	}
}

func TestTickMarksAFailedEventPastTheAttemptLimit(t *testing.T) {
	st := newFakeStore()
	seedCompletedJob(st)
	st.events[0].AttemptCount = eventMaxAttempts - 1
	addDestination(st, "dest", "template")
	q := &fakeQueue{err: errors.New("redis down")}
	d := newDispatcherForTest(st, q)

	d.Tick(context.Background())

	if len(st.records) != 1 || st.records[0].Status != "FAILED" {
		t.Fatalf("records = %+v, want one FAILED", st.records)
	}
	if st.records[0].NextAttemptAt.Valid {
		t.Error("a permanently failed event kept a retry time")
	}
}

func TestBackoffGrowsAndCaps(t *testing.T) {
	for attempt, want := range map[int32]time.Duration{
		1:  eventBackoffInitial,
		2:  2 * eventBackoffInitial,
		3:  4 * eventBackoffInitial,
		30: eventBackoffMax,
	} {
		if got := backoff(attempt); got != want {
			t.Errorf("backoff(%d) = %s, want %s", attempt, got, want)
		}
	}
}

func int64ptr(v int64) *int64 { return &v }
