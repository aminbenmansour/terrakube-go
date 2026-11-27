package notification

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// The dispatcher drains notification_outbox, mirroring upstream's
// NotificationOutboxPollerJob and NotificationDispatchService: claim one row,
// deliver it outside any transaction, then record the outcome. Claims use
// SELECT ... FOR UPDATE SKIP LOCKED, so every API replica can run the poller
// without a leader election.
const (
	// maxAttempts bounds redelivery of one outbox row (upstream's MAX_ATTEMPTS).
	maxAttempts = 3
	// attemptLease is how long a claimed row may run before the sweep treats the
	// claiming replica as dead.
	attemptLease = 5 * time.Minute
	// pollerBatch bounds how many rows one tick drains.
	pollerBatch = 200
	// failureSummaryLength is the tail of the failing step's output kept for a
	// failure payload, matching JobFailureSummaryService.
	failureSummaryLength = 500
)

// Dispatcher delivers the rows UpdateJobStatus enqueues.
type Dispatcher struct {
	Store     store.Querier
	Deliverer *Deliverer
	PublicURL string
	Logger    *slog.Logger
	Now       func() time.Time
	Interval  time.Duration
}

// NewDispatcher builds a dispatcher with the production clock.
func NewDispatcher(st store.Querier, deliverer *Deliverer, publicURL string, logger *slog.Logger) *Dispatcher {
	return &Dispatcher{
		Store:     st,
		Deliverer: deliverer,
		PublicURL: strings.TrimRight(publicURL, "/"),
		Logger:    logger,
		Now:       time.Now,
		Interval:  10 * time.Second,
	}
}

// Run ticks until ctx is cancelled.
func (d *Dispatcher) Run(ctx context.Context) {
	interval := d.Interval
	if interval <= 0 {
		interval = 10 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	d.Tick(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.Tick(ctx)
		}
	}
}

// Tick sweeps abandoned rows and drains every due one.
func (d *Dispatcher) Tick(ctx context.Context) {
	now := d.now()
	d.sweep(ctx, now)
	for range pollerBatch {
		row, err := d.Store.ClaimNotificationOutbox(ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			return
		}
		if err != nil {
			if ctx.Err() == nil {
				d.logger().Warn("claim notification outbox", "error", err)
			}
			return
		}
		d.process(ctx, row)
	}
}

func (d *Dispatcher) sweep(ctx context.Context, now time.Time) {
	cutoff := pgtype.Timestamptz{Time: now.Add(-attemptLease), Valid: true}
	reclaimed, err := d.Store.ReclaimStuckNotificationOutbox(ctx, store.ReclaimStuckNotificationOutboxParams{
		LastAttemptAt: cutoff,
		AttemptCount:  maxAttempts,
	})
	if err != nil && ctx.Err() == nil {
		d.logger().Warn("sweep stuck notification outbox", "error", err)
	}
	message := "Delivery attempt did not complete within the lease; the claiming instance likely crashed"
	failed, err := d.Store.FailStuckNotificationOutbox(ctx, store.FailStuckNotificationOutboxParams{
		LastAttemptAt: cutoff,
		AttemptCount:  maxAttempts,
		LastError:     &message,
	})
	if err != nil && ctx.Err() == nil {
		d.logger().Warn("fail stuck notification outbox", "error", err)
	}
	if reclaimed > 0 || failed > 0 {
		d.logger().Warn("notification outbox sweep", "reclaimed", reclaimed, "failed", failed)
	}
}

// process delivers one claimed row and records the outcome. The row is already
// SENDING and its attempt counted, so every path records a result.
func (d *Dispatcher) process(ctx context.Context, row store.NotificationOutbox) {
	status, lastError, next := d.deliver(ctx, row)

	rows, err := d.Store.RecordNotificationOutbox(ctx, store.RecordNotificationOutboxParams{
		ID:            row.ID,
		LastAttemptAt: row.LastAttemptAt,
		Status:        status,
		LastError:     lastError,
		NextAttemptAt: next,
	})
	if err != nil && ctx.Err() == nil {
		d.logger().Error("record notification outbox result", "outbox_id", row.ID, "error", err)
	} else if err == nil && rows == 0 {
		// The lease expired and the sweep reclaimed the row while the HTTP call
		// was running; its result must not overwrite the newer attempt.
		d.logger().Warn("discarded stale notification delivery result", "outbox_id", row.ID)
	}
}

