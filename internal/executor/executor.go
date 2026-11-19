// Package executor consumes Terraform/OpenTofu jobs from the Redis queue and
// runs them against the API state backend. It implements the minimal M3
// execution slice: configuration download, workspace variables, init, plan,
// saved-plan apply, live logs, final logs, state history and job status.
package executor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

const (
	// Group is the Redis consumer group the executor uses on the job stream.
	Group = "terrakube-executor"

	planStepNumber     int32 = 100
	approvalStepNumber int32 = 150
	applyStepNumber    int32 = 200

	maxLogStreamEntries int64 = 200000

	// maxJobAttempts bounds how often a transient failure is retried before
	// the job is marked failed. Attempts are carried in the queue message.
	maxJobAttempts = 3
)

// errJobCancelled reports that the job was cancelled while it was running. It
// is not a failure: the API already stored the cancelled status.
var errJobCancelled = errors.New("job cancelled")

// JobLocker takes an exclusive lease on a job so two executors cannot run the
// same job at the same time.
type JobLocker interface {
	TryLock(ctx context.Context, jobID int64) (release func(), ok bool, err error)
}

// Deps are the collaborators the executor needs.
type Deps struct {
	Store               store.Querier
	Storage             storage.Backend
	Queue               JobQueue
	Locker              JobLocker
	Layout              tfstate.Layout
	Logger              *slog.Logger
	PublicURL           string
	InternalTokenSecret string
	Binary              string
	WorkDir             string
	// JobTimeout bounds a single job, including configuration acquisition and
	// every Terraform command. Zero disables the timeout.
	JobTimeout time.Duration
}

// Executor processes jobs from the Redis stream.
type Executor struct {
	deps     Deps
	consumer string

	// retryDelay is the pause before republishing a transiently failed job.
	// It is a field so tests can run without sleeping.
	retryDelay func(attempt int) time.Duration
	// cancelPollInterval is how often a running job's status is checked for
	// cancellation.
	cancelPollInterval time.Duration
	// claimInterval is how often the executor looks for messages left behind
	// by a dead executor.
	claimInterval time.Duration
}

// New builds an executor.
func New(deps Deps) *Executor {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	host, _ := os.Hostname()
	return &Executor{
		deps:               deps,
		consumer:           fmt.Sprintf("%s-%d-%s", host, os.Getpid(), uuid.NewString()[:8]),
		retryDelay:         defaultRetryDelay,
		cancelPollInterval: 2 * time.Second,
		claimInterval:      time.Minute,
	}
}

func defaultRetryDelay(attempt int) time.Duration {
	delay := time.Second << attempt
	if delay > 8*time.Second {
		delay = 8 * time.Second
	}
	return delay
}

// Run consumes jobs until the context is cancelled.
func (e *Executor) Run(ctx context.Context) error {
	if runningAsRoot() {
		e.deps.Logger.Warn("executor is running as root; run it as an unprivileged user so a job cannot modify the host")
	}
	if err := cleanupStaleWorkspaces(e.deps.WorkDir, e.staleWorkspaceAge()); err != nil {
		e.deps.Logger.Warn("could not clean stale job workspaces", "error", err)
	}

	if err := e.deps.Queue.EnsureGroup(ctx, queue.JobsStream, Group); err != nil {
		return err
	}
	e.deps.Logger.Info("executor consuming jobs", "stream", queue.JobsStream, "group", Group, "consumer", e.consumer)

	// Pick up jobs that were delivered before a previous executor stopped.
	e.reclaimStale(ctx)
	lastClaim := time.Now()

	for {
		if time.Since(lastClaim) >= e.claimInterval {
			e.reclaimStale(ctx)
			lastClaim = time.Now()
		}
		messages, err := e.deps.Queue.ReadGroup(ctx, queue.JobsStream, Group, e.consumer, 1, time.Second)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			e.deps.Logger.Error("read job queue failed", "error", err)
			time.Sleep(time.Second)
			continue
		}
		for _, message := range messages {
			e.handle(ctx, message)
		}
	}
}

