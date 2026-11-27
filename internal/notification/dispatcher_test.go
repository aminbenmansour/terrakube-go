package notification

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type dispatcherStore struct {
	store.Querier
	queued        []store.NotificationOutbox
	jobs          map[int64]store.Job
	workspaces    map[string]store.Workspace
	organizations map[string]store.Organization
	configs       map[string]store.NotificationConfiguration
	steps         map[int64][]store.Step
	recorded      []store.RecordNotificationOutboxParams
}

func (s *dispatcherStore) ClaimNotificationOutbox(context.Context) (store.NotificationOutbox, error) {
	if len(s.queued) == 0 {
		return store.NotificationOutbox{}, pgx.ErrNoRows
	}
	row := s.queued[0]
	s.queued = s.queued[1:]
	return row, nil
}

func (s *dispatcherStore) ReclaimStuckNotificationOutbox(context.Context, store.ReclaimStuckNotificationOutboxParams) (int64, error) {
	return 0, nil
}

func (s *dispatcherStore) FailStuckNotificationOutbox(context.Context, store.FailStuckNotificationOutboxParams) (int64, error) {
	return 0, nil
}

func (s *dispatcherStore) RecordNotificationOutbox(_ context.Context, arg store.RecordNotificationOutboxParams) (int64, error) {
	s.recorded = append(s.recorded, arg)
	return 1, nil
}

func (s *dispatcherStore) GetJob(_ context.Context, id int64) (store.Job, error) {
	job, ok := s.jobs[id]
	if !ok {
		return store.Job{}, pgx.ErrNoRows
	}
	return job, nil
}

func (s *dispatcherStore) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	workspace, ok := s.workspaces[id]
	if !ok {
		return store.Workspace{}, pgx.ErrNoRows
	}
	return workspace, nil
}

func (s *dispatcherStore) GetOrganization(_ context.Context, id string) (store.Organization, error) {
	organization, ok := s.organizations[id]
	if !ok {
		return store.Organization{}, pgx.ErrNoRows
	}
	return organization, nil
}

func (s *dispatcherStore) GetNotificationConfiguration(_ context.Context, id string) (store.NotificationConfiguration, error) {
	config, ok := s.configs[id]
	if !ok {
		return store.NotificationConfiguration{}, pgx.ErrNoRows
	}
	return config, nil
}

func (s *dispatcherStore) ListStepsByJob(_ context.Context, jobID int64) ([]store.Step, error) {
	return s.steps[jobID], nil
}

const (
	testOrgID       = "11111111-1111-4111-8111-111111111111"
	testWorkspaceID = "22222222-2222-4222-8222-222222222222"
	testConfigID    = "33333333-3333-4333-8333-333333333333"
	testRowID       = "44444444-4444-4444-8444-444444444444"
)

