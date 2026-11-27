// Package api implements the Terrakube management API at /api/v1.
//
// Upstream Terrakube serves this surface through Elide, so requests and
// responses are JSON:API 1.0 documents with camelCase attribute names.
package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/httpapi"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/notification"
	"github.com/aminbenmansour/terrakube-go/internal/rsql"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Deps are the collaborators the management API needs.
type Deps struct {
	Store store.Querier
	// Lists runs the filtered list queries. New falls back to Store when it
	// implements store.Lister, and the list handlers need one of the two.
	Lists     store.Lister
	Logger    *slog.Logger
	PublicURL string
	AuthToken string
	// Verifiers accept user bearer tokens (Dex, PATs). Empty leaves only the
	// development token, or no authentication at all when that is empty too.
	Verifiers []auth.Verifier
	// PATs backs /pat/v1. Nil answers 503 there.
	PATs PATService
	// AccessTokens backs /access-token/v1. Nil answers 503 there.
	AccessTokens AccessTokenStore
	// TeamTokens backs the team token routes under /access-token/v1. Nil
	// answers 503 there.
	TeamTokens TeamTokenService
	// AdminGroup is the instance owner group. Empty means TERRAKUBE_ADMIN.
	AdminGroup string
	// CORSOrigins lists the origins allowed to call the API from a browser.
	// Empty allows any origin, which is the historical behaviour; a UI served
	// by this API shares its origin and needs no entry.
	CORSOrigins []string
	// VCSAllowLocal permits VCS API URLs on loopback/link-local addresses.
	// Off by default. See config.Config.VCSAllowLocal.
	VCSAllowLocal bool
	// Deliverer posts notification test sends. Nil uses the SSRF-blocking
	// deliverer the HTTP senders default to.
	Deliverer *notification.Deliverer
}

// API holds the HTTP handlers.
type API struct {
	deps Deps
}

// MountPrefix is where cmd/api mounts this router. The UUID check needs it
// because chi leaves the mount prefix in URL.Path, and a r.Use middleware runs
// before routing has resolved anything else.
const MountPrefix = "/api/v1"

