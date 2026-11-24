package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aminbenmansour/terrakube-go/internal/store"
)

func grantsAPI(t *testing.T) (*stubQuerier, http.Handler) {
	t.Helper()
	q := &stubQuerier{
		workspaces: []store.Workspace{{ID: capturedRelWS, OrganizationID: capturedRelOrg}},
		projects:   []store.Project{{ID: capturedRelP1, OrganizationID: capturedRelOrg}},
		access:     []store.Access{captureAccessRow()},
		projectAccess: []store.ProjectAccess{
			captureProjectAccessRow(),
		},
	}
	return q, newTestAPI(t, q)
}

func doGrant(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, request)
	return rec
}

// A create that sends only a name gets Elide's primitive defaults: destroyJob
// true, everything else false, and a null role.
func TestWorkspaceAccessCreateDefaults(t *testing.T) {
	_, handler := grantsAPI(t)
	rec := doGrant(t, handler, http.MethodPost,
		"/organization/"+capturedRelOrg+"/workspace/"+capturedRelWS+"/access",
		`{"data":{"type":"access","attributes":{"name":"new-group"}}}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{
		`"name":"new-group"`, `"destroyJob":true`, `"manageState":false`,
		`"manageWorkspace":false`, `"planJob":false`, `"approveJob":false`, `"role":null`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("create body %s does not contain %s", body, want)
		}
	}
}

func TestWorkspaceAccessCreateWithoutNameConflicts(t *testing.T) {
	_, handler := grantsAPI(t)
	rec := doGrant(t, handler, http.MethodPost,
		"/organization/"+capturedRelOrg+"/workspace/"+capturedRelWS+"/access",
		`{"data":{"type":"access","attributes":{"role":"read"}}}`)

	if rec.Code != http.StatusConflict {
		t.Errorf("status = %d, want 409: %s", rec.Code, rec.Body.String())
	}
}

// A PATCH that names one attribute leaves the rest of the row alone, because
// the UI's role picker sends only the role.
func TestWorkspaceAccessUpdateKeepsUnspecifiedAttributes(t *testing.T) {
	q, handler := grantsAPI(t)
	rec := doGrant(t, handler, http.MethodPatch,
		"/organization/"+capturedRelOrg+"/workspace/"+capturedRelWS+"/access/"+capturedAccessID,
		`{"data":{"type":"access","id":"`+capturedAccessID+`","attributes":{"role":"read"}}}`)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204: %s", rec.Code, rec.Body.String())
	}
	if q.updatedAccess == nil {
		t.Fatal("no update recorded")
	}
	got := *q.updatedAccess
	if got.Role == nil || *got.Role != "read" {
		t.Errorf("role = %v, want read", got.Role)
	}
	if !got.ManageWorkspace || !got.ManageState || !got.PlanJob || got.ManageJob || got.ApproveJob {
		t.Errorf("flags = %+v, want the row's flags preserved", got)
	}
}

// A grant may only be reached through its own workspace; the id alone is not
// enough, and upstream answers the same 404.
func TestWorkspaceAccessRejectsAnotherWorkspace(t *testing.T) {
	q, handler := grantsAPI(t)
	other := "00000000-0000-4000-8000-0000000000ff"
	q.workspaces = append(q.workspaces, store.Workspace{ID: other, OrganizationID: capturedRelOrg})
	rec := doGrant(t, handler, http.MethodGet,
		"/organization/"+capturedRelOrg+"/workspace/"+other+"/access/"+capturedAccessID, "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Unknown identifier "+capturedAccessID+" for access") {
		t.Errorf("detail = %s, want upstream's identifier message", rec.Body.String())
	}
}

// project_access.name is nullable upstream, so an unnamed grant is valid.
func TestProjectAccessCreateAllowsMissingName(t *testing.T) {
	_, handler := grantsAPI(t)
	rec := doGrant(t, handler, http.MethodPost,
		"/organization/"+capturedRelOrg+"/project/"+capturedRelP1+"/projectAccess",
		`{"data":{"type":"project_access","attributes":{"planJob":true}}}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{`"name":null`, `"destroyJob":true`, `"planJob":true`, `"type":"project_access"`} {
		if !strings.Contains(body, want) {
			t.Errorf("create body %s does not contain %s", body, want)
		}
	}
}

func TestProjectAccessRejectsAnotherOrganization(t *testing.T) {
	q, handler := grantsAPI(t)
	q.projects = []store.Project{{ID: capturedRelP1, OrganizationID: "70000000-0000-0000-0000-000000009999"}}
	rec := doGrant(t, handler, http.MethodGet,
		"/organization/"+capturedRelOrg+"/project/"+capturedRelP1+"/projectAccess", "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Unknown identifier "+capturedRelP1+" for project") {
		t.Errorf("detail = %s, want upstream's identifier message", rec.Body.String())
	}
}

func TestProjectAccessDelete(t *testing.T) {
	_, handler := grantsAPI(t)
	path := "/organization/" + capturedRelOrg + "/project/" + capturedRelP1 + "/projectAccess/" + capturedProjectAccess
	if rec := doGrant(t, handler, http.MethodDelete, path, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204: %s", rec.Code, rec.Body.String())
	}
}
