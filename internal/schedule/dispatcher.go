package schedule

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// JobQueue is the slice of the queue the dispatcher publishes to.
type JobQueue interface {
	EnqueueJob(ctx context.Context, jobID int64) (string, error)
}

// Step numbers mirror the executor's; a scheduled run is a normal plan,
// approval and apply.
const (
	planStepNumber     int32 = 100
	approvalStepNumber int32 = 150
	applyStepNumber    int32 = 200
)

const (
	// maxScanWindow bounds the catch-up scan after the process was paused, so
	// a long sleep cannot fire hours of runs at once.
	maxScanWindow = 2 * time.Minute
	// firedRetention keeps the duplicate guard small.
	firedRetention = 24 * time.Hour
)

// Dispatcher turns matching cron ticks into queued runs. It scans the seconds
// elapsed since the previous tick, so a cron whose seconds field is not zero
// still fires without depending on the ticker's phase.
type Dispatcher struct {
	Store    store.Querier
	Queue    JobQueue
	Logger   *slog.Logger
	Now      func() time.Time
	Interval time.Duration

	lastTick time.Time
	fired    map[string]time.Time
	warned   map[string]string
}

// NewDispatcher builds a dispatcher with the production clock.
func NewDispatcher(st store.Querier, q JobQueue, logger *slog.Logger) *Dispatcher {
	return &Dispatcher{
		Store:    st,
		Queue:    q,
		Logger:   logger,
		Now:      time.Now,
		Interval: 20 * time.Second,
		fired:    map[string]time.Time{},
		warned:   map[string]string{},
	}
}

