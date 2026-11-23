// Package pat issues and validates the personal access tokens the upstream UI
// and CLI use. The token is an HMAC-signed JWT; only its metadata row is stored,
// and revocation is a soft delete the verifier checks on every request.
package pat

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	// Issuer is both the iss and aud claim upstream signs.
	Issuer = "Terrakube"
	// SourceAPI is the upstream source for tokens minted from the UI/API.
	SourceAPI = "API"
)

var (
	ErrInactive        = errors.New("pat: token is revoked, unknown or expired")
	ErrNotFound        = errors.New("pat: personal access token not found")
	ErrForbidden       = errors.New("pat: user is not authorized to delete this token")
	ErrNonInteractive  = errors.New("pat: only an interactive user can create tokens")
	ErrInvalidLifetime = errors.New("pat: token lifetime cannot be negative")
)

// Store is the subset of the generated queries the service needs.
type Store interface {
	CreatePat(ctx context.Context, arg store.CreatePatParams) (store.Pat, error)
	GetPat(ctx context.Context, id string) (store.Pat, error)
	ListPatsByCreator(ctx context.Context, createdBy *string) ([]store.Pat, error)
	SoftDeletePat(ctx context.Context, arg store.SoftDeletePatParams) error
	StampPatLastUsed(ctx context.Context, id string) error
}

// Service mints, verifies and revokes tokens.
type Service struct {
	store Store
	key   []byte
	alg   jwt.SigningMethod
	owner string
}

// New builds the service. The secret is base64url and must decode to at least
// 32 bytes, like upstream's PatSecret.
func New(store Store, secret, owner string) (*Service, error) {
	key, err := decodeSecret(secret)
	if err != nil {
		return nil, err
	}
	if len(key) < 32 {
		return nil, errors.New("pat: secret must decode to at least 32 bytes")
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
	return &Service{store: store, key: key, alg: alg, owner: strings.TrimSpace(owner)}, nil
}

// Token is the metadata the UI lists. The signed token is never stored.
type Token struct {
	ID          string
	Days        *int32
	Description *string
	Deleted     bool
	Source      string
	LastUsedAt  *time.Time
	CreatedDate *time.Time
	UpdatedDate *time.Time
	CreatedBy   *string
	UpdatedBy   *string
}

// Create stores the metadata row and returns the signed token.
func (s *Service) Create(ctx context.Context, user auth.Identity, days int, description string) (string, error) {
	if user.Issuer == Issuer {
		return "", ErrNonInteractive
	}
	if days < 0 {
		return "", ErrInvalidLifetime
	}

	id := uuid.NewString()
	daysValue := int32(days)
	row, err := s.store.CreatePat(ctx, store.CreatePatParams{
		ID:          id,
		Days:        &daysValue,
		Description: &description,
		Source:      SourceAPI,
		CreatedBy:   stringPtr(user.Email),
		UpdatedBy:   stringPtr(user.Email),
	})
	if err != nil {
		return "", fmt.Errorf("pat: create metadata: %w", err)
	}

	name := user.Name
	if name == "" {
		name = user.Email
	}
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":            Issuer,
		"sub":            name + " (Token)",
		"aud":            []string{Issuer},
		"jti":            row.ID,
		"email":          user.Email,
		"email_verified": true,
		"name":           name + " (Token)",
		"groups":         user.Groups,
		"iat":            now.Unix(),
	}
	if days > 0 {
		claims["exp"] = now.AddDate(0, 0, days).Unix()
	}

	token, err := jwt.NewWithClaims(s.alg, claims).SignedString(s.key)
	if err != nil {
		// Do not leave a row behind that no client can ever use.
		_ = s.store.SoftDeletePat(ctx, store.SoftDeletePatParams{ID: row.ID, UpdatedBy: stringPtr(user.Email)})
		return "", fmt.Errorf("pat: sign token: %w", err)
	}
	return token, nil
}