// reclaimStale hands messages left by a dead executor to this one. The idle
// window must exceed the job timeout, otherwise a live worker running a long
// plan would have its message stolen and the job would run twice.
func (e *Executor) reclaimStale(ctx context.Context) {
	messages, err := e.deps.Queue.ClaimStale(ctx, queue.JobsStream, Group, e.consumer, e.staleClaimIdle(), 10)
	if err != nil {
		if ctx.Err() == nil {
			e.deps.Logger.Warn("could not claim stale jobs", "error", err)
		}
		return
	}
	for _, message := range messages {
		e.handle(ctx, message)
	}
}

// staleClaimIdle is how long a delivered message must sit untouched before it
// is treated as abandoned.
func (e *Executor) staleClaimIdle() time.Duration {
	if e.deps.JobTimeout > 0 {
		return e.deps.JobTimeout + 15*time.Minute
	}
	return 24 * time.Hour
}

// staleWorkspaceAge keeps a job directory for one timeout window plus an hour
// so a concurrently running executor is never cleaned up underneath.
func (e *Executor) staleWorkspaceAge() time.Duration {
	if e.deps.JobTimeout > 0 {
		return e.deps.JobTimeout + time.Hour
	}
	return 24 * time.Hour
}

func (e *Executor) handle(ctx context.Context, message queue.Message) {
	jobID, err := strconv.ParseInt(fmt.Sprint(message.Values["job_id"]), 10, 64)
	if err != nil {
		e.deps.Logger.Error("job message has no usable job_id", "message_id", message.ID, "error", err)
		e.ack(ctx, message.ID)
		return
	}
	attempt := parseAttempt(message.Values["attempt"])

	release, ok, lockErr := e.lockJob(ctx, jobID)
	if lockErr == nil && !ok {
		// Another executor owns this job; it will finish it. Acknowledging the
		// duplicate stops it from being redelivered forever.
		e.deps.Logger.Info("job is already claimed by another executor", "job_id", jobID)
		e.ack(ctx, message.ID)
		return
	}
	if lockErr != nil {
		err = transient(lockErr)
	} else {
		defer release()
		err = e.processJob(ctx, jobID)
	}

	if ctx.Err() != nil {
		// Shutting down: leave the message pending so another executor can
		// reclaim it. The interrupted step stays "running" and is resumed.
		return
	}

	switch {
	case err == nil, errors.Is(err, errJobCancelled):
		e.ack(ctx, message.ID)
	case isTransient(err) && attempt+1 < maxJobAttempts:
		e.retry(ctx, message, jobID, attempt, err)
	default:
		if isTransient(err) {
			e.failExhausted(ctx, jobID, err)
		}
		e.ack(ctx, message.ID)
	}
}

// lockJob takes the job's lease. Without a locker the executor runs
// unguarded, which is what the single-process development setup does.
func (e *Executor) lockJob(ctx context.Context, jobID int64) (func(), bool, error) {
	if e.deps.Locker == nil {
		return func() {}, true, nil
	}
	return e.deps.Locker.TryLock(ctx, jobID)
}

func (e *Executor) ack(ctx context.Context, messageID string) {
	if err := e.deps.Queue.Ack(ctx, queue.JobsStream, Group, messageID); err != nil {
		e.deps.Logger.Warn("ack job message failed", "message_id", messageID, "error", err)
	}
}

// retry republishes a transiently failed job with an incremented attempt. The
// original message is only acknowledged after the retry is queued; a publish
// failure leaves it pending for ClaimStale instead of losing the job.
func (e *Executor) retry(ctx context.Context, message queue.Message, jobID int64, attempt int, cause error) {
	delay := e.retryDelay(attempt)
	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return
		}
	}
	_, err := e.deps.Queue.Publish(ctx, queue.JobsStream, map[string]any{
		"job_id":  strconv.FormatInt(jobID, 10),
		"attempt": attempt + 1,
	})
	if err != nil {
		e.deps.Logger.Error("could not republish failed job; leaving message pending",
			"job_id", jobID, "attempt", attempt+1, "error", err, "cause", cause)
		return
	}
	e.deps.Logger.Warn("retrying job after transient failure",
		"job_id", jobID, "attempt", attempt+1, "error", cause)
	e.ack(ctx, message.ID)
}

