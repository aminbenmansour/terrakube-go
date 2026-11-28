package store

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
)

// TestFederatedQueries covers the credential and claim round trips the Settings
// > Federated credentials screen needs, including the foreign key that refuses
// to delete a credential which still has claims.
func TestFederatedQueries(t *testing.T) {
	q := New(dbtest.Pool(t))
	ctx := t.Context()

	federated, err := q.CreateFederated(ctx, CreateFederatedParams{
		ID: uuid.NewString(), Name: "capture-fed", IssuerUrl: strPtr("https://issuer.test"),
		Audience: strPtr("terrakube"), CreatedBy: strPtr("admin@example.com"),
	})
	if err != nil {
		t.Fatalf("create federated: %v", err)
	}
	list, err := q.ListFederated(ctx)
	if err != nil || len(list) != 1 || list[0].ID != federated.ID {
		t.Fatalf("list federated = %+v/%v", list, err)
	}

	claim, err := q.CreateFederatedClaim(ctx, CreateFederatedClaimParams{
		ID: uuid.NewString(), ClaimKey: "repository", ClaimValue: "acme/infra",
		FederatedID: federated.ID,
	})
	if err != nil {
		t.Fatalf("create claim: %v", err)
	}
	claims, err := q.ListFederatedClaims(ctx, federated.ID)
	if err != nil || len(claims) != 1 || claims[0].ID != claim.ID {
		t.Fatalf("list claims = %+v/%v", claims, err)
	}
	if updated, err := q.UpdateFederatedClaim(ctx, UpdateFederatedClaimParams{
		ID: claim.ID, ClaimKey: "repository", ClaimValue: "acme/*",
	}); err != nil || updated.ClaimValue != "acme/*" {
		t.Fatalf("update claim = %+v/%v", updated, err)
	}

	var pgErr *pgconn.PgError
	if err := q.DeleteFederated(ctx, federated.ID); !errors.As(err, &pgErr) || pgErr.Code != "23503" {
		t.Fatalf("delete federated with claims = %v, want foreign key violation", err)
	}

	if err := q.DeleteFederatedClaim(ctx, claim.ID); err != nil {
		t.Fatalf("delete claim: %v", err)
	}
	if updated, err := q.UpdateFederated(ctx, UpdateFederatedParams{
		ID: federated.ID, Name: "capture-fed", IssuerUrl: strPtr("https://issuer.test"),
		Audience: strPtr("terrakube-ci"), UpdatedBy: strPtr("admin@example.com"),
	}); err != nil || *updated.Audience != "terrakube-ci" {
		t.Fatalf("update federated = %+v/%v", updated, err)
	}

	if err := q.DeleteFederated(ctx, federated.ID); err != nil {
		t.Fatalf("delete federated: %v", err)
	}
	if _, err := q.GetFederated(ctx, federated.ID); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("get deleted federated = %v, want no rows", err)
	}
}