// deliver returns the terminal or retryable state for one attempt.
func (d *Dispatcher) deliver(ctx context.Context, row store.NotificationOutbox) (string, *string, pgtype.Timestamptz) {
	job, err := d.Store.GetJob(ctx, row.JobID)
	if errors.Is(err, pgx.ErrNoRows) {
		return d.fail(row, "the job no longer exists")
	}
	if err != nil {
		return d.retryOrFail(ctx, row, &DeliveryError{Message: "read job: " + err.Error(), Retryable: true})
	}

	// A transient "failed" usually means the API lost the executor's dispatch
	// response; the executor's own callback then corrects the status. Never
	// send a notification that no longer describes the run.
	if row.JobStatus != nil && job.Status != *row.JobStatus {
		message := fmt.Sprintf("Job status changed from %s to %s before delivery; suppressed as a false alarm", *row.JobStatus, job.Status)
		return "SKIPPED", &message, pgtype.Timestamptz{}
	}

	config, err := d.Store.GetNotificationConfiguration(ctx, row.ConfigurationID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "SKIPPED", strPtr("the notification configuration was deleted before delivery"), pgtype.Timestamptz{}
	}
	if err != nil {
		return d.retryOrFail(ctx, row, &DeliveryError{Message: "read notification configuration: " + err.Error(), Retryable: true})
	}

	renderContext, err := d.renderContext(ctx, job, config)
	if err != nil {
		return d.fail(row, err.Error())
	}
	payload, err := BuildPayload(config.ChannelType, renderContext)
	if err != nil {
		return d.fail(row, err.Error())
	}

	err = d.Deliverer.Deliver(ctx, Configuration{
		ChannelType:    config.ChannelType,
		DestinationURL: config.DestinationUrl,
		SigningSecret:  config.SigningSecret,
	}, payload)
	if err == nil {
		if row.AttemptCount > 1 {
			d.logger().Info("notification delivered after a retry", "outbox_id", row.ID, "attempts", row.AttemptCount)
		}
		return "SENT", nil, pgtype.Timestamptz{}
	}

	var deliveryErr *DeliveryError
	if errors.As(err, &deliveryErr) {
		return d.retryOrFail(ctx, row, deliveryErr)
	}
	// An unclassified error (e.g. a payload bug) would retry forever.
	return d.fail(row, err.Error())
}

func (d *Dispatcher) retryOrFail(ctx context.Context, row store.NotificationOutbox, err *DeliveryError) (string, *string, pgtype.Timestamptz) {
	message := err.Error()
	if !err.Retryable || row.AttemptCount >= maxAttempts {
		d.logger().Warn("notification permanently failed",
			"outbox_id", row.ID, "attempts", row.AttemptCount, "error", message)
		return "FAILED", &message, pgtype.Timestamptz{}
	}
	delay := err.RetryAfter
	if delay <= 0 {
		delay = backoff(row.AttemptCount)
	}
	next := pgtype.Timestamptz{Time: d.now().Add(delay), Valid: true}
	d.logger().Warn("notification failed; will retry",
		"outbox_id", row.ID, "attempts", row.AttemptCount, "retry_at", next.Time, "error", message)
	return "PENDING", &message, next
}

func (d *Dispatcher) fail(row store.NotificationOutbox, message string) (string, *string, pgtype.Timestamptz) {
	d.logger().Warn("notification permanently failed", "outbox_id", row.ID, "error", message)
	return "FAILED", &message, pgtype.Timestamptz{}
}

// renderContext assembles the data a payload builder needs. The URLs point at
// the UI, which is what upstream's io.terrakube.ui.url does.
func (d *Dispatcher) renderContext(ctx context.Context, job store.Job, config store.NotificationConfiguration) (Context, error) {
	workspace, err := d.Store.GetWorkspace(ctx, job.WorkspaceID)
	if err != nil {
		return Context{}, fmt.Errorf("read workspace: %w", err)
	}
	organization, err := d.Store.GetOrganization(ctx, job.OrganizationID)
	if err != nil {
		return Context{}, fmt.Errorf("read organization: %w", err)
	}

	workspaceURL := fmt.Sprintf("%s/organizations/%s/workspaces/%s", d.PublicURL, organization.ID, workspace.ID)
	renderContext := Context{
		OrganizationName:  organization.Name,
		WorkspaceName:     workspace.Name,
		JobID:             job.ID,
		Status:            job.Status,
		RunURL:            fmt.Sprintf("%s/runs/%d", workspaceURL, job.ID),
		ConfigurationName: config.Name,
		WorkspaceURL:      workspaceURL,
		MessageStyle:      config.MessageStyle,
	}
	if job.CommitID != nil {
		renderContext.CommitID = *job.CommitID
	}
	if job.Status == "failed" {
		renderContext.FailureReason = d.failureSummary(ctx, job)
	}
	return renderContext, nil
}

// failureSummary keeps the tail of the failing step's console output, the same
// material JobFailureSummaryService reads. A job can reach failed before any
// step is marked failed, so the highest-numbered step is the fallback.
func (d *Dispatcher) failureSummary(ctx context.Context, job store.Job) string {
	steps, err := d.Store.ListStepsByJob(ctx, job.ID)
	if err != nil {
		d.logger().Warn("could not read job steps for the failure summary", "job_id", job.ID, "error", err)
		return ""
	}
	var failed, last *store.Step
	for i := range steps {
		step := &steps[i]
		if last == nil || step.StepNumber > last.StepNumber {
			last = step
		}
		if step.Status != nil && *step.Status == "failed" &&
			(failed == nil || step.StepNumber > failed.StepNumber) {
			failed = step
		}
	}
	chosen := failed
	if chosen == nil {
		chosen = last
	}
	if chosen == nil || chosen.Output == nil {
		return ""
	}
	output := strings.TrimSpace(*chosen.Output)
	if utf8.RuneCountInString(output) > failureSummaryLength {
		runes := []rune(output)
		output = strings.TrimSpace(string(runes[len(runes)-failureSummaryLength:]))
	}
	return output
}

// backoff mirrors upstream: one minute after the first failure, five after
// every later one. A server-supplied Retry-After takes precedence.
func backoff(attempt int32) time.Duration {
	if attempt <= 1 {
		return time.Minute
	}
	return 5 * time.Minute
}

func (d *Dispatcher) now() time.Time {
	if d.Now == nil {
		return time.Now()
	}
	return d.Now()
}

func (d *Dispatcher) logger() *slog.Logger {
	if d.Logger == nil {
		return slog.Default()
	}
	return d.Logger
}

func strPtr(v string) *string { return &v }