// Verify checks the signature, issuer, expiry and the live metadata row, then
// stamps last use (at most hourly, enforced by the query). It never returns the
// token in an error.
func (s *Service) Verify(ctx context.Context, raw string) (*auth.Identity, error) {
	if s == nil || len(s.key) == 0 {
		return nil, ErrInactive
	}
	claims := jwt.MapClaims{}
	// Audience is not checked: upstream's PAT decoder validates the signature,
	// issuer and expiry only.
	_, err := jwt.NewParser(
		jwt.WithValidMethods([]string{s.alg.Alg()}),
		jwt.WithIssuer(Issuer),
		jwt.WithLeeway(60*time.Second),
	).ParseWithClaims(strings.TrimSpace(raw), claims, func(*jwt.Token) (any, error) {
		return s.key, nil
	})
	if err != nil {
		return nil, ErrInactive
	}

	id, _ := claims["jti"].(string)
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInactive
	}
	row, err := s.store.GetPat(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && row.Deleted) {
		return nil, ErrInactive
	}
	if err != nil {
		return nil, fmt.Errorf("pat: load metadata: %w", err)
	}
	// A failed stamp must not fail the request; the token stays valid.
	_ = s.store.StampPatLastUsed(ctx, id)

	subject, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)
	return &auth.Identity{
		Issuer:  Issuer,
		Subject: subject,
		Email:   email,
		Name:    name,
		Groups:  stringSlice(claims["groups"]),
	}, nil
}

// List returns the caller's active tokens.
func (s *Service) List(ctx context.Context, email string) ([]Token, error) {
	rows, err := s.store.ListPatsByCreator(ctx, stringPtr(email))
	if err != nil {
		return nil, fmt.Errorf("pat: list: %w", err)
	}
	tokens := make([]Token, 0, len(rows))
	for _, row := range rows {
		if row.Deleted {
			continue
		}
		tokens = append(tokens, Token{
			ID:          row.ID,
			Days:        row.Days,
			Description: row.Description,
			Deleted:     row.Deleted,
			Source:      row.Source,
			LastUsedAt:  timePtr(row.LastUsedAt),
			CreatedDate: timePtr(row.CreatedDate),
			UpdatedDate: timePtr(row.UpdatedDate),
			CreatedBy:   row.CreatedBy,
			UpdatedBy:   row.UpdatedBy,
		})
	}
	return tokens, nil
}

// Delete revokes a token. The caller must own it, or be the configured
// instance owner (by email or group, as upstream does).
func (s *Service) Delete(ctx context.Context, id string, user auth.Identity) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ErrNotFound
	}
	row, err := s.store.GetPat(ctx, parsed.String())
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && row.Deleted) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("pat: load metadata: %w", err)
	}
	if !s.canDelete(user, row) {
		return ErrForbidden
	}
	if err := s.store.SoftDeletePat(ctx, store.SoftDeletePatParams{ID: row.ID, UpdatedBy: stringPtr(user.Email)}); err != nil {
		return fmt.Errorf("pat: revoke: %w", err)
	}
	return nil
}

func (s *Service) canDelete(user auth.Identity, row store.Pat) bool {
	if row.CreatedBy != nil && strings.EqualFold(user.Email, *row.CreatedBy) {
		return true
	}
	if s.owner == "" {
		return false
	}
	if strings.EqualFold(user.Email, s.owner) {
		return true
	}
	for _, group := range user.Groups {
		if group == s.owner {
			return true
		}
	}
	return false
}

func stringPtr(value string) *string { return &value }

func timePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	utc := value.Time.UTC().Truncate(time.Millisecond)
	return &utc
}

func stringSlice(value any) []string {
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				values = append(values, text)
			}
		}
		return values
	default:
		return nil
	}
}

func decodeSecret(secret string) ([]byte, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("pat: secret is not configured")
	}
	key, err := base64.RawURLEncoding.DecodeString(secret)
	if err == nil {
		return key, nil
	}
	key, err = base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("pat: decode secret: %w", err)
	}
	return key, nil
}
