package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aminbenmansour/terrakube-go/internal/auth"
)

func TestOrganizationScope(t *testing.T) {
	tests := []struct {
		name  string
		ids   []string
		scope string
		args  int
	}{
		{"owner sees everything", nil, "", 0},
		{"no group sees nothing", []string{}, "FALSE", 0},
		{"member is limited to its organizations", []string{"a", "b"}, "organization_id = ANY($1)", 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope, args := organizationScope("organization_id", test.ids)
			if scope != test.scope || len(args) != test.args {
				t.Fatalf("scope = %q args = %#v, want %q with %d args", scope, args, test.scope, test.args)
			}
		})
	}
}

func TestOwnerRule(t *testing.T) {
	resolver := &Resolver{adminGroup: "TERRAKUBE_ADMIN"}
	tests := []struct {
		name     string
		identity auth.Identity
		want     bool
	}{
		{"the development token owns the stack", auth.Identity{Issuer: auth.DevIssuer}, true},
		{"the admin group owns the instance", auth.Identity{Issuer: "dex", Groups: []string{"TERRAKUBE_ADMIN"}}, true},
		{"a team member does not", auth.Identity{Issuer: "dex", Groups: []string{"readers"}}, false},
		{"a token without groups does not", auth.Identity{Issuer: "dex"}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := resolver.isOwner(&test.identity); got != test.want {
				t.Errorf("isOwner = %v, want %v", got, test.want)
			}
		})
	}
}

type groupVerifier struct{}

func (groupVerifier) Verify(_ context.Context, raw string) (*auth.Identity, error) {
	groups, ok := strings.CutPrefix(raw, "groups:")
	if !ok {
		return nil, errors.New("unknown token")
	}
	return &auth.Identity{Issuer: "dex", Email: "user@example.com", Groups: strings.Split(groups, ",")}, nil
}

func runAuthorizedQuery(t *testing.T, handler http.Handler, token, query string) map[string]any {
	t.Helper()
	body, err := json.Marshal(map[string]any{"query": query})
	if err != nil {
		t.Fatalf("encode request: %v", err)
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/graphql/api/v1", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
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

func TestOrganizationQueryIsScopedToTheCaller(t *testing.T) {
	deps := testDeps(t)
	deps.AuthToken = "dev"
	deps.Verifiers = []auth.Verifier{groupVerifier{}}
	handler := New(deps)

	ctx := context.Background()
	pool := deps.DB.(*pgxpool.Pool)

	mine, other := uuid.NewString(), uuid.NewString()
	for _, org := range []struct{ id, name string }{
		{mine, "gql-mine-" + mine[:8]},
		{other, "gql-other-" + other[:8]},
	} {
		if _, err := pool.Exec(ctx, `INSERT INTO organization (id, name) VALUES ($1, $2)`, org.id, org.name); err != nil {
			t.Fatalf("insert organization: %v", err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM organization WHERE id = ANY($1)`, []string{mine, other})
	})
	if _, err := pool.Exec(ctx,
		`INSERT INTO team (id, name, organization_id, role) VALUES ($1, 'readers', $2, 'read')`,
		uuid.NewString(), mine); err != nil {
		t.Fatalf("insert team: %v", err)
	}

	data := runAuthorizedQuery(t, handler, "groups:readers",
		`{ organization { edges { node { id } } } }`)
	edges := data["organization"].(map[string]any)["edges"].([]any)
	if len(edges) != 1 {
		t.Fatalf("edges = %#v, want only the caller's organization", edges)
	}
	if id := edges[0].(map[string]any)["node"].(map[string]any)["id"]; id != mine {
		t.Errorf("organization = %v, want %s", id, mine)
	}

	data = runAuthorizedQuery(t, handler, "groups:strangers",
		`{ organization { edges { node { id } } } }`)
	if edges := data["organization"].(map[string]any)["edges"].([]any); len(edges) != 0 {
		t.Errorf("edges = %#v, want none for a token without groups", edges)
	}
}
