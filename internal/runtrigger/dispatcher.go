// Package runtrigger turns a completed apply into queued runs on the
// destination workspaces of its enabled run triggers.
//
// The outbox row is written by UpdateJobStatus in the same statement that marks
// a job completed, so a crash between completion and dispatch cannot lose it.
// This poller drains those rows: it qualifies the run, fans out to the edges and
// records the outcome, retrying with backoff when a dependent could not be
// queued.
package runtrigger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// JobQueue is the slice of the queue the dispatcher publishes to.
type JobQueue interface {
	EnqueueJob(ctx context.Context, jobID int64) (string, error)
}

// A triggered run is a normal plan, approval and apply.
const (
	planStepNumber     int32 = 100
	approvalStepNumber int32 = 150
	applyStepNumber    int32 = 200
)

const (
	// maxCascadeDepth is the runtime net against a cycle that slipped past the
	// write-time check, mirroring io.terrakube.run-trigger.max-cascade-depth.
	maxCascadeDepth = 10
	// eventMaxAttempts bounds redelivery of one completed job's event.
	eventMaxAttempts = 10
	// eventLease is how long a claimed event may run before the sweep treats the
	// claiming replica as dead.
	eventLease = 60 * time.Second
	// pollerBatch bounds how many events one tick drains.
	pollerBatch = 200
	// eventBackoffInitial and eventBackoffMax bound the retry delay.
	eventBackoffInitial = 5 * time.Second
	eventBackoffMax     = 15 * time.Minute
)

// Dispatcher drains run trigger events. The caller must hold the leader lease.
type Dispatcher struct {
	Store    store.Querier
	Queue    JobQueue
	Logger   *slog.Logger
	Now      func() time.Time
	Interval time.Duration
}

