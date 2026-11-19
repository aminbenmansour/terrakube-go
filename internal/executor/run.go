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

	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	planFileName  = "tfplan"
	maxLogBytes   = 10 << 20
	logLinePrefix = ""
)

func (e *Executor) runPlan(ctx context.Context, job store.Job, ws store.Workspace, step store.Step) error {
	jobCtx, stopWatching := e.watchJob(ctx, job.ID)
	defer stopWatching()
	jobCtx, stopTimeout := e.jobContext(jobCtx)
	defer stopTimeout()

	binary := e.binary(ws)
	prepared, err := e.prepareWorkspace(jobCtx, job, ws, step)
	if err != nil {
		if e.timedOut(jobCtx) {
			return e.timeoutError(err)
		}
		if jobCtx.Err() != nil && ctx.Err() == nil {
			return errJobCancelled
		}
		return err
	}
	defer func() { _ = os.RemoveAll(prepared.workDir) }()

	sink := newLogSink(ctx, e.deps.Queue, e.deps.Logger, job.ID, step.ID, newRedactor(prepared.secrets))
	if !e.startStep(ctx, job, step) {
		return errJobCancelled
	}
	sink.WriteString("Executing terraform init\n")

	if err := e.runCommand(jobCtx, prepared.runDir, prepared.env, sink, binary, "init", "-input=false", "-no-color"); err != nil {
		return e.handleCommandError(ctx, jobCtx, job, step, sink, "terraform init", err)
	}

	sink.WriteString("Executing terraform plan\n")
	args := []string{"plan", "-input=false", "-no-color", "-detailed-exitcode", "-out=" + planFileName}
	if job.Command != nil && strings.EqualFold(*job.Command, "destroy") {
		args = []string{"plan", "-destroy", "-input=false", "-no-color", "-detailed-exitcode", "-out=" + planFileName}
	}
	err = e.runCommand(jobCtx, prepared.runDir, prepared.env, sink, binary, args...)
	changes := false
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 2 && !e.timedOut(jobCtx) && jobCtx.Err() == nil {
			// -detailed-exitcode uses 2 to signal a non-empty diff.
			changes = true
		} else {
			return e.handleCommandError(ctx, jobCtx, job, step, sink, "terraform plan", err)
		}
	}

	if err := e.savePlan(ctx, job, step, prepared.runDir); err != nil {
		_ = e.completeStep(ctx, job, step, sink, "failed")
		return fmt.Errorf("save plan artifact: %w", err)
	}
	if _, err := e.deps.Store.UpdateJobPlanChanges(ctx, store.UpdateJobPlanChangesParams{
		ID:          job.ID,
		PlanChanges: changes,
	}); err != nil {
		_ = e.completeStep(ctx, job, step, sink, "failed")
		return err
	}

	return e.completeStep(ctx, job, step, sink, "completed")
}

func (e *Executor) runApply(ctx context.Context, job store.Job, ws store.Workspace, step store.Step) error {
	jobCtx, stopWatching := e.watchJob(ctx, job.ID)
	defer stopWatching()
	jobCtx, stopTimeout := e.jobContext(jobCtx)
	defer stopTimeout()

	binary := e.binary(ws)
	prepared, err := e.prepareWorkspace(jobCtx, job, ws, step)
	if err != nil {
		if e.timedOut(jobCtx) {
			return e.timeoutError(err)
		}
		if jobCtx.Err() != nil && ctx.Err() == nil {
			return errJobCancelled
		}
		return err
	}
	defer func() { _ = os.RemoveAll(prepared.workDir) }()

	sink := newLogSink(ctx, e.deps.Queue, e.deps.Logger, job.ID, step.ID, newRedactor(prepared.secrets))
	if !e.startStep(ctx, job, step) {
		return errJobCancelled
	}
	sink.WriteString("Executing terraform init\n")

	if err := e.runCommand(jobCtx, prepared.runDir, prepared.env, sink, binary, "init", "-input=false", "-no-color"); err != nil {
		return e.handleCommandError(ctx, jobCtx, job, step, sink, "terraform init", err)
	}

	sink.WriteString("Executing terraform apply\n")
	hasPlan := false
	if planStep, err := e.deps.Store.GetStepByJobAndNumber(ctx, store.GetStepByJobAndNumberParams{
		JobID:      job.ID,
		StepNumber: planStepNumber,
	}); err == nil {
		hasPlan = e.downloadPlan(jobCtx, job, planStep, prepared.runDir) == nil
	}
	args := []string{"apply", "-auto-approve", "-input=false", "-no-color"}
	if hasPlan {
		args = append(args, planFileName)
	} else {
		sink.WriteString("Saved plan not available; applying a fresh plan\n")
	}
	if err := e.runCommand(jobCtx, prepared.runDir, prepared.env, sink, binary, args...); err != nil {
		return e.handleCommandError(ctx, jobCtx, job, step, sink, "terraform apply", err)
	}

	if err := e.archiveHistory(ctx, job, ws); err != nil {
		// The apply itself succeeded; a history failure should not mark the run
		// failed, but it must be visible.
		e.deps.Logger.Warn("could not archive state history", "job_id", job.ID, "error", err)
	}
	return e.completeStep(ctx, job, step, sink, "completed")
}

