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

// globalVarStub keeps the writes so the tests can prove the handler stored what
// the UI sent and preserved what it did not send (sensitive on PATCH).
type globalVarStub struct {
	store.Querier
	orgs    []store.Organization
	vars    []store.Globalvar
	created *store.CreateGlobalVariableParams
	updated *store.UpdateGlobalVariableParams
	deleted string
}

func (s *globalVarStub) GetOrganization(_ context.Context, id string) (store.Organization, error) {
	for _, org := range s.orgs {
		if org.ID == id {
			return org, nil
		}
	}
	return store.Organization{}, pgx.ErrNoRows
}

func (s *globalVarStub) ListGlobalVariablesByOrganization(_ context.Context, orgID string) ([]store.Globalvar, error) {
	var out []store.Globalvar
	for _, variable := range s.vars {
		if variable.OrganizationID == orgID {
			out = append(out, variable)
		}
	}
	return out, nil
}

func (s *globalVarStub) GetGlobalVariable(_ context.Context, id string) (store.Globalvar, error) {
	for _, variable := range s.vars {
		if variable.ID == id {
			return variable, nil
		}
	}
	return store.Globalvar{}, pgx.ErrNoRows
}

func (s *globalVarStub) CreateGlobalVariable(_ context.Context, arg store.CreateGlobalVariableParams) (store.Globalvar, error) {
	s.created = &arg
	return store.Globalvar{
		ID: arg.ID, VariableKey: arg.VariableKey, VariableValue: arg.VariableValue,
		VariableDescription: arg.VariableDescription, VariableCategory: arg.VariableCategory,
		Sensitive: arg.Sensitive, Hcl: arg.Hcl, OrganizationID: arg.OrganizationID,
	}, nil
}

func (s *globalVarStub) UpdateGlobalVariable(_ context.Context, arg store.UpdateGlobalVariableParams) (store.Globalvar, error) {
	s.updated = &arg
	return store.Globalvar{
		ID: arg.ID, VariableKey: arg.VariableKey, VariableValue: arg.VariableValue,
		VariableDescription: arg.VariableDescription, VariableCategory: arg.VariableCategory,
		Sensitive: arg.Sensitive, Hcl: arg.Hcl, OrganizationID: relOrgID,
	}, nil
}

func (s *globalVarStub) DeleteGlobalVariable(_ context.Context, id string) error {
	s.deleted = id
	return nil
}

const globalVarID = "88888888-8888-4888-8888-888888888888"

func newGlobalVarStub() *globalVarStub {
	return &globalVarStub{
		orgs: []store.Organization{{ID: relOrgID, Name: "simple-governance"}},
		vars: []store.Globalvar{{
			ID: globalVarID, VariableKey: "REGION", VariableValue: "eu", VariableDescription: "d",
			VariableCategory: "ENV", Sensitive: true, OrganizationID: relOrgID,
		}},
	}
}

func serveGlobalVar(t *testing.T, s *globalVarStub, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	newTestAPI(t, s).ServeHTTP(rec, req)
	return rec
}

func TestCreateGlobalVariableHidesASensitiveValue(t *testing.T) {
	s := newGlobalVarStub()
	rec := serveGlobalVar(t, s, http.MethodPost, "/organization/"+relOrgID+"/globalvar",
		`{"data":{"type":"globalvar","attributes":{"key":"TOKEN","value":"secret","description":"d","sensitive":true,"hcl":false,"category":"ENV"}}}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	if s.created == nil || s.created.VariableValue != "secret" || s.created.Sensitive != true {
		t.Fatalf("created = %+v, want the value stored", s.created)
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, present := doc.Data.Attributes["value"]; present {
		t.Errorf("attributes = %v, want no value for a sensitive variable", doc.Data.Attributes)
	}
	if got := string(doc.Data.Relationships["organization"].Data); !strings.Contains(got, relOrgID) {
		t.Errorf("relationships.organization = %s, want %s", got, relOrgID)
	}
}

func TestCreateGlobalVariableRejectsALowercaseCategory(t *testing.T) {
	rec := serveGlobalVar(t, newGlobalVarStub(), http.MethodPost, "/organization/"+relOrgID+"/globalvar",
		`{"data":{"type":"globalvar","attributes":{"key":"REGION","value":"eu","description":"d","category":"env"}}}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Errors) == 0 || doc.Errors[0].Detail != "Invalid value: env" {
		t.Errorf("errors = %v, want Invalid value: env", doc.Errors)
	}
}

// Upstream answers a missing required column with its generic conflict, not a
// 422 that names the attribute.
func TestCreateGlobalVariableWithoutAKeyConflicts(t *testing.T) {
	rec := serveGlobalVar(t, newGlobalVarStub(), http.MethodPost, "/organization/"+relOrgID+"/globalvar",
		`{"data":{"type":"globalvar","attributes":{"value":"eu","description":"d","category":"ENV"}}}`)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body = %s", rec.Code, rec.Body.String())
	}
}

func TestPatchGlobalVariableKeepsSensitiveAndAnswersNoContent(t *testing.T) {
	s := newGlobalVarStub()
	rec := serveGlobalVar(t, s, http.MethodPatch, "/organization/"+relOrgID+"/globalvar/"+globalVarID,
		`{"data":{"type":"globalvar","id":"`+globalVarID+`","attributes":{"key":"REGION","value":"us","description":"after","hcl":true,"category":"TERRAFORM"}}}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
	if s.updated == nil || !s.updated.Sensitive {
		t.Errorf("updated = %+v, want sensitive preserved", s.updated)
	}
	if s.updated.VariableValue != "us" || s.updated.VariableCategory != "TERRAFORM" || !s.updated.Hcl {
		t.Errorf("updated = %+v, want the sent values", s.updated)
	}
}

func TestGlobalVariableFromAnotherOrganizationIsNotFound(t *testing.T) {
	rec := serveGlobalVar(t, newGlobalVarStub(), http.MethodGet,
		"/organization/99999999-9999-4999-8999-999999999999/globalvar/"+globalVarID, "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestListGlobalVariablesHidesSensitiveValues(t *testing.T) {
	rec := serveGlobalVar(t, newGlobalVarStub(), http.MethodGet, "/organization/"+relOrgID+"/globalvar", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		Data []struct {
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 1 {
		t.Fatalf("data = %v, want one variable", doc.Data)
	}
	if _, present := doc.Data[0].Attributes["value"]; present {
		t.Errorf("attributes = %v, want no value", doc.Data[0].Attributes)
	}
}

func TestDeleteGlobalVariableAnswersNoContent(t *testing.T) {
	s := newGlobalVarStub()
	rec := serveGlobalVar(t, s, http.MethodDelete, "/organization/"+relOrgID+"/globalvar/"+globalVarID, "")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if s.deleted != globalVarID {
		t.Errorf("deleted = %q, want %s", s.deleted, globalVarID)
	}
}
