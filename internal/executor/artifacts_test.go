package executor

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

func TestCompleteStepArchivesLog(t *testing.T) {
	ctx := context.Background()
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	st := &artifactStore{}
	e := New(Deps{Store: st, Storage: backend, Layout: layout, Logger: testLogger()})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	sink.WriteString("line one\nline two\n")
	err := e.completeStep(ctx, store.Job{ID: 9, OrganizationID: "org-1"}, store.Step{ID: "step-1"}, sink, "completed")
	if err != nil {
		t.Fatalf("completeStep: %v", err)
	}

	got := string(backend.objects[layout.Output("org-1", 9, "step-1")])
	if got != "line one\nline two\n" {
		t.Fatalf("archived log = %q", got)
	}
	if st.stepStatus != "completed" {
		t.Errorf("step status = %q, want completed", st.stepStatus)
	}
	if st.stepOutput == nil || *st.stepOutput != "line one\nline two\n" {
		t.Errorf("inline step output = %v", st.stepOutput)
	}
}

func TestCompleteStepLogUploadFailureFailsStep(t *testing.T) {
	ctx := context.Background()
	e := New(Deps{
		Store:   &artifactStore{},
		Storage: failingBackend{err: errors.New("storage down")},
		Layout:  tfstate.NewLayout("local"),
		Logger:  testLogger(),
	})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	sink.WriteString("output\n")
	err := e.completeStep(ctx, store.Job{ID: 9}, store.Step{ID: "step-1"}, sink, "completed")
	if err == nil || !strings.Contains(err.Error(), "archive step log") {
		t.Fatalf("error = %v, want archive failure", err)
	}
	if st := e.deps.Store.(*artifactStore); st.stepStatus != "failed" {
		t.Errorf("step status = %q, want failed", st.stepStatus)
	}
}

func TestCompleteStepCancelledKeepsStatusWhenLogUploadFails(t *testing.T) {
	ctx := context.Background()
	st := &artifactStore{}
	e := New(Deps{
		Store:   st,
		Storage: failingBackend{err: errors.New("storage down")},
		Layout:  tfstate.NewLayout("local"),
		Logger:  testLogger(),
	})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	if err := e.completeStep(ctx, store.Job{ID: 9}, store.Step{ID: "step-1"}, sink, "cancelled"); err != nil {
		t.Fatalf("completeStep: %v", err)
	}
	if st.stepStatus != "cancelled" {
		t.Errorf("step status = %q, want cancelled", st.stepStatus)
	}
}

func TestArchiveHistoryDoesNotCreateRowWhenUploadFails(t *testing.T) {
	ctx := context.Background()
	layout := tfstate.NewLayout("local")
	mem := newMemoryBackend()
	state := []byte(`{"version":4,"serial":3,"lineage":"lineage-1"}`)
	if err := mem.Put(ctx, layout.CurrentState("org-1", "ws-1"), bytes.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("seed state: %v", err)
	}
	backend := &failingKeyBackend{memoryBackend: mem, failSuffix: ".raw.json"}
	st := &artifactStore{}
	e := New(Deps{Store: st, Storage: backend, Layout: layout, Logger: testLogger(), PublicURL: "https://example.invalid"})

	err := e.archiveHistory(ctx, store.Job{ID: 9, OrganizationID: "org-1", WorkspaceID: "ws-1"}, store.Workspace{ID: "ws-1"})
	if err == nil {
		t.Fatal("archiveHistory succeeded with a failing raw-state upload")
	}
	if st.created != nil {
		t.Fatalf("created history row %q before all state objects were stored", st.created.ID)
	}
}

func TestStateJSONURL(t *testing.T) {
	job := store.Job{OrganizationID: "org-1"}
	ws := store.Workspace{ID: "ws-1"}
	want := "https://terrakube.example.com/tfstate/v1/organization/org-1/workspace/ws-1/state/h-1.json"
	cases := map[string]string{
		"":                               "",
		"/":                              "",
		"https://terrakube.example.com":  want,
		"https://terrakube.example.com/": want,
	}
	for publicURL, expected := range cases {
		e := New(Deps{PublicURL: publicURL})
		if got := e.stateJSONURL(job, ws, "h-1"); got != expected {
			t.Errorf("stateJSONURL(%q) = %q, want %q", publicURL, got, expected)
		}
	}
}

