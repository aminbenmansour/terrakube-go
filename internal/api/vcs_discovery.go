package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
	vcsprovider "github.com/aminbenmansour/terrakube-go/internal/vcs"
)

// NewVCS builds the root /vcs/v1 router: the workspace picker's repository and
// branch discovery, which upstream mounts outside /api/v1.
func NewVCS(deps Deps) http.Handler {
	a := &API{deps: deps}
	r := chi.NewRouter()
	r.Use(httpapi.CORS(deps.CORSOrigins))
	r.Use(httpapi.LimitBody(jsonapi.MaxBodyBytes))
	r.Use(a.authenticate)
	r.Route("/{vcsID}", func(r chi.Router) {
		r.Use(a.vcsDiscoveryScope)
		r.Get("/groups", a.listVCSGroups)
		r.Get("/repositories", a.listVCSRepositories)
		r.Get("/branches", a.listVCSBranches)
		r.Post("/test", a.testVCSConnection)
	})
	return r
}

type vcsProviderKey struct{}

// vcsDiscoveryScope loads the connection once and requires the permission
// upstream's VcsRepositoryAccessService.hasViewPermission checks: manage VCS or
// manage workspaces in the connection's organization.
func (a *API) vcsDiscoveryScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := param(r, "vcsID")
		if _, err := uuid.Parse(id); err != nil {
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+id)
			return
		}
		provider, err := a.deps.Store.GetVCS(r.Context(), id)
		if err != nil {
			if !a.writeStoreError(w, "get vcs provider", err) {
				a.internalError(w, "get vcs provider", err)
			}
			return
		}
		access, err := a.resolveOrganizationAuthorization(r, provider.OrganizationID)
		if err != nil {
			a.internalError(w, "resolve vcs access", err)
			return
		}
		if !access.Owner && !(access.Permissions.ManageVcs || access.Permissions.ManageWorkspace) {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), vcsProviderKey{}, provider)))
	})
}

func (a *API) listVCSGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := vcsprovider.Groups(r.Context(), a.discoveryConnection(r.Context()))
	if err != nil {
		a.writeVCSDiscoveryError(w, err)
		return
	}
	writeRawJSON(w, http.StatusOK, groups)
}

func (a *API) listVCSRepositories(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	if group == "" {
		jsonapi.WriteError(w, http.StatusBadRequest, "Missing Parameter", "A group is required")
		return
	}
	page, ok := discoveryPage(w, r)
	if !ok {
		return
	}
	repositories, err := vcsprovider.Repositories(r.Context(), a.discoveryConnection(r.Context()),
		group, r.URL.Query().Get("search"), page)
	if err != nil {
		a.writeVCSDiscoveryError(w, err)
		return
	}
	writeRawJSON(w, http.StatusOK, repositories)
}

func (a *API) listVCSBranches(w http.ResponseWriter, r *http.Request) {
	page, ok := discoveryPage(w, r)
	if !ok {
		return
	}
	branches, err := vcsprovider.Branches(r.Context(), a.discoveryConnection(r.Context()),
		r.URL.Query().Get("repository"), page)
	if err != nil {
		a.writeVCSDiscoveryError(w, err)
		return
	}
	writeRawJSON(w, http.StatusOK, branches)
}

func (a *API) testVCSConnection(w http.ResponseWriter, r *http.Request) {
	message, err := vcsprovider.Test(r.Context(), a.discoveryConnection(r.Context()))
	if err != nil {
		a.writeVCSDiscoveryError(w, err)
		return
	}
	writeRawJSON(w, http.StatusOK, map[string]any{"ok": true, "message": message})
}

func (a *API) discoveryConnection(ctx context.Context) vcsprovider.Connection {
	provider := ctx.Value(vcsProviderKey{}).(store.Vc)
	return vcsprovider.Connection{
		Type:           provider.VcsType,
		APIURL:         strValue(provider.ApiUrl),
		AccessToken:    strValue(provider.AccessToken),
		PrivateKey:     strValue(provider.PrivateKey),
		ConnectionType: strValue(provider.ConnectionType),
		AllowLocal:     a.deps.VCSAllowLocal,
	}
}

func discoveryPage(w http.ResponseWriter, r *http.Request) (int, bool) {
	raw := r.URL.Query().Get("page")
	if raw == "" {
		return 1, true
	}
	page, err := strconv.Atoi(raw)
	if err != nil || page < 1 || page > 1000 {
		jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+raw)
		return 0, false
	}
	return page, true
}

func (a *API) writeVCSDiscoveryError(w http.ResponseWriter, err error) {
	var discoveryErr *vcsprovider.Error
	if errors.As(err, &discoveryErr) {
		jsonapi.WriteError(w, discoveryErr.Status, http.StatusText(discoveryErr.Status), discoveryErr.Message)
		return
	}
	a.internalError(w, "vcs discovery", err)
}

func writeRawJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
