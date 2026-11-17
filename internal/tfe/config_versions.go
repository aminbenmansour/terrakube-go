package tfe

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	maxConfigurationBytes = 100 << 20
	uploadTicketTTL       = time.Hour
)

var (
	errUploadTicketInvalid = errors.New("invalid upload ticket")
	errUploadTicketExpired = errors.New("upload ticket expired")
)

type configurationVersionAttributes struct {
	AutoQueueRuns *bool `json:"auto-queue-runs"`
	Speculative   *bool `json:"speculative"`
}

func (a *API) createConfigurationVersion(w http.ResponseWriter, r *http.Request) {
	ws, ok := a.workspaceByID(w, r)
	if !ok {
		return
	}

	data, err := jsonapi.DecodeRequest(r)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	attrs, err := jsonapi.DecodeAttributes[configurationVersionAttributes](data)
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	autoQueueRuns := false
	if attrs.AutoQueueRuns != nil {
		autoQueueRuns = *attrs.AutoQueueRuns
	}
	speculative := false
	if attrs.Speculative != nil {
		speculative = *attrs.Speculative
	}

	content, err := a.deps.Store.CreateContent(r.Context(), store.CreateContentParams{
		ID:            uuid.NewString(),
		WorkspaceID:   stringPointer(ws.ID),
		AutoQueueRuns: autoQueueRuns,
		Speculative:   &speculative,
		Status:        stringPointer("pending"),
		Source:        stringPointer("tfe-api"),
	})
	if err != nil {
		a.internalError(w, "create configuration version", err)
		return
	}

	resource := configurationVersionResource(content)
	resource.Attributes["upload-url"] = a.configurationUploadURL(content.ID)
	jsonapi.WriteResource(w, http.StatusCreated, resource)
}

func (a *API) getConfigurationVersion(w http.ResponseWriter, r *http.Request) {
	content, err := a.deps.Store.GetContent(r.Context(), chi.URLParam(r, "configurationID"))
	if err != nil {
		if !a.writeStoreError(w, "get configuration version", err) {
			a.internalError(w, "get configuration version", err)
		}
		return
	}
	jsonapi.WriteResource(w, http.StatusOK, configurationVersionResource(content))
}

func (a *API) uploadConfigurationVersion(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "configurationID")
	ticket := r.URL.Query().Get("ticket")
	if ticket == "" {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "An upload ticket is required.")
		return
	}
	if err := a.validateUploadTicket(contentID, ticket); err != nil {
		if errors.Is(err, errUploadTicketExpired) {
			jsonapi.WriteError(w, http.StatusGone, "Gone", "The configuration version upload ticket has expired.")
			return
		}
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "The configuration version upload ticket is invalid.")
		return
	}
	if a.deps.Storage == nil {
		a.internalError(w, "upload configuration version", errors.New("storage is not configured"))
		return
	}

	content, err := a.deps.Store.GetContent(r.Context(), contentID)
	if err != nil {
		if !a.writeStoreError(w, "get configuration version", err) {
			a.internalError(w, "get configuration version", err)
		}
		return
	}
	if content.Status == nil || *content.Status != "pending" {
		jsonapi.WriteError(w, http.StatusConflict, "Conflict", "The configuration version has already been uploaded or is not pending.")
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxConfigurationBytes))
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			jsonapi.WriteError(w, http.StatusRequestEntityTooLarge, "Request Entity Too Large", "The configuration tarball exceeds the upload limit.")
			return
		}
		a.internalError(w, "read configuration upload", err)
		return
	}
	if len(body) == 0 {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Configuration", "The configuration tarball is empty.")
		return
	}

	if err := a.deps.Storage.Put(r.Context(), a.deps.Layout.Content(contentID), bytes.NewReader(body), int64(len(body)), "application/gzip"); err != nil {
		a.internalError(w, "store configuration tarball", err)
		return
	}
	if _, err := a.deps.Store.UpdateContentStatus(r.Context(), store.UpdateContentStatusParams{
		ID:     contentID,
		Status: stringPointer("uploaded"),
	}); err != nil {
		a.internalError(w, "mark configuration uploaded", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) downloadConfigurationVersion(w http.ResponseWriter, r *http.Request) {
	contentID := chi.URLParam(r, "configurationID")
	if a.deps.Storage == nil {
		a.internalError(w, "download configuration version", errors.New("storage is not configured"))
		return
	}
	if _, err := a.deps.Store.GetContent(r.Context(), contentID); err != nil {
		if !a.writeStoreError(w, "get configuration version", err) {
			a.internalError(w, "get configuration version", err)
		}
		return
	}

	rc, err := a.deps.Storage.Get(r.Context(), a.deps.Layout.Content(contentID))
	if errors.Is(err, storage.ErrNotFound) {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The configuration tarball has not been uploaded.")
		return
	}
	if err != nil {
		a.internalError(w, "read configuration tarball", err)
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, rc); err != nil {
		a.deps.Logger.Warn("configuration download interrupted", "error", err, "content_id", contentID)
	}
}

func configurationVersionResource(content store.Content) jsonapi.Resource {
	attributes := map[string]any{
		"error":           nil,
		"error-message":   nil,
		"source":          nil,
		"status":          nil,
		"speculative":     false,
		"auto-queue-runs": content.AutoQueueRuns,
	}
	if content.Source != nil {
		attributes["source"] = *content.Source
	}
	if content.Status != nil {
		attributes["status"] = *content.Status
	}
	if content.Speculative != nil {
		attributes["speculative"] = *content.Speculative
	}
	return jsonapi.Resource{
		Type:       "configuration-versions",
		ID:         content.ID,
		Attributes: attributes,
	}
}

func (a *API) configurationUploadURL(contentID string) string {
	ticket := a.signUploadTicket(contentID, time.Now().Add(uploadTicketTTL))
	return fmt.Sprintf("%s/remote/tfe/v2/configuration-versions/%s?ticket=%s",
		strings.TrimRight(a.deps.PublicURL, "/"), contentID, url.QueryEscape(ticket))
}

func (a *API) signUploadTicket(contentID string, expires time.Time) string {
	expiry := expires.Unix()
	payload := fmt.Sprintf("%s.%d", contentID, expiry)
	mac := hmac.New(sha256.New, []byte(a.deps.UploadTicketSecret))
	_, _ = mac.Write([]byte(payload))
	return payload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (a *API) validateUploadTicket(contentID, ticket string) error {
	parts := strings.Split(ticket, ".")
	if len(parts) != 3 || parts[0] != contentID {
		return errUploadTicketInvalid
	}
	expiry, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return errUploadTicketInvalid
	}
	expected := a.signUploadTicket(contentID, time.Unix(expiry, 0))
	if subtle.ConstantTimeCompare([]byte(ticket), []byte(expected)) != 1 {
		return errUploadTicketInvalid
	}
	if time.Now().Unix() > expiry {
		return errUploadTicketExpired
	}
	return nil
}