type failingKeyBackend struct {
	*memoryBackend
	failSuffix string
}

func (b *failingKeyBackend) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	if strings.HasSuffix(key, b.failSuffix) {
		return errors.New("storage down")
	}
	return b.memoryBackend.Put(ctx, key, r, size, contentType)
}

func TestFinishApplyFailsWhenHistoryArchiveFails(t *testing.T) {
	ctx := context.Background()
	backend := newMemoryBackend()
	layout := tfstate.NewLayout("local")
	state := []byte(`{"version":4,"serial":3,"lineage":"lineage-1"}`)
	if err := backend.Put(ctx, layout.CurrentState("org-1", "ws-1"), bytes.NewReader(state), int64(len(state)), "application/json"); err != nil {
		t.Fatalf("seed state: %v", err)
	}
	st := &artifactStore{historyErr: errors.New("history database down")}
	e := New(Deps{Store: st, Storage: backend, Layout: layout, Logger: testLogger()})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	err := e.finishApply(ctx, store.Job{ID: 9, OrganizationID: "org-1", WorkspaceID: "ws-1"}, store.Workspace{ID: "ws-1"}, store.Step{ID: "step-1"}, sink)
	if err == nil || !strings.Contains(err.Error(), "archive state history") {
		t.Fatalf("error = %v, want archive state history", err)
	}
	if st.stepStatus != "failed" {
		t.Errorf("step status = %q, want failed", st.stepStatus)
	}
}

func TestFinishApplyCompletesStep(t *testing.T) {
	ctx := context.Background()
	st := &artifactStore{}
	e := New(Deps{Store: st, Logger: testLogger()})

	sink := newLogSink(ctx, nil, testLogger(), 9, "step-1", nil)
	if err := e.finishApply(ctx, store.Job{ID: 9, OrganizationID: "org-1"}, store.Workspace{ID: "ws-1"}, store.Step{ID: "step-1"}, sink); err != nil {
		t.Fatalf("finishApply: %v", err)
	}
	if st.stepStatus != "completed" {
		t.Errorf("step status = %q, want completed", st.stepStatus)
	}
}

type artifactStore struct {
	store.Querier
	stepStatus string
	stepOutput *string
	jobOutput  *string
	historyErr error
	created    *store.CreateHistoryParams
}

func (s *artifactStore) UpdateStepStatus(_ context.Context, arg store.UpdateStepStatusParams) (store.Step, error) {
	if arg.Status != nil {
		s.stepStatus = *arg.Status
	}
	return store.Step{ID: arg.ID}, nil
}

func (s *artifactStore) UpdateStepOutput(_ context.Context, arg store.UpdateStepOutputParams) error {
	s.stepOutput = arg.Output
	return nil
}

func (s *artifactStore) UpdateJobOutput(_ context.Context, arg store.UpdateJobOutputParams) error {
	s.jobOutput = arg.Output
	return nil
}

func (s *artifactStore) GetHistoryByJobReference(context.Context, store.GetHistoryByJobReferenceParams) (store.History, error) {
	if s.historyErr != nil {
		return store.History{}, s.historyErr
	}
	return store.History{}, pgx.ErrNoRows
}

func (s *artifactStore) CreateHistory(_ context.Context, arg store.CreateHistoryParams) (store.History, error) {
	s.created = &arg
	return store.History{ID: arg.ID}, nil
}

type failingBackend struct{ err error }

func (f failingBackend) Put(context.Context, string, io.Reader, int64, string) error { return f.err }
func (f failingBackend) Get(context.Context, string) (io.ReadCloser, error)          { return nil, f.err }
func (f failingBackend) Delete(context.Context, string) error                        { return f.err }
func (f failingBackend) Stat(context.Context, string) (storage.Object, error) {
	return storage.Object{}, f.err
}
func (f failingBackend) Health(context.Context) error { return nil }

var _ storage.Backend = failingBackend{}
