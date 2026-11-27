package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/rbac"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/teamtoken"
)

// TeamTokenService issues and revokes team-scoped access tokens. It is also a
// verifier, so a team token authenticates REST and GraphQL with its group.
type TeamTokenService interface {
	auth.Verifier
	Create(ctx context.Context, user auth.Identity, group string, days, hours, minutes int, description string) (string, error)
	List(ctx context.Context, user auth.Identity) ([]teamtoken.Token, error)
	Delete(ctx context.Context, id string, user auth.Identity) error
}

// AccessTokenStore is the slice of the generated queries the permission
// endpoints need.
type AccessTokenStore interface {
	ListTeamsByOrganizationAndNames(ctx context.Context, arg store.ListTeamsByOrganizationAndNamesParams) ([]store.Team, error)
	ListAccessByWorkspaceAndNames(ctx context.Context, arg store.ListAccessByWorkspaceAndNamesParams) ([]store.Access, error)
	ListProjectAccessByProjectAndNames(ctx context.Context, arg store.ListProjectAccessByProjectAndNamesParams) ([]store.ProjectAccess, error)
	GetWorkspace(ctx context.Context, id string) (store.Workspace, error)
}

// NewAccessToken builds the /access-token/v1 router the UI calls for its
// permission checks. Like upstream it answers plain JSON with 202, not
// JSON:API, and it is mounted outside /api/v1.
func NewAccessToken(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	a := &API{deps: deps}

	r := chi.NewRouter()
	r.Use(httpapi.CORS(deps.CORSOrigins))
	r.Use(a.authenticate)
	// Only the two permission routes read the store; current-teams echoes the
	// token and must keep working when no database is configured.
	r.Get("/teams/current-teams", a.currentTeams)
	r.Group(func(r chi.Router) {
		r.Use(a.requireAccessTokenStore)
		r.Get("/teams/permissions/organization/{organizationID}", a.organizationPermissions)
		r.Get("/teams/permissions/organization/{organizationID}/workspace/{workspaceID}", a.workspacePermissions)
	})
	r.Group(func(r chi.Router) {
		r.Use(a.requireTeamTokens)
		r.Get("/teams", a.listTeamTokens)
		r.Post("/teams", a.createTeamToken)
		r.Delete("/teams/{tokenID}", a.deleteTeamToken)
	})
	return r
}

func (a *API) requireTeamTokens(next http.Handler) http.Handler {
	if a.deps.TeamTokens != nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		jsonapi.WriteError(w, http.StatusServiceUnavailable, "Service Unavailable", "Team tokens are not configured.")
	})
}

type createTeamTokenRequest struct {
	Group       string `json:"group"`
	Description string `json:"description"`
	Days        int    `json:"days"`
	Hours       int    `json:"hours"`
	Minutes     int    `json:"minutes"`
}

func (a *API) createTeamToken(w http.ResponseWriter, r *http.Request) {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
		return
	}
	request := createTeamTokenRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Bad Request", "The request body is not valid JSON.")
		return
	}

	token, err := a.deps.TeamTokens.Create(r.Context(), *identity, request.Group,
		request.Days, request.Hours, request.Minutes, request.Description)
	switch {
	case errors.Is(err, teamtoken.ErrForbidden):
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "You are not a member of this team.")
		return
	case errors.Is(err, teamtoken.ErrNonInteractive):
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "Only an interactive user can create tokens.")
		return
	case errors.Is(err, teamtoken.ErrInvalidLifetime):
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Token lifetime values cannot be negative.")
		return
	case err != nil:
		a.internalError(w, "create team token", err)
		return
	}
	writeAccessTokenJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (a *API) listTeamTokens(w http.ResponseWriter, r *http.Request) {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
		return
	}
	tokens, err := a.deps.TeamTokens.List(r.Context(), *identity)
	if err != nil {
		a.internalError(w, "list team tokens", err)
		return
	}
	// Upstream answers 202 here, and with [] rather than null when empty.
	writeAccessTokenJSON(w, http.StatusAccepted, tokens)
}

func (a *API) deleteTeamToken(w http.ResponseWriter, r *http.Request) {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
		return
	}
	err := a.deps.TeamTokens.Delete(r.Context(), chi.URLParam(r, "tokenID"), *identity)
	switch {
	case errors.Is(err, teamtoken.ErrNotFound):
		jsonapi.WriteError(w, http.StatusBadRequest, "Bad Request", "The requested token does not exist.")
		return
	case errors.Is(err, teamtoken.ErrForbidden):
		jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", "You are not a member of this team.")
		return
	case err != nil:
		a.internalError(w, "delete team token", err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (a *API) requireAccessTokenStore(next http.Handler) http.Handler {
	if a.deps.AccessTokens != nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		jsonapi.WriteError(w, http.StatusServiceUnavailable, "Service Unavailable", "Team permissions are not configured.")
	})
}

