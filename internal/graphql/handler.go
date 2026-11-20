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
	"time"

	gql "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/aminbenmansour/terrakube-go/internal/jobstatus"
	"github.com/aminbenmansour/terrakube-go/internal/jsonapi"
)

// EventSource streams job status changes to subscriptions.
type EventSource interface {
	SubscribeJobStatus(ctx context.Context) (<-chan jobstatus.Event, error)
}

// Deps are the collaborators the GraphQL endpoint needs.
type Deps struct {
	DB        DB
	Events    EventSource
	Logger    *slog.Logger
	AuthToken string
}

// New builds the /graphql/api/v1 handler.
func New(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}

	srv := handler.New(NewExecutableSchema(Config{Resolvers: &Resolver{db: deps.DB, events: deps.Events}}))
	srv.AddTransport(transport.POST{})
	srv.SetErrorPresenter(deps.errorPresenter)

	if deps.AuthToken == "" {
		deps.Logger.Warn("GraphQL authentication is disabled: set TERRAKUBE_AUTH_TOKEN to require a bearer token")
		return srv
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !bearerMatches(r.Header.Get("Authorization"), deps.AuthToken) {
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

// NewSubscriptions builds the /subscriptions handler. The UI opens a
// graphql-ws connection here and sends the bearer token as a connection
// parameter, not as a header, so it is checked after the socket upgrade.
func NewSubscriptions(deps Deps) http.Handler {
	if deps.Logger == nil {
		deps.Logger = slog.Default()
	}

	srv := handler.New(NewExecutableSchema(Config{Resolvers: &Resolver{db: deps.DB, events: deps.Events}}))
	srv.SetErrorPresenter(deps.errorPresenter)
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		InitFunc: func(ctx context.Context, payload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
			if deps.AuthToken == "" {
				return ctx, nil, nil
			}
			if !bearerMatches(payload.Authorization(), deps.AuthToken) {
				return nil, nil, errors.New("unauthorized")
			}
			return ctx, nil, nil
		},
	})
	return srv
}

func bearerMatches(header, token string) bool {
	raw, ok := strings.CutPrefix(header, "Bearer ")
	return ok && subtle.ConstantTimeCompare([]byte(raw), []byte(token)) == 1
}
