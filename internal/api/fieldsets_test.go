package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

// The UI builds breadcrumbs with fields[organization]=name; returning every
// attribute is a JSON:API violation and the replay of the capture pins it.
func TestSparseFieldsetsRestrictAttributes(t *testing.T) {
	rec := serveNested(t, newNestedStub(), http.MethodGet,
		"/organization/"+orgID+"?fields[organization]=name", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var doc nestedDocument
	if err := decodeJSON(rec, &doc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(doc.Data.Attributes) != 1 || doc.Data.Attributes["name"] != "demo" {
		t.Errorf("attributes = %v, want only name", doc.Data.Attributes)
	}
}

// A fields parameter for another type removes the attributes of this one, the
// way a live 2.34.3 does.
func TestSparseFieldsetsForAnotherTypeOmitAttributes(t *testing.T) {
	rec := serveNested(t, newNestedStub(), http.MethodGet,
		"/organization/"+orgID+"?fields[job]=status", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var raw struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, present := raw.Data["attributes"]; present {
		t.Errorf("data = %s, want no attributes member", rec.Body.String())
	}
}
