package query

import (
	"reflect"
	"testing"
)

func TestParseRSQLComparisons(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  Node
	}{
		{
			name:  "equality with quotes",
			input: `lastJobStatus=='running'`,
			want:  &Compare{Field: "lastJobStatus", Op: OpEq, Values: []string{"running"}},
		},
		{
			name:  "equality without quotes",
			input: `lastJobStatus==running`,
			want:  &Compare{Field: "lastJobStatus", Op: OpEq, Values: []string{"running"}},
		},
		{
			name:  "double quotes are accepted",
			input: `name=="Acme"`,
			want:  &Compare{Field: "name", Op: OpEq, Values: []string{"Acme"}},
		},
		{
			name:  "not equal",
			input: `status!='cancelled'`,
			want:  &Compare{Field: "status", Op: OpNotEq, Values: []string{"cancelled"}},
		},
		{
			name:  "relationship path",
			input: `project.id=='d9b58bd3'`,
			want:  &Compare{Field: "project.id", Op: OpEq, Values: []string{"d9b58bd3"}},
		},
		{
			name:  "in list",
			input: `workspaceTag.tagId=in=('a','b')`,
			want:  &Compare{Field: "workspaceTag.tagId", Op: OpIn, Values: []string{"a", "b"}},
		},
		{
			name:  "out list",
			input: `id=out=('a','b')`,
			want:  &Compare{Field: "id", Op: OpOut, Values: []string{"a", "b"}},
		},
		{
			name:  "isnull",
			input: `project=isnull=true`,
			want:  &Compare{Field: "project", Op: OpIsNull, Values: []string{"true"}},
		},
		{
			name:  "numeric comparison",
			input: `id=gt=42`,
			want:  &Compare{Field: "id", Op: OpGt, Values: []string{"42"}},
		},
		{
			name:  "case insensitive wildcard",
			input: `name=ini=*prod*`,
			want:  &Compare{Field: "name", Op: OpIni, Values: []string{"*prod*"}},
		},
		{
			name:  "quoted value keeps delimiters",
			input: `name=='Acme, Inc. (EU)'`,
			want:  &Compare{Field: "name", Op: OpEq, Values: []string{"Acme, Inc. (EU)"}},
		},
		{
			name:  "quoted value escapes a quote",
			input: `name=='O''Brien'`,
			want:  &Compare{Field: "name", Op: OpEq, Values: []string{"O'Brien"}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseRSQL(tc.input)
			if err != nil {
				t.Fatalf("ParseRSQL(%q): %v", tc.input, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("ParseRSQL(%q) = %#v, want %#v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseRSQLBooleanStructure(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  Node
	}{
		{
			name:  "comma is or",
			input: `(name=ini=*prod*,description=ini=*prod*)`,
			want: &Or{Nodes: []Node{
				&Compare{Field: "name", Op: OpIni, Values: []string{"*prod*"}},
				&Compare{Field: "description", Op: OpIni, Values: []string{"*prod*"}},
			}},
		},
		{
			name:  "semicolon is and",
			input: `lastJobStatus=='running';project.id=='x'`,
			want: &And{Nodes: []Node{
				&Compare{Field: "lastJobStatus", Op: OpEq, Values: []string{"running"}},
				&Compare{Field: "project.id", Op: OpEq, Values: []string{"x"}},
			}},
		},
		{
			name:  "and binds tighter than or",
			input: `a==1;b==2,c==3`,
			want: &Or{Nodes: []Node{
				&And{Nodes: []Node{
					&Compare{Field: "a", Op: OpEq, Values: []string{"1"}},
					&Compare{Field: "b", Op: OpEq, Values: []string{"2"}},
				}},
				&Compare{Field: "c", Op: OpEq, Values: []string{"3"}},
			}},
		},
		{
			name:  "parentheses override precedence",
			input: `(a==1,b==2);c==3`,
			want: &And{Nodes: []Node{
				&Or{Nodes: []Node{
					&Compare{Field: "a", Op: OpEq, Values: []string{"1"}},
					&Compare{Field: "b", Op: OpEq, Values: []string{"2"}},
				}},
				&Compare{Field: "c", Op: OpEq, Values: []string{"3"}},
			}},
		},
		{
			name:  "three way and",
			input: `a==1;b==2;c==3`,
			want: &And{Nodes: []Node{
				&Compare{Field: "a", Op: OpEq, Values: []string{"1"}},
				&Compare{Field: "b", Op: OpEq, Values: []string{"2"}},
				&Compare{Field: "c", Op: OpEq, Values: []string{"3"}},
			}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseRSQL(tc.input)
			if err != nil {
				t.Fatalf("ParseRSQL(%q): %v", tc.input, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("ParseRSQL(%q) = %#v, want %#v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseRSQLEmpty(t *testing.T) {
	for _, input := range []string{"", "   "} {
		got, err := ParseRSQL(input)
		if err != nil {
			t.Fatalf("ParseRSQL(%q): %v", input, err)
		}
		if got != nil {
			t.Fatalf("ParseRSQL(%q) = %#v, want nil", input, got)
		}
	}
}

func TestParseRSQLErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"unbalanced open paren", `(a==1`},
		{"unbalanced close paren", `a==1)`},
		{"unknown operator", `a=foo=1`},
		{"missing value", `a==`},
		{"missing operator", `a 1`},
		{"empty field", `==1`},
		{"empty in list", `a=in=()`},
		{"trailing junk", `a==1 b==2`},
		{"dangling conjunction", `a==1;`},
		{"unterminated quote", `a=='x`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := ParseRSQL(tc.input); err == nil {
				t.Fatalf("ParseRSQL(%q) = %#v, want an error", tc.input, got)
			}
		})
	}
}

// The expressions the upstream UI actually builds (see
// docs/upstream/ui-operations.md) must all parse.
func TestParseRSQLAcceptsUIExpressions(t *testing.T) {
	expressions := []string{
		`(name=ini='*prod*',description=ini='*prod*')`,
		`project.id=isnull=true`,
		`project.id=='<uuid>'`,
		`(lastJobStatus=isnull=true,lastJobStatus=='NeverExecuted')`,
		`(lastJobStatus=='waitingApproval',job.status=='waitingApproval')`,
		`lastJobStatus=='running'`,
		`(policyComplianceStatus=isnull=true,policyComplianceStatus=='UNKNOWN')`,
		`workspaceTag.tagId=='<uuid>'`,
		`workspaceTag.tagId=='<uuid>';workspaceTag.value=='prod'`,
		`id=in=('<uuid1>','<uuid2>')`,
		`id=isnull=true`,
		`name=="<org>"`,
	}
	for _, expression := range expressions {
		if _, err := ParseRSQL(expression); err != nil {
			t.Errorf("ParseRSQL(%q): %v", expression, err)
		}
	}
}
