package tfoutput

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aminbenmansour/terrakube-go/internal/queue"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// mutableQuerier lets the test finish a job while the stream is open, which is
// how a real run ends: the last log line arrives, then the job goes terminal.
type mutableQuerier struct {
	*stubQuerier
	mu     sync.Mutex
	status string
}

func (m *mutableQuerier) GetJob(ctx context.Context, id int64) (store.Job, error) {
	job, err := m.stubQuerier.GetJob(ctx, id)
	m.mu.Lock()
	job.Status = m.status
	m.mu.Unlock()
	return job, err
}

func (m *mutableQuerier) finish(status string) {
	m.mu.Lock()
	m.status = status
	m.mu.Unlock()
}

// TestStepLogStreamDeliversLiveOutput is the "live logs appear during a run"
// contract, without a browser: a record written to the job's stream reaches an
// open SSE client, and the stream ends once the job is terminal.
func TestStepLogStreamDeliversLiveOutput(t *testing.T) {
	redisURL := os.Getenv("TERRAKUBE_TEST_REDIS_URL")
	if redisURL == "" {
		t.Skip("TERRAKUBE_TEST_REDIS_URL is not set")
	}
	redisQueue, err := queue.New(redisURL)
	if err != nil {
		t.Fatalf("queue.New: %v", err)
	}
	t.Cleanup(func() { _ = redisQueue.Close() })

	const jobID = 990001 // not a real job; the store is a stub
	ctx := context.Background()
	stream := streamKey(jobID)
	t.Cleanup(func() { _ = redisQueue.DeleteStream(context.Background(), stream) })

	status := "running"
	q := &mutableQuerier{
		stubQuerier: &stubQuerier{
			job:  store.Job{ID: jobID, OrganizationID: "org-1", Status: status},
			step: store.Step{ID: "step-1", JobID: jobID, Status: &status},
		},
		status: status,
	}
	server := httptest.NewServer(New(Deps{Store: q, Queue: redisQueue}))
	defer server.Close()

	response, err := http.Get(server.URL + "/organization/org-1/job/990001/step/step-1/stream")
	if err != nil {
		t.Fatalf("open stream: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", response.StatusCode)
	}
	if got := response.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("Content-Type = %q", got)
	}

	// A record for another step must not reach this client, so it is published
	// first and the assertion below expects the matching line.
	if _, err := redisQueue.Publish(ctx, stream, map[string]any{"stepId": "step-other", "output": "wrong step\n"}); err != nil {
		t.Fatalf("publish other step: %v", err)
	}
	if _, err := redisQueue.Publish(ctx, stream, map[string]any{"stepId": "step-1", "output": "initializing providers\n"}); err != nil {
		t.Fatalf("publish: %v", err)
	}

	frames := make(chan string, 16)
	go func() {
		defer close(frames)
		scanner := bufio.NewScanner(response.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		var frame strings.Builder
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				if frame.Len() > 0 {
					frames <- frame.String()
					frame.Reset()
				}
				continue
			}
			frame.WriteString(line)
			frame.WriteString("\n")
		}
	}()

	first := nextFrame(t, frames)
	if !strings.Contains(first, "data: initializing providers") {
		t.Fatalf("first frame = %q, want the matching step's output", first)
	}

	// Finishing the job ends the stream, which is what tells the UI to stop
	// waiting for output.
	q.finish("completed")
	for {
		frame := nextFrame(t, frames)
		if strings.Contains(frame, "event: end") {
			if !strings.Contains(frame, `"reason":"complete"`) {
				t.Fatalf("end frame = %q", frame)
			}
			return
		}
	}
}

func nextFrame(t *testing.T, frames <-chan string) string {
	t.Helper()
	select {
	case frame, ok := <-frames:
		if !ok {
			t.Fatal("stream closed before the expected frame")
		}
		return frame
	case <-time.After(15 * time.Second):
		t.Fatal("timed out waiting for a frame")
	}
	return ""
}
