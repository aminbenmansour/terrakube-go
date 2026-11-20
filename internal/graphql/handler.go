// Package graphql serves the Elide-compatible GraphQL endpoint the upstream UI
// uses at /graphql/api/v1. The schema only declares the entities and fields the
// UI reads; see docs/upstream/ui-operations.md.
package graphql

import (
	"context"
	"crypto/subtle"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
)

// Deps are the collaborators the GraphQL endpoint needs.
type Deps struct {
	DB        DB
	Logger    *slog.Logger
	AuthToken string
}

// New builds the /graphql/api/v1 handler.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}

	srv := handler.New(NewExecutableSchema(Config{Resolvers: &Resolver{db: deps.DB}}))
	srv.AddTransport(transport.POST{})
	srv.SetErrorPresenter(deps.errorPresenter)

	if deps.AuthToken == "" {
		deps.Logger.Warn("GraphQL authentication is disabled: set TERRAKUBE_AUTH_TOKEN to require a bearer token")
		return srv
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) ||
			subtle.ConstantTimeCompare([]byte(strings.TrimPrefix(header, prefix)), []byte(deps.AuthToken)) != 1 {
			w.Header().Set("WWW-Authenticate", "Bearer")
			jsonapi.WriteError(w, http.StatusUnauthorized, "Unauthorized", "A valid bearer token is required.")
			return
		}
		srv.ServeHTTP(w, r)
	})
}

// errorPresenter logs resolver failures and hides their text. Query validation
// errors keep their message: they describe the client's document, not our data.
func (deps Deps) errorPresenter(ctx context.Context, err error) *gqlerror.Error {
	presented := gql.DefaultErrorPresenter(ctx, err)
	var gqlErr *gqlerror.Error
	if errors.As(err, &gqlErr) && len(gqlErr.Path) == 0 {
		return presented
	}
	deps.Logger.Error("graphql request failed", "error", err)
	presented.Message = "The request could not be completed."
	return presented
}