func newDispatcherFixture() (*dispatcherStore, store.NotificationOutbox) {
	job := store.Job{
		ID: 7, OrganizationID: testOrgID, WorkspaceID: testWorkspaceID,
		Status: "completed", CommitID: strPtr("abc123"),
	}
	st := &dispatcherStore{
		jobs:          map[int64]store.Job{7: job},
		workspaces:    map[string]store.Workspace{testWorkspaceID: {ID: testWorkspaceID, Name: "networking"}},
		organizations: map[string]store.Organization{testOrgID: {ID: testOrgID, Name: "acme"}},
		steps:         map[int64][]store.Step{},
		configs: map[string]store.NotificationConfiguration{testConfigID: {
			ID: testConfigID, Name: "alerts", ChannelType: "WEBHOOK",
			MessageStyle: StyleDetailed, SigningSecret: strPtr("sign-me"),
		}},
	}
	row := store.NotificationOutbox{
		ID: testRowID, JobID: 7, ConfigurationID: testConfigID,
		Status: "SENDING", JobStatus: strPtr("completed"), AttemptCount: 1,
		LastAttemptAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	return st, row
}

func dispatcherWith(st *dispatcherStore, destination string) *Dispatcher {
	st.configs[testConfigID] = store.NotificationConfiguration{
		ID: testConfigID, Name: "alerts", ChannelType: "WEBHOOK",
		DestinationUrl: destination, MessageStyle: StyleDetailed, SigningSecret: strPtr("sign-me"),
	}
	now := time.Date(2026, 11, 27, 12, 0, 0, 0, time.UTC)
	return &Dispatcher{
		Store: st, Deliverer: NewDeliverer(false), PublicURL: "https://ui.example.com",
		Now: func() time.Time { return now },
	}
}

func TestDispatcherSendsAndRecordsSuccess(t *testing.T) {
	var body string
	var header http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		body = string(raw)
		header = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	st, row := newDispatcherFixture()
	dispatcherWith(st, server.URL).process(t.Context(), row)

	if len(st.recorded) != 1 {
		t.Fatalf("recorded = %d results, want 1", len(st.recorded))
	}
	recorded := st.recorded[0]
	if recorded.Status != "SENT" || recorded.LastError != nil || recorded.NextAttemptAt.Valid {
		t.Fatalf("recorded = %+v, want SENT", recorded)
	}
	for _, want := range []string{
		`"organization":"acme"`, `"workspace":"networking"`, `"status":"completed"`,
		`"jobId":7`, `"runUrl":"https://ui.example.com/organizations/` + testOrgID,
		`"commitId":"abc123"`, `"configurationName":"alerts"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("payload missing %q: %s", want, body)
		}
	}
	if got := header.Get("X-Terrakube-Signature"); got != Sign("sign-me", body) {
		t.Fatalf("signature = %q, want the webhook signature", got)
	}
}

func TestDispatcherRetriesTransientFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "120")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	st, row := newDispatcherFixture()
	dispatcherWith(st, server.URL).process(t.Context(), row)

	recorded := st.recorded[0]
	if recorded.Status != "PENDING" || recorded.LastError == nil {
		t.Fatalf("recorded = %+v, want a retry", recorded)
	}
	want := time.Date(2026, 11, 27, 12, 2, 0, 0, time.UTC)
	if !recorded.NextAttemptAt.Valid || !recorded.NextAttemptAt.Time.Equal(want) {
		t.Fatalf("retry at %v, want the Retry-After delay %v", recorded.NextAttemptAt, want)
	}
}

func TestDispatcherFailsPermanentRejection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	st, row := newDispatcherFixture()
	dispatcherWith(st, server.URL).process(t.Context(), row)

	if got := st.recorded[0].Status; got != "FAILED" {
		t.Fatalf("status = %q, want FAILED", got)
	}
	if st.recorded[0].NextAttemptAt.Valid {
		t.Fatal("a permanent rejection was scheduled for a retry")
	}
}

func TestDispatcherFailsAfterMaxAttempts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	st, row := newDispatcherFixture()
	row.AttemptCount = maxAttempts
	dispatcherWith(st, server.URL).process(t.Context(), row)

	if got := st.recorded[0].Status; got != "FAILED" {
		t.Fatalf("status = %q, want FAILED at the attempt limit", got)
	}
}

func TestDispatcherSkipsSupersededJobStatus(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
	}))
	defer server.Close()

	st, row := newDispatcherFixture()
	job := st.jobs[7]
	job.Status = "completed"
	row.JobStatus = strPtr("failed") // the row was written for a transient failure
	st.jobs[7] = job
	dispatcherWith(st, server.URL).process(t.Context(), row)

	if called {
		t.Fatal("delivered a notification whose job status moved on")
	}
	if got := st.recorded[0].Status; got != "SKIPPED" {
		t.Fatalf("status = %q, want SKIPPED", got)
	}
}

func TestDispatcherIncludesTheFailureSummary(t *testing.T) {
	var body string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		body = string(raw)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	st, row := newDispatcherFixture()
	job := st.jobs[7]
	job.Status = "failed"
	row.JobStatus = strPtr("failed")
	st.jobs[7] = job
	st.steps[7] = []store.Step{
		{StepNumber: 100, Status: strPtr("completed"), Output: strPtr("plan ok")},
		{StepNumber: 200, Status: strPtr("failed"), Output: strPtr("Error: invalid provider configuration")},
	}
	dispatcherWith(st, server.URL).process(t.Context(), row)

	if !strings.Contains(body, "Error: invalid provider configuration") {
		t.Fatalf("payload missing the failure summary: %s", body)
	}
}
