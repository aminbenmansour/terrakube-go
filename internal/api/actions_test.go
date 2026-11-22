package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// actionStub keeps writes so the tests can prove PATCH preserved the attributes
// it did not receive, and that the slug id is not mistaken for a UUID.
type actionStub struct {
	store.Querier
	actions []store.Action
	created *store.CreateActionParams
	updated *store.UpdateActionParams
	deleted string
}

func (s *actionStub) ListActions(context.Context) ([]store.Action, error) {
	return s.actions, nil
}

func (s *actionStub) GetAction(_ context.Context, id string) (store.Action, error) {
	for _, action := range s.actions {
		if action.ID == id {
			return action, nil
		}
	}
	return store.Action{}, pgx.ErrNoRows
}

func (s *actionStub) CreateAction(_ context.Context, arg store.CreateActionParams) (store.Action, error) {
	s.created = &arg
	return store.Action{
		ID: arg.ID, Name: arg.Name, Label: arg.Label, Category: arg.Category,
		Description: arg.Description, Type: arg.Type, Action: arg.Action,
		DisplayCriteria: arg.DisplayCriteria, Version: arg.Version, Active: arg.Active,
	}, nil
}

func (s *actionStub) UpdateAction(_ context.Context, arg store.UpdateActionParams) (store.Action, error) {
	s.updated = &arg
	return store.Action{ID: arg.ID}, nil
}

func (s *actionStub) DeleteAction(_ context.Context, id string) error {
	s.deleted = id
	return nil
}

func newActionStub() *actionStub {
	return &actionStub{actions: []store.Action{
		{ID: "terrakube.open-documentation", Name: "Open Documentation", Type: "Workspace/ResourceDrawer/Action", Active: true},
		{ID: "terrakube.resource-details", Name: "Resource Details", Type: "Workspace/ResourceDrawer/Tab", Active: true},
		{ID: "terrakube.open-ai", Name: "Open AI", Type: "Workspace/Action", Active: false},
	}}
}

func TestListActionsFiltersTheDrawerExpression(t *testing.T) {
	stub := newActionStub()
	handler := newTestAPI(t, stub)

	for _, tc := range []struct {
		name  string
		query string
		want  int
	}{
		{"all", "", 3},
		{"drawer", "?filter[action]=active==true%3Btype=in=('Workspace/ResourceDrawer/Action','Workspace/ResourceDrawer/Tab')", 2},
		{"workspace actions", "?filter[action]=active==true%3Btype=in=('Workspace/Action')", 0},
		{"inactive", "?filter[action]=active==false", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/action"+tc.query, nil))
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
			}
			var doc struct {
				Data []json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if len(doc.Data) != tc.want {
				t.Fatalf("data = %d resources, want %d", len(doc.Data), tc.want)
			}
		})
	}
}

func TestGetActionAcceptsSlugID(t *testing.T) {
	handler := newTestAPI(t, newActionStub())

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/action/terrakube.open-ai", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"id":"terrakube.open-ai"`) {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestCreateActionRequiresAnIDAndAttributes(t *testing.T) {
	stub := newActionStub()
	handler := newTestAPI(t, stub)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/action",
		strings.NewReader(`{"data":{"type":"action","attributes":{"name":"x"}}}`)))
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "No id provided") {
		t.Fatalf("missing id: status = %d, body = %s", rec.Code, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/action",
		strings.NewReader(`{"data":{"type":"action","id":"capture.temp","attributes":{"name":"Top"}}}`)))
	if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "label is required") {
		t.Fatalf("missing label: status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if stub.created != nil {
		t.Fatal("invalid payload reached the store")
	}
}

func TestCreateActionStoresTheCatalogRow(t *testing.T) {
	stub := newActionStub()
	handler := newTestAPI(t, stub)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/action", strings.NewReader(
		`{"data":{"type":"action","id":"capture.temp","attributes":{"name":"Top","label":"Top","category":"General","type":"Workspace/Action","action":"dGVzdA==","displayCriteria":"[]","version":"1.0.0"}}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if stub.created == nil || !stub.created.Active {
		t.Fatalf("created = %+v, want active by default", stub.created)
	}
}

func TestUpdateActionKeepsOmittedAttributes(t *testing.T) {
	stub := newActionStub()
	handler := newTestAPI(t, stub)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/action/terrakube.open-ai",
		strings.NewReader(`{"data":{"type":"action","id":"terrakube.open-ai","attributes":{"active":true}}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if stub.updated == nil || stub.updated.Name != "Open AI" || !stub.updated.Active {
		t.Fatalf("updated = %+v, want the name preserved and active set", stub.updated)
	}
}

func TestDeleteAction(t *testing.T) {
	stub := newActionStub()
	handler := newTestAPI(t, stub)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/action/terrakube.open-ai", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	if stub.deleted != "terrakube.open-ai" {
		t.Fatalf("deleted = %q", stub.deleted)
	}
}
