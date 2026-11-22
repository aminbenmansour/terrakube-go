package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

const (
	collOrgID    = "11111111-1111-4111-8111-111111111111"
	collWSID     = "22222222-2222-4222-8222-222222222222"
	collOtherOrg = "33333333-3333-4333-8333-333333333333"
	collID       = "44444444-4444-4444-8444-444444444444"
	collItemID   = "55555555-5555-4555-8555-555555555555"
	collRefID    = "66666666-6666-4666-8666-666666666666"
)

// collectionStub answers only the collection methods; the embedded interface
// panics for anything else, which would mean the handler grew a dependency.
type collectionStub struct {
	store.Querier
	org          store.Organization
	workspaces   []store.Workspace
	collections  []store.Collection
	items        []store.Item
	references   []store.Reference
	createdItem  *store.CreateItemParams
	createdRef   *store.CreateReferenceParams
	deletedRefID string
}

func (s *collectionStub) GetOrganization(context.Context, string) (store.Organization, error) {
	return s.org, nil
}

func (s *collectionStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *collectionStub) ListCollectionsByOrganization(context.Context, string) ([]store.Collection, error) {
	return s.collections, nil
}

func (s *collectionStub) GetCollection(_ context.Context, id string) (store.Collection, error) {
	for _, collection := range s.collections {
		if collection.ID == id {
			return collection, nil
		}
	}
	return store.Collection{}, pgx.ErrNoRows
}

func (s *collectionStub) CreateCollection(_ context.Context, arg store.CreateCollectionParams) (store.Collection, error) {
	return store.Collection{ID: arg.ID, Name: arg.Name, Description: arg.Description, Priority: arg.Priority, OrganizationID: arg.OrganizationID}, nil
}

func (s *collectionStub) CreateItem(_ context.Context, arg store.CreateItemParams) (store.Item, error) {
	s.createdItem = &arg
	return store.Item{ID: arg.ID, ItemKey: arg.ItemKey, ItemValue: arg.ItemValue, ItemCategory: arg.ItemCategory, CollectionID: arg.CollectionID}, nil
}

func (s *collectionStub) ListItemsByCollection(context.Context, string) ([]store.Item, error) {
	return s.items, nil
}

func (s *collectionStub) GetItem(_ context.Context, id string) (store.Item, error) {
	for _, item := range s.items {
		if item.ID == id {
			return item, nil
		}
	}
	return store.Item{}, pgx.ErrNoRows
}

func (s *collectionStub) ListReferencesByCollection(context.Context, string) ([]store.Reference, error) {
	return s.references, nil
}

func (s *collectionStub) ListReferencesByWorkspace(_ context.Context, workspaceID string) ([]store.Reference, error) {
	var out []store.Reference
	for _, reference := range s.references {
		if reference.WorkspaceID == workspaceID {
			out = append(out, reference)
		}
	}
	return out, nil
}

func (s *collectionStub) GetReference(_ context.Context, id string) (store.Reference, error) {
	for _, reference := range s.references {
		if reference.ID == id {
			return reference, nil
		}
	}
	return store.Reference{}, pgx.ErrNoRows
}

func (s *collectionStub) CreateReference(_ context.Context, arg store.CreateReferenceParams) (store.Reference, error) {
	s.createdRef = &arg
	return store.Reference{ID: arg.ID, Description: arg.Description, WorkspaceID: arg.WorkspaceID, CollectionID: arg.CollectionID}, nil
}

func (s *collectionStub) DeleteReference(_ context.Context, id string) error {
	s.deletedRefID = id
	return nil
}

func newCollectionAPI(t *testing.T, q store.Querier) http.Handler {
	t.Helper()
	return newTestAPI(t, q)
}

func TestCollectionListCarriesItsItemsAndReferences(t *testing.T) {
	q := &collectionStub{
		collections: []store.Collection{{ID: collID, Name: "shared", Priority: 10, OrganizationID: collOrgID}},
		items:       []store.Item{{ID: collItemID, ItemKey: "KEY", ItemValue: "v", ItemCategory: "ENV", CollectionID: collID}},
		references:  []store.Reference{{ID: collRefID, WorkspaceID: collWSID, CollectionID: collID}},
	}
	rec := httptest.NewRecorder()
	newCollectionAPI(t, q).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+collOrgID+"/collection", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{`"type":"collection"`, `"name":"shared"`, `"item":{"data":[{"type":"item","id":"` + collItemID + `"}]}`, `"reference":{"data":[{"type":"reference","id":"` + collRefID + `"}]}`} {
		if !strings.Contains(body, want) {
			t.Errorf("body does not contain %s: %s", want, body)
		}
	}
}

