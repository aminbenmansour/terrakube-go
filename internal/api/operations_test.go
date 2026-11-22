package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const atomicMedia = `application/vnd.api+json;ext="https://jsonapi.org/ext/atomic"`

func serveAtomic(t *testing.T, s *workspaceStub, body string, atomic bool) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(body))
	if atomic {
		req.Header.Set("Content-Type", atomicMedia)
		req.Header.Set("Accept", atomicMedia)
	} else {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}
	rec := httptest.NewRecorder()
	newTestAPI(t, s).ServeHTTP(rec, req)
	return rec
}

func decodeAtomicResults(t *testing.T, rec *httptest.ResponseRecorder) []struct {
	Data map[string]any `json:"data"`
} {
	t.Helper()
	var doc struct {
		Results []struct {
			Data map[string]any `json:"data"`
		} `json:"atomic:results"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v; body = %s", err, rec.Body.String())
	}
	return doc.Results
}

// The create screen sends op:add with a client-side lid and the same
// attributes and relationships POST /workspace accepts.
func TestAtomicOperationCreatesAWorkspace(t *testing.T) {
	s := newWorkspaceStub()
	rec := serveAtomic(t, s, `{"atomic:operations":[{"op":"add","href":"/organization/`+relOrgID+`/workspace",
		"data":{"type":"workspace","lid":"11111111-1111-4111-8111-111111111111",
		"attributes":{"name":"from-ui","source":"https://github.com/x/y.git","branch":"main","terraformVersion":"1.9.0"},
		"relationships":{"project":{"data":{"type":"project","id":"`+relProjectA+`"}}}}}]}`, true)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	results := decodeAtomicResults(t, rec)
	if len(results) != 1 || results[0].Data == nil {
		t.Fatalf("results = %v, want one created resource", results)
	}
	if id, _ := results[0].Data["id"].(string); id == "" {
		t.Errorf("data.id = %v, want the stored id", results[0].Data["id"])
	}
	if s.created == nil || s.created.Name != "from-ui" || s.created.ProjectID == nil || *s.created.ProjectID != relProjectA {
		t.Errorf("created = %+v, want the workspace from the operation", s.created)
	}
}

func TestAtomicOperationUpdatesAWorkspace(t *testing.T) {
	s := newWorkspaceStub()
	rec := serveAtomic(t, s, `{"atomic:operations":[{"op":"update","href":"/organization/`+relOrgID+`/workspace/`+relWorkspace+`",
		"data":{"type":"workspace","id":"`+relWorkspace+`","attributes":{"description":"after","sharedIds":"`+relWorkspace+`"}}}]}`, true)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	results := decodeAtomicResults(t, rec)
	if len(results) != 1 || results[0].Data != nil {
		t.Fatalf("results = %v, want one null result", results)
	}
	if s.updated == nil || s.updated.Description == nil || *s.updated.Description != "after" {
		t.Fatalf("updated = %+v, want the new description", s.updated)
	}
	if s.updated.SharedIds == nil || *s.updated.SharedIds != relWorkspace {
		t.Errorf("sharedIds = %v, want %s", s.updated.SharedIds, relWorkspace)
	}
}

// Upstream asks for the atomic media type on both headers.
func TestAtomicOperationRequiresTheAtomicMediaType(t *testing.T) {
	rec := serveAtomic(t, newWorkspaceStub(), `{"atomic:operations":[]}`, false)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want 415; body = %s", rec.Code, rec.Body.String())
	}
}

func TestAtomicOperationRejectsAnUnsupportedOp(t *testing.T) {
	rec := serveAtomic(t, newWorkspaceStub(),
		`{"atomic:operations":[{"op":"remove","href":"/organization/`+relOrgID+`/workspace/`+relWorkspace+`"}]}`, true)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
}

func TestAtomicOperationOnAnotherOrganizationsWorkspaceIsNotFound(t *testing.T) {
	rec := serveAtomic(t, newWorkspaceStub(),
		`{"atomic:operations":[{"op":"update","href":"/organization/`+foreignOrg+`/workspace/`+relWorkspace+`",
		"data":{"type":"workspace","id":"`+relWorkspace+`","attributes":{"description":"x"}}}]}`, true)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body = %s", rec.Code, rec.Body.String())
	}
}

func TestAtomicOperationWithNoOperationsReturnsNoResults(t *testing.T) {
	rec := serveAtomic(t, newWorkspaceStub(), `{"atomic:operations":[]}`, true)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"atomic:results":[]`) {
		t.Errorf("body = %s, want an empty atomic:results array", rec.Body.String())
	}
}