// New builds the /api/v1 router.
func New(deps Deps) http.Handler {
	if deps.Lists == nil {
		deps.Lists, _ = deps.Store.(store.Lister)
	}
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.AuthToken == "" && len(deps.Verifiers) == 0 {
		deps.Logger.Warn("management API authentication is disabled: set TERRAKUBE_AUTH_TOKEN or TERRAKUBE_DEX_ENABLED to require a bearer token")
	}

	a := &API{deps: deps}

	r := chi.NewRouter()
	r.Use(httpapi.CORS(deps.CORSOrigins))
	r.Use(httpapi.LimitBody(jsonapi.MaxBodyBytes))
	r.Use(a.authenticate)
	r.Use(requireUUIDParams)

	// JSON:API atomic operations; the UI creates and updates workspaces here.
	r.Post("/operations", a.atomicOperations)

	// The action catalog fills the workspace resource drawer and the Settings >
	// Actions editor.
	r.Get("/action", a.listActions)
	r.With(a.requireOwner).Post("/action", a.createAction)
	r.Route("/action/{actionID}", func(r chi.Router) {
		r.Get("/", a.getAction)
		r.With(a.requireOwner).Patch("/", a.updateAction)
		r.With(a.requireOwner).Delete("/", a.deleteAction)
	})

	// Organizations
	r.Get("/organization", a.listOrganizations)
	r.With(a.requireOwner).Post("/organization", a.createOrganization)
	r.Route("/organization/{organizationID}", func(r chi.Router) {
		r.Use(a.organizationAccess)
		r.Get("/", a.getOrganization)
		r.With(a.requireOwner).Patch("/", a.updateOrganization)
		r.With(a.requireOwner).Delete("/", a.deleteOrganization)

		// Projects and workspaces scoped to an organization.
		r.Get("/project", a.listProjects)
		r.With(a.requireWrite(canManageWorkspace)).Post("/project", a.createProject)
		// Project access grants: Settings > Projects > Access.
		r.Route("/project/{projectID}", func(r chi.Router) {
			r.Use(a.projectInOrganization)
			r.Get("/projectAccess", a.listProjectAccess)
			r.With(a.requireWrite(canManageWorkspace)).Post("/projectAccess", a.createProjectAccess)
			r.Route("/projectAccess/{projectAccessID}", func(r chi.Router) {
				r.Get("/", a.getProjectAccess)
				r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateProjectAccess)
				r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteProjectAccess)
			})
		})
		r.Get("/workspace", a.listWorkspaces)
		r.With(a.requireWrite(canManageWorkspace)).Post("/workspace", a.createWorkspace)

		// Runs.
		r.Get("/job", a.listJobs)
		r.Route("/job/{jobID}", func(r chi.Router) {
			r.Get("/", a.getJob)
			r.With(a.requireWrite(canManageJob)).Patch("/", a.updateJob)
		})

		// Tags: Settings > Tags and the workspace tag editor.
		r.Get("/tag", a.listTags)
		r.With(a.requireWrite(canManageWorkspace)).Post("/tag", a.createTag)
		r.Route("/tag/{tagID}", func(r chi.Router) {
			r.Get("/", a.getTag)
			r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateTag)
			r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteTag)
		})

		// VCS providers: Settings > VCS and the workspace source picker.
		r.Get("/vcs", a.listVCS)
		r.With(a.requireWrite(canManageVcs)).Post("/vcs", a.createVCS)
		r.Route("/vcs/{vcsID}", func(r chi.Router) {
			r.Get("/", a.getVCS)
			r.With(a.requireWrite(canManageVcs)).Patch("/", a.updateVCS)
			r.With(a.requireWrite(canManageVcs)).Delete("/", a.deleteVCS)
		})

		// Notifications: Settings > Notifications, org-wide and per workspace.
		r.Get("/notificationConfiguration", a.listNotificationConfigurations)
		r.With(a.requireWrite(canManageWorkspace)).Post("/notificationConfiguration", a.createNotificationConfiguration)
		r.Get("/notificationConfiguration/{configurationID}", a.getNotificationConfiguration)

		// Teams: Settings > Teams and the access pickers.
		r.Get("/team", a.listTeams)
		r.With(a.requireOwner).Post("/team", a.createTeam)
		r.Route("/team/{teamID}", func(r chi.Router) {
			r.Get("/", a.getTeam)
			r.With(a.requireOwner).Patch("/", a.updateTeam)
			r.With(a.requireOwner).Delete("/", a.deleteTeam)
		})

		// Job-flow templates: Settings > Templates and the workspace default.
		r.Get("/template", a.listTemplates)
		r.With(a.requireWrite(canManageTemplate)).Post("/template", a.createTemplate)
		r.Route("/template/{templateID}", func(r chi.Router) {
			r.Get("/", a.getTemplate)
			r.With(a.requireWrite(canManageTemplate)).Patch("/", a.updateTemplate)
			r.With(a.requireWrite(canManageTemplate)).Delete("/", a.deleteTemplate)
		})

		// Private registry providers: Settings > Registry and the provider
		// picker. Versions and platform implementations hang off the provider.
		r.Get("/provider", a.listProviders)
		r.With(a.requireWrite(canManageProvider)).Post("/provider", a.createProvider)
		r.Route("/provider/{providerID}", func(r chi.Router) {
			r.Get("/", a.getProvider)
			r.With(a.requireWrite(canManageProvider)).Patch("/", a.updateProvider)
			r.With(a.requireWrite(canManageProvider)).Delete("/", a.deleteProvider)
			r.Get("/version", a.listProviderVersions)
			r.With(a.requireWrite(canManageProvider)).Post("/version", a.createProviderVersion)
			r.Route("/version/{versionID}", func(r chi.Router) {
				r.Get("/", a.getProviderVersion)
				r.With(a.requireWrite(canManageProvider)).Patch("/", a.updateProviderVersion)
				r.With(a.requireWrite(canManageProvider)).Delete("/", a.deleteProviderVersion)
				r.Get("/implementation", a.listImplementations)
				r.With(a.requireWrite(canManageProvider)).Post("/implementation", a.createImplementation)
				r.Route("/implementation/{implementationID}", func(r chi.Router) {
					r.Get("/", a.getImplementation)
					r.With(a.requireWrite(canManageProvider)).Patch("/", a.updateImplementation)
					r.With(a.requireWrite(canManageProvider)).Delete("/", a.deleteImplementation)
				})
			})
		})

		// Private registry modules: Settings > Registry and the module picker.
		// Versions are created by the refresh poller and administered here.
		r.Get("/module", a.listModules)
		r.With(a.requireWrite(canManageModule)).Post("/module", a.createModule)
		r.Route("/module/{moduleID}", func(r chi.Router) {
			r.Get("/", a.getModule)
			r.With(a.requireWrite(canManageModule)).Patch("/", a.updateModule)
			r.With(a.requireWrite(canManageModule)).Delete("/", a.deleteModule)
			r.Get("/version", a.listModuleVersions)
			r.Route("/version/{versionID}", func(r chi.Router) {
				r.Get("/", a.getModuleVersion)
				r.With(a.requireWrite(canManageModule)).Patch("/", a.updateModuleVersion)
			})
		})

		// SSH keys and self-hosted agent pools for workspace connections.
		r.Get("/ssh", a.listSSHKeys)
		r.With(a.requireWrite(canManageVcs)).Post("/ssh", a.createSSHKey)
		r.Route("/ssh/{sshID}", func(r chi.Router) {
			r.Get("/", a.getSSHKey)
			r.With(a.requireWrite(canManageVcs)).Patch("/", a.updateSSHKey)
			r.With(a.requireWrite(canManageVcs)).Delete("/", a.deleteSSHKey)
		})
		r.Get("/agent", a.listAgents)
		r.With(a.requireOwner).Post("/agent", a.createAgent)
		r.With(a.requireOwner).Delete("/agent/{agentID}", a.deleteAgent)

		// Organization-level variables inherited by every workspace.
		r.Get("/globalvar", a.listGlobalVariables)
		r.With(a.requireOwner).Post("/globalvar", a.createGlobalVariable)
		r.Route("/globalvar/{globalVarID}", func(r chi.Router) {
			r.Get("/", a.getGlobalVariable)
			r.With(a.requireOwner).Patch("/", a.updateGlobalVariable)
			r.With(a.requireOwner).Delete("/", a.deleteGlobalVariable)
		})
	})

	// Projects
	r.Route("/project/{projectID}", func(r chi.Router) {
		r.Use(a.projectScope)
		r.Get("/", a.getProject)
		r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateProject)
		r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteProject)
		r.Get("/workspace", a.listProjectWorkspaces)
	})

	// Variable collections: Settings > Collections, and the read-only
	// "Collection Variables" table on the workspace page.
	r.Route("/organization/{organizationID}/collection", func(r chi.Router) {
		r.Use(a.organizationAccess)
		r.Get("/", a.listCollections)
		r.With(a.requireWrite(canManageCollection)).Post("/", a.createCollection)
		r.Route("/{collectionID}", func(r chi.Router) {
			r.Get("/", a.getCollection)
			r.With(a.requireWrite(canManageCollection)).Patch("/", a.updateCollection)
			r.With(a.requireWrite(canManageCollection)).Delete("/", a.deleteCollection)
			r.Get("/item", a.listCollectionItems)
			r.With(a.requireWrite(canManageCollection)).Post("/item", a.createCollectionItem)
			r.Route("/item/{itemID}", func(r chi.Router) {
				r.With(a.requireWrite(canManageCollection)).Patch("/", a.updateCollectionItem)
				r.With(a.requireWrite(canManageCollection)).Delete("/", a.deleteCollectionItem)
			})
			r.Get("/reference", a.listCollectionReferences)
			r.With(a.requireWrite(canManageCollection)).Post("/reference", a.createCollectionReference)
			r.With(a.requireWrite(canManageCollection)).Delete("/reference/{referenceID}", a.deleteCollectionReference)
		})
	})
	r.With(a.referenceScope).Get("/reference/{referenceID}/collection", a.getReferenceCollection)

	// Notification configurations are addressed by id at the root because the
	// UI's editor and trigger sub-resource calls are not organization-nested.
	r.Route("/notification_configuration/{configurationID}", func(r chi.Router) {
		r.Use(a.notificationScope)
		r.Get("/", a.getNotificationConfiguration)
		r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateNotificationConfiguration)
		r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteNotificationConfiguration)
		r.With(a.requireWrite(canManageWorkspace)).Post("/triggers", a.createNotificationTrigger)
		r.With(a.requireWrite(canManageWorkspace)).Delete("/triggers/{triggerID}", a.deleteNotificationTrigger)
		r.With(a.requireWrite(canManageWorkspace)).Patch("/relationships/templates", a.replaceNotificationTemplates)
	})

	// Workspace run triggers: the workspace page's Run Triggers tab calls the
	// collection at the root, filtered by either end.
	r.Get("/runTrigger", a.listRunTriggers)
	r.Post("/runTrigger", a.createRunTrigger)
	r.Route("/runTrigger/{runTriggerID}", func(r chi.Router) {
		r.Get("/", a.getRunTrigger)
		r.Patch("/", a.updateRunTrigger)
		r.Delete("/", a.deleteRunTrigger)
	})

	// Workspaces
	r.Route("/workspace/{workspaceID}", func(r chi.Router) {
		r.Use(a.workspaceScope)
		r.Get("/", a.getWorkspace)
		r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateWorkspace)
		r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteWorkspace)
		r.Get("/variable", a.listVariables)
		r.With(a.requireWrite(canManageWorkspace)).Post("/variable", a.createVariable)
	})

	// The React UI reaches workspaces and their variables through the
	// organization: /organization/:org/workspace/:ws[/variable/:id].
	r.Route("/organization/{organizationID}/workspace/{workspaceID}", func(r chi.Router) {
		r.Use(a.workspaceInOrganization)
		r.Get("/", a.getWorkspace)
		r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateWorkspace)
		r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteWorkspace)
		r.Get("/variable", a.listVariables)
		r.With(a.requireWrite(canManageWorkspace)).Post("/variable", a.createVariable)
		// Schedules: the workspace page's Schedules tab.
		r.Get("/schedule", a.listSchedules)
		r.With(a.requireWrite(canManageWorkspace)).Post("/schedule", a.createSchedule)
		r.Route("/schedule/{scheduleID}", func(r chi.Router) {
			r.Get("/", a.getSchedule)
			r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateSchedule)
			r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteSchedule)
		})
		// Webhook: the workspace Webhook settings screen (#147).
		r.With(a.requireWrite(canManageWorkspace)).Post("/webhook", a.createWebhook)
		r.Route("/webhook/{webhookID}", func(r chi.Router) {
			r.Get("/", a.getWebhook)
			r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateWebhook)
			r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteWebhook)
			r.Get("/events", a.listWebhookEvents)
			r.With(a.requireWrite(canManageWorkspace)).Post("/events", a.createWebhookEvent)
			r.Route("/events/{eventID}", func(r chi.Router) {
				r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateWebhookEvent)
				r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteWebhookEvent)
			})
		})
		r.With(a.requireWrite(canManageWorkspace)).Patch("/relationships/{relation}", a.updateWorkspaceRelationship)
		// Workspace access grants: Settings > Workspace > Access.
		r.Get("/access", a.listWorkspaceAccess)
		r.With(a.requireWrite(canManageWorkspace)).Post("/access", a.createWorkspaceAccess)
		r.Route("/access/{accessID}", func(r chi.Router) {
			r.Get("/", a.getWorkspaceAccess)
			r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateWorkspaceAccess)
			r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteWorkspaceAccess)
		})
		r.Get("/workspaceTag", a.listWorkspaceTags)
		r.With(a.requireWrite(canManageWorkspace)).Post("/workspaceTag", a.createWorkspaceTag)
		r.Get("/notificationConfiguration", a.listWorkspaceNotificationConfigurations)
		r.With(a.requireWrite(canManageWorkspace)).Post("/notificationConfiguration", a.createWorkspaceNotificationConfiguration)
		r.Route("/workspaceTag/{workspaceTagID}", func(r chi.Router) {
			r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateWorkspaceTag)
			r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteWorkspaceTag)
		})
		r.Route("/variable/{variableID}", func(r chi.Router) {
			r.Use(a.variableInWorkspace)
			r.Get("/", a.getVariable)
			r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateVariable)
			r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteVariable)
		})
	})

	// Variables
	r.Route("/variable/{variableID}", func(r chi.Router) {
		r.Use(a.variableScope)
		r.Get("/", a.getVariable)
		r.With(a.requireWrite(canManageWorkspace)).Patch("/", a.updateVariable)
		r.With(a.requireWrite(canManageWorkspace)).Delete("/", a.deleteVariable)
	})

	// Elide answers an unknown path with a JSON:API error rather than an empty
	// 404, and the UI shows that detail to the user.
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "Unknown collection "+unknownCollection(r.URL.Path))
	})

	return r
}

