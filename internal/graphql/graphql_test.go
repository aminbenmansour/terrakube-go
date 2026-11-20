package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aminbenmansour/terrakube-go/internal/db"
	"github.com/aminbenmansour/terrakube-go/migrations"
)

func TestCollectionSpecBuild(t *testing.T) {
	filter := `name=ini=*prod*;project.id=isnull=false`
	sort := "-lastJobDate,name"
	first := StringOrInt("25")
	after := StringOrInt("50")
	args := listArgs{Filter: &filter, First: &first, After: &after, Sort: &sort}

	sqlText, countText, sqlArgs, offset, err := workspaceSpec.build("organization_id = $1", []any{"org-1"}, args)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if offset != 50 {
		t.Fatalf("offset = %d, want 50", offset)
	}
	if want := "SELECT count(*) FROM workspace WHERE organization_id = $1 AND (lower(name) LIKE $2 ESCAPE '\\' AND project_id IS NOT NULL)"; countText != want {
		t.Fatalf("count sql = %q, want %q", countText, want)
	}
	for _, fragment := range []string{
		"FROM workspace WHERE organization_id = $1",
		"ORDER BY last_job_date DESC, name ASC",
		"LIMIT $3 OFFSET $4",
	} {
		if !strings.Contains(sqlText, fragment) {
			t.Fatalf("sql %q is missing %q", sqlText, fragment)
		}
	}
	wantArgs := []any{"org-1", "%prod%", int64(25), int64(50)}
	if len(sqlArgs) != len(wantArgs) {
		t.Fatalf("args = %#v, want %#v", sqlArgs, wantArgs)
	}
	for i := range wantArgs {
		if sqlArgs[i] != wantArgs[i] {
			t.Fatalf("args[%d] = %#v, want %#v", i, sqlArgs[i], wantArgs[i])
		}
	}
}

func TestCollectionSpecRejectsInjection(t *testing.T) {
	cases := []listArgs{
		{Filter: ptr(`name==x;id) OR true--==1`)},
		{Filter: ptr(`secret=='x'`)},
		{Sort: ptr("-name;DROP TABLE workspace")},
		{Sort: ptr("-notAColumn")},
		{Filter: ptr(`name==`)},
	}
	for _, args := range cases {
		if _, _, _, _, err := workspaceSpec.build("", nil, args); err == nil {
			t.Fatalf("build(%#v) was accepted", args)
		}
	}
}

func TestCollectionSpecScopesByIDs(t *testing.T) {
	_, countText, sqlArgs, _, err := organizationSpec.build("", nil, listArgs{IDs: []string{"a", "b"}})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if want := "SELECT count(*) FROM organization WHERE id = ANY($1)"; countText != want {
		t.Fatalf("count sql = %q, want %q", countText, want)
	}
	if len(sqlArgs) != 3 {
		t.Fatalf("args = %#v", sqlArgs)
	}
}

func TestStringOrIntAcceptsBothForms(t *testing.T) {
	for value, want := range map[any]string{"10000": "10000", 20: "20", int64(20): "20", 20.0: "20", json.Number("7"): "7"} {
		var scalar StringOrInt
		if err := scalar.UnmarshalGQL(value); err != nil {
			t.Fatalf("UnmarshalGQL(%#v): %v", value, err)
		}
		if got := string(scalar); got != want {
			t.Fatalf("UnmarshalGQL(%#v) = %q, want %q", value, got, want)
		}
	}
}

func TestSchemaLoads(t *testing.T) {
	if NewExecutableSchema(Config{Resolvers: &Resolver{}}) == nil {
		t.Fatal("schema did not load")
	}
}

func ptr[T any](value T) *T { return &value }