// Run ticks until ctx is cancelled. The caller must hold the leader lease.
func (d *Dispatcher) Run(ctx context.Context) {
	interval := d.Interval
	if interval <= 0 {
		interval = 20 * time.Second
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

// Tick fires every schedule whose expression matches a second since the
// previous tick. The first tick only records the clock: replaying the previous
// interval would launch runs on every restart.
func (d *Dispatcher) Tick(ctx context.Context) {
	if d.Now == nil {
		d.Now = time.Now
	}
	now := d.Now().Truncate(time.Second)
	from := d.lastTick
	d.lastTick = now
	if from.IsZero() || !now.After(from) {
		return
	}
	if now.Sub(from) > maxScanWindow {
		from = now.Add(-maxScanWindow)
	}

	schedules, err := d.Store.ListEnabledSchedules(ctx)
	if err != nil {
		if ctx.Err() == nil {
			d.logger().Warn("list enabled schedules", "error", err)
		}
		return
	}
	for _, schedule := range schedules {
		if ctx.Err() != nil {
			return
		}
		if schedule.Cron == nil || *schedule.Cron == "" {
			continue
		}
		cron, err := Parse(*schedule.Cron)
		if err != nil {
			d.warnOnce(schedule.ID, *schedule.Cron, err)
			continue
		}
		for at := from.Add(time.Second); !at.After(now); at = at.Add(time.Second) {
			if !cron.Matches(at) {
				continue
			}
			if last, ok := d.fired[schedule.ID]; ok && last.Equal(at) {
				continue
			}
			d.fired[schedule.ID] = at
			if err := d.fire(ctx, schedule); err != nil && ctx.Err() == nil {
				d.logger().Warn("dispatch schedule", "schedule_id", schedule.ID,
					"cron", cron.String(), "error", err)
			}
		}
	}
	d.pruneFired(now)
}

// fire creates and queues the job for one tick. A locked workspace and a
// template that has left the organization are skipped, as upstream's trigger
// does; anything else is reported and left for the next tick.
func (d *Dispatcher) fire(ctx context.Context, schedule store.Schedule) error {
	ws, err := d.Store.GetWorkspace(ctx, schedule.WorkspaceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			d.logger().Warn("schedule workspace is gone", "schedule_id", schedule.ID)
			return nil
		}
		return fmt.Errorf("get workspace: %w", err)
	}
	if ws.Locked {
		d.logger().Info("workspace is locked; skipping scheduled run",
			"schedule_id", schedule.ID, "workspace_id", ws.ID)
		return nil
	}

	var (
		tcl               *string
		templateReference *string
	)
	if schedule.TemplateReference != nil && *schedule.TemplateReference != "" {
		template, err := d.Store.GetTemplate(ctx, *schedule.TemplateReference)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				d.logger().Warn("schedule template is gone", "schedule_id", schedule.ID)
				return nil
			}
			return fmt.Errorf("get template: %w", err)
		}
		if template.OrganizationID != ws.OrganizationID {
			d.logger().Warn("schedule template belongs to another organization",
				"schedule_id", schedule.ID, "template_id", template.ID)
			return nil
		}
		tcl = template.Tcl
		templateReference = schedule.TemplateReference
	} else {
		tcl = schedule.Tcl
	}

	job, err := d.Store.CreateJob(ctx, store.CreateJobParams{
		OrganizationID:    ws.OrganizationID,
		WorkspaceID:       ws.ID,
		Status:            "pending",
		Via:               "SCHEDULE",
		Refresh:           true,
		PlanChanges:       true,
		CreatedBy:         strptr("serviceAccount"),
		Tcl:               tcl,
		TemplateReference: templateReference,
	})
	if err != nil {
		return fmt.Errorf("create job: %w", err)
	}

	steps := []struct {
		number int32
		name   string
	}{
		{planStepNumber, "Terraform Plan"},
		{approvalStepNumber, "Approve Plan"},
		{applyStepNumber, "Terraform Apply"},
	}
	for _, step := range steps {
		if _, err := d.Store.CreateStep(ctx, store.CreateStepParams{
			ID:         uuid.NewString(),
			JobID:      job.ID,
			StepNumber: step.number,
			Name:       strptr(step.name),
			Status:     strptr("pending"),
			LogStatus:  strptr("pending"),
		}); err != nil {
			d.failJob(ctx, job.ID)
			return fmt.Errorf("create step: %w", err)
		}
	}
	_ = d.Store.SetWorkspaceLastJob(ctx, store.SetWorkspaceLastJobParams{
		ID:            ws.ID,
		LastJobStatus: "pending",
	})

	if d.Queue == nil {
		d.failJob(ctx, job.ID)
		return errors.New("run queue is not configured")
	}
	if _, err := d.Queue.EnqueueJob(ctx, job.ID); err != nil {
		d.failJob(ctx, job.ID)
		return fmt.Errorf("enqueue job: %w", err)
	}
	d.logger().Info("scheduled run queued", "schedule_id", schedule.ID,
		"workspace_id", ws.ID, "job_id", job.ID)
	return nil
}

// failJob marks a run that could not be dispatched so it never sits pending.
func (d *Dispatcher) failJob(ctx context.Context, jobID int64) {
	if _, err := d.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{ID: jobID, Status: "failed"}); err != nil {
		d.logger().Error("mark scheduled job failed", "job_id", jobID, "error", err)
	}
}

// warnOnce logs an unparseable expression when it changes, not on every tick.
func (d *Dispatcher) warnOnce(scheduleID, expression string, err error) {
	if d.warned[scheduleID] == expression {
		return
	}
	d.warned[scheduleID] = expression
	d.logger().Warn("schedule cron is not supported; it will not fire",
		"schedule_id", scheduleID, "cron", expression, "error", err)
}

func (d *Dispatcher) pruneFired(now time.Time) {
	if len(d.fired) == 0 {
		return
	}
	for id, at := range d.fired {
		if now.Sub(at) > firedRetention {
			delete(d.fired, id)
		}
	}
}

func (d *Dispatcher) logger() *slog.Logger {
	if d.Logger == nil {
		return slog.Default()
	}
	return d.Logger
}

func strptr(v string) *string { return &v }
