package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/rbac"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// authorization is the caller's standing for one request. Owner bypasses every
// check; Member means a team or grant row names one of the caller's groups.
type authorization struct {
	Owner       bool
	Member      bool
	Permissions rbac.Permissions
}

type authorizationKey struct{}

func withAuthorization(ctx context.Context, access authorization) context.Context {
	return context.WithValue(ctx, authorizationKey{}, access)
}

func authorizationFrom(ctx context.Context) (authorization, bool) {
	access, ok := ctx.Value(authorizationKey{}).(authorization)
	return access, ok
}

// callerIsOwner reports the instance owner, and treats a server without
// authentication configured as one, which keeps the open development mode.
func (a *API) callerIsOwner(r *http.Request) bool {
	identity, ok := auth.IdentityFromContext(r.Context())
	return !ok || a.isInstanceOwner(identity)
}

// organizationAccess resolves the caller's team permissions once for every
// organization-scoped route and answers a non-member the way Elide does.
func (a *API) organizationAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.callerIsOwner(r) {
			next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), authorization{Owner: true, Member: true})))
			return
		}
		identity, _ := auth.IdentityFromContext(r.Context())
		access, err := rbac.ResolveOrganization(r.Context(), a.deps.Store, param(r, "organizationID"), identity.Groups)
		if err != nil {
			a.internalError(w, "resolve organization access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), authorization{
			Member:      true,
			Permissions: access.Permissions,
		})))
	})
}

// requireWrite checks a write permission for the route it wraps. Reads pass:
// the organization middleware already required membership for them.
func (a *API) requireWrite(check func(rbac.Permissions) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isWriteMethod(r.Method) {
				next.ServeHTTP(w, r)
				return
			}
			if a.callerIsOwner(r) {
				next.ServeHTTP(w, r)
				return
			}
			access, ok := authorizationFrom(r.Context())
			if !ok || !access.Member || !check(access.Permissions) {
				writePermissionDenied(w, r.Method)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// requireOwner guards the resources upstream reserves for the instance owner:
// organizations, teams, agents, global variables and the action catalog.
func (a *API) requireOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.callerIsOwner(r) {
			next.ServeHTTP(w, r)
			return
		}
		writePermissionDenied(w, r.Method)
	})
}

func isWriteMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	}
	return true
}

// writePermissionDenied matches the detail Elide sends. The status and title
// are added for JSON:API clients; the UI reads the detail.
func writePermissionDenied(w http.ResponseWriter, method string) {
	jsonapi.WriteError(w, http.StatusForbidden, "Forbidden", permissionDeniedDetail(method))
}

func permissionDeniedDetail(method string) string {
	switch method {
	case http.MethodGet, http.MethodHead:
		return "ReadPermission Denied"
	case http.MethodPost:
		return "CreatePermission Denied"
	case http.MethodDelete:
		return "DeletePermission Denied"
	default:
		return "UpdatePermission Denied"
	}
}

// resolveWorkspaceAuthorization adds a workspace's grant rows to the access
// already resolved for its organization, or resolves the organization first
// on the root-level workspace and variable routes.
func (a *API) resolveWorkspaceAuthorization(r *http.Request, ws store.Workspace) (authorization, error) {
	if a.callerIsOwner(r) {
		return authorization{Owner: true, Member: true}, nil
	}
	identity, _ := auth.IdentityFromContext(r.Context())

	base, ok := authorizationFrom(r.Context())
	if !ok {
		org, err := rbac.ResolveOrganization(r.Context(), a.deps.Store, ws.OrganizationID, identity.Groups)
		if err != nil {
			return authorization{}, err
		}
		if !org.Member {
			return authorization{}, nil
		}
		base = authorization{Member: true, Permissions: org.Permissions}
	}

	access, err := rbac.ApplyWorkspaceGrants(r.Context(), a.deps.Store,
		rbac.Access{Member: base.Member, Permissions: base.Permissions},
		ws.OrganizationID, ws, identity.Groups)
	if err != nil {
		return authorization{}, err
	}
	return authorization{Member: access.Member, Permissions: access.Permissions}, nil
}

// resolveProjectAuthorization is the project counterpart.
func (a *API) resolveProjectAuthorization(r *http.Request, project store.Project) (authorization, error) {
	if a.callerIsOwner(r) {
		return authorization{Owner: true, Member: true}, nil
	}
	identity, _ := auth.IdentityFromContext(r.Context())

	base, ok := authorizationFrom(r.Context())
	if !ok {
		org, err := rbac.ResolveOrganization(r.Context(), a.deps.Store, project.OrganizationID, identity.Groups)
		if err != nil {
			return authorization{}, err
		}
		if !org.Member {
			return authorization{}, nil
		}
		base = authorization{Member: true, Permissions: org.Permissions}
	}

	access, err := rbac.ApplyProjectGrants(r.Context(), a.deps.Store,
		rbac.Access{Member: base.Member, Permissions: base.Permissions},
		project.OrganizationID, project, identity.Groups)
	if err != nil {
		return authorization{}, err
	}
	return authorization{Member: access.Member, Permissions: access.Permissions}, nil
}

// workspaceScope guards the root-level workspace routes, which do not carry
// the organization in the path.
func (a *API) workspaceScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := a.deps.Store.GetWorkspace(r.Context(), param(r, "workspaceID"))
		if err != nil {
			if !a.writeStoreError(w, "get workspace", err) {
				a.internalError(w, "get workspace", err)
			}
			return
		}
		access, err := a.resolveWorkspaceAuthorization(r, ws)
		if err != nil {
			a.internalError(w, "resolve workspace access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), access)))
	})
}

