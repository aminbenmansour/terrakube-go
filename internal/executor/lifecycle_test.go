package executor

import (
	"context"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

func TestAckDeletesTheJobMessage(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "completed"}}
	q := &fakeQueue{}
	e := newTestExecutor(st, q)

	e.ack(context.Background(), "1-1")

	if len(q.acks) != 1 {
		t.Fatalf("acks = %d, want 1", len(q.acks))
	}
	if len(q.deleted) != 1 || q.deleted[0] != "1-1" {
		t.Fatalf("deleted = %v, want [1-1]", q.deleted)
	}
}

func TestHandleDropsLiveLogsWhenTheJobEnds(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "completed"}}
	q := &fakeQueue{}
	e := newTestExecutor(st, q)

	e.handle(context.Background(), jobMessage("1-1", 7, 0))

	if len(q.dropped) != 1 || q.dropped[0] != "7" {
		t.Fatalf("dropped streams = %v, want [7]", q.dropped)
	}
}

func TestHandleKeepsLiveLogsOnShutdown(t *testing.T) {
	st := &reliabilityStore{job: store.Job{ID: 7, Status: "running"}}
	q := &fakeQueue{}
	e := newTestExecutor(st, q)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	e.handle(ctx, jobMessage("1-1", 7, 0))

	if len(q.dropped) != 0 {
		t.Fatalf("dropped streams = %v, want none; the resumed job needs them", q.dropped)
	}
}
