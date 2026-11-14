// Package queue wraps Redis Streams, which Terrakube Go uses for the run queue
// and for fan-out of live executor logs.
package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Queue is a thin Redis Streams client.
type Queue struct {
	client *redis.Client
}

// New parses a redis:// URL and returns a client. It does not connect eagerly;
// call Ping during startup to fail fast.
func New(redisURL string) (*Queue, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("queue: parse redis URL: %w", err)
	}
	return &Queue{client: redis.NewClient(opts)}, nil
}

// Ping verifies connectivity.
func (q *Queue) Ping(ctx context.Context) error {
	if err := q.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("queue: ping redis: %w", err)
	}
	return nil
}

// Close releases the connection pool.
func (q *Queue) Close() error {
	return q.client.Close()
}

// Publish appends an entry to a stream and returns its message ID.
func (q *Queue) Publish(ctx context.Context, stream string, values map[string]any) (string, error) {
	id, err := q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Result()
	if err != nil {
		return "", fmt.Errorf("queue: publish to %s: %w", stream, err)
	}
	return id, nil
}

// EnsureGroup creates a consumer group if it does not exist yet.
func (q *Queue) EnsureGroup(ctx context.Context, stream, group string) error {
	err := q.client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("queue: ensure group %s on %s: %w", group, stream, err)
	}
	return nil
}

// Message is one stream entry.
type Message struct {
	ID     string
	Values map[string]any
}

// ReadGroup reads up to count new entries for a consumer, blocking up to block.
func (q *Queue) ReadGroup(ctx context.Context, stream, group, consumer string, count int64, block time.Duration) ([]Message, error) {
	streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{stream, ">"},
		Count:    count,
		Block:    block,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("queue: read group %s on %s: %w", group, stream, err)
	}

	var messages []Message
	for _, s := range streams {
		for _, m := range s.Messages {
			messages = append(messages, Message{ID: m.ID, Values: m.Values})
		}
	}
	return messages, nil
}

// Ack acknowledges a processed message.
func (q *Queue) Ack(ctx context.Context, stream, group, id string) error {
	if err := q.client.XAck(ctx, stream, group, id).Err(); err != nil {
		return fmt.Errorf("queue: ack %s on %s: %w", id, stream, err)
	}
	return nil
}

// Client exposes the underlying client for callers that need a primitive the
// wrapper does not cover yet (for example additional stream trimming).
func (q *Queue) Client() *redis.Client {
	return q.client
}
