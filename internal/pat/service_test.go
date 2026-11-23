package pat

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

type fakeStore struct {
	pats    map[string]store.Pat
	stamped []string
	deleted []string
}

func newFakeStore() *fakeStore { return &fakeStore{pats: map[string]store.Pat{}} }

func (f *fakeStore) CreatePat(_ context.Context, arg store.CreatePatParams) (store.Pat, error) {
	row := store.Pat{
		ID:          arg.ID,
		Days:        arg.Days,
		Description: arg.Description,
		Source:      arg.Source,
		Deleted:     false,
		CreatedBy:   arg.CreatedBy,
		UpdatedBy:   arg.UpdatedBy,
		CreatedDate: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		UpdatedDate: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
	}
	f.pats[row.ID] = row
	return row, nil
}

func (f *fakeStore) GetPat(_ context.Context, id string) (store.Pat, error) {
	row, ok := f.pats[id]
	if !ok {
		return store.Pat{}, pgx.ErrNoRows
	}
	return row, nil
}

func (f *fakeStore) ListPatsByCreator(_ context.Context, createdBy *string) ([]store.Pat, error) {
	var rows []store.Pat
	for _, row := range f.pats {
		if !row.Deleted && row.CreatedBy != nil && createdBy != nil && *row.CreatedBy == *createdBy {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (f *fakeStore) SoftDeletePat(_ context.Context, arg store.SoftDeletePatParams) error {
	row, ok := f.pats[arg.ID]
	if !ok {
		return pgx.ErrNoRows
	}
	row.Deleted = true
	row.UpdatedBy = arg.UpdatedBy
	f.pats[row.ID] = row
	f.deleted = append(f.deleted, row.ID)
	return nil
}

func (f *fakeStore) StampPatLastUsed(_ context.Context, id string) error {
	f.stamped = append(f.stamped, id)
	return nil
}

func testSecret() string {
	return base64.RawURLEncoding.EncodeToString([]byte(strings.Repeat("k", 32)))
}

func newService(t *testing.T, owner string) (*Service, *fakeStore) {
	t.Helper()
	store := newFakeStore()
	service, err := New(store, testSecret(), owner)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return service, store
}

func dexUser() auth.Identity {
	return auth.Identity{
		Issuer: "https://dex.example/dex",
		Email:  "user@example.com",
		Name:   "Example User",
		Groups: []string{"terrakube-admins"},
	}
}

func TestCreateSignsTheUpstreamToken(t *testing.T) {
	service, store := newService(t, "")
	token, err := service.Create(context.Background(), dexUser(), 30, "ci")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return service.key, nil })
	if err != nil || !parsed.Valid {
		t.Fatalf("parse token: %v", err)
	}
	if claims["iss"] != Issuer || claims["email"] != "user@example.com" || claims["name"] != "Example User (Token)" {
		t.Errorf("claims = %v", claims)
	}
	if claims["sub"] != "Example User (Token)" {
		t.Errorf("sub = %v", claims["sub"])
	}
	if claims["jti"] == nil || store.pats[claims["jti"].(string)].ID == "" {
		t.Errorf("jti %v has no pat row", claims["jti"])
	}
	if claims["exp"] == nil {
		t.Error("30-day token has no exp")
	}
	audience, _ := claims.GetAudience()
	if len(audience) != 1 || audience[0] != Issuer {
		t.Errorf("aud = %v, want [%s]", audience, Issuer)
	}
}

func TestCreateWithoutExpiryHasNoExp(t *testing.T) {
	service, _ := newService(t, "")
	token, err := service.Create(context.Background(), dexUser(), 0, "forever")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return service.key, nil }); err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, ok := claims["exp"]; ok {
		t.Error("days=0 token has an exp claim")
	}
}

func TestCreateRejectsNonInteractiveUsersAndBadLifetimes(t *testing.T) {
	service, _ := newService(t, "")

	pat := dexUser()
	pat.Issuer = Issuer
	if _, err := service.Create(context.Background(), pat, 1, ""); !errors.Is(err, ErrNonInteractive) {
		t.Errorf("Create(PAT principal) = %v, want ErrNonInteractive", err)
	}
	if _, err := service.Create(context.Background(), dexUser(), -1, ""); !errors.Is(err, ErrInvalidLifetime) {
		t.Errorf("Create(-1 days) = %v, want ErrInvalidLifetime", err)
	}
}

