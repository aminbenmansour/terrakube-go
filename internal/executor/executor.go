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
	claimIdle                 = 30 * time.Second
)

// Deps are the collaborators the executor needs.
type Deps struct {
	Store               store.Querier
	Storage             storage.Backend
	Queue               *queue.Queue
	Layout              tfstate.Layout
	Logger              *slog.Logger
	PublicURL           string
	InternalTokenSecret string
	Binary              string
	WorkDir             string
}

// Executor processes jobs from the Redis stream.
type Executor struct {
	deps     Deps
	consumer string
}

// New builds an executor.
func New(deps Deps) *Executor {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	host, _ := os.Hostname()
	return &Executor{
		deps:     deps,
		consumer: fmt.Sprintf("%s-%d-%s", host, os.Getpid(), uuid.NewString()[:8]),
	}
}

// Run consumes jobs until the context is cancelled.
func (e *Executor) Run(ctx context.Context) error {
	if err := e.deps.Queue.EnsureGroup(ctx, queue.JobsStream, Group); err != nil {
		return err
	}
	e.deps.Logger.Info("executor consuming jobs", "stream", queue.JobsStream, "group", Group, "consumer", e.consumer)

	// Pick up jobs that were delivered before a previous executor stopped.
	if stale, err := e.deps.Queue.ClaimStale(ctx, queue.JobsStream, Group, e.consumer, claimIdle, 10); err != nil {
		e.deps.Logger.Warn("could not claim stale jobs", "error", err)
	} else {
		for _, message := range stale {
			e.handle(ctx, message)
		}
	}

	for {
		messages, err := e.deps.Queue.ReadGroup(ctx, queue.JobsStream, Group, e.consumer, 1, 5*time.Second)
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

func (e *Executor) handle(ctx context.Context, message queue.Message) {
	jobID, err := strconv.ParseInt(fmt.Sprint(message.Values["job_id"]), 10, 64)
	if err != nil {
		e.deps.Logger.Error("job message has no usable job_id", "message_id", message.ID, "error", err)
	} else if err := e.processJob(ctx, jobID); err != nil {
		e.deps.Logger.Error("job failed", "job_id", jobID, "error", err)
	}
	if err := e.deps.Queue.Ack(ctx, queue.JobsStream, Group, message.ID); err != nil {
		e.deps.Logger.Warn("ack job message failed", "message_id", message.ID, "error", err)
	}
}

// processJob advances a job through its pending steps.
func (e *Executor) processJob(ctx context.Context, jobID int64) error {
	job, err := e.deps.Store.GetJob(ctx, jobID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("job %d does not exist", jobID)
	}
	if err != nil {
		return err
	}
	if isTerminalStatus(job.Status) {
		return nil
	}
	ws, err := e.deps.Store.GetWorkspace(ctx, job.WorkspaceID)
	if err != nil {
		return fmt.Errorf("load workspace: %w", err)
	}

	for {
		// Steps and approval change the job row, so always work from a fresh
		// copy instead of the snapshot loaded when the message arrived.
		current, err := e.deps.Store.GetJob(ctx, job.ID)
		if err != nil {
			return err
		}
		job = current

		step, err := e.deps.Store.GetPendingStepByJob(ctx, job.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			return e.finishJob(ctx, job, ws, "completed")
		}
		if err != nil {
			return err
		}

		switch step.StepNumber {
		case planStepNumber:
			if err := e.runPlan(ctx, job, ws, step); err != nil {
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

func stringPointer(value string) *string {
	return &value
}