func (a *API) projectScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project, err := a.deps.Store.GetProject(r.Context(), param(r, "projectID"))
		if err != nil {
			if !a.writeStoreError(w, "get project", err) {
				a.internalError(w, "get project", err)
			}
			return
		}
		access, err := a.resolveProjectAuthorization(r, project)
		if err != nil {
			a.internalError(w, "resolve project access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), access)))
	})
}

func (a *API) variableScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		variable, err := a.deps.Store.GetVariable(r.Context(), param(r, "variableID"))
		if err != nil {
			if !a.writeStoreError(w, "get variable", err) {
				a.internalError(w, "get variable", err)
			}
			return
		}
		ws, err := a.deps.Store.GetWorkspace(r.Context(), variable.WorkspaceID)
		if err != nil {
			if !a.writeStoreError(w, "get workspace", err) {
				a.internalError(w, "get workspace", err)
			}
			return
		}
		access, err := a.resolveWorkspaceAuthorization(r, ws)
		if err != nil {
			a.internalError(w, "resolve workspace access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), access)))
	})
}

func (a *API) notificationScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config, err := a.deps.Store.GetNotificationConfiguration(r.Context(), param(r, "configurationID"))
		if err != nil {
			if !a.writeStoreError(w, "get notification configuration", err) {
				a.internalError(w, "get notification configuration", err)
			}
			return
		}

		var access authorization
		if config.WorkspaceID != nil && *config.WorkspaceID != "" {
			ws, err := a.deps.Store.GetWorkspace(r.Context(), *config.WorkspaceID)
			if err != nil {
				if !a.writeStoreError(w, "get workspace", err) {
					a.internalError(w, "get workspace", err)
				}
				return
			}
			access, err = a.resolveWorkspaceAuthorization(r, ws)
		} else {
			access, err = a.resolveOrganizationAuthorization(r, config.OrganizationID)
		}
		if err != nil {
			a.internalError(w, "resolve notification access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), access)))
	})
}

func (a *API) referenceScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reference, err := a.deps.Store.GetReference(r.Context(), param(r, "referenceID"))
		if err != nil {
			if !a.writeStoreError(w, "get reference", err) {
				a.internalError(w, "get reference", err)
			}
			return
		}
		collection, err := a.deps.Store.GetCollection(r.Context(), reference.CollectionID)
		if err != nil {
			if !a.writeStoreError(w, "get collection", err) {
				a.internalError(w, "get collection", err)
			}
			return
		}
		access, err := a.resolveOrganizationAuthorization(r, collection.OrganizationID)
		if err != nil {
			a.internalError(w, "resolve collection access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), access)))
	})
}

// resolveOrganizationAuthorization is the root-route counterpart of the
// organization middleware.
func (a *API) resolveOrganizationAuthorization(r *http.Request, organizationID string) (authorization, error) {
	if a.callerIsOwner(r) {
		return authorization{Owner: true, Member: true}, nil
	}
	identity, _ := auth.IdentityFromContext(r.Context())
	access, err := rbac.ResolveOrganization(r.Context(), a.deps.Store, organizationID, identity.Groups)
	if err != nil {
		return authorization{}, err
	}
	return authorization{Member: access.Member, Permissions: access.Permissions}, nil
}

// projectInOrganization guards the project routes nested under an
// organization and resolves the project's grant rows.
func (a *API) projectInOrganization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project, err := a.deps.Store.GetProject(r.Context(), param(r, "projectID"))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+param(r, "projectID")+"] for project")
				return
			}
			if !a.writeStoreError(w, "get project", err) {
				a.internalError(w, "get project", err)
			}
			return
		}
		if project.OrganizationID != param(r, "organizationID") {
			jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown identifier ["+project.ID+"] for project")
			return
		}
		access, err := a.resolveProjectAuthorization(r, project)
		if err != nil {
			a.internalError(w, "resolve project access", err)
			return
		}
		if !access.Member {
			writePermissionDenied(w, r.Method)
			return
		}
		next.ServeHTTP(w, r.WithContext(withAuthorization(r.Context(), access)))
	})
}

// The permission set equivalents of rbac's flag checks, for requireWrite.
func canManageWorkspace(p rbac.Permissions) bool  { return p.ManageWorkspace }
func canManageJob(p rbac.Permissions) bool        { return p.ManageJob }
func canManageVcs(p rbac.Permissions) bool        { return p.ManageVcs }
func canManageTemplate(p rbac.Permissions) bool   { return p.ManageTemplate }
func canManageCollection(p rbac.Permissions) bool { return p.ManageCollection }

// authorizeWorkspaceWrite guards the atomic operations document, which carries
// its organization in the href rather than in the route.
func (a *API) authorizeWorkspaceWrite(w http.ResponseWriter, r *http.Request, organizationID string, ws *store.Workspace) bool {
	if a.callerIsOwner(r) {
		return true
	}
	var (
		access authorization
		err    error
	)
	if ws != nil {
		access, err = a.resolveWorkspaceAuthorization(r, *ws)
	} else {
		access, err = a.resolveOrganizationAuthorization(r, organizationID)
	}
	if err != nil {
		a.internalError(w, "resolve workspace access", err)
		return false
	}
	if !access.Member || !access.Permissions.ManageWorkspace {
		writePermissionDenied(w, r.Method)
		return false
	}
	return true
}
