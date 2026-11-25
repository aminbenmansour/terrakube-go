package rsql

import (
	"errors"
	"strings"
	"testing"
)

// dump renders an AST as an s-expression so the cases stay readable.
func dump(n Node) string {
	switch v := n.(type) {
	case And:
		return "(and " + dump(v.Left) + " " + dump(v.Right) + ")"
	case Or:
		return "(or " + dump(v.Left) + " " + dump(v.Right) + ")"
	case Compare:
		return "(" + strings.Join(append([]string{string(v.Op), v.Attr}, v.Values...), " ") + ")"
	}
	return "?"
}

func TestParse(t *testing.T) {
	// The expressions are the ones the upstream UI sends, verbatim from
	// docs/upstream/ui-operations.md.
	for _, tc := range []struct {
		raw  string
		want string
	}{
		{"name==prod", "(== name prod)"},
		{"name==prod*", "(== name prod*)"},
		{"name=ini='*prod*'", "(=ini= name *prod*)"},
		{"(name=ini='*prod*',description=ini='*prod*')", "(or (=ini= name *prod*) (=ini= description *prod*))"},
		{"project.id=isnull=true", "(=isnull= project.id true)"},
		{"project.id=='<uuid>'", "(== project.id <uuid>)"},
		{"(lastJobStatus=isnull=true,lastJobStatus=='NeverExecuted')",
			"(or (=isnull= lastJobStatus true) (== lastJobStatus NeverExecuted))"},
		{"(lastJobStatus=='waitingApproval',job.status=='waitingApproval')",
			"(or (== lastJobStatus waitingApproval) (== job.status waitingApproval))"},
		{"lastJobStatus=='running'", "(== lastJobStatus running)"},
		{"(policyComplianceStatus=isnull=true,policyComplianceStatus=='UNKNOWN')",
			"(or (=isnull= policyComplianceStatus true) (== policyComplianceStatus UNKNOWN))"},
		{"workspaceTag.tagId=='<uuid>'", "(== workspaceTag.tagId <uuid>)"},
		{"workspaceTag.tagId=='<uuid>';workspaceTag.value=='prod'",
			"(and (== workspaceTag.tagId <uuid>) (== workspaceTag.value prod))"},
		{"id=in=('<uuid1>','<uuid2>')", "(=in= id <uuid1> <uuid2>)"},
		{"id=isnull=true", "(=isnull= id true)"},
		{`name=="<org>"`, "(== name <org>)"},
		// Precedence: AND binds tighter than OR.
		{"a==1;b==2,c==3", "(or (and (== a 1) (== b 2)) (== c 3))"},
		{"a==1,b==2;c==3", "(or (== a 1) (and (== b 2) (== c 3)))"},
		// Quoted values keep delimiters and wildcards literal.
		{"name=='a;b,c'", "(== name a;b,c)"},
		{"name=='it\\'s'", "(== name it's)"},
		{"name=='*prod*'", "(== name *prod*)"},
		// Whitespace around tokens is allowed.
		{" name == prod ; status =in= ( a , b ) ", "(and (== name prod) (=in= status a b))"},
	} {
		t.Run(tc.raw, func(t *testing.T) {
			node, err := Parse(tc.raw)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.raw, err)
			}
			if got := dump(node); got != tc.want {
				t.Errorf("Parse(%q) = %s, want %s", tc.raw, got, tc.want)
			}
		})
	}
}

func TestParseRejects(t *testing.T) {
	for _, raw := range []string{
		"",
		"   ",
		"name",
		"name==",
		"name=in=()",
		"name=in=(a",
		"(name==x",
		"name==x;",
		"name==x)",
		"==x",
		"name=='unterminated",
		"name==a b",
	} {
		t.Run(raw, func(t *testing.T) {
			if _, err := Parse(raw); !errors.Is(err, ErrSyntax) {
				t.Fatalf("Parse(%q) error = %v, want ErrSyntax", raw, err)
			}
		})
	}
}