type graphQLResponse struct {
	Data   map[string]any `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

const workspacePageQuery = `query WorkspacePage($organizationIds: [String], $first: StringOrInt, $after: StringOrInt, $filter: String, $sort: String, $statusAll: String) {
  organization(ids: $organizationIds) {
    edges {
      node {
        name
        workspace(first: $first, after: $after, filter: $filter, sort: $sort) {
          edges {
            node {
              id
              name
              lastJobStatus
              lastJobDate
              workspaceTag { edges { node { tagId value } } }
              project { edges { node { id name } } }
            }
          }
          pageInfo { startCursor endCursor hasNextPage hasPreviousPage totalRecords }
        }
        statusAll: workspace(first: "1", filter: $statusAll) { pageInfo { totalRecords } }
      }
    }
  }
}`

// TestWorkspacePageEndToEnd runs the shape of the UI's workspace page query
// against a real database: RSQL filtering, sorting, offset paging, the
// relationship connections and the aliased facet count.
func TestWorkspacePageEndToEnd(t *testing.T) {
	deps := testDeps(t)
	ctx := context.Background()

	orgID := uuid.NewString()
	projectID := uuid.NewString()
	if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx, `INSERT INTO organization (id, name) VALUES ($1, $2)`, orgID, "gql-test-"+orgID[:8]); err != nil {
		t.Fatalf("insert organization: %v", err)
	}
	t.Cleanup(func() {
		_, _ = deps.DB.(*pgxpool.Pool).Exec(context.Background(), `DELETE FROM organization WHERE id = $1`, orgID)
	})
	if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx, `INSERT INTO project (id, name, organization_id) VALUES ($1, 'core', $2)`, projectID, orgID); err != nil {
		t.Fatalf("insert project: %v", err)
	}

	names := []string{"prod-web", "prod-api", "dev-api"}
	for i, name := range names {
		workspaceID := uuid.NewString()
		if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx,
			`INSERT INTO workspace (id, name, organization_id, project_id, last_job_status) VALUES ($1, $2, $3, $4, $5)`,
			workspaceID, name, orgID, projectID, []string{"success", "failed", "NeverExecuted"}[i]); err != nil {
			t.Fatalf("insert workspace: %v", err)
		}
		if i == 0 {
			tagID := uuid.NewString()
			if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx, `INSERT INTO tag (id, name, organization_id) VALUES ($1, 'env', $2)`, tagID, orgID); err != nil {
				t.Fatalf("insert tag: %v", err)
			}
			if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx, `INSERT INTO workspacetag (id, tag_id, tag_value, workspace_id) VALUES ($1, $2, 'prod', $3)`, uuid.NewString(), tagID, workspaceID); err != nil {
				t.Fatalf("insert workspacetag: %v", err)
			}
		}
	}

	handler := New(deps)

	page := runQuery(t, handler, workspacePageQuery, map[string]any{
		"organizationIds": []string{orgID},
		"first":           "10",
		"filter":          `(name=ini='*prod*',description=ini='*prod*')`,
		"sort":            "name",
		"statusAll":       `name=ini='*prod*'`,
	})

	connection := page["organization"].(map[string]any)["edges"].([]any)[0].(map[string]any)["node"].(map[string]any)
	workspaces := connection["workspace"].(map[string]any)
	edges := workspaces["edges"].([]any)
	if len(edges) != 2 {
		t.Fatalf("got %d workspaces, want 2: %#v", len(edges), edges)
	}
	first := edges[0].(map[string]any)["node"].(map[string]any)
	second := edges[1].(map[string]any)["node"].(map[string]any)
	if first["name"] != "prod-api" || second["name"] != "prod-web" {
		t.Fatalf("sort order = %v, %v", first["name"], second["name"])
	}
	if first["lastJobStatus"] != "failed" {
		t.Fatalf("lastJobStatus = %v", first["lastJobStatus"])
	}
	project := first["project"].(map[string]any)["edges"].([]any)[0].(map[string]any)["node"].(map[string]any)
	if project["name"] != "core" {
		t.Fatalf("project = %#v", project)
	}

	info := workspaces["pageInfo"].(map[string]any)
	if info["totalRecords"].(float64) != 2 || info["hasNextPage"] != false || info["endCursor"] != "1" {
		t.Fatalf("pageInfo = %#v", info)
	}
	if got := connection["statusAll"].(map[string]any)["pageInfo"].(map[string]any)["totalRecords"].(float64); got != 2 {
		t.Fatalf("statusAll totalRecords = %v, want 2", got)
	}

	// Second page of one row, using the cursor the first page returned.
	next := runQuery(t, handler, workspacePageQuery, map[string]any{
		"organizationIds": []string{orgID},
		"first":           "1",
		"after":           "1",
		"filter":          `name=ini='*prod*'`,
		"sort":            "name",
	})
	nextConnection := next["organization"].(map[string]any)["edges"].([]any)[0].(map[string]any)["node"].(map[string]any)["workspace"].(map[string]any)
	nextEdges := nextConnection["edges"].([]any)
	if len(nextEdges) != 1 || nextEdges[0].(map[string]any)["node"].(map[string]any)["name"] != "prod-web" {
		t.Fatalf("second page = %#v", nextEdges)
	}
	nextInfo := nextConnection["pageInfo"].(map[string]any)
	if nextInfo["hasPreviousPage"] != true || nextInfo["hasNextPage"] != false {
		t.Fatalf("second pageInfo = %#v", nextInfo)
	}
}

func TestFilterByTagAndJobStatus(t *testing.T) {
	deps := testDeps(t)
	ctx := context.Background()

	orgID := uuid.NewString()
	if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx, `INSERT INTO organization (id, name) VALUES ($1, $2)`, orgID, "gql-test-"+orgID[:8]); err != nil {
		t.Fatalf("insert organization: %v", err)
	}
	t.Cleanup(func() {
		_, _ = deps.DB.(*pgxpool.Pool).Exec(context.Background(), `DELETE FROM organization WHERE id = $1`, orgID)
	})

	tagID := uuid.NewString()
	if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx, `INSERT INTO tag (id, name, organization_id) VALUES ($1, 'env', $2)`, tagID, orgID); err != nil {
		t.Fatalf("insert tag: %v", err)
	}

	taggedID, otherID := uuid.NewString(), uuid.NewString()
	for _, ws := range []struct{ id, name, status string }{
		{taggedID, "tagged", "running"},
		{otherID, "other", "failed"},
	} {
		if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx,
			`INSERT INTO workspace (id, name, organization_id, last_job_status) VALUES ($1, $2, $3, $4)`,
			ws.id, ws.name, orgID, ws.status); err != nil {
			t.Fatalf("insert workspace: %v", err)
		}
	}
	if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx,
		`INSERT INTO workspacetag (id, tag_id, tag_value, workspace_id) VALUES ($1, $2, 'prod', $3)`,
		uuid.NewString(), tagID, taggedID); err != nil {
		t.Fatalf("insert workspacetag: %v", err)
	}

	if _, err := deps.DB.(*pgxpool.Pool).Exec(ctx,
		`INSERT INTO job (organization_id, workspace_id, status) VALUES ($1, $2, 'failed')`, orgID, otherID); err != nil {
		t.Fatalf("insert job: %v", err)
	}

	handler := New(deps)
	for _, tc := range []struct {
		filter string
		want   string
	}{
		{`workspaceTag.tagId=='` + tagID + `'`, "tagged"},
		{`job.status=='failed'`, "other"},
		{`(lastJobStatus=='failed',job.status=='failed')`, "other"},
		{`lastJobStatus=='running'`, "tagged"},
	} {
		query := `query Scoped($filter: String) { organization(ids: ["` + orgID + `"]) { edges { node { workspace(filter: $filter) { edges { node { name } } pageInfo { totalRecords } } } } } }`
		page := runQuery(t, handler, query, map[string]any{"filter": tc.filter})
		node := page["organization"].(map[string]any)["edges"].([]any)[0].(map[string]any)["node"].(map[string]any)
		edges := node["workspace"].(map[string]any)["edges"].([]any)
		if len(edges) != 1 || edges[0].(map[string]any)["node"].(map[string]any)["name"] != tc.want {
			t.Fatalf("filter %q = %#v, want %s", tc.filter, edges, tc.want)
		}
	}
}

func TestResolverErrorsAreNotLeaked(t *testing.T) {
	handler := New(Deps{DB: failingDB{}, Logger: slog.New(slog.NewTextHandler(io.Discard, nil))})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/graphql/api/v1",
		strings.NewReader(`{"query":"{ organization { edges { node { id } } } }"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, request)

	var response graphQLResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(response.Errors) != 1 {
		t.Fatalf("errors = %#v", response.Errors)
	}
	if response.Errors[0].Message != "The request could not be completed." {
		t.Fatalf("message = %q", response.Errors[0].Message)
	}
}