// failExhausted marks a job failed once its retry budget is spent.
func (e *Executor) failExhausted(ctx context.Context, jobID int64, cause error) {
	job, err := e.deps.Store.GetJob(ctx, jobID)
	if err != nil {
		e.deps.Logger.Error("could not load job after retries", "job_id", jobID, "error", err)
		return
	}
	if isTerminalStatus(job.Status) {
		return
	}
	if _, err := e.deps.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{
		ID:     jobID,
		Status: "failed",
	}); err != nil {
		e.deps.Logger.Error("could not mark job failed after retries", "job_id", jobID, "error", err)
		return
	}
	e.deps.Logger.Error("job failed after retries", "job_id", jobID, "error", cause)
}

func parseAttempt(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case string:
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return 0
}

// transientError marks failures worth retrying (database, storage or network).
// Terraform failures are permanent and are handled by failJob instead.
type transientError struct{ err error }

func (t transientError) Error() string { return t.err.Error() }
func (t transientError) Unwrap() error { return t.err }

func transient(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, errJobCancelled) {
		return err
	}
	return transientError{err: err}
}

func isTransient(err error) bool {
	var target transientError
	return errors.As(err, &target)
}

// processJob advances a job through its pending steps.
func (e *Executor) processJob(ctx context.Context, jobID int64) error {
	job, err := e.deps.Store.GetJob(ctx, jobID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("job %d does not exist", jobID)
	}
	if err != nil {
		return transient(err)
	}
	if isTerminalStatus(job.Status) {
		return nil
	}
	ws, err := e.deps.Store.GetWorkspace(ctx, job.WorkspaceID)
	if err != nil {
		return transient(fmt.Errorf("load workspace: %w", err))
	}

	for {
		// Steps and approval change the job row, so always work from a fresh
		// copy instead of the snapshot loaded when the message arrived.
		current, err := e.deps.Store.GetJob(ctx, job.ID)
		if err != nil {
			return transient(err)
		}
		job = current
		if isCancelledStatus(job.Status) {
			return errJobCancelled
		}

		step, err := e.deps.Store.GetPendingStepByJob(ctx, job.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			return e.finishJob(ctx, job, ws, "completed")
		}
		if err != nil {
			return transient(err)
		}

		switch step.StepNumber {
		case planStepNumber:
			if err := e.runPlan(ctx, job, ws, step); err != nil {
				if errors.Is(err, errJobCancelled) || ctx.Err() != nil || isTransient(err) {
					return err
				}
				return e.failJob(ctx, job, ws, step, err)
			}
			job, _ = e.deps.Store.GetJob(ctx, job.ID)
			if e.isSpeculative(ctx, job) {
				// Upstream's speculative "Plan" template has no approval or
				// apply step; the run finishes with the plan.
				return e.finishJob(ctx, job, ws, "completed")
			}
			if job.AutoApply {
				if err := e.approveAutomatically(ctx, job); err != nil {
					return err
				}
				continue
			}
			return e.waitForApproval(ctx, job, ws)
		case approvalStepNumber:
			if job.ApprovedAt.Valid {
				_, _ = e.deps.Store.UpdateStepStatus(ctx, store.UpdateStepStatusParams{
					ID:     step.ID,
					Status: stringPointer("completed"),
				})
				job, _ = e.deps.Store.GetJob(ctx, job.ID)
				continue
			}
			return e.waitForApproval(ctx, job, ws)
		case applyStepNumber:
			if !job.ApprovedAt.Valid {
				return e.waitForApproval(ctx, job, ws)
			}
			if err := e.runApply(ctx, job, ws, step); err != nil {
				if errors.Is(err, errJobCancelled) || ctx.Err() != nil || isTransient(err) {
					return err
				}
				return e.failJob(ctx, job, ws, step, err)
			}
			return e.finishJob(ctx, job, ws, "completed")
		default:
			e.deps.Logger.Warn("skipping unknown step", "job_id", job.ID, "step_number", step.StepNumber)
			_, _ = e.deps.Store.UpdateStepStatus(ctx, store.UpdateStepStatusParams{
				ID:     step.ID,
				Status: stringPointer("completed"),
			})
		}
	}
}

