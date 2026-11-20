// Package api implements the Terrakube management API at /api/v1.
//
// Upstream Terrakube serves this surface through Elide, so requests and
// responses are JSON:API 1.0 documents with camelCase attribute names.
package api

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// Deps are the collaborators the management API needs.
type Deps struct {
	Store     store.Querier
	Logger    *slog.Logger
	PublicURL string
	AuthToken string
	// CORSOrigins lists the origins allowed to call the API from a browser.
	// Empty allows any origin, which is the historical behaviour; a UI served
	// by this API shares its origin and needs no entry.
	CORSOrigins []string
}

// API holds the HTTP handlers.
type API struct {
	deps Deps
}

var errProjectOutsideOrganization = errors.New("project does not belong to organization")

// New builds the /api/v1 router.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}
	if deps.AuthToken == "" {
		deps.Logger.Warn("management API authentication is disabled: set TERRAKUBE_AUTH_TOKEN to require a bearer token")
	}

	a := &API{deps: deps}

	r := chi.NewRouter()
	r.Use(corsMiddleware(deps.CORSOrigins))
	r.Use(a.authenticate)

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
	})

	// Projects
	r.Route("/project/{projectID}", func(r chi.Router) {
		r.Get("/", a.getProject)
		r.Patch("/", a.updateProject)
		r.Delete("/", a.deleteProject)
		r.Get("/workspace", a.listProjectWorkspaces)
	})

	// Workspaces
	r.Route("/workspace/{workspaceID}", func(r chi.Router) {
		r.Get("/", a.getWorkspace)
		r.Patch("/", a.updateWorkspace)
		r.Delete("/", a.deleteWorkspace)
		r.Get("/variable", a.listVariables)
		r.Post("/variable", a.createVariable)
	})

	// Variables
	r.Route("/variable/{variableID}", func(r chi.Router) {
		r.Get("/", a.getVariable)
		r.Patch("/", a.updateVariable)
		r.Delete("/", a.deleteVariable)
	})

	return r
}

func corsMiddleware(origins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		allowed[strings.TrimSpace(origin)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			switch {
			case len(allowed) == 0:
				w.Header().Set("Access-Control-Allow-Origin", "*")
			case allowed[origin]:
				// Echoing the origin keeps credentials and any future
				// Vary-sensitive caching correct.
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (a *API) authenticate(next http.Handler) http.Handler {
	if a.deps.AuthToken == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) ||
			subtle.ConstantTimeCompare([]byte(strings.TrimPrefix(header, prefix)), []byte(a.deps.AuthToken)) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
			return
		}
		next.ServeHTTP(w, r)
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
