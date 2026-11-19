package executor

import (
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
