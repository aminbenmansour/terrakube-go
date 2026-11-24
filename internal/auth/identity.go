package auth

import (
	"context"
	"crypto/subtle"
	"strings"
)

// DevIssuer marks the identity behind the shared development token.
const DevIssuer = "development"

// Identity is the authenticated caller behind a bearer token.
type Identity struct {
	Issuer  string
	Subject string
	Email   string
	Name    string
	Groups  []string
}

type identityKey struct{}

// WithIdentity attaches the caller to a request context.
func WithIdentity(ctx context.Context, identity *Identity) context.Context {
	return context.WithValue(ctx, identityKey{}, identity)
}

// IdentityFromContext returns the caller the middleware authenticated.
func IdentityFromContext(ctx context.Context) (*Identity, bool) {
	identity, ok := ctx.Value(identityKey{}).(*Identity)
	return identity, ok && identity != nil
}

// Authenticate resolves an Authorization header into an identity. The
// development token is checked first, then each verifier in order; a verifier
// that rejects the token does not stop the next one from trying.
func Authenticate(ctx context.Context, header, devToken string, verifiers ...Verifier) (*Identity, bool) {
	raw, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return nil, false
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}
	if devToken != "" && subtle.ConstantTimeCompare([]byte(raw), []byte(devToken)) == 1 {
		return &Identity{Issuer: DevIssuer, Subject: "dev", Email: "dev@localhost", Name: "dev"}, true
	}
	for _, verifier := range verifiers {
		if verifier == nil {
			continue
		}
		identity, err := verifier.Verify(ctx, raw)
		if err == nil && identity != nil {
			return identity, true
		}
	}
	return nil, false
}
