package executor

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

func TestRedactorReplacesSensitiveValues(t *testing.T) {
	r := newRedactor([]string{"super-secret-value", "super-secret-value", ""})
	got := r.redact("token=super-secret-value done")
	if got != "token=[REDACTED] done" {
		t.Fatalf("redact = %q, want token=[REDACTED] done", got)
	}
}

func TestRedactorPrefersLongestValue(t *testing.T) {
	r := newRedactor([]string{"abc", "abcdef"})
	if got := r.redact("xabcdef"); got != "x[REDACTED]" {
		t.Fatalf("redact = %q, want x[REDACTED]", got)
	}
}

func TestNilRedactorIsSafe(t *testing.T) {
	var r *redactor
	if got := r.redact("unchanged"); got != "unchanged" {
		t.Fatalf("redact = %q, want unchanged", got)
	}
}

func TestLogSinkRedactsSecrets(t *testing.T) {
	pub := &captureLogPublisher{}
	sink := newLogSink(context.Background(), pub, nil, 4, "step-1", newRedactor([]string{"hunter2"}))
	sink.WriteString("password=hunter2\n")
	sink.Flush()

	if strings.Contains(sink.String(), "hunter2") {
		t.Fatalf("stored output leaks the secret: %q", sink.String())
	}
	if len(pub.lines) != 1 || strings.Contains(pub.lines[0], "hunter2") {
		t.Fatalf("published lines = %v, want one redacted line", pub.lines)
	}
	if !strings.Contains(sink.String(), "[REDACTED]") {
		t.Fatalf("stored output was not redacted: %q", sink.String())
	}
}

func TestJobContextTimesOut(t *testing.T) {
	e := New(Deps{JobTimeout: 20 * time.Millisecond})
	ctx, cancel := e.jobContext(context.Background())
	defer cancel()

	<-ctx.Done()
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("job context error = %v, want deadline exceeded", ctx.Err())
	}
}

func TestJobContextWithoutTimeoutOnlyCancels(t *testing.T) {
	e := New(Deps{})
	ctx, cancel := e.jobContext(context.Background())
	cancel()
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatalf("job context error = %v, want canceled", ctx.Err())
	}
}

func TestCleanupStaleWorkspaces(t *testing.T) {
	root := t.TempDir()
	old := filepath.Join(root, "terrakube-job-1-old")
	recent := filepath.Join(root, "terrakube-job-2-new")
	other := filepath.Join(root, "keep-me")
	for _, dir := range []string{old, recent, other} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatalf("create %s: %v", dir, err)
		}
	}
	past := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(old, past, past); err != nil {
		t.Fatalf("age old workspace: %v", err)
	}
	if err := os.Chtimes(other, past, past); err != nil {
		t.Fatalf("age unrelated directory: %v", err)
	}

	if err := cleanupStaleWorkspaces(root, time.Hour); err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	if _, err := os.Stat(old); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stale workspace still exists: %v", err)
	}
	if _, err := os.Stat(recent); err != nil {
		t.Fatalf("recent workspace was removed: %v", err)
	}
	if _, err := os.Stat(other); err != nil {
		t.Fatalf("unrelated directory was removed: %v", err)
	}
}

func TestEnvironmentSeparatesSecrets(t *testing.T) {
	st := &variableStore{
		vars: []store.Variable{
			{VariableKey: "FOO", VariableValue: "env-value", VariableCategory: "ENV"},
			{VariableKey: "NAME", VariableValue: "var-value", VariableCategory: "TERRAFORM"},
			{VariableKey: "TOKEN", VariableValue: "super-secret-value", VariableCategory: "ENV", Sensitive: true},
			{VariableKey: "SKIP", VariableValue: "nope", VariableCategory: "ENV", Incomplete: true},
		},
		globals: []store.Globalvar{
			// Same key and category as the workspace variable: the workspace wins.
			{VariableKey: "FOO", VariableValue: "global-foo", VariableCategory: "ENV", OrganizationID: "org-1"},
			{VariableKey: "ONLY", VariableValue: "global-value", VariableCategory: "ENV", OrganizationID: "org-1"},
			{VariableKey: "REGION", VariableValue: "eu", VariableCategory: "TERRAFORM", OrganizationID: "org-1"},
			{VariableKey: "SECRET", VariableValue: "global-secret", VariableCategory: "ENV", Sensitive: true, OrganizationID: "org-1"},
		},
	}
	e := New(Deps{
		Store:               st,
		InternalTokenSecret: base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{0x42}, 32)),
		Logger:              testLogger(),
	})

	env, secrets, err := e.environment(context.Background(),
		store.Job{ID: 7, OrganizationID: "org-1", WorkspaceID: "ws-1"},
		store.Workspace{ID: "ws-1"},
		store.Step{ID: "step-1"})
	if err != nil {
		t.Fatalf("environment: %v", err)
	}

	joined := strings.Join(env, "\n")
	for _, want := range []string{
		"FOO=env-value", "TF_VAR_NAME=var-value", "TOKEN=super-secret-value",
		"ONLY=global-value", "TF_VAR_REGION=eu", "SECRET=global-secret",
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("environment is missing %q:\n%s", want, joined)
		}
	}
	if strings.Contains(joined, "SKIP=") {
		t.Errorf("incomplete variable was injected:\n%s", joined)
	}
	for _, want := range []string{"super-secret-value", "global-secret"} {
		if !contains(secrets, want) {
			t.Errorf("secret %q is not marked secret: %v", want, secrets)
		}
	}
	if n := strings.Count(joined, "FOO="); n != 1 {
		t.Errorf("FOO was injected %d times, want the workspace value only:\n%s", n, joined)
	}
	password := envValue(env, "TF_HTTP_PASSWORD")
	if password == "" {
		t.Fatal("internal token was not injected")
	}
	if !contains(secrets, password) {
		t.Errorf("internal token is not marked secret")
	}
}