// unknownCollection names the segment upstream blames: the collection after a
// /{collection}/{id} prefix, otherwise the first path segment.
func unknownCollection(path string) string {
	segments := strings.Split(strings.Trim(strings.TrimPrefix(path, MountPrefix), "/"), "/")
	if len(segments) > 2 {
		if _, err := uuid.Parse(segments[1]); err == nil {
			return segments[2]
		}
	}
	if len(segments) > 0 && segments[0] != "" {
		return segments[0]
	}
	return path
}

// requireUUIDParams answers a malformed id the way upstream does (400 "Invalid
// value: x"). Without it the value reaches the store, which reports no rows, and
// a client bug looks like a legitimate 404.
//
// It reads the path rather than chi's params because r.Use middlewares run
// before routing resolves them; every resource route here is
// /{collection}/{id}[/...], so the id is always the second segment. The router
// is mounted under MountPrefix, which chi leaves in URL.Path.
func requireUUIDParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, MountPrefix)
		segments := strings.Split(strings.Trim(path, "/"), "/")
		// Actions are keyed by the editor's slug ("terrakube.open-ai"), not a UUID.
		if len(segments) > 0 && segments[0] == "action" {
			next.ServeHTTP(w, r)
			return
		}
		if len(segments) > 1 && segments[1] != "" {
			if _, err := uuid.Parse(segments[1]); err != nil {
				jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Value", "Invalid value: "+segments[1])
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (a *API) authenticate(next http.Handler) http.Handler {
	if a.deps.AuthToken == "" && len(a.deps.Verifiers) == 0 {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, ok := auth.Authenticate(r.Context(), r.Header.Get("Authorization"), a.deps.AuthToken, a.deps.Verifiers...)
		if !ok {
			w.Header().Set("WWW-Authenticate", "Bearer")
			jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
			return
		}
		next.ServeHTTP(w, r.WithContext(auth.WithIdentity(r.Context(), identity)))
	})
}

// writeStoreError maps common PostgreSQL errors onto JSON:API responses. It
// returns false when the error is not recognised, so callers can fall back to a
// generic 500.
func (a *API) writeStoreError(w http.ResponseWriter, operation string, err error) bool {
	if errors.Is(err, pgx.ErrNoRows) {
		jsonapi.WriteError(w, http.StatusNotFound, "Not Found", "The requested resource does not exist.")
		return true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			jsonapi.WriteError(w, http.StatusConflict, "Conflict", "A resource with the same name already exists.")
			return true
		case "23503": // foreign_key_violation
			jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Reference", "A referenced resource does not exist.")
			return true
		case "23502": // not_null_violation
			jsonapi.WriteUnprocessable(w, http.StatusUnprocessableEntity, "Missing Attribute", "A required attribute is missing.", pgErr.ColumnName)
			return true
		}
	}

	return false
}

