// Package tfstate implements the executor-facing Terraform state backend at
// /tfstate/v1. It mirrors the upstream TerraformStateController: current state
// pull/push/delete, archived JSON state, and plan artifacts.
package tfstate

import (
	"bytes"
	"crypto/md5"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const defaultMaxStateSizeBytes int64 = 104857600 // 100 MiB, upstream default

// Deps are the collaborators the state backend needs.
type Deps struct {
	Store             store.Querier
	Storage           storage.Backend
	Logger            *slog.Logger
	Layout            Layout
	Auth              *internalauth.Validator
	MaxStateSizeBytes int64
}

// Handler serves the /tfstate/v1 routes.
type Handler struct {
	deps Deps
}

// New builds the state backend router.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.MaxStateSizeBytes <= 0 {
		deps.MaxStateSizeBytes = defaultMaxStateSizeBytes
	}

	h := &Handler{deps: deps}
	r := chi.NewRouter()
	r.Route("/organization/{organizationID}/workspace/{workspaceID}", func(r chi.Router) {
		r.Use(h.authenticate)

		r.Get("/backend/terraform.tfstate", h.getBackendState)
		r.Post("/backend/terraform.tfstate", h.updateBackendState)
		r.Put("/backend/terraform.tfstate", h.updateBackendState)
		r.Delete("/backend/terraform.tfstate", h.deleteBackendState)

		r.Get("/state/terraform.tfstate", h.getCurrentState)
		r.Get("/state/{stateFilename}", h.getStateJSON)
		r.Get("/jobId/{jobID}/step/{stepID}/terraform.tfstate", h.getPlan)
	})

	return r
}

func (h *Handler) getBackendState(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requestClaims(w, r)
	if !ok {
		return
	}
	organizationID, workspaceID := pathIDs(r)
	if !h.authorizeBackend(w, claims, organizationID, workspaceID) {
		return
	}

	key := h.deps.Layout.CurrentState(organizationID, workspaceID)
	obj, err := h.deps.Storage.Stat(r.Context(), key)
	if errors.Is(err, storage.ErrNotFound) || (err == nil && obj.Size == 0) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		h.serverError(w, "stat current state", err)
		return
	}

	rc, err := h.deps.Storage.Get(r.Context(), key)
	if errors.Is(err, storage.ErrNotFound) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		h.serverError(w, "read current state", err)
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, rc); err != nil {
		h.deps.Logger.Warn("state backend response interrupted", "error", err)
	}
}

