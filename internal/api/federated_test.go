package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	federatedID      = "61616161-6161-4161-8161-616161616161"
	federatedClaimID = "62626262-6262-4262-8262-626262626262"
)

type federatedStub struct {
	store.Querier
	federated []store.Federated
	claims    []store.FederatedClaim

	created      *store.CreateFederatedParams
	updated      *store.UpdateFederatedParams
	createdClaim *store.CreateFederatedClaimParams
	updatedClaim *store.UpdateFederatedClaimParams
	deleted      []string
}

func (s *federatedStub) ListFederated(context.Context) ([]store.Federated, error) {
	return s.federated, nil
}

func (s *federatedStub) ListAllFederatedClaims(context.Context) ([]store.FederatedClaim, error) {
	return s.claims, nil
}

func (s *federatedStub) GetFederated(_ context.Context, id string) (store.Federated, error) {
	for _, federated := range s.federated {
		if federated.ID == id {
			return federated, nil
		}
	}
	return store.Federated{}, pgx.ErrNoRows
}

func (s *federatedStub) CreateFederated(_ context.Context, arg store.CreateFederatedParams) (store.Federated, error) {
	s.created = &arg
	return store.Federated{ID: arg.ID, Name: arg.Name, IssuerUrl: arg.IssuerUrl, Audience: arg.Audience}, nil
}

func (s *federatedStub) UpdateFederated(_ context.Context, arg store.UpdateFederatedParams) (store.Federated, error) {
	s.updated = &arg
	return store.Federated{ID: arg.ID, Name: arg.Name, IssuerUrl: arg.IssuerUrl, Audience: arg.Audience}, nil
}

// DeleteFederated answers the foreign key violation upstream reports when a
// credential still has claims.
func (s *federatedStub) DeleteFederated(_ context.Context, id string) error {
	for _, claim := range s.claims {
		if claim.FederatedID == id {
			return &pgconn.PgError{Code: "23503"}
		}
	}
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *federatedStub) ListFederatedClaims(_ context.Context, id string) ([]store.FederatedClaim, error) {
	var out []store.FederatedClaim
	for _, claim := range s.claims {
		if claim.FederatedID == id {
			out = append(out, claim)
		}
	}
	return out, nil
}

func (s *federatedStub) GetFederatedClaim(_ context.Context, id string) (store.FederatedClaim, error) {
	for _, claim := range s.claims {
		if claim.ID == id {
			return claim, nil
		}
	}
	return store.FederatedClaim{}, pgx.ErrNoRows
}

func (s *federatedStub) CreateFederatedClaim(_ context.Context, arg store.CreateFederatedClaimParams) (store.FederatedClaim, error) {
	s.createdClaim = &arg
	return store.FederatedClaim{ID: arg.ID, ClaimKey: arg.ClaimKey, ClaimValue: arg.ClaimValue, FederatedID: arg.FederatedID}, nil
}

func (s *federatedStub) UpdateFederatedClaim(_ context.Context, arg store.UpdateFederatedClaimParams) (store.FederatedClaim, error) {
	s.updatedClaim = &arg
	return store.FederatedClaim{ID: arg.ID, ClaimKey: arg.ClaimKey, ClaimValue: arg.ClaimValue}, nil
}

func (s *federatedStub) DeleteFederatedClaim(_ context.Context, id string) error {
	s.deleted = append(s.deleted, id)
	return nil
}

func federatedTestStore() *federatedStub {
	return &federatedStub{
		federated: []store.Federated{{ID: federatedID, Name: "github", IssuerUrl: strPtr("https://issuer.test"), Audience: strPtr("terrakube")}},
		claims:    []store.FederatedClaim{{ID: federatedClaimID, ClaimKey: "repository", ClaimValue: "acme/infra", FederatedID: federatedID}},
	}
}

