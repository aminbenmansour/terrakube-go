package jsonapi

import (
	"net/http/httptest"
	"testing"
)

func TestParseFieldsetsReadsEveryType(t *testing.T) {
	r := httptest.NewRequest("GET", "/organization?fields[organization]=name,description&fields[workspace]=name&fields[job]=", nil)
	sets := ParseFieldsets(r)

	if len(sets) != 3 {
		t.Fatalf("sets = %v, want three types", sets)
	}
	if !sets["organization"]["name"] || !sets["organization"]["description"] || len(sets["organization"]) != 2 {
		t.Errorf("organization = %v, want name and description", sets["organization"])
	}
	if !sets["workspace"]["name"] || len(sets["workspace"]) != 1 {
		t.Errorf("workspace = %v, want name", sets["workspace"])
	}
	if names, present := sets["job"]; !present || len(names) != 0 {
		t.Errorf("job = %v (present %v), want a present, empty set", names, present)
	}
}

func TestRestrictKeepsOnlyTheRequestedAttributes(t *testing.T) {
	doc := Document{
		Data: Resource{
			Type: "organization",
			ID:   "1",
			Attributes: map[string]any{
				"name": "azure", "description": "sample", "icon": nil,
			},
			Relationships: map[string]Relationship{"project": LinkList("project")},
		},
	}
	sets := Fieldsets{"organization": {"name": true}}

	restricted := sets.Restrict(doc)
	data := restricted.Data.(Resource)

	if len(data.Attributes) != 1 || data.Attributes["name"] != "azure" {
		t.Errorf("attributes = %v, want only name", data.Attributes)
	}
	if data.ID != "1" || data.Type != "organization" || data.Relationships == nil {
		t.Errorf("id, type or relationships changed: %+v", data)
	}
}

// Upstream omits attributes entirely for a type that has no fields parameter
// once any fields parameter is present, and for an empty parameter list.
func TestRestrictDropsAttributesForUnrequestedAndEmptySets(t *testing.T) {
	doc := Document{
		Data: []Resource{
			{Type: "organization", ID: "1", Attributes: map[string]any{"name": "azure"}},
			{Type: "workspace", ID: "2", Attributes: map[string]any{"name": "ws"}},
		},
	}
	restricted := Fieldsets{"workspace": {}}.Restrict(doc)
	data := restricted.Data.([]Resource)

	if data[0].Attributes != nil {
		t.Errorf("organization attributes = %v, want omitted", data[0].Attributes)
	}
	if data[1].Attributes != nil {
		t.Errorf("workspace attributes = %v, want omitted", data[1].Attributes)
	}
}

func TestRestrictFiltersIncludedResources(t *testing.T) {
	doc := Document{
		Data:     Resource{Type: "workspace", ID: "2", Attributes: map[string]any{"name": "ws"}},
		Included: []Resource{{Type: "variable", ID: "3", Attributes: map[string]any{"key": "FOO", "value": "bar"}}},
	}
	restricted := Fieldsets{"workspace": {"name": true}, "variable": {"key": true}}.Restrict(doc)

	if len(restricted.Included) != 1 || len(restricted.Included[0].Attributes) != 1 ||
		restricted.Included[0].Attributes["key"] != "FOO" {
		t.Errorf("included = %v, want only variable.key", restricted.Included)
	}
}

func TestRestrictLeavesDocumentsWithoutFieldsetsAlone(t *testing.T) {
	doc := Document{Data: Resource{Type: "organization", ID: "1", Attributes: map[string]any{"name": "azure"}}}
	restricted := Fieldsets(nil).Restrict(doc)
	if data := restricted.Data.(Resource); len(data.Attributes) != 1 {
		t.Errorf("attributes = %v, want untouched", data.Attributes)
	}
}