// NewDispatcher builds a dispatcher with the production clock.
func NewDispatcher(st store.Querier, q JobQueue, logger *slog.Logger) *Dispatcher {
	return &Dispatcher{
		Store:    st,
		Queue:    q,
		Logger:   logger,
		Now:      time.Now,
		Interval: 10 * time.Second,
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

// Tick sweeps abandoned events and processes every due one.
func (d *Dispatcher) Tick(ctx context.Context) {
	now := d.now()
	d.sweep(ctx, now)
	for range pollerBatch {
		event, err := d.Store.ClaimRunTriggerEvent(ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			return
		}
		if err != nil {
			if ctx.Err() == nil {
				d.logger().Warn("claim run trigger event", "error", err)
			}
			return
		}
		d.process(ctx, event, now)
	}
}

func (d *Dispatcher) sweep(ctx context.Context, now time.Time) {
	cutoff := pgtype.Timestamptz{Time: now.Add(-eventLease), Valid: true}
	reclaimed, err := d.Store.ReclaimStuckRunTriggerEvents(ctx, store.ReclaimStuckRunTriggerEventsParams{
		LastAttemptAt: cutoff,
		AttemptCount:  eventMaxAttempts,
	})
	if err != nil && ctx.Err() == nil {
		d.logger().Warn("sweep stuck run trigger events", "error", err)
	}
	message := "Processing attempt did not complete within the lease; the claiming instance likely crashed"
	failed, err := d.Store.FailStuckRunTriggerEvents(ctx, store.FailStuckRunTriggerEventsParams{
		LastAttemptAt: cutoff,
		AttemptCount:  eventMaxAttempts,
		LastError:     &message,
	})
	if err != nil && ctx.Err() == nil {
		d.logger().Warn("fail stuck run trigger events", "error", err)
	}
	if reclaimed > 0 || failed > 0 {
		d.logger().Warn("run trigger event sweep", "reclaimed", reclaimed, "failed", failed)
	}
}

// process runs one claimed event and records its outcome. A dispatch failure is
// retried, because a partial fan-out must not mark the event done.
func (d *Dispatcher) process(ctx context.Context, event store.RunTriggerEvent, now time.Time) {
	err := d.dispatch(ctx, event.JobID)

	status := "PROCESSED"
	var lastError *string
	next := pgtype.Timestamptz{}
	if err != nil {
		message := err.Error()
		lastError = &message
		if event.AttemptCount >= eventMaxAttempts {
			status = "FAILED"
			d.logger().Error("run trigger event permanently failed",
				"event_id", event.ID, "job_id", event.JobID,
				"attempts", event.AttemptCount, "error", message)
		} else {
			status = "PENDING"
			next = pgtype.Timestamptz{Time: now.Add(backoff(event.AttemptCount)), Valid: true}
			d.logger().Warn("run trigger event failed; will retry",
				"event_id", event.ID, "job_id", event.JobID,
				"attempts", event.AttemptCount, "retry_at", next.Time, "error", message)
		}
	} else if event.AttemptCount > 1 {
		d.logger().Info("run trigger event processed after a retry",
			"event_id", event.ID, "job_id", event.JobID, "attempts", event.AttemptCount)
	}

	rows, err := d.Store.RecordRunTriggerEvent(ctx, store.RecordRunTriggerEventParams{
		ID:            event.ID,
		LastAttemptAt: event.LastAttemptAt,
		Status:        status,
		LastError:     lastError,
		NextAttemptAt: next,
	})
	if err != nil && ctx.Err() == nil {
		d.logger().Error("record run trigger event result", "event_id", event.ID, "error", err)
	} else if err == nil && rows == 0 {
		// The lease expired and the sweep reclaimed the row while this attempt
		// was running; its result must not overwrite the newer attempt.
		d.logger().Warn("discarded stale run trigger event result", "event_id", event.ID)
	}
}

// dispatch fans out from one completed job. A missing job or an unchanged run
// is a success: there is nothing to retry.
func (d *Dispatcher) dispatch(ctx context.Context, jobID int64) error {
	job, err := d.Store.GetJob(ctx, jobID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get completed job: %w", err)
	}
	if job.Status != "completed" {
		return nil
	}

	changed, err := d.changedState(ctx, job)
	if err != nil || !changed {
		return err
	}

	if job.CascadeDepth+1 > maxCascadeDepth {
		d.logger().Warn("run trigger cascade limit reached; not dispatching further",
			"job_id", job.ID, "depth", job.CascadeDepth, "limit", maxCascadeDepth)
		return nil
	}

	triggers, err := d.Store.ListEnabledRunTriggersBySource(ctx, job.WorkspaceID)
	if err != nil {
		return fmt.Errorf("list run triggers: %w", err)
	}

	// One destination missing a template or failing to queue must not cost the
	// dependents after it, but the event still has to be retried.
	var partial bool
	for _, trigger := range triggers {
		if err := d.enqueue(ctx, trigger, job); err != nil {
			partial = true
			d.logger().Error("could not enqueue triggered run",
				"job_id", job.ID, "trigger_id", trigger.ID,
				"destination_workspace_id", trigger.DestinationWorkspaceID, "error", err)
		}
	}
	if partial {
		return fmt.Errorf("one or more run triggers failed to dispatch for job %d", job.ID)
	}
	return nil
}

// changedState reports whether the run actually wrote remote state.
//
// Our executor only models the plan/approval/apply flow: a run that reached a
// completed apply step changed state, and a plan-only or failed run never has
// one. When the apply produced a history row, its serial and md5 are compared
// with the previous state so a no-change apply does not fan out.
func (d *Dispatcher) changedState(ctx context.Context, job store.Job) (bool, error) {
	steps, err := d.Store.ListStepsByJob(ctx, job.ID)
	if err != nil {
		return false, fmt.Errorf("list job steps: %w", err)
	}
	applied := false
	for _, step := range steps {
		if step.StepNumber == applyStepNumber && step.Status != nil && *step.Status == "completed" {
			applied = true
			break
		}
	}
	if !applied {
		return false, nil
	}

	reference := strconv.FormatInt(job.ID, 10)
	history, err := d.Store.GetHistoryByJobReference(ctx, store.GetHistoryByJobReferenceParams{
		WorkspaceID:  job.WorkspaceID,
		JobReference: &reference,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// No state was archived (for example a destroy), so treat it as changed.
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("get job history: %w", err)
	}

	previous, err := d.Store.GetHistoryBefore(ctx, store.GetHistoryBeforeParams{
		WorkspaceID: job.WorkspaceID,
		CreatedDate: history.CreatedDate,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// First state record for the workspace.
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("get previous history: %w", err)
	}

	// Legacy rows written before real serial/md5 were captured.
	if history.Serial <= 1 && history.Md5 == "0" {
		return true, nil
	}
	return history.Serial > previous.Serial || (history.Md5 != "0" && history.Md5 != previous.Md5), nil
}

// enqueue creates the dependent run, or re-queues one a previous attempt left
// pending. The (destination, triggering job) pair is the idempotency key.
func (d *Dispatcher) enqueue(ctx context.Context, trigger store.WorkspaceRunTrigger, job store.Job) error {
	existing, err := d.Store.GetJobByWorkspaceAndTrigger(ctx, store.GetJobByWorkspaceAndTriggerParams{
		WorkspaceID:      trigger.DestinationWorkspaceID,
		TriggeredByJobID: &job.ID,
	})
	switch {
	case err == nil:
		if existing.Status != "pending" {
			d.logger().Info("triggered run already exists; skipping duplicate",
				"job_id", job.ID, "dependent_job_id", existing.ID, "status", existing.Status)
			return nil
		}
		if err := d.enqueueJob(ctx, existing.ID); err != nil {
			return err
		}
		d.logger().Info("re-queued previously created triggered run",
			"job_id", job.ID, "dependent_job_id", existing.ID)
		return nil
	case !errors.Is(err, pgx.ErrNoRows):
		return fmt.Errorf("look up existing triggered run: %w", err)
	}

	workspace, err := d.Store.GetWorkspace(ctx, trigger.DestinationWorkspaceID)
	if errors.Is(err, pgx.ErrNoRows) {
		d.logger().Warn("run trigger destination is gone; skipping",
			"trigger_id", trigger.ID, "destination_workspace_id", trigger.DestinationWorkspaceID)
		return nil
	}
	if err != nil {
		return fmt.Errorf("get destination workspace: %w", err)
	}

	templateID := trigger.TemplateID
	if templateID == nil || *templateID == "" {
		templateID = workspace.DefaultTemplate
	}
	if templateID == nil || *templateID == "" {
		d.logger().Warn("run trigger destination has no template; skipping",
			"trigger_id", trigger.ID, "workspace_id", workspace.ID)
		return nil
	}
	template, err := d.Store.GetTemplate(ctx, *templateID)
	if errors.Is(err, pgx.ErrNoRows) {
		d.logger().Warn("run trigger template is gone; skipping",
			"trigger_id", trigger.ID, "template_id", *templateID)
		return nil
	}
	if err != nil {
		return fmt.Errorf("get trigger template: %w", err)
	}
	if template.OrganizationID != workspace.OrganizationID {
		d.logger().Warn("run trigger template belongs to another organization; skipping",
			"trigger_id", trigger.ID, "template_id", template.ID)
		return nil
	}

	depth := job.CascadeDepth + 1
	created, err := d.Store.CreateJob(ctx, store.CreateJobParams{
		OrganizationID:    workspace.OrganizationID,
		WorkspaceID:       workspace.ID,
		Status:            "pending",
		Via:               "RUN_TRIGGER",
		Refresh:           true,
		PlanChanges:       true,
		CreatedBy:         strptr("serviceAccount"),
		Tcl:               template.Tcl,
		TemplateReference: templateID,
		TriggeredByJobID:  &job.ID,
		CascadeDepth:      depth,
	})
	if err != nil {
		return fmt.Errorf("create triggered run: %w", err)
	}

	for _, step := range []struct {
		number int32
		name   string
	}{
		{planStepNumber, "Terraform Plan"},
		{approvalStepNumber, "Approve Plan"},
		{applyStepNumber, "Terraform Apply"},
	} {
		if _, err := d.Store.CreateStep(ctx, store.CreateStepParams{
			ID:         uuid.NewString(),
			JobID:      created.ID,
			StepNumber: step.number,
			Name:       strptr(step.name),
			Status:     strptr("pending"),
			LogStatus:  strptr("pending"),
		}); err != nil {
			// The executor would treat a partial flow as complete, so remove the
			// row and let the event retry instead.
			_ = d.Store.DeleteJob(ctx, created.ID)
			return fmt.Errorf("create triggered run step: %w", err)
		}
	}

	_ = d.Store.SetWorkspaceLastJob(ctx, store.SetWorkspaceLastJobParams{
		ID:            workspace.ID,
		LastJobStatus: "pending",
	})
	if err := d.enqueueJob(ctx, created.ID); err != nil {
		// The row stays pending; a retried event finds and queues it again.
		return err
	}
	d.logger().Info("queued run triggered by completed job",
		"job_id", job.ID, "dependent_job_id", created.ID,
		"workspace_id", workspace.ID, "depth", depth)
	return nil
}

func (d *Dispatcher) enqueueJob(ctx context.Context, jobID int64) error {
	if d.Queue == nil {
		return errors.New("run queue is not configured")
	}
	if _, err := d.Queue.EnqueueJob(ctx, jobID); err != nil {
		return fmt.Errorf("enqueue triggered run %d: %w", jobID, err)
	}
	return nil
}

// backoff is exponential without jitter: events are rare enough that a
// correlated retry is harmless, and a deterministic delay is testable.
func backoff(attempt int32) time.Duration {
	delay := eventBackoffInitial
	for i := int32(1); i < attempt && delay < eventBackoffMax; i++ {
		delay *= 2
	}
	if delay > eventBackoffMax {
		return eventBackoffMax
	}
	return delay
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

func strptr(v string) *string { return &v }
