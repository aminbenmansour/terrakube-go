package schedule

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type fakeStore struct {
	store.Querier
	schedules []store.Schedule
	workspace store.Workspace
	templates []store.Template
	jobs      []store.CreateJobParams
	steps     []store.CreateStepParams
	failed    []int64
	lastJob   []string
	nextID    int64
}

func (f *fakeStore) ListEnabledSchedules(context.Context) ([]store.Schedule, error) {
	return f.schedules, nil
}

func (f *fakeStore) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	if id == f.workspace.ID {
		return f.workspace, nil
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (f *fakeStore) GetTemplate(_ context.Context, id string) (store.Template, error) {
	for _, template := range f.templates {
		if template.ID == id {
			return template, nil
		}
	}
	return store.Template{}, pgx.ErrNoRows
}

func (f *fakeStore) CreateJob(_ context.Context, arg store.CreateJobParams) (store.Job, error) {
	f.jobs = append(f.jobs, arg)
	f.nextID++
	return store.Job{ID: f.nextID, OrganizationID: arg.OrganizationID, WorkspaceID: arg.WorkspaceID, Status: arg.Status}, nil
}

func (f *fakeStore) CreateStep(_ context.Context, arg store.CreateStepParams) (store.Step, error) {
	f.steps = append(f.steps, arg)
	return store.Step{ID: arg.ID, JobID: arg.JobID, StepNumber: arg.StepNumber}, nil
}

func (f *fakeStore) SetWorkspaceLastJob(_ context.Context, arg store.SetWorkspaceLastJobParams) error {
	f.lastJob = append(f.lastJob, arg.LastJobStatus)
	return nil
}

func (f *fakeStore) UpdateJobStatus(_ context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	f.failed = append(f.failed, arg.ID)
	return store.Job{ID: arg.ID, Status: arg.Status}, nil
}

type fakeQueue struct {
	enqueued []int64
	err      error
}

func (f *fakeQueue) EnqueueJob(_ context.Context, jobID int64) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	f.enqueued = append(f.enqueued, jobID)
	return "1-0", nil
}

func newDispatcherFixture(expression string) (*Dispatcher, *fakeStore, *fakeQueue) {
	ws := store.Workspace{
		ID:             "ws-1",
		OrganizationID: "org-1",
		Name:           "scheduled",
	}
	fs := &fakeStore{
		workspace: ws,
		schedules: []store.Schedule{{
			ID:          "schedule-1",
			Cron:        &expression,
			Enabled:     true,
			WorkspaceID: ws.ID,
		}},
	}
	fq := &fakeQueue{}
	d := NewDispatcher(fs, fq, slog.New(slog.NewTextHandler(io.Discard, nil)))
	return d, fs, fq
}

func tickAt(d *Dispatcher, at time.Time) {
	d.Now = func() time.Time { return at }
	d.Tick(context.Background())
}

func TestTickFiresOncePerMatchingSecond(t *testing.T) {
	d, fs, fq := newDispatcherFixture("0 * * * * ?")

	start := time.Date(2026, 11, 26, 9, 59, 50, 0, time.UTC)
	tickAt(d, start)                     // records the clock only
	tickAt(d, start.Add(10*time.Second)) // 10:00:00 falls inside the window
	tickAt(d, start.Add(20*time.Second)) // 10:00:10, no match
	if len(fs.jobs) != 1 {
		t.Fatalf("jobs = %d, want 1; body %+v", len(fs.jobs), fs.jobs)
	}
	if len(fq.enqueued) != 1 {
		t.Fatalf("enqueued = %v, want one job", fq.enqueued)
	}
	if len(fs.steps) != 3 {
		t.Fatalf("steps = %d, want plan, approval and apply", len(fs.steps))
	}
	if fs.jobs[0].Via != "SCHEDULE" || fs.jobs[0].Status != "pending" {
		t.Errorf("job = via %q status %q", fs.jobs[0].Via, fs.jobs[0].Status)
	}
	if fs.jobs[0].CreatedBy == nil || *fs.jobs[0].CreatedBy != "serviceAccount" {
		t.Errorf("createdBy = %v", fs.jobs[0].CreatedBy)
	}
	if len(fs.lastJob) != 1 || fs.lastJob[0] != "pending" {
		t.Errorf("last job status = %v", fs.lastJob)
	}

	// The next minute fires again; the same second never fires twice.
	tickAt(d, start.Add(70*time.Second))
	if len(fs.jobs) != 2 {
		t.Fatalf("jobs after the next minute = %d, want 2", len(fs.jobs))
	}
}

func TestTickSkipsLockedWorkspaceAndUnparseableCron(t *testing.T) {
	d, fs, _ := newDispatcherFixture("0 * * * * ?")
	fs.workspace.Locked = true

	start := time.Date(2026, 11, 26, 9, 59, 50, 0, time.UTC)
	tickAt(d, start)
	tickAt(d, start.Add(10*time.Second))
	if len(fs.jobs) != 0 {
		t.Fatalf("jobs = %d, want none for a locked workspace", len(fs.jobs))
	}

	// An expression Quartz refuses is logged once and never fires.
	bad := "0 0 0 L * ?"
	d, fs, _ = newDispatcherFixture(bad)
	tickAt(d, start)
	tickAt(d, start.Add(10*time.Second))
	tickAt(d, start.Add(20*time.Second))
	if len(fs.jobs) != 0 {
		t.Fatalf("jobs = %d, want none for an unsupported cron", len(fs.jobs))
	}
}

func TestTickUsesTheTemplateOfTheWorkspaceOrganization(t *testing.T) {
	expression := "0 * * * * ?"
	d, fs, _ := newDispatcherFixture(expression)
	template := "template-1"
	fs.schedules[0].TemplateReference = &template
	fs.templates = []store.Template{{ID: template, OrganizationID: "org-1", Tcl: strptr("terraform apply -auto-approve")}}

	start := time.Date(2026, 11, 26, 9, 59, 50, 0, time.UTC)
	tickAt(d, start)
	tickAt(d, start.Add(10*time.Second))
	if len(fs.jobs) != 1 {
		t.Fatalf("jobs = %d, want 1", len(fs.jobs))
	}
	if got := fs.jobs[0].TemplateReference; got == nil || *got != template {
		t.Errorf("templateReference = %v, want %s", got, template)
	}
	if got := fs.jobs[0].Tcl; got == nil || *got != "terraform apply -auto-approve" {
		t.Errorf("tcl = %v, want the template body", got)
	}

	// A template that moved to another organization is refused.
	d, fs, _ = newDispatcherFixture(expression)
	fs.schedules[0].TemplateReference = &template
	fs.templates = []store.Template{{ID: template, OrganizationID: "org-2"}}
	tickAt(d, start)
	tickAt(d, start.Add(10*time.Second))
	if len(fs.jobs) != 0 {
		t.Fatalf("jobs = %d, want none for a foreign template", len(fs.jobs))
	}
}

func TestTickMarksAJobFailedWhenItCannotBeQueued(t *testing.T) {
	d, fs, fq := newDispatcherFixture("0 * * * * ?")
	fq.err = errors.New("redis is down")

	start := time.Date(2026, 11, 26, 9, 59, 50, 0, time.UTC)
	tickAt(d, start)
	tickAt(d, start.Add(10*time.Second))
	if len(fs.jobs) != 1 {
		t.Fatalf("jobs = %d, want 1", len(fs.jobs))
	}
	if len(fs.failed) != 1 || fs.failed[0] != fs.nextID {
		t.Errorf("failed = %v, want job %d", fs.failed, fs.nextID)
	}
}
