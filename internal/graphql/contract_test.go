package graphql

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

var fixtureValue = regexp.MustCompile(`\$\{[A-Za-z_][A-Za-z0-9_]*\}`)

// TestUpstreamUIDocumentsMatchTheSchema validates every document captured from
// the upstream UI against our schema. It is the contract test for the GraphQL
// surface: a field the UI reads and we drop fails here, without a browser.
func TestUpstreamUIDocumentsMatchTheSchema(t *testing.T) {
	source, err := os.ReadFile("schema.graphqls")
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	schema, err := gqlparser.LoadSchema(&ast.Source{Name: "schema.graphqls", Input: string(source)})
	if err != nil {
		t.Fatalf("load schema: %v", err)
	}

	fixtures, err := filepath.Glob("../../docs/upstream/ui-operations/*.graphql")
	if err != nil {
		t.Fatalf("glob fixtures: %v", err)
	}
	if len(fixtures) == 0 {
		t.Fatal("no UI fixtures found")
	}

	for _, fixture := range fixtures {
		t.Run(filepath.Base(fixture), func(t *testing.T) {
			document, err := os.ReadFile(fixture)
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}
			// The UI interpolates ids and page sizes into the document text;
			// any value works for a shape check.
			normalized := fixtureValue.ReplaceAllString(string(document), "00000000-0000-0000-0000-000000000000")
			if _, err := gqlparser.LoadQuery(schema, normalized); err != nil {
				t.Fatalf("the upstream UI document no longer matches the schema: %v", err)
			}
		})
	}
}
