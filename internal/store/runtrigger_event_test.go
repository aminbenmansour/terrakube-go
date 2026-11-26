package store

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
)

// TestUpdateJobStatusWritesOneEventForAQualifyingCompletion covers the atomic
// completion/outbox statement: only a completed job with an enabled edge to a
// live destination gets a row, and a repeated completion does not duplicate it.
func TestUpdateJobStatusWritesOneEventForAQualifyingCompletion(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "run-trigger-events", ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}
	source, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
		ID: uuid.NewString(), Name: "source", IacType: "terraform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create source workspace: %v", err)
	}
	destination, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
		ID: uuid.NewString(), Name: "destination", IacType: "terraform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create destination workspace: %v", err)
	}
	edge, err := q.CreateRunTrigger(ctx, CreateRunTriggerParams{
		ID:                     uuid.NewString(),
		SourceWorkspaceID:      source.ID,
		DestinationWorkspaceID: destination.ID,
		OrganizationID:         org.ID,
		Enabled:                true,
	})
	if err != nil {
		t.Fatalf("create run trigger: %v", err)
	}

	newJob := func() Job {
		job, err := q.CreateJob(ctx, CreateJobParams{
			OrganizationID: org.ID, WorkspaceID: source.ID, Status: "pending", Via: "UI",
			Refresh: true, PlanChanges: true,
		})
		if err != nil {
			t.Fatalf("create job: %v", err)
		}
		return job
	}
	countEvents := func(jobID int64) int64 {
		t.Helper()
		var count int64
		if err := q.db.QueryRow(ctx, "SELECT count(*) FROM run_trigger_event WHERE job_id = $1", jobID).Scan(&count); err != nil {
			t.Fatalf("count events: %v", err)
		}
		return count
	}

	job := newJob()
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: job.ID, Status: "running"}); err != nil {
		t.Fatalf("mark running: %v", err)
	}
	if count := countEvents(job.ID); count != 0 {
		t.Fatalf("events after running = %d, want 0", count)
	}

	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: job.ID, Status: "completed"}); err != nil {
		t.Fatalf("mark completed: %v", err)
	}
	if count := countEvents(job.ID); count != 1 {
		t.Fatalf("events after completion = %d, want 1", count)
	}
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: job.ID, Status: "completed"}); err != nil {
		t.Fatalf("mark completed again: %v", err)
	}
	if count := countEvents(job.ID); count != 1 {
		t.Fatalf("events after repeated completion = %d, want 1", count)
	}

	// A job with no inbound-visible edge, a disabled edge and a soft-deleted
	// destination all stay eventless.
	if _, err := q.UpdateRunTrigger(ctx, UpdateRunTriggerParams{
		ID: edge.ID, SourceWorkspaceID: source.ID, DestinationWorkspaceID: destination.ID,
		OrganizationID: org.ID, Enabled: false,
	}); err != nil {
		t.Fatalf("disable edge: %v", err)
	}
	disabled := newJob()
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: disabled.ID, Status: "completed"}); err != nil {
		t.Fatalf("complete disabled job: %v", err)
	}
	if count := countEvents(disabled.ID); count != 0 {
		t.Errorf("events with a disabled edge = %d, want 0", count)
	}

	if _, err := q.UpdateRunTrigger(ctx, UpdateRunTriggerParams{
		ID: edge.ID, SourceWorkspaceID: source.ID, DestinationWorkspaceID: destination.ID,
		OrganizationID: org.ID, Enabled: true,
	}); err != nil {
		t.Fatalf("re-enable edge: %v", err)
	}
	if err := q.SoftDeleteWorkspace(ctx, destination.ID); err != nil {
		t.Fatalf("soft-delete destination: %v", err)
	}
	deletedDestination := newJob()
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: deletedDestination.ID, Status: "completed"}); err != nil {
		t.Fatalf("complete job with a deleted destination: %v", err)
	}
	if count := countEvents(deletedDestination.ID); count != 0 {
		t.Errorf("events with a deleted destination = %d, want 0", count)
	}
}

