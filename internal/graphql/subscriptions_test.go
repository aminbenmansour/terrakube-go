package graphql

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/aminbenmansour/terrakube-go/internal/jobstatus"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// TestSubscriptionOverGraphQLWS drives the real transport with the UI's
// protocol and payload shape, so the graphql-ws wiring is covered and not just
// the resolver.
func TestSubscriptionOverGraphQLWS(t *testing.T) {
	source := &fakeEvents{events: make(chan jobstatus.Event, 1)}
	server := httptest.NewServer(NewSubscriptions(Deps{
		DB:        failingDB{},
		Events:    source,
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken: "secret",
	}))
	defer server.Close()

	conn := dialSubscriptions(t, server.URL, "Bearer secret")
	defer conn.Close()

	writeMessage(t, conn, map[string]any{
		"type":    "connection_init",
		"payload": map[string]any{"Authorization": "Bearer secret"},
	})
	if message := readMessage(t, conn); message["type"] != "connection_ack" {
		t.Fatalf("init response = %#v, want connection_ack", message)
	}

	writeMessage(t, conn, map[string]any{
		"id":   "1",
		"type": "subscribe",
		"payload": map[string]any{
			"query":     `subscription OnJobStatusChanged($workspaceId: ID!) { jobStatusChanged(workspaceId: $workspaceId) { jobId workspaceId status } }`,
			"variables": map[string]any{"workspaceId": "ws-1"},
		},
	})

	// Give the subscription a moment to register before the event is published,
	// then push one event for another workspace and one for this one.
	time.Sleep(100 * time.Millisecond)
	source.events <- jobstatus.Event{JobID: 1, WorkspaceID: "ws-2", Status: "pending"}
	source.events <- jobstatus.Event{JobID: 2, WorkspaceID: "ws-1", Status: "running"}

	message := readMessage(t, conn)
	if message["type"] != "next" {
		t.Fatalf("message = %#v, want next", message)
	}
	payload, _ := json.Marshal(message["payload"])
	var decoded struct {
		Data struct {
			Event struct {
				JobID       int    `json:"jobId"`
				WorkspaceID string `json:"workspaceId"`
				Status      string `json:"status"`
			} `json:"jobStatusChanged"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("decode payload %s: %v", payload, err)
	}
	if decoded.Data.Event.JobID != 2 || decoded.Data.Event.WorkspaceID != "ws-1" || decoded.Data.Event.Status != "running" {
		t.Fatalf("payload = %s", payload)
	}
}

func TestSubscriptionRejectsABadConnectionToken(t *testing.T) {
	server := httptest.NewServer(NewSubscriptions(Deps{
		DB:        failingDB{},
		Events:    &fakeEvents{events: make(chan jobstatus.Event)},
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		AuthToken: "secret",
	}))
	defer server.Close()

	conn := dialSubscriptions(t, server.URL, "Bearer secret")
	defer conn.Close()

	writeMessage(t, conn, map[string]any{
		"type":    "connection_init",
		"payload": map[string]any{"Authorization": "Bearer wrong"},
	})

	// graphql-transport-ws answers a refused init by closing the socket, so a
	// read error is the rejection.
	_ = conn.SetReadDeadline(time.Now().Add(time.Second))
	var message map[string]any
	if err := conn.ReadJSON(&message); err == nil {
		t.Fatalf("connection stayed open after a bad token: %#v", message)
	}
}

func dialSubscriptions(t *testing.T, serverURL, _ string) *websocket.Conn {
	t.Helper()
	header := http.Header{}
	header.Set("Sec-WebSocket-Protocol", "graphql-transport-ws")
	conn, response, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(serverURL, "http"), header)
	if err != nil {
		if response != nil {
			t.Fatalf("dial: %v (HTTP %d)", err, response.StatusCode)
		}
		t.Fatalf("dial: %v", err)
	}
	return conn
}

func writeMessage(t *testing.T, conn *websocket.Conn, message map[string]any) {
	t.Helper()
	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := conn.WriteJSON(message); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func readMessage(t *testing.T, conn *websocket.Conn) map[string]any {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var message map[string]any
	if err := conn.ReadJSON(&message); err != nil {
		t.Fatalf("read: %v", err)
	}
	return message
}

// statusQuerier is just enough of the store to change a job's status.
type statusQuerier struct {
	store.Querier
	job store.Job
}

func (s statusQuerier) UpdateJobStatus(context.Context, store.UpdateJobStatusParams) (store.Job, error) {
	return s.job, nil
}

// TestSubscriptionSeesAStatusWrittenThroughTheStore closes the whole loop the
// UI depends on: the executor or API updates a job row, the wrapper publishes
// it to Redis, and a browser holding a subscription receives it.
func TestSubscriptionSeesAStatusWrittenThroughTheStore(t *testing.T) {
	redisURL := os.Getenv("TERRAKUBE_TEST_REDIS_URL")
	if redisURL == "" {
		t.Skip("TERRAKUBE_TEST_REDIS_URL is not set")
	}
	bus, err := jobstatus.NewBus(redisURL)
	if err != nil {
		t.Fatalf("NewBus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Close() })

	server := httptest.NewServer(NewSubscriptions(Deps{
		DB:     failingDB{},
		Events: bus,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}))
	defer server.Close()

	conn := dialSubscriptions(t, server.URL, "")
	defer conn.Close()
	writeMessage(t, conn, map[string]any{"type": "connection_init"})
	if message := readMessage(t, conn); message["type"] != "connection_ack" {
		t.Fatalf("init response = %#v", message)
	}
	writeMessage(t, conn, map[string]any{
		"id":   "1",
		"type": "subscribe",
		"payload": map[string]any{
			"query":     `subscription($id: ID!) { jobStatusChanged(workspaceId: $id) { jobId status } }`,
			"variables": map[string]any{"id": "ws-live"},
		},
	})
	time.Sleep(100 * time.Millisecond)

	wrapped := jobstatus.Wrap(statusQuerier{job: store.Job{
		ID:             4242,
		WorkspaceID:    "ws-live",
		OrganizationID: "org-live",
		Status:         "completed",
	}}, bus.Publish)
	if _, err := wrapped.UpdateJobStatus(context.Background(), store.UpdateJobStatusParams{ID: 4242, Status: "completed"}); err != nil {
		t.Fatalf("UpdateJobStatus: %v", err)
	}

	message := readMessage(t, conn)
	if message["type"] != "next" {
		t.Fatalf("message = %#v, want next", message)
	}
	payload, _ := json.Marshal(message["payload"])
	if !strings.Contains(string(payload), `"jobId":4242`) || !strings.Contains(string(payload), `"status":"completed"`) {
		t.Fatalf("payload = %s", payload)
	}
}
