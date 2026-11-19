package executor

import (
	"context"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/queue"
)

// JobQueue is the subset of the Redis Streams client the executor needs. It is
// an interface so reliability behavior (ack, retry, shutdown) can be tested
// without Redis.
type JobQueue interface {
	EnsureGroup(ctx context.Context, stream, group string) error
	ClaimStale(ctx context.Context, stream, group, consumer string, minIdle time.Duration, count int64) ([]queue.Message, error)
	ReadGroup(ctx context.Context, stream, group, consumer string, count int64, block time.Duration) ([]queue.Message, error)
	Ack(ctx context.Context, stream, group, id string) error
	Delete(ctx context.Context, stream, id string) error
	DeleteStream(ctx context.Context, stream string) error
	Publish(ctx context.Context, stream string, values map[string]any) (string, error)
	PublishCapped(ctx context.Context, stream string, values map[string]any, maxLen int64) (string, error)
}

// LogPublisher fans live log lines out to Redis.
type LogPublisher interface {
	PublishCapped(ctx context.Context, stream string, values map[string]any, maxLen int64) (string, error)
}

var (
	_ JobQueue     = (*queue.Queue)(nil)
	_ LogPublisher = (*queue.Queue)(nil)
)
