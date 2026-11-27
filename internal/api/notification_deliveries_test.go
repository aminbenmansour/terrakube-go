package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const notificationDeliveryID = "55555555-5555-4555-8555-555555555555"

func newDeliveryStub() *stubQuerier {
	q := newNotificationStub()
	q.deliveries = []store.ListNotificationDeliveriesByWorkspaceRow{{
		ID: notificationDeliveryID, JobID: 42, ConfigurationName: "alerts",
		ChannelType: "WEBHOOK", Status: "SENT", AttemptCount: 1,
		CreatedDate: pgtype.Timestamptz{Time: time.Date(2026, 11, 27, 12, 0, 0, 0, time.UTC), Valid: true},
	}}
	return q
}

func TestNotificationDeliveryHistory(t *testing.T) {
	q := newDeliveryStub()
	h := notificationSendAPI(q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspace/"+notificationWorkspaceID+"/deliveries?limit=5", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	var deliveries []struct {
		ID                string  `json:"id"`
		JobID             int64   `json:"jobId"`
		ConfigurationName string  `json:"configurationName"`
		ChannelType       string  `json:"channelType"`
		Status            string  `json:"status"`
		AttemptCount      int32   `json:"attemptCount"`
		LastAttemptAt     *string `json:"lastAttemptAt"`
		LastError         *string `json:"lastError"`
		CreatedDate       string  `json:"createdDate"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &deliveries); err != nil {
		t.Fatalf("decode: %v; body = %s", err, rec.Body.String())
	}
	if len(deliveries) != 1 {
		t.Fatalf("deliveries = %d, want 1", len(deliveries))
	}
	got := deliveries[0]
	if got.ID != notificationDeliveryID || got.JobID != 42 || got.ConfigurationName != "alerts" ||
		got.ChannelType != "WEBHOOK" || got.Status != "SENT" || got.AttemptCount != 1 {
		t.Fatalf("delivery = %+v", got)
	}
	if got.LastAttemptAt != nil || got.LastError != nil {
		t.Fatalf("empty fields = %v/%v, want nulls", got.LastAttemptAt, got.LastError)
	}
	if got.CreatedDate != "2026-11-27T12:00:00Z" {
		t.Fatalf("createdDate = %q", got.CreatedDate)
	}
}

func TestNotificationDeliveryHistoryRejectsABadLimit(t *testing.T) {
	h := notificationSendAPI(newDeliveryStub())
	for _, limit := range []string{"nope", "0", "-2"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet,
			"/workspace/"+notificationWorkspaceID+"/deliveries?limit="+limit, nil))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("limit %q status = %d, want 400", limit, rec.Code)
		}
	}
}

func TestNotificationDeliveryHistoryUnknownWorkspace(t *testing.T) {
	h := notificationSendAPI(newDeliveryStub())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspace/"+notificationMissingID+"/deliveries", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestNotificationDeliveryRetry(t *testing.T) {
	q := newDeliveryStub()
	q.deliveryExists = map[string]bool{notificationDeliveryID: true}
	q.rearmable = map[string]bool{notificationDeliveryID: true}
	h := notificationSendAPI(q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost,
		"/workspace/"+notificationWorkspaceID+"/deliveries/"+notificationDeliveryID+"/retry", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body = %s", rec.Code, rec.Body.String())
	}
	if len(q.rearmed) != 1 || q.rearmed[0] != notificationDeliveryID {
		t.Fatalf("rearmed = %v, want the failed delivery", q.rearmed)
	}
}

func TestNotificationDeliveryRetryRejections(t *testing.T) {
	q := newDeliveryStub()
	q.deliveryExists = map[string]bool{notificationDeliveryID: true}
	q.rearmable = map[string]bool{} // exists but is not FAILED
	h := notificationSendAPI(q)

	tests := []struct {
		name string
		path string
		want int
	}{
		{"not failed", "/workspace/" + notificationWorkspaceID + "/deliveries/" + notificationDeliveryID + "/retry", http.StatusConflict},
		{"unknown delivery", "/workspace/" + notificationWorkspaceID + "/deliveries/" + notificationMissingID + "/retry", http.StatusNotFound},
		{"invalid delivery", "/workspace/" + notificationWorkspaceID + "/deliveries/nope/retry", http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, test.path, nil))
			if rec.Code != test.want {
				t.Fatalf("status = %d, want %d", rec.Code, test.want)
			}
		})
	}
}