// currentTeams echoes the token's groups, which the UI treats as the teams the
// user belongs to. Upstream returns 202 here too.
func (a *API) currentTeams(w http.ResponseWriter, r *http.Request) {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
		return
	}
	groups := identity.Groups
	if groups == nil {
		groups = []string{}
	}
	writeAccessTokenJSON(w, http.StatusAccepted, map[string]any{"groups": groups})
}

func (a *API) organizationPermissions(w http.ResponseWriter, r *http.Request) {
	identity, orgID, ok := a.permissionCaller(w, r)
	if !ok {
		return
	}
	permissions, ok := a.organizationPermissionSet(w, r, identity, orgID)
	if !ok {
		return
	}
	writeAccessTokenJSON(w, http.StatusAccepted, permissions)
}

func (a *API) workspacePermissions(w http.ResponseWriter, r *http.Request) {
	identity, orgID, ok := a.permissionCaller(w, r)
	if !ok {
		return
	}
	// Parsing and re-serialising canonicalises forms Go accepts but Postgres
	// would reject, such as urn:uuid:.
	workspaceUUID, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		jsonapi.WriteError(w, http.StatusBadRequest, "Bad Request", "The workspace id is not a UUID.")
		return
	}
	workspaceID := workspaceUUID.String()
	permissions, ok := a.organizationPermissionSet(w, r, identity, orgID)
	if !ok {
		return
	}

	// Upstream looks the workspace up by id alone and treats a miss as "no
	// extra permissions", so a missing workspace is not an error here.
	workspace, err := a.deps.AccessTokens.GetWorkspace(r.Context(), workspaceID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeAccessTokenJSON(w, http.StatusAccepted, permissions)
		return
	case err != nil:
		a.internalError(w, "load workspace permissions", err)
		return
	}

	groups := identity.Groups
	accessRows, err := a.deps.AccessTokens.ListAccessByWorkspaceAndNames(r.Context(),
		store.ListAccessByWorkspaceAndNamesParams{WorkspaceID: workspace.ID, Names: groups})
	if err != nil {
		a.internalError(w, "load workspace access", err)
		return
	}
	for _, row := range accessRows {
		permissions.ApplyAccess(rbac.FlagsFromAccess(row))
	}

	if workspace.ProjectID != nil && *workspace.ProjectID != "" {
		projectRows, err := a.deps.AccessTokens.ListProjectAccessByProjectAndNames(r.Context(),
			store.ListProjectAccessByProjectAndNamesParams{ProjectID: *workspace.ProjectID, Names: groups})
		if err != nil {
			a.internalError(w, "load project access", err)
			return
		}
		for _, row := range projectRows {
			permissions.ApplyProjectAccess(rbac.FlagsFromProjectAccess(row))
		}
	}
	writeAccessTokenJSON(w, http.StatusAccepted, permissions)
}

// permissionCaller checks the identity and the organization id shared by both
// permission endpoints.
func (a *API) permissionCaller(w http.ResponseWriter, r *http.Request) (*auth.Identity, string, bool) {
	identity, ok := auth.IdentityFromContext(r.Context())
	if !ok {
		jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
		return nil, "", false
	}
	orgUUID, err := uuid.Parse(chi.URLParam(r, "organizationID"))
	if err != nil {
		// Upstream answers 500 here; a 400 says the same without a stack trace.
		jsonapi.WriteError(w, http.StatusBadRequest, "Bad Request", "The organization id is not a UUID.")
		return nil, "", false
	}
	return identity, orgUUID.String(), true
}

// organizationPermissionSet is the team-derived half of both endpoints.
func (a *API) organizationPermissionSet(w http.ResponseWriter, r *http.Request, identity *auth.Identity, orgID string) (rbac.Permissions, bool) {
	var permissions rbac.Permissions
	if len(identity.Groups) > 0 {
		teams, err := a.deps.AccessTokens.ListTeamsByOrganizationAndNames(r.Context(),
			store.ListTeamsByOrganizationAndNamesParams{OrganizationID: orgID, Names: identity.Groups})
		if err != nil {
			a.internalError(w, "load organization permissions", err)
			return permissions, false
		}
		for _, team := range teams {
			permissions.ApplyTeam(rbac.FlagsFromTeam(team))
		}
	}
	if a.isInstanceOwner(identity) {
		permissions.GrantOwner()
	}
	return permissions, true
}

// isInstanceOwner mirrors upstream's groups.contains(instanceOwner), plus the
// development token, which has no Dex groups but owns a local stack.
func (a *API) isInstanceOwner(identity *auth.Identity) bool {
	if identity.Issuer == auth.DevIssuer {
		return true
	}
	return slices.Contains(identity.Groups, a.adminGroup())
}

func (a *API) adminGroup() string {
	if a.deps.AdminGroup == "" {
		return "TERRAKUBE_ADMIN"
	}
	return a.deps.AdminGroup
}

func writeAccessTokenJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	// The same URL answers differently per token; keep shared caches out.
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
