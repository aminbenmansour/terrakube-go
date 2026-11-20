// Package jobstatus publishes run status changes so the React UI's GraphQL
// subscriptions can push them, and adapts the store so every status update is
// published without touching the call sites.
package jobstatus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Channel is the Redis pub/sub channel carrying status events. It is an
// internal contract, not part of the upstream HTTP API.
const Channel = "terrakube:job-status"

// Event is what a subscriber receives. JobID is Int in the GraphQL schema, so
// it fits in an int.
type Event struct {
	JobID          int64  `json:"jobId"`
	WorkspaceID    string `json:"workspaceId"`
	OrganizationID string `json:"organizationId"`
	Status         string `json:"status"`
}

// Bus publishes and subscribes to status events.
type Bus struct {
	client *redis.Client
}

// NewBus parses a redis:// URL. The connection is opened lazily.
func NewBus(redisURL string) (*Bus, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("jobstatus: parse redis URL: %w", err)
	}
	return &Bus{client: redis.NewClient(opts)}, nil
}

// Close releases the connection pool.
func (b *Bus) Close() error { return b.client.Close() }

// Publish sends one event. A subscriber that is not connected loses the event;
// this is live-update traffic, not durable state.
func (b *Bus) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("jobstatus: encode event: %w", err)
	}
	if err := b.client.Publish(ctx, Channel, payload).Err(); err != nil {
		return fmt.Errorf("jobstatus: publish: %w", err)
	}
	return nil
}

// SubscribeJobStatus returns a channel of events. It is closed when ctx is
// cancelled or Redis goes away.
func (b *Bus) SubscribeJobStatus(ctx context.Context) (<-chan Event, error) {
	pubsub := b.client.Subscribe(ctx, Channel)
	if _, err := pubsub.Receive(ctx); err != nil {
		_ = pubsub.Close()
		return nil, fmt.Errorf("jobstatus: subscribe: %w", err)
	}

	events := make(chan Event)
	go func() {
		defer close(events)
		defer pubsub.Close()
		for {
			message, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				return
			}
			var event Event
			if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
				continue
			}
			select {
			case events <- event:
			case <-ctx.Done():
				return
			}
		}
	}()
	return events, nil
}

// Publisher is the part of the bus the store wrapper needs, so a test can
// supply a recorder.
type Publisher func(ctx context.Context, event Event) error

// Wrap returns a Querier that publishes the job after every successful status
// update. Publishing is best effort: a Redis outage must not fail a run.
func Wrap(querier store.Querier, publish Publisher) store.Querier {
	return publishingStore{Querier: querier, publish: publish}
}

type publishingStore struct {
	store.Querier
	publish Publisher
}

func (s publishingStore) UpdateJobStatus(ctx context.Context, arg store.UpdateJobStatusParams) (store.Job, error) {
	job, err := s.Querier.UpdateJobStatus(ctx, arg)
	if err != nil {
		return job, err
	}
	_ = s.publish(ctx, Event{
		JobID:          job.ID,
		WorkspaceID:    job.WorkspaceID,
		OrganizationID: job.OrganizationID,
		Status:         job.Status,
	})
	return job, nil
}
