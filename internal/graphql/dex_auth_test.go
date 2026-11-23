package graphql

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/jobstatus"
)

type stubDexVerifier struct {
	token string
	err   error
}

func (s stubDexVerifier) Verify(_ context.Context, raw string) error {
	if s.err != nil {
		return s.err
	}
	if raw != s.token {
		return errors.New("invalid token")
	}
	return nil
}

func TestDexTokenIsAcceptedOnGraphQL(t *testing.T) {
	handler := New(Deps{
		DB:     failingDB{},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Dex:    stubDexVerifier{token: "dex-token"},
	})

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/graphql/api/v1", strings.NewReader(`{"query":"{ __typename }"}`)))
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"errors"`) {
		t.Errorf("401 body = %s, want a JSON:API errors document", recorder.Body.String())
	}

	request := httptest.NewRequest(http.MethodPost, "/graphql/api/v1", strings.NewReader(`{"query":"{ __typename }"}`))
	request.Header.Set("Authorization", "Bearer dex-token")
	request.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("with a valid Dex token status = %d, want 200; body = %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"__typename"`) {
		t.Errorf("body = %s, want the query result", recorder.Body.String())
	}
}

func TestSubscriptionAcceptsADexConnectionToken(t *testing.T) {
	server := httptest.NewServer(NewSubscriptions(Deps{
		DB:     failingDB{},
		Events: &fakeEvents{events: make(chan jobstatus.Event)},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Dex:    stubDexVerifier{token: "dex-token"},
	}))
	defer server.Close()

	conn := dialSubscriptions(t, server.URL, "")
	defer conn.Close()

	writeMessage(t, conn, map[string]any{
		"type":    "connection_init",
		"payload": map[string]any{"Authorization": "Bearer dex-token"},
	})
	if message := readMessage(t, conn); message["type"] != "connection_ack" {
		t.Fatalf("init response = %#v, want connection_ack", message)
	}
}
