package executor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	planFileName  = "tfplan"
	maxLogBytes   = 10 << 20
	logLinePrefix = ""
)

func (e *Executor) runPlan(ctx context.Context, job store.Job, ws store.Workspace, step store.Step) error {
	binary := e.binary(ws)
	workDir, runDir, env, err := e.prepareWorkspace(ctx, job, ws, step)
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(workDir) }()

	sink := newLogSink(ctx, e.deps.Queue, e.deps.Logger, job.ID, step.ID)
	e.startStep(ctx, job, step)
	sink.WriteString("Executing terraform init\n")

	if err := e.runCommand(ctx, runDir, env, sink, binary, "init", "-input=false", "-no-color"); err != nil {
		e.completeStep(ctx, job, step, sink, "failed")
		return fmt.Errorf("terraform init: %w", err)
	}

	sink.WriteString("Executing terraform plan\n")
	args := []string{"plan", "-input=false", "-no-color", "-detailed-exitcode", "-out=" + planFileName}
	if job.Command != nil && strings.EqualFold(*job.Command, "destroy") {
		args = []string{"plan", "-destroy", "-input=false", "-no-color", "-detailed-exitcode", "-out=" + planFileName}
	}
	err = e.runCommand(ctx, runDir, env, sink, binary, args...)
	changes := false
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 2 {
			// -detailed-exitcode uses 2 to signal a non-empty diff.
			changes = true
		} else {
			e.completeStep(ctx, job, step, sink, "failed")
			return fmt.Errorf("terraform plan: %w", err)
		}
	}

	if err := e.savePlan(ctx, job, step, runDir); err != nil {
		e.completeStep(ctx, job, step, sink, "failed")
		return fmt.Errorf("save plan artifact: %w", err)
	}
	if _, err := e.deps.Store.UpdateJobPlanChanges(ctx, store.UpdateJobPlanChangesParams{
		ID:          job.ID,
		PlanChanges: changes,
	}); err != nil {
		e.completeStep(ctx, job, step, sink, "failed")
		return err
	}

	e.completeStep(ctx, job, step, sink, "completed")
	return nil
}

func (e *Executor) runApply(ctx context.Context, job store.Job, ws store.Workspace, step store.Step) error {
	binary := e.binary(ws)
	workDir, runDir, env, err := e.prepareWorkspace(ctx, job, ws, step)
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(workDir) }()

	sink := newLogSink(ctx, e.deps.Queue, e.deps.Logger, job.ID, step.ID)
	e.startStep(ctx, job, step)
	sink.WriteString("Executing terraform init\n")

	if err := e.runCommand(ctx, runDir, env, sink, binary, "init", "-input=false", "-no-color"); err != nil {
		e.completeStep(ctx, job, step, sink, "failed")
		return fmt.Errorf("terraform init: %w", err)
	}

	sink.WriteString("Executing terraform apply\n")
	hasPlan := false
	if planStep, err := e.deps.Store.GetStepByJobAndNumber(ctx, store.GetStepByJobAndNumberParams{
		JobID:      job.ID,
		StepNumber: planStepNumber,
	}); err == nil {
		hasPlan = e.downloadPlan(ctx, job, planStep, runDir) == nil
	}
	args := []string{"apply", "-auto-approve", "-input=false", "-no-color"}
	if hasPlan {
		args = append(args, planFileName)
	} else {
		sink.WriteString("Saved plan not available; applying a fresh plan\n")
	}
	if err := e.runCommand(ctx, runDir, env, sink, binary, args...); err != nil {
		e.completeStep(ctx, job, step, sink, "failed")
		return fmt.Errorf("terraform apply: %w", err)
	}

	if err := e.archiveHistory(ctx, job, ws); err != nil {
		// The apply itself succeeded; a history failure should not mark the run
		// failed, but it must be visible.
		e.deps.Logger.Warn("could not archive state history", "job_id", job.ID, "error", err)
	}
	e.completeStep(ctx, job, step, sink, "completed")
	return nil
}

func (e *Executor) startStep(ctx context.Context, job store.Job, step store.Step) {
	if _, err := e.deps.Store.UpdateJobStatus(ctx, store.UpdateJobStatusParams{
		ID:     job.ID,
		Status: "running",
	}); err != nil {
		e.deps.Logger.Warn("could not mark job running", "job_id", job.ID, "error", err)
	}
	if _, err := e.deps.Store.UpdateStepStatus(ctx, store.UpdateStepStatusParams{
		ID:     step.ID,
		Status: stringPointer("running"),
	}); err != nil {
		e.deps.Logger.Warn("could not mark step running", "step_id", step.ID, "error", err)
	}
}