func TestVerifyAcceptsAnActiveTokenAndStampsUse(t *testing.T) {
	service, store := newService(t, "")
	token, err := service.Create(context.Background(), dexUser(), 1, "ci")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	identity, err := service.Verify(context.Background(), token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if identity.Email != "user@example.com" || identity.Name != "Example User (Token)" || len(identity.Groups) != 1 {
		t.Errorf("identity = %+v", identity)
	}
	if len(store.stamped) != 1 {
		t.Errorf("stamped = %v, want one last-used update", store.stamped)
	}
}

func TestVerifyRejectsRevokedExpiredAndForeignTokens(t *testing.T) {
	service, store := newService(t, "")
	token, err := service.Create(context.Background(), dexUser(), 1, "ci")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return service.key, nil })
	patID := claims["jti"].(string)

	row := store.pats[patID]
	row.Deleted = true
	store.pats[patID] = row
	if _, err := service.Verify(context.Background(), token); !errors.Is(err, ErrInactive) {
		t.Errorf("Verify(revoked) = %v, want ErrInactive", err)
	}
	row.Deleted = false
	store.pats[patID] = row

	expired := jwt.NewWithClaims(service.alg, jwt.MapClaims{
		"iss": Issuer, "jti": patID, "exp": time.Now().Add(-time.Hour).Unix(),
	})
	raw, err := expired.SignedString(service.key)
	if err != nil {
		t.Fatalf("sign expired: %v", err)
	}
	if _, err := service.Verify(context.Background(), raw); err == nil {
		t.Error("Verify(expired) = nil, want an error")
	}

	other, _ := New(store, base64.RawURLEncoding.EncodeToString([]byte(strings.Repeat("z", 32))), "")
	if _, err := other.Verify(context.Background(), token); err == nil {
		t.Error("Verify(wrong secret) = nil, want an error")
	}
}

func TestListReturnsMetadataOnlyForTheCreator(t *testing.T) {
	service, store := newService(t, "")
	if _, err := service.Create(context.Background(), dexUser(), 1, "one"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, err := service.Create(context.Background(), auth.Identity{Issuer: "dex", Email: "other@example.com", Name: "Other"}, 2, "two"); err != nil {
		t.Fatalf("Create: %v", err)
	}

	tokens, err := service.List(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tokens) != 1 || tokens[0].Description == nil || *tokens[0].Description != "one" {
		t.Fatalf("tokens = %+v, want only the caller's token", tokens)
	}
	if tokens[0].CreatedDate == nil || tokens[0].CreatedDate.IsZero() {
		t.Error("CreatedDate is missing")
	}
	_ = store
}

func TestDeleteEnforcesOwnership(t *testing.T) {
	service, _ := newService(t, "")
	token, err := service.Create(context.Background(), dexUser(), 1, "ci")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return service.key, nil })
	id := claims["jti"].(string)

	other := auth.Identity{Issuer: "dex", Email: "other@example.com"}
	if err := service.Delete(context.Background(), id, other); !errors.Is(err, ErrForbidden) {
		t.Errorf("Delete(other) = %v, want ErrForbidden", err)
	}
	if err := service.Delete(context.Background(), id, dexUser()); err != nil {
		t.Errorf("Delete(owner) = %v", err)
	}
	if _, err := service.Verify(context.Background(), token); !errors.Is(err, ErrInactive) {
		t.Errorf("Verify(revoked) = %v, want ErrInactive", err)
	}
	if err := service.Delete(context.Background(), "not-a-uuid", dexUser()); !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete(bad id) = %v, want ErrNotFound", err)
	}
}

func TestInstanceOwnerCanRevokeAnyToken(t *testing.T) {
	service, _ := newService(t, "admin@example.com")
	token, err := service.Create(context.Background(), dexUser(), 1, "ci")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return service.key, nil })

	admin := auth.Identity{Issuer: "dex", Email: "admin@example.com"}
	if err := service.Delete(context.Background(), claims["jti"].(string), admin); err != nil {
		t.Errorf("Delete(instance owner) = %v", err)
	}
}
