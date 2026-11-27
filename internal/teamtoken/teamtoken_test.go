package teamtoken

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
	tokens  map[string]store.TeamToken
	stamped []string
	deleted []string
}

func newFakeStore() *fakeStore { return &fakeStore{tokens: map[string]store.TeamToken{}} }

func (f *fakeStore) CreateTeamToken(_ context.Context, arg store.CreateTeamTokenParams) (store.TeamToken, error) {
	row := store.TeamToken{
		ID: arg.ID, Days: arg.Days, Hours: arg.Hours, Minutes: arg.Minutes,
		Description: arg.Description, GroupName: arg.GroupName,
		CreatedDate: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		UpdatedDate: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
		CreatedBy:   arg.CreatedBy, UpdatedBy: arg.UpdatedBy,
	}
	f.tokens[row.ID] = row
	return row, nil
}

func (f *fakeStore) GetTeamToken(_ context.Context, id string) (store.TeamToken, error) {
	row, ok := f.tokens[id]
	if !ok {
		return store.TeamToken{}, pgx.ErrNoRows
	}
	return row, nil
}

func (f *fakeStore) ListTeamTokensByGroups(_ context.Context, groups []string) ([]store.TeamToken, error) {
	var rows []store.TeamToken
	for _, row := range f.tokens {
		if row.GroupName != nil && contains(groups, *row.GroupName) {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func (f *fakeStore) SoftDeleteTeamToken(_ context.Context, arg store.SoftDeleteTeamTokenParams) error {
	row, ok := f.tokens[arg.ID]
	if !ok {
		return pgx.ErrNoRows
	}
	row.Deleted = true
	row.UpdatedBy = arg.UpdatedBy
	f.tokens[row.ID] = row
	f.deleted = append(f.deleted, row.ID)
	return nil
}

func (f *fakeStore) StampTeamTokenLastUsed(_ context.Context, id string) error {
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

func dexUser(groups ...string) auth.Identity {
	return auth.Identity{
		Issuer: "https://dex.example/dex",
		Email:  "user@example.com",
		Name:   "Example User",
		Groups: groups,
	}
}

func TestCreateSignsTheUpstreamToken(t *testing.T) {
	service, store := newService(t, "")
	token, err := service.Create(context.Background(), dexUser("platform"), "platform", 1, 2, 30, "ci runner")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return service.key, nil })
	if err != nil || !parsed.Valid {
		t.Fatalf("parse token: %v", err)
	}
	if claims["iss"] != Issuer || claims["sub"] != "platform (Team Token)" || claims["name"] != "platform (Token)" {
		t.Errorf("claims = %v", claims)
	}
	if claims["email"] != "user@example.com" || claims["description"] != "ci runner" {
		t.Errorf("claims = %v", claims)
	}
	groups := stringSlice(claims["groups"])
	if len(groups) != 1 || groups[0] != "platform" {
		t.Errorf("groups = %v", claims["groups"])
	}
	id, _ := claims["jti"].(string)
	row, ok := store.tokens[id]
	if !ok || row.GroupName == nil || *row.GroupName != "platform" {
		t.Fatalf("jti %q has no team_token row: %+v", id, store.tokens)
	}
	if row.CreatedBy == nil || *row.CreatedBy != "user@example.com" {
		t.Errorf("created by = %v", row.CreatedBy)
	}
	exp, err := parsed.Claims.GetExpirationTime()
	if err != nil || exp == nil {
		t.Fatalf("exp = %v/%v", exp, err)
	}
	if lifetime := time.Until(exp.Time); lifetime < 26*time.Hour || lifetime > 28*time.Hour {
		t.Errorf("lifetime = %v, want about 26 hours", lifetime)
	}

	forever, err := service.Create(context.Background(), dexUser("platform"), "platform", 0, 0, 0, "no expiry")
	if err != nil {
		t.Fatalf("Create without expiry: %v", err)
	}
	foreverClaims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(forever, foreverClaims, func(*jwt.Token) (any, error) { return service.key, nil }); err != nil {
		t.Fatalf("parse non-expiring token: %v", err)
	}
	if _, ok := foreverClaims["exp"]; ok {
		t.Errorf("non-expiring token has exp = %v", foreverClaims["exp"])
	}
}

func TestCreateRequiresMembershipAndAnInteractiveUser(t *testing.T) {
	service, store := newService(t, "")
	ctx := context.Background()

	if _, err := service.Create(ctx, dexUser("other"), "platform", 1, 0, 0, ""); !errors.Is(err, ErrForbidden) {
		t.Fatalf("foreign group error = %v, want ErrForbidden", err)
	}
	if len(store.tokens) != 0 {
		t.Fatalf("foreign group stored %d rows", len(store.tokens))
	}
	if _, err := service.Create(ctx, dexUser("platform"), "", 1, 0, 0, ""); !errors.Is(err, ErrForbidden) {
		t.Fatalf("empty group error = %v, want ErrForbidden", err)
	}

	tokenUser := dexUser("platform")
	tokenUser.Issuer = Issuer
	if _, err := service.Create(ctx, tokenUser, "platform", 1, 0, 0, ""); !errors.Is(err, ErrNonInteractive) {
		t.Fatalf("token issuer error = %v, want ErrNonInteractive", err)
	}
	if _, err := service.Create(ctx, dexUser("platform"), "platform", -1, 0, 0, ""); !errors.Is(err, ErrInvalidLifetime) {
		t.Fatalf("negative lifetime error = %v, want ErrInvalidLifetime", err)
	}
}

func TestVerifyReturnsTheGroupAndStampsLastUse(t *testing.T) {
	service, fake := newService(t, "")
	ctx := context.Background()

	token, err := service.Create(ctx, dexUser("platform"), "platform", 1, 0, 0, "ci")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	identity, err := service.Verify(ctx, token)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if identity.Issuer != Issuer || identity.Email != "user@example.com" ||
		len(identity.Groups) != 1 || identity.Groups[0] != "platform" {
		t.Fatalf("identity = %+v", identity)
	}
	if len(fake.stamped) != 1 {
		t.Errorf("stamped = %v, want one stamp", fake.stamped)
	}

	// REST and GraphQL both authenticate through this helper, so accepting the
	// token here is what gives it its team's permissions on either surface.
	authenticated, ok := auth.Authenticate(ctx, "Bearer "+token, "", service)
	if !ok || len(authenticated.Groups) != 1 || authenticated.Groups[0] != "platform" {
		t.Fatalf("authenticate = %+v/%v", authenticated, ok)
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return service.key, nil }); err != nil {
		t.Fatalf("parse token: %v", err)
	}
	id, _ := claims["jti"].(string)
	if err := fake.SoftDeleteTeamToken(ctx, store.SoftDeleteTeamTokenParams{ID: id}); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if _, err := service.Verify(ctx, token); !errors.Is(err, ErrInactive) {
		t.Fatalf("revoked token error = %v, want ErrInactive", err)
	}
}

func TestVerifyRejectsAnExpiredToken(t *testing.T) {
	service, store := newService(t, "")
	ctx := context.Background()
	token, err := service.Create(ctx, dexUser("platform"), "platform", 0, 0, 0, "no expiry")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return service.key, nil }); err != nil {
		t.Fatalf("parse token: %v", err)
	}
	claims["exp"] = time.Now().Add(-2 * time.Minute).Unix()
	expired, err := jwt.NewWithClaims(service.alg, claims).SignedString(service.key)
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}
	if _, err := service.Verify(ctx, expired); !errors.Is(err, ErrInactive) {
		t.Fatalf("expired token error = %v, want ErrInactive", err)
	}

	id, _ := claims["jti"].(string)
	if _, ok := store.tokens[id]; !ok {
		t.Fatalf("expiry test lost its metadata row")
	}
}