func (e *Executor) completeStep(ctx context.Context, job store.Job, step store.Step, sink *logSink, status string) {
	sink.Flush()
	output := sink.String()
	if err := e.deps.Store.UpdateStepOutput(ctx, store.UpdateStepOutputParams{
		ID:     step.ID,
		Output: stringPointer(output),
	}); err != nil {
		e.deps.Logger.Warn("could not save step output", "step_id", step.ID, "error", err)
	}
	if err := e.deps.Store.UpdateJobOutput(ctx, store.UpdateJobOutputParams{
		ID:     job.ID,
		Output: stringPointer(output),
	}); err != nil {
		e.deps.Logger.Warn("could not save job output", "job_id", job.ID, "error", err)
	}
	if _, err := e.deps.Store.UpdateStepStatus(ctx, store.UpdateStepStatusParams{
		ID:     step.ID,
		Status: stringPointer(status),
	}); err != nil {
		e.deps.Logger.Warn("could not save step status", "step_id", step.ID, "error", err)
	}
}

func (e *Executor) runCommand(ctx context.Context, dir string, env []string, sink *logSink, binary string, args ...string) error {
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = sink
	cmd.Stderr = sink
	err := cmd.Run()
	sink.Flush()
	return err
}

func (e *Executor) savePlan(ctx context.Context, job store.Job, step store.Step, runDir string) error {
	if e.deps.Storage == nil {
		return nil
	}
	planPath := filepath.Join(runDir, planFileName)
	file, err := os.Open(planPath)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	key := e.deps.Layout.Plan(job.OrganizationID, job.WorkspaceID, fmt.Sprintf("%d", job.ID), step.ID)
	return e.deps.Storage.Put(ctx, key, file, info.Size(), "application/octet-stream")
}

func (e *Executor) downloadPlan(ctx context.Context, job store.Job, planStep store.Step, runDir string) error {
	if e.deps.Storage == nil {
		return storage.ErrNotFound
	}
	key := e.deps.Layout.Plan(job.OrganizationID, job.WorkspaceID, fmt.Sprintf("%d", job.ID), planStep.ID)
	rc, err := e.deps.Storage.Get(ctx, key)
	if err != nil {
		return err
	}
	defer rc.Close()

	file, err := os.OpenFile(filepath.Join(runDir, planFileName), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o640)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.ReadFrom(rc)
	return err
}

func (e *Executor) binary(ws store.Workspace) string {
	if e.deps.Binary != "" {
		return e.deps.Binary
	}
	if strings.Contains(strings.ToLower(ws.IacType), "tofu") {
		return "tofu"
	}
	if path, err := exec.LookPath("terraform"); err == nil {
		return path
	}
	return "tofu"
}

// logSink fans command output out to Redis and keeps the final log for the
// step row.
type logSink struct {
	ctx    context.Context
	queue  *queue.Queue
	logger *slog.Logger
	jobID  int64
	stepID string

	mu      sync.Mutex
	pending string
	output  strings.Builder
	trunc   bool
}

func newLogSink(ctx context.Context, q *queue.Queue, logger *slog.Logger, jobID int64, stepID string) *logSink {
	return &logSink{ctx: ctx, queue: q, logger: logger, jobID: jobID, stepID: stepID}
}

func (s *logSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending += string(p)
	for {
		index := strings.IndexByte(s.pending, '\n')
		if index < 0 {
			break
		}
		line := s.pending[:index+1]
		s.pending = s.pending[index+1:]
		s.emit(line)
	}
	return len(p), nil
}

func (s *logSink) WriteString(line string) {
	_, _ = s.Write([]byte(line))
}

// Flush publishes a trailing line that was not newline terminated.
func (s *logSink) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pending != "" {
		s.emit(s.pending)
		s.pending = ""
	}
}

func (s *logSink) emit(line string) {
	if !s.trunc {
		if s.output.Len()+len(line) > maxLogBytes {
			s.trunc = true
			s.output.WriteString("\n[terrakube: log truncated]\n")
		} else {
			s.output.WriteString(line)
		}
	}
	if s.queue == nil {
		return
	}
	if _, err := s.queue.PublishCapped(s.ctx, fmt.Sprintf("%d", s.jobID), map[string]any{
		"stepId": s.stepID,
		"output": logLinePrefix + line,
	}, maxLogStreamEntries); err != nil && s.logger != nil {
		s.logger.Warn("could not publish live log line", "job_id", s.jobID, "error", err)
	}
}

func (s *logSink) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.output.String()
}