func (h *Handler) updateBackendState(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requestClaims(w, r)
	if !ok {
		return
	}
	organizationID, workspaceID := pathIDs(r)
	if !h.authorizeBackend(w, claims, organizationID, workspaceID) {
		return
	}

	if r.ContentLength > h.deps.MaxStateSizeBytes {
		h.deps.Logger.Warn("state upload rejected: content length exceeds limit",
			"organization_id", organizationID, "workspace_id", workspaceID,
			"content_length", r.ContentLength, "limit", h.deps.MaxStateSizeBytes)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, h.deps.MaxStateSizeBytes))
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			h.deps.Logger.Warn("state upload rejected: body exceeds limit",
				"organization_id", organizationID, "workspace_id", workspaceID,
				"limit", h.deps.MaxStateSizeBytes)
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		h.serverError(w, "read state upload", err)
		return
	}
	if len(body) == 0 {
		h.deps.Logger.Warn("state upload rejected: empty body",
			"organization_id", organizationID, "workspace_id", workspaceID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if contentMD5 := strings.TrimSpace(r.Header.Get("Content-MD5")); contentMD5 != "" && !matchesMD5(body, contentMD5) {
		h.deps.Logger.Warn("state upload rejected: Content-MD5 mismatch",
			"organization_id", organizationID, "workspace_id", workspaceID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := h.deps.Layout.CurrentState(organizationID, workspaceID)
	if err := h.deps.Storage.Put(r.Context(), key, bytes.NewReader(body), int64(len(body)), "application/json"); err != nil {
		h.serverError(w, "save current state", err)
		return
	}

	h.deps.Logger.Info("saved HTTP backend state",
		"organization_id", organizationID,
		"workspace_id", workspaceID,
		"job_id", claims.JobID,
		"step_id", claims.EffectiveStepID(),
		"bytes", len(body))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteBackendState(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requestClaims(w, r)
	if !ok {
		return
	}
	organizationID, workspaceID := pathIDs(r)
	if !h.authorizeBackend(w, claims, organizationID, workspaceID) {
		return
	}

	key := h.deps.Layout.CurrentState(organizationID, workspaceID)
	if err := h.deps.Storage.Delete(r.Context(), key); err != nil && !errors.Is(err, storage.ErrNotFound) {
		h.serverError(w, "delete current state", err)
		return
	}
	h.deps.Logger.Info("deleted HTTP backend state",
		"organization_id", organizationID,
		"workspace_id", workspaceID,
		"job_id", claims.JobID,
		"step_id", claims.EffectiveStepID())
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getCurrentState(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requestClaims(w, r)
	if !ok {
		return
	}
	organizationID, workspaceID := pathIDs(r)
	if !h.authorizeReadState(w, r, claims, organizationID, workspaceID) {
		return
	}
	h.writeObject(w, r, h.deps.Layout.CurrentState(organizationID, workspaceID), "application/json")
}

func (h *Handler) getStateJSON(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requestClaims(w, r)
	if !ok {
		return
	}
	organizationID, workspaceID := pathIDs(r)

	filename := chi.URLParam(r, "stateFilename")
	if !strings.HasSuffix(filename, ".json") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	filename = strings.TrimSuffix(filename, ".json")
	if filename == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !h.authorizeReadState(w, r, claims, organizationID, workspaceID) {
		return
	}
	h.writeObject(w, r, h.deps.Layout.JSONState(organizationID, workspaceID, filename), "application/json")
}

func (h *Handler) getPlan(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requestClaims(w, r)
	if !ok {
		return
	}
	organizationID, workspaceID := pathIDs(r)
	if !h.authorizeManageState(w, claims, organizationID, workspaceID) {
		return
	}
	key := h.deps.Layout.Plan(organizationID, workspaceID, chi.URLParam(r, "jobID"), chi.URLParam(r, "stepID"))
	h.writeObject(w, r, key, "application/octet-stream")
}

func (h *Handler) writeObject(w http.ResponseWriter, r *http.Request, key, contentType string) {
	rc, err := h.deps.Storage.Get(r.Context(), key)
	if errors.Is(err, storage.ErrNotFound) {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		return
	}
	if err != nil {
		h.serverError(w, "read state object", err)
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, rc); err != nil {
		h.deps.Logger.Warn("state response interrupted", "key", key, "error", err)
	}
}

func (h *Handler) authorizeBackend(w http.ResponseWriter, claims *internalauth.Claims, organizationID, workspaceID string) bool {
	if internalauth.Classify(claims) != internalauth.TypeWorkspaceExecutor {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	if !strings.EqualFold(claims.OrganizationID, organizationID) || !strings.EqualFold(claims.WorkspaceID, workspaceID) {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}

func (h *Handler) authorizeReadState(w http.ResponseWriter, r *http.Request, claims *internalauth.Claims, organizationID, workspaceID string) bool {
	switch internalauth.Classify(claims) {
	case internalauth.TypeExecutorService:
		return true
	case internalauth.TypeWorkspaceExecutor:
	default:
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	if !strings.EqualFold(claims.OrganizationID, organizationID) {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	if strings.EqualFold(claims.WorkspaceID, workspaceID) {
		return true
	}

	target, err := h.deps.Store.GetWorkspace(r.Context(), workspaceID)
	if err != nil || !strings.EqualFold(target.OrganizationID, organizationID) {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	if target.GlobalRemoteState || sharedIDsContain(target.SharedIds, claims.WorkspaceID) {
		return true
	}
	w.WriteHeader(http.StatusForbidden)
	return false
}

func (h *Handler) authorizeManageState(w http.ResponseWriter, claims *internalauth.Claims, organizationID, workspaceID string) bool {
	switch internalauth.Classify(claims) {
	case internalauth.TypeExecutorService:
		return true
	case internalauth.TypeWorkspaceExecutor:
		if strings.EqualFold(claims.OrganizationID, organizationID) && strings.EqualFold(claims.WorkspaceID, workspaceID) {
			return true
		}
	}
	w.WriteHeader(http.StatusForbidden)
	return false
}

func (h *Handler) requestClaims(w http.ResponseWriter, r *http.Request) (*internalauth.Claims, bool) {
	claims, ok := claimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}
	return claims, true
}

func (h *Handler) serverError(w http.ResponseWriter, operation string, err error) {
	h.deps.Logger.Error("state backend request failed", "operation", operation, "error", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func pathIDs(r *http.Request) (string, string) {
	return chi.URLParam(r, "organizationID"), chi.URLParam(r, "workspaceID")
}

func matchesMD5(data []byte, expectedBase64 string) bool {
	sum := md5.Sum(data)
	computed := base64.StdEncoding.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(computed), []byte(expectedBase64)) == 1
}

func sharedIDsContain(shared *string, workspaceID string) bool {
	if shared == nil || *shared == "" {
		return false
	}
	for _, id := range strings.Split(*shared, ",") {
		if strings.EqualFold(strings.TrimSpace(id), workspaceID) {
			return true
		}
	}
	return false
}
