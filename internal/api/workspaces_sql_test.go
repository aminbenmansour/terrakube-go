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

	// A tag binding and jobs back the relationship filters.
	tag, err := q.CreateTag(t.Context(), store.CreateTagParams{
		ID: uuid.NewString(), Name: strPtr("env"), OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	for i, value := range []string{"prod", "dev"} {
		if _, err := q.CreateWorkspaceTag(t.Context(), store.CreateWorkspaceTagParams{
			ID: uuid.NewString(), TagID: tag.ID, TagValue: strPtr(value), WorkspaceID: rows[i].ID,
		}); err != nil {
			t.Fatalf("bind tag %q: %v", value, err)
		}
	}
	for i, status := range []string{"waitingApproval", "failed"} {
		if _, err := q.CreateJob(t.Context(), store.CreateJobParams{
			OrganizationID: org.ID, WorkspaceID: rows[i].ID, Status: status,
			Via: "UI", AutoApply: false, Refresh: true, PlanChanges: true,
		}); err != nil {
			t.Fatalf("create %s job: %v", status, err)
		}
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

func TestWorkspaceListSortsAndFilters(t *testing.T) {
	q, org, _, _ := seedWorkspaces(t)
	handler := newTestAPI(t, q)

	tags, err := q.ListTagsByOrganization(t.Context(), org.ID)
	if err != nil || len(tags) != 1 {
		t.Fatalf("list tags = %v, %v", tags, err)
	}
	tagID := tags[0].ID

	for _, tc := range []struct {
		query string
		want  []string
	}{
		{"", []string{"alpha", "bravo", "charlie", "delta", "echo"}},
		{"sort=-name", []string{"echo", "delta", "charlie", "bravo", "alpha"}},
		{"sort=-name&page[size]=2&page[number]=2", []string{"charlie", "bravo"}},
		{"filter[workspace]=name==bravo", []string{"bravo"}},
		{"filter[workspace]=name==ch*", []string{"charlie"}},
		{"filter[workspace]=(name=ini='*a*',name=ini='*e*')",
			[]string{"alpha", "bravo", "charlie", "delta", "echo"}},
		{"filter[workspace]=project.id=isnull=true", []string{"echo"}},
		{"filter[workspace]=project.id=isnull=false", []string{"alpha", "bravo", "charlie", "delta"}},
		{"filter[workspace]=workspaceTag.tagId=='" + tagID + "'", []string{"alpha", "bravo"}},
		{"filter[workspace]=workspaceTag.tagId!='" + tagID + "'", []string{"charlie", "delta", "echo"}},
		{"filter[workspace]=workspaceTag.value=ini='PROD'", []string{"alpha"}},
		{"filter[workspace]=job.status==waitingApproval", []string{"alpha"}},
		{"filter[workspace]=job.status!=failed", []string{"alpha", "charlie", "delta", "echo"}},
		{"filter[workspace]=job.status=isnull=true",
			[]string{"alpha", "bravo", "charlie", "delta", "echo"}},
		{"filter[workspace]=workspaceTag.tagId=hasmember='" + tagID + "'", []string{"alpha", "bravo"}},
		{"filter[workspace]=workspaceTag.tagId=hasnomember='" + tagID + "'",
			[]string{"charlie", "delta", "echo"}},
		{"filter[workspace]=workspaceTag=isempty=false", []string{"alpha", "bravo"}},
		{"filter[workspace]=job=isempty=true", []string{"charlie", "delta", "echo"}},
		{"page[size]=3&page[number]=2&page[totals]", []string{"delta", "echo"}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			got, meta := getNames(t, handler, "/organization/"+org.ID+"/workspace?"+tc.query)
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
			if tc.query == "page[size]=3&page[number]=2&page[totals]" && meta.TotalRecords != 5 {
				t.Errorf("totalRecords = %d, want 5", meta.TotalRecords)
			}
		})
	}
}

func TestProjectWorkspaceListSortsAndFilters(t *testing.T) {
	q, _, project, _ := seedWorkspaces(t)
	handler := newTestAPI(t, q)

	// "echo" belongs to the organization only and "gone" is soft-deleted.
	for _, tc := range []struct {
		query string
		want  []string
	}{
		{"", []string{"alpha", "bravo", "charlie", "delta"}},
		{"sort=-name", []string{"delta", "charlie", "bravo", "alpha"}},
		{"page[size]=2&page[number]=2", []string{"charlie", "delta"}},
		{"filter[workspace]=name==ch*", []string{"charlie"}},
		{"filter[workspace]=id=isnull=false", []string{"alpha", "bravo", "charlie", "delta"}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			target := "/project/" + project.ID + "/workspace?" + tc.query
			got, _ := getNames(t, handler, target)
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
		})
	}
}
