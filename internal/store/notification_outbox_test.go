package store

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
)

// TestUpdateJobStatusEnqueuesNotifications covers the atomic enqueue: one row
// per active in-scope configuration whose trigger names the new status, the
// template filter narrows it, a "failed" row is held briefly, and a repeated
// status is idempotent.
func TestUpdateJobStatusEnqueuesNotifications(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "notify-enqueue", ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}
	workspace, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
		ID: uuid.NewString(), Name: "notify", IacType: "terraform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	template, err := q.CreateTemplate(ctx, CreateTemplateParams{
		ID: uuid.NewString(), OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create template: %v", err)
	}

	newConfig := func(name string, workspaceScoped, active bool, status string, templates ...string) NotificationConfiguration {
		t.Helper()
		var workspaceID *string
		if workspaceScoped {
			workspaceID = &workspace.ID
		}
		config, err := q.CreateNotificationConfiguration(ctx, CreateNotificationConfigurationParams{
			ID: uuid.NewString(), Name: name, OrganizationID: org.ID, WorkspaceID: workspaceID,
			ChannelType: "WEBHOOK", DestinationUrl: "https://example.com/hook",
			Active: active, MessageStyle: "DETAILED",
		})
		if err != nil {
			t.Fatalf("create configuration %s: %v", name, err)
		}
		if _, err := q.CreateNotificationTrigger(ctx, CreateNotificationTriggerParams{
			ID: uuid.NewString(), JobStatus: status, ConfigurationID: config.ID,
		}); err != nil {
			t.Fatalf("create trigger %s: %v", name, err)
		}
		for _, templateID := range templates {
			if err := q.AddNotificationConfigurationTemplate(ctx, AddNotificationConfigurationTemplateParams{
				ConfigurationID: config.ID, TemplateID: templateID,
			}); err != nil {
				t.Fatalf("add template %s: %v", name, err)
			}
		}
		return config
	}

	orgCompleted := newConfig("org-completed", false, true, "completed")
	workspaceFailed := newConfig("workspace-failed", true, true, "failed")
	newConfig("inactive", false, false, "completed")
	templateBound := newConfig("template-bound", false, true, "completed", template.ID)
	newConfig("other-status", false, true, "cancelled")

	newJob := func(reference *string) Job {
		t.Helper()
		job, err := q.CreateJob(ctx, CreateJobParams{
			OrganizationID: org.ID, WorkspaceID: workspace.ID, Status: "pending", Via: "UI",
			Refresh: true, PlanChanges: true, TemplateReference: reference,
		})
		if err != nil {
			t.Fatalf("create job: %v", err)
		}
		return job
	}
	rowsFor := func(jobID int64) []NotificationOutbox {
		t.Helper()
		rows, err := q.db.Query(ctx,
			"SELECT * FROM notification_outbox WHERE job_id = $1 ORDER BY configuration_id", jobID)
		if err != nil {
			t.Fatalf("list outbox: %v", err)
		}
		defer rows.Close()
		var found []NotificationOutbox
		for rows.Next() {
			var row NotificationOutbox
			if err := rows.Scan(&row.ID, &row.JobID, &row.ConfigurationID, &row.Payload, &row.Status,
				&row.JobStatus, &row.AttemptCount, &row.LastAttemptAt, &row.NextAttemptAt, &row.LastError,
				&row.CreatedDate, &row.UpdatedDate, &row.CreatedBy, &row.UpdatedBy); err != nil {
				t.Fatalf("scan outbox: %v", err)
			}
			found = append(found, row)
		}
		return found
	}

	// A completed run reaches the org-wide config and the config whose template
	// filter matches it; the inactive and other-status configs stay out.
	job := newJob(&template.ID)
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: job.ID, Status: "completed"}); err != nil {
		t.Fatalf("complete job: %v", err)
	}
	rows := rowsFor(job.ID)
	if len(rows) != 2 {
		t.Fatalf("outbox rows after completion = %d, want 2", len(rows))
	}
	seen := map[string]bool{rows[0].ConfigurationID: true, rows[1].ConfigurationID: true}
	if !seen[orgCompleted.ID] || !seen[templateBound.ID] {
		t.Fatalf("unexpected configurations enqueued: %v", seen)
	}
	if rows[0].JobStatus == nil || *rows[0].JobStatus != "completed" {
		t.Fatalf("job status = %v, want completed", rows[0].JobStatus)
	}
	if rows[0].NextAttemptAt.Valid {
		t.Fatalf("completed notification is delayed until %v", rows[0].NextAttemptAt.Time)
	}

	// A second identical status must not duplicate the delivery.
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: job.ID, Status: "completed"}); err != nil {
		t.Fatalf("complete job again: %v", err)
	}
	if got := len(rowsFor(job.ID)); got != 2 {
		t.Fatalf("outbox rows after repeated completion = %d, want 2", got)
	}

	// A failed run uses the workspace-scoped config and is held for the grace
	// window so a corrected status can suppress it.
	failing := newJob(nil)
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: failing.ID, Status: "failed"}); err != nil {
		t.Fatalf("fail job: %v", err)
	}
	rows = rowsFor(failing.ID)
	if len(rows) != 1 || rows[0].ConfigurationID != workspaceFailed.ID {
		t.Fatalf("failed outbox = %+v, want one row for the workspace config", rows)
	}
	if !rows[0].NextAttemptAt.Valid || rows[0].NextAttemptAt.Time.Before(time.Now().Add(30*time.Second)) {
		t.Fatalf("failed notification is not held for the grace window: %v", rows[0].NextAttemptAt)
	}

	// A config whose template filter names another template does not fire.
	other := newJob(pgtemplate())
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: other.ID, Status: "completed"}); err != nil {
		t.Fatalf("complete untemplated job: %v", err)
	}
	rows = rowsFor(other.ID)
	if len(rows) != 1 || rows[0].ConfigurationID != orgCompleted.ID {
		t.Fatalf("outbox for an unmatched template = %+v, want only the org-wide config", rows)
	}
}

