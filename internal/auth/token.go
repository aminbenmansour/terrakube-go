// Package auth validates the internal JWTs Terrakube uses between the API,
// executor and registry. The token format mirrors the upstream Java
// implementation: HMAC-signed JWTs with the TerrakubeInternal issuer and
// workspace/job/step claims.
package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Issuer is the issuer of every internal Terrakube token.
	Issuer = "TerrakubeInternal"
	// SubjectExecutor identifies the executor service token.
	SubjectExecutor = "TerrakubeInternal (EXECUTOR)"
	// SubjectRegistry identifies the registry service token.
	SubjectRegistry = "TerrakubeInternal (REGISTRY)"
)

type claimsContextKey struct{}

// WithClaims attaches validated claims to a request context.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

// ClaimsFromContext returns the claims a middleware attached to the request.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*Claims)
	return claims, ok
}

// Type classifies an internal token. It mirrors InternalTokenClassifier.
type Type int

const (
	TypeNotInternal Type = iota
	TypeWorkspaceExecutor
	TypeExecutorService
	TypeRegistryService
	TypeUnknownInternal
)

func (t Type) String() string {
	switch t {
	case TypeWorkspaceExecutor:
		return "workspace-executor"
	case TypeExecutorService:
		return "executor-service"
	case TypeRegistryService:
		return "registry-service"
	case TypeUnknownInternal:
		return "unknown-internal"
	default:
		return "not-internal"
	}
}

// Claims are the JWT claims used by the API and executor.
type Claims struct {
	jwt.RegisteredClaims
	OrganizationID string `json:"organizationId,omitempty"`
	WorkspaceID    string `json:"workspaceId,omitempty"`
	JobID          string `json:"jobId,omitempty"`
	StepID         string `json:"stepId,omitempty"`
}

// Classify returns the internal token type using the same rules as upstream.
func Classify(claims *Claims) Type {
	if claims == nil || claims.Issuer != Issuer {
		return TypeNotInternal
	}
	if strings.TrimSpace(claims.WorkspaceID) != "" {
		return TypeWorkspaceExecutor
	}
	switch claims.Subject {
	case SubjectExecutor:
		return TypeExecutorService
	case SubjectRegistry:
		return TypeRegistryService
	default:
		return TypeUnknownInternal
	}
}

// EffectiveStepID returns the stepId claim, falling back to the upstream
// jti format "jobId-stepId" used by the executor.
func (c *Claims) EffectiveStepID() string {
	if c == nil {
		return ""
	}
	if step := strings.TrimSpace(c.StepID); step != "" {
		return step
	}
	jti := strings.TrimSpace(c.ID)
	// Two standard 36-character UUIDs separated by a hyphen have length 73
	// with the separator at index 36.
	if len(jti) >= 73 && jti[36] == '-' {
		return jti[37:]
	}
	return ""
}

// Validator verifies internal JWTs with the configured shared secret.
type Validator struct {
	key []byte
	alg jwt.SigningMethod
}

// NewValidator builds a validator from the base64url-encoded shared secret.
// The HMAC algorithm follows the decoded key length, matching upstream.
func NewValidator(secret string) (*Validator, error) {
	key, err := decodeSecret(secret)
	if err != nil {
		return nil, err
	}
	if len(key) < 32 {
		return nil, fmt.Errorf("auth: internal token secret must decode to at least 32 bytes")
	}

	var alg jwt.SigningMethod
	switch {
	case len(key) >= 64:
		alg = jwt.SigningMethodHS512
	case len(key) >= 48:
		alg = jwt.SigningMethodHS384
	default:
		alg = jwt.SigningMethodHS256
	}
	return &Validator{key: key, alg: alg}, nil
}

// Parse verifies the signature and registered claims, then returns the
// Terrakube claims.
func (v *Validator) Parse(raw string) (*Claims, error) {
	if v == nil {
		return nil, errors.New("auth: validator is not configured")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("auth: token is empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method != v.alg {
			return nil, fmt.Errorf("auth: unexpected signing method %q", token.Method.Alg())
		}
		return v.key, nil
	}, jwt.WithValidMethods([]string{v.alg.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("auth: parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("auth: token is not valid")
	}
	return claims, nil
}

// Sign mints an internal token. The executor will use the same helper once it
// lands; tests use it to exercise the API without a running executor.
func Sign(secret string, claims Claims) (string, error) {
	key, err := decodeSecret(secret)
	if err != nil {
		return "", err
	}
	if len(key) < 32 {
		return "", fmt.Errorf("auth: internal token secret must decode to at least 32 bytes")
	}

	var alg jwt.SigningMethod
	switch {
	case len(key) >= 64:
		alg = jwt.SigningMethodHS512
	case len(key) >= 48:
		alg = jwt.SigningMethodHS384
	default:
		alg = jwt.SigningMethodHS256
	}

	now := time.Now()
	if claims.Issuer == "" {
		claims.Issuer = Issuer
	}
	if claims.Subject == "" {
		claims.Subject = SubjectExecutor
	}
	if len(claims.Audience) == 0 {
		claims.Audience = jwt.ClaimStrings{Issuer}
	}
	if claims.IssuedAt == nil {
		claims.IssuedAt = jwt.NewNumericDate(now)
	}
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(30 * time.Minute))
	}
	if claims.ID == "" && claims.JobID != "" && claims.StepID != "" {
		claims.ID = claims.JobID + "-" + claims.StepID
	}

	return jwt.NewWithClaims(alg, claims).SignedString(key)
}

// ExtractToken returns the token from an Authorization header. Both the
// executor's Basic username/password form and a Bearer token are accepted.
func ExtractToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	if len(header) >= 6 && strings.EqualFold(header[:6], "Basic ") {
		raw := strings.TrimSpace(header[6:])
		if raw == "" {
			return ""
		}
		decoded, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return ""
		}
		colon := strings.IndexByte(string(decoded), ':')
		if colon <= 0 {
			return ""
		}
		return strings.TrimSpace(string(decoded)[colon+1:])
	}
	if len(header) >= 7 && strings.EqualFold(header[:7], "Bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return ""
}

func decodeSecret(secret string) ([]byte, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("auth: internal token secret is not configured")
	}
	key, err := base64.RawURLEncoding.DecodeString(secret)
	if err == nil {
		return key, nil
	}
	key, err = base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("auth: decode internal token secret: %w", err)
	}
	return key, nil
}