func federatedServe(t *testing.T, q *federatedStub, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	handler := New(Deps{Store: q})
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestFederatedRoutes(t *testing.T) {
	q := federatedTestStore()

	rec := federatedServe(t, q, http.MethodPost, "/federated",
		`{"data":{"type":"federated","attributes":{"name":"github","issuerUrl":"https://issuer.test","audience":"terrakube"}}}`)
	if rec.Code != http.StatusCreated || q.created == nil || q.created.Name != "github" {
		t.Fatalf("create federated = %d/%v", rec.Code, q.created)
	}
	if !strings.Contains(rec.Body.String(), `"claims":{"data":[]}`) {
		t.Fatalf("create body = %s", rec.Body.String())
	}

	rec = federatedServe(t, q, http.MethodPost, "/federated", `{"data":{"type":"federated","attributes":{"issuerUrl":"https://issuer.test","audience":"a"}}}`)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), `"code":"NotBlank"`) ||
		!strings.Contains(rec.Body.String(), `"pointer":"/data/attributes/name"`) {
		t.Fatalf("missing name = %d/%s", rec.Code, rec.Body.String())
	}

	rec = federatedServe(t, q, http.MethodGet, "/federated", "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"`+federatedID+`"`) {
		t.Fatalf("list federated = %d/%s", rec.Code, rec.Body.String())
	}

	rec = federatedServe(t, q, http.MethodGet, "/federated/"+federatedID, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"`+federatedClaimID+`"`) {
		t.Fatalf("get federated = %d/%s", rec.Code, rec.Body.String())
	}

	rec = federatedServe(t, q, http.MethodPatch, "/federated/"+federatedID,
		`{"data":{"type":"federated","id":"`+federatedID+`","attributes":{"audience":"terrakube-ci"}}}`)
	if rec.Code != http.StatusNoContent || q.updated == nil || stringValue(q.updated.Audience) != "terrakube-ci" {
		t.Fatalf("patch federated = %d/%v", rec.Code, q.updated)
	}

	rec = federatedServe(t, q, http.MethodDelete, "/federated/"+federatedID, "")
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "conflicts with data that already exists") {
		t.Fatalf("delete federated with claims = %d/%s", rec.Code, rec.Body.String())
	}

	empty := &federatedStub{federated: []store.Federated{{ID: federatedID, Name: "github"}}}
	rec = federatedServe(t, empty, http.MethodDelete, "/federated/"+federatedID, "")
	if rec.Code != http.StatusNoContent || len(empty.deleted) != 1 {
		t.Fatalf("delete federated = %d/%v", rec.Code, empty.deleted)
	}

	rec = federatedServe(t, empty, http.MethodGet, "/federated/00000000-0000-4000-8000-000000000099", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing federated = %d/%s", rec.Code, rec.Body.String())
	}
}

func TestFederatedClaimRoutes(t *testing.T) {
	q := federatedTestStore()
	base := "/federated/" + federatedID + "/claims"

	rec := federatedServe(t, q, http.MethodPost, base, `{"data":{"type":"federated_claim","attributes":{"claimValue":"x"}}}`)
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), `"property":"claimKey"`) {
		t.Fatalf("missing claim key = %d/%s", rec.Code, rec.Body.String())
	}

	rec = federatedServe(t, q, http.MethodPost, base,
		`{"data":{"type":"federated_claim","attributes":{"claimKey":"ref","claimValue":"refs/heads/main"}}}`)
	if rec.Code != http.StatusCreated || q.createdClaim == nil || q.createdClaim.FederatedID != federatedID {
		t.Fatalf("create claim = %d/%v", rec.Code, q.createdClaim)
	}

	rec = federatedServe(t, q, http.MethodGet, base, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"id":"`+federatedClaimID+`"`) {
		t.Fatalf("list claims = %d/%s", rec.Code, rec.Body.String())
	}

	rec = federatedServe(t, q, http.MethodGet, base+"/"+federatedClaimID, "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"claimKey":"repository"`) {
		t.Fatalf("get claim = %d/%s", rec.Code, rec.Body.String())
	}

	rec = federatedServe(t, q, http.MethodPatch, base+"/"+federatedClaimID,
		`{"data":{"type":"federated_claim","id":"`+federatedClaimID+`","attributes":{"claimValue":"acme/*"}}}`)
	if rec.Code != http.StatusNoContent || q.updatedClaim == nil || q.updatedClaim.ClaimValue != "acme/*" {
		t.Fatalf("patch claim = %d/%v", rec.Code, q.updatedClaim)
	}

	rec = federatedServe(t, q, http.MethodDelete, base+"/"+federatedClaimID, "")
	if rec.Code != http.StatusNoContent || len(q.deleted) != 1 || q.deleted[0] != federatedClaimID {
		t.Fatalf("delete claim = %d/%v", rec.Code, q.deleted)
	}

	other := &federatedStub{
		federated: []store.Federated{{ID: federatedID, Name: "github"}},
		claims:    []store.FederatedClaim{{ID: federatedClaimID, ClaimKey: "k", ClaimValue: "v", FederatedID: "63636363-6363-4363-8363-636363636363"}},
	}
	rec = federatedServe(t, other, http.MethodGet, base+"/"+federatedClaimID, "")
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "for claims") {
		t.Fatalf("foreign claim = %d/%s", rec.Code, rec.Body.String())
	}
}

// federatedCaptureStore stands in for the rows the federated captures were
// recorded against. Ids come from the captures, so re-recording needs no edits.
func federatedCaptureStore(captures map[string]upstreamCapture, withClaims bool) *federatedStub {
	created := captureID(captures, "federated-create")
	claimID := captureID(captures, "federated-claim-create")
	q := &federatedStub{
		federated: []store.Federated{{
			ID: created, Name: "capture-fed",
			IssuerUrl: strPtr("https://token.actions.githubusercontent.com"),
			Audience:  strPtr("terrakube"),
		}},
	}
	if !withClaims {
		return q
	}
	q.claims = []store.FederatedClaim{{
		ID: claimID, ClaimKey: "repository", ClaimValue: "acme/infra", FederatedID: created,
	}}
	// The foreign-claim capture addresses a second credential, so both rows
	// exist and the claim answers 404 for the wrong owner rather than a missing
	// credential.
	parts := strings.Split(strings.Trim(captures["federated-claim-foreign"].path, "/"), "/")
	if len(parts) == 6 {
		q.federated = append(q.federated, store.Federated{
			ID: parts[3], Name: "capture-fed-other", IssuerUrl: strPtr("https://issuer.test/fed-a"),
			Audience: strPtr("a"),
		})
		q.claims = append(q.claims, store.FederatedClaim{
			ID: parts[5], ClaimKey: "k", ClaimValue: "v", FederatedID: created,
		})
	}
	return q
}

// TestFederatedTokenIsRefused documents the exchange gap: an external token from
// a registered issuer is not verified here, so it answers 401 like any other
// unknown bearer token. Upstream verifies it through its Dex resolver, which is
// #46; until then a credential only records who would be allowed to sign in.
func TestFederatedTokenIsRefused(t *testing.T) {
	handler := New(Deps{Store: federatedTestStore(), AuthToken: "dev-token"})
	req := httptest.NewRequest(http.MethodGet, "/organization", nil)
	req.Header.Set("Authorization", "Bearer external-issuer-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("federated token status = %d, want 401 until the Dex exchange lands", rec.Code)
	}
}
