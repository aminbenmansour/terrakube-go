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

const (
	templateOrgID     = "31313131-3131-4131-8131-313131313131"
	templateOtherID   = "32323232-3232-4232-8232-323232323232"
	templateID        = "33333333-3333-4333-8333-333333333333"
	templateForeignID = "34343434-3434-4434-8434-343434343434"
)

// templateStub answers only the template methods; the embedded interface panics
// for anything else, which would mean the handler grew a dependency.
type templateStub struct {
	store.Querier
	org       store.Organization
	templates []store.Template
	created   *store.CreateTemplateParams
	updated   *store.UpdateTemplateParams
	deleted   string
}

func (s *templateStub) GetOrganization(_ context.Context, id string) (store.Organization, error) {
	if id != s.org.ID {
		return store.Organization{}, pgx.ErrNoRows
	}
	return s.org, nil
}

func (s *templateStub) ListTemplatesByOrganization(context.Context, string) ([]store.Template, error) {
	return s.templates, nil
}

func (s *templateStub) GetTemplate(_ context.Context, id string) (store.Template, error) {
	for _, template := range s.templates {
		if template.ID == id {
			return template, nil
		}
	}
	return store.Template{}, pgx.ErrNoRows
}

func (s *templateStub) CreateTemplate(_ context.Context, arg store.CreateTemplateParams) (store.Template, error) {
	s.created = &arg
	return store.Template{
		ID: arg.ID, Name: arg.Name, Description: arg.Description, Version: arg.Version,
		Tcl: arg.Tcl, OrganizationID: arg.OrganizationID,
	}, nil
}

func (s *templateStub) UpdateTemplate(_ context.Context, arg store.UpdateTemplateParams) (store.Template, error) {
	s.updated = &arg
	return store.Template{ID: arg.ID, Name: arg.Name, Description: arg.Description}, nil
}

func (s *templateStub) DeleteTemplate(_ context.Context, id string) error {
	s.deleted = id
	return nil
}

func newTemplateStub() *templateStub {
	return &templateStub{
		org: store.Organization{ID: templateOrgID},
		templates: []store.Template{{
			ID: templateID, Name: strPtr("Plan"), Description: strPtr("plan"),
			Version: strPtr("1.0.0"), Tcl: strPtr("ZmxvdzoK"), OrganizationID: templateOrgID,
		}},
	}
}

func templateTestAPI(q store.Querier) http.Handler {
	return New(Deps{Store: q})
}

func TestTemplateJSONAPI(t *testing.T) {
	q := newTemplateStub()
	h := templateTestAPI(q)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+templateOrgID+"/template", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d; body = %s", rec.Code, rec.Body.String())
	}
	var list struct {
		Data []struct {
			Type       string         `json:"type"`
			Attributes map[string]any `json:"attributes"`
			Rel        map[string]any `json:"relationships"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 1 || list.Data[0].Type != "template" {
		t.Fatalf("list = %+v", list.Data)
	}
	if list.Data[0].Attributes["name"] != "Plan" || list.Data[0].Attributes["version"] != "1.0.0" {
		t.Errorf("attributes = %+v", list.Data[0].Attributes)
	}
	orgRel, _ := list.Data[0].Rel["organization"].(map[string]any)
	if orgRel == nil {
		t.Errorf("organization relationship missing: %+v", list.Data[0].Rel)
	}

	// Upstream accepts a template with every attribute missing.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+templateOrgID+"/template",
		strings.NewReader(`{"data":{"type":"template","attributes":{}}}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("empty create = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.created == nil || q.created.Name != nil || q.created.OrganizationID != templateOrgID {
		t.Errorf("empty create params = %+v", q.created)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/organization/"+templateOrgID+"/template",
		strings.NewReader(`{"data":{"type":"template","attributes":{"name":"Apply","description":"apply","version":"2.0.0","tcl":"ZmxvdzoK"}}}`)))
	if rec.Code != http.StatusCreated || q.created == nil || q.created.Name == nil || *q.created.Name != "Apply" {
		t.Fatalf("create = %d; params = %+v", rec.Code, q.created)
	}

	// A patch that leaves an attribute out keeps it; null clears it.
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/organization/"+templateOrgID+"/template/"+templateID,
		strings.NewReader(`{"data":{"type":"template","id":"`+templateID+`","attributes":{"description":"renamed","name":null}}}`)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("patch = %d; body = %s", rec.Code, rec.Body.String())
	}
	if q.updated == nil || q.updated.Description == nil || *q.updated.Description != "renamed" ||
		q.updated.Name != nil || q.updated.Version == nil || *q.updated.Version != "1.0.0" {
		t.Errorf("patch params = %+v", q.updated)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/organization/"+templateOrgID+"/template/"+templateID, nil))
	if rec.Code != http.StatusNoContent || q.deleted != templateID {
		t.Fatalf("delete = %d, deleted %q", rec.Code, q.deleted)
	}

	// A template id from another organization is not readable through this one.
	other := newTemplateStub()
	other.org = store.Organization{ID: templateOtherID}
	other.templates = append(other.templates, store.Template{ID: templateForeignID, OrganizationID: templateOtherID})
	rec = httptest.NewRecorder()
	templateTestAPI(other).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+templateOtherID+"/template/"+templateID, nil))
	if rec.Code != http.StatusNotFound {
		t.Errorf("cross-org get = %d, want 404", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization/"+templateOrgID+"/template/not-a-uuid", nil))
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "Invalid value: not-a-uuid") {
		t.Errorf("bad id = %d; body = %s", rec.Code, rec.Body.String())
	}
}
