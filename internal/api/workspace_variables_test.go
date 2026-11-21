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

// nestedStub is a tiny in-memory store for the routes the UI reaches through
// /organization/:org/workspace/:ws. It exists because those routes must prove a
// workspace and a variable belong to the organization in the path; a stub that
// answers regardless would hide exactly the bug worth testing.
type nestedStub struct {
	store.Querier
	orgs       []store.Organization
	workspaces []store.Workspace
	variables  []store.Variable
}

func (s *nestedStub) GetOrganization(_ context.Context, id string) (store.Organization, error) {
	for _, org := range s.orgs {
		if org.ID == id {
			return org, nil
		}
	}
	return store.Organization{}, pgx.ErrNoRows
}

func (s *nestedStub) GetWorkspace(_ context.Context, id string) (store.Workspace, error) {
	for _, ws := range s.workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return store.Workspace{}, pgx.ErrNoRows
}

func (s *nestedStub) ListVariablesByWorkspace(_ context.Context, workspaceID string) ([]store.Variable, error) {
	var out []store.Variable
	for _, v := range s.variables {
		if v.WorkspaceID == workspaceID {
			out = append(out, v)
		}
	}
	return out, nil
}

func (s *nestedStub) GetVariable(_ context.Context, id string) (store.Variable, error) {
	for _, v := range s.variables {
		if v.ID == id {
			return v, nil
		}
	}
	return store.Variable{}, pgx.ErrNoRows
}

func (s *nestedStub) CreateVariable(_ context.Context, arg store.CreateVariableParams) (store.Variable, error) {
	v := store.Variable{
		ID:                  arg.ID,
		VariableKey:         arg.VariableKey,
		VariableValue:       arg.VariableValue,
		VariableDescription: arg.VariableDescription,
		VariableCategory:    arg.VariableCategory,
		Sensitive:           arg.Sensitive,
		Hcl:                 arg.Hcl,
		Incomplete:          arg.Incomplete,
		WorkspaceID:         arg.WorkspaceID,
	}
	s.variables = append(s.variables, v)
	return v, nil
}

func (s *nestedStub) UpdateVariable(_ context.Context, arg store.UpdateVariableParams) (store.Variable, error) {
	for i, v := range s.variables {
		if v.ID != arg.ID {
			continue
		}
		s.variables[i].VariableKey = arg.VariableKey
		s.variables[i].VariableValue = arg.VariableValue
		s.variables[i].VariableDescription = arg.VariableDescription
		s.variables[i].VariableCategory = arg.VariableCategory
		s.variables[i].Sensitive = arg.Sensitive
		s.variables[i].Hcl = arg.Hcl
		s.variables[i].Incomplete = arg.Incomplete
		return s.variables[i], nil
	}
	return store.Variable{}, pgx.ErrNoRows
}

func (s *nestedStub) DeleteVariable(_ context.Context, id string) error {
	for i, v := range s.variables {
		if v.ID == id {
			s.variables = append(s.variables[:i], s.variables[i+1:]...)
			return nil
		}
	}
	return pgx.ErrNoRows
}

const (
	orgID   = "11111111-1111-4111-8111-111111111111"
	wsID    = "22222222-2222-4222-8222-222222222222"
	otherWS = "33333333-3333-4333-8333-333333333333"
	varID   = "44444444-4444-4444-8444-444444444444"
)

func newNestedStub() *nestedStub {
	return &nestedStub{
		orgs: []store.Organization{{ID: orgID, Name: "demo", ExecutionMode: "remote"}},
		workspaces: []store.Workspace{
			{ID: wsID, OrganizationID: orgID, Name: "primary", ExecutionMode: strPtr("remote")},
			{ID: otherWS, OrganizationID: orgID, Name: "namesake", ExecutionMode: strPtr("remote")},
		},
		variables: []store.Variable{
			{ID: varID, VariableKey: "FOO", VariableValue: "bar", VariableCategory: "TERRAFORM", WorkspaceID: wsID},
			{ID: "55555555-5555-4555-8555-555555555555", VariableKey: "BAZ", VariableValue: "qux", VariableCategory: "ENV", WorkspaceID: otherWS},
		},
	}
}

type nestedDocument struct {
	Data struct {
		Type          string         `json:"type"`
		ID            string         `json:"id"`
		Attributes    map[string]any `json:"attributes"`
		Relationships map[string]struct {
			Data json.RawMessage `json:"data"`
		} `json:"relationships"`
	} `json:"data"`
	Included []struct {
		Type       string         `json:"type"`
		ID         string         `json:"id"`
		Attributes map[string]any `json:"attributes"`
	} `json:"included"`
	Errors []struct {
		Detail string `json:"detail"`
	} `json:"errors"`
}

