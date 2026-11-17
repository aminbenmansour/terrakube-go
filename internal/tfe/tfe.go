// Package tfe implements the Terraform/OpenTofu TFE v2 API surface at
// /remote/tfe/v2. Upstream Terrakube serves this through RemoteTfeController;
// the JSON:API shapes here mirror that controller and its model classes.
package tfe

import (
	"context"
	"crypto/subtle"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/storage"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	"github.com/aminbenmansour/terrakube-go/internal/tfstate"
)

// JobQueue publishes created runs to the executor.
type JobQueue interface {
	EnqueueJob(ctx context.Context, jobID int64) (string, error)
}

// Deps are the collaborators the TFE API needs.
type Deps struct {
	Store              store.Querier
	Storage            storage.Backend
	Queue              JobQueue
	Layout             tfstate.Layout
	Logger             *slog.Logger
	PublicURL          string
	AuthToken          string
	UploadTicketSecret string
	InternalAuth       *internalauth.Validator
	DexClientID        string
	DexIssuerURL       string
	TerraformLogin     bool
}

// API holds the TFE handlers.
type API struct {
	deps Deps
}

// New builds the /remote/tfe/v2 router.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.UploadTicketSecret == "" {
		deps.UploadTicketSecret = deps.AuthToken
	}
	if deps.UploadTicketSecret == "" {
		deps.UploadTicketSecret = uuid.NewString()
	}

	a := &API{deps: deps}
	r := chi.NewRouter()

	// Upstream's ping handler has no Principal and intentionally returns 404
	// with TFP headers; go-tfe uses the headers to learn the API version.
	r.Get("/ping", a.ping)

	// Configuration uploads authenticate with a short-lived signed ticket, not
	// the bearer token, matching the upstream upload URL. The download endpoint
	// is also unauthenticated upstream because the executor fetches it.
	r.Put("/configuration-versions/{configurationID}", a.uploadConfigurationVersion)
	r.Get("/configuration-versions/{configurationID}/terraformContent.tar.gz", a.downloadConfigurationVersion)

	r.Group(func(r chi.Router) {
		r.Use(a.authenticate)

		r.Get("/organizations/{organizationName}", a.getOrganization)
		r.Get("/organizations/{organizationName}/entitlement-set", a.getEntitlementSet)
		r.Get("/organizations/{organizationName}/capacity", a.getCapacity)

		r.Get("/organizations/{organizationName}/workspaces", a.listWorkspaces)
		r.Post("/organizations/{organizationName}/workspaces", a.createWorkspace)
		r.Get("/organizations/{organizationName}/workspaces/{workspaceName}", a.getWorkspace)
		r.Patch("/workspaces/{workspaceID}", a.updateWorkspace)
		r.Post("/workspaces/{workspaceID}/actions/lock", a.lockWorkspace)
		r.Post("/workspaces/{workspaceID}/actions/unlock", a.unlockWorkspace)
		r.Get("/workspaces/{workspaceID}/tag-bindings", a.listTagBindings)
		r.Patch("/workspaces/{workspaceID}/tag-bindings", a.updateTagBindings)
		r.Post("/workspaces/{workspaceID}/relationships/tags", a.updateWorkspaceTags)

		r.Post("/workspaces/{workspaceID}/state-versions", a.createStateVersion)
		r.Get("/state-versions/{historyID}", a.getStateVersion)
		r.Get("/workspaces/{workspaceID}/current-state-version", a.getCurrentStateVersion)

		r.Post("/workspaces/{workspaceID}/configuration-versions", a.createConfigurationVersion)
		r.Get("/configuration-versions/{configurationID}", a.getConfigurationVersion)

		r.Post("/runs", a.createRun)
		r.Get("/runs/{runID}", a.getRun)
		r.Get("/runs/{runID}/run-events", a.getRunEvents)
		r.Get("/runs/{runID}/policy-checks", a.getRunPolicyChecks)
		r.Post("/runs/{runID}/actions/apply", a.applyRun)
		r.Post("/runs/{runID}/actions/discard", a.discardRun)
		r.Get("/workspaces/{workspaceID}/runs", a.listWorkspaceRuns)
		r.Get("/organizations/{organizationName}/runs/queue", a.getRunQueue)
		r.Get("/plans/{planID}", a.getPlan)
		r.Get("/plans/logs/{planID}/{logToken}", a.getPlanLogs)
		r.Get("/applies/{applyID}", a.getApply)
		r.Get("/applies/logs/{applyID}/{logToken}", a.getApplyLogs)
		r.Get("/policy-checks/{checkID}", a.getPolicyCheck)

		r.Get("/organizations/{organizationName}/projects", a.listProjects)
		r.Post("/organizations/{organizationName}/projects", a.createProject)
	})

	return r
}

func (a *API) authenticate(next http.Handler) http.Handler {
	if a.deps.AuthToken == "" && a.deps.InternalAuth == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := internalauth.ExtractToken(r.Header.Get("Authorization"))
		if raw == "" {
			a.unauthorized(w)
			return
		}

		if a.deps.AuthToken != "" &&
			subtle.ConstantTimeCompare([]byte(raw), []byte(a.deps.AuthToken)) == 1 {
			next.ServeHTTP(w, r)
			return
		}

		if a.deps.InternalAuth != nil {
			claims, err := a.deps.InternalAuth.Parse(raw)
			if err == nil && internalauth.Classify(claims) != internalauth.TypeNotInternal {
				next.ServeHTTP(w, r)
				return
			}
		}

		a.unauthorized(w)
	})
}

func (a *API) unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
}

func (a *API) internalError(w http.ResponseWriter, operation string, err error) {
	a.deps.Logger.Error("tfe request failed", "operation", operation, "error", err)
	jsonapi.WriteError(w, http.StatusInternalServerError, "Internal Server Error", "The request could not be completed.")
}

// writeStoreError maps the common store errors onto TFE responses. It returns
// false when the caller should fall back to a generic 500.
func (a *API) writeStoreError(w http.ResponseWriter, operation string, err error) bool {
	if errors.Is(err, pgx.ErrNoRows) {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return true
	}
	return false
}

func (a *API) organizationByName(w http.ResponseWriter, r *http.Request) (store.Organization, bool) {
	org, err := a.deps.Store.GetOrganizationByName(r.Context(), chi.URLParam(r, "organizationName"))
	if err != nil {
		if !a.writeStoreError(w, "get organization", err) {
			a.internalError(w, "get organization", err)
		}
		return store.Organization{}, false
	}
	return org, true
}

func (a *API) workspaceByID(w http.ResponseWriter, r *http.Request) (store.Workspace, bool) {
	ws, err := a.deps.Store.GetWorkspace(r.Context(), chi.URLParam(r, "workspaceID"))
	if err != nil {
		if !a.writeStoreError(w, "get workspace", err) {
			a.internalError(w, "get workspace", err)
		}
		return store.Workspace{}, false
	}
	return ws, true
}