func TestAuthenticationIsRequiredWhenConfigured(t *testing.T) {
	handler := New(Deps{DB: failingDB{}, Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), AuthToken: "secret"})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/graphql/api/v1", strings.NewReader(`{"query":"{ __typename }"}`))
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", recorder.Code)
	}
}

type failingDB struct{}

func (failingDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	panic("unreachable")
}

func (failingDB) QueryRow(context.Context, string, ...any) pgx.Row { return failingRow{} }

type failingRow struct{}

func (failingRow) Scan(...any) error {
	return errors.New(`pq: relation "organization" does not exist`)
}

func runQuery(t *testing.T, handler http.Handler, query string, variables map[string]any) map[string]any {
	t.Helper()
	body, err := json.Marshal(map[string]any{"query": query, "variables": variables})
	if err != nil {
		t.Fatalf("encode request: %v", err)
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/graphql/api/v1", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(recorder, request)

	var response graphQLResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if len(response.Errors) != 0 {
		t.Fatalf("query failed: %#v", response.Errors)
	}
	return response.Data
}

func testDeps(t *testing.T) Deps {
	t.Helper()
	url := os.Getenv("TERRAKUBE_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TERRAKUBE_TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := db.Migrate(ctx, pool, migrations.Files, "."); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return Deps{DB: pool, Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
}
