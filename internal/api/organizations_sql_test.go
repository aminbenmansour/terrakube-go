package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
	"github.com/aminbenmansour/terrakube-go/internal/rsql"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// seedOrganizations inserts rows whose lowercase names sort the same way under
// every collation.
func seedOrganizations(t *testing.T) (*store.Queries, []store.Organization) {
	t.Helper()
	q := store.New(dbtest.Pool(t))
	names := []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}
	rows := make([]store.Organization, 0, len(names))
	for _, name := range names {
		rows = append(rows, createListOrganization(t, q, name))
	}
	return q, rows
}

func TestOrganizationListSortsAndPages(t *testing.T) {
	q, _ := seedOrganizations(t)
	handler := newTestAPI(t, q)

	for _, tc := range []struct {
		query string
		want  []string
		total int
	}{
		{"", []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}, 0},
		{"sort=-name", []string{"foxtrot", "echo", "delta", "charlie", "bravo", "alpha"}, 0},
		{"sort=-name&page[size]=2&page[number]=2", []string{"delta", "charlie"}, 0},
		{"page[size]=3&page[number]=3&page[totals]", []string{}, 6},
	} {
		t.Run(tc.query, func(t *testing.T) {
			got, meta := getNames(t, handler, "/organization?"+tc.query)
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
			if tc.total != 0 && meta.TotalRecords != tc.total {
				t.Errorf("totalRecords = %d, want %d", meta.TotalRecords, tc.total)
			}
		})
	}
}

func TestOrganizationListSQLFilters(t *testing.T) {
	q, _ := seedOrganizations(t)
	handler := newTestAPI(t, q)

	for _, tc := range []struct {
		query string
		want  []string
	}{
		// Equality is case-sensitive, as Elide's == is.
		{"filter[organization]=name==bravo", []string{"bravo"}},
		{"filter[organization]=name==BRAVO", []string{}},
		{"filter[organization]=name=ini='*BRAVO*'", []string{"bravo"}},
		{"filter[organization]=name==ch*", []string{"charlie"}},
		{"filter[organization]=name!=alpha%3Bname!=bravo", []string{"charlie", "delta", "echo", "foxtrot"}},
		{"filter[organization]=name==nobody", []string{}},
		{"page[size]=2&page[number]=3", []string{"echo", "foxtrot"}},
		{"page[size]=2&page[number]=4", []string{}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			got, _ := getNames(t, handler, "/organization?"+tc.query)
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
		})
	}

	t.Run("page totals count the filtered set", func(t *testing.T) {
		_, meta := getNames(t, handler,
			"/organization?filter[organization]=name==a*&page[size]=1&page[totals]")
		if meta.TotalRecords != 1 || meta.TotalPages != 1 {
			t.Errorf("meta = %+v, want totalRecords 1, totalPages 1", meta)
		}
	})
}

func TestOrganizationListRejectsMalformedFilter(t *testing.T) {
	handler := newTestAPI(t, &stubQuerier{})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/organization?filter[organization]=name==", nil))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	if !json.Valid(rec.Body.Bytes()) {
		t.Errorf("body is not JSON: %s", rec.Body.Bytes())
	}
}

// TestOrganizationListPrefiltersByAllowedIDs covers the membership restriction:
// it belongs in the query so totals and pages count only visible rows.
func TestOrganizationListPrefiltersByAllowedIDs(t *testing.T) {
	q, rows := seedOrganizations(t)
	visible := []string{rows[1].ID, rows[3].ID}

	items, err := q.ListOrganizations(t.Context(), []any{visible}, nil, nil, nil, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	got := make([]string, 0, len(items))
	for _, org := range items {
		got = append(got, org.Name)
	}
	if want := []string{"bravo", "delta"}; !slices.Equal(got, want) {
		t.Errorf("names = %v, want %v", got, want)
	}

	none, err := q.ListOrganizations(t.Context(), []any{[]string{}}, nil, nil, nil, 0)
	if err != nil {
		t.Fatalf("list with an empty allow list: %v", err)
	}
	if len(none) != 0 {
		t.Errorf("names = %v, want none", none)
	}

	all, err := q.ListOrganizations(t.Context(), []any{nil}, nil, nil, nil, 0)
	if err != nil {
		t.Fatalf("list as owner: %v", err)
	}
	if len(all) != len(rows) {
		t.Errorf("owner sees %d organizations, want %d", len(all), len(rows))
	}
}

// TestOrganizationListEscapesWildcards keeps a name containing % or _ a literal.
func TestOrganizationListEscapesWildcards(t *testing.T) {
	q := store.New(dbtest.Pool(t))
	for _, name := range []string{"a%b", "a_b", "axb"} {
		createListOrganization(t, q, name)
	}

	for filter, want := range map[string]string{
		"name==a%*": "a%b",
		"name==a_*": "a_b",
	} {
		node, err := rsql.Parse(filter)
		if err != nil {
			t.Fatalf("parse %q: %v", filter, err)
		}
		items, err := q.ListOrganizations(t.Context(), []any{nil}, node, nil, nil, 0)
		if err != nil {
			t.Fatalf("list %q: %v", filter, err)
		}
		if len(items) != 1 || items[0].Name != want {
			t.Errorf("%s = %+v, want only %s", filter, items, want)
		}
	}
}

type pageMeta struct {
	Number       int `json:"number"`
	Limit        int `json:"limit"`
	TotalPages   int `json:"totalPages"`
	TotalRecords int `json:"totalRecords"`
}

func getNames(t *testing.T, handler http.Handler, target string) ([]string, pageMeta) {
	t.Helper()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("%s: status = %d; body = %s", target, rec.Code, rec.Body.String())
	}
	var doc struct {
		Data []struct {
			Attributes struct {
				Name string `json:"name"`
			} `json:"attributes"`
		} `json:"data"`
		Meta struct {
			Page pageMeta `json:"page"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("%s: decode: %v", target, err)
	}
	names := make([]string, 0, len(doc.Data))
	for _, resource := range doc.Data {
		names = append(names, resource.Attributes.Name)
	}
	return names, doc.Meta.Page
}

func createListOrganization(t *testing.T, q *store.Queries, name string) store.Organization {
	t.Helper()
	org, err := q.CreateOrganization(t.Context(), store.CreateOrganizationParams{
		ID: uuid.NewString(), Name: name, ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization %q: %v", name, err)
	}
	return org
}