// TestRunTriggerEventClaimRecordAndSweep covers the poller's SQL: one claim at a
// time, a stale result is discarded, retries reclaim abandoned rows and rows at
// the attempt limit fail for good.
func TestRunTriggerEventClaimRecordAndSweep(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "run-trigger-claim", ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}
	workspace, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
		ID: uuid.NewString(), Name: "source", IacType: "terraform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	newJob := func() int64 {
		job, err := q.CreateJob(ctx, CreateJobParams{
			OrganizationID: org.ID, WorkspaceID: workspace.ID, Status: "completed", Via: "UI",
			Refresh: true, PlanChanges: true,
		})
		if err != nil {
			t.Fatalf("create job: %v", err)
		}
		return job.ID
	}
	insertEvent := func(id string, jobID int64) {
		if _, err := q.db.Exec(ctx,
			"INSERT INTO run_trigger_event (id, job_id, status) VALUES ($1, $2, 'PENDING')", id, jobID); err != nil {
			t.Fatalf("insert event: %v", err)
		}
	}

	insertEvent("first", newJob())
	claimed, err := q.ClaimRunTriggerEvent(ctx)
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if claimed.ID != "first" || claimed.Status != "PROCESSING" || claimed.AttemptCount != 1 {
		t.Fatalf("claimed = %+v, want first/ PROCESSING/1", claimed)
	}
	if _, err := q.ClaimRunTriggerEvent(ctx); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("second claim err = %v, want ErrNoRows", err)
	}

	rows, err := q.RecordRunTriggerEvent(ctx, RecordRunTriggerEventParams{
		ID: claimed.ID, LastAttemptAt: claimed.LastAttemptAt, Status: "PROCESSED",
	})
	if err != nil || rows != 1 {
		t.Fatalf("record = %d, %v; want 1", rows, err)
	}
	// A stale attempt (its last_attempt_at moved on) must not clobber the result.
	rows, err = q.RecordRunTriggerEvent(ctx, RecordRunTriggerEventParams{
		ID: claimed.ID, LastAttemptAt: pgtype.Timestamptz{Time: time.Now().Add(-time.Hour), Valid: true},
		Status: "PENDING",
	})
	if err != nil || rows != 0 {
		t.Fatalf("stale record = %d, %v; want 0", rows, err)
	}

	// A future retry is not due yet.
	insertEvent("later", newJob())
	if _, err := q.db.Exec(ctx,
		"UPDATE run_trigger_event SET next_attempt_at = now() + interval '1 hour' WHERE id = 'later'"); err != nil {
		t.Fatalf("defer event: %v", err)
	}
	if _, err := q.ClaimRunTriggerEvent(ctx); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("claim of a deferred event err = %v, want ErrNoRows", err)
	}

	// An abandoned PROCESSING row is reclaimed, or failed once out of attempts.
	abandoned := time.Now().Add(-2 * time.Minute)
	cutoff := pgtype.Timestamptz{Time: time.Now().Add(-time.Minute), Valid: true}
	if _, err := q.db.Exec(ctx,
		"UPDATE run_trigger_event SET status = 'PROCESSING', attempt_count = 1, last_attempt_at = $1 WHERE id = 'later'",
		abandoned); err != nil {
		t.Fatalf("abandon event: %v", err)
	}
	reclaimed, err := q.ReclaimStuckRunTriggerEvents(ctx, ReclaimStuckRunTriggerEventsParams{
		LastAttemptAt: cutoff, AttemptCount: 10,
	})
	if err != nil || reclaimed != 1 {
		t.Fatalf("reclaim = %d, %v; want 1", reclaimed, err)
	}

	if _, err := q.db.Exec(ctx,
		"UPDATE run_trigger_event SET status = 'PROCESSING', attempt_count = 10, last_attempt_at = $1 WHERE id = 'later'",
		abandoned); err != nil {
		t.Fatalf("exhaust event: %v", err)
	}
	message := "crashed mid-dispatch"
	failed, err := q.FailStuckRunTriggerEvents(ctx, FailStuckRunTriggerEventsParams{
		LastAttemptAt: cutoff, AttemptCount: 10, LastError: &message,
	})
	if err != nil || failed != 1 {
		t.Fatalf("fail stuck = %d, %v; want 1", failed, err)
	}
	var status, lastError string
	if err := q.db.QueryRow(ctx,
		"SELECT status, last_error FROM run_trigger_event WHERE id = 'later'").Scan(&status, &lastError); err != nil {
		t.Fatalf("read failed event: %v", err)
	}
	if status != "FAILED" || lastError != message {
		t.Errorf("event = %s/%q, want FAILED/%q", status, lastError, message)
	}
}
