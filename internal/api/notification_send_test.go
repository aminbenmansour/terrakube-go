package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/notification"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

func notificationSendAPI(q store.Querier) http.Handler {
	return NewNotification(Deps{Store: q, Deliverer: notification.NewDeliverer(false)})
}

func captureJSONServer(t *testing.T, status int) (*httptest.Server, *string, *http.Header) {
	t.Helper()
	var body string
	header := http.Header{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		body = string(raw)
		header = r.Header.Clone()
		w.WriteHeader(status)
	}))
	t.Cleanup(server.Close)
	return server, &body, &header
}

func TestNotificationTestSendSavedConfiguration(t *testing.T) {
	server, body, header := captureJSONServer(t, http.StatusOK)
	q := newNotificationStub()
	q.org.Name = "acme"
	q.notifications[0].DestinationUrl = server.URL + "/hook"
	q.notifications[0].SigningSecret = strPtr("sign-me")
	before := len(q.notifications)

	rec := httptest.NewRecorder()
	notificationSendAPI(q).ServeHTTP(rec,
		httptest.NewRequest(http.MethodPost, "/configuration/"+notificationConfigID+"/test", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	for _, want := range []string{`"organization":"acme"`, `"workspace":"(test notification)"`,
		`"configurationName":"alerts"`, `"status":"completed"`} {
		if !strings.Contains(*body, want) {
			t.Fatalf("payload missing %q: %s", want, *body)
		}
	}
	if got := (*header).Get("X-Terrakube-Signature"); got != notification.Sign("sign-me", *body) {
		t.Fatalf("signature = %q, want the webhook signature", got)
	}
	if len(q.notifications) != before {
		t.Fatal("send test persisted a configuration")
	}
}

func TestNotificationTestSendUsesTheWorkspaceName(t *testing.T) {
	server, body, _ := captureJSONServer(t, http.StatusOK)
	q := newNotificationStub()
	q.workspaces[0].Name = "networking"
	q.notifications[0].WorkspaceID = strPtr(notificationWorkspaceID)
	q.notifications[0].DestinationUrl = server.URL + "/hook"

	rec := httptest.NewRecorder()
	notificationSendAPI(q).ServeHTTP(rec,
		httptest.NewRequest(http.MethodPost, "/configuration/"+notificationConfigID+"/test", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(*body, `"workspace":"networking"`) {
		t.Fatalf("payload does not name the workspace: %s", *body)
	}
}

func TestNotificationTestSendAnswersBadGatewayWhenDeliveryFails(t *testing.T) {
	server, _, _ := captureJSONServer(t, http.StatusInternalServerError)
	q := newNotificationStub()
	q.notifications[0].DestinationUrl = server.URL + "/hook"

	rec := httptest.NewRecorder()
	notificationSendAPI(q).ServeHTTP(rec,
		httptest.NewRequest(http.MethodPost, "/configuration/"+notificationConfigID+"/test", nil))
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", rec.Code)
	}
}

func TestNotificationTestSendUnknownConfiguration(t *testing.T) {
	q := newNotificationStub()
	for _, id := range []string{notificationMissingID, "not-a-uuid"} {
		rec := httptest.NewRecorder()
		notificationSendAPI(q).ServeHTTP(rec,
			httptest.NewRequest(http.MethodPost, "/configuration/"+id+"/test", nil))
		want := http.StatusNotFound
		if id == "not-a-uuid" {
			want = http.StatusBadRequest
		}
		if rec.Code != want {
			t.Fatalf("id %s status = %d, want %d", id, rec.Code, want)
		}
	}
}

func TestNotificationAdHocTestSend(t *testing.T) {
	server, body, header := captureJSONServer(t, http.StatusOK)
	q := newNotificationStub()
	q.org.Name = "acme"

	// The local destination is why the test API uses the permissive deliverer.
	request := `{"channelType":"WEBHOOK","destinationUrl":"` + server.URL + `/hook","signingSecret":"ad-hoc"}`
	rec := httptest.NewRecorder()
	notificationSendAPI(q).ServeHTTP(rec,
		httptest.NewRequest(http.MethodPost, "/organization/"+notificationOrgID+"/configuration/test", strings.NewReader(request)))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	for _, want := range []string{`"organization":"acme"`, `"workspace":"(test notification)"`,
		`"configurationName":"Test Notification"`} {
		if !strings.Contains(*body, want) {
			t.Fatalf("payload missing %q: %s", want, *body)
		}
	}
	if got := (*header).Get("X-Terrakube-Signature"); got != notification.Sign("ad-hoc", *body) {
		t.Fatalf("signature = %q, want the ad-hoc signature", got)
	}
	if len(q.notifications) != 1 {
		t.Fatal("ad-hoc send persisted a configuration")
	}
}

func TestNotificationAdHocTestSendRejectsBadRequests(t *testing.T) {
	q := newNotificationStub()
	h := notificationSendAPI(q)
	tests := []struct {
		name string
		body string
		want int
	}{
		{"invalid json", `{`, http.StatusBadRequest},
		{"missing channel", `{"destinationUrl":"https://example.com/hook"}`, http.StatusBadRequest},
		{"invalid channel", `{"channelType":"EMAIL","destinationUrl":"https://example.com/hook"}`, http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec,
				httptest.NewRequest(http.MethodPost, "/organization/"+notificationOrgID+"/configuration/test", strings.NewReader(test.body)))
			if rec.Code != test.want {
				t.Fatalf("status = %d, want %d", rec.Code, test.want)
			}
		})
	}
}

type missingOrganizationStub struct{ *stubQuerier }

func (missingOrganizationStub) GetOrganization(context.Context, string) (store.Organization, error) {
	return store.Organization{}, pgx.ErrNoRows
}

func TestNotificationAdHocTestSendUnknownOrganization(t *testing.T) {
	q := missingOrganizationStub{newNotificationStub()}
	rec := httptest.NewRecorder()
	notificationSendAPI(q).ServeHTTP(rec,
		httptest.NewRequest(http.MethodPost, "/organization/"+notificationMissingID+"/configuration/test",
			strings.NewReader(`{"channelType":"WEBHOOK","destinationUrl":"https://example.com/hook"}`)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}