// TestNotificationOutboxClaimRecordAndSweep covers the poller's SQL: one due
// row at a time, a stale result is discarded, retries reclaim abandoned rows,
// retry re-arms only failed rows, and the workspace read is scoped.
func TestNotificationOutboxClaimRecordAndSweep(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "notify-claim", ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}
	workspace, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
		ID: uuid.NewString(), Name: "notify", IacType: "terraform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create workspace: %v", err)
	}
	config, err := q.CreateNotificationConfiguration(ctx, CreateNotificationConfigurationParams{
		ID: uuid.NewString(), Name: "alerts", OrganizationID: org.ID, WorkspaceID: &workspace.ID,
		ChannelType: "WEBHOOK", DestinationUrl: "https://example.com/hook",
		Active: true, MessageStyle: "DETAILED",
	})
	if err != nil {
		t.Fatalf("create configuration: %v", err)
	}
	if _, err := q.CreateNotificationTrigger(ctx, CreateNotificationTriggerParams{
		ID: uuid.NewString(), JobStatus: "completed", ConfigurationID: config.ID,
	}); err != nil {
		t.Fatalf("create trigger: %v", err)
	}
	job, err := q.CreateJob(ctx, CreateJobParams{
		OrganizationID: org.ID, WorkspaceID: workspace.ID, Status: "pending", Via: "UI",
		Refresh: true, PlanChanges: true,
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}
	if _, err := q.UpdateJobStatus(ctx, UpdateJobStatusParams{ID: job.ID, Status: "completed"}); err != nil {
		t.Fatalf("complete job: %v", err)
	}

	row, err := q.ClaimNotificationOutbox(ctx)
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if row.Status != "SENDING" || row.AttemptCount != 1 || !row.LastAttemptAt.Valid {
		t.Fatalf("claimed row = %+v, want SENDING with one attempt", row)
	}
	if _, err := q.ClaimNotificationOutbox(ctx); err == nil {
		t.Fatal("claimed a row that was already in flight")
	}

	// A result keyed on a stale attempt is discarded.
	stale := row.LastAttemptAt.Time.Add(-time.Minute)
	rows, err := q.RecordNotificationOutbox(ctx, RecordNotificationOutboxParams{
		ID: row.ID, LastAttemptAt: pgtypeTime(stale), Status: "SENT",
	})
	if err != nil || rows != 0 {
		t.Fatalf("stale result rows = %d, err = %v; want 0", rows, err)
	}

	message := "boom"
	next := pgtypeTime(time.Now().Add(time.Minute))
	rows, err = q.RecordNotificationOutbox(ctx, RecordNotificationOutboxParams{
		ID: row.ID, LastAttemptAt: row.LastAttemptAt, Status: "PENDING",
		LastError: &message, NextAttemptAt: next,
	})
	if err != nil || rows != 1 {
		t.Fatalf("record retry rows = %d, err = %v; want 1", rows, err)
	}
	if _, err := q.ClaimNotificationOutbox(ctx); err == nil {
		t.Fatal("claimed a row before its retry time")
	}

	// The sweep reclaims a row whose attempt outlived its lease, and fails one
	// at the attempt limit.
	if _, err := q.db.Exec(ctx,
		"UPDATE notification_outbox SET status = 'SENDING', next_attempt_at = NULL, last_attempt_at = now() - interval '10 minutes' WHERE id = $1", row.ID); err != nil {
		t.Fatalf("age the row: %v", err)
	}
	reclaimed, err := q.ReclaimStuckNotificationOutbox(ctx, ReclaimStuckNotificationOutboxParams{
		LastAttemptAt: pgtypeTime(time.Now().Add(-5 * time.Minute)), AttemptCount: 3,
	})
	if err != nil || reclaimed != 1 {
		t.Fatalf("reclaimed = %d, err = %v; want 1", reclaimed, err)
	}
	if _, err := q.db.Exec(ctx,
		"UPDATE notification_outbox SET status = 'SENDING', attempt_count = 3, last_attempt_at = now() - interval '10 minutes' WHERE id = $1", row.ID); err != nil {
		t.Fatalf("exhaust the row: %v", err)
	}
	failed, err := q.FailStuckNotificationOutbox(ctx, FailStuckNotificationOutboxParams{
		LastAttemptAt: pgtypeTime(time.Now().Add(-5 * time.Minute)), AttemptCount: 3,
		LastError: &message,
	})
	if err != nil || failed != 1 {
		t.Fatalf("failed = %d, err = %v; want 1", failed, err)
	}

	// Retry re-arms a failed row only; the workspace read stays scoped.
	exists, err := q.NotificationOutboxExistsForWorkspace(ctx, NotificationOutboxExistsForWorkspaceParams{
		ID: row.ID, WorkspaceID: workspace.ID,
	})
	if err != nil || !exists {
		t.Fatalf("exists = %v, err = %v; want true", exists, err)
	}
	exists, err = q.NotificationOutboxExistsForWorkspace(ctx, NotificationOutboxExistsForWorkspaceParams{
		ID: row.ID, WorkspaceID: uuid.NewString(),
	})
	if err != nil || exists {
		t.Fatalf("exists for another workspace = %v, err = %v; want false", exists, err)
	}
	rearmed, err := q.RearmNotificationOutbox(ctx, row.ID)
	if err != nil || rearmed != 1 {
		t.Fatalf("rearmed = %d, err = %v; want 1", rearmed, err)
	}
	rearmed, err = q.RearmNotificationOutbox(ctx, row.ID)
	if err != nil || rearmed != 0 {
		t.Fatalf("rearmed a non-failed row = %d, err = %v; want 0", rearmed, err)
	}

	deliveries, err := q.ListNotificationDeliveriesByWorkspace(ctx, ListNotificationDeliveriesByWorkspaceParams{
		WorkspaceID: workspace.ID, Limit: 10,
	})
	if err != nil || len(deliveries) != 1 {
		t.Fatalf("deliveries = %+v, err = %v; want 1", deliveries, err)
	}
	if deliveries[0].ConfigurationName != "alerts" || deliveries[0].ChannelType != "WEBHOOK" {
		t.Fatalf("delivery view = %+v, want the configuration name and channel", deliveries[0])
	}
}

func pgtypeTime(t time.Time) pgtype.Timestamptz { return pgtype.Timestamptz{Time: t, Valid: true} }

func pgtemplate() *string { id := uuid.NewString(); return &id }
