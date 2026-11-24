package tfstate

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"strings"

	internalauth "github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type claimsContextKey struct{}
type externalUserContextKey struct{}

func withClaims(ctx context.Context, claims *internalauth.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

func claimsFromContext(ctx context.Context) (*internalauth.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*internalauth.Claims)
	return claims, ok
}

func withExternalUser(ctx context.Context) context.Context {
	return context.WithValue(ctx, externalUserContextKey{}, true)
}

// isExternalUser reports a request authenticated with the configured user
// token rather than an internal executor JWT. External users may read state
// (the Terraform CLI downloads hosted state with its own token) but never
// write the backend.
func isExternalUser(ctx context.Context) bool {
	external, _ := ctx.Value(externalUserContextKey{}).(bool)
	return external
}

// authenticate validates the executor JWT and, for workspace-scoped tokens,
// confirms that the token still points at an active step. Archive uploads are
// intentionally outside this middleware because upstream authenticates them
// with the unguessable archive ID instead.
func (h *Handler) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := internalauth.ExtractToken(r.Header.Get("Authorization"))
		if raw == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// A user token (the development bearer token, later a Dex JWT) can read
		// hosted state. Backend state writes still require a workspace token,
		// which the handlers enforce.
		if h.deps.AuthToken != "" && subtle.ConstantTimeCompare([]byte(raw), []byte(h.deps.AuthToken)) == 1 {
			next.ServeHTTP(w, r.WithContext(withExternalUser(r.Context())))
			return
		}

		if h.deps.Auth == nil {
			h.deps.Logger.Error("state backend authentication is not configured; failing closed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := h.deps.Auth.Parse(raw)
		if err != nil {
			h.deps.Logger.Info("state backend token rejected", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch internalauth.Classify(claims) {
		case internalauth.TypeWorkspaceExecutor:
			if err := ValidateActiveStep(r.Context(), h.deps.Store, claims); err != nil {
				h.deps.Logger.Info("state backend runner token rejected", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		case internalauth.TypeExecutorService:
			// The executor service token is allowed for read/manage state
			// operations, but backend state writes still require a workspace
			// token (checked by the handlers).
		default:
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r.WithContext(withClaims(r.Context(), claims)))
	})
}

// ValidateActiveStep reports whether a workspace token still points at a step
// that has not finished. The TFE and log APIs use it so a runner token expires
// when its job ends, like the state backend.
func ValidateActiveStep(ctx context.Context, q store.Querier, claims *internalauth.Claims) error {
	stepID := claims.EffectiveStepID()
	if stepID == "" {
		return errors.New("runner token has no stepId")
	}

	step, err := q.GetStep(ctx, stepID)
	if err != nil {
		return fmt.Errorf("load step %s: %w", stepID, err)
	}
	if !isActiveStepStatus(step.Status) {
		return fmt.Errorf("step %s is not active", stepID)
	}
	return nil
}

func isActiveStepStatus(status *string) bool {
	if status == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(*status)) {
	case "running", "queue", "pending", "waitingapproval", "approved":
		return true
	default:
		return false
	}
}