func TestListAndDeleteAreScopedToTheCallersTeams(t *testing.T) {
	service, store := newService(t, "owner@example.com")
	ctx := context.Background()

	if _, err := service.Create(ctx, dexUser("platform"), "platform", 1, 0, 0, "platform ci"); err != nil {
		t.Fatalf("Create platform: %v", err)
	}
	if _, err := service.Create(ctx, dexUser("database"), "database", 1, 0, 0, "database ci"); err != nil {
		t.Fatalf("Create database: %v", err)
	}

	tokens, err := service.List(ctx, dexUser("platform"))
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tokens) != 1 || tokens[0].Group == nil || *tokens[0].Group != "platform" {
		t.Fatalf("list for platform = %+v", tokens)
	}
	if tokens[0].ID == "" || tokens[0].Description == nil || *tokens[0].Description != "platform ci" {
		t.Fatalf("listed token = %+v", tokens[0])
	}

	if err := service.Delete(ctx, tokens[0].ID, dexUser("database")); !errors.Is(err, ErrForbidden) {
		t.Fatalf("foreign delete error = %v, want ErrForbidden", err)
	}
	if err := service.Delete(ctx, tokens[0].ID, dexUser("platform")); err != nil {
		t.Fatalf("member delete: %v", err)
	}
	if len(store.deleted) != 1 {
		t.Fatalf("deleted = %v", store.deleted)
	}
	if err := service.Delete(ctx, tokens[0].ID, dexUser("platform")); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second delete error = %v, want ErrNotFound", err)
	}

	database, err := service.List(ctx, dexUser("database"))
	if err != nil || len(database) != 1 {
		t.Fatalf("list for database = %+v/%v", database, err)
	}
	if err := service.Delete(ctx, database[0].ID, dexUser("owner@example.com")); err != nil {
		t.Fatalf("instance owner delete: %v", err)
	}
	if err := service.Delete(ctx, "not-a-uuid", dexUser("platform")); !errors.Is(err, ErrNotFound) {
		t.Fatalf("bad id error = %v, want ErrNotFound", err)
	}
}