func TestCreateCollectionRequiresName(t *testing.T) {
	q := &collectionStub{org: store.Organization{ID: collOrgID}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/organization/"+collOrgID+"/collection", strings.NewReader(
		`{"data":{"type":"collection","attributes":{"description":"d"}}}`))
	newCollectionAPI(t, q).ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "conflicts with data that already exists") {
		t.Errorf("body = %s, want upstream's conflict detail", rec.Body.String())
	}
}

func TestCreateCollectionDefaultsPriorityToZero(t *testing.T) {
	q := &collectionStub{org: store.Organization{ID: collOrgID}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/organization/"+collOrgID+"/collection", strings.NewReader(
		`{"data":{"type":"collection","attributes":{"name":"shared"}}}`))
	newCollectionAPI(t, q).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"priority":0`) {
		t.Errorf("body = %s, want upstream's zero priority when none is sent", rec.Body.String())
	}
}

func TestCreateCollectionItemValidatesKeyValueAndCategory(t *testing.T) {
	q := &collectionStub{
		org:         store.Organization{ID: collOrgID},
		collections: []store.Collection{{ID: collID, Name: "shared", OrganizationID: collOrgID}},
	}
	handler := newCollectionAPI(t, q)

	for _, tc := range []struct {
		name   string
		body   string
		status int
		detail string
	}{
		{"missing key", `{"data":{"type":"item","attributes":{"value":"v","category":"ENV"}}}`, http.StatusConflict, "conflicts with data that already exists"},
		{"invalid category", `{"data":{"type":"item","attributes":{"key":"K","value":"v","category":"env"}}}`, http.StatusBadRequest, "Invalid value: env"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/organization/"+collOrgID+"/collection/"+collID+"/item", strings.NewReader(tc.body))
			handler.ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, tc.status, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tc.detail) {
				t.Errorf("body = %s, want %q", rec.Body.String(), tc.detail)
			}
		})
	}
}

func TestCreateCollectionItemStoresTheValue(t *testing.T) {
	q := &collectionStub{
		org:         store.Organization{ID: collOrgID},
		collections: []store.Collection{{ID: collID, Name: "shared", OrganizationID: collOrgID}},
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/organization/"+collOrgID+"/collection/"+collID+"/item", strings.NewReader(
		`{"data":{"type":"item","attributes":{"key":"KEY","value":"v","category":"TERRAFORM","hcl":true}}}`))
	newCollectionAPI(t, q).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if q.createdItem == nil || q.createdItem.ItemCategory != "TERRAFORM" || !q.createdItem.Hcl {
		t.Fatalf("stored %+v, want the category and hcl flag from the request", q.createdItem)
	}
}

func TestCreateReferenceRejectsAnotherOrganization(t *testing.T) {
	q := &collectionStub{
		org:         store.Organization{ID: collOrgID},
		collections: []store.Collection{{ID: collID, Name: "shared", OrganizationID: collOrgID}},
		workspaces:  []store.Workspace{{ID: collWSID, OrganizationID: collOtherOrg}},
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/organization/"+collOrgID+"/collection/"+collID+"/reference", strings.NewReader(
		`{"data":{"type":"reference","relationships":{"workspace":{"data":{"type":"workspace","id":"`+collWSID+`"}}}}}`))
	newCollectionAPI(t, q).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
	if q.createdRef != nil {
		t.Errorf("stored %+v, want no reference", q.createdRef)
	}
}

func TestReferenceCollectionIncludesItems(t *testing.T) {
	q := &collectionStub{
		collections: []store.Collection{{ID: collID, Name: "shared", Priority: 10, OrganizationID: collOrgID}},
		items:       []store.Item{{ID: collItemID, ItemKey: "KEY", ItemValue: "v", ItemCategory: "ENV", CollectionID: collID}},
		references:  []store.Reference{{ID: collRefID, WorkspaceID: collWSID, CollectionID: collID}},
	}
	rec := httptest.NewRecorder()
	newCollectionAPI(t, q).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/reference/"+collRefID+"/collection?include=item", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"included":[{"type":"item","id":"`+collItemID+`"`) {
		t.Errorf("body = %s, want the collection's items in included", rec.Body.String())
	}
}

func TestWorkspaceIncludeReference(t *testing.T) {
	q := &collectionStub{
		workspaces: []store.Workspace{{ID: collWSID, OrganizationID: collOrgID}},
		references: []store.Reference{{ID: collRefID, WorkspaceID: collWSID, CollectionID: collID}},
	}
	rec := httptest.NewRecorder()
	newCollectionAPI(t, q).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/workspace/"+collWSID+"?include=reference", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"reference":{"data":[{"type":"reference","id":"`+collRefID+`"}]}`) {
		t.Errorf("body = %s, want the reference relationship", body)
	}
	if !strings.Contains(body, `"included":[{"type":"reference","id":"`+collRefID+`"`) {
		t.Errorf("body = %s, want the reference in included", body)
	}
}
