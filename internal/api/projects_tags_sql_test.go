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

func TestProjectListMatchesInMemory(t *testing.T) {
	q := store.New(dbtest.Pool(t))
	org := createListOrganization(t, q, "projects")

	rows := make([]store.Project, 0, 4)
	for _, name := range []string{"alpha", "bravo", "charlie", "delta"} {
		project, err := q.CreateProject(t.Context(), store.CreateProjectParams{
			ID: uuid.NewString(), Name: name, OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("create project %q: %v", name, err)
		}
		rows = append(rows, project)
	}
	handler := newTestAPI(t, q)

	for _, query := range []string{
		"",
		"sort=-name",
		"sort=-name&page[size]=2&page[number]=2",
		"filter[project]=name==bravo",
		"filter[project]=name==ch*",
		"page[size]=3&page[number]=2&page[totals]",
	} {
		t.Run(query, func(t *testing.T) {
			got, meta := getNames(t, handler, "/organization/"+org.ID+"/project?"+query)

			q := parseList(newListRequest("/organization/"+org.ID+"/project?"+query), "project")
			want, total := filterSortPage(rows, q, func(p store.Project) string { return p.Name })
			wantNames := make([]string, 0, len(want))
			for _, project := range want {
				wantNames = append(wantNames, project.Name)
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

func TestTagListMatchesInMemory(t *testing.T) {
	q := store.New(dbtest.Pool(t))
	org := createListOrganization(t, q, "tags")

	// A tag may be nameless; the in-memory sort treated that as the empty
	// string, and the SQL order must do the same.
	rows := make([]store.Tag, 0, 4)
	for _, name := range []*string{nil, strPtr("alpha"), strPtr("Bravo"), strPtr("charlie")} {
		tag, err := q.CreateTag(t.Context(), store.CreateTagParams{
			ID: uuid.NewString(), Name: name, OrganizationID: org.ID,
		})
		if err != nil {
			t.Fatalf("create tag %v: %v", name, err)
		}
		rows = append(rows, tag)
	}
	handler := newTestAPI(t, q)

	for _, query := range []string{
		"",
		"sort=-name",
		"sort=-name&page[size]=2&page[number]=2",
		"filter[tag]=name==bravo",
		"filter[tag]=name==ch*",
		"page[size]=3&page[number]=2&page[totals]",
	} {
		t.Run(query, func(t *testing.T) {
			got, meta := getNames(t, handler, "/organization/"+org.ID+"/tag?"+query)

			q := parseList(newListRequest("/organization/"+org.ID+"/tag?"+query), "tag")
			want, total := filterSortPage(rows, q, func(tag store.Tag) string { return strValue(tag.Name) })
			wantNames := make([]string, 0, len(want))
			for _, tag := range want {
				wantNames = append(wantNames, strValue(tag.Name))
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
		{"filter[tag]=name==ENVIRONMENT", []string{"environment"}},
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
