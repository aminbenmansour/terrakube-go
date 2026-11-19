package executor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// fakeLocker grants a job to the first caller and refuses the second, which is
// what the Postgres advisory lock does across executors.
type fakeLocker struct {
	held     map[int64]bool
	releases int
}

func newFakeLocker() *fakeLocker { return &fakeLocker{held: map[int64]bool{}} }

func (l *fakeLocker) TryLock(_ context.Context, jobID int64) (func(), bool, error) {
	if l.held[jobID] {
		return nil, false, nil
	}
	l.held[jobID] = true
	return func() {
		l.releases++
		delete(l.held, jobID)
	}, true, nil
}

func TestHandleSkipsAJobHeldByAnotherExecutor(t *testing.T) {
	st := &reliabilityStore{
		job:        store.Job{ID: 7, Status: "running"},
		getJobErrs: []error{errors.New("the held job must not be loaded")},
	}
	q := &fakeQueue{}
	locker := newFakeLocker()
	locker.held[7] = true

	e := New(Deps{Store: st, Queue: q, Logger: testLogger(), Locker: locker})
	e.retryDelay = func(int) time.Duration { return 0 }

	e.handle(context.Background(), jobMessage("1-1", 7, 0))

	if got := len(q.acks); got != 1 {
		t.Fatalf("acks = %d, want 1", got)
	}
	if got := len(q.published); got != 0 {
		t.Fatalf("published = %d, want 0", got)
	}
	if st.getJobCalls != 0 {
		t.Fatalf("job was loaded %d times; a held job must not be processed", st.getJobCalls)
	}
	if locker.releases != 0 {
		t.Fatalf("releases = %d, want 0; the lock belongs to the other executor", locker.releases)
	}
}

func TestHandleReleasesTheJobLock(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "completed"}}
	q := &fakeQueue{}
	locker := newFakeLocker()

	e := New(Deps{Store: st, Queue: q, Logger: testLogger(), Locker: locker})
	e.handle(context.Background(), jobMessage("1-1", 7, 0))

	if locker.releases != 1 {
		t.Fatalf("releases = %d, want 1", locker.releases)
	}
	if len(locker.held) != 0 {
		t.Fatalf("job is still held after the run")
	}
}

func TestHandleWithoutALockerRunsTheJob(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "completed"}}
	q := &fakeQueue{}
	e := newTestExecutor(st, q)

	e.handle(context.Background(), jobMessage("1-1", 7, 0))

	if got := len(q.acks); got != 1 {
		t.Fatalf("acks = %d, want 1", got)
	}
	if st.getJobCalls == 0 {
		t.Fatal("job was not loaded without a locker")
	}
}

func TestReclaimStaleHandlesClaimedJobs(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "completed"}}
	q := &fakeQueue{stale: []queue.Message{jobMessage("1-1", 7, 0)}}
	e := newTestExecutor(st, q)

	e.reclaimStale(context.Background())

	if got := len(q.acks); got != 1 {
		t.Fatalf("acks = %d, want 1", got)
	}
	if st.getJobCalls == 0 {
		t.Fatal("reclaimed job was not processed")
	}
}

func TestStaleClaimIdleOutlivesTheJobTimeout(t *testing.T) {
	cases := []struct {
		timeout time.Duration
		want    time.Duration
	}{
		{2 * time.Hour, 2*time.Hour + 15*time.Minute},
		{30 * time.Minute, 45 * time.Minute},
		{0, 24 * time.Hour},
	}
	for _, tc := range cases {
		e := New(Deps{JobTimeout: tc.timeout, Logger: testLogger()})
		if got := e.staleClaimIdle(); got != tc.want {
			t.Errorf("staleClaimIdle(%s) = %s, want %s", tc.timeout, got, tc.want)
		}
	}
}
