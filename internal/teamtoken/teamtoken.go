// Package teamtoken issues and validates the team-scoped access tokens the
// upstream UI lists under Settings > Teams. The token is an HMAC-signed JWT
// carrying one group; only its metadata row is stored, and revocation is a soft
// delete the verifier checks on every request.
package teamtoken

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
)

var (
	ErrInactive        = errors.New("teamtoken: token is revoked, unknown or expired")
	ErrNotFound        = errors.New("teamtoken: team token not found")
	ErrForbidden       = errors.New("teamtoken: caller is not a member of this team")
	ErrNonInteractive  = errors.New("teamtoken: only an interactive user can create tokens")
	ErrInvalidLifetime = errors.New("teamtoken: token lifetime cannot be negative")
)

// Store is the subset of the generated queries the service needs.
type Store interface {
	CreateTeamToken(ctx context.Context, arg store.CreateTeamTokenParams) (store.TeamToken, error)
	GetTeamToken(ctx context.Context, id string) (store.TeamToken, error)
	ListTeamTokensByGroups(ctx context.Context, groups []string) ([]store.TeamToken, error)
	SoftDeleteTeamToken(ctx context.Context, arg store.SoftDeleteTeamTokenParams) error
	StampTeamTokenLastUsed(ctx context.Context, id string) error
}

// Service mints, verifies and revokes team tokens.
type Service struct {
	store Store
	key   []byte
	alg   jwt.SigningMethod
	owner string
}

// New builds the service. The secret is the same base64url PAT secret upstream
// reuses for team tokens and must decode to at least 32 bytes.
func New(store Store, secret, owner string) (*Service, error) {
	key, err := decodeSecret(secret)
	if err != nil {
		return nil, err
	}
	if len(key) < 32 {
		return nil, errors.New("teamtoken: secret must decode to at least 32 bytes")
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
	ID          string     `json:"id"`
	Days        *int32     `json:"days"`
	Hours       *int32     `json:"hours"`
	Minutes     *int32     `json:"minutes"`
	Group       *string    `json:"group"`
	Description *string    `json:"description"`
	Deleted     bool       `json:"deleted"`
	LastUsedAt  *time.Time `json:"lastUsedAt"`
	CreatedDate *time.Time `json:"createdDate"`
	UpdatedDate *time.Time `json:"updatedDate"`
	CreatedBy   *string    `json:"createdBy"`
	UpdatedBy   *string    `json:"updatedBy"`
}

// Create stores the metadata row and returns the signed token. The caller must
// belong to the group, so a user cannot mint a token that escalates their
// access.
func (s *Service) Create(ctx context.Context, user auth.Identity, group string, days, hours, minutes int, description string) (string, error) {
	if user.Issuer == Issuer {
		return "", ErrNonInteractive
	}
	group = strings.TrimSpace(group)
	if group == "" || !contains(user.Groups, group) {
		return "", ErrForbidden
	}
	if days < 0 || hours < 0 || minutes < 0 {
		return "", ErrInvalidLifetime
	}

	id := uuid.NewString()
	daysValue, hoursValue, minutesValue := int32(days), int32(hours), int32(minutes)
	row, err := s.store.CreateTeamToken(ctx, store.CreateTeamTokenParams{
		ID:          id,
		Days:        &daysValue,
		Hours:       &hoursValue,
		Minutes:     &minutesValue,
		Description: &description,
		GroupName:   &group,
		CreatedBy:   stringPtr(user.Email),
		UpdatedBy:   stringPtr(user.Email),
	})
	if err != nil {
		return "", fmt.Errorf("teamtoken: create metadata: %w", err)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":            Issuer,
		"sub":            group + " (Team Token)",
		"aud":            []string{Issuer},
		"jti":            row.ID,
		"email":          user.Email,
		"description":    description,
		"email_verified": true,
		"name":           group + " (Token)",
		"groups":         []string{group},
		"iat":            now.Unix(),
	}
	if lifetime := time.Duration(days)*24*time.Hour + time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute; lifetime > 0 {
		claims["exp"] = now.Add(lifetime).Unix()
	}

	token, err := jwt.NewWithClaims(s.alg, claims).SignedString(s.key)
	if err != nil {
		// Do not leave a row behind that no client can ever use.
		_ = s.store.SoftDeleteTeamToken(ctx, store.SoftDeleteTeamTokenParams{ID: row.ID, UpdatedBy: stringPtr(user.Email)})
		return "", fmt.Errorf("teamtoken: sign token: %w", err)
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
	row, err := s.store.GetTeamToken(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && row.Deleted) {
		return nil, ErrInactive
	}
	if err != nil {
		return nil, fmt.Errorf("teamtoken: load metadata: %w", err)
	}
	// A failed stamp must not fail the request; the token stays valid.
	_ = s.store.StampTeamTokenLastUsed(ctx, id)

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

// List returns the active tokens of the caller's teams.
func (s *Service) List(ctx context.Context, user auth.Identity) ([]Token, error) {
	rows, err := s.store.ListTeamTokensByGroups(ctx, user.Groups)
	if err != nil {
		return nil, fmt.Errorf("teamtoken: list: %w", err)
	}
	tokens := make([]Token, 0, len(rows))
	for _, row := range rows {
		if row.Deleted {
			continue
		}
		tokens = append(tokens, Token{
			ID:          row.ID,
			Days:        row.Days,
			Hours:       row.Hours,
			Minutes:     row.Minutes,
			Group:       row.GroupName,
			Description: row.Description,
			Deleted:     row.Deleted,
			LastUsedAt:  timePtr(row.LastUsedAt),
			CreatedDate: timePtr(row.CreatedDate),
			UpdatedDate: timePtr(row.UpdatedDate),
			CreatedBy:   row.CreatedBy,
			UpdatedBy:   row.UpdatedBy,
		})
	}
	return tokens, nil
}

// Delete revokes a token. The caller must belong to the token's team, or be the
// configured instance owner (by email or group, as upstream does).
func (s *Service) Delete(ctx context.Context, id string, user auth.Identity) error {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ErrNotFound
	}
	row, err := s.store.GetTeamToken(ctx, parsed.String())
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && row.Deleted) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("teamtoken: load metadata: %w", err)
	}
	if !s.canDelete(user, row) {
		return ErrForbidden
	}
	if err := s.store.SoftDeleteTeamToken(ctx, store.SoftDeleteTeamTokenParams{ID: row.ID, UpdatedBy: stringPtr(user.Email)}); err != nil {
		return fmt.Errorf("teamtoken: revoke: %w", err)
	}
	return nil
}

func (s *Service) canDelete(user auth.Identity, row store.TeamToken) bool {
	if row.GroupName != nil && contains(user.Groups, *row.GroupName) {
		return true
	}
	if s.owner == "" {
		return false
	}
	if strings.EqualFold(user.Email, s.owner) {
		return true
	}
	return contains(user.Groups, s.owner)
}

func decodeSecret(secret string) ([]byte, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("teamtoken: secret is not configured")
	}
	key, err := base64.RawURLEncoding.DecodeString(secret)
	if err == nil {
		return key, nil
	}
	key, err = base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("teamtoken: decode secret: %w", err)
	}
	return key, nil
}

func stringSlice(value any) []string {
	values, ok := value.([]any)
	if !ok {
		if strings, ok := value.([]string); ok {
			return strings
		}
		return nil
	}
	out := make([]string, 0, len(values))
	for _, entry := range values {
		if text, ok := entry.(string); ok {
			out = append(out, text)
		}
	}
	return out
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
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
	time := value.Time
	return &time
}
