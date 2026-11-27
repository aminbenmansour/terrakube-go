package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// notificationDeliveryView is the raw shape upstream's
// NotificationDeliveryView serializes; the delivery history panel reads it
// directly rather than through JSON:API.
type notificationDeliveryView struct {
	ID                string  `json:"id"`
	JobID             int64   `json:"jobId"`
	ConfigurationName string  `json:"configurationName"`
	ChannelType       string  `json:"channelType"`
	Status            string  `json:"status"`
	AttemptCount      int32   `json:"attemptCount"`
	LastAttemptAt     any     `json:"lastAttemptAt"`
	LastError         *string `json:"lastError"`
	CreatedDate       any     `json:"createdDate"`
}

func (a *API) listNotificationDeliveries(w http.ResponseWriter, r *http.Request) {
	limit := int64(10)
	if raw := r.URL.Query().Get("limit"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || parsed < 1 {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+raw)
			return
		}
		limit = min(parsed, 100)
	}
	rows, err := a.deps.Store.ListNotificationDeliveriesByWorkspace(r.Context(), store.ListNotificationDeliveriesByWorkspaceParams{
		WorkspaceID: param(r, "workspaceID"),
		Limit:       int32(limit),
	})
	if err != nil {
		if !a.writeStoreError(w, "list notification deliveries", err) {
			a.internalError(w, "list notification deliveries", err)
		}
		return
	}
	views := make([]notificationDeliveryView, 0, len(rows))
	for _, row := range rows {
		views = append(views, notificationDeliveryView{
			ID:                row.ID,
			JobID:             row.JobID,
			ConfigurationName: row.ConfigurationName,
			ChannelType:       row.ChannelType,
			Status:            row.Status,
			AttemptCount:      row.AttemptCount,
			LastAttemptAt:     timestamp(row.LastAttemptAt),
			LastError:         row.LastError,
			CreatedDate:       timestamp(row.CreatedDate),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(views)
}

// retryNotificationDelivery re-arms a terminal FAILED row so the poller picks
// it up on its next tick. The delivery id is checked against the workspace in
// the path, so a caller cannot retry another workspace's delivery by guessing
// its id.
func (a *API) retryNotificationDelivery(w http.ResponseWriter, r *http.Request) {
	deliveryID := param(r, "deliveryID")
	if _, err := uuid.Parse(deliveryID); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+deliveryID)
		return
	}
	exists, err := a.deps.Store.NotificationOutboxExistsForWorkspace(r.Context(), store.NotificationOutboxExistsForWorkspaceParams{
		ID:          deliveryID,
		WorkspaceID: param(r, "workspaceID"),
	})
	if err != nil {
		a.internalError(w, "check notification delivery", err)
		return
	}
	if !exists {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return
	}
	rows, err := a.deps.Store.RearmNotificationOutbox(r.Context(), deliveryID)
	if err != nil {
		a.internalError(w, "retry notification delivery", err)
		return
	}
	if rows == 0 {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "Only a failed delivery can be retried.")
		return
	}
	w.WriteHeader(http.StatusOK)
}
