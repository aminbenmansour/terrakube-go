package executor

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/config"
	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
	"github.com/aminbenmansour/terrakube-go/migrations"
)

// TestExecutorEndToEnd runs a job against a real Postgres and Redis. It is
// skipped unless TERRAKUBE_TEST_DATABASE_URL and TERRAKUBE_TEST_REDIS_URL are
// set; CI provides both as service containers. The Terraform binary is a fake
// script, so the test needs no cloud credentials and is deterministic.
func TestExecutorEndToEnd(t *testing.T) {
	databaseURL := os.Getenv("TERRAKUBE_TEST_DATABASE_URL")
	redisURL := os.Getenv("TERRAKUBE_TEST_REDIS_URL")
	if databaseURL == "" || redisURL == "" {
		t.Skip("set TERRAKUBE_TEST_DATABASE_URL and TERRAKUBE_TEST_REDIS_URL to run the executor end-to-end test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	defer pool.Close()
	if err := db.Migrate(ctx, pool, migrations.Files, "."); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	st := store.New(pool)

	q, err := queue.New(redisURL)
	if err != nil {
		t.Fatalf("connect redis: %v", err)
	}
	defer q.Close()

	backend, err := storage.New(ctx, config.Storage{Type: config.StorageTypeLocal, LocalPath: t.TempDir()})
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	layout := tfstate.NewLayout("local")

	orgID := uuid.NewString()
	if _, err := st.CreateOrganization(ctx, store.CreateOrganizationParams{
		ID:            orgID,
		Name:          "executor-e2e-" + orgID[:8],
		ExecutionMode: "remote",
	}); err != nil {
		t.Fatalf("create organization: %v", err)
	}
	projectID := uuid.NewString()
	if _, err := st.CreateProject(ctx, store.CreateProjectParams{
		ID:             projectID,
		Name:           "executor-e2e",
		OrganizationID: orgID,
	}); err != nil {
		t.Fatalf("create project: %v", err)
	}
	wsID := uuid.NewString()
	if _, err := st.CreateWorkspace(ctx, store.CreateWorkspaceParams{
		ID:             wsID,
		Name:           "executor-e2e",
		OrganizationID: orgID,
		ProjectID:      &projectID,
		IacType:        "opentofu",
	}); err != nil {
		t.Fatalf("create workspace: %v", err)
	}

	// Seed a state so the apply has something to archive into history. The
	// fake Terraform binary does not speak the HTTP backend protocol.
	state := []byte(`{"version":4,"serial":3,"lineage":"lineage-e2e"}`)
	if err := backend.Put(ctx, layout.CurrentState(orgID, wsID), bytes.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("seed state: %v", err)
	}

	callsFile := filepath.Join(t.TempDir(), "calls")

	for _, variable := range []struct {
		key       string
		value     string
		category  string
		sensitive bool
	}{
		{key: "E2E_ENV", value: "env-value", category: "ENV"},
		{key: "E2E_VAR", value: "var-value", category: "TERRAFORM"},
		{key: "E2E_SECRET", value: "super-secret-value", category: "ENV", sensitive: true},
		{key: "E2E_CALLS_FILE", value: callsFile, category: "ENV"},
	} {
		if _, err := st.CreateVariable(ctx, store.CreateVariableParams{
			ID:               uuid.NewString(),
			VariableKey:      variable.key,
			VariableValue:    variable.value,
			VariableCategory: variable.category,
			Sensitive:        variable.sensitive,
			WorkspaceID:      wsID,
		}); err != nil {
			t.Fatalf("create variable %s: %v", variable.key, err)
		}
	}

	contentID := uuid.NewString()
	archive := tarGz(t, map[string]string{
		"main.tf": `resource "terraform_data" "example" {}`,
	})
	if err := backend.Put(ctx, layout.Content(contentID), bytes.NewReader(archive), int64(len(archive)), "application/gzip"); err != nil {
		t.Fatalf("store configuration: %v", err)
	}
	source := "https://example.invalid/remote/tfe/v2/configuration-versions/" + contentID + "/terraformContent.tar.gz"

	job, err := st.CreateJob(ctx, store.CreateJobParams{
		OrganizationID: orgID,
		WorkspaceID:    wsID,
		Status:         "pending",
		Via:            "cli",
		AutoApply:      true,
		Refresh:        true,
	})
	if err != nil {
		t.Fatalf("create job: %v", err)
	}
	if _, err := st.SetJobOverrideSource(ctx, store.SetJobOverrideSourceParams{
		ID:             job.ID,
		OverrideSource: &source,
		OverrideBranch: stringPointer("remote-content"),
	}); err != nil {
		t.Fatalf("set configuration source: %v", err)
	}

	// The plan step starts as "running" to prove a reclaimed job resumes the
	// interrupted step instead of treating the job as finished.
	for _, spec := range []struct {
		number int32
		name   string
		status string
	}{
		{planStepNumber, "Plan", "running"},
		{approvalStepNumber, "Approve", "pending"},
		{applyStepNumber, "Apply", "pending"},
	} {
		if _, err := st.CreateStep(ctx, store.CreateStepParams{
			ID:         uuid.NewString(),
			JobID:      job.ID,
			StepNumber: spec.number,
			Name:       stringPointer(spec.name),
			Status:     stringPointer(spec.status),
			LogStatus:  stringPointer("pending"),
		}); err != nil {
			t.Fatalf("create step %d: %v", spec.number, err)
		}
	}

	binary := filepath.Join(t.TempDir(), "tofu")
	script := "#!/bin/sh\n" +
		"echo \"$1\" >> \"$E2E_CALLS_FILE\"\n" +
		"echo \"fake tofu $*\"\n" +
		"echo \"env E2E_ENV=$E2E_ENV\"\n" +
		"echo \"var E2E_VAR=$TF_VAR_E2E_VAR\"\n" +
		"echo \"secret E2E_SECRET=$E2E_SECRET\"\n" +
		"case \"$1\" in\n" +
		"  plan) : > tfplan ;;\n" +
		"esac\n" +
		"exit 0\n"
	if err := os.WriteFile(binary, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake terraform: %v", err)
	}

	e := New(Deps{
		Store:     st,
		Storage:   backend,
		Queue:     q,
		Layout:    layout,
		Logger:    testLogger(),
		PublicURL: "https://example.invalid",
		Binary:    binary,
		WorkDir:   t.TempDir(),
	})
	if _, err := q.EnqueueJob(ctx, job.ID); err != nil {
		t.Fatalf("enqueue job: %v", err)
	}

	runCtx, stop := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- e.Run(runCtx) }()
	defer func() {
		stop()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("executor did not stop")
		}
	}()

	deadline := time.Now().Add(60 * time.Second)
	var final store.Job
	for {
		final, err = st.GetJob(ctx, job.ID)
		if err != nil {
			t.Fatalf("load job: %v", err)
		}
		if isTerminalStatus(final.Status) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("job %d did not finish; last status %q", job.ID, final.Status)
		}
		time.Sleep(50 * time.Millisecond)
	}
	if final.Status != "completed" {
		t.Fatalf("job status = %q, want completed (output: %v)", final.Status, deref(final.Output))
	}

	steps, err := st.ListStepsByJob(ctx, job.ID)
	if err != nil {
		t.Fatalf("list steps: %v", err)
	}
	if len(steps) != 3 {
		t.Fatalf("steps = %d, want 3", len(steps))
	}
	for _, step := range steps {
		if deref(step.Status) != "completed" {
			t.Errorf("step %d status = %q, want completed", step.StepNumber, deref(step.Status))
		}
		if step.StepNumber == planStepNumber && step.Output == nil {
			t.Error("plan step has no output")
		}
		// The approval step has no command output, so it has no archived log.
		if step.StepNumber == approvalStepNumber {
			continue
		}
		rc, err := backend.Get(ctx, layout.Output(orgID, job.ID, step.ID))
		if err != nil {
			t.Fatalf("read archived log for step %d: %v", step.StepNumber, err)
		}
		data, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			t.Fatalf("read archived log body: %v", err)
		}
		if !bytes.Contains(data, []byte("fake tofu")) {
			t.Errorf("archived log for step %d = %q", step.StepNumber, data)
		}
	}

	if final.Output == nil || !bytes.Contains([]byte(*final.Output), []byte("fake tofu")) {
		t.Fatalf("job output does not contain the fake terraform output: %v", deref(final.Output))
	}
	output := deref(final.Output)
	for _, want := range []string{"env E2E_ENV=env-value", "var E2E_VAR=var-value"} {
		if !strings.Contains(output, want) {
			t.Errorf("job output is missing %q:\n%s", want, output)
		}
	}
	if strings.Contains(output, "super-secret-value") {
		t.Errorf("job output leaks a sensitive variable:\n%s", output)
	}
	if !strings.Contains(output, "[REDACTED]") {
		t.Errorf("job output does not show a redacted secret:\n%s", output)
	}

	jobReference := strconv.FormatInt(job.ID, 10)
	history, err := st.GetHistoryByJobReference(ctx, store.GetHistoryByJobReferenceParams{
		WorkspaceID:  wsID,
		JobReference: &jobReference,
	})
	if err != nil {
		t.Fatalf("load history: %v", err)
	}
	wantURL := "https://example.invalid/tfstate/v1/organization/" + orgID + "/workspace/" + wsID + "/state/" + history.ID + ".json"
	if history.Output == nil || *history.Output != wantURL {
		t.Errorf("history output = %v, want %q", history.Output, wantURL)
	}
	if history.Serial != 3 || history.Lineage != "lineage-e2e" {
		t.Errorf("history serial/lineage = %d/%q", history.Serial, history.Lineage)
	}
	sum := md5.Sum(state)
	if history.Md5 != hex.EncodeToString(sum[:]) {
		t.Errorf("history md5 = %q", history.Md5)
	}

	calls, err := os.ReadFile(callsFile)
	if err != nil {
		t.Fatalf("read terraform calls: %v", err)
	}
	if got, want := strings.Fields(string(calls)), []string{"init", "plan", "init", "apply"}; strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("terraform calls = %v, want %v", got, want)
	}

	logs, err := q.ReadStream(ctx, strconv.FormatInt(job.ID, 10), "0", 100, time.Second)
	if err != nil {
		t.Fatalf("read live logs: %v", err)
	}
	if len(logs) == 0 {
		t.Fatal("no live log entries were published")
	}
	for _, entry := range logs {
		if strings.Contains(fmt.Sprint(entry.Values["output"]), "super-secret-value") {
			t.Fatalf("live log leaks a sensitive variable: %v", entry.Values)
		}
	}
}

func deref(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
