package jsonapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteResource(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteResource(rec, http.StatusCreated, Resource{
		Type:       "organization",
		ID:         "org-1",
		Attributes: map[string]any{"name": "simple"},
	})

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, MediaType) {
		t.Errorf("Content-Type = %q, want %q", ct, MediaType)
	}

	var doc map[string]map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if doc["data"]["type"] != "organization" {
		t.Errorf("type = %v", doc["data"]["type"])
	}
	if doc["data"]["id"] != "org-1" {
		t.Errorf("id = %v", doc["data"]["id"])
	}
}

func TestRelationshipMarshalling(t *testing.T) {
	resource := Resource{
		Type: "runs",
		ID:   "run-1",
		Relationships: map[string]Relationship{
			"workspace":  Link("workspaces", "ws-1"),
			"run-events": LinkList("run-events"),
			"tags":       LinkList("tags", "a", "b"),
			"templates":  NullLink(),
		},
	}
	raw, err := json.Marshal(resource)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	body := string(raw)
	for _, want := range []string{
		`"workspace":{"data":{"type":"workspaces","id":"ws-1"}}`,
		`"run-events":{"data":[]}`,
		`"tags":{"data":[{"type":"tags","id":"a"},{"type":"tags","id":"b"}]}`,
		`"templates":{"data":null}`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("resource = %s, missing %s", body, want)
		}
	}
}

// The three shapes below were captured from upstream 2.34.3; see
// docs/upstream/rest/organization-list.json.
func TestWriteCollectionPage(t *testing.T) {
	cases := []struct {
		name string
		page Page
		want map[string]any
	}{
		{"no page parameter", Page{}, nil},
		{"page number only", Page{Number: 2, Limit: 10000}, map[string]any{"number": 2, "limit": 10000}},
		{"page totals", Page{Number: 1, Limit: 2, Totals: true}, map[string]any{"number": 1, "limit": 2, "totalPages": 2, "totalRecords": 3}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			WriteCollection(rec, http.StatusOK, nil, 3, tc.page)

			var doc struct {
				Data []any          `json:"data"`
				Meta map[string]any `json:"meta"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if doc.Data == nil || len(doc.Data) != 0 {
				t.Errorf("data = %v, want empty array", doc.Data)
			}
			if tc.want == nil {
				if doc.Meta != nil {
					t.Fatalf("meta = %v, want none", doc.Meta)
				}
				return
			}
			got, _ := doc.Meta["page"].(map[string]any)
			if len(got) != len(tc.want) {
				t.Fatalf("meta.page = %v, want %v", got, tc.want)
			}
			for k, want := range tc.want {
				if got[k] != float64(want.(int)) {
					t.Errorf("meta.page.%s = %v, want %v", k, got[k], want)
				}
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusNotFound, "Not Found", "nope")

	var doc struct {
		Errors []Error `json:"errors"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Errors) != 1 {
		t.Fatalf("errors = %d, want 1", len(doc.Errors))
	}
	if doc.Errors[0].Status != "404" {
		t.Errorf("status = %q, want 404", doc.Errors[0].Status)
	}
}

func TestDecodeRequest(t *testing.T) {
	body := `{"data":{"type":"workspace","attributes":{"name":"demo"},"relationships":{"project":{"data":{"type":"project","id":"p1"}}}}}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	data, err := DecodeRequest(req)
	if err != nil {
		t.Fatalf("DecodeRequest() error = %v", err)
	}
	if data.Type != "workspace" {
		t.Errorf("type = %q", data.Type)
	}
	if !HasAttribute(data, "name") {
		t.Error("HasAttribute(name) = false")
	}
	if HasAttribute(data, "missing") {
		t.Error("HasAttribute(missing) = true")
	}

	id, present, err := RelationshipID(data, "project")
	if err != nil {
		t.Fatalf("RelationshipID() error = %v", err)
	}
	if !present || id != "p1" {
		t.Errorf("RelationshipID() = (%q, %v), want (p1, true)", id, present)
	}
}

func TestRelationshipIDNull(t *testing.T) {
	body := `{"data":{"type":"workspace","relationships":{"project":{"data":null}}}}`
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))

	data, err := DecodeRequest(req)
	if err != nil {
		t.Fatalf("DecodeRequest() error = %v", err)
	}
	id, present, err := RelationshipID(data, "project")
	if err != nil {
		t.Fatalf("RelationshipID() error = %v", err)
	}
	if !present || id != "" {
		t.Errorf("RelationshipID() = (%q, %v), want (\"\", true)", id, present)
	}
}

func TestRelationshipList(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"data":[{"type":"template","id":"t1"},{"type":"template","id":"t2"}]}`))
	ids, err := RelationshipList(req)
	if err != nil {
		t.Fatalf("RelationshipList() error = %v", err)
	}
	if len(ids) != 2 || ids[1].ID != "t2" {
		t.Errorf("RelationshipList() = %v", ids)
	}

	// An explicit null clears the relationship; a missing data member is an
	// error rather than a silent clear.
	req = httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"data":null}`))
	if ids, err := RelationshipList(req); err != nil || ids != nil {
		t.Errorf("null data = (%v, %v), want (nil, nil)", ids, err)
	}
	req = httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{}`))
	if _, err := RelationshipList(req); err == nil {
		t.Error("missing data should be an error")
	}
}

func TestRelationshipIdentifiers(t *testing.T) {
	body := `{"data":{"type":"workspaces","relationships":{"tags":{"data":[{"type":"tags","attributes":{"name":"env"}},{"type":"tags","attributes":{"name":"team"}}]},"project":{"data":{"type":"projects","id":"p1"}}}}}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	data, err := DecodeRequest(req)
	if err != nil {
		t.Fatalf("DecodeRequest() error = %v", err)
	}
	ids, present, err := RelationshipIdentifiers(data, "tags")
	if err != nil {
		t.Fatalf("RelationshipIdentifiers() error = %v", err)
	}
	if !present || len(ids) != 2 {
		t.Fatalf("RelationshipIdentifiers() = (%v, %v), want 2 present", ids, present)
	}
	if string(ids[0].Attributes["name"]) != `"env"` {
		t.Errorf("first tag name = %s", ids[0].Attributes["name"])
	}

	// The same document still decodes the to-one project relationship.
	id, present, err := RelationshipID(data, "project")
	if err != nil || !present || id != "p1" {
		t.Errorf("RelationshipID(project) = (%q, %v, %v)", id, present, err)
	}
}

func TestDecodeRequestRejectsBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{"))
	if _, err := DecodeRequest(req); err == nil {
		t.Fatal("DecodeRequest() error = nil, want error")
	}
}