func decodeJSON(rec *httptest.ResponseRecorder, out any) error {
	return json.Unmarshal(rec.Body.Bytes(), out)
}

func serveNested(t *testing.T, s *nestedStub, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	newTestAPI(t, s).ServeHTTP(rec, req)
	return rec
}

// The UI loads a workspace through the organization-scoped path with
// include=job,variable,... and reads the variable list out of `included`
// (Details.tsx), so listVariables alone cannot make the Variables tab work.
func TestNestedWorkspaceGetIncludesVariables(t *testing.T) {
	rec := serveNested(t, newNestedStub(), http.MethodGet,
		"/organization/"+orgID+"/workspace/"+wsID+"?include=variable,organization", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc.Data.ID != wsID {
		t.Errorf("data.id = %q, want %q", doc.Data.ID, wsID)
	}

	var ids []struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}
	if err := json.Unmarshal(doc.Data.Relationships["variable"].Data, &ids); err != nil {
		t.Fatalf("decode relationships.variable: %v", err)
	}
	if len(ids) != 1 || ids[0].ID != varID || ids[0].Type != "variable" {
		t.Errorf("relationships.variable.data = %v, want one variable %s", ids, varID)
	}

	var includedVariable, includedOrganization bool
	for _, inc := range doc.Included {
		switch inc.Type {
		case "variable":
			includedVariable = inc.ID == varID && inc.Attributes["key"] == "FOO"
		case "organization":
			includedOrganization = inc.ID == orgID && inc.Attributes["name"] == "demo"
		}
	}
	if !includedVariable {
		t.Errorf("included = %v, want the workspace's variable %s", doc.Included, varID)
	}
	if !includedOrganization {
		t.Errorf("included = %v, want the workspace's organization", doc.Included)
	}
}

func TestNestedWorkspaceInAnotherOrganizationIsNotFound(t *testing.T) {
	s := newNestedStub()
	s.workspaces[0].OrganizationID = "99999999-9999-4999-8999-999999999999"

	rec := serveNested(t, s, http.MethodGet, "/organization/"+orgID+"/workspace/"+wsID, "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestNestedVariablesListIsScopedToTheWorkspace(t *testing.T) {
	rec := serveNested(t, newNestedStub(), http.MethodGet, "/organization/"+orgID+"/workspace/"+wsID+"/variable", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	var doc struct {
		Data []struct {
			ID         string         `json:"id"`
			Attributes map[string]any `json:"attributes"`
		} `json:"data"`
	}
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data) != 1 || doc.Data[0].ID != varID {
		t.Fatalf("data = %v, want only the variable of this workspace", doc.Data)
	}
}

func TestNestedVariableCreateUsesTheWorkspaceInThePath(t *testing.T) {
	s := newNestedStub()
	rec := serveNested(t, s, http.MethodPost, "/organization/"+orgID+"/workspace/"+wsID+"/variable",
		`{"data":{"type":"variable","attributes":{"key":"NEW","value":"x","category":"ENV"}}}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", rec.Code, rec.Body.String())
	}
	var created store.Variable
	for _, v := range s.variables {
		if v.VariableKey == "NEW" {
			created = v
		}
	}
	if created.WorkspaceID != wsID {
		t.Errorf("created workspace = %q, want %q", created.WorkspaceID, wsID)
	}
}

// A variable id from a sibling workspace must not be reachable through this
// workspace's path: the handler verifies the workspace, the patch must verify
// the variable too.
func TestNestedVariableFromAnotherWorkspaceIsNotFound(t *testing.T) {
	rec := serveNested(t, newNestedStub(), http.MethodPatch,
		"/organization/"+orgID+"/workspace/"+wsID+"/variable/55555555-5555-4555-8555-555555555555",
		`{"data":{"type":"variable","id":"55555555-5555-4555-8555-555555555555","attributes":{"key":"HACKED"}}}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestNestedVariableDeleteRemovesOnlyThatVariable(t *testing.T) {
	s := newNestedStub()
	rec := serveNested(t, s, http.MethodDelete, "/organization/"+orgID+"/workspace/"+wsID+"/variable/"+varID, "")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body = %s", rec.Code, rec.Body.String())
	}
	if len(s.variables) != 1 || s.variables[0].ID == varID {
		t.Errorf("variables after delete = %v, want only the sibling workspace's", s.variables)
	}
}