func (a *API) internalError(w http.ResponseWriter, operation string, err error) {
	a.deps.Logger.Error("api request failed", "operation", operation, "error", err)
	jsonapi.WriteError(w, http.StatusInternalServerError, "Internal Server Error", "The request could not be completed.")
}

// listQuery is the supported subset of Elide's collection parameters.
type listQuery struct {
	Filter     rsql.Node
	Sorts      []store.SortKey
	PageSize   int
	PageNumber int
	// Paged records that the request carried a page parameter at all, because
	// Elide emits meta.page only then.
	Paged  bool
	Totals bool
}

// elideDefaultPageLimit is Elide's limit when page[number] is sent alone.
const elideDefaultPageLimit = 10000

// page is the jsonapi.Page a collection response reports.
func (q listQuery) page() jsonapi.Page {
	if !q.Paged {
		return jsonapi.Page{}
	}
	return jsonapi.Page{Number: q.PageNumber, Limit: q.PageSize, Totals: q.Totals}
}

// pageWindow maps page[size]/page[number] onto LIMIT/OFFSET. A nil limit means
// the request carried no page size, which Elide answers with every record.
func (q listQuery) pageWindow() (*int64, int64) {
	if !q.Paged || q.PageSize == 0 {
		return nil, 0
	}
	limit := int64(q.PageSize)
	return &limit, int64(q.PageNumber-1) * limit
}

