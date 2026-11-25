package api

import (
	"slices"
	"testing"

	"github.com/google/uuid"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
	"github.com/aminbenmansour/terrakube-go/internal/store"
)

// seedWorkspaces creates four workspaces in a project, one more in the
// organization without a project, and one soft-deleted workspace that no list
// may return.
func seedWorkspaces(t *testing.T) (*store.Queries, store.Organization, store.Project, []store.Workspace) {
	t.Helper()
	q := store.New(dbtest.Pool(t))
	org := createListOrganization(t, q, "workspaces")
	project, err := q.CreateProject(t.Context(), store.CreateProjectParams{
		ID: uuid.NewString(), Name: "platform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	rows := make([]store.Workspace, 0, 6)
	for _, name := range []string{"alpha", "bravo", "charlie", "delta", "echo"} {
		arg := store.CreateWorkspaceParams{
			ID: uuid.NewString(), Name: name, IacType: "terraform", OrganizationID: org.ID,
		}
		if name != "echo" {
			arg.ProjectID = &project.ID
		}
		ws, err := q.CreateWorkspace(t.Context(), arg)
		if err != nil {
			t.Fatalf("create workspace %q: %v", name, err)
		}
		rows = append(rows, ws)
	}

	gone, err := q.CreateWorkspace(t.Context(), store.CreateWorkspaceParams{
		ID: uuid.NewString(), Name: "gone", IacType: "terraform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create workspace gone: %v", err)
	}
	if err := q.SoftDeleteWorkspace(t.Context(), gone.ID); err != nil {
		t.Fatalf("soft delete: %v", err)
	}
	return q, org, project, rows
}

func TestWorkspaceListMatchesInMemory(t *testing.T) {
	q, org, _, rows := seedWorkspaces(t)
	handler := newTestAPI(t, q)

	for _, query := range []string{
		"",
		"sort=-name",
		"sort=-name&page[size]=2&page[number]=2",
		"filter[workspace]=name==bravo",
		"filter[workspace]=name==ch*",
		"page[size]=3&page[number]=2&page[totals]",
	} {
		t.Run(query, func(t *testing.T) {
			got, meta := getNames(t, handler, "/organization/"+org.ID+"/workspace?"+query)

			q := parseList(newListRequest("/organization/"+org.ID+"/workspace?"+query), "workspace")
			want, total := filterSortPage(rows, q, func(ws store.Workspace) string { return ws.Name })
			wantNames := make([]string, 0, len(want))
			for _, ws := range want {
				wantNames = append(wantNames, ws.Name)
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

func TestProjectWorkspaceListMatchesInMemory(t *testing.T) {
	q, _, project, rows := seedWorkspaces(t)
	handler := newTestAPI(t, q)

	// The project list holds the four workspaces assigned to it; "echo" belongs
	// to the organization only and "gone" is soft-deleted.
	inProject := rows[:4]
	for _, query := range []string{
		"",
		"sort=-name",
		"page[size]=2&page[number]=2",
		"filter[workspace]=name==ch*",
	} {
		t.Run(query, func(t *testing.T) {
			target := "/project/" + project.ID + "/workspace?" + query
			got, _ := getNames(t, handler, target)

			q := parseList(newListRequest(target), "workspace")
			want, _ := filterSortPage(inProject, q, func(ws store.Workspace) string { return ws.Name })
			wantNames := make([]string, 0, len(want))
			for _, ws := range want {
				wantNames = append(wantNames, ws.Name)
			}
			if !slices.Equal(got, wantNames) {
				t.Fatalf("names = %v, want %v", got, wantNames)
			}
		})
	}
}