func (e *Executor) approveAutomatically(ctx context.Context, job store.Job) error {
	approval, err := e.deps.Store.GetStepByJobAndNumber(ctx, store.GetStepByJobAndNumberParams{
		JobID:      job.ID,
		StepNumber: approvalStepNumber,
	})
	if err == nil {
		if _, err := e.deps.Store.UpdateStepStatus(ctx, store.UpdateStepStatusParams{
			ID:     approval.ID,
			Status: stringPointer("completed"),
		}); err != nil {
			return err
		}
	}
	_, err = e.deps.Store.ApproveJob(ctx, store.ApproveJobParams{
		ID:         job.ID,
		ApprovedBy: stringPointer("auto-apply"),
	})
	return err
}

func (e *Executor) isSpeculative(ctx context.Context, job store.Job) bool {
	if job.Command != nil && strings.EqualFold(*job.Command, "destroy") {
		return false
	}
	contentID := contentIDFromSource(job.OverrideSource)
	if contentID == "" {
		return false
	}
	content, err := e.deps.Store.GetContent(ctx, contentID)
	if err != nil {
		return false
	}
	return content.Speculative != nil && *content.Speculative
}

func (e *Executor) waitForApproval(ctx context.Context, job store.Job, ws store.Workspace) error {
	if _, err := e.deps.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{
		ID:     job.ID,
		Status: "waitingApproval",
	}); err != nil {
		return err
	}
	e.setWorkspaceLastJob(ctx, ws, "waitingApproval")
	return nil
}

func (e *Executor) finishJob(ctx context.Context, job store.Job, ws store.Workspace, status string) error {
	if _, err := e.deps.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{
		ID:     job.ID,
		Status: status,
	}); err != nil {
		return err
	}
	e.setWorkspaceLastJob(ctx, ws, status)
	return nil
}

func (e *Executor) failJob(ctx context.Context, job store.Job, ws store.Workspace, step store.Step, cause error) error {
	e.deps.Logger.Error("step failed", "job_id", job.ID, "step", step.StepNumber, "error", cause)
	_, _ = e.deps.Store.UpdateStepStatus(ctx, store.UpdateStepStatusParams{
		ID:     step.ID,
		Status: stringPointer("failed"),
	})
	if _, err := e.deps.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{
		ID:     job.ID,
		Status: "failed",
	}); err != nil {
		return err
	}
	e.setWorkspaceLastJob(ctx, ws, "failed")
	return cause
}

func (e *Executor) setWorkspaceLastJob(ctx context.Context, ws store.Workspace, status string) {
	if err := e.deps.Store.SetWorkspaceLastJob(ctx, store.SetWorkspaceLastJobParams{
		ID:            ws.ID,
		LastJobStatus: status,
	}); err != nil {
		e.deps.Logger.Warn("could not update workspace last job", "workspace_id", ws.ID, "error", err)
	}
}

func (e *Executor) markStep(ctx context.Context, step store.Step, status string, output string) error {
	if _, err := e.deps.Store.UpdateStepStatus(ctx, store.UpdateStepStatusParams{
		ID:     step.ID,
		Status: stringPointer(status),
	}); err != nil {
		return err
	}
	if output != "" {
		if err := e.deps.Store.UpdateStepOutput(ctx, store.UpdateStepOutputParams{
			ID:     step.ID,
			Output: stringPointer(output),
		}); err != nil {
			return err
		}
	}
	return nil
}

func isTerminalStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "completed", "failed", "cancelled", "canceled", "nochanges", "rejected", "notexecuted", "neverexecuted":
		return true
	default:
		return false
	}
}

func isCancelledStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "cancelled", "canceled":
		return true
	default:
		return false
	}
}

// watchJob returns a context that is cancelled when the job is cancelled in the
// API. Terraform commands run under it, so a discard stops the child process
// instead of letting it finish and overwrite the cancelled status.
func (e *Executor) watchJob(ctx context.Context, jobID int64) (context.Context, context.CancelFunc) {
	jobCtx, cancel := context.WithCancel(ctx)
	interval := e.cancelPollInterval
	if interval <= 0 {
		interval = 2 * time.Second
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-jobCtx.Done():
				return
			case <-ticker.C:
				job, err := e.deps.Store.GetJob(jobCtx, jobID)
				if err != nil {
					continue
				}
				if isCancelledStatus(job.Status) {
					cancel()
					return
				}
				if isTerminalStatus(job.Status) {
					return
				}
			}
		}
	}()
	return jobCtx, cancel
}

func stringPointer(value string) *string {
	return &value
}