// parseInclude turns a comma-separated include parameter into a set. Unknown
// relationships are ignored: the UI sends every relationship it knows about,
// including ones whose tables we do not have yet.
func parseInclude(r *http.Request) map[string]bool {
	want := map[string]bool{}
	for _, name := range strings.Split(r.URL.Query().Get("include"), ",") {
		if name = strings.TrimSpace(name); name != "" {
			want[name] = true
		}
	}
	return want
}

// parseList reads filter[type], sort and page[*] from the query string. A
// filter the RSQL parser rejects is an error the handler answers with 400.
func parseList(r *http.Request, filterType string) (listQuery, error) {
	values := r.URL.Query()

	q := listQuery{PageNumber: 1}
	if raw := values.Get("filter[" + filterType + "]"); raw != "" {
		filter, err := rsql.Parse(raw)
		if err != nil {
			return q, err
		}
		q.Filter = filter
	}

	for _, key := range strings.Split(values.Get("sort"), ",") {
		key = strings.TrimSpace(key)
		desc := strings.HasPrefix(key, "-")
		key = strings.TrimPrefix(key, "-")
		if key != "" {
			q.Sorts = append(q.Sorts, store.SortKey{Column: key, Desc: desc})
		}
	}

	if raw := values.Get("page[size]"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			q.PageSize = n
		}
	}
	if raw := values.Get("page[number]"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			q.PageNumber = n
		}
	}
	_, q.Totals = values["page[totals]"]
	_, sized := values["page[size]"]
	_, numbered := values["page[number]"]
	q.Paged = q.Totals || sized || numbered
	if q.Paged && q.PageSize == 0 {
		q.PageSize = elideDefaultPageLimit
	}

	return q, nil
}

func param(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

func strValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func strPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func boolValue(v *bool, fallback bool) bool {
	if v == nil {
		return fallback
	}
	return *v
}

func requireType(w http.ResponseWriter, data *jsonapi.RequestData, want string) bool {
	if data.Type != want {
		jsonapi.WriteError(w, http.StatusConflict, "Invalid Type", fmt.Sprintf("Expected resource type %q, got %q.", want, data.Type))
		return false
	}
	return true
}
