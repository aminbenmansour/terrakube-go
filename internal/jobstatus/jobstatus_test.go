package jobstatus

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type stubQuerier struct {
	store.Querier
	job store.Job
	err error
}

func (s stubQuerier) UpdateJobStatus(context.Context, store.UpdateJobStatusParams) (store.Job, error) {
	return s.job, s.err
}

func TestWrapPublishesEveryStatusChange(t *testing.T) {
	published := make(chan Event, 1)
	job := store.Job{ID: 7, WorkspaceID: "ws-1", OrganizationID: "org-1", Status: "running"}
	querier := Wrap(stubQuerier{job: job}, func(_ context.Context, event Event) error {
		published <- event
		return nil
	})

	if _, err := querier.UpdateJobStatus(context.Background(), store.UpdateJobStatusParams{ID: 7, Status: "running"}); err != nil {
		t.Fatalf("UpdateJobStatus: %v", err)
	}
	select {
	case event := <-published:
		if event.JobID != 7 || event.WorkspaceID != "ws-1" || event.OrganizationID != "org-1" || event.Status != "running" {
			t.Fatalf("event = %#v", event)
		}
	default:
		t.Fatal("no event was published")
	}
}

func TestWrapIgnoresAPublishFailure(t *testing.T) {
	querier := Wrap(stubQuerier{job: store.Job{ID: 1}}, func(context.Context, Event) error {
		return errors.New("redis is down")
	})
	if _, err := querier.UpdateJobStatus(context.Background(), store.UpdateJobStatusParams{ID: 1}); err != nil {
		t.Fatalf("a publish failure must not fail the update: %v", err)
	}
}

func TestWrapDoesNotPublishAFailedUpdate(t *testing.T) {
	querier := Wrap(stubQuerier{err: errors.New("constraint")}, func(context.Context, Event) error {
		t.Fatal("published after a failed update")
		return nil
	})
	if _, err := querier.UpdateJobStatus(context.Background(), store.UpdateJobStatusParams{ID: 1}); err == nil {
		t.Fatal("error was swallowed")
	}
}

func TestBusRoundTrip(t *testing.T) {
	url := os.Getenv("TERRAKUBE_TEST_REDIS_URL")
	if url == "" {
		t.Skip("TERRAKUBE_TEST_REDIS_URL is not set")
	}
	bus, err := NewBus(url)
	if err != nil {
		t.Fatalf("NewBus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	events, err := bus.SubscribeJobStatus(ctx)
	if err != nil {
		t.Fatalf("SubscribeJobStatus: %v", err)
	}
	if err := bus.Publish(ctx, Event{JobID: 42, WorkspaceID: "ws", OrganizationID: "org", Status: "completed"}); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	select {
	case event := <-events:
		if event.JobID != 42 || event.Status != "completed" || event.WorkspaceID != "ws" {
			t.Fatalf("event = %#v", event)
		}
	case <-ctx.Done():
		t.Fatal("no event arrived")
	}
}