// startStep marks the job and step running. It returns false when the job was
// cancelled while the workspace was being prepared, so a discarded run is not
// resurrected by the executor.
func (e *Executor) startStep(ctx context.Context, job store.Job, step store.Step) bool {
	if current, err := e.deps.Store.GetJob(ctx, job.ID); err == nil && isCancelledStatus(current.Status) {
		return false
	}
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
	return true
}

func (e *Executor) completeStep(ctx context.Context, job store.Job, step store.Step, sink *logSink, status string) error {
	sink.Flush()
	output := sink.String()

	effective := status
	archiveErr := e.archiveStepLog(ctx, job, step, output)
	if archiveErr != nil {
		e.deps.Logger.Warn("could not archive step log", "step_id", step.ID, "error", archiveErr)
		if status == "completed" {
			effective = "failed"
		}
	}

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
		Status: stringPointer(effective),
	}); err != nil {
		e.deps.Logger.Warn("could not save step status", "step_id", step.ID, "error", err)
	}

	if archiveErr != nil && status == "completed" {
		return fmt.Errorf("archive step log: %w", archiveErr)
	}
	return nil
}

// archiveStepLog writes the final console output under the upstream object
// key. The inline step.output column remains for the TFE log endpoints.
func (e *Executor) archiveStepLog(ctx context.Context, job store.Job, step store.Step, output string) error {
	if e.deps.Storage == nil {
		return nil
	}
	key := e.deps.Layout.Output(job.OrganizationID, job.ID, step.ID)
	return e.deps.Storage.Put(ctx, key, strings.NewReader(output), int64(len(output)), "text/plain")
}

// inheritedEnvKeys are the host variables a Terraform command receives. The
// executor's own environment holds database, Redis, storage and signing secrets
// that untrusted configuration must not be able to read.
var inheritedEnvKeys = []string{
	"PATH", "HOME", "TMPDIR", "TMP", "TEMP", "TZ", "LANG", "LC_ALL",
	"SSL_CERT_FILE", "SSL_CERT_DIR", "HTTP_PROXY", "HTTPS_PROXY", "NO_PROXY",
	"http_proxy", "https_proxy", "no_proxy", "SYSTEMROOT",
}

func baseEnvironment() []string {
	env := make([]string, 0, len(inheritedEnvKeys))
	for _, key := range inheritedEnvKeys {
		if value, ok := os.LookupEnv(key); ok {
			env = append(env, key+"="+value)
		}
	}
	return env
}

// handleCommandError records a failed command and classifies it. A timeout is
// a failed step, a job cancellation is a cancelled step, and anything else is
// an actionable Terraform error.
func (e *Executor) handleCommandError(ctx, jobCtx context.Context, job store.Job, step store.Step, sink *logSink, label string, err error) error {
	if ctx.Err() != nil {
		// Shutting down: leave the step running so the reclaimed job resumes it.
		return ctx.Err()
	}
	if e.timedOut(jobCtx) {
		sink.WriteString(fmt.Sprintf("Job timed out after %s\n", e.deps.JobTimeout))
		_ = e.completeStep(ctx, job, step, sink, "failed")
		return e.timeoutError(err)
	}
	if jobCtx.Err() != nil && ctx.Err() == nil {
		_ = e.completeStep(ctx, job, step, sink, "cancelled")
		return errJobCancelled
	}
	_ = e.completeStep(ctx, job, step, sink, "failed")
	return fmt.Errorf("%s: %w", label, err)
}

// jobContext bounds a job by the configured timeout. Without a timeout it only
// follows the parent context.
func (e *Executor) jobContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if e.deps.JobTimeout > 0 {
		return context.WithTimeout(ctx, e.deps.JobTimeout)
	}
	return context.WithCancel(ctx)
}

func (e *Executor) timedOut(ctx context.Context) bool {
	return errors.Is(ctx.Err(), context.DeadlineExceeded)
}

func (e *Executor) timeoutError(cause error) error {
	return fmt.Errorf("job timed out after %s: %w", e.deps.JobTimeout, cause)
}

func (e *Executor) runCommand(ctx context.Context, dir string, env []string, sink *logSink, binary string, args ...string) error {
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	cmd.Env = append(baseEnvironment(), env...)
	cmd.Stdout = sink
	cmd.Stderr = sink
	configureProcessGroup(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		sink.Flush()
		return err
	case <-ctx.Done():
		killProcessGroup(cmd.Process)
		<-done
		sink.Flush()
		return ctx.Err()
	}
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
	ctx      context.Context
	queue    LogPublisher
	logger   *slog.Logger
	jobID    int64
	stepID   string
	redactor *redactor

	mu      sync.Mutex
	pending string
	output  strings.Builder
	trunc   bool
}

func newLogSink(ctx context.Context, q LogPublisher, logger *slog.Logger, jobID int64, stepID string, redact *redactor) *logSink {
	return &logSink{ctx: ctx, queue: q, logger: logger, jobID: jobID, stepID: stepID, redactor: redact}
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
	line = s.redactor.redact(line)
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
