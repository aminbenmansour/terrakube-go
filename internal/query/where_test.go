package query

import (
	"reflect"
	"testing"
)

var testFields = map[string]Field{
	"name":               {Column: "name"},
	"lastJobStatus":      {Column: "last_job_status"},
	"lastJobDate":        {Column: "last_job_date"},
	"project.id":         {Column: "project_id"},
	"workspaceTag.tagId": {Column: "workspace_tag_id"},
	"id":                 {Column: "id"},
}

func mustParse(t *testing.T, expression string) Node {
	t.Helper()
	node, err := ParseRSQL(expression)
	if err != nil {
		t.Fatalf("ParseRSQL(%q): %v", expression, err)
	}
	return node
}

func TestWhereOperators(t *testing.T) {
	cases := []struct {
		name       string
		expression string
		wantSQL    string
		wantArgs   []any
	}{
		{
			name:       "equality",
			expression: `lastJobStatus=='running'`,
			wantSQL:    "last_job_status = $1",
			wantArgs:   []any{"running"},
		},
		{
			name:       "inequality",
			expression: `lastJobStatus!='cancelled'`,
			wantSQL:    "last_job_status <> $1",
			wantArgs:   []any{"cancelled"},
		},
		{
			name:       "greater than",
			expression: `id=gt=42`,
			wantSQL:    "id > $1",
			wantArgs:   []any{"42"},
		},
		{
			name:       "in list",
			expression: `id=in=('a','b')`,
			wantSQL:    "id IN ($1, $2)",
			wantArgs:   []any{"a", "b"},
		},
		{
			name:       "not in list",
			expression: `id=out=('a','b')`,
			wantSQL:    "id NOT IN ($1, $2)",
			wantArgs:   []any{"a", "b"},
		},
		{
			name:       "is null",
			expression: `project.id=isnull=true`,
			wantSQL:    "project_id IS NULL",
			wantArgs:   nil,
		},
		{
			name:       "is not null",
			expression: `project.id=isnull=false`,
			wantSQL:    "project_id IS NOT NULL",
			wantArgs:   nil,
		},
		{
			name:       "case insensitive wildcard",
			expression: `name=ini=*prod*`,
			wantSQL:    `lower(name) LIKE $1 ESCAPE '\'`,
			wantArgs:   []any{"%prod%"},
		},
		{
			name:       "case insensitive exact",
			expression: `name=ini=Acme`,
			wantSQL:    `lower(name) LIKE $1 ESCAPE '\'`,
			wantArgs:   []any{"acme"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sql, args, err := Where(mustParse(t, tc.expression), testFields, 0)
			if err != nil {
				t.Fatalf("Where: %v", err)
			}
			if sql != tc.wantSQL {
				t.Fatalf("sql = %q, want %q", sql, tc.wantSQL)
			}
			if !reflect.DeepEqual(args, tc.wantArgs) {
				t.Fatalf("args = %#v, want %#v", args, tc.wantArgs)
			}
		})
	}
}

func TestWhereEscapesLikeWildcards(t *testing.T) {
	sql, args, err := Where(mustParse(t, `name=ini=*100%_done*`), testFields, 0)
	if err != nil {
		t.Fatalf("Where: %v", err)
	}
	if want := `lower(name) LIKE $1 ESCAPE '\'`; sql != want {
		t.Fatalf("sql = %q, want %q", sql, want)
	}
	if got := args[0].(string); got != `%100\%\_done%` {
		t.Fatalf("pattern = %q", got)
	}
}

func TestWhereCombinesAndOr(t *testing.T) {
	sql, args, err := Where(mustParse(t, `(name=ini=*prod*,description=ini=*prod*);lastJobStatus=='running'`), map[string]Field{
		"name":          {Column: "name"},
		"description":   {Column: "description"},
		"lastJobStatus": {Column: "last_job_status"},
	}, 0)
	if err != nil {
		t.Fatalf("Where: %v", err)
	}
	want := `((lower(name) LIKE $1 ESCAPE '\' OR lower(description) LIKE $2 ESCAPE '\') AND last_job_status = $3)`
	if sql != want {
		t.Fatalf("sql = %q, want %q", sql, want)
	}
	wantArgs := []any{"%prod%", "%prod%", "running"}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("args = %#v, want %#v", args, wantArgs)
	}
}

func TestWhereRespectsArgumentOffset(t *testing.T) {
	sql, args, err := Where(mustParse(t, `id=in=('a','b')`), testFields, 3)
	if err != nil {
		t.Fatalf("Where: %v", err)
	}
	if want := "id IN ($4, $5)"; sql != want {
		t.Fatalf("sql = %q, want %q", sql, want)
	}
	if !reflect.DeepEqual(args, []any{"a", "b"}) {
		t.Fatalf("args = %#v", args)
	}
}

func TestWhereRejectsUnknownField(t *testing.T) {
	_, _, err := Where(mustParse(t, `secret=ini=*x*`), testFields, 0)
	if err == nil {
		t.Fatal("unknown field was accepted")
	}
}

func TestWhereNilNode(t *testing.T) {
	sql, args, err := Where(nil, testFields, 0)
	if err != nil || sql != "" || args != nil {
		t.Fatalf("Where(nil) = (%q, %#v, %v)", sql, args, err)
	}
}

func TestWhereRendersCorrelatedPredicates(t *testing.T) {
	fields := map[string]Field{
		"id":            {Column: "id"},
		"lastJobStatus": {Column: "last_job_status"},
		"job.status":    {Predicate: "EXISTS (SELECT 1 FROM job WHERE job.workspace_id = workspace.id AND job.status = ?)"},
	}
	node := mustParse(t, `(lastJobStatus=='waitingApproval',job.status=='waitingApproval')`)
	sql, args, err := Where(node, fields, 0)
	if err != nil {
		t.Fatalf("Where: %v", err)
	}
	want := `(last_job_status = $1 OR EXISTS (SELECT 1 FROM job WHERE job.workspace_id = workspace.id AND job.status = $2))`
	if sql != want {
		t.Fatalf("sql = %q, want %q", sql, want)
	}
	if !reflect.DeepEqual(args, []any{"waitingApproval", "waitingApproval"}) {
		t.Fatalf("args = %#v", args)
	}
}

func TestWhereRejectsPredicateWithUnsupportedOperator(t *testing.T) {
	fields := map[string]Field{
		"job.status": {Predicate: "EXISTS (SELECT 1 FROM job WHERE job.status = ?)"},
	}
	if _, _, err := Where(mustParse(t, `job.status=in=('a','b')`), fields, 0); err == nil {
		t.Fatal("in list on a predicate was accepted")
	}
}
