package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// seedOrganizations inserts rows whose lowercase names sort the same way under
// every collation, so the SQL order and filterSortPage agree.
func seedOrganizations(t *testing.T) (*store.Queries, []store.Organization) {
	t.Helper()
	q := store.New(dbtest.Pool(t))
	names := []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}
	rows := make([]store.Organization, 0, len(names))
	for _, name := range names {
		org, err := q.CreateOrganization(t.Context(), store.CreateOrganizationParams{
			ID: uuid.NewString(), Name: name, ExecutionMode: "remote",
		})
		if err != nil {
			t.Fatalf("create organization %q: %v", name, err)
		}
		rows = append(rows, org)
	}
	return q, rows
}

// TestOrganizationListMatchesInMemory checks the SQL window against the helper
// the other list endpoints still use: the endpoint was cut over without a
// behaviour change.
func TestOrganizationListMatchesInMemory(t *testing.T) {
	q, rows := seedOrganizations(t)
	handler := newTestAPI(t, q)

	for _, query := range []string{
		"",
		"sort=-name",
		"sort=-name&page[size]=2&page[number]=2",
		"filter[organization]=name==bravo",
		"filter[organization]=name==ch*",
		"page[size]=3&page[number]=3&page[totals]",
	} {
		t.Run(query, func(t *testing.T) {
			got, meta := getOrganizationNames(t, handler, "/organization?"+query)

			q := parseList(httptest.NewRequest(http.MethodGet, "/organization?"+query, nil), "organization")
			want, total := filterSortPage(rows, q, func(o store.Organization) string { return o.Name })
			wantNames := make([]string, 0, len(want))
			for _, org := range want {
				wantNames = append(wantNames, org.Name)
			}
			if !slices.Equal(got, wantNames) {
				t.Fatalf("names = %v, want %v", got, wantNames)
			}
			if q.Totals && meta.TotalRecords != total {
				t.Errorf("totalRecords = %d, want %d", meta.TotalRecords, total)
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
		// Equality is case-insensitive, as the in-memory filter was.
		{"filter[organization]=name==BRAVO", []string{"bravo"}},
		{"filter[organization]=name==ch*", []string{"charlie"}},
		{"filter[organization]=name==nobody", []string{}},
		{"page[size]=2&page[number]=3", []string{"echo", "foxtrot"}},
		{"page[size]=2&page[number]=4", []string{}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			got, _ := getOrganizationNames(t, handler, "/organization?"+tc.query)
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
		})
	}

	t.Run("page totals count the filtered set", func(t *testing.T) {
		_, meta := getOrganizationNames(t, handler,
			"/organization?filter[organization]=name==a*&page[size]=1&page[totals]")
		if meta.TotalRecords != 1 || meta.TotalPages != 1 {
			t.Errorf("meta = %+v, want totalRecords 1, totalPages 1", meta)
		}
	})
}

// TestOrganizationListPrefiltersByAllowedIDs covers the membership restriction
// that used to run in Go: it belongs in the query so totals and pages count
// only visible rows.
func TestOrganizationListPrefiltersByAllowedIDs(t *testing.T) {
	q, rows := seedOrganizations(t)
	visible := []string{rows[1].ID, rows[3].ID}

	items, err := q.ListOrganizationsPage(t.Context(), store.ListOrganizationsPageParams{
		AllowedIds: visible,
	})
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

	none, err := q.ListOrganizationsPage(t.Context(), store.ListOrganizationsPageParams{
		AllowedIds: []string{},
	})
	if err != nil {
		t.Fatalf("list with an empty allow list: %v", err)
	}
	if len(none) != 0 {
		t.Errorf("names = %v, want none", none)
	}
}

// TestOrganizationListEscapesWildcards keeps a name containing % or _ a literal
// when it is used as a prefix.
func TestOrganizationListEscapesWildcards(t *testing.T) {
	q := store.New(dbtest.Pool(t))
	for _, name := range []string{"a%b", "a_b", "axb"} {
		if _, err := q.CreateOrganization(t.Context(), store.CreateOrganizationParams{
			ID: uuid.NewString(), Name: name, ExecutionMode: "remote",
		}); err != nil {
			t.Fatalf("create %q: %v", name, err)
		}
	}

	pattern := prefixPattern("a%")
	items, err := q.ListOrganizationsPage(t.Context(), store.ListOrganizationsPageParams{
		NamePrefix: &pattern,
	})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 || items[0].Name != "a%b" {
		t.Errorf("items = %+v, want only a%%b", items)
	}
}

type pageMeta struct {
	Number       int `json:"number"`
	Limit        int `json:"limit"`
	TotalPages   int `json:"totalPages"`
	TotalRecords int `json:"totalRecords"`
}

func getOrganizationNames(t *testing.T, handler http.Handler, target string) ([]string, pageMeta) {
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