func TestRunCommandDoesNotInheritExecutorSecrets(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "dump.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nenv > env.txt\n"), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	t.Setenv("TERRAKUBE_DATABASE_URL", "postgres://user:secret@host/db")

	e := New(Deps{})
	sink := newLogSink(context.Background(), nil, nil, 1, "step-1", nil)
	if err := e.runCommand(context.Background(), dir, []string{"E2E_MARKER=job-value"}, sink, script); err != nil {
		t.Fatalf("run command: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "env.txt"))
	if err != nil {
		t.Fatalf("read environment: %v", err)
	}
	env := string(data)
	if strings.Contains(env, "TERRAKUBE_DATABASE_URL") {
		t.Errorf("command inherited the executor environment:\n%s", env)
	}
	if !strings.Contains(env, "E2E_MARKER=job-value") {
		t.Errorf("command did not receive the job environment:\n%s", env)
	}
	if !strings.Contains(env, "PATH=") {
		t.Errorf("command did not receive PATH:\n%s", env)
	}
}

func TestCommandErrorDuringShutdownKeepsStepRunning(t *testing.T) {
	st := &cancelRunStore{}
	e := New(Deps{Store: st, Logger: testLogger()})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sink := newLogSink(context.Background(), nil, nil, 1, "step-1", nil)
	err := e.handleCommandError(ctx, ctx, store.Job{ID: 1}, store.Step{ID: "step-1"}, sink, "terraform init", context.Canceled)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context canceled", err)
	}
	if got := st.stepStatus(); got != "" {
		t.Fatalf("step status = %q; a shutdown must leave the step running for resume", got)
	}
}

func TestRunPlanTimesOut(t *testing.T) {
	ctx := context.Background()
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	archive := tarGz(t, map[string]string{"main.tf": `resource "terraform_data" "x" {}`})
	if err := backend.Put(ctx, layout.Content("cv-timeout"), bytes.NewReader(archive), int64(len(archive)), "application/gzip"); err != nil {
		t.Fatalf("seed configuration: %v", err)
	}

	script := filepath.Join(t.TempDir(), "tofu")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nexec sleep 10\n"), 0o755); err != nil {
		t.Fatalf("write fake tofu: %v", err)
	}

	st := &cancelRunStore{
		job: store.Job{
			ID:             11,
			OrganizationID: "org-1",
			WorkspaceID:    "ws-1",
			Status:         "running",
			OverrideSource: ptr("https://host/remote/tfe/v2/configuration-versions/cv-timeout/terraformContent.tar.gz"),
		},
		workspace: store.Workspace{ID: "ws-1"},
	}
	e := New(Deps{
		Store:      st,
		Storage:    backend,
		Layout:     layout,
		Logger:     testLogger(),
		Binary:     script,
		WorkDir:    t.TempDir(),
		JobTimeout: 75 * time.Millisecond,
	})
	e.cancelPollInterval = 5 * time.Millisecond

	start := time.Now()
	err := e.runPlan(ctx, st.job, st.workspace, store.Step{ID: "step-1", JobID: 11, StepNumber: 100, Status: ptr("pending")})
	elapsed := time.Since(start)

	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("runPlan error = %v, want a timeout error", err)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("runPlan took %s; the terraform command was not stopped", elapsed)
	}
	if got := st.stepStatus(); got != "failed" {
		t.Fatalf("step status = %q, want failed", got)
	}
}

type captureLogPublisher struct {
	lines []string
}

func (p *captureLogPublisher) PublishCapped(_ context.Context, _ string, values map[string]any, _ int64) (string, error) {
	p.lines = append(p.lines, fmt.Sprint(values["output"]))
	return "1-1", nil
}

type variableStore struct {
	store.Querier
	vars    []store.Variable
	globals []store.Globalvar
}

func (s *variableStore) ListVariablesByWorkspace(context.Context, string) ([]store.Variable, error) {
	return s.vars, nil
}

func (s *variableStore) ListGlobalVariablesByOrganization(_ context.Context, orgID string) ([]store.Globalvar, error) {
	var out []store.Globalvar
	for _, variable := range s.globals {
		if variable.OrganizationID == orgID {
			out = append(out, variable)
		}
	}
	return out, nil
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func envValue(env []string, key string) string {
	prefix := key + "="
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			return strings.TrimPrefix(entry, prefix)
		}
	}
	return ""
}
