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
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Deps are the collaborators the management API needs.
type Deps struct {
	Store     store.Querier
	Logger    *slog.Logger
	PublicURL string
	AuthToken string
	// Verifiers accept user bearer tokens (Dex, PATs). Empty leaves only the
	// development token, or no authentication at all when that is empty too.
	Verifiers []auth.Verifier
	// PATs backs /pat/v1. Nil answers 503 there.
	PATs PATService
	// CORSOrigins lists the origins allowed to call the API from a browser.
	// Empty allows any origin, which is the historical behaviour; a UI served
	// by this API shares its origin and needs no entry.
	CORSOrigins []string
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
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.AuthToken == "" && len(deps.Verifiers) == 0 {
		deps.Logger.Warn("management API authentication is disabled: set TERRAKUBE_AUTH_TOKEN or TERRAKUBE_DEX_ENABLED to require a bearer token")
	}

	a := &API{deps: deps}

	r := chi.NewRouter()
	r.Use(httpapi.CORS(deps.CORSOrigins))
	r.Use(a.authenticate)
	r.Use(requireUUIDParams)

	// JSON:API atomic operations; the UI creates and updates workspaces here.
	r.Post("/operations", a.atomicOperations)

	// The action catalog fills the workspace resource drawer and the Settings >
	// Actions editor.
	r.Get("/action", a.listActions)
	r.Post("/action", a.createAction)
	r.Route("/action/{actionID}", func(r chi.Router) {
		r.Get("/", a.getAction)
		r.Patch("/", a.updateAction)
		r.Delete("/", a.deleteAction)
	})

	// Organizations
	r.Get("/organization", a.listOrganizations)
	r.Post("/organization", a.createOrganization)
	r.Route("/organization/{organizationID}", func(r chi.Router) {
		r.Get("/", a.getOrganization)
		r.Patch("/", a.updateOrganization)
		r.Delete("/", a.deleteOrganization)

		// Projects and workspaces scoped to an organization.
		r.Get("/project", a.listProjects)
		r.Post("/project", a.createProject)
		r.Get("/workspace", a.listWorkspaces)
		r.Post("/workspace", a.createWorkspace)

		// Runs.
		r.Get("/job", a.listJobs)
		r.Route("/job/{jobID}", func(r chi.Router) {
			r.Get("/", a.getJob)
			r.Patch("/", a.updateJob)
		})

		// Tags: Settings > Tags and the workspace tag editor.
		r.Get("/tag", a.listTags)
		r.Post("/tag", a.createTag)
		r.Route("/tag/{tagID}", func(r chi.Router) {
			r.Get("/", a.getTag)
			r.Patch("/", a.updateTag)
			r.Delete("/", a.deleteTag)
		})

		// VCS providers: Settings > VCS and the workspace source picker.
		r.Get("/vcs", a.listVCS)
		r.Post("/vcs", a.createVCS)
		r.Route("/vcs/{vcsID}", func(r chi.Router) {
			r.Get("/", a.getVCS)
			r.Patch("/", a.updateVCS)
			r.Delete("/", a.deleteVCS)
		})

		// Notifications: Settings > Notifications, org-wide and per workspace.
		r.Get("/notificationConfiguration", a.listNotificationConfigurations)
		r.Post("/notificationConfiguration", a.createNotificationConfiguration)
		r.Get("/notificationConfiguration/{configurationID}", a.getNotificationConfiguration)

		// Teams: Settings > Teams and the access pickers.
		r.Get("/team", a.listTeams)
		r.Post("/team", a.createTeam)
		r.Route("/team/{teamID}", func(r chi.Router) {
			r.Get("/", a.getTeam)
			r.Patch("/", a.updateTeam)
			r.Delete("/", a.deleteTeam)
		})

		// Job-flow templates: Settings > Templates and the workspace default.
		r.Get("/template", a.listTemplates)
		r.Post("/template", a.createTemplate)
		r.Route("/template/{templateID}", func(r chi.Router) {
			r.Get("/", a.getTemplate)
			r.Patch("/", a.updateTemplate)
			r.Delete("/", a.deleteTemplate)
		})

		// SSH keys and self-hosted agent pools for workspace connections.
		r.Get("/ssh", a.listSSHKeys)
		r.Post("/ssh", a.createSSHKey)
		r.Route("/ssh/{sshID}", func(r chi.Router) {
			r.Get("/", a.getSSHKey)
			r.Patch("/", a.updateSSHKey)
			r.Delete("/", a.deleteSSHKey)
		})
		r.Get("/agent", a.listAgents)
		r.Post("/agent", a.createAgent)
		r.Delete("/agent/{agentID}", a.deleteAgent)

		// Organization-level variables inherited by every workspace.
		r.Get("/globalvar", a.listGlobalVariables)
		r.Post("/globalvar", a.createGlobalVariable)
		r.Route("/globalvar/{globalVarID}", func(r chi.Router) {
			r.Get("/", a.getGlobalVariable)
			r.Patch("/", a.updateGlobalVariable)
			r.Delete("/", a.deleteGlobalVariable)
		})
	})

	// Projects
	r.Route("/project/{projectID}", func(r chi.Router) {
		r.Get("/", a.getProject)
		r.Patch("/", a.updateProject)
		r.Delete("/", a.deleteProject)
		r.Get("/workspace", a.listProjectWorkspaces)
	})

	// Variable collections: Settings > Collections, and the read-only
	// "Collection Variables" table on the workspace page.
	r.Route("/organization/{organizationID}/collection", func(r chi.Router) {
		r.Get("/", a.listCollections)
		r.Post("/", a.createCollection)
		r.Route("/{collectionID}", func(r chi.Router) {
			r.Get("/", a.getCollection)
			r.Patch("/", a.updateCollection)
			r.Delete("/", a.deleteCollection)
			r.Get("/item", a.listCollectionItems)
			r.Post("/item", a.createCollectionItem)
			r.Route("/item/{itemID}", func(r chi.Router) {
				r.Patch("/", a.updateCollectionItem)
				r.Delete("/", a.deleteCollectionItem)
			})
			r.Get("/reference", a.listCollectionReferences)
			r.Post("/reference", a.createCollectionReference)
			r.Delete("/reference/{referenceID}", a.deleteCollectionReference)
		})
	})
	r.Get("/reference/{referenceID}/collection", a.getReferenceCollection)

	// Notification configurations are addressed by id at the root because the
	// UI's editor and trigger sub-resource calls are not organization-nested.
	r.Route("/notification_configuration/{configurationID}", func(r chi.Router) {
		r.Get("/", a.getNotificationConfiguration)
		r.Patch("/", a.updateNotificationConfiguration)
		r.Delete("/", a.deleteNotificationConfiguration)
		r.Post("/triggers", a.createNotificationTrigger)
		r.Delete("/triggers/{triggerID}", a.deleteNotificationTrigger)
		r.Patch("/relationships/templates", a.replaceNotificationTemplates)
	})

	// Workspaces
	r.Route("/workspace/{workspaceID}", func(r chi.Router) {
		r.Get("/", a.getWorkspace)
		r.Patch("/", a.updateWorkspace)
		r.Delete("/", a.deleteWorkspace)
		r.Get("/variable", a.listVariables)
		r.Post("/variable", a.createVariable)
	})

	// The React UI reaches workspaces and their variables through the
	// organization: /organization/:org/workspace/:ws[/variable/:id].
	r.Route("/organization/{organizationID}/workspace/{workspaceID}", func(r chi.Router) {
		r.Use(a.workspaceInOrganization)
		r.Get("/", a.getWorkspace)
		r.Patch("/", a.updateWorkspace)
		r.Delete("/", a.deleteWorkspace)
		r.Get("/variable", a.listVariables)
		r.Post("/variable", a.createVariable)
		r.Patch("/relationships/{relation}", a.updateWorkspaceRelationship)
		r.Get("/workspaceTag", a.listWorkspaceTags)
		r.Post("/workspaceTag", a.createWorkspaceTag)
		r.Get("/notificationConfiguration", a.listWorkspaceNotificationConfigurations)
		r.Post("/notificationConfiguration", a.createWorkspaceNotificationConfiguration)
		r.Route("/workspaceTag/{workspaceTagID}", func(r chi.Router) {
			r.Patch("/", a.updateWorkspaceTag)
			r.Delete("/", a.deleteWorkspaceTag)
		})
		r.Route("/variable/{variableID}", func(r chi.Router) {
			r.Use(a.variableInWorkspace)
			r.Get("/", a.getVariable)
			r.Patch("/", a.updateVariable)
			r.Delete("/", a.deleteVariable)
		})
	})

	// Variables
	r.Route("/variable/{variableID}", func(r chi.Router) {
		r.Get("/", a.getVariable)
		r.Patch("/", a.updateVariable)
		r.Delete("/", a.deleteVariable)
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
	FilterValue  string
	FilterPrefix bool
	SortDesc     bool
	PageSize     int
	PageNumber   int
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

// parseList reads filter[type], sort and page[*] from the query string. Unknown
// filters are ignored until they are explicitly implemented.
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

func parseList(r *http.Request, filterType string) listQuery {
	values := r.URL.Query()

	q := listQuery{PageNumber: 1}
	if raw := values.Get("filter[" + filterType + "]"); raw != "" {
		if _, value, ok := strings.Cut(raw, "=="); ok {
			value = strings.TrimSpace(value)
			if strings.HasSuffix(value, "*") {
				q.FilterPrefix = true
				value = strings.TrimSuffix(value, "*")
			}
			q.FilterValue = strings.Trim(value, `"`)
		}
	}

	switch values.Get("sort") {
	case "-name":
		q.SortDesc = true
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

	return q
}

// filterSortPage applies the in-memory filter, sort and page parameters. The
// collections here are small; this keeps the first milestone free of dynamic
// SQL while preserving the observable behaviour.
func filterSortPage[T any](items []T, q listQuery, name func(T) string) ([]T, int) {
	filtered := make([]T, 0, len(items))
	for _, item := range items {
		n := strings.ToLower(name(item))
		filter := strings.ToLower(q.FilterValue)
		switch {
		case q.FilterValue == "":
			filtered = append(filtered, item)
		case q.FilterPrefix && strings.HasPrefix(n, filter):
			filtered = append(filtered, item)
		case !q.FilterPrefix && n == filter:
			filtered = append(filtered, item)
		}
	}

	// Insertion sort keeps the helper dependency-free and the lists are small.
	for i := 1; i < len(filtered); i++ {
		for j := i; j > 0; j-- {
			less := strings.ToLower(name(filtered[j])) < strings.ToLower(name(filtered[j-1]))
			if q.SortDesc {
				less = !less
			}
			if !less {
				break
			}
			filtered[j], filtered[j-1] = filtered[j-1], filtered[j]
		}
	}

	total := len(filtered)
	if q.PageSize == 0 {
		return filtered, total
	}
	start := (q.PageNumber - 1) * q.PageSize
	if start >= total {
		return []T{}, total
	}
	end := start + q.PageSize
	if end > total {
		end = total
	}
	return filtered[start:end], total
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
