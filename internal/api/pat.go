package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/pat"
)

// PATService issues personal access tokens. Nil leaves /pat/v1 unavailable.
type PATService interface {
	auth.Verifier
	Create(ctx context.Context, user auth.Identity, days int, description string) (string, error)
	List(ctx context.Context, email string) ([]pat.Token, error)
	Delete(ctx context.Context, id string, user auth.Identity) error
}

// NewPAT builds the upstream /pat/v1 router the UI's Settings > Tokens screen
// calls. It is not JSON:API: upstream returns plain objects here.
func NewPAT(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	a := &API{deps: deps}

	r := chi.NewRouter()
	r.Use(httpapi.LimitBody(jsonapi.MaxBodyBytes))
	r.Use(a.authenticate)
	r.Use(a.requirePATService)
	r.Post("/", a.createPAT)
	r.Get("/", a.listPATs)
	r.Delete("/{tokenID}", a.deletePAT)
	return r
}

func (a *API) requirePATService(next http.Handler) http.Handler {
	if a.deps.PATs != nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		jsonapi.WriteError(w, http.StatusServiceUnavailable, "Service Unavailable", "Personal access tokens are not configured.")
	})
}

type createPATRequest struct {
	Days        int    `json:"days"`
	Description string `json:"description"`
}

func (a *API) createPAT(w http.ResponseWriter, r *http.Request) {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
		return
	}

	request := createPATRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Bad Request", "The request body is not valid JSON.")
		return
	}

	token, err := a.deps.PATs.Create(r.Context(), *identity, request.Days, request.Description)
	switch {
	case errors.Is(err, pat.ErrNonInteractive):
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "A personal access token cannot create another one.")
		return
	case errors.Is(err, pat.ErrInvalidLifetime):
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "days cannot be negative.")
		return
	case err != nil:
		a.internalError(w, "create pat", err)
		return
	}

	writePlainJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (a *API) listPATs(w http.ResponseWriter, r *http.Request) {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
		return
	}

	tokens, err := a.deps.PATs.List(r.Context(), identity.Email)
	if err != nil {
		a.internalError(w, "list pats", err)
		return
	}
	// Upstream answers 202 here, and with [] rather than null when empty.
	writePlainJSON(w, http.StatusAccepted, tokens)
}

func (a *API) deletePAT(w http.ResponseWriter, r *http.Request) {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
		return
	}

	err := a.deps.PATs.Delete(r.Context(), chi.URLParam(r, "tokenID"), *identity)
	switch {
	case errors.Is(err, pat.ErrNotFound):
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested token does not exist.")
		return
	case errors.Is(err, pat.ErrForbidden):
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "This token belongs to another user.")
		return
	case err != nil:
		a.internalError(w, "delete pat", err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func writePlainJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
