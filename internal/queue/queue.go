// Package queue wraps Redis Streams, which Terrakube Go uses for the run queue
// and for fan-out of live executor logs.
package queue

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// JobsStream is the Redis Stream the API publishes run IDs to and the executor
// consumes. It is an internal contract, not part of the upstream HTTP API.
const JobsStream = "terrakube:jobs"

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

// Depth reports how many entries a stream currently holds.
func (q *Queue) Depth(ctx context.Context, stream string) (int64, error) {
	depth, err := q.client.XLen(ctx, stream).Result()
	if err != nil {
		return 0, fmt.Errorf("queue: depth of %s: %w", stream, err)
	}
	return depth, nil
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

// PublishCapped appends to a stream and trims it to roughly maxLen entries.
// It is used for live log fan-out so a long run cannot grow Redis without
// bound.
func (q *Queue) PublishCapped(ctx context.Context, stream string, values map[string]any, maxLen int64) (string, error) {
	id, err := q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		MaxLen: maxLen,
		Approx: true,
		Values: values,
	}).Result()
	if err != nil {
		return "", fmt.Errorf("queue: publish to %s: %w", stream, err)
	}
	return id, nil
}

// ClaimStale takes over entries that were delivered to a consumer that stopped
// before acknowledging them. It returns nil when nothing was pending.
func (q *Queue) ClaimStale(ctx context.Context, stream, group, consumer string, minIdle time.Duration, count int64) ([]Message, error) {
	entries, _, err := q.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream:   stream,
		Group:    group,
		Consumer: consumer,
		MinIdle:  minIdle,
		Start:    "0-0",
		Count:    count,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("queue: claim pending on %s: %w", stream, err)
	}
	messages := make([]Message, 0, len(entries))
	for _, entry := range entries {
		messages = append(messages, Message{ID: entry.ID, Values: entry.Values})
	}
	return messages, nil
}

// EnqueueJob publishes a run to the executor queue. The message carries only
// the job ID; the executor reads the rest of the job from Postgres.
func (q *Queue) EnqueueJob(ctx context.Context, jobID int64) (string, error) {
	return q.Publish(ctx, JobsStream, map[string]any{
		"job_id": strconv.FormatInt(jobID, 10),
	})
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

// ReadStream reads entries after lastID from a plain (non-group) stream.
// It returns nil when the block timeout expires.
func (q *Queue) ReadStream(ctx context.Context, stream, lastID string, count int64, block time.Duration) ([]Message, error) {
	streams, err := q.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{stream, lastID},
		Count:   count,
		Block:   block,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("queue: read %s: %w", stream, err)
	}

	var messages []Message
	for _, s := range streams {
		for _, m := range s.Messages {
			messages = append(messages, Message{ID: m.ID, Values: m.Values})
		}
	}
	return messages, nil
}

// FirstStreamID returns the oldest retained entry ID, or "" when empty.
func (q *Queue) FirstStreamID(ctx context.Context, stream string) (string, error) {
	entries, err := q.client.XRangeN(ctx, stream, "-", "+", 1).Result()
	if err != nil {
		return "", fmt.Errorf("queue: first entry of %s: %w", stream, err)
	}
	if len(entries) == 0 {
		return "", nil
	}
	return entries[0].ID, nil
}

// HeadStream returns up to count of the oldest entries in chronological order.
func (q *Queue) HeadStream(ctx context.Context, stream string, count int64) ([]Message, error) {
	entries, err := q.client.XRangeN(ctx, stream, "-", "+", count).Result()
	if err != nil {
		return nil, fmt.Errorf("queue: head %s: %w", stream, err)
	}
	messages := make([]Message, 0, len(entries))
	for _, entry := range entries {
		messages = append(messages, Message{ID: entry.ID, Values: entry.Values})
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

// Delete removes an entry from the stream. Acknowledgement only clears the
// pending list; without this the entry stays in Redis forever.
func (q *Queue) Delete(ctx context.Context, stream, id string) error {
	if err := q.client.XDel(ctx, stream, id).Err(); err != nil {
		return fmt.Errorf("queue: delete %s on %s: %w", id, stream, err)
	}
	return nil
}

// DeleteStream removes a whole stream, used for the per-job live log streams.
func (q *Queue) DeleteStream(ctx context.Context, stream string) error {
	if err := q.client.Del(ctx, stream).Err(); err != nil {
		return fmt.Errorf("queue: delete stream %s: %w", stream, err)
	}
	return nil
}
