package api

import (
	"slices"
	"testing"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// The project and tag lists share the organization list's window parameters, so
// these cases mirror the ones in organizations_sql_test.go.

func TestProjectListSortsAndFilters(t *testing.T) {
	q := store.New(dbtest.Pool(t))
	org := createListOrganization(t, q, "projects")

	for _, name := range []string{"alpha", "bravo", "charlie", "delta"} {
		if _, err := q.CreateProject(t.Context(), store.CreateProjectParams{
			ID: uuid.NewString(), Name: name, OrganizationID: org.ID,
		}); err != nil {
			t.Fatalf("create project %q: %v", name, err)
		}
	}
	handler := newTestAPI(t, q)

	for _, tc := range []struct {
		query string
		want  []string
	}{
		{"", []string{"alpha", "bravo", "charlie", "delta"}},
		{"sort=-name", []string{"delta", "charlie", "bravo", "alpha"}},
		{"sort=-name&page[size]=2&page[number]=2", []string{"bravo", "alpha"}},
		{"filter[project]=name==bravo", []string{"bravo"}},
		{"filter[project]=name==ch*", []string{"charlie"}},
		{"page[size]=3&page[number]=2&page[totals]", []string{"delta"}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			got, meta := getNames(t, handler, "/organization/"+org.ID+"/project?"+tc.query)
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
			if tc.query == "page[size]=3&page[number]=2&page[totals]" && meta.TotalRecords != 4 {
				t.Errorf("totalRecords = %d, want 4", meta.TotalRecords)
			}
		})
	}
}

func TestTagListSortsAndFilters(t *testing.T) {
	q := store.New(dbtest.Pool(t))
	org := createListOrganization(t, q, "tags")

	// A tag may be nameless; Postgres sorts NULL after the names ascending and
	// before them descending.
	for _, name := range []*string{nil, strPtr("alpha"), strPtr("Bravo"), strPtr("charlie")} {
		if _, err := q.CreateTag(t.Context(), store.CreateTagParams{
			ID: uuid.NewString(), Name: name, OrganizationID: org.ID,
		}); err != nil {
			t.Fatalf("create tag %v: %v", name, err)
		}
	}
	handler := newTestAPI(t, q)

	for _, tc := range []struct {
		query string
		want  []string
	}{
		{"", []string{"alpha", "Bravo", "charlie", ""}},
		{"sort=-name", []string{"", "charlie", "Bravo", "alpha"}},
		{"sort=-name&page[size]=2&page[number]=2", []string{"Bravo", "alpha"}},
		{"filter[tag]=name==ch*", []string{"charlie"}},
		{"filter[tag]=name=ini=bravo", []string{"Bravo"}},
		{"page[size]=3&page[number]=2&page[totals]", []string{""}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			got, meta := getNames(t, handler, "/organization/"+org.ID+"/tag?"+tc.query)
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
			if tc.query == "page[size]=3&page[number]=2&page[totals]" && meta.TotalRecords != 4 {
				t.Errorf("totalRecords = %d, want 4", meta.TotalRecords)
			}
		})
	}
}

func TestTagListSQLFilters(t *testing.T) {
	q := store.New(dbtest.Pool(t))
	org := createListOrganization(t, q, "tag-filters")
	for _, name := range []string{"environment", "team", "tenant"} {
		if _, err := q.CreateTag(t.Context(), store.CreateTagParams{
			ID: uuid.NewString(), Name: strPtr(name), OrganizationID: org.ID,
		}); err != nil {
			t.Fatalf("create tag %q: %v", name, err)
		}
	}
	handler := newTestAPI(t, q)

	for _, tc := range []struct {
		query string
		want  []string
	}{
		{"filter[tag]=name==environment", []string{"environment"}},
		{"filter[tag]=name==ENVIRONMENT", []string{}},
		{"filter[tag]=name=ini='ENVIRONMENT'", []string{"environment"}},
		{"filter[tag]=name==te*", []string{"team", "tenant"}},
		{"filter[tag]=name==nobody", []string{}},
		{"page[size]=2&page[number]=2", []string{"tenant"}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			got, _ := getNames(t, handler, "/organization/"+org.ID+"/tag?"+tc.query)
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
		})
	}
}
