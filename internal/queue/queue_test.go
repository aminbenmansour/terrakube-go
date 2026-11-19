package queue

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func testQueue(t *testing.T) *Queue {
	t.Helper()
	url := os.Getenv("TERRAKUBE_TEST_REDIS_URL")
	if url == "" {
		t.Skip("TERRAKUBE_TEST_REDIS_URL is not set")
	}
	q, err := New(url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = q.Close() })
	return q
}

func TestHeadStreamReturnsTheOldestEntries(t *testing.T) {
	q := testQueue(t)
	ctx := context.Background()
	stream := fmt.Sprintf("test-head-%d", time.Now().UnixNano())
	defer func() { _ = q.DeleteStream(ctx, stream) }()

	for i := 0; i < 3; i++ {
		if _, err := q.Publish(ctx, stream, map[string]any{"n": i}); err != nil {
			t.Fatalf("publish %d: %v", i, err)
		}
	}

	messages, err := q.HeadStream(ctx, stream, 2)
	if err != nil {
		t.Fatalf("head: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("messages = %d, want 2", len(messages))
	}
	if got := fmt.Sprint(messages[0].Values["n"]); got != "0" {
		t.Fatalf("first message = %s, want the oldest (0)", got)
	}
}

func TestDeleteRemovesAMessage(t *testing.T) {
	q := testQueue(t)
	ctx := context.Background()
	stream := fmt.Sprintf("test-del-%d", time.Now().UnixNano())
	defer func() { _ = q.DeleteStream(ctx, stream) }()

	id, err := q.Publish(ctx, stream, map[string]any{"n": 1})
	if err != nil {
		t.Fatalf("publish: %v", err)
	}
	if err := q.Delete(ctx, stream, id); err != nil {
		t.Fatalf("delete: %v", err)
	}

	messages, err := q.HeadStream(ctx, stream, 10)
	if err != nil {
		t.Fatalf("head: %v", err)
	}
	if len(messages) != 0 {
		t.Fatalf("messages = %d, want 0 after delete", len(messages))
	}
}
