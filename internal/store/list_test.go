package store

import (
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/aminbenmansour/terrakube-go/internal/dbtest"
	"github.com/aminbenmansour/terrakube-go/internal/rsql"
)

// seedWorkspaces creates four live workspaces and one soft-deleted one. Two
// carry a job date and a status, so the comparison and null cases have rows to
// select.
func seedWorkspaces(t *testing.T) (*Queries, []Workspace) {
	t.Helper()
	pool := dbtest.Pool(t)
	q := New(pool)
	ctx := t.Context()

	org, err := q.CreateOrganization(ctx, CreateOrganizationParams{
		ID: uuid.NewString(), Name: "filters", ExecutionMode: "remote",
	})
	if err != nil {
		t.Fatalf("create organization: %v", err)
	}
	project, err := q.CreateProject(ctx, CreateProjectParams{
		ID: uuid.NewString(), Name: "platform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	rows := make([]Workspace, 0, 5)
	for _, seed := range []struct {
		name   string
		parent *string
		date   *time.Time
		status string
	}{
		{"prod-api", &project.ID, ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)), "success"},
		{"Prod-Web", nil, ptr(time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)), "failed"},
		{"staging", &project.ID, nil, "NeverExecuted"},
		{"other", nil, nil, "NeverExecuted"},
	} {
		ws, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
			ID: uuid.NewString(), Name: seed.name, IacType: "terraform",
			OrganizationID: org.ID, ProjectID: seed.parent,
		})
		if err != nil {
			t.Fatalf("create workspace %q: %v", seed.name, err)
		}
		if seed.date != nil || seed.status != "NeverExecuted" {
			if _, err := pool.Exec(ctx,
				`UPDATE workspace SET last_job_date = $2, last_job_status = $3 WHERE id = $1`,
				ws.ID, seed.date, seed.status); err != nil {
				t.Fatalf("set job columns for %q: %v", seed.name, err)
			}
			if ws, err = q.GetWorkspace(ctx, ws.ID); err != nil {
				t.Fatalf("reload workspace %q: %v", seed.name, err)
			}
		}
		rows = append(rows, ws)
	}

	// Tags and jobs back the relationship filters: prod-api carries the prod
	// tag and a job waiting for approval, staging carries the dev tag, and the
	// other workspace only has a discarded failed job.
	tag, err := q.CreateTag(ctx, CreateTagParams{
		ID: uuid.NewString(), Name: ptr("env"), OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	for i, value := range []string{"prod", "dev"} {
		if _, err := q.CreateWorkspaceTag(ctx, CreateWorkspaceTagParams{
			ID: uuid.NewString(), TagID: tag.ID, TagValue: ptr(value), WorkspaceID: rows[i*2].ID,
		}); err != nil {
			t.Fatalf("bind tag %q: %v", value, err)
		}
	}
	for i, status := range []string{"waitingApproval", "failed"} {
		if _, err := q.CreateJob(ctx, CreateJobParams{
			OrganizationID: org.ID, WorkspaceID: rows[i].ID, Status: status,
			Via: "UI", AutoApply: false, Refresh: true, PlanChanges: true,
		}); err != nil {
			t.Fatalf("create %s job: %v", status, err)
		}
	}
	discarded, err := q.CreateJob(ctx, CreateJobParams{
		OrganizationID: org.ID, WorkspaceID: rows[3].ID, Status: "failed",
		Via: "UI", AutoApply: false, Refresh: true, PlanChanges: true,
	})
	if err != nil {
		t.Fatalf("create discarded job: %v", err)
	}
	if _, err := q.db.Exec(ctx, `UPDATE job SET deleted = true WHERE id = $1`, discarded.ID); err != nil {
		t.Fatalf("discard job: %v", err)
	}

	gone, err := q.CreateWorkspace(ctx, CreateWorkspaceParams{
		ID: uuid.NewString(), Name: "gone", IacType: "terraform", OrganizationID: org.ID,
	})
	if err != nil {
		t.Fatalf("create soft-deleted workspace: %v", err)
	}
	if err := q.SoftDeleteWorkspace(ctx, gone.ID); err != nil {
		t.Fatalf("soft delete workspace: %v", err)
	}
	return q, rows
}

func ptr[T any](v T) *T { return &v }

func workspaceNames(rows []Workspace) []string {
	names := make([]string, 0, len(rows))
	for _, ws := range rows {
		names = append(names, ws.Name)
	}
	return names
}

func TestListPageFilters(t *testing.T) {
	q, rows := seedWorkspaces(t)
	scope := []any{rows[0].OrganizationID}
	ids := []string{rows[0].ID, rows[2].ID}

	for _, tc := range []struct {
		name   string
		filter string
		sorts  []SortKey
		want   []string
	}{
		{"no filter", "", nil, []string{"other", "prod-api", "Prod-Web", "staging"}},
		{"case sensitive equality", "name==prod-api", nil, []string{"prod-api"}},
		{"case sensitive miss", "name==PROD-API", nil, nil},
		{"case insensitive wildcard", "name=ini='*PROD*'", nil, []string{"prod-api", "Prod-Web"}},
		{"prefix wildcard", "name==prod*", nil, []string{"prod-api"}},
		{"negated wildcard", "name!=prod*", nil, []string{"other", "Prod-Web", "staging"}},
		{"set membership", "id=in=('" + ids[0] + "','" + ids[1] + "')", nil, []string{"prod-api", "staging"}},
		{"set exclusion", "id=out=('" + ids[0] + "')", nil, []string{"other", "Prod-Web", "staging"}},
		{"status set", "lastJobStatus=in=('success','failed')", nil, []string{"prod-api", "Prod-Web"}},
		{"null project", "project.id=isnull=true", nil, []string{"other", "Prod-Web"}},
		{"non-null project", "project.id=isnull=false", nil, []string{"prod-api", "staging"}},
		{"date comparison", "lastJobDate=gt=2026-01-15T00:00:00Z", nil, []string{"Prod-Web"}},
		{"and", "name=ini='*prod*';lastJobStatus==failed", nil, []string{"Prod-Web"}},
		{"or", "name==staging,name==other", nil, []string{"other", "staging"}},
		{"unsupported attribute widens", "workspaceTag.colour=='x'", nil,
			[]string{"other", "prod-api", "Prod-Web", "staging"}},
		{"descending", "", []SortKey{{Column: "name", Desc: true}},
			[]string{"staging", "Prod-Web", "prod-api", "other"}},
		{"unknown sort falls back to name", "", []SortKey{{Column: "bogus"}},
			[]string{"other", "prod-api", "Prod-Web", "staging"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var filter rsql.Node
			if tc.filter != "" {
				parsed, err := rsql.Parse(tc.filter)
				if err != nil {
					t.Fatalf("parse %q: %v", tc.filter, err)
				}
				filter = parsed
			}

			got, err := q.ListPage[Workspace](t.Context(), WorkspaceList, scope, filter, tc.sorts, nil, 0)
			if err != nil {
				t.Fatalf("ListPage: %v", err)
			}
			if names := workspaceNames(got); !slices.Equal(names, tc.want) {
				t.Errorf("names = %v, want %v", names, tc.want)
			}

			total, err := q.CountList(t.Context(), WorkspaceList, scope, filter)
			if err != nil {
				t.Fatalf("CountList: %v", err)
			}
			if want := int64(len(tc.want)); total != want {
				t.Errorf("count = %d, want %d", total, want)
			}
		})
	}
}

func TestListPageRelationshipFilters(t *testing.T) {
	q, rows := seedWorkspaces(t)
	scope := []any{rows[0].OrganizationID}
	ctx := t.Context()

	// The tag id comes from the binding seeded on prod-api.
	var tagID string
	if err := q.db.QueryRow(ctx, `SELECT tag_id FROM workspacetag WHERE workspace_id = $1`, rows[0].ID).Scan(&tagID); err != nil {
		t.Fatalf("read tag id: %v", err)
	}

	for _, tc := range []struct {
		name   string
		filter string
		want   []string
	}{
		{"tag binding", "workspaceTag.tagId=='" + tagID + "'", []string{"prod-api", "staging"}},
		{"tag exclusion includes untagged", "workspaceTag.tagId!='" + tagID + "'", []string{"other", "Prod-Web"}},
		{"tag value case insensitive", "workspaceTag.value=ini='PROD'", []string{"prod-api"}},
		{"tag value set", "workspaceTag.value=in=('dev')", []string{"staging"}},
		{"tag value exclusion", "workspaceTag.value=out=('dev')", []string{"other", "prod-api", "Prod-Web"}},
		{"job status", "job.status==waitingApproval", []string{"prod-api"}},
		{"discarded job is ignored", "job.status==failed", []string{"Prod-Web"}},
		{"job status exclusion", "job.status!=failed", []string{"other", "prod-api", "staging"}},
		{"unsupported operator widens", "job.status=isnull=true",
			[]string{"other", "prod-api", "Prod-Web", "staging"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			node, err := rsql.Parse(tc.filter)
			if err != nil {
				t.Fatalf("parse %q: %v", tc.filter, err)
			}
			got, err := q.ListPage[Workspace](ctx, WorkspaceList, scope, node, nil, nil, 0)
			if err != nil {
				t.Fatalf("ListPage: %v", err)
			}
			if names := workspaceNames(got); !slices.Equal(names, tc.want) {
				t.Errorf("names = %v, want %v", names, tc.want)
			}
		})
	}
}

func TestListPageSortsByJobDate(t *testing.T) {
	q, rows := seedWorkspaces(t)
	scope := []any{rows[0].OrganizationID}

	// Postgres puts NULLs last ascending and first descending; the NULL rows
	// are ordered by the id tiebreaker, so only the dated rows are asserted.
	got, err := q.ListPage[Workspace](t.Context(), WorkspaceList, scope, nil,
		[]SortKey{{Column: "lastJobDate"}}, nil, 0)
	if err != nil {
		t.Fatalf("ListPage: %v", err)
	}
	if names := workspaceNames(got); len(names) != 4 || names[0] != "prod-api" || names[1] != "Prod-Web" {
		t.Errorf("ascending names = %v, want the dated rows first", names)
	}
	if got[2].LastJobDate.Valid || got[3].LastJobDate.Valid {
		t.Errorf("ascending order should end with the undated rows: %v", workspaceNames(got))
	}

	got, err = q.ListPage[Workspace](t.Context(), WorkspaceList, scope, nil,
		[]SortKey{{Column: "lastJobDate", Desc: true}}, nil, 0)
	if err != nil {
		t.Fatalf("ListPage: %v", err)
	}
	if names := workspaceNames(got); len(names) != 4 || names[2] != "Prod-Web" || names[3] != "prod-api" {
		t.Errorf("descending names = %v, want newest first", names)
	}
}

func TestListPagePages(t *testing.T) {
	q, rows := seedWorkspaces(t)
	scope := []any{rows[0].OrganizationID}

	limit, offset := int64(2), int64(1)
	got, err := q.ListPage[Workspace](t.Context(), WorkspaceList, scope, nil, nil, &limit, offset)
	if err != nil {
		t.Fatalf("ListPage: %v", err)
	}
	if names := workspaceNames(got); !slices.Equal(names, []string{"prod-api", "Prod-Web"}) {
		t.Errorf("names = %v", names)
	}
}

func TestListPageRejectsBadValues(t *testing.T) {
	q, rows := seedWorkspaces(t)
	scope := []any{rows[0].OrganizationID}

	for _, filter := range []string{
		"lastJobDate=gt=yesterday",
		"lastJobDate=in=(2026-01-01T00:00:00Z,not-a-date)",
	} {
		node, err := rsql.Parse(filter)
		if err != nil {
			t.Fatalf("parse %q: %v", filter, err)
		}
		if _, err := q.ListPage[Workspace](t.Context(), WorkspaceList, scope, node, nil, nil, 0); !errors.Is(err, rsql.ErrSyntax) {
			t.Errorf("ListPage(%q) error = %v, want ErrSyntax", filter, err)
		}
	}

	node, err := rsql.Parse("lastJobStatus=isnull=maybe")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, err := q.ListPage[Workspace](t.Context(), WorkspaceList, scope, node, nil, nil, 0); !errors.Is(err, rsql.ErrSyntax) {
		t.Errorf("ListPage(=isnull=maybe) error = %v, want ErrSyntax", err)
	}
}

// TestModelColumnOrder pins the assumption ListPage scans on: sqlc emits model
// fields in table-column order, so SELECT * lines up with RowToStructByPos.
func TestModelColumnOrder(t *testing.T) {
	pool := dbtest.Pool(t)

	for _, tc := range []struct {
		table string
		model any
	}{
		{"organization", Organization{}},
		{"project", Project{}},
		{"tag", Tag{}},
		{"workspace", Workspace{}},
	} {
		t.Run(tc.table, func(t *testing.T) {
			rows, err := pool.Query(t.Context(),
				`SELECT column_name FROM information_schema.columns
				 WHERE table_name = $1 ORDER BY ordinal_position`, tc.table)
			if err != nil {
				t.Fatalf("read columns: %v", err)
			}
			columns, err := collectStrings(rows)
			if err != nil {
				t.Fatalf("collect columns: %v", err)
			}

			typ := reflect.TypeOf(tc.model)
			fields := make([]string, 0, typ.NumField())
			for i := 0; i < typ.NumField(); i++ {
				fields = append(fields, snakeCase(typ.Field(i).Name))
			}
			if !slices.Equal(columns, fields) {
				t.Errorf("table columns = %v, model fields = %v", columns, fields)
			}
		})
	}
}

func collectStrings(rows pgx.Rows) ([]string, error) {
	defer rows.Close()
	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

// snakeCase turns a Go field name into its column name: OrganizationID becomes
// organization_id and ModuleSshKey becomes module_ssh_key.
func snakeCase(s string) string {
	runes := []rune(s)
	out := make([]rune, 0, len(runes)+4)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && nextLower) {
				out = append(out, '_')
			}
		}
		out = append(out, unicode.ToLower(r))
	}
	return string(out)
}
